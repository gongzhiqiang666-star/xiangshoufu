package handler

import (
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// SettlementWalletHandler 沉淀钱包处理器
type SettlementWalletHandler struct {
	settlementService *service.SettlementWalletService
}

// NewSettlementWalletHandler 创建沉淀钱包处理器
func NewSettlementWalletHandler(settlementService *service.SettlementWalletService) *SettlementWalletHandler {
	return &SettlementWalletHandler{
		settlementService: settlementService,
	}
}

// EnableSettlementWalletRequest 开通沉淀钱包请求
type EnableSettlementWalletRequest struct {
	AgentID int64 `json:"agent_id" binding:"required"`
	Ratio   int   `json:"ratio" binding:"required,min=1,max=100"` // 沉淀比例(1-100)
}

// EnableSettlementWallet 开通沉淀钱包
// @Summary 开通沉淀钱包
// @Description PC端管理员为代理商开通沉淀钱包
// @Tags 沉淀钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body EnableSettlementWalletRequest true "开通请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/enable [post]
func (h *SettlementWalletHandler) EnableSettlementWallet(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req EnableSettlementWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.EnableSettlementWalletRequest{
		AgentID:   req.AgentID,
		Ratio:     req.Ratio,
		EnabledBy: userID,
	}

	if err := h.settlementService.EnableSettlementWallet(serviceReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "开通成功")
}

// DisableSettlementWallet 关闭沉淀钱包
// @Summary 关闭沉淀钱包
// @Description PC端管理员关闭代理商的沉淀钱包
// @Tags 沉淀钱包
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id path int true "代理商ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/disable/{agent_id} [post]
func (h *SettlementWalletHandler) DisableSettlementWallet(c *gin.Context) {
	agentIDStr := c.Param("agent_id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的代理商ID")
		return
	}

	if err := h.settlementService.DisableSettlementWallet(agentID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "关闭成功")
}

// UpdateRatioRequest 更新沉淀比例请求
type UpdateRatioRequest struct {
	Ratio int `json:"ratio" binding:"required,min=1,max=100"`
}

// UpdateSettlementRatio 更新沉淀比例
// @Summary 更新沉淀比例
// @Description 更新代理商的沉淀比例
// @Tags 沉淀钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id path int true "代理商ID"
// @Param request body UpdateRatioRequest true "比例请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/ratio/{agent_id} [put]
func (h *SettlementWalletHandler) UpdateSettlementRatio(c *gin.Context) {
	agentIDStr := c.Param("agent_id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的代理商ID")
		return
	}

	var req UpdateRatioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.settlementService.UpdateSettlementRatio(agentID, req.Ratio); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

// GetSummary 获取沉淀钱包汇总
// @Summary 获取沉淀钱包汇总
// @Description 获取当前代理商的沉淀钱包汇总信息
// @Tags 沉淀钱包
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/summary [get]
func (h *SettlementWalletHandler) GetSummary(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	summary, err := h.settlementService.GetSettlementWalletSummary(agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, summary)
}

// GetSubordinateBalances 获取下级余额明细
// @Summary 获取下级余额明细
// @Description 获取所有直属下级的可用余额明细
// @Tags 沉淀钱包
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/subordinates [get]
func (h *SettlementWalletHandler) GetSubordinateBalances(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	balances, err := h.settlementService.GetSubordinateBalances(agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list": balances,
	})
}

// UseSettlementRequest 使用沉淀款请求
type UseSettlementRequest struct {
	Amount int64  `json:"amount" binding:"required,min=100"` // 最少1元(100分)
	Remark string `json:"remark"`
}

// UseSettlement 使用沉淀款
// @Summary 使用沉淀款
// @Description 使用沉淀钱包中的额度
// @Tags 沉淀钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UseSettlementRequest true "使用请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/use [post]
func (h *SettlementWalletHandler) UseSettlement(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	userID := middleware.GetCurrentUserID(c)

	var req UseSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.UseSettlementRequest{
		AgentID:   agentID,
		Amount:    req.Amount,
		Remark:    req.Remark,
		CreatedBy: userID,
	}

	usage, err := h.settlementService.UseSettlement(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, gin.H{
		"usage_no": usage.UsageNo,
	}, "使用成功")
}

// ReturnSettlementRequest 归还沉淀款请求
type ReturnSettlementRequest struct {
	Amount int64  `json:"amount" binding:"required,min=1"` // 分
	Remark string `json:"remark"`
}

// ReturnSettlement 归还沉淀款
// @Summary 归还沉淀款
// @Description 归还沉淀钱包中的使用额度
// @Tags 沉淀钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ReturnSettlementRequest true "归还请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/return [post]
func (h *SettlementWalletHandler) ReturnSettlement(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	userID := middleware.GetCurrentUserID(c)

	var req ReturnSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.ReturnSettlementRequest{
		AgentID:   agentID,
		Amount:    req.Amount,
		Remark:    req.Remark,
		CreatedBy: userID,
	}

	usage, err := h.settlementService.ReturnSettlement(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, gin.H{
		"usage_no": usage.UsageNo,
	}, "归还成功")
}

// GetUsageList 获取使用记录列表
// @Summary 获取使用记录列表
// @Description 获取沉淀款使用/归还记录
// @Tags 沉淀钱包
// @Produce json
// @Security ApiKeyAuth
// @Param usage_type query int false "类型: 1=使用 2=归还"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settlement-wallet/usages [get]
func (h *SettlementWalletHandler) GetUsageList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var usageType *int16
	if typeStr := c.Query("usage_type"); typeStr != "" {
		t, _ := strconv.Atoi(typeStr)
		t16 := int16(t)
		usageType = &t16
	}

	list, total, err := h.settlementService.GetUsageList(agentID, usageType, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// RegisterSettlementWalletRoutes 注册沉淀钱包路由
func RegisterSettlementWalletRoutes(r *gin.RouterGroup, h *SettlementWalletHandler, authService *service.AuthService) {
	settlement := r.Group("/settlement-wallet")
	settlement.Use(middleware.AuthMiddleware(authService))
	{
		settlement.POST("/enable", h.EnableSettlementWallet)
		settlement.POST("/disable/:agent_id", h.DisableSettlementWallet)
		settlement.PUT("/ratio/:agent_id", h.UpdateSettlementRatio)

		settlement.GET("/summary", h.GetSummary)
		settlement.GET("/subordinates", h.GetSubordinateBalances)

		settlement.POST("/use", h.UseSettlement)
		settlement.POST("/return", h.ReturnSettlement)
		settlement.GET("/usages", h.GetUsageList)
	}
}
