package handler

import (
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// ChargingWalletHandler 充值钱包处理器
type ChargingWalletHandler struct {
	chargingService *service.ChargingWalletService
}

// NewChargingWalletHandler 创建充值钱包处理器
func NewChargingWalletHandler(chargingService *service.ChargingWalletService) *ChargingWalletHandler {
	return &ChargingWalletHandler{
		chargingService: chargingService,
	}
}

// GetWalletConfig 获取钱包配置
// @Summary 获取钱包配置
// @Description 获取当前代理商的充值钱包和沉淀钱包配置
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/config [get]
func (h *ChargingWalletHandler) GetWalletConfig(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	config, err := h.chargingService.GetWalletConfig(agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, config)
}

// GetWalletConfigByAgent 获取指定代理商的钱包配置(管理员)
// @Summary 获取指定代理商的钱包配置
// @Description 管理员获取指定代理商的充值钱包和沉淀钱包配置
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id path int true "代理商ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/config/{agent_id} [get]
func (h *ChargingWalletHandler) GetWalletConfigByAgent(c *gin.Context) {
	agentIDStr := c.Param("agent_id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的代理商ID")
		return
	}

	config, err := h.chargingService.GetWalletConfig(agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, config)
}

// EnableChargingWalletRequest 开通充值钱包请求
type EnableChargingWalletRequest struct {
	AgentID int64 `json:"agent_id" binding:"required"`
	Limit   int64 `json:"limit"` // 限额(分)
}

// EnableChargingWallet 开通充值钱包
// @Summary 开通充值钱包
// @Description PC端管理员为代理商开通充值钱包
// @Tags 充值钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body EnableChargingWalletRequest true "开通请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/enable [post]
func (h *ChargingWalletHandler) EnableChargingWallet(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req EnableChargingWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.EnableChargingWalletRequest{
		AgentID:   req.AgentID,
		Limit:     req.Limit,
		EnabledBy: userID,
	}

	if err := h.chargingService.EnableChargingWallet(serviceReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "开通成功")
}

// DisableChargingWallet 关闭充值钱包
// @Summary 关闭充值钱包
// @Description PC端管理员关闭代理商的充值钱包
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id path int true "代理商ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/disable/{agent_id} [post]
func (h *ChargingWalletHandler) DisableChargingWallet(c *gin.Context) {
	agentIDStr := c.Param("agent_id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的代理商ID")
		return
	}

	if err := h.chargingService.DisableChargingWallet(agentID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "关闭成功")
}

// GetSummary 获取充值钱包汇总
// @Summary 获取充值钱包汇总
// @Description 获取当前代理商充值钱包的余额和发放统计
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/summary [get]
func (h *ChargingWalletHandler) GetSummary(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	summary, err := h.chargingService.GetChargingWalletSummary(agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, summary)
}

// CreateDepositRequest 申请充值请求
type CreateDepositRequest struct {
	Amount        int64  `json:"amount" binding:"required,min=100"` // 最少1元(100分)
	PaymentMethod int16  `json:"payment_method" binding:"required"` // 1=银行转账 2=微信 3=支付宝
	PaymentRef    string `json:"payment_ref"`
	Remark        string `json:"remark"`
}

// CreateDeposit 申请充值
// @Summary 申请充值
// @Description 代理商申请充值钱包充值
// @Tags 充值钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateDepositRequest true "充值请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/deposits [post]
func (h *ChargingWalletHandler) CreateDeposit(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	userID := middleware.GetCurrentUserID(c)

	var req CreateDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.CreateDepositRequest{
		AgentID:       agentID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentRef:    req.PaymentRef,
		Remark:        req.Remark,
		CreatedBy:     userID,
	}

	deposit, err := h.chargingService.CreateDeposit(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, gin.H{
		"deposit_no": deposit.DepositNo,
	}, "充值申请已提交")
}

// GetDepositList 获取充值记录
// @Summary 获取充值记录
// @Description 获取当前代理商的充值记录列表
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Param status query int false "状态筛选"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/deposits [get]
func (h *ChargingWalletHandler) GetDepositList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var status *int16
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		s16 := int16(s)
		status = &s16
	}

	list, total, err := h.chargingService.GetDepositList(agentID, status, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// GetPendingDeposits 获取待审核充值记录(管理员)
// @Summary 获取待审核充值记录
// @Description 管理员获取所有待审核的充值申请
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/deposits/pending [get]
func (h *ChargingWalletHandler) GetPendingDeposits(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pending := int16(0)
	list, total, err := h.chargingService.GetDepositList(0, &pending, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// ConfirmDeposit 确认充值
// @Summary 确认充值
// @Description 管理员确认充值申请
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "充值记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/deposits/{id}/confirm [post]
func (h *ChargingWalletHandler) ConfirmDeposit(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.chargingService.ConfirmDeposit(id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "确认成功")
}

// RejectDepositRequest 拒绝充值请求
type RejectDepositRequest struct {
	Reason string `json:"reason"`
}

// RejectDeposit 拒绝充值
// @Summary 拒绝充值
// @Description 管理员拒绝充值申请
// @Tags 充值钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "充值记录ID"
// @Param request body RejectDepositRequest true "拒绝原因"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/deposits/{id}/reject [post]
func (h *ChargingWalletHandler) RejectDeposit(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req RejectDepositRequest
	c.ShouldBindJSON(&req)

	if err := h.chargingService.RejectDeposit(id, userID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "拒绝成功")
}

// IssueRewardRequest 发放奖励请求
type IssueRewardRequest struct {
	ToAgentID int64  `json:"to_agent_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,min=1"` // 分
	Remark    string `json:"remark"`
}

// IssueReward 发放奖励
// @Summary 发放奖励
// @Description 代理商从充值钱包发放奖励给下级
// @Tags 充值钱包
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body IssueRewardRequest true "发放请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/rewards [post]
func (h *ChargingWalletHandler) IssueReward(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	userID := middleware.GetCurrentUserID(c)

	var req IssueRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.IssueRewardRequest{
		FromAgentID: agentID,
		ToAgentID:   req.ToAgentID,
		Amount:      req.Amount,
		Remark:      req.Remark,
		CreatedBy:   userID,
	}

	reward, err := h.chargingService.IssueReward(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, gin.H{
		"reward_no": reward.RewardNo,
	}, "发放成功")
}

// GetRewardList 获取奖励记录
// @Summary 获取奖励记录
// @Description 获取当前代理商的奖励发放/接收记录
// @Tags 充值钱包
// @Produce json
// @Security ApiKeyAuth
// @Param direction query string true "方向: from=我发放的, to=我收到的"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/charging-wallet/rewards [get]
func (h *ChargingWalletHandler) GetRewardList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	direction := c.DefaultQuery("direction", "from")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.chargingService.GetRewardList(agentID, direction, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// RegisterChargingWalletRoutes 注册充值钱包路由
func RegisterChargingWalletRoutes(r *gin.RouterGroup, h *ChargingWalletHandler, authService *service.AuthService) {
	charging := r.Group("/charging-wallet")
	charging.Use(middleware.AuthMiddleware(authService))
	{
		charging.GET("/config", h.GetWalletConfig)
		charging.GET("/config/:agent_id", h.GetWalletConfigByAgent)
		charging.POST("/enable", h.EnableChargingWallet)
		charging.POST("/disable/:agent_id", h.DisableChargingWallet)

		charging.GET("/summary", h.GetSummary)

		charging.POST("/deposits", h.CreateDeposit)
		charging.GET("/deposits", h.GetDepositList)
		charging.GET("/deposits/pending", h.GetPendingDeposits)
		charging.POST("/deposits/:id/confirm", h.ConfirmDeposit)
		charging.POST("/deposits/:id/reject", h.RejectDeposit)

		charging.POST("/rewards", h.IssueReward)
		charging.GET("/rewards", h.GetRewardList)
	}
}
