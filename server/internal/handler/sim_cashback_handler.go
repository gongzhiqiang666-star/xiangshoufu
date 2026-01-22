package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// SimCashbackHandler 流量费返现Handler
type SimCashbackHandler struct {
	simCashbackService *service.SimCashbackService
}

// NewSimCashbackHandler 创建流量费返现Handler
func NewSimCashbackHandler(simCashbackService *service.SimCashbackService) *SimCashbackHandler {
	return &SimCashbackHandler{
		simCashbackService: simCashbackService,
	}
}

// GetCashbackRecordsRequest 获取返现记录请求
type GetCashbackRecordsRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// GetCashbackRecords 获取返现记录列表
// @Summary 获取流量费返现记录
// @Description 获取当前代理商的流量费返现记录列表
// @Tags 流量费返现
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sim-cashback/records [get]
func (h *SimCashbackHandler) GetCashbackRecords(c *gin.Context) {
	var req GetCashbackRecordsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	agentID := getCurrentAgentID(c)
	offset := (req.Page - 1) * req.PageSize

	records, total, err := h.simCashbackService.GetCashbackRecords(agentID, req.PageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	// 转换为前端友好的格式
	list := make([]gin.H, 0, len(records))
	for _, record := range records {
		list = append(list, gin.H{
			"id":              record.ID,
			"terminal_sn":     record.TerminalSN,
			"sim_fee_count":   record.SimFeeCount,
			"sim_fee_amount":  record.SimFeeAmount,
			"cashback_tier":   record.CashbackTier,
			"tier_name":       getCashbackTierName(record.CashbackTier),
			"self_cashback":   record.SelfCashback,
			"actual_cashback": record.ActualCashback,
			"wallet_status":   record.WalletStatus,
			"status_name":     getWalletStatusName(record.WalletStatus),
			"created_at":      record.CreatedAt,
			"processed_at":    record.ProcessedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetCashbackStats 获取返现统计
// @Summary 获取流量费返现统计
// @Description 获取当前代理商的流量费返现统计数据
// @Tags 流量费返现
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sim-cashback/stats [get]
func (h *SimCashbackHandler) GetCashbackStats(c *gin.Context) {
	agentID := getCurrentAgentID(c)

	stats, err := h.simCashbackService.GetCashbackStatsByAgent(agentID)
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
			"total_cashback":       stats.TotalCashback,
			"total_cashback_yuan":  float64(stats.TotalCashback) / 100,
			"first_tier_cashback":  stats.FirstTierCashback,
			"second_tier_cashback": stats.SecondTierCashback,
			"third_tier_cashback":  stats.ThirdTierCashback,
			"total_count":          stats.TotalCount,
		},
	})
}

// GetCashbackTiers 获取返现档次说明
// @Summary 获取返现档次说明
// @Description 获取流量费返现的三档说明
// @Tags 流量费返现
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sim-cashback/tiers [get]
func (h *SimCashbackHandler) GetCashbackTiers(c *gin.Context) {
	tiers := []gin.H{
		{
			"tier":        models.SimCashbackTierFirst,
			"name":        "首次",
			"description": "终端首次缴纳流量费",
		},
		{
			"tier":        models.SimCashbackTierSecond,
			"name":        "第2次",
			"description": "终端第二次缴纳流量费",
		},
		{
			"tier":        models.SimCashbackTierThirdPlus,
			"name":        "第3次及以后",
			"description": "终端第三次及后续缴纳流量费",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"tiers": tiers,
			"rule":  "流量费返现按级差计算：实际返现 = 自身配置 - 下级配置",
		},
	})
}

// GetCashbackRecordDetail 获取返现记录详情
// @Summary 获取返现记录详情
// @Tags 流量费返现
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sim-cashback/records/{id} [get]
func (h *SimCashbackHandler) GetCashbackRecordDetail(c *gin.Context) {
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
	record, err := h.simCashbackService.GetCashbackRecordByID(id, agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "记录不存在或无权访问",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":              record.ID,
			"terminal_sn":     record.TerminalSN,
			"channel_id":      record.ChannelID,
			"agent_id":        record.AgentID,
			"sim_fee_count":   record.SimFeeCount,
			"sim_fee_amount":  record.SimFeeAmount,
			"cashback_tier":   record.CashbackTier,
			"tier_name":       getCashbackTierName(record.CashbackTier),
			"self_cashback":   record.SelfCashback,
			"upper_cashback":  record.UpperCashback,
			"actual_cashback": record.ActualCashback,
			"source_agent_id": record.SourceAgentID,
			"wallet_type":     record.WalletType,
			"wallet_status":   record.WalletStatus,
			"status_name":     getWalletStatusName(record.WalletStatus),
			"created_at":      record.CreatedAt,
			"processed_at":    record.ProcessedAt,
		},
	})
}

// RegisterSimCashbackRoutes 注册流量费返现路由
func RegisterSimCashbackRoutes(r *gin.RouterGroup, h *SimCashbackHandler) {
	simCashback := r.Group("/sim-cashback")
	{
		simCashback.GET("/records", h.GetCashbackRecords)
		simCashback.GET("/records/:id", h.GetCashbackRecordDetail)
		simCashback.GET("/stats", h.GetCashbackStats)
		simCashback.GET("/tiers", h.GetCashbackTiers)
	}
}

// getCashbackTierName 获取返现档次名称
func getCashbackTierName(tier int16) string {
	switch tier {
	case models.SimCashbackTierFirst:
		return "首次"
	case models.SimCashbackTierSecond:
		return "第2次"
	case models.SimCashbackTierThirdPlus:
		return "第3次及以后"
	default:
		return "未知"
	}
}

// getWalletStatusName 获取钱包状态名称
func getWalletStatusName(status int16) string {
	switch status {
	case 0:
		return "待入账"
	case 1:
		return "已入账"
	default:
		return "未知"
	}
}
