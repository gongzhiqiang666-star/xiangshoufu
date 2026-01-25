package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================
// 奖励模块 - 与通道解耦的独立奖励配置
// ============================================================

// TimeType 时间类型
type TimeType string

const (
	TimeTypeDays   TimeType = "days"   // 按天数
	TimeTypeMonths TimeType = "months" // 按自然月
)

// DimensionType 维度类型
type DimensionType string

const (
	DimensionTypeAmount DimensionType = "amount" // 按金额
	DimensionTypeCount  DimensionType = "count"  // 按笔数
)

// RewardProgressStatus 奖励进度状态
type RewardProgressStatus string

const (
	RewardProgressStatusActive     RewardProgressStatus = "active"     // 进行中
	RewardProgressStatusCompleted  RewardProgressStatus = "completed"  // 已完成
	RewardProgressStatusTerminated RewardProgressStatus = "terminated" // 已终止
)

// StageRewardStatus 阶段奖励状态
type StageRewardStatus string

const (
	StageRewardStatusPending    StageRewardStatus = "pending"     // 待检查
	StageRewardStatusAchieved   StageRewardStatus = "achieved"    // 已达标
	StageRewardStatusFailed     StageRewardStatus = "failed"      // 未达标
	StageRewardStatusGapBlocked StageRewardStatus = "gap_blocked" // 被断档阻断
	StageRewardStatusSettled    StageRewardStatus = "settled"     // 已结算
)

// ============================================================
// 奖励政策模版
// ============================================================

// RewardPolicyTemplate 奖励政策模版
type RewardPolicyTemplate struct {
	ID            int64         `json:"id" gorm:"primaryKey"`
	Name          string        `json:"name" gorm:"size:100;not null"`           // 模版名称
	TimeType      TimeType      `json:"time_type" gorm:"size:20;not null"`       // 时间类型
	DimensionType DimensionType `json:"dimension_type" gorm:"size:20;not null"`  // 维度类型
	TransTypes    string        `json:"trans_types" gorm:"size:100"`             // 交易类型，逗号分隔
	AmountMin     int64         `json:"amount_min" gorm:"default:0"`             // 交易金额下限（分）
	AmountMax     *int64        `json:"amount_max"`                              // 交易金额上限（分），NULL表示无上限
	AllowGap      bool          `json:"allow_gap" gorm:"default:false"`          // 断档开关
	Enabled       bool          `json:"enabled" gorm:"default:true"`             // 是否启用
	Description   string        `json:"description" gorm:"type:text"`            // 描述说明
	CreatedAt     time.Time     `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"default:now()"`

	// 关联数据（不存储在数据库）
	Stages []*RewardStage `json:"stages,omitempty" gorm:"-"`
}

// TableName 表名
func (RewardPolicyTemplate) TableName() string {
	return "reward_policy_templates"
}

// ============================================================
// 奖励阶段配置
// ============================================================

// RewardStage 奖励阶段配置
type RewardStage struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	TemplateID   int64     `json:"template_id" gorm:"not null;index"`      // 模版ID
	StageOrder   int       `json:"stage_order" gorm:"not null"`            // 阶段顺序（从1开始）
	StartValue   int       `json:"start_value" gorm:"not null"`            // 开始值（天数或月份）
	EndValue     int       `json:"end_value" gorm:"not null"`              // 结束值（天数或月份）
	TargetValue  int64     `json:"target_value" gorm:"not null"`           // 达标值（金额分或笔数）
	RewardAmount int64     `json:"reward_amount" gorm:"not null"`          // 奖励金额（分）
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (RewardStage) TableName() string {
	return "reward_stages"
}

// ============================================================
// 代理商奖励比例配置
// ============================================================

// AgentRewardRate 代理商奖励比例配置
type AgentRewardRate struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	AgentID    int64     `json:"agent_id" gorm:"not null;uniqueIndex"`     // 代理商ID
	RewardRate float64   `json:"reward_rate" gorm:"type:decimal(5,4)"`     // 奖励比例（0.10表示10%）
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:now()"`
}

// TableName 表名
func (AgentRewardRate) TableName() string {
	return "agent_reward_rates"
}

// ============================================================
// 终端奖励进度
// ============================================================

// TemplateSnapshot 模版快照（用于JSONB存储）
type TemplateSnapshot struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	TimeType      TimeType       `json:"time_type"`
	DimensionType DimensionType  `json:"dimension_type"`
	TransTypes    string         `json:"trans_types"`
	AmountMin     int64          `json:"amount_min"`
	AmountMax     *int64         `json:"amount_max"`
	AllowGap      bool           `json:"allow_gap"`
	Stages        []*RewardStage `json:"stages"`
}

// Scan 实现sql.Scanner接口
func (t *TemplateSnapshot) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into TemplateSnapshot", value)
	}

	return json.Unmarshal(bytes, t)
}

// Value 实现driver.Valuer接口
func (t TemplateSnapshot) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// TerminalRewardProgress 终端奖励进度
type TerminalRewardProgress struct {
	ID                int64                `json:"id" gorm:"primaryKey"`
	TerminalSN        string               `json:"terminal_sn" gorm:"size:50;not null;index"`      // 终端SN
	TerminalID        *int64               `json:"terminal_id"`                                    // 终端ID
	TemplateID        int64                `json:"template_id" gorm:"not null;index"`              // 模版ID
	TemplateSnapshot  TemplateSnapshot     `json:"template_snapshot" gorm:"type:jsonb;not null"`   // 政策快照
	BindAgentID       int64                `json:"bind_agent_id" gorm:"not null;index"`            // 绑定时的代理商ID
	BindTime          time.Time            `json:"bind_time" gorm:"not null;index"`                // 绑定时间
	CurrentStage      int                  `json:"current_stage" gorm:"default:1"`                 // 当前阶段
	LastAchievedStage int                  `json:"last_achieved_stage" gorm:"default:0"`           // 最后达标阶段
	Status            RewardProgressStatus `json:"status" gorm:"size:20;default:'active';index"`   // 状态
	CompletedAt       *time.Time           `json:"completed_at"`                                   // 完成时间
	TerminatedAt      *time.Time           `json:"terminated_at"`                                  // 终止时间
	CreatedAt         time.Time            `json:"created_at" gorm:"default:now()"`
	UpdatedAt         time.Time            `json:"updated_at" gorm:"default:now()"`

	// 关联数据
	StageRewards []*TerminalStageReward `json:"stage_rewards,omitempty" gorm:"-"`
}

// TableName 表名
func (TerminalRewardProgress) TableName() string {
	return "terminal_reward_progress"
}

// ============================================================
// 终端阶段奖励记录
// ============================================================

// TerminalStageReward 终端阶段奖励记录
type TerminalStageReward struct {
	ID           int64             `json:"id" gorm:"primaryKey"`
	ProgressID   int64             `json:"progress_id" gorm:"not null;index"`        // 进度ID
	TerminalSN   string            `json:"terminal_sn" gorm:"size:50;not null;index"` // 终端SN
	StageOrder   int               `json:"stage_order" gorm:"not null"`              // 阶段顺序
	StageStart   time.Time         `json:"stage_start" gorm:"not null"`              // 阶段开始时间
	StageEnd     time.Time         `json:"stage_end" gorm:"not null;index"`          // 阶段结束时间
	TargetValue  int64             `json:"target_value" gorm:"not null"`             // 目标值
	ActualValue  int64             `json:"actual_value" gorm:"default:0"`            // 实际值
	IsAchieved   bool              `json:"is_achieved" gorm:"default:false"`         // 是否达标
	RewardAmount *int64            `json:"reward_amount"`                            // 应发奖励金额（分）
	Status       StageRewardStatus `json:"status" gorm:"size:20;not null;index"`     // 状态
	GapBlocked   bool              `json:"gap_blocked" gorm:"default:false"`         // 是否被断档阻断
	SettledAt    *time.Time        `json:"settled_at"`                               // 结算时间
	CreatedAt    time.Time         `json:"created_at" gorm:"default:now()"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"default:now()"`

	// 关联数据
	Distributions []*RewardDistribution `json:"distributions,omitempty" gorm:"-"`
}

// TableName 表名
func (TerminalStageReward) TableName() string {
	return "terminal_stage_rewards"
}

// ============================================================
// 奖励发放记录（级差分配）
// ============================================================

// RewardDistribution 奖励发放记录
type RewardDistribution struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	StageRewardID  int64     `json:"stage_reward_id" gorm:"not null;index"`      // 阶段奖励ID
	TerminalSN     string    `json:"terminal_sn" gorm:"size:50;not null;index"`  // 终端SN
	AgentID        int64     `json:"agent_id" gorm:"not null;index"`             // 代理商ID
	AgentLevel     int       `json:"agent_level" gorm:"not null"`                // 层级
	RewardRate     float64   `json:"reward_rate" gorm:"type:decimal(5,4)"`       // 奖励比例
	RewardAmount   int64     `json:"reward_amount" gorm:"not null"`              // 奖励金额（分）
	WalletRecordID *int64    `json:"wallet_record_id"`                           // 钱包记录ID
	WalletStatus   int16     `json:"wallet_status" gorm:"default:0"`             // 0-待入账 1-已入账
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 表名
func (RewardDistribution) TableName() string {
	return "reward_distributions"
}

// ============================================================
// 奖励池溢出异常日志
// ============================================================

// AgentChainInfo 代理商链信息（用于JSONB存储）
type AgentChainInfo struct {
	AgentID    int64   `json:"agent_id"`
	AgentName  string  `json:"agent_name"`
	Level      int     `json:"level"`
	RewardRate float64 `json:"reward_rate"`
}

// AgentChain 代理商链
type AgentChain []AgentChainInfo

// Scan 实现sql.Scanner接口
func (a *AgentChain) Scan(value interface{}) error {
	if value == nil {
		*a = make(AgentChain, 0)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into AgentChain", value)
	}

	return json.Unmarshal(bytes, a)
}

// Value 实现driver.Valuer接口
func (a AgentChain) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

// RewardOverflowLog 奖励池溢出异常日志
type RewardOverflowLog struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	TerminalSN    string     `json:"terminal_sn" gorm:"size:50;not null;index"`   // 终端SN
	StageRewardID *int64     `json:"stage_reward_id"`                             // 阶段奖励ID
	AgentChain    AgentChain `json:"agent_chain" gorm:"type:jsonb;not null"`      // 代理商链
	TotalRate     float64    `json:"total_rate" gorm:"type:decimal(5,4)"`         // 总比例
	RewardAmount  int64      `json:"reward_amount" gorm:"not null"`               // 原应发奖励金额
	ErrorMessage  string     `json:"error_message" gorm:"type:text"`              // 错误信息
	Resolved      bool       `json:"resolved" gorm:"default:false;index"`         // 是否已解决
	ResolvedAt    *time.Time `json:"resolved_at"`                                 // 解决时间
	ResolvedBy    string     `json:"resolved_by" gorm:"size:50"`                  // 解决人
	CreatedAt     time.Time  `json:"created_at" gorm:"default:now();index"`
}

// TableName 表名
func (RewardOverflowLog) TableName() string {
	return "reward_overflow_logs"
}

// ============================================================
// 请求/响应DTO
// ============================================================

// CreateRewardTemplateRequest 创建奖励模版请求
type CreateRewardTemplateRequest struct {
	Name          string                `json:"name" binding:"required"`
	TimeType      TimeType              `json:"time_type" binding:"required,oneof=days months"`
	DimensionType DimensionType         `json:"dimension_type" binding:"required,oneof=amount count"`
	TransTypes    string                `json:"trans_types"`
	AmountMin     int64                 `json:"amount_min"`
	AmountMax     *int64                `json:"amount_max"`
	AllowGap      bool                  `json:"allow_gap"`
	Description   string                `json:"description"`
	Stages        []CreateStageRequest  `json:"stages" binding:"required,min=1,dive"`
}

// CreateStageRequest 创建阶段请求
type CreateStageRequest struct {
	StageOrder   int   `json:"stage_order" binding:"required,min=1"`
	StartValue   int   `json:"start_value" binding:"required,min=1"`
	EndValue     int   `json:"end_value" binding:"required,min=1"`
	TargetValue  int64 `json:"target_value" binding:"required,min=1"`
	RewardAmount int64 `json:"reward_amount" binding:"required,min=1"`
}

// UpdateRewardTemplateRequest 更新奖励模版请求
type UpdateRewardTemplateRequest struct {
	Name          string                `json:"name"`
	TransTypes    string                `json:"trans_types"`
	AmountMin     int64                 `json:"amount_min"`
	AmountMax     *int64                `json:"amount_max"`
	AllowGap      *bool                 `json:"allow_gap"`
	Enabled       *bool                 `json:"enabled"`
	Description   string                `json:"description"`
	Stages        []CreateStageRequest  `json:"stages"`
}

// RewardTemplateListItem 奖励模版列表项
type RewardTemplateListItem struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	TimeType      TimeType      `json:"time_type"`
	DimensionType DimensionType `json:"dimension_type"`
	TransTypes    string        `json:"trans_types"`
	AllowGap      bool          `json:"allow_gap"`
	Enabled       bool          `json:"enabled"`
	StageCount    int           `json:"stage_count"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// RewardTemplateDetail 奖励模版详情
type RewardTemplateDetail struct {
	RewardPolicyTemplate
	Stages []*RewardStage `json:"stages"`
}

// AgentRewardRateRequest 代理商奖励比例请求
type AgentRewardRateRequest struct {
	AgentID    int64   `json:"agent_id" binding:"required"`
	RewardRate float64 `json:"reward_rate" binding:"required,min=0,max=1"`
}

// TerminalRewardProgressDetail 终端奖励进度详情
type TerminalRewardProgressDetail struct {
	TerminalRewardProgress
	TemplateName string                 `json:"template_name"`
	AgentName    string                 `json:"agent_name"`
	StageRewards []*TerminalStageReward `json:"stage_rewards"`
}
