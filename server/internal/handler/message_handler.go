package handler

import (
	"net/http"
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageRepo *repository.GormMessageRepository
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(messageRepo *repository.GormMessageRepository) *MessageHandler {
	return &MessageHandler{
		messageRepo: messageRepo,
	}
}

// GetMessageList 获取消息列表
// @Summary 获取消息列表
// @Description 获取当前代理商的消息列表
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages [get]
func (h *MessageHandler) GetMessageList(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	messages, err := h.messageRepo.FindByAgentID(agentID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为前端友好格式
	list := make([]gin.H, 0, len(messages))
	for _, m := range messages {
		list = append(list, gin.H{
			"id":           m.ID,
			"title":        m.Title,
			"content":      m.Content,
			"message_type": m.MessageType,
			"is_read":      m.IsRead,
			"created_at":   m.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUnreadCount 获取未读消息数
// @Summary 获取未读消息数
// @Description 获取当前代理商的未读消息数量
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/unread-count [get]
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	messages, err := h.messageRepo.FindUnreadByAgentID(agentID)
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
			"count": len(messages),
		},
	})
}

// MarkAsRead 标记消息已读
// @Summary 标记消息已读
// @Description 将指定消息标记为已读
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "消息ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/{id}/read [put]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.messageRepo.MarkAsRead(id); err != nil {
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

// MarkAllAsRead 标记全部已读
// @Summary 标记全部已读
// @Description 将当前代理商的所有消息标记为已读
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/read-all [put]
func (h *MessageHandler) MarkAllAsRead(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	if err := h.messageRepo.MarkAllAsRead(agentID); err != nil {
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

// RegisterMessageRoutes 注册消息中心路由
func RegisterMessageRoutes(r *gin.RouterGroup, h *MessageHandler, authService *service.AuthService) {
	messages := r.Group("/messages")
	messages.Use(middleware.AuthMiddleware(authService))
	{
		messages.GET("", h.GetMessageList)
		messages.GET("/unread-count", h.GetUnreadCount)
		messages.PUT("/:id/read", h.MarkAsRead)
		messages.PUT("/read-all", h.MarkAllAsRead)
	}
}
