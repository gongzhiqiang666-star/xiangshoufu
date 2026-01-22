package handler

import (
	"net/http"
	"strconv"
	"time"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// MerchantHandler 商户处理器
type MerchantHandler struct {
	merchantRepo    *repository.GormMerchantRepository
	transactionRepo *repository.GormTransactionRepository
	merchantService *service.MerchantService
	agentRepo       *repository.GormAgentRepository
	channelRepo     *repository.GormChannelRepository
	terminalRepo    *repository.GormTerminalRepository
}

// NewMerchantHandler 创建商户处理器
func NewMerchantHandler(
	merchantRepo *repository.GormMerchantRepository,
	transactionRepo *repository.GormTransactionRepository,
	merchantService *service.MerchantService,
) *MerchantHandler {
	return &MerchantHandler{
		merchantRepo:    merchantRepo,
		transactionRepo: transactionRepo,
		merchantService: merchantService,
	}
}

// SetAgentRepo 设置代理商仓储（延迟注入，避免循环依赖）
func (h *MerchantHandler) SetAgentRepo(agentRepo *repository.GormAgentRepository) {
	h.agentRepo = agentRepo
}

// SetChannelRepo 设置通道仓储
func (h *MerchantHandler) SetChannelRepo(channelRepo *repository.GormChannelRepository) {
	h.channelRepo = channelRepo
}

// SetTerminalRepo 设置终端仓储
func (h *MerchantHandler) SetTerminalRepo(terminalRepo *repository.GormTerminalRepository) {
	h.terminalRepo = terminalRepo
}

// GetMerchantList 获取商户列表
// @Summary 获取商户列表
// @Description 获取当前代理商的商户列表
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态"
// @Param merchant_type query string false "商户类型"
// @Param is_direct query bool false "是否直营"
// @Param channel_id query int false "通道ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants [get]
func (h *MerchantHandler) GetMerchantList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	// 解析参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	params := repository.MerchantQueryParams{
		AgentID:      agentID,
		Keyword:      c.Query("keyword"),
		MerchantType: c.Query("merchant_type"),
		Limit:        pageSize,
		Offset:       offset,
	}

	// 解析status
	if s := c.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			t := int16(v)
			params.Status = &t
		}
	}

	// 解析is_direct
	if d := c.Query("is_direct"); d != "" {
		if v, err := strconv.ParseBool(d); err == nil {
			params.IsDirect = &v
		}
	}

	// 解析channel_id
	if cid := c.Query("channel_id"); cid != "" {
		if v, err := strconv.ParseInt(cid, 10, 64); err == nil {
			params.ChannelID = &v
		}
	}

	merchants, total, err := h.merchantRepo.FindByParams(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(merchants))
	for _, m := range merchants {
		list = append(list, gin.H{
			"id":              m.ID,
			"merchant_no":     m.MerchantNo,
			"merchant_name":   m.MerchantName,
			"agent_id":        m.AgentID,
			"terminal_sn":     m.TerminalSN,
			"status":          m.Status,
			"status_name":     getMerchantStatusName(m.Status),
			"approve_status":  m.ApproveStatus,
			"approve_name":    getMerchantApproveStatusName(m.ApproveStatus),
			"mcc":             m.MCC,
			"credit_rate":     m.CreditRate,
			"debit_rate":      m.DebitRate,
			"merchant_type":   m.MerchantType,
			"is_direct":       m.IsDirect,
			"owner_type":      getOwnerTypeName(m.IsDirect),
			"activated_at":    m.ActivatedAt,
			"registered_phone": maskPhone(m.RegisteredPhone),
			"created_at":      m.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMerchantDetail 获取商户详情
// @Summary 获取商户详情
// @Description 获取指定商户的详细信息
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id} [get]
func (h *MerchantHandler) GetMerchantDetail(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	merchant, err := h.merchantRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "商户不存在",
		})
		return
	}

	// 验证权限
	if merchant.AgentID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权访问该商户",
		})
		return
	}

	// 格式化时间字段
	var activatedAtStr, createdAtStr, updatedAtStr *string
	if merchant.ActivatedAt != nil {
		s := merchant.ActivatedAt.Format("2006-01-02 15:04:05")
		activatedAtStr = &s
	}
	createdAt := merchant.CreatedAt.Format("2006-01-02 15:04:05")
	createdAtStr = &createdAt
	updatedAt := merchant.UpdatedAt.Format("2006-01-02 15:04:05")
	updatedAtStr = &updatedAt

	// 获取代理商信息
	var agentName *string
	var agentLevel *int
	if h.agentRepo != nil {
		if agent, err := h.agentRepo.FindByID(merchant.AgentID); err == nil && agent != nil {
			agentName = &agent.AgentName
			agentLevel = &agent.Level
		}
	}

	// 获取通道信息
	var channelName *string
	if h.channelRepo != nil {
		if channel, err := h.channelRepo.FindByID(merchant.ChannelID); err == nil && channel != nil {
			channelName = &channel.ChannelName
		}
	}

	// 获取本月交易统计
	var monthAmount, monthCount *int64
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if stats, err := h.merchantRepo.GetMerchantTransStats(merchant.ID, &monthStart, &now); err == nil && stats != nil {
		monthAmount = &stats.TotalAmount
		monthCount = &stats.TotalCount
	}

	// 获取终端数量
	var terminalCount *int64
	if h.terminalRepo != nil {
		if count, err := h.terminalRepo.CountByMerchantNo(merchant.MerchantNo); err == nil {
			terminalCount = &count
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":              merchant.ID,
			"merchant_no":     merchant.MerchantNo,
			"merchant_name":   merchant.MerchantName,
			"agent_id":        merchant.AgentID,
			"agent_name":      agentName,
			"agent_level":     agentLevel,
			"channel_id":      merchant.ChannelID,
			"channel_name":    channelName,
			"terminal_sn":     merchant.TerminalSN,
			"status":          merchant.Status,
			"status_name":     getMerchantStatusName(merchant.Status),
			"approve_status":  merchant.ApproveStatus,
			"approve_name":    getMerchantApproveStatusName(merchant.ApproveStatus),
			"legal_name":      merchant.LegalName,
			"legal_id_card":   maskIDCard(merchant.LegalIDCard),
			"mcc":             merchant.MCC,
			"credit_rate":     merchant.CreditRate,
			"debit_rate":      merchant.DebitRate,
			"merchant_type":   merchant.MerchantType,
			"is_direct":       merchant.IsDirect,
			"activated_at":    activatedAtStr,
			"registered_phone": maskPhone(merchant.RegisteredPhone),
			"register_remark": merchant.RegisterRemark,
			"month_amount":    monthAmount,
			"month_count":     monthCount,
			"terminal_count":  terminalCount,
			"created_at":      createdAtStr,
			"updated_at":      updatedAtStr,
		},
	})
}

// GetMerchantStats 获取商户统计
// @Summary 获取商户统计
// @Description 获取当前代理商的商户统计数据
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/stats [get]
func (h *MerchantHandler) GetMerchantStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	stats, err := h.merchantRepo.GetAgentMerchantStats(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_count":    stats.TotalCount,
			"active_count":   stats.ActiveCount,
			"pending_count":  stats.PendingCount,
			"disabled_count": stats.DisabledCount,
		},
	})
}

// GetMerchantTransactions 获取商户交易
// @Summary 获取商户交易
// @Description 获取指定商户的交易记录
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id}/transactions [get]
func (h *MerchantHandler) GetMerchantTransactions(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	idStr := c.Param("id")
	merchantID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	// 验证权限
	merchant, err := h.merchantRepo.FindByID(merchantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "商户不存在",
		})
		return
	}
	if merchant.AgentID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权访问该商户",
		})
		return
	}

	// 解析参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		if t, err := time.Parse("2006-01-02", st); err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		if t, err := time.Parse("2006-01-02", et); err == nil {
			endOfDay := t.Add(24 * time.Hour)
			endTime = &endOfDay
		}
	}

	transactions, total, err := h.transactionRepo.FindByMerchantID(merchantID, startTime, endTime, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(transactions))
	for _, tx := range transactions {
		list = append(list, gin.H{
			"id":              tx.ID,
			"order_no":        tx.OrderNo,
			"trade_no":        tx.TradeNo,
			"trade_type":      tx.TradeType,
			"trade_type_name": getTradeTypeName(tx.TradeType),
			"pay_type":        tx.PayType,
			"pay_type_name":   getPayTypeName(tx.PayType),
			"amount":          tx.Amount,
			"amount_yuan":     float64(tx.Amount) / 100,
			"fee":             tx.Fee,
			"fee_yuan":        float64(tx.Fee) / 100,
			"trade_time":      tx.TradeTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// getMerchantStatusName 商户状态名称
func getMerchantStatusName(status int16) string {
	switch status {
	case 1:
		return "正常"
	case 2:
		return "禁用"
	default:
		return "未知"
	}
}

// getMerchantApproveStatusName 审核状态名称
func getMerchantApproveStatusName(status int16) string {
	switch status {
	case 1:
		return "待审核"
	case 2:
		return "已通过"
	case 3:
		return "已拒绝"
	default:
		return "未知"
	}
}

// maskIDCard 遮掩身份证号
func maskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return idCard
	}
	return idCard[:6] + "********" + idCard[len(idCard)-4:]
}

// maskPhone 遮掩手机号
func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// getOwnerTypeName 获取归属类型名称
func getOwnerTypeName(isDirect bool) string {
	if isDirect {
		return "direct"
	}
	return "team"
}

// getMerchantTypeName 商户类型名称
func getMerchantTypeName(merchantType string) string {
	switch merchantType {
	case "loyal":
		return "忠诚商户"
	case "quality":
		return "优质商户"
	case "potential":
		return "潜力商户"
	case "normal":
		return "一般商户"
	case "low_active":
		return "低活跃"
	case "inactive":
		return "30天无交易"
	default:
		return "未知"
	}
}

// ==================== 新增API端点 ====================

// CreateMerchant 创建商户
// @Summary 创建商户
// @Description 创建新商户
// @Tags 商户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateMerchantRequest true "商户信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants [post]
func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	var req service.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置代理商ID为当前用户
	req.AgentID = middleware.GetCurrentAgentID(c)

	merchant, err := h.merchantService.CreateMerchant(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":          merchant.ID,
			"merchant_no": merchant.MerchantNo,
		},
	})
}

// UpdateMerchant 更新商户
// @Summary 更新商户
// @Description 更新商户信息
// @Tags 商户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Param request body service.UpdateMerchantRequest true "更新信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id} [put]
func (h *MerchantHandler) UpdateMerchant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req service.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.merchantService.UpdateMerchant(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteMerchant 删除商户
// @Summary 删除商户
// @Description 删除指定商户
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id} [delete]
func (h *MerchantHandler) DeleteMerchant(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.merchantService.DeleteMerchant(id, agentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdateMerchantStatus 更新商户状态
// @Summary 更新商户状态
// @Description 启用/禁用商户
// @Tags 商户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Param request body service.UpdateStatusRequest true "状态信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id}/status [put]
func (h *MerchantHandler) UpdateMerchantStatus(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req service.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.merchantService.UpdateStatus(id, agentID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
	})
}

// UpdateMerchantRate 更新商户费率
// @Summary 更新商户费率
// @Description 修改商户贷记卡/借记卡费率
// @Tags 商户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Param request body service.UpdateRateRequest true "费率信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id}/rate [put]
func (h *MerchantHandler) UpdateMerchantRate(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req service.UpdateRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.merchantService.UpdateRate(id, agentID, req.CreditRate, req.DebitRate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "费率更新成功",
	})
}

// RegisterMerchant 商户登记
// @Summary 商户登记
// @Description 登记商户联系方式
// @Tags 商户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "商户ID"
// @Param request body service.RegisterMerchantRequest true "登记信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id}/register [post]
func (h *MerchantHandler) RegisterMerchant(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var req service.RegisterMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.merchantService.RegisterMerchant(id, agentID, req.Phone, req.Remark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登记成功",
	})
}

// GetExtendedStats 获取扩展统计
// @Summary 获取扩展统计
// @Description 获取商户扩展统计数据（包含直营/团队、商户类型分布）
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/stats/extended [get]
func (h *MerchantHandler) GetExtendedStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	stats, err := h.merchantService.GetExtendedStats(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// ExportMerchants 导出商户
// @Summary 导出商户
// @Description 导出商户列表为Excel
// @Tags 商户管理
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态"
// @Param merchant_type query string false "商户类型"
// @Param is_direct query bool false "是否直营"
// @Success 200 {file} binary
// @Router /api/v1/merchants/export [get]
func (h *MerchantHandler) ExportMerchants(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	params := repository.MerchantQueryParams{
		AgentID:      agentID,
		Keyword:      c.Query("keyword"),
		MerchantType: c.Query("merchant_type"),
	}

	// 解析status
	if s := c.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			t := int16(v)
			params.Status = &t
		}
	}

	// 解析is_direct
	if d := c.Query("is_direct"); d != "" {
		if v, err := strconv.ParseBool(d); err == nil {
			params.IsDirect = &v
		}
	}

	merchants, err := h.merchantRepo.ExportMerchants(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 生成CSV格式数据
	csvData := "商户编号,商户名称,终端SN,状态,商户类型,归属类型,贷记卡费率,借记卡费率,创建时间\n"
	for _, m := range merchants {
		ownerType := "直营"
		if !m.IsDirect {
			ownerType = "团队"
		}
		csvData += m.MerchantNo + "," +
			m.MerchantName + "," +
			m.TerminalSN + "," +
			getMerchantStatusName(m.Status) + "," +
			getMerchantTypeName(m.MerchantType) + "," +
			ownerType + "," +
			m.CreditRate + "," +
			m.DebitRate + "," +
			m.CreatedAt.Format("2006-01-02 15:04:05") + "\n"
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=merchants_"+time.Now().Format("20060102150405")+".csv")
	c.Header("Content-Transfer-Encoding", "binary")

	c.String(http.StatusOK, "\xEF\xBB\xBF"+csvData) // UTF-8 BOM for Excel compatibility
}

// RegisterMerchantRoutes 注册商户路由
func RegisterMerchantRoutes(r *gin.RouterGroup, h *MerchantHandler, authService *service.AuthService) {
	merchants := r.Group("/merchants")
	merchants.Use(middleware.AuthMiddleware(authService))
	{
		// 列表和统计
		merchants.GET("", h.GetMerchantList)
		merchants.GET("/stats", h.GetMerchantStats)
		merchants.GET("/stats/extended", h.GetExtendedStats)
		merchants.GET("/export", h.ExportMerchants)

		// CRUD操作
		merchants.POST("", h.CreateMerchant)
		merchants.GET("/:id", h.GetMerchantDetail)
		merchants.PUT("/:id", h.UpdateMerchant)
		merchants.DELETE("/:id", h.DeleteMerchant)

		// 商户操作
		merchants.PUT("/:id/status", h.UpdateMerchantStatus)
		merchants.PUT("/:id/rate", h.UpdateMerchantRate)
		merchants.POST("/:id/register", h.RegisterMerchant)

		// 交易记录
		merchants.GET("/:id/transactions", h.GetMerchantTransactions)
	}
}
