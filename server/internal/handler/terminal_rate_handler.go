package handler

import (
	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// TerminalRateHandler 终端费率处理器
type TerminalRateHandler struct {
	terminalService *service.TerminalService
}

// NewTerminalRateHandler 创建终端费率处理器
func NewTerminalRateHandler(terminalService *service.TerminalService) *TerminalRateHandler {
	return &TerminalRateHandler{
		terminalService: terminalService,
	}
}

// UpdateTerminalRateRequest 更新终端费率请求
type UpdateTerminalRateRequest struct {
	CreditRate   int `json:"credit_rate" binding:"required"`   // 贷记卡费率（万分比）
	DebitRate    int `json:"debit_rate" binding:"required"`    // 借记卡费率（万分比）
	DebitCap     int `json:"debit_cap"`                        // 借记卡封顶（分）
	UnionpayRate int `json:"unionpay_rate"`                    // 云闪付费率（万分比）
	WechatRate   int `json:"wechat_rate"`                      // 微信费率（万分比）
	AlipayRate   int `json:"alipay_rate"`                      // 支付宝费率（万分比）
}

// UpdateTerminalRate 更新单个终端费率
// @Summary 更新终端费率
// @Description 调整指定机具的结算价（费率）
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sn path string true "终端SN"
// @Param request body UpdateTerminalRateRequest true "费率信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/{sn}/rate [put]
func (h *TerminalRateHandler) UpdateTerminalRate(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	userID := middleware.GetCurrentUserID(c)

	terminalSN := c.Param("sn")
	if terminalSN == "" {
		response.BadRequest(c, "终端SN不能为空")
		return
	}

	var req UpdateTerminalRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.CreditRate < 0 || req.CreditRate > 1000 {
		response.BadRequest(c, "贷记卡费率范围无效（0-1000万分比）")
		return
	}
	if req.DebitRate < 0 || req.DebitRate > 1000 {
		response.BadRequest(c, "借记卡费率范围无效（0-1000万分比）")
		return
	}

	batchReq := &service.BatchSetRateRequest{
		TerminalSNs:  []string{terminalSN},
		AgentID:      agentID,
		CreditRate:   req.CreditRate,
		DebitRate:    req.DebitRate,
		DebitCap:     req.DebitCap,
		UnionpayRate: req.UnionpayRate,
		WechatRate:   req.WechatRate,
		AlipayRate:   req.AlipayRate,
		UpdatedBy:    userID,
	}

	result, err := h.terminalService.BatchSetRate(batchReq)
	if err != nil {
		response.InternalError(c, "费率更新失败: "+err.Error())
		return
	}

	if result.FailedCount > 0 {
		response.ErrorWithData(c, 400, result.Errors[0], gin.H{
			"success_count": result.SuccessCount,
			"failed_count":  result.FailedCount,
			"errors":        result.Errors,
		})
		return
	}

	response.Success(c, gin.H{
		"terminal_sn": terminalSN,
		"credit_rate": req.CreditRate,
		"debit_rate":  req.DebitRate,
	})
}

// GetTerminalRate 获取终端费率
// @Summary 获取终端费率
// @Description 获取指定机具的当前费率设置
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param sn path string true "终端SN"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/{sn}/rate [get]
func (h *TerminalRateHandler) GetTerminalRate(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	terminalSN := c.Param("sn")
	if terminalSN == "" {
		response.BadRequest(c, "终端SN不能为空")
		return
	}

	policy, err := h.terminalService.GetTerminalPolicy(terminalSN, agentID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"terminal_sn":   policy.TerminalSN,
		"credit_rate":   policy.CreditRate,
		"debit_rate":    policy.DebitRate,
		"debit_cap":     policy.DebitCap,
		"unionpay_rate": policy.UnionpayRate,
		"wechat_rate":   policy.WechatRate,
		"alipay_rate":   policy.AlipayRate,
		"is_synced":     policy.IsSynced,
	})
}

// RegisterTerminalRateRoutes 注册终端费率路由
func RegisterTerminalRateRoutes(r *gin.RouterGroup, h *TerminalRateHandler, authService *service.AuthService) {
	terminals := r.Group("/terminals")
	terminals.Use(middleware.AuthMiddleware(authService))
	{
		terminals.GET("/:sn/rate", h.GetTerminalRate)
		terminals.PUT("/:sn/rate", h.UpdateTerminalRate)
	}
}
