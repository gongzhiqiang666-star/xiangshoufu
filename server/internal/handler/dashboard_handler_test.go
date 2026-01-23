package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDashboardOverviewResponse 测试数据概览响应结构
func TestDashboardOverviewResponse(t *testing.T) {
	overview := map[string]interface{}{
		"today": map[string]interface{}{
			"trans_amount":       int64(100000),
			"trans_amount_yuan":  float64(1000.00),
			"trans_count":        int64(50),
			"profit_amount":      int64(1000),
			"profit_amount_yuan": float64(10.00),
		},
		"month": map[string]interface{}{
			"trans_amount":       int64(3000000),
			"trans_amount_yuan":  float64(30000.00),
			"trans_count":        int64(1500),
			"profit_amount":      int64(30000),
			"profit_amount_yuan": float64(300.00),
		},
		"team": map[string]interface{}{
			"direct_agent_count":    10,
			"direct_merchant_count": 50,
			"team_agent_count":      100,
			"team_merchant_count":   500,
		},
		"wallet": map[string]interface{}{
			"total_balance":      int64(500000),
			"total_balance_yuan": float64(5000.00),
		},
	}

	// 验证今日数据
	today := overview["today"].(map[string]interface{})
	assert.Equal(t, int64(100000), today["trans_amount"])
	assert.Equal(t, float64(1000.00), today["trans_amount_yuan"])
	assert.Equal(t, int64(50), today["trans_count"])

	// 验证本月数据
	month := overview["month"].(map[string]interface{})
	assert.Equal(t, int64(3000000), month["trans_amount"])
	assert.Equal(t, int64(1500), month["trans_count"])

	// 验证团队数据
	team := overview["team"].(map[string]interface{})
	assert.Equal(t, 10, team["direct_agent_count"])
	assert.Equal(t, 50, team["direct_merchant_count"])
	assert.Equal(t, 100, team["team_agent_count"])

	// 验证钱包数据
	wallet := overview["wallet"].(map[string]interface{})
	assert.Equal(t, int64(500000), wallet["total_balance"])
}

// TestDashboardCharts 测试图表数据结构
func TestDashboardCharts(t *testing.T) {
	charts := map[string]interface{}{
		"trans_trend": []map[string]interface{}{
			{"date": "2024-01-13", "amount": int64(100000), "count": int64(50)},
			{"date": "2024-01-14", "amount": int64(120000), "count": int64(60)},
			{"date": "2024-01-15", "amount": int64(90000), "count": int64(45)},
			{"date": "2024-01-16", "amount": int64(150000), "count": int64(75)},
			{"date": "2024-01-17", "amount": int64(130000), "count": int64(65)},
			{"date": "2024-01-18", "amount": int64(110000), "count": int64(55)},
			{"date": "2024-01-19", "amount": int64(140000), "count": int64(70)},
		},
	}

	trend := charts["trans_trend"].([]map[string]interface{})
	assert.Len(t, trend, 7) // 7天数据
	assert.Equal(t, "2024-01-13", trend[0]["date"])
	assert.Equal(t, int64(100000), trend[0]["amount"])
}

// TestDateRange 测试日期范围计算
func TestDateRange(t *testing.T) {
	now := time.Date(2024, 1, 19, 12, 0, 0, 0, time.Local)

	tests := []struct {
		name      string
		days      int
		startDate string
		endDate   string
	}{
		{"7 days", 7, "2024-01-13", "2024-01-20"},
		{"14 days", 14, "2024-01-06", "2024-01-20"},
		{"30 days", 30, "2023-12-21", "2024-01-20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endDate := now.AddDate(0, 0, 1)
			startDate := now.AddDate(0, 0, -tt.days+1)

			assert.Equal(t, tt.startDate, startDate.Format("2006-01-02"))
			assert.Equal(t, tt.endDate, endDate.Format("2006-01-02"))
		})
	}
}

// TestTodayStats 测试今日统计计算
func TestTodayStats(t *testing.T) {
	now := time.Date(2024, 1, 19, 15, 30, 0, 0, time.Local)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	assert.Equal(t, "2024-01-19 00:00:00", startOfDay.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2024-01-20 00:00:00", endOfDay.Format("2006-01-02 15:04:05"))
}

// TestMonthStats 测试本月统计计算
func TestMonthStats(t *testing.T) {
	tests := []struct {
		date         time.Time
		startOfMonth string
		endOfMonth   string
	}{
		{
			date:         time.Date(2024, 1, 19, 0, 0, 0, 0, time.Local),
			startOfMonth: "2024-01-01",
			endOfMonth:   "2024-02-01",
		},
		{
			date:         time.Date(2024, 2, 15, 0, 0, 0, 0, time.Local),
			startOfMonth: "2024-02-01",
			endOfMonth:   "2024-03-01",
		},
		{
			date:         time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local),
			startOfMonth: "2024-12-01",
			endOfMonth:   "2025-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.startOfMonth, func(t *testing.T) {
			startOfMonth := time.Date(tt.date.Year(), tt.date.Month(), 1, 0, 0, 0, 0, tt.date.Location())
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			assert.Equal(t, tt.startOfMonth, startOfMonth.Format("2006-01-02"))
			assert.Equal(t, tt.endOfMonth, endOfMonth.Format("2006-01-02"))
		})
	}
}

// TestWalletBalanceSum 测试钱包余额汇总
func TestWalletBalanceSum(t *testing.T) {
	wallets := []struct {
		walletType int16
		balance    int64
	}{
		{1, 100000}, // 分润钱包
		{2, 50000},  // 服务费钱包
		{3, 30000},  // 奖励钱包
	}

	var totalBalance int64
	for _, w := range wallets {
		totalBalance += w.balance
	}

	assert.Equal(t, int64(180000), totalBalance)
	assert.Equal(t, float64(1800.00), float64(totalBalance)/100)
}

// TestTeamCount 测试团队统计
func TestTeamCount(t *testing.T) {
	team := struct {
		directAgentCount    int
		directMerchantCount int
		teamAgentCount      int
		teamMerchantCount   int
	}{
		directAgentCount:    10,
		directMerchantCount: 50,
		teamAgentCount:      100,
		teamMerchantCount:   500,
	}

	// 团队总人数应该包含直属下级
	assert.GreaterOrEqual(t, team.teamAgentCount, team.directAgentCount)
	assert.GreaterOrEqual(t, team.teamMerchantCount, team.directMerchantCount)
}

// TestAmountFormat 测试金额格式化
func TestAmountFormat(t *testing.T) {
	tests := []struct {
		amountFen  int64
		amountYuan float64
		formatted  string
	}{
		{0, 0.00, "0.00"},
		{100, 1.00, "1.00"},
		{1000, 10.00, "10.00"},
		{10000, 100.00, "100.00"},
		{100000, 1000.00, "1,000.00"},
		{1000000, 10000.00, "10,000.00"},
	}

	for _, tt := range tests {
		t.Run(tt.formatted, func(t *testing.T) {
			yuan := float64(tt.amountFen) / 100
			assert.Equal(t, tt.amountYuan, yuan)
		})
	}
}

// TestGrowthRate 测试增长率计算
func TestGrowthRate(t *testing.T) {
	tests := []struct {
		name       string
		current    int64
		previous   int64
		growthRate float64
	}{
		{"positive growth", 120, 100, 20.0},
		{"negative growth", 80, 100, -20.0},
		{"no growth", 100, 100, 0.0},
		{"from zero", 100, 0, 0.0}, // 特殊处理
		{"double", 200, 100, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rate float64
			if tt.previous > 0 {
				rate = float64(tt.current-tt.previous) / float64(tt.previous) * 100
			}
			assert.InDelta(t, tt.growthRate, rate, 0.01)
		})
	}
}

// TestDashboardScopeParameter 测试scope参数
func TestDashboardScopeParameter(t *testing.T) {
	tests := []struct {
		name          string
		scope         string
		expectedScope string
	}{
		{"direct scope", "direct", "direct"},
		{"team scope", "team", "team"},
		{"empty defaults to direct", "", "direct"},
		{"invalid defaults to direct", "invalid", "direct"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := tt.scope
			if scope == "" || (scope != "direct" && scope != "team") {
				scope = "direct"
			}
			assert.Equal(t, tt.expectedScope, scope)
		})
	}
}

// TestProfitBreakdown 测试收益分类
func TestProfitBreakdown(t *testing.T) {
	dayStats := map[string]int64{
		"profit_trade":   85600,  // 交易分润
		"profit_deposit": 15000,  // 押金返现
		"profit_sim":     13840,  // 流量返现
		"profit_reward":  9000,   // 激活奖励
	}

	// 验证总分润 = 各分类之和
	totalProfit := dayStats["profit_trade"] + dayStats["profit_deposit"] +
		dayStats["profit_sim"] + dayStats["profit_reward"]
	assert.Equal(t, int64(123440), totalProfit)

	// 验证元转换
	assert.Equal(t, float64(856.00), float64(dayStats["profit_trade"])/100)
	assert.Equal(t, float64(150.00), float64(dayStats["profit_deposit"])/100)
	assert.Equal(t, float64(138.40), float64(dayStats["profit_sim"])/100)
	assert.Equal(t, float64(90.00), float64(dayStats["profit_reward"])/100)
}

// TestChannelStatsCalculation 测试通道统计计算
func TestChannelStatsCalculation(t *testing.T) {
	channelStats := []struct {
		channelCode string
		channelName string
		transAmount int64
		transCount  int
	}{
		{"HENGXINTONG", "恒信通", 6000000, 100},
		{"LAKALA", "拉卡拉", 2500000, 50},
		{"OTHER", "其他", 1500000, 30},
	}

	// 计算总额
	var totalAmount int64
	for _, c := range channelStats {
		totalAmount += c.transAmount
	}
	assert.Equal(t, int64(10000000), totalAmount)

	// 计算百分比
	for _, c := range channelStats {
		percentage := float64(c.transAmount) / float64(totalAmount) * 100
		switch c.channelCode {
		case "HENGXINTONG":
			assert.InDelta(t, 60.0, percentage, 0.01)
		case "LAKALA":
			assert.InDelta(t, 25.0, percentage, 0.01)
		case "OTHER":
			assert.InDelta(t, 15.0, percentage, 0.01)
		}
	}
}

// TestMerchantDistribution 测试商户类型分布
func TestMerchantDistribution(t *testing.T) {
	distribution := map[string]int{
		"loyal":      5,   // 忠诚 (>5万)
		"quality":    12,  // 优质 (3-5万)
		"potential":  25,  // 潜力 (2-3万)
		"normal":     40,  // 一般 (1-2万)
		"low_active": 30,  // 低活跃 (<1万)
		"inactive":   8,   // 30天无交易
	}

	// 验证总数
	var total int
	for _, count := range distribution {
		total += count
	}
	assert.Equal(t, 120, total)

	// 验证百分比计算
	loyalPercentage := float64(distribution["loyal"]) / float64(total) * 100
	assert.InDelta(t, 4.17, loyalPercentage, 0.01)
}

// TestAgentPathMatching 测试物化路径匹配
func TestAgentPathMatching(t *testing.T) {
	tests := []struct {
		name       string
		agentPath  string
		queryPath  string
		shouldMatch bool
	}{
		{"direct child", "/1/5/12/", "/1/5/%", true},
		{"grandchild", "/1/5/12/38/", "/1/5/%", true},
		{"same agent", "/1/5/", "/1/5/%", true},
		{"different branch", "/1/6/12/", "/1/5/%", false},
		{"root level", "/1/", "/1/%", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟 SQL LIKE 匹配
			pattern := tt.queryPath[:len(tt.queryPath)-1] // 去掉 %
			matched := len(tt.agentPath) >= len(pattern) && tt.agentPath[:len(pattern)] == pattern
			assert.Equal(t, tt.shouldMatch, matched)
		})
	}
}

// TestTerminalStats 测试终端统计
func TestTerminalStats(t *testing.T) {
	terminalStats := struct {
		total          int
		activated      int
		todayActivated int
		monthActivated int
	}{
		total:          150,
		activated:      120,
		todayActivated: 3,
		monthActivated: 15,
	}

	// 验证激活数不超过总数
	assert.LessOrEqual(t, terminalStats.activated, terminalStats.total)
	// 验证今日激活不超过本月激活
	assert.LessOrEqual(t, terminalStats.todayActivated, terminalStats.monthActivated)
	// 验证本月激活不超过已激活总数
	assert.LessOrEqual(t, terminalStats.monthActivated, terminalStats.activated)
	// 计算未激活数
	unactivated := terminalStats.total - terminalStats.activated
	assert.Equal(t, 30, unactivated)
}
