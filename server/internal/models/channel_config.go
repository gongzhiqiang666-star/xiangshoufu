package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================
// 通道费率配置
// ============================================================

// ChannelRateConfig 通道费率配置
type ChannelRateConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ChannelID   int64     `json:"channel_id" gorm:"not null;index"`
	RateCode    string    `json:"rate_code" gorm:"size:32;not null"`    // 费率编码 (CREDIT/DEBIT/WECHAT等)
	RateName    string    `json:"rate_name" gorm:"size:64;not null"`    // 费率名称
	MinRate     string    `json:"min_rate" gorm:"type:decimal(10,4)"`   // 最低成本（通道底价）
	MaxRate     string    `json:"max_rate" gorm:"type:decimal(10,4)"`   // 最高限制
	DefaultRate string    `json:"default_rate" gorm:"type:decimal(10,4)"` // 默认费率
	MaxHighRate *string   `json:"max_high_rate" gorm:"size:10"`         // 高调费率上限
	MaxD0Extra  *int64    `json:"max_d0_extra"`                         // P+0加价上限（分）
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	Status      int16     `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (ChannelRateConfig) TableName() string {
	return "channel_rate_configs"
}

// ============================================================
// 通道流量费返现档位
// ============================================================

// ChannelSimCashbackTier 通道流量费返现档位
type ChannelSimCashbackTier struct {
	ID                int64     `json:"id" gorm:"primaryKey"`
	ChannelID         int64     `json:"channel_id" gorm:"not null;index"`
	BrandCode         string    `json:"brand_code" gorm:"size:32;default:''"` // 品牌编码（空=通用）
	TierOrder         int       `json:"tier_order" gorm:"not null"`           // 档位序号 1=首次, 2=第2次...
	TierName          string    `json:"tier_name" gorm:"size:64;not null"`    // 档位名称
	IsLastTier        bool      `json:"is_last_tier" gorm:"default:false"`    // 是否最后档(N次及以后)
	MaxCashbackAmount int64     `json:"max_cashback_amount" gorm:"not null"`  // 返现上限（分）
	DefaultCashback   int64     `json:"default_cashback" gorm:"default:0"`    // 默认返现（分）
	SimFeeAmount      int64     `json:"sim_fee_amount" gorm:"not null"`       // 流量费金额（分）
	Status            int16     `json:"status" gorm:"default:1"`
	CreatedAt         time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (ChannelSimCashbackTier) TableName() string {
	return "channel_sim_cashback_tiers"
}

// ============================================================
// 流量费返现N档配置类型（用于结算价和政策模板）
// ============================================================

// SimCashbackItem 流量费返现配置项
type SimCashbackItem struct {
	TierOrder      int   `json:"tier_order"`      // 档位序号
	CashbackAmount int64 `json:"cashback_amount"` // 返现金额（分）
}

// SimCashbacks 流量费返现配置列表
type SimCashbacks []SimCashbackItem

// Scan 实现sql.Scanner接口
func (s *SimCashbacks) Scan(value interface{}) error {
	if value == nil {
		*s = make(SimCashbacks, 0)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into SimCashbacks", value)
	}

	if len(bytes) == 0 || string(bytes) == "[]" {
		*s = make(SimCashbacks, 0)
		return nil
	}

	return json.Unmarshal(bytes, s)
}

// Value 实现driver.Valuer接口
func (s SimCashbacks) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// GetCashbackAmount 获取指定档位的返现金额
func (s SimCashbacks) GetCashbackAmount(tierOrder int) int64 {
	for _, item := range s {
		if item.TierOrder == tierOrder {
			return item.CashbackAmount
		}
	}
	// 超出配置档位，查找最后一档（tier_order最大的）
	if len(s) > 0 {
		maxTier := s[0]
		for _, item := range s {
			if item.TierOrder > maxTier.TierOrder {
				maxTier = item
			}
		}
		return maxTier.CashbackAmount
	}
	return 0
}

// SetCashbackAmount 设置指定档位的返现金额
func (s *SimCashbacks) SetCashbackAmount(tierOrder int, amount int64) {
	for i, item := range *s {
		if item.TierOrder == tierOrder {
			(*s)[i].CashbackAmount = amount
			return
		}
	}
	// 不存在则添加
	*s = append(*s, SimCashbackItem{TierOrder: tierOrder, CashbackAmount: amount})
}

// ============================================================
// 通道完整配置（聚合响应）
// ============================================================

// ChannelFullConfig 通道完整配置
type ChannelFullConfig struct {
	ChannelID        int64                    `json:"channel_id"`
	ChannelCode      string                   `json:"channel_code"`
	ChannelName      string                   `json:"channel_name"`
	RateConfigs      []ChannelRateConfig      `json:"rate_configs"`
	DepositTiers     []ChannelDepositTier     `json:"deposit_tiers"`
	SimCashbackTiers []ChannelSimCashbackTier `json:"sim_cashback_tiers"`
}

// ============================================================
// API请求/响应DTO
// ============================================================

// CreateChannelRateConfigRequest 创建费率配置请求
type CreateChannelRateConfigRequest struct {
	RateCode    string `json:"rate_code" binding:"required"`
	RateName    string `json:"rate_name" binding:"required"`
	MinRate     string `json:"min_rate" binding:"required"`
	MaxRate     string `json:"max_rate" binding:"required"`
	DefaultRate string `json:"default_rate"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateChannelRateConfigRequest 更新费率配置请求
type UpdateChannelRateConfigRequest struct {
	RateName    string  `json:"rate_name"`
	MinRate     string  `json:"min_rate"`
	MaxRate     string  `json:"max_rate"`
	DefaultRate string  `json:"default_rate"`
	MaxHighRate *string `json:"max_high_rate"` // 高调费率上限
	MaxD0Extra  *int64  `json:"max_d0_extra"`  // P+0加价上限（分）
	SortOrder   int     `json:"sort_order"`
	Status      *int16  `json:"status"`
}

// UpdateChannelDepositTierRequest 更新押金档位请求
type UpdateChannelDepositTierRequest struct {
	MaxCashbackAmount int64  `json:"max_cashback_amount"`
	DefaultCashback   int64  `json:"default_cashback"`
	Status            *int16 `json:"status"`
}

// BatchSetSimCashbackTiersRequest 批量设置流量费返现档位请求
type BatchSetSimCashbackTiersRequest struct {
	Tiers []SimCashbackTierInput `json:"tiers" binding:"required"`
}

// SimCashbackTierInput 流量费返现档位输入
type SimCashbackTierInput struct {
	TierOrder         int    `json:"tier_order" binding:"required"`
	TierName          string `json:"tier_name" binding:"required"`
	IsLastTier        bool   `json:"is_last_tier"`
	MaxCashbackAmount int64  `json:"max_cashback_amount" binding:"required"`
	DefaultCashback   int64  `json:"default_cashback"`
	SimFeeAmount      int64  `json:"sim_fee_amount" binding:"required"`
}

// ChannelRateConfigResponse 费率配置响应
type ChannelRateConfigResponse struct {
	ID          int64   `json:"id"`
	ChannelID   int64   `json:"channel_id"`
	RateCode    string  `json:"rate_code"`
	RateName    string  `json:"rate_name"`
	MinRate     string  `json:"min_rate"`
	MaxRate     string  `json:"max_rate"`
	DefaultRate string  `json:"default_rate"`
	MaxHighRate *string `json:"max_high_rate"` // 高调费率上限
	MaxD0Extra  *int64  `json:"max_d0_extra"`  // P+0加价上限（分）
	SortOrder   int     `json:"sort_order"`
	Status      int16   `json:"status"`
}

// ChannelDepositTierResponse 押金档位响应
type ChannelDepositTierResponse struct {
	ID                int64  `json:"id"`
	ChannelID         int64  `json:"channel_id"`
	BrandCode         string `json:"brand_code"`
	TierCode          string `json:"tier_code"`
	DepositAmount     int64  `json:"deposit_amount"`
	TierName          string `json:"tier_name"`
	MaxCashbackAmount int64  `json:"max_cashback_amount"`
	DefaultCashback   int64  `json:"default_cashback"`
	SortOrder         int    `json:"sort_order"`
	Status            int16  `json:"status"`
}

// ChannelSimCashbackTierResponse 流量费返现档位响应
type ChannelSimCashbackTierResponse struct {
	ID                int64  `json:"id"`
	ChannelID         int64  `json:"channel_id"`
	BrandCode         string `json:"brand_code"`
	TierOrder         int    `json:"tier_order"`
	TierName          string `json:"tier_name"`
	IsLastTier        bool   `json:"is_last_tier"`
	MaxCashbackAmount int64  `json:"max_cashback_amount"`
	DefaultCashback   int64  `json:"default_cashback"`
	SimFeeAmount      int64  `json:"sim_fee_amount"`
	Status            int16  `json:"status"`
}

// ============================================================
// 校验相关函数
// ============================================================

// ValidateRateRange 校验费率是否在通道允许范围内
func ValidateRateRange(rate string, minRate string, maxRate string) error {
	// 转换为数值比较
	rateVal := ParseRateToFloat(rate)
	minVal := ParseRateToFloat(minRate)
	maxVal := ParseRateToFloat(maxRate)

	if rateVal < minVal {
		return fmt.Errorf("费率 %s 不能低于通道成本 %s", rate, minRate)
	}
	if rateVal > maxVal {
		return fmt.Errorf("费率 %s 不能超过通道上限 %s", rate, maxRate)
	}
	return nil
}

// ValidateSettlementRate 校验结算价费率（必须 >= 上级费率）
func ValidateSettlementRate(rate, upperRate, channelMinRate, channelMaxRate string) error {
	if err := ValidateRateRange(rate, channelMinRate, channelMaxRate); err != nil {
		return err
	}

	rateVal := ParseRateToFloat(rate)
	upperVal := ParseRateToFloat(upperRate)

	if rateVal < upperVal {
		return fmt.Errorf("费率 %s 不能低于上级费率 %s", rate, upperRate)
	}
	return nil
}

// ValidateCashbackAmount 校验返现金额是否在通道允许范围内
func ValidateCashbackAmount(cashback, maxCashback int64) error {
	if cashback > maxCashback {
		return fmt.Errorf("返现金额 %d 不能超过通道上限 %d", cashback, maxCashback)
	}
	if cashback < 0 {
		return fmt.Errorf("返现金额不能为负数")
	}
	return nil
}

// ValidateSettlementCashback 校验结算价返现（必须 <= 上级返现）
func ValidateSettlementCashback(cashback, upperCashback, channelMax int64) error {
	if err := ValidateCashbackAmount(cashback, channelMax); err != nil {
		return err
	}
	if cashback > upperCashback {
		return fmt.Errorf("返现金额 %d 不能超过上级配置 %d", cashback, upperCashback)
	}
	return nil
}

// ParseRateToFloat 将费率字符串转换为浮点数
func ParseRateToFloat(rate string) float64 {
	var val float64
	fmt.Sscanf(rate, "%f", &val)
	return val
}

// ValidateHighRate 校验高调费率是否在通道上限内
func ValidateHighRate(highRate string, maxHighRate *string, rateName string) error {
	if maxHighRate == nil || *maxHighRate == "" {
		return nil // 无上限限制
	}
	highRateVal := ParseRateToFloat(highRate)
	maxVal := ParseRateToFloat(*maxHighRate)
	if highRateVal > maxVal {
		return fmt.Errorf("%s高调费率 %s 超过通道上限 %s", rateName, highRate, *maxHighRate)
	}
	return nil
}

// ValidateD0Extra 校验P+0加价是否在通道上限内
func ValidateD0Extra(extraFee int64, maxD0Extra *int64, rateName string) error {
	if maxD0Extra == nil {
		return nil // 无上限限制
	}
	if extraFee > *maxD0Extra {
		return fmt.Errorf("%s P+0加价 %d分 超过通道上限 %d分", rateName, extraFee, *maxD0Extra)
	}
	return nil
}

// ============================================================
// 通道配置变更影响检查相关类型
// ============================================================

// ConfigChangeImpact 配置变更影响
type ConfigChangeImpact struct {
	AffectedTemplates   []AffectedTemplate   `json:"affected_templates"`
	AffectedSettlements []AffectedSettlement `json:"affected_settlements"`
	TotalAffectedAgents int                  `json:"total_affected_agents"`
}

// AffectedTemplate 受影响的政策模版
type AffectedTemplate struct {
	TemplateID   int64  `json:"template_id"`
	TemplateName string `json:"template_name"`
	Issue        string `json:"issue"`
}

// AffectedSettlement 受影响的结算价
type AffectedSettlement struct {
	SettlementID int64  `json:"settlement_id"`
	AgentID      int64  `json:"agent_id"`
	AgentName    string `json:"agent_name"`
	Issue        string `json:"issue"`
}

// CheckRateConfigChangeImpactRequest 检查费率配置变更影响请求
type CheckRateConfigChangeImpactRequest struct {
	NewMinRate     string  `json:"new_min_rate"`
	NewMaxRate     string  `json:"new_max_rate"`
	NewMaxHighRate *string `json:"new_max_high_rate"`
	NewMaxD0Extra  *int64  `json:"new_max_d0_extra"`
}

// ChannelConfigChangeLog 通道配置变更记录
type ChannelConfigChangeLog struct {
	ID                  int64     `json:"id" gorm:"primaryKey"`
	ChannelID           int64     `json:"channel_id" gorm:"not null;index"`
	ChangeType          string    `json:"change_type" gorm:"size:50;not null"`
	RateCode            string    `json:"rate_code" gorm:"size:32"`
	OldValue            JSONMap   `json:"old_value" gorm:"type:jsonb"`
	NewValue            JSONMap   `json:"new_value" gorm:"type:jsonb"`
	AffectedTemplates   int       `json:"affected_templates" gorm:"default:0"`
	AffectedSettlements int       `json:"affected_settlements" gorm:"default:0"`
	AffectedAgents      int       `json:"affected_agents" gorm:"default:0"`
	ImpactDetails       JSONMap   `json:"impact_details" gorm:"type:jsonb"`
	OperatorID          int64     `json:"operator_id" gorm:"not null"`
	OperatorName        string    `json:"operator_name" gorm:"size:100"`
	CreatedAt           time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (ChannelConfigChangeLog) TableName() string {
	return "channel_config_change_logs"
}
