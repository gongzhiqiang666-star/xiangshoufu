package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
)

// PriceChangeLogHandler 调价记录处理器
type PriceChangeLogHandler struct {
	service *service.PriceChangeLogService
}

// NewPriceChangeLogHandler 创建调价记录处理器
func NewPriceChangeLogHandler(service *service.PriceChangeLogService) *PriceChangeLogHandler {
	return &PriceChangeLogHandler{service: service}
}

// List 获取调价记录列表
// @Summary 获取调价记录列表
// @Tags 调价记录
// @Accept json
// @Produce json
// @Param agent_id query int64 false "代理商ID"
// @Param channel_id query int64 false "通道ID"
// @Param change_type query int16 false "变更类型"
// @Param config_type query int16 false "配置类型"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.PriceChangeLogListResponse
// @Router /api/v1/price-change-logs [get]
func (h *PriceChangeLogHandler) List(c *gin.Context) {
	var req models.PriceChangeLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	resp, err := h.service.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetByID 获取调价记录详情
// @Summary 获取调价记录详情
// @Tags 调价记录
// @Accept json
// @Produce json
// @Param id path int64 true "调价记录ID"
// @Success 200 {object} models.PriceChangeLog
// @Router /api/v1/price-change-logs/{id} [get]
func (h *PriceChangeLogHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	log, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "调价记录不存在"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// ListByAgent 按代理商获取调价记录
// @Summary 按代理商获取调价记录
// @Tags 调价记录
// @Accept json
// @Produce json
// @Param agent_id path int64 true "代理商ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.PriceChangeLogListResponse
// @Router /api/v1/agents/{agent_id}/price-change-logs [get]
func (h *PriceChangeLogHandler) ListByAgent(c *gin.Context) {
	agentIDStr := c.Param("agent_id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的代理商ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := h.service.ListByAgent(agentID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
