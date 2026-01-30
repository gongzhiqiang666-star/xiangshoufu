package handler

import (
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// TerminalDistributeHandler 终端下发Handler
type TerminalDistributeHandler struct {
	distributeService *service.TerminalDistributeService
}

// NewTerminalDistributeHandler 创建终端下发Handler
func NewTerminalDistributeHandler(distributeService *service.TerminalDistributeService) *TerminalDistributeHandler {
	return &TerminalDistributeHandler{
		distributeService: distributeService,
	}
}

// DistributeTerminalRequest 终端下发请求
type DistributeTerminalRequest struct {
	ToAgentID        int64  `json:"to_agent_id" binding:"required"`       // 接收方代理商ID
	TerminalSN       string `json:"terminal_sn" binding:"required"`       // 终端SN
	ChannelID        int64  `json:"channel_id" binding:"required"`        // 通道ID
	GoodsPrice       int64  `json:"goods_price" binding:"required"`       // 货款金额（分）
	DeductionType    int16  `json:"deduction_type" binding:"required"`    // 1:一次性付款 2:分期代扣
	DeductionPeriods int    `json:"deduction_periods"`                    // 分期期数
	Remark           string `json:"remark"`                               // 备注
}

// DistributeTerminal 终端下发
// @Summary 终端下发
// @Description 将终端下发给下级代理商，APP端不能跨级，PC端可以跨级
// @Tags 终端管理
// @Accept json
// @Produce json
// @Param request body DistributeTerminalRequest true "终端下发请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminal/distribute [post]
func (h *TerminalDistributeHandler) DistributeTerminal(c *gin.Context) {
	var req DistributeTerminalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户信息
	fromAgentID := getCurrentAgentID(c)
	createdBy := getCurrentUserID(c)
	source := getRequestSource(c)

	serviceReq := &service.DistributeTerminalRequest{
		FromAgentID:      fromAgentID,
		ToAgentID:        req.ToAgentID,
		TerminalSN:       req.TerminalSN,
		ChannelID:        req.ChannelID,
		GoodsPrice:       req.GoodsPrice,
		DeductionType:    req.DeductionType,
		DeductionPeriods: req.DeductionPeriods,
		Source:           source,
		Remark:           req.Remark,
		CreatedBy:        createdBy,
	}

	distribute, err := h.distributeService.DistributeTerminal(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, distribute, "下发成功，等待接收方确认")
}

// ConfirmDistribute 确认接收终端
// @Summary 确认接收终端
// @Description 接收方确认接收终端下发
// @Tags 终端管理
// @Produce json
// @Param id path int true "下发记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminal/distribute/{id}/confirm [post]
func (h *TerminalDistributeHandler) ConfirmDistribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	confirmedBy := getCurrentUserID(c)

	if err := h.distributeService.ConfirmDistribute(id, confirmedBy); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "确认成功")
}

// RejectDistribute 拒绝终端下发
// @Summary 拒绝终端下发
// @Description 接收方拒绝终端下发
// @Tags 终端管理
// @Produce json
// @Param id path int true "下发记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminal/distribute/{id}/reject [post]
func (h *TerminalDistributeHandler) RejectDistribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	confirmedBy := getCurrentUserID(c)

	if err := h.distributeService.RejectDistribute(id, confirmedBy); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "已拒绝")
}

// CancelDistribute 取消终端下发
// @Summary 取消终端下发
// @Description 下发方取消终端下发
// @Tags 终端管理
// @Produce json
// @Param id path int true "下发记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminal/distribute/{id}/cancel [post]
func (h *TerminalDistributeHandler) CancelDistribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	cancelBy := getCurrentAgentID(c)

	if err := h.distributeService.CancelDistribute(id, cancelBy); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "已取消")
}

// GetDistributeListRequest 获取下发列表请求
type GetDistributeListRequest struct {
	Direction string  `form:"direction" binding:"required,oneof=from to"` // from:我下发的 to:下发给我的
	Status    []int16 `form:"status"`                                     // 状态筛选
	Page      int     `form:"page,default=1"`
	PageSize  int     `form:"page_size,default=20"`
}

// GetDistributeList 获取下发列表
// @Summary 获取下发列表
// @Description 获取终端下发列表，支持按方向筛选
// @Tags 终端管理
// @Produce json
// @Param direction query string true "方向: from(我下发的) to(下发给我的)"
// @Param status query []int false "状态筛选"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminal/distribute [get]
func (h *TerminalDistributeHandler) GetDistributeList(c *gin.Context) {
	var req GetDistributeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	agentID := getCurrentAgentID(c)
	offset := (req.Page - 1) * req.PageSize

	list, total, err := h.distributeService.GetDistributeList(agentID, req.Direction, req.Status, req.PageSize, offset)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, req.Page, req.PageSize)
}

// RegisterTerminalDistributeRoutes 注册终端下发路由
func RegisterTerminalDistributeRoutes(r *gin.RouterGroup, h *TerminalDistributeHandler) {
	terminal := r.Group("/terminal")
	{
		// 终端下发
		terminal.POST("/distribute", h.DistributeTerminal)
		terminal.GET("/distribute", h.GetDistributeList)
		terminal.POST("/distribute/:id/confirm", h.ConfirmDistribute)
		terminal.POST("/distribute/:id/reject", h.RejectDistribute)
		terminal.POST("/distribute/:id/cancel", h.CancelDistribute)
	}
}

// getRequestSource 获取请求来源
func getRequestSource(c *gin.Context) int16 {
	// 根据User-Agent或特定header判断来源
	userAgent := c.GetHeader("User-Agent")
	xSource := c.GetHeader("X-Source")

	// 如果明确指定了来源
	if xSource == "app" {
		return models.TerminalDistributeSourceApp
	}
	if xSource == "pc" {
		return models.TerminalDistributeSourcePC
	}

	// 根据User-Agent判断
	if contains(userAgent, "Mobile") || contains(userAgent, "Android") || contains(userAgent, "iOS") {
		return models.TerminalDistributeSourceApp
	}

	return models.TerminalDistributeSourcePC
}

// contains 判断字符串是否包含
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TerminalDistributeStatusDesc 终端下发状态描述
func TerminalDistributeStatusDesc(status int16) string {
	switch status {
	case models.TerminalDistributeStatusPending:
		return "待确认"
	case models.TerminalDistributeStatusConfirmed:
		return "已确认"
	case models.TerminalDistributeStatusRejected:
		return "已拒绝"
	case models.TerminalDistributeStatusCancelled:
		return "已取消"
	default:
		return "未知状态"
	}
}
