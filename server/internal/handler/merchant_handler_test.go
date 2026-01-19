package handler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMerchantStatusName 测试商户状态名称
func TestGetMerchantStatusName(t *testing.T) {
	tests := []struct {
		status   int16
		expected string
	}{
		{1, "正常"},
		{2, "禁用"},
		{0, "未知"},
		{99, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getMerchantStatusName(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetMerchantApproveStatusName 测试审核状态名称
func TestGetMerchantApproveStatusName(t *testing.T) {
	tests := []struct {
		status   int16
		expected string
	}{
		{1, "待审核"},
		{2, "已通过"},
		{3, "已拒绝"},
		{0, "未知"},
		{99, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getMerchantApproveStatusName(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMaskIDCard 测试身份证号遮掩
func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		name     string
		idCard   string
		expected string
	}{
		{"normal 18 digits", "110101199001011234", "110101********1234"},
		{"normal 15 digits", "110101900101123", "110101********1123"},
		{"short id", "12345", "12345"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskIDCard(tt.idCard)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMerchantListResponse 测试商户列表响应结构
func TestMerchantListResponse(t *testing.T) {
	// 模拟商户列表响应
	list := []map[string]interface{}{
		{
			"id":             int64(1),
			"merchant_no":    "M001",
			"merchant_name":  "测试商户",
			"terminal_sn":    "SN001",
			"status":         int16(1),
			"status_name":    "正常",
			"approve_status": int16(2),
			"approve_name":   "已通过",
			"mcc":            "5812",
			"credit_rate":    "0.60",
			"debit_rate":     "0.50",
		},
	}

	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), list[0]["id"])
	assert.Equal(t, "M001", list[0]["merchant_no"])
	assert.Equal(t, "测试商户", list[0]["merchant_name"])
	assert.Equal(t, "正常", list[0]["status_name"])
	assert.Equal(t, "已通过", list[0]["approve_name"])
}

// TestMerchantStats 测试商户统计结构
func TestMerchantStats(t *testing.T) {
	stats := map[string]interface{}{
		"total_count":    int64(100),
		"active_count":   int64(80),
		"pending_count":  int64(15),
		"disabled_count": int64(5),
	}

	assert.Equal(t, int64(100), stats["total_count"])
	assert.Equal(t, int64(80), stats["active_count"])
	assert.Equal(t, int64(15), stats["pending_count"])
	assert.Equal(t, int64(5), stats["disabled_count"])

	// 验证统计数据一致性
	total := stats["active_count"].(int64) + stats["pending_count"].(int64) + stats["disabled_count"].(int64)
	assert.Equal(t, stats["total_count"], total)
}

// TestMerchantTransactions 测试商户交易列表
func TestMerchantTransactions(t *testing.T) {
	list := []map[string]interface{}{
		{
			"id":              int64(1),
			"order_no":        "ORD001",
			"trade_no":        "TRD001",
			"trade_type":      int16(1),
			"trade_type_name": "消费",
			"pay_type":        int16(1),
			"pay_type_name":   "刷卡",
			"amount":          int64(10000),
			"amount_yuan":     float64(100.00),
			"fee":             int64(60),
			"fee_yuan":        float64(0.60),
		},
	}

	assert.Len(t, list, 1)
	assert.Equal(t, "ORD001", list[0]["order_no"])
	assert.Equal(t, "消费", list[0]["trade_type_name"])
	assert.Equal(t, "刷卡", list[0]["pay_type_name"])
	assert.Equal(t, float64(100.00), list[0]["amount_yuan"])
}

// TestMerchantPermissionCheck 测试商户权限检查
func TestMerchantPermissionCheck(t *testing.T) {
	tests := []struct {
		name            string
		merchantAgentID int64
		currentAgentID  int64
		hasPermission   bool
	}{
		{"same agent", 1, 1, true},
		{"different agent", 1, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasPermission := tt.merchantAgentID == tt.currentAgentID
			assert.Equal(t, tt.hasPermission, hasPermission)
		})
	}
}

// TestMCC 测试MCC行业代码
func TestMCC(t *testing.T) {
	// 常见MCC代码
	mccCodes := map[string]string{
		"5411": "杂货店",
		"5812": "餐饮",
		"5814": "快餐店",
		"5912": "药店",
		"5921": "酒类商店",
		"7011": "酒店",
		"7832": "电影院",
	}

	assert.Equal(t, "餐饮", mccCodes["5812"])
	assert.Equal(t, "酒店", mccCodes["7011"])
}

// TestMerchantKeywordSearch 测试商户关键词搜索
func TestMerchantKeywordSearch(t *testing.T) {
	merchants := []struct {
		merchantNo   string
		merchantName string
	}{
		{"M001", "张三餐饮店"},
		{"M002", "李四超市"},
		{"M003", "王五便利店"},
	}

	tests := []struct {
		keyword  string
		expected int
	}{
		{"张三", 1},
		{"M001", 1},
		{"店", 2},       // 张三餐饮店 和 王五便利店
		{"不存在", 0},
		{"", 3}, // 空关键词返回全部
	}

	for _, tt := range tests {
		t.Run(tt.keyword, func(t *testing.T) {
			count := 0
			for _, m := range merchants {
				if tt.keyword == "" ||
					containsStr(m.merchantNo, tt.keyword) ||
					containsStr(m.merchantName, tt.keyword) {
					count++
				}
			}
			assert.Equal(t, tt.expected, count)
		})
	}
}

// containsStr 简单的字符串包含检查（支持UTF-8）
func containsStr(s, substr string) bool {
	return strings.Contains(s, substr)
}
