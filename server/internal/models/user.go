package models

import (
	"time"
)

// User 用户模型（用于登录认证）
type User struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	Username    string     `json:"username" gorm:"size:50;uniqueIndex"`
	Password    string     `json:"-" gorm:"size:100"` // 密码不返回给前端
	Salt        string     `json:"-" gorm:"size:32"`
	AgentID     int64      `json:"agent_id" gorm:"index"`
	RoleType    int16      `json:"role_type" gorm:"default:1"`    // 1普通用户 2管理员
	Status      int16      `json:"status" gorm:"default:1"`       // 1正常 2禁用
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip" gorm:"size:50"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// 用户角色类型
const (
	UserRoleTypeNormal int16 = 1 // 普通用户
	UserRoleTypeAdmin  int16 = 2 // 管理员
)

// 用户状态
const (
	UserStatusActive   int16 = 1 // 正常
	UserStatusDisabled int16 = 2 // 禁用
)

// UserSession 用户会话信息（存储在JWT中）
type UserSession struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	AgentID  int64  `json:"agent_id"`
	RoleType int16  `json:"role_type"`
}

// RefreshToken 刷新令牌
type RefreshToken struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index"`
	Token     string    `json:"token" gorm:"size:64;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// LoginLog 登录日志
type LoginLog struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index"`
	Username  string    `json:"username" gorm:"size:50"`
	LoginIP   string    `json:"login_ip" gorm:"size:50"`
	UserAgent string    `json:"user_agent" gorm:"size:255"`
	Status    int16     `json:"status"` // 1成功 2失败
	FailMsg   string    `json:"fail_msg" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// 登录状态
const (
	LoginStatusSuccess int16 = 1
	LoginStatusFailed  int16 = 2
)
