package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"xiangshoufu/internal/async"
	"xiangshoufu/internal/repository"
)

// ProfitService 分润计算服务
type ProfitService struct {
	transactionRepo   repository.TransactionRepository
	profitRepo        repository.ProfitRecordRepository
	walletRepo        repository.WalletRepository
	walletLogRepo     repository.WalletLogRepository
	agentRepo         repository.AgentRepository
	agentPolicyRepo   repository.AgentPolicyRepository
	messageService    *MessageService
	queue             async.MessageQueue
	rateStagingService *RateStagingService // 费率阶梯服务
	goodsDeductionService *GoodsDeductionService // 货款代扣服务（保留兼容）
	deductionService *DeductionService // 统一代扣服务
	settlementPriceService *SettlementPriceService // 结算价服务（用于获取高调/P+0配置）
}

// NewProfitService 创建分润服务
func NewProfitService(
	transactionRepo repository.TransactionRepository,
	profitRepo repository.ProfitRecordRepository,
	walletRepo repository.WalletRepository,
	walletLogRepo repository.WalletLogRepository,
	agentRepo repository.AgentRepository,
	agentPolicyRepo repository.AgentPolicyRepository,
	messageService *MessageService,
	queue async.MessageQueue,
) *ProfitService {
	return &ProfitService{
		transactionRepo: transactionRepo,
		profitRepo:      profitRepo,
		walletRepo:      walletRepo,
		walletLogRepo:   walletLogRepo,
		agentRepo:       agentRepo,
		agentPolicyRepo: agentPolicyRepo,
		messageService:  messageService,
		queue:           queue,
	}
}

// SetRateStagingService 设置费率阶梯服务（延迟注入，避免循环依赖）
func (s *ProfitService) SetRateStagingService(rss *RateStagingService) {
	s.rateStagingService = rss
}

// SetGoodsDeductionService 设置货款代扣服务（延迟注入，避免循环依赖）
func (s *ProfitService) SetGoodsDeductionService(gds *GoodsDeductionService) {
	s.goodsDeductionService = gds
}

// SetDeductionService 设置统一代扣服务（延迟注入，避免循环依赖）
func (s *ProfitService) SetDeductionService(ds *DeductionService) {
	s.deductionService = ds
}

// SetSettlementPriceService 设置结算价服务（延迟注入，避免循环依赖）
func (s *ProfitService) SetSettlementPriceService(sps *SettlementPriceService) {
	s.settlementPriceService = sps
}

// ProcessMessage 处理分润计算消息
func (s *ProfitService) ProcessMessage(msgBytes []byte) error {
	var msg ProfitMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return fmt.Errorf("unmarshal profit message failed: %w", err)
	}

	return s.CalculateProfit(msg.TransactionID)
}

// CalculateProfit 计算单笔交易的分润
func (s *ProfitService) CalculateProfit(txID int64) error {
	// 1. 获取交易信息
	tx, err := s.transactionRepo.FindByID(txID)
	if err != nil {
		return fmt.Errorf("find transaction failed: %w", err)
	}
	if tx == nil {
		return fmt.Errorf("transaction not found: %d", txID)
	}

	// 2. 检查是否已计算过分润
	if tx.ProfitStatus == 1 {
		log.Printf("[ProfitService] Transaction already calculated: %d", txID)
		return nil
	}

	// 3. 获取直属代理商
	agent, err := s.agentRepo.FindByID(tx.AgentID)
	if err != nil || agent == nil {
		return fmt.Errorf("find agent failed: %d", tx.AgentID)
	}

	// 4. 获取代理商链（从直属代理商向上遍历）
	ancestors, err := s.agentRepo.FindAncestors(agent.ID)
	if err != nil {
		return fmt.Errorf("find ancestors failed: %w", err)
	}

	// 5. 计算每一级的分润
	profitRecords := make([]*repository.ProfitRecord, 0)
	walletUpdates := make(map[int64]int64) // walletID -> 分润金额

	// 构建代理商链：直属代理商 + 所有上级
	agentChain := append([]*repository.Agent{agent}, ancestors...)

	for i := 0; i < len(agentChain); i++ {
		currentAgent := agentChain[i]

		// 获取当前代理商的结算价（费率）
		selfRate, err := s.getAgentRate(currentAgent.ID, tx.ChannelID, tx.CardType)
		if err != nil {
			log.Printf("[ProfitService] Get agent rate failed: %v", err)
			continue
		}

		// 获取下级费率（如果有下级）
		var lowerRate float64
		if i == 0 {
			// 直属代理商：下级费率 = 商户费率（交易费率）
			lowerRate, _ = strconv.ParseFloat(tx.Rate, 64)
		} else {
			// 非直属：下级费率 = 下级代理商的结算价
			lowerAgent := agentChain[i-1]
			lowerRate, _ = s.getAgentRate(lowerAgent.ID, tx.ChannelID, tx.CardType)
		}

		// 计算费率差
		rateDiff := lowerRate - selfRate
		if rateDiff <= 0 {
			continue // 没有分润空间
		}

		// 计算分润金额（单位：分）
		// 分润 = 交易金额 * 费率差 / 100
		profitAmount := int64(float64(tx.Amount) * rateDiff / 100)
		if profitAmount <= 0 {
			continue
		}

		// 创建分润记录
		record := &repository.ProfitRecord{
			TransactionID:    tx.ID,
			OrderNo:          tx.OrderNo,
			AgentID:          currentAgent.ID,
			ProfitType:       1, // 交易分润
			TradeAmount:      tx.Amount,
			SelfRate:         fmt.Sprintf("%.4f", selfRate),
			LowerRate:        fmt.Sprintf("%.4f", lowerRate),
			RateDiff:         fmt.Sprintf("%.4f", rateDiff),
			ProfitAmount:     profitAmount,
			SourceMerchantID: tx.MerchantID,
			SourceAgentID:    tx.AgentID,
			ChannelID:        tx.ChannelID,
			WalletType:       1, // 分润钱包
			WalletStatus:     0, // 待入账
			CreatedAt:        time.Now(),
		}

		// 计算高调分润（如果交易有高调）
		if tx.HighRate != "" && tx.HighRate != "0" {
			highRateProfit, selfHighRate, lowerHighRate := s.calculateHighRateProfit(tx, currentAgent.ID, agentChain, i)
			if highRateProfit > 0 {
				record.HighRateProfit = highRateProfit
				record.HighRateSelf = selfHighRate
				record.HighRateLower = lowerHighRate
				record.ProfitAmount += highRateProfit // 累加到总分润
			}
		}

		// 计算P+0分润（如果交易有D0费用）
		if tx.D0Fee > 0 {
			d0ExtraProfit, selfD0Extra, lowerD0Extra := s.calculateD0ExtraProfit(tx, currentAgent.ID, agentChain, i)
			if d0ExtraProfit > 0 {
				record.D0ExtraProfit = d0ExtraProfit
				record.D0ExtraSelf = selfD0Extra
				record.D0ExtraLower = lowerD0Extra
				record.ProfitAmount += d0ExtraProfit // 累加到总分润
			}
		}

		profitRecords = append(profitRecords, record)

		// 获取钱包ID
		wallet, err := s.walletRepo.FindByAgentAndType(currentAgent.ID, tx.ChannelID, 1)
		if err == nil && wallet != nil {
			walletUpdates[wallet.ID] = record.ProfitAmount
		}
	}

	// 6. 批量写入分润记录
	if len(profitRecords) > 0 {
		if err := s.profitRepo.BatchCreate(profitRecords); err != nil {
			return fmt.Errorf("batch create profit records failed: %w", err)
		}
	}

	// 7. 批量更新钱包余额
	if len(walletUpdates) > 0 {
		if err := s.walletRepo.BatchUpdateBalance(walletUpdates); err != nil {
			return fmt.Errorf("batch update wallet failed: %w", err)
		}
	}

	// 7.1 触发代扣冻结（替代原实时扣款）
	// 优先使用统一代扣服务，如果未注入则使用旧的货款代扣服务
	if s.deductionService != nil && len(profitRecords) > 0 {
		for _, record := range profitRecords {
			// 触发该代理商的代扣冻结
			frozen, err := s.deductionService.FreezeOnIncome(
				record.AgentID,
				record.ChannelID,
				int16(record.WalletType), // 分润钱包
				record.ProfitAmount,
			)
			if err != nil {
				log.Printf("[ProfitService] Trigger deduction freeze failed for agent %d: %v", record.AgentID, err)
			} else if frozen > 0 {
				log.Printf("[ProfitService] Deduction freeze triggered: agent=%d, frozen=%d", record.AgentID, frozen)
			}
		}
	} else if s.goodsDeductionService != nil && len(profitRecords) > 0 {
		// 兼容旧的货款代扣服务
		for _, record := range profitRecords {
			deducted, err := s.goodsDeductionService.TriggerRealtimeDeduction(
				record.AgentID,
				record.ChannelID,
				int16(record.WalletType),
				record.ProfitAmount,
				"profit_income",
				&tx.ID,
				&record.ID,
			)
			if err != nil {
				log.Printf("[ProfitService] Trigger goods deduction failed for agent %d: %v", record.AgentID, err)
			} else if deducted > 0 {
				log.Printf("[ProfitService] Goods deduction triggered: agent=%d, deducted=%d", record.AgentID, deducted)
			}
		}
	}

	// 8. 更新交易分润状态
	if err := s.transactionRepo.UpdateProfitStatus(tx.ID, 1); err != nil {
		return fmt.Errorf("update profit status failed: %w", err)
	}

	// 9. 发送消息通知
	s.sendProfitNotifications(profitRecords)

	log.Printf("[ProfitService] Calculated profit for transaction %d, records: %d", txID, len(profitRecords))
	return nil
}

// getAgentRate 获取代理商的结算费率
func (s *ProfitService) getAgentRate(agentID, channelID int64, cardType int16) (float64, error) {
	policy, err := s.agentPolicyRepo.FindByAgentAndChannel(agentID, channelID)
	if err != nil || policy == nil {
		return 0, fmt.Errorf("policy not found for agent %d, channel %d", agentID, channelID)
	}

	// 根据卡类型返回对应费率
	var baseRate float64
	switch cardType {
	case 1: // 借记卡
		baseRate, _ = strconv.ParseFloat(policy.DebitRate, 64)
	case 2: // 贷记卡
		baseRate, _ = strconv.ParseFloat(policy.CreditRate, 64)
	default:
		baseRate, _ = strconv.ParseFloat(policy.CreditRate, 64)
	}

	// 应用费率阶梯调整（如果配置了RateStagingService）
	if s.rateStagingService != nil {
		adjustedRate, err := s.rateStagingService.GetAgentRateAdjustment(agentID, channelID, baseRate, cardType)
		if err == nil && adjustedRate != nil {
			return adjustedRate.AdjustedRate, nil
		}
	}

	return baseRate, nil
}

// getMerchantAdjustedRate 获取商户调整后的费率（应用商户入网时间的费率阶梯）
func (s *ProfitService) getMerchantAdjustedRate(merchantID, channelID int64, baseRate float64, cardType int16) float64 {
	if s.rateStagingService == nil {
		return baseRate
	}

	adjustment, err := s.rateStagingService.GetMerchantRateAdjustment(merchantID, channelID, baseRate, cardType)
	if err != nil {
		return baseRate
	}

	return adjustment.AdjustedRate
}

// calculateHighRateProfit 计算高调分润
// 高调分润 = 交易金额 × (下级高调费率 - 自身高调费率) / 100
func (s *ProfitService) calculateHighRateProfit(tx *repository.Transaction, agentID int64, agentChain []*repository.Agent, idx int) (profit int64, selfRate string, lowerRate string) {
	if s.settlementPriceService == nil {
		return 0, "0", "0"
	}

	// 确定费率类型
	rateType := s.getRateTypeFromCardType(tx.CardType)

	// 获取自身高调费率
	selfHighRate, err := s.settlementPriceService.GetAgentHighRate(agentID, tx.ChannelID, "", rateType)
	if err != nil {
		selfHighRate = "0"
	}

	// 获取下级高调费率
	var lowerHighRate string
	if idx == 0 {
		// 直属代理商：下级高调费率 = 交易的实际高调费率
		lowerHighRate = tx.HighRate
	} else {
		// 非直属：下级高调费率 = 下级代理商的配置
		lowerAgent := agentChain[idx-1]
		lowerHighRate, _ = s.settlementPriceService.GetAgentHighRate(lowerAgent.ID, tx.ChannelID, "", rateType)
	}

	// 计算费率差
	selfRateFloat, _ := strconv.ParseFloat(selfHighRate, 64)
	lowerRateFloat, _ := strconv.ParseFloat(lowerHighRate, 64)

	rateDiff := lowerRateFloat - selfRateFloat
	if rateDiff <= 0 {
		return 0, selfHighRate, lowerHighRate
	}

	// 计算高调分润金额
	profit = int64(float64(tx.Amount) * rateDiff / 100)
	return profit, selfHighRate, lowerHighRate
}

// calculateD0ExtraProfit 计算P+0分润（差额分配模式）
// 直属代理商：获得上级给自己配置的全部金额
// 中间/上级代理商：获得（上级给自己的 - 自己给下级的）
func (s *ProfitService) calculateD0ExtraProfit(tx *repository.Transaction, agentID int64, agentChain []*repository.Agent, idx int) (profit int64, selfExtra int64, lowerExtra int64) {
	if s.settlementPriceService == nil {
		return 0, 0, 0
	}

	// 确定费率类型
	rateType := s.getRateTypeFromCardType(tx.CardType)

	// 获取自身P+0加价配置（上级给当前代理商配置的金额）
	selfD0Extra, err := s.settlementPriceService.GetAgentD0Extra(agentID, tx.ChannelID, "", rateType)
	if err != nil {
		selfD0Extra = 0
	}

	// 获取下级P+0加价配置
	var lowerD0Extra int64
	if idx == 0 {
		// 直属代理商（最底层）：获得全部配置金额
		lowerD0Extra = 0
		profit = selfD0Extra
	} else {
		// 中间/上级代理商：获得差额
		lowerAgent := agentChain[idx-1]
		lowerD0Extra, _ = s.settlementPriceService.GetAgentD0Extra(lowerAgent.ID, tx.ChannelID, "", rateType)
		profit = selfD0Extra - lowerD0Extra
	}

	if profit < 0 {
		profit = 0
	}

	return profit, selfD0Extra, lowerD0Extra
}

// getRateTypeFromCardType 根据卡类型获取费率类型编码
func (s *ProfitService) getRateTypeFromCardType(cardType int16) string {
	switch cardType {
	case 1:
		return "DEBIT" // 借记卡
	case 2:
		return "CREDIT" // 贷记卡
	default:
		return "CREDIT"
	}
}

// sendProfitNotifications 发送分润通知
func (s *ProfitService) sendProfitNotifications(records []*repository.ProfitRecord) {
	if s.messageService == nil {
		return
	}

	for _, record := range records {
		msg := &NotificationMessage{
			AgentID:     record.AgentID,
			MessageType: 1, // 分润通知
			Title:       "交易分润到账",
			Content:     fmt.Sprintf("您获得交易分润 ¥%.2f", float64(record.ProfitAmount)/100),
			RelatedID:   record.ID,
			RelatedType: "profit_record",
		}

		msgBytes, _ := json.Marshal(msg)
		if err := s.queue.Publish(async.TopicNotification, msgBytes); err != nil {
			log.Printf("[ProfitService] Send notification failed: %v", err)
		}
	}
}

// RevokeProfit 撤销分润（退款时调用）
func (s *ProfitService) RevokeProfit(txID int64, reason string) error {
	// 1. 获取该交易的所有分润记录
	records, err := s.profitRepo.FindByTransactionID(txID)
	if err != nil {
		return fmt.Errorf("find profit records failed: %w", err)
	}

	// 2. 扣减钱包余额
	walletDeductions := make(map[int64]int64)
	for _, record := range records {
		if record.IsRevoked {
			continue
		}
		wallet, err := s.walletRepo.FindByAgentAndType(record.AgentID, record.ChannelID, int16(record.WalletType))
		if err == nil && wallet != nil {
			walletDeductions[wallet.ID] = -record.ProfitAmount // 负值表示扣减
		}
	}

	// 3. 批量扣减钱包
	if len(walletDeductions) > 0 {
		if err := s.walletRepo.BatchUpdateBalance(walletDeductions); err != nil {
			return fmt.Errorf("deduct wallet failed: %w", err)
		}
	}

	// 4. 标记分润记录为已撤销
	if err := s.profitRepo.RevokeByTransactionID(txID, reason); err != nil {
		return fmt.Errorf("revoke profit records failed: %w", err)
	}

	// 5. 更新交易退款状态
	if err := s.transactionRepo.UpdateRefundStatus(txID, 1); err != nil {
		return fmt.Errorf("update refund status failed: %w", err)
	}

	// 6. 发送撤销通知
	for _, record := range records {
		msg := &NotificationMessage{
			AgentID:     record.AgentID,
			MessageType: 5, // 退款撤销
			Title:       "分润撤销通知",
			Content:     fmt.Sprintf("交易已退款，分润 ¥%.2f 已扣回", float64(record.ProfitAmount)/100),
			RelatedID:   record.ID,
			RelatedType: "profit_record",
		}

		msgBytes, _ := json.Marshal(msg)
		s.queue.Publish(async.TopicNotification, msgBytes)
	}

	log.Printf("[ProfitService] Revoked profit for transaction %d, records: %d", txID, len(records))
	return nil
}
