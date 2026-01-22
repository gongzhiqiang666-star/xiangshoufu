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

// SimCashbackService 流量费返现服务
// 业务规则Q30：流量费返现三档（首次/2次/2+N次），按级差计算
type SimCashbackService struct {
	terminalRepo    repository.TerminalRepository
	policyRepo      repository.SimCashbackPolicyRepository
	recordRepo      repository.SimCashbackRecordRepository
	deviceFeeRepo   repository.DeviceFeeRepository
	walletRepo      repository.WalletRepository
	walletLogRepo   repository.WalletLogRepository
	agentRepo       repository.AgentRepository
	agentPolicyRepo repository.AgentPolicyRepository
	messageService  *MessageService
	queue           async.MessageQueue
}

// NewSimCashbackService 创建流量费返现服务
func NewSimCashbackService(
	terminalRepo repository.TerminalRepository,
	policyRepo repository.SimCashbackPolicyRepository,
	recordRepo repository.SimCashbackRecordRepository,
	deviceFeeRepo repository.DeviceFeeRepository,
	walletRepo repository.WalletRepository,
	walletLogRepo repository.WalletLogRepository,
	agentRepo repository.AgentRepository,
	agentPolicyRepo repository.AgentPolicyRepository,
	messageService *MessageService,
	queue async.MessageQueue,
) *SimCashbackService {
	return &SimCashbackService{
		terminalRepo:    terminalRepo,
		policyRepo:      policyRepo,
		recordRepo:      recordRepo,
		deviceFeeRepo:   deviceFeeRepo,
		walletRepo:      walletRepo,
		walletLogRepo:   walletLogRepo,
		agentRepo:       agentRepo,
		agentPolicyRepo: agentPolicyRepo,
		messageService:  messageService,
		queue:           queue,
	}
}

// ProcessSimFee 处理流量费缴费并计算返现
// 业务规则：
// - 流量费返现三档：首次/2次/2+N次
// - 按级差计算：每级代理商获得的返现 = 自身配置的返现 - 下级配置的返现
func (s *SimCashbackService) ProcessSimFee(deviceFee *models.DeviceFee) error {
	// 1. 获取终端信息
	terminal, err := s.terminalRepo.FindBySN(deviceFee.TerminalSN)
	if err != nil || terminal == nil {
		return fmt.Errorf("终端不存在: %s", deviceFee.TerminalSN)
	}

	// 2. 更新终端的流量费缴费次数
	newSimFeeCount := terminal.SimFeeCount + 1
	if err := s.terminalRepo.UpdateSimFeeCount(terminal.ID, newSimFeeCount); err != nil {
		log.Printf("[SimCashbackService] Update terminal sim fee count failed: %v", err)
	}

	// 3. 确定返现档次
	cashbackTier := models.GetCashbackTier(newSimFeeCount)

	// 4. 获取直属代理商
	if terminal.OwnerAgentID == 0 {
		log.Printf("[SimCashbackService] Terminal has no owner: %s", deviceFee.TerminalSN)
		return nil
	}

	agent, err := s.agentRepo.FindByID(terminal.OwnerAgentID)
	if err != nil || agent == nil {
		return fmt.Errorf("代理商不存在: %d", terminal.OwnerAgentID)
	}

	// 5. 获取代理商链（向上遍历）
	ancestors, err := s.agentRepo.FindAncestors(agent.ID)
	if err != nil {
		return fmt.Errorf("获取上级代理商失败: %w", err)
	}

	// 构建代理商链：直属代理商 + 所有上级
	agentChain := append([]*repository.Agent{agent}, ancestors...)

	// 6. 计算每一级的返现（按级差计算）
	cashbackRecords := make([]*models.SimCashbackRecord, 0)
	walletUpdates := make(map[int64]int64) // walletID -> 返现金额

	for i := 0; i < len(agentChain); i++ {
		currentAgent := agentChain[i]

		// 获取当前代理商的返现配置
		selfCashback, err := s.getAgentCashback(currentAgent.ID, deviceFee.ChannelID, terminal.BrandCode, cashbackTier)
		if err != nil {
			log.Printf("[SimCashbackService] Get agent cashback failed: %v", err)
			continue
		}

		// 获取下级的返现配置
		var lowerCashback int64
		if i == 0 {
			// 直属代理商：下级返现为0（商户没有返现）
			lowerCashback = 0
		} else {
			// 非直属：下级返现 = 下级代理商的返现配置
			lowerAgent := agentChain[i-1]
			lowerCashback, _ = s.getAgentCashback(lowerAgent.ID, deviceFee.ChannelID, terminal.BrandCode, cashbackTier)
		}

		// 计算级差返现
		actualCashback := selfCashback - lowerCashback
		if actualCashback <= 0 {
			continue // 没有返现空间
		}

		// 创建返现记录
		record := &models.SimCashbackRecord{
			DeviceFeeID:    deviceFee.ID,
			TerminalSN:     deviceFee.TerminalSN,
			ChannelID:      deviceFee.ChannelID,
			AgentID:        currentAgent.ID,
			SimFeeCount:    newSimFeeCount,
			SimFeeAmount:   deviceFee.FeeAmount,
			CashbackTier:   cashbackTier,
			SelfCashback:   selfCashback,
			UpperCashback:  lowerCashback,
			ActualCashback: actualCashback,
			SourceAgentID:  terminal.OwnerAgentID,
			WalletType:     models.WalletTypeService, // 服务费钱包（P0修复：原为4充值钱包，应为2服务费钱包）
			WalletStatus:   0,                        // 待入账
			CreatedAt:      time.Now(),
		}
		cashbackRecords = append(cashbackRecords, record)

		// 获取钱包并记录更新
		wallet, err := s.walletRepo.FindByAgentAndType(currentAgent.ID, deviceFee.ChannelID, models.WalletTypeService)
		if err == nil && wallet != nil {
			walletUpdates[wallet.ID] = actualCashback
		}
	}

	// 7. 批量创建返现记录
	if len(cashbackRecords) > 0 {
		if err := s.recordRepo.BatchCreate(cashbackRecords); err != nil {
			return fmt.Errorf("批量创建返现记录失败: %w", err)
		}
	}

	// 8. 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := s.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			return fmt.Errorf("批量更新钱包余额失败: %w", err)
		}

		// 更新返现记录状态为已入账
		for _, record := range cashbackRecords {
			s.recordRepo.UpdateWalletStatus(record.ID, 1)
		}
	}

	// 9. 更新流量费记录的返现状态
	totalCashback := int64(0)
	for _, record := range cashbackRecords {
		totalCashback += record.ActualCashback
	}
	if err := s.deviceFeeRepo.UpdateCashbackStatus(deviceFee.ID, 1, totalCashback); err != nil {
		log.Printf("[SimCashbackService] Update device fee cashback status failed: %v", err)
	}

	// 10. 发送返现通知
	s.sendCashbackNotifications(cashbackRecords)

	log.Printf("[SimCashbackService] Processed sim fee: terminal=%s, count=%d, tier=%d, records=%d",
		deviceFee.TerminalSN, newSimFeeCount, cashbackTier, len(cashbackRecords))

	return nil
}

// getAgentCashback 获取代理商的流量费返现配置
func (s *SimCashbackService) getAgentCashback(agentID, channelID int64, brandCode string, tier int16) (int64, error) {
	// 获取代理商的政策模板ID
	policy, err := s.agentPolicyRepo.FindByAgentAndChannel(agentID, channelID)
	if err != nil || policy == nil {
		return 0, fmt.Errorf("代理商政策不存在: agent=%d, channel=%d", agentID, channelID)
	}

	// 获取流量费返现政策
	cashbackPolicy, err := s.policyRepo.FindByTemplateAndChannel(policy.TemplateID, channelID, brandCode)
	if err != nil || cashbackPolicy == nil {
		return 0, fmt.Errorf("流量费返现政策不存在: template=%d", policy.TemplateID)
	}

	// 根据档次返回返现金额
	switch tier {
	case models.SimCashbackTierFirst:
		return cashbackPolicy.FirstTimeCashback, nil
	case models.SimCashbackTierSecond:
		return cashbackPolicy.SecondTimeCashback, nil
	case models.SimCashbackTierThirdPlus:
		return cashbackPolicy.ThirdPlusCashback, nil
	default:
		return 0, fmt.Errorf("未知的返现档次: %d", tier)
	}
}

// sendCashbackNotifications 发送返现通知
func (s *SimCashbackService) sendCashbackNotifications(records []*models.SimCashbackRecord) {
	if s.messageService == nil || s.queue == nil {
		return
	}

	for _, record := range records {
		tierName := getTierName(record.CashbackTier)
		msg := &NotificationMessage{
			AgentID:     record.AgentID,
			MessageType: 4, // 流量返现
			Title:       "流量费返现到账",
			Content: fmt.Sprintf("终端%s%s缴费，获得返现 ¥%.2f",
				record.TerminalSN, tierName, float64(record.ActualCashback)/100),
			RelatedID:   record.ID,
			RelatedType: "sim_cashback_record",
		}

		msgBytes, _ := json.Marshal(msg)
		if err := s.queue.Publish(async.TopicNotification, msgBytes); err != nil {
			log.Printf("[SimCashbackService] Send notification failed: %v", err)
		}
	}
}

// getTierName 获取档次名称
func getTierName(tier int16) string {
	switch tier {
	case models.SimCashbackTierFirst:
		return "首次"
	case models.SimCashbackTierSecond:
		return "第2次"
	case models.SimCashbackTierThirdPlus:
		return "第3次及以后"
	default:
		return ""
	}
}

// GetCashbackRecords 获取返现记录
func (s *SimCashbackService) GetCashbackRecords(agentID int64, limit, offset int) ([]*models.SimCashbackRecord, int64, error) {
	return s.recordRepo.FindByAgent(agentID, limit, offset)
}

// GetCashbackStats 获取返现统计
type CashbackStats struct {
	TotalCashback      int64 `json:"total_cashback"`       // 总返现金额
	FirstTierCashback  int64 `json:"first_tier_cashback"`  // 首次返现金额
	SecondTierCashback int64 `json:"second_tier_cashback"` // 第2次返现金额
	ThirdTierCashback  int64 `json:"third_tier_cashback"`  // 第3次及以后返现金额
	TotalCount         int64 `json:"total_count"`          // 总返现次数
}

// GetCashbackStatsByAgent 获取代理商返现统计（需要在repository中添加聚合查询）
func (s *SimCashbackService) GetCashbackStatsByAgent(agentID int64) (*CashbackStats, error) {
	// 这里简化处理，实际应该在repository中实现聚合查询
	records, _, err := s.recordRepo.FindByAgent(agentID, 10000, 0)
	if err != nil {
		return nil, err
	}

	stats := &CashbackStats{}
	for _, record := range records {
		stats.TotalCashback += record.ActualCashback
		stats.TotalCount++

		switch record.CashbackTier {
		case models.SimCashbackTierFirst:
			stats.FirstTierCashback += record.ActualCashback
		case models.SimCashbackTierSecond:
			stats.SecondTierCashback += record.ActualCashback
		case models.SimCashbackTierThirdPlus:
			stats.ThirdTierCashback += record.ActualCashback
		}
	}

	return stats, nil
}
