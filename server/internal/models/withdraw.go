package models

import (
	"time"
)

// 提现状态
const (
	WithdrawStatusPending   int16 = 0 // 待审核
	WithdrawStatusApproved  int16 = 1 // 已审核（待打款）
	WithdrawStatusPaid      int16 = 2 // 已打款
	WithdrawStatusRejected  int16 = 3 // 已拒绝
	WithdrawStatusFailed    int16 = 4 // 打款失败
	WithdrawStatusCancelled int16 = 5 // 已取消
)

// GetWithdrawStatusName 获取提现状态名称
func GetWithdrawStatusName(status int16) string {
	switch status {
	case WithdrawStatusPending:
		return "待审核"
	case WithdrawStatusApproved:
		return "已审核"
	case WithdrawStatusPaid:
		return "已打款"
	case WithdrawStatusRejected:
		return "已拒绝"
	case WithdrawStatusFailed:
		return "打款失败"
	case WithdrawStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// WithdrawRecord 提现记录
type WithdrawRecord struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	WithdrawNo   string `json:"withdraw_no" gorm:"size:50;uniqueIndex"` // 提现单号
	AgentID      int64  `json:"agent_id" gorm:"index"`                  // 代理商ID
	WalletID     int64  `json:"wallet_id" gorm:"index"`                 // 钱包ID
	WalletType   int16  `json:"wallet_type"`                            // 钱包类型
	ChannelID    int64  `json:"channel_id"`                             // 支付通道ID
	TaxChannelID *int64 `json:"tax_channel_id"`                         // 税筹通道ID

	// 金额信息
	Amount        int64 `json:"amount"`         // 提现金额(分)
	TaxFee        int64 `json:"tax_fee"`        // 税费(分)
	FixedFee      int64 `json:"fixed_fee"`      // 固定手续费(分)
	ActualAmount  int64 `json:"actual_amount"`  // 实际到账金额(分)

	// 结算卡信息（冗余存储，防止修改后影响历史记录）
	BankName    string `json:"bank_name" gorm:"size:100"`    // 开户银行
	BankAccount string `json:"bank_account" gorm:"size:50"`  // 银行卡号（加密存储）
	AccountName string `json:"account_name" gorm:"size:50"`  // 开户名

	// 状态信息
	Status       int16  `json:"status" gorm:"default:0"`
	RejectReason string `json:"reject_reason" gorm:"size:500"` // 拒绝原因
	FailReason   string `json:"fail_reason" gorm:"size:500"`   // 失败原因

	// 审核信息
	AuditedBy   *int64     `json:"audited_by"`
	AuditedAt   *time.Time `json:"audited_at"`
	AuditRemark string     `json:"audit_remark" gorm:"size:500"`

	// 打款信息
	PaidAt     *time.Time `json:"paid_at"`
	PaidRef    string     `json:"paid_ref" gorm:"size:100"` // 打款流水号
	PaidRemark string     `json:"paid_remark" gorm:"size:500"`

	// 时间戳
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (WithdrawRecord) TableName() string {
	return "withdraw_records"
}

// WithdrawStats 提现统计
type WithdrawStats struct {
	TotalCount     int64 `json:"total_count"`      // 总提现次数
	TotalAmount    int64 `json:"total_amount"`     // 总提现金额
	PendingCount   int64 `json:"pending_count"`    // 待审核数量
	PendingAmount  int64 `json:"pending_amount"`   // 待审核金额
	PaidCount      int64 `json:"paid_count"`       // 已打款数量
	PaidAmount     int64 `json:"paid_amount"`      // 已打款金额
	RejectedCount  int64 `json:"rejected_count"`   // 已拒绝数量
	RejectedAmount int64 `json:"rejected_amount"`  // 已拒绝金额
}
