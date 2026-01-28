package models

import "time"

// TerminalType 终端类型
type TerminalType struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ChannelID   int64     `json:"channel_id" gorm:"not null;index"`
	ChannelCode string    `json:"channel_code" gorm:"type:varchar(50);not null;index"`
	BrandCode   string    `json:"brand_code" gorm:"type:varchar(50);not null"`
	BrandName   string    `json:"brand_name" gorm:"type:varchar(100);not null"`
	ModelCode   string    `json:"model_code" gorm:"type:varchar(50);not null"`
	ModelName   string    `json:"model_name" gorm:"type:varchar(100)"`
	Description string    `json:"description" gorm:"type:text"`
	Status      int16     `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Channel *Channel `json:"channel,omitempty" gorm:"foreignKey:ChannelID"`
}

// TableName 表名
func (TerminalType) TableName() string {
	return "terminal_types"
}

// TerminalTypeStatus 终端类型状态
const (
	TerminalTypeStatusDisabled = 0 // 禁用
	TerminalTypeStatusEnabled  = 1 // 启用
)

// IsEnabled 是否启用
func (t *TerminalType) IsEnabled() bool {
	return t.Status == TerminalTypeStatusEnabled
}

// FullName 完整名称（品牌 - 型号）
func (t *TerminalType) FullName() string {
	if t.ModelName != "" {
		return t.BrandName + " - " + t.ModelName
	}
	return t.BrandName + " - " + t.ModelCode
}
