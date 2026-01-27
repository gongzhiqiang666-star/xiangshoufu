package models

import "time"

// Agent 代理商模型（用于测试和业务逻辑）
type Agent struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100"`
	ParentID  int64     `json:"parent_id" gorm:"index"`
	Level     int       `json:"level" gorm:"default:1"`
	Status    int16     `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (Agent) TableName() string {
	return "agents"
}
