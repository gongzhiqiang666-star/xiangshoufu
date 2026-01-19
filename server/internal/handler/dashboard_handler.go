package handler

import (
	"net/http"
	"time"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 数据看板处理器
type DashboardHandler struct {
	transactionRepo *repository.GormTransactionRepository
	profitRepo      *repository.GormProfitRecordRepository
	agentRepo       *repository.GormAgentRepository
	walletRepo      *repository.GormWalletRepository
}

// NewDashboardHandler 创建数据看板处理器
func NewDashboardHandler(
	transactionRepo *repository.GormTransactionRepository,
	profitRepo *repository.GormProfitRecordRepository,
	agentRepo *repository.GormAgentRepository,
	walletRepo *repository.GormWalletRepository,
) *DashboardHandler {
	return &DashboardHandler{
		transactionRepo: transactionRepo,
		profitRepo:      profitRepo,
		agentRepo:       agentRepo,
		walletRepo:      walletRepo,
	}
}

// GetDashboardOverview 获取数据概览
// @Summary 获取数据概览
// @Description 获取首页数据概览，包括今日交易、分润、团队等
// @Tags 数据看板
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/overview [get]
func (h *DashboardHandler) GetDashboardOverview(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	now := time.Now()

	// 今日交易统计
	todayTrans, _ := h.transactionRepo.GetAgentDailyStats(agentID, now)

	// 今日分润统计
	todayProfit, _ := h.profitRepo.GetAgentDailyProfitStats(agentID, now)

	// 本月交易统计
	monthTrans, _ := h.transactionRepo.GetAgentMonthlyStats(agentID, now)

	// 本月分润统计
	monthProfit, _ := h.profitRepo.GetAgentMonthlyProfitStats(agentID, now)

	// 代理商信息
	agent, _ := h.agentRepo.FindByIDFull(agentID)

	// 钱包余额
	wallets, _ := h.walletRepo.FindByAgentID(agentID)
	var totalBalance int64
	for _, w := range wallets {
		totalBalance += w.Balance
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"today": gin.H{
				"trans_amount":      todayTrans.TotalAmount,
				"trans_amount_yuan": float64(todayTrans.TotalAmount) / 100,
				"trans_count":       todayTrans.TotalCount,
				"profit_amount":     todayProfit.TotalAmount,
				"profit_amount_yuan": float64(todayProfit.TotalAmount) / 100,
			},
			"month": gin.H{
				"trans_amount":      monthTrans.TotalAmount,
				"trans_amount_yuan": float64(monthTrans.TotalAmount) / 100,
				"trans_count":       monthTrans.TotalCount,
				"profit_amount":     monthProfit.TotalAmount,
				"profit_amount_yuan": float64(monthProfit.TotalAmount) / 100,
			},
			"team": gin.H{
				"direct_agent_count":    agent.DirectAgentCount,
				"direct_merchant_count": agent.DirectMerchantCount,
				"team_agent_count":      agent.TeamAgentCount,
				"team_merchant_count":   agent.TeamMerchantCount,
			},
			"wallet": gin.H{
				"total_balance":      totalBalance,
				"total_balance_yuan": float64(totalBalance) / 100,
			},
		},
	})
}

// GetDashboardCharts 获取图表数据
// @Summary 获取图表数据
// @Description 获取交易和分润趋势图表数据
// @Tags 数据看板
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/dashboard/charts [get]
func (h *DashboardHandler) GetDashboardCharts(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	// 最近7天趋势
	endDate := time.Now().AddDate(0, 0, 1)
	startDate := time.Now().AddDate(0, 0, -6)

	trend, _ := h.transactionRepo.GetTransactionTrend(agentID, startDate, endDate)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"trans_trend": trend,
		},
	})
}

// RegisterDashboardRoutes 注册数据看板路由
func RegisterDashboardRoutes(r *gin.RouterGroup, h *DashboardHandler, authService *service.AuthService) {
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(authService))
	{
		dashboard.GET("/overview", h.GetDashboardOverview)
		dashboard.GET("/charts", h.GetDashboardCharts)
	}
}
