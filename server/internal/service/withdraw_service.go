package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/pkg/crypto"
)

// WithdrawService 提现服务
type WithdrawService struct {
	withdrawRepo    *repository.GormWithdrawRepository
	walletRepo      *repository.GormWalletRepository
	walletLogRepo   *repository.GormWalletLogRepository
	agentRepo       *repository.GormAgentRepository
	taxChannelRepo  *repository.GormTaxChannelRepository
}

// NewWithdrawService 创建提现服务
func NewWithdrawService(
	withdrawRepo *repository.GormWithdrawRepository,
	walletRepo *repository.GormWalletRepository,
	walletLogRepo *repository.GormWalletLogRepository,
	agentRepo *repository.GormAgentRepository,
	taxChannelRepo *repository.GormTaxChannelRepository,
) *WithdrawService {
	return &WithdrawService{
		withdrawRepo:   withdrawRepo,
		walletRepo:     walletRepo,
		walletLogRepo:  walletLogRepo,
		agentRepo:      agentRepo,
		taxChannelRepo: taxChannelRepo,
	}
}

// CreateWithdrawRequest 创建提现请求
type CreateWithdrawRequest struct {
	AgentID  int64 `json:"-"`
	WalletID int64 `json:"wallet_id" binding:"required"`
	Amount   int64 `json:"amount" binding:"required,min=100"` // 分，最少1元
}

// WithdrawDetailResponse 提现详情响应
type WithdrawDetailResponse struct {
	ID            int64      `json:"id"`
	WithdrawNo    string     `json:"withdraw_no"`
	WalletType    int16      `json:"wallet_type"`
	WalletTypeName string    `json:"wallet_type_name"`
	Amount        int64      `json:"amount"`
	AmountYuan    float64    `json:"amount_yuan"`
	TaxFee        int64      `json:"tax_fee"`
	TaxFeeYuan    float64    `json:"tax_fee_yuan"`
	FixedFee      int64      `json:"fixed_fee"`
	FixedFeeYuan  float64    `json:"fixed_fee_yuan"`
	ActualAmount  int64      `json:"actual_amount"`
	ActualYuan    float64    `json:"actual_yuan"`
	BankName      string     `json:"bank_name"`
	BankAccount   string     `json:"bank_account"` // 脱敏显示
	AccountName   string     `json:"account_name"`
	Status        int16      `json:"status"`
	StatusName    string     `json:"status_name"`
	RejectReason  string     `json:"reject_reason,omitempty"`
	FailReason    string     `json:"fail_reason,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	PaidRef       string     `json:"paid_ref,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// CreateWithdraw 创建提现申请
func (s *WithdrawService) CreateWithdraw(req *CreateWithdrawRequest) (*models.WithdrawRecord, error) {
	// 1. 检查钱包
	wallet, err := s.walletRepo.FindByID(req.WalletID)
	if err != nil || wallet == nil {
		return nil, errors.New("钱包不存在")
	}

	// 2. 验证归属
	if wallet.AgentID != req.AgentID {
		return nil, errors.New("无权操作该钱包")
	}

	// 3. 检查可用余额
	availableBalance := wallet.Balance - wallet.FrozenAmount
	if availableBalance < req.Amount {
		return nil, errors.New("可用余额不足")
	}

	// 4. 检查提现门槛
	if req.Amount < wallet.WithdrawThreshold {
		return nil, fmt.Errorf("提现金额不能低于%d元", wallet.WithdrawThreshold/100)
	}

	// 5. 奖励钱包提现需检查上级充值钱包余额
	if wallet.WalletType == models.WalletTypeReward {
		if err := s.checkParentChargingWalletBalance(req.AgentID, req.Amount); err != nil {
			return nil, err
		}
	}

	// 6. 获取代理商结算卡信息
	agent, err := s.agentRepo.FindByIDFull(req.AgentID)
	if err != nil || agent == nil {
		return nil, errors.New("代理商信息不存在")
	}

	if agent.BankAccount == "" {
		return nil, errors.New("请先设置结算卡信息")
	}

	// 7. 获取税筹通道配置
	var taxChannelID *int64
	var taxFee, fixedFee int64
	taxChannel, err := s.taxChannelRepo.GetTaxChannelForWithdrawal(wallet.ChannelID, wallet.WalletType)
	if err == nil && taxChannel != nil {
		taxChannelID = &taxChannel.ID
		// 计算税费和固定费用
		taxFee = int64(float64(req.Amount) * taxChannel.TaxRate)
		fixedFee = taxChannel.FixedFee
	}

	// 8. 计算实际到账金额
	actualAmount := req.Amount - taxFee - fixedFee
	if actualAmount <= 0 {
		return nil, errors.New("提现金额太小，扣除费用后不足")
	}

	// 9. 冻结金额
	if err := s.walletRepo.FreezeBalance(req.WalletID, req.Amount); err != nil {
		return nil, fmt.Errorf("冻结金额失败: %w", err)
	}

	// 10. 加密银行卡号
	encryptedAccount, _ := crypto.EncryptPhone(agent.BankAccount)

	// 11. 创建提现记录
	record := &models.WithdrawRecord{
		WithdrawNo:   repository.GenerateWithdrawNo(),
		AgentID:      req.AgentID,
		WalletID:     req.WalletID,
		WalletType:   wallet.WalletType,
		ChannelID:    wallet.ChannelID,
		TaxChannelID: taxChannelID,
		Amount:       req.Amount,
		TaxFee:       taxFee,
		FixedFee:     fixedFee,
		ActualAmount: actualAmount,
		BankName:     agent.BankName,
		BankAccount:  encryptedAccount,
		AccountName:  agent.ContactName,
		Status:       models.WithdrawStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.withdrawRepo.Create(record); err != nil {
		// 回滚冻结
		s.walletRepo.UnfreezeBalance(req.WalletID, req.Amount)
		return nil, fmt.Errorf("创建提现记录失败: %w", err)
	}

	// 12. 记录钱包流水
	walletLog := &repository.WalletLog{
		WalletID:      req.WalletID,
		AgentID:       req.AgentID,
		WalletType:    wallet.WalletType,
		LogType:       WalletLogTypeWithdrawFreeze,
		Amount:        -req.Amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance,
		RefType:       "withdraw",
		RefID:         record.ID,
		Remark:        fmt.Sprintf("提现申请，金额%.2f元，单号%s", float64(req.Amount)/100, record.WithdrawNo),
		CreatedAt:     time.Now(),
	}
	s.walletLogRepo.Create(walletLog)

	log.Printf("[WithdrawService] Created withdraw: agent=%d, wallet=%d, amount=%d, no=%s",
		req.AgentID, req.WalletID, req.Amount, record.WithdrawNo)

	return record, nil
}

// GetWithdrawList 获取提现记录列表
func (s *WithdrawService) GetWithdrawList(agentID int64, status *int16, page, pageSize int) ([]*WithdrawDetailResponse, int64, error) {
	offset := (page - 1) * pageSize
	records, total, err := s.withdrawRepo.FindByAgentID(agentID, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*WithdrawDetailResponse, 0, len(records))
	for _, r := range records {
		list = append(list, s.toDetailResponse(r))
	}

	return list, total, nil
}

// GetWithdrawDetail 获取提现详情
func (s *WithdrawService) GetWithdrawDetail(agentID, withdrawID int64) (*WithdrawDetailResponse, error) {
	record, err := s.withdrawRepo.FindByID(withdrawID)
	if err != nil || record == nil {
		return nil, errors.New("提现记录不存在")
	}

	if record.AgentID != agentID {
		return nil, errors.New("无权查看此记录")
	}

	return s.toDetailResponse(record), nil
}

// GetWithdrawStats 获取提现统计
func (s *WithdrawService) GetWithdrawStats(agentID int64) (*models.WithdrawStats, error) {
	return s.withdrawRepo.GetStatsByAgent(agentID)
}

// CancelWithdraw 取消提现
func (s *WithdrawService) CancelWithdraw(agentID, withdrawID int64) error {
	record, err := s.withdrawRepo.FindByID(withdrawID)
	if err != nil || record == nil {
		return errors.New("提现记录不存在")
	}

	if record.AgentID != agentID {
		return errors.New("无权操作此记录")
	}

	if record.Status != models.WithdrawStatusPending {
		return errors.New("只能取消待审核的提现")
	}

	// 解冻金额
	if err := s.walletRepo.UnfreezeBalance(record.WalletID, record.Amount); err != nil {
		return fmt.Errorf("解冻金额失败: %w", err)
	}

	// 更新状态
	record.Status = models.WithdrawStatusCancelled
	record.UpdatedAt = time.Now()
	if err := s.withdrawRepo.Update(record); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 记录流水
	wallet, _ := s.walletRepo.FindByID(record.WalletID)
	if wallet != nil {
		walletLog := &repository.WalletLog{
			WalletID:      record.WalletID,
			AgentID:       agentID,
			WalletType:    record.WalletType,
			LogType:       WalletLogTypeWithdrawReturn,
			Amount:        record.Amount,
			BalanceBefore: wallet.Balance - record.Amount,
			BalanceAfter:  wallet.Balance,
			RefType:       "withdraw_cancel",
			RefID:         record.ID,
			Remark:        fmt.Sprintf("取消提现，金额%.2f元，单号%s", float64(record.Amount)/100, record.WithdrawNo),
			CreatedAt:     time.Now(),
		}
		s.walletLogRepo.Create(walletLog)
	}

	log.Printf("[WithdrawService] Cancelled withdraw: id=%d, no=%s", withdrawID, record.WithdrawNo)

	return nil
}

// ApproveWithdraw 审核通过提现（管理员）
func (s *WithdrawService) ApproveWithdraw(withdrawID, operatorID int64, remark string) error {
	record, err := s.withdrawRepo.FindByID(withdrawID)
	if err != nil || record == nil {
		return errors.New("提现记录不存在")
	}

	if record.Status != models.WithdrawStatusPending {
		return errors.New("只能审核待审核的提现")
	}

	now := time.Now()
	record.Status = models.WithdrawStatusApproved
	record.AuditedBy = &operatorID
	record.AuditedAt = &now
	record.AuditRemark = remark
	record.UpdatedAt = now

	if err := s.withdrawRepo.Update(record); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	log.Printf("[WithdrawService] Approved withdraw: id=%d, no=%s, by=%d", withdrawID, record.WithdrawNo, operatorID)

	// TODO: 调用税筹通道API自动打款
	// s.processPayment(record)

	return nil
}

// RejectWithdraw 拒绝提现（管理员）
func (s *WithdrawService) RejectWithdraw(withdrawID, operatorID int64, reason string) error {
	record, err := s.withdrawRepo.FindByID(withdrawID)
	if err != nil || record == nil {
		return errors.New("提现记录不存在")
	}

	if record.Status != models.WithdrawStatusPending {
		return errors.New("只能审核待审核的提现")
	}

	// 解冻金额
	if err := s.walletRepo.UnfreezeBalance(record.WalletID, record.Amount); err != nil {
		return fmt.Errorf("解冻金额失败: %w", err)
	}

	now := time.Now()
	record.Status = models.WithdrawStatusRejected
	record.AuditedBy = &operatorID
	record.AuditedAt = &now
	record.RejectReason = reason
	record.UpdatedAt = now

	if err := s.withdrawRepo.Update(record); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 记录流水
	wallet, _ := s.walletRepo.FindByID(record.WalletID)
	if wallet != nil {
		walletLog := &repository.WalletLog{
			WalletID:      record.WalletID,
			AgentID:       record.AgentID,
			WalletType:    record.WalletType,
			LogType:       WalletLogTypeWithdrawReturn,
			Amount:        record.Amount,
			BalanceBefore: wallet.Balance - record.Amount,
			BalanceAfter:  wallet.Balance,
			RefType:       "withdraw_reject",
			RefID:         record.ID,
			Remark:        fmt.Sprintf("提现被拒绝，金额%.2f元，原因：%s", float64(record.Amount)/100, reason),
			CreatedAt:     time.Now(),
		}
		s.walletLogRepo.Create(walletLog)
	}

	log.Printf("[WithdrawService] Rejected withdraw: id=%d, no=%s, reason=%s", withdrawID, record.WithdrawNo, reason)

	return nil
}

// ConfirmPaid 确认打款（管理员）
func (s *WithdrawService) ConfirmPaid(withdrawID, operatorID int64, paidRef, remark string) error {
	record, err := s.withdrawRepo.FindByID(withdrawID)
	if err != nil || record == nil {
		return errors.New("提现记录不存在")
	}

	if record.Status != models.WithdrawStatusApproved {
		return errors.New("只能确认已审核的提现")
	}

	// 扣减钱包余额（实际出款）
	if err := s.walletRepo.DeductFrozenBalance(record.WalletID, record.Amount); err != nil {
		return fmt.Errorf("扣减余额失败: %w", err)
	}

	now := time.Now()
	record.Status = models.WithdrawStatusPaid
	record.PaidAt = &now
	record.PaidRef = paidRef
	record.PaidRemark = remark
	record.UpdatedAt = now

	if err := s.withdrawRepo.Update(record); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 记录流水
	wallet, _ := s.walletRepo.FindByID(record.WalletID)
	balanceAfter := int64(0)
	if wallet != nil {
		balanceAfter = wallet.Balance
	}
	walletLog := &repository.WalletLog{
		WalletID:      record.WalletID,
		AgentID:       record.AgentID,
		WalletType:    record.WalletType,
		LogType:       WalletLogTypeWithdrawSuccess,
		Amount:        -record.Amount,
		BalanceBefore: balanceAfter + record.Amount,
		BalanceAfter:  balanceAfter,
		RefType:       "withdraw_paid",
		RefID:         record.ID,
		Remark:        fmt.Sprintf("提现成功，金额%.2f元，实际到账%.2f元", float64(record.Amount)/100, float64(record.ActualAmount)/100),
		CreatedAt:     time.Now(),
	}
	s.walletLogRepo.Create(walletLog)

	log.Printf("[WithdrawService] Paid withdraw: id=%d, no=%s, ref=%s", withdrawID, record.WithdrawNo, paidRef)

	return nil
}

// GetPendingList 获取待审核列表（管理员）
func (s *WithdrawService) GetPendingList(page, pageSize int) ([]*WithdrawDetailResponse, int64, error) {
	offset := (page - 1) * pageSize
	records, total, err := s.withdrawRepo.FindPending(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*WithdrawDetailResponse, 0, len(records))
	for _, r := range records {
		list = append(list, s.toDetailResponse(r))
	}

	return list, total, nil
}

// checkParentChargingWalletBalance 检查上级充值钱包余额
func (s *WithdrawService) checkParentChargingWalletBalance(agentID int64, amount int64) error {
	agent, err := s.agentRepo.FindByIDFull(agentID)
	if err != nil || agent == nil {
		return errors.New("代理商信息不存在")
	}

	if agent.ParentID == 0 {
		return errors.New("顶级代理商无法从奖励钱包提现")
	}

	parentChargingWallet, err := s.walletRepo.FindByAgentAndType(agent.ParentID, 0, models.WalletTypeCharging)
	if err != nil || parentChargingWallet == nil {
		return errors.New("上级充值钱包不存在，无法提现")
	}

	parentAvailable := parentChargingWallet.Balance - parentChargingWallet.FrozenAmount
	if parentAvailable < amount {
		return fmt.Errorf("上级充值钱包余额不足，无法提现。上级可用余额：%.2f元，提现金额：%.2f元",
			float64(parentAvailable)/100, float64(amount)/100)
	}

	return nil
}

// toDetailResponse 转换为详情响应
func (s *WithdrawService) toDetailResponse(r *models.WithdrawRecord) *WithdrawDetailResponse {
	// 脱敏银行卡号
	maskedAccount := r.BankAccount
	if crypto.IsEncrypted(r.BankAccount) {
		decrypted, err := crypto.DecryptPhone(r.BankAccount)
		if err == nil {
			maskedAccount = maskBankAccount(decrypted)
		}
	} else if len(r.BankAccount) > 8 {
		maskedAccount = maskBankAccount(r.BankAccount)
	}

	return &WithdrawDetailResponse{
		ID:             r.ID,
		WithdrawNo:     r.WithdrawNo,
		WalletType:     r.WalletType,
		WalletTypeName: models.WalletTypeName(r.WalletType),
		Amount:         r.Amount,
		AmountYuan:     float64(r.Amount) / 100,
		TaxFee:         r.TaxFee,
		TaxFeeYuan:     float64(r.TaxFee) / 100,
		FixedFee:       r.FixedFee,
		FixedFeeYuan:   float64(r.FixedFee) / 100,
		ActualAmount:   r.ActualAmount,
		ActualYuan:     float64(r.ActualAmount) / 100,
		BankName:       r.BankName,
		BankAccount:    maskedAccount,
		AccountName:    r.AccountName,
		Status:         r.Status,
		StatusName:     models.GetWithdrawStatusName(r.Status),
		RejectReason:   r.RejectReason,
		FailReason:     r.FailReason,
		PaidAt:         r.PaidAt,
		PaidRef:        r.PaidRef,
		CreatedAt:      r.CreatedAt,
	}
}

// maskBankAccount 银行卡号脱敏
func maskBankAccount(account string) string {
	if len(account) <= 8 {
		return account
	}
	return account[:4] + "****" + account[len(account)-4:]
}
