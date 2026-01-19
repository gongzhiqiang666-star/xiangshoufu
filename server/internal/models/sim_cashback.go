package models

import (
	"time"
)

// SimCashbackPolicy 流量费返现政策（三档：首次/2次/2+N次）
type SimCashbackPolicy struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	TemplateID         int64     `json:"template_id" gorm:"not null;index"`    // 政策模板ID
	ChannelID          int64     `json:"channel_id" gorm:"not null;index"`     // 通道ID
	BrandCode          string    `json:"brand_code" gorm:"size:32"`            // 品牌编码
	FirstTimeCashback  int64     `json:"first_time_cashback" gorm:"not null"`  // 首次返现金额（分）
	SecondTimeCashback int64     `json:"second_time_cashback" gorm:"not null"` // 第2次返现金额（分）
	ThirdPlusCashback  int64     `json:"third_plus_cashback" gorm:"not null"`  // 第3次及以后返现金额（分）
	SimFeeAmount       int64     `json:"sim_fee_amount" gorm:"not null"`       // 流量费金额（分）
	Status             int16     `json:"status" gorm:"default:1"`              // 1:启用 0:禁用
	CreatedAt          time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"default:now()"`
}

func (SimCashbackPolicy) TableName() string {
	return "sim_cashback_policies"
}

// SimCashbackRecord 流量费返现记录
type SimCashbackRecord struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	DeviceFeeID    int64      `json:"device_fee_id" gorm:"not null;index"`       // 关联流量费记录ID
	TerminalSN     string     `json:"terminal_sn" gorm:"size:50;not null;index"` // 终端SN
	ChannelID      int64      `json:"channel_id" gorm:"not null"`                // 通道ID
	AgentID        int64      `json:"agent_id" gorm:"not null;index"`            // 获得返现的代理商
	SimFeeCount    int        `json:"sim_fee_count" gorm:"not null"`             // 当前是第几次缴费
	SimFeeAmount   int64      `json:"sim_fee_amount" gorm:"not null"`            // 流量费金额（分）
	CashbackTier   int16      `json:"cashback_tier" gorm:"not null"`             // 返现档次 1:首次 2:第2次 3:第3次及以后
	SelfCashback   int64      `json:"self_cashback" gorm:"not null"`             // 自身返现金额（分）
	UpperCashback  int64      `json:"upper_cashback" gorm:"not null"`            // 上级应返金额（分）
	ActualCashback int64      `json:"actual_cashback" gorm:"not null"`           // 实际返现金额（级差）（分）
	SourceAgentID  int64      `json:"source_agent_id"`                           // 下级代理商ID（级差来源）
	WalletType     int16      `json:"wallet_type" gorm:"default:1"`              // 钱包类型
	WalletStatus   int16      `json:"wallet_status" gorm:"default:0"`            // 0:待入账 1:已入账
	CreatedAt      time.Time  `json:"created_at" gorm:"default:now()"`
	ProcessedAt    *time.Time `json:"processed_at"`
}

func (SimCashbackRecord) TableName() string {
	return "sim_cashback_records"
}

// SimCashbackTier 返现档次
const (
	SimCashbackTierFirst     = 1 // 首次
	SimCashbackTierSecond    = 2 // 第2次
	SimCashbackTierThirdPlus = 3 // 第3次及以后
)

// GetCashbackTier 根据缴费次数获取返现档次
func GetCashbackTier(simFeeCount int) int16 {
	switch {
	case simFeeCount == 1:
		return SimCashbackTierFirst
	case simFeeCount == 2:
		return SimCashbackTierSecond
	default:
		return SimCashbackTierThirdPlus
	}
}
