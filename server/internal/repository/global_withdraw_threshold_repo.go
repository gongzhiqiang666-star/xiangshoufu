package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GlobalWithdrawThresholdRepository 全局提现门槛仓库接口
type GlobalWithdrawThresholdRepository interface {
	// GetAll 获取所有门槛配置
	GetAll() ([]models.GlobalWithdrawThreshold, error)
	// GetByWalletType 获取指定钱包类型的所有门槛配置
	GetByWalletType(walletType int16) ([]models.GlobalWithdrawThreshold, error)
	// GetThreshold 获取指定钱包类型和通道的门槛
	GetThreshold(walletType int16, channelID int64) (*models.GlobalWithdrawThreshold, error)
	// UpsertBatch 批量更新或插入门槛配置
	UpsertBatch(thresholds []models.GlobalWithdrawThreshold) error
	// Delete 删除指定门槛配置
	Delete(id int64) error
	// DeleteByChannel 删除指定通道的门槛配置
	DeleteByChannel(channelID int64) error
}

// GormGlobalWithdrawThresholdRepository GORM实现
type GormGlobalWithdrawThresholdRepository struct {
	db *gorm.DB
}

// NewGlobalWithdrawThresholdRepository 创建仓库实例
func NewGlobalWithdrawThresholdRepository(db *gorm.DB) GlobalWithdrawThresholdRepository {
	return &GormGlobalWithdrawThresholdRepository{db: db}
}

// GetAll 获取所有门槛配置
func (r *GormGlobalWithdrawThresholdRepository) GetAll() ([]models.GlobalWithdrawThreshold, error) {
	var thresholds []models.GlobalWithdrawThreshold
	err := r.db.Order("wallet_type ASC, channel_id ASC").Find(&thresholds).Error
	return thresholds, err
}

// GetByWalletType 获取指定钱包类型的所有门槛配置
func (r *GormGlobalWithdrawThresholdRepository) GetByWalletType(walletType int16) ([]models.GlobalWithdrawThreshold, error) {
	var thresholds []models.GlobalWithdrawThreshold
	err := r.db.Where("wallet_type = ?", walletType).Order("channel_id ASC").Find(&thresholds).Error
	return thresholds, err
}

// GetThreshold 获取指定钱包类型和通道的门槛
// 优先返回特定通道的门槛，如果没有则返回通用门槛(channel_id=0)
func (r *GormGlobalWithdrawThresholdRepository) GetThreshold(walletType int16, channelID int64) (*models.GlobalWithdrawThreshold, error) {
	var threshold models.GlobalWithdrawThreshold

	// 先查询特定通道的门槛
	err := r.db.Where("wallet_type = ? AND channel_id = ?", walletType, channelID).First(&threshold).Error
	if err == nil {
		return &threshold, nil
	}

	// 如果没有特定通道门槛，查询通用门槛
	if channelID != 0 {
		err = r.db.Where("wallet_type = ? AND channel_id = 0", walletType).First(&threshold).Error
		if err == nil {
			return &threshold, nil
		}
	}

	return nil, err
}

// UpsertBatch 批量更新或插入门槛配置
func (r *GormGlobalWithdrawThresholdRepository) UpsertBatch(thresholds []models.GlobalWithdrawThreshold) error {
	if len(thresholds) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "wallet_type"}, {Name: "channel_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"threshold_amount", "updated_at"}),
	}).Create(&thresholds).Error
}

// Delete 删除指定门槛配置
func (r *GormGlobalWithdrawThresholdRepository) Delete(id int64) error {
	return r.db.Delete(&models.GlobalWithdrawThreshold{}, id).Error
}

// DeleteByChannel 删除指定通道的所有门槛配置
func (r *GormGlobalWithdrawThresholdRepository) DeleteByChannel(channelID int64) error {
	return r.db.Where("channel_id = ?", channelID).Delete(&models.GlobalWithdrawThreshold{}).Error
}
