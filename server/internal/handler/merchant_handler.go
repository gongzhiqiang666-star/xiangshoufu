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
}

// NewMerchantHandler 创建商户处理器
func NewMerchantHandler(
	merchantRepo *repository.GormMerchantRepository,
	transactionRepo *repository.GormTransactionRepository,
) *MerchantHandler {
	return &MerchantHandler{
		merchantRepo:    merchantRepo,
		transactionRepo: transactionRepo,
	}
}

// GetMerchantList 获取商户列表
// @Summary 获取商户列表
// @Description 获取当前代理商的商户列表
// @Tags 商户管理
// @Produce json
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态"
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

	keyword := c.Query("keyword")

	var status *int16
	if s := c.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			t := int16(v)
			status = &t
		}
	}

	merchants, total, err := h.merchantRepo.FindByAgentID(agentID, keyword, status, pageSize, offset)
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
			"id":             m.ID,
			"merchant_no":    m.MerchantNo,
			"merchant_name":  m.MerchantName,
			"terminal_sn":    m.TerminalSN,
			"status":         m.Status,
			"status_name":    getMerchantStatusName(m.Status),
			"approve_status": m.ApproveStatus,
			"approve_name":   getMerchantApproveStatusName(m.ApproveStatus),
			"mcc":            m.MCC,
			"credit_rate":    m.CreditRate,
			"debit_rate":     m.DebitRate,
			"created_at":     m.CreatedAt,
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

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":             merchant.ID,
			"merchant_no":    merchant.MerchantNo,
			"merchant_name":  merchant.MerchantName,
			"terminal_sn":    merchant.TerminalSN,
			"status":         merchant.Status,
			"status_name":    getMerchantStatusName(merchant.Status),
			"approve_status": merchant.ApproveStatus,
			"approve_name":   getMerchantApproveStatusName(merchant.ApproveStatus),
			"legal_name":     merchant.LegalName,
			"legal_id_card":  maskIDCard(merchant.LegalIDCard),
			"mcc":            merchant.MCC,
			"credit_rate":    merchant.CreditRate,
			"debit_rate":     merchant.DebitRate,
			"created_at":     merchant.CreatedAt,
			"updated_at":     merchant.UpdatedAt,
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

// RegisterMerchantRoutes 注册商户路由
func RegisterMerchantRoutes(r *gin.RouterGroup, h *MerchantHandler, authService *service.AuthService) {
	merchants := r.Group("/merchants")
	merchants.Use(middleware.AuthMiddleware(authService))
	{
		merchants.GET("", h.GetMerchantList)
		merchants.GET("/stats", h.GetMerchantStats)
		merchants.GET("/:id", h.GetMerchantDetail)
		merchants.GET("/:id/transactions", h.GetMerchantTransactions)
	}
}
