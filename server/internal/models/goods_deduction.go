package models

import (
	"time"
)

// GoodsDeduction 货款代扣
// 终端划拨时设置的货款代扣，与代扣管理模块独立
// 业务规则:
//   1. 扣款规则: 分润优先 - 先扣分润钱包，扣完再扣服务费钱包
//   2. 扣款时机: 实时扣款 - 钱包入账时立即触发扣款
//   3. 部分扣款: 有多少扣多少 - 余额不足时部分扣除，剩余下次继续扣
//   4. 扣款上限: 无上限 - 每次入账时全部扣除，直到扣完为止
//   5. 扣款优先级: 货款代扣 > 上级代扣 > 伙伴代扣
//   6. 与提现关系: 待扣金额占用钱包余额，影响可提现金额
type GoodsDeduction struct {
	ID              int64      `json:"id" gorm:"primaryKey"`
	DeductionNo     string     `json:"deduction_no" gorm:"size:64;uniqueIndex"`  // 代扣编号
	FromAgentID     int64      `json:"from_agent_id" gorm:"not null;index"`      // 上级代理商ID（扣款方/发起方）
	ToAgentID       int64      `json:"to_agent_id" gorm:"not null;index"`        // 下级代理商ID（被扣款方/接收方）
	TotalAmount     int64      `json:"total_amount" gorm:"not null"`             // 代扣总金额（分）
	DeductedAmount  int64      `json:"deducted_amount" gorm:"default:0"`         // 已扣金额（分）
	RemainingAmount int64      `json:"remaining_amount" gorm:"not null"`         // 剩余金额（分）
	DeductionSource int16      `json:"deduction_source" gorm:"not null;default:3"` // 扣款来源: 1=分润 2=服务费 3=两者
	TerminalCount   int        `json:"terminal_count" gorm:"default:0"`          // 终端数量
	UnitPrice       int64      `json:"unit_price" gorm:"default:0"`              // 单价（分）
	Status          int16      `json:"status" gorm:"not null;default:1;index"`   // 状态: 1=待接收 2=进行中 3=已完成 4=已拒绝
	AgreementSigned bool       `json:"agreement_signed" gorm:"default:false"`    // 是否签署协议
	AgreementURL    string     `json:"agreement_url" gorm:"size:500"`            // 协议文件URL
	DistributeID    *int64     `json:"distribute_id" gorm:"index"`               // 关联的终端划拨ID
	Remark          string     `json:"remark" gorm:"size:500"`                   // 备注
	CreatedBy       int64      `json:"created_by"`                               // 创建人
	AcceptedAt      *time.Time `json:"accepted_at"`                              // 接收时间
	CompletedAt     *time.Time `json:"completed_at"`                             // 完成时间
	CreatedAt       time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:now()"`

	// 关联数据（非数据库字段）
	FromAgentName string                     `json:"from_agent_name" gorm:"-"` // 发起方名称
	ToAgentName   string                     `json:"to_agent_name" gorm:"-"`   // 接收方名称
	Terminals     []*GoodsDeductionTerminal  `json:"terminals" gorm:"-"`       // 关联终端列表
	Details       []*GoodsDeductionDetail    `json:"details" gorm:"-"`         // 扣款明细列表
}

func (GoodsDeduction) TableName() string {
	return "goods_deductions"
}

// GoodsDeductionStatus 货款代扣状态
const (
	GoodsDeductionStatusPendingAccept int16 = 1 // 待接收
	GoodsDeductionStatusInProgress    int16 = 2 // 进行中
	GoodsDeductionStatusCompleted     int16 = 3 // 已完成
	GoodsDeductionStatusRejected      int16 = 4 // 已拒绝
)

// GetGoodsDeductionStatusName 获取货款代扣状态名称
func GetGoodsDeductionStatusName(status int16) string {
	switch status {
	case GoodsDeductionStatusPendingAccept:
		return "待接收"
	case GoodsDeductionStatusInProgress:
		return "进行中"
	case GoodsDeductionStatusCompleted:
		return "已完成"
	case GoodsDeductionStatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}

// GoodsDeductionSource 扣款来源
const (
	GoodsDeductionSourceProfit     int16 = 1 // 仅分润钱包
	GoodsDeductionSourceServiceFee int16 = 2 // 仅服务费钱包
	GoodsDeductionSourceBoth       int16 = 3 // 两者都扣（优先分润）
)

// GetGoodsDeductionSourceName 获取扣款来源名称
func GetGoodsDeductionSourceName(source int16) string {
	switch source {
	case GoodsDeductionSourceProfit:
		return "分润钱包"
	case GoodsDeductionSourceServiceFee:
		return "服务费钱包"
	case GoodsDeductionSourceBoth:
		return "分润+服务费"
	default:
		return "未知"
	}
}

// GoodsDeductionDetail 货款代扣明细
// 每次钱包入账触发的实时扣款记录
type GoodsDeductionDetail struct {
	ID                   int64     `json:"id" gorm:"primaryKey"`
	DeductionID          int64     `json:"deduction_id" gorm:"not null;index"`         // 关联货款代扣ID
	DeductionNo          string    `json:"deduction_no" gorm:"size:64"`                // 代扣编号（冗余）
	Amount               int64     `json:"amount" gorm:"not null"`                     // 本次扣款金额（分）
	WalletType           int16     `json:"wallet_type" gorm:"not null"`                // 扣款钱包类型: 1=分润 2=服务费
	ChannelID            *int64    `json:"channel_id"`                                 // 通道ID
	WalletBalanceBefore  int64     `json:"wallet_balance_before" gorm:"not null"`      // 扣款前余额（分）
	WalletBalanceAfter   int64     `json:"wallet_balance_after" gorm:"not null"`       // 扣款后余额（分）
	CumulativeDeducted   int64     `json:"cumulative_deducted" gorm:"not null"`        // 累计已扣金额（分）
	RemainingAfter       int64     `json:"remaining_after" gorm:"not null"`            // 扣款后剩余待扣（分）
	TriggerType          string    `json:"trigger_type" gorm:"size:32"`                // 触发类型
	TriggerTransactionID *int64    `json:"trigger_transaction_id" gorm:"index"`        // 触发扣款的交易ID
	TriggerProfitID      *int64    `json:"trigger_profit_id"`                          // 触发扣款的分润记录ID
	CreatedAt            time.Time `json:"created_at" gorm:"default:now();index"`

	// 关联数据（非数据库字段）
	WalletTypeName string `json:"wallet_type_name" gorm:"-"` // 钱包类型名称
	ChannelName    string `json:"channel_name" gorm:"-"`     // 通道名称
}

func (GoodsDeductionDetail) TableName() string {
	return "goods_deduction_details"
}

// 触发类型常量
const (
	GoodsDeductionTriggerTypeProfit     = "profit_income"      // 分润入账
	GoodsDeductionTriggerTypeServiceFee = "service_fee_income" // 服务费入账
)

// GoodsDeductionTerminal 货款代扣终端关联
type GoodsDeductionTerminal struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	DeductionID int64     `json:"deduction_id" gorm:"not null;index"` // 关联货款代扣ID
	TerminalID  int64     `json:"terminal_id" gorm:"not null;index"`  // 终端ID
	TerminalSN  string    `json:"terminal_sn" gorm:"size:50"`         // 终端SN（冗余）
	UnitPrice   int64     `json:"unit_price" gorm:"not null"`         // 单价（分）
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
}

func (GoodsDeductionTerminal) TableName() string {
	return "goods_deduction_terminals"
}

// GoodsDeductionNotification 货款代扣通知
type GoodsDeductionNotification struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	DeductionID int64      `json:"deduction_id" gorm:"not null;index"` // 关联货款代扣ID
	DetailID    *int64     `json:"detail_id" gorm:"index"`             // 关联扣款明细
	AgentID     int64      `json:"agent_id" gorm:"not null;index"`     // 接收通知的代理商
	NotifyType  int16      `json:"notify_type" gorm:"not null"`        // 通知类型: 1=待接收 2=扣款通知 3=完成通知
	Title       string     `json:"title" gorm:"size:100"`              // 通知标题
	Content     string     `json:"content" gorm:"size:500"`            // 通知内容
	IsRead      bool       `json:"is_read" gorm:"default:false"`       // 是否已读
	ReadAt      *time.Time `json:"read_at"`                            // 阅读时间
	CreatedAt   time.Time  `json:"created_at" gorm:"default:now()"`
}

func (GoodsDeductionNotification) TableName() string {
	return "goods_deduction_notifications"
}

// GoodsDeductionNotifyType 通知类型
const (
	GoodsDeductionNotifyTypePending   int16 = 1 // 待接收通知
	GoodsDeductionNotifyTypeDeduction int16 = 2 // 扣款通知
	GoodsDeductionNotifyTypeCompleted int16 = 3 // 完成通知
)

// GetGoodsDeductionNotifyTypeName 获取通知类型名称
func GetGoodsDeductionNotifyTypeName(notifyType int16) string {
	switch notifyType {
	case GoodsDeductionNotifyTypePending:
		return "待接收"
	case GoodsDeductionNotifyTypeDeduction:
		return "扣款通知"
	case GoodsDeductionNotifyTypeCompleted:
		return "完成通知"
	default:
		return "未知"
	}
}

// CreateGoodsDeductionRequest 创建货款代扣请求
type CreateGoodsDeductionRequest struct {
	ToAgentID       int64                          `json:"to_agent_id" binding:"required"`       // 下级代理商ID
	UnitPrice       int64                          `json:"unit_price" binding:"required,gt=0"`   // 单价（分）
	DeductionSource int16                          `json:"deduction_source" binding:"required"`  // 扣款来源
	Terminals       []CreateGoodsDeductionTerminal `json:"terminals" binding:"required,dive"`    // 终端列表
	AgreementURL    string                         `json:"agreement_url"`                        // 协议文件URL
	Remark          string                         `json:"remark"`                               // 备注
	DistributeID    *int64                         `json:"distribute_id"`                        // 关联的终端划拨ID
}

// CreateGoodsDeductionTerminal 创建货款代扣终端
type CreateGoodsDeductionTerminal struct {
	TerminalID int64  `json:"terminal_id" binding:"required"` // 终端ID
	TerminalSN string `json:"terminal_sn"`                    // 终端SN
	UnitPrice  int64  `json:"unit_price"`                     // 单价（分），可选，默认使用请求中的统一单价
}

// GoodsDeductionListResponse 货款代扣列表响应
type GoodsDeductionListResponse struct {
	ID              int64      `json:"id"`
	DeductionNo     string     `json:"deduction_no"`
	FromAgentID     int64      `json:"from_agent_id"`
	FromAgentName   string     `json:"from_agent_name"`
	ToAgentID       int64      `json:"to_agent_id"`
	ToAgentName     string     `json:"to_agent_name"`
	TotalAmount     int64      `json:"total_amount"`      // 总金额（分）
	TotalAmountYuan float64    `json:"total_amount_yuan"` // 总金额（元）
	DeductedAmount  int64      `json:"deducted_amount"`
	RemainingAmount int64      `json:"remaining_amount"`
	DeductionSource int16      `json:"deduction_source"`
	SourceName      string     `json:"source_name"`
	TerminalCount   int        `json:"terminal_count"`
	UnitPrice       int64      `json:"unit_price"`
	Status          int16      `json:"status"`
	StatusName      string     `json:"status_name"`
	Progress        float64    `json:"progress"` // 进度百分比
	CreatedAt       time.Time  `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

// ToListResponse 转换为列表响应
func (g *GoodsDeduction) ToListResponse() *GoodsDeductionListResponse {
	progress := float64(0)
	if g.TotalAmount > 0 {
		progress = float64(g.DeductedAmount) / float64(g.TotalAmount) * 100
	}

	return &GoodsDeductionListResponse{
		ID:              g.ID,
		DeductionNo:     g.DeductionNo,
		FromAgentID:     g.FromAgentID,
		FromAgentName:   g.FromAgentName,
		ToAgentID:       g.ToAgentID,
		ToAgentName:     g.ToAgentName,
		TotalAmount:     g.TotalAmount,
		TotalAmountYuan: float64(g.TotalAmount) / 100,
		DeductedAmount:  g.DeductedAmount,
		RemainingAmount: g.RemainingAmount,
		DeductionSource: g.DeductionSource,
		SourceName:      GetGoodsDeductionSourceName(g.DeductionSource),
		TerminalCount:   g.TerminalCount,
		UnitPrice:       g.UnitPrice,
		Status:          g.Status,
		StatusName:      GetGoodsDeductionStatusName(g.Status),
		Progress:        progress,
		CreatedAt:       g.CreatedAt,
		AcceptedAt:      g.AcceptedAt,
		CompletedAt:     g.CompletedAt,
	}
}

// GoodsDeductionSummary 货款代扣统计汇总
type GoodsDeductionSummary struct {
	TotalCount       int64 `json:"total_count"`        // 总代扣数
	PendingCount     int64 `json:"pending_count"`      // 待接收数
	InProgressCount  int64 `json:"in_progress_count"`  // 进行中数
	CompletedCount   int64 `json:"completed_count"`    // 已完成数
	TotalAmount      int64 `json:"total_amount"`       // 代扣总金额（分）
	DeductedAmount   int64 `json:"deducted_amount"`    // 已扣总金额（分）
	RemainingAmount  int64 `json:"remaining_amount"`   // 剩余待扣总金额（分）
}
