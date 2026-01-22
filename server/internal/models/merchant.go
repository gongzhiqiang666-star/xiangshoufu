package models

import "time"

// Merchant 商户模型
type Merchant struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	MerchantNo    string     `json:"merchant_no" gorm:"size:64;uniqueIndex"`
	MerchantName  string     `json:"merchant_name" gorm:"size:100"`
	AgentID       int64      `json:"agent_id" gorm:"index"`
	ChannelID     int64      `json:"channel_id"`
	TerminalSN    string     `json:"terminal_sn" gorm:"size:50;index"`
	Status        int16      `json:"status" gorm:"default:1"`         // 1正常 2禁用
	ApproveStatus int16      `json:"approve_status" gorm:"default:1"` // 1待审核 2已通过 3已拒绝
	LegalName     string     `json:"legal_name" gorm:"size:50"`
	LegalIDCard   string     `json:"legal_id_card" gorm:"size:18"`
	MCC           string     `json:"mcc" gorm:"size:10"`
	CreditRate    string     `json:"credit_rate" gorm:"type:decimal(10,4)"`
	DebitRate     string     `json:"debit_rate" gorm:"type:decimal(10,4)"`
	// 新增字段
	MerchantType    string     `json:"merchant_type" gorm:"size:20;default:'normal';index"` // loyal/quality/potential/normal/low_active/inactive
	IsDirect        bool       `json:"is_direct" gorm:"default:true;index"`                 // true=直营 false=团队
	ActivatedAt     *time.Time `json:"activated_at"`                                        // 激活时间(首次交易时间)
	RegisteredPhone string     `json:"registered_phone" gorm:"size:100"`                    // 登记手机号(加密存储)
	RegisterRemark  string     `json:"register_remark" gorm:"size:500"`                     // 登记备注
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// 商户类型常量
const (
	MerchantTypeLoyal     = "loyal"      // 忠诚商户: 月均交易 > 5万
	MerchantTypeQuality   = "quality"    // 优质商户: 3万 ≤ 月均交易 < 5万
	MerchantTypePotential = "potential"  // 潜力商户: 2万 ≤ 月均交易 < 3万
	MerchantTypeNormal    = "normal"     // 一般商户: 1万 ≤ 月均交易 < 2万
	MerchantTypeLowActive = "low_active" // 低活跃: 0 < 月均交易 < 1万
	MerchantTypeInactive  = "inactive"   // 30天无交易
)

// 商户状态常量
const (
	MerchantStatusActive   int16 = 1 // 正常
	MerchantStatusDisabled int16 = 2 // 禁用
)

// 审核状态常量
const (
	MerchantApproveStatusPending  int16 = 1 // 待审核
	MerchantApproveStatusApproved int16 = 2 // 已通过
	MerchantApproveStatusRejected int16 = 3 // 已拒绝
)

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
