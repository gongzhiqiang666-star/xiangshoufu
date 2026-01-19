package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"xiangshoufu/internal/async"
	"xiangshoufu/internal/cache"
	"xiangshoufu/internal/channel"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// CallbackHandler 统一回调处理器
type CallbackHandler struct {
	factory          *channel.AdapterFactory
	cache            cache.Cache
	queue            async.MessageQueue
	callbackRepo     repository.RawCallbackRepository
	metricsCollector *MetricsCollector
}

// NewCallbackHandler 创建回调处理器
func NewCallbackHandler(
	factory *channel.AdapterFactory,
	cache cache.Cache,
	queue async.MessageQueue,
	callbackRepo repository.RawCallbackRepository,
) *CallbackHandler {
	return &CallbackHandler{
		factory:          factory,
		cache:            cache,
		queue:            queue,
		callbackRepo:     callbackRepo,
		metricsCollector: NewMetricsCollector(),
	}
}

// CallbackResponse 回调响应
type CallbackResponse struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message,omitempty" example:"success"`
}

// HandleCallback 统一回调入口
// @Summary 支付通道回调
// @Description 接收各支付通道的回调通知，支持8个通道：HENGXINTONG(恒信通)、LAKALA(拉卡拉)、YEAHKA(乐刷)、SUIXINGFU(随行付)、LIANLIAN(连连)、SANDPAY(杉德)、FUIOU(富友)、HEEPAY(汇付天下)
// @Tags 回调
// @Accept json
// @Produce json
// @Param channel_code path string true "通道编码" Enums(HENGXINTONG, LAKALA, YEAHKA, SUIXINGFU, LIANLIAN, SANDPAY, FUIOU, HEEPAY)
// @Param body body object true "回调数据（各通道格式不同）"
// @Success 200 {object} CallbackResponse "成功响应"
// @Failure 500 {object} CallbackResponse "服务器错误"
// @Router /callback/{channel_code} [post]
func (h *CallbackHandler) HandleCallback(c *gin.Context) {
	startTime := time.Now()
	channelCode := c.Param("channel_code")

	// 1. 获取适配器
	adapter, err := h.factory.GetAdapter(channelCode)
	if err != nil {
		log.Printf("[Callback] Unknown channel: %s", channelCode)
		h.metricsCollector.RecordFailure(channelCode, "unknown_channel")
		c.JSON(http.StatusOK, CallbackResponse{Code: 0}) // 返回成功，避免通道重试
		return
	}

	// 2. 读取请求体
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[Callback] Read body failed: %v", err)
		h.metricsCollector.RecordFailure(channelCode, "read_body_error")
		c.JSON(http.StatusOK, CallbackResponse{Code: 0})
		return
	}

	// 3. 验证签名
	signVerified, err := adapter.VerifySign(rawBody)
	if err != nil {
		log.Printf("[Callback] Verify sign error: %v", err)
	}
	if !signVerified {
		log.Printf("[Callback] Sign verification failed for channel: %s", channelCode)
		h.metricsCollector.RecordSignFailure(channelCode)
		// 签名失败仍返回成功，防止通道重试，但记录告警
		c.JSON(http.StatusOK, CallbackResponse{Code: 0})
		return
	}

	// 4. 解析回调类型
	actionType, err := adapter.ParseActionType(rawBody)
	if err != nil {
		log.Printf("[Callback] Parse action type failed: %v", err)
		h.metricsCollector.RecordFailure(channelCode, "parse_action_error")
		c.JSON(http.StatusOK, CallbackResponse{Code: 0})
		return
	}

	// 5. 生成幂等键
	idempotentKey, err := adapter.ParseIdempotentKey(rawBody)
	if err != nil {
		log.Printf("[Callback] Parse idempotent key failed: %v", err)
		h.metricsCollector.RecordFailure(channelCode, "parse_key_error")
		c.JSON(http.StatusOK, CallbackResponse{Code: 0})
		return
	}

	// 6. 幂等检查（使用本地缓存）
	if h.cache.Exists(idempotentKey) {
		log.Printf("[Callback] Duplicate callback ignored: %s", idempotentKey)
		h.metricsCollector.RecordDuplicate(channelCode)
		c.JSON(http.StatusOK, CallbackResponse{Code: 0})
		return
	}

	// 7. 保存原始数据到数据库
	callbackLog := &models.RawCallbackLog{
		ChannelCode:   channelCode,
		ActionType:    string(actionType),
		RawRequest:    string(rawBody),
		SignVerified:  signVerified,
		ProcessStatus: models.ProcessStatusPending,
		IdempotentKey: idempotentKey,
		ClientIP:      c.ClientIP(),
		ReceivedAt:    time.Now(),
		CreatedDate:   time.Now(),
	}

	if err := h.callbackRepo.Create(callbackLog); err != nil {
		log.Printf("[Callback] Save callback log failed: %v", err)
		h.metricsCollector.RecordFailure(channelCode, "db_error")
		// 数据库写入失败，返回错误让通道重试
		c.JSON(http.StatusInternalServerError, CallbackResponse{Code: -1, Message: "internal error"})
		return
	}

	// 8. 设置幂等缓存（24小时过期）
	h.cache.Set(idempotentKey, callbackLog.ID, 24*time.Hour)

	// 9. 发送到异步队列处理
	queueMsg := &QueueMessage{
		CallbackLogID: callbackLog.ID,
		ChannelCode:   channelCode,
		ActionType:    string(actionType),
		RawBody:       rawBody,
	}
	msgBytes, _ := json.Marshal(queueMsg)
	if err := h.queue.Publish(async.TopicRawCallback, msgBytes); err != nil {
		log.Printf("[Callback] Publish to queue failed: %v", err)
		// 队列发送失败不影响响应，后续定时任务会兜底处理
	}

	// 10. 记录成功指标
	h.metricsCollector.RecordSuccess(channelCode, time.Since(startTime))

	// 11. 返回成功
	c.JSON(http.StatusOK, CallbackResponse{Code: 0})
}

// QueueMessage 队列消息
type QueueMessage struct {
	CallbackLogID int64  `json:"callback_log_id"`
	ChannelCode   string `json:"channel_code"`
	ActionType    string `json:"action_type"`
	RawBody       []byte `json:"raw_body"`
}

// GetMetrics 获取监控指标
func (h *CallbackHandler) GetMetrics() map[string]*ChannelStats {
	return h.metricsCollector.GetAllStats()
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
	stats map[string]*ChannelStats
}

// ChannelStats 通道统计
type ChannelStats struct {
	TotalCount     int64   `json:"total_count"`
	SuccessCount   int64   `json:"success_count"`
	FailCount      int64   `json:"fail_count"`
	DuplicateCount int64   `json:"duplicate_count"`
	SignFailCount  int64   `json:"sign_fail_count"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	latencySum     int64
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		stats: make(map[string]*ChannelStats),
	}
}

func (m *MetricsCollector) getOrCreateStats(channelCode string) *ChannelStats {
	if _, exists := m.stats[channelCode]; !exists {
		m.stats[channelCode] = &ChannelStats{}
	}
	return m.stats[channelCode]
}

func (m *MetricsCollector) RecordSuccess(channelCode string, latency time.Duration) {
	stats := m.getOrCreateStats(channelCode)
	stats.TotalCount++
	stats.SuccessCount++
	stats.latencySum += latency.Milliseconds()
	if stats.SuccessCount > 0 {
		stats.AvgLatencyMs = float64(stats.latencySum) / float64(stats.SuccessCount)
	}
}

func (m *MetricsCollector) RecordFailure(channelCode, reason string) {
	stats := m.getOrCreateStats(channelCode)
	stats.TotalCount++
	stats.FailCount++
}

func (m *MetricsCollector) RecordDuplicate(channelCode string) {
	stats := m.getOrCreateStats(channelCode)
	stats.TotalCount++
	stats.DuplicateCount++
}

func (m *MetricsCollector) RecordSignFailure(channelCode string) {
	stats := m.getOrCreateStats(channelCode)
	stats.TotalCount++
	stats.SignFailCount++
}

func (m *MetricsCollector) GetAllStats() map[string]*ChannelStats {
	return m.stats
}
