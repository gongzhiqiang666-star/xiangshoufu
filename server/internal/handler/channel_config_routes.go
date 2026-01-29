package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterChannelConfigRoutes 注册通道配置相关路由
func RegisterChannelConfigRoutes(router *gin.RouterGroup, handler *ChannelConfigHandler) {
	channels := router.Group("/channels")
	{
		// 费率配置
		channels.GET("/:channelId/rate-configs", handler.GetRateConfigs)
		channels.POST("/:channelId/rate-configs", handler.CreateRateConfig)
		channels.PUT("/:channelId/rate-configs/:configId", handler.UpdateRateConfig)
		channels.DELETE("/:channelId/rate-configs/:configId", handler.DeleteRateConfig)

		// 押金档位
		channels.GET("/:channelId/deposit-tiers", handler.GetDepositTiers)
		channels.PUT("/:channelId/deposit-tiers/:tierId", handler.UpdateDepositTier)

		// 流量费返现档位
		channels.GET("/:channelId/sim-cashback-tiers", handler.GetSimCashbackTiers)
		channels.POST("/:channelId/sim-cashback-tiers/batch", handler.BatchSetSimCashbackTiers)

		// 通道完整配置
		channels.GET("/:channelId/full-config", handler.GetFullConfig)
	}
}
