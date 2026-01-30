package repository

import (
	"context"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// ChannelConfigRepository 通道配置仓库接口
type ChannelConfigRepository interface {
	// 费率配置
	GetRateConfigs(ctx context.Context, channelID int64) ([]models.ChannelRateConfig, error)
	GetRateConfigByCode(ctx context.Context, channelID int64, rateCode string) (*models.ChannelRateConfig, error)
	GetRateConfigByID(ctx context.Context, id int64) (*models.ChannelRateConfig, error)
	CreateRateConfig(ctx context.Context, config *models.ChannelRateConfig) error
	UpdateRateConfig(ctx context.Context, config *models.ChannelRateConfig) error
	DeleteRateConfig(ctx context.Context, id int64) error

	// 押金档位
	GetDepositTiers(ctx context.Context, channelID int64) ([]models.ChannelDepositTier, error)
	GetDepositTierByID(ctx context.Context, id int64) (*models.ChannelDepositTier, error)
	UpdateDepositTier(ctx context.Context, tier *models.ChannelDepositTier) error

	// 流量费返现档位
	GetSimCashbackTiers(ctx context.Context, channelID int64, brandCode string) ([]models.ChannelSimCashbackTier, error)
	BatchSetSimCashbackTiers(ctx context.Context, channelID int64, brandCode string, tiers []models.ChannelSimCashbackTier) error
	DeleteSimCashbackTiers(ctx context.Context, channelID int64, brandCode string) error

	// 通道完整配置
	GetFullConfig(ctx context.Context, channelID int64) (*models.ChannelFullConfig, error)

	// 影响检查相关
	GetPolicyTemplatesByChannel(ctx context.Context, channelID int64) ([]models.PolicyTemplateComplete, error)
	GetSettlementPricesByChannel(ctx context.Context, channelID int64) ([]models.SettlementPrice, error)
}

// GormChannelConfigRepository GORM实现的通道配置仓库
type GormChannelConfigRepository struct {
	db *gorm.DB
}

// NewGormChannelConfigRepository 创建通道配置仓库实例
func NewGormChannelConfigRepository(db *gorm.DB) *GormChannelConfigRepository {
	return &GormChannelConfigRepository{db: db}
}

// ============================================================
// 费率配置方法
// ============================================================

// GetRateConfigs 获取通道的所有费率配置
func (r *GormChannelConfigRepository) GetRateConfigs(ctx context.Context, channelID int64) ([]models.ChannelRateConfig, error) {
	var configs []models.ChannelRateConfig
	err := r.db.WithContext(ctx).
		Where("channel_id = ?", channelID).
		Order("sort_order ASC, id ASC").
		Find(&configs).Error
	return configs, err
}

// GetRateConfigByCode 根据费率编码获取配置
func (r *GormChannelConfigRepository) GetRateConfigByCode(ctx context.Context, channelID int64, rateCode string) (*models.ChannelRateConfig, error) {
	var config models.ChannelRateConfig
	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND rate_code = ?", channelID, rateCode).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateRateConfig 创建费率配置
func (r *GormChannelConfigRepository) CreateRateConfig(ctx context.Context, config *models.ChannelRateConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// UpdateRateConfig 更新费率配置
func (r *GormChannelConfigRepository) UpdateRateConfig(ctx context.Context, config *models.ChannelRateConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// DeleteRateConfig 删除费率配置
func (r *GormChannelConfigRepository) DeleteRateConfig(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.ChannelRateConfig{}, id).Error
}

// ============================================================
// 押金档位方法
// ============================================================

// GetDepositTiers 获取通道的所有押金档位
func (r *GormChannelConfigRepository) GetDepositTiers(ctx context.Context, channelID int64) ([]models.ChannelDepositTier, error) {
	var tiers []models.ChannelDepositTier
	err := r.db.WithContext(ctx).
		Where("channel_id = ?", channelID).
		Order("sort_order ASC, deposit_amount ASC").
		Find(&tiers).Error
	return tiers, err
}

// GetDepositTierByID 根据ID获取押金档位
func (r *GormChannelConfigRepository) GetDepositTierByID(ctx context.Context, id int64) (*models.ChannelDepositTier, error) {
	var tier models.ChannelDepositTier
	err := r.db.WithContext(ctx).First(&tier, id).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

// UpdateDepositTier 更新押金档位
func (r *GormChannelConfigRepository) UpdateDepositTier(ctx context.Context, tier *models.ChannelDepositTier) error {
	return r.db.WithContext(ctx).Save(tier).Error
}

// ============================================================
// 流量费返现档位方法
// ============================================================

// GetSimCashbackTiers 获取通道的流量费返现档位
func (r *GormChannelConfigRepository) GetSimCashbackTiers(ctx context.Context, channelID int64, brandCode string) ([]models.ChannelSimCashbackTier, error) {
	var tiers []models.ChannelSimCashbackTier
	query := r.db.WithContext(ctx).Where("channel_id = ?", channelID)
	if brandCode != "" {
		query = query.Where("brand_code = ? OR brand_code = ''", brandCode)
	}
	err := query.Order("tier_order ASC").Find(&tiers).Error
	return tiers, err
}

// BatchSetSimCashbackTiers 批量设置流量费返现档位
func (r *GormChannelConfigRepository) BatchSetSimCashbackTiers(ctx context.Context, channelID int64, brandCode string, tiers []models.ChannelSimCashbackTier) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有档位
		if err := tx.Where("channel_id = ? AND brand_code = ?", channelID, brandCode).
			Delete(&models.ChannelSimCashbackTier{}).Error; err != nil {
			return err
		}

		// 批量插入新档位
		if len(tiers) > 0 {
			for i := range tiers {
				tiers[i].ChannelID = channelID
				tiers[i].BrandCode = brandCode
			}
			if err := tx.Create(&tiers).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteSimCashbackTiers 删除流量费返现档位
func (r *GormChannelConfigRepository) DeleteSimCashbackTiers(ctx context.Context, channelID int64, brandCode string) error {
	return r.db.WithContext(ctx).
		Where("channel_id = ? AND brand_code = ?", channelID, brandCode).
		Delete(&models.ChannelSimCashbackTier{}).Error
}

// ============================================================
// 通道完整配置方法
// ============================================================

// GetFullConfig 获取通道完整配置
func (r *GormChannelConfigRepository) GetFullConfig(ctx context.Context, channelID int64) (*models.ChannelFullConfig, error) {
	// 获取通道信息
	var channel models.Channel
	if err := r.db.WithContext(ctx).First(&channel, channelID).Error; err != nil {
		return nil, err
	}

	// 获取费率配置
	rateConfigs, err := r.GetRateConfigs(ctx, channelID)
	if err != nil {
		return nil, err
	}

	// 获取押金档位
	depositTiers, err := r.GetDepositTiers(ctx, channelID)
	if err != nil {
		return nil, err
	}

	// 获取流量费返现档位（通用）
	simCashbackTiers, err := r.GetSimCashbackTiers(ctx, channelID, "")
	if err != nil {
		return nil, err
	}

	return &models.ChannelFullConfig{
		ChannelID:        channel.ID,
		ChannelCode:      channel.ChannelCode,
		ChannelName:      channel.ChannelName,
		RateConfigs:      rateConfigs,
		DepositTiers:     depositTiers,
		SimCashbackTiers: simCashbackTiers,
	}, nil
}

// GetRateConfigByID 根据ID获取费率配置
func (r *GormChannelConfigRepository) GetRateConfigByID(ctx context.Context, id int64) (*models.ChannelRateConfig, error) {
	var config models.ChannelRateConfig
	err := r.db.WithContext(ctx).First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetPolicyTemplatesByChannel 获取通道的政策模版列表
func (r *GormChannelConfigRepository) GetPolicyTemplatesByChannel(ctx context.Context, channelID int64) ([]models.PolicyTemplateComplete, error) {
	var templates []models.PolicyTemplateComplete
	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND status = 1", channelID).
		Find(&templates).Error
	return templates, err
}

// GetSettlementPricesByChannel 获取通道的结算价列表
func (r *GormChannelConfigRepository) GetSettlementPricesByChannel(ctx context.Context, channelID int64) ([]models.SettlementPrice, error) {
	var prices []models.SettlementPrice
	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND status = 1", channelID).
		Find(&prices).Error
	return prices, err
}
