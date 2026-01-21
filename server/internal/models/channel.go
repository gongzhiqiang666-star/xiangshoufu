package models

import "time"

// Channel 支付通道
type Channel struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ChannelCode string    `json:"channel_code" gorm:"size:32;not null;uniqueIndex"` // 通道编码
	ChannelName string    `json:"channel_name" gorm:"size:64;not null"`             // 通道名称
	Description string    `json:"description" gorm:"type:text"`                     // 通道描述
	Status      int16     `json:"status" gorm:"default:1"`                          // 1:启用 0:禁用
	Priority    int       `json:"priority" gorm:"default:0"`                        // 优先级
	Config      string    `json:"config" gorm:"type:jsonb"`                         // 通道配置
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:now()"`
}

func (Channel) TableName() string {
	return "channels"
}
