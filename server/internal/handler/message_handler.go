package handler

import (
	"strconv"
	"strings"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageRepo *repository.GormMessageRepository
}

func NewMessageHandler(messageRepo *repository.GormMessageRepository) *MessageHandler {
	return &MessageHandler{
		messageRepo: messageRepo,
	}
}

// GetMessageList 获取消息列表
// @Summary 获取消息列表
// @Description 获取当前代理商的消息列表，支持按类型和分类筛选
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param type query string false "消息类型ID，多个用逗号分隔（如1,2,3）"
// @Param category query string false "消息分类（all/profit/register/consumption/system）"
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

	var types []int16
	typeStr := c.Query("type")
	category := c.Query("category")

	if category != "" && category != "all" {
		types = models.GetMessageTypesByCategory(category)
	} else if typeStr != "" {
		typeStrs := strings.Split(typeStr, ",")
		for _, ts := range typeStrs {
			if t, err := strconv.ParseInt(strings.TrimSpace(ts), 10, 16); err == nil {
				types = append(types, int16(t))
			}
		}
	}

	var messages []*models.Message
	var total int64
	var err error

	if len(types) > 0 {
		messages, err = h.messageRepo.FindByAgentIDAndTypes(agentID, types, pageSize, offset)
		if err == nil {
			total, err = h.messageRepo.CountByAgentIDAndTypes(agentID, types)
		}
	} else {
		messages, err = h.messageRepo.FindByAgentID(agentID, pageSize, offset)
		if err == nil {
			total, err = h.messageRepo.CountByAgentIDAndTypes(agentID, nil)
		}
	}

	if err != nil {
		response.InternalError(c, "查询失败: "+err.Error())
		return
	}

	list := make([]gin.H, 0, len(messages))
	for _, m := range messages {
		list = append(list, gin.H{
			"id":           m.ID,
			"title":        m.Title,
			"content":      m.Content,
			"message_type": m.MessageType,
			"type_name":    models.GetMessageTypeName(m.MessageType),
			"is_read":      m.IsRead,
			"expire_at":    m.ExpireAt,
			"created_at":   m.CreatedAt,
		})
	}

	response.SuccessPage(c, list, total, page, pageSize)
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
		response.InternalError(c, "查询失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"count": len(messages),
	})
}

// GetMessageStats 获取消息分类统计
// @Summary 获取消息分类统计
// @Description 获取当前代理商的消息分类统计信息
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/stats [get]
func (h *MessageHandler) GetMessageStats(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)

	stats, err := h.messageRepo.GetStatsByAgentID(agentID)
	if err != nil {
		response.InternalError(c, "查询失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"total":             stats.Total,
		"unread_total":      stats.UnreadTotal,
		"profit_count":      stats.ProfitCount,
		"register_count":    stats.RegisterCount,
		"consumption_count": stats.ConsumptionCount,
		"system_count":      stats.SystemCount,
	})
}

// GetMessageDetail 获取消息详情
// @Summary 获取消息详情
// @Description 获取消息详情并自动标记为已读
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "消息ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/{id} [get]
func (h *MessageHandler) GetMessageDetail(c *gin.Context) {
	agentID := middleware.GetCurrentAgentID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	message, err := h.messageRepo.FindByID(id)
	if err != nil {
		response.NotFound(c, "消息不存在")
		return
	}

	if message.AgentID != agentID {
		response.Forbidden(c, "无权查看此消息")
		return
	}

	if !message.IsRead {
		_ = h.messageRepo.MarkAsRead(id)
		message.IsRead = true
	}

	response.Success(c, gin.H{
		"id":           message.ID,
		"title":        message.Title,
		"content":      message.Content,
		"message_type": message.MessageType,
		"type_name":    models.GetMessageTypeName(message.MessageType),
		"is_read":      message.IsRead,
		"related_id":   message.RelatedID,
		"related_type": message.RelatedType,
		"expire_at":    message.ExpireAt,
		"created_at":   message.CreatedAt,
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
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.messageRepo.MarkAsRead(id); err != nil {
		response.InternalError(c, "操作失败: "+err.Error())
		return
	}

	response.SuccessMessage(c, "success")
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
		response.InternalError(c, "操作失败: "+err.Error())
		return
	}

	response.SuccessMessage(c, "success")
}

// GetMessageTypes 获取消息类型列表
// @Summary 获取消息类型列表
// @Description 获取所有消息类型及分类信息
// @Tags 消息中心
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/types [get]
func (h *MessageHandler) GetMessageTypes(c *gin.Context) {
	types := []gin.H{
		{"value": models.MessageTypeProfit, "label": "交易分润", "category": "profit"},
		{"value": models.MessageTypeActivation, "label": "激活奖励", "category": "profit"},
		{"value": models.MessageTypeDeposit, "label": "押金返现", "category": "profit"},
		{"value": models.MessageTypeSimCashback, "label": "流量返现", "category": "profit"},
		{"value": models.MessageTypeRefund, "label": "退款撤销", "category": "system"},
		{"value": models.MessageTypeAnnouncement, "label": "系统公告", "category": "system"},
		{"value": models.MessageTypeNewAgent, "label": "新代理注册", "category": "register"},
		{"value": models.MessageTypeTransaction, "label": "交易通知", "category": "consumption"},
	}

	categories := []gin.H{
		{"value": "all", "label": "全部"},
		{"value": "profit", "label": "分润"},
		{"value": "register", "label": "注册"},
		{"value": "consumption", "label": "消费"},
		{"value": "system", "label": "系统"},
	}

	response.Success(c, gin.H{
		"types":      types,
		"categories": categories,
	})
}

func RegisterMessageRoutes(r *gin.RouterGroup, h *MessageHandler, authService *service.AuthService) {
	messages := r.Group("/messages")
	messages.Use(middleware.AuthMiddleware(authService))
	{
		messages.GET("", h.GetMessageList)
		messages.GET("/unread-count", h.GetUnreadCount)
		messages.GET("/stats", h.GetMessageStats)
		messages.GET("/types", h.GetMessageTypes)
		messages.GET("/:id", h.GetMessageDetail)
		messages.PUT("/:id/read", h.MarkAsRead)
		messages.PUT("/read-all", h.MarkAllAsRead)
	}
}
