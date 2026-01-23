package repository

import (
	"fmt"
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormAgentStatsRepository 代理商统计数据仓库
type GormAgentStatsRepository struct {
	db *gorm.DB
}

// NewGormAgentStatsRepository 创建代理商统计数据仓库
func NewGormAgentStatsRepository(db *gorm.DB) *GormAgentStatsRepository {
	return &GormAgentStatsRepository{db: db}
}

// ========================================
// 每日统计相关方法
// ========================================

// GetDailyStats 获取指定日期的统计数据
func (r *GormAgentStatsRepository) GetDailyStats(agentID int64, date time.Time, scope string) (*models.AgentDailyStats, error) {
	var stats models.AgentDailyStats
	dateStr := date.Format("2006-01-02")
	err := r.db.Where("agent_id = ? AND stat_date = ? AND scope = ?", agentID, dateStr, scope).First(&stats).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &stats, err
}

// GetDailyStatsRange 获取日期范围内的统计数据
func (r *GormAgentStatsRepository) GetDailyStatsRange(agentID int64, startDate, endDate time.Time, scope string) ([]models.AgentDailyStats, error) {
	var stats []models.AgentDailyStats
	err := r.db.Where("agent_id = ? AND stat_date BETWEEN ? AND ? AND scope = ?",
		agentID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), scope).
		Order("stat_date ASC").
		Find(&stats).Error
	return stats, err
}

// UpsertDailyStats 更新或插入每日统计
func (r *GormAgentStatsRepository) UpsertDailyStats(stats *models.AgentDailyStats) error {
	stats.UpdatedAt = time.Now()
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "agent_id"}, {Name: "stat_date"}, {Name: "scope"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"trans_amount", "trans_count",
			"profit_trade", "profit_deposit", "profit_sim", "profit_reward", "profit_total",
			"merchant_count", "merchant_new", "terminal_total", "terminal_activated", "terminal_new_activated",
			"updated_at",
		}),
	}).Create(stats).Error
}

// GetTodayStats 获取今日统计(带昨日对比)
func (r *GormAgentStatsRepository) GetTodayStats(agentID int64, scope string) (*models.DayStats, *models.DayStats, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	var todayStats, yesterdayStats models.AgentDailyStats

	// 今日统计
	err := r.db.Where("agent_id = ? AND stat_date = ? AND scope = ?", agentID, today, scope).First(&todayStats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	// 昨日统计
	err = r.db.Where("agent_id = ? AND stat_date = ? AND scope = ?", agentID, yesterday, scope).First(&yesterdayStats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	todayResult := &models.DayStats{
		TransAmount:     todayStats.TransAmount,
		TransAmountYuan: float64(todayStats.TransAmount) / 100,
		TransCount:      todayStats.TransCount,
		ProfitTotal:     todayStats.ProfitTotal,
		ProfitTotalYuan: float64(todayStats.ProfitTotal) / 100,
		ProfitTrade:     todayStats.ProfitTrade,
		ProfitDeposit:   todayStats.ProfitDeposit,
		ProfitSim:       todayStats.ProfitSim,
		ProfitReward:    todayStats.ProfitReward,
	}

	yesterdayResult := &models.DayStats{
		TransAmount:     yesterdayStats.TransAmount,
		TransAmountYuan: float64(yesterdayStats.TransAmount) / 100,
		TransCount:      yesterdayStats.TransCount,
		ProfitTotal:     yesterdayStats.ProfitTotal,
		ProfitTotalYuan: float64(yesterdayStats.ProfitTotal) / 100,
		ProfitTrade:     yesterdayStats.ProfitTrade,
		ProfitDeposit:   yesterdayStats.ProfitDeposit,
		ProfitSim:       yesterdayStats.ProfitSim,
		ProfitReward:    yesterdayStats.ProfitReward,
	}

	return todayResult, yesterdayResult, nil
}

// GetTrendData 获取趋势数据
func (r *GormAgentStatsRepository) GetTrendData(agentID int64, days int, scope string) ([]models.TrendPoint, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -days+1).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	var stats []models.AgentDailyStats
	err := r.db.Where("agent_id = ? AND stat_date BETWEEN ? AND ? AND scope = ?",
		agentID, startDate, endDate, scope).
		Order("stat_date ASC").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	// 构建日期到数据的映射
	dataMap := make(map[string]*models.AgentDailyStats)
	for i := range stats {
		dateStr := stats[i].StatDate.Format("2006-01-02")
		dataMap[dateStr] = &stats[i]
	}

	// 填充完整的日期序列
	result := make([]models.TrendPoint, 0, days)
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -days+1+i)
		dateStr := date.Format("2006-01-02")
		displayDate := date.Format("01-02")

		point := models.TrendPoint{Date: displayDate}
		if s, ok := dataMap[dateStr]; ok {
			point.TransAmount = s.TransAmount
			point.TransAmountYuan = float64(s.TransAmount) / 100
			point.TransCount = s.TransCount
			point.ProfitTotal = s.ProfitTotal
			point.ProfitTotalYuan = float64(s.ProfitTotal) / 100
		}
		result = append(result, point)
	}

	return result, nil
}

// ========================================
// 每月统计相关方法
// ========================================

// GetMonthlyStats 获取指定月份的统计数据
func (r *GormAgentStatsRepository) GetMonthlyStats(agentID int64, month string, scope string) (*models.AgentMonthlyStats, error) {
	var stats models.AgentMonthlyStats
	err := r.db.Where("agent_id = ? AND stat_month = ? AND scope = ?", agentID, month, scope).First(&stats).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &stats, err
}

// UpsertMonthlyStats 更新或插入每月统计
func (r *GormAgentStatsRepository) UpsertMonthlyStats(stats *models.AgentMonthlyStats) error {
	stats.UpdatedAt = time.Now()
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "agent_id"}, {Name: "stat_month"}, {Name: "scope"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"trans_amount", "trans_count",
			"profit_trade", "profit_deposit", "profit_sim", "profit_reward", "profit_total",
			"merchant_count", "merchant_new", "terminal_total", "terminal_activated",
			"updated_at",
		}),
	}).Create(stats).Error
}

// GetCurrentMonthStats 获取本月统计
func (r *GormAgentStatsRepository) GetCurrentMonthStats(agentID int64, scope string) (*models.PeriodStats, error) {
	month := time.Now().Format("2006-01")
	stats, err := r.GetMonthlyStats(agentID, month, scope)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &models.PeriodStats{}, nil
	}

	return &models.PeriodStats{
		TransAmount:     stats.TransAmount,
		TransAmountYuan: float64(stats.TransAmount) / 100,
		TransCount:      stats.TransCount,
		ProfitTotal:     stats.ProfitTotal,
		ProfitTotalYuan: float64(stats.ProfitTotal) / 100,
		MerchantNew:     stats.MerchantNew,
	}, nil
}

// ========================================
// 通道统计相关方法
// ========================================

// GetChannelStats 获取通道交易占比
func (r *GormAgentStatsRepository) GetChannelStats(agentID int64, startDate, endDate time.Time, scope string) ([]models.ChannelStats, error) {
	var results []models.ChannelStats

	// 根据scope确定查询范围
	var agentCondition string
	var args []interface{}

	if scope == models.StatScopeDirect {
		agentCondition = "t.agent_id = ?"
		args = append(args, agentID)
	} else {
		// 团队范围：查询所有下级代理商的交易
		agentCondition = "t.agent_id IN (SELECT id FROM agents WHERE agent_path LIKE ?)"
		// 获取当前代理商的agent_path
		var agent struct{ AgentPath string }
		r.db.Table("agents").Select("agent_path").Where("id = ?", agentID).Scan(&agent)
		args = append(args, agent.AgentPath+"%")
	}

	args = append(args, startDate, endDate)

	query := fmt.Sprintf(`
		SELECT
			t.channel_id,
			t.channel_code,
			COALESCE(c.channel_name, t.channel_code) as channel_name,
			SUM(t.amount) as trans_amount,
			COUNT(*) as trans_count
		FROM transactions t
		LEFT JOIN channels c ON t.channel_id = c.id
		WHERE %s AND t.trade_time BETWEEN ? AND ?
		GROUP BY t.channel_id, t.channel_code, c.channel_name
		ORDER BY trans_amount DESC
	`, agentCondition)

	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 计算占比
	var totalAmount int64
	for _, r := range results {
		totalAmount += r.TransAmount
	}
	if totalAmount > 0 {
		for i := range results {
			results[i].Percentage = float64(results[i].TransAmount) * 100 / float64(totalAmount)
		}
	}

	return results, nil
}

// ========================================
// 商户分布统计
// ========================================

// GetMerchantDistribution 获取商户类型分布
func (r *GormAgentStatsRepository) GetMerchantDistribution(agentID int64, scope string) ([]models.MerchantDistribution, error) {
	var results []models.MerchantDistribution

	// 根据scope确定查询范围
	var condition string
	var args []interface{}

	if scope == models.StatScopeDirect {
		condition = "agent_id = ?"
		args = append(args, agentID)
	} else {
		condition = "agent_id IN (SELECT id FROM agents WHERE agent_path LIKE ?)"
		var agent struct{ AgentPath string }
		r.db.Table("agents").Select("agent_path").Where("id = ?", agentID).Scan(&agent)
		args = append(args, agent.AgentPath+"%")
	}

	query := fmt.Sprintf(`
		SELECT
			COALESCE(merchant_type, 'normal') as merchant_type,
			COUNT(*) as count
		FROM merchants
		WHERE %s AND status = 1
		GROUP BY merchant_type
		ORDER BY count DESC
	`, condition)

	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 添加类型名称和计算占比
	typeNames := map[string]string{
		"loyal":      "忠诚商户",
		"quality":    "优质商户",
		"potential":  "潜力商户",
		"normal":     "一般商户",
		"low_active": "低活跃",
		"inactive":   "30天无交易",
	}

	var totalCount int
	for _, r := range results {
		totalCount += r.Count
	}

	for i := range results {
		if name, ok := typeNames[results[i].MerchantType]; ok {
			results[i].TypeName = name
		} else {
			results[i].TypeName = results[i].MerchantType
		}
		if totalCount > 0 {
			results[i].Percentage = float64(results[i].Count) * 100 / float64(totalCount)
		}
	}

	return results, nil
}

// ========================================
// 排名相关方法
// ========================================

// GetAgentRanking 获取下级代理商排名
func (r *GormAgentStatsRepository) GetAgentRanking(agentID int64, period string, rankBy string, limit int) ([]models.AgentRanking, error) {
	var results []models.AgentRanking

	// 确定时间范围
	now := time.Now()
	var startDate, endDate string
	var prevStartDate, prevEndDate string

	switch period {
	case "day":
		startDate = now.Format("2006-01-02")
		endDate = startDate
		prevStartDate = now.AddDate(0, 0, -1).Format("2006-01-02")
		prevEndDate = prevStartDate
	case "week":
		startDate = now.AddDate(0, 0, -6).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
		prevStartDate = now.AddDate(0, 0, -13).Format("2006-01-02")
		prevEndDate = now.AddDate(0, 0, -7).Format("2006-01-02")
	default: // month
		startDate = now.Format("2006-01") + "-01"
		endDate = now.Format("2006-01-02")
		prevMonth := now.AddDate(0, -1, 0)
		prevStartDate = prevMonth.Format("2006-01") + "-01"
		prevEndDate = prevMonth.Format("2006-01-02")
	}

	// 确定排名字段
	var valueField string
	switch rankBy {
	case "profit":
		valueField = "profit_total"
	case "terminal":
		valueField = "terminal_new_activated"
	default: // trans_amount
		valueField = "trans_amount"
	}

	// 查询当期排名
	query := fmt.Sprintf(`
		WITH current_stats AS (
			SELECT
				agent_id,
				SUM(%s) as value
			FROM agent_daily_stats
			WHERE agent_id IN (SELECT id FROM agents WHERE parent_id = ?)
			  AND stat_date BETWEEN ? AND ?
			  AND scope = 'direct'
			GROUP BY agent_id
		),
		prev_stats AS (
			SELECT
				agent_id,
				SUM(%s) as value
			FROM agent_daily_stats
			WHERE agent_id IN (SELECT id FROM agents WHERE parent_id = ?)
			  AND stat_date BETWEEN ? AND ?
			  AND scope = 'direct'
			GROUP BY agent_id
		)
		SELECT
			ROW_NUMBER() OVER (ORDER BY c.value DESC) as rank,
			a.id as agent_id,
			a.agent_name,
			a.agent_no,
			COALESCE(c.value, 0) as value,
			COALESCE(c.value, 0) - COALESCE(p.value, 0) as change
		FROM agents a
		LEFT JOIN current_stats c ON a.id = c.agent_id
		LEFT JOIN prev_stats p ON a.id = p.agent_id
		WHERE a.parent_id = ?
		ORDER BY c.value DESC NULLS LAST
		LIMIT ?
	`, valueField, valueField)

	err := r.db.Raw(query, agentID, startDate, endDate, agentID, prevStartDate, prevEndDate, agentID, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 计算变化率和金额(元)
	for i := range results {
		results[i].ValueYuan = float64(results[i].Value) / 100
		prevValue := results[i].Value - results[i].Change
		if prevValue > 0 {
			results[i].ChangeRate = float64(results[i].Change) * 100 / float64(prevValue)
		}
	}

	return results, nil
}

// GetMerchantRanking 获取商户排名
func (r *GormAgentStatsRepository) GetMerchantRanking(agentID int64, merchantType string, scope string, limit int) ([]models.MerchantRanking, error) {
	var results []models.MerchantRanking

	now := time.Now()
	monthStart := now.Format("2006-01") + "-01"
	monthEnd := now.Format("2006-01-02")

	// 根据scope确定查询范围
	var merchantCondition string
	var args []interface{}

	if scope == models.StatScopeDirect {
		merchantCondition = "m.agent_id = ?"
		args = append(args, agentID)
	} else {
		merchantCondition = "m.agent_id IN (SELECT id FROM agents WHERE agent_path LIKE ?)"
		var agent struct{ AgentPath string }
		r.db.Table("agents").Select("agent_path").Where("id = ?", agentID).Scan(&agent)
		args = append(args, agent.AgentPath+"%")
	}

	// 添加商户类型筛选
	if merchantType != "" && merchantType != "all" {
		merchantCondition += " AND m.merchant_type = ?"
		args = append(args, merchantType)
	}

	args = append(args, monthStart, monthEnd, limit)

	query := fmt.Sprintf(`
		SELECT
			ROW_NUMBER() OVER (ORDER BY month_amount DESC) as rank,
			m.id as merchant_id,
			CONCAT(LEFT(m.merchant_name, 2), '***') as merchant_name,
			COALESCE(m.merchant_type, 'normal') as merchant_type,
			COALESCE(total.amount, 0) as total_amount,
			COALESCE(month.amount, 0) as month_amount
		FROM merchants m
		LEFT JOIN (
			SELECT merchant_id, SUM(amount) as amount
			FROM transactions
			GROUP BY merchant_id
		) total ON m.id = total.merchant_id
		LEFT JOIN (
			SELECT merchant_id, SUM(amount) as amount
			FROM transactions
			WHERE trade_time BETWEEN ? AND ?
			GROUP BY merchant_id
		) month ON m.id = month.merchant_id
		WHERE %s AND m.status = 1
		ORDER BY month_amount DESC NULLS LAST
		LIMIT ?
	`, merchantCondition)

	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for i := range results {
		results[i].MonthAmountYuan = float64(results[i].MonthAmount) / 100
	}

	return results, nil
}

// ========================================
// 终端统计方法
// ========================================

// GetTerminalStats 获取终端统计
func (r *GormAgentStatsRepository) GetTerminalStats(agentID int64, scope string) (*models.TerminalStats, error) {
	var condition string
	var args []interface{}

	if scope == models.StatScopeDirect {
		condition = "owner_agent_id = ?"
		args = append(args, agentID)
	} else {
		condition = "owner_agent_id IN (SELECT id FROM agents WHERE agent_path LIKE ?)"
		var agent struct{ AgentPath string }
		r.db.Table("agents").Select("agent_path").Where("id = ?", agentID).Scan(&agent)
		args = append(args, agent.AgentPath+"%")
	}

	var stats models.TerminalStats

	// 总数
	r.db.Table("terminals").Where(condition, args...).Count(new(int64))
	var total int64
	r.db.Table("terminals").Where(condition, args...).Count(&total)
	stats.Total = int(total)

	// 已激活
	var activated int64
	r.db.Table("terminals").Where(condition+" AND status = 3", args...).Count(&activated)
	stats.Activated = int(activated)

	// 今日激活
	today := time.Now().Format("2006-01-02")
	var todayActivated int64
	r.db.Table("terminals").Where(condition+" AND DATE(activate_time) = ?", append(args, today)...).Count(&todayActivated)
	stats.TodayActivated = int(todayActivated)

	// 本月激活
	monthStart := time.Now().Format("2006-01") + "-01"
	var monthActivated int64
	r.db.Table("terminals").Where(condition+" AND activate_time >= ?", append(args, monthStart)...).Count(&monthActivated)
	stats.MonthActivated = int(monthActivated)

	return &stats, nil
}

// ========================================
// 最近交易方法
// ========================================

// GetRecentTransactions 获取最近交易列表
func (r *GormAgentStatsRepository) GetRecentTransactions(agentID int64, limit int) ([]models.RecentTransaction, error) {
	var results []models.RecentTransaction

	query := `
		SELECT
			t.id,
			COALESCE(m.merchant_name, '未知商户') as merchant_name,
			t.pay_type,
			t.amount,
			t.trade_time
		FROM transactions t
		LEFT JOIN merchants m ON t.merchant_id = m.id
		WHERE t.agent_id = ?
		ORDER BY t.trade_time DESC
		LIMIT ?
	`

	err := r.db.Raw(query, agentID, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 添加支付类型名称和时间描述
	payTypeNames := map[int16]string{
		1: "刷卡",
		2: "微信",
		3: "支付宝",
		4: "云闪付",
		5: "银联二维码",
	}

	now := time.Now()
	for i := range results {
		results[i].AmountYuan = float64(results[i].Amount) / 100
		if name, ok := payTypeNames[results[i].PayType]; ok {
			results[i].PayTypeName = name
		} else {
			results[i].PayTypeName = "其他"
		}
		results[i].TimeAgo = formatTimeAgo(results[i].TradeTime, now)
	}

	return results, nil
}

// formatTimeAgo 格式化时间为"X分钟前"
func formatTimeAgo(t time.Time, now time.Time) string {
	diff := now.Sub(t)
	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d分钟前", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d小时前", int(diff.Hours()))
	} else {
		return t.Format("01-02 15:04")
	}
}
