package repository

import (
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// SettlementPriceRepository 结算价仓库接口
type SettlementPriceRepository interface {
	Create(price *models.SettlementPrice) error
	Update(price *models.SettlementPrice) error
	GetByID(id int64) (*models.SettlementPrice, error)
	GetByAgentAndChannel(agentID, channelID int64, brandCode string) (*models.SettlementPrice, error)
	List(req *models.SettlementPriceListRequest) ([]models.SettlementPrice, int64, error)
	ListByAgent(agentID int64) ([]models.SettlementPrice, error)
	Delete(id int64) error
}

// GormSettlementPriceRepository GORM实现
type GormSettlementPriceRepository struct {
	db *gorm.DB
}

// NewGormSettlementPriceRepository 创建结算价仓库
func NewGormSettlementPriceRepository(db *gorm.DB) *GormSettlementPriceRepository {
	return &GormSettlementPriceRepository{db: db}
}

// Create 创建结算价
func (r *GormSettlementPriceRepository) Create(price *models.SettlementPrice) error {
	return r.db.Create(price).Error
}

// Update 更新结算价
func (r *GormSettlementPriceRepository) Update(price *models.SettlementPrice) error {
	price.UpdatedAt = time.Now()
	return r.db.Save(price).Error
}

// GetByID 根据ID获取结算价
func (r *GormSettlementPriceRepository) GetByID(id int64) (*models.SettlementPrice, error) {
	var price models.SettlementPrice
	err := r.db.First(&price, id).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

// GetByAgentAndChannel 根据代理商ID和通道ID获取结算价
func (r *GormSettlementPriceRepository) GetByAgentAndChannel(agentID, channelID int64, brandCode string) (*models.SettlementPrice, error) {
	var price models.SettlementPrice
	query := r.db.Where("agent_id = ? AND channel_id = ? AND status = 1", agentID, channelID)
	if brandCode != "" {
		query = query.Where("brand_code = ?", brandCode)
	} else {
		query = query.Where("brand_code = '' OR brand_code IS NULL")
	}
	err := query.First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

// List 获取结算价列表
func (r *GormSettlementPriceRepository) List(req *models.SettlementPriceListRequest) ([]models.SettlementPrice, int64, error) {
	var prices []models.SettlementPrice
	var total int64

	query := r.db.Model(&models.SettlementPrice{})

	if req.AgentID != nil {
		query = query.Where("agent_id = ?", *req.AgentID)
	}
	if req.ChannelID != nil {
		query = query.Where("channel_id = ?", *req.ChannelID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	err = query.Order("updated_at DESC").Offset(offset).Limit(req.PageSize).Find(&prices).Error
	if err != nil {
		return nil, 0, err
	}

	return prices, total, nil
}

// ListByAgent 获取代理商的所有结算价
func (r *GormSettlementPriceRepository) ListByAgent(agentID int64) ([]models.SettlementPrice, error) {
	var prices []models.SettlementPrice
	err := r.db.Where("agent_id = ? AND status = 1", agentID).Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return prices, nil
}

// Delete 删除结算价（软删除，设置status=0）
func (r *GormSettlementPriceRepository) Delete(id int64) error {
	return r.db.Model(&models.SettlementPrice{}).Where("id = ?", id).Update("status", 0).Error
}
