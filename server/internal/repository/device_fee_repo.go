package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// GormDeviceFeeRepository GORM实现的流量费仓库
type GormDeviceFeeRepository struct {
	db *gorm.DB
}

// NewGormDeviceFeeRepository 创建仓库
func NewGormDeviceFeeRepository(db *gorm.DB) *GormDeviceFeeRepository {
	return &GormDeviceFeeRepository{db: db}
}

// Create 创建流量费记录
func (r *GormDeviceFeeRepository) Create(fee *models.DeviceFee) error {
	return r.db.Create(fee).Error
}

// Update 更新流量费记录
func (r *GormDeviceFeeRepository) Update(fee *models.DeviceFee) error {
	fee.UpdatedAt = time.Now()
	return r.db.Save(fee).Error
}

// FindByOrderNo 根据订单号查找
func (r *GormDeviceFeeRepository) FindByOrderNo(orderNo string) (*models.DeviceFee, error) {
	var fee models.DeviceFee
	err := r.db.Where("order_no = ?", orderNo).First(&fee).Error
	if err != nil {
		return nil, err
	}
	return &fee, nil
}

// FindPendingCashback 查找待返现的记录
func (r *GormDeviceFeeRepository) FindPendingCashback(limit int) ([]*models.DeviceFee, error) {
	var fees []*models.DeviceFee
	err := r.db.Where("cashback_status = ?", 0).
		Order("charging_time ASC").
		Limit(limit).
		Find(&fees).Error
	return fees, err
}

// UpdateCashbackStatus 更新返现状态
func (r *GormDeviceFeeRepository) UpdateCashbackStatus(id int64, status int16, amount int64) error {
	return r.db.Model(&models.DeviceFee{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"cashback_status": status,
			"cashback_amount": amount,
			"updated_at":      time.Now(),
		}).Error
}

// 确保实现了接口
var _ DeviceFeeRepository = (*GormDeviceFeeRepository)(nil)

// GormRateChangeRepository GORM实现的费率变更仓库
type GormRateChangeRepository struct {
	db *gorm.DB
}

// NewGormRateChangeRepository 创建仓库
func NewGormRateChangeRepository(db *gorm.DB) *GormRateChangeRepository {
	return &GormRateChangeRepository{db: db}
}

// Create 创建费率变更记录
func (r *GormRateChangeRepository) Create(change *models.RateChange) error {
	return r.db.Create(change).Error
}

// FindPendingSync 查找待同步的记录
func (r *GormRateChangeRepository) FindPendingSync(limit int) ([]*models.RateChange, error) {
	var changes []*models.RateChange
	err := r.db.Where("sync_status = ?", 0).
		Order("received_at ASC").
		Limit(limit).
		Find(&changes).Error
	return changes, err
}

// UpdateSyncStatus 更新同步状态
func (r *GormRateChangeRepository) UpdateSyncStatus(id int64, status int16) error {
	return r.db.Model(&models.RateChange{}).
		Where("id = ?", id).
		Update("sync_status", status).Error
}

// 确保实现了接口
var _ RateChangeRepository = (*GormRateChangeRepository)(nil)
