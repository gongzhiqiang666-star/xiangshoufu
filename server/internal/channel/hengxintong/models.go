package hengxintong

import "time"

// 恒信通回调请求模型

// BaseRequest 基础请求字段
type BaseRequest struct {
	Action    string `json:"action"`    // 动作类型
	Sign      string `json:"sign"`      // 签名
	BrandCode string `json:"brandCode"` // 品牌编号
	Tusn      string `json:"tusn"`      // 机具号
}

// MerchantIncomeRequest 商户入网回调请求
type MerchantIncomeRequest struct {
	BaseRequest

	MerchantNo        string `json:"merchantNo"`        // 商户号
	ApproveStatus     string `json:"approveStatus"`     // 审核状态 1-审核中 2-审核通过 3-审核失败 4-商户停用
	BindStatus        string `json:"bindStatus"`        // 绑定状态 0-未绑定 1-已绑定
	PosBusinessActive string `json:"posBusinessActive"` // POS业务开通状态 1-成功 2-失败

	// 法人信息
	LegalName       string `json:"legalName"`       // 法人姓名
	LegalNo         string `json:"legalNo"`         // 法人身份证
	IdCardStartDate string `json:"idCardStartDate"` // 身份证有效期开始
	IdCardEndDate   string `json:"idCardEndDate"`   // 身份证有效期结束
	IdCardFrontUrl  string `json:"idCardFrontUrl"`  // 身份证正面
	IdCardBackUrl   string `json:"idCardBackUrl"`   // 身份证反面

	// 结算信息
	SettleCardNo   string `json:"settleCardNo"`   // 结算卡号
	SettleBankName string `json:"settleBankName"` // 开户行
	SettleCardUrl  string `json:"settleCardUrl"`  // 结算卡照片

	// 经营信息
	DistrictCode string `json:"districtCode"` // 经营地区
	Address      string `json:"address"`      // 详细地址
	Mcc          string `json:"mcc"`          // MCC码

	// 费率信息
	CreditCardFeeRate  string `json:"creditCardFeeRate"`  // 贷记卡费率
	DebitCardFeeRate   string `json:"debitCardFeeRate"`   // 借记卡费率
	DebitCardFeeMax    string `json:"debitCardFeeMax"`    // 借记卡封顶
	AlipayFeeRate      string `json:"alipayFeeRate"`      // 支付宝费率
	CloudCreditFeeRate string `json:"cloudCreditFeeRate"` // 云闪付贷记卡费率
	WxPayFeeRate       string `json:"wxPayFeeRate"`       // 微信费率
	JdPayFeeRate       string `json:"jdPayFeeRate"`       // 京东白条费率
}

// TerminalBindRequest 终端绑定/解绑回调请求
type TerminalBindRequest struct {
	BaseRequest

	AgentId    string `json:"agentId"`    // 代理商ID
	Status     string `json:"status"`     // 绑定状态 1-绑定成功 2-解绑成功
	MerchantNo string `json:"merchantNo"` // 商户号
}

// DeviceFeeRequest 流量费/服务费扣费回调请求
type DeviceFeeRequest struct {
	BaseRequest

	MerchantNo     string `json:"merchantNo"`     // 商户号
	AgentId        string `json:"agentId"`        // 代理商ID
	ChargingAmount string `json:"chargingAmount"` // 扣费金额（分）
	Type           string `json:"type"`           // 1-服务费 2-通讯费
	OrderNo        string `json:"orderNo"`        // 外部订单号
	ChargingTime   string `json:"chargingTime"`   // 扣款时间 yyyy-MM-dd HH:mm:ss
}

// TransactionRequest 交易回调请求
type TransactionRequest struct {
	BaseRequest

	TransTime      string `json:"transTime"`      // 交易时间 yyyy-MM-dd HH:mm:ss
	OrderNo        string `json:"orderNo"`        // 订单号
	TransCardType  string `json:"transCardType"`  // 卡类型 00-借记卡 01-贷记卡 061-微信 062-支付宝 063-银联 065-苹果
	CardNo         string `json:"cardNo"`         // 卡号
	Amount         string `json:"amount"`         // 交易金额（分）
	TransactionFee string `json:"transactionFee"` // 交易费率（%）
	FeeExt         string `json:"feeExt"`         // D0手续费（分）
	MerchantNo     string `json:"merchantNo"`     // 商户号
	AgentId        string `json:"agentId"`        // 代理商ID
	HighRate       string `json:"highRate"`       // 调价费率（%）
}

// RateChangeRequest 费率修改回调请求
type RateChangeRequest struct {
	BaseRequest

	MerchantNo              string `json:"merchantNo"`              // 商户号
	AgentId                 string `json:"agentId"`                 // 代理商ID
	CreditCardFeeRate       string `json:"creditCardFeeRate"`       // 贷记卡费率
	CreditCardExtraRate     string `json:"creditCardExtraRate"`     // 贷记卡额外费率
	DebitCardFeeRate        string `json:"debitCardFeeRate"`        // 借记卡费率
	AlipayFeeRate           string `json:"alipayFeeRate"`           // 支付宝费率
	UnionpayPayFeeRate      string `json:"unionpayPayFeeRate"`      // 银联费率
	WxPayFeeRate            string `json:"wxPayFeeRate"`            // 微信费率
	CreditAdditionFeeRate   string `json:"creditAdditionFeeRate"`   // 贷记卡调价费率
	UnionpayAdditionFeeRate string `json:"unionpayAdditionFeeRate"` // 银联调价费率
	AlipayAdditionFeeRate   string `json:"alipayAdditionFeeRate"`   // 支付宝调价费率
	WxAdditionFeeRate       string `json:"wxAdditionFeeRate"`       // 微信调价费率
}

// ParseTime 解析恒信通时间格式
func ParseTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
}
