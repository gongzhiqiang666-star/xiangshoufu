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
}

// NewAgentHandler 创建代理商处理器
func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
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

// RegisterAgentRoutes 注册代理商路由
func RegisterAgentRoutes(r *gin.RouterGroup, h *AgentHandler, authService *service.AuthService) {
	agents := r.Group("/agents")
	agents.Use(middleware.AuthMiddleware(authService))
	{
		agents.GET("/detail", h.GetAgentDetail)
		agents.GET("/subordinates", h.GetSubordinateList)
		agents.GET("/team-tree", h.GetTeamTree)
		agents.GET("/stats", h.GetAgentStats)
		agents.PUT("/profile", h.UpdateProfile)
		agents.GET("/invite-code", h.GetInviteCode)
	}
}
