package models

import (
	"time"
)

// 调账状态
const (
	AdjustmentStatusEffective int16 = 1 // 已生效
	AdjustmentStatusPending   int16 = 2 // 待审批（预留）
	AdjustmentStatusRejected  int16 = 3 // 已驳回（预留）
)

// WalletAdjustment 钱包调账记录
type WalletAdjustment struct {
	ID            int64  `json:"id" gorm:"primaryKey"`
	AdjustmentNo  string `json:"adjustment_no" gorm:"size:50;uniqueIndex"` // 调账单号
	AgentID       int64  `json:"agent_id" gorm:"index"`                    // 代理商ID
	WalletID      int64  `json:"wallet_id" gorm:"index"`                   // 钱包ID
	WalletType    int16  `json:"wallet_type"`                              // 钱包类型
	ChannelID     int64  `json:"channel_id" gorm:"default:0"`              // 通道ID
	Amount        int64  `json:"amount"`                                   // 调账金额(分)
	BalanceBefore int64  `json:"balance_before"`                           // 调账前余额
	BalanceAfter  int64  `json:"balance_after"`                            // 调账后余额
	Reason        string `json:"reason" gorm:"size:500"`                   // 调账原因
	OperatorID    int64  `json:"operator_id"`                              // 操作人ID
	OperatorName  string `json:"operator_name" gorm:"size:50"`             // 操作人名称

	// 状态
	Status       int16      `json:"status" gorm:"default:1"`      // 1已生效 2待审批 3已驳回
	ApprovedBy   *int64     `json:"approved_by"`                  // 审批人ID
	ApprovedAt   *time.Time `json:"approved_at"`                  // 审批时间
	RejectReason string     `json:"reject_reason" gorm:"size:500"` // 驳回原因

	// 关联
	WalletLogID *int64 `json:"wallet_log_id"` // 关联的钱包流水ID

	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (WalletAdjustment) TableName() string {
	return "wallet_adjustments"
}

// GetAdjustmentStatusName 获取调账状态名称
func GetAdjustmentStatusName(status int16) string {
	switch status {
	case AdjustmentStatusEffective:
		return "已生效"
	case AdjustmentStatusPending:
		return "待审批"
	case AdjustmentStatusRejected:
		return "已驳回"
	default:
		return "未知"
	}
}

// GetAdjustmentTypeName 获取调账类型名称（充入/扣减）
func GetAdjustmentTypeName(amount int64) string {
	if amount >= 0 {
		return "充入"
	}
	return "扣减"
}
