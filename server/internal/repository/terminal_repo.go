package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// TerminalRepository 终端仓库接口
type TerminalRepository interface {
	Create(terminal *models.Terminal) error
	Update(terminal *models.Terminal) error
	FindByID(id int64) (*models.Terminal, error)
	FindBySN(terminalSN string) (*models.Terminal, error)
	FindByOwner(ownerAgentID int64, status []int16, limit, offset int) ([]*models.Terminal, int64, error)
	UpdateOwner(id int64, newOwnerID int64) error
	UpdateStatus(id int64, status int16) error
	UpdateSimFeeCount(id int64, count int) error
}

// TerminalDistributeRepository 终端下发仓库接口
type TerminalDistributeRepository interface {
	Create(distribute *models.TerminalDistribute) error
	Update(distribute *models.TerminalDistribute) error
	FindByID(id int64) (*models.TerminalDistribute, error)
	FindByDistributeNo(distributeNo string) (*models.TerminalDistribute, error)
	FindByFromAgent(fromAgentID int64, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error)
	FindByToAgent(toAgentID int64, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error)
	UpdateStatus(id int64, status int16, confirmedBy *int64) error
}

// SimCashbackPolicyRepository 流量费返现政策仓库接口
type SimCashbackPolicyRepository interface {
	Create(policy *models.SimCashbackPolicy) error
	Update(policy *models.SimCashbackPolicy) error
	FindByID(id int64) (*models.SimCashbackPolicy, error)
	FindByTemplateAndChannel(templateID, channelID int64, brandCode string) (*models.SimCashbackPolicy, error)
	FindByChannel(channelID int64) ([]*models.SimCashbackPolicy, error)
}

// SimCashbackRecordRepository 流量费返现记录仓库接口
type SimCashbackRecordRepository interface {
	Create(record *models.SimCashbackRecord) error
	BatchCreate(records []*models.SimCashbackRecord) error
	FindByDeviceFeeID(deviceFeeID int64) ([]*models.SimCashbackRecord, error)
	FindByAgent(agentID int64, limit, offset int) ([]*models.SimCashbackRecord, int64, error)
	UpdateWalletStatus(id int64, status int16) error
}

// GormTerminalRepository GORM实现
type GormTerminalRepository struct {
	db *gorm.DB
}

func NewGormTerminalRepository(db *gorm.DB) *GormTerminalRepository {
	return &GormTerminalRepository{db: db}
}

func (r *GormTerminalRepository) Create(terminal *models.Terminal) error {
	return r.db.Create(terminal).Error
}

func (r *GormTerminalRepository) Update(terminal *models.Terminal) error {
	return r.db.Save(terminal).Error
}

func (r *GormTerminalRepository) FindByID(id int64) (*models.Terminal, error) {
	var terminal models.Terminal
	err := r.db.First(&terminal, id).Error
	if err != nil {
		return nil, err
	}
	return &terminal, nil
}

func (r *GormTerminalRepository) FindBySN(terminalSN string) (*models.Terminal, error) {
	var terminal models.Terminal
	err := r.db.Where("terminal_sn = ?", terminalSN).First(&terminal).Error
	if err != nil {
		return nil, err
	}
	return &terminal, nil
}

func (r *GormTerminalRepository) FindByOwner(ownerAgentID int64, status []int16, limit, offset int) ([]*models.Terminal, int64, error) {
	var terminals []*models.Terminal
	var total int64

	query := r.db.Model(&models.Terminal{}).Where("owner_agent_id = ?", ownerAgentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&terminals).Error
	return terminals, total, err
}

func (r *GormTerminalRepository) UpdateOwner(id int64, newOwnerID int64) error {
	return r.db.Model(&models.Terminal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"owner_agent_id": newOwnerID,
			"status":         models.TerminalStatusAllocated,
			"updated_at":     time.Now(),
		}).Error
}

func (r *GormTerminalRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.Terminal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *GormTerminalRepository) UpdateSimFeeCount(id int64, count int) error {
	now := time.Now()
	return r.db.Model(&models.Terminal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"sim_fee_count":   count,
			"last_sim_fee_at": &now,
			"updated_at":      now,
		}).Error
}

var _ TerminalRepository = (*GormTerminalRepository)(nil)

// CountByOwnerAndStatus 按代理商和状态统计终端数量
func (r *GormTerminalRepository) CountByOwnerAndStatus(ownerAgentID int64, status int16, count *int64) error {
	return r.db.Model(&models.Terminal{}).
		Where("owner_agent_id = ? AND status = ?", ownerAgentID, status).
		Count(count).Error
}

// GormTerminalDistributeRepository GORM实现
type GormTerminalDistributeRepository struct {
	db *gorm.DB
}

func NewGormTerminalDistributeRepository(db *gorm.DB) *GormTerminalDistributeRepository {
	return &GormTerminalDistributeRepository{db: db}
}

func (r *GormTerminalDistributeRepository) Create(distribute *models.TerminalDistribute) error {
	return r.db.Create(distribute).Error
}

func (r *GormTerminalDistributeRepository) Update(distribute *models.TerminalDistribute) error {
	return r.db.Save(distribute).Error
}

func (r *GormTerminalDistributeRepository) FindByID(id int64) (*models.TerminalDistribute, error) {
	var distribute models.TerminalDistribute
	err := r.db.First(&distribute, id).Error
	if err != nil {
		return nil, err
	}
	return &distribute, nil
}

func (r *GormTerminalDistributeRepository) FindByDistributeNo(distributeNo string) (*models.TerminalDistribute, error) {
	var distribute models.TerminalDistribute
	err := r.db.Where("distribute_no = ?", distributeNo).First(&distribute).Error
	if err != nil {
		return nil, err
	}
	return &distribute, nil
}

func (r *GormTerminalDistributeRepository) FindByFromAgent(fromAgentID int64, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error) {
	var distributes []*models.TerminalDistribute
	var total int64

	query := r.db.Model(&models.TerminalDistribute{}).Where("from_agent_id = ?", fromAgentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&distributes).Error
	return distributes, total, err
}

func (r *GormTerminalDistributeRepository) FindByToAgent(toAgentID int64, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error) {
	var distributes []*models.TerminalDistribute
	var total int64

	query := r.db.Model(&models.TerminalDistribute{}).Where("to_agent_id = ?", toAgentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&distributes).Error
	return distributes, total, err
}

func (r *GormTerminalDistributeRepository) UpdateStatus(id int64, status int16, confirmedBy *int64) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"confirmed_at": &now,
	}
	if confirmedBy != nil {
		updates["confirmed_by"] = confirmedBy
	}
	return r.db.Model(&models.TerminalDistribute{}).Where("id = ?", id).Updates(updates).Error
}

var _ TerminalDistributeRepository = (*GormTerminalDistributeRepository)(nil)

// GormSimCashbackPolicyRepository GORM实现
type GormSimCashbackPolicyRepository struct {
	db *gorm.DB
}

func NewGormSimCashbackPolicyRepository(db *gorm.DB) *GormSimCashbackPolicyRepository {
	return &GormSimCashbackPolicyRepository{db: db}
}

func (r *GormSimCashbackPolicyRepository) Create(policy *models.SimCashbackPolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormSimCashbackPolicyRepository) Update(policy *models.SimCashbackPolicy) error {
	return r.db.Save(policy).Error
}

func (r *GormSimCashbackPolicyRepository) FindByID(id int64) (*models.SimCashbackPolicy, error) {
	var policy models.SimCashbackPolicy
	err := r.db.First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormSimCashbackPolicyRepository) FindByTemplateAndChannel(templateID, channelID int64, brandCode string) (*models.SimCashbackPolicy, error) {
	var policy models.SimCashbackPolicy
	query := r.db.Where("template_id = ? AND channel_id = ? AND status = 1", templateID, channelID)
	if brandCode != "" {
		query = query.Where("brand_code = ?", brandCode)
	}
	err := query.First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormSimCashbackPolicyRepository) FindByChannel(channelID int64) ([]*models.SimCashbackPolicy, error) {
	var policies []*models.SimCashbackPolicy
	err := r.db.Where("channel_id = ? AND status = 1", channelID).Find(&policies).Error
	return policies, err
}

var _ SimCashbackPolicyRepository = (*GormSimCashbackPolicyRepository)(nil)

// GormSimCashbackRecordRepository GORM实现
type GormSimCashbackRecordRepository struct {
	db *gorm.DB
}

func NewGormSimCashbackRecordRepository(db *gorm.DB) *GormSimCashbackRecordRepository {
	return &GormSimCashbackRecordRepository{db: db}
}

func (r *GormSimCashbackRecordRepository) Create(record *models.SimCashbackRecord) error {
	return r.db.Create(record).Error
}

func (r *GormSimCashbackRecordRepository) BatchCreate(records []*models.SimCashbackRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.CreateInBatches(records, 100).Error
}

func (r *GormSimCashbackRecordRepository) FindByDeviceFeeID(deviceFeeID int64) ([]*models.SimCashbackRecord, error) {
	var records []*models.SimCashbackRecord
	err := r.db.Where("device_fee_id = ?", deviceFeeID).Find(&records).Error
	return records, err
}

func (r *GormSimCashbackRecordRepository) FindByAgent(agentID int64, limit, offset int) ([]*models.SimCashbackRecord, int64, error) {
	var records []*models.SimCashbackRecord
	var total int64

	query := r.db.Model(&models.SimCashbackRecord{}).Where("agent_id = ?", agentID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error
	return records, total, err
}

func (r *GormSimCashbackRecordRepository) UpdateWalletStatus(id int64, status int16) error {
	now := time.Now()
	return r.db.Model(&models.SimCashbackRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"wallet_status": status,
			"processed_at":  &now,
		}).Error
}

var _ SimCashbackRecordRepository = (*GormSimCashbackRecordRepository)(nil)
