package hengxintong

import (
	"encoding/json"
	"testing"

	"xiangshoufu/internal/channel"
)

func TestParseActionType(t *testing.T) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	tests := []struct {
		name     string
		input    string
		expected channel.ActionType
		wantErr  bool
	}{
		{
			name:     "merchant income",
			input:    `{"action":"merc_income","sign":"xxx"}`,
			expected: channel.ActionMerchantIncome,
			wantErr:  false,
		},
		{
			name:     "terminal bind",
			input:    `{"action":"sn_bind","sign":"xxx"}`,
			expected: channel.ActionTerminalBind,
			wantErr:  false,
		},
		{
			name:     "device fee",
			input:    `{"action":"sn_device_fee","sign":"xxx"}`,
			expected: channel.ActionDeviceFee,
			wantErr:  false,
		},
		{
			name:     "transaction",
			input:    `{"action":"pos_order","sign":"xxx"}`,
			expected: channel.ActionTransaction,
			wantErr:  false,
		},
		{
			name:     "rate change",
			input:    `{"action":"merc_rate_update","sign":"xxx"}`,
			expected: channel.ActionRateChange,
			wantErr:  false,
		},
		{
			name:     "unknown action",
			input:    `{"action":"unknown","sign":"xxx"}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.ParseActionType([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %s but got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestParseTransaction(t *testing.T) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	input := `{
		"action": "pos_order",
		"sign": "xxx",
		"brandCode": "HXT001",
		"tusn": "SN12345678",
		"transTime": "2024-01-15 10:30:00",
		"orderNo": "ORDER123456",
		"transCardType": "01",
		"cardNo": "6228480402564890018",
		"amount": "10000",
		"transactionFee": "0.60",
		"feeExt": "300",
		"merchantNo": "M12345678",
		"agentId": "A001",
		"highRate": "0.05"
	}`

	result, err := adapter.ParseTransaction([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证字段映射
	if result.ChannelCode != "HENGXINTONG" {
		t.Errorf("expected channel code HENGXINTONG but got %s", result.ChannelCode)
	}
	if result.BrandCode != "HXT001" {
		t.Errorf("expected brand code HXT001 but got %s", result.BrandCode)
	}
	if result.TerminalSN != "SN12345678" {
		t.Errorf("expected terminal SN SN12345678 but got %s", result.TerminalSN)
	}
	if result.OrderNo != "ORDER123456" {
		t.Errorf("expected order no ORDER123456 but got %s", result.OrderNo)
	}
	if result.Amount != 10000 {
		t.Errorf("expected amount 10000 but got %d", result.Amount)
	}
	if result.CardType != channel.CardTypeCredit {
		t.Errorf("expected card type credit but got %s", result.CardType)
	}
	if result.FeeRate != "0.60" {
		t.Errorf("expected fee rate 0.60 but got %s", result.FeeRate)
	}
	if result.D0Fee != 300 {
		t.Errorf("expected D0 fee 300 but got %d", result.D0Fee)
	}
	if result.HighRate != "0.05" {
		t.Errorf("expected high rate 0.05 but got %s", result.HighRate)
	}
	// 验证卡号脱敏
	if result.CardNo == "6228480402564890018" {
		t.Error("card number should be masked")
	}
}

func TestParseDeviceFee(t *testing.T) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	input := `{
		"action": "sn_device_fee",
		"sign": "xxx",
		"brandCode": "HXT001",
		"tusn": "SN12345678",
		"merchantNo": "M12345678",
		"agentId": "A001",
		"chargingAmount": "3600",
		"type": "2",
		"orderNo": "FEE123456",
		"chargingTime": "2024-01-15 10:30:00"
	}`

	result, err := adapter.ParseDeviceFee([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OrderNo != "FEE123456" {
		t.Errorf("expected order no FEE123456 but got %s", result.OrderNo)
	}
	if result.FeeType != 2 {
		t.Errorf("expected fee type 2 but got %d", result.FeeType)
	}
	if result.FeeAmount != 3600 {
		t.Errorf("expected fee amount 3600 but got %d", result.FeeAmount)
	}
}

func TestParseIdempotentKey(t *testing.T) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "transaction",
			input:    `{"action":"pos_order","orderNo":"ORDER123"}`,
			contains: "HENGXINTONG:pos_order:ORDER123",
		},
		{
			name:     "device fee",
			input:    `{"action":"sn_device_fee","orderNo":"FEE456"}`,
			contains: "HENGXINTONG:sn_device_fee:FEE456",
		},
		{
			name:     "terminal bind",
			input:    `{"action":"sn_bind","tusn":"SN123","merchantNo":"M456","status":"1"}`,
			contains: "HENGXINTONG:sn_bind:SN123_M456_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.ParseIdempotentKey([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.contains {
				t.Errorf("expected %s but got %s", tt.contains, result)
			}
		})
	}
}

func TestParseRateChange(t *testing.T) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	input := `{
		"action": "merc_rate_update",
		"sign": "xxx",
		"brandCode": "HXT001",
		"tusn": "SN12345678",
		"merchantNo": "M12345678",
		"agentId": "A001",
		"creditCardFeeRate": "0.60",
		"debitCardFeeRate": "0.50",
		"alipayFeeRate": "0.38",
		"wxPayFeeRate": "0.38",
		"unionpayPayFeeRate": "0.38",
		"creditAdditionFeeRate": "0.05"
	}`

	result, err := adapter.ParseRateChange([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.CreditRate != "0.60" {
		t.Errorf("expected credit rate 0.60 but got %s", result.CreditRate)
	}
	if result.DebitRate != "0.50" {
		t.Errorf("expected debit rate 0.50 but got %s", result.DebitRate)
	}
	if result.AlipayRate != "0.38" {
		t.Errorf("expected alipay rate 0.38 but got %s", result.AlipayRate)
	}
	if result.CreditAdditionRate != "0.05" {
		t.Errorf("expected credit addition rate 0.05 but got %s", result.CreditAdditionRate)
	}
}

func TestMapCardType(t *testing.T) {
	tests := []struct {
		input    string
		expected channel.CardType
	}{
		{"00", channel.CardTypeDebit},
		{"01", channel.CardTypeCredit},
		{"061", channel.CardTypeWechat},
		{"062", channel.CardTypeAlipay},
		{"063", channel.CardTypeUnionpay},
		{"065", channel.CardTypeApplePay},
		{"unknown", channel.CardTypeDebit}, // 默认值
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapCardType(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, result)
			}
		})
	}
}

func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"110101199001011234", "110101********1234"},
		{"12345", "12345"}, // 短于10位不脱敏
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := maskIDCard(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, result)
			}
		})
	}
}

func TestMaskBankCard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"6228480402564890018", "622848*********0018"},
		{"12345", "12345"}, // 短于10位不脱敏
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := maskBankCard(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, result)
			}
		})
	}
}

func TestParseMerchantIncome(t *testing.T) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	input := `{
		"action": "merc_income",
		"sign": "xxx",
		"brandCode": "HXT001",
		"tusn": "SN12345678",
		"merchantNo": "M12345678",
		"approveStatus": "2",
		"bindStatus": "1",
		"posBusinessActive": "1",
		"legalName": "张三",
		"legalNo": "110101199001011234",
		"settleCardNo": "6228480402564890018",
		"settleBankName": "中国银行",
		"districtCode": "110101",
		"address": "北京市东城区xxx街道",
		"mcc": "5411",
		"creditCardFeeRate": "0.60",
		"debitCardFeeRate": "0.50"
	}`

	result, err := adapter.ParseMerchantIncome([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.MerchantNo != "M12345678" {
		t.Errorf("expected merchant no M12345678 but got %s", result.MerchantNo)
	}
	if result.ApproveStatus != 2 {
		t.Errorf("expected approve status 2 but got %d", result.ApproveStatus)
	}
	if result.LegalName != "张三" {
		t.Errorf("expected legal name 张三 but got %s", result.LegalName)
	}
	// 验证身份证脱敏
	if result.LegalIDCard == "110101199001011234" {
		t.Error("ID card should be masked")
	}
	// 验证银行卡脱敏
	if result.SettleCardNo == "6228480402564890018" {
		t.Error("bank card should be masked")
	}
}

func TestAdapterImplementsInterface(t *testing.T) {
	// 编译时检查接口实现
	var _ channel.ChannelAdapter = (*Adapter)(nil)
	var _ channel.ConfigurableAdapter = (*Adapter)(nil)
}

// BenchmarkParseTransaction 性能测试
func BenchmarkParseTransaction(b *testing.B) {
	adapter, _ := NewAdapter(&channel.ChannelConfig{})
	input := []byte(`{
		"action": "pos_order",
		"sign": "xxx",
		"brandCode": "HXT001",
		"tusn": "SN12345678",
		"transTime": "2024-01-15 10:30:00",
		"orderNo": "ORDER123456",
		"transCardType": "01",
		"cardNo": "6228480402564890018",
		"amount": "10000",
		"transactionFee": "0.60",
		"feeExt": "300",
		"merchantNo": "M12345678",
		"agentId": "A001",
		"highRate": "0.05"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.ParseTransaction(input)
	}
}

// TestVerifySign 签名验证测试（需要提供测试用公钥）
func TestVerifySign(t *testing.T) {
	// 无公钥配置时应返回true
	adapter, _ := NewAdapter(&channel.ChannelConfig{})

	input := `{"action":"pos_order","orderNo":"123","sign":"xxx"}`
	verified, err := adapter.VerifySign([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !verified {
		t.Error("expected verified=true when no public key configured")
	}
}

// TestJSONMarshal 测试统一数据模型的JSON序列化
func TestJSONMarshal(t *testing.T) {
	tx := &channel.UnifiedTransaction{
		ChannelCode: "HENGXINTONG",
		OrderNo:     "ORDER123",
		Amount:      10000,
		CardType:    channel.CardTypeCredit,
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result channel.UnifiedTransaction
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if result.ChannelCode != tx.ChannelCode {
		t.Errorf("expected %s but got %s", tx.ChannelCode, result.ChannelCode)
	}
	if result.Amount != tx.Amount {
		t.Errorf("expected %d but got %d", tx.Amount, result.Amount)
	}
}
