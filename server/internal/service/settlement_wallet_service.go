package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// SettlementWalletService 沉淀钱包服务
type SettlementWalletService struct {
	settlementRepo *repository.GormSettlementWalletRepository
	chargingRepo   *repository.GormChargingWalletRepository
	walletRepo     *repository.GormWalletRepository
	walletLogRepo  *repository.GormWalletLogRepository
	agentRepo      *repository.GormAgentRepository
}

// NewSettlementWalletService 创建沉淀钱包服务
func NewSettlementWalletService(
	settlementRepo *repository.GormSettlementWalletRepository,
	chargingRepo *repository.GormChargingWalletRepository,
	walletRepo *repository.GormWalletRepository,
	walletLogRepo *repository.GormWalletLogRepository,
	agentRepo *repository.GormAgentRepository,
) *SettlementWalletService {
	return &SettlementWalletService{
		settlementRepo: settlementRepo,
		chargingRepo:   chargingRepo,
		walletRepo:     walletRepo,
		walletLogRepo:  walletLogRepo,
		agentRepo:      agentRepo,
	}
}

// ========== 钱包配置 ==========

// EnableSettlementWalletRequest 开通沉淀钱包请求
type EnableSettlementWalletRequest struct {
	AgentID   int64 `json:"agent_id" binding:"required"`
	Ratio     int   `json:"ratio" binding:"required,min=1,max=100"` // 沉淀比例(1-100)
	EnabledBy int64 `json:"enabled_by" binding:"required"`
}

// EnableSettlementWallet 开通沉淀钱包(PC端操作)
func (s *SettlementWalletService) EnableSettlementWallet(req *EnableSettlementWalletRequest) error {
	// 获取现有配置
	config, err := s.chargingRepo.GetConfig(req.AgentID)
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	now := time.Now()

	if config == nil {
		config = &models.AgentWalletConfig{
			AgentID:                 req.AgentID,
			SettlementWalletEnabled: true,
			SettlementRatio:         req.Ratio,
			EnabledBy:               &req.EnabledBy,
			EnabledAt:               &now,
			CreatedAt:               now,
			UpdatedAt:               now,
		}
	} else {
		config.SettlementWalletEnabled = true
		config.SettlementRatio = req.Ratio
		if config.EnabledBy == nil {
			config.EnabledBy = &req.EnabledBy
			config.EnabledAt = &now
		}
		config.UpdatedAt = now
	}

	if err := s.chargingRepo.SaveConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	log.Printf("[SettlementWalletService] Enabled settlement wallet for agent %d, ratio=%d%%", req.AgentID, req.Ratio)
	return nil
}

// DisableSettlementWallet 关闭沉淀钱包
func (s *SettlementWalletService) DisableSettlementWallet(agentID int64) error {
	config, err := s.chargingRepo.GetConfig(agentID)
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	if config == nil {
		return nil
	}

	// 检查是否有待归还的使用记录
	pending, _ := s.settlementRepo.GetPendingReturn(agentID)
	if len(pending) > 0 {
		return fmt.Errorf("还有待归还的沉淀款项，请先归还后再关闭")
	}

	config.SettlementWalletEnabled = false
	config.UpdatedAt = time.Now()

	if err := s.chargingRepo.SaveConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	log.Printf("[SettlementWalletService] Disabled settlement wallet for agent %d", agentID)
	return nil
}

// UpdateSettlementRatio 更新沉淀比例
func (s *SettlementWalletService) UpdateSettlementRatio(agentID int64, ratio int) error {
	if ratio < 1 || ratio > 100 {
		return fmt.Errorf("沉淀比例必须在1-100之间")
	}

	config, err := s.chargingRepo.GetConfig(agentID)
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	if config == nil {
		return fmt.Errorf("沉淀钱包未开通")
	}

	config.SettlementRatio = ratio
	config.UpdatedAt = time.Now()

	if err := s.chargingRepo.SaveConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	log.Printf("[SettlementWalletService] Updated settlement ratio for agent %d to %d%%", agentID, ratio)
	return nil
}

// ========== 沉淀钱包汇总 ==========

// SettlementWalletSummary 沉淀钱包汇总
type SettlementWalletSummary struct {
	SubordinateTotalBalance   int64   `json:"subordinate_total_balance"`    // 下级未提现总额(分)
	SubordinateTotalBalanceYuan float64 `json:"subordinate_total_balance_yuan"`
	SettlementRatio           int     `json:"settlement_ratio"`             // 沉淀比例
	AvailableAmount           int64   `json:"available_amount"`             // 可用沉淀额度(分)
	AvailableAmountYuan       float64 `json:"available_amount_yuan"`
	UsedAmount                int64   `json:"used_amount"`                  // 已使用沉淀额(分)
	UsedAmountYuan            float64 `json:"used_amount_yuan"`
	PendingReturnAmount       int64   `json:"pending_return_amount"`        // 待归还金额(分)
	PendingReturnAmountYuan   float64 `json:"pending_return_amount_yuan"`
}

// GetSettlementWalletSummary 获取沉淀钱包汇总
func (s *SettlementWalletService) GetSettlementWalletSummary(agentID int64) (*SettlementWalletSummary, error) {
	// 获取配置
	config, err := s.chargingRepo.GetConfig(agentID)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	ratio := 30 // 默认30%
	if config != nil && config.SettlementRatio > 0 {
		ratio = config.SettlementRatio
	}

	// 获取下级未提现余额汇总
	subordinateTotal, _ := s.settlementRepo.GetSubordinateUnwithdrawBalance(agentID)

	// 计算可用沉淀额度
	availableAmount := subordinateTotal * int64(ratio) / 100

	// 获取已使用和待归还金额
	usedAmount := int64(0)
	pendingReturnAmount := int64(0)

	pending, _ := s.settlementRepo.GetPendingReturn(agentID)
	for _, p := range pending {
		pendingReturnAmount += p.Amount
	}

	// 获取沉淀钱包余额作为已使用金额
	wallet, _ := s.walletRepo.FindByAgentAndType(agentID, 0, models.WalletTypeSettlement)
	if wallet != nil {
		usedAmount = wallet.TotalIncome
	}

	return &SettlementWalletSummary{
		SubordinateTotalBalance:     subordinateTotal,
		SubordinateTotalBalanceYuan: float64(subordinateTotal) / 100,
		SettlementRatio:             ratio,
		AvailableAmount:             availableAmount,
		AvailableAmountYuan:         float64(availableAmount) / 100,
		UsedAmount:                  usedAmount,
		UsedAmountYuan:              float64(usedAmount) / 100,
		PendingReturnAmount:         pendingReturnAmount,
		PendingReturnAmountYuan:     float64(pendingReturnAmount) / 100,
	}, nil
}

// GetSubordinateBalances 获取下级余额明细
func (s *SettlementWalletService) GetSubordinateBalances(agentID int64) ([]SubordinateBalanceInfo, error) {
	details, err := s.settlementRepo.GetSubordinateBalanceDetails(agentID)
	if err != nil {
		return nil, fmt.Errorf("获取下级余额失败: %w", err)
	}

	list := make([]SubordinateBalanceInfo, 0, len(details))
	for _, d := range details {
		list = append(list, SubordinateBalanceInfo{
			AgentID:              d.AgentID,
			AgentName:            d.AgentName,
			AvailableBalance:     d.AvailableBalance,
			AvailableBalanceYuan: float64(d.AvailableBalance) / 100,
		})
	}

	return list, nil
}

// SubordinateBalanceInfo 下级余额信息
type SubordinateBalanceInfo struct {
	AgentID              int64   `json:"agent_id"`
	AgentName            string  `json:"agent_name"`
	AvailableBalance     int64   `json:"available_balance"`
	AvailableBalanceYuan float64 `json:"available_balance_yuan"`
}

// ========== 使用沉淀款 ==========

// UseSettlementRequest 使用沉淀款请求
type UseSettlementRequest struct {
	AgentID   int64  `json:"agent_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,min=100"` // 最少1元
	Remark    string `json:"remark"`
	CreatedBy int64  `json:"created_by" binding:"required"`
}

// UseSettlement 使用沉淀款
func (s *SettlementWalletService) UseSettlement(req *UseSettlementRequest) (*models.SettlementWalletUsage, error) {
	// 检查是否开通沉淀钱包
	config, err := s.chargingRepo.GetConfig(req.AgentID)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	if config == nil || !config.SettlementWalletEnabled {
		return nil, fmt.Errorf("沉淀钱包未开通")
	}

	// 获取可用沉淀额度
	summary, err := s.GetSettlementWalletSummary(req.AgentID)
	if err != nil {
		return nil, fmt.Errorf("获取汇总失败: %w", err)
	}

	remainingAvailable := summary.AvailableAmount - summary.UsedAmount
	if remainingAvailable < req.Amount {
		return nil, fmt.Errorf("可用沉淀额度不足，当前可用%.2f元", float64(remainingAvailable)/100)
	}

	// 获取下级余额明细（记录来源）
	details, _ := s.settlementRepo.GetSubordinateBalanceDetails(req.AgentID)
	sourceDetailsJSON, _ := json.Marshal(details)

	// 生成使用单号
	usageNo := fmt.Sprintf("SWU%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	usage := &models.SettlementWalletUsage{
		UsageNo:       usageNo,
		AgentID:       req.AgentID,
		Amount:        req.Amount,
		UsageType:     models.SettlementUsageTypeUse,
		SourceDetails: string(sourceDetailsJSON),
		Status:        models.SettlementUsageStatusNormal,
		Remark:        req.Remark,
		CreatedBy:     req.CreatedBy,
		CreatedAt:     time.Now(),
	}

	if err := s.settlementRepo.CreateUsage(usage); err != nil {
		return nil, fmt.Errorf("创建使用记录失败: %w", err)
	}

	// 增加沉淀钱包余额
	wallet, err := s.getOrCreateSettlementWallet(req.AgentID)
	if err != nil {
		log.Printf("[SettlementWalletService] Failed to get settlement wallet: %v", err)
	} else {
		if err := s.walletRepo.UpdateBalance(wallet.ID, req.Amount); err != nil {
			log.Printf("[SettlementWalletService] Failed to update wallet balance: %v", err)
		}

		// 记录流水
		now := time.Now()
		walletLog := &repository.WalletLog{
			WalletID:      wallet.ID,
			AgentID:       req.AgentID,
			WalletType:    models.WalletTypeSettlement,
			LogType:       WalletLogTypeSettlementUse,
			Amount:        req.Amount,
			BalanceBefore: wallet.Balance,
			BalanceAfter:  wallet.Balance + req.Amount,
			RefType:       "settlement_use",
			RefID:         usage.ID,
			Remark:        fmt.Sprintf("使用沉淀款，单号%s，金额%.2f元", usageNo, float64(req.Amount)/100),
			CreatedAt:     now,
		}
		s.walletLogRepo.Create(walletLog)
	}

	log.Printf("[SettlementWalletService] Used settlement: %s, agent=%d, amount=%d", usageNo, req.AgentID, req.Amount)
	return usage, nil
}

// ReturnSettlementRequest 归还沉淀款请求
type ReturnSettlementRequest struct {
	AgentID   int64  `json:"agent_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,min=1"` // 分
	Remark    string `json:"remark"`
	CreatedBy int64  `json:"created_by" binding:"required"`
}

// ReturnSettlement 归还沉淀款
func (s *SettlementWalletService) ReturnSettlement(req *ReturnSettlementRequest) (*models.SettlementWalletUsage, error) {
	// 获取沉淀钱包余额
	wallet, err := s.walletRepo.FindByAgentAndType(req.AgentID, 0, models.WalletTypeSettlement)
	if err != nil || wallet == nil {
		return nil, fmt.Errorf("沉淀钱包不存在")
	}

	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("沉淀钱包余额不足")
	}

	// 生成归还单号
	usageNo := fmt.Sprintf("SWR%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	usage := &models.SettlementWalletUsage{
		UsageNo:   usageNo,
		AgentID:   req.AgentID,
		Amount:    req.Amount,
		UsageType: models.SettlementUsageTypeReturn,
		Status:    models.SettlementUsageStatusNormal,
		Remark:    req.Remark,
		CreatedBy: req.CreatedBy,
		CreatedAt: time.Now(),
	}

	now := time.Now()
	usage.ReturnedAt = &now

	if err := s.settlementRepo.CreateUsage(usage); err != nil {
		return nil, fmt.Errorf("创建归还记录失败: %w", err)
	}

	// 扣除沉淀钱包余额
	if err := s.walletRepo.UpdateBalance(wallet.ID, -req.Amount); err != nil {
		return nil, fmt.Errorf("扣除钱包余额失败: %w", err)
	}

	// 记录流水
	walletLog := &repository.WalletLog{
		WalletID:      wallet.ID,
		AgentID:       req.AgentID,
		WalletType:    models.WalletTypeSettlement,
		LogType:       WalletLogTypeSettlementReturn,
		Amount:        -req.Amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance - req.Amount,
		RefType:       "settlement_return",
		RefID:         usage.ID,
		Remark:        fmt.Sprintf("归还沉淀款，单号%s，金额%.2f元", usageNo, float64(req.Amount)/100),
		CreatedAt:     now,
	}
	s.walletLogRepo.Create(walletLog)

	log.Printf("[SettlementWalletService] Returned settlement: %s, agent=%d, amount=%d", usageNo, req.AgentID, req.Amount)
	return usage, nil
}

// GetUsageList 获取使用记录列表
func (s *SettlementWalletService) GetUsageList(agentID int64, usageType *int16, page, pageSize int) ([]*UsageInfo, int64, error) {
	offset := (page - 1) * pageSize
	usages, total, err := s.settlementRepo.GetUsagesByAgent(agentID, usageType, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取使用记录失败: %w", err)
	}

	list := make([]*UsageInfo, 0, len(usages))
	for _, u := range usages {
		list = append(list, s.toUsageInfo(u))
	}

	return list, total, nil
}

// UsageInfo 使用记录信息
type UsageInfo struct {
	ID             int64      `json:"id"`
	UsageNo        string     `json:"usage_no"`
	AgentID        int64      `json:"agent_id"`
	Amount         int64      `json:"amount"`
	AmountYuan     float64    `json:"amount_yuan"`
	UsageType      int16      `json:"usage_type"`
	UsageTypeName  string     `json:"usage_type_name"`
	Status         int16      `json:"status"`
	StatusName     string     `json:"status_name"`
	ReturnDeadline *time.Time `json:"return_deadline"`
	ReturnedAt     *time.Time `json:"returned_at"`
	Remark         string     `json:"remark"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (s *SettlementWalletService) toUsageInfo(u *models.SettlementWalletUsage) *UsageInfo {
	return &UsageInfo{
		ID:             u.ID,
		UsageNo:        u.UsageNo,
		AgentID:        u.AgentID,
		Amount:         u.Amount,
		AmountYuan:     float64(u.Amount) / 100,
		UsageType:      u.UsageType,
		UsageTypeName:  getUsageTypeName(u.UsageType),
		Status:         u.Status,
		StatusName:     getUsageStatusName(u.Status),
		ReturnDeadline: u.ReturnDeadline,
		ReturnedAt:     u.ReturnedAt,
		Remark:         u.Remark,
		CreatedAt:      u.CreatedAt,
	}
}

func getUsageTypeName(usageType int16) string {
	switch usageType {
	case models.SettlementUsageTypeUse:
		return "使用"
	case models.SettlementUsageTypeReturn:
		return "归还"
	default:
		return "未知"
	}
}

func getUsageStatusName(status int16) string {
	switch status {
	case models.SettlementUsageStatusNormal:
		return "正常"
	case models.SettlementUsageStatusToReturn:
		return "待归还"
	default:
		return "未知"
	}
}

// ========== 辅助方法 ==========

// getOrCreateSettlementWallet 获取或创建沉淀钱包
func (s *SettlementWalletService) getOrCreateSettlementWallet(agentID int64) (*repository.Wallet, error) {
	wallet, err := s.walletRepo.FindByAgentAndType(agentID, 0, models.WalletTypeSettlement)
	if err == nil && wallet != nil {
		return wallet, nil
	}

	// 创建新钱包 - 这里简化处理
	newWallet := &repository.Wallet{
		AgentID:           agentID,
		ChannelID:         0,
		WalletType:        models.WalletTypeSettlement,
		WithdrawThreshold: 0, // 沉淀钱包不能直接提现
	}
	return newWallet, nil
}

// 新增流水类型
const (
	WalletLogTypeSettlementUse    int16 = 11 // 使用沉淀款
	WalletLogTypeSettlementReturn int16 = 12 // 归还沉淀款
)
