package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAgentRankingPeriod 测试代理商排名时间周期
func TestAgentRankingPeriod(t *testing.T) {
	tests := []struct {
		name           string
		period         string
		expectedPeriod string
	}{
		{"day period", "day", "day"},
		{"week period", "week", "week"},
		{"month period", "month", "month"},
		{"empty defaults to month", "", "month"},
		{"invalid defaults to month", "year", "month"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period := tt.period
			if period == "" || (period != "day" && period != "week" && period != "month") {
				period = "month"
			}
			assert.Equal(t, tt.expectedPeriod, period)
		})
	}
}

// TestRankingLimit 测试排名数量限制
func TestRankingLimit(t *testing.T) {
	tests := []struct {
		name          string
		requestLimit  int
		expectedLimit int
	}{
		{"default limit", 0, 10},
		{"custom limit 5", 5, 5},
		{"custom limit 20", 20, 20},
		{"max limit exceeded", 100, 50},
		{"negative defaults", -1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := tt.requestLimit
			if limit <= 0 {
				limit = 10
			} else if limit > 50 {
				limit = 50
			}
			assert.Equal(t, tt.expectedLimit, limit)
		})
	}
}

// TestAgentRankingData 测试代理商排名数据结构
func TestAgentRankingData(t *testing.T) {
	ranking := []struct {
		rank       int
		agentID    int64
		agentName  string
		agentNo    string
		value      int64
		valueYuan  float64
		change     int64
		changeRate float64
	}{
		{1, 101, "代理商A", "AGT001", 5200000, 52000.00, 500000, 10.64},
		{2, 102, "代理商B", "AGT002", 4800000, 48000.00, -200000, -4.00},
		{3, 103, "代理商C", "AGT003", 4500000, 45000.00, 100000, 2.27},
	}

	// 验证排名顺序（按value降序）
	for i := 0; i < len(ranking)-1; i++ {
		assert.Greater(t, ranking[i].value, ranking[i+1].value)
		assert.Equal(t, i+1, ranking[i].rank)
	}

	// 验证元转换
	for _, r := range ranking {
		assert.Equal(t, float64(r.value)/100, r.valueYuan)
	}

	// 验证变化率计算
	for _, r := range ranking {
		previousValue := r.value - r.change
		if previousValue > 0 {
			expectedRate := float64(r.change) / float64(previousValue) * 100
			assert.InDelta(t, expectedRate, r.changeRate, 0.1)
		}
	}
}

// TestMerchantRankingFilters 测试商户排名筛选
func TestMerchantRankingFilters(t *testing.T) {
	merchantTypes := []string{"all", "loyal", "quality", "potential", "normal", "low_active", "inactive"}

	for _, mt := range merchantTypes {
		t.Run(mt, func(t *testing.T) {
			// 验证类型有效
			valid := false
			for _, validType := range merchantTypes {
				if mt == validType {
					valid = true
					break
				}
			}
			assert.True(t, valid)
		})
	}
}

// TestMerchantRankingData 测试商户排名数据结构
func TestMerchantRankingData(t *testing.T) {
	ranking := []struct {
		rank           int
		merchantID     int64
		merchantName   string
		merchantNo     string
		merchantType   string
		typeName       string
		totalAmount    int64
		monthAmount    int64
		lastTransTime  string
	}{
		{1, 201, "张***店", "M001", "loyal", "忠诚客户", 15000000, 5200000, "2026-01-23 10:30:00"},
		{2, 202, "李***市", "M002", "quality", "优质客户", 12000000, 4800000, "2026-01-23 09:15:00"},
		{3, 203, "王***行", "M003", "potential", "潜力客户", 8000000, 3500000, "2026-01-22 16:45:00"},
	}

	// 验证排名顺序（按monthAmount降序）
	for i := 0; i < len(ranking)-1; i++ {
		assert.Greater(t, ranking[i].monthAmount, ranking[i+1].monthAmount)
	}

	// 验证商户名脱敏
	for _, r := range ranking {
		assert.Contains(t, r.merchantName, "***")
	}

	// 验证类型名称映射
	typeNames := map[string]string{
		"loyal":      "忠诚客户",
		"quality":    "优质客户",
		"potential":  "潜力客户",
		"normal":     "一般客户",
		"low_active": "低活跃客户",
		"inactive":   "无交易客户",
	}
	for _, r := range ranking {
		assert.Equal(t, typeNames[r.merchantType], r.typeName)
	}
}

// TestAnalyticsSummary 测试分析汇总数据
func TestAnalyticsSummary(t *testing.T) {
	summary := struct {
		transAmount      int64
		transCount       int
		profitTotal      int64
		merchantTotal    int
		merchantActive   int
		terminalTotal    int
		terminalActivated int
		avgTransAmount   int64
		successRate      float64
	}{
		transAmount:       123456000,
		transCount:        1500,
		profitTotal:       1234560,
		merchantTotal:     120,
		merchantActive:    85,
		terminalTotal:     150,
		terminalActivated: 120,
		avgTransAmount:    82304,
		successRate:       99.2,
	}

	// 验证平均交易额计算
	expectedAvg := summary.transAmount / int64(summary.transCount)
	assert.Equal(t, expectedAvg, summary.avgTransAmount)

	// 验证活跃商户比例
	activeRate := float64(summary.merchantActive) / float64(summary.merchantTotal) * 100
	assert.InDelta(t, 70.83, activeRate, 0.01)

	// 验证终端激活率
	activatedRate := float64(summary.terminalActivated) / float64(summary.terminalTotal) * 100
	assert.Equal(t, float64(80), activatedRate)
}

// TestRankByDimension 测试排名维度
func TestRankByDimension(t *testing.T) {
	validDimensions := []string{"trans_amount", "profit", "terminal"}

	tests := []struct {
		name              string
		rankBy            string
		expectedDimension string
	}{
		{"trans_amount", "trans_amount", "trans_amount"},
		{"profit", "profit", "profit"},
		{"terminal", "terminal", "terminal"},
		{"empty defaults to trans_amount", "", "trans_amount"},
		{"invalid defaults to trans_amount", "invalid", "trans_amount"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rankBy := tt.rankBy
			valid := false
			for _, d := range validDimensions {
				if rankBy == d {
					valid = true
					break
				}
			}
			if !valid {
				rankBy = "trans_amount"
			}
			assert.Equal(t, tt.expectedDimension, rankBy)
		})
	}
}
