package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
)

// RewardHandler 奖励处理器
type RewardHandler struct {
	rewardService *service.RewardService
}

// NewRewardHandler 创建奖励处理器
func NewRewardHandler(rewardService *service.RewardService) *RewardHandler {
	return &RewardHandler{
		rewardService: rewardService,
	}
}

// ============================================================
// 奖励政策模版管理
// ============================================================

// GetRewardTemplates 获取奖励模版列表
// @Summary 获取奖励模版列表
// @Description 获取奖励政策模版列表，支持按启用状态筛选
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param enabled query bool false "是否启用"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/templates [get]
func (h *RewardHandler) GetRewardTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var enabled *bool
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		e := enabledStr == "true" || enabledStr == "1"
		enabled = &e
	}

	list, total, err := h.rewardService.GetRewardTemplateList(enabled, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRewardTemplateDetail 获取奖励模版详情
// @Summary 获取奖励模版详情
// @Description 获取指定奖励模版的详细信息，包含阶段配置
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "模版ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/templates/{id} [get]
func (h *RewardHandler) GetRewardTemplateDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	detail, err := h.rewardService.GetRewardTemplateDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "模版不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    detail,
	})
}

// CreateRewardTemplate 创建奖励模版
// @Summary 创建奖励模版
// @Description 创建新的奖励政策模版，包含阶段配置
// @Tags 奖励管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateRewardTemplateRequest true "奖励模版"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/templates [post]
func (h *RewardHandler) CreateRewardTemplate(c *gin.Context) {
	var req models.CreateRewardTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := h.rewardService.CreateRewardTemplate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    result,
	})
}

// UpdateRewardTemplate 更新奖励模版
// @Summary 更新奖励模版
// @Description 更新指定奖励模版的配置
// @Tags 奖励管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "模版ID"
// @Param request body models.UpdateRewardTemplateRequest true "奖励模版"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/templates/{id} [put]
func (h *RewardHandler) UpdateRewardTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req models.UpdateRewardTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := h.rewardService.UpdateRewardTemplate(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    result,
	})
}

// DeleteRewardTemplate 删除奖励模版
// @Summary 删除奖励模版
// @Description 删除指定奖励模版
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "模版ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/templates/{id} [delete]
func (h *RewardHandler) DeleteRewardTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.rewardService.DeleteRewardTemplate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdateRewardTemplateStatus 更新模版启用状态
// @Summary 更新模版启用状态
// @Description 启用或禁用奖励模版
// @Tags 奖励管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "模版ID"
// @Param request body map[string]bool true "启用状态"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/templates/{id}/status [put]
func (h *RewardHandler) UpdateRewardTemplateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.rewardService.UpdateRewardTemplateEnabled(id, req.Enabled); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// ============================================================
// 代理商奖励比例配置
// ============================================================

// GetAgentRewardRate 获取代理商奖励比例
// @Summary 获取代理商奖励比例
// @Description 获取指定代理商的奖励比例配置
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "代理商ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/agents/{id}/rate [get]
func (h *RewardHandler) GetAgentRewardRate(c *gin.Context) {
	idStr := c.Param("id")
	agentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的代理商ID",
		})
		return
	}

	rate, err := h.rewardService.GetAgentRewardRate(agentID)
	if err != nil {
		// 未配置时返回默认值
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"agent_id":    agentID,
				"reward_rate": 0,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    rate,
	})
}

// SetAgentRewardRate 设置代理商奖励比例
// @Summary 设置代理商奖励比例
// @Description 设置指定代理商的奖励比例配置
// @Tags 奖励管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "代理商ID"
// @Param request body models.AgentRewardRateRequest true "奖励比例"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/agents/{id}/rate [put]
func (h *RewardHandler) SetAgentRewardRate(c *gin.Context) {
	idStr := c.Param("id")
	agentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的代理商ID",
		})
		return
	}

	var req models.AgentRewardRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	req.AgentID = agentID

	if err := h.rewardService.SetAgentRewardRate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
	})
}

// ============================================================
// 终端奖励进度
// ============================================================

// GetTerminalRewardProgress 获取终端奖励进度
// @Summary 获取终端奖励进度
// @Description 获取指定终端的奖励进度详情
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param terminal_sn path string true "终端SN"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/terminals/{terminal_sn}/progress [get]
func (h *RewardHandler) GetTerminalRewardProgress(c *gin.Context) {
	terminalSN := c.Param("terminal_sn")
	if terminalSN == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "终端SN不能为空",
		})
		return
	}

	progress, err := h.rewardService.GetTerminalRewardProgress(terminalSN)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "未找到奖励进度: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    progress,
	})
}

// InitTerminalRewardProgress 初始化终端奖励进度
// @Summary 初始化终端奖励进度
// @Description 为终端初始化奖励进度（终端绑定时调用）
// @Tags 奖励管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body map[string]interface{} true "初始化参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/terminals/progress [post]
func (h *RewardHandler) InitTerminalRewardProgress(c *gin.Context) {
	var req struct {
		TerminalSN string `json:"terminal_sn" binding:"required"`
		TerminalID *int64 `json:"terminal_id"`
		AgentID    int64  `json:"agent_id" binding:"required"`
		TemplateID int64  `json:"template_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	progress, err := h.rewardService.InitTerminalRewardProgress(req.TerminalSN, req.TerminalID, req.AgentID, req.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "初始化成功",
		"data":    progress,
	})
}

// ============================================================
// 溢出日志管理
// ============================================================

// GetOverflowLogs 获取溢出日志列表
// @Summary 获取溢出日志列表
// @Description 获取未解决的奖励池溢出日志
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/overflow-logs [get]
func (h *RewardHandler) GetOverflowLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.rewardService.GetUnresolvedOverflowLogs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ResolveOverflowLog 解决溢出日志
// @Summary 解决溢出日志
// @Description 标记溢出日志为已解决
// @Tags 奖励管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "日志ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rewards/overflow-logs/{id}/resolve [post]
func (h *RewardHandler) ResolveOverflowLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	// 获取操作人
	operatorName := middleware.GetCurrentUsername(c)
	if operatorName == "" {
		operatorName = "system"
	}

	if err := h.rewardService.ResolveOverflowLog(id, operatorName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "解决成功",
	})
}

// ============================================================
// 路由注册
// ============================================================

// RegisterRewardRoutes 注册奖励管理路由
func RegisterRewardRoutes(r *gin.RouterGroup, h *RewardHandler, authService *service.AuthService) {
	rewards := r.Group("/rewards")
	rewards.Use(middleware.AuthMiddleware(authService))
	{
		// 奖励模版管理
		rewards.GET("/templates", h.GetRewardTemplates)
		rewards.POST("/templates", h.CreateRewardTemplate)
		rewards.GET("/templates/:id", h.GetRewardTemplateDetail)
		rewards.PUT("/templates/:id", h.UpdateRewardTemplate)
		rewards.DELETE("/templates/:id", h.DeleteRewardTemplate)
		rewards.PUT("/templates/:id/status", h.UpdateRewardTemplateStatus)

		// 代理商奖励比例
		rewards.GET("/agents/:id/rate", h.GetAgentRewardRate)
		rewards.PUT("/agents/:id/rate", h.SetAgentRewardRate)

		// 终端奖励进度
		rewards.GET("/terminals/:terminal_sn/progress", h.GetTerminalRewardProgress)
		rewards.POST("/terminals/progress", h.InitTerminalRewardProgress)

		// 溢出日志
		rewards.GET("/overflow-logs", h.GetOverflowLogs)
		rewards.POST("/overflow-logs/:id/resolve", h.ResolveOverflowLog)
	}
}
