package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// GoodsDeductionHandler 货款代扣Handler
type GoodsDeductionHandler struct {
	goodsDeductionService *service.GoodsDeductionService
}

// NewGoodsDeductionHandler 创建货款代扣Handler
func NewGoodsDeductionHandler(goodsDeductionService *service.GoodsDeductionService) *GoodsDeductionHandler {
	return &GoodsDeductionHandler{
		goodsDeductionService: goodsDeductionService,
	}
}

// CreateGoodsDeductionRequest API请求
type CreateGoodsDeductionRequest struct {
	ToAgentID       int64                                   `json:"to_agent_id" binding:"required"`      // 下级代理商ID
	UnitPrice       int64                                   `json:"unit_price" binding:"required,gt=0"`  // 单价（分）
	DeductionSource int16                                   `json:"deduction_source" binding:"required"` // 扣款来源: 1=分润 2=服务费 3=两者
	Terminals       []CreateGoodsDeductionTerminalRequest   `json:"terminals" binding:"required,dive"`   // 终端列表
	AgreementURL    string                                  `json:"agreement_url"`                       // 协议文件URL
	Remark          string                                  `json:"remark"`                              // 备注
	DistributeID    *int64                                  `json:"distribute_id"`                       // 关联的终端划拨ID
}

// CreateGoodsDeductionTerminalRequest 终端请求
type CreateGoodsDeductionTerminalRequest struct {
	TerminalID int64  `json:"terminal_id" binding:"required"` // 终端ID
	TerminalSN string `json:"terminal_sn"`                    // 终端SN
	UnitPrice  int64  `json:"unit_price"`                     // 单价（分）
}

// CreateGoodsDeduction 创建货款代扣
// @Summary 创建货款代扣
// @Description 终端划拨时创建货款代扣，需要下级接收确认后生效
// @Tags 货款代扣
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateGoodsDeductionRequest true "创建货款代扣请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction [post]
func (h *GoodsDeductionHandler) CreateGoodsDeduction(c *gin.Context) {
	var req CreateGoodsDeductionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前代理商ID
	fromAgentID := getCurrentAgentID(c)
	if fromAgentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	createdBy := getCurrentUserID(c)

	// 转换终端列表
	terminals := make([]models.CreateGoodsDeductionTerminal, 0, len(req.Terminals))
	for _, t := range req.Terminals {
		terminals = append(terminals, models.CreateGoodsDeductionTerminal{
			TerminalID: t.TerminalID,
			TerminalSN: t.TerminalSN,
			UnitPrice:  t.UnitPrice,
		})
	}

	serviceReq := &models.CreateGoodsDeductionRequest{
		ToAgentID:       req.ToAgentID,
		UnitPrice:       req.UnitPrice,
		DeductionSource: req.DeductionSource,
		Terminals:       terminals,
		AgreementURL:    req.AgreementURL,
		Remark:          req.Remark,
		DistributeID:    req.DistributeID,
	}

	deduction, err := h.goodsDeductionService.CreateGoodsDeduction(serviceReq, fromAgentID, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功，等待下级接收确认",
		"data":    deduction,
	})
}

// AcceptGoodsDeduction 接收货款代扣
// @Summary 接收货款代扣
// @Description 下级代理商接收货款代扣，接收后开始生效
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "货款代扣ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/{id}/accept [post]
func (h *GoodsDeductionHandler) AcceptGoodsDeduction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	if err := h.goodsDeductionService.AcceptGoodsDeduction(id, agentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "接收成功",
	})
}

// RejectGoodsDeduction 拒绝货款代扣
// @Summary 拒绝货款代扣
// @Description 下级代理商拒绝货款代扣，拒绝后终端划拨也会失败
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "货款代扣ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/{id}/reject [post]
func (h *GoodsDeductionHandler) RejectGoodsDeduction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	if err := h.goodsDeductionService.RejectGoodsDeduction(id, agentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "拒绝成功",
	})
}

// GetGoodsDeduction 获取货款代扣详情
// @Summary 获取货款代扣详情
// @Description 获取货款代扣详情，包括扣款进度和终端列表
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "货款代扣ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/{id} [get]
func (h *GoodsDeductionHandler) GetGoodsDeduction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	deduction, err := h.goodsDeductionService.GetGoodsDeductionByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 计算进度
	progress := float64(0)
	if deduction.TotalAmount > 0 {
		progress = float64(deduction.DeductedAmount) / float64(deduction.TotalAmount) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":               deduction.ID,
			"deduction_no":     deduction.DeductionNo,
			"from_agent_id":    deduction.FromAgentID,
			"from_agent_name":  deduction.FromAgentName,
			"to_agent_id":      deduction.ToAgentID,
			"to_agent_name":    deduction.ToAgentName,
			"total_amount":     deduction.TotalAmount,
			"total_amount_yuan": float64(deduction.TotalAmount) / 100,
			"deducted_amount":  deduction.DeductedAmount,
			"deducted_amount_yuan": float64(deduction.DeductedAmount) / 100,
			"remaining_amount": deduction.RemainingAmount,
			"remaining_amount_yuan": float64(deduction.RemainingAmount) / 100,
			"deduction_source": deduction.DeductionSource,
			"source_name":      models.GetGoodsDeductionSourceName(deduction.DeductionSource),
			"terminal_count":   deduction.TerminalCount,
			"unit_price":       deduction.UnitPrice,
			"unit_price_yuan":  float64(deduction.UnitPrice) / 100,
			"status":           deduction.Status,
			"status_name":      models.GetGoodsDeductionStatusName(deduction.Status),
			"progress":         progress,
			"agreement_signed": deduction.AgreementSigned,
			"agreement_url":    deduction.AgreementURL,
			"remark":           deduction.Remark,
			"terminals":        deduction.Terminals,
			"created_at":       deduction.CreatedAt,
			"accepted_at":      deduction.AcceptedAt,
			"completed_at":     deduction.CompletedAt,
		},
	})
}

// GetSentList 获取我发起的货款代扣列表
// @Summary 获取我发起的货款代扣列表
// @Description 获取当前代理商发起的货款代扣列表
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param status query int false "状态筛选: 1=待接收 2=进行中 3=已完成 4=已拒绝"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/sent [get]
func (h *GoodsDeductionHandler) GetSentList(c *gin.Context) {
	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 解析状态筛选
	var statusFilter []int16
	if s := c.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			statusFilter = []int16{int16(v)}
		}
	}

	list, total, err := h.goodsDeductionService.GetSentList(agentID, statusFilter, page, pageSize)
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
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetReceivedList 获取我接收的货款代扣列表
// @Summary 获取我接收的货款代扣列表
// @Description 获取当前代理商接收的货款代扣列表
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param status query int false "状态筛选: 1=待接收 2=进行中 3=已完成 4=已拒绝"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/received [get]
func (h *GoodsDeductionHandler) GetReceivedList(c *gin.Context) {
	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 解析状态筛选
	var statusFilter []int16
	if s := c.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			statusFilter = []int16{int16(v)}
		}
	}

	list, total, err := h.goodsDeductionService.GetReceivedList(agentID, statusFilter, page, pageSize)
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
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDeductionDetails 获取扣款明细列表
// @Summary 获取扣款明细列表
// @Description 获取货款代扣的扣款明细记录
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "货款代扣ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/{id}/details [get]
func (h *GoodsDeductionHandler) GetDeductionDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	details, total, err := h.goodsDeductionService.GetDeductionDetails(id, page, pageSize)
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
			"list":      details,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSummary 获取货款代扣统计
// @Summary 获取货款代扣统计
// @Description 获取货款代扣的统计汇总信息
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param type query string false "类型: sent=我发起的 received=我接收的" default(received)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/summary [get]
func (h *GoodsDeductionHandler) GetSummary(c *gin.Context) {
	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	summaryType := c.DefaultQuery("type", "received")
	isSent := summaryType == "sent"

	summary, err := h.goodsDeductionService.GetSummary(agentID, isSent)
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
			"total_count":           summary.TotalCount,
			"pending_count":         summary.PendingCount,
			"in_progress_count":     summary.InProgressCount,
			"completed_count":       summary.CompletedCount,
			"total_amount":          summary.TotalAmount,
			"total_amount_yuan":     float64(summary.TotalAmount) / 100,
			"deducted_amount":       summary.DeductedAmount,
			"deducted_amount_yuan":  float64(summary.DeductedAmount) / 100,
			"remaining_amount":      summary.RemainingAmount,
			"remaining_amount_yuan": float64(summary.RemainingAmount) / 100,
		},
	})
}

// GetNotifications 获取货款代扣通知列表
// @Summary 获取货款代扣通知列表
// @Description 获取货款代扣相关的通知列表
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param is_read query bool false "是否已读筛选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/notifications [get]
func (h *GoodsDeductionHandler) GetNotifications(c *gin.Context) {
	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 解析已读筛选
	var isRead *bool
	if s := c.Query("is_read"); s != "" {
		v := s == "true" || s == "1"
		isRead = &v
	}

	notifications, total, err := h.goodsDeductionService.GetNotifications(agentID, isRead, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 获取未读数量
	unreadCount, _ := h.goodsDeductionService.GetUnreadNotificationCount(agentID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":         notifications,
			"total":        total,
			"unread_count": unreadCount,
			"page":         page,
			"page_size":    pageSize,
		},
	})
}

// MarkNotificationAsRead 标记通知为已读
// @Summary 标记通知为已读
// @Description 标记单条货款代扣通知为已读
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "通知ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/notifications/{id}/read [post]
func (h *GoodsDeductionHandler) MarkNotificationAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.goodsDeductionService.MarkNotificationAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "操作失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// MarkAllNotificationsAsRead 标记所有通知为已读
// @Summary 标记所有通知为已读
// @Description 标记当前代理商的所有货款代扣通知为已读
// @Tags 货款代扣
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goods-deduction/notifications/read-all [post]
func (h *GoodsDeductionHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	agentID := getCurrentAgentID(c)
	if agentID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		return
	}

	if err := h.goodsDeductionService.MarkAllNotificationsAsRead(agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "操作失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// RegisterGoodsDeductionRoutes 注册货款代扣路由
func RegisterGoodsDeductionRoutes(r *gin.RouterGroup, h *GoodsDeductionHandler) {
	gd := r.Group("/goods-deduction")
	{
		// 创建货款代扣
		gd.POST("", h.CreateGoodsDeduction)

		// 我发起的/我接收的列表
		gd.GET("/sent", h.GetSentList)
		gd.GET("/received", h.GetReceivedList)

		// 统计汇总
		gd.GET("/summary", h.GetSummary)

		// 通知相关
		gd.GET("/notifications", h.GetNotifications)
		gd.POST("/notifications/:id/read", h.MarkNotificationAsRead)
		gd.POST("/notifications/read-all", h.MarkAllNotificationsAsRead)

		// 单个货款代扣操作
		gd.GET("/:id", h.GetGoodsDeduction)
		gd.GET("/:id/details", h.GetDeductionDetails)
		gd.POST("/:id/accept", h.AcceptGoodsDeduction)
		gd.POST("/:id/reject", h.RejectGoodsDeduction)
	}
}
