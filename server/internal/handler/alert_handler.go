package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
)

// AlertHandler 告警配置处理器
type AlertHandler struct {
	configRepo   repository.AlertConfigRepository
	logRepo      repository.AlertLogRepository
	alertService *service.AlertService
}

// NewAlertHandler 创建告警配置处理器
func NewAlertHandler(
	configRepo repository.AlertConfigRepository,
	logRepo repository.AlertLogRepository,
	alertService *service.AlertService,
) *AlertHandler {
	return &AlertHandler{
		configRepo:   configRepo,
		logRepo:      logRepo,
		alertService: alertService,
	}
}

// RegisterRoutes 注册路由
func (h *AlertHandler) RegisterRoutes(rg *gin.RouterGroup) {
	configGroup := rg.Group("/alert-configs")
	{
		configGroup.GET("", h.ListConfigs)
		configGroup.POST("", h.CreateConfig)
		configGroup.GET("/:id", h.GetConfig)
		configGroup.PUT("/:id", h.UpdateConfig)
		configGroup.DELETE("/:id", h.DeleteConfig)
		configGroup.PUT("/:id/enable", h.EnableConfig)
		configGroup.POST("/:id/test", h.TestConfig)
	}

	logGroup := rg.Group("/alert-logs")
	{
		logGroup.GET("", h.ListAlertLogs)
		logGroup.GET("/:id", h.GetAlertLog)
	}
}

// AlertConfigResponse 告警配置响应
type AlertConfigResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	ChannelType     int16     `json:"channel_type"`
	ChannelTypeName string    `json:"channel_type_name"`
	WebhookURL      string    `json:"webhook_url"`
	EmailAddresses  string    `json:"email_addresses"`
	IsEnabled       bool      `json:"is_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListConfigs 告警配置列表
// @Summary 获取告警配置列表
// @Tags 告警配置
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs [get]
func (h *AlertHandler) ListConfigs(c *gin.Context) {
	configs, err := h.configRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	var list []AlertConfigResponse
	for _, config := range configs {
		list = append(list, AlertConfigResponse{
			ID:              config.ID,
			Name:            config.Name,
			ChannelType:     config.ChannelType,
			ChannelTypeName: models.GetAlertChannelName(config.ChannelType),
			WebhookURL:      h.maskWebhookURL(config.WebhookURL),
			EmailAddresses:  config.EmailAddresses,
			IsEnabled:       config.IsEnabled,
			CreatedAt:       config.CreatedAt,
			UpdatedAt:       config.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
	})
}

// CreateAlertConfigRequest 创建告警配置请求
type CreateAlertConfigRequest struct {
	Name           string `json:"name" binding:"required"`
	ChannelType    int16  `json:"channel_type" binding:"required,oneof=1 2 3"`
	WebhookURL     string `json:"webhook_url"`
	WebhookSecret  string `json:"webhook_secret"`
	EmailAddresses string `json:"email_addresses"`
	EmailSMTPHost  string `json:"email_smtp_host"`
	EmailSMTPPort  int    `json:"email_smtp_port"`
	EmailUsername  string `json:"email_username"`
	EmailPassword  string `json:"email_password"`
}

// CreateConfig 创建告警配置
// @Summary 创建告警配置
// @Tags 告警配置
// @Accept json
// @Param body body CreateAlertConfigRequest true "配置信息"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs [post]
func (h *AlertHandler) CreateConfig(c *gin.Context) {
	var req CreateAlertConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 验证通道类型与配置
	if req.ChannelType == models.AlertChannelDingTalk || req.ChannelType == models.AlertChannelWeChatWork {
		if req.WebhookURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook地址不能为空"})
			return
		}
	}
	if req.ChannelType == models.AlertChannelEmail {
		if req.EmailAddresses == "" || req.EmailSMTPHost == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件配置不完整"})
			return
		}
	}

	// 获取当前用户ID（从上下文获取）
	var createdBy *int64
	if userID, exists := c.Get("user_id"); exists {
		id := userID.(int64)
		createdBy = &id
	}

	config := &models.AlertConfig{
		Name:           req.Name,
		ChannelType:    req.ChannelType,
		WebhookURL:     req.WebhookURL,
		WebhookSecret:  req.WebhookSecret,
		EmailAddresses: req.EmailAddresses,
		EmailSMTPHost:  req.EmailSMTPHost,
		EmailSMTPPort:  req.EmailSMTPPort,
		EmailUsername:  req.EmailUsername,
		EmailPassword:  req.EmailPassword,
		IsEnabled:      true,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if config.EmailSMTPPort == 0 {
		config.EmailSMTPPort = 465
	}

	if err := h.configRepo.Create(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    gin.H{"id": config.ID},
	})
}

// GetConfig 获取告警配置详情
// @Summary 获取告警配置详情
// @Tags 告警配置
// @Param id path int true "配置ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs/{id} [get]
func (h *AlertHandler) GetConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	config, err := h.configRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	// 不返回密码
	config.EmailPassword = ""

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config,
	})
}

// UpdateAlertConfigRequest 更新告警配置请求
type UpdateAlertConfigRequest struct {
	Name           string `json:"name"`
	WebhookURL     string `json:"webhook_url"`
	WebhookSecret  string `json:"webhook_secret"`
	EmailAddresses string `json:"email_addresses"`
	EmailSMTPHost  string `json:"email_smtp_host"`
	EmailSMTPPort  int    `json:"email_smtp_port"`
	EmailUsername  string `json:"email_username"`
	EmailPassword  string `json:"email_password"`
}

// UpdateConfig 更新告警配置
// @Summary 更新告警配置
// @Tags 告警配置
// @Param id path int true "配置ID"
// @Param body body UpdateAlertConfigRequest true "配置信息"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs/{id} [put]
func (h *AlertHandler) UpdateConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	config, err := h.configRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	var req UpdateAlertConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 更新配置
	if req.Name != "" {
		config.Name = req.Name
	}
	if req.WebhookURL != "" {
		config.WebhookURL = req.WebhookURL
	}
	if req.WebhookSecret != "" {
		config.WebhookSecret = req.WebhookSecret
	}
	if req.EmailAddresses != "" {
		config.EmailAddresses = req.EmailAddresses
	}
	if req.EmailSMTPHost != "" {
		config.EmailSMTPHost = req.EmailSMTPHost
	}
	if req.EmailSMTPPort > 0 {
		config.EmailSMTPPort = req.EmailSMTPPort
	}
	if req.EmailUsername != "" {
		config.EmailUsername = req.EmailUsername
	}
	if req.EmailPassword != "" {
		config.EmailPassword = req.EmailPassword
	}

	if err := h.configRepo.Update(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteConfig 删除告警配置
// @Summary 删除告警配置
// @Tags 告警配置
// @Param id path int true "配置ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs/{id} [delete]
func (h *AlertHandler) DeleteConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.configRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// EnableConfigRequest 启用/禁用配置请求
type EnableConfigRequest struct {
	IsEnabled bool `json:"is_enabled"`
}

// EnableConfig 启用/禁用告警配置
// @Summary 启用/禁用告警配置
// @Tags 告警配置
// @Param id path int true "配置ID"
// @Param body body EnableConfigRequest true "启用状态"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs/{id}/enable [put]
func (h *AlertHandler) EnableConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req EnableConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.configRepo.UpdateEnabled(id, req.IsEnabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	status := "启用"
	if !req.IsEnabled {
		status = "禁用"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "配置已" + status,
	})
}

// TestConfig 测试告警配置
// @Summary 测试告警配置
// @Tags 告警配置
// @Param id path int true "配置ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-configs/{id}/test [post]
func (h *AlertHandler) TestConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.alertService.TestAlert(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "测试失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "测试消息已发送，请检查接收端",
	})
}

// ListAlertLogsRequest 告警日志列表请求
type ListAlertLogsRequest struct {
	JobName   string `form:"job_name"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// ListAlertLogs 告警日志列表
// @Summary 获取告警日志列表
// @Tags 告警配置
// @Param job_name query string false "任务名称"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-logs [get]
func (h *AlertHandler) ListAlertLogs(c *gin.Context) {
	var req ListAlertLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var logs []*models.AlertLog
	var total int64
	var err error

	if req.StartDate != "" && req.EndDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = endDate.Add(24 * time.Hour)

		logs, err = h.logRepo.FindByDateRange(startDate, endDate, req.PageSize, offset)
		if err == nil {
			total, _ = h.logRepo.CountByDateRange(startDate, endDate)
		}
	} else {
		logs, err = h.logRepo.FindByJobName(req.JobName, req.PageSize, offset)
		if err == nil {
			total, _ = h.logRepo.CountByJobName(req.JobName)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetAlertLog 告警日志详情
// @Summary 获取告警日志详情
// @Tags 告警配置
// @Param id path int true "日志ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/alert-logs/{id} [get]
func (h *AlertHandler) GetAlertLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	log, err := h.logRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": log,
	})
}

// maskWebhookURL 隐藏Webhook URL中的敏感信息
func (h *AlertHandler) maskWebhookURL(url string) string {
	if len(url) <= 30 {
		return url
	}
	return url[:30] + "***"
}
