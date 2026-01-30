package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"
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
		response.BadRequest(c, "无效的通道ID")
		return
	}

	rateTypes, err := h.channelService.GetRateTypes(channelID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, rateTypes)
}

// GetChannelList 获取通道列表
// GET /api/admin/channels
func (h *ChannelHandler) GetChannelList(c *gin.Context) {
	channels, err := h.channelService.GetEnabledChannels()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, channels)
}

// GetChannelDetail 获取通道详情
// GET /api/admin/channels/:channelId
func (h *ChannelHandler) GetChannelDetail(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通道ID")
		return
	}

	channel, err := h.channelService.GetChannelByID(channelID)
	if err != nil {
		response.NotFound(c, "通道不存在")
		return
	}

	response.Success(c, channel)
}
