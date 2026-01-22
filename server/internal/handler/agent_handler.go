package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// AgentHandler 代理商处理器
type AgentHandler struct {
	agentService *service.AgentService
	auditService *service.AuditService
}

// NewAgentHandler 创建代理商处理器
func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

// SetAuditService 设置审计服务
func (h *AgentHandler) SetAuditService(auditService *service.AuditService) {
	h.auditService = auditService
}

// GetAgentDetail 获取代理商详情
// @Summary 获取代理商详情
// @Description 获取当前登录代理商或指定代理商的详细信息
// @Tags 代理商管理
// @Produce json
// @Security ApiKeyAuth
// @Param id query int false "代理商ID（不传则获取当前登录代理商）"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/detail [get]
func (h *AgentHandler) GetAgentDetail(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	// 如果指定了ID，则查询指定代理商（需要权限检查）
	if idStr := c.Query("id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的ID",
			})
			return
		}

		// 检查是否有权限查看（必须是下级或自己）
		if id != agentID {
			isSubordinate, err := h.agentService.IsSubordinate(agentID, id)
			if err != nil || !isSubordinate {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "无权查看该代理商信息",
				})
				return
			}
		}
		agentID = id
	}

	detail, err := h.agentService.GetAgentDetail(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    detail,
	})
}

// GetSubordinateList 获取下级代理商列表
// @Summary 获取下级代理商列表
// @Description 获取当前代理商的直属下级列表
// @Tags 代理商管理
// @Produce json
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/subordinates [get]
func (h *AgentHandler) GetSubordinateList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	var req service.SubordinateListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	req.AgentID = agentID
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	list, total, err := h.agentService.GetSubordinateList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetTeamTree 获取团队层级树
// @Summary 获取团队层级树
// @Description 获取代理商团队的层级树结构
// @Tags 代理商管理
// @Produce json
// @Security ApiKeyAuth
// @Param max_depth query int false "最大深度（默认3）"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/team-tree [get]
func (h *AgentHandler) GetTeamTree(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	maxDepth := 3
	if depthStr := c.Query("max_depth"); depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil && d > 0 && d <= 5 {
			maxDepth = d
		}
	}

	tree, err := h.agentService.GetTeamTree(agentID, maxDepth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tree,
	})
}

// GetAgentStats 获取代理商统计
// @Summary 获取代理商统计
// @Description 获取代理商的交易、分润、团队等统计数据
// @Tags 代理商管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/stats [get]
func (h *AgentHandler) GetAgentStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	stats, err := h.agentService.GetAgentStats(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// UpdateAgentProfileRequest 更新资料请求
type UpdateAgentProfileRequest struct {
	AgentName    string `json:"agent_name"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	BankName     string `json:"bank_name"`
	BankAccount  string `json:"bank_account"`
	BankCardNo   string `json:"bank_card_no"`
}

// UpdateProfile 更新代理商资料
// @Summary 更新代理商资料
// @Description 更新当前代理商的基本资料
// @Tags 代理商管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateAgentProfileRequest true "更新资料请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/profile [put]
func (h *AgentHandler) UpdateProfile(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	var req UpdateAgentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	serviceReq := &service.UpdateAgentProfileRequest{
		AgentID:      agentID,
		AgentName:    req.AgentName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankCardNo:   req.BankCardNo,
	}

	if err := h.agentService.UpdateAgentProfile(serviceReq); err != nil {
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

// GetInviteCode 获取邀请码
// @Summary 获取邀请码
// @Description 获取当前代理商的邀请码和二维码
// @Tags 代理商管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/invite-code [get]
func (h *AgentHandler) GetInviteCode(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	inviteCode, qrCodeURL, err := h.agentService.GetInviteCode(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"invite_code": inviteCode,
			"qr_code_url": qrCodeURL,
		},
	})
}

// SetCustomInviteCodeRequest 设置自定义邀请码请求
type SetCustomInviteCodeRequest struct {
	InviteCode string `json:"invite_code" binding:"required,min=4,max=12,alphanum"`
}

// SetCustomInviteCode 设置自定义邀请码（靓号）
// @Summary 设置自定义邀请码
// @Description 设置当前代理商的自定义邀请码（靓号），支持4-12位字母数字
// @Tags 代理商管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body SetCustomInviteCodeRequest true "设置邀请码请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/invite-code [put]
func (h *AgentHandler) SetCustomInviteCode(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	var req SetCustomInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	serviceReq := &service.SetCustomInviteCodeRequest{
		AgentID:    agentID,
		InviteCode: req.InviteCode,
	}

	if err := h.agentService.SetCustomInviteCode(serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "邀请码设置成功",
	})
}

// CreateAgentRequest 创建代理商请求
type CreateAgentRequest struct {
	AgentName    string `json:"agent_name" binding:"required"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	IDCardNo     string `json:"id_card_no"`
	BankName     string `json:"bank_name"`
	BankAccount  string `json:"bank_account"`
	BankCardNo   string `json:"bank_card_no"`
	ParentID     int64  `json:"parent_id"`
}

// CreateAgent 创建代理商
// @Summary 创建代理商
// @Description 创建新的代理商（下级代理）
// @Tags 代理商管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateAgentRequest true "创建代理商请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents [post]
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	operatorID := middleware.GetCurrentAgentID(c)

	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 如果没有指定上级，默认为当前操作人
	if req.ParentID == 0 {
		req.ParentID = operatorID
	}

	// 验证权限：只能创建自己的下级
	if req.ParentID != operatorID {
		isSubordinate, err := h.agentService.IsSubordinate(operatorID, req.ParentID)
		if err != nil || !isSubordinate {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权在该代理商下创建下级",
			})
			return
		}
	}

	serviceReq := &service.CreateAgentRequest{
		AgentName:    req.AgentName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		IDCardNo:     req.IDCardNo,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankCardNo:   req.BankCardNo,
		ParentID:     req.ParentID,
	}

	agent, err := h.agentService.CreateAgent(serviceReq, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 记录创建代理商审计日志
	if h.auditService != nil {
		auditCtx := service.NewAuditContextFromGin(c)
		h.auditService.LogAgentCreate(auditCtx, agent.ID, agent.AgentName)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    agent,
	})
}

// UpdateAgentStatusRequest 更新状态请求
type UpdateAgentStatusRequest struct {
	Status int16 `json:"status" binding:"required,oneof=1 2"`
}

// UpdateAgentStatus 更新代理商状态
// @Summary 更新代理商状态
// @Description 启用或禁用代理商
// @Tags 代理商管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "代理商ID"
// @Param request body UpdateAgentStatusRequest true "更新状态请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/{id}/status [put]
func (h *AgentHandler) UpdateAgentStatus(c *gin.Context) {
	operatorID := middleware.GetCurrentAgentID(c)

	agentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	// 检查权限：只能操作自己的下级
	if agentID != operatorID {
		isSubordinate, err := h.agentService.IsSubordinate(operatorID, agentID)
		if err != nil || !isSubordinate {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权操作该代理商",
			})
			return
		}
	}

	var req UpdateAgentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	serviceReq := &service.UpdateAgentStatusRequest{
		AgentID: agentID,
		Status:  req.Status,
	}

	if err := h.agentService.UpdateAgentStatus(serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 记录代理商状态变更审计日志
	if h.auditService != nil {
		auditCtx := service.NewAuditContextFromGin(c)
		// 获取代理商信息用于审计
		detail, _ := h.agentService.GetAgentDetail(agentID)
		agentName := ""
		if detail != nil {
			agentName = detail.AgentName
		}
		oldStatus := int16(1)
		if req.Status == 1 {
			oldStatus = 2
		}
		h.auditService.LogAgentStatusChange(auditCtx, agentID, agentName, oldStatus, req.Status)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// SearchAgents 搜索代理商
// @Summary 搜索代理商
// @Description 全局搜索代理商（用于选择器等场景）
// @Tags 代理商管理
// @Produce json
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agents/search [get]
func (h *AgentHandler) SearchAgents(c *gin.Context) {
	var req service.SearchAgentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	list, total, err := h.agentService.SearchAgents(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// RegisterAgentRoutes 注册代理商路由
func RegisterAgentRoutes(r *gin.RouterGroup, h *AgentHandler, authService *service.AuthService) {
	agents := r.Group("/agents")
	agents.Use(middleware.AuthMiddleware(authService))
	{
		agents.POST("", h.CreateAgent)
		agents.GET("/detail", h.GetAgentDetail)
		agents.GET("/subordinates", h.GetSubordinateList)
		agents.GET("/team-tree", h.GetTeamTree)
		agents.GET("/stats", h.GetAgentStats)
		agents.PUT("/profile", h.UpdateProfile)
		agents.GET("/invite-code", h.GetInviteCode)
		agents.PUT("/invite-code", h.SetCustomInviteCode)
		agents.GET("/search", h.SearchAgents)
		agents.PUT("/:id/status", h.UpdateAgentStatus)
	}
}
