package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// RawCallbackLog 类型别名，方便外部使用
type RawCallbackLog = models.RawCallbackLog

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

// FindArchivedLogs 查找需要归档的日志（超过指定天数）
func (r *GormRawCallbackRepository) FindArchivedLogs(retentionDays int, limit int) ([]*models.RawCallbackLog, error) {
	var logs []*models.RawCallbackLog
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	err := r.db.Where("received_at < ? AND process_status = ?", cutoffTime, models.ProcessStatusSuccess).
		Order("received_at ASC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// DeleteArchivedLogs 删除已归档的日志
func (r *GormRawCallbackRepository) DeleteArchivedLogs(ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.Delete(&models.RawCallbackLog{}, ids)
	return result.RowsAffected, result.Error
}

// CountArchivedLogs 统计需要归档的日志数量
func (r *GormRawCallbackRepository) CountArchivedLogs(retentionDays int) (int64, error) {
	var count int64
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	err := r.db.Model(&models.RawCallbackLog{}).
		Where("received_at < ? AND process_status = ?", cutoffTime, models.ProcessStatusSuccess).
		Count(&count).Error
	return count, err
}

// GetDB 获取数据库连接（用于分区管理等原生SQL操作）
func (r *GormRawCallbackRepository) GetDB() *gorm.DB {
	return r.db
}

// 确保实现了接口
var _ RawCallbackRepository = (*GormRawCallbackRepository)(nil)
