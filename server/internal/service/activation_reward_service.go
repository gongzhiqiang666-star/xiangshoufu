package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/async"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// ActivationRewardService 激活奖励服务
// 业务规则Q33：激活奖励 - 按入网时间+交易量条件触发奖励，入奖励钱包
type ActivationRewardService struct {
	terminalRepo              repository.TerminalRepository
	merchantRepo              repository.MerchantRepository
	transactionRepo           repository.TransactionRepository
	rewardPolicyRepo          *repository.GormActivationRewardPolicyRepository
	agentRewardPolicyRepo     *repository.GormAgentActivationRewardPolicyRepository
	rewardRecordRepo          *repository.GormActivationRewardRecordRepository
	walletRepo                repository.WalletRepository
	walletLogRepo             repository.WalletLogRepository
	agentRepo                 repository.AgentRepository
	agentPolicyRepo           repository.AgentPolicyRepository
	messageService            *MessageService
	queue                     async.MessageQueue
}

// NewActivationRewardService 创建激活奖励服务
func NewActivationRewardService(
	terminalRepo repository.TerminalRepository,
	merchantRepo repository.MerchantRepository,
	transactionRepo repository.TransactionRepository,
	rewardPolicyRepo *repository.GormActivationRewardPolicyRepository,
	agentRewardPolicyRepo *repository.GormAgentActivationRewardPolicyRepository,
	rewardRecordRepo *repository.GormActivationRewardRecordRepository,
	walletRepo repository.WalletRepository,
	walletLogRepo repository.WalletLogRepository,
	agentRepo repository.AgentRepository,
	agentPolicyRepo repository.AgentPolicyRepository,
	messageService *MessageService,
	queue async.MessageQueue,
) *ActivationRewardService {
	return &ActivationRewardService{
		terminalRepo:          terminalRepo,
		merchantRepo:          merchantRepo,
		transactionRepo:       transactionRepo,
		rewardPolicyRepo:      rewardPolicyRepo,
		agentRewardPolicyRepo: agentRewardPolicyRepo,
		rewardRecordRepo:      rewardRecordRepo,
		walletRepo:            walletRepo,
		walletLogRepo:         walletLogRepo,
		agentRepo:             agentRepo,
		agentPolicyRepo:       agentPolicyRepo,
		messageService:        messageService,
		queue:                 queue,
	}
}

// TerminalRewardCheckRequest 终端奖励检查请求
type TerminalRewardCheckRequest struct {
	TerminalID  int64     `json:"terminal_id"`
	TerminalSN  string    `json:"terminal_sn"`
	MerchantID  int64     `json:"merchant_id"`
	ChannelID   int64     `json:"channel_id"`
	TradeAmount int64     `json:"trade_amount"`  // 累计交易量（分）
	RegisterAt  time.Time `json:"register_at"`   // 入网时间
	CheckDate   time.Time `json:"check_date"`    // 检查日期
}

// CheckAndProcessReward 检查并处理激活奖励
// 业务规则：
// - 根据入网天数和交易量判断是否满足奖励条件
// - 按级差计算：每级代理商获得的奖励 = 自身配置的奖励 - 下级配置的奖励
// - 入奖励钱包（wallet_type = 3）
func (s *ActivationRewardService) CheckAndProcessReward(req *TerminalRewardCheckRequest) error {
	log.Printf("[ActivationRewardService] Checking reward: terminal=%s, trade_amount=%d",
		req.TerminalSN, req.TradeAmount)

	// 1. 计算入网天数
	registerDays := int(req.CheckDate.Sub(req.RegisterAt).Hours() / 24)
	if registerDays < 0 {
		registerDays = 0
	}

	// 2. 获取终端信息
	terminal, err := s.terminalRepo.FindBySN(req.TerminalSN)
	if err != nil || terminal == nil {
		return fmt.Errorf("终端不存在: %s", req.TerminalSN)
	}

	// 3. 获取直属代理商
	if terminal.OwnerAgentID == 0 {
		log.Printf("[ActivationRewardService] Terminal has no owner: %s", req.TerminalSN)
		return nil
	}

	agent, err := s.agentRepo.FindByID(terminal.OwnerAgentID)
	if err != nil || agent == nil {
		return fmt.Errorf("代理商不存在: %d", terminal.OwnerAgentID)
	}

	// 4. 获取代理商政策
	agentPolicy, err := s.agentPolicyRepo.FindByAgentAndChannel(agent.ID, req.ChannelID)
	if err != nil || agentPolicy == nil {
		log.Printf("[ActivationRewardService] Agent policy not found: agent=%d, channel=%d", agent.ID, req.ChannelID)
		return nil
	}

	// 5. 获取激活奖励政策
	policies, err := s.rewardPolicyRepo.FindByTemplateID(agentPolicy.TemplateID)
	if err != nil || len(policies) == 0 {
		log.Printf("[ActivationRewardService] No reward policies found for template: %d", agentPolicy.TemplateID)
		return nil
	}

	// 6. 查找匹配的奖励政策
	var matchedPolicy *models.ActivationRewardPolicy
	for _, policy := range policies {
		if registerDays >= policy.MinRegisterDays &&
			registerDays <= policy.MaxRegisterDays &&
			req.TradeAmount >= policy.TargetAmount {
			matchedPolicy = policy
			break // 取优先级最高的匹配政策
		}
	}

	if matchedPolicy == nil {
		log.Printf("[ActivationRewardService] No matching policy: terminal=%s, days=%d, amount=%d",
			req.TerminalSN, registerDays, req.TradeAmount)
		return nil
	}

	// 7. 检查是否已发放过该奖励
	existingRecord, _ := s.rewardRecordRepo.FindByPolicyAndTerminal(matchedPolicy.ID, terminal.ID, req.CheckDate)
	if existingRecord != nil {
		log.Printf("[ActivationRewardService] Reward already processed: policy=%d, terminal=%d",
			matchedPolicy.ID, terminal.ID)
		return nil
	}

	// 8. 获取代理商链（向上遍历）
	ancestors, err := s.agentRepo.FindAncestors(agent.ID)
	if err != nil {
		return fmt.Errorf("获取上级代理商失败: %w", err)
	}

	// 构建代理商链：直属代理商 + 所有上级
	agentChain := append([]*repository.Agent{agent}, ancestors...)

	// 9. 计算每一级的奖励（按级差计算）
	rewardRecords := make([]*models.ActivationRewardRecord, 0)
	walletUpdates := make(map[int64]int64) // walletID -> 奖励金额

	for i := 0; i < len(agentChain); i++ {
		currentAgent := agentChain[i]

		// 获取当前代理商的奖励配置
		selfReward, err := s.getAgentReward(currentAgent.ID, req.ChannelID, matchedPolicy)
		if err != nil {
			log.Printf("[ActivationRewardService] Get agent reward failed: %v", err)
			continue
		}

		// 获取下级的奖励配置
		var lowerReward int64
		if i == 0 {
			// 直属代理商：下级奖励为0（商户没有奖励）
			lowerReward = 0
		} else {
			// 非直属：下级奖励 = 下级代理商的奖励配置
			lowerAgent := agentChain[i-1]
			lowerReward, _ = s.getAgentReward(lowerAgent.ID, req.ChannelID, matchedPolicy)
		}

		// 计算级差奖励
		actualReward := selfReward - lowerReward
		if actualReward <= 0 {
			continue // 没有奖励空间
		}

		// 获取下级代理商ID（级差来源）
		var sourceAgentID *int64
		if i > 0 {
			srcID := agentChain[i-1].ID
			sourceAgentID = &srcID
		}

		// 创建奖励记录
		record := &models.ActivationRewardRecord{
			PolicyID:      matchedPolicy.ID,
			TerminalID:    terminal.ID,
			TerminalSN:    req.TerminalSN,
			MerchantID:    req.MerchantID,
			ChannelID:     req.ChannelID,
			AgentID:       currentAgent.ID,
			RegisterDays:  registerDays,
			TradeAmount:   req.TradeAmount,
			TargetAmount:  matchedPolicy.TargetAmount,
			SelfReward:    selfReward,
			UpperReward:   lowerReward,
			ActualReward:  actualReward,
			SourceAgentID: sourceAgentID,
			WalletType:    models.WalletTypeReward, // 奖励钱包
			WalletStatus:  0,                       // 待入账
			CheckDate:     req.CheckDate,
			CreatedAt:     time.Now(),
		}
		rewardRecords = append(rewardRecords, record)

		// 获取钱包并记录更新
		wallet, err := s.walletRepo.FindByAgentAndType(currentAgent.ID, req.ChannelID, models.WalletTypeReward)
		if err == nil && wallet != nil {
			walletUpdates[wallet.ID] = actualReward
		}
	}

	// 10. 批量创建奖励记录
	if len(rewardRecords) > 0 {
		if err := s.rewardRecordRepo.BatchCreate(rewardRecords); err != nil {
			return fmt.Errorf("批量创建激活奖励记录失败: %w", err)
		}
	}

	// 11. 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := s.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			return fmt.Errorf("批量更新钱包余额失败: %w", err)
		}

		// 更新奖励记录状态为已入账
		for _, record := range rewardRecords {
			s.rewardRecordRepo.UpdateWalletStatus(record.ID, 1)
		}
	}

	// 12. 发送奖励通知
	s.sendRewardNotifications(rewardRecords, matchedPolicy.RewardName)

	log.Printf("[ActivationRewardService] Processed reward: terminal=%s, policy=%s, records=%d",
		req.TerminalSN, matchedPolicy.RewardName, len(rewardRecords))

	return nil
}

// getAgentReward 获取代理商的激活奖励配置
func (s *ActivationRewardService) getAgentReward(agentID, channelID int64, policy *models.ActivationRewardPolicy) (int64, error) {
	// 优先从代理商个性化配置获取
	agentPolicies, err := s.agentRewardPolicyRepo.FindByAgentAndChannel(agentID, channelID)
	if err == nil && len(agentPolicies) > 0 {
		// 查找匹配的奖励政策
		for _, ap := range agentPolicies {
			if ap.MinRegisterDays == policy.MinRegisterDays &&
				ap.MaxRegisterDays == policy.MaxRegisterDays &&
				ap.TargetAmount == policy.TargetAmount {
				return ap.RewardAmount, nil
			}
		}
	}

	// 使用模板配置
	return policy.RewardAmount, nil
}

// sendRewardNotifications 发送奖励通知
func (s *ActivationRewardService) sendRewardNotifications(records []*models.ActivationRewardRecord, rewardName string) {
	if s.messageService == nil || s.queue == nil {
		return
	}

	for _, record := range records {
		msg := &NotificationMessage{
			AgentID:     record.AgentID,
			MessageType: 6, // 激活奖励
			Title:       "激活奖励到账",
			Content: fmt.Sprintf("终端%s达成「%s」，获得奖励 ¥%.2f",
				record.TerminalSN, rewardName, float64(record.ActualReward)/100),
			RelatedID:   record.ID,
			RelatedType: "activation_reward_record",
		}

		msgBytes, _ := json.Marshal(msg)
		if err := s.queue.Publish(async.TopicNotification, msgBytes); err != nil {
			log.Printf("[ActivationRewardService] Send notification failed: %v", err)
		}
	}
}

// GetRewardRecords 获取奖励记录
func (s *ActivationRewardService) GetRewardRecords(agentID int64, limit, offset int) ([]*models.ActivationRewardRecord, int64, error) {
	records, err := s.rewardRecordRepo.FindPendingRecords(10000)
	if err != nil {
		return nil, 0, err
	}

	// 过滤当前代理商的记录
	agentRecords := make([]*models.ActivationRewardRecord, 0)
	for _, r := range records {
		if r.AgentID == agentID {
			agentRecords = append(agentRecords, r)
		}
	}

	total := int64(len(agentRecords))
	if offset >= len(agentRecords) {
		return []*models.ActivationRewardRecord{}, total, nil
	}
	end := offset + limit
	if end > len(agentRecords) {
		end = len(agentRecords)
	}
	return agentRecords[offset:end], total, nil
}

// ActivationRewardStats 激活奖励统计
type ActivationRewardStats struct {
	TotalReward    int64 `json:"total_reward"`    // 总奖励金额
	TotalCount     int64 `json:"total_count"`     // 总奖励次数
	ThisMonthReward int64 `json:"this_month_reward"` // 本月奖励
	ThisMonthCount  int64 `json:"this_month_count"`  // 本月次数
}

// GetRewardStatsByAgent 获取代理商奖励统计
func (s *ActivationRewardService) GetRewardStatsByAgent(agentID int64) (*ActivationRewardStats, error) {
	records, _, err := s.GetRewardRecords(agentID, 10000, 0)
	if err != nil {
		return nil, err
	}

	stats := &ActivationRewardStats{}
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for _, record := range records {
		stats.TotalReward += record.ActualReward
		stats.TotalCount++

		if record.CreatedAt.After(thisMonth) {
			stats.ThisMonthReward += record.ActualReward
			stats.ThisMonthCount++
		}
	}

	return stats, nil
}

// BatchCheckTerminalRewards 批量检查终端激活奖励（定时任务调用）
func (s *ActivationRewardService) BatchCheckTerminalRewards(channelID int64, checkDate time.Time) error {
	log.Printf("[ActivationRewardService] Starting batch reward check: channel=%d, date=%s",
		channelID, checkDate.Format("2006-01-02"))

	// 获取该通道下所有激活奖励政策
	policies, err := s.rewardPolicyRepo.FindActiveByChannelID(channelID)
	if err != nil || len(policies) == 0 {
		log.Printf("[ActivationRewardService] No active policies for channel: %d", channelID)
		return nil
	}

	// 获取最大入网天数范围
	maxDays := 0
	for _, p := range policies {
		if p.MaxRegisterDays > maxDays {
			maxDays = p.MaxRegisterDays
		}
	}

	// 计算入网时间范围
	minRegisterTime := checkDate.AddDate(0, 0, -maxDays)

	// 获取符合条件的终端列表
	terminals, err := s.terminalRepo.FindActivatedAfter(channelID, minRegisterTime)
	if err != nil {
		return fmt.Errorf("获取终端列表失败: %w", err)
	}

	log.Printf("[ActivationRewardService] Found %d terminals to check", len(terminals))

	// 逐个检查终端
	processedCount := 0
	for _, terminal := range terminals {
		// 跳过没有商户ID的终端
		if terminal.MerchantID == nil {
			continue
		}

		// 跳过没有激活时间的终端
		if terminal.ActivatedAt == nil {
			continue
		}

		// 获取终端累计交易量
		tradeAmount, err := s.getTerminalTradeAmount(terminal.TerminalSN)
		if err != nil {
			log.Printf("[ActivationRewardService] Get trade amount failed for terminal %s: %v",
				terminal.TerminalSN, err)
			continue
		}

		// 检查并处理奖励
		req := &TerminalRewardCheckRequest{
			TerminalID:  terminal.ID,
			TerminalSN:  terminal.TerminalSN,
			MerchantID:  *terminal.MerchantID,
			ChannelID:   channelID,
			TradeAmount: tradeAmount,
			RegisterAt:  *terminal.ActivatedAt,
			CheckDate:   checkDate,
		}

		if err := s.CheckAndProcessReward(req); err != nil {
			log.Printf("[ActivationRewardService] Process reward failed for terminal %s: %v",
				terminal.TerminalSN, err)
		} else {
			processedCount++
		}
	}

	log.Printf("[ActivationRewardService] Batch check completed: processed=%d/%d",
		processedCount, len(terminals))

	return nil
}

// getTerminalTradeAmount 获取终端累计交易量
// P0修复：实现实际的交易量查询逻辑
func (s *ActivationRewardService) getTerminalTradeAmount(terminalSN string) (int64, error) {
	// 调用交易仓库获取终端累计交易量
	return s.transactionRepo.GetTerminalTotalTradeAmount(terminalSN)
}
