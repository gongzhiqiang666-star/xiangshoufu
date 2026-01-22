package jobs

import (
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// JobManagerService 任务管理服务
// 负责任务的注册、管理和调度集成
type JobManagerService struct {
	configRepo      repository.JobConfigRepository
	logRepo         repository.JobExecutionLogRepository
	failCounterRepo repository.JobFailCounterRepository
	alertService    Alerter
	scheduler       *Scheduler
	wrappers        map[string]*JobWrapper
}

// NewJobManagerService 创建任务管理服务
func NewJobManagerService(
	configRepo repository.JobConfigRepository,
	logRepo repository.JobExecutionLogRepository,
	failCounterRepo repository.JobFailCounterRepository,
	alertService Alerter,
) *JobManagerService {
	return &JobManagerService{
		configRepo:      configRepo,
		logRepo:         logRepo,
		failCounterRepo: failCounterRepo,
		alertService:    alertService,
		wrappers:        make(map[string]*JobWrapper),
	}
}

// SetScheduler 设置调度器
func (s *JobManagerService) SetScheduler(scheduler *Scheduler) {
	s.scheduler = scheduler
}

// RegisterJob 注册任务
// jobName: 任务名称
// jobFunc: 任务执行函数
// defaultInterval: 默认执行间隔
func (s *JobManagerService) RegisterJob(jobName string, jobFunc JobFunc, defaultInterval time.Duration) *JobWrapper {
	// 创建任务包装器
	wrapper := NewJobWrapper(
		jobName,
		jobFunc,
		s.logRepo,
		s.configRepo,
		s.failCounterRepo,
		s.alertService,
	)

	s.wrappers[jobName] = wrapper

	// 如果有调度器，添加到调度器
	if s.scheduler != nil {
		// 尝试从配置获取间隔
		interval := defaultInterval
		if config, err := s.configRepo.FindByName(jobName); err == nil && config.IntervalSeconds > 0 {
			interval = time.Duration(config.IntervalSeconds) * time.Second
		}

		s.scheduler.AddJob(jobName, interval, wrapper.Run)
	}

	log.Printf("[JobManagerService] Registered job: %s", jobName)
	return wrapper
}

// GetWrapper 获取任务包装器
func (s *JobManagerService) GetWrapper(jobName string) *JobWrapper {
	return s.wrappers[jobName]
}

// GetAllWrappers 获取所有任务包装器
func (s *JobManagerService) GetAllWrappers() map[string]*JobWrapper {
	return s.wrappers
}

// TriggerJob 手动触发任务
func (s *JobManagerService) TriggerJob(jobName string) error {
	wrapper, ok := s.wrappers[jobName]
	if !ok {
		return nil
	}

	go wrapper.RunManual()
	return nil
}

// GetJobStatus 获取任务状态
func (s *JobManagerService) GetJobStatus(jobName string) (*JobStatusInfo, error) {
	wrapper, ok := s.wrappers[jobName]
	if !ok {
		return nil, nil
	}

	config, _ := s.configRepo.FindByName(jobName)
	latestLog, _ := s.logRepo.FindLatestByJobName(jobName)

	status := &JobStatusInfo{
		JobName:   jobName,
		IsRunning: wrapper.IsRunning(),
	}

	if config != nil {
		status.IsEnabled = config.IsEnabled
		status.IntervalSeconds = config.IntervalSeconds
	}

	if latestLog != nil {
		status.LastRunAt = &latestLog.StartedAt
		status.LastStatus = latestLog.Status
		status.LastDurationMs = latestLog.DurationMs
	}

	return status, nil
}

// JobStatusInfo 任务状态信息
type JobStatusInfo struct {
	JobName         string     `json:"job_name"`
	IsEnabled       bool       `json:"is_enabled"`
	IsRunning       bool       `json:"is_running"`
	IntervalSeconds int        `json:"interval_seconds"`
	LastRunAt       *time.Time `json:"last_run_at"`
	LastStatus      int16      `json:"last_status"`
	LastDurationMs  int64      `json:"last_duration_ms"`
}

// AdaptJobFunc 将旧版任务函数适配为新版JobFunc
// 用于兼容现有的不返回结果的任务
func AdaptJobFunc(oldJobFunc func()) JobFunc {
	return func() (*models.JobExecutionResult, error) {
		oldJobFunc()
		return &models.JobExecutionResult{
			SuccessCount: 1,
		}, nil
	}
}

// AdaptJobFuncWithCounts 将返回计数的任务函数适配为JobFunc
func AdaptJobFuncWithCounts(jobFunc func() (processed, success, fail int, err error)) JobFunc {
	return func() (*models.JobExecutionResult, error) {
		processed, success, fail, err := jobFunc()
		return &models.JobExecutionResult{
			ProcessedCount: processed,
			SuccessCount:   success,
			FailCount:      fail,
		}, err
	}
}

// AdaptJobFuncWithError 将可能返回错误的任务函数适配为JobFunc
func AdaptJobFuncWithError(jobFunc func() error) JobFunc {
	return func() (*models.JobExecutionResult, error) {
		err := jobFunc()
		if err != nil {
			return &models.JobExecutionResult{
				FailCount:    1,
				ErrorMessage: err.Error(),
			}, err
		}
		return &models.JobExecutionResult{
			SuccessCount: 1,
		}, nil
	}
}
