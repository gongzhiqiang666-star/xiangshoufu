package handler

import (
	"net/http"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// WalletHandler 钱包处理器
type WalletHandler struct {
	walletService *service.WalletService
}

// NewWalletHandler 创建钱包处理器
func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// GetWalletList 获取钱包列表
// @Summary 获取钱包列表
// @Description 获取当前代理商的所有钱包及余额
// @Tags 钱包管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallets [get]
func (h *WalletHandler) GetWalletList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	list, err := h.walletService.GetWalletList(agentID)
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
			"list": list,
		},
	})
}

// GetWalletSummary 获取钱包汇总
// @Summary 获取钱包汇总
// @Description 获取所有钱包的汇总统计
// @Tags 钱包管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallets/summary [get]
func (h *WalletHandler) GetWalletSummary(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	summary, err := h.walletService.GetWalletSummary(agentID)
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
		"data":    summary,
	})
}

// GetWalletLogs 获取钱包流水
// @Summary 获取钱包流水
// @Description 获取钱包流水明细
// @Tags 钱包管理
// @Produce json
// @Security ApiKeyAuth
// @Param wallet_id query int false "钱包ID"
// @Param log_type query int false "流水类型"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallets/logs [get]
func (h *WalletHandler) GetWalletLogs(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	var req service.GetWalletLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	logs, total, err := h.walletService.GetWalletLogs(agentID, &req)
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
			"list":      logs,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// WithdrawRequest 提现请求
type WalletWithdrawRequest struct {
	WalletID int64 `json:"wallet_id" binding:"required"`
	Amount   int64 `json:"amount" binding:"required,min=100"` // 分
}

// Withdraw 申请提现
// @Summary 申请提现
// @Description 从指定钱包申请提现
// @Tags 钱包管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body WalletWithdrawRequest true "提现请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallets/withdraw [post]
func (h *WalletHandler) Withdraw(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	userID := middleware.GetCurrentUserID(c)

	var req WalletWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	serviceReq := &service.WithdrawRequest{
		AgentID:   agentID,
		WalletID:  req.WalletID,
		Amount:    req.Amount,
		CreatedBy: userID,
	}

	if err := h.walletService.Withdraw(serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "提现申请已提交",
	})
}

// RegisterWalletRoutes 注册钱包路由
func RegisterWalletRoutes(r *gin.RouterGroup, h *WalletHandler, authService *service.AuthService) {
	wallets := r.Group("/wallets")
	wallets.Use(middleware.AuthMiddleware(authService))
	{
		wallets.GET("", h.GetWalletList)
		wallets.GET("/summary", h.GetWalletSummary)
		wallets.GET("/logs", h.GetWalletLogs)
		wallets.POST("/withdraw", h.Withdraw)
	}
}
