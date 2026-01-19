package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// GormRawCallbackRepository GORM实现的原始回调仓库
type GormRawCallbackRepository struct {
	db *gorm.DB
}

// NewGormRawCallbackRepository 创建仓库
func NewGormRawCallbackRepository(db *gorm.DB) *GormRawCallbackRepository {
	return &GormRawCallbackRepository{db: db}
}

// Create 创建回调日志
func (r *GormRawCallbackRepository) Create(log *models.RawCallbackLog) error {
	return r.db.Create(log).Error
}

// Update 更新回调日志
func (r *GormRawCallbackRepository) Update(log *models.RawCallbackLog) error {
	return r.db.Save(log).Error
}

// FindByID 根据ID查找
func (r *GormRawCallbackRepository) FindByID(id int64) (*models.RawCallbackLog, error) {
	var log models.RawCallbackLog
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// FindByIdempotentKey 根据幂等键查找
func (r *GormRawCallbackRepository) FindByIdempotentKey(key string) (*models.RawCallbackLog, error) {
	var log models.RawCallbackLog
	err := r.db.Where("idempotent_key = ?", key).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// FindPendingLogs 查找待处理的日志
func (r *GormRawCallbackRepository) FindPendingLogs(limit int) ([]*models.RawCallbackLog, error) {
	var logs []*models.RawCallbackLog
	err := r.db.Where("process_status = ?", models.ProcessStatusPending).
		Order("received_at ASC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// FindFailedLogs 查找失败的日志（用于重试）
func (r *GormRawCallbackRepository) FindFailedLogs(maxRetry int, limit int) ([]*models.RawCallbackLog, error) {
	var logs []*models.RawCallbackLog
	err := r.db.Where("process_status = ? AND retry_count < ?", models.ProcessStatusFailed, maxRetry).
		Order("received_at ASC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// UpdateStatus 更新处理状态
func (r *GormRawCallbackRepository) UpdateStatus(id int64, status int16, errorMsg string) error {
	updates := map[string]interface{}{
		"process_status": status,
		"error_message":  errorMsg,
	}
	if status == models.ProcessStatusSuccess {
		now := time.Now()
		updates["processed_at"] = &now
	}
	return r.db.Model(&models.RawCallbackLog{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementRetryCount 增加重试次数
func (r *GormRawCallbackRepository) IncrementRetryCount(id int64) error {
	return r.db.Model(&models.RawCallbackLog{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

// 确保实现了接口
var _ RawCallbackRepository = (*GormRawCallbackRepository)(nil)
