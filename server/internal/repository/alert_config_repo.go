package repository

import (
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// AlertConfigRepository 告警配置仓库接口
type AlertConfigRepository interface {
	FindAll() ([]*models.AlertConfig, error)
	FindByID(id int64) (*models.AlertConfig, error)
	FindEnabled() ([]*models.AlertConfig, error)
	FindByChannelType(channelType int16) ([]*models.AlertConfig, error)
	Create(config *models.AlertConfig) error
	Update(config *models.AlertConfig) error
	Delete(id int64) error
	UpdateEnabled(id int64, isEnabled bool) error
}

// GormAlertConfigRepository 告警配置仓库GORM实现
type GormAlertConfigRepository struct {
	db *gorm.DB
}

// NewGormAlertConfigRepository 创建告警配置仓库
func NewGormAlertConfigRepository(db *gorm.DB) *GormAlertConfigRepository {
	return &GormAlertConfigRepository{db: db}
}

// FindAll 查询所有告警配置
func (r *GormAlertConfigRepository) FindAll() ([]*models.AlertConfig, error) {
	var configs []*models.AlertConfig
	err := r.db.Order("id ASC").Find(&configs).Error
	return configs, err
}

// FindByID 根据ID查询配置
func (r *GormAlertConfigRepository) FindByID(id int64) (*models.AlertConfig, error) {
	var config models.AlertConfig
	err := r.db.Where("id = ?", id).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FindEnabled 查询所有启用的配置
func (r *GormAlertConfigRepository) FindEnabled() ([]*models.AlertConfig, error) {
	var configs []*models.AlertConfig
	err := r.db.Where("is_enabled = ?", true).Order("id ASC").Find(&configs).Error
	return configs, err
}

// FindByChannelType 根据通道类型查询配置
func (r *GormAlertConfigRepository) FindByChannelType(channelType int16) ([]*models.AlertConfig, error) {
	var configs []*models.AlertConfig
	err := r.db.Where("channel_type = ? AND is_enabled = ?", channelType, true).
		Order("id ASC").Find(&configs).Error
	return configs, err
}

// Create 创建告警配置
func (r *GormAlertConfigRepository) Create(config *models.AlertConfig) error {
	return r.db.Create(config).Error
}

// Update 更新告警配置
func (r *GormAlertConfigRepository) Update(config *models.AlertConfig) error {
	config.UpdatedAt = time.Now()
	return r.db.Save(config).Error
}

// Delete 删除告警配置
func (r *GormAlertConfigRepository) Delete(id int64) error {
	return r.db.Delete(&models.AlertConfig{}, id).Error
}

// UpdateEnabled 更新启用状态
func (r *GormAlertConfigRepository) UpdateEnabled(id int64, isEnabled bool) error {
	return r.db.Model(&models.AlertConfig{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_enabled": isEnabled,
			"updated_at": time.Now(),
		}).Error
}

// Ensure interface compliance
var _ AlertConfigRepository = (*GormAlertConfigRepository)(nil)

// AlertLogRepository 告警记录仓库接口
type AlertLogRepository interface {
	Create(log *models.AlertLog) error
	Update(log *models.AlertLog) error
	FindByID(id int64) (*models.AlertLog, error)
	FindPending(limit int) ([]*models.AlertLog, error)
	FindByJobName(jobName string, limit, offset int) ([]*models.AlertLog, error)
	FindByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.AlertLog, error)
	CountByJobName(jobName string) (int64, error)
	CountByDateRange(startDate, endDate time.Time) (int64, error)
	UpdateSendStatus(id int64, status int16, errorMsg string) error
	DeleteOlderThan(days int) (int64, error)
}

// GormAlertLogRepository 告警记录仓库GORM实现
type GormAlertLogRepository struct {
	db *gorm.DB
}

// NewGormAlertLogRepository 创建告警记录仓库
func NewGormAlertLogRepository(db *gorm.DB) *GormAlertLogRepository {
	return &GormAlertLogRepository{db: db}
}

// Create 创建告警记录
func (r *GormAlertLogRepository) Create(log *models.AlertLog) error {
	return r.db.Create(log).Error
}

// Update 更新告警记录
func (r *GormAlertLogRepository) Update(log *models.AlertLog) error {
	return r.db.Save(log).Error
}

// FindByID 根据ID查询记录
func (r *GormAlertLogRepository) FindByID(id int64) (*models.AlertLog, error) {
	var log models.AlertLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// FindPending 查询待发送的告警
func (r *GormAlertLogRepository) FindPending(limit int) ([]*models.AlertLog, error) {
	var logs []*models.AlertLog
	err := r.db.Where("send_status = ?", models.AlertSendStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// FindByJobName 根据任务名称查询告警记录
func (r *GormAlertLogRepository) FindByJobName(jobName string, limit, offset int) ([]*models.AlertLog, error) {
	var logs []*models.AlertLog
	query := r.db.Model(&models.AlertLog{})
	if jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}

// FindByDateRange 根据日期范围查询告警记录
func (r *GormAlertLogRepository) FindByDateRange(startDate, endDate time.Time, limit, offset int) ([]*models.AlertLog, error) {
	var logs []*models.AlertLog
	err := r.db.Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&logs).Error
	return logs, err
}

// CountByJobName 统计任务告警数量
func (r *GormAlertLogRepository) CountByJobName(jobName string) (int64, error) {
	var count int64
	query := r.db.Model(&models.AlertLog{})
	if jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	err := query.Count(&count).Error
	return count, err
}

// CountByDateRange 统计日期范围内告警数量
func (r *GormAlertLogRepository) CountByDateRange(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.AlertLog{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// UpdateSendStatus 更新发送状态
func (r *GormAlertLogRepository) UpdateSendStatus(id int64, status int16, errorMsg string) error {
	updates := map[string]interface{}{
		"send_status": status,
	}
	if status == models.AlertSendStatusSent {
		now := time.Now()
		updates["send_at"] = &now
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}
	return r.db.Model(&models.AlertLog{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteOlderThan 删除指定天数之前的记录
func (r *GormAlertLogRepository) DeleteOlderThan(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := r.db.Where("created_at < ?", cutoff).Delete(&models.AlertLog{})
	return result.RowsAffected, result.Error
}

// Ensure interface compliance
var _ AlertLogRepository = (*GormAlertLogRepository)(nil)

// JobFailCounterRepository 任务失败计数仓库接口
type JobFailCounterRepository interface {
	FindByJobName(jobName string) (*models.JobFailCounter, error)
	CreateOrUpdate(counter *models.JobFailCounter) error
	IncrementFail(jobName string) (*models.JobFailCounter, error)
	ResetOnSuccess(jobName string) error
	UpdateLastAlert(jobName string) error
}

// GormJobFailCounterRepository 任务失败计数仓库GORM实现
type GormJobFailCounterRepository struct {
	db *gorm.DB
}

// NewGormJobFailCounterRepository 创建任务失败计数仓库
func NewGormJobFailCounterRepository(db *gorm.DB) *GormJobFailCounterRepository {
	return &GormJobFailCounterRepository{db: db}
}

// FindByJobName 根据任务名称查询计数
func (r *GormJobFailCounterRepository) FindByJobName(jobName string) (*models.JobFailCounter, error) {
	var counter models.JobFailCounter
	err := r.db.Where("job_name = ?", jobName).First(&counter).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

// CreateOrUpdate 创建或更新计数
func (r *GormJobFailCounterRepository) CreateOrUpdate(counter *models.JobFailCounter) error {
	counter.UpdatedAt = time.Now()
	return r.db.Save(counter).Error
}

// IncrementFail 增加失败次数
func (r *GormJobFailCounterRepository) IncrementFail(jobName string) (*models.JobFailCounter, error) {
	now := time.Now()

	// 尝试查找现有记录
	var counter models.JobFailCounter
	err := r.db.Where("job_name = ?", jobName).First(&counter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新记录
			counter = models.JobFailCounter{
				JobName:          jobName,
				ConsecutiveFails: 1,
				LastFailAt:       &now,
				UpdatedAt:        now,
			}
			if err := r.db.Create(&counter).Error; err != nil {
				return nil, err
			}
			return &counter, nil
		}
		return nil, err
	}

	// 更新现有记录
	counter.ConsecutiveFails++
	counter.LastFailAt = &now
	counter.UpdatedAt = now
	if err := r.db.Save(&counter).Error; err != nil {
		return nil, err
	}
	return &counter, nil
}

// ResetOnSuccess 成功时重置计数
func (r *GormJobFailCounterRepository) ResetOnSuccess(jobName string) error {
	now := time.Now()

	// 尝试更新，如果不存在则创建
	result := r.db.Model(&models.JobFailCounter{}).
		Where("job_name = ?", jobName).
		Updates(map[string]interface{}{
			"consecutive_fails": 0,
			"last_success_at":   now,
			"updated_at":        now,
		})

	if result.RowsAffected == 0 {
		// 记录不存在，创建新记录
		counter := models.JobFailCounter{
			JobName:          jobName,
			ConsecutiveFails: 0,
			LastSuccessAt:    &now,
			UpdatedAt:        now,
		}
		return r.db.Create(&counter).Error
	}

	return result.Error
}

// UpdateLastAlert 更新最后告警时间
func (r *GormJobFailCounterRepository) UpdateLastAlert(jobName string) error {
	now := time.Now()
	return r.db.Model(&models.JobFailCounter{}).
		Where("job_name = ?", jobName).
		Updates(map[string]interface{}{
			"last_alert_at": now,
			"updated_at":    now,
		}).Error
}

// Ensure interface compliance
var _ JobFailCounterRepository = (*GormJobFailCounterRepository)(nil)
