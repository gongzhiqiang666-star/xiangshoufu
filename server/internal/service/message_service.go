package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// MessageService 消息通知服务
type MessageService struct {
	messageRepo repository.MessageRepository
	agentRepo   repository.AgentRepository
	pushConfig  *PushConfig
	httpClient  *http.Client
}

// PushConfig 推送配置
type PushConfig struct {
	Enabled      bool   `json:"enabled"`
	Provider     string `json:"provider"` // jiguang, getui
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
	WebhookURL   string `json:"webhook_url"` // 企业微信/钉钉告警webhook
}

// NotificationMessage 通知消息
type NotificationMessage struct {
	AgentID     int64  `json:"agent_id"`
	MessageType int16  `json:"message_type"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	RelatedID   int64  `json:"related_id"`
	RelatedType string `json:"related_type"`
}

// SendMessageRequest 发送消息请求（管理端）
type SendMessageRequest struct {
	Title       string  `json:"title" binding:"required"`
	Content     string  `json:"content" binding:"required"`
	MessageType int16   `json:"message_type" binding:"required"`
	ExpireDays  int     `json:"expire_days"`   // 1-30天，默认3天
	SendScope   string  `json:"send_scope"`    // all/agents/level
	AgentIDs    []int64 `json:"agent_ids"`     // 指定代理商ID列表
	Level       int     `json:"level"`         // 指定层级
}

// NewMessageService 创建消息服务
func NewMessageService(messageRepo repository.MessageRepository, pushConfig *PushConfig) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		pushConfig:  pushConfig,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetAgentRepo 设置代理商仓库（用于按层级发送消息）
func (s *MessageService) SetAgentRepo(agentRepo repository.AgentRepository) {
	s.agentRepo = agentRepo
}

// ProcessMessage 处理通知消息
func (s *MessageService) ProcessMessage(msgBytes []byte) error {
	var msg NotificationMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return fmt.Errorf("unmarshal notification message failed: %w", err)
	}

	return s.SendNotification(&msg)
}

// SendNotification 发送通知
func (s *MessageService) SendNotification(msg *NotificationMessage) error {
	// 1. 创建站内消息
	expireAt := time.Now().Add(72 * time.Hour) // 3天后过期
	message := &models.Message{
		AgentID:     msg.AgentID,
		MessageType: msg.MessageType,
		Title:       msg.Title,
		Content:     msg.Content,
		IsRead:      false,
		IsPushed:    false,
		RelatedID:   msg.RelatedID,
		RelatedType: msg.RelatedType,
		ExpireAt:    &expireAt,
		CreatedAt:   time.Now(),
	}

	if err := s.messageRepo.Create(message); err != nil {
		return fmt.Errorf("create message failed: %w", err)
	}

	// 2. 推送通知到APP（如果启用）
	if s.pushConfig != nil && s.pushConfig.Enabled {
		if err := s.pushToApp(msg); err != nil {
			log.Printf("[MessageService] Push to app failed: %v", err)
			// 推送失败不影响站内消息
		} else {
			// 更新推送状态
			message.IsPushed = true
		}
	}

	return nil
}

// pushToApp 推送到APP
func (s *MessageService) pushToApp(msg *NotificationMessage) error {
	switch s.pushConfig.Provider {
	case "jiguang":
		return s.pushViaJiguang(msg)
	case "getui":
		return s.pushViaGetui(msg)
	default:
		return fmt.Errorf("unknown push provider: %s", s.pushConfig.Provider)
	}
}

// pushViaJiguang 通过极光推送
func (s *MessageService) pushViaJiguang(msg *NotificationMessage) error {
	// 极光推送API
	url := "https://api.jpush.cn/v3/push"

	payload := map[string]interface{}{
		"platform": "all",
		"audience": map[string]interface{}{
			"alias": []string{fmt.Sprintf("agent_%d", msg.AgentID)},
		},
		"notification": map[string]interface{}{
			"alert": msg.Content,
			"android": map[string]interface{}{
				"alert": msg.Content,
				"title": msg.Title,
			},
			"ios": map[string]interface{}{
				"alert": map[string]interface{}{
					"title": msg.Title,
					"body":  msg.Content,
				},
				"sound": "default",
			},
		},
		"options": map[string]interface{}{
			"time_to_live": 86400, // 24小时
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.SetBasicAuth(s.pushConfig.AppKey, s.pushConfig.MasterSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jiguang push failed with status: %d", resp.StatusCode)
	}

	return nil
}

// pushViaGetui 通过个推推送
func (s *MessageService) pushViaGetui(msg *NotificationMessage) error {
	// 个推推送API（简化实现）
	log.Printf("[MessageService] Getui push: agent=%d, title=%s", msg.AgentID, msg.Title)
	return nil
}

// SendAlert 发送告警到企业微信/钉钉
func (s *MessageService) SendAlert(alertMsg string) error {
	if s.pushConfig == nil || s.pushConfig.WebhookURL == "" {
		return nil
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("[分润系统告警]\n%s\n时间: %s", alertMsg, time.Now().Format("2006-01-02 15:04:05")),
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := s.httpClient.Post(s.pushConfig.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook failed with status: %d", resp.StatusCode)
	}

	log.Printf("[MessageService] Alert sent: %s", alertMsg)
	return nil
}

// BatchSendNotification 批量发送通知
func (s *MessageService) BatchSendNotification(msgs []*NotificationMessage) error {
	messages := make([]*models.Message, 0, len(msgs))
	expireAt := time.Now().Add(72 * time.Hour)

	for _, msg := range msgs {
		message := &models.Message{
			AgentID:     msg.AgentID,
			MessageType: msg.MessageType,
			Title:       msg.Title,
			Content:     msg.Content,
			IsRead:      false,
			IsPushed:    false,
			RelatedID:   msg.RelatedID,
			RelatedType: msg.RelatedType,
			ExpireAt:    &expireAt,
			CreatedAt:   time.Now(),
		}
		messages = append(messages, message)
	}

	return s.messageRepo.BatchCreate(messages)
}

// GetUnreadCount 获取未读消息数
func (s *MessageService) GetUnreadCount(agentID int64) (int, error) {
	messages, err := s.messageRepo.FindUnreadByAgentID(agentID)
	if err != nil {
		return 0, err
	}
	return len(messages), nil
}

// GetMessages 获取消息列表
func (s *MessageService) GetMessages(agentID int64, page, pageSize int) ([]*models.Message, error) {
	offset := (page - 1) * pageSize
	return s.messageRepo.FindByAgentID(agentID, pageSize, offset)
}

// GetMessagesByTypes 按类型获取消息列表
func (s *MessageService) GetMessagesByTypes(agentID int64, types []int16, page, pageSize int) ([]*models.Message, int64, error) {
	offset := (page - 1) * pageSize
	messages, err := s.messageRepo.FindByAgentIDAndTypes(agentID, types, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.messageRepo.CountByAgentIDAndTypes(agentID, types)
	if err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}

// GetMessagesByCategory 按分类获取消息列表
func (s *MessageService) GetMessagesByCategory(agentID int64, category string, page, pageSize int) ([]*models.Message, int64, error) {
	types := models.GetMessageTypesByCategory(category)
	return s.GetMessagesByTypes(agentID, types, page, pageSize)
}

// GetMessageStats 获取消息统计
func (s *MessageService) GetMessageStats(agentID int64) (*repository.MessageStats, error) {
	return s.messageRepo.GetStatsByAgentID(agentID)
}

// MarkAsRead 标记消息已读
func (s *MessageService) MarkAsRead(messageID int64) error {
	return s.messageRepo.MarkAsRead(messageID)
}

// MarkAllAsRead 标记所有消息已读
func (s *MessageService) MarkAllAsRead(agentID int64) error {
	return s.messageRepo.MarkAllAsRead(agentID)
}

// CleanupExpiredMessages 清理过期消息
func (s *MessageService) CleanupExpiredMessages() (int64, error) {
	return s.messageRepo.DeleteExpired()
}

// ============================================================
// 管理端方法
// ============================================================

// AdminSendMessage 管理员发送消息
func (s *MessageService) AdminSendMessage(req *SendMessageRequest, agentIDs []int64) error {
	// 验证有效期（1-30天）
	expireDays := req.ExpireDays
	if expireDays <= 0 || expireDays > 30 {
		expireDays = 3 // 默认3天
	}
	expireAt := time.Now().Add(time.Duration(expireDays) * 24 * time.Hour)

	// 创建消息列表
	messages := make([]*models.Message, 0, len(agentIDs))
	for _, agentID := range agentIDs {
		message := &models.Message{
			AgentID:     agentID,
			MessageType: req.MessageType,
			Title:       req.Title,
			Content:     req.Content,
			IsRead:      false,
			IsPushed:    false,
			ExpireAt:    &expireAt,
			CreatedAt:   time.Now(),
		}
		messages = append(messages, message)
	}

	// 批量创建
	if err := s.messageRepo.BatchCreate(messages); err != nil {
		return fmt.Errorf("batch create messages failed: %w", err)
	}

	log.Printf("[MessageService] Admin sent message to %d agents: %s", len(agentIDs), req.Title)
	return nil
}

// AdminGetAllMessages 管理员获取所有消息
func (s *MessageService) AdminGetAllMessages(page, pageSize int) ([]*models.Message, int64, error) {
	offset := (page - 1) * pageSize
	messages, err := s.messageRepo.FindAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.messageRepo.CountAll()
	if err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}

// AdminGetMessageByID 管理员获取消息详情
func (s *MessageService) AdminGetMessageByID(id int64) (*models.Message, error) {
	return s.messageRepo.FindByID(id)
}

// AdminDeleteMessage 管理员删除消息
func (s *MessageService) AdminDeleteMessage(id int64) error {
	return s.messageRepo.Delete(id)
}
