package models

import "time"

// Merchant 商户模型
type Merchant struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	MerchantNo    string    `json:"merchant_no" gorm:"size:64;uniqueIndex"`
	MerchantName  string    `json:"merchant_name" gorm:"size:100"`
	AgentID       int64     `json:"agent_id" gorm:"index"`
	ChannelID     int64     `json:"channel_id"`
	TerminalSN    string    `json:"terminal_sn" gorm:"size:50;index"`
	Status        int16     `json:"status" gorm:"default:1"`        // 1正常 2禁用
	ApproveStatus int16     `json:"approve_status" gorm:"default:1"` // 1待审核 2已通过 3已拒绝
	LegalName     string    `json:"legal_name" gorm:"size:50"`
	LegalIDCard   string    `json:"legal_id_card" gorm:"size:18"`
	MCC           string    `json:"mcc" gorm:"size:10"`
	CreditRate    string    `json:"credit_rate" gorm:"type:decimal(10,4)"`
	DebitRate     string    `json:"debit_rate" gorm:"type:decimal(10,4)"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 表名
func (Merchant) TableName() string {
	return "merchants"
}

// PolicyTemplate 政策模板模型
type PolicyTemplate struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	TemplateName string    `json:"template_name" gorm:"size:100"`
	ChannelID    int64     `json:"channel_id"`
	IsDefault    bool      `json:"is_default" gorm:"default:false"`
	CreditRate   string    `json:"credit_rate" gorm:"type:decimal(10,4)"`
	DebitRate    string    `json:"debit_rate" gorm:"type:decimal(10,4)"`
	DebitCap     string    `json:"debit_cap" gorm:"type:decimal(10,2)"`
	UnionpayRate string    `json:"unionpay_rate" gorm:"type:decimal(10,4)"`
	WechatRate   string    `json:"wechat_rate" gorm:"type:decimal(10,4)"`
	AlipayRate   string    `json:"alipay_rate" gorm:"type:decimal(10,4)"`
	Status       int16     `json:"status" gorm:"default:1"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 表名
func (PolicyTemplate) TableName() string {
	return "policy_templates"
}
