package handler

import (
	"strconv"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminMessageHandler struct {
	messageService *service.MessageService
	agentRepo      repository.AgentRepository
}

func NewAdminMessageHandler(messageService *service.MessageService, agentRepo repository.AgentRepository) *AdminMessageHandler {
	messageService.SetAgentRepo(agentRepo)
	return &AdminMessageHandler{
		messageService: messageService,
		agentRepo:      agentRepo,
	}
}

// SendMessage 发送消息
// @Summary 发送消息
// @Description 管理员发送消息到指定代理商
// @Tags 管理端-消息管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.SendMessageRequest true "发送消息请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/messages [post]
func (h *AdminMessageHandler) SendMessage(c *gin.Context) {
	var req service.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.MessageType != models.MessageTypeAnnouncement {
		response.BadRequest(c, "管理员只能发送系统公告类型的消息")
		return
	}

	var agentIDs []int64
	switch req.SendScope {
	case "all":
		agents, err := h.getAllAgentIDs()
		if err != nil {
			response.InternalError(c, "获取代理商列表失败: "+err.Error())
			return
		}
		agentIDs = agents
	case "agents":
		if len(req.AgentIDs) == 0 {
			response.BadRequest(c, "请指定代理商ID列表")
			return
		}
		agentIDs = req.AgentIDs
	case "level":
		if req.Level <= 0 || req.Level > 5 {
			response.BadRequest(c, "层级必须在1-5之间")
			return
		}
		agents, err := h.getAgentIDsByLevel(req.Level)
		if err != nil {
			response.InternalError(c, "获取代理商列表失败: "+err.Error())
			return
		}
		agentIDs = agents
	default:
		response.BadRequest(c, "无效的发送范围，支持: all/agents/level")
		return
	}

	if len(agentIDs) == 0 {
		response.BadRequest(c, "没有可发送的代理商")
		return
	}

	if err := h.messageService.AdminSendMessage(&req, agentIDs); err != nil {
		response.InternalError(c, "发送消息失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"sent_count": len(agentIDs),
	})
}

// GetMessageList 获取消息列表
// @Summary 获取消息列表
// @Description 管理员获取所有消息列表
// @Tags 管理端-消息管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/messages [get]
func (h *AdminMessageHandler) GetMessageList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	messages, total, err := h.messageService.AdminGetAllMessages(page, pageSize)
	if err != nil {
		response.InternalError(c, "查询失败: "+err.Error())
		return
	}

	list := make([]gin.H, 0, len(messages))
	for _, m := range messages {
		list = append(list, gin.H{
			"id":           m.ID,
			"agent_id":     m.AgentID,
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

// GetMessageDetail 获取消息详情
// @Summary 获取消息详情
// @Description 管理员获取消息详情
// @Tags 管理端-消息管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "消息ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/messages/{id} [get]
func (h *AdminMessageHandler) GetMessageDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	message, err := h.messageService.AdminGetMessageByID(id)
	if err != nil {
		response.NotFound(c, "消息不存在")
		return
	}

	response.Success(c, gin.H{
		"id":           message.ID,
		"agent_id":     message.AgentID,
		"title":        message.Title,
		"content":      message.Content,
		"message_type": message.MessageType,
		"type_name":    models.GetMessageTypeName(message.MessageType),
		"is_read":      message.IsRead,
		"is_pushed":    message.IsPushed,
		"related_id":   message.RelatedID,
		"related_type": message.RelatedType,
		"expire_at":    message.ExpireAt,
		"created_at":   message.CreatedAt,
	})
}

// DeleteMessage 删除消息
// @Summary 删除消息
// @Description 管理员删除消息
// @Tags 管理端-消息管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "消息ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/messages/{id} [delete]
func (h *AdminMessageHandler) DeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.messageService.AdminDeleteMessage(id); err != nil {
		response.InternalError(c, "删除失败: "+err.Error())
		return
	}

	response.SuccessMessage(c, "success")
}

// GetMessageTypes 获取消息类型列表
// @Summary 获取消息类型列表
// @Description 获取所有消息类型
// @Tags 管理端-消息管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/messages/types [get]
func (h *AdminMessageHandler) GetMessageTypes(c *gin.Context) {
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

	response.Success(c, types)
}

func (h *AdminMessageHandler) getAllAgentIDs() ([]int64, error) {
	if agentRepo, ok := h.agentRepo.(*repository.GormAgentRepository); ok {
		return agentRepo.GetAllAgentIDs()
	}
	return []int64{}, nil
}

func (h *AdminMessageHandler) getAgentIDsByLevel(level int) ([]int64, error) {
	if agentRepo, ok := h.agentRepo.(*repository.GormAgentRepository); ok {
		return agentRepo.GetAgentIDsByLevel(level)
	}
	return []int64{}, nil
}

func RegisterAdminMessageRoutes(r *gin.RouterGroup, h *AdminMessageHandler, authService *service.AuthService) {
	messages := r.Group("/admin/messages")
	messages.Use(middleware.AuthMiddleware(authService))
	messages.Use(middleware.AdminMiddleware())
	{
		messages.POST("", h.SendMessage)
		messages.GET("", h.GetMessageList)
		messages.GET("/types", h.GetMessageTypes)
		messages.GET("/:id", h.GetMessageDetail)
		messages.DELETE("/:id", h.DeleteMessage)
	}
}
