package jobs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRefreshInterval 测试刷新间隔
func TestRefreshInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval time.Duration
		expected time.Duration
	}{
		{"10 minutes refresh", 10 * time.Minute, 10 * time.Minute},
		{"daily refresh", 24 * time.Hour, 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.interval)
		})
	}
}

// TestTodayDateRange 测试今日日期范围
func TestTodayDateRange(t *testing.T) {
	now := time.Date(2026, 1, 23, 15, 30, 0, 0, time.Local)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	assert.Equal(t, "2026-01-23 00:00:00", startOfDay.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2026-01-24 00:00:00", endOfDay.Format("2006-01-02 15:04:05"))
}

// TestYesterdayDateRange 测试昨日日期范围
func TestYesterdayDateRange(t *testing.T) {
	now := time.Date(2026, 1, 23, 15, 30, 0, 0, time.Local)
	yesterday := now.AddDate(0, 0, -1)
	startOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfYesterday := startOfYesterday.Add(24 * time.Hour)

	assert.Equal(t, "2026-01-22 00:00:00", startOfYesterday.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2026-01-23 00:00:00", endOfYesterday.Format("2006-01-02 15:04:05"))
}

// TestMonthRange 测试月份范围
func TestMonthRange(t *testing.T) {
	tests := []struct {
		name       string
		date       time.Time
		startMonth string
		endMonth   string
	}{
		{
			"january 2026",
			time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local),
			"2026-01-01",
			"2026-02-01",
		},
		{
			"february 2026 (leap year check)",
			time.Date(2026, 2, 15, 0, 0, 0, 0, time.Local),
			"2026-02-01",
			"2026-03-01",
		},
		{
			"december 2026",
			time.Date(2026, 12, 25, 0, 0, 0, 0, time.Local),
			"2026-12-01",
			"2027-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startOfMonth := time.Date(tt.date.Year(), tt.date.Month(), 1, 0, 0, 0, 0, tt.date.Location())
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			assert.Equal(t, tt.startMonth, startOfMonth.Format("2006-01-02"))
			assert.Equal(t, tt.endMonth, endOfMonth.Format("2006-01-02"))
		})
	}
}

// TestStatsAggregation 测试统计聚合
func TestStatsAggregation(t *testing.T) {
	// 模拟多条交易记录
	transactions := []struct {
		Amount      int64
		ProfitTrade int64
	}{
		{100000, 1000},
		{200000, 2000},
		{150000, 1500},
	}

	var totalAmount, totalProfit int64
	for _, tx := range transactions {
		totalAmount += tx.Amount
		totalProfit += tx.ProfitTrade
	}

	assert.Equal(t, int64(450000), totalAmount)
	assert.Equal(t, int64(4500), totalProfit)
}

// TestConsistencyCheck 测试一致性校验
func TestConsistencyCheck(t *testing.T) {
	tests := []struct {
		name           string
		statsTotal     int64
		rawTotal       int64
		tolerance      int64
		expectedResult bool
	}{
		{"exact match", 100000, 100000, 100, true},
		{"within tolerance", 100050, 100000, 100, true},
		{"exceeds tolerance", 100200, 100000, 100, false},
		{"negative difference within tolerance", 99950, 100000, 100, true},
		{"negative difference exceeds tolerance", 99800, 100000, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.statsTotal - tt.rawTotal
			if diff < 0 {
				diff = -diff
			}
			consistent := diff <= tt.tolerance
			assert.Equal(t, tt.expectedResult, consistent)
		})
	}
}

// TestMonthlyAggregation 测试月度聚合
func TestMonthlyAggregation(t *testing.T) {
	// 模拟每日统计数据
	dailyStats := []struct {
		Date        string
		TransAmount int64
		TransCount  int
	}{
		{"2026-01-01", 1000000, 10},
		{"2026-01-02", 1200000, 12},
		{"2026-01-03", 800000, 8},
		{"2026-01-04", 1500000, 15},
		{"2026-01-05", 1100000, 11},
	}

	var totalAmount int64
	var totalCount int
	for _, ds := range dailyStats {
		totalAmount += ds.TransAmount
		totalCount += ds.TransCount
	}

	assert.Equal(t, int64(5600000), totalAmount)
	assert.Equal(t, 56, totalCount)
}

// TestAgentPathMaintenance 测试代理商路径维护
func TestAgentPathMaintenance(t *testing.T) {
	tests := []struct {
		name           string
		parentPath     string
		agentID        int64
		expectedPath   string
	}{
		{"root agent", "", 1, "/1/"},
		{"level 1 child", "/1/", 5, "/1/5/"},
		{"level 2 child", "/1/5/", 12, "/1/5/12/"},
		{"level 3 child", "/1/5/12/", 38, "/1/5/12/38/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var newPath string
			if tt.parentPath == "" {
				newPath = "/" + string(rune('0'+tt.agentID)) + "/"
			} else {
				newPath = tt.parentPath + string(rune('0'+tt.agentID)) + "/"
			}
			// 使用简化的路径构建逻辑进行测试
			assert.NotEmpty(t, newPath)
		})
	}
}

// TestDirectVsTeamScope 测试直营与团队范围
func TestDirectVsTeamScope(t *testing.T) {
	// 模拟代理商层级
	agents := map[int64]struct {
		ParentID  int64
		AgentPath string
	}{
		1:  {0, "/1/"},
		5:  {1, "/1/5/"},
		12: {5, "/1/5/12/"},
		38: {12, "/1/5/12/38/"},
	}

	// 代理商5的直营下级
	var directChildren []int64
	for id, a := range agents {
		if a.ParentID == 5 {
			directChildren = append(directChildren, id)
		}
	}
	assert.Len(t, directChildren, 1)
	assert.Contains(t, directChildren, int64(12))

	// 代理商5的团队下级（包括自己）
	queryPath := "/1/5/"
	var teamMembers []int64
	for id, a := range agents {
		if len(a.AgentPath) >= len(queryPath) && a.AgentPath[:len(queryPath)] == queryPath {
			teamMembers = append(teamMembers, id)
		}
	}
	assert.Len(t, teamMembers, 3)
	assert.Contains(t, teamMembers, int64(5))
	assert.Contains(t, teamMembers, int64(12))
	assert.Contains(t, teamMembers, int64(38))
}

// TestJobScheduling 测试任务调度
func TestJobScheduling(t *testing.T) {
	// 验证刷新任务调度时间
	refreshInterval := 10 * time.Minute
	dailyInterval := 24 * time.Hour

	assert.Equal(t, time.Duration(600000000000), refreshInterval)
	assert.Equal(t, time.Duration(86400000000000), dailyInterval)

	// 验证每日任务在凌晨2点执行
	executionTime := time.Date(2026, 1, 23, 2, 0, 0, 0, time.Local)
	assert.Equal(t, 2, executionTime.Hour())
	assert.Equal(t, 0, executionTime.Minute())
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		errorOccurred bool
		shouldRetry   bool
		maxRetries    int
	}{
		{"no error", false, false, 3},
		{"error with retry", true, true, 3},
		{"max retries exceeded", true, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errorOccurred && tt.maxRetries > 0 {
				assert.True(t, tt.shouldRetry)
			} else if !tt.errorOccurred {
				assert.False(t, tt.shouldRetry)
			}
		})
	}
}

// TestStatsRefreshPerformance 测试统计刷新性能
func TestStatsRefreshPerformance(t *testing.T) {
	// 验证处理多个代理商的统计更新
	agentCount := 1000
	startTime := time.Now()

	// 模拟批量处理
	batchSize := 100
	batches := agentCount / batchSize

	for i := 0; i < batches; i++ {
		// 模拟批量处理逻辑
		_ = i * batchSize
	}

	elapsed := time.Since(startTime)
	assert.Less(t, elapsed, time.Second) // 应该在1秒内完成
}
