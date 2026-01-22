package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// GormAuditLogRepository 审计日志仓储
type GormAuditLogRepository struct {
	db *gorm.DB
}

// NewGormAuditLogRepository 创建审计日志仓储
func NewGormAuditLogRepository(db *gorm.DB) *GormAuditLogRepository {
	return &GormAuditLogRepository{db: db}
}

// Create 创建审计日志
func (r *GormAuditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// AuditLogQueryParams 查询参数
type AuditLogQueryParams struct {
	LogType     *models.AuditLogType
	LogLevel    *models.AuditLogLevel
	UserID      *int64
	AgentID     *int64
	TargetType  string
	TargetID    *int64
	StartTime   *time.Time
	EndTime     *time.Time
	Keyword     string
	Limit       int
	Offset      int
}

// FindByParams 根据参数查询审计日志
func (r *GormAuditLogRepository) FindByParams(params AuditLogQueryParams) ([]*models.AuditLog, int64, error) {
	var logs []*models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{})

	if params.LogType != nil {
		query = query.Where("log_type = ?", *params.LogType)
	}
	if params.LogLevel != nil {
		query = query.Where("log_level = ?", *params.LogLevel)
	}
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if params.AgentID != nil {
		query = query.Where("agent_id = ?", *params.AgentID)
	}
	if params.TargetType != "" {
		query = query.Where("target_type = ?", params.TargetType)
	}
	if params.TargetID != nil {
		query = query.Where("target_id = ?", *params.TargetID)
	}
	if params.StartTime != nil {
		query = query.Where("created_at >= ?", *params.StartTime)
	}
	if params.EndTime != nil {
		query = query.Where("created_at <= ?", *params.EndTime)
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("description LIKE ? OR username LIKE ? OR agent_name LIKE ?", keyword, keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}

	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindByID 根据ID查询
func (r *GormAuditLogRepository) FindByID(id int64) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := r.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// CountByType 按类型统计
func (r *GormAuditLogRepository) CountByType(startTime, endTime time.Time) (map[models.AuditLogType]int64, error) {
	type result struct {
		LogType models.AuditLogType
		Count   int64
	}

	var results []result
	if err := r.db.Model(&models.AuditLog{}).
		Select("log_type, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("log_type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[models.AuditLogType]int64)
	for _, r := range results {
		countMap[r.LogType] = r.Count
	}

	return countMap, nil
}

// DeleteOldLogs 删除旧日志（保留指定天数）
func (r *GormAuditLogRepository) DeleteOldLogs(retentionDays int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	result := r.db.Where("created_at < ?", cutoff).Delete(&models.AuditLog{})
	return result.RowsAffected, result.Error
}

// GetDB 获取数据库连接
func (r *GormAuditLogRepository) GetDB() *gorm.DB {
	return r.db
}
