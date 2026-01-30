package handler

import (
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// ChannelConfigHandler 通道配置处理器
type ChannelConfigHandler struct {
	svc service.ChannelConfigService
}

// NewChannelConfigHandler 创建通道配置处理器实例
func NewChannelConfigHandler(svc service.ChannelConfigService) *ChannelConfigHandler {
	return &ChannelConfigHandler{svc: svc}
}

// ============================================================
// 费率配置接口
// ============================================================

// GetRateConfigs 获取费率配置列表
// @Summary 获取通道费率配置列表
// @Tags 通道配置
// @Produce json
// @Param id path int true "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/rate-configs [get]
func (h *ChannelConfigHandler) GetRateConfigs(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	configs, err := h.svc.GetRateConfigs(c.Request.Context(), channelID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, configs)
}

// CreateRateConfig 创建费率配置
// @Summary 创建通道费率配置
// @Tags 通道配置
// @Accept json
// @Produce json
// @Param id path int true "通道ID"
// @Param body body models.CreateChannelRateConfigRequest true "费率配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/rate-configs [post]
func (h *ChannelConfigHandler) CreateRateConfig(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	var req models.CreateChannelRateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config, err := h.svc.CreateRateConfig(c.Request.Context(), channelID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, config)
}

// UpdateRateConfig 更新费率配置
// @Summary 更新通道费率配置
// @Tags 通道配置
// @Accept json
// @Produce json
// @Param id path int true "通道ID"
// @Param configId path int true "配置ID"
// @Param body body models.UpdateChannelRateConfigRequest true "费率配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/rate-configs/{configId} [put]
func (h *ChannelConfigHandler) UpdateRateConfig(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	configID, err := strconv.ParseInt(c.Param("configId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的配置ID")
		return
	}

	var req models.UpdateChannelRateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateRateConfig(c.Request.Context(), channelID, configID, &req); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

// DeleteRateConfig 删除费率配置
// @Summary 删除通道费率配置
// @Tags 通道配置
// @Param id path int true "通道ID"
// @Param configId path int true "配置ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/rate-configs/{configId} [delete]
func (h *ChannelConfigHandler) DeleteRateConfig(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	configID, err := strconv.ParseInt(c.Param("configId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的配置ID")
		return
	}

	if err := h.svc.DeleteRateConfig(c.Request.Context(), channelID, configID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

// ============================================================
// 押金档位接口
// ============================================================

// GetDepositTiers 获取押金档位列表
// @Summary 获取通道押金档位列表
// @Tags 通道配置
// @Produce json
// @Param id path int true "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/deposit-tiers [get]
func (h *ChannelConfigHandler) GetDepositTiers(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	tiers, err := h.svc.GetDepositTiers(c.Request.Context(), channelID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tiers)
}

// UpdateDepositTier 更新押金档位
// @Summary 更新通道押金档位
// @Tags 通道配置
// @Accept json
// @Produce json
// @Param id path int true "通道ID"
// @Param tierId path int true "档位ID"
// @Param body body models.UpdateChannelDepositTierRequest true "押金档位"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/deposit-tiers/{tierId} [put]
func (h *ChannelConfigHandler) UpdateDepositTier(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	tierID, err := strconv.ParseInt(c.Param("tierId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的档位ID")
		return
	}

	var req models.UpdateChannelDepositTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateDepositTier(c.Request.Context(), channelID, tierID, &req); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

// ============================================================
// 流量费返现档位接口
// ============================================================

// GetSimCashbackTiers 获取流量费返现档位列表
// @Summary 获取通道流量费返现档位列表
// @Tags 通道配置
// @Produce json
// @Param id path int true "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/sim-cashback-tiers [get]
func (h *ChannelConfigHandler) GetSimCashbackTiers(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	tiers, err := h.svc.GetSimCashbackTiers(c.Request.Context(), channelID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tiers)
}

// BatchSetSimCashbackTiers 批量设置流量费返现档位
// @Summary 批量设置通道流量费返现档位
// @Tags 通道配置
// @Accept json
// @Produce json
// @Param id path int true "通道ID"
// @Param body body models.BatchSetSimCashbackTiersRequest true "档位列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/sim-cashback-tiers/batch [post]
func (h *ChannelConfigHandler) BatchSetSimCashbackTiers(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	var req models.BatchSetSimCashbackTiersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchSetSimCashbackTiers(c.Request.Context(), channelID, &req); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "设置成功")
}

// ============================================================
// 通道完整配置接口
// ============================================================

// GetFullConfig 获取通道完整配置
// @Summary 获取通道完整配置（费率+押金+流量费）
// @Tags 通道配置
// @Produce json
// @Param id path int true "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/channels/{id}/full-config [get]
func (h *ChannelConfigHandler) GetFullConfig(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	config, err := h.svc.GetFullConfig(c.Request.Context(), channelID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, config)
}
