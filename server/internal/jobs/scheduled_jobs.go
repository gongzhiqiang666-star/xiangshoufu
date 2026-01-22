package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm"
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
	callbackRepo  *repository.GormRawCallbackRepository
	archivePath   string
	retentionDays int
	batchSize     int
	running       bool
	mu            sync.Mutex
}

// NewDataArchiverJob 创建数据归档任务
func NewDataArchiverJob(callbackRepo *repository.GormRawCallbackRepository, archivePath string) *DataArchiverJob {
	// 确保归档目录存在
	if archivePath != "" {
		os.MkdirAll(archivePath, 0755)
	}
	return &DataArchiverJob{
		callbackRepo:  callbackRepo,
		archivePath:   archivePath,
		retentionDays: 90, // 保留90天
		batchSize:     1000,
	}
}

// Run 执行任务（每天凌晨4点执行）
func (j *DataArchiverJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[DataArchiverJob] Already running, skip")
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
	log.Printf("[DataArchiverJob] Started")

	// 1. 统计需要归档的数据量
	totalCount, err := j.callbackRepo.CountArchivedLogs(j.retentionDays)
	if err != nil {
		log.Printf("[DataArchiverJob] Count archived logs failed: %v", err)
		return
	}

	if totalCount == 0 {
		log.Printf("[DataArchiverJob] No data to archive")
		return
	}

	log.Printf("[DataArchiverJob] Found %d logs to archive", totalCount)

	archivedCount := int64(0)
	deletedCount := int64(0)

	// 2. 分批处理归档
	for {
		// 查询待归档数据
		logs, err := j.callbackRepo.FindArchivedLogs(j.retentionDays, j.batchSize)
		if err != nil {
			log.Printf("[DataArchiverJob] Find archived logs failed: %v", err)
			break
		}

		if len(logs) == 0 {
			break
		}

		// 3. 导出到归档文件（如果配置了归档路径）
		if j.archivePath != "" {
			archiveFile := filepath.Join(j.archivePath, fmt.Sprintf("callbacks_%s.jsonl", time.Now().Format("20060102_150405")))
			if err := j.exportToFile(logs, archiveFile); err != nil {
				log.Printf("[DataArchiverJob] Export to file failed: %v", err)
				break
			}
			archivedCount += int64(len(logs))
		}

		// 4. 删除已归档的数据
		ids := make([]int64, len(logs))
		for i, l := range logs {
			ids[i] = l.ID
		}

		deleted, err := j.callbackRepo.DeleteArchivedLogs(ids)
		if err != nil {
			log.Printf("[DataArchiverJob] Delete archived logs failed: %v", err)
			break
		}
		deletedCount += deleted

		// 进度日志
		log.Printf("[DataArchiverJob] Progress: archived=%d, deleted=%d", archivedCount, deletedCount)
	}

	log.Printf("[DataArchiverJob] Completed: archived=%d, deleted=%d, took=%v",
		archivedCount, deletedCount, time.Since(startTime))
}

// exportToFile 导出日志到文件（JSONL格式）
func (j *DataArchiverJob) exportToFile(logs []*repository.RawCallbackLog, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open archive file failed: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, l := range logs {
		if err := encoder.Encode(l); err != nil {
			return fmt.Errorf("encode log failed: %w", err)
		}
	}

	return nil
}

// PartitionManagerJob 分区管理定时任务
type PartitionManagerJob struct {
	db      *gorm.DB
	running bool
	mu      sync.Mutex
}

// NewPartitionManagerJob 创建分区管理任务
func NewPartitionManagerJob(db *gorm.DB) *PartitionManagerJob {
	return &PartitionManagerJob{db: db}
}

// Run 执行任务（每月1号凌晨1点执行）
func (j *PartitionManagerJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[PartitionManagerJob] Already running, skip")
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
	log.Printf("[PartitionManagerJob] Started")

	// 创建未来3个月的分区（确保始终有足够的分区）
	for i := 1; i <= 3; i++ {
		futureMonth := time.Now().AddDate(0, i, 0)
		if err := j.createMonthlyPartition(futureMonth); err != nil {
			log.Printf("[PartitionManagerJob] Create partition for %s failed: %v",
				futureMonth.Format("2006-01"), err)
		} else {
			log.Printf("[PartitionManagerJob] Created/verified partition for %s",
				futureMonth.Format("2006-01"))
		}
	}

	log.Printf("[PartitionManagerJob] Completed, took=%v", time.Since(startTime))
}

// createMonthlyPartition 创建月度分区
func (j *PartitionManagerJob) createMonthlyPartition(month time.Time) error {
	if j.db == nil {
		log.Printf("[PartitionManagerJob] Database connection not available, skip partition creation")
		return nil
	}

	year := month.Year()
	mon := int(month.Month())

	// 分区表名
	partitionName := fmt.Sprintf("raw_callback_logs_%d_%02d", year, mon)

	// 计算分区范围
	startDate := time.Date(year, time.Month(mon), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// 检查分区是否已存在
	var exists int64
	checkSQL := `
		SELECT COUNT(*) FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relname = ? AND n.nspname = 'public'
	`
	if err := j.db.Raw(checkSQL, partitionName).Scan(&exists).Error; err != nil {
		return fmt.Errorf("check partition exists failed: %w", err)
	}

	if exists > 0 {
		log.Printf("[PartitionManagerJob] Partition %s already exists", partitionName)
		return nil
	}

	// 创建分区（使用原生SQL）
	createSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s PARTITION OF raw_callback_logs
		FOR VALUES FROM ('%s') TO ('%s')
	`, partitionName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	if err := j.db.Exec(createSQL).Error; err != nil {
		// 如果分区已存在或表不是分区表，忽略错误
		log.Printf("[PartitionManagerJob] Create partition SQL result: %v", err)
		return nil
	}

	log.Printf("[PartitionManagerJob] Successfully created partition: %s", partitionName)
	return nil
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
