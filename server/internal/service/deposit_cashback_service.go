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

// DepositCashbackService 押金返现服务
// 业务规则Q32：押金返现 - 商户押金收取后返现给代理商，入服务费钱包
type DepositCashbackService struct {
	terminalRepo               repository.TerminalRepository
	merchantRepo               repository.MerchantRepository
	depositPolicyRepo          *repository.GormDepositCashbackPolicyRepository
	agentDepositPolicyRepo     *repository.GormAgentDepositCashbackPolicyRepository
	depositRecordRepo          *repository.GormDepositCashbackRecordRepository
	walletRepo                 repository.WalletRepository
	walletLogRepo              repository.WalletLogRepository
	agentRepo                  repository.AgentRepository
	agentPolicyRepo            repository.AgentPolicyRepository
	messageService             *MessageService
	queue                      async.MessageQueue
}

// NewDepositCashbackService 创建押金返现服务
func NewDepositCashbackService(
	terminalRepo repository.TerminalRepository,
	merchantRepo repository.MerchantRepository,
	depositPolicyRepo *repository.GormDepositCashbackPolicyRepository,
	agentDepositPolicyRepo *repository.GormAgentDepositCashbackPolicyRepository,
	depositRecordRepo *repository.GormDepositCashbackRecordRepository,
	walletRepo repository.WalletRepository,
	walletLogRepo repository.WalletLogRepository,
	agentRepo repository.AgentRepository,
	agentPolicyRepo repository.AgentPolicyRepository,
	messageService *MessageService,
	queue async.MessageQueue,
) *DepositCashbackService {
	return &DepositCashbackService{
		terminalRepo:           terminalRepo,
		merchantRepo:           merchantRepo,
		depositPolicyRepo:      depositPolicyRepo,
		agentDepositPolicyRepo: agentDepositPolicyRepo,
		depositRecordRepo:      depositRecordRepo,
		walletRepo:             walletRepo,
		walletLogRepo:          walletLogRepo,
		agentRepo:              agentRepo,
		agentPolicyRepo:        agentPolicyRepo,
		messageService:         messageService,
		queue:                  queue,
	}
}

// DepositCashbackRequest 押金返现请求
type DepositCashbackRequest struct {
	TerminalID    int64 `json:"terminal_id"`    // 终端ID
	TerminalSN    string `json:"terminal_sn"`   // 终端SN
	MerchantID    int64 `json:"merchant_id"`    // 商户ID
	ChannelID     int64 `json:"channel_id"`     // 通道ID
	DepositAmount int64 `json:"deposit_amount"` // 押金金额（分）
}

// ProcessDepositCashback 处理押金返现
// 业务规则：
// - 押金返现按级差计算：每级代理商获得的返现 = 自身配置的返现 - 下级配置的返现
// - 入服务费钱包（wallet_type = 2）
func (s *DepositCashbackService) ProcessDepositCashback(req *DepositCashbackRequest) error {
	log.Printf("[DepositCashbackService] Processing deposit cashback: terminal=%s, deposit=%d",
		req.TerminalSN, req.DepositAmount)

	// 1. 获取终端信息
	terminal, err := s.terminalRepo.FindBySN(req.TerminalSN)
	if err != nil || terminal == nil {
		return fmt.Errorf("终端不存在: %s", req.TerminalSN)
	}

	// 2. 验证押金金额
	if req.DepositAmount <= 0 {
		log.Printf("[DepositCashbackService] No deposit amount for terminal: %s", req.TerminalSN)
		return nil
	}

	// 3. 获取直属代理商
	if terminal.OwnerAgentID == 0 {
		log.Printf("[DepositCashbackService] Terminal has no owner: %s", req.TerminalSN)
		return nil
	}

	agent, err := s.agentRepo.FindByID(terminal.OwnerAgentID)
	if err != nil || agent == nil {
		return fmt.Errorf("代理商不存在: %d", terminal.OwnerAgentID)
	}

	// 4. 获取代理商链（向上遍历）
	ancestors, err := s.agentRepo.FindAncestors(agent.ID)
	if err != nil {
		return fmt.Errorf("获取上级代理商失败: %w", err)
	}

	// 构建代理商链：直属代理商 + 所有上级
	agentChain := append([]*repository.Agent{agent}, ancestors...)

	// 5. 计算每一级的返现（按级差计算）
	cashbackRecords := make([]*models.DepositCashbackRecord, 0)
	walletUpdates := make(map[int64]int64) // walletID -> 返现金额

	for i := 0; i < len(agentChain); i++ {
		currentAgent := agentChain[i]

		// 获取当前代理商的押金返现配置
		selfCashback, err := s.getAgentDepositCashback(currentAgent.ID, req.ChannelID, req.DepositAmount)
		if err != nil {
			log.Printf("[DepositCashbackService] Get agent deposit cashback failed: %v", err)
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
			lowerCashback, _ = s.getAgentDepositCashback(lowerAgent.ID, req.ChannelID, req.DepositAmount)
		}

		// 计算级差返现
		actualCashback := selfCashback - lowerCashback
		if actualCashback <= 0 {
			continue // 没有返现空间
		}

		// 获取下级代理商ID（级差来源）
		var sourceAgentID *int64
		if i > 0 {
			srcID := agentChain[i-1].ID
			sourceAgentID = &srcID
		}

		// 创建返现记录
		record := &models.DepositCashbackRecord{
			TerminalID:     terminal.ID,
			TerminalSN:     req.TerminalSN,
			MerchantID:     req.MerchantID,
			ChannelID:      req.ChannelID,
			AgentID:        currentAgent.ID,
			DepositAmount:  req.DepositAmount,
			SelfCashback:   selfCashback,
			UpperCashback:  lowerCashback,
			ActualCashback: actualCashback,
			SourceAgentID:  sourceAgentID,
			WalletType:     models.WalletTypeService, // 服务费钱包
			WalletStatus:   0,                        // 待入账
			TriggerType:    models.DepositTriggerTypeAuto,
			CreatedAt:      time.Now(),
		}
		cashbackRecords = append(cashbackRecords, record)

		// 获取钱包并记录更新
		wallet, err := s.walletRepo.FindByAgentAndType(currentAgent.ID, req.ChannelID, models.WalletTypeService)
		if err == nil && wallet != nil {
			walletUpdates[wallet.ID] = actualCashback
		}
	}

	// 6. 批量创建返现记录
	if len(cashbackRecords) > 0 {
		if err := s.depositRecordRepo.BatchCreate(cashbackRecords); err != nil {
			return fmt.Errorf("批量创建押金返现记录失败: %w", err)
		}
	}

	// 7. 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := s.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			return fmt.Errorf("批量更新钱包余额失败: %w", err)
		}

		// 更新返现记录状态为已入账
		for _, record := range cashbackRecords {
			s.depositRecordRepo.UpdateWalletStatus(record.ID, 1)
		}
	}

	// 8. 发送返现通知
	s.sendDepositCashbackNotifications(cashbackRecords)

	log.Printf("[DepositCashbackService] Processed deposit cashback: terminal=%s, deposit=%d, records=%d",
		req.TerminalSN, req.DepositAmount, len(cashbackRecords))

	return nil
}

// getAgentDepositCashback 获取代理商的押金返现配置
func (s *DepositCashbackService) getAgentDepositCashback(agentID, channelID, depositAmount int64) (int64, error) {
	// 优先从代理商个性化配置获取
	agentPolicy, err := s.agentDepositPolicyRepo.FindByAgentChannelAndAmount(agentID, channelID, depositAmount)
	if err == nil && agentPolicy != nil {
		return agentPolicy.CashbackAmount, nil
	}

	// 从政策模板获取
	policy, err := s.agentPolicyRepo.FindByAgentAndChannel(agentID, channelID)
	if err != nil || policy == nil {
		return 0, fmt.Errorf("代理商政策不存在: agent=%d, channel=%d", agentID, channelID)
	}

	templatePolicy, err := s.depositPolicyRepo.FindByTemplateAndAmount(policy.TemplateID, depositAmount)
	if err != nil || templatePolicy == nil {
		return 0, fmt.Errorf("押金返现政策不存在: template=%d, deposit=%d", policy.TemplateID, depositAmount)
	}

	return templatePolicy.CashbackAmount, nil
}

// sendDepositCashbackNotifications 发送押金返现通知
func (s *DepositCashbackService) sendDepositCashbackNotifications(records []*models.DepositCashbackRecord) {
	if s.messageService == nil || s.queue == nil {
		return
	}

	for _, record := range records {
		msg := &NotificationMessage{
			AgentID:     record.AgentID,
			MessageType: 5, // 押金返现
			Title:       "押金返现到账",
			Content: fmt.Sprintf("终端%s押金返现 ¥%.2f 已到账",
				record.TerminalSN, float64(record.ActualCashback)/100),
			RelatedID:   record.ID,
			RelatedType: "deposit_cashback_record",
		}

		msgBytes, _ := json.Marshal(msg)
		if err := s.queue.Publish(async.TopicNotification, msgBytes); err != nil {
			log.Printf("[DepositCashbackService] Send notification failed: %v", err)
		}
	}
}

// GetDepositCashbackRecords 获取押金返现记录
func (s *DepositCashbackService) GetDepositCashbackRecords(agentID int64, limit, offset int) ([]*models.DepositCashbackRecord, int64, error) {
	records, err := s.depositRecordRepo.FindByTerminalID(agentID)
	if err != nil {
		return nil, 0, err
	}
	// 简化处理，实际应该分页查询
	total := int64(len(records))
	if offset >= len(records) {
		return []*models.DepositCashbackRecord{}, total, nil
	}
	end := offset + limit
	if end > len(records) {
		end = len(records)
	}
	return records[offset:end], total, nil
}

// DepositCashbackStats 押金返现统计
type DepositCashbackStats struct {
	TotalCashback   int64 `json:"total_cashback"`   // 总返现金额
	TotalCount      int64 `json:"total_count"`      // 总返现次数
	Deposit99Count  int64 `json:"deposit_99_count"` // 99元押金次数
	Deposit199Count int64 `json:"deposit_199_count"`// 199元押金次数
	Deposit299Count int64 `json:"deposit_299_count"`// 299元押金次数
}

// GetDepositCashbackStatsByAgent 获取代理商押金返现统计
func (s *DepositCashbackService) GetDepositCashbackStatsByAgent(agentID int64) (*DepositCashbackStats, error) {
	records, err := s.depositRecordRepo.FindByTerminalID(agentID)
	if err != nil {
		return nil, err
	}

	stats := &DepositCashbackStats{}
	for _, record := range records {
		stats.TotalCashback += record.ActualCashback
		stats.TotalCount++

		switch record.DepositAmount {
		case 9900:
			stats.Deposit99Count++
		case 19900:
			stats.Deposit199Count++
		case 29900:
			stats.Deposit299Count++
		}
	}

	return stats, nil
}
