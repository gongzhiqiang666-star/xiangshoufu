package channel

import (
	"time"
)

// ActionType 回调动作类型
type ActionType string

const (
	ActionMerchantIncome ActionType = "merc_income"      // 商户入网
	ActionTerminalBind   ActionType = "sn_bind"          // 终端绑定/解绑
	ActionDeviceFee      ActionType = "sn_device_fee"    // 流量费/服务费
	ActionTransaction    ActionType = "pos_order"        // 交易回调
	ActionRateChange     ActionType = "merc_rate_update" // 费率修改
)

// ChannelAdapter 通道适配器接口
// 每个支付通道需要实现此接口
type ChannelAdapter interface {
	// GetChannelCode 获取通道编码
	GetChannelCode() string

	// GetChannelName 获取通道名称
	GetChannelName() string

	// VerifySign 验证签名
	// rawBody: 原始请求体
	// 返回: 是否验证通过, 错误信息
	VerifySign(rawBody []byte) (bool, error)

	// ParseActionType 解析回调类型
	// rawBody: 原始请求体
	// 返回: 动作类型, 错误信息
	ParseActionType(rawBody []byte) (ActionType, error)

	// ParseIdempotentKey 生成幂等键
	// rawBody: 原始请求体
	// 返回: 幂等键 (格式: channel_code:action_type:business_key)
	ParseIdempotentKey(rawBody []byte) (string, error)

	// ParseMerchantIncome 解析商户入网回调
	ParseMerchantIncome(rawBody []byte) (*UnifiedMerchantIncome, error)

	// ParseTerminalBind 解析终端绑定/解绑回调
	ParseTerminalBind(rawBody []byte) (*UnifiedTerminalBind, error)

	// ParseDeviceFee 解析流量费/服务费回调
	ParseDeviceFee(rawBody []byte) (*UnifiedDeviceFee, error)

	// ParseTransaction 解析交易回调
	ParseTransaction(rawBody []byte) (*UnifiedTransaction, error)

	// ParseRateChange 解析费率变更回调
	ParseRateChange(rawBody []byte) (*UnifiedRateChange, error)

	// UpdateMerchantRate 更新商户费率（主动调用通道API）
	// 返回: 通道流水号, 错误信息
	UpdateMerchantRate(req *RateUpdateRequest) (*RateUpdateResponse, error)

	// SupportsRateUpdate 是否支持费率实时更新
	SupportsRateUpdate() bool
}

// RateUpdateRequest 费率更新请求
type RateUpdateRequest struct {
	MerchantNo   string  `json:"merchant_no"`   // 商户号
	TerminalSN   string  `json:"terminal_sn"`   // 终端SN
	CreditRate   float64 `json:"credit_rate"`   // 贷记卡费率（如0.006表示0.6%）
	DebitRate    float64 `json:"debit_rate"`    // 借记卡费率
	DebitCap     int64   `json:"debit_cap"`     // 借记卡封顶（分）
	WechatRate   float64 `json:"wechat_rate"`   // 微信费率
	AlipayRate   float64 `json:"alipay_rate"`   // 支付宝费率
	UnionpayRate float64 `json:"unionpay_rate"` // 云闪付费率
}

// RateUpdateResponse 费率更新响应
type RateUpdateResponse struct {
	Success     bool   `json:"success"`      // 是否成功
	ChannelCode string `json:"channel_code"` // 通道编码
	TradeNo     string `json:"trade_no"`     // 通道流水号
	Message     string `json:"message"`      // 返回消息
}

// ========== 统一数据模型 ==========

// UnifiedMerchantIncome 统一商户入网数据
type UnifiedMerchantIncome struct {
	ChannelCode string `json:"channel_code"` // 通道编码
	BrandCode   string `json:"brand_code"`   // 品牌编号

	// 商户信息
	MerchantNo string `json:"merchant_no"` // 商户号
	TerminalSN string `json:"terminal_sn"` // 机具SN号
	AgentID    string `json:"agent_id"`    // 代理商ID

	// 审核状态
	ApproveStatus int `json:"approve_status"` // 1-审核中 2-审核通过 3-审核失败 4-商户停用
	BindStatus    int `json:"bind_status"`    // 0-未绑定 1-已绑定
	PosActive     int `json:"pos_active"`     // POS业务开通状态 1-成功 2-失败

	// 法人信息
	LegalName       string `json:"legal_name"`         // 法人姓名
	LegalIDCard     string `json:"legal_id_card"`      // 法人身份证（脱敏存储）
	IDCardStartDate string `json:"id_card_start_date"` // 身份证有效期开始
	IDCardEndDate   string `json:"id_card_end_date"`   // 身份证有效期结束

	// 结算信息
	SettleCardNo   string `json:"settle_card_no"`   // 结算卡号
	SettleBankName string `json:"settle_bank_name"` // 开户行

	// 经营信息
	DistrictCode string `json:"district_code"` // 经营地区
	Address      string `json:"address"`       // 详细地址
	MCC          string `json:"mcc"`           // MCC码

	// 费率信息（百分比，如0.60表示0.60%）
	CreditRate   string `json:"credit_rate"`   // 贷记卡费率
	DebitRate    string `json:"debit_rate"`    // 借记卡费率
	DebitCap     string `json:"debit_cap"`     // 借记卡封顶
	AlipayRate   string `json:"alipay_rate"`   // 支付宝费率
	WechatRate   string `json:"wechat_rate"`   // 微信费率
	UnionpayRate string `json:"unionpay_rate"` // 云闪付费率
	JDRate       string `json:"jd_rate"`       // 京东白条费率

	// 扩展字段
	ExtData map[string]interface{} `json:"ext_data"` // 通道特有字段
}

// UnifiedTerminalBind 统一终端绑定数据
type UnifiedTerminalBind struct {
	ChannelCode string `json:"channel_code"` // 通道编码
	BrandCode   string `json:"brand_code"`   // 品牌编号

	TerminalSN string `json:"terminal_sn"` // 机具SN号
	MerchantNo string `json:"merchant_no"` // 商户号
	AgentID    string `json:"agent_id"`    // 代理商ID

	BindStatus int `json:"bind_status"` // 1-绑定成功 2-解绑成功

	// 扩展字段
	ExtData map[string]interface{} `json:"ext_data"`
}

// UnifiedDeviceFee 统一流量费/服务费数据
type UnifiedDeviceFee struct {
	ChannelCode string `json:"channel_code"` // 通道编码
	BrandCode   string `json:"brand_code"`   // 品牌编号

	TerminalSN string `json:"terminal_sn"` // 机具SN号
	MerchantNo string `json:"merchant_no"` // 商户号
	AgentID    string `json:"agent_id"`    // 代理商ID

	OrderNo      string    `json:"order_no"`      // 订单号
	FeeType      int       `json:"fee_type"`      // 1-服务费 2-流量费/通讯费
	FeeAmount    int64     `json:"fee_amount"`    // 扣费金额（分）
	ChargingTime time.Time `json:"charging_time"` // 扣款时间

	// 扩展字段
	ExtData map[string]interface{} `json:"ext_data"`
}

// CardType 卡类型
type CardType string

const (
	CardTypeDebit    CardType = "debit"    // 借记卡
	CardTypeCredit   CardType = "credit"   // 贷记卡
	CardTypeWechat   CardType = "wechat"   // 微信
	CardTypeAlipay   CardType = "alipay"   // 支付宝
	CardTypeUnionpay CardType = "unionpay" // 银联云闪付
	CardTypeApplePay CardType = "applepay" // 苹果支付
)

// UnifiedTransaction 统一交易数据
type UnifiedTransaction struct {
	ChannelCode string `json:"channel_code"` // 通道编码
	BrandCode   string `json:"brand_code"`   // 品牌编号

	// 关联信息
	TerminalSN string `json:"terminal_sn"` // 机具SN号
	MerchantNo string `json:"merchant_no"` // 商户号
	AgentID    string `json:"agent_id"`    // 代理商ID

	// 交易信息
	OrderNo   string    `json:"order_no"`   // 订单号
	TransTime time.Time `json:"trans_time"` // 交易时间
	Amount    int64     `json:"amount"`     // 交易金额（分）
	CardType  CardType  `json:"card_type"`  // 卡类型
	CardNo    string    `json:"card_no"`    // 卡号（脱敏）

	// 费率信息
	FeeRate  string `json:"fee_rate"`  // 交易费率（%）
	D0Fee    int64  `json:"d0_fee"`    // D0手续费（分）
	HighRate string `json:"high_rate"` // 调价费率（%）

	// 扩展字段
	ExtData map[string]interface{} `json:"ext_data"`
}

// UnifiedRateChange 统一费率变更数据
type UnifiedRateChange struct {
	ChannelCode string `json:"channel_code"` // 通道编码
	BrandCode   string `json:"brand_code"`   // 品牌编号

	TerminalSN string `json:"terminal_sn"` // 机具SN号
	MerchantNo string `json:"merchant_no"` // 商户号
	AgentID    string `json:"agent_id"`    // 代理商ID

	// 费率信息（百分比）
	CreditRate      string `json:"credit_rate"`       // 贷记卡费率
	CreditExtraRate string `json:"credit_extra_rate"` // 贷记卡额外费率
	DebitRate       string `json:"debit_rate"`        // 借记卡费率
	AlipayRate      string `json:"alipay_rate"`       // 支付宝费率
	WechatRate      string `json:"wechat_rate"`       // 微信费率
	UnionpayRate    string `json:"unionpay_rate"`     // 银联费率

	// 调价费率
	CreditAdditionRate   string `json:"credit_addition_rate"`   // 贷记卡调价费率
	UnionpayAdditionRate string `json:"unionpay_addition_rate"` // 银联调价费率
	AlipayAdditionRate   string `json:"alipay_addition_rate"`   // 支付宝调价费率
	WechatAdditionRate   string `json:"wechat_addition_rate"`   // 微信调价费率

	// 扩展字段
	ExtData map[string]interface{} `json:"ext_data"`
}
