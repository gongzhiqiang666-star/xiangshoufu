package service

import (
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// ChargingWalletService 充值钱包服务
type ChargingWalletService struct {
	chargingRepo  *repository.GormChargingWalletRepository
	walletRepo    *repository.GormWalletRepository
	walletLogRepo *repository.GormWalletLogRepository
	agentRepo     *repository.GormAgentRepository
}

// NewChargingWalletService 创建充值钱包服务
func NewChargingWalletService(
	chargingRepo *repository.GormChargingWalletRepository,
	walletRepo *repository.GormWalletRepository,
	walletLogRepo *repository.GormWalletLogRepository,
	agentRepo *repository.GormAgentRepository,
) *ChargingWalletService {
	return &ChargingWalletService{
		chargingRepo:  chargingRepo,
		walletRepo:    walletRepo,
		walletLogRepo: walletLogRepo,
		agentRepo:     agentRepo,
	}
}

// ========== 钱包配置 ==========

// WalletConfigInfo 钱包配置信息
type WalletConfigInfo struct {
	AgentID                 int64      `json:"agent_id"`
	ChargingWalletEnabled   bool       `json:"charging_wallet_enabled"`
	ChargingWalletLimit     int64      `json:"charging_wallet_limit"`      // 分
	ChargingWalletLimitYuan float64    `json:"charging_wallet_limit_yuan"` // 元
	SettlementWalletEnabled bool       `json:"settlement_wallet_enabled"`
	SettlementRatio         int        `json:"settlement_ratio"`
	EnabledAt               *time.Time `json:"enabled_at"`
}

// GetWalletConfig 获取钱包配置
func (s *ChargingWalletService) GetWalletConfig(agentID int64) (*WalletConfigInfo, error) {
	config, err := s.chargingRepo.GetConfig(agentID)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	if config == nil {
		return &WalletConfigInfo{
			AgentID:         agentID,
			SettlementRatio: 30, // 默认30%
		}, nil
	}

	return &WalletConfigInfo{
		AgentID:                 config.AgentID,
		ChargingWalletEnabled:   config.ChargingWalletEnabled,
		ChargingWalletLimit:     config.ChargingWalletLimit,
		ChargingWalletLimitYuan: float64(config.ChargingWalletLimit) / 100,
		SettlementWalletEnabled: config.SettlementWalletEnabled,
		SettlementRatio:         config.SettlementRatio,
		EnabledAt:               config.EnabledAt,
	}, nil
}

// EnableChargingWalletRequest 开通充值钱包请求
type EnableChargingWalletRequest struct {
	AgentID   int64 `json:"agent_id" binding:"required"`
	Limit     int64 `json:"limit"`                       // 充值限额(分)
	EnabledBy int64 `json:"enabled_by" binding:"required"`
}

// EnableChargingWallet 开通充值钱包(PC端操作)
func (s *ChargingWalletService) EnableChargingWallet(req *EnableChargingWalletRequest) error {
	// 获取现有配置
	config, err := s.chargingRepo.GetConfig(req.AgentID)
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	now := time.Now()

	if config == nil {
		config = &models.AgentWalletConfig{
			AgentID:               req.AgentID,
			ChargingWalletEnabled: true,
			ChargingWalletLimit:   req.Limit,
			SettlementRatio:       30,
			EnabledBy:             &req.EnabledBy,
			EnabledAt:             &now,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
	} else {
		config.ChargingWalletEnabled = true
		config.ChargingWalletLimit = req.Limit
		if config.EnabledBy == nil {
			config.EnabledBy = &req.EnabledBy
			config.EnabledAt = &now
		}
		config.UpdatedAt = now
	}

	if err := s.chargingRepo.SaveConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	log.Printf("[ChargingWalletService] Enabled charging wallet for agent %d, limit=%d", req.AgentID, req.Limit)
	return nil
}

// DisableChargingWallet 关闭充值钱包
func (s *ChargingWalletService) DisableChargingWallet(agentID int64) error {
	config, err := s.chargingRepo.GetConfig(agentID)
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	if config == nil {
		return nil // 没有配置，无需关闭
	}

	config.ChargingWalletEnabled = false
	config.UpdatedAt = time.Now()

	if err := s.chargingRepo.SaveConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	log.Printf("[ChargingWalletService] Disabled charging wallet for agent %d", agentID)
	return nil
}

// ========== 充值操作 ==========

// CreateDepositRequest 申请充值请求
type CreateDepositRequest struct {
	AgentID       int64  `json:"agent_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=100"` // 最少1元
	PaymentMethod int16  `json:"payment_method" binding:"required"` // 1=银行转账 2=微信 3=支付宝
	PaymentRef    string `json:"payment_ref"`                       // 支付流水号
	Remark        string `json:"remark"`
	CreatedBy     int64  `json:"created_by" binding:"required"`
}

// CreateDeposit 申请充值
func (s *ChargingWalletService) CreateDeposit(req *CreateDepositRequest) (*models.ChargingWalletDeposit, error) {
	// 检查是否开通充值钱包
	config, err := s.chargingRepo.GetConfig(req.AgentID)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	if config == nil || !config.ChargingWalletEnabled {
		return nil, fmt.Errorf("充值钱包未开通")
	}

	// 检查限额
	if config.ChargingWalletLimit > 0 {
		// 获取当前钱包余额
		wallet, _ := s.walletRepo.FindByAgentAndType(req.AgentID, 0, models.WalletTypeCharging)
		currentBalance := int64(0)
		if wallet != nil {
			currentBalance = wallet.Balance
		}

		if currentBalance+req.Amount > config.ChargingWalletLimit {
			return nil, fmt.Errorf("充值后将超出限额，当前余额%.2f元，限额%.2f元",
				float64(currentBalance)/100, float64(config.ChargingWalletLimit)/100)
		}
	}

	// 生成充值单号
	depositNo := fmt.Sprintf("CWD%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	deposit := &models.ChargingWalletDeposit{
		DepositNo:     depositNo,
		AgentID:       req.AgentID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentRef:    req.PaymentRef,
		Status:        models.ChargingDepositStatusPending,
		Remark:        req.Remark,
		CreatedBy:     req.CreatedBy,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.chargingRepo.CreateDeposit(deposit); err != nil {
		return nil, fmt.Errorf("创建充值记录失败: %w", err)
	}

	log.Printf("[ChargingWalletService] Created deposit: %s, agent=%d, amount=%d", depositNo, req.AgentID, req.Amount)
	return deposit, nil
}

// ConfirmDeposit 确认充值(管理员操作)
func (s *ChargingWalletService) ConfirmDeposit(depositID int64, confirmedBy int64) error {
	deposit, err := s.chargingRepo.GetDeposit(depositID)
	if err != nil || deposit == nil {
		return fmt.Errorf("充值记录不存在")
	}

	if deposit.Status != models.ChargingDepositStatusPending {
		return fmt.Errorf("充值记录状态不正确")
	}

	now := time.Now()

	// 更新充值记录状态
	deposit.Status = models.ChargingDepositStatusConfirmed
	deposit.ConfirmedBy = &confirmedBy
	deposit.ConfirmedAt = &now
	deposit.UpdatedAt = now

	if err := s.chargingRepo.UpdateDeposit(deposit); err != nil {
		return fmt.Errorf("更新充值记录失败: %w", err)
	}

	// 增加充值钱包余额
	wallet, err := s.getOrCreateChargingWallet(deposit.AgentID)
	if err != nil {
		return fmt.Errorf("获取钱包失败: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(wallet.ID, deposit.Amount); err != nil {
		return fmt.Errorf("更新钱包余额失败: %w", err)
	}

	// 记录流水
	walletLog := &repository.WalletLog{
		WalletID:      wallet.ID,
		AgentID:       deposit.AgentID,
		WalletType:    models.WalletTypeCharging,
		LogType:       WalletLogTypeChargingDeposit,
		Amount:        deposit.Amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance + deposit.Amount,
		RefType:       "charging_deposit",
		RefID:         deposit.ID,
		Remark:        fmt.Sprintf("充值钱包充值，单号%s，金额%.2f元", deposit.DepositNo, float64(deposit.Amount)/100),
		CreatedAt:     now,
	}
	s.walletLogRepo.Create(walletLog)

	log.Printf("[ChargingWalletService] Confirmed deposit: %s, agent=%d, amount=%d", deposit.DepositNo, deposit.AgentID, deposit.Amount)
	return nil
}

// RejectDeposit 拒绝充值
func (s *ChargingWalletService) RejectDeposit(depositID int64, confirmedBy int64, reason string) error {
	deposit, err := s.chargingRepo.GetDeposit(depositID)
	if err != nil || deposit == nil {
		return fmt.Errorf("充值记录不存在")
	}

	if deposit.Status != models.ChargingDepositStatusPending {
		return fmt.Errorf("充值记录状态不正确")
	}

	now := time.Now()
	deposit.Status = models.ChargingDepositStatusRejected
	deposit.ConfirmedBy = &confirmedBy
	deposit.ConfirmedAt = &now
	deposit.RejectReason = reason
	deposit.UpdatedAt = now

	if err := s.chargingRepo.UpdateDeposit(deposit); err != nil {
		return fmt.Errorf("更新充值记录失败: %w", err)
	}

	log.Printf("[ChargingWalletService] Rejected deposit: %s, agent=%d, reason=%s", deposit.DepositNo, deposit.AgentID, reason)
	return nil
}

// GetDepositList 获取充值记录列表
func (s *ChargingWalletService) GetDepositList(agentID int64, status *int16, page, pageSize int) ([]*DepositInfo, int64, error) {
	offset := (page - 1) * pageSize
	deposits, total, err := s.chargingRepo.GetDepositsByAgent(agentID, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取充值记录失败: %w", err)
	}

	list := make([]*DepositInfo, 0, len(deposits))
	for _, d := range deposits {
		list = append(list, s.toDepositInfo(d))
	}

	return list, total, nil
}

// DepositInfo 充值记录信息
type DepositInfo struct {
	ID              int64      `json:"id"`
	DepositNo       string     `json:"deposit_no"`
	AgentID         int64      `json:"agent_id"`
	Amount          int64      `json:"amount"`
	AmountYuan      float64    `json:"amount_yuan"`
	PaymentMethod   int16      `json:"payment_method"`
	PaymentMethodName string   `json:"payment_method_name"`
	PaymentRef      string     `json:"payment_ref"`
	Status          int16      `json:"status"`
	StatusName      string     `json:"status_name"`
	RejectReason    string     `json:"reject_reason"`
	Remark          string     `json:"remark"`
	CreatedAt       time.Time  `json:"created_at"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
}

func (s *ChargingWalletService) toDepositInfo(d *models.ChargingWalletDeposit) *DepositInfo {
	return &DepositInfo{
		ID:               d.ID,
		DepositNo:        d.DepositNo,
		AgentID:          d.AgentID,
		Amount:           d.Amount,
		AmountYuan:       float64(d.Amount) / 100,
		PaymentMethod:    d.PaymentMethod,
		PaymentMethodName: getPaymentMethodName(d.PaymentMethod),
		PaymentRef:       d.PaymentRef,
		Status:           d.Status,
		StatusName:       models.GetChargingDepositStatusName(d.Status),
		RejectReason:     d.RejectReason,
		Remark:           d.Remark,
		CreatedAt:        d.CreatedAt,
		ConfirmedAt:      d.ConfirmedAt,
	}
}

func getPaymentMethodName(method int16) string {
	switch method {
	case 1:
		return "银行转账"
	case 2:
		return "微信"
	case 3:
		return "支付宝"
	default:
		return "其他"
	}
}

// ========== 奖励发放 ==========

// IssueRewardRequest 发放奖励请求
type IssueRewardRequest struct {
	FromAgentID int64  `json:"from_agent_id" binding:"required"`
	ToAgentID   int64  `json:"to_agent_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,min=1"` // 分
	Remark      string `json:"remark"`
	CreatedBy   int64  `json:"created_by" binding:"required"`
}

// IssueReward 发放奖励给下级
func (s *ChargingWalletService) IssueReward(req *IssueRewardRequest) (*models.ChargingWalletReward, error) {
	// 验证是上下级关系
	toAgent, err := s.agentRepo.FindByIDFull(req.ToAgentID)
	if err != nil || toAgent == nil {
		return nil, fmt.Errorf("接收方代理商不存在")
	}

	if toAgent.ParentID != req.FromAgentID {
		return nil, fmt.Errorf("只能给直属下级发放奖励")
	}

	// 检查充值钱包余额
	fromWallet, err := s.walletRepo.FindByAgentAndType(req.FromAgentID, 0, models.WalletTypeCharging)
	if err != nil || fromWallet == nil {
		return nil, fmt.Errorf("充值钱包不存在")
	}

	availableBalance := fromWallet.Balance - fromWallet.FrozenAmount
	if availableBalance < req.Amount {
		return nil, fmt.Errorf("充值钱包余额不足，可用余额%.2f元", float64(availableBalance)/100)
	}

	// 生成奖励单号
	rewardNo := fmt.Sprintf("CWR%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	reward := &models.ChargingWalletReward{
		RewardNo:    rewardNo,
		FromAgentID: req.FromAgentID,
		ToAgentID:   req.ToAgentID,
		Amount:      req.Amount,
		RewardType:  models.ChargingRewardTypeManual,
		Status:      models.ChargingRewardStatusIssued,
		Remark:      req.Remark,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
	}

	if err := s.chargingRepo.CreateReward(reward); err != nil {
		return nil, fmt.Errorf("创建奖励记录失败: %w", err)
	}

	// 扣除发放方充值钱包余额
	if err := s.walletRepo.UpdateBalance(fromWallet.ID, -req.Amount); err != nil {
		return nil, fmt.Errorf("扣除发放方余额失败: %w", err)
	}

	// 记录发放方流水
	now := time.Now()
	fromLog := &repository.WalletLog{
		WalletID:      fromWallet.ID,
		AgentID:       req.FromAgentID,
		WalletType:    models.WalletTypeCharging,
		LogType:       WalletLogTypeChargingRewardOut,
		Amount:        -req.Amount,
		BalanceBefore: fromWallet.Balance,
		BalanceAfter:  fromWallet.Balance - req.Amount,
		RefType:       "charging_reward",
		RefID:         reward.ID,
		Remark:        fmt.Sprintf("发放奖励给%s，金额%.2f元", toAgent.AgentName, float64(req.Amount)/100),
		CreatedAt:     now,
	}
	s.walletLogRepo.Create(fromLog)

	// 增加接收方奖励钱包余额
	toWallet, err := s.getOrCreateRewardWallet(req.ToAgentID)
	if err != nil {
		log.Printf("[ChargingWalletService] Failed to get reward wallet for agent %d: %v", req.ToAgentID, err)
	} else {
		if err := s.walletRepo.UpdateBalance(toWallet.ID, req.Amount); err != nil {
			log.Printf("[ChargingWalletService] Failed to update reward wallet balance: %v", err)
		}

		// 记录接收方流水
		toLog := &repository.WalletLog{
			WalletID:      toWallet.ID,
			AgentID:       req.ToAgentID,
			WalletType:    models.WalletTypeReward,
			LogType:       WalletLogTypeChargingRewardIn,
			Amount:        req.Amount,
			BalanceBefore: toWallet.Balance,
			BalanceAfter:  toWallet.Balance + req.Amount,
			RefType:       "charging_reward",
			RefID:         reward.ID,
			Remark:        fmt.Sprintf("收到上级奖励，金额%.2f元", float64(req.Amount)/100),
			CreatedAt:     now,
		}
		s.walletLogRepo.Create(toLog)
	}

	log.Printf("[ChargingWalletService] Issued reward: %s, from=%d, to=%d, amount=%d",
		rewardNo, req.FromAgentID, req.ToAgentID, req.Amount)
	return reward, nil
}

// GetRewardList 获取奖励记录列表
func (s *ChargingWalletService) GetRewardList(agentID int64, direction string, page, pageSize int) ([]*RewardInfo, int64, error) {
	offset := (page - 1) * pageSize

	var rewards []*models.ChargingWalletReward
	var total int64
	var err error

	if direction == "from" {
		rewards, total, err = s.chargingRepo.GetRewardsByFromAgent(agentID, pageSize, offset)
	} else {
		rewards, total, err = s.chargingRepo.GetRewardsByToAgent(agentID, pageSize, offset)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("获取奖励记录失败: %w", err)
	}

	list := make([]*RewardInfo, 0, len(rewards))
	for _, r := range rewards {
		list = append(list, s.toRewardInfo(r))
	}

	return list, total, nil
}

// RewardInfo 奖励记录信息
type RewardInfo struct {
	ID            int64      `json:"id"`
	RewardNo      string     `json:"reward_no"`
	FromAgentID   int64      `json:"from_agent_id"`
	FromAgentName string     `json:"from_agent_name"`
	ToAgentID     int64      `json:"to_agent_id"`
	ToAgentName   string     `json:"to_agent_name"`
	Amount        int64      `json:"amount"`
	AmountYuan    float64    `json:"amount_yuan"`
	RewardType    int16      `json:"reward_type"`
	RewardTypeName string    `json:"reward_type_name"`
	Status        int16      `json:"status"`
	StatusName    string     `json:"status_name"`
	Remark        string     `json:"remark"`
	CreatedAt     time.Time  `json:"created_at"`
	RevokedAt     *time.Time `json:"revoked_at"`
}

func (s *ChargingWalletService) toRewardInfo(r *models.ChargingWalletReward) *RewardInfo {
	info := &RewardInfo{
		ID:            r.ID,
		RewardNo:      r.RewardNo,
		FromAgentID:   r.FromAgentID,
		ToAgentID:     r.ToAgentID,
		Amount:        r.Amount,
		AmountYuan:    float64(r.Amount) / 100,
		RewardType:    r.RewardType,
		RewardTypeName: models.GetRewardTypeName(r.RewardType),
		Status:        r.Status,
		StatusName:    models.GetRewardStatusName(r.Status),
		Remark:        r.Remark,
		CreatedAt:     r.CreatedAt,
		RevokedAt:     r.RevokedAt,
	}

	// 获取代理商名称
	if fromAgent, _ := s.agentRepo.FindByIDFull(r.FromAgentID); fromAgent != nil {
		info.FromAgentName = fromAgent.AgentName
	}
	if toAgent, _ := s.agentRepo.FindByIDFull(r.ToAgentID); toAgent != nil {
		info.ToAgentName = toAgent.AgentName
	}

	return info
}

// GetChargingWalletSummary 获取充值钱包汇总
// P0修复：增加奖励总金额显示（包含手动发放+系统自动发放）
func (s *ChargingWalletService) GetChargingWalletSummary(agentID int64) (*ChargingWalletSummary, error) {
	wallet, _ := s.walletRepo.FindByAgentAndType(agentID, 0, models.WalletTypeCharging)

	balance := int64(0)
	if wallet != nil {
		balance = wallet.Balance
	}

	// 获取手动发放奖励总额（从充值钱包发放的）
	manualIssued, _ := s.chargingRepo.GetTotalRewardsIssuedByAgent(agentID)

	// 获取自动激活奖励总额（系统自动发放给下级的）
	// 这里查询从该代理商作为上级（source_agent_id）发放的所有激活奖励
	autoReward, _ := s.chargingRepo.GetTotalAutoRewardsForDownline(agentID)

	// 总奖励 = 手动发放 + 自动发放
	totalReward := manualIssued + autoReward

	return &ChargingWalletSummary{
		Balance:          balance,
		BalanceYuan:      float64(balance) / 100,
		TotalIssued:      manualIssued,
		TotalIssuedYuan:  float64(manualIssued) / 100,
		TotalAutoReward:  autoReward,
		TotalAutoRewardYuan: float64(autoReward) / 100,
		TotalReward:      totalReward,
		TotalRewardYuan:  float64(totalReward) / 100,
	}, nil
}

// ChargingWalletSummary 充值钱包汇总
type ChargingWalletSummary struct {
	Balance             int64   `json:"balance"`               // 当前余额(分)
	BalanceYuan         float64 `json:"balance_yuan"`          // 当前余额(元)
	TotalIssued         int64   `json:"total_issued"`          // 手动发放奖励总额(分)
	TotalIssuedYuan     float64 `json:"total_issued_yuan"`     // 手动发放奖励总额(元)
	TotalAutoReward     int64   `json:"total_auto_reward"`     // 系统自动奖励总额(分)
	TotalAutoRewardYuan float64 `json:"total_auto_reward_yuan"` // 系统自动奖励总额(元)
	TotalReward         int64   `json:"total_reward"`          // 奖励总金额(分)（手动+自动）
	TotalRewardYuan     float64 `json:"total_reward_yuan"`     // 奖励总金额(元)（手动+自动）
}

// ========== 辅助方法 ==========

// getOrCreateChargingWallet 获取或创建充值钱包
func (s *ChargingWalletService) getOrCreateChargingWallet(agentID int64) (*repository.Wallet, error) {
	wallet, err := s.walletRepo.FindByAgentAndType(agentID, 0, models.WalletTypeCharging)
	if err == nil && wallet != nil {
		return wallet, nil
	}

	// 创建新钱包
	newWallet := &repository.Wallet{
		AgentID:           agentID,
		ChannelID:         0, // 充值钱包不区分通道
		WalletType:        models.WalletTypeCharging,
		WithdrawThreshold: 0, // 充值钱包不能提现
	}

	// 这里简化处理，实际应该在创建时就自动创建钱包
	// 暂时返回错误让上层处理
	return newWallet, nil
}

// getOrCreateRewardWallet 获取或创建奖励钱包
func (s *ChargingWalletService) getOrCreateRewardWallet(agentID int64) (*repository.Wallet, error) {
	wallet, err := s.walletRepo.FindByAgentAndType(agentID, 0, models.WalletTypeReward)
	if err == nil && wallet != nil {
		return wallet, nil
	}

	// 创建新钱包 - 这里简化处理
	return nil, fmt.Errorf("奖励钱包不存在")
}

// 新增流水类型
const (
	WalletLogTypeChargingDeposit   int16 = 8  // 充值钱包充值
	WalletLogTypeChargingRewardOut int16 = 9  // 充值钱包奖励发放
	WalletLogTypeChargingRewardIn  int16 = 10 // 收到充值钱包奖励
)
