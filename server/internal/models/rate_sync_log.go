package models

import (
	"time"
)

// RateSyncStatus 费率同步状态
type RateSyncStatus int

const (
	RateSyncStatusPending    RateSyncStatus = 0 // 待同步
	RateSyncStatusSyncing    RateSyncStatus = 1 // 同步中
	RateSyncStatusSuccess    RateSyncStatus = 2 // 同步成功
	RateSyncStatusFailed     RateSyncStatus = 3 // 同步失败
)

// RateSyncLog 费率同步日志
type RateSyncLog struct {
	ID         int64  `gorm:"primaryKey"`
	MerchantID int64  `gorm:"not null"`
	MerchantNo string `gorm:"size:64;not null"`
	TerminalSN string `gorm:"size:64"`
	ChannelCode string `gorm:"size:32;not null"`
	AgentID    int64  `gorm:"not null"`

	// 原费率
	OldCreditRate   *float64 `gorm:"type:decimal(10,4)"`
	OldDebitRate    *float64 `gorm:"type:decimal(10,4)"`
	OldDebitCap     *int64
	OldWechatRate   *float64 `gorm:"type:decimal(10,4)"`
	OldAlipayRate   *float64 `gorm:"type:decimal(10,4)"`
	OldUnionpayRate *float64 `gorm:"type:decimal(10,4)"`

	// 新费率
	NewCreditRate   *float64 `gorm:"type:decimal(10,4)"`
	NewDebitRate    *float64 `gorm:"type:decimal(10,4)"`
	NewDebitCap     *int64
	NewWechatRate   *float64 `gorm:"type:decimal(10,4)"`
	NewAlipayRate   *float64 `gorm:"type:decimal(10,4)"`
	NewUnionpayRate *float64 `gorm:"type:decimal(10,4)"`

	// 同步状态
	SyncStatus     RateSyncStatus `gorm:"not null;default:0"`
	ChannelTradeNo string         `gorm:"size:128"`
	ErrorMessage   string         `gorm:"type:text"`
	RetryCount     int            `gorm:"not null;default:0"`
	MaxRetries     int            `gorm:"not null;default:3"`
	NextRetryAt    *time.Time

	// 时间戳
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	SyncedAt  *time.Time
}

// TableName 表名
func (RateSyncLog) TableName() string {
	return "rate_sync_logs"
}

// CanRetry 是否可以重试
func (r *RateSyncLog) CanRetry() bool {
	return r.SyncStatus == RateSyncStatusFailed && r.RetryCount < r.MaxRetries
}

// MarkSyncing 标记为同步中
func (r *RateSyncLog) MarkSyncing() {
	r.SyncStatus = RateSyncStatusSyncing
	r.UpdatedAt = time.Now()
}

// MarkSuccess 标记为同步成功
func (r *RateSyncLog) MarkSuccess(tradeNo string) {
	r.SyncStatus = RateSyncStatusSuccess
	r.ChannelTradeNo = tradeNo
	now := time.Now()
	r.SyncedAt = &now
	r.UpdatedAt = now
}

// MarkFailed 标记为同步失败
func (r *RateSyncLog) MarkFailed(errMsg string) {
	r.SyncStatus = RateSyncStatusFailed
	r.ErrorMessage = errMsg
	r.RetryCount++
	r.UpdatedAt = time.Now()

	// 计算下次重试时间（指数退避：1分钟、5分钟、15分钟）
	if r.RetryCount < r.MaxRetries {
		var delay time.Duration
		switch r.RetryCount {
		case 1:
			delay = 1 * time.Minute
		case 2:
			delay = 5 * time.Minute
		default:
			delay = 15 * time.Minute
		}
		nextRetry := time.Now().Add(delay)
		r.NextRetryAt = &nextRetry
	}
}
