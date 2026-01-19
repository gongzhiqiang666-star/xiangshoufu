package hengxintong

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"xiangshoufu/internal/channel"
)

// Adapter 恒信通适配器
type Adapter struct {
	publicKey *rsa.PublicKey
	config    *channel.ChannelConfig
}

// NewAdapter 创建恒信通适配器
func NewAdapter(config *channel.ChannelConfig) (*Adapter, error) {
	adapter := &Adapter{
		config: config,
	}

	if config.PublicKey != "" {
		pubKey, err := parsePublicKey(config.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("parse public key failed: %w", err)
		}
		adapter.publicKey = pubKey
	}

	return adapter, nil
}

// parsePublicKey 解析RSA公钥
func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	// 如果不包含PEM头，添加
	if !strings.Contains(publicKeyPEM, "-----BEGIN") {
		publicKeyPEM = "-----BEGIN PUBLIC KEY-----\n" + publicKeyPEM + "\n-----END PUBLIC KEY-----"
	}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// 尝试解析为PKCS1格式
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

// GetChannelCode 获取通道编码
func (a *Adapter) GetChannelCode() string {
	return channel.ChannelCodeHengxintong
}

// GetChannelName 获取通道名称
func (a *Adapter) GetChannelName() string {
	return "恒信通"
}

// VerifySign 验证签名
func (a *Adapter) VerifySign(rawBody []byte) (bool, error) {
	if a.publicKey == nil {
		// 如果没有配置公钥，跳过验签（开发环境）
		return true, nil
	}

	// 解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal(rawBody, &data); err != nil {
		return false, fmt.Errorf("parse json failed: %w", err)
	}

	// 获取sign字段
	signStr, ok := data["sign"].(string)
	if !ok || signStr == "" {
		return false, errors.New("sign field not found")
	}

	// 移除sign字段
	delete(data, "sign")

	// 构建待签名字符串：按key字典序排列
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		v := data[k]
		if v == nil {
			continue
		}
		// 转换为字符串
		var strVal string
		switch val := v.(type) {
		case string:
			strVal = val
		case float64:
			strVal = strconv.FormatFloat(val, 'f', -1, 64)
		default:
			strVal = fmt.Sprintf("%v", val)
		}
		if strVal != "" {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, strVal))
		}
	}

	signContent := strings.Join(pairs, "&")

	// 计算SHA256哈希
	hash := sha256.Sum256([]byte(signContent))

	// Base64解码签名
	signature, err := base64.StdEncoding.DecodeString(signStr)
	if err != nil {
		return false, fmt.Errorf("decode sign failed: %w", err)
	}

	// RSA验签
	err = rsa.VerifyPKCS1v15(a.publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return false, nil // 验签失败，不返回错误
	}

	return true, nil
}

// ParseActionType 解析回调类型
func (a *Adapter) ParseActionType(rawBody []byte) (channel.ActionType, error) {
	var base BaseRequest
	if err := json.Unmarshal(rawBody, &base); err != nil {
		return "", fmt.Errorf("parse action type failed: %w", err)
	}

	switch base.Action {
	case "merc_income":
		return channel.ActionMerchantIncome, nil
	case "sn_bind":
		return channel.ActionTerminalBind, nil
	case "sn_device_fee":
		return channel.ActionDeviceFee, nil
	case "pos_order":
		return channel.ActionTransaction, nil
	case "merc_rate_update":
		return channel.ActionRateChange, nil
	default:
		return "", fmt.Errorf("unknown action type: %s", base.Action)
	}
}

// ParseIdempotentKey 生成幂等键
func (a *Adapter) ParseIdempotentKey(rawBody []byte) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(rawBody, &data); err != nil {
		return "", fmt.Errorf("parse json failed: %w", err)
	}

	action, _ := data["action"].(string)
	channelCode := a.GetChannelCode()

	// 根据不同的action类型，选择不同的业务唯一键
	var bizKey string
	switch action {
	case "merc_income":
		// 商户入网：商户号 + 审核状态
		merchantNo, _ := data["merchantNo"].(string)
		approveStatus, _ := data["approveStatus"].(string)
		bizKey = merchantNo + "_" + approveStatus
	case "sn_bind":
		// 终端绑定：机具号 + 商户号 + 状态
		tusn, _ := data["tusn"].(string)
		merchantNo, _ := data["merchantNo"].(string)
		status, _ := data["status"].(string)
		bizKey = tusn + "_" + merchantNo + "_" + status
	case "sn_device_fee":
		// 流量费：订单号
		orderNo, _ := data["orderNo"].(string)
		bizKey = orderNo
	case "pos_order":
		// 交易：订单号
		orderNo, _ := data["orderNo"].(string)
		bizKey = orderNo
	case "merc_rate_update":
		// 费率变更：商户号 + 时间戳（取前10位）
		merchantNo, _ := data["merchantNo"].(string)
		// 使用商户号作为业务键，允许同一商户多次费率变更
		bizKey = merchantNo
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	return fmt.Sprintf("%s:%s:%s", channelCode, action, bizKey), nil
}

// ParseMerchantIncome 解析商户入网回调
func (a *Adapter) ParseMerchantIncome(rawBody []byte) (*channel.UnifiedMerchantIncome, error) {
	var req MerchantIncomeRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return nil, fmt.Errorf("parse merchant income failed: %w", err)
	}

	approveStatus, _ := strconv.Atoi(req.ApproveStatus)
	bindStatus, _ := strconv.Atoi(req.BindStatus)
	posActive, _ := strconv.Atoi(req.PosBusinessActive)

	return &channel.UnifiedMerchantIncome{
		ChannelCode:     a.GetChannelCode(),
		BrandCode:       req.BrandCode,
		MerchantNo:      req.MerchantNo,
		TerminalSN:      req.Tusn,
		ApproveStatus:   approveStatus,
		BindStatus:      bindStatus,
		PosActive:       posActive,
		LegalName:       req.LegalName,
		LegalIDCard:     maskIDCard(req.LegalNo),
		IDCardStartDate: req.IdCardStartDate,
		IDCardEndDate:   req.IdCardEndDate,
		SettleCardNo:    maskBankCard(req.SettleCardNo),
		SettleBankName:  req.SettleBankName,
		DistrictCode:    req.DistrictCode,
		Address:         req.Address,
		MCC:             req.Mcc,
		CreditRate:      req.CreditCardFeeRate,
		DebitRate:       req.DebitCardFeeRate,
		DebitCap:        req.DebitCardFeeMax,
		AlipayRate:      req.AlipayFeeRate,
		WechatRate:      req.WxPayFeeRate,
		UnionpayRate:    req.CloudCreditFeeRate,
		JDRate:          req.JdPayFeeRate,
		ExtData: map[string]interface{}{
			"idCardFrontUrl": req.IdCardFrontUrl,
			"idCardBackUrl":  req.IdCardBackUrl,
			"settleCardUrl":  req.SettleCardUrl,
		},
	}, nil
}

// ParseTerminalBind 解析终端绑定/解绑回调
func (a *Adapter) ParseTerminalBind(rawBody []byte) (*channel.UnifiedTerminalBind, error) {
	var req TerminalBindRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return nil, fmt.Errorf("parse terminal bind failed: %w", err)
	}

	bindStatus, _ := strconv.Atoi(req.Status)

	return &channel.UnifiedTerminalBind{
		ChannelCode: a.GetChannelCode(),
		BrandCode:   req.BrandCode,
		TerminalSN:  req.Tusn,
		MerchantNo:  req.MerchantNo,
		AgentID:     req.AgentId,
		BindStatus:  bindStatus,
	}, nil
}

// ParseDeviceFee 解析流量费/服务费回调
func (a *Adapter) ParseDeviceFee(rawBody []byte) (*channel.UnifiedDeviceFee, error) {
	var req DeviceFeeRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return nil, fmt.Errorf("parse device fee failed: %w", err)
	}

	feeType, _ := strconv.Atoi(req.Type)
	feeAmount, _ := strconv.ParseInt(req.ChargingAmount, 10, 64)
	chargingTime, _ := ParseTime(req.ChargingTime)

	return &channel.UnifiedDeviceFee{
		ChannelCode:  a.GetChannelCode(),
		BrandCode:    req.BrandCode,
		TerminalSN:   req.Tusn,
		MerchantNo:   req.MerchantNo,
		AgentID:      req.AgentId,
		OrderNo:      req.OrderNo,
		FeeType:      feeType,
		FeeAmount:    feeAmount,
		ChargingTime: chargingTime,
	}, nil
}

// ParseTransaction 解析交易回调
func (a *Adapter) ParseTransaction(rawBody []byte) (*channel.UnifiedTransaction, error) {
	var req TransactionRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return nil, fmt.Errorf("parse transaction failed: %w", err)
	}

	amount, _ := strconv.ParseInt(req.Amount, 10, 64)
	d0Fee, _ := strconv.ParseInt(req.FeeExt, 10, 64)
	transTime, _ := ParseTime(req.TransTime)

	// 映射卡类型
	cardType := mapCardType(req.TransCardType)

	return &channel.UnifiedTransaction{
		ChannelCode: a.GetChannelCode(),
		BrandCode:   req.BrandCode,
		TerminalSN:  req.Tusn,
		MerchantNo:  req.MerchantNo,
		AgentID:     req.AgentId,
		OrderNo:     req.OrderNo,
		TransTime:   transTime,
		Amount:      amount,
		CardType:    cardType,
		CardNo:      maskBankCard(req.CardNo),
		FeeRate:     req.TransactionFee,
		D0Fee:       d0Fee,
		HighRate:    req.HighRate,
	}, nil
}

// ParseRateChange 解析费率变更回调
func (a *Adapter) ParseRateChange(rawBody []byte) (*channel.UnifiedRateChange, error) {
	var req RateChangeRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return nil, fmt.Errorf("parse rate change failed: %w", err)
	}

	return &channel.UnifiedRateChange{
		ChannelCode:          a.GetChannelCode(),
		BrandCode:            req.BrandCode,
		TerminalSN:           req.Tusn,
		MerchantNo:           req.MerchantNo,
		AgentID:              req.AgentId,
		CreditRate:           req.CreditCardFeeRate,
		CreditExtraRate:      req.CreditCardExtraRate,
		DebitRate:            req.DebitCardFeeRate,
		AlipayRate:           req.AlipayFeeRate,
		WechatRate:           req.WxPayFeeRate,
		UnionpayRate:         req.UnionpayPayFeeRate,
		CreditAdditionRate:   req.CreditAdditionFeeRate,
		UnionpayAdditionRate: req.UnionpayAdditionFeeRate,
		AlipayAdditionRate:   req.AlipayAdditionFeeRate,
		WechatAdditionRate:   req.WxAdditionFeeRate,
	}, nil
}

// mapCardType 映射恒信通卡类型到统一类型
func mapCardType(transCardType string) channel.CardType {
	switch transCardType {
	case "00":
		return channel.CardTypeDebit
	case "01":
		return channel.CardTypeCredit
	case "061":
		return channel.CardTypeWechat
	case "062":
		return channel.CardTypeAlipay
	case "063":
		return channel.CardTypeUnionpay
	case "065":
		return channel.CardTypeApplePay
	default:
		return channel.CardTypeDebit // 默认借记卡
	}
}

// maskIDCard 脱敏身份证号（保留前6位和后4位）
func maskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return idCard
	}
	return idCard[:6] + "********" + idCard[len(idCard)-4:]
}

// maskBankCard 脱敏银行卡号（保留前6位和后4位）
func maskBankCard(cardNo string) string {
	if len(cardNo) < 10 {
		return cardNo
	}
	maskLen := len(cardNo) - 10
	mask := strings.Repeat("*", maskLen)
	return cardNo[:6] + mask + cardNo[len(cardNo)-4:]
}

// Configure 配置适配器
func (a *Adapter) Configure(config *channel.ChannelConfig) error {
	a.config = config
	if config.PublicKey != "" {
		pubKey, err := parsePublicKey(config.PublicKey)
		if err != nil {
			return fmt.Errorf("parse public key failed: %w", err)
		}
		a.publicKey = pubKey
	}
	return nil
}

// 确保Adapter实现了ChannelAdapter接口
var _ channel.ChannelAdapter = (*Adapter)(nil)
var _ channel.ConfigurableAdapter = (*Adapter)(nil)
