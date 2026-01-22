package models

import (
	"time"
)

// 税筹通道扣费类型
const (
	TaxFeeTypePayment    int16 = 1 // 付款扣(充值时扣除)
	TaxFeeTypeWithdrawal int16 = 2 // 出款扣(提现时扣除)
)

// 税筹通道状态
const (
	TaxChannelStatusDisabled int16 = 0 // 禁用
	TaxChannelStatusEnabled  int16 = 1 // 启用
)

// TaxChannel 税筹通道
type TaxChannel struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ChannelCode string    `json:"channel_code" gorm:"column:channel_code;uniqueIndex;size:32"`
	ChannelName string    `json:"channel_name" gorm:"column:channel_name;size:100"`
	FeeType     int16     `json:"fee_type" gorm:"column:fee_type"`           // 1=付款扣 2=出款扣
	TaxRate     float64   `json:"tax_rate" gorm:"column:tax_rate;type:decimal(5,4)"` // 税率 如0.09表示9%
	FixedFee    int64     `json:"fixed_fee" gorm:"column:fixed_fee"`         // 固定费用(分)
	ApiURL      string    `json:"api_url" gorm:"column:api_url;size:255"`
	ApiKey      string    `json:"-" gorm:"column:api_key;size:255"`
	ApiSecret   string    `json:"-" gorm:"column:api_secret;size:255"`
	Status      int16     `json:"status" gorm:"column:status"`
	Remark      string    `json:"remark" gorm:"column:remark;size:500"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (TaxChannel) TableName() string {
	return "tax_channels"
}

// ChannelTaxMapping 支付通道与税筹通道关联
type ChannelTaxMapping struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	ChannelID    int64     `json:"channel_id" gorm:"column:channel_id"`
	WalletType   int16     `json:"wallet_type" gorm:"column:wallet_type"` // 钱包类型
	TaxChannelID int64     `json:"tax_channel_id" gorm:"column:tax_channel_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (ChannelTaxMapping) TableName() string {
	return "channel_tax_mappings"
}

// GetFeeTypeName 获取扣费类型名称
func GetFeeTypeName(feeType int16) string {
	switch feeType {
	case TaxFeeTypePayment:
		return "付款扣"
	case TaxFeeTypeWithdrawal:
		return "出款扣"
	default:
		return "未知"
	}
}

// GetTaxChannelStatusName 获取状态名称
func GetTaxChannelStatusName(status int16) string {
	switch status {
	case TaxChannelStatusEnabled:
		return "启用"
	case TaxChannelStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}
