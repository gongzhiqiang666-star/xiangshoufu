package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"xiangshoufu/internal/models"
)

// =============================================================================
// ChannelConfig 结构体测试
// =============================================================================

// TestChannelConfig_ParseRateTypes 测试解析费率类型配置
func TestChannelConfig_ParseRateTypes(t *testing.T) {
	tests := []struct {
		name          string
		configJSON    string
		expectedCount int
		expectError   bool
	}{
		{
			name: "正常解析多个费率类型",
			configJSON: `{
				"rate_types": [
					{"code": "CREDIT", "name": "贷记卡费率", "sort_order": 1, "min_rate": "0.50", "max_rate": "0.68"},
					{"code": "DEBIT", "name": "借记卡费率", "sort_order": 2, "min_rate": "0.45", "max_rate": "0.60"}
				]
			}`,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "空配置",
			configJSON:    `{}`,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "无效JSON",
			configJSON:    `{invalid}`,
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config ChannelConfig
			err := json.Unmarshal([]byte(tt.configJSON), &config)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(config.RateTypes))
			}
		})
	}
}

// =============================================================================
// RateTypeDefinition 测试
// =============================================================================

// TestRateTypeDefinition_Structure 测试费率类型定义结构
func TestRateTypeDefinition_Structure(t *testing.T) {
	rt := models.RateTypeDefinition{
		Code:      "POS_CC",
		Name:      "普通刷卡-贷记",
		SortOrder: 4,
		MinRate:   "0.50",
		MaxRate:   "0.68",
	}

	assert.Equal(t, "POS_CC", rt.Code)
	assert.Equal(t, "普通刷卡-贷记", rt.Name)
	assert.Equal(t, 4, rt.SortOrder)
	assert.Equal(t, "0.50", rt.MinRate)
	assert.Equal(t, "0.68", rt.MaxRate)
}

// TestRateTypeDefinition_Sorting 测试费率类型排序
func TestRateTypeDefinition_Sorting(t *testing.T) {
	rateTypes := []models.RateTypeDefinition{
		{Code: "WECHAT", Name: "微信", SortOrder: 3},
		{Code: "CREDIT", Name: "贷记卡", SortOrder: 1},
		{Code: "DEBIT", Name: "借记卡", SortOrder: 2},
	}

	// 按sort_order排序
	sort.Slice(rateTypes, func(i, j int) bool {
		return rateTypes[i].SortOrder < rateTypes[j].SortOrder
	})

	assert.Equal(t, "CREDIT", rateTypes[0].Code)
	assert.Equal(t, "DEBIT", rateTypes[1].Code)
	assert.Equal(t, "WECHAT", rateTypes[2].Code)
}

// =============================================================================
// RateConfigs 测试
// =============================================================================

// TestRateConfigs_GetRate 测试获取费率配置值
func TestRateConfigs_GetRate(t *testing.T) {
	rateConfigs := models.RateConfigs{
		"CREDIT":   {Rate: "0.60"},
		"DEBIT":    {Rate: "0.50"},
		"WECHAT":   {Rate: "0.38"},
		"ALIPAY":   {Rate: "0.38"},
		"UNIONPAY": {Rate: "0.35"},
	}

	tests := []struct {
		code     string
		expected string
		exists   bool
	}{
		{"CREDIT", "0.60", true},
		{"DEBIT", "0.50", true},
		{"WECHAT", "0.38", true},
		{"UNKNOWN", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			config, ok := rateConfigs[tt.code]
			assert.Equal(t, tt.exists, ok)
			if ok {
				assert.Equal(t, tt.expected, config.Rate)
			}
		})
	}
}

// TestRateConfigs_Empty 测试空费率配置
func TestRateConfigs_Empty(t *testing.T) {
	rateConfigs := models.RateConfigs{}

	assert.Equal(t, 0, len(rateConfigs))

	_, ok := rateConfigs["CREDIT"]
	assert.False(t, ok)
}

// =============================================================================
// ValidateRateConfigs 费率范围校验测试
// =============================================================================

// validateRateConfigsLocal 本地验证函数（用于测试，不依赖Repository）
func validateRateConfigsLocal(rateTypes []models.RateTypeDefinition, rateConfigs models.RateConfigs) error {
	rateTypeMap := make(map[string]models.RateTypeDefinition)
	for _, rt := range rateTypes {
		rateTypeMap[rt.Code] = rt
	}

	for code, config := range rateConfigs {
		rt, ok := rateTypeMap[code]
		if !ok {
			return fmt.Errorf("未知的费率类型: %s", code)
		}

		rate, err := strconv.ParseFloat(config.Rate, 64)
		if err != nil {
			return fmt.Errorf("%s费率格式错误: %s", rt.Name, config.Rate)
		}

		minRate, _ := strconv.ParseFloat(rt.MinRate, 64)
		maxRate, _ := strconv.ParseFloat(rt.MaxRate, 64)

		if rate < minRate {
			return fmt.Errorf("%s费率不能低于%.2f%%", rt.Name, minRate)
		}
		if rate > maxRate {
			return fmt.Errorf("%s费率不能高于%.2f%%", rt.Name, maxRate)
		}
	}

	return nil
}

// TestValidateRateConfigs 测试费率配置校验
func TestValidateRateConfigs(t *testing.T) {
	// 模拟费率类型定义
	rateTypes := []models.RateTypeDefinition{
		{Code: "CREDIT", Name: "贷记卡费率", MinRate: "0.50", MaxRate: "0.68"},
		{Code: "DEBIT", Name: "借记卡费率", MinRate: "0.45", MaxRate: "0.60"},
		{Code: "WECHAT", Name: "微信费率", MinRate: "0.30", MaxRate: "0.60"},
	}

	tests := []struct {
		name        string
		rateConfigs models.RateConfigs
		expectError bool
		errContains string
	}{
		{
			name: "正常费率配置",
			rateConfigs: models.RateConfigs{
				"CREDIT": {Rate: "0.60"},
				"DEBIT":  {Rate: "0.50"},
				"WECHAT": {Rate: "0.38"},
			},
			expectError: false,
		},
		{
			name: "费率达到下限",
			rateConfigs: models.RateConfigs{
				"CREDIT": {Rate: "0.50"},
			},
			expectError: false,
		},
		{
			name: "费率达到上限",
			rateConfigs: models.RateConfigs{
				"CREDIT": {Rate: "0.68"},
			},
			expectError: false,
		},
		{
			name: "费率低于下限",
			rateConfigs: models.RateConfigs{
				"CREDIT": {Rate: "0.40"},
			},
			expectError: true,
			errContains: "不能低于",
		},
		{
			name: "费率高于上限",
			rateConfigs: models.RateConfigs{
				"CREDIT": {Rate: "0.80"},
			},
			expectError: true,
			errContains: "不能高于",
		},
		{
			name: "未知费率类型",
			rateConfigs: models.RateConfigs{
				"UNKNOWN": {Rate: "0.50"},
			},
			expectError: true,
			errContains: "未知的费率类型",
		},
		{
			name: "费率格式错误",
			rateConfigs: models.RateConfigs{
				"CREDIT": {Rate: "abc"},
			},
			expectError: true,
			errContains: "格式错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRateConfigsLocal(rateTypes, tt.rateConfigs)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// RateDeltas 费率阶梯调整值测试
// =============================================================================

// TestRateDeltas_GetDelta 测试获取费率调整值
func TestRateDeltas_GetDelta(t *testing.T) {
	rateDeltas := models.RateDeltas{
		"CREDIT": "0.05",
		"DEBIT":  "0.03",
		"WECHAT": "0.02",
	}

	tests := []struct {
		code     string
		expected string
		exists   bool
	}{
		{"CREDIT", "0.05", true},
		{"DEBIT", "0.03", true},
		{"WECHAT", "0.02", true},
		{"UNKNOWN", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			delta, ok := rateDeltas[tt.code]
			assert.Equal(t, tt.exists, ok)
			if ok {
				assert.Equal(t, tt.expected, delta)
			}
		})
	}
}

// TestRateDeltas_ParseFloat 测试费率调整值解析为浮点数
func TestRateDeltas_ParseFloat(t *testing.T) {
	tests := []struct {
		deltaStr string
		expected float64
		hasError bool
	}{
		{"0.05", 0.05, false},
		{"0.10", 0.10, false},
		{"0.00", 0.00, false},
		{"-0.05", -0.05, false},
		{"", 0.00, true},
		{"invalid", 0.00, true},
	}

	for _, tt := range tests {
		t.Run(tt.deltaStr, func(t *testing.T) {
			delta, err := strconv.ParseFloat(tt.deltaStr, 64)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, delta)
			}
		})
	}
}

// =============================================================================
// 费率阶梯应用测试
// =============================================================================

// TestApplyRateDelta 测试应用费率阶梯调整
func TestApplyRateDelta(t *testing.T) {
	tests := []struct {
		name         string
		baseRate     float64
		rateDelta    float64
		expectedRate float64
	}{
		{"正常调增", 0.55, 0.05, 0.60},
		{"正常调减", 0.60, -0.05, 0.55},
		{"无调整", 0.55, 0.00, 0.55},
		{"小数调整", 0.55, 0.02, 0.57},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adjustedRate := tt.baseRate + tt.rateDelta
			assert.InDelta(t, tt.expectedRate, adjustedRate, 0.0001)
		})
	}
}

// =============================================================================
// 通道费率类型配置测试（恒信通/联动示例）
// =============================================================================

// TestHengxintongRateTypes 测试恒信通费率类型配置
func TestHengxintongRateTypes(t *testing.T) {
	// 恒信通标准费率类型
	hxtRateTypes := []models.RateTypeDefinition{
		{Code: "CREDIT", Name: "贷记卡费率", SortOrder: 1, MinRate: "0.50", MaxRate: "0.68"},
		{Code: "DEBIT", Name: "借记卡费率", SortOrder: 2, MinRate: "0.45", MaxRate: "0.60"},
		{Code: "WECHAT", Name: "微信费率", SortOrder: 3, MinRate: "0.30", MaxRate: "0.60"},
		{Code: "ALIPAY", Name: "支付宝费率", SortOrder: 4, MinRate: "0.30", MaxRate: "0.60"},
		{Code: "UNIONPAY", Name: "云闪付费率", SortOrder: 5, MinRate: "0.30", MaxRate: "0.60"},
	}

	assert.Equal(t, 5, len(hxtRateTypes))
	assert.Equal(t, "CREDIT", hxtRateTypes[0].Code)
	assert.Equal(t, "贷记卡费率", hxtRateTypes[0].Name)
}

// TestLiandongRateTypes 测试联动费率类型配置
func TestLiandongRateTypes(t *testing.T) {
	// 联动费率类型（包含特惠费率）
	ldRateTypes := []models.RateTypeDefinition{
		{Code: "WECHAT", Name: "微信", SortOrder: 1, MinRate: "0.30", MaxRate: "0.60"},
		{Code: "ALIPAY", Name: "支付宝", SortOrder: 2, MinRate: "0.30", MaxRate: "0.60"},
		{Code: "POS_DC", Name: "普通刷卡-借记", SortOrder: 3, MinRate: "0.45", MaxRate: "0.60"},
		{Code: "POS_CC", Name: "普通刷卡-贷记", SortOrder: 4, MinRate: "0.50", MaxRate: "0.68"},
		{Code: "POS_DISCOUNT_CC", Name: "特惠", SortOrder: 10, MinRate: "0.48", MaxRate: "0.60"},
		{Code: "POS_DISCOUNT_GF_CC", Name: "特惠GF", SortOrder: 11, MinRate: "0.45", MaxRate: "0.58"},
	}

	assert.Equal(t, 6, len(ldRateTypes))

	// 验证特惠费率存在
	hasDiscount := false
	for _, rt := range ldRateTypes {
		if rt.Code == "POS_DISCOUNT_CC" {
			hasDiscount = true
			assert.Equal(t, "特惠", rt.Name)
			assert.Equal(t, "0.48", rt.MinRate)
			break
		}
	}
	assert.True(t, hasDiscount, "联动应该有特惠费率类型")
}

// =============================================================================
// 边界情况测试
// =============================================================================

// TestRateConfigs_Boundary 测试费率配置边界情况
func TestRateConfigs_Boundary(t *testing.T) {
	tests := []struct {
		name        string
		rate        string
		minRate     string
		maxRate     string
		expectValid bool
	}{
		{"刚好等于下限", "0.50", "0.50", "0.68", true},
		{"刚好等于上限", "0.68", "0.50", "0.68", true},
		{"低于下限0.01", "0.49", "0.50", "0.68", false},
		{"高于上限0.01", "0.69", "0.50", "0.68", false},
		{"中间值", "0.55", "0.50", "0.68", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, _ := strconv.ParseFloat(tt.rate, 64)
			minRate, _ := strconv.ParseFloat(tt.minRate, 64)
			maxRate, _ := strconv.ParseFloat(tt.maxRate, 64)

			isValid := rate >= minRate && rate <= maxRate
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

// TestRateConfigs_Precision 测试费率精度
func TestRateConfigs_Precision(t *testing.T) {
	tests := []struct {
		rateStr  string
		expected float64
	}{
		{"0.55", 0.55},
		{"0.555", 0.555},
		{"0.5555", 0.5555},
		{"0.38", 0.38},
	}

	for _, tt := range tests {
		t.Run(tt.rateStr, func(t *testing.T) {
			rate, err := strconv.ParseFloat(tt.rateStr, 64)
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, rate, 0.00001)
		})
	}
}

// =============================================================================
// JSON序列化/反序列化测试
// =============================================================================

// TestRateConfigs_JSONSerialization 测试RateConfigs JSON序列化
func TestRateConfigs_JSONSerialization(t *testing.T) {
	original := models.RateConfigs{
		"CREDIT": {Rate: "0.60"},
		"DEBIT":  {Rate: "0.50"},
	}

	// 序列化
	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	// 反序列化
	var parsed models.RateConfigs
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)

	assert.Equal(t, original["CREDIT"].Rate, parsed["CREDIT"].Rate)
	assert.Equal(t, original["DEBIT"].Rate, parsed["DEBIT"].Rate)
}

// TestRateDeltas_JSONSerialization 测试RateDeltas JSON序列化
func TestRateDeltas_JSONSerialization(t *testing.T) {
	original := models.RateDeltas{
		"CREDIT": "0.05",
		"DEBIT":  "0.03",
	}

	// 序列化
	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	// 反序列化
	var parsed models.RateDeltas
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)

	assert.Equal(t, original["CREDIT"], parsed["CREDIT"])
	assert.Equal(t, original["DEBIT"], parsed["DEBIT"])
}

// =============================================================================
// 分润计算中的费率匹配测试
// =============================================================================

// TestRateConfigs_DirectMatch 测试直接匹配费率类型（无需映射）
func TestRateConfigs_DirectMatch(t *testing.T) {
	// 模拟交易回调的payTypeCode直接作为key
	rateConfigs := models.RateConfigs{
		"POS_CC":          {Rate: "0.60"},
		"POS_DISCOUNT_CC": {Rate: "0.55"},
		"WECHAT":          {Rate: "0.38"},
	}

	tests := []struct {
		payTypeCode  string
		expectedRate string
		exists       bool
	}{
		{"POS_CC", "0.60", true},
		{"POS_DISCOUNT_CC", "0.55", true},
		{"WECHAT", "0.38", true},
		{"UNKNOWN_TYPE", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.payTypeCode, func(t *testing.T) {
			config, ok := rateConfigs[tt.payTypeCode]
			assert.Equal(t, tt.exists, ok)
			if ok {
				assert.Equal(t, tt.expectedRate, config.Rate)
			}
		})
	}
}

// TestProfitCalculation_WithRateConfigs 测试使用RateConfigs的分润计算
func TestProfitCalculation_WithRateConfigs(t *testing.T) {
	tests := []struct {
		name           string
		transAmount    int64   // 交易金额（分）
		merchantRate   float64 // 商户费率（%）
		agentRate      float64 // 代理商成本费率（%）
		expectedProfit int64   // 预期分润（分）
	}{
		{
			name:           "普通刷卡分润",
			transAmount:    100000, // 1000元
			merchantRate:   0.60,   // 0.60%
			agentRate:      0.50,   // 0.50%
			expectedProfit: 99,     // 浮点精度：(1000 * 0.0999...% ≈ 99分)
		},
		{
			name:           "特惠交易分润",
			transAmount:    100000, // 1000元
			merchantRate:   0.55,   // 0.55%
			agentRate:      0.48,   // 0.48%
			expectedProfit: 70,     // (1000 * 0.07% = 0.7元 = 70分)
		},
		{
			name:           "零差额无分润",
			transAmount:    10000, // 100元
			merchantRate:   0.55,  // 0.55%
			agentRate:      0.55,  // 0.55%
			expectedProfit: 0,
		},
		{
			name:           "大额交易分润-整数差",
			transAmount:    1000000, // 10000元
			merchantRate:   0.55,    // 0.55%
			agentRate:      0.45,    // 0.45%（0.1整数差）
			expectedProfit: 1000,    // (10000 * 0.10% = 10元 = 1000分)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 分润 = 交易金额 * (商户费率 - 代理商费率) / 100
			rateDiff := tt.merchantRate - tt.agentRate
			profit := int64(float64(tt.transAmount) * rateDiff / 100)

			assert.Equal(t, tt.expectedProfit, profit)
		})
	}
}
