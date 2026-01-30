package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// GormWalletAdjustmentRepository GORM实现的钱包调账仓库
type GormWalletAdjustmentRepository struct {
	db *gorm.DB
}

// NewGormWalletAdjustmentRepository 创建仓库
func NewGormWalletAdjustmentRepository(db *gorm.DB) *GormWalletAdjustmentRepository {
	return &GormWalletAdjustmentRepository{db: db}
}

// Create 创建调账记录
func (r *GormWalletAdjustmentRepository) Create(adjustment *models.WalletAdjustment) error {
	return r.db.Create(adjustment).Error
}

// GetByID 根据ID获取调账记录
func (r *GormWalletAdjustmentRepository) GetByID(id int64) (*models.WalletAdjustment, error) {
	var adjustment models.WalletAdjustment
	err := r.db.First(&adjustment, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &adjustment, err
}

// GetByAdjustmentNo 根据调账单号获取调账记录
func (r *GormWalletAdjustmentRepository) GetByAdjustmentNo(adjustmentNo string) (*models.WalletAdjustment, error) {
	var adjustment models.WalletAdjustment
	err := r.db.Where("adjustment_no = ?", adjustmentNo).First(&adjustment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &adjustment, err
}

// WalletAdjustmentQueryParams 调账查询参数
type WalletAdjustmentQueryParams struct {
	AgentID    int64
	WalletType *int16
	ChannelID  *int64
	Status     *int16
	StartTime  *time.Time
	EndTime    *time.Time
	Limit      int
	Offset     int
}

// List 查询调账列表
func (r *GormWalletAdjustmentRepository) List(params *WalletAdjustmentQueryParams) ([]*models.WalletAdjustment, int64, error) {
	query := r.db.Model(&models.WalletAdjustment{})

	if params.AgentID > 0 {
		query = query.Where("agent_id = ?", params.AgentID)
	}
	if params.WalletType != nil {
		query = query.Where("wallet_type = ?", *params.WalletType)
	}
	if params.ChannelID != nil {
		query = query.Where("channel_id = ?", *params.ChannelID)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.StartTime != nil {
		query = query.Where("created_at >= ?", *params.StartTime)
	}
	if params.EndTime != nil {
		query = query.Where("created_at < ?", *params.EndTime)
	}

	var total int64
	query.Count(&total)

	var adjustments []*models.WalletAdjustment
	err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&adjustments).Error
	return adjustments, total, err
}

// Update 更新调账记录
func (r *GormWalletAdjustmentRepository) Update(adjustment *models.WalletAdjustment) error {
	return r.db.Save(adjustment).Error
}

// UpdateWalletLogID 更新关联的钱包流水ID
func (r *GormWalletAdjustmentRepository) UpdateWalletLogID(id int64, walletLogID int64) error {
	return r.db.Model(&models.WalletAdjustment{}).
		Where("id = ?", id).
		Update("wallet_log_id", walletLogID).Error
}

// GetDB 获取数据库连接（用于事务）
func (r *GormWalletAdjustmentRepository) GetDB() *gorm.DB {
	return r.db
}
