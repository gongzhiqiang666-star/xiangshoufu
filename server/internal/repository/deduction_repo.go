package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// DeductionPlanRepository 代扣计划仓库接口
type DeductionPlanRepository interface {
	Create(plan *models.DeductionPlan) error
	Update(plan *models.DeductionPlan) error
	FindByID(id int64) (*models.DeductionPlan, error)
	FindByPlanNo(planNo string) (*models.DeductionPlan, error)
	FindByDeductee(deducteeID int64, status []int16, limit, offset int) ([]*models.DeductionPlan, int64, error)
	FindByDeductor(deductorID int64, status []int16, limit, offset int) ([]*models.DeductionPlan, int64, error)
	FindActivePlans(limit int) ([]*models.DeductionPlan, error)
	UpdateStatus(id int64, status int16) error
	UpdateDeductedAmount(id int64, amount int64, currentPeriod int) error
}

// DeductionRecordRepository 代扣记录仓库接口
type DeductionRecordRepository interface {
	Create(record *models.DeductionRecord) error
	BatchCreate(records []*models.DeductionRecord) error
	FindByID(id int64) (*models.DeductionRecord, error)
	FindByPlanID(planID int64) ([]*models.DeductionRecord, error)
	FindPendingRecords(scheduledBefore time.Time, limit int) ([]*models.DeductionRecord, error)
	UpdateStatus(id int64, status int16, actualAmount int64, walletDetails string, failReason string) error
}

// DeductionChainRepository 代扣链仓库接口
type DeductionChainRepository interface {
	Create(chain *models.DeductionChain) error
	FindByID(id int64) (*models.DeductionChain, error)
	FindByDistributeID(distributeID int64) (*models.DeductionChain, error)
	UpdateStatus(id int64, status int16) error
}

// DeductionChainItemRepository 代扣链节点仓库接口
type DeductionChainItemRepository interface {
	Create(item *models.DeductionChainItem) error
	BatchCreate(items []*models.DeductionChainItem) error
	FindByChainID(chainID int64) ([]*models.DeductionChainItem, error)
	UpdatePlanID(id int64, planID int64) error
	UpdateStatus(id int64, status int16) error
}

// GormDeductionPlanRepository GORM实现
type GormDeductionPlanRepository struct {
	db *gorm.DB
}

func NewGormDeductionPlanRepository(db *gorm.DB) *GormDeductionPlanRepository {
	return &GormDeductionPlanRepository{db: db}
}

func (r *GormDeductionPlanRepository) Create(plan *models.DeductionPlan) error {
	return r.db.Create(plan).Error
}

func (r *GormDeductionPlanRepository) Update(plan *models.DeductionPlan) error {
	return r.db.Save(plan).Error
}

func (r *GormDeductionPlanRepository) FindByID(id int64) (*models.DeductionPlan, error) {
	var plan models.DeductionPlan
	err := r.db.First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *GormDeductionPlanRepository) FindByPlanNo(planNo string) (*models.DeductionPlan, error) {
	var plan models.DeductionPlan
	err := r.db.Where("plan_no = ?", planNo).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *GormDeductionPlanRepository) FindByDeductee(deducteeID int64, status []int16, limit, offset int) ([]*models.DeductionPlan, int64, error) {
	var plans []*models.DeductionPlan
	var total int64

	query := r.db.Model(&models.DeductionPlan{}).Where("deductee_id = ?", deducteeID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&plans).Error
	return plans, total, err
}

func (r *GormDeductionPlanRepository) FindByDeductor(deductorID int64, status []int16, limit, offset int) ([]*models.DeductionPlan, int64, error) {
	var plans []*models.DeductionPlan
	var total int64

	query := r.db.Model(&models.DeductionPlan{}).Where("deductor_id = ?", deductorID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&plans).Error
	return plans, total, err
}

func (r *GormDeductionPlanRepository) FindActivePlans(limit int) ([]*models.DeductionPlan, error) {
	var plans []*models.DeductionPlan
	err := r.db.Where("status = ?", models.DeductionPlanStatusActive).
		Where("remaining_amount > 0").
		Order("created_at ASC").
		Limit(limit).
		Find(&plans).Error
	return plans, err
}

func (r *GormDeductionPlanRepository) UpdateStatus(id int64, status int16) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if status == models.DeductionPlanStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return r.db.Model(&models.DeductionPlan{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormDeductionPlanRepository) UpdateDeductedAmount(id int64, amount int64, currentPeriod int) error {
	return r.db.Model(&models.DeductionPlan{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deducted_amount":  gorm.Expr("deducted_amount + ?", amount),
			"remaining_amount": gorm.Expr("remaining_amount - ?", amount),
			"current_period":   currentPeriod,
			"updated_at":       time.Now(),
		}).Error
}

var _ DeductionPlanRepository = (*GormDeductionPlanRepository)(nil)

// GormDeductionRecordRepository GORM实现
type GormDeductionRecordRepository struct {
	db *gorm.DB
}

func NewGormDeductionRecordRepository(db *gorm.DB) *GormDeductionRecordRepository {
	return &GormDeductionRecordRepository{db: db}
}

func (r *GormDeductionRecordRepository) Create(record *models.DeductionRecord) error {
	return r.db.Create(record).Error
}

func (r *GormDeductionRecordRepository) BatchCreate(records []*models.DeductionRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.CreateInBatches(records, 100).Error
}

func (r *GormDeductionRecordRepository) FindByID(id int64) (*models.DeductionRecord, error) {
	var record models.DeductionRecord
	err := r.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *GormDeductionRecordRepository) FindByPlanID(planID int64) ([]*models.DeductionRecord, error) {
	var records []*models.DeductionRecord
	err := r.db.Where("plan_id = ?", planID).Order("period_num ASC").Find(&records).Error
	return records, err
}

func (r *GormDeductionRecordRepository) FindPendingRecords(scheduledBefore time.Time, limit int) ([]*models.DeductionRecord, error) {
	var records []*models.DeductionRecord
	err := r.db.Where("status = ?", models.DeductionRecordStatusPending).
		Where("scheduled_at <= ?", scheduledBefore).
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

func (r *GormDeductionRecordRepository) UpdateStatus(id int64, status int16, actualAmount int64, walletDetails string, failReason string) error {
	now := time.Now()
	return r.db.Model(&models.DeductionRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"actual_amount":  actualAmount,
			"wallet_details": walletDetails,
			"fail_reason":    failReason,
			"deducted_at":    &now,
		}).Error
}

var _ DeductionRecordRepository = (*GormDeductionRecordRepository)(nil)

// GormDeductionChainRepository GORM实现
type GormDeductionChainRepository struct {
	db *gorm.DB
}

func NewGormDeductionChainRepository(db *gorm.DB) *GormDeductionChainRepository {
	return &GormDeductionChainRepository{db: db}
}

func (r *GormDeductionChainRepository) Create(chain *models.DeductionChain) error {
	return r.db.Create(chain).Error
}

func (r *GormDeductionChainRepository) FindByID(id int64) (*models.DeductionChain, error) {
	var chain models.DeductionChain
	err := r.db.First(&chain, id).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

func (r *GormDeductionChainRepository) FindByDistributeID(distributeID int64) (*models.DeductionChain, error) {
	var chain models.DeductionChain
	err := r.db.Where("distribute_id = ?", distributeID).First(&chain).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

func (r *GormDeductionChainRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.DeductionChain{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

var _ DeductionChainRepository = (*GormDeductionChainRepository)(nil)

// GormDeductionChainItemRepository GORM实现
type GormDeductionChainItemRepository struct {
	db *gorm.DB
}

func NewGormDeductionChainItemRepository(db *gorm.DB) *GormDeductionChainItemRepository {
	return &GormDeductionChainItemRepository{db: db}
}

func (r *GormDeductionChainItemRepository) Create(item *models.DeductionChainItem) error {
	return r.db.Create(item).Error
}

func (r *GormDeductionChainItemRepository) BatchCreate(items []*models.DeductionChainItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.CreateInBatches(items, 100).Error
}

func (r *GormDeductionChainItemRepository) FindByChainID(chainID int64) ([]*models.DeductionChainItem, error) {
	var items []*models.DeductionChainItem
	err := r.db.Where("chain_id = ?", chainID).Order("level ASC").Find(&items).Error
	return items, err
}

func (r *GormDeductionChainItemRepository) UpdatePlanID(id int64, planID int64) error {
	return r.db.Model(&models.DeductionChainItem{}).
		Where("id = ?", id).
		Update("plan_id", planID).Error
}

func (r *GormDeductionChainItemRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.DeductionChainItem{}).
		Where("id = ?", id).
		Update("status", status).Error
}

var _ DeductionChainItemRepository = (*GormDeductionChainItemRepository)(nil)
