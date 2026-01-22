package jobs

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// Alerter 告警服务接口
// 用于打破 jobs 和 service 包之间的循环依赖
type Alerter interface {
	CheckAndSendConsecutiveFailAlert(jobName string, threshold int) error
	SendAlert(req *models.AlertRequest) error
}

// JobFunc 任务执行函数类型
// 返回执行结果和错误
type JobFunc func() (*models.JobExecutionResult, error)

// JobWrapper 任务包装器
// 统一处理日志记录、重试、告警等逻辑
type JobWrapper struct {
	jobName         string
	jobFunc         JobFunc
	logRepo         repository.JobExecutionLogRepository
	configRepo      repository.JobConfigRepository
	failCounterRepo repository.JobFailCounterRepository
	alertService    Alerter
	retryStrategy   RetryStrategy
	config          *models.JobConfig
	running         bool
	mu              sync.Mutex
}

// JobWrapperOption 任务包装器选项
type JobWrapperOption func(*JobWrapper)

// WithRetryStrategy 设置重试策略
func WithRetryStrategy(strategy RetryStrategy) JobWrapperOption {
	return func(w *JobWrapper) {
		w.retryStrategy = strategy
	}
}

// WithConfig 设置任务配置
func WithConfig(config *models.JobConfig) JobWrapperOption {
	return func(w *JobWrapper) {
		w.config = config
	}
}

// NewJobWrapper 创建任务包装器
func NewJobWrapper(
	jobName string,
	jobFunc JobFunc,
	logRepo repository.JobExecutionLogRepository,
	configRepo repository.JobConfigRepository,
	failCounterRepo repository.JobFailCounterRepository,
	alertService Alerter,
	opts ...JobWrapperOption,
) *JobWrapper {
	w := &JobWrapper{
		jobName:         jobName,
		jobFunc:         jobFunc,
		logRepo:         logRepo,
		configRepo:      configRepo,
		failCounterRepo: failCounterRepo,
		alertService:    alertService,
		retryStrategy:   NewExponentialBackoffStrategy(DefaultRetryIntervals),
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Run 执行任务（由调度器调用）
func (w *JobWrapper) Run() {
	w.RunWithTriggerType(models.JobTriggerTypeAuto)
}

// RunManual 手动触发执行
func (w *JobWrapper) RunManual() {
	w.RunWithTriggerType(models.JobTriggerTypeManual)
}

// RunWithTriggerType 带触发类型执行
func (w *JobWrapper) RunWithTriggerType(triggerType int16) {
	// 防止重复执行
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		log.Printf("[%s] Already running, skip", w.jobName)
		return
	}
	w.running = true
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.running = false
		w.mu.Unlock()
	}()

	// 加载最新配置
	if err := w.loadConfig(); err != nil {
		log.Printf("[%s] Load config failed: %v", w.jobName, err)
	}

	// 检查是否启用
	if w.config != nil && !w.config.IsEnabled {
		log.Printf("[%s] Job is disabled, skip", w.jobName)
		return
	}

	// 创建执行日志
	logEntry := w.startLog(triggerType)

	// 执行任务（带重试）
	result, err := w.executeWithRetry(logEntry)

	// 结束日志
	w.endLog(logEntry, result, err)

	// 处理执行结果
	if err != nil {
		w.handleFailure(logEntry, err)
	} else {
		w.handleSuccess()
	}
}

// loadConfig 加载任务配置
func (w *JobWrapper) loadConfig() error {
	if w.configRepo == nil {
		return nil
	}

	config, err := w.configRepo.FindByName(w.jobName)
	if err != nil {
		return err
	}

	w.config = config

	// 根据配置更新重试策略
	if config.MaxRetries > 0 {
		intervals := CreateRetryIntervalsFromConfig(config.RetryInterval, config.MaxRetries)
		w.retryStrategy = NewExponentialBackoffStrategy(intervals)
	}

	return nil
}

// startLog 开始执行日志
func (w *JobWrapper) startLog(triggerType int16) *models.JobExecutionLog {
	logEntry := &models.JobExecutionLog{
		JobName:     w.jobName,
		StartedAt:   time.Now(),
		Status:      models.JobStatusRunning,
		TriggerType: triggerType,
		CreatedAt:   time.Now(),
	}

	if w.logRepo != nil {
		if err := w.logRepo.Create(logEntry); err != nil {
			log.Printf("[%s] Create log entry failed: %v", w.jobName, err)
		}
	}

	log.Printf("[%s] Started at %s", w.jobName, logEntry.StartedAt.Format("2006-01-02 15:04:05"))
	return logEntry
}

// executeWithRetry 带重试执行任务
func (w *JobWrapper) executeWithRetry(logEntry *models.JobExecutionLog) (*models.JobExecutionResult, error) {
	maxRetries := 3
	if w.config != nil && w.config.MaxRetries > 0 {
		maxRetries = w.config.MaxRetries
	}

	var lastErr error
	var lastResult *models.JobExecutionResult

	for retryCount := 0; retryCount <= maxRetries; retryCount++ {
		// 非首次执行，等待重试间隔
		if retryCount > 0 {
			interval := w.retryStrategy.GetRetryInterval(retryCount)
			log.Printf("[%s] Retry %d/%d, waiting %v...", w.jobName, retryCount, maxRetries, interval)
			time.Sleep(interval)
		}

		// 执行任务
		result, err := w.safeExecute()
		lastResult = result
		logEntry.RetryCount = retryCount

		if err == nil {
			return result, nil
		}

		lastErr = err
		log.Printf("[%s] Execution failed (attempt %d/%d): %v", w.jobName, retryCount+1, maxRetries+1, err)

		// 检查是否应该继续重试
		if !w.retryStrategy.ShouldRetry(retryCount, maxRetries) {
			break
		}
	}

	return lastResult, lastErr
}

// safeExecute 安全执行任务（捕获panic）
func (w *JobWrapper) safeExecute() (result *models.JobExecutionResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("panic: %v\nstack: %s", r, stack)
			log.Printf("[%s] Panic recovered: %v", w.jobName, r)
		}
	}()

	return w.jobFunc()
}

// endLog 结束执行日志
func (w *JobWrapper) endLog(logEntry *models.JobExecutionLog, result *models.JobExecutionResult, err error) {
	now := time.Now()
	logEntry.EndedAt = &now
	logEntry.DurationMs = now.Sub(logEntry.StartedAt).Milliseconds()

	if err != nil {
		logEntry.Status = models.JobStatusFailed
		logEntry.ErrorMessage = err.Error()
		// 尝试获取堆栈信息
		if stackErr, ok := err.(interface{ Stack() string }); ok {
			logEntry.ErrorStack = stackErr.Stack()
		}
	} else {
		logEntry.Status = models.JobStatusSuccess
	}

	// 更新处理计数
	if result != nil {
		logEntry.ProcessedCount = result.ProcessedCount
		logEntry.SuccessCount = result.SuccessCount
		logEntry.FailCount = result.FailCount
		if result.ErrorMessage != "" && logEntry.ErrorMessage == "" {
			logEntry.ErrorMessage = result.ErrorMessage
		}
	}

	// 更新日志
	if w.logRepo != nil {
		if err := w.logRepo.Update(logEntry); err != nil {
			log.Printf("[%s] Update log entry failed: %v", w.jobName, err)
		}
	}

	log.Printf("[%s] Completed: status=%s, duration=%dms, processed=%d, success=%d, fail=%d",
		w.jobName,
		models.GetJobStatusName(logEntry.Status),
		logEntry.DurationMs,
		logEntry.ProcessedCount,
		logEntry.SuccessCount,
		logEntry.FailCount)
}

// handleFailure 处理任务失败
func (w *JobWrapper) handleFailure(logEntry *models.JobExecutionLog, err error) {
	// 增加失败计数
	if w.failCounterRepo != nil {
		counter, counterErr := w.failCounterRepo.IncrementFail(w.jobName)
		if counterErr != nil {
			log.Printf("[%s] Increment fail counter failed: %v", w.jobName, counterErr)
		} else {
			log.Printf("[%s] Consecutive fails: %d", w.jobName, counter.ConsecutiveFails)

			// 检查是否需要发送连续失败告警
			if w.alertService != nil && w.config != nil {
				threshold := w.config.AlertThreshold
				if threshold <= 0 {
					threshold = 3
				}
				if counter.ConsecutiveFails >= threshold {
					if alertErr := w.alertService.CheckAndSendConsecutiveFailAlert(w.jobName, threshold); alertErr != nil {
						log.Printf("[%s] Send consecutive fail alert failed: %v", w.jobName, alertErr)
					}
				}
			}
		}
	}
}

// handleSuccess 处理任务成功
func (w *JobWrapper) handleSuccess() {
	// 重置失败计数
	if w.failCounterRepo != nil {
		if err := w.failCounterRepo.ResetOnSuccess(w.jobName); err != nil {
			log.Printf("[%s] Reset fail counter failed: %v", w.jobName, err)
		}
	}
}

// IsRunning 检查任务是否正在运行
func (w *JobWrapper) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// GetConfig 获取任务配置
func (w *JobWrapper) GetConfig() *models.JobConfig {
	return w.config
}

// GetJobName 获取任务名称
func (w *JobWrapper) GetJobName() string {
	return w.jobName
}
