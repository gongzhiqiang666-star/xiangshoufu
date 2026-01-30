package handler

import (
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

// TaxChannelHandler 税筹通道处理器
type TaxChannelHandler struct {
	taxChannelService *service.TaxChannelService
}

// NewTaxChannelHandler 创建税筹通道处理器
func NewTaxChannelHandler(taxChannelService *service.TaxChannelService) *TaxChannelHandler {
	return &TaxChannelHandler{
		taxChannelService: taxChannelService,
	}
}

// ========== 税筹通道管理 ==========

// CreateTaxChannelRequest 创建税筹通道请求
type CreateTaxChannelRequest struct {
	ChannelCode string  `json:"channel_code" binding:"required"`
	ChannelName string  `json:"channel_name" binding:"required"`
	FeeType     int16   `json:"fee_type" binding:"required,oneof=1 2"` // 1=付款扣 2=出款扣
	TaxRate     float64 `json:"tax_rate" binding:"required,min=0,max=1"`
	FixedFee    int64   `json:"fixed_fee"`
	ApiURL      string  `json:"api_url"`
	ApiKey      string  `json:"api_key"`
	ApiSecret   string  `json:"api_secret"`
	Remark      string  `json:"remark"`
}

// CreateTaxChannel 创建税筹通道
// @Summary 创建税筹通道
// @Description 创建新的税筹通道配置
// @Tags 税筹通道
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateTaxChannelRequest true "创建请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels [post]
func (h *TaxChannelHandler) CreateTaxChannel(c *gin.Context) {
	var req CreateTaxChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.CreateTaxChannelRequest{
		ChannelCode: req.ChannelCode,
		ChannelName: req.ChannelName,
		FeeType:     req.FeeType,
		TaxRate:     req.TaxRate,
		FixedFee:    req.FixedFee,
		ApiURL:      req.ApiURL,
		ApiKey:      req.ApiKey,
		ApiSecret:   req.ApiSecret,
		Remark:      req.Remark,
	}

	taxChannel, err := h.taxChannelService.CreateTaxChannel(serviceReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, gin.H{"id": taxChannel.ID}, "创建成功")
}

// UpdateTaxChannelRequest 更新税筹通道请求
type UpdateTaxChannelRequest struct {
	ChannelName string  `json:"channel_name"`
	FeeType     int16   `json:"fee_type"`
	TaxRate     float64 `json:"tax_rate"`
	FixedFee    int64   `json:"fixed_fee"`
	ApiURL      string  `json:"api_url"`
	ApiKey      string  `json:"api_key"`
	ApiSecret   string  `json:"api_secret"`
	Status      *int16  `json:"status"`
	Remark      string  `json:"remark"`
}

// UpdateTaxChannel 更新税筹通道
// @Summary 更新税筹通道
// @Description 更新税筹通道配置
// @Tags 税筹通道
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "税筹通道ID"
// @Param request body UpdateTaxChannelRequest true "更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/{id} [put]
func (h *TaxChannelHandler) UpdateTaxChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateTaxChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.UpdateTaxChannelRequest{
		ID:          id,
		ChannelName: req.ChannelName,
		FeeType:     req.FeeType,
		TaxRate:     req.TaxRate,
		FixedFee:    req.FixedFee,
		ApiURL:      req.ApiURL,
		ApiKey:      req.ApiKey,
		ApiSecret:   req.ApiSecret,
		Status:      req.Status,
		Remark:      req.Remark,
	}

	if _, err := h.taxChannelService.UpdateTaxChannel(serviceReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

// GetTaxChannel 获取税筹通道详情
// @Summary 获取税筹通道详情
// @Description 获取税筹通道的详细信息
// @Tags 税筹通道
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "税筹通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/{id} [get]
func (h *TaxChannelHandler) GetTaxChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	taxChannel, err := h.taxChannelService.GetTaxChannel(id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, taxChannel)
}

// GetTaxChannelList 获取税筹通道列表
// @Summary 获取税筹通道列表
// @Description 获取所有税筹通道配置
// @Tags 税筹通道
// @Produce json
// @Security ApiKeyAuth
// @Param status query int false "状态: 0=禁用 1=启用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels [get]
func (h *TaxChannelHandler) GetTaxChannelList(c *gin.Context) {
	var status *int16
	if statusStr := c.Query("status"); statusStr != "" {
		s, err := strconv.Atoi(statusStr)
		if err == nil {
			s16 := int16(s)
			status = &s16
		}
	}

	list, err := h.taxChannelService.GetTaxChannelList(status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": len(list),
	})
}

// DeleteTaxChannel 删除税筹通道
// @Summary 删除税筹通道
// @Description 删除税筹通道配置
// @Tags 税筹通道
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "税筹通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/{id} [delete]
func (h *TaxChannelHandler) DeleteTaxChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.taxChannelService.DeleteTaxChannel(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

// ========== 通道-税筹通道映射 ==========

// SetChannelTaxMappingRequest 设置通道税筹映射请求
type SetChannelTaxMappingRequest struct {
	ChannelID    int64 `json:"channel_id" binding:"required"`
	WalletType   int16 `json:"wallet_type" binding:"required,min=1,max=5"`
	TaxChannelID int64 `json:"tax_channel_id" binding:"required"`
}

// SetChannelTaxMapping 设置通道税筹映射
// @Summary 设置通道税筹映射
// @Description 设置支付通道与税筹通道的关联
// @Tags 税筹通道
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body SetChannelTaxMappingRequest true "映射请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/mappings [post]
func (h *TaxChannelHandler) SetChannelTaxMapping(c *gin.Context) {
	var req SetChannelTaxMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	serviceReq := &service.SetChannelTaxMappingRequest{
		ChannelID:    req.ChannelID,
		WalletType:   req.WalletType,
		TaxChannelID: req.TaxChannelID,
	}

	if err := h.taxChannelService.SetChannelTaxMapping(serviceReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "设置成功")
}

// GetChannelTaxMappings 获取通道税筹映射
// @Summary 获取通道税筹映射
// @Description 获取支付通道的税筹通道关联配置
// @Tags 税筹通道
// @Produce json
// @Security ApiKeyAuth
// @Param channel_id path int true "支付通道ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/mappings/channel/{channel_id} [get]
func (h *TaxChannelHandler) GetChannelTaxMappings(c *gin.Context) {
	channelIDStr := c.Param("channel_id")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	mappings, err := h.taxChannelService.GetChannelTaxMappings(channelID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": mappings})
}

// DeleteChannelTaxMapping 删除通道税筹映射
// @Summary 删除通道税筹映射
// @Description 删除支付通道与税筹通道的关联
// @Tags 税筹通道
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "映射ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/mappings/{id} [delete]
func (h *TaxChannelHandler) DeleteChannelTaxMapping(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.taxChannelService.DeleteChannelTaxMapping(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

// ========== 税费计算 ==========

// CalculateTaxRequest 计算税费请求
type CalculateTaxRequest struct {
	ChannelID  int64 `json:"channel_id" binding:"required"`
	WalletType int16 `json:"wallet_type" binding:"required,min=1,max=5"`
	Amount     int64 `json:"amount" binding:"required,min=1"` // 金额(分)
}

// CalculateWithdrawalTax 计算提现税费
// @Summary 计算提现税费
// @Description 根据通道和钱包类型计算提现税费
// @Tags 税筹通道
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CalculateTaxRequest true "计算请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tax-channels/calculate [post]
func (h *TaxChannelHandler) CalculateWithdrawalTax(c *gin.Context) {
	var req CalculateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.taxChannelService.CalculateWithdrawalTax(req.ChannelID, req.WalletType, req.Amount)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// RegisterTaxChannelRoutes 注册税筹通道路由
func RegisterTaxChannelRoutes(r *gin.RouterGroup, h *TaxChannelHandler, authService *service.AuthService) {
	taxChannel := r.Group("/tax-channels")
	taxChannel.Use(middleware.AuthMiddleware(authService))
	{
		// 税筹通道CRUD
		taxChannel.POST("", h.CreateTaxChannel)
		taxChannel.GET("", h.GetTaxChannelList)
		taxChannel.GET("/:id", h.GetTaxChannel)
		taxChannel.PUT("/:id", h.UpdateTaxChannel)
		taxChannel.DELETE("/:id", h.DeleteTaxChannel)

		// 通道-税筹通道映射
		taxChannel.POST("/mappings", h.SetChannelTaxMapping)
		taxChannel.GET("/mappings/channel/:channel_id", h.GetChannelTaxMappings)
		taxChannel.DELETE("/mappings/:id", h.DeleteChannelTaxMapping)

		// 税费计算
		taxChannel.POST("/calculate", h.CalculateWithdrawalTax)
	}
}
