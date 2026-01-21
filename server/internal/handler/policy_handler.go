package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// PolicyHandler 政策处理器
type PolicyHandler struct {
	policyTemplateRepo *repository.GormPolicyTemplateRepository
	agentPolicyRepo    *repository.GormAgentPolicyRepository
	policyService      *service.PolicyService
}

// NewPolicyHandler 创建政策处理器
func NewPolicyHandler(
	policyTemplateRepo *repository.GormPolicyTemplateRepository,
	agentPolicyRepo *repository.GormAgentPolicyRepository,
	policyService *service.PolicyService,
) *PolicyHandler {
	return &PolicyHandler{
		policyTemplateRepo: policyTemplateRepo,
		agentPolicyRepo:    agentPolicyRepo,
		policyService:      policyService,
	}
}

// GetPolicyTemplates 获取政策模板列表
// @Summary 获取政策模板列表
// @Description 获取系统的政策模板列表
// @Tags 政策管理
// @Produce json
// @Security ApiKeyAuth
// @Param channel_id query int false "通道ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/templates [get]
func (h *PolicyHandler) GetPolicyTemplates(c *gin.Context) {
	channelID, _ := strconv.ParseInt(c.DefaultQuery("channel_id", "1"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	list, total, err := h.policyService.GetPolicyTemplateList(channelID, page, pageSize)
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

// GetPolicyTemplateDetail 获取政策模板详情（含4块政策配置）
// @Summary 获取政策模板详情
// @Description 获取指定政策模板的详细信息，包含费率、押金返现、流量返现、激活奖励配置
// @Tags 政策管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "模板ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/templates/{id} [get]
func (h *PolicyHandler) GetPolicyTemplateDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	detail, err := h.policyService.GetPolicyTemplateDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "模板不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    detail,
	})
}

// CreatePolicyTemplate 创建政策模板
// @Summary 创建政策模板
// @Description 创建新的政策模板，包含4块政策配置
// @Tags 政策管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreatePolicyTemplateRequest true "政策模板"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/templates [post]
func (h *PolicyHandler) CreatePolicyTemplate(c *gin.Context) {
	var req service.CreatePolicyTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := h.policyService.CreatePolicyTemplate(&req)
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

// UpdatePolicyTemplate 更新政策模板
// @Summary 更新政策模板
// @Description 更新指定政策模板的配置
// @Tags 政策管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "模板ID"
// @Param request body service.CreatePolicyTemplateRequest true "政策模板"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/templates/{id} [put]
func (h *PolicyHandler) UpdatePolicyTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req service.CreatePolicyTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := h.policyService.UpdatePolicyTemplate(id, &req)
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

// GetMyPolicies 获取我的政策（只读）
// @Summary 获取我的政策
// @Description 获取当前代理商的完整政策配置（只读，不可修改）
// @Tags 政策管理
// @Produce json
// @Security ApiKeyAuth
// @Param channel_id query int false "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/my [get]
func (h *PolicyHandler) GetMyPolicies(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	channelID, _ := strconv.ParseInt(c.DefaultQuery("channel_id", "1"), 10, 64)

	policy, err := h.policyService.GetAgentPolicy(agentID, channelID)
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
		"data":    policy,
	})
}

// AssignAgentPolicy 分配政策给代理商
// @Summary 分配政策给代理商
// @Description 给指定代理商分配政策模板和个性化配置
// @Tags 政策管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "代理商ID"
// @Param request body service.AssignAgentPolicyRequest true "政策分配"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/{id}/policies [post]
func (h *PolicyHandler) AssignAgentPolicy(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的代理商ID",
		})
		return
	}

	var req service.AssignAgentPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	req.AgentID = agentID

	operatorID := middleware.GetCurrentAgentID(c)
	if err := h.policyService.AssignAgentPolicy(&req, operatorID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "分配成功",
	})
}

// GetSubordinatePolicy 获取下级代理商政策（APP用）
// @Summary 获取下级代理商政策
// @Description 获取指定下级代理商的政策配置
// @Tags 政策管理-APP
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "下级代理商ID"
// @Param channel_id query int false "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/subordinates/{id}/policy [get]
func (h *PolicyHandler) GetSubordinatePolicy(c *gin.Context) {
	subordinateIDStr := c.Param("id")
	subordinateID, err := strconv.ParseInt(subordinateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的代理商ID",
		})
		return
	}

	channelID, _ := strconv.ParseInt(c.DefaultQuery("channel_id", "1"), 10, 64)

	// 获取下级政策
	policy, err := h.policyService.GetAgentPolicy(subordinateID, channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 获取当前代理商的政策限制
	operatorID := middleware.GetCurrentAgentID(c)
	limits, _ := h.policyService.GetPolicyLimits(operatorID, channelID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"policy": policy,
			"limits": limits,
		},
	})
}

// UpdateSubordinatePolicy 更新下级代理商政策（APP用）
// @Summary 更新下级代理商政策
// @Description 更新指定下级代理商的政策配置（在自己政策范围内调整）
// @Tags 政策管理-APP
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "下级代理商ID"
// @Param request body service.UpdateSubordinatePolicyRequest true "政策配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/subordinates/{id}/policy [put]
func (h *PolicyHandler) UpdateSubordinatePolicy(c *gin.Context) {
	subordinateIDStr := c.Param("id")
	subordinateID, err := strconv.ParseInt(subordinateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的代理商ID",
		})
		return
	}

	var req service.UpdateSubordinatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	req.SubordinateID = subordinateID

	operatorID := middleware.GetCurrentAgentID(c)
	if err := h.policyService.UpdateSubordinatePolicy(operatorID, &req); err != nil {
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

// GetPolicyLimits 获取政策限制（当前代理商可设置的范围）
// @Summary 获取政策限制
// @Description 获取当前代理商可为下级设置的政策范围
// @Tags 政策管理-APP
// @Produce json
// @Security ApiKeyAuth
// @Param channel_id query int false "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/limits [get]
func (h *PolicyHandler) GetPolicyLimits(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	channelID, _ := strconv.ParseInt(c.DefaultQuery("channel_id", "1"), 10, 64)

	limits, err := h.policyService.GetPolicyLimits(agentID, channelID)
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
		"data":    limits,
	})
}

// RegisterPolicyRoutes 注册政策路由
func RegisterPolicyRoutes(r *gin.RouterGroup, h *PolicyHandler, authService *service.AuthService) {
	policies := r.Group("/policies")
	policies.Use(middleware.AuthMiddleware(authService))
	{
		// 政策模板管理
		policies.GET("/templates", h.GetPolicyTemplates)
		policies.POST("/templates", h.CreatePolicyTemplate)
		policies.GET("/templates/:id", h.GetPolicyTemplateDetail)
		policies.PUT("/templates/:id", h.UpdatePolicyTemplate)

		// 我的政策（只读）
		policies.GET("/my", h.GetMyPolicies)

		// 政策限制
		policies.GET("/limits", h.GetPolicyLimits)
	}

	// 代理商政策分配
	agents := r.Group("/agents")
	agents.Use(middleware.AuthMiddleware(authService))
	{
		agents.POST("/:id/policies", h.AssignAgentPolicy)
	}

	// 下级代理商政策管理（APP端用）
	subordinates := r.Group("/subordinates")
	subordinates.Use(middleware.AuthMiddleware(authService))
	{
		subordinates.GET("/:id/policy", h.GetSubordinatePolicy)
		subordinates.PUT("/:id/policy", h.UpdateSubordinatePolicy)
	}
}
