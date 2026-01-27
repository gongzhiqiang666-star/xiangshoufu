package repository

import (
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// PriceChangeLogRepository 调价记录仓库接口
type PriceChangeLogRepository interface {
	Create(log *models.PriceChangeLog) error
	GetByID(id int64) (*models.PriceChangeLog, error)
	List(req *models.PriceChangeLogListRequest) ([]models.PriceChangeLog, int64, error)
	ListByAgent(agentID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error)
	ListBySettlementPrice(settlementPriceID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error)
	ListByRewardSetting(rewardSettingID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error)
}

// GormPriceChangeLogRepository GORM实现
type GormPriceChangeLogRepository struct {
	db *gorm.DB
}

// NewGormPriceChangeLogRepository 创建调价记录仓库
func NewGormPriceChangeLogRepository(db *gorm.DB) *GormPriceChangeLogRepository {
	return &GormPriceChangeLogRepository{db: db}
}

// Create 创建调价记录
func (r *GormPriceChangeLogRepository) Create(log *models.PriceChangeLog) error {
	return r.db.Create(log).Error
}

// GetByID 根据ID获取调价记录
func (r *GormPriceChangeLogRepository) GetByID(id int64) (*models.PriceChangeLog, error) {
	var log models.PriceChangeLog
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// List 获取调价记录列表
func (r *GormPriceChangeLogRepository) List(req *models.PriceChangeLogListRequest) ([]models.PriceChangeLog, int64, error) {
	var logs []models.PriceChangeLog
	var total int64

	query := r.db.Model(&models.PriceChangeLog{})

	if req.AgentID != nil {
		query = query.Where("agent_id = ?", *req.AgentID)
	}
	if req.ChannelID != nil {
		query = query.Where("channel_id = ?", *req.ChannelID)
	}
	if req.ChangeType != nil {
		query = query.Where("change_type = ?", *req.ChangeType)
	}
	if req.ConfigType != nil {
		query = query.Where("config_type = ?", *req.ConfigType)
	}
	if req.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if req.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			query = query.Where("created_at < ?", endTime.AddDate(0, 0, 1))
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ListByAgent 按代理商获取调价记录
func (r *GormPriceChangeLogRepository) ListByAgent(agentID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error) {
	var logs []models.PriceChangeLog
	var total int64

	query := r.db.Model(&models.PriceChangeLog{}).Where("agent_id = ?", agentID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ListBySettlementPrice 按结算价ID获取调价记录
func (r *GormPriceChangeLogRepository) ListBySettlementPrice(settlementPriceID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error) {
	var logs []models.PriceChangeLog
	var total int64

	query := r.db.Model(&models.PriceChangeLog{}).Where("settlement_price_id = ?", settlementPriceID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ListByRewardSetting 按奖励配置ID获取调价记录
func (r *GormPriceChangeLogRepository) ListByRewardSetting(rewardSettingID int64, page, pageSize int) ([]models.PriceChangeLog, int64, error) {
	var logs []models.PriceChangeLog
	var total int64

	query := r.db.Model(&models.PriceChangeLog{}).Where("reward_setting_id = ?", rewardSettingID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
