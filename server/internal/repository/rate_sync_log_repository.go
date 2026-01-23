package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// RateSyncLogRepository 费率同步日志仓库接口
type RateSyncLogRepository interface {
	Create(ctx context.Context, log *models.RateSyncLog) error
	Update(ctx context.Context, log *models.RateSyncLog) error
	GetByID(ctx context.Context, id int64) (*models.RateSyncLog, error)
	GetPendingRetries(ctx context.Context, limit int) ([]*models.RateSyncLog, error)
	GetByMerchantID(ctx context.Context, merchantID int64, page, pageSize int) ([]*models.RateSyncLog, int64, error)
}

// GormRateSyncLogRepository GORM实现
type GormRateSyncLogRepository struct {
	db *gorm.DB
}

// NewGormRateSyncLogRepository 创建仓库实例
func NewGormRateSyncLogRepository(db *gorm.DB) *GormRateSyncLogRepository {
	return &GormRateSyncLogRepository{db: db}
}

// Create 创建记录
func (r *GormRateSyncLogRepository) Create(ctx context.Context, log *models.RateSyncLog) error {
	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(log).Error
}

// Update 更新记录
func (r *GormRateSyncLogRepository) Update(ctx context.Context, log *models.RateSyncLog) error {
	log.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(log).Error
}

// GetByID 根据ID获取
func (r *GormRateSyncLogRepository) GetByID(ctx context.Context, id int64) (*models.RateSyncLog, error) {
	var log models.RateSyncLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetPendingRetries 获取待重试的记录
func (r *GormRateSyncLogRepository) GetPendingRetries(ctx context.Context, limit int) ([]*models.RateSyncLog, error) {
	var logs []*models.RateSyncLog
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("sync_status = ?", models.RateSyncStatusFailed).
		Where("retry_count < max_retries").
		Where("next_retry_at <= ?", now).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// GetByMerchantID 根据商户ID获取记录（分页）
func (r *GormRateSyncLogRepository) GetByMerchantID(ctx context.Context, merchantID int64, page, pageSize int) ([]*models.RateSyncLog, int64, error) {
	var logs []*models.RateSyncLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RateSyncLog{}).Where("merchant_id = ?", merchantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
