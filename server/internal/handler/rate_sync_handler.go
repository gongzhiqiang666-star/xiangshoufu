package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/service"
)

// RateSyncHandler 费率同步处理器
type RateSyncHandler struct {
	rateSyncService *service.RateSyncService
}

// NewRateSyncHandler 创建费率同步处理器
func NewRateSyncHandler(rateSyncService *service.RateSyncService) *RateSyncHandler {
	return &RateSyncHandler{
		rateSyncService: rateSyncService,
	}
}

// RateSyncLogResponse 费率同步日志响应
type RateSyncLogResponse struct {
	ID              int64   `json:"id"`
	MerchantID      int64   `json:"merchant_id"`
	MerchantNo      string  `json:"merchant_no"`
	TerminalSN      string  `json:"terminal_sn"`
	ChannelCode     string  `json:"channel_code"`
	OldCreditRate   float64 `json:"old_credit_rate"`
	OldDebitRate    float64 `json:"old_debit_rate"`
	NewCreditRate   float64 `json:"new_credit_rate"`
	NewDebitRate    float64 `json:"new_debit_rate"`
	SyncStatus      int     `json:"sync_status"`
	SyncStatusName  string  `json:"sync_status_name"`
	ChannelTradeNo  string  `json:"channel_trade_no"`
	ErrorMessage    string  `json:"error_message"`
	CreatedAt       string  `json:"created_at"`
	SyncedAt        string  `json:"synced_at"`
}

// GetSyncLogs 获取费率同步日志列表
// @Summary 获取费率同步日志
// @Description 查询商户的费率同步记录
// @Tags 费率同步
// @Accept json
// @Produce json
// @Param merchant_id query int false "商户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} Response
// @Router /api/v1/rate-sync/logs [get]
func (h *RateSyncHandler) GetSyncLogs(c *gin.Context) {
	merchantIDStr := c.Query("merchant_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var merchantID int64
	if merchantIDStr != "" {
		merchantID, _ = strconv.ParseInt(merchantIDStr, 10, 64)
	}

	logs, total, err := h.rateSyncService.GetSyncLogsByMerchant(c.Request.Context(), merchantID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换响应
	items := make([]RateSyncLogResponse, 0, len(logs))
	for _, log := range logs {
		item := RateSyncLogResponse{
			ID:          log.ID,
			MerchantID:  log.MerchantID,
			MerchantNo:  log.MerchantNo,
			TerminalSN:  log.TerminalSN,
			ChannelCode: log.ChannelCode,
			SyncStatus:  int(log.SyncStatus),
			SyncStatusName: getSyncStatusName(int(log.SyncStatus)),
			ChannelTradeNo: log.ChannelTradeNo,
			ErrorMessage:   log.ErrorMessage,
			CreatedAt:      log.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if log.OldCreditRate != nil {
			item.OldCreditRate = *log.OldCreditRate
		}
		if log.OldDebitRate != nil {
			item.OldDebitRate = *log.OldDebitRate
		}
		if log.NewCreditRate != nil {
			item.NewCreditRate = *log.NewCreditRate
		}
		if log.NewDebitRate != nil {
			item.NewDebitRate = *log.NewDebitRate
		}
		if log.SyncedAt != nil {
			item.SyncedAt = log.SyncedAt.Format("2006-01-02 15:04:05")
		}

		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items": items,
			"total": total,
			"page":  page,
			"page_size": pageSize,
		},
	})
}

// GetSyncLogDetail 获取费率同步日志详情
// @Summary 获取费率同步日志详情
// @Description 查询单条费率同步记录详情
// @Tags 费率同步
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} Response
// @Router /api/v1/rate-sync/logs/{id} [get]
func (h *RateSyncHandler) GetSyncLogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	log, err := h.rateSyncService.GetSyncLogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "记录不存在",
		})
		return
	}

	item := RateSyncLogResponse{
		ID:             log.ID,
		MerchantID:    log.MerchantID,
		MerchantNo:    log.MerchantNo,
		TerminalSN:    log.TerminalSN,
		ChannelCode:   log.ChannelCode,
		SyncStatus:    int(log.SyncStatus),
		SyncStatusName: getSyncStatusName(int(log.SyncStatus)),
		ChannelTradeNo: log.ChannelTradeNo,
		ErrorMessage:  log.ErrorMessage,
		CreatedAt:     log.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if log.OldCreditRate != nil {
		item.OldCreditRate = *log.OldCreditRate
	}
	if log.OldDebitRate != nil {
		item.OldDebitRate = *log.OldDebitRate
	}
	if log.NewCreditRate != nil {
		item.NewCreditRate = *log.NewCreditRate
	}
	if log.NewDebitRate != nil {
		item.NewDebitRate = *log.NewDebitRate
	}
	if log.SyncedAt != nil {
		item.SyncedAt = log.SyncedAt.Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    item,
	})
}

// getSyncStatusName 获取同步状态名称
func getSyncStatusName(status int) string {
	switch status {
	case 0:
		return "待同步"
	case 1:
		return "同步中"
	case 2:
		return "同步成功"
	case 3:
		return "同步失败"
	default:
		return "未知"
	}
}
