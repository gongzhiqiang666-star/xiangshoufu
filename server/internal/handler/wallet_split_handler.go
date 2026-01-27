package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
)

// WalletSplitHandler 钱包拆分配置处理器
type WalletSplitHandler struct {
	splitService *service.WalletSplitService
}

// NewWalletSplitHandler 创建处理器
func NewWalletSplitHandler(splitService *service.WalletSplitService) *WalletSplitHandler {
	return &WalletSplitHandler{splitService: splitService}
}

// GetSplitConfig 获取代理商钱包拆分配置
// @Summary 获取代理商钱包拆分配置
// @Tags 钱包拆分配置
// @Produce json
// @Param id path int true "代理商ID"
// @Success 200 {object} models.AgentWalletSplitConfig
// @Router /api/v1/agents/{id}/wallet-split [get]
func (h *WalletSplitHandler) GetSplitConfig(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的代理商ID"})
		return
	}

	config, err := h.splitService.GetSplitConfig(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config,
	})
}

// SetSplitConfigRequest 设置拆分配置请求
type SetSplitConfigRequest struct {
	SplitByChannel bool `json:"split_by_channel"`
}

// SetSplitConfig 设置代理商钱包拆分配置
// @Summary 设置代理商钱包拆分配置
// @Tags 钱包拆分配置
// @Accept json
// @Produce json
// @Param id path int true "代理商ID"
// @Param request body SetSplitConfigRequest true "拆分配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/{id}/wallet-split [put]
func (h *WalletSplitHandler) SetSplitConfig(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的代理商ID"})
		return
	}

	var req SetSplitConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取当前操作者ID（从JWT或session中获取）
	configuredBy := getOperatorID(c)

	serviceReq := &service.SetSplitConfigRequest{
		AgentID:        agentID,
		SplitByChannel: req.SplitByChannel,
		ConfiguredBy:   configuredBy,
	}

	if err := h.splitService.SetSplitConfig(serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
	})
}

// CheckSplitStatus 检查代理商是否按通道拆分
// @Summary 检查代理商是否按通道拆分
// @Tags 钱包拆分配置
// @Produce json
// @Param id path int true "代理商ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/{id}/wallet-split/status [get]
func (h *WalletSplitHandler) CheckSplitStatus(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的代理商ID"})
		return
	}

	isSplit, err := h.splitService.IsSplitByChannel(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"split_by_channel": isSplit,
		},
	})
}

// ============================================================
// 提现门槛配置
// ============================================================

// GetWithdrawThresholds 获取政策模版提现门槛
// @Summary 获取政策模版提现门槛配置
// @Tags 提现门槛配置
// @Produce json
// @Param id path int true "政策模版ID"
// @Success 200 {array} models.PolicyWithdrawThreshold
// @Router /api/v1/policy-templates/{id}/withdraw-thresholds [get]
func (h *WalletSplitHandler) GetWithdrawThresholds(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的政策模版ID"})
		return
	}

	thresholds, err := h.splitService.GetWithdrawThresholds(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": thresholds,
	})
}

// SetWithdrawThresholdRequest 设置提现门槛请求
type SetWithdrawThresholdRequest struct {
	WalletType      int16 `json:"wallet_type" binding:"required,min=1,max=3"`
	ChannelID       int64 `json:"channel_id"`
	ThresholdAmount int64 `json:"threshold_amount" binding:"required,min=100"`
}

// SetWithdrawThreshold 设置政策模版提现门槛
// @Summary 设置政策模版提现门槛
// @Tags 提现门槛配置
// @Accept json
// @Produce json
// @Param id path int true "政策模版ID"
// @Param request body SetWithdrawThresholdRequest true "门槛配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policy-templates/{id}/withdraw-thresholds [put]
func (h *WalletSplitHandler) SetWithdrawThreshold(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的政策模版ID"})
		return
	}

	var req SetWithdrawThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	serviceReq := &service.SetWithdrawThresholdRequest{
		TemplateID:      templateID,
		WalletType:      req.WalletType,
		ChannelID:       req.ChannelID,
		ThresholdAmount: req.ThresholdAmount,
	}

	if err := h.splitService.SetWithdrawThreshold(serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
	})
}

// BatchSetWithdrawThresholdsRequest 批量设置提现门槛请求
type BatchSetWithdrawThresholdsRequest struct {
	Thresholds []SetWithdrawThresholdRequest `json:"thresholds" binding:"required"`
}

// BatchSetWithdrawThresholds 批量设置政策模版提现门槛
// @Summary 批量设置政策模版提现门槛
// @Tags 提现门槛配置
// @Accept json
// @Produce json
// @Param id path int true "政策模版ID"
// @Param request body BatchSetWithdrawThresholdsRequest true "门槛配置列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policy-templates/{id}/withdraw-thresholds/batch [put]
func (h *WalletSplitHandler) BatchSetWithdrawThresholds(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的政策模版ID"})
		return
	}

	var req BatchSetWithdrawThresholdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 转换请求
	serviceReqs := make([]*service.SetWithdrawThresholdRequest, 0, len(req.Thresholds))
	for _, t := range req.Thresholds {
		serviceReqs = append(serviceReqs, &service.SetWithdrawThresholdRequest{
			TemplateID:      templateID,
			WalletType:      t.WalletType,
			ChannelID:       t.ChannelID,
			ThresholdAmount: t.ThresholdAmount,
		})
	}

	if err := h.splitService.BatchSetWithdrawThresholds(templateID, serviceReqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量设置成功",
	})
}

// RegisterWalletSplitRoutes 注册钱包拆分配置路由
func RegisterWalletSplitRoutes(rg *gin.RouterGroup, h *WalletSplitHandler, authService *service.AuthService) {
	// 代理商钱包拆分配置（需要管理员或上级权限）
	agents := rg.Group("/agents")
	agents.Use(middleware.AuthMiddleware(authService))
	{
		agents.GET("/:id/wallet-split", h.GetSplitConfig)
		agents.PUT("/:id/wallet-split", h.SetSplitConfig)
		agents.GET("/:id/wallet-split/status", h.CheckSplitStatus)
	}

	// 政策模版提现门槛配置（仅管理员）
	policies := rg.Group("/policy-templates")
	policies.Use(middleware.AuthMiddleware(authService))
	{
		policies.GET("/:id/withdraw-thresholds", h.GetWithdrawThresholds)
		policies.PUT("/:id/withdraw-thresholds", h.SetWithdrawThreshold)
		policies.PUT("/:id/withdraw-thresholds/batch", h.BatchSetWithdrawThresholds)
	}
}
