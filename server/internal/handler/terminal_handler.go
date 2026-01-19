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
}

// NewTerminalHandler 创建终端处理器
func NewTerminalHandler(
	terminalRepo *repository.GormTerminalRepository,
	transactionRepo *repository.GormTransactionRepository,
) *TerminalHandler {
	return &TerminalHandler{
		terminalRepo:    terminalRepo,
		transactionRepo: transactionRepo,
	}
}

// GetTerminalList 获取终端列表
// @Summary 获取终端列表
// @Description 获取当前代理商的终端列表
// @Tags 终端管理
// @Produce json
// @Security ApiKeyAuth
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/terminals [get]
func (h *TerminalHandler) GetTerminalList(c *gin.Context) {
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

	var statusFilter []int16
	if s := c.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 16); err == nil {
			statusFilter = []int16{int16(v)}
		}
	}

	terminals, total, err := h.terminalRepo.FindByOwner(agentID, statusFilter, pageSize, offset)
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
		list = append(list, gin.H{
			"id":           t.ID,
			"terminal_sn":  t.TerminalSN,
			"channel_id":   t.ChannelID,
			"channel_code": t.ChannelCode,
			"brand_code":   t.BrandCode,
			"model_code":   t.ModelCode,
			"merchant_no":  t.MerchantNo,
			"status":       t.Status,
			"status_name":  getTerminalStatusName(t.Status),
			"activated_at": t.ActivatedAt,
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

	// 各状态统计
	var pendingCount, allocatedCount, boundCount, activatedCount int64

	h.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusPending, &pendingCount)
	h.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusAllocated, &allocatedCount)
	h.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusBound, &boundCount)
	h.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusActivated, &activatedCount)

	total := pendingCount + allocatedCount + boundCount + activatedCount

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":           total,
			"pending_count":   pendingCount,
			"allocated_count": allocatedCount,
			"bound_count":     boundCount,
			"activated_count": activatedCount,
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
		terminals.GET("/:sn", h.GetTerminalDetail)
	}
}
