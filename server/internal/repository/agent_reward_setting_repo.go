package repository

import (
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// AgentRewardSettingRepository 代理商奖励配置仓库接口
type AgentRewardSettingRepository interface {
	Create(setting *models.AgentRewardSetting) error
	Update(setting *models.AgentRewardSetting) error
	GetByID(id int64) (*models.AgentRewardSetting, error)
	GetByAgentID(agentID int64) (*models.AgentRewardSetting, error)
	List(page, pageSize int) ([]models.AgentRewardSetting, int64, error)
	Delete(id int64) error
}

// GormAgentRewardSettingRepository GORM实现
type GormAgentRewardSettingRepository struct {
	db *gorm.DB
}

// NewGormAgentRewardSettingRepository 创建代理商奖励配置仓库
func NewGormAgentRewardSettingRepository(db *gorm.DB) *GormAgentRewardSettingRepository {
	return &GormAgentRewardSettingRepository{db: db}
}

// Create 创建代理商奖励配置
func (r *GormAgentRewardSettingRepository) Create(setting *models.AgentRewardSetting) error {
	return r.db.Create(setting).Error
}

// Update 更新代理商奖励配置
func (r *GormAgentRewardSettingRepository) Update(setting *models.AgentRewardSetting) error {
	setting.UpdatedAt = time.Now()
	return r.db.Save(setting).Error
}

// GetByID 根据ID获取代理商奖励配置
func (r *GormAgentRewardSettingRepository) GetByID(id int64) (*models.AgentRewardSetting, error) {
	var setting models.AgentRewardSetting
	err := r.db.First(&setting, id).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetByAgentID 根据代理商ID获取奖励配置
func (r *GormAgentRewardSettingRepository) GetByAgentID(agentID int64) (*models.AgentRewardSetting, error) {
	var setting models.AgentRewardSetting
	err := r.db.Where("agent_id = ? AND status = 1", agentID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// List 获取代理商奖励配置列表
func (r *GormAgentRewardSettingRepository) List(page, pageSize int) ([]models.AgentRewardSetting, int64, error) {
	var settings []models.AgentRewardSetting
	var total int64

	query := r.db.Model(&models.AgentRewardSetting{}).Where("status = 1")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&settings).Error
	if err != nil {
		return nil, 0, err
	}

	return settings, total, nil
}

// Delete 删除代理商奖励配置（软删除）
func (r *GormAgentRewardSettingRepository) Delete(id int64) error {
	return r.db.Model(&models.AgentRewardSetting{}).Where("id = ?", id).Update("status", 0).Error
}
