package service

import (
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// GoodsDeductionService 货款代扣服务
// 业务规则：
//   1. 扣款规则: 分润优先 - 先扣分润钱包，扣完再扣服务费钱包
//   2. 扣款时机: 实时扣款 - 钱包入账时立即触发扣款
//   3. 部分扣款: 有多少扣多少 - 余额不足时部分扣除，剩余下次继续扣
//   4. 扣款上限: 无上限 - 每次入账时全部扣除，直到扣完为止
//   5. 扣款优先级: 货款代扣 > 上级代扣 > 伙伴代扣
//   6. 与提现关系: 待扣金额占用钱包余额，影响可提现金额
//   7. 接收确认: 下级拒绝接收则终端划拨失败
//   8. 修改/取消: 设置后不可修改，只能扣完或线下协商
type GoodsDeductionService struct {
	deductionRepo     repository.GoodsDeductionRepository
	detailRepo        repository.GoodsDeductionDetailRepository
	terminalRepo      repository.GoodsDeductionTerminalRepository
	notificationRepo  repository.GoodsDeductionNotificationRepository
	walletRepo        repository.WalletRepository
	walletLogRepo     repository.WalletLogRepository
	agentRepo         repository.AgentRepository
}

// NewGoodsDeductionService 创建货款代扣服务
func NewGoodsDeductionService(
	deductionRepo repository.GoodsDeductionRepository,
	detailRepo repository.GoodsDeductionDetailRepository,
	terminalRepo repository.GoodsDeductionTerminalRepository,
	notificationRepo repository.GoodsDeductionNotificationRepository,
	walletRepo repository.WalletRepository,
	walletLogRepo repository.WalletLogRepository,
	agentRepo repository.AgentRepository,
) *GoodsDeductionService {
	return &GoodsDeductionService{
		deductionRepo:    deductionRepo,
		detailRepo:       detailRepo,
		terminalRepo:     terminalRepo,
		notificationRepo: notificationRepo,
		walletRepo:       walletRepo,
		walletLogRepo:    walletLogRepo,
		agentRepo:        agentRepo,
	}
}

// CreateGoodsDeduction 创建货款代扣（终端划拨时调用）
func (s *GoodsDeductionService) CreateGoodsDeduction(req *models.CreateGoodsDeductionRequest, fromAgentID int64, createdBy int64) (*models.GoodsDeduction, error) {
	// 验证发起方代理商
	fromAgent, err := s.agentRepo.FindByID(fromAgentID)
	if err != nil || fromAgent == nil {
		return nil, fmt.Errorf("发起方代理商不存在: %d", fromAgentID)
	}

	// 验证接收方代理商
	toAgent, err := s.agentRepo.FindByID(req.ToAgentID)
	if err != nil || toAgent == nil {
		return nil, fmt.Errorf("接收方代理商不存在: %d", req.ToAgentID)
	}

	// 验证接收方是发起方的下级
	if toAgent.ParentID != fromAgentID {
		return nil, fmt.Errorf("接收方必须是发起方的直接下级")
	}

	// 验证终端列表
	if len(req.Terminals) == 0 {
		return nil, fmt.Errorf("终端列表不能为空")
	}

	// 验证扣款来源
	if req.DeductionSource < 1 || req.DeductionSource > 3 {
		return nil, fmt.Errorf("无效的扣款来源")
	}

	// 计算总金额
	terminalCount := len(req.Terminals)
	totalAmount := req.UnitPrice * int64(terminalCount)
	if totalAmount <= 0 {
		return nil, fmt.Errorf("代扣总金额必须大于0")
	}

	// 生成代扣编号
	deductionNo := fmt.Sprintf("GD%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	// 创建货款代扣
	deduction := &models.GoodsDeduction{
		DeductionNo:     deductionNo,
		FromAgentID:     fromAgentID,
		ToAgentID:       req.ToAgentID,
		TotalAmount:     totalAmount,
		DeductedAmount:  0,
		RemainingAmount: totalAmount,
		DeductionSource: req.DeductionSource,
		TerminalCount:   terminalCount,
		UnitPrice:       req.UnitPrice,
		Status:          models.GoodsDeductionStatusPendingAccept,
		AgreementSigned: req.AgreementURL != "",
		AgreementURL:    req.AgreementURL,
		DistributeID:    req.DistributeID,
		Remark:          req.Remark,
		CreatedBy:       createdBy,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.deductionRepo.Create(deduction); err != nil {
		return nil, fmt.Errorf("创建货款代扣失败: %w", err)
	}

	// 创建终端关联
	terminals := make([]*models.GoodsDeductionTerminal, 0, len(req.Terminals))
	for _, t := range req.Terminals {
		unitPrice := req.UnitPrice
		if t.UnitPrice > 0 {
			unitPrice = t.UnitPrice
		}
		terminal := &models.GoodsDeductionTerminal{
			DeductionID: deduction.ID,
			TerminalID:  t.TerminalID,
			TerminalSN:  t.TerminalSN,
			UnitPrice:   unitPrice,
			CreatedAt:   time.Now(),
		}
		terminals = append(terminals, terminal)
	}

	if err := s.terminalRepo.BatchCreate(terminals); err != nil {
		return nil, fmt.Errorf("创建终端关联失败: %w", err)
	}

	// 发送待接收通知给接收方
	notification := &models.GoodsDeductionNotification{
		DeductionID: deduction.ID,
		AgentID:     req.ToAgentID,
		NotifyType:  models.GoodsDeductionNotifyTypePending,
		Title:       "货款代扣待接收",
		Content:     fmt.Sprintf("您有一笔货款代扣待接收，金额：%.2f元，终端数量：%d台", float64(totalAmount)/100, terminalCount),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}
	s.notificationRepo.Create(notification)

	log.Printf("[GoodsDeductionService] Created goods deduction: %s, from: %d, to: %d, amount: %d, terminals: %d",
		deductionNo, fromAgentID, req.ToAgentID, totalAmount, terminalCount)

	return deduction, nil
}

// AcceptGoodsDeduction 接收货款代扣
func (s *GoodsDeductionService) AcceptGoodsDeduction(deductionID int64, agentID int64) error {
	deduction, err := s.deductionRepo.FindByID(deductionID)
	if err != nil || deduction == nil {
		return fmt.Errorf("货款代扣不存在: %d", deductionID)
	}

	// 验证是否为接收方
	if deduction.ToAgentID != agentID {
		return fmt.Errorf("无权操作此货款代扣")
	}

	// 验证状态
	if deduction.Status != models.GoodsDeductionStatusPendingAccept {
		return fmt.Errorf("货款代扣状态不允许接收")
	}

	// 更新状态为进行中
	if err := s.deductionRepo.UpdateStatus(deductionID, models.GoodsDeductionStatusInProgress); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	log.Printf("[GoodsDeductionService] Accepted goods deduction: %d, agent: %d", deductionID, agentID)

	return nil
}

// RejectGoodsDeduction 拒绝货款代扣
func (s *GoodsDeductionService) RejectGoodsDeduction(deductionID int64, agentID int64) error {
	deduction, err := s.deductionRepo.FindByID(deductionID)
	if err != nil || deduction == nil {
		return fmt.Errorf("货款代扣不存在: %d", deductionID)
	}

	// 验证是否为接收方
	if deduction.ToAgentID != agentID {
		return fmt.Errorf("无权操作此货款代扣")
	}

	// 验证状态
	if deduction.Status != models.GoodsDeductionStatusPendingAccept {
		return fmt.Errorf("货款代扣状态不允许拒绝")
	}

	// 更新状态为已拒绝
	if err := s.deductionRepo.UpdateStatus(deductionID, models.GoodsDeductionStatusRejected); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	log.Printf("[GoodsDeductionService] Rejected goods deduction: %d, agent: %d", deductionID, agentID)

	return nil
}

// GetGoodsDeductionByID 获取货款代扣详情
func (s *GoodsDeductionService) GetGoodsDeductionByID(deductionID int64) (*models.GoodsDeduction, error) {
	deduction, err := s.deductionRepo.FindByID(deductionID)
	if err != nil {
		return nil, fmt.Errorf("查询货款代扣失败: %w", err)
	}

	if deduction == nil {
		return nil, fmt.Errorf("货款代扣不存在")
	}

	// 填充代理商名称
	s.fillAgentNames(deduction)

	// 获取终端列表
	terminals, err := s.terminalRepo.FindByDeductionID(deductionID)
	if err == nil {
		deduction.Terminals = terminals
	}

	return deduction, nil
}

// GetSentList 获取我发起的货款代扣列表
func (s *GoodsDeductionService) GetSentList(agentID int64, status []int16, page, pageSize int) ([]*models.GoodsDeductionListResponse, int64, error) {
	offset := (page - 1) * pageSize
	deductions, total, err := s.deductionRepo.FindByFromAgent(agentID, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询失败: %w", err)
	}

	// 转换为响应格式
	result := make([]*models.GoodsDeductionListResponse, 0, len(deductions))
	for _, d := range deductions {
		s.fillAgentNames(d)
		result = append(result, d.ToListResponse())
	}

	return result, total, nil
}

// GetReceivedList 获取我接收的货款代扣列表
func (s *GoodsDeductionService) GetReceivedList(agentID int64, status []int16, page, pageSize int) ([]*models.GoodsDeductionListResponse, int64, error) {
	offset := (page - 1) * pageSize
	deductions, total, err := s.deductionRepo.FindByToAgent(agentID, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询失败: %w", err)
	}

	// 转换为响应格式
	result := make([]*models.GoodsDeductionListResponse, 0, len(deductions))
	for _, d := range deductions {
		s.fillAgentNames(d)
		result = append(result, d.ToListResponse())
	}

	return result, total, nil
}

// GetDeductionDetails 获取扣款明细列表
func (s *GoodsDeductionService) GetDeductionDetails(deductionID int64, page, pageSize int) ([]*models.GoodsDeductionDetail, int64, error) {
	offset := (page - 1) * pageSize
	details, total, err := s.detailRepo.FindByDeductionID(deductionID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询扣款明细失败: %w", err)
	}

	// 填充钱包类型名称
	for _, d := range details {
		d.WalletTypeName = models.WalletTypeName(d.WalletType)
	}

	return details, total, nil
}

// GetSummary 获取货款代扣统计
func (s *GoodsDeductionService) GetSummary(agentID int64, isSent bool) (*models.GoodsDeductionSummary, error) {
	if isSent {
		return s.deductionRepo.GetSummaryByFromAgent(agentID)
	}
	return s.deductionRepo.GetSummaryByToAgent(agentID)
}

// TriggerRealtimeDeduction 触发实时扣款（钱包入账时调用）
// 业务规则：
//   1. 扣款优先级：货款代扣 > 上级代扣 > 伙伴代扣
//   2. 钱包优先级：分润钱包 > 服务费钱包
//   3. 每次入账时全部扣除，直到扣完为止
func (s *GoodsDeductionService) TriggerRealtimeDeduction(
	agentID int64,
	channelID int64,
	walletType int16,
	incomeAmount int64,
	triggerType string,
	triggerTransactionID *int64,
	triggerProfitID *int64,
) (int64, error) {
	// 获取该代理商所有进行中的货款代扣
	deductions, err := s.deductionRepo.FindActiveByToAgent(agentID)
	if err != nil {
		return 0, fmt.Errorf("查询货款代扣失败: %w", err)
	}

	if len(deductions) == 0 {
		return 0, nil // 无待扣款
	}

	// 获取钱包余额
	wallet, err := s.walletRepo.FindByAgentAndType(agentID, channelID, walletType)
	if err != nil || wallet == nil {
		return 0, nil // 钱包不存在
	}

	// 可扣金额 = 入账金额（刚入账的钱可以全部用于扣款）
	availableAmount := incomeAmount
	if availableAmount <= 0 {
		return 0, nil
	}

	var totalDeducted int64

	// 遍历所有进行中的货款代扣
	for _, deduction := range deductions {
		if availableAmount <= 0 {
			break
		}

		// 检查扣款来源是否匹配
		if !s.isWalletTypeAllowed(deduction.DeductionSource, walletType) {
			continue
		}

		// 计算本次可扣金额
		deductAmount := deduction.RemainingAmount
		if deductAmount > availableAmount {
			deductAmount = availableAmount
		}

		if deductAmount <= 0 {
			continue
		}

		// 执行扣款
		deducted, err := s.executeDeduction(deduction, wallet, deductAmount, triggerType, triggerTransactionID, triggerProfitID)
		if err != nil {
			log.Printf("[GoodsDeductionService] Deduction failed for %d: %v", deduction.ID, err)
			continue
		}

		totalDeducted += deducted
		availableAmount -= deducted
	}

	return totalDeducted, nil
}

// isWalletTypeAllowed 检查钱包类型是否允许扣款
func (s *GoodsDeductionService) isWalletTypeAllowed(deductionSource int16, walletType int16) bool {
	switch deductionSource {
	case models.GoodsDeductionSourceProfit:
		return walletType == models.WalletTypeProfit
	case models.GoodsDeductionSourceServiceFee:
		return walletType == models.WalletTypeServiceFee
	case models.GoodsDeductionSourceBoth:
		return walletType == models.WalletTypeProfit || walletType == models.WalletTypeServiceFee
	default:
		return false
	}
}

// executeDeduction 执行单笔扣款
func (s *GoodsDeductionService) executeDeduction(
	deduction *models.GoodsDeduction,
	wallet *repository.Wallet,
	amount int64,
	triggerType string,
	triggerTransactionID *int64,
	triggerProfitID *int64,
) (int64, error) {
	// 检查钱包余额
	if wallet.Balance < amount {
		amount = wallet.Balance
	}

	if amount <= 0 {
		return 0, nil
	}

	// 扣减钱包余额
	if err := s.walletRepo.UpdateBalance(wallet.ID, -amount); err != nil {
		return 0, fmt.Errorf("扣减钱包余额失败: %w", err)
	}

	// 更新货款代扣已扣金额
	if err := s.deductionRepo.UpdateDeductedAmount(deduction.ID, amount); err != nil {
		// 回滚钱包余额
		s.walletRepo.UpdateBalance(wallet.ID, amount)
		return 0, fmt.Errorf("更新已扣金额失败: %w", err)
	}

	// 计算新的累计已扣和剩余待扣
	newDeducted := deduction.DeductedAmount + amount
	newRemaining := deduction.RemainingAmount - amount

	// 创建扣款明细
	detail := &models.GoodsDeductionDetail{
		DeductionID:          deduction.ID,
		DeductionNo:          deduction.DeductionNo,
		Amount:               amount,
		WalletType:           wallet.WalletType,
		ChannelID:            &wallet.ChannelID,
		WalletBalanceBefore:  wallet.Balance,
		WalletBalanceAfter:   wallet.Balance - amount,
		CumulativeDeducted:   newDeducted,
		RemainingAfter:       newRemaining,
		TriggerType:          triggerType,
		TriggerTransactionID: triggerTransactionID,
		TriggerProfitID:      triggerProfitID,
		CreatedAt:            time.Now(),
	}

	if err := s.detailRepo.Create(detail); err != nil {
		log.Printf("[GoodsDeductionService] Create detail failed: %v", err)
	}

	// 记录钱包流水
	walletLog := &repository.WalletLog{
		WalletID:      wallet.ID,
		AgentID:       wallet.AgentID,
		WalletType:    wallet.WalletType,
		LogType:       7, // 货款代扣
		Amount:        -amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance - amount,
		RefType:       "goods_deduction_detail",
		RefID:         detail.ID,
		Remark:        fmt.Sprintf("货款代扣 - %s", deduction.DeductionNo),
		CreatedAt:     time.Now(),
	}
	s.walletLogRepo.Create(walletLog)

	// 将扣款金额转入发起方钱包
	if err := s.transferToFromAgent(deduction.FromAgentID, wallet.ChannelID, amount, detail.ID); err != nil {
		log.Printf("[GoodsDeductionService] Transfer to from agent failed: %v", err)
	}

	// 发送扣款通知
	s.sendDeductionNotification(deduction, detail, amount)

	// 检查是否已完成
	if newRemaining <= 0 {
		s.deductionRepo.UpdateStatus(deduction.ID, models.GoodsDeductionStatusCompleted)
		s.sendCompletionNotification(deduction)
	}

	log.Printf("[GoodsDeductionService] Deduction executed: deduction=%d, amount=%d, remaining=%d",
		deduction.ID, amount, newRemaining)

	return amount, nil
}

// transferToFromAgent 将扣款金额转入发起方钱包
func (s *GoodsDeductionService) transferToFromAgent(agentID int64, channelID int64, amount int64, detailID int64) error {
	// 转入发起方的分润钱包
	wallet, err := s.walletRepo.FindByAgentAndType(agentID, channelID, models.WalletTypeProfit)
	if err != nil || wallet == nil {
		return fmt.Errorf("发起方钱包不存在")
	}

	// 增加余额
	if err := s.walletRepo.UpdateBalance(wallet.ID, amount); err != nil {
		return fmt.Errorf("增加发起方余额失败: %w", err)
	}

	// 记录钱包流水
	walletLog := &repository.WalletLog{
		WalletID:      wallet.ID,
		AgentID:       agentID,
		WalletType:    wallet.WalletType,
		LogType:       7, // 货款代扣
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance + amount,
		RefType:       "goods_deduction_detail",
		RefID:         detailID,
		Remark:        "货款代扣收款",
		CreatedAt:     time.Now(),
	}
	return s.walletLogRepo.Create(walletLog)
}

// sendDeductionNotification 发送扣款通知
func (s *GoodsDeductionService) sendDeductionNotification(deduction *models.GoodsDeduction, detail *models.GoodsDeductionDetail, amount int64) {
	notification := &models.GoodsDeductionNotification{
		DeductionID: deduction.ID,
		DetailID:    &detail.ID,
		AgentID:     deduction.ToAgentID,
		NotifyType:  models.GoodsDeductionNotifyTypeDeduction,
		Title:       "货款代扣扣款通知",
		Content:     fmt.Sprintf("您的货款代扣已扣款%.2f元，累计已扣%.2f元，剩余%.2f元", float64(amount)/100, float64(detail.CumulativeDeducted)/100, float64(detail.RemainingAfter)/100),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}
	s.notificationRepo.Create(notification)
}

// sendCompletionNotification 发送完成通知
func (s *GoodsDeductionService) sendCompletionNotification(deduction *models.GoodsDeduction) {
	// 通知接收方
	notification := &models.GoodsDeductionNotification{
		DeductionID: deduction.ID,
		AgentID:     deduction.ToAgentID,
		NotifyType:  models.GoodsDeductionNotifyTypeCompleted,
		Title:       "货款代扣已完成",
		Content:     fmt.Sprintf("您的货款代扣已全部扣完，总金额%.2f元", float64(deduction.TotalAmount)/100),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}
	s.notificationRepo.Create(notification)

	// 通知发起方
	notificationFrom := &models.GoodsDeductionNotification{
		DeductionID: deduction.ID,
		AgentID:     deduction.FromAgentID,
		NotifyType:  models.GoodsDeductionNotifyTypeCompleted,
		Title:       "货款代扣已完成",
		Content:     fmt.Sprintf("您发起的货款代扣已全部扣完，总金额%.2f元", float64(deduction.TotalAmount)/100),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}
	s.notificationRepo.Create(notificationFrom)
}

// fillAgentNames 填充代理商名称
func (s *GoodsDeductionService) fillAgentNames(deduction *models.GoodsDeduction) {
	if fromAgent, err := s.agentRepo.FindByID(deduction.FromAgentID); err == nil && fromAgent != nil {
		deduction.FromAgentName = fromAgent.AgentName
	}
	if toAgent, err := s.agentRepo.FindByID(deduction.ToAgentID); err == nil && toAgent != nil {
		deduction.ToAgentName = toAgent.AgentName
	}
}

// GetPendingDeductionAmount 获取代理商待扣货款金额（影响可提现金额）
func (s *GoodsDeductionService) GetPendingDeductionAmount(agentID int64) (int64, error) {
	summary, err := s.deductionRepo.GetSummaryByToAgent(agentID)
	if err != nil {
		return 0, err
	}
	return summary.RemainingAmount, nil
}

// GetNotifications 获取通知列表
func (s *GoodsDeductionService) GetNotifications(agentID int64, isRead *bool, page, pageSize int) ([]*models.GoodsDeductionNotification, int64, error) {
	offset := (page - 1) * pageSize
	return s.notificationRepo.FindByAgentID(agentID, isRead, pageSize, offset)
}

// GetUnreadNotificationCount 获取未读通知数量
func (s *GoodsDeductionService) GetUnreadNotificationCount(agentID int64) (int64, error) {
	return s.notificationRepo.FindUnreadCount(agentID)
}

// MarkNotificationAsRead 标记通知为已读
func (s *GoodsDeductionService) MarkNotificationAsRead(notificationID int64) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func (s *GoodsDeductionService) MarkAllNotificationsAsRead(agentID int64) error {
	return s.notificationRepo.MarkAllAsRead(agentID)
}
