package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// AgentChannelRepository 代理商通道仓库接口
type AgentChannelRepository interface {
	Create(ac *models.AgentChannel) error
	Update(ac *models.AgentChannel) error
	Upsert(ac *models.AgentChannel) error
	FindByAgentID(agentID int64) ([]*models.AgentChannel, error)
	FindByAgentIDWithInfo(agentID int64) ([]*models.AgentChannelWithInfo, error)
	FindEnabledByAgentID(agentID int64) ([]*models.AgentChannel, error)
	FindByAgentAndChannel(agentID, channelID int64) (*models.AgentChannel, error)
	Enable(agentID, channelID int64, enabledBy int64) error
	Disable(agentID, channelID int64, disabledBy int64) error
	SetVisibility(agentID, channelID int64, isVisible bool) error
	Delete(agentID, channelID int64) error
	BatchEnable(agentID int64, channelIDs []int64, enabledBy int64) error
	BatchDisable(agentID int64, channelIDs []int64, disabledBy int64) error
}

// GormAgentChannelRepository GORM实现
type GormAgentChannelRepository struct {
	db *gorm.DB
}

// NewGormAgentChannelRepository 创建代理商通道仓库
func NewGormAgentChannelRepository(db *gorm.DB) *GormAgentChannelRepository {
	return &GormAgentChannelRepository{db: db}
}

func (r *GormAgentChannelRepository) Create(ac *models.AgentChannel) error {
	return r.db.Create(ac).Error
}

func (r *GormAgentChannelRepository) Update(ac *models.AgentChannel) error {
	ac.UpdatedAt = time.Now()
	return r.db.Save(ac).Error
}

func (r *GormAgentChannelRepository) Upsert(ac *models.AgentChannel) error {
	return r.db.Where("agent_id = ? AND channel_id = ?", ac.AgentID, ac.ChannelID).
		Assign(map[string]interface{}{
			"is_enabled": ac.IsEnabled,
			"is_visible": ac.IsVisible,
			"updated_at": time.Now(),
		}).FirstOrCreate(ac).Error
}

func (r *GormAgentChannelRepository) FindByAgentID(agentID int64) ([]*models.AgentChannel, error) {
	var channels []*models.AgentChannel
	err := r.db.Where("agent_id = ?", agentID).
		Order("channel_id ASC").
		Find(&channels).Error
	return channels, err
}

func (r *GormAgentChannelRepository) FindByAgentIDWithInfo(agentID int64) ([]*models.AgentChannelWithInfo, error) {
	var results []*models.AgentChannelWithInfo
	err := r.db.Table("agent_channels ac").
		Select("ac.*, c.channel_code, c.channel_name").
		Joins("LEFT JOIN channels c ON ac.channel_id = c.id").
		Where("ac.agent_id = ?", agentID).
		Order("c.priority DESC, ac.channel_id ASC").
		Scan(&results).Error
	return results, err
}

func (r *GormAgentChannelRepository) FindEnabledByAgentID(agentID int64) ([]*models.AgentChannel, error) {
	var channels []*models.AgentChannel
	err := r.db.Where("agent_id = ? AND is_enabled = true", agentID).
		Order("channel_id ASC").
		Find(&channels).Error
	return channels, err
}

func (r *GormAgentChannelRepository) FindByAgentAndChannel(agentID, channelID int64) (*models.AgentChannel, error) {
	var ac models.AgentChannel
	err := r.db.Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		First(&ac).Error
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *GormAgentChannelRepository) Enable(agentID, channelID int64, enabledBy int64) error {
	now := time.Now()
	return r.db.Model(&models.AgentChannel{}).
		Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Updates(map[string]interface{}{
			"is_enabled": true,
			"enabled_at": &now,
			"enabled_by": enabledBy,
			"updated_at": now,
		}).Error
}

func (r *GormAgentChannelRepository) Disable(agentID, channelID int64, disabledBy int64) error {
	now := time.Now()
	return r.db.Model(&models.AgentChannel{}).
		Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Updates(map[string]interface{}{
			"is_enabled":  false,
			"disabled_at": &now,
			"disabled_by": disabledBy,
			"updated_at":  now,
		}).Error
}

func (r *GormAgentChannelRepository) SetVisibility(agentID, channelID int64, isVisible bool) error {
	return r.db.Model(&models.AgentChannel{}).
		Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Updates(map[string]interface{}{
			"is_visible": isVisible,
			"updated_at": time.Now(),
		}).Error
}

func (r *GormAgentChannelRepository) Delete(agentID, channelID int64) error {
	return r.db.Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Delete(&models.AgentChannel{}).Error
}

func (r *GormAgentChannelRepository) BatchEnable(agentID int64, channelIDs []int64, enabledBy int64) error {
	if len(channelIDs) == 0 {
		return nil
	}
	now := time.Now()
	return r.db.Model(&models.AgentChannel{}).
		Where("agent_id = ? AND channel_id IN ?", agentID, channelIDs).
		Updates(map[string]interface{}{
			"is_enabled": true,
			"enabled_at": &now,
			"enabled_by": enabledBy,
			"updated_at": now,
		}).Error
}

func (r *GormAgentChannelRepository) BatchDisable(agentID int64, channelIDs []int64, disabledBy int64) error {
	if len(channelIDs) == 0 {
		return nil
	}
	now := time.Now()
	return r.db.Model(&models.AgentChannel{}).
		Where("agent_id = ? AND channel_id IN ?", agentID, channelIDs).
		Updates(map[string]interface{}{
			"is_enabled":  false,
			"disabled_at": &now,
			"disabled_by": disabledBy,
			"updated_at":  now,
		}).Error
}

// InitAgentChannels 初始化代理商通道配置（为代理商创建所有可用通道的配置）
func (r *GormAgentChannelRepository) InitAgentChannels(agentID int64, channelIDs []int64, enabledBy int64) error {
	now := time.Now()
	for _, channelID := range channelIDs {
		ac := &models.AgentChannel{
			AgentID:   agentID,
			ChannelID: channelID,
			IsEnabled: true,
			IsVisible: true,
			EnabledAt: &now,
			EnabledBy: &enabledBy,
			CreatedAt: now,
			UpdatedAt: now,
		}
		// 使用Upsert避免重复
		if err := r.Upsert(ac); err != nil {
			return err
		}
	}
	return nil
}

var _ AgentChannelRepository = (*GormAgentChannelRepository)(nil)
