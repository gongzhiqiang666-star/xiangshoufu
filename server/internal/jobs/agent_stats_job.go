package jobs

import (
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// AgentStatsRefreshJob 代理商统计刷新任务
// 每10分钟增量刷新今日汇总数据
type AgentStatsRefreshJob struct {
	db *gorm.DB
}

// NewAgentStatsRefreshJob 创建代理商统计刷新任务
func NewAgentStatsRefreshJob(db *gorm.DB) *AgentStatsRefreshJob {
	return &AgentStatsRefreshJob{db: db}
}

// Name 任务名称
func (j *AgentStatsRefreshJob) Name() string {
	return "agent_stats_refresh"
}

// Run 执行任务
func (j *AgentStatsRefreshJob) Run() {
	log.Println("[AgentStatsRefreshJob] 开始刷新今日统计数据...")
	startTime := time.Now()

	today := time.Now().Format("2006-01-02")

	// 获取所有代理商
	var agents []struct {
		ID        int64
		AgentPath string
	}
	if err := j.db.Table("agents").Select("id, agent_path").Find(&agents).Error; err != nil {
		log.Printf("[AgentStatsRefreshJob] 获取代理商列表失败: %v", err)
		return
	}

	successCount := 0
	for _, agent := range agents {
		// 刷新直营数据
		if err := j.refreshAgentDailyStats(agent.ID, agent.AgentPath, today, models.StatScopeDirect); err != nil {
			log.Printf("[AgentStatsRefreshJob] 刷新代理商 %d 直营数据失败: %v", agent.ID, err)
		} else {
			successCount++
		}

		// 刷新团队数据
		if err := j.refreshAgentDailyStats(agent.ID, agent.AgentPath, today, models.StatScopeTeam); err != nil {
			log.Printf("[AgentStatsRefreshJob] 刷新代理商 %d 团队数据失败: %v", agent.ID, err)
		}
	}

	elapsed := time.Since(startTime)
	log.Printf("[AgentStatsRefreshJob] 完成刷新，成功 %d/%d 个代理商，耗时 %v", successCount, len(agents), elapsed)
}

// refreshAgentDailyStats 刷新单个代理商的每日统计
func (j *AgentStatsRefreshJob) refreshAgentDailyStats(agentID int64, agentPath string, date string, scope string) error {
	stats := &models.AgentDailyStats{
		AgentID:   agentID,
		Scope:     scope,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 解析日期
	statDate, _ := time.Parse("2006-01-02", date)
	stats.StatDate = statDate

	// 确定查询条件
	var agentCondition string
	var merchantCondition string
	var terminalCondition string

	if scope == models.StatScopeDirect {
		agentCondition = fmt.Sprintf("agent_id = %d", agentID)
		merchantCondition = fmt.Sprintf("agent_id = %d", agentID)
		terminalCondition = fmt.Sprintf("owner_agent_id = %d", agentID)
	} else {
		// 团队范围：使用物化路径查询
		if agentPath == "" {
			agentPath = fmt.Sprintf("/%d/", agentID)
		}
		agentCondition = fmt.Sprintf("agent_id IN (SELECT id FROM agents WHERE agent_path LIKE '%s%%')", agentPath)
		merchantCondition = fmt.Sprintf("agent_id IN (SELECT id FROM agents WHERE agent_path LIKE '%s%%')", agentPath)
		terminalCondition = fmt.Sprintf("owner_agent_id IN (SELECT id FROM agents WHERE agent_path LIKE '%s%%')", agentPath)
	}

	// 1. 交易统计
	var transStats struct {
		TotalAmount int64
		TotalCount  int
	}
	j.db.Table("transactions").
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count").
		Where(agentCondition+" AND DATE(trade_time) = ?", date).
		Scan(&transStats)
	stats.TransAmount = transStats.TotalAmount
	stats.TransCount = transStats.TotalCount

	// 2. 分润统计(按类型分)
	var profitStats []struct {
		ProfitType  int16
		TotalAmount int64
	}
	j.db.Table("profit_records").
		Select("profit_type, COALESCE(SUM(profit_amount), 0) as total_amount").
		Where(agentCondition+" AND DATE(created_at) = ? AND is_revoked = false", date).
		Group("profit_type").
		Scan(&profitStats)

	for _, p := range profitStats {
		switch p.ProfitType {
		case models.ProfitTypeTrade:
			stats.ProfitTrade = p.TotalAmount
		case models.ProfitTypeDeposit:
			stats.ProfitDeposit = p.TotalAmount
		case models.ProfitTypeSim:
			stats.ProfitSim = p.TotalAmount
		case models.ProfitTypeReward:
			stats.ProfitReward = p.TotalAmount
		}
	}
	stats.ProfitTotal = stats.ProfitTrade + stats.ProfitDeposit + stats.ProfitSim + stats.ProfitReward

	// 3. 商户统计
	var merchantCount int64
	j.db.Table("merchants").Where(merchantCondition+" AND status = 1").Count(&merchantCount)
	stats.MerchantCount = int(merchantCount)

	var merchantNew int64
	j.db.Table("merchants").Where(merchantCondition+" AND DATE(created_at) = ?", date).Count(&merchantNew)
	stats.MerchantNew = int(merchantNew)

	// 4. 终端统计
	var terminalTotal int64
	j.db.Table("terminals").Where(terminalCondition).Count(&terminalTotal)
	stats.TerminalTotal = int(terminalTotal)

	var terminalActivated int64
	j.db.Table("terminals").Where(terminalCondition + " AND status = 3").Count(&terminalActivated)
	stats.TerminalActivated = int(terminalActivated)

	var terminalNewActivated int64
	j.db.Table("terminals").Where(terminalCondition+" AND DATE(activate_time) = ?", date).Count(&terminalNewActivated)
	stats.TerminalNewActivated = int(terminalNewActivated)

	// 5. 插入或更新记录
	return j.db.Exec(`
		INSERT INTO agent_daily_stats (
			agent_id, stat_date, scope,
			trans_amount, trans_count,
			profit_trade, profit_deposit, profit_sim, profit_reward, profit_total,
			merchant_count, merchant_new, terminal_total, terminal_activated, terminal_new_activated,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (agent_id, stat_date, scope) DO UPDATE SET
			trans_amount = EXCLUDED.trans_amount,
			trans_count = EXCLUDED.trans_count,
			profit_trade = EXCLUDED.profit_trade,
			profit_deposit = EXCLUDED.profit_deposit,
			profit_sim = EXCLUDED.profit_sim,
			profit_reward = EXCLUDED.profit_reward,
			profit_total = EXCLUDED.profit_total,
			merchant_count = EXCLUDED.merchant_count,
			merchant_new = EXCLUDED.merchant_new,
			terminal_total = EXCLUDED.terminal_total,
			terminal_activated = EXCLUDED.terminal_activated,
			terminal_new_activated = EXCLUDED.terminal_new_activated,
			updated_at = EXCLUDED.updated_at
	`, stats.AgentID, stats.StatDate, stats.Scope,
		stats.TransAmount, stats.TransCount,
		stats.ProfitTrade, stats.ProfitDeposit, stats.ProfitSim, stats.ProfitReward, stats.ProfitTotal,
		stats.MerchantCount, stats.MerchantNew, stats.TerminalTotal, stats.TerminalActivated, stats.TerminalNewActivated,
		stats.CreatedAt, stats.UpdatedAt).Error
}

// AgentStatsDailyJob 每日统计汇总任务
// 每天凌晨2点全量刷新昨日数据并汇总到月表
type AgentStatsDailyJob struct {
	db *gorm.DB
}

// NewAgentStatsDailyJob 创建每日统计汇总任务
func NewAgentStatsDailyJob(db *gorm.DB) *AgentStatsDailyJob {
	return &AgentStatsDailyJob{db: db}
}

// Name 任务名称
func (j *AgentStatsDailyJob) Name() string {
	return "agent_stats_daily"
}

// Run 执行任务
func (j *AgentStatsDailyJob) Run() {
	log.Println("[AgentStatsDailyJob] 开始执行每日统计汇总...")
	startTime := time.Now()

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	month := time.Now().AddDate(0, 0, -1).Format("2006-01")

	// 1. 刷新昨日统计
	refreshJob := NewAgentStatsRefreshJob(j.db)
	j.refreshDayStats(refreshJob, yesterday)

	// 2. 汇总到月表
	j.aggregateMonthlyStats(month)

	elapsed := time.Since(startTime)
	log.Printf("[AgentStatsDailyJob] 完成每日统计汇总，耗时 %v", elapsed)
}

// refreshDayStats 刷新指定日期的统计
func (j *AgentStatsDailyJob) refreshDayStats(refreshJob *AgentStatsRefreshJob, date string) {
	var agents []struct {
		ID        int64
		AgentPath string
	}
	if err := j.db.Table("agents").Select("id, agent_path").Find(&agents).Error; err != nil {
		log.Printf("[AgentStatsDailyJob] 获取代理商列表失败: %v", err)
		return
	}

	for _, agent := range agents {
		refreshJob.refreshAgentDailyStats(agent.ID, agent.AgentPath, date, models.StatScopeDirect)
		refreshJob.refreshAgentDailyStats(agent.ID, agent.AgentPath, date, models.StatScopeTeam)
	}
}

// aggregateMonthlyStats 汇总月度统计
func (j *AgentStatsDailyJob) aggregateMonthlyStats(month string) {
	log.Printf("[AgentStatsDailyJob] 汇总 %s 月度统计...", month)

	// 从日表汇总到月表
	err := j.db.Exec(`
		INSERT INTO agent_monthly_stats (
			agent_id, stat_month, scope,
			trans_amount, trans_count,
			profit_trade, profit_deposit, profit_sim, profit_reward, profit_total,
			merchant_count, merchant_new, terminal_total, terminal_activated,
			created_at, updated_at
		)
		SELECT
			agent_id,
			? as stat_month,
			scope,
			SUM(trans_amount),
			SUM(trans_count),
			SUM(profit_trade),
			SUM(profit_deposit),
			SUM(profit_sim),
			SUM(profit_reward),
			SUM(profit_total),
			MAX(merchant_count),
			SUM(merchant_new),
			MAX(terminal_total),
			MAX(terminal_activated),
			NOW(),
			NOW()
		FROM agent_daily_stats
		WHERE TO_CHAR(stat_date, 'YYYY-MM') = ?
		GROUP BY agent_id, scope
		ON CONFLICT (agent_id, stat_month, scope) DO UPDATE SET
			trans_amount = EXCLUDED.trans_amount,
			trans_count = EXCLUDED.trans_count,
			profit_trade = EXCLUDED.profit_trade,
			profit_deposit = EXCLUDED.profit_deposit,
			profit_sim = EXCLUDED.profit_sim,
			profit_reward = EXCLUDED.profit_reward,
			profit_total = EXCLUDED.profit_total,
			merchant_count = EXCLUDED.merchant_count,
			merchant_new = EXCLUDED.merchant_new,
			terminal_total = EXCLUDED.terminal_total,
			terminal_activated = EXCLUDED.terminal_activated,
			updated_at = NOW()
	`, month, month).Error

	if err != nil {
		log.Printf("[AgentStatsDailyJob] 汇总月度统计失败: %v", err)
	}
}

// StatsConsistencyChecker 统计一致性校验任务
type StatsConsistencyChecker struct {
	db *gorm.DB
}

// NewStatsConsistencyChecker 创建统计一致性校验任务
func NewStatsConsistencyChecker(db *gorm.DB) *StatsConsistencyChecker {
	return &StatsConsistencyChecker{db: db}
}

// Name 任务名称
func (j *StatsConsistencyChecker) Name() string {
	return "stats_consistency_checker"
}

// Run 执行任务
func (j *StatsConsistencyChecker) Run() {
	log.Println("[StatsConsistencyChecker] 开始校验统计数据一致性...")

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 校验交易金额
	j.checkTransactionAmount(yesterday)

	// 校验分润金额
	j.checkProfitAmount(yesterday)

	log.Println("[StatsConsistencyChecker] 一致性校验完成")
}

// checkTransactionAmount 校验交易金额一致性
func (j *StatsConsistencyChecker) checkTransactionAmount(date string) {
	// 从汇总表获取总额
	var statsTotal int64
	j.db.Table("agent_daily_stats").
		Select("COALESCE(SUM(trans_amount), 0)").
		Where("stat_date = ? AND scope = 'direct'", date).
		Scan(&statsTotal)

	// 从原始表计算总额
	var rawTotal int64
	j.db.Table("transactions").
		Select("COALESCE(SUM(amount), 0)").
		Where("DATE(trade_time) = ?", date).
		Scan(&rawTotal)

	diff := abs(statsTotal - rawTotal)
	if diff > 100 { // 允许1元误差
		log.Printf("[StatsConsistencyChecker] 交易金额不一致! 汇总表=%d分, 原始表=%d分, 差异=%d分",
			statsTotal, rawTotal, diff)
		// TODO: 发送告警通知
	}
}

// checkProfitAmount 校验分润金额一致性
func (j *StatsConsistencyChecker) checkProfitAmount(date string) {
	// 从汇总表获取总额
	var statsTotal int64
	j.db.Table("agent_daily_stats").
		Select("COALESCE(SUM(profit_total), 0)").
		Where("stat_date = ? AND scope = 'direct'", date).
		Scan(&statsTotal)

	// 从原始表计算总额
	var rawTotal int64
	j.db.Table("profit_records").
		Select("COALESCE(SUM(profit_amount), 0)").
		Where("DATE(created_at) = ? AND is_revoked = false", date).
		Scan(&rawTotal)

	diff := abs(statsTotal - rawTotal)
	if diff > 100 { // 允许1元误差
		log.Printf("[StatsConsistencyChecker] 分润金额不一致! 汇总表=%d分, 原始表=%d分, 差异=%d分",
			statsTotal, rawTotal, diff)
		// TODO: 发送告警通知
	}
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
