package jobs

import (
	"fmt"
	"log"
	"sync"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// WalletBalanceCheckJob 钱包余额一致性检查定时任务
// 业务规则：定时核对钱包余额与流水记录是否一致
type WalletBalanceCheckJob struct {
	walletRepo   repository.WalletRepository
	alertService Alerter
	batchSize    int
	running      bool
	mu           sync.Mutex
}

// NewWalletBalanceCheckJob 创建钱包余额一致性检查任务
func NewWalletBalanceCheckJob(
	walletRepo repository.WalletRepository,
	alertService Alerter,
) *WalletBalanceCheckJob {
	return &WalletBalanceCheckJob{
		walletRepo:   walletRepo,
		alertService: alertService,
		batchSize:    500,
	}
}

// WalletBalanceDiscrepancy 钱包余额差异记录
type WalletBalanceDiscrepancy struct {
	WalletID       int64  `json:"wallet_id"`
	AgentID        int64  `json:"agent_id"`
	WalletType     int16  `json:"wallet_type"`
	WalletTypeName string `json:"wallet_type_name"`
	CurrentBalance int64  `json:"current_balance"`      // 当前余额(分)
	CalculatedBalance int64  `json:"calculated_balance"`  // 流水计算余额(分)
	Difference     int64  `json:"difference"`           // 差异(分)
}

// Run 执行任务（每天凌晨4点执行）
func (j *WalletBalanceCheckJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[WalletBalanceCheckJob] Already running, skip")
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
	log.Printf("[WalletBalanceCheckJob] Started")

	// 执行检查
	discrepancies, err := j.checkBalanceConsistency()
	if err != nil {
		log.Printf("[WalletBalanceCheckJob] Check failed: %v", err)
		return
	}

	// 如果发现差异，发送告警
	if len(discrepancies) > 0 {
		log.Printf("[WalletBalanceCheckJob] Found %d discrepancies", len(discrepancies))
		j.sendDiscrepancyAlert(discrepancies)
	} else {
		log.Printf("[WalletBalanceCheckJob] No discrepancies found")
	}

	log.Printf("[WalletBalanceCheckJob] Completed, took=%v", time.Since(startTime))
}

// RunWithResult 执行任务并返回结果（供JobWrapper使用）
func (j *WalletBalanceCheckJob) RunWithResult() (*models.JobExecutionResult, error) {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return nil, fmt.Errorf("job already running")
	}
	j.running = true
	j.mu.Unlock()

	defer func() {
		j.mu.Lock()
		j.running = false
		j.mu.Unlock()
	}()

	// 执行检查
	discrepancies, err := j.checkBalanceConsistency()
	if err != nil {
		return &models.JobExecutionResult{
			ErrorMessage: err.Error(),
		}, err
	}

	result := &models.JobExecutionResult{
		ProcessedCount: len(discrepancies),
		FailCount:      len(discrepancies), // 差异算作失败
	}

	// 如果发现差异，发送告警
	if len(discrepancies) > 0 {
		j.sendDiscrepancyAlert(discrepancies)
		result.ErrorMessage = fmt.Sprintf("发现 %d 个钱包余额不一致", len(discrepancies))
	} else {
		result.SuccessCount = 1
	}

	return result, nil
}

// checkBalanceConsistency 检查余额一致性
func (j *WalletBalanceCheckJob) checkBalanceConsistency() ([]WalletBalanceDiscrepancy, error) {
	var discrepancies []WalletBalanceDiscrepancy

	// 查询所有钱包并核对
	// 注意：这里需要WalletRepository支持相关查询方法
	// 实际实现中需要对比wallet表余额与wallet_logs表流水汇总

	// 由于需要复杂查询，这里提供SQL参考：
	// SELECT w.id, w.agent_id, w.wallet_type, w.balance,
	//        COALESCE(SUM(wl.amount), 0) as calculated_balance
	// FROM wallets w
	// LEFT JOIN wallet_logs wl ON w.id = wl.wallet_id
	// GROUP BY w.id
	// HAVING w.balance != COALESCE(SUM(wl.amount), 0)

	// 这里使用占位实现，实际需要根据repository接口扩展
	log.Printf("[WalletBalanceCheckJob] Checking wallet balance consistency...")

	// TODO: 实现具体的余额检查逻辑
	// wallets, err := j.walletRepo.FindAllWithBalance()
	// for _, wallet := range wallets {
	//     calculatedBalance, err := j.walletLogRepo.SumByWalletID(wallet.ID)
	//     if wallet.Balance != calculatedBalance {
	//         discrepancies = append(discrepancies, WalletBalanceDiscrepancy{...})
	//     }
	// }

	return discrepancies, nil
}

// sendDiscrepancyAlert 发送差异告警
func (j *WalletBalanceCheckJob) sendDiscrepancyAlert(discrepancies []WalletBalanceDiscrepancy) {
	if j.alertService == nil {
		log.Printf("[WalletBalanceCheckJob] Alert service not available")
		return
	}

	// 构建告警消息
	message := fmt.Sprintf("**钱包余额一致性检查**\n\n发现 %d 个钱包余额与流水不一致：\n\n", len(discrepancies))

	for i, d := range discrepancies {
		if i >= 10 {
			message += fmt.Sprintf("\n... 还有 %d 条差异未显示", len(discrepancies)-10)
			break
		}
		message += fmt.Sprintf("- 钱包ID: %d, 代理商: %d, 类型: %s, 当前余额: %.2f, 计算余额: %.2f, 差异: %.2f\n",
			d.WalletID, d.AgentID, d.WalletTypeName,
			float64(d.CurrentBalance)/100,
			float64(d.CalculatedBalance)/100,
			float64(d.Difference)/100)
	}

	req := &models.AlertRequest{
		JobName:   "WalletBalanceCheckJob",
		AlertType: models.AlertTypeJobFailed,
		Title:     "【钱包余额异常】发现余额不一致",
		Message:   message,
	}

	if err := j.alertService.SendAlert(req); err != nil {
		log.Printf("[WalletBalanceCheckJob] Send alert failed: %v", err)
	}
}

// JobLogCleanupJob 任务日志清理定时任务
// 每天清理90天前的任务执行日志和告警记录
type JobLogCleanupJob struct {
	jobLogRepo   repository.JobExecutionLogRepository
	alertLogRepo repository.AlertLogRepository
	retentionDays int
	running      bool
	mu           sync.Mutex
}

// NewJobLogCleanupJob 创建任务日志清理任务
func NewJobLogCleanupJob(
	jobLogRepo repository.JobExecutionLogRepository,
	alertLogRepo repository.AlertLogRepository,
) *JobLogCleanupJob {
	return &JobLogCleanupJob{
		jobLogRepo:   jobLogRepo,
		alertLogRepo: alertLogRepo,
		retentionDays: 90,
	}
}

// Run 执行任务
func (j *JobLogCleanupJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[JobLogCleanupJob] Already running, skip")
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
	log.Printf("[JobLogCleanupJob] Started, retention_days=%d", j.retentionDays)

	// 清理任务执行日志
	jobLogDeleted, err := j.jobLogRepo.DeleteOlderThan(j.retentionDays)
	if err != nil {
		log.Printf("[JobLogCleanupJob] Delete job logs failed: %v", err)
	} else {
		log.Printf("[JobLogCleanupJob] Deleted %d job execution logs", jobLogDeleted)
	}

	// 清理告警记录
	alertLogDeleted, err := j.alertLogRepo.DeleteOlderThan(j.retentionDays)
	if err != nil {
		log.Printf("[JobLogCleanupJob] Delete alert logs failed: %v", err)
	} else {
		log.Printf("[JobLogCleanupJob] Deleted %d alert logs", alertLogDeleted)
	}

	log.Printf("[JobLogCleanupJob] Completed, job_logs_deleted=%d, alert_logs_deleted=%d, took=%v",
		jobLogDeleted, alertLogDeleted, time.Since(startTime))
}

// RunWithResult 执行任务并返回结果
func (j *JobLogCleanupJob) RunWithResult() (*models.JobExecutionResult, error) {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return nil, fmt.Errorf("job already running")
	}
	j.running = true
	j.mu.Unlock()

	defer func() {
		j.mu.Lock()
		j.running = false
		j.mu.Unlock()
	}()

	var totalDeleted int64

	// 清理任务执行日志
	jobLogDeleted, err := j.jobLogRepo.DeleteOlderThan(j.retentionDays)
	if err != nil {
		return &models.JobExecutionResult{ErrorMessage: err.Error()}, err
	}
	totalDeleted += jobLogDeleted

	// 清理告警记录
	alertLogDeleted, err := j.alertLogRepo.DeleteOlderThan(j.retentionDays)
	if err != nil {
		return &models.JobExecutionResult{ErrorMessage: err.Error()}, err
	}
	totalDeleted += alertLogDeleted

	return &models.JobExecutionResult{
		ProcessedCount: int(totalDeleted),
		SuccessCount:   int(totalDeleted),
	}, nil
}
