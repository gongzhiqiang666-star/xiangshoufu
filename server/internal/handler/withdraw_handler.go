package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"
)

// WithdrawHandler 提现处理器
type WithdrawHandler struct {
	withdrawService *service.WithdrawService
}

// NewWithdrawHandler 创建提现处理器
func NewWithdrawHandler(withdrawService *service.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{
		withdrawService: withdrawService,
	}
}

// CreateWithdraw 创建提现申请
// @Summary 创建提现申请
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body service.CreateWithdrawRequest true "提现请求"
// @Success 200 {object} response.Response
// @Router /api/v1/withdraw [post]
func (h *WithdrawHandler) CreateWithdraw(c *gin.Context) {
	agentID := c.GetInt64("agent_id")
	if agentID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.AgentID = agentID

	record, err := h.withdrawService.CreateWithdraw(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"withdraw_no": record.WithdrawNo,
		"amount":      record.Amount,
		"actual":      record.ActualAmount,
		"tax_fee":     record.TaxFee,
		"fixed_fee":   record.FixedFee,
	})
}

// GetWithdrawList 获取提现记录列表
// @Summary 获取提现记录列表
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param status query int false "状态筛选"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/withdraw/list [get]
func (h *WithdrawHandler) GetWithdrawList(c *gin.Context) {
	agentID := c.GetInt64("agent_id")
	if agentID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var status *int16
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			st := int16(s)
			status = &st
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.withdrawService.GetWithdrawList(agentID, status, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// GetWithdrawDetail 获取提现详情
// @Summary 获取提现详情
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "提现记录ID"
// @Success 200 {object} response.Response
// @Router /api/v1/withdraw/{id} [get]
func (h *WithdrawHandler) GetWithdrawDetail(c *gin.Context) {
	agentID := c.GetInt64("agent_id")
	if agentID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	detail, err := h.withdrawService.GetWithdrawDetail(agentID, id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, detail)
}

// GetWithdrawStats 获取提现统计
// @Summary 获取提现统计
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Router /api/v1/withdraw/stats [get]
func (h *WithdrawHandler) GetWithdrawStats(c *gin.Context) {
	agentID := c.GetInt64("agent_id")
	if agentID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	stats, err := h.withdrawService.GetWithdrawStats(agentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// CancelWithdraw 取消提现
// @Summary 取消提现
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "提现记录ID"
// @Success 200 {object} response.Response
// @Router /api/v1/withdraw/{id}/cancel [post]
func (h *WithdrawHandler) CancelWithdraw(c *gin.Context) {
	agentID := c.GetInt64("agent_id")
	if agentID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.withdrawService.CancelWithdraw(agentID, id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ================== 管理员接口 ==================

// GetPendingList 获取待审核列表（管理员）
// @Summary 获取待审核提现列表
// @Tags 提现管理-管理员
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/withdraw/pending [get]
func (h *WithdrawHandler) GetPendingList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.withdrawService.GetPendingList(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// ApproveWithdrawRequest 审核通过请求
type ApproveWithdrawRequest struct {
	Remark string `json:"remark"`
}

// ApproveWithdraw 审核通过提现（管理员）
// @Summary 审核通过提现
// @Tags 提现管理-管理员
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "提现记录ID"
// @Param request body ApproveWithdrawRequest true "审核信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/withdraw/{id}/approve [post]
func (h *WithdrawHandler) ApproveWithdraw(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req ApproveWithdrawRequest
	c.ShouldBindJSON(&req)

	if err := h.withdrawService.ApproveWithdraw(id, userID, req.Remark); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// RejectWithdrawRequest 拒绝请求
type RejectWithdrawRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RejectWithdraw 拒绝提现（管理员）
// @Summary 拒绝提现
// @Tags 提现管理-管理员
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "提现记录ID"
// @Param request body RejectWithdrawRequest true "拒绝原因"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/withdraw/{id}/reject [post]
func (h *WithdrawHandler) RejectWithdraw(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req RejectWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请填写拒绝原因")
		return
	}

	if err := h.withdrawService.RejectWithdraw(id, userID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ConfirmPaidRequest 确认打款请求
type ConfirmPaidRequest struct {
	PaidRef string `json:"paid_ref" binding:"required"`
	Remark  string `json:"remark"`
}

// ConfirmPaid 确认打款（管理员）
// @Summary 确认打款
// @Tags 提现管理-管理员
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "提现记录ID"
// @Param request body ConfirmPaidRequest true "打款信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/withdraw/{id}/paid [post]
func (h *WithdrawHandler) ConfirmPaid(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req ConfirmPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请填写打款流水号")
		return
	}

	if err := h.withdrawService.ConfirmPaid(id, userID, req.PaidRef, req.Remark); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
