package jobs

import (
	"log"
	"sync"
	"time"

	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
)

// Note: service.MerchantService is used by MerchantTypeCalculatorJob

// ProfitCalculatorJob 分润计算定时任务（兜底重试）
type ProfitCalculatorJob struct {
	transactionRepo repository.TransactionRepository
	profitService   *service.ProfitService
	batchSize       int
	running         bool
	mu              sync.Mutex
}

// NewProfitCalculatorJob 创建分润计算任务
func NewProfitCalculatorJob(
	transactionRepo repository.TransactionRepository,
	profitService *service.ProfitService,
) *ProfitCalculatorJob {
	return &ProfitCalculatorJob{
		transactionRepo: transactionRepo,
		profitService:   profitService,
		batchSize:       500,
	}
}

// Run 执行任务（每5分钟执行一次）
func (j *ProfitCalculatorJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[ProfitCalculatorJob] Already running, skip")
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
	log.Printf("[ProfitCalculatorJob] Started")

	// 查询待计算的交易
	transactions, err := j.transactionRepo.FindUnprocessedProfit(j.batchSize)
	if err != nil {
		log.Printf("[ProfitCalculatorJob] Find unprocessed failed: %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("[ProfitCalculatorJob] No pending transactions")
		return
	}

	successCount := 0
	failCount := 0

	for _, tx := range transactions {
		if err := j.profitService.CalculateProfit(tx.ID); err != nil {
			log.Printf("[ProfitCalculatorJob] Calculate failed for tx %d: %v", tx.ID, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("[ProfitCalculatorJob] Completed: success=%d, fail=%d, took=%v",
		successCount, failCount, time.Since(startTime))
}

// CallbackRetryJob 回调重试定时任务
type CallbackRetryJob struct {
	callbackRepo repository.RawCallbackRepository
	processor    *service.CallbackProcessor
	maxRetry     int
	batchSize    int
	running      bool
	mu           sync.Mutex
}

// NewCallbackRetryJob 创建回调重试任务
func NewCallbackRetryJob(
	callbackRepo repository.RawCallbackRepository,
	processor *service.CallbackProcessor,
) *CallbackRetryJob {
	return &CallbackRetryJob{
		callbackRepo: callbackRepo,
		processor:    processor,
		maxRetry:     3,
		batchSize:    100,
	}
}

// Run 执行任务（每5分钟执行一次）
func (j *CallbackRetryJob) Run() {
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
	log.Printf("[CallbackRetryJob] Started")

	// 查询失败且未超过最大重试次数的回调
	logs, err := j.callbackRepo.FindFailedLogs(j.maxRetry, j.batchSize)
	if err != nil {
		log.Printf("[CallbackRetryJob] Find failed logs error: %v", err)
		return
	}

	if len(logs) == 0 {
		log.Printf("[CallbackRetryJob] No failed callbacks to retry")
		return
	}

	successCount := 0
	failCount := 0

	for _, logEntry := range logs {
		err := j.processor.ProcessCallback(
			logEntry.ID,
			logEntry.ChannelCode,
			logEntry.ActionType,
			[]byte(logEntry.RawRequest),
		)
		if err != nil {
			log.Printf("[CallbackRetryJob] Retry failed for log %d: %v", logEntry.ID, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("[CallbackRetryJob] Completed: success=%d, fail=%d, took=%v",
		successCount, failCount, time.Since(startTime))
}

// MessageCleanupJob 消息清理定时任务
type MessageCleanupJob struct {
	messageService *service.MessageService
}

// NewMessageCleanupJob 创建消息清理任务
func NewMessageCleanupJob(messageService *service.MessageService) *MessageCleanupJob {
	return &MessageCleanupJob{
		messageService: messageService,
	}
}

// Run 执行任务（每天凌晨3点执行）
func (j *MessageCleanupJob) Run() {
	startTime := time.Now()
	log.Printf("[MessageCleanupJob] Started")

	count, err := j.messageService.CleanupExpiredMessages()
	if err != nil {
		log.Printf("[MessageCleanupJob] Cleanup failed: %v", err)
		return
	}

	log.Printf("[MessageCleanupJob] Completed: deleted=%d, took=%v", count, time.Since(startTime))
}

// DataArchiverJob 数据归档定时任务
type DataArchiverJob struct {
	callbackRepo  repository.RawCallbackRepository
	archivePath   string
	retentionDays int
}

// NewDataArchiverJob 创建数据归档任务
func NewDataArchiverJob(callbackRepo repository.RawCallbackRepository, archivePath string) *DataArchiverJob {
	return &DataArchiverJob{
		callbackRepo:  callbackRepo,
		archivePath:   archivePath,
		retentionDays: 90, // 保留90天
	}
}

// Run 执行任务（每天凌晨4点执行）
func (j *DataArchiverJob) Run() {
	startTime := time.Now()
	log.Printf("[DataArchiverJob] Started")

	// TODO: 实现数据归档逻辑
	// 1. 查询超过retentionDays的数据
	// 2. 导出到archivePath目录
	// 3. 删除已归档的数据

	log.Printf("[DataArchiverJob] Completed, took=%v", time.Since(startTime))
}

// PartitionManagerJob 分区管理定时任务
type PartitionManagerJob struct {
	// TODO: 添加数据库连接
}

// NewPartitionManagerJob 创建分区管理任务
func NewPartitionManagerJob() *PartitionManagerJob {
	return &PartitionManagerJob{}
}

// Run 执行任务（每月1号凌晨1点执行）
func (j *PartitionManagerJob) Run() {
	startTime := time.Now()
	log.Printf("[PartitionManagerJob] Started")

	// TODO: 实现自动创建下个月分区的逻辑
	// CREATE TABLE raw_callback_logs_YYYY_MM PARTITION OF raw_callback_logs
	// FOR VALUES FROM ('YYYY-MM-01') TO ('YYYY-MM+1-01')

	log.Printf("[PartitionManagerJob] Completed, took=%v", time.Since(startTime))
}

// MerchantTypeCalculatorJob 商户类型计算定时任务
// 每天凌晨2点执行，根据最近30天交易额计算商户类型
type MerchantTypeCalculatorJob struct {
	merchantRepo    *repository.GormMerchantRepository
	merchantService *service.MerchantService
	batchSize       int
	running         bool
	mu              sync.Mutex
}

// NewMerchantTypeCalculatorJob 创建商户类型计算任务
func NewMerchantTypeCalculatorJob(
	merchantRepo *repository.GormMerchantRepository,
	merchantService *service.MerchantService,
) *MerchantTypeCalculatorJob {
	return &MerchantTypeCalculatorJob{
		merchantRepo:    merchantRepo,
		merchantService: merchantService,
		batchSize:       500,
	}
}

// Run 执行任务
func (j *MerchantTypeCalculatorJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[MerchantTypeCalculatorJob] Already running, skip")
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
	log.Printf("[MerchantTypeCalculatorJob] Started")

	// 统计总商户数
	totalMerchants, err := j.merchantRepo.CountAllActiveMerchants()
	if err != nil {
		log.Printf("[MerchantTypeCalculatorJob] Count merchants failed: %v", err)
		return
	}

	if totalMerchants == 0 {
		log.Printf("[MerchantTypeCalculatorJob] No active merchants to process")
		return
	}

	successCount := 0
	failCount := 0
	offset := 0

	// 分批处理所有商户
	for {
		merchantIDs, err := j.merchantRepo.FindAllMerchantIDs(j.batchSize, offset)
		if err != nil {
			log.Printf("[MerchantTypeCalculatorJob] Find merchant IDs failed: %v", err)
			break
		}

		if len(merchantIDs) == 0 {
			break
		}

		// 逐个计算商户类型
		for _, merchantID := range merchantIDs {
			_, err := j.merchantService.CalculateMerchantType(merchantID)
			if err != nil {
				log.Printf("[MerchantTypeCalculatorJob] Calculate type failed for merchant %d: %v", merchantID, err)
				failCount++
			} else {
				successCount++
			}
		}

		offset += j.batchSize

		// 进度日志（每处理1000个记录一次）
		if offset%1000 == 0 {
			log.Printf("[MerchantTypeCalculatorJob] Progress: %d/%d merchants processed", offset, totalMerchants)
		}
	}

	log.Printf("[MerchantTypeCalculatorJob] Completed: total=%d, success=%d, fail=%d, took=%v",
		totalMerchants, successCount, failCount, time.Since(startTime))
}
