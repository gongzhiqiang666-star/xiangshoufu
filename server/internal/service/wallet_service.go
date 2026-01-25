package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// WalletService 钱包服务
type WalletService struct {
	walletRepo    *repository.GormWalletRepository
	walletLogRepo *repository.GormWalletLogRepository
	agentRepo     *repository.GormAgentRepository
}

// NewWalletService 创建钱包服务
func NewWalletService(
	walletRepo *repository.GormWalletRepository,
	walletLogRepo *repository.GormWalletLogRepository,
	agentRepo *repository.GormAgentRepository,
) *WalletService {
	return &WalletService{
		walletRepo:    walletRepo,
		walletLogRepo: walletLogRepo,
		agentRepo:     agentRepo,
	}
}

// WalletInfoDetail 钱包详细信息
type WalletInfoDetail struct {
	ID                int64   `json:"id"`
	ChannelID         int64   `json:"channel_id"`
	ChannelName       string  `json:"channel_name"`
	WalletType        int16   `json:"wallet_type"`
	WalletTypeName    string  `json:"wallet_type_name"`
	Balance           int64   `json:"balance"`             // 分
	BalanceYuan       float64 `json:"balance_yuan"`        // 元
	FrozenAmount      int64   `json:"frozen_amount"`       // 冻结金额
	TotalIncome       int64   `json:"total_income"`        // 累计收入
	TotalWithdraw     int64   `json:"total_withdraw"`      // 累计提现
	WithdrawThreshold int64   `json:"withdraw_threshold"`  // 提现门槛
	CanWithdraw       bool    `json:"can_withdraw"`        // 是否可提现
}

// GetWalletList 获取钱包列表
func (s *WalletService) GetWalletList(agentID int64) ([]*WalletInfoDetail, error) {
	wallets, err := s.walletRepo.FindByAgentID(agentID)
	if err != nil {
		return nil, fmt.Errorf("查询钱包失败: %w", err)
	}

	list := make([]*WalletInfoDetail, 0, len(wallets))
	for _, w := range wallets {
		info := &WalletInfoDetail{
			ID:                w.ID,
			ChannelID:         w.ChannelID,
			ChannelName:       getChannelName(w.ChannelID),
			WalletType:        w.WalletType,
			WalletTypeName:    getWalletTypeNameStr(w.WalletType),
			Balance:           w.Balance,
			BalanceYuan:       float64(w.Balance) / 100,
			FrozenAmount:      w.FrozenAmount,
			TotalIncome:       w.TotalIncome,
			TotalWithdraw:     w.TotalWithdraw,
			WithdrawThreshold: w.WithdrawThreshold,
			CanWithdraw:       w.Balance-w.FrozenAmount >= w.WithdrawThreshold,
		}
		list = append(list, info)
	}

	return list, nil
}

// WalletSummary 钱包汇总
type WalletSummary struct {
	TotalBalance      int64   `json:"total_balance"`       // 总余额（分）
	TotalBalanceYuan  float64 `json:"total_balance_yuan"`  // 总余额（元）
	TotalFrozen       int64   `json:"total_frozen"`        // 总冻结
	TotalIncome       int64   `json:"total_income"`        // 总收入
	TotalWithdraw     int64   `json:"total_withdraw"`      // 总提现
	AvailableBalance  int64   `json:"available_balance"`   // 可用余额
	WalletCount       int     `json:"wallet_count"`        // 钱包数量
}

// GetWalletSummary 获取钱包汇总
func (s *WalletService) GetWalletSummary(agentID int64) (*WalletSummary, error) {
	wallets, err := s.walletRepo.FindByAgentID(agentID)
	if err != nil {
		return nil, fmt.Errorf("查询钱包失败: %w", err)
	}

	summary := &WalletSummary{
		WalletCount: len(wallets),
	}

	for _, w := range wallets {
		summary.TotalBalance += w.Balance
		summary.TotalFrozen += w.FrozenAmount
		summary.TotalIncome += w.TotalIncome
		summary.TotalWithdraw += w.TotalWithdraw
	}

	summary.AvailableBalance = summary.TotalBalance - summary.TotalFrozen
	summary.TotalBalanceYuan = float64(summary.TotalBalance) / 100

	return summary, nil
}

// WalletLogInfo 钱包流水信息
type WalletLogInfo struct {
	ID            int64     `json:"id"`
	LogType       int16     `json:"log_type"`
	LogTypeName   string    `json:"log_type_name"`
	Amount        int64     `json:"amount"`         // 分（可为负）
	AmountYuan    float64   `json:"amount_yuan"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	RefType       string    `json:"ref_type"`
	RefID         int64     `json:"ref_id"`
	Remark        string    `json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetWalletLogsRequest 获取钱包流水请求
type GetWalletLogsRequest struct {
	WalletID  int64      `form:"wallet_id"`
	LogType   *int16     `form:"log_type"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"`
	Page      int        `form:"page,default=1"`
	PageSize  int        `form:"page_size,default=20"`
}

// GetWalletLogs 获取钱包流水
func (s *WalletService) GetWalletLogs(agentID int64, req *GetWalletLogsRequest) ([]*WalletLogInfo, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	logs, total, err := s.walletLogRepo.FindByAgentID(agentID, req.WalletID, req.LogType, req.StartTime, req.EndTime, req.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询流水失败: %w", err)
	}

	list := make([]*WalletLogInfo, 0, len(logs))
	for _, l := range logs {
		list = append(list, &WalletLogInfo{
			ID:            l.ID,
			LogType:       l.LogType,
			LogTypeName:   getLogTypeName(l.LogType),
			Amount:        l.Amount,
			AmountYuan:    float64(l.Amount) / 100,
			BalanceBefore: l.BalanceBefore,
			BalanceAfter:  l.BalanceAfter,
			RefType:       l.RefType,
			RefID:         l.RefID,
			Remark:        l.Remark,
			CreatedAt:     l.CreatedAt,
		})
	}

	return list, total, nil
}

// 钱包流水类型
const (
	WalletLogTypeProfitIn        int16 = 1  // 分润入账
	WalletLogTypeWithdrawFreeze  int16 = 2  // 提现冻结
	WalletLogTypeWithdrawSuccess int16 = 3  // 提现成功
	WalletLogTypeWithdrawReturn  int16 = 4  // 提现退回
	WalletLogTypeAdjust          int16 = 5  // 调账
	WalletLogTypeDeduction       int16 = 6  // 代扣
	WalletLogTypeCashback        int16 = 7  // 返现（押金/流量费）
	WalletLogTypeActivationReward int16 = 11 // 激活奖励入账
)

// getWalletTypeNameStr 获取钱包类型名称
func getWalletTypeNameStr(walletType int16) string {
	switch walletType {
	case 1:
		return "分润钱包"
	case 2:
		return "服务费钱包"
	case 3:
		return "奖励钱包"
	default:
		return "未知钱包"
	}
}

// getLogTypeName 获取流水类型名称
func getLogTypeName(logType int16) string {
	switch logType {
	case WalletLogTypeProfitIn:
		return "分润入账"
	case WalletLogTypeWithdrawFreeze:
		return "提现冻结"
	case WalletLogTypeWithdrawSuccess:
		return "提现成功"
	case WalletLogTypeWithdrawReturn:
		return "提现退回"
	case WalletLogTypeAdjust:
		return "调账"
	case WalletLogTypeDeduction:
		return "代扣"
	case WalletLogTypeCashback:
		return "返现"
	case WalletLogTypeChargingDeposit:
		return "充值钱包充值"
	case WalletLogTypeChargingRewardOut:
		return "奖励发放"
	case WalletLogTypeChargingRewardIn:
		return "收到奖励"
	case WalletLogTypeActivationReward:
		return "激活奖励"
	default:
		return "未知"
	}
}

// getChannelName 获取通道名称（简化版）
func getChannelName(channelID int64) string {
	channelNames := map[int64]string{
		1: "恒信通",
		2: "拉卡拉",
		3: "乐刷",
		4: "随行付",
		5: "连连支付",
		6: "杉德支付",
		7: "富友支付",
		8: "汇付天下",
	}
	if name, ok := channelNames[channelID]; ok {
		return name
	}
	return "未知通道"
}

// WithdrawRequest 提现请求
type WithdrawRequest struct {
	AgentID   int64 `json:"-"`
	WalletID  int64 `json:"wallet_id" binding:"required"`
	Amount    int64 `json:"amount" binding:"required,min=100"` // 分，最少1元
	CreatedBy int64 `json:"-"`
}

// Withdraw 申请提现
func (s *WalletService) Withdraw(req *WithdrawRequest) error {
	// 检查钱包
	wallet, err := s.walletRepo.FindByID(req.WalletID)
	if err != nil || wallet == nil {
		return errors.New("钱包不存在")
	}

	// 验证归属
	if wallet.AgentID != req.AgentID {
		return errors.New("无权操作该钱包")
	}

	// 检查可用余额
	availableBalance := wallet.Balance - wallet.FrozenAmount
	if availableBalance < req.Amount {
		return errors.New("可用余额不足")
	}

	// 检查提现门槛
	if req.Amount < wallet.WithdrawThreshold {
		return fmt.Errorf("提现金额不能低于%d元", wallet.WithdrawThreshold/100)
	}

	// P0修复：奖励钱包提现需检查上级充值钱包余额
	// 业务规则：奖励钱包的资金来源于上级代理商的充值钱包，提现需上级充值钱包有足够余额
	if wallet.WalletType == models.WalletTypeReward {
		if err := s.checkParentChargingWalletBalance(req.AgentID, req.Amount); err != nil {
			return err
		}
	}

	// 冻结金额
	if err := s.walletRepo.FreezeBalance(req.WalletID, req.Amount); err != nil {
		return fmt.Errorf("冻结金额失败: %w", err)
	}

	// 记录流水
	walletLog := &repository.WalletLog{
		WalletID:      req.WalletID,
		AgentID:       req.AgentID,
		WalletType:    wallet.WalletType,
		LogType:       WalletLogTypeWithdrawFreeze,
		Amount:        -req.Amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance, // 冻结不改变余额
		RefType:       "withdraw",
		Remark:        fmt.Sprintf("提现申请，金额%.2f元", float64(req.Amount)/100),
	}
	s.walletLogRepo.Create(walletLog)

	log.Printf("[WalletService] Withdraw request: agent=%d, wallet=%d, amount=%d", req.AgentID, req.WalletID, req.Amount)

	// TODO: 创建提现记录，等待审核/自动打款

	return nil
}

// AddRewardWalletBalance 添加奖励钱包余额
// 返回钱包流水记录ID
func (s *WalletService) AddRewardWalletBalance(agentID int64, amount int64, remark string) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("金额必须大于0")
	}

	// 获取奖励钱包（wallet_type=3）
	wallet, err := s.walletRepo.FindByAgentAndType(agentID, 0, models.WalletTypeReward)
	if err != nil {
		// 钱包不存在时返回错误，需要先创建钱包
		return 0, fmt.Errorf("奖励钱包不存在，请先创建: %w", err)
	}

	// 更新余额
	balanceBefore := wallet.Balance
	balanceAfter := balanceBefore + amount

	if err := s.walletRepo.UpdateBalance(wallet.ID, amount); err != nil {
		return 0, fmt.Errorf("更新钱包余额失败: %w", err)
	}

	// 创建流水记录
	walletLog := &repository.WalletLog{
		WalletID:      wallet.ID,
		AgentID:       agentID,
		WalletType:    models.WalletTypeReward,
		LogType:       1, // 入账
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		RefType:       "stage_reward",
		Remark:        remark,
		CreatedAt:     time.Now(),
	}

	if err := s.walletLogRepo.Create(walletLog); err != nil {
		log.Printf("[WalletService] 创建钱包流水失败: %v", err)
		// 流水记录失败不影响主流程
	}

	log.Printf("[WalletService] 奖励入账成功: AgentID=%d, Amount=%d, Remark=%s", agentID, amount, remark)
	return walletLog.ID, nil
}

// checkParentChargingWalletBalance 检查上级充值钱包余额是否足够
// 业务规则：奖励钱包提现时，需要确保上级充值钱包余额 >= 提现金额
func (s *WalletService) checkParentChargingWalletBalance(agentID int64, amount int64) error {
	// 获取代理商信息
	agent, err := s.agentRepo.FindByIDFull(agentID)
	if err != nil || agent == nil {
		return errors.New("代理商信息不存在")
	}

	// 检查是否有上级
	if agent.ParentID == 0 {
		return errors.New("顶级代理商无法从奖励钱包提现")
	}

	// 获取上级充值钱包
	parentChargingWallet, err := s.walletRepo.FindByAgentAndType(agent.ParentID, 0, models.WalletTypeCharging)
	if err != nil || parentChargingWallet == nil {
		return errors.New("上级充值钱包不存在，无法提现")
	}

	// 检查上级充值钱包可用余额
	parentAvailable := parentChargingWallet.Balance - parentChargingWallet.FrozenAmount
	if parentAvailable < amount {
		return fmt.Errorf("上级充值钱包余额不足，无法提现。上级可用余额：%.2f元，提现金额：%.2f元",
			float64(parentAvailable)/100, float64(amount)/100)
	}

	return nil
}
