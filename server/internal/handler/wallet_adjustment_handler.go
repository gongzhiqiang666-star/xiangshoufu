package handler

import (
	"strconv"
	"time"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// WalletAdjustmentHandler 钱包调账处理器
type WalletAdjustmentHandler struct {
	adjustmentService *service.WalletAdjustmentService
}

// NewWalletAdjustmentHandler 创建钱包调账处理器
func NewWalletAdjustmentHandler(adjustmentService *service.WalletAdjustmentService) *WalletAdjustmentHandler {
	return &WalletAdjustmentHandler{
		adjustmentService: adjustmentService,
	}
}

// CreateAdjustmentRequest 创建调账请求
type CreateAdjustmentRequest struct {
	AgentID    int64  `json:"agent_id" binding:"required"`
	WalletType int16  `json:"wallet_type" binding:"required"` // 1分润 2服务费 3奖励 4充值 5沉淀
	ChannelID  int64  `json:"channel_id"`                     // 0表示不区分通道
	Amount     int64  `json:"amount" binding:"required"`      // 正数充入，负数扣减（分）
	Reason     string `json:"reason" binding:"required"`
}

// CreateAdjustment 创建调账
// @Summary 创建调账
// @Description 管理员手动调账（充入或扣减指定钱包余额）
// @Tags 钱包调账
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateAdjustmentRequest true "调账请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallet-adjustments [post]
func (h *WalletAdjustmentHandler) CreateAdjustment(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	userName := middleware.GetCurrentUsername(c)

	var req CreateAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Amount == 0 {
		response.BadRequest(c, "调账金额不能为0")
		return
	}

	serviceReq := &service.CreateAdjustmentRequest{
		AgentID:      req.AgentID,
		WalletType:   req.WalletType,
		ChannelID:    req.ChannelID,
		Amount:       req.Amount,
		Reason:       req.Reason,
		OperatorID:   userID,
		OperatorName: userName,
	}

	adjustment, err := h.adjustmentService.CreateAdjustment(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, adjustment, "调账成功")
}

// GetAdjustmentList 获取调账列表
// @Summary 获取调账列表
// @Description 查询调账记录列表
// @Tags 钱包调账
// @Produce json
// @Security ApiKeyAuth
// @Param agent_id query int false "代理商ID"
// @Param wallet_type query int false "钱包类型"
// @Param channel_id query int false "通道ID"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallet-adjustments [get]
func (h *WalletAdjustmentHandler) GetAdjustmentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := &service.AdjustmentListParams{
		Page:     page,
		PageSize: pageSize,
	}

	if agentIDStr := c.Query("agent_id"); agentIDStr != "" {
		agentID, _ := strconv.ParseInt(agentIDStr, 10, 64)
		params.AgentID = agentID
	}

	if walletTypeStr := c.Query("wallet_type"); walletTypeStr != "" {
		wt, _ := strconv.Atoi(walletTypeStr)
		wt16 := int16(wt)
		params.WalletType = &wt16
	}

	if channelIDStr := c.Query("channel_id"); channelIDStr != "" {
		chID, _ := strconv.ParseInt(channelIDStr, 10, 64)
		params.ChannelID = &chID
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			params.StartTime = &t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			t = t.AddDate(0, 0, 1)
			params.EndTime = &t
		}
	}

	list, total, err := h.adjustmentService.GetAdjustmentList(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// GetAdjustmentDetail 获取调账详情
// @Summary 获取调账详情
// @Description 根据ID获取调账记录详情
// @Tags 钱包调账
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "调账记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wallet-adjustments/{id} [get]
func (h *WalletAdjustmentHandler) GetAdjustmentDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	adjustment, err := h.adjustmentService.GetAdjustmentDetail(id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, adjustment)
}

// RegisterWalletAdjustmentRoutes 注册钱包调账路由
func RegisterWalletAdjustmentRoutes(r *gin.RouterGroup, h *WalletAdjustmentHandler, authService *service.AuthService) {
	adjustments := r.Group("/wallet-adjustments")
	adjustments.Use(middleware.AuthMiddleware(authService))
	{
		adjustments.POST("", h.CreateAdjustment)
		adjustments.GET("", h.GetAdjustmentList)
		adjustments.GET("/:id", h.GetAdjustmentDetail)
	}
}
