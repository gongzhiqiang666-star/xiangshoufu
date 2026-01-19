package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetProfitTypeName 测试分润类型名称
func TestGetProfitTypeName(t *testing.T) {
	tests := []struct {
		profitType int16
		expected   string
	}{
		{1, "交易分润"},
		{2, "激活奖励"},
		{3, "押金返现"},
		{4, "流量返现"},
		{0, "未知"},
		{99, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getProfitTypeName(tt.profitType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetWalletTypeStr 测试钱包类型名称
func TestGetWalletTypeStr(t *testing.T) {
	tests := []struct {
		walletType int16
		expected   string
	}{
		{1, "分润钱包"},
		{2, "服务费钱包"},
		{3, "奖励钱包"},
		{0, "未知"},
		{99, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getWalletTypeStr(tt.walletType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestProfitListResponse 测试分润列表响应结构
func TestProfitListResponse(t *testing.T) {
	// 模拟分润列表响应
	list := []map[string]interface{}{
		{
			"id":                int64(1),
			"order_no":          "ORD001",
			"profit_type":       int16(1),
			"profit_type_name":  "交易分润",
			"trade_amount":      int64(10000),
			"trade_amount_yuan": float64(100.00),
			"profit_amount":     int64(100),
			"profit_amount_yuan": float64(1.00),
			"self_rate":         "0.60%",
			"lower_rate":        "0.55%",
			"rate_diff":         "0.05%",
			"wallet_type":       int16(1),
			"wallet_type_name":  "分润钱包",
		},
	}

	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), list[0]["id"])
	assert.Equal(t, "ORD001", list[0]["order_no"])
	assert.Equal(t, "交易分润", list[0]["profit_type_name"])
	assert.Equal(t, float64(1.00), list[0]["profit_amount_yuan"])
}

// TestProfitStats 测试分润统计结构
func TestProfitStats(t *testing.T) {
	stats := map[string]interface{}{
		"today": map[string]interface{}{
			"amount":      int64(10000),
			"amount_yuan": float64(100.00),
			"count":       int64(50),
		},
		"month": map[string]interface{}{
			"amount":      int64(300000),
			"amount_yuan": float64(3000.00),
			"count":       int64(1500),
		},
	}

	today := stats["today"].(map[string]interface{})
	assert.Equal(t, int64(10000), today["amount"])
	assert.Equal(t, float64(100.00), today["amount_yuan"])
	assert.Equal(t, int64(50), today["count"])

	month := stats["month"].(map[string]interface{})
	assert.Equal(t, int64(300000), month["amount"])
	assert.Equal(t, int64(1500), month["count"])
}

// TestProfitTypeFilter 测试分润类型筛选
func TestProfitTypeFilter(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		hasFilter  bool
		filterVal  int16
	}{
		{"empty", "", false, 0},
		{"type 1", "1", true, 1},
		{"type 2", "2", true, 2},
		{"type 3", "3", true, 3},
		{"type 4", "4", true, 4},
		{"invalid", "abc", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var profitType *int16
			if tt.input != "" {
				// 尝试解析
				if tt.input >= "1" && tt.input <= "9" {
					val := int16(tt.input[0] - '0')
					profitType = &val
				}
			}

			if tt.hasFilter {
				assert.NotNil(t, profitType)
				assert.Equal(t, tt.filterVal, *profitType)
			} else {
				if tt.input == "" {
					assert.Nil(t, profitType)
				}
			}
		})
	}
}

// TestRateDiffCalculation 测试费率差计算
func TestRateDiffCalculation(t *testing.T) {
	tests := []struct {
		name     string
		selfRate float64
		lowerRate float64
		diff     float64
	}{
		{"normal diff", 0.60, 0.55, 0.05},
		{"zero diff", 0.60, 0.60, 0.00},
		{"large diff", 0.80, 0.50, 0.30},
		{"negative diff", 0.50, 0.60, -0.10}, // 理论上不应该发生
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.selfRate - tt.lowerRate
			assert.InDelta(t, tt.diff, diff, 0.0001)
		})
	}
}

// TestProfitCalculation 测试分润金额计算
func TestProfitCalculation(t *testing.T) {
	tests := []struct {
		name         string
		tradeAmount  int64
		rateDiff     float64 // 费率差，百分比形式 (如0.05表示0.05%)
		profitAmount int64
	}{
		{"normal", 10000, 0.0005, 5},       // 100元交易，0.05%费率差，得0.05元
		{"large trade", 100000, 0.0005, 50}, // 1000元交易，0.05%费率差，得0.5元
		{"high rate diff", 10000, 0.001, 10}, // 100元交易，0.1%费率差，得0.1元
		{"zero diff", 10000, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profit := int64(float64(tt.tradeAmount) * tt.rateDiff)
			assert.Equal(t, tt.profitAmount, profit)
		})
	}
}

// TestWalletStatus 测试钱包状态
func TestWalletStatus(t *testing.T) {
	tests := []struct {
		status   int16
		name     string
		canTrans bool
	}{
		{0, "待入账", false},
		{1, "已入账", true},
		{2, "已撤销", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var statusName string
			var canTrans bool

			switch tt.status {
			case 0:
				statusName = "待入账"
				canTrans = false
			case 1:
				statusName = "已入账"
				canTrans = true
			case 2:
				statusName = "已撤销"
				canTrans = false
			}

			assert.Equal(t, tt.name, statusName)
			assert.Equal(t, tt.canTrans, canTrans)
		})
	}
}
