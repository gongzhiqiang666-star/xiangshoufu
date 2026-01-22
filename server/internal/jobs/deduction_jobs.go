package jobs

import (
	"log"
	"sync"
	"time"

	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
)

// DeductionJob 代扣定时任务
// 业务规则Q7：每日扣款 - 每天固定时间检查余额并扣款
type DeductionJob struct {
	deductionService *service.DeductionService
	running          bool
	mu               sync.Mutex
}

// NewDeductionJob 创建代扣定时任务
func NewDeductionJob(deductionService *service.DeductionService) *DeductionJob {
	return &DeductionJob{
		deductionService: deductionService,
	}
}

// Run 执行每日代扣任务
// 建议每天8:00执行，可通过调度器配置
func (j *DeductionJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[DeductionJob] Already running, skip")
		return
	}
	j.running = true
	j.mu.Unlock()

	defer func() {
		j.mu.Lock()
		j.running = false
		j.mu.Unlock()
	}()

	startTime := time.Now()
	log.Printf("[DeductionJob] Started daily deduction at %s", startTime.Format("2006-01-02 15:04:05"))

	if err := j.deductionService.ExecuteDailyDeduction(); err != nil {
		log.Printf("[DeductionJob] Daily deduction failed: %v", err)
	}

	log.Printf("[DeductionJob] Completed, took=%v", time.Since(startTime))
}

// SimCashbackJob 流量费返现定时任务（兜底处理未返现的流量费）
type SimCashbackJob struct {
	simCashbackService *service.SimCashbackService
	deviceFeeRepo      repository.DeviceFeeRepository
	batchSize          int
	running            bool
	mu                 sync.Mutex
}

// NewSimCashbackJob 创建流量费返现任务
func NewSimCashbackJob(simCashbackService *service.SimCashbackService, deviceFeeRepo repository.DeviceFeeRepository) *SimCashbackJob {
	return &SimCashbackJob{
		simCashbackService: simCashbackService,
		deviceFeeRepo:      deviceFeeRepo,
		batchSize:          100,
	}
}

// Run 执行流量费返现任务（每10分钟执行一次，作为兜底）
func (j *SimCashbackJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	defer func() {
		j.mu.Lock()
		j.running = false
		j.mu.Unlock()
	}()

	startTime := time.Now()
	log.Printf("[SimCashbackJob] Started")

	// 查询待返现的流量费记录并处理
	if j.deviceFeeRepo != nil {
		pendingFees, err := j.deviceFeeRepo.FindPendingCashback(j.batchSize)
		if err != nil {
			log.Printf("[SimCashbackJob] Find pending cashback failed: %v", err)
		} else {
			successCount := 0
			failCount := 0
			for _, fee := range pendingFees {
				if err := j.simCashbackService.ProcessSimFee(fee); err != nil {
					log.Printf("[SimCashbackJob] Process fee %d failed: %v", fee.ID, err)
					failCount++
				} else {
					successCount++
				}
			}
			log.Printf("[SimCashbackJob] Processed %d fees, success=%d, fail=%d", len(pendingFees), successCount, failCount)
		}
	}

	log.Printf("[SimCashbackJob] Completed, took=%v", time.Since(startTime))
}

// SetupDeductionJobs 设置代扣相关的定时任务
func SetupDeductionJobs(scheduler *Scheduler, deductionService *service.DeductionService, simCashbackService *service.SimCashbackService, deviceFeeRepo repository.DeviceFeeRepository) {
	// 每日代扣任务 - 每天执行（通过24小时间隔模拟）
	// 实际生产环境建议使用cron表达式调度，每天8:00执行
	deductionJob := NewDeductionJob(deductionService)
	scheduler.AddJob("daily_deduction", 24*time.Hour, deductionJob.Run)

	// 流量费返现兜底任务 - 每10分钟执行
	if simCashbackService != nil {
		simCashbackJob := NewSimCashbackJob(simCashbackService, deviceFeeRepo)
		scheduler.AddJob("sim_cashback_fallback", 10*time.Minute, simCashbackJob.Run)
	}

	log.Printf("[Jobs] Deduction related jobs registered")
}

// 添加每日定时任务调度支持
// DailyScheduler 每日定时调度器（用于精确控制每天固定时间执行）
type DailyScheduler struct {
	hour   int
	minute int
	job    func()
	stopCh chan struct{}
}

// NewDailyScheduler 创建每日定时调度器
func NewDailyScheduler(hour, minute int, job func()) *DailyScheduler {
	return &DailyScheduler{
		hour:   hour,
		minute: minute,
		job:    job,
		stopCh: make(chan struct{}),
	}
}

// Start 启动调度器
func (d *DailyScheduler) Start() {
	go d.run()
	log.Printf("[DailyScheduler] Started, will run at %02d:%02d every day", d.hour, d.minute)
}

// Stop 停止调度器
func (d *DailyScheduler) Stop() {
	close(d.stopCh)
}

// run 调度循环
func (d *DailyScheduler) run() {
	for {
		// 计算下次执行时间
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), d.hour, d.minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}

		waitDuration := next.Sub(now)
		log.Printf("[DailyScheduler] Next run at %s, wait %v", next.Format("2006-01-02 15:04:05"), waitDuration)

		select {
		case <-d.stopCh:
			return
		case <-time.After(waitDuration):
			log.Printf("[DailyScheduler] Executing scheduled job at %s", time.Now().Format("2006-01-02 15:04:05"))
			d.job()
		}
	}
}
