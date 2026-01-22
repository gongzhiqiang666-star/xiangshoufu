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
	FindActivatedAfter(channelID int64, activatedAfter time.Time) ([]*models.Terminal, error)
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
	FindByID(id int64) (*models.SimCashbackRecord, error)
	FindByDeviceFeeID(deviceFeeID int64) ([]*models.SimCashbackRecord, error)
	FindByAgent(agentID int64, limit, offset int) ([]*models.SimCashbackRecord, int64, error)
	UpdateWalletStatus(id int64, status int16) error
	FindPending(limit int) ([]*models.SimCashbackRecord, error)
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

// FindActivatedAfter 查找指定通道在某时间后激活的终端
func (r *GormTerminalRepository) FindActivatedAfter(channelID int64, activatedAfter time.Time) ([]*models.Terminal, error) {
	var terminals []*models.Terminal
	err := r.db.Where("channel_id = ? AND status = ? AND activated_at >= ?",
		channelID, models.TerminalStatusActivated, activatedAfter).
		Find(&terminals).Error
	return terminals, err
}

// CountByOwnerAndStatus 按代理商和状态统计终端数量
func (r *GormTerminalRepository) CountByOwnerAndStatus(ownerAgentID int64, status int16, count *int64) error {
	return r.db.Model(&models.Terminal{}).
		Where("owner_agent_id = ? AND status = ?", ownerAgentID, status).
		Count(count).Error
}

// CountByMerchantNo 按商户号统计终端数量
func (r *GormTerminalRepository) CountByMerchantNo(merchantNo string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Terminal{}).
		Where("merchant_no = ?", merchantNo).
		Count(&count).Error
	return count, err
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

func (r *GormSimCashbackRecordRepository) FindByID(id int64) (*models.SimCashbackRecord, error) {
	var record models.SimCashbackRecord
	err := r.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
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

func (r *GormSimCashbackRecordRepository) FindPending(limit int) ([]*models.SimCashbackRecord, error) {
	var records []*models.SimCashbackRecord
	err := r.db.Where("wallet_status = 0").
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

var _ SimCashbackRecordRepository = (*GormSimCashbackRecordRepository)(nil)

// TerminalRecallRepository 终端回拨仓库接口
type TerminalRecallRepository interface {
	Create(recall *models.TerminalRecall) error
	Update(recall *models.TerminalRecall) error
	FindByID(id int64) (*models.TerminalRecall, error)
	FindByRecallNo(recallNo string) (*models.TerminalRecall, error)
	FindByFromAgent(fromAgentID int64, status []int16, limit, offset int) ([]*models.TerminalRecall, int64, error)
	FindByToAgent(toAgentID int64, status []int16, limit, offset int) ([]*models.TerminalRecall, int64, error)
	UpdateStatus(id int64, status int16, confirmedBy *int64) error
}

// GormTerminalRecallRepository GORM实现
type GormTerminalRecallRepository struct {
	db *gorm.DB
}

func NewGormTerminalRecallRepository(db *gorm.DB) *GormTerminalRecallRepository {
	return &GormTerminalRecallRepository{db: db}
}

func (r *GormTerminalRecallRepository) Create(recall *models.TerminalRecall) error {
	return r.db.Create(recall).Error
}

func (r *GormTerminalRecallRepository) Update(recall *models.TerminalRecall) error {
	return r.db.Save(recall).Error
}

func (r *GormTerminalRecallRepository) FindByID(id int64) (*models.TerminalRecall, error) {
	var recall models.TerminalRecall
	err := r.db.First(&recall, id).Error
	if err != nil {
		return nil, err
	}
	return &recall, nil
}

func (r *GormTerminalRecallRepository) FindByRecallNo(recallNo string) (*models.TerminalRecall, error) {
	var recall models.TerminalRecall
	err := r.db.Where("recall_no = ?", recallNo).First(&recall).Error
	if err != nil {
		return nil, err
	}
	return &recall, nil
}

func (r *GormTerminalRecallRepository) FindByFromAgent(fromAgentID int64, status []int16, limit, offset int) ([]*models.TerminalRecall, int64, error) {
	var recalls []*models.TerminalRecall
	var total int64

	query := r.db.Model(&models.TerminalRecall{}).Where("from_agent_id = ?", fromAgentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&recalls).Error
	return recalls, total, err
}

func (r *GormTerminalRecallRepository) FindByToAgent(toAgentID int64, status []int16, limit, offset int) ([]*models.TerminalRecall, int64, error) {
	var recalls []*models.TerminalRecall
	var total int64

	query := r.db.Model(&models.TerminalRecall{}).Where("to_agent_id = ?", toAgentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&recalls).Error
	return recalls, total, err
}

func (r *GormTerminalRecallRepository) UpdateStatus(id int64, status int16, confirmedBy *int64) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"confirmed_at": &now,
	}
	if confirmedBy != nil {
		updates["confirmed_by"] = confirmedBy
	}
	return r.db.Model(&models.TerminalRecall{}).Where("id = ?", id).Updates(updates).Error
}

var _ TerminalRecallRepository = (*GormTerminalRecallRepository)(nil)

// TerminalImportRecordRepository 终端入库记录仓库接口
type TerminalImportRecordRepository interface {
	Create(record *models.TerminalImportRecord) error
	FindByID(id int64) (*models.TerminalImportRecord, error)
	FindByOwner(ownerAgentID int64, limit, offset int) ([]*models.TerminalImportRecord, int64, error)
}

// GormTerminalImportRecordRepository GORM实现
type GormTerminalImportRecordRepository struct {
	db *gorm.DB
}

func NewGormTerminalImportRecordRepository(db *gorm.DB) *GormTerminalImportRecordRepository {
	return &GormTerminalImportRecordRepository{db: db}
}

func (r *GormTerminalImportRecordRepository) Create(record *models.TerminalImportRecord) error {
	return r.db.Create(record).Error
}

func (r *GormTerminalImportRecordRepository) FindByID(id int64) (*models.TerminalImportRecord, error) {
	var record models.TerminalImportRecord
	err := r.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *GormTerminalImportRecordRepository) FindByOwner(ownerAgentID int64, limit, offset int) ([]*models.TerminalImportRecord, int64, error) {
	var records []*models.TerminalImportRecord
	var total int64

	query := r.db.Model(&models.TerminalImportRecord{}).Where("owner_agent_id = ?", ownerAgentID)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error
	return records, total, err
}

var _ TerminalImportRecordRepository = (*GormTerminalImportRecordRepository)(nil)

// CountActivatedByDate 按日期统计激活终端数量
func (r *GormTerminalRepository) CountActivatedByDate(ownerAgentID int64, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Terminal{}).
		Where("owner_agent_id = ? AND status = ? AND activated_at >= ? AND activated_at < ?",
			ownerAgentID, models.TerminalStatusActivated, startDate, endDate).
		Count(&count).Error
	return count, err
}

// BatchCreate 批量创建终端
func (r *GormTerminalRepository) BatchCreate(terminals []*models.Terminal) error {
	if len(terminals) == 0 {
		return nil
	}
	return r.db.CreateInBatches(terminals, 100).Error
}

// FindBySNs 批量查询终端
func (r *GormTerminalRepository) FindBySNs(sns []string) ([]*models.Terminal, error) {
	var terminals []*models.Terminal
	err := r.db.Where("terminal_sn IN ?", sns).Find(&terminals).Error
	return terminals, err
}

// ========== 终端政策相关 ==========

// FindPolicyBySN 查询终端政策
func (r *GormTerminalRepository) FindPolicyBySN(sn string) (*models.TerminalPolicy, error) {
	var policy models.TerminalPolicy
	err := r.db.Where("terminal_sn = ?", sn).First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// SavePolicy 保存终端政策（创建或更新）
func (r *GormTerminalRepository) SavePolicy(policy *models.TerminalPolicy) error {
	if policy.ID == 0 {
		return r.db.Create(policy).Error
	}
	return r.db.Save(policy).Error
}

// FindPoliciesByAgent 查询代理商的所有终端政策
func (r *GormTerminalRepository) FindPoliciesByAgent(agentID int64, limit, offset int) ([]*models.TerminalPolicy, int64, error) {
	var policies []*models.TerminalPolicy
	var total int64

	query := r.db.Model(&models.TerminalPolicy{}).Where("agent_id = ?", agentID)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&policies).Error
	return policies, total, err
}

// FindUnsyncedPolicies 查询未同步的政策
func (r *GormTerminalRepository) FindUnsyncedPolicies(limit int) ([]*models.TerminalPolicy, error) {
	var policies []*models.TerminalPolicy
	err := r.db.Where("is_synced = ?", false).
		Order("updated_at ASC").
		Limit(limit).
		Find(&policies).Error
	return policies, err
}

// UpdatePolicySyncStatus 更新政策同步状态
func (r *GormTerminalRepository) UpdatePolicySyncStatus(id int64, isSynced bool, syncError string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_synced":  isSynced,
		"sync_error": syncError,
	}
	if isSynced {
		updates["synced_at"] = &now
	}
	return r.db.Model(&models.TerminalPolicy{}).Where("id = ?", id).Updates(updates).Error
}
