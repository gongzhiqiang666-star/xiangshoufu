package models

import (
	"time"
)

// ============================================================
// 钱包类型常量
// ============================================================
const (
	WalletTypeProfit     int16 = 1 // 分润钱包
	WalletTypeService    int16 = 2 // 服务费钱包
	WalletTypeServiceFee int16 = 2 // 服务费钱包（别名）
	WalletTypeReward     int16 = 3 // 奖励钱包
	WalletTypeCharging   int16 = 4 // 充值钱包
	WalletTypeSettlement int16 = 5 // 沉淀钱包
)

// WalletTypeName 获取钱包类型名称
func WalletTypeName(walletType int16) string {
	switch walletType {
	case WalletTypeProfit:
		return "分润钱包"
	case WalletTypeService:
		return "服务费钱包"
	case WalletTypeReward:
		return "奖励钱包"
	case WalletTypeCharging:
		return "充值钱包"
	case WalletTypeSettlement:
		return "沉淀钱包"
	default:
		return "未知钱包"
	}
}

// ============================================================
// 押金返现政策
// ============================================================

// DepositCashbackPolicy 押金返现政策（模板级配置）
type DepositCashbackPolicy struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	TemplateID     int64     `json:"template_id" gorm:"not null;index"`     // 政策模板ID
	ChannelID      int64     `json:"channel_id" gorm:"not null;index"`      // 通道ID
	BrandCode      string    `json:"brand_code" gorm:"size:32"`             // 品牌编码
	DepositAmount  int64     `json:"deposit_amount" gorm:"not null"`        // 押金金额（分）
	CashbackAmount int64     `json:"cashback_amount" gorm:"not null"`       // 返现金额（分）
	Status         int16     `json:"status" gorm:"default:1"`               // 1:启用 0:禁用
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:now()"`
}

func (DepositCashbackPolicy) TableName() string {
	return "deposit_cashback_policies"
}

// DepositCashbackRecord 押金返现记录
type DepositCashbackRecord struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	TerminalID     int64      `json:"terminal_id" gorm:"not null;index"`       // 终端ID
	TerminalSN     string     `json:"terminal_sn" gorm:"size:50;not null"`     // 终端SN
	MerchantID     int64      `json:"merchant_id" gorm:"not null;index"`       // 商户ID
	ChannelID      int64      `json:"channel_id" gorm:"not null"`              // 通道ID
	AgentID        int64      `json:"agent_id" gorm:"not null;index"`          // 获得返现的代理商
	DepositAmount  int64      `json:"deposit_amount" gorm:"not null"`          // 押金金额（分）
	SelfCashback   int64      `json:"self_cashback" gorm:"not null"`           // 自身返现配置金额（分）
	UpperCashback  int64      `json:"upper_cashback" gorm:"not null"`          // 上级应返金额（分）
	ActualCashback int64      `json:"actual_cashback" gorm:"not null"`         // 实际返现金额（级差）（分）
	SourceAgentID  *int64     `json:"source_agent_id"`                         // 下级代理商ID（级差来源）
	WalletType     int16      `json:"wallet_type" gorm:"default:2"`            // 钱包类型：2-服务费钱包
	WalletStatus   int16      `json:"wallet_status" gorm:"default:0"`          // 0:待入账 1:已入账
	TriggerType    int16      `json:"trigger_type" gorm:"default:1"`           // 触发类型：1-押金扣款 2-手动触发
	CreatedAt      time.Time  `json:"created_at" gorm:"default:now()"`
	ProcessedAt    *time.Time `json:"processed_at"`
}

func (DepositCashbackRecord) TableName() string {
	return "deposit_cashback_records"
}

// DepositCashbackTriggerType 触发类型
const (
	DepositTriggerTypeAuto   = 1 // 押金扣款自动触发
	DepositTriggerTypeManual = 2 // 手动触发
)

// ============================================================
// 激活奖励政策
// ============================================================

// ActivationRewardPolicy 激活奖励政策（模板级配置）
type ActivationRewardPolicy struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	TemplateID      int64     `json:"template_id" gorm:"not null;index"`     // 政策模板ID
	ChannelID       int64     `json:"channel_id" gorm:"not null;index"`      // 通道ID
	BrandCode       string    `json:"brand_code" gorm:"size:32"`             // 品牌编码
	RewardName      string    `json:"reward_name" gorm:"size:100;not null"`  // 奖励名称
	MinRegisterDays int       `json:"min_register_days" gorm:"default:0"`    // 最小入网天数
	MaxRegisterDays int       `json:"max_register_days" gorm:"default:30"`   // 最大入网天数
	TargetAmount    int64     `json:"target_amount" gorm:"not null"`         // 目标交易量（分）
	RewardAmount    int64     `json:"reward_amount" gorm:"not null"`         // 奖励金额（分）
	RewardType      int16     `json:"reward_type" gorm:"default:1"`          // 奖励类型：1-固定金额 2-交易量比例
	IsCumulative    bool      `json:"is_cumulative" gorm:"default:false"`    // 是否累计
	Priority        int       `json:"priority" gorm:"default:0"`             // 优先级
	Status          int16     `json:"status" gorm:"default:1"`               // 1:启用 0:禁用
	CreatedAt       time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"default:now()"`
}

func (ActivationRewardPolicy) TableName() string {
	return "activation_reward_policies"
}

// ActivationRewardRecord 激活奖励记录
type ActivationRewardRecord struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	PolicyID      int64      `json:"policy_id" gorm:"not null;index"`        // 奖励政策ID
	TerminalID    int64      `json:"terminal_id" gorm:"not null;index"`      // 终端ID
	TerminalSN    string     `json:"terminal_sn" gorm:"size:50;not null"`    // 终端SN
	MerchantID    int64      `json:"merchant_id" gorm:"not null;index"`      // 商户ID
	ChannelID     int64      `json:"channel_id" gorm:"not null"`             // 通道ID
	AgentID       int64      `json:"agent_id" gorm:"not null;index"`         // 获得奖励的代理商
	RegisterDays  int        `json:"register_days" gorm:"not null"`          // 入网天数
	TradeAmount   int64      `json:"trade_amount" gorm:"not null"`           // 达成交易量（分）
	TargetAmount  int64      `json:"target_amount" gorm:"not null"`          // 目标交易量（分）
	SelfReward    int64      `json:"self_reward" gorm:"not null"`            // 自身奖励配置金额（分）
	UpperReward   int64      `json:"upper_reward" gorm:"not null"`           // 上级应返金额（分）
	ActualReward  int64      `json:"actual_reward" gorm:"not null"`          // 实际奖励金额（级差）（分）
	SourceAgentID *int64     `json:"source_agent_id"`                        // 下级代理商ID（级差来源）
	WalletType    int16      `json:"wallet_type" gorm:"default:3"`           // 钱包类型：3-奖励钱包
	WalletStatus  int16      `json:"wallet_status" gorm:"default:0"`         // 0:待入账 1:已入账
	CheckDate     time.Time  `json:"check_date" gorm:"type:date;not null"`   // 检查日期
	CreatedAt     time.Time  `json:"created_at" gorm:"default:now()"`
	ProcessedAt   *time.Time `json:"processed_at"`
}

func (ActivationRewardRecord) TableName() string {
	return "activation_reward_records"
}

// ActivationRewardType 奖励类型
const (
	RewardTypeFixed      = 1 // 固定金额
	RewardTypePercentage = 2 // 交易量比例
)

// ============================================================
// 费率阶梯政策（代理商调价）
// ============================================================

// RateStagePolicy 费率阶梯政策
type RateStagePolicy struct {
	ID                int64     `json:"id" gorm:"primaryKey"`
	TemplateID        int64     `json:"template_id" gorm:"not null;index"`      // 政策模板ID
	ChannelID         int64     `json:"channel_id" gorm:"not null;index"`       // 通道ID
	BrandCode         string    `json:"brand_code" gorm:"size:32"`              // 品牌编码
	StageName         string    `json:"stage_name" gorm:"size:100;not null"`    // 阶梯名称
	ApplyTo           int16     `json:"apply_to" gorm:"not null;default:1"`     // 应用对象：1-商户 2-代理商
	MinDays           int       `json:"min_days" gorm:"default:0"`              // 最小入网天数
	MaxDays           int       `json:"max_days" gorm:"not null"`               // 最大入网天数（-1表示无限）
	CreditRateDelta   string    `json:"credit_rate_delta" gorm:"type:decimal(10,4);default:0"`   // 贷记卡费率调整值
	DebitRateDelta    string    `json:"debit_rate_delta" gorm:"type:decimal(10,4);default:0"`    // 借记卡费率调整值
	UnionpayRateDelta string    `json:"unionpay_rate_delta" gorm:"type:decimal(10,4);default:0"` // 云闪付费率调整值
	WechatRateDelta   string    `json:"wechat_rate_delta" gorm:"type:decimal(10,4);default:0"`   // 微信费率调整值
	AlipayRateDelta   string    `json:"alipay_rate_delta" gorm:"type:decimal(10,4);default:0"`   // 支付宝费率调整值
	Priority          int       `json:"priority" gorm:"default:0"`              // 优先级
	Status            int16     `json:"status" gorm:"default:1"`                // 1:启用 0:禁用
	CreatedAt         time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"default:now()"`
}

func (RateStagePolicy) TableName() string {
	return "rate_stage_policies"
}

// RateStageApplyTo 应用对象
const (
	RateStageApplyToMerchant = 1 // 商户
	RateStageApplyToAgent    = 2 // 代理商
)

// ============================================================
// 代理商个性化政策配置
// ============================================================

// AgentDepositCashbackPolicy 代理商押金返现政策
type AgentDepositCashbackPolicy struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	AgentID        int64     `json:"agent_id" gorm:"not null;index"`        // 代理商ID
	ChannelID      int64     `json:"channel_id" gorm:"not null;index"`      // 通道ID
	DepositAmount  int64     `json:"deposit_amount" gorm:"not null"`        // 押金金额（分）
	CashbackAmount int64     `json:"cashback_amount" gorm:"not null"`       // 返现金额（分）
	Status         int16     `json:"status" gorm:"default:1"`               // 1:启用 0:禁用
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:now()"`
}

func (AgentDepositCashbackPolicy) TableName() string {
	return "agent_deposit_cashback_policies"
}

// AgentSimCashbackPolicy 代理商流量卡返现政策
type AgentSimCashbackPolicy struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	AgentID            int64     `json:"agent_id" gorm:"not null;index"`      // 代理商ID
	ChannelID          int64     `json:"channel_id" gorm:"not null;index"`    // 通道ID
	BrandCode          string    `json:"brand_code" gorm:"size:32"`           // 品牌编码
	FirstTimeCashback  int64     `json:"first_time_cashback" gorm:"not null"` // 首次返现金额（分）
	SecondTimeCashback int64     `json:"second_time_cashback" gorm:"not null"`// 第2次返现金额（分）
	ThirdPlusCashback  int64     `json:"third_plus_cashback" gorm:"not null"` // 第3次及以后返现金额（分）
	Status             int16     `json:"status" gorm:"default:1"`             // 1:启用 0:禁用
	CreatedAt          time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"default:now()"`
}

func (AgentSimCashbackPolicy) TableName() string {
	return "agent_sim_cashback_policies"
}

// AgentActivationRewardPolicy 代理商激活奖励政策
type AgentActivationRewardPolicy struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	AgentID         int64     `json:"agent_id" gorm:"not null;index"`       // 代理商ID
	ChannelID       int64     `json:"channel_id" gorm:"not null;index"`     // 通道ID
	BrandCode       string    `json:"brand_code" gorm:"size:32"`            // 品牌编码
	RewardName      string    `json:"reward_name" gorm:"size:100;not null"` // 奖励名称
	MinRegisterDays int       `json:"min_register_days" gorm:"default:0"`   // 最小入网天数
	MaxRegisterDays int       `json:"max_register_days" gorm:"default:30"`  // 最大入网天数
	TargetAmount    int64     `json:"target_amount" gorm:"not null"`        // 目标交易量（分）
	RewardAmount    int64     `json:"reward_amount" gorm:"not null"`        // 奖励金额（分）
	Priority        int       `json:"priority" gorm:"default:0"`            // 优先级
	Status          int16     `json:"status" gorm:"default:1"`              // 1:启用 0:禁用
	CreatedAt       time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"default:now()"`
}

func (AgentActivationRewardPolicy) TableName() string {
	return "agent_activation_reward_policies"
}

// ============================================================
// 政策模板扩展（包含4块政策）
// ============================================================

// PolicyTemplateComplete 完整政策模板（包含4块政策配置）
type PolicyTemplateComplete struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	TemplateName string    `json:"template_name" gorm:"size:100;not null"`
	ChannelID    int64     `json:"channel_id" gorm:"not null;index"`
	IsDefault    bool      `json:"is_default" gorm:"default:false"`

	// 成本（费率）
	CreditRate   string    `json:"credit_rate" gorm:"type:decimal(10,4)"`   // 贷记卡费率
	DebitRate    string    `json:"debit_rate" gorm:"type:decimal(10,4)"`    // 借记卡费率
	DebitCap     string    `json:"debit_cap" gorm:"type:decimal(10,2)"`     // 借记卡封顶
	UnionpayRate string    `json:"unionpay_rate" gorm:"type:decimal(10,4)"` // 云闪付费率
	WechatRate   string    `json:"wechat_rate" gorm:"type:decimal(10,4)"`   // 微信费率
	AlipayRate   string    `json:"alipay_rate" gorm:"type:decimal(10,4)"`   // 支付宝费率

	Status    int16     `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`

	// 关联的政策配置（非数据库字段）
	DepositCashbackPolicies   []DepositCashbackPolicy   `json:"deposit_cashback_policies" gorm:"-"`
	SimCashbackPolicies       []SimCashbackPolicy       `json:"sim_cashback_policies" gorm:"-"`
	ActivationRewardPolicies  []ActivationRewardPolicy  `json:"activation_reward_policies" gorm:"-"`
	RateStagePolicies         []RateStagePolicy         `json:"rate_stage_policies" gorm:"-"`
}

func (PolicyTemplateComplete) TableName() string {
	return "policy_templates"
}

// ============================================================
// 代理商完整政策配置（用于API响应）
// ============================================================

// AgentPolicyComplete 代理商完整政策配置
type AgentPolicyComplete struct {
	AgentID   int64  `json:"agent_id"`
	ChannelID int64  `json:"channel_id"`

	// 成本（费率）
	CreditRate   string `json:"credit_rate"`
	DebitRate    string `json:"debit_rate"`
	DebitCap     string `json:"debit_cap"`
	UnionpayRate string `json:"unionpay_rate"`
	WechatRate   string `json:"wechat_rate"`
	AlipayRate   string `json:"alipay_rate"`

	// 押金返现配置
	DepositCashbacks []DepositCashbackConfig `json:"deposit_cashbacks"`

	// 流量卡返现配置
	SimCashback *SimCashbackConfig `json:"sim_cashback"`

	// 激活奖励配置
	ActivationRewards []ActivationRewardConfig `json:"activation_rewards"`
}

// DepositCashbackConfig 押金返现配置
type DepositCashbackConfig struct {
	DepositAmount  int64 `json:"deposit_amount"`  // 押金金额（分）
	CashbackAmount int64 `json:"cashback_amount"` // 返现金额（分）
}

// SimCashbackConfig 流量卡返现配置
type SimCashbackConfig struct {
	FirstTimeCashback  int64 `json:"first_time_cashback"`  // 首次返现金额（分）
	SecondTimeCashback int64 `json:"second_time_cashback"` // 第2次返现金额（分）
	ThirdPlusCashback  int64 `json:"third_plus_cashback"`  // 第3次及以后返现金额（分）
}

// ActivationRewardConfig 激活奖励配置
type ActivationRewardConfig struct {
	RewardName      string `json:"reward_name"`       // 奖励名称
	MinRegisterDays int    `json:"min_register_days"` // 最小入网天数
	MaxRegisterDays int    `json:"max_register_days"` // 最大入网天数
	TargetAmount    int64  `json:"target_amount"`     // 目标交易量（分）
	RewardAmount    int64  `json:"reward_amount"`     // 奖励金额（分）
}

// ============================================================
// 政策限制（用于下级调整政策时的范围限制）
// ============================================================

// PolicyLimits 政策限制（上级的政策配置，下级不能超过）
type PolicyLimits struct {
	// 费率限制（下级费率不能低于上级）
	MinCreditRate   string `json:"min_credit_rate"`
	MinDebitRate    string `json:"min_debit_rate"`
	MinUnionpayRate string `json:"min_unionpay_rate"`
	MinWechatRate   string `json:"min_wechat_rate"`
	MinAlipayRate   string `json:"min_alipay_rate"`

	// 押金返现限制（下级返现不能高于上级）
	MaxDepositCashbacks []DepositCashbackConfig `json:"max_deposit_cashbacks"`

	// 流量卡返现限制（下级返现不能高于上级）
	MaxSimCashback *SimCashbackConfig `json:"max_sim_cashback"`

	// 激活奖励限制（下级奖励不能高于上级）
	MaxActivationRewards []ActivationRewardConfig `json:"max_activation_rewards"`
}
