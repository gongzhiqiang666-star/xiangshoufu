package models

import (
	"time"
)

// Terminal 终端/机具
type Terminal struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	TerminalSN   string     `json:"terminal_sn" gorm:"size:50;uniqueIndex"` // 终端序列号
	ChannelID    int64      `json:"channel_id" gorm:"not null;index"`       // 所属通道
	ChannelCode  string     `json:"channel_code" gorm:"size:32"`            // 通道编码
	BrandCode    string     `json:"brand_code" gorm:"size:32"`              // 品牌编码
	ModelCode    string     `json:"model_code" gorm:"size:32"`              // 型号编码
	OwnerAgentID int64      `json:"owner_agent_id" gorm:"index"`            // 当前所属代理商
	MerchantID   *int64     `json:"merchant_id" gorm:"index"`               // 绑定的商户
	MerchantNo   string     `json:"merchant_no" gorm:"size:64;index"`       // 商户号
	Status       int16      `json:"status" gorm:"default:1"`                // 1:待分配 2:已分配 3:已绑定 4:已激活 5:已解绑 6:已回收
	ActivatedAt  *time.Time `json:"activated_at"`                           // 激活时间
	BoundAt      *time.Time `json:"bound_at"`                               // 绑定时间
	SimFeeCount  int        `json:"sim_fee_count" gorm:"default:0"`         // 流量费缴费次数
	LastSimFeeAt *time.Time `json:"last_sim_fee_at"`                        // 最后缴费时间
	CreatedAt    time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:now()"`
}

func (Terminal) TableName() string {
	return "terminals"
}

// TerminalStatus 终端状态
const (
	TerminalStatusPending   = 1 // 待分配
	TerminalStatusAllocated = 2 // 已分配
	TerminalStatusBound     = 3 // 已绑定
	TerminalStatusActivated = 4 // 已激活
	TerminalStatusUnbound   = 5 // 已解绑
	TerminalStatusRecycled  = 6 // 已回收
)

// TerminalDistribute 终端下发记录
type TerminalDistribute struct {
	ID              int64      `json:"id" gorm:"primaryKey"`
	DistributeNo    string     `json:"distribute_no" gorm:"size:64;uniqueIndex"`  // 下发单号
	FromAgentID     int64      `json:"from_agent_id" gorm:"not null;index"`       // 下发方代理商
	ToAgentID       int64      `json:"to_agent_id" gorm:"not null;index"`         // 接收方代理商
	TerminalSN      string     `json:"terminal_sn" gorm:"size:50;not null;index"` // 终端SN
	ChannelID       int64      `json:"channel_id" gorm:"not null"`                // 通道ID
	IsCrossLevel    bool       `json:"is_cross_level" gorm:"default:false"`       // 是否跨级下发
	CrossLevelPath  string     `json:"cross_level_path" gorm:"size:500"`          // 跨级路径 /A/B/C/
	GoodsPrice      int64      `json:"goods_price" gorm:"not null"`               // 货款金额（分）
	DeductionType   int16      `json:"deduction_type" gorm:"not null"`            // 1:一次性付款 2:分期代扣
	DeductionPlanID *int64     `json:"deduction_plan_id"`                         // 关联代扣计划ID
	ChainID         *int64     `json:"chain_id"`                                  // 关联代扣链ID（跨级时）
	Status          int16      `json:"status" gorm:"default:1"`                   // 1:待确认 2:已确认 3:已拒绝 4:已取消
	Source          int16      `json:"source" gorm:"not null"`                    // 1:APP 2:PC
	Remark          string     `json:"remark" gorm:"size:255"`
	CreatedBy       int64      `json:"created_by"`
	ConfirmedBy     *int64     `json:"confirmed_by"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:now()"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
}

func (TerminalDistribute) TableName() string {
	return "terminal_distributes"
}

// TerminalDistributeStatus 终端下发状态
const (
	TerminalDistributeStatusPending   = 1 // 待确认
	TerminalDistributeStatusConfirmed = 2 // 已确认
	TerminalDistributeStatusRejected  = 3 // 已拒绝
	TerminalDistributeStatusCancelled = 4 // 已取消
)

// TerminalDistributeSource 终端下发来源
const (
	TerminalDistributeSourceApp = 1 // APP端（不能跨级）
	TerminalDistributeSourcePC  = 2 // PC端（可以跨级）
)

// TerminalDistributeDeductionType 货款扣款方式
const (
	TerminalDistributeDeductionOneTime     = 1 // 一次性付款
	TerminalDistributeDeductionInstallment = 2 // 分期代扣
)

// TerminalRecall 终端回拨记录
type TerminalRecall struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	RecallNo       string     `json:"recall_no" gorm:"size:64;uniqueIndex"`     // 回拨单号
	FromAgentID    int64      `json:"from_agent_id" gorm:"not null;index"`      // 回拨方代理商（当前持有者）
	ToAgentID      int64      `json:"to_agent_id" gorm:"not null;index"`        // 接收方代理商（上级）
	TerminalSN     string     `json:"terminal_sn" gorm:"size:50;not null;index"`// 终端SN
	ChannelID      int64      `json:"channel_id" gorm:"not null"`               // 通道ID
	IsCrossLevel   bool       `json:"is_cross_level" gorm:"default:false"`      // 是否跨级回拨
	CrossLevelPath string     `json:"cross_level_path" gorm:"size:500"`         // 跨级路径
	Status         int16      `json:"status" gorm:"default:1"`                  // 1:待确认 2:已确认 3:已拒绝 4:已取消
	Source         int16      `json:"source" gorm:"not null"`                   // 1:APP 2:PC
	Remark         string     `json:"remark" gorm:"size:255"`
	CreatedBy      int64      `json:"created_by"`
	ConfirmedBy    *int64     `json:"confirmed_by"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:now()"`
	ConfirmedAt    *time.Time `json:"confirmed_at"`
}

func (TerminalRecall) TableName() string {
	return "terminal_recalls"
}

// TerminalRecallStatus 终端回拨状态
const (
	TerminalRecallStatusPending   = 1 // 待确认
	TerminalRecallStatusConfirmed = 2 // 已确认
	TerminalRecallStatusRejected  = 3 // 已拒绝
	TerminalRecallStatusCancelled = 4 // 已取消
)

// TerminalImportRecord 终端入库记录
type TerminalImportRecord struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ImportNo    string    `json:"import_no" gorm:"size:64;uniqueIndex"`   // 入库批次号
	ChannelID   int64     `json:"channel_id" gorm:"not null;index"`       // 通道ID
	ChannelCode string    `json:"channel_code" gorm:"size:32"`            // 通道编码
	BrandCode   string    `json:"brand_code" gorm:"size:32"`              // 品牌编码
	ModelCode   string    `json:"model_code" gorm:"size:32"`              // 型号编码
	TotalCount  int       `json:"total_count" gorm:"not null"`            // 导入总数
	SuccessCount int      `json:"success_count" gorm:"default:0"`         // 成功数
	FailedCount int       `json:"failed_count" gorm:"default:0"`          // 失败数
	FailedSNs   string    `json:"failed_sns" gorm:"type:text"`            // 失败SN列表（JSON数组）
	OwnerAgentID int64    `json:"owner_agent_id" gorm:"not null;index"`   // 入库代理商
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
}

func (TerminalImportRecord) TableName() string {
	return "terminal_import_records"
}

// TerminalPolicy 终端政策设置
type TerminalPolicy struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	TerminalSN string    `json:"terminal_sn" gorm:"size:50;uniqueIndex"` // 终端SN
	ChannelID  int64     `json:"channel_id" gorm:"not null;index"`       // 通道ID
	AgentID    int64     `json:"agent_id" gorm:"not null;index"`         // 所属代理商

	// 费率设置 (万分比)
	CreditRate    int `json:"credit_rate" gorm:"default:0"`    // 贷记卡费率，如 55 表示万分之55
	DebitRate     int `json:"debit_rate" gorm:"default:0"`     // 借记卡费率
	DebitCap      int `json:"debit_cap" gorm:"default:0"`      // 借记卡封顶(分)
	UnionpayRate  int `json:"unionpay_rate" gorm:"default:0"`  // 银联云闪付费率
	WechatRate    int `json:"wechat_rate" gorm:"default:0"`    // 微信扫码费率
	AlipayRate    int `json:"alipay_rate" gorm:"default:0"`    // 支付宝扫码费率

	// SIM卡费用设置 (分)
	FirstSimFee         int `json:"first_sim_fee" gorm:"default:0"`          // 首次流量费
	NonFirstSimFee      int `json:"non_first_sim_fee" gorm:"default:0"`      // 非首次流量费
	SimFeeIntervalDays  int `json:"sim_fee_interval_days" gorm:"default:0"`  // 非首次间隔天数

	// 押金设置 (分)
	DepositAmount int `json:"deposit_amount" gorm:"default:0"` // 押金金额，0表示无押金

	// 状态
	IsSynced   bool       `json:"is_synced" gorm:"default:false"`  // 是否已同步到通道
	SyncedAt   *time.Time `json:"synced_at"`                       // 同步时间
	SyncError  string     `json:"sync_error" gorm:"size:500"`      // 同步错误信息

	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

func (TerminalPolicy) TableName() string {
	return "terminal_policies"
}

// 预设费率选项 (万分比)
var RateOptions = []int{53, 54, 55, 56, 57, 58, 59, 60}

// 预设首次流量费选项 (分)
var FirstSimFeeOptions = []int{4800, 6000, 6900, 7900, 8900, 9900}

// 预设非首次流量费选项 (分)
var NonFirstSimFeeOptions = []int{4800, 6000, 6900, 7900, 8900, 9900}

// 预设非首次间隔天数选项
var SimFeeIntervalDaysOptions = []int{180, 210, 240, 270, 300, 330, 360}

// 预设押金选项 (分)
var DepositOptions = []int{0, 9900, 19900, 29900}
