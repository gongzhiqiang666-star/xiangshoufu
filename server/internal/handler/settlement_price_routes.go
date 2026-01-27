package handler

import (
	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
)

// RegisterSettlementPriceRoutes 注册结算价路由
func RegisterSettlementPriceRoutes(rg *gin.RouterGroup, h *SettlementPriceHandler, authService *service.AuthService) {
	group := rg.Group("/settlement-prices")
	group.Use(middleware.AuthMiddleware(authService))
	{
		group.GET("", h.List)
		group.GET("/:id", h.GetByID)
		group.POST("", h.Create)
		group.PUT("/:id/rate", h.UpdateRate)
		group.PUT("/:id/deposit", h.UpdateDeposit)
		group.PUT("/:id/sim", h.UpdateSim)
		group.GET("/:id/change-logs", h.GetChangeLogs)
	}
}

// RegisterAgentRewardSettingRoutes 注册代理商奖励配置路由
func RegisterAgentRewardSettingRoutes(rg *gin.RouterGroup, h *AgentRewardSettingHandler, authService *service.AuthService) {
	group := rg.Group("/reward-settings")
	group.Use(middleware.AuthMiddleware(authService))
	{
		group.GET("", h.List)
		group.GET("/:id", h.GetByID)
		group.POST("", h.Create)
		group.PUT("/:id/activation", h.UpdateActivation)
		group.GET("/:id/change-logs", h.GetChangeLogs)
	}
}

// RegisterPriceChangeLogRoutes 注册调价记录路由
func RegisterPriceChangeLogRoutes(rg *gin.RouterGroup, h *PriceChangeLogHandler, authService *service.AuthService) {
	group := rg.Group("/price-change-logs")
	group.Use(middleware.AuthMiddleware(authService))
	{
		group.GET("", h.List)
		group.GET("/:id", h.GetByID)
	}

	// 代理商调价记录
	agentGroup := rg.Group("/agents")
	agentGroup.Use(middleware.AuthMiddleware(authService))
	{
		agentGroup.GET("/:agent_id/price-change-logs", h.ListByAgent)
	}
}
