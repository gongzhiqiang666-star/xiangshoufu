package handler

import (
	"net/http"
	"strconv"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// GlobalWithdrawThresholdHandler 全局提现门槛处理器
type GlobalWithdrawThresholdHandler struct {
	svc *service.GlobalWithdrawThresholdService
}

// NewGlobalWithdrawThresholdHandler 创建处理器实例
func NewGlobalWithdrawThresholdHandler(svc *service.GlobalWithdrawThresholdService) *GlobalWithdrawThresholdHandler {
	return &GlobalWithdrawThresholdHandler{svc: svc}
}

// GetThresholds 获取所有门槛配置
// @Summary 获取提现门槛配置
// @Description 获取所有钱包类型的提现门槛配置，包括通用门槛和按通道门槛
// @Tags 提现门槛
// @Produce json
// @Success 200 {object} service.ThresholdListResponse
// @Router /api/v1/withdraw-thresholds [get]
func (h *GlobalWithdrawThresholdHandler) GetThresholds(c *gin.Context) {
	response, err := h.svc.GetAllThresholds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取门槛配置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    response,
	})
}

// SetGeneralThresholds 设置通用门槛
// @Summary 设置通用提现门槛
// @Description 设置各钱包类型的通用提现门槛（适用于所有通道）
// @Tags 提现门槛
// @Accept json
// @Produce json
// @Param request body service.SetGeneralThresholdRequest true "通用门槛配置"
// @Success 200 {object} map[string]string
// @Router /api/v1/withdraw-thresholds/general [put]
func (h *GlobalWithdrawThresholdHandler) SetGeneralThresholds(c *gin.Context) {
	var req service.SetGeneralThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := h.svc.SetGeneralThresholds(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "设置门槛失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置成功", "data": nil})
}

// SetChannelThresholds 设置通道门槛
// @Summary 设置通道提现门槛
// @Description 设置指定通道的提现门槛，优先级高于通用门槛
// @Tags 提现门槛
// @Accept json
// @Produce json
// @Param request body service.SetChannelThresholdRequest true "通道门槛配置"
// @Success 200 {object} map[string]string
// @Router /api/v1/withdraw-thresholds/channel [put]
func (h *GlobalWithdrawThresholdHandler) SetChannelThresholds(c *gin.Context) {
	var req service.SetChannelThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.ChannelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "通道ID不能为空"})
		return
	}

	if err := h.svc.SetChannelThresholds(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "设置门槛失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置成功", "data": nil})
}

// DeleteChannelThreshold 删除通道门槛
// @Summary 删除通道提现门槛
// @Description 删除指定通道的提现门槛配置，删除后将使用通用门槛
// @Tags 提现门槛
// @Produce json
// @Param channel_id path int true "通道ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/withdraw-thresholds/channel/{channel_id} [delete]
func (h *GlobalWithdrawThresholdHandler) DeleteChannelThreshold(c *gin.Context) {
	channelIDStr := c.Param("channel_id")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil || channelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "通道ID无效"})
		return
	}

	if err := h.svc.DeleteChannelThreshold(channelID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除门槛失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
