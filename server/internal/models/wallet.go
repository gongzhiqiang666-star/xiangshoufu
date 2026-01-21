package models

import (
	"time"
)

// 注意: 钱包类型常量定义在 policy.go 中
// WalletTypeProfit     int16 = 1 // 分润钱包
// WalletTypeService    int16 = 2 // 服务费钱包
// WalletTypeServiceFee int16 = 2 // 服务费钱包（别名）
// WalletTypeReward     int16 = 3 // 奖励钱包
// WalletTypeCharging   int16 = 4 // 充值钱包
// WalletTypeSettlement int16 = 5 // 沉淀钱包

// GetWalletTypeName 获取钱包类型名称（使用 WalletTypeName）
func GetWalletTypeName(walletType int16) string {
	return WalletTypeName(walletType)
}

// AgentWalletConfig 代理商特殊钱包配置
type AgentWalletConfig struct {
	ID int64 `json:"id" gorm:"primaryKey"`

	AgentID int64 `json:"agent_id" gorm:"uniqueIndex"`

	// 充值钱包配置
	ChargingWalletEnabled bool  `json:"charging_wallet_enabled" gorm:"default:false"`
	ChargingWalletLimit   int64 `json:"charging_wallet_limit" gorm:"default:0"` // 充值钱包限额(分)

	// 沉淀钱包配置
	SettlementWalletEnabled bool `json:"settlement_wallet_enabled" gorm:"default:false"`
	SettlementRatio         int  `json:"settlement_ratio" gorm:"default:30"` // 沉淀比例(百分比)

	// 审计字段
	EnabledBy *int64     `json:"enabled_by"`
	EnabledAt *time.Time `json:"enabled_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (AgentWalletConfig) TableName() string {
	return "agent_wallet_configs"
}

// 充值记录状态
const (
	ChargingDepositStatusPending   int16 = 0 // 待确认
	ChargingDepositStatusConfirmed int16 = 1 // 已确认
	ChargingDepositStatusRejected  int16 = 2 // 已拒绝
)

// GetChargingDepositStatusName 获取充值状态名称
func GetChargingDepositStatusName(status int16) string {
	switch status {
	case ChargingDepositStatusPending:
		return "待确认"
	case ChargingDepositStatusConfirmed:
		return "已确认"
	case ChargingDepositStatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}

// ChargingWalletDeposit 充值钱包充值记录
type ChargingWalletDeposit struct {
	ID            int64  `json:"id" gorm:"primaryKey"`
	DepositNo     string `json:"deposit_no" gorm:"size:50;uniqueIndex"`
	AgentID       int64  `json:"agent_id" gorm:"index"`
	Amount        int64  `json:"amount"`                        // 充值金额(分)
	PaymentMethod int16  `json:"payment_method" gorm:"default:1"` // 1=银行转账 2=微信 3=支付宝
	PaymentRef    string `json:"payment_ref" gorm:"size:100"`     // 支付流水号
	Status        int16  `json:"status" gorm:"default:0"`

	// 审核信息
	ConfirmedBy  *int64     `json:"confirmed_by"`
	ConfirmedAt  *time.Time `json:"confirmed_at"`
	RejectReason string     `json:"reject_reason" gorm:"size:500"`

	Remark    string    `json:"remark" gorm:"size:500"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (ChargingWalletDeposit) TableName() string {
	return "charging_wallet_deposits"
}

// 奖励类型
const (
	ChargingRewardTypeManual int16 = 1 // 手动发放
	ChargingRewardTypeAuto   int16 = 2 // 自动发放(政策触发)
)

// 奖励状态
const (
	ChargingRewardStatusIssued  int16 = 1 // 已发放
	ChargingRewardStatusRevoked int16 = 2 // 已撤销
)

// ChargingWalletReward 充值钱包奖励发放记录
type ChargingWalletReward struct {
	ID          int64 `json:"id" gorm:"primaryKey"`
	RewardNo    string `json:"reward_no" gorm:"size:50;uniqueIndex"`
	FromAgentID int64 `json:"from_agent_id" gorm:"index"`
	ToAgentID   int64 `json:"to_agent_id" gorm:"index"`
	Amount      int64 `json:"amount"`                          // 奖励金额(分)
	RewardType  int16 `json:"reward_type" gorm:"default:1"`    // 1=手动发放 2=自动发放
	PolicyID    *int64 `json:"policy_id"`                       // 关联的奖励政策ID

	// 状态
	Status       int16      `json:"status" gorm:"default:1"`
	RevokedAt    *time.Time `json:"revoked_at"`
	RevokeReason string     `json:"revoke_reason" gorm:"size:500"`

	Remark    string    `json:"remark" gorm:"size:500"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (ChargingWalletReward) TableName() string {
	return "charging_wallet_rewards"
}

// GetRewardTypeName 获取奖励类型名称
func GetRewardTypeName(rewardType int16) string {
	switch rewardType {
	case ChargingRewardTypeManual:
		return "手动发放"
	case ChargingRewardTypeAuto:
		return "自动发放"
	default:
		return "未知"
	}
}

// GetRewardStatusName 获取奖励状态名称
func GetRewardStatusName(status int16) string {
	switch status {
	case ChargingRewardStatusIssued:
		return "已发放"
	case ChargingRewardStatusRevoked:
		return "已撤销"
	default:
		return "未知"
	}
}

// 沉淀钱包使用类型
const (
	SettlementUsageTypeUse    int16 = 1 // 使用
	SettlementUsageTypeReturn int16 = 2 // 归还
)

// 沉淀钱包使用状态
const (
	SettlementUsageStatusNormal    int16 = 1 // 正常
	SettlementUsageStatusToReturn  int16 = 2 // 待归还
)

// SettlementWalletUsage 沉淀钱包使用记录
type SettlementWalletUsage struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	UsageNo   string `json:"usage_no" gorm:"size:50;uniqueIndex"`
	AgentID   int64  `json:"agent_id" gorm:"index"`
	Amount    int64  `json:"amount"`                       // 使用金额(分)
	UsageType int16  `json:"usage_type" gorm:"default:1"`  // 1=使用 2=归还

	// 来源明细(哪些下级的钱)
	SourceDetails string `json:"source_details" gorm:"type:jsonb"`

	// 状态
	Status         int16      `json:"status" gorm:"default:1"`
	ReturnDeadline *time.Time `json:"return_deadline"`
	ReturnedAt     *time.Time `json:"returned_at"`

	Remark    string    `json:"remark" gorm:"size:500"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (SettlementWalletUsage) TableName() string {
	return "settlement_wallet_usages"
}
