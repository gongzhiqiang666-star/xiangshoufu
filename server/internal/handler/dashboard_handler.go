package handler

import (
	"strconv"
	"time"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 数据看板处理器
type DashboardHandler struct {
	transactionRepo *repository.GormTransactionRepository
	profitRepo      *repository.GormProfitRecordRepository
	agentRepo       *repository.GormAgentRepository
	walletRepo      *repository.GormWalletRepository
	statsRepo       *repository.GormAgentStatsRepository
}

// NewDashboardHandler 创建数据看板处理器
func NewDashboardHandler(
	transactionRepo *repository.GormTransactionRepository,
	profitRepo *repository.GormProfitRecordRepository,
	agentRepo *repository.GormAgentRepository,
	walletRepo *repository.GormWalletRepository,
	statsRepo *repository.GormAgentStatsRepository,
) *DashboardHandler {
	return &DashboardHandler{
		transactionRepo: transactionRepo,
		profitRepo:      profitRepo,
		agentRepo:       agentRepo,
		walletRepo:      walletRepo,
		statsRepo:       statsRepo,
	}
}

// GetDashboardOverview 获取数据概览
// @Summary 获取数据概览
// @Description 获取首页数据概览，包括今日交易、分润、团队等
// @Tags 数据看板
// @Produce json
// @Param scope query string false "统计范围: direct=直营, team=团队" default(direct)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/overview [get]
func (h *DashboardHandler) GetDashboardOverview(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	scope := c.DefaultQuery("scope", models.StatScopeDirect)

	// 验证scope参数
	if scope != models.StatScopeDirect && scope != models.StatScopeTeam {
		scope = models.StatScopeDirect
	}

	// 从汇总表获取今日/昨日统计
	todayStats, yesterdayStats, _ := h.statsRepo.GetTodayStats(agentID, scope)

	// 获取本月统计
	monthStats, _ := h.statsRepo.GetCurrentMonthStats(agentID, scope)

	// 代理商信息
	agent, _ := h.agentRepo.FindByIDFull(agentID)

	// 终端统计
	terminalStats, _ := h.statsRepo.GetTerminalStats(agentID, scope)

	// 钱包余额
	wallets, _ := h.walletRepo.FindByAgentID(agentID)
	var totalBalance int64
	for _, w := range wallets {
		totalBalance += w.Balance
	}

	response.Success(c, gin.H{
		"today": gin.H{
			"trans_amount":       todayStats.TransAmount,
			"trans_amount_yuan":  todayStats.TransAmountYuan,
			"trans_count":        todayStats.TransCount,
			"profit_total":       todayStats.ProfitTotal,
			"profit_total_yuan":  todayStats.ProfitTotalYuan,
			"profit_trade":       todayStats.ProfitTrade,
			"profit_deposit":     todayStats.ProfitDeposit,
			"profit_sim":         todayStats.ProfitSim,
			"profit_reward":      todayStats.ProfitReward,
		},
		"yesterday": gin.H{
			"trans_amount":       yesterdayStats.TransAmount,
			"trans_amount_yuan":  yesterdayStats.TransAmountYuan,
			"trans_count":        yesterdayStats.TransCount,
			"profit_total":       yesterdayStats.ProfitTotal,
			"profit_total_yuan":  yesterdayStats.ProfitTotalYuan,
		},
		"month": gin.H{
			"trans_amount":       monthStats.TransAmount,
			"trans_amount_yuan":  monthStats.TransAmountYuan,
			"trans_count":        monthStats.TransCount,
			"profit_total":       monthStats.ProfitTotal,
			"profit_total_yuan":  monthStats.ProfitTotalYuan,
			"merchant_new":       monthStats.MerchantNew,
		},
		"team": gin.H{
			"direct_agent_count":    agent.DirectAgentCount,
			"direct_merchant_count": agent.DirectMerchantCount,
			"team_agent_count":      agent.TeamAgentCount,
			"team_merchant_count":   agent.TeamMerchantCount,
		},
		"terminal": terminalStats,
		"wallet": gin.H{
			"total_balance":      totalBalance,
			"total_balance_yuan": float64(totalBalance) / 100,
		},
	})
}

// GetDashboardCharts 获取图表数据
// @Summary 获取图表数据
// @Description 获取交易和分润趋势图表数据
// @Tags 数据看板
// @Produce json
// @Param days query int false "天数(7/15/30)" default(7)
// @Param scope query string false "统计范围" default(direct)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/charts [get]
func (h *DashboardHandler) GetDashboardCharts(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	scope := c.DefaultQuery("scope", models.StatScopeDirect)
	daysStr := c.DefaultQuery("days", "7")
	days, _ := strconv.Atoi(daysStr)
	if days <= 0 || days > 30 {
		days = 7
	}

	// 从汇总表获取趋势数据
	trend, _ := h.statsRepo.GetTrendData(agentID, days, scope)

	response.Success(c, gin.H{
		"trans_trend": trend,
	})
}

// GetProfitTrend 获取收益趋势
// @Summary 获取收益趋势
// @Description 获取分润收益趋势图表数据
// @Tags 数据看板
// @Produce json
// @Param days query int false "天数(7/30)" default(7)
// @Param scope query string false "统计范围" default(direct)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/profit-trend [get]
func (h *DashboardHandler) GetProfitTrend(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	scope := c.DefaultQuery("scope", models.StatScopeDirect)
	daysStr := c.DefaultQuery("days", "7")
	days, _ := strconv.Atoi(daysStr)
	if days <= 0 || days > 30 {
		days = 7
	}

	trend, _ := h.statsRepo.GetTrendData(agentID, days, scope)

	response.Success(c, gin.H{
		"profit_trend": trend,
	})
}

// GetChannelStats 获取通道统计
// @Summary 获取通道统计
// @Description 获取各通道交易占比
// @Tags 数据看板
// @Produce json
// @Param scope query string false "统计范围" default(direct)
// @Param period query string false "时间范围(day/week/month)" default(month)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/channel-stats [get]
func (h *DashboardHandler) GetChannelStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	scope := c.DefaultQuery("scope", models.StatScopeDirect)
	period := c.DefaultQuery("period", "month")

	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = now
	case "week":
		startDate = now.AddDate(0, 0, -6)
		endDate = now
	default: // month
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	}

	stats, _ := h.statsRepo.GetChannelStats(agentID, startDate, endDate, scope)

	response.Success(c, gin.H{
		"channel_stats": stats,
	})
}

// GetMerchantDistribution 获取商户类型分布
// @Summary 获取商户类型分布
// @Description 获取商户按类型的分布统计
// @Tags 数据看板
// @Produce json
// @Param scope query string false "统计范围" default(direct)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/merchant-distribution [get]
func (h *DashboardHandler) GetMerchantDistribution(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	scope := c.DefaultQuery("scope", models.StatScopeDirect)

	distribution, _ := h.statsRepo.GetMerchantDistribution(agentID, scope)

	response.Success(c, gin.H{
		"distribution": distribution,
	})
}

// GetRecentTransactions 获取最近交易
// @Summary 获取最近交易
// @Description 获取最近交易列表
// @Tags 数据看板
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/recent-transactions [get]
func (h *DashboardHandler) GetRecentTransactions(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	transactions, _ := h.statsRepo.GetRecentTransactions(agentID, limit)

	response.Success(c, gin.H{
		"transactions": transactions,
	})
}

// RegisterDashboardRoutes 注册数据看板路由
func RegisterDashboardRoutes(r *gin.RouterGroup, h *DashboardHandler, authService *service.AuthService) {
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(authService))
	{
		dashboard.GET("/overview", h.GetDashboardOverview)
		dashboard.GET("/charts", h.GetDashboardCharts)
		dashboard.GET("/profit-trend", h.GetProfitTrend)
		dashboard.GET("/channel-stats", h.GetChannelStats)
		dashboard.GET("/merchant-distribution", h.GetMerchantDistribution)
		dashboard.GET("/recent-transactions", h.GetRecentTransactions)
	}
}
