package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// TerminalHandler 终端处理器
type TerminalHandler struct {
	terminalRepo    *repository.GormTerminalRepository
	transactionRepo *repository.GormTransactionRepository
	terminalService *service.TerminalService
}

// NewTerminalHandler 创建终端处理器
func NewTerminalHandler(
	terminalRepo *repository.GormTerminalRepository,
	transactionRepo *repository.GormTransactionRepository,
	terminalService *service.TerminalService,
) *TerminalHandler {
	return &TerminalHandler{
		terminalRepo:    terminalRepo,
		transactionRepo: transactionRepo,
		terminalService: terminalService,
	}
}

// GetTerminalList 获取终端列表
// @Summary 获取终端列表
// @Description 获取当前代理商的终端列表，支持多条件筛选
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param status query int false "状态（兼容旧版）"
// @Param channel_id query int false "通道ID"
// @Param brand_code query string false "品牌编码"
// @Param model_code query string false "型号编码"
// @Param status_group query string false "状态分组: all/unstock/stocked/unbound/inactive/active"
// @Param keyword query string false "搜索关键词（终端SN或商户号）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals [get]
func (h *TerminalHandler) GetTerminalList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 解析筛选参数
	var channelID *int64
	if cid := c.Query("channel_id"); cid != "" {
		if v, err := strconv.ParseInt(cid, 10, 64); err == nil {
			channelID = &v
		}
	}

	statusGroup := c.Query("status_group")
	// 兼容旧版status参数
	if statusGroup == "" {
		if s := c.Query("status"); s != "" {
			if v, err := strconv.ParseInt(s, 10, 16); err == nil {
				switch int16(v) {
				case models.TerminalStatusPending:
					statusGroup = "unstock"
				case models.TerminalStatusAllocated:
					statusGroup = "stocked"
				case models.TerminalStatusBound:
					statusGroup = "inactive"
				case models.TerminalStatusActivated:
					statusGroup = "active"
				}
			}
		}
	}

	params := repository.TerminalFilterParams{
		OwnerAgentID: agentID,
		ChannelID:    channelID,
		BrandCode:    c.Query("brand_code"),
		ModelCode:    c.Query("model_code"),
		StatusGroup:  statusGroup,
		Keyword:      c.Query("keyword"),
		Limit:        pageSize,
		Offset:       offset,
	}

	terminals, total, err := h.terminalRepo.FindByOwnerWithFilter(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(terminals))
	for _, t := range terminals {
		// 计算状态分组
		terminalStatusGroup := getTerminalStatusGroup(t)
		list = append(list, gin.H{
			"id":           t.ID,
			"terminal_sn":  t.TerminalSN,
			"channel_id":   t.ChannelID,
			"channel_code": t.ChannelCode,
			"brand_code":   t.BrandCode,
			"model_code":   t.ModelCode,
			"merchant_id":  t.MerchantID,
			"merchant_no":  t.MerchantNo,
			"status":       t.Status,
			"status_name":  getTerminalStatusName(t.Status),
			"status_group": terminalStatusGroup,
			"activated_at": t.ActivatedAt,
			"bound_at":     t.BoundAt,
			"created_at":   t.CreatedAt,
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

// getTerminalStatusGroup 获取终端状态分组
func getTerminalStatusGroup(t *models.Terminal) string {
	switch t.Status {
	case models.TerminalStatusPending:
		return "unstock"
	case models.TerminalStatusAllocated:
		if t.MerchantID == nil {
			return "unbound"
		}
		return "stocked"
	case models.TerminalStatusBound:
		if t.ActivatedAt == nil {
			return "inactive"
		}
		return "active"
	case models.TerminalStatusActivated:
		return "active"
	default:
		return "other"
	}
}

// GetTerminalDetail 获取终端详情
// @Summary 获取终端详情
// @Description 获取指定终端的详细信息
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param sn path string true "终端SN"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/{sn} [get]
func (h *TerminalHandler) GetTerminalDetail(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	sn := c.Param("sn")

	terminal, err := h.terminalRepo.FindBySN(sn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "终端不存在",
		})
		return
	}

	// 验证权限
	if terminal.OwnerAgentID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权访问该终端",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":             terminal.ID,
			"terminal_sn":    terminal.TerminalSN,
			"channel_id":     terminal.ChannelID,
			"channel_code":   terminal.ChannelCode,
			"brand_code":     terminal.BrandCode,
			"model_code":     terminal.ModelCode,
			"merchant_id":    terminal.MerchantID,
			"merchant_no":    terminal.MerchantNo,
			"status":         terminal.Status,
			"status_name":    getTerminalStatusName(terminal.Status),
			"sim_fee_count":  terminal.SimFeeCount,
			"last_sim_fee_at": terminal.LastSimFeeAt,
			"activated_at":   terminal.ActivatedAt,
			"bound_at":       terminal.BoundAt,
			"created_at":     terminal.CreatedAt,
			"updated_at":     terminal.UpdatedAt,
		},
	})
}

// GetTerminalStats 获取终端统计
// @Summary 获取终端统计
// @Description 获取当前代理商的终端统计数据
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/stats [get]
func (h *TerminalHandler) GetTerminalStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	stats, err := h.terminalService.GetTerminalStats(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// ImportTerminalsRequest 终端入库请求
type ImportTerminalsRequest struct {
	ChannelID   int64    `json:"channel_id" binding:"required"`   // 通道ID
	ChannelCode string   `json:"channel_code"`                     // 通道编码
	BrandCode   string   `json:"brand_code"`                       // 品牌编码
	ModelCode   string   `json:"model_code"`                       // 型号编码
	SNList      []string `json:"sn_list" binding:"required,min=1"` // SN列表
}

// ImportTerminals 终端入库
// @Summary 终端入库
// @Description 批量导入终端SN
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ImportTerminalsRequest true "终端入库请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/import [post]
func (h *TerminalHandler) ImportTerminals(c *gin.Context) {
	var req ImportTerminalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := middleware.GetCurrentAgentID(c)
	userID := getCurrentUserID(c)

	serviceReq := &service.ImportTerminalsRequest{
		ChannelID:    req.ChannelID,
		ChannelCode:  req.ChannelCode,
		BrandCode:    req.BrandCode,
		ModelCode:    req.ModelCode,
		SNList:       req.SNList,
		OwnerAgentID: agentID,
		CreatedBy:    userID,
	}

	result, err := h.terminalService.ImportTerminals(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "入库完成",
		"data":    result,
	})
}

// RecallTerminalRequest 终端回拨请求
type RecallTerminalRequest struct {
	ToAgentID  int64  `json:"to_agent_id" binding:"required"` // 接收方代理商ID
	TerminalSN string `json:"terminal_sn" binding:"required"` // 终端SN
	ChannelID  int64  `json:"channel_id"`                     // 通道ID
	Remark     string `json:"remark"`                         // 备注
}

// RecallTerminal 终端回拨
// @Summary 终端回拨
// @Description 将终端回拨给上级代理商，APP端不能跨级，PC端可以跨级
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body RecallTerminalRequest true "终端回拨请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/recall [post]
func (h *TerminalHandler) RecallTerminal(c *gin.Context) {
	var req RecallTerminalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := middleware.GetCurrentAgentID(c)
	userID := getCurrentUserID(c)
	source := getRequestSource(c)

	serviceReq := &service.RecallTerminalRequest{
		FromAgentID: agentID,
		ToAgentID:   req.ToAgentID,
		TerminalSN:  req.TerminalSN,
		ChannelID:   req.ChannelID,
		Source:      source,
		Remark:      req.Remark,
		CreatedBy:   userID,
	}

	recall, err := h.terminalService.RecallTerminal(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "回拨成功，等待接收方确认",
		"data":    recall,
	})
}

// BatchRecallRequest 批量回拨请求
type BatchRecallRequest struct {
	ToAgentID   int64    `json:"to_agent_id" binding:"required"`   // 接收方代理商ID
	TerminalSNs []string `json:"terminal_sns" binding:"required,min=1"` // 终端SN列表
	Remark      string   `json:"remark"`                           // 备注
}

// BatchRecallTerminals 批量回拨终端
// @Summary 批量回拨终端
// @Description 批量将终端回拨给上级代理商
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchRecallRequest true "批量回拨请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/batch-recall [post]
func (h *TerminalHandler) BatchRecallTerminals(c *gin.Context) {
	var req BatchRecallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := middleware.GetCurrentAgentID(c)
	userID := getCurrentUserID(c)
	source := getRequestSource(c)

	serviceReq := &service.BatchRecallRequest{
		TerminalSNs: req.TerminalSNs,
		ToAgentID:   req.ToAgentID,
		FromAgentID: agentID,
		Source:      source,
		Remark:      req.Remark,
		CreatedBy:   userID,
	}

	result, err := h.terminalService.BatchRecallTerminals(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量回拨完成",
		"data":    result,
	})
}

// ConfirmRecall 确认接收终端回拨
// @Summary 确认接收终端回拨
// @Description 接收方确认接收终端回拨
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回拨记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/recall/{id}/confirm [post]
func (h *TerminalHandler) ConfirmRecall(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	confirmedBy := getCurrentUserID(c)

	if err := h.terminalService.ConfirmRecall(id, confirmedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "确认成功",
	})
}

// RejectRecall 拒绝终端回拨
// @Summary 拒绝终端回拨
// @Description 接收方拒绝终端回拨
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回拨记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/recall/{id}/reject [post]
func (h *TerminalHandler) RejectRecall(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	confirmedBy := getCurrentUserID(c)

	if err := h.terminalService.RejectRecall(id, confirmedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已拒绝",
	})
}

// CancelRecall 取消终端回拨
// @Summary 取消终端回拨
// @Description 回拨方取消终端回拨
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回拨记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/recall/{id}/cancel [post]
func (h *TerminalHandler) CancelRecall(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	cancelBy := middleware.GetCurrentAgentID(c)

	if err := h.terminalService.CancelRecall(id, cancelBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已取消",
	})
}

// GetRecallList 获取回拨列表
// @Summary 获取回拨列表
// @Description 获取终端回拨列表
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param direction query string true "方向: from(我回拨的) to(回拨给我的)"
// @Param status query []int false "状态筛选"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/recall [get]
func (h *TerminalHandler) GetRecallList(c *gin.Context) {
	direction := c.Query("direction")
	if direction != "from" && direction != "to" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "direction 参数必须是 from 或 to",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	agentID := middleware.GetCurrentAgentID(c)

	list, total, err := h.terminalService.GetRecallList(agentID, direction, nil, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
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

// getTerminalStatusName 终端状态名称
func getTerminalStatusName(status int16) string {
	switch status {
	case models.TerminalStatusPending:
		return "待分配"
	case models.TerminalStatusAllocated:
		return "已分配"
	case models.TerminalStatusBound:
		return "已绑定"
	case models.TerminalStatusActivated:
		return "已激活"
	case models.TerminalStatusUnbound:
		return "已解绑"
	case models.TerminalStatusRecycled:
		return "已回收"
	default:
		return "未知"
	}
}

// RegisterTerminalRoutes 注册终端路由
func RegisterTerminalRoutes(r *gin.RouterGroup, h *TerminalHandler, authService *service.AuthService) {
	terminals := r.Group("/terminals")
	terminals.Use(middleware.AuthMiddleware(authService))
	{
		terminals.GET("", h.GetTerminalList)
		terminals.GET("/stats", h.GetTerminalStats)
		terminals.GET("/filter-options", h.GetFilterOptions) // 新增：筛选选项API
		terminals.GET("/:sn", h.GetTerminalDetail)
		terminals.GET("/:sn/flow-logs", h.GetTerminalFlowLogs) // 新增：流动记录API

		// 终端入库
		terminals.POST("/import", h.ImportTerminals)

		// 终端回拨
		terminals.POST("/recall", h.RecallTerminal)
		terminals.POST("/batch-recall", h.BatchRecallTerminals)
		terminals.GET("/recall", h.GetRecallList)
		terminals.POST("/recall/:id/confirm", h.ConfirmRecall)
		terminals.POST("/recall/:id/reject", h.RejectRecall)
		terminals.POST("/recall/:id/cancel", h.CancelRecall)

		// 终端政策设置
		terminals.GET("/policy-options", h.GetPolicyOptions)
		terminals.GET("/:sn/policy", h.GetTerminalPolicy)
		terminals.POST("/batch-set-rate", h.BatchSetRate)
		terminals.POST("/batch-set-sim", h.BatchSetSimFee)
		terminals.POST("/batch-set-deposit", h.BatchSetDeposit)
	}
}

// GetPolicyOptions 获取政策选项
// @Summary 获取政策选项
// @Description 获取费率、流量费、押金的预设选项
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/policy-options [get]
func (h *TerminalHandler) GetPolicyOptions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"rate_options":              models.RateOptions,
			"first_sim_fee_options":     models.FirstSimFeeOptions,
			"non_first_sim_fee_options": models.NonFirstSimFeeOptions,
			"sim_fee_interval_options":  models.SimFeeIntervalDaysOptions,
			"deposit_options":           models.DepositOptions,
		},
	})
}

// GetTerminalPolicy 获取终端政策
// @Summary 获取终端政策
// @Description 获取指定终端的政策设置
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param sn path string true "终端SN"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/{sn}/policy [get]
func (h *TerminalHandler) GetTerminalPolicy(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	sn := c.Param("sn")

	policy, err := h.terminalService.GetTerminalPolicy(sn, agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    policy,
	})
}

// BatchSetRateRequest 批量设置费率请求
type BatchSetRateRequest struct {
	TerminalSNs  []string `json:"terminal_sns" binding:"required,min=1"` // 终端SN列表
	CreditRate   int      `json:"credit_rate" binding:"required,min=1"`  // 贷记卡费率(万分比)
	DebitRate    int      `json:"debit_rate"`                            // 借记卡费率
	DebitCap     int      `json:"debit_cap"`                             // 借记卡封顶(分)
	UnionpayRate int      `json:"unionpay_rate"`                         // 银联云闪付费率
	WechatRate   int      `json:"wechat_rate"`                           // 微信扫码费率
	AlipayRate   int      `json:"alipay_rate"`                           // 支付宝扫码费率
}

// BatchSetRate 批量设置费率
// @Summary 批量设置费率
// @Description 批量设置终端费率
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchSetRateRequest true "批量设置费率请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/batch-set-rate [post]
func (h *TerminalHandler) BatchSetRate(c *gin.Context) {
	var req BatchSetRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := middleware.GetCurrentAgentID(c)
	userID := getCurrentUserID(c)

	result, err := h.terminalService.BatchSetRate(&service.BatchSetRateRequest{
		TerminalSNs:  req.TerminalSNs,
		AgentID:      agentID,
		CreditRate:   req.CreditRate,
		DebitRate:    req.DebitRate,
		DebitCap:     req.DebitCap,
		UnionpayRate: req.UnionpayRate,
		WechatRate:   req.WechatRate,
		AlipayRate:   req.AlipayRate,
		UpdatedBy:    userID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
		"data":    result,
	})
}

// BatchSetSimFeeRequest 批量设置SIM卡费用请求
type BatchSetSimFeeRequest struct {
	TerminalSNs        []string `json:"terminal_sns" binding:"required,min=1"` // 终端SN列表
	FirstSimFee        int      `json:"first_sim_fee" binding:"min=0"`         // 首次流量费(分)
	NonFirstSimFee     int      `json:"non_first_sim_fee" binding:"min=0"`     // 非首次流量费(分)
	SimFeeIntervalDays int      `json:"sim_fee_interval_days" binding:"min=0"` // 非首次间隔天数
}

// BatchSetSimFee 批量设置SIM卡费用
// @Summary 批量设置SIM卡费用
// @Description 批量设置终端SIM卡费用
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchSetSimFeeRequest true "批量设置SIM卡费用请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/batch-set-sim [post]
func (h *TerminalHandler) BatchSetSimFee(c *gin.Context) {
	var req BatchSetSimFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := middleware.GetCurrentAgentID(c)
	userID := getCurrentUserID(c)

	result, err := h.terminalService.BatchSetSimFee(&service.BatchSetSimFeeRequest{
		TerminalSNs:        req.TerminalSNs,
		AgentID:            agentID,
		FirstSimFee:        req.FirstSimFee,
		NonFirstSimFee:     req.NonFirstSimFee,
		SimFeeIntervalDays: req.SimFeeIntervalDays,
		UpdatedBy:          userID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
		"data":    result,
	})
}

// BatchSetDepositRequest 批量设置押金请求
type BatchSetDepositRequest struct {
	TerminalSNs   []string `json:"terminal_sns" binding:"required,min=1"` // 终端SN列表
	DepositAmount int      `json:"deposit_amount" binding:"min=0"`        // 押金金额(分)，0表示无押金
}

// BatchSetDeposit 批量设置押金
// @Summary 批量设置押金
// @Description 批量设置终端押金
// @Tags 终端管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body BatchSetDepositRequest true "批量设置押金请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/batch-set-deposit [post]
func (h *TerminalHandler) BatchSetDeposit(c *gin.Context) {
	var req BatchSetDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := middleware.GetCurrentAgentID(c)
	userID := getCurrentUserID(c)

	result, err := h.terminalService.BatchSetDeposit(&service.BatchSetDepositRequest{
		TerminalSNs:   req.TerminalSNs,
		AgentID:       agentID,
		DepositAmount: req.DepositAmount,
		UpdatedBy:     userID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
		"data":    result,
	})
}

// GetFilterOptions 获取筛选选项
// @Summary 获取筛选选项
// @Description 获取终端列表的筛选选项（通道、终端类型、状态分组）
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param channel_id query int false "通道ID（用于获取该通道下的终端类型）"
// @Param brand_code query string false "品牌编码（用于获取状态分组统计）"
// @Param model_code query string false "型号编码（用于获取状态分组统计）"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/filter-options [get]
func (h *TerminalHandler) GetFilterOptions(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	// 解析筛选参数
	var channelID *int64
	if cid := c.Query("channel_id"); cid != "" {
		if v, err := strconv.ParseInt(cid, 10, 64); err == nil {
			channelID = &v
		}
	}
	brandCode := c.Query("brand_code")
	modelCode := c.Query("model_code")

	// 获取通道列表
	channels, err := h.terminalRepo.GetChannelList(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取通道列表失败: " + err.Error(),
		})
		return
	}

	// 获取终端类型列表
	terminalTypes, err := h.terminalRepo.GetTerminalTypes(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取终端类型失败: " + err.Error(),
		})
		return
	}

	// 获取状态分组统计
	statusGroups, err := h.terminalRepo.GetStatusGroupCounts(agentID, channelID, brandCode, modelCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取状态统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"channels":       channels,
			"terminal_types": terminalTypes,
			"status_groups":  statusGroups,
		},
	})
}

// GetTerminalFlowLogs 获取终端流动记录
// @Summary 获取终端流动记录
// @Description 获取指定终端的流动记录（下发、回拨、绑定、解绑、激活）
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param sn path string true "终端SN"
// @Param log_type query string false "日志类型: all/distribute/recall"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals/{sn}/flow-logs [get]
func (h *TerminalHandler) GetTerminalFlowLogs(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	sn := c.Param("sn")

	// 验证终端权限
	terminal, err := h.terminalRepo.FindBySN(sn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "终端不存在",
		})
		return
	}

	if terminal.OwnerAgentID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权访问该终端",
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

	logType := c.DefaultQuery("log_type", "all")

	logs, total, err := h.terminalRepo.GetTerminalFlowLogs(sn, logType, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取流动记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"terminal": gin.H{
				"terminal_sn":  terminal.TerminalSN,
				"channel_id":   terminal.ChannelID,
				"channel_code": terminal.ChannelCode,
				"brand_code":   terminal.BrandCode,
				"model_code":   terminal.ModelCode,
			},
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
