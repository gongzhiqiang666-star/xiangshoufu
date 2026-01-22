package models

import (
	"time"
)

// 告警通道类型
const (
	AlertChannelDingTalk   int16 = 1 // 钉钉
	AlertChannelWeChatWork int16 = 2 // 企业微信
	AlertChannelEmail      int16 = 3 // 邮件
)

// 告警类型
const (
	AlertTypeJobFailed       int16 = 1 // 任务失败
	AlertTypeConsecutiveFail int16 = 2 // 连续失败
	AlertTypeJobTimeout      int16 = 3 // 任务超时
)

// 告警发送状态
const (
	AlertSendStatusPending int16 = 0 // 待发送
	AlertSendStatusSent    int16 = 1 // 已发送
	AlertSendStatusFailed  int16 = 2 // 发送失败
)

// GetAlertChannelName 获取告警通道名称
func GetAlertChannelName(channelType int16) string {
	switch channelType {
	case AlertChannelDingTalk:
		return "钉钉"
	case AlertChannelWeChatWork:
		return "企业微信"
	case AlertChannelEmail:
		return "邮件"
	default:
		return "未知"
	}
}

// GetAlertTypeName 获取告警类型名称
func GetAlertTypeName(alertType int16) string {
	switch alertType {
	case AlertTypeJobFailed:
		return "任务失败"
	case AlertTypeConsecutiveFail:
		return "连续失败"
	case AlertTypeJobTimeout:
		return "任务超时"
	default:
		return "未知"
	}
}

// GetAlertSendStatusName 获取发送状态名称
func GetAlertSendStatusName(status int16) string {
	switch status {
	case AlertSendStatusPending:
		return "待发送"
	case AlertSendStatusSent:
		return "已发送"
	case AlertSendStatusFailed:
		return "发送失败"
	default:
		return "未知"
	}
}

// AlertConfig 告警配置
type AlertConfig struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"size:100"`                  // 配置名称
	ChannelType    int16     `json:"channel_type"`                          // 1钉钉 2企微 3邮件
	WebhookURL     string    `json:"webhook_url" gorm:"size:500"`           // Webhook地址
	WebhookSecret  string    `json:"webhook_secret" gorm:"size:200"`        // Webhook密钥（钉钉签名）
	EmailAddresses string    `json:"email_addresses" gorm:"type:text"`      // 邮箱地址(逗号分隔)
	EmailSMTPHost  string    `json:"email_smtp_host" gorm:"size:100"`       // SMTP服务器
	EmailSMTPPort  int       `json:"email_smtp_port" gorm:"default:465"`    // SMTP端口
	EmailUsername  string    `json:"email_username" gorm:"size:100"`        // SMTP用户名
	EmailPassword  string    `json:"-" gorm:"size:200"`                     // SMTP密码（不返回给前端）
	IsEnabled      bool      `json:"is_enabled" gorm:"default:true"`
	CreatedBy      *int64    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (AlertConfig) TableName() string {
	return "alert_configs"
}

// GetChannelTypeName 获取通道类型名称
func (c *AlertConfig) GetChannelTypeName() string {
	return GetAlertChannelName(c.ChannelType)
}

// AlertLog 告警记录
type AlertLog struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	JobName      string     `json:"job_name" gorm:"size:100;index"`       // 任务名称
	AlertType    int16      `json:"alert_type"`                           // 1任务失败 2连续失败 3任务超时
	ChannelType  int16      `json:"channel_type"`                         // 1钉钉 2企微 3邮件
	ConfigID     *int64     `json:"config_id"`                            // 关联的告警配置ID
	Title        string     `json:"title" gorm:"size:200"`                // 告警标题
	Message      string     `json:"message" gorm:"type:text"`             // 告警内容
	SendStatus   int16      `json:"send_status" gorm:"default:0;index"`   // 0待发送 1已发送 2发送失败
	SendAt       *time.Time `json:"send_at"`                              // 发送时间
	ErrorMessage string     `json:"error_message" gorm:"type:text"`       // 发送失败原因
	CreatedAt    time.Time  `json:"created_at" gorm:"default:now();index"`
}

// TableName 表名
func (AlertLog) TableName() string {
	return "alert_logs"
}

// GetAlertTypeName 获取告警类型名称
func (l *AlertLog) GetAlertTypeName() string {
	return GetAlertTypeName(l.AlertType)
}

// GetChannelTypeName 获取通道类型名称
func (l *AlertLog) GetChannelTypeName() string {
	return GetAlertChannelName(l.ChannelType)
}

// GetSendStatusName 获取发送状态名称
func (l *AlertLog) GetSendStatusName() string {
	return GetAlertSendStatusName(l.SendStatus)
}

// AlertRequest 告警请求
type AlertRequest struct {
	JobName      string // 任务名称
	AlertType    int16  // 告警类型
	Title        string // 告警标题
	Message      string // 告警内容
	ErrorMessage string // 错误信息
}
