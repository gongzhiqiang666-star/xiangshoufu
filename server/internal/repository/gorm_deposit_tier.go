package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// GormChannelDepositTierRepository 通道押金档位仓库GORM实现
type GormChannelDepositTierRepository struct {
	db *gorm.DB
}

// NewGormChannelDepositTierRepository 创建通道押金档位仓库
func NewGormChannelDepositTierRepository(db *gorm.DB) *GormChannelDepositTierRepository {
	return &GormChannelDepositTierRepository{db: db}
}

// Create 创建押金档位
func (r *GormChannelDepositTierRepository) Create(tier *models.ChannelDepositTier) error {
	return r.db.Create(tier).Error
}

// Update 更新押金档位
func (r *GormChannelDepositTierRepository) Update(tier *models.ChannelDepositTier) error {
	return r.db.Save(tier).Error
}

// Delete 删除押金档位
func (r *GormChannelDepositTierRepository) Delete(id int64) error {
	return r.db.Delete(&models.ChannelDepositTier{}, id).Error
}

// FindByID 根据ID查找
func (r *GormChannelDepositTierRepository) FindByID(id int64) (*models.ChannelDepositTier, error) {
	var tier models.ChannelDepositTier
	err := r.db.First(&tier, id).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

// FindByChannelID 根据通道ID查找所有档位
func (r *GormChannelDepositTierRepository) FindByChannelID(channelID int64) ([]*models.ChannelDepositTier, error) {
	var tiers []*models.ChannelDepositTier
	err := r.db.Where("channel_id = ? AND status = 1", channelID).
		Order("sort_order ASC, deposit_amount ASC").
		Find(&tiers).Error
	return tiers, err
}

// FindByChannelAndBrand 根据通道ID和品牌编码查找
func (r *GormChannelDepositTierRepository) FindByChannelAndBrand(channelID int64, brandCode string) ([]*models.ChannelDepositTier, error) {
	var tiers []*models.ChannelDepositTier
	err := r.db.Where("channel_id = ? AND brand_code = ? AND status = 1", channelID, brandCode).
		Order("sort_order ASC, deposit_amount ASC").
		Find(&tiers).Error
	return tiers, err
}

// FindByTierCode 根据档位编码查找
func (r *GormChannelDepositTierRepository) FindByTierCode(channelID int64, brandCode string, tierCode string) (*models.ChannelDepositTier, error) {
	var tier models.ChannelDepositTier
	err := r.db.Where("channel_id = ? AND brand_code = ? AND tier_code = ?", channelID, brandCode, tierCode).
		First(&tier).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}
