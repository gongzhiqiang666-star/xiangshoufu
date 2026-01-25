package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"xiangshoufu/internal/service"
)

// ChannelHandler 通道处理器
type ChannelHandler struct {
	channelService *service.ChannelService
}

// NewChannelHandler 创建通道处理器
func NewChannelHandler(channelService *service.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
	}
}

// GetRateTypes 获取通道费率类型列表
// GET /api/admin/channels/:channelId/rate-types
func (h *ChannelHandler) GetRateTypes(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的通道ID",
		})
		return
	}

	rateTypes, err := h.channelService.GetRateTypes(channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rateTypes,
	})
}

// GetChannelList 获取通道列表
// GET /api/admin/channels
func (h *ChannelHandler) GetChannelList(c *gin.Context) {
	channels, err := h.channelService.GetEnabledChannels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": channels,
	})
}

// GetChannelDetail 获取通道详情
// GET /api/admin/channels/:channelId
func (h *ChannelHandler) GetChannelDetail(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的通道ID",
		})
		return
	}

	channel, err := h.channelService.GetChannelByID(channelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "通道不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": channel,
	})
}
