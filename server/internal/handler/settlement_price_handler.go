package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"
)

// SettlementPriceHandler 结算价处理器
type SettlementPriceHandler struct {
	service       *service.SettlementPriceService
	changeLogSvc  *service.PriceChangeLogService
}

// NewSettlementPriceHandler 创建结算价处理器
func NewSettlementPriceHandler(
	service *service.SettlementPriceService,
	changeLogSvc *service.PriceChangeLogService,
) *SettlementPriceHandler {
	return &SettlementPriceHandler{
		service:      service,
		changeLogSvc: changeLogSvc,
	}
}

// List 获取结算价列表
// @Summary 获取结算价列表
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param agent_id query int64 false "代理商ID"
// @Param channel_id query int64 false "通道ID"
// @Param status query int16 false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.SettlementPriceListResponse
// @Router /api/v1/settlement-prices [get]
func (h *SettlementPriceHandler) List(c *gin.Context) {
	var req models.SettlementPriceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetByID 获取结算价详情
// @Summary 获取结算价详情
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices/{id} [get]
func (h *SettlementPriceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	price, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "结算价不存在")
		return
	}

	response.Success(c, price)
}

// Create 创建结算价
// @Summary 创建结算价
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param body body models.CreateSettlementPriceRequest true "创建结算价请求"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices [post]
func (h *SettlementPriceHandler) Create(c *gin.Context) {
	var req models.CreateSettlementPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.AgentID <= 0 {
		response.BadRequest(c, "请选择代理商")
		return
	}
	if req.ChannelID <= 0 {
		response.BadRequest(c, "请选择通道")
		return
	}
	if req.TemplateID == nil || *req.TemplateID <= 0 {
		response.BadRequest(c, "请选择政策模板")
		return
	}

	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)

	price, err := h.service.CreateFromTemplate(
		req.AgentID,
		req.ChannelID,
		req.TemplateID,
		req.BrandCode,
		nil,
		operatorID,
		operatorName,
		source,
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, price)
}

// UpdateRate 更新费率
// @Summary 更新费率
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Param body body models.UpdateRateRequest true "更新费率请求"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices/{id}/rate [put]
func (h *SettlementPriceHandler) UpdateRate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req models.UpdateRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)
	ipAddress := c.ClientIP()

	price, err := h.service.UpdateRate(id, &req, operatorID, operatorName, source, ipAddress)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, price)
}

// UpdateDeposit 更新押金返现
// @Summary 更新押金返现
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Param body body models.UpdateDepositCashbackRequest true "更新押金返现请求"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices/{id}/deposit [put]
func (h *SettlementPriceHandler) UpdateDeposit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req models.UpdateDepositCashbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)
	ipAddress := c.ClientIP()

	price, err := h.service.UpdateDepositCashback(id, &req, operatorID, operatorName, source, ipAddress)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, price)
}

// UpdateSim 更新流量费返现
// @Summary 更新流量费返现
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Param body body models.UpdateSimCashbackRequest true "更新流量费返现请求"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices/{id}/sim [put]
func (h *SettlementPriceHandler) UpdateSim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req models.UpdateSimCashbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)
	ipAddress := c.ClientIP()

	price, err := h.service.UpdateSimCashback(id, &req, operatorID, operatorName, source, ipAddress)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, price)
}

// GetChangeLogs 获取结算价调价记录
// @Summary 获取结算价调价记录
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} models.PriceChangeLogListResponse
// @Router /api/v1/settlement-prices/{id}/change-logs [get]
func (h *SettlementPriceHandler) GetChangeLogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := h.changeLogSvc.ListBySettlementPrice(id, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// 辅助函数：获取操作者ID
func getOperatorID(c *gin.Context) int64 {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int64); ok {
			return id
		}
	}
	return 0
}

// 辅助函数：获取操作者名称
func getOperatorName(c *gin.Context) string {
	if userName, exists := c.Get("user_name"); exists {
		if name, ok := userName.(string); ok {
			return name
		}
	}
	return ""
}

// 辅助函数：获取操作来源
func getSource(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	if userAgent != "" {
		// 简单判断是否为移动端
		if len(userAgent) > 0 && (userAgent[0] == 'D' || userAgent[0] == 'o') {
			return "APP"
		}
	}
	return "PC"
}

// UpdateHighRate 更新高调费率配置
// @Summary 更新高调费率配置
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Param body body service.UpdateHighRateRequest true "更新高调费率请求"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices/{id}/high-rate [put]
func (h *SettlementPriceHandler) UpdateHighRate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req service.UpdateHighRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)
	ipAddress := c.ClientIP()

	price, err := h.service.UpdateHighRate(id, &req, operatorID, operatorName, source, ipAddress)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, price)
}

// UpdateD0Extra 更新P+0加价配置
// @Summary 更新P+0加价配置
// @Tags 结算价管理
// @Accept json
// @Produce json
// @Param id path int64 true "结算价ID"
// @Param body body service.UpdateD0ExtraRequest true "更新P+0加价请求"
// @Success 200 {object} models.SettlementPrice
// @Router /api/v1/settlement-prices/{id}/d0-extra [put]
func (h *SettlementPriceHandler) UpdateD0Extra(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req service.UpdateD0ExtraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := getOperatorID(c)
	operatorName := getOperatorName(c)
	source := getSource(c)
	ipAddress := c.ClientIP()

	price, err := h.service.UpdateD0Extra(id, &req, operatorID, operatorName, source, ipAddress)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, price)
}
