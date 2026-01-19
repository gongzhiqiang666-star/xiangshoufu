package models

import (
	"time"
)

// DeductionPlan 代扣计划
// 支持伙伴代扣（任意代理商之间，不限层级关系）
type DeductionPlan struct {
	ID              int64      `json:"id" gorm:"primaryKey"`
	PlanNo          string     `json:"plan_no" gorm:"size:64;uniqueIndex"` // 计划编号
	DeductorID      int64      `json:"deductor_id" gorm:"not null;index"`  // 扣款方代理商ID
	DeducteeID      int64      `json:"deductee_id" gorm:"not null;index"`  // 被扣款方代理商ID
	PlanType        int16      `json:"plan_type" gorm:"not null"`          // 1:货款代扣 2:伙伴代扣 3:押金代扣
	TotalAmount     int64      `json:"total_amount" gorm:"not null"`       // 总金额（分）
	DeductedAmount  int64      `json:"deducted_amount" gorm:"default:0"`   // 已扣金额（分）
	RemainingAmount int64      `json:"remaining_amount" gorm:"not null"`   // 剩余金额（分）
	TotalPeriods    int        `json:"total_periods" gorm:"not null"`      // 总期数
	CurrentPeriod   int        `json:"current_period" gorm:"default:0"`    // 当前期数
	PeriodAmount    int64      `json:"period_amount" gorm:"not null"`      // 每期金额（分）
	Status          int16      `json:"status" gorm:"default:1"`            // 1:进行中 2:已完成 3:已暂停 4:已取消
	RelatedType     string     `json:"related_type" gorm:"size:32"`        // 关联类型: terminal_distribute, partner_loan
	RelatedID       int64      `json:"related_id"`                         // 关联ID
	Remark          string     `json:"remark" gorm:"size:255"`             // 备注
	CreatedBy       int64      `json:"created_by"`                         // 创建人
	CreatedAt       time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:now()"`
	CompletedAt     *time.Time `json:"completed_at"`
}

func (DeductionPlan) TableName() string {
	return "deduction_plans"
}

// DeductionPlanStatus 代扣计划状态
const (
	DeductionPlanStatusActive    = 1 // 进行中
	DeductionPlanStatusCompleted = 2 // 已完成
	DeductionPlanStatusPaused    = 3 // 已暂停
	DeductionPlanStatusCancelled = 4 // 已取消
)

// DeductionPlanType 代扣计划类型
const (
	DeductionPlanTypeGoods   = 1 // 货款代扣
	DeductionPlanTypePartner = 2 // 伙伴代扣
	DeductionPlanTypeDeposit = 3 // 押金代扣
)

// DeductionRecord 代扣记录
type DeductionRecord struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	PlanID        int64      `json:"plan_id" gorm:"not null;index"`      // 代扣计划ID
	PlanNo        string     `json:"plan_no" gorm:"size:64;index"`       // 计划编号
	DeductorID    int64      `json:"deductor_id" gorm:"not null"`        // 扣款方
	DeducteeID    int64      `json:"deductee_id" gorm:"not null;index"`  // 被扣款方
	PeriodNum     int        `json:"period_num" gorm:"not null"`         // 期数
	Amount        int64      `json:"amount" gorm:"not null"`             // 应扣金额（分）
	ActualAmount  int64      `json:"actual_amount" gorm:"default:0"`     // 实扣金额（分）
	Status        int16      `json:"status" gorm:"default:0"`            // 0:待扣款 1:成功 2:部分成功 3:失败
	WalletDetails string     `json:"wallet_details" gorm:"type:jsonb"`   // 钱包扣款明细JSON
	FailReason    string     `json:"fail_reason" gorm:"size:255"`        // 失败原因
	ScheduledAt   time.Time  `json:"scheduled_at" gorm:"not null;index"` // 计划扣款时间
	DeductedAt    *time.Time `json:"deducted_at"`                        // 实际扣款时间
	CreatedAt     time.Time  `json:"created_at" gorm:"default:now()"`
}

func (DeductionRecord) TableName() string {
	return "deduction_records"
}

// DeductionRecordStatus 代扣记录状态
const (
	DeductionRecordStatusPending        = 0 // 待扣款
	DeductionRecordStatusSuccess        = 1 // 成功
	DeductionRecordStatusPartialSuccess = 2 // 部分成功
	DeductionRecordStatusFailed         = 3 // 失败
)

// WalletDeductDetail 钱包扣款明细
type WalletDeductDetail struct {
	WalletID      int64  `json:"wallet_id"`
	WalletType    int16  `json:"wallet_type"`
	WalletName    string `json:"wallet_name"`
	BalanceBefore int64  `json:"balance_before"`
	DeductAmount  int64  `json:"deduct_amount"`
	BalanceAfter  int64  `json:"balance_after"`
}

// DeductionChain 代扣链（用于跨级下发时生成多级代扣）
type DeductionChain struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	ChainNo      string    `json:"chain_no" gorm:"size:64;uniqueIndex"` // 代扣链编号
	DistributeID int64     `json:"distribute_id" gorm:"not null;index"` // 终端下发记录ID
	TerminalSN   string    `json:"terminal_sn" gorm:"size:50;index"`    // 终端SN
	TotalLevels  int       `json:"total_levels" gorm:"not null"`        // 总层级数
	TotalAmount  int64     `json:"total_amount" gorm:"not null"`        // 总金额
	Status       int16     `json:"status" gorm:"default:1"`             // 1:进行中 2:已完成 3:已取消
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:now()"`
}

func (DeductionChain) TableName() string {
	return "deduction_chains"
}

// DeductionChainItem 代扣链节点（A→B→C 中的每个节点）
type DeductionChainItem struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	ChainID     int64  `json:"chain_id" gorm:"not null;index"`      // 所属代扣链
	ChainNo     string `json:"chain_no" gorm:"size:64;index"`       // 代扣链编号
	Level       int    `json:"level" gorm:"not null"`               // 层级（1,2,3...）
	FromAgentID int64  `json:"from_agent_id" gorm:"not null;index"` // 扣款方
	ToAgentID   int64  `json:"to_agent_id" gorm:"not null;index"`   // 收款方
	PlanID      int64  `json:"plan_id" gorm:"index"`                // 关联的代扣计划ID
	Amount      int64  `json:"amount" gorm:"not null"`              // 代扣金额
	Status      int16  `json:"status" gorm:"default:0"`             // 0:待处理 1:已生成计划 2:已完成
}

func (DeductionChainItem) TableName() string {
	return "deduction_chain_items"
}
