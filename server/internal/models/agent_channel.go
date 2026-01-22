package models

import "time"

// AgentChannel 代理商通道配置
// 记录代理商可以使用哪些通道，APP端只显示已开通的通道
type AgentChannel struct {
	ID         int64      `json:"id" gorm:"primaryKey"`
	AgentID    int64      `json:"agent_id" gorm:"not null;index"`      // 代理商ID
	ChannelID  int64      `json:"channel_id" gorm:"not null;index"`    // 通道ID
	IsEnabled  bool       `json:"is_enabled" gorm:"default:true"`      // 是否启用
	IsVisible  bool       `json:"is_visible" gorm:"default:true"`      // 对代理商是否可见
	EnabledAt  *time.Time `json:"enabled_at"`                          // 启用时间
	DisabledAt *time.Time `json:"disabled_at"`                         // 禁用时间
	EnabledBy  *int64     `json:"enabled_by"`                          // 启用人（用户ID）
	DisabledBy *int64     `json:"disabled_by"`                         // 禁用人（用户ID）
	Remark     string     `json:"remark" gorm:"type:text"`             // 备注
	CreatedAt  time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:now()"`

	// 关联字段（非数据库字段）
	Channel *Channel `json:"channel,omitempty" gorm:"-"` // 通道信息
}

func (AgentChannel) TableName() string {
	return "agent_channels"
}

// AgentChannelWithInfo 带通道信息的代理商通道配置
type AgentChannelWithInfo struct {
	AgentChannel
	ChannelCode string `json:"channel_code"` // 通道编码
	ChannelName string `json:"channel_name"` // 通道名称
}
