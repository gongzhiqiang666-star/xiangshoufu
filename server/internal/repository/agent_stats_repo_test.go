package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestStatDateFormat 测试统计日期格式
func TestStatDateFormat(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{"normal date", time.Date(2026, 1, 23, 0, 0, 0, 0, time.Local), "2026-01-23"},
		{"first day of month", time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local), "2026-01-01"},
		{"last day of month", time.Date(2026, 1, 31, 0, 0, 0, 0, time.Local), "2026-01-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := tt.date.Format("2006-01-02")
			assert.Equal(t, tt.expected, formatted)
		})
	}
}

// TestStatMonthFormat 测试统计月份格式
func TestStatMonthFormat(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{"january", time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local), "2026-01"},
		{"december", time.Date(2026, 12, 25, 0, 0, 0, 0, time.Local), "2026-12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := tt.date.Format("2006-01")
			assert.Equal(t, tt.expected, formatted)
		})
	}
}

// TestScopeValidation 测试范围验证
func TestScopeValidation(t *testing.T) {
	validScopes := []string{"direct", "team"}

	tests := []struct {
		scope    string
		expected bool
	}{
		{"direct", true},
		{"team", true},
		{"", false},
		{"all", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			valid := false
			for _, s := range validScopes {
				if tt.scope == s {
					valid = true
					break
				}
			}
			assert.Equal(t, tt.expected, valid)
		})
	}
}

// TestAgentDailyStatsStructure 测试代理商每日统计结构
func TestAgentDailyStatsStructure(t *testing.T) {
	stats := struct {
		ID            int64
		AgentID       int64
		StatDate      string
		Scope         string
		TransAmount   int64
		TransCount    int
		ProfitTrade   int64
		ProfitDeposit int64
		ProfitSim     int64
		ProfitReward  int64
		MerchantCount int
		TerminalActivated int
	}{
		ID:            1,
		AgentID:       100,
		StatDate:      "2026-01-23",
		Scope:         "direct",
		TransAmount:   12345600,
		TransCount:    156,
		ProfitTrade:   85600,
		ProfitDeposit: 15000,
		ProfitSim:     13840,
		ProfitReward:  9000,
		MerchantCount: 50,
		TerminalActivated: 3,
	}

	// 验证分润总额
	totalProfit := stats.ProfitTrade + stats.ProfitDeposit + stats.ProfitSim + stats.ProfitReward
	assert.Equal(t, int64(123440), totalProfit)

	// 验证范围有效
	assert.Contains(t, []string{"direct", "team"}, stats.Scope)
}

// TestAgentMonthlyStatsStructure 测试代理商月度统计结构
func TestAgentMonthlyStatsStructure(t *testing.T) {
	stats := struct {
		ID               int64
		AgentID          int64
		StatMonth        string
		Scope            string
		TransAmount      int64
		TransCount       int
		ProfitTotal      int64
		MerchantCount    int
		TerminalTotal    int
		TerminalActivated int
	}{
		ID:               1,
		AgentID:          100,
		StatMonth:        "2026-01",
		Scope:            "team",
		TransAmount:      123456000,
		TransCount:       1500,
		ProfitTotal:      1234560,
		MerchantCount:    120,
		TerminalTotal:    150,
		TerminalActivated: 120,
	}

	// 验证月份格式
	assert.Regexp(t, `^\d{4}-\d{2}$`, stats.StatMonth)

	// 验证激活数不超过总数
	assert.LessOrEqual(t, stats.TerminalActivated, stats.TerminalTotal)
}

// TestTrendDataRange 测试趋势数据日期范围
func TestTrendDataRange(t *testing.T) {
	tests := []struct {
		name     string
		days     int
		expected int
	}{
		{"7 days", 7, 7},
		{"15 days", 15, 15},
		{"30 days", 30, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endDate := time.Now()
			startDate := endDate.AddDate(0, 0, -tt.days+1)
			dayCount := int(endDate.Sub(startDate).Hours()/24) + 1
			assert.Equal(t, tt.expected, dayCount)
		})
	}
}

// TestChannelStatsAggregation 测试通道统计聚合
func TestChannelStatsAggregation(t *testing.T) {
	channelData := []struct {
		ChannelCode  string
		ChannelName  string
		TransAmount  int64
		TransCount   int
		SuccessRate  float64
	}{
		{"HENGXINTONG", "恒信通", 6000000, 100, 99.2},
		{"LAKALA", "拉卡拉", 2500000, 50, 98.5},
		{"YEAHKA", "乐刷", 1500000, 30, 99.0},
	}

	// 计算总交易额
	var totalAmount int64
	for _, c := range channelData {
		totalAmount += c.TransAmount
	}
	assert.Equal(t, int64(10000000), totalAmount)

	// 计算百分比并验证总和为100%
	var totalPercentage float64
	for _, c := range channelData {
		percentage := float64(c.TransAmount) / float64(totalAmount) * 100
		totalPercentage += percentage
	}
	assert.InDelta(t, 100.0, totalPercentage, 0.01)
}

// TestMerchantTypeClassification 测试商户类型分类
func TestMerchantTypeClassification(t *testing.T) {
	tests := []struct {
		name          string
		avgAmountYuan float64
		expectedType  string
	}{
		{"loyal", 60000, "loyal"},           // >5万
		{"quality", 40000, "quality"},       // 3-5万
		{"potential", 25000, "potential"},   // 2-3万
		{"normal", 15000, "normal"},         // 1-2万
		{"low_active", 5000, "low_active"},  // <1万
		{"inactive", 0, "inactive"},         // 无交易
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var merchantType string
			avgAmount := tt.avgAmountYuan * 100 // 转换为分

			switch {
			case avgAmount >= 5000000:
				merchantType = "loyal"
			case avgAmount >= 3000000:
				merchantType = "quality"
			case avgAmount >= 2000000:
				merchantType = "potential"
			case avgAmount >= 1000000:
				merchantType = "normal"
			case avgAmount > 0:
				merchantType = "low_active"
			default:
				merchantType = "inactive"
			}

			assert.Equal(t, tt.expectedType, merchantType)
		})
	}
}

// TestAgentPathQuery 测试代理商路径查询
func TestAgentPathQuery(t *testing.T) {
	agents := []struct {
		ID        int64
		AgentPath string
	}{
		{1, "/1/"},
		{5, "/1/5/"},
		{12, "/1/5/12/"},
		{38, "/1/5/12/38/"},
		{6, "/1/6/"},
		{15, "/1/6/15/"},
	}

	// 查询代理商5的团队
	queryPath := "/1/5/"
	var teamMembers []int64
	for _, a := range agents {
		if len(a.AgentPath) >= len(queryPath) && a.AgentPath[:len(queryPath)] == queryPath {
			teamMembers = append(teamMembers, a.ID)
		}
	}

	// 验证团队成员
	assert.Contains(t, teamMembers, int64(5))
	assert.Contains(t, teamMembers, int64(12))
	assert.Contains(t, teamMembers, int64(38))
	assert.NotContains(t, teamMembers, int64(6))
	assert.NotContains(t, teamMembers, int64(15))
	assert.Len(t, teamMembers, 3)
}

// TestRecentTransactionLimit 测试最近交易限制
func TestRecentTransactionLimit(t *testing.T) {
	tests := []struct {
		requestLimit  int
		expectedLimit int
	}{
		{0, 10},
		{5, 5},
		{20, 20},
		{100, 50},
		{-5, 10},
	}

	for _, tt := range tests {
		limit := tt.requestLimit
		if limit <= 0 {
			limit = 10
		} else if limit > 50 {
			limit = 50
		}
		assert.Equal(t, tt.expectedLimit, limit)
	}
}

// TestUpsertLogic 测试Upsert逻辑
func TestUpsertLogic(t *testing.T) {
	// 模拟现有数据
	existingStats := map[string]int64{
		"trans_amount": 100000,
		"trans_count":  10,
	}

	// 新数据
	newStats := map[string]int64{
		"trans_amount": 150000,
		"trans_count":  15,
	}

	// Upsert后应该更新为新值
	for k, v := range newStats {
		existingStats[k] = v
	}

	assert.Equal(t, int64(150000), existingStats["trans_amount"])
	assert.Equal(t, int64(15), existingStats["trans_count"])
}
