package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================
// 结算价相关常量
// ============================================================

// ChangeType 变更类型
type ChangeType int16

const (
	ChangeTypeInit       ChangeType = 1 // 初始化（从模板创建）
	ChangeTypeRate       ChangeType = 2 // 费率调整
	ChangeTypeDeposit    ChangeType = 3 // 押金返现调整
	ChangeTypeSim        ChangeType = 4 // 流量费返现调整
	ChangeTypeActivation ChangeType = 5 // 激活奖励调整
	ChangeTypeBatch      ChangeType = 6 // 批量调整
	ChangeTypeSync       ChangeType = 7 // 模板同步
)

// ChangeTypeName 获取变更类型名称
func ChangeTypeName(ct ChangeType) string {
	switch ct {
	case ChangeTypeInit:
		return "初始化"
	case ChangeTypeRate:
		return "费率调整"
	case ChangeTypeDeposit:
		return "押金返现调整"
	case ChangeTypeSim:
		return "流量费返现调整"
	case ChangeTypeActivation:
		return "激活奖励调整"
	case ChangeTypeBatch:
		return "批量调整"
	case ChangeTypeSync:
		return "模板同步"
	default:
		return "未知类型"
	}
}

// ConfigType 配置类型
type ConfigType int16

const (
	ConfigTypeSettlement ConfigType = 1 // 结算价
	ConfigTypeReward     ConfigType = 2 // 奖励配置
)

// ConfigTypeName 获取配置类型名称
func ConfigTypeName(ct ConfigType) string {
	switch ct {
	case ConfigTypeSettlement:
		return "结算价"
	case ConfigTypeReward:
		return "奖励配置"
	default:
		return "未知类型"
	}
}

// OperatorType 操作者类型
type OperatorType int16

const (
	OperatorTypeAdmin OperatorType = 1 // 管理员
	OperatorTypeAgent OperatorType = 2 // 代理商
)

// ============================================================
// 押金返现配置JSONB类型
// ============================================================

// DepositCashbackItem 押金返现配置项
type DepositCashbackItem struct {
	DepositAmount  int64 `json:"deposit_amount"`  // 押金金额（分）
	CashbackAmount int64 `json:"cashback_amount"` // 返现金额（分）
}

// DepositCashbacks 押金返现配置列表
type DepositCashbacks []DepositCashbackItem

// Scan 实现sql.Scanner接口
func (d *DepositCashbacks) Scan(value interface{}) error {
	if value == nil {
		*d = make(DepositCashbacks, 0)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into DepositCashbacks", value)
	}

	if len(bytes) == 0 || string(bytes) == "[]" {
		*d = make(DepositCashbacks, 0)
		return nil
	}

	return json.Unmarshal(bytes, d)
}

// Value 实现driver.Valuer接口
func (d DepositCashbacks) Value() (driver.Value, error) {
	if d == nil {
		return "[]", nil
	}
	return json.Marshal(d)
}

// ============================================================
// 激活奖励配置JSONB类型
// ============================================================

// ActivationRewardItem 激活奖励配置项
type ActivationRewardItem struct {
	RewardName      string `json:"reward_name"`       // 奖励名称
	MinRegisterDays int    `json:"min_register_days"` // 最小入网天数
	MaxRegisterDays int    `json:"max_register_days"` // 最大入网天数
	TargetAmount    int64  `json:"target_amount"`     // 目标交易量（分）
	RewardAmount    int64  `json:"reward_amount"`     // 奖励金额（分）
	Priority        int    `json:"priority"`          // 优先级
}

// ActivationRewards 激活奖励配置列表
type ActivationRewards []ActivationRewardItem

// Scan 实现sql.Scanner接口
func (a *ActivationRewards) Scan(value interface{}) error {
	if value == nil {
		*a = make(ActivationRewards, 0)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into ActivationRewards", value)
	}

	if len(bytes) == 0 || string(bytes) == "[]" {
		*a = make(ActivationRewards, 0)
		return nil
	}

	return json.Unmarshal(bytes, a)
}

// Value 实现driver.Valuer接口
func (a ActivationRewards) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

// ============================================================
// 通道结算价表
// ============================================================

// SettlementPrice 通道结算价
type SettlementPrice struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	AgentID   int64  `json:"agent_id" gorm:"not null;index"`
	ChannelID int64  `json:"channel_id" gorm:"not null;index"`
	TemplateID *int64 `json:"template_id" gorm:"index"`
	BrandCode string `json:"brand_code" gorm:"size:32;default:''"`

	// 费率配置（动态JSONB）
	RateConfigs RateConfigs `json:"rate_configs" gorm:"type:jsonb;default:'{}'"`

	// 旧字段保留兼容
	CreditRate   *string `json:"credit_rate" gorm:"type:decimal(10,4)"`
	DebitRate    *string `json:"debit_rate" gorm:"type:decimal(10,4)"`
	DebitCap     *string `json:"debit_cap" gorm:"type:decimal(10,2)"`
	UnionpayRate *string `json:"unionpay_rate" gorm:"type:decimal(10,4)"`
	WechatRate   *string `json:"wechat_rate" gorm:"type:decimal(10,4)"`
	AlipayRate   *string `json:"alipay_rate" gorm:"type:decimal(10,4)"`

	// 押金返现配置（JSONB数组）
	DepositCashbacks DepositCashbacks `json:"deposit_cashbacks" gorm:"type:jsonb;default:'[]'"`

	// 流量费返现配置
	SimFirstCashback     int64 `json:"sim_first_cashback" gorm:"default:0"`
	SimSecondCashback    int64 `json:"sim_second_cashback" gorm:"default:0"`
	SimThirdPlusCashback int64 `json:"sim_third_plus_cashback" gorm:"default:0"`

	// 元数据
	Version     int        `json:"version" gorm:"default:1"`
	Status      int16      `json:"status" gorm:"default:1"`
	EffectiveAt time.Time  `json:"effective_at" gorm:"default:now()"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:now()"`
	CreatedBy   *int64     `json:"created_by"`
	UpdatedBy   *int64     `json:"updated_by"`
}

// TableName 表名
func (SettlementPrice) TableName() string {
	return "settlement_prices"
}

// ============================================================
// 代理商奖励配置表
// ============================================================

// AgentRewardSetting 代理商奖励配置
type AgentRewardSetting struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
	AgentID    int64  `json:"agent_id" gorm:"not null;uniqueIndex"`
	TemplateID *int64 `json:"template_id" gorm:"index"`

	// 奖励金额（差额分配模式）
	RewardAmount int64 `json:"reward_amount" gorm:"default:0"`

	// 激活奖励配置（JSONB数组）
	ActivationRewards ActivationRewards `json:"activation_rewards" gorm:"type:jsonb;default:'[]'"`

	// 元数据
	Version     int        `json:"version" gorm:"default:1"`
	Status      int16      `json:"status" gorm:"default:1"`
	EffectiveAt time.Time  `json:"effective_at" gorm:"default:now()"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:now()"`
	CreatedBy   *int64     `json:"created_by"`
	UpdatedBy   *int64     `json:"updated_by"`
}

// TableName 表名
func (AgentRewardSetting) TableName() string {
	return "agent_reward_settings"
}

// ============================================================
// 调价记录表
// ============================================================

// PriceChangeLog 调价记录
type PriceChangeLog struct {
	ID                int64  `json:"id" gorm:"primaryKey"`
	AgentID           int64  `json:"agent_id" gorm:"not null;index"`
	ChannelID         *int64 `json:"channel_id" gorm:"index"`
	SettlementPriceID *int64 `json:"settlement_price_id" gorm:"index"`
	RewardSettingID   *int64 `json:"reward_setting_id" gorm:"index"`

	// 变更类型
	ChangeType ChangeType `json:"change_type" gorm:"not null"`
	ConfigType ConfigType `json:"config_type" gorm:"not null"`

	// 变更内容
	FieldName     string `json:"field_name" gorm:"size:100"`
	OldValue      string `json:"old_value" gorm:"type:text"`
	NewValue      string `json:"new_value" gorm:"type:text"`
	ChangeSummary string `json:"change_summary" gorm:"size:500"`

	// 完整快照
	SnapshotBefore JSONMap `json:"snapshot_before" gorm:"type:jsonb"`
	SnapshotAfter  JSONMap `json:"snapshot_after" gorm:"type:jsonb"`

	// 操作信息
	OperatorType OperatorType `json:"operator_type" gorm:"default:1"`
	OperatorID   int64        `json:"operator_id" gorm:"not null;index"`
	OperatorName string       `json:"operator_name" gorm:"size:100"`
	Source       string       `json:"source" gorm:"size:20;default:'PC'"`
	IPAddress    string       `json:"ip_address" gorm:"size:50"`
	UserAgent    string       `json:"user_agent" gorm:"size:500"`

	CreatedAt time.Time `json:"created_at" gorm:"default:now();index"`
}

// TableName 表名
func (PriceChangeLog) TableName() string {
	return "price_change_logs"
}

// JSONMap 通用JSON Map类型
type JSONMap map[string]interface{}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into JSONMap", value)
	}

	if len(bytes) == 0 || string(bytes) == "{}" {
		*j = make(JSONMap)
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// ============================================================
// 请求/响应DTO
// ============================================================

// CreateSettlementPriceRequest 创建结算价请求
type CreateSettlementPriceRequest struct {
	AgentID    int64  `json:"agent_id" binding:"required"`
	ChannelID  int64  `json:"channel_id" binding:"required"`
	TemplateID *int64 `json:"template_id"`
	BrandCode  string `json:"brand_code"`
}

// UpdateRateRequest 更新费率请求
type UpdateRateRequest struct {
	RateConfigs  RateConfigs `json:"rate_configs"`
	CreditRate   *string     `json:"credit_rate"`
	DebitRate    *string     `json:"debit_rate"`
	DebitCap     *string     `json:"debit_cap"`
	UnionpayRate *string     `json:"unionpay_rate"`
	WechatRate   *string     `json:"wechat_rate"`
	AlipayRate   *string     `json:"alipay_rate"`
}

// UpdateDepositCashbackRequest 更新押金返现请求
type UpdateDepositCashbackRequest struct {
	DepositCashbacks DepositCashbacks `json:"deposit_cashbacks" binding:"required"`
}

// UpdateSimCashbackRequest 更新流量费返现请求
type UpdateSimCashbackRequest struct {
	SimFirstCashback     int64 `json:"sim_first_cashback"`
	SimSecondCashback    int64 `json:"sim_second_cashback"`
	SimThirdPlusCashback int64 `json:"sim_third_plus_cashback"`
}

// UpdateActivationRewardRequest 更新激活奖励请求
type UpdateActivationRewardRequest struct {
	ActivationRewards ActivationRewards `json:"activation_rewards" binding:"required"`
}

// SettlementPriceListRequest 结算价列表请求
type SettlementPriceListRequest struct {
	AgentID   *int64 `form:"agent_id"`
	ChannelID *int64 `form:"channel_id"`
	Status    *int16 `form:"status"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// SettlementPriceListResponse 结算价列表响应
type SettlementPriceListResponse struct {
	List  []SettlementPriceItem `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// SettlementPriceItem 结算价列表项
type SettlementPriceItem struct {
	ID          int64            `json:"id"`
	AgentID     int64            `json:"agent_id"`
	AgentName   string           `json:"agent_name"`
	ChannelID   int64            `json:"channel_id"`
	ChannelName string           `json:"channel_name"`
	BrandCode   string           `json:"brand_code"`
	RateConfigs RateConfigs      `json:"rate_configs"`
	DepositCashbacks DepositCashbacks `json:"deposit_cashbacks"`
	SimFirstCashback     int64  `json:"sim_first_cashback"`
	SimSecondCashback    int64  `json:"sim_second_cashback"`
	SimThirdPlusCashback int64  `json:"sim_third_plus_cashback"`
	Version     int              `json:"version"`
	Status      int16            `json:"status"`
	EffectiveAt time.Time        `json:"effective_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// PriceChangeLogListRequest 调价记录列表请求
type PriceChangeLogListRequest struct {
	AgentID    *int64      `form:"agent_id"`
	ChannelID  *int64      `form:"channel_id"`
	ChangeType *ChangeType `form:"change_type"`
	ConfigType *ConfigType `form:"config_type"`
	StartDate  string      `form:"start_date"`
	EndDate    string      `form:"end_date"`
	Page       int         `form:"page,default=1"`
	PageSize   int         `form:"page_size,default=20"`
}

// PriceChangeLogListResponse 调价记录列表响应
type PriceChangeLogListResponse struct {
	List  []PriceChangeLogItem `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// PriceChangeLogItem 调价记录列表项
type PriceChangeLogItem struct {
	ID             int64      `json:"id"`
	AgentID        int64      `json:"agent_id"`
	AgentName      string     `json:"agent_name"`
	ChannelID      *int64     `json:"channel_id"`
	ChannelName    string     `json:"channel_name"`
	ChangeType     ChangeType `json:"change_type"`
	ChangeTypeName string     `json:"change_type_name"`
	ConfigType     ConfigType `json:"config_type"`
	ConfigTypeName string     `json:"config_type_name"`
	FieldName      string     `json:"field_name"`
	OldValue       string     `json:"old_value"`
	NewValue       string     `json:"new_value"`
	ChangeSummary  string     `json:"change_summary"`
	OperatorName   string     `json:"operator_name"`
	Source         string     `json:"source"`
	CreatedAt      time.Time  `json:"created_at"`
}

// AgentRewardSettingRequest 代理商奖励配置请求
type AgentRewardSettingRequest struct {
	AgentID      int64  `json:"agent_id" binding:"required"`
	TemplateID   *int64 `json:"template_id"`
	RewardAmount int64  `json:"reward_amount"`
}

// AgentRewardSettingResponse 代理商奖励配置响应
type AgentRewardSettingResponse struct {
	ID                int64             `json:"id"`
	AgentID           int64             `json:"agent_id"`
	AgentName         string            `json:"agent_name"`
	TemplateID        *int64            `json:"template_id"`
	TemplateName      string            `json:"template_name"`
	RewardAmount      int64             `json:"reward_amount"`
	ActivationRewards ActivationRewards `json:"activation_rewards"`
	Version           int               `json:"version"`
	Status            int16             `json:"status"`
	EffectiveAt       time.Time         `json:"effective_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
