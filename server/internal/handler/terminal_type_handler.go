package handler

import (
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// TerminalTypeHandler 终端类型处理器
type TerminalTypeHandler struct {
	svc *service.TerminalTypeService
}

// NewTerminalTypeHandler 创建处理器
func NewTerminalTypeHandler(svc *service.TerminalTypeService) *TerminalTypeHandler {
	return &TerminalTypeHandler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *TerminalTypeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/terminal-types")
	{
		group.GET("", h.List)
		group.POST("", h.Create)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.PATCH("/:id/status", h.UpdateStatus)
		group.DELETE("/:id", h.Delete)
		group.GET("/by-channel/:channel_id", h.ListByChannel)
	}
}

// List 获取终端类型列表
// @Summary 获取终端类型列表
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param channel_id query int false "通道ID"
// @Param status query int false "状态：1启用 0禁用"
// @Param keyword query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} Response
// @Router /api/admin/terminal-types [get]
func (h *TerminalTypeHandler) List(c *gin.Context) {
	var req service.TerminalTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	list, total, err := h.svc.List(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "获取列表失败: "+err.Error())
		return
	}

	response.SuccessPage(c, h.svc.ToResponseList(list), total, req.Page, req.PageSize)
}

// Create 创建终端类型
// @Summary 创建终端类型
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param body body service.CreateTerminalTypeRequest true "创建参数"
// @Success 200 {object} Response
// @Router /api/admin/terminal-types [post]
func (h *TerminalTypeHandler) Create(c *gin.Context) {
	var req service.CreateTerminalTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	terminalType, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, h.svc.ToResponse(terminalType), "创建成功")
}

// GetByID 获取终端类型详情
// @Summary 获取终端类型详情
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param id path int true "终端类型ID"
// @Success 200 {object} Response
// @Router /api/admin/terminal-types/{id} [get]
func (h *TerminalTypeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	terminalType, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, h.svc.ToResponse(terminalType))
}

// Update 更新终端类型
// @Summary 更新终端类型
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param id path int true "终端类型ID"
// @Param body body service.UpdateTerminalTypeRequest true "更新参数"
// @Success 200 {object} Response
// @Router /api/admin/terminal-types/{id} [put]
func (h *TerminalTypeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req service.UpdateTerminalTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	terminalType, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, h.svc.ToResponse(terminalType), "更新成功")
}

// UpdateStatus 更新状态
// @Summary 更新终端类型状态
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param id path int true "终端类型ID"
// @Param body body UpdateStatusRequest true "状态参数"
// @Success 200 {object} Response
// @Router /api/admin/terminal-types/{id}/status [patch]
func (h *TerminalTypeHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Status int16 `json:"status" binding:"oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	statusText := "禁用"
	if req.Status == models.TerminalTypeStatusEnabled {
		statusText = "启用"
	}

	response.SuccessMessage(c, statusText+"成功")
}

// Delete 删除终端类型
// @Summary 删除终端类型
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param id path int true "终端类型ID"
// @Success 200 {object} Response
// @Router /api/admin/terminal-types/{id} [delete]
func (h *TerminalTypeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

// ListByChannel 根据通道ID获取终端类型列表（用于下拉选择）
// @Summary 根据通道获取终端类型列表
// @Tags 终端类型
// @Accept json
// @Produce json
// @Param channel_id path int true "通道ID"
// @Success 200 {object} Response
// @Router /api/admin/terminal-types/by-channel/{channel_id} [get]
func (h *TerminalTypeHandler) ListByChannel(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("channel_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	list, err := h.svc.ListByChannelID(c.Request.Context(), channelID)
	if err != nil {
		response.InternalError(c, "获取列表失败: "+err.Error())
		return
	}

	response.Success(c, h.svc.ToResponseList(list))
}
