package service

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"xiangshoufu/internal/async"
)

// MetricsService 监控指标服务
type MetricsService struct {
	channelStats   map[string]*ChannelMetrics
	queueStats     *QueueMetrics
	systemStats    *SystemMetrics
	alertService   *MessageService
	queue          *async.MemoryQueue
	alertThreshold *AlertThreshold
	mu             sync.RWMutex
}

// ChannelMetrics 通道指标
type ChannelMetrics struct {
	TotalCount     int64   `json:"total_count"`
	SuccessCount   int64   `json:"success_count"`
	FailCount      int64   `json:"fail_count"`
	DuplicateCount int64   `json:"duplicate_count"`
	SignFailCount  int64   `json:"sign_fail_count"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	latencySum     int64
	LastMinuteRT   []int64 `json:"-"` // 最近一分钟的响应时间
}

// QueueMetrics 队列指标
type QueueMetrics struct {
	RawCallbackLength  int `json:"raw_callback_length"`
	ProfitCalcLength   int `json:"profit_calc_length"`
	NotificationLength int `json:"notification_length"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	StartTime      time.Time `json:"start_time"`
	Uptime         string    `json:"uptime"`
	GoroutineCount int       `json:"goroutine_count"`
	MemoryUsageMB  float64   `json:"memory_usage_mb"`
}

// AlertThreshold 告警阈值
type AlertThreshold struct {
	SuccessRateMin   float64 `json:"success_rate_min"`    // 最低成功率，默认95%
	QueueLengthMax   int     `json:"queue_length_max"`    // 队列最大长度，默认500
	LatencyP99Max    int64   `json:"latency_p99_max"`     // P99延迟最大值(ms)，默认500
	SignFailCountMax int64   `json:"sign_fail_count_max"` // 签名失败最大次数/分钟，默认10
}

// DefaultAlertThreshold 默认告警阈值
func DefaultAlertThreshold() *AlertThreshold {
	return &AlertThreshold{
		SuccessRateMin:   0.95,
		QueueLengthMax:   500,
		LatencyP99Max:    500,
		SignFailCountMax: 10,
	}
}

// NewMetricsService 创建监控服务
func NewMetricsService(alertService *MessageService, queue *async.MemoryQueue) *MetricsService {
	return &MetricsService{
		channelStats:   make(map[string]*ChannelMetrics),
		queueStats:     &QueueMetrics{},
		systemStats:    &SystemMetrics{StartTime: time.Now()},
		alertService:   alertService,
		queue:          queue,
		alertThreshold: DefaultAlertThreshold(),
	}
}

// RecordSuccess 记录成功请求
func (m *MetricsService) RecordSuccess(channelCode string, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := m.getOrCreateStats(channelCode)
	atomic.AddInt64(&stats.TotalCount, 1)
	atomic.AddInt64(&stats.SuccessCount, 1)
	stats.latencySum += latency.Milliseconds()
	if stats.SuccessCount > 0 {
		stats.AvgLatencyMs = float64(stats.latencySum) / float64(stats.SuccessCount)
	}
}

// RecordFailure 记录失败请求
func (m *MetricsService) RecordFailure(channelCode, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := m.getOrCreateStats(channelCode)
	atomic.AddInt64(&stats.TotalCount, 1)
	atomic.AddInt64(&stats.FailCount, 1)
}

// RecordDuplicate 记录重复请求
func (m *MetricsService) RecordDuplicate(channelCode string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := m.getOrCreateStats(channelCode)
	atomic.AddInt64(&stats.TotalCount, 1)
	atomic.AddInt64(&stats.DuplicateCount, 1)
}

// RecordSignFailure 记录签名验证失败
func (m *MetricsService) RecordSignFailure(channelCode string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := m.getOrCreateStats(channelCode)
	atomic.AddInt64(&stats.TotalCount, 1)
	atomic.AddInt64(&stats.SignFailCount, 1)
}

// getOrCreateStats 获取或创建通道统计
func (m *MetricsService) getOrCreateStats(channelCode string) *ChannelMetrics {
	if _, exists := m.channelStats[channelCode]; !exists {
		m.channelStats[channelCode] = &ChannelMetrics{}
	}
	return m.channelStats[channelCode]
}

// GetAllMetrics 获取所有指标
func (m *MetricsService) GetAllMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 更新队列指标
	if m.queue != nil {
		lengths := m.queue.GetAllQueueLengths()
		m.queueStats.RawCallbackLength = lengths[async.TopicRawCallback]
		m.queueStats.ProfitCalcLength = lengths[async.TopicProfitCalc]
		m.queueStats.NotificationLength = lengths[async.TopicNotification]
	}

	// 更新系统指标
	m.systemStats.Uptime = time.Since(m.systemStats.StartTime).String()

	return map[string]interface{}{
		"channels": m.channelStats,
		"queue":    m.queueStats,
		"system":   m.systemStats,
	}
}

// CheckAndAlert 检查指标并发送告警
func (m *MetricsService) CheckAndAlert() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for channelCode, stats := range m.channelStats {
		// 检查成功率
		if stats.TotalCount > 100 { // 至少100个请求才告警
			successRate := float64(stats.SuccessCount) / float64(stats.TotalCount)
			if successRate < m.alertThreshold.SuccessRateMin {
				m.sendAlert(fmt.Sprintf("通道 %s 成功率过低: %.2f%% (阈值: %.2f%%)",
					channelCode, successRate*100, m.alertThreshold.SuccessRateMin*100))
			}
		}

		// 检查签名失败次数
		if stats.SignFailCount > m.alertThreshold.SignFailCountMax {
			m.sendAlert(fmt.Sprintf("通道 %s 签名验证失败次数过多: %d (阈值: %d)",
				channelCode, stats.SignFailCount, m.alertThreshold.SignFailCountMax))
		}

		// 检查延迟
		if stats.AvgLatencyMs > float64(m.alertThreshold.LatencyP99Max) {
			m.sendAlert(fmt.Sprintf("通道 %s 平均延迟过高: %.2fms (阈值: %dms)",
				channelCode, stats.AvgLatencyMs, m.alertThreshold.LatencyP99Max))
		}
	}

	// 检查队列长度
	if m.queue != nil {
		lengths := m.queue.GetAllQueueLengths()
		for topic, length := range lengths {
			if length > m.alertThreshold.QueueLengthMax {
				m.sendAlert(fmt.Sprintf("队列 %s 积压过多: %d (阈值: %d)",
					topic, length, m.alertThreshold.QueueLengthMax))
			}
		}
	}
}

// sendAlert 发送告警
func (m *MetricsService) sendAlert(msg string) {
	log.Printf("[ALERT] %s", msg)
	if m.alertService != nil {
		m.alertService.SendAlert(msg)
	}
}

// ResetStats 重置统计（通常每天重置一次）
func (m *MetricsService) ResetStats() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for channelCode := range m.channelStats {
		m.channelStats[channelCode] = &ChannelMetrics{}
	}

	log.Printf("[MetricsService] Stats reset")
}

// AlertCheckerJob 告警检查定时任务
type AlertCheckerJob struct {
	metricsService *MetricsService
}

// NewAlertCheckerJob 创建告警检查任务
func NewAlertCheckerJob(metricsService *MetricsService) *AlertCheckerJob {
	return &AlertCheckerJob{
		metricsService: metricsService,
	}
}

// Run 执行任务（每分钟执行一次）
func (j *AlertCheckerJob) Run() {
	j.metricsService.CheckAndAlert()
}
