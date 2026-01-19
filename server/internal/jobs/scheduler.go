package jobs

import (
	"log"
	"sync"
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	jobs    map[string]*ScheduledJob
	running bool
	stopCh  chan struct{}
	mu      sync.RWMutex
}

// ScheduledJob 定时任务
type ScheduledJob struct {
	Name     string
	Interval time.Duration
	RunFunc  func()
	LastRun  time.Time
	NextRun  time.Time
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs:   make(map[string]*ScheduledJob),
		stopCh: make(chan struct{}),
	}
}

// AddJob 添加定时任务
func (s *Scheduler) AddJob(name string, interval time.Duration, runFunc func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[name] = &ScheduledJob{
		Name:     name,
		Interval: interval,
		RunFunc:  runFunc,
		NextRun:  time.Now().Add(interval),
	}

	log.Printf("[Scheduler] Added job: %s, interval: %v", name, interval)
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.run()
	log.Printf("[Scheduler] Started with %d jobs", len(s.jobs))
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopCh)
	log.Printf("[Scheduler] Stopped")
}

// run 调度循环
func (s *Scheduler) run() {
	ticker := time.NewTicker(10 * time.Second) // 每10秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case now := <-ticker.C:
			s.checkAndRunJobs(now)
		}
	}
}

// checkAndRunJobs 检查并执行到期的任务
func (s *Scheduler) checkAndRunJobs(now time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, job := range s.jobs {
		if now.After(job.NextRun) {
			go s.runJob(job)
		}
	}
}

// runJob 执行单个任务
func (s *Scheduler) runJob(job *ScheduledJob) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Scheduler] Job %s panicked: %v", job.Name, r)
		}
	}()

	startTime := time.Now()
	job.RunFunc()
	job.LastRun = startTime
	job.NextRun = startTime.Add(job.Interval)
}

// GetJobStatus 获取任务状态
func (s *Scheduler) GetJobStatus() map[string]JobStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := make(map[string]JobStatus)
	for name, job := range s.jobs {
		status[name] = JobStatus{
			Name:     job.Name,
			Interval: job.Interval.String(),
			LastRun:  job.LastRun,
			NextRun:  job.NextRun,
		}
	}
	return status
}

// JobStatus 任务状态
type JobStatus struct {
	Name     string    `json:"name"`
	Interval string    `json:"interval"`
	LastRun  time.Time `json:"last_run"`
	NextRun  time.Time `json:"next_run"`
}
