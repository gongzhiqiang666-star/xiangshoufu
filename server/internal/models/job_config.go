package models

import (
	"time"
)

// 任务执行状态
const (
	JobStatusSuccess   int16 = 1 // 成功
	JobStatusFailed    int16 = 2 // 失败
	JobStatusRunning   int16 = 3 // 运行中
)

// 任务触发类型
const (
	JobTriggerTypeAuto   int16 = 1 // 自动触发
	JobTriggerTypeManual int16 = 2 // 手动触发
)

// GetJobStatusName 获取任务状态名称
func GetJobStatusName(status int16) string {
	switch status {
	case JobStatusSuccess:
		return "成功"
	case JobStatusFailed:
		return "失败"
	case JobStatusRunning:
		return "运行中"
	default:
		return "未知"
	}
}

// GetJobTriggerTypeName 获取触发类型名称
func GetJobTriggerTypeName(triggerType int16) string {
	switch triggerType {
	case JobTriggerTypeAuto:
		return "自动触发"
	case JobTriggerTypeManual:
		return "手动触发"
	default:
		return "未知"
	}
}

// JobConfig 定时任务配置
type JobConfig struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	JobName         string    `json:"job_name" gorm:"size:100;uniqueIndex"`   // 任务名称（唯一标识）
	JobDesc         string    `json:"job_desc" gorm:"size:255"`               // 任务描述
	CronExpr        string    `json:"cron_expr" gorm:"size:50"`               // Cron表达式（预留）
	IntervalSeconds int       `json:"interval_seconds" gorm:"default:300"`    // 执行间隔(秒)
	IsEnabled       bool      `json:"is_enabled" gorm:"default:true"`         // 是否启用
	MaxRetries      int       `json:"max_retries" gorm:"default:3"`           // 最大重试次数
	RetryInterval   int       `json:"retry_interval" gorm:"default:60"`       // 初始重试间隔(秒)
	AlertThreshold  int       `json:"alert_threshold" gorm:"default:3"`       // 连续失败N次告警
	TimeoutSeconds  int       `json:"timeout_seconds" gorm:"default:3600"`    // 任务超时时间(秒)
	CreatedAt       time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (JobConfig) TableName() string {
	return "job_configs"
}

// JobExecutionLog 任务执行日志
type JobExecutionLog struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	JobName        string     `json:"job_name" gorm:"size:100;index"`         // 任务名称
	StartedAt      time.Time  `json:"started_at" gorm:"index"`                // 开始时间
	EndedAt        *time.Time `json:"ended_at"`                               // 结束时间
	DurationMs     int64      `json:"duration_ms"`                            // 执行耗时(毫秒)
	Status         int16      `json:"status" gorm:"default:3;index"`          // 1成功 2失败 3运行中
	ProcessedCount int        `json:"processed_count" gorm:"default:0"`       // 处理条数
	SuccessCount   int        `json:"success_count" gorm:"default:0"`         // 成功条数
	FailCount      int        `json:"fail_count" gorm:"default:0"`            // 失败条数
	ErrorMessage   string     `json:"error_message" gorm:"type:text"`         // 错误信息
	ErrorStack     string     `json:"error_stack" gorm:"type:text"`           // 错误堆栈
	RetryCount     int        `json:"retry_count" gorm:"default:0"`           // 当前重试次数
	TriggerType    int16      `json:"trigger_type" gorm:"default:1"`          // 1自动触发 2手动触发
	CreatedAt      time.Time  `json:"created_at" gorm:"default:now();index"`
}

// TableName 表名
func (JobExecutionLog) TableName() string {
	return "job_execution_logs"
}

// GetStatusName 获取状态名称
func (l *JobExecutionLog) GetStatusName() string {
	return GetJobStatusName(l.Status)
}

// GetTriggerTypeName 获取触发类型名称
func (l *JobExecutionLog) GetTriggerTypeName() string {
	return GetJobTriggerTypeName(l.TriggerType)
}

// JobFailCounter 任务失败计数
type JobFailCounter struct {
	ID               int64      `json:"id" gorm:"primaryKey"`
	JobName          string     `json:"job_name" gorm:"size:100;uniqueIndex"` // 任务名称
	ConsecutiveFails int        `json:"consecutive_fails" gorm:"default:0"`   // 连续失败次数
	LastFailAt       *time.Time `json:"last_fail_at"`                         // 最后失败时间
	LastSuccessAt    *time.Time `json:"last_success_at"`                      // 最后成功时间
	LastAlertAt      *time.Time `json:"last_alert_at"`                        // 最后告警时间
	UpdatedAt        time.Time  `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (JobFailCounter) TableName() string {
	return "job_fail_counters"
}

// JobExecutionResult 任务执行结果（用于任务返回）
type JobExecutionResult struct {
	ProcessedCount int    // 处理条数
	SuccessCount   int    // 成功条数
	FailCount      int    // 失败条数
	ErrorMessage   string // 错误信息
}
