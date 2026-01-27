package models

import (
	"time"
)

// DeductionPlan 代扣计划（统一代扣管理）
// 合并货款代扣和伙伴代扣，统一使用定时扣款+冻结机制
// 业务规则:
//   - 接收确认: 下级需确认后代扣才生效并开始冻结
//   - 冻结时机: 接收确认后开始冻结现有余额，后续入账时继续冻结
//   - 扣款频率: 每天一次（8:00执行）
//   - 扣款优先级: 按创建时间先后（FIFO）
//   - 冻结上限: 冻结金额 ≤ 剩余待扣金额
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
	Status          int16      `json:"status" gorm:"default:0"`            // 0:待接收 1:进行中 2:已完成 3:已暂停 4:已取消 5:已拒绝
	RelatedType     string     `json:"related_type" gorm:"size:32"`        // 关联类型: terminal_distribute, partner_loan
	RelatedID       int64      `json:"related_id"`                         // 关联ID
	Remark          string     `json:"remark" gorm:"size:255"`             // 备注
	CreatedBy       int64      `json:"created_by"`                         // 创建人
	CreatedAt       time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:now()"`
	CompletedAt     *time.Time `json:"completed_at"`

	// 新增：接收确认相关
	NeedAccept bool       `json:"need_accept" gorm:"default:false"` // 是否需要接收确认
	AcceptedAt *time.Time `json:"accepted_at"`                      // 接收时间

	// 新增：冻结金额跟踪
	FrozenAmount int64 `json:"frozen_amount" gorm:"default:0"` // 已冻结金额（分）

	// 新增：扣款来源配置
	DeductionSource int16 `json:"deduction_source" gorm:"default:3"` // 1=分润 2=服务费 3=两者

	// 关联数据（非数据库字段）
	DeductorName string             `json:"deductor_name" gorm:"-"` // 扣款方名称
	DeducteeName string             `json:"deductee_name" gorm:"-"` // 被扣款方名称
	Records      []*DeductionRecord `json:"records" gorm:"-"`       // 扣款记录列表
}

func (DeductionPlan) TableName() string {
	return "deduction_plans"
}

// DeductionPlanStatus 代扣计划状态
const (
	DeductionPlanStatusPendingAccept = 0 // 待接收
	DeductionPlanStatusActive        = 1 // 进行中
	DeductionPlanStatusCompleted     = 2 // 已完成
	DeductionPlanStatusPaused        = 3 // 已暂停
	DeductionPlanStatusCancelled     = 4 // 已取消
	DeductionPlanStatusRejected      = 5 // 已拒绝
)

// GetDeductionPlanStatusName 获取代扣计划状态名称
func GetDeductionPlanStatusName(status int16) string {
	switch status {
	case DeductionPlanStatusPendingAccept:
		return "待接收"
	case DeductionPlanStatusActive:
		return "进行中"
	case DeductionPlanStatusCompleted:
		return "已完成"
	case DeductionPlanStatusPaused:
		return "已暂停"
	case DeductionPlanStatusCancelled:
		return "已取消"
	case DeductionPlanStatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}

// DeductionSource 扣款来源
const (
	DeductionSourceProfit     int16 = 1 // 仅分润钱包
	DeductionSourceServiceFee int16 = 2 // 仅服务费钱包
	DeductionSourceBoth       int16 = 3 // 两者都扣（优先分润）
)

// GetDeductionSourceName 获取扣款来源名称
func GetDeductionSourceName(source int16) string {
	switch source {
	case DeductionSourceProfit:
		return "分润钱包"
	case DeductionSourceServiceFee:
		return "服务费钱包"
	case DeductionSourceBoth:
		return "分润+服务费"
	default:
		return "未知"
	}
}

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

// DeductionFreezeLog 代扣冻结明细
// 记录每次冻结操作的详情
type DeductionFreezeLog struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	PlanID       int64     `json:"plan_id" gorm:"not null;index"`   // 代扣计划ID
	AgentID      int64     `json:"agent_id" gorm:"not null;index"`  // 被扣款方代理商ID
	WalletID     int64     `json:"wallet_id" gorm:"not null"`       // 钱包ID
	WalletType   int16     `json:"wallet_type" gorm:"not null"`     // 钱包类型
	ChannelID    int64     `json:"channel_id"`                      // 通道ID
	FreezeAmount int64     `json:"freeze_amount" gorm:"not null"`   // 本次冻结金额（分）
	TotalFrozen  int64     `json:"total_frozen" gorm:"not null"`    // 累计冻结金额（分）
	TriggerType  string    `json:"trigger_type" gorm:"size:32"`     // 触发类型: accept/income
	TriggerRefID int64     `json:"trigger_ref_id"`                  // 触发来源ID
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`
}

func (DeductionFreezeLog) TableName() string {
	return "deduction_freeze_logs"
}

// DeductionFreezeTriggerType 冻结触发类型
const (
	DeductionFreezeTriggerTypeAccept = "accept" // 接收确认时冻结
	DeductionFreezeTriggerTypeIncome = "income" // 入账时冻结
)

// DeductionPlanListResponse 代扣计划列表响应
type DeductionPlanListResponse struct {
	ID              int64      `json:"id"`
	PlanNo          string     `json:"plan_no"`
	DeductorID      int64      `json:"deductor_id"`
	DeductorName    string     `json:"deductor_name"`
	DeducteeID      int64      `json:"deductee_id"`
	DeducteeName    string     `json:"deductee_name"`
	PlanType        int16      `json:"plan_type"`
	PlanTypeName    string     `json:"plan_type_name"`
	TotalAmount     int64      `json:"total_amount"`
	TotalAmountYuan float64    `json:"total_amount_yuan"`
	DeductedAmount  int64      `json:"deducted_amount"`
	RemainingAmount int64      `json:"remaining_amount"`
	FrozenAmount    int64      `json:"frozen_amount"`
	DeductionSource int16      `json:"deduction_source"`
	SourceName      string     `json:"source_name"`
	TotalPeriods    int        `json:"total_periods"`
	CurrentPeriod   int        `json:"current_period"`
	Status          int16      `json:"status"`
	StatusName      string     `json:"status_name"`
	Progress        float64    `json:"progress"`
	NeedAccept      bool       `json:"need_accept"`
	CreatedAt       time.Time  `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

// ToListResponse 转换为列表响应
func (p *DeductionPlan) ToListResponse() *DeductionPlanListResponse {
	progress := float64(0)
	if p.TotalAmount > 0 {
		progress = float64(p.DeductedAmount) / float64(p.TotalAmount) * 100
	}

	planTypeName := "未知"
	switch p.PlanType {
	case DeductionPlanTypeGoods:
		planTypeName = "货款代扣"
	case DeductionPlanTypePartner:
		planTypeName = "伙伴代扣"
	case DeductionPlanTypeDeposit:
		planTypeName = "押金代扣"
	}

	return &DeductionPlanListResponse{
		ID:              p.ID,
		PlanNo:          p.PlanNo,
		DeductorID:      p.DeductorID,
		DeductorName:    p.DeductorName,
		DeducteeID:      p.DeducteeID,
		DeducteeName:    p.DeducteeName,
		PlanType:        p.PlanType,
		PlanTypeName:    planTypeName,
		TotalAmount:     p.TotalAmount,
		TotalAmountYuan: float64(p.TotalAmount) / 100,
		DeductedAmount:  p.DeductedAmount,
		RemainingAmount: p.RemainingAmount,
		FrozenAmount:    p.FrozenAmount,
		DeductionSource: p.DeductionSource,
		SourceName:      GetDeductionSourceName(p.DeductionSource),
		TotalPeriods:    p.TotalPeriods,
		CurrentPeriod:   p.CurrentPeriod,
		Status:          p.Status,
		StatusName:      GetDeductionPlanStatusName(p.Status),
		Progress:        progress,
		NeedAccept:      p.NeedAccept,
		CreatedAt:       p.CreatedAt,
		AcceptedAt:      p.AcceptedAt,
		CompletedAt:     p.CompletedAt,
	}
}

// CreateDeductionPlanWithAcceptRequest 创建需要接收确认的代扣计划请求
type CreateDeductionPlanWithAcceptRequest struct {
	DeductorID      int64  `json:"deductor_id" binding:"required"`      // 扣款方代理商ID
	DeducteeID      int64  `json:"deductee_id" binding:"required"`      // 被扣款方代理商ID
	PlanType        int16  `json:"plan_type" binding:"required"`        // 计划类型
	TotalAmount     int64  `json:"total_amount" binding:"required,gt=0"` // 总金额（分）
	TotalPeriods    int    `json:"total_periods" binding:"required,gt=0"` // 总期数
	DeductionSource int16  `json:"deduction_source" binding:"required"` // 扣款来源
	Remark          string `json:"remark"`                              // 备注
}

// DeductionSummary 代扣统计汇总
type DeductionSummary struct {
	TotalCount        int64 `json:"total_count"`         // 总代扣数
	PendingCount      int64 `json:"pending_count"`       // 待接收数
	InProgressCount   int64 `json:"in_progress_count"`   // 进行中数
	CompletedCount    int64 `json:"completed_count"`     // 已完成数
	TotalAmount       int64 `json:"total_amount"`        // 代扣总金额（分）
	DeductedAmount    int64 `json:"deducted_amount"`     // 已扣总金额（分）
	RemainingAmount   int64 `json:"remaining_amount"`    // 剩余待扣总金额（分）
	TotalFrozenAmount int64 `json:"total_frozen_amount"` // 总冻结金额（分）
}
