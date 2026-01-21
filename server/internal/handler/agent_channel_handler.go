package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// AgentChannelHandler 代理商通道处理器
type AgentChannelHandler struct {
	agentChannelService *service.AgentChannelService
}

// NewAgentChannelHandler 创建代理商通道处理器
func NewAgentChannelHandler(agentChannelService *service.AgentChannelService) *AgentChannelHandler {
	return &AgentChannelHandler{
		agentChannelService: agentChannelService,
	}
}

// GetAgentChannels 获取代理商通道列表
// @Summary 获取代理商通道列表
// @Description 获取代理商已配置的通道列表
// @Tags 代理商通道管理
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id query int false "代理商ID（不传则获取当前登录代理商）"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels [get]
func (h *AgentChannelHandler) GetAgentChannels(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	// 如果指定了代理商ID，则查询指定代理商
	if idStr := c.Query("agent_id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的代理商ID",
			})
			return
		}
		agentID = id
	}

	channels, err := h.agentChannelService.GetAgentChannels(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取通道列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    channels,
	})
}

// GetEnabledChannels 获取已启用的通道列表（用于APP端）
// @Summary 获取已启用的通道列表
// @Description 获取代理商已启用且可见的通道列表
// @Tags 代理商通道管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/enabled [get]
func (h *AgentChannelHandler) GetEnabledChannels(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	channels, err := h.agentChannelService.GetEnabledChannels(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取通道列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    channels,
	})
}

// EnableChannelRequest 启用通道请求
type EnableChannelRequest struct {
	AgentID   int64 `json:"agent_id" binding:"required"`
	ChannelID int64 `json:"channel_id" binding:"required"`
}

// EnableChannel 启用代理商通道
// @Summary 启用代理商通道
// @Description 为代理商启用指定通道
// @Tags 代理商通道管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body EnableChannelRequest true "启用通道请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/enable [post]
func (h *AgentChannelHandler) EnableChannel(c *gin.Context) {
	var req EnableChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	operatorID := middleware.GetCurrentUserID(c)
	if err := h.agentChannelService.EnableChannel(req.AgentID, req.ChannelID, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "启用通道失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "启用成功",
	})
}

// DisableChannel 禁用代理商通道
// @Summary 禁用代理商通道
// @Description 为代理商禁用指定通道
// @Tags 代理商通道管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body EnableChannelRequest true "禁用通道请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/disable [post]
func (h *AgentChannelHandler) DisableChannel(c *gin.Context) {
	var req EnableChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	operatorID := middleware.GetCurrentUserID(c)
	if err := h.agentChannelService.DisableChannel(req.AgentID, req.ChannelID, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "禁用通道失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "禁用成功",
	})
}

// SetChannelVisibilityRequest 设置通道可见性请求
type SetChannelVisibilityRequest struct {
	AgentID   int64 `json:"agent_id" binding:"required"`
	ChannelID int64 `json:"channel_id" binding:"required"`
	IsVisible bool  `json:"is_visible"`
}

// SetChannelVisibility 设置通道可见性
// @Summary 设置通道可见性
// @Description 设置代理商某个通道是否对APP端可见
// @Tags 代理商通道管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body SetChannelVisibilityRequest true "设置可见性请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/visibility [post]
func (h *AgentChannelHandler) SetChannelVisibility(c *gin.Context) {
	var req SetChannelVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := h.agentChannelService.SetChannelVisibility(req.AgentID, req.ChannelID, req.IsVisible); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "设置可见性失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
	})
}

// BatchEnableChannelsRequest 批量启用通道请求
type BatchEnableChannelsRequest struct {
	AgentID    int64   `json:"agent_id" binding:"required"`
	ChannelIDs []int64 `json:"channel_ids" binding:"required"`
}

// BatchEnableChannels 批量启用通道
// @Summary 批量启用通道
// @Description 为代理商批量启用多个通道
// @Tags 代理商通道管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchEnableChannelsRequest true "批量启用请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/batch-enable [post]
func (h *AgentChannelHandler) BatchEnableChannels(c *gin.Context) {
	var req BatchEnableChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	operatorID := middleware.GetCurrentUserID(c)
	if err := h.agentChannelService.BatchEnableChannels(req.AgentID, req.ChannelIDs, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量启用失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量启用成功",
	})
}

// BatchDisableChannels 批量禁用通道
// @Summary 批量禁用通道
// @Description 为代理商批量禁用多个通道
// @Tags 代理商通道管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchEnableChannelsRequest true "批量禁用请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/batch-disable [post]
func (h *AgentChannelHandler) BatchDisableChannels(c *gin.Context) {
	var req BatchEnableChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	operatorID := middleware.GetCurrentUserID(c)
	if err := h.agentChannelService.BatchDisableChannels(req.AgentID, req.ChannelIDs, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量禁用失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量禁用成功",
	})
}

// GetAgentChannelStats 获取代理商通道统计
// @Summary 获取代理商通道统计
// @Description 获取代理商通道的统计信息
// @Tags 代理商通道管理
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id query int false "代理商ID（不传则获取当前登录代理商）"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/stats [get]
func (h *AgentChannelHandler) GetAgentChannelStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	if idStr := c.Query("agent_id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的代理商ID",
			})
			return
		}
		agentID = id
	}

	stats, err := h.agentChannelService.GetAgentChannelStats(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// InitAgentChannelsRequest 初始化代理商通道请求
type InitAgentChannelsRequest struct {
	AgentID int64 `json:"agent_id" binding:"required"`
}

// InitAgentChannels 初始化代理商通道配置
// @Summary 初始化代理商通道配置
// @Description 为代理商初始化所有可用通道的配置
// @Tags 代理商通道管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body InitAgentChannelsRequest true "初始化请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-channels/init [post]
func (h *AgentChannelHandler) InitAgentChannels(c *gin.Context) {
	var req InitAgentChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	operatorID := middleware.GetCurrentUserID(c)
	if err := h.agentChannelService.InitAgentChannels(req.AgentID, operatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "初始化失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "初始化成功",
	})
}

// RegisterAgentChannelRoutes 注册代理商通道路由
func RegisterAgentChannelRoutes(router *gin.RouterGroup, handler *AgentChannelHandler, authService *service.AuthService) {
	channels := router.Group("/agent-channels")
	channels.Use(middleware.AuthMiddleware(authService))
	{
		channels.GET("", handler.GetAgentChannels)
		channels.GET("/enabled", handler.GetEnabledChannels)
		channels.GET("/stats", handler.GetAgentChannelStats)
		channels.POST("/enable", handler.EnableChannel)
		channels.POST("/disable", handler.DisableChannel)
		channels.POST("/visibility", handler.SetChannelVisibility)
		channels.POST("/batch-enable", handler.BatchEnableChannels)
		channels.POST("/batch-disable", handler.BatchDisableChannels)
		channels.POST("/init", handler.InitAgentChannels)
	}
}
