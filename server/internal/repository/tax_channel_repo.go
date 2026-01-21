package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// TaxChannelRepository 税筹通道仓库接口
type TaxChannelRepository interface {
	// 税筹通道
	Create(taxChannel *models.TaxChannel) error
	Update(taxChannel *models.TaxChannel) error
	Delete(id int64) error
	GetByID(id int64) (*models.TaxChannel, error)
	GetByCode(code string) (*models.TaxChannel, error)
	GetAll(status *int16) ([]*models.TaxChannel, error)

	// 通道-税筹通道关联
	CreateMapping(mapping *models.ChannelTaxMapping) error
	UpdateMapping(mapping *models.ChannelTaxMapping) error
	DeleteMapping(id int64) error
	GetMappingsByChannel(channelID int64) ([]*models.ChannelTaxMapping, error)
	GetMappingByChannelAndWallet(channelID int64, walletType int16) (*models.ChannelTaxMapping, error)
	GetTaxChannelForWithdrawal(channelID int64, walletType int16) (*models.TaxChannel, error)
}

// GormTaxChannelRepository GORM实现
type GormTaxChannelRepository struct {
	db *gorm.DB
}

// NewGormTaxChannelRepository 创建税筹通道仓库
func NewGormTaxChannelRepository(db *gorm.DB) *GormTaxChannelRepository {
	return &GormTaxChannelRepository{db: db}
}

// ========== 税筹通道 ==========

// Create 创建税筹通道
func (r *GormTaxChannelRepository) Create(taxChannel *models.TaxChannel) error {
	return r.db.Create(taxChannel).Error
}

// Update 更新税筹通道
func (r *GormTaxChannelRepository) Update(taxChannel *models.TaxChannel) error {
	return r.db.Save(taxChannel).Error
}

// Delete 删除税筹通道
func (r *GormTaxChannelRepository) Delete(id int64) error {
	return r.db.Delete(&models.TaxChannel{}, id).Error
}

// GetByID 根据ID获取税筹通道
func (r *GormTaxChannelRepository) GetByID(id int64) (*models.TaxChannel, error) {
	var taxChannel models.TaxChannel
	err := r.db.First(&taxChannel, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &taxChannel, err
}

// GetByCode 根据编码获取税筹通道
func (r *GormTaxChannelRepository) GetByCode(code string) (*models.TaxChannel, error) {
	var taxChannel models.TaxChannel
	err := r.db.Where("channel_code = ?", code).First(&taxChannel).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &taxChannel, err
}

// GetAll 获取所有税筹通道
func (r *GormTaxChannelRepository) GetAll(status *int16) ([]*models.TaxChannel, error) {
	query := r.db.Model(&models.TaxChannel{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var taxChannels []*models.TaxChannel
	err := query.Order("id ASC").Find(&taxChannels).Error
	return taxChannels, err
}

// ========== 通道-税筹通道关联 ==========

// CreateMapping 创建关联
func (r *GormTaxChannelRepository) CreateMapping(mapping *models.ChannelTaxMapping) error {
	return r.db.Create(mapping).Error
}

// UpdateMapping 更新关联
func (r *GormTaxChannelRepository) UpdateMapping(mapping *models.ChannelTaxMapping) error {
	return r.db.Save(mapping).Error
}

// DeleteMapping 删除关联
func (r *GormTaxChannelRepository) DeleteMapping(id int64) error {
	return r.db.Delete(&models.ChannelTaxMapping{}, id).Error
}

// GetMappingsByChannel 获取支付通道的所有税筹通道关联
func (r *GormTaxChannelRepository) GetMappingsByChannel(channelID int64) ([]*models.ChannelTaxMapping, error) {
	var mappings []*models.ChannelTaxMapping
	err := r.db.Where("channel_id = ?", channelID).Find(&mappings).Error
	return mappings, err
}

// GetMappingByChannelAndWallet 获取特定支付通道和钱包类型的税筹通道关联
func (r *GormTaxChannelRepository) GetMappingByChannelAndWallet(channelID int64, walletType int16) (*models.ChannelTaxMapping, error) {
	var mapping models.ChannelTaxMapping
	err := r.db.Where("channel_id = ? AND wallet_type = ?", channelID, walletType).First(&mapping).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mapping, err
}

// GetTaxChannelForWithdrawal 获取提现使用的税筹通道
func (r *GormTaxChannelRepository) GetTaxChannelForWithdrawal(channelID int64, walletType int16) (*models.TaxChannel, error) {
	// 先查找特定通道和钱包类型的映射
	mapping, err := r.GetMappingByChannelAndWallet(channelID, walletType)
	if err != nil {
		return nil, err
	}

	var taxChannelID int64
	if mapping != nil {
		taxChannelID = mapping.TaxChannelID
	} else {
		// 没有映射，使用默认税筹通道
		defaultChannel, err := r.GetByCode("DEFAULT")
		if err != nil || defaultChannel == nil {
			return nil, nil
		}
		return defaultChannel, nil
	}

	// 获取税筹通道
	return r.GetByID(taxChannelID)
}

// TaxChannelWithMapping 税筹通道及其映射信息
type TaxChannelWithMapping struct {
	TaxChannel *models.TaxChannel
	Mapping    *models.ChannelTaxMapping
}

// GetAllWithMappings 获取所有税筹通道及其关联的支付通道
func (r *GormTaxChannelRepository) GetAllWithMappings() ([]*TaxChannelWithMapping, error) {
	var taxChannels []*models.TaxChannel
	if err := r.db.Order("id ASC").Find(&taxChannels).Error; err != nil {
		return nil, err
	}

	result := make([]*TaxChannelWithMapping, 0, len(taxChannels))
	for _, tc := range taxChannels {
		result = append(result, &TaxChannelWithMapping{
			TaxChannel: tc,
		})
	}
	return result, nil
}

var _ TaxChannelRepository = (*GormTaxChannelRepository)(nil)
