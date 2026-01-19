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

// ProfitHandler 分润处理器
type ProfitHandler struct {
	profitRepo *repository.GormProfitRecordRepository
}

// NewProfitHandler 创建分润处理器
func NewProfitHandler(profitRepo *repository.GormProfitRecordRepository) *ProfitHandler {
	return &ProfitHandler{
		profitRepo: profitRepo,
	}
}

// GetProfitList 获取分润列表
// @Summary 获取分润列表
// @Description 获取当前代理商的分润记录列表
// @Tags 分润管理
// @Produce json
// @Security ApiKeyAuth
// @Param profit_type query int false "分润类型"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profits [get]
func (h *ProfitHandler) GetProfitList(c *gin.Context) {
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

	var profitType *int16
	if pt := c.Query("profit_type"); pt != "" {
		if v, err := strconv.ParseInt(pt, 10, 16); err == nil {
			t := int16(v)
			profitType = &t
		}
	}

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

	records, total, err := h.profitRepo.FindByAgentID(agentID, profitType, startTime, endTime, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(records))
	for _, r := range records {
		list = append(list, gin.H{
			"id":               r.ID,
			"order_no":         r.OrderNo,
			"profit_type":      r.ProfitType,
			"profit_type_name": getProfitTypeName(r.ProfitType),
			"trade_amount":     r.TradeAmount,
			"trade_amount_yuan": float64(r.TradeAmount) / 100,
			"self_rate":        r.SelfRate,
			"lower_rate":       r.LowerRate,
			"rate_diff":        r.RateDiff,
			"profit_amount":    r.ProfitAmount,
			"profit_amount_yuan": float64(r.ProfitAmount) / 100,
			"wallet_type":      r.WalletType,
			"wallet_type_name": getWalletTypeStr(r.WalletType),
			"wallet_status":    r.WalletStatus,
			"created_at":       r.CreatedAt,
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

// GetProfitStats 获取分润统计
// @Summary 获取分润统计
// @Description 获取今日和本月的分润统计
// @Tags 分润管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profits/stats [get]
func (h *ProfitHandler) GetProfitStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	now := time.Now()

	todayStats, _ := h.profitRepo.GetAgentDailyProfitStats(agentID, now)
	monthStats, _ := h.profitRepo.GetAgentMonthlyProfitStats(agentID, now)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"today": gin.H{
				"amount":      todayStats.TotalAmount,
				"amount_yuan": float64(todayStats.TotalAmount) / 100,
				"count":       todayStats.TotalCount,
			},
			"month": gin.H{
				"amount":      monthStats.TotalAmount,
				"amount_yuan": float64(monthStats.TotalAmount) / 100,
				"count":       monthStats.TotalCount,
			},
		},
	})
}

// getProfitTypeName 分润类型名称
func getProfitTypeName(profitType int16) string {
	switch profitType {
	case 1:
		return "交易分润"
	case 2:
		return "激活奖励"
	case 3:
		return "押金返现"
	case 4:
		return "流量返现"
	default:
		return "未知"
	}
}

// getWalletTypeStr 钱包类型名称
func getWalletTypeStr(walletType int16) string {
	switch walletType {
	case 1:
		return "分润钱包"
	case 2:
		return "服务费钱包"
	case 3:
		return "奖励钱包"
	default:
		return "未知"
	}
}

// RegisterProfitRoutes 注册分润路由
func RegisterProfitRoutes(r *gin.RouterGroup, h *ProfitHandler, authService *service.AuthService) {
	profits := r.Group("/profits")
	profits.Use(middleware.AuthMiddleware(authService))
	{
		profits.GET("", h.GetProfitList)
		profits.GET("/stats", h.GetProfitStats)
	}
}
