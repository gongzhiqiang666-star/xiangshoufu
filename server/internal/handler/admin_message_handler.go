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

// AdminMessageHandler 管理端消息处理器
type AdminMessageHandler struct {
	messageService *service.MessageService
	agentRepo      repository.AgentRepository
}

// NewAdminMessageHandler 创建管理端消息处理器
func NewAdminMessageHandler(messageService *service.MessageService, agentRepo repository.AgentRepository) *AdminMessageHandler {
	// 设置代理商仓库用于按层级发送
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证消息类型（管理员只能发送系统公告）
	if req.MessageType != models.MessageTypeAnnouncement {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "管理员只能发送系统公告类型的消息",
		})
		return
	}

	// 获取目标代理商列表
	var agentIDs []int64
	switch req.SendScope {
	case "all":
		// 发送给所有代理商
		agents, err := h.getAllAgentIDs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取代理商列表失败: " + err.Error(),
			})
			return
		}
		agentIDs = agents
	case "agents":
		// 发送给指定代理商
		if len(req.AgentIDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请指定代理商ID列表",
			})
			return
		}
		agentIDs = req.AgentIDs
	case "level":
		// 发送给指定层级的代理商
		if req.Level <= 0 || req.Level > 5 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "层级必须在1-5之间",
			})
			return
		}
		agents, err := h.getAgentIDsByLevel(req.Level)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取代理商列表失败: " + err.Error(),
			})
			return
		}
		agentIDs = agents
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的发送范围，支持: all/agents/level",
		})
		return
	}

	if len(agentIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "没有可发送的代理商",
		})
		return
	}

	// 发送消息
	if err := h.messageService.AdminSendMessage(&req, agentIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送消息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"sent_count": len(agentIDs),
		},
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	message, err := h.messageService.AdminGetMessageByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "消息不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
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
		},
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := h.messageService.AdminDeleteMessage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
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

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    types,
	})
}

// getAllAgentIDs 获取所有代理商ID（简化实现）
func (h *AdminMessageHandler) getAllAgentIDs() ([]int64, error) {
	// 这里需要调用代理商仓库获取所有代理商ID
	// 简化实现：返回空列表，实际需要从数据库查询
	// TODO: 实现从agentRepo获取所有代理商ID
	return []int64{}, nil
}

// getAgentIDsByLevel 根据层级获取代理商ID（简化实现）
func (h *AdminMessageHandler) getAgentIDsByLevel(level int) ([]int64, error) {
	// 这里需要调用代理商仓库获取指定层级的代理商ID
	// 简化实现：返回空列表，实际需要从数据库查询
	// TODO: 实现从agentRepo获取指定层级的代理商ID
	return []int64{}, nil
}

// RegisterAdminMessageRoutes 注册管理端消息路由
func RegisterAdminMessageRoutes(r *gin.RouterGroup, h *AdminMessageHandler, authService *service.AuthService) {
	messages := r.Group("/admin/messages")
	messages.Use(middleware.AuthMiddleware(authService))
	messages.Use(middleware.AdminMiddleware()) // 管理员权限验证
	{
		messages.POST("", h.SendMessage)
		messages.GET("", h.GetMessageList)
		messages.GET("/types", h.GetMessageTypes)
		messages.GET("/:id", h.GetMessageDetail)
		messages.DELETE("/:id", h.DeleteMessage)
	}
}
