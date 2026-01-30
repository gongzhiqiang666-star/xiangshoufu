package handler

import (
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler 数据分析处理器
type AnalyticsHandler struct {
	statsRepo *repository.GormAgentStatsRepository
}

// NewAnalyticsHandler 创建数据分析处理器
func NewAnalyticsHandler(statsRepo *repository.GormAgentStatsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{
		statsRepo: statsRepo,
	}
}

// GetAgentRanking 获取代理商排名
// @Summary 获取代理商排名
// @Description 获取下级代理商业绩排名
// @Tags 数据分析
// @Produce json
// @Param period query string false "时间范围(day/week/month)" default(month)
// @Param rank_by query string false "排名维度(trans_amount/profit/terminal)" default(trans_amount)
// @Param limit query int false "数量限制" default(20)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/analytics/agent-ranking [get]
func (h *AnalyticsHandler) GetAgentRanking(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	period := c.DefaultQuery("period", "month")
	rankBy := c.DefaultQuery("rank_by", "trans_amount")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	ranking, err := h.statsRepo.GetAgentRanking(agentID, period, rankBy, limit)
	if err != nil {
		response.InternalError(c, "获取排名失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"ranking": ranking,
		"period":  period,
		"rank_by": rankBy,
	})
}

// GetMerchantRanking 获取商户排名
// @Summary 获取商户排名
// @Description 获取商户交易额排名
// @Tags 数据分析
// @Produce json
// @Param merchant_type query string false "商户类型筛选(all/loyal/quality/potential/normal/low_active/inactive)" default(all)
// @Param scope query string false "统计范围(direct/team)" default(direct)
// @Param limit query int false "数量限制" default(20)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/analytics/merchant-ranking [get]
func (h *AnalyticsHandler) GetMerchantRanking(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	merchantType := c.DefaultQuery("merchant_type", "all")
	scope := c.DefaultQuery("scope", models.StatScopeDirect)
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	ranking, err := h.statsRepo.GetMerchantRanking(agentID, merchantType, scope, limit)
	if err != nil {
		response.InternalError(c, "获取商户排名失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"ranking":       ranking,
		"merchant_type": merchantType,
		"scope":         scope,
	})
}

// GetAnalyticsSummary 获取分析汇总
// @Summary 获取分析汇总
// @Description 获取数据分析页面的汇总数据
// @Tags 数据分析
// @Produce json
// @Param period query string false "时间范围(day/week/month)" default(month)
// @Param scope query string false "统计范围" default(direct)
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/analytics/summary [get]
func (h *AnalyticsHandler) GetAnalyticsSummary(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	scope := c.DefaultQuery("scope", models.StatScopeDirect)
	period := c.DefaultQuery("period", "month")

	var days int
	switch period {
	case "day":
		days = 1
	case "week":
		days = 7
	default:
		days = 30
	}

	// 获取趋势数据
	trend, _ := h.statsRepo.GetTrendData(agentID, days, scope)

	// 获取通道统计 (本月)
	// 需要从 trend 中汇总
	var totalTransAmount, totalProfitTotal int64
	var totalTransCount int
	for _, t := range trend {
		totalTransAmount += t.TransAmount
		totalProfitTotal += t.ProfitTotal
		totalTransCount += t.TransCount
	}

	// 获取商户分布
	distribution, _ := h.statsRepo.GetMerchantDistribution(agentID, scope)

	// 获取终端统计
	terminalStats, _ := h.statsRepo.GetTerminalStats(agentID, scope)

	response.Success(c, gin.H{
		"summary": gin.H{
			"trans_amount":      totalTransAmount,
			"trans_amount_yuan": float64(totalTransAmount) / 100,
			"trans_count":       totalTransCount,
			"profit_total":      totalProfitTotal,
			"profit_total_yuan": float64(totalProfitTotal) / 100,
		},
		"trend":        trend,
		"distribution": distribution,
		"terminal":     terminalStats,
		"period":       period,
		"scope":        scope,
	})
}

// RegisterAnalyticsRoutes 注册数据分析路由
func RegisterAnalyticsRoutes(r *gin.RouterGroup, h *AnalyticsHandler, authService *service.AuthService) {
	analytics := r.Group("/analytics")
	analytics.Use(middleware.AuthMiddleware(authService))
	{
		analytics.GET("/agent-ranking", h.GetAgentRanking)
		analytics.GET("/merchant-ranking", h.GetMerchantRanking)
		analytics.GET("/summary", h.GetAnalyticsSummary)
	}
}
