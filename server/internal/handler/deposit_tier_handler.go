package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"
)

// DepositTierHandler 押金档位处理器
type DepositTierHandler struct {
	service *service.DepositTierService
}

// NewDepositTierHandler 创建押金档位处理器
func NewDepositTierHandler(service *service.DepositTierService) *DepositTierHandler {
	return &DepositTierHandler{service: service}
}

// RegisterRoutes 注册路由
func (h *DepositTierHandler) RegisterRoutes(rg *gin.RouterGroup) {
	depositTiers := rg.Group("/deposit-tiers")
	{
		depositTiers.GET("", h.List)
		depositTiers.GET("/:id", h.GetByID)
		depositTiers.POST("", h.Create)
		depositTiers.PUT("/:id", h.Update)
		depositTiers.DELETE("/:id", h.Delete)
	}

	// 通道下的押金档位
	channels := rg.Group("/channels")
	{
		channels.GET("/:channelId/deposit-tiers", h.ListByChannel)
	}
}

// List 获取押金档位列表
// @Summary 获取押金档位列表
// @Tags 押金档位
// @Accept json
// @Produce json
// @Param channel_id query int false "通道ID"
// @Param brand_code query string false "品牌编码"
// @Success 200 {object} Response
// @Router /v1/deposit-tiers [get]
func (h *DepositTierHandler) List(c *gin.Context) {
	channelIDStr := c.Query("channel_id")
	brandCode := c.Query("brand_code")

	if channelIDStr == "" {
		response.BadRequest(c, "channel_id is required")
		return
	}

	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid channel_id")
		return
	}

	var tiers interface{}
	if brandCode != "" {
		tiers, err = h.service.GetByChannelAndBrand(channelID, brandCode)
	} else {
		tiers, err = h.service.GetByChannelID(channelID)
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tiers)
}

// ListByChannel 根据通道ID获取押金档位列表
// @Summary 根据通道ID获取押金档位列表
// @Tags 押金档位
// @Accept json
// @Produce json
// @Param channelId path int true "通道ID"
// @Param brand_code query string false "品牌编码"
// @Success 200 {object} Response
// @Router /v1/channels/{channelId}/deposit-tiers [get]
func (h *DepositTierHandler) ListByChannel(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	brandCode := c.Query("brand_code")

	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid channel_id")
		return
	}

	var tiers interface{}
	if brandCode != "" {
		tiers, err = h.service.GetByChannelAndBrand(channelID, brandCode)
	} else {
		tiers, err = h.service.GetByChannelID(channelID)
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tiers)
}

// GetByID 根据ID获取押金档位
// @Summary 根据ID获取押金档位
// @Tags 押金档位
// @Accept json
// @Produce json
// @Param id path int true "档位ID"
// @Success 200 {object} Response
// @Router /v1/deposit-tiers/{id} [get]
func (h *DepositTierHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	tier, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, tier)
}

// Create 创建押金档位
// @Summary 创建押金档位
// @Tags 押金档位
// @Accept json
// @Produce json
// @Param body body service.CreateDepositTierRequest true "创建请求"
// @Success 200 {object} Response
// @Router /v1/deposit-tiers [post]
func (h *DepositTierHandler) Create(c *gin.Context) {
	var req service.CreateDepositTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tier, err := h.service.Create(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, tier, "创建成功")
}

// Update 更新押金档位
// @Summary 更新押金档位
// @Tags 押金档位
// @Accept json
// @Produce json
// @Param id path int true "档位ID"
// @Param body body service.UpdateDepositTierRequest true "更新请求"
// @Success 200 {object} Response
// @Router /v1/deposit-tiers/{id} [put]
func (h *DepositTierHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdateDepositTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tier, err := h.service.Update(id, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, tier, "更新成功")
}

// Delete 删除押金档位
// @Summary 删除押金档位
// @Tags 押金档位
// @Accept json
// @Produce json
// @Param id path int true "档位ID"
// @Success 200 {object} Response
// @Router /v1/deposit-tiers/{id} [delete]
func (h *DepositTierHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}
