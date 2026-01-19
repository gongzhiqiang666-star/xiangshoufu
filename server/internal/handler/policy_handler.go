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
}

// NewPolicyHandler 创建政策处理器
func NewPolicyHandler(
	policyTemplateRepo *repository.GormPolicyTemplateRepository,
	agentPolicyRepo *repository.GormAgentPolicyRepository,
) *PolicyHandler {
	return &PolicyHandler{
		policyTemplateRepo: policyTemplateRepo,
		agentPolicyRepo:    agentPolicyRepo,
	}
}

// GetPolicyTemplates 获取政策模板列表
// @Summary 获取政策模板列表
// @Description 获取系统的政策模板列表
// @Tags 政策管理
// @Produce json
// @Security ApiKeyAuth
// @Param channel_id query int false "通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/templates [get]
func (h *PolicyHandler) GetPolicyTemplates(c *gin.Context) {
	var channelID *int64
	if cid := c.Query("channel_id"); cid != "" {
		if v, err := strconv.ParseInt(cid, 10, 64); err == nil {
			channelID = &v
		}
	}

	var status int16 = 1 // 只查启用的
	templates, err := h.policyTemplateRepo.FindAll(channelID, &status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(templates))
	for _, t := range templates {
		list = append(list, gin.H{
			"id":            t.ID,
			"template_name": t.TemplateName,
			"channel_id":    t.ChannelID,
			"is_default":    t.IsDefault,
			"credit_rate":   t.CreditRate,
			"debit_rate":    t.DebitRate,
			"debit_cap":     t.DebitCap,
			"unionpay_rate": t.UnionpayRate,
			"wechat_rate":   t.WechatRate,
			"alipay_rate":   t.AlipayRate,
			"status":        t.Status,
			"created_at":    t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": list,
		},
	})
}

// GetPolicyTemplateDetail 获取政策模板详情
// @Summary 获取政策模板详情
// @Description 获取指定政策模板的详细信息
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

	template, err := h.policyTemplateRepo.FindByID(id)
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
		"data": gin.H{
			"id":            template.ID,
			"template_name": template.TemplateName,
			"channel_id":    template.ChannelID,
			"is_default":    template.IsDefault,
			"credit_rate":   template.CreditRate,
			"debit_rate":    template.DebitRate,
			"debit_cap":     template.DebitCap,
			"unionpay_rate": template.UnionpayRate,
			"wechat_rate":   template.WechatRate,
			"alipay_rate":   template.AlipayRate,
			"status":        template.Status,
			"created_at":    template.CreatedAt,
		},
	})
}

// GetMyPolicies 获取我的政策
// @Summary 获取我的政策
// @Description 获取当前代理商的政策配置
// @Tags 政策管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/policies/my [get]
func (h *PolicyHandler) GetMyPolicies(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	policies, err := h.agentPolicyRepo.FindByAgentID(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(policies))
	for _, p := range policies {
		list = append(list, gin.H{
			"id":          p.ID,
			"channel_id":  p.ChannelID,
			"template_id": p.TemplateID,
			"credit_rate": p.CreditRate,
			"debit_rate":  p.DebitRate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": list,
		},
	})
}

// RegisterPolicyRoutes 注册政策路由
func RegisterPolicyRoutes(r *gin.RouterGroup, h *PolicyHandler, authService *service.AuthService) {
	policies := r.Group("/policies")
	policies.Use(middleware.AuthMiddleware(authService))
	{
		policies.GET("/templates", h.GetPolicyTemplates)
		policies.GET("/templates/:id", h.GetPolicyTemplateDetail)
		policies.GET("/my", h.GetMyPolicies)
	}
}
