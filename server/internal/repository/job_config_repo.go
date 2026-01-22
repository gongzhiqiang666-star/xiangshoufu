package repository

import (
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// JobConfigRepository 任务配置仓库接口
type JobConfigRepository interface {
	FindAll() ([]*models.JobConfig, error)
	FindByName(jobName string) (*models.JobConfig, error)
	FindEnabled() ([]*models.JobConfig, error)
	Create(config *models.JobConfig) error
	Update(config *models.JobConfig) error
	UpdateEnabled(jobName string, isEnabled bool) error
}

// GormJobConfigRepository 任务配置仓库GORM实现
type GormJobConfigRepository struct {
	db *gorm.DB
}

// NewGormJobConfigRepository 创建任务配置仓库
func NewGormJobConfigRepository(db *gorm.DB) *GormJobConfigRepository {
	return &GormJobConfigRepository{db: db}
}

// FindAll 查询所有任务配置
func (r *GormJobConfigRepository) FindAll() ([]*models.JobConfig, error) {
	var configs []*models.JobConfig
	err := r.db.Order("id ASC").Find(&configs).Error
	return configs, err
}

// FindByName 根据任务名称查询配置
func (r *GormJobConfigRepository) FindByName(jobName string) (*models.JobConfig, error) {
	var config models.JobConfig
	err := r.db.Where("job_name = ?", jobName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FindEnabled 查询所有启用的任务配置
func (r *GormJobConfigRepository) FindEnabled() ([]*models.JobConfig, error) {
	var configs []*models.JobConfig
	err := r.db.Where("is_enabled = ?", true).Order("id ASC").Find(&configs).Error
	return configs, err
}

// Create 创建任务配置
func (r *GormJobConfigRepository) Create(config *models.JobConfig) error {
	return r.db.Create(config).Error
}

// Update 更新任务配置
func (r *GormJobConfigRepository) Update(config *models.JobConfig) error {
	config.UpdatedAt = time.Now()
	return r.db.Save(config).Error
}

// UpdateEnabled 更新任务启用状态
func (r *GormJobConfigRepository) UpdateEnabled(jobName string, isEnabled bool) error {
	return r.db.Model(&models.JobConfig{}).
		Where("job_name = ?", jobName).
		Updates(map[string]interface{}{
			"is_enabled": isEnabled,
			"updated_at": time.Now(),
		}).Error
}

// Ensure interface compliance
var _ JobConfigRepository = (*GormJobConfigRepository)(nil)
