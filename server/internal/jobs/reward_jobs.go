package jobs

import (
	"log"
	"time"

	"xiangshoufu/internal/service"
)

// RewardJobRunner 奖励定时任务执行器
type RewardJobRunner struct {
	rewardService *service.RewardService
	batchSize     int
}

// NewRewardJobRunner 创建奖励任务执行器
func NewRewardJobRunner(rewardService *service.RewardService) *RewardJobRunner {
	return &RewardJobRunner{
		rewardService: rewardService,
		batchSize:     500, // 每批处理500个
	}
}

// ProcessStageRewards 处理阶段奖励（定时任务入口）
// 建议执行时间：凌晨2点
func (r *RewardJobRunner) ProcessStageRewards() {
	startTime := time.Now()
	log.Printf("[RewardJob] 开始处理阶段奖励...")

	successCount, failCount, err := r.rewardService.ProcessPendingStageRewards(r.batchSize)
	if err != nil {
		log.Printf("[RewardJob] 处理阶段奖励出错: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("[RewardJob] 阶段奖励处理完成: 成功=%d, 失败=%d, 耗时=%v", successCount, failCount, duration)
}

// SettleRewardDistributions 结算奖励分配（入账到钱包）
// 建议执行时间：每10分钟
func (r *RewardJobRunner) SettleRewardDistributions() {
	startTime := time.Now()

	successCount, failCount, err := r.rewardService.SettleRewardDistributions(r.batchSize)
	if err != nil {
		log.Printf("[RewardJob] 结算奖励分配出错: %v", err)
	}

	if successCount > 0 || failCount > 0 {
		duration := time.Since(startTime)
		log.Printf("[RewardJob] 奖励分配结算完成: 成功=%d, 失败=%d, 耗时=%v", successCount, failCount, duration)
	}
}

// RegisterRewardJobs 注册奖励相关定时任务
func RegisterRewardJobs(scheduler *Scheduler, runner *RewardJobRunner) {
	// 每天凌晨2点处理阶段奖励
	scheduler.AddJob("reward_stage_processor", 24*time.Hour, runner.ProcessStageRewards)

	// 每10分钟结算奖励分配
	scheduler.AddJob("reward_distribution_settle", 10*time.Minute, runner.SettleRewardDistributions)

	log.Printf("[RewardJob] 奖励定时任务已注册")
}
