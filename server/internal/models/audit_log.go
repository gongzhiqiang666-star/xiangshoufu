// Package models 审计日志模型
// 记录敏感操作，满足三级等保审计要求
package models

import "time"

// AuditLogType 审计日志类型
type AuditLogType int16

const (
	AuditLogTypeLogin        AuditLogType = 1  // 登录
	AuditLogTypeLogout       AuditLogType = 2  // 登出
	AuditLogTypePasswordChange AuditLogType = 3 // 密码修改
	AuditLogTypeDataExport   AuditLogType = 4  // 数据导出
	AuditLogTypeDataDelete   AuditLogType = 5  // 数据删除
	AuditLogTypeConfigChange AuditLogType = 6  // 配置变更
	AuditLogTypeWithdraw     AuditLogType = 7  // 提现操作
	AuditLogTypeTransfer     AuditLogType = 8  // 转账操作
	AuditLogTypeAgentCreate  AuditLogType = 9  // 创建代理商
	AuditLogTypeAgentDisable AuditLogType = 10 // 禁用代理商
	AuditLogTypePolicyChange AuditLogType = 11 // 政策变更
	AuditLogTypeRateChange   AuditLogType = 12 // 费率变更
	AuditLogTypeTerminalOp   AuditLogType = 13 // 终端操作
	AuditLogTypeDeduction    AuditLogType = 14 // 代扣操作
	AuditLogTypeReward       AuditLogType = 15 // 奖励发放
)

// AuditLogLevel 审计日志级别
type AuditLogLevel int16

const (
	AuditLogLevelInfo     AuditLogLevel = 1 // 信息
	AuditLogLevelWarning  AuditLogLevel = 2 // 警告
	AuditLogLevelCritical AuditLogLevel = 3 // 严重
)

// AuditLog 审计日志
type AuditLog struct {
	ID           int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	LogType      AuditLogType  `gorm:"type:smallint;not null;index" json:"log_type"`
	LogLevel     AuditLogLevel `gorm:"type:smallint;not null;default:1" json:"log_level"`
	UserID       int64         `gorm:"index" json:"user_id"`
	Username     string        `gorm:"size:64" json:"username"`
	AgentID      int64         `gorm:"index" json:"agent_id"`
	AgentName    string        `gorm:"size:128" json:"agent_name"`
	TargetType   string        `gorm:"size:32" json:"target_type"`   // 操作目标类型：agent/merchant/terminal等
	TargetID     int64         `gorm:"index" json:"target_id"`       // 操作目标ID
	TargetName   string        `gorm:"size:128" json:"target_name"`  // 操作目标名称
	Action       string        `gorm:"size:64;not null" json:"action"` // 操作动作
	Description  string        `gorm:"size:512" json:"description"`  // 操作描述
	OldValue     string        `gorm:"type:text" json:"old_value"`   // 变更前的值（JSON）
	NewValue     string        `gorm:"type:text" json:"new_value"`   // 变更后的值（JSON）
	IP           string        `gorm:"size:64" json:"ip"`
	UserAgent    string        `gorm:"size:256" json:"user_agent"`
	RequestPath  string        `gorm:"size:256" json:"request_path"`
	RequestMethod string       `gorm:"size:16" json:"request_method"`
	Result       int16         `gorm:"type:smallint;default:1" json:"result"` // 1成功 2失败
	ErrorMsg     string        `gorm:"size:512" json:"error_msg"`
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// GetLogTypeName 获取日志类型名称
func GetLogTypeName(logType AuditLogType) string {
	names := map[AuditLogType]string{
		AuditLogTypeLogin:        "用户登录",
		AuditLogTypeLogout:       "用户登出",
		AuditLogTypePasswordChange: "密码修改",
		AuditLogTypeDataExport:   "数据导出",
		AuditLogTypeDataDelete:   "数据删除",
		AuditLogTypeConfigChange: "配置变更",
		AuditLogTypeWithdraw:     "提现操作",
		AuditLogTypeTransfer:     "转账操作",
		AuditLogTypeAgentCreate:  "创建代理商",
		AuditLogTypeAgentDisable: "禁用代理商",
		AuditLogTypePolicyChange: "政策变更",
		AuditLogTypeRateChange:   "费率变更",
		AuditLogTypeTerminalOp:   "终端操作",
		AuditLogTypeDeduction:    "代扣操作",
		AuditLogTypeReward:       "奖励发放",
	}
	if name, ok := names[logType]; ok {
		return name
	}
	return "未知操作"
}

// GetLogLevelName 获取日志级别名称
func GetLogLevelName(level AuditLogLevel) string {
	switch level {
	case AuditLogLevelInfo:
		return "信息"
	case AuditLogLevelWarning:
		return "警告"
	case AuditLogLevelCritical:
		return "严重"
	default:
		return "未知"
	}
}
