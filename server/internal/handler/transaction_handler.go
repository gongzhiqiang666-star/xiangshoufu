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

// TransactionHandler 交易处理器
type TransactionHandler struct {
	transactionRepo *repository.GormTransactionRepository
}

// NewTransactionHandler 创建交易处理器
func NewTransactionHandler(transactionRepo *repository.GormTransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: transactionRepo,
	}
}

// GetTransactionList 获取交易列表
// @Summary 获取交易列表
// @Description 获取当前代理商的交易流水列表
// @Tags 交易管理
// @Produce json
// @Security ApiKeyAuth
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param trade_type query int false "交易类型"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetTransactionList(c *gin.Context) {
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

	var tradeType *int16
	if tt := c.Query("trade_type"); tt != "" {
		if v, err := strconv.ParseInt(tt, 10, 16); err == nil {
			t := int16(v)
			tradeType = &t
		}
	}

	transactions, total, err := h.transactionRepo.FindByAgentID(agentID, startTime, endTime, tradeType, pageSize, offset)
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
			"id":            tx.ID,
			"order_no":      tx.OrderNo,
			"trade_no":      tx.TradeNo,
			"terminal_sn":   tx.TerminalSN,
			"trade_type":    tx.TradeType,
			"trade_type_name": getTradeTypeName(tx.TradeType),
			"pay_type":      tx.PayType,
			"pay_type_name": getPayTypeName(tx.PayType),
			"card_type":     tx.CardType,
			"amount":        tx.Amount,
			"amount_yuan":   float64(tx.Amount) / 100,
			"fee":           tx.Fee,
			"fee_yuan":      float64(tx.Fee) / 100,
			"rate":          tx.Rate,
			"card_no":       tx.CardNo,
			"trade_time":    tx.TradeTime,
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

// GetTransactionStats 获取交易统计
// @Summary 获取交易统计
// @Description 获取今日和本月的交易统计
// @Tags 交易管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions/stats [get]
func (h *TransactionHandler) GetTransactionStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	now := time.Now()

	todayStats, _ := h.transactionRepo.GetAgentDailyStats(agentID, now)
	monthStats, _ := h.transactionRepo.GetAgentMonthlyStats(agentID, now)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"today": gin.H{
				"amount":      todayStats.TotalAmount,
				"amount_yuan": float64(todayStats.TotalAmount) / 100,
				"count":       todayStats.TotalCount,
				"fee":         todayStats.TotalFee,
			},
			"month": gin.H{
				"amount":      monthStats.TotalAmount,
				"amount_yuan": float64(monthStats.TotalAmount) / 100,
				"count":       monthStats.TotalCount,
				"fee":         monthStats.TotalFee,
			},
		},
	})
}

// GetTransactionTrend 获取交易趋势
// @Summary 获取交易趋势
// @Description 获取最近N天的交易趋势数据
// @Tags 交易管理
// @Produce json
// @Security ApiKeyAuth
// @Param days query int false "天数（默认7）"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions/trend [get]
func (h *TransactionHandler) GetTransactionTrend(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days <= 0 || days > 30 {
		days = 7
	}

	endDate := time.Now().AddDate(0, 0, 1)
	startDate := time.Now().AddDate(0, 0, -days+1)

	trend, err := h.transactionRepo.GetTransactionTrend(agentID, startDate, endDate)
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
			"trend": trend,
			"days":  days,
		},
	})
}

// getTradeTypeName 交易类型名称
func getTradeTypeName(tradeType int16) string {
	switch tradeType {
	case 1:
		return "消费"
	case 2:
		return "撤销"
	case 3:
		return "退货"
	default:
		return "未知"
	}
}

// getPayTypeName 支付类型名称
func getPayTypeName(payType int16) string {
	switch payType {
	case 1:
		return "刷卡"
	case 2:
		return "微信"
	case 3:
		return "支付宝"
	case 4:
		return "云闪付"
	default:
		return "未知"
	}
}

// RegisterTransactionRoutes 注册交易路由
func RegisterTransactionRoutes(r *gin.RouterGroup, h *TransactionHandler, authService *service.AuthService) {
	transactions := r.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(authService))
	{
		transactions.GET("", h.GetTransactionList)
		transactions.GET("/stats", h.GetTransactionStats)
		transactions.GET("/trend", h.GetTransactionTrend)
	}
}
