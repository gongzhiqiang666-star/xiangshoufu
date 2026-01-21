package jobs

import (
	"log"
	"sync"
	"time"

	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
)

// RewardCheckJob 激活奖励检查定时任务
// 业务规则Q33：每日凌晨检查所有终端的激活奖励条件
type RewardCheckJob struct {
	terminalRepo  repository.TerminalRepository
	channelRepo   repository.ChannelRepository
	rewardService *service.ActivationRewardService
	batchSize     int
	running       bool
	mu            sync.Mutex
}

// NewRewardCheckJob 创建激活奖励检查任务
func NewRewardCheckJob(
	terminalRepo repository.TerminalRepository,
	channelRepo repository.ChannelRepository,
	rewardService *service.ActivationRewardService,
) *RewardCheckJob {
	return &RewardCheckJob{
		terminalRepo:  terminalRepo,
		channelRepo:   channelRepo,
		rewardService: rewardService,
		batchSize:     500,
	}
}

// Run 执行任务（每天凌晨2点执行）
func (j *RewardCheckJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[RewardCheckJob] Already running, skip")
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
	checkDate := time.Now().Truncate(24 * time.Hour) // 当天0点
	log.Printf("[RewardCheckJob] Started for date: %s", checkDate.Format("2006-01-02"))

	// 获取所有启用的通道
	channels, err := j.channelRepo.FindAllActive()
	if err != nil {
		log.Printf("[RewardCheckJob] Find channels failed: %v", err)
		return
	}

	totalProcessed := 0
	totalErrors := 0

	// 逐个通道处理
	for _, channel := range channels {
		log.Printf("[RewardCheckJob] Processing channel: %d (%s)", channel.ID, channel.ChannelCode)

		err := j.rewardService.BatchCheckTerminalRewards(channel.ID, checkDate)
		if err != nil {
			log.Printf("[RewardCheckJob] Channel %d check failed: %v", channel.ID, err)
			totalErrors++
		} else {
			totalProcessed++
		}
	}

	log.Printf("[RewardCheckJob] Completed: channels_processed=%d, errors=%d, took=%v",
		totalProcessed, totalErrors, time.Since(startTime))
}

// DepositCashbackJob 押金返现处理定时任务
// 业务规则Q32：检查待处理的押金返现记录并入账
type DepositCashbackJob struct {
	depositRecordRepo *repository.GormDepositCashbackRecordRepository
	walletRepo        repository.WalletRepository
	batchSize         int
	running           bool
	mu                sync.Mutex
}

// NewDepositCashbackJob 创建押金返现处理任务
func NewDepositCashbackJob(
	depositRecordRepo *repository.GormDepositCashbackRecordRepository,
	walletRepo repository.WalletRepository,
) *DepositCashbackJob {
	return &DepositCashbackJob{
		depositRecordRepo: depositRecordRepo,
		walletRepo:        walletRepo,
		batchSize:         200,
	}
}

// Run 执行任务（每10分钟执行一次）
func (j *DepositCashbackJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[DepositCashbackJob] Already running, skip")
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
	log.Printf("[DepositCashbackJob] Started")

	// 查询待处理的押金返现记录
	records, err := j.depositRecordRepo.FindPendingRecords(j.batchSize)
	if err != nil {
		log.Printf("[DepositCashbackJob] Find pending records failed: %v", err)
		return
	}

	if len(records) == 0 {
		log.Printf("[DepositCashbackJob] No pending records")
		return
	}

	successCount := 0
	failCount := 0

	// 批量处理
	walletUpdates := make(map[int64]int64)
	for _, record := range records {
		// 获取钱包
		wallet, err := j.walletRepo.FindByAgentAndType(record.AgentID, record.ChannelID, record.WalletType)
		if err != nil {
			log.Printf("[DepositCashbackJob] Find wallet failed for agent %d: %v", record.AgentID, err)
			failCount++
			continue
		}

		walletUpdates[wallet.ID] += record.ActualCashback
		successCount++
	}

	// 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := j.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			log.Printf("[DepositCashbackJob] Batch update balance failed: %v", err)
		} else {
			// 更新记录状态为已入账
			for _, record := range records {
				j.depositRecordRepo.UpdateWalletStatus(record.ID, 1)
			}
		}
	}

	log.Printf("[DepositCashbackJob] Completed: success=%d, fail=%d, took=%v",
		successCount, failCount, time.Since(startTime))
}

// ActivationRewardSettleJob 激活奖励入账定时任务
// 将待入账的激活奖励记录批量入账到奖励钱包
type ActivationRewardSettleJob struct {
	rewardRecordRepo *repository.GormActivationRewardRecordRepository
	walletRepo       repository.WalletRepository
	batchSize        int
	running          bool
	mu               sync.Mutex
}

// NewActivationRewardSettleJob 创建激活奖励入账任务
func NewActivationRewardSettleJob(
	rewardRecordRepo *repository.GormActivationRewardRecordRepository,
	walletRepo repository.WalletRepository,
) *ActivationRewardSettleJob {
	return &ActivationRewardSettleJob{
		rewardRecordRepo: rewardRecordRepo,
		walletRepo:       walletRepo,
		batchSize:        200,
	}
}

// Run 执行任务（每10分钟执行一次）
func (j *ActivationRewardSettleJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[ActivationRewardSettleJob] Already running, skip")
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
	log.Printf("[ActivationRewardSettleJob] Started")

	// 查询待处理的激活奖励记录
	records, err := j.rewardRecordRepo.FindPendingRecords(j.batchSize)
	if err != nil {
		log.Printf("[ActivationRewardSettleJob] Find pending records failed: %v", err)
		return
	}

	if len(records) == 0 {
		log.Printf("[ActivationRewardSettleJob] No pending records")
		return
	}

	successCount := 0
	failCount := 0

	// 批量处理
	walletUpdates := make(map[int64]int64)
	for _, record := range records {
		// 获取奖励钱包
		wallet, err := j.walletRepo.FindByAgentAndType(record.AgentID, record.ChannelID, record.WalletType)
		if err != nil {
			log.Printf("[ActivationRewardSettleJob] Find wallet failed for agent %d: %v", record.AgentID, err)
			failCount++
			continue
		}

		walletUpdates[wallet.ID] += record.ActualReward
		successCount++
	}

	// 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := j.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			log.Printf("[ActivationRewardSettleJob] Batch update balance failed: %v", err)
		} else {
			// 更新记录状态为已入账
			for _, record := range records {
				j.rewardRecordRepo.UpdateWalletStatus(record.ID, 1)
			}
		}
	}

	log.Printf("[ActivationRewardSettleJob] Completed: success=%d, fail=%d, took=%v",
		successCount, failCount, time.Since(startTime))
}

// SimCashbackSettleJob 流量费返现入账定时任务
type SimCashbackSettleJob struct {
	simRecordRepo repository.SimCashbackRecordRepository
	walletRepo    repository.WalletRepository
	batchSize     int
	running       bool
	mu            sync.Mutex
}

// NewSimCashbackSettleJob 创建流量费返现入账任务
func NewSimCashbackSettleJob(
	simRecordRepo repository.SimCashbackRecordRepository,
	walletRepo repository.WalletRepository,
) *SimCashbackSettleJob {
	return &SimCashbackSettleJob{
		simRecordRepo: simRecordRepo,
		walletRepo:    walletRepo,
		batchSize:     200,
	}
}

// Run 执行任务（每10分钟执行一次）
func (j *SimCashbackSettleJob) Run() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		log.Printf("[SimCashbackSettleJob] Already running, skip")
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
	log.Printf("[SimCashbackSettleJob] Started")

	// 查询待处理的流量费返现记录
	records, err := j.simRecordRepo.FindPending(j.batchSize)
	if err != nil {
		log.Printf("[SimCashbackSettleJob] Find pending records failed: %v", err)
		return
	}

	if len(records) == 0 {
		log.Printf("[SimCashbackSettleJob] No pending records")
		return
	}

	successCount := 0
	failCount := 0

	// 批量处理
	walletUpdates := make(map[int64]int64)
	for _, record := range records {
		// 获取服务费钱包
		wallet, err := j.walletRepo.FindByAgentAndType(record.AgentID, record.ChannelID, record.WalletType)
		if err != nil {
			log.Printf("[SimCashbackSettleJob] Find wallet failed for agent %d: %v", record.AgentID, err)
			failCount++
			continue
		}

		walletUpdates[wallet.ID] += record.ActualCashback
		successCount++
	}

	// 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := j.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			log.Printf("[SimCashbackSettleJob] Batch update balance failed: %v", err)
		} else {
			// 更新记录状态为已入账
			for _, record := range records {
				j.simRecordRepo.UpdateWalletStatus(record.ID, 1)
			}
		}
	}

	log.Printf("[SimCashbackSettleJob] Completed: success=%d, fail=%d, took=%v",
		successCount, failCount, time.Since(startTime))
}
