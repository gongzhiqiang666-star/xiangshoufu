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
	CountByTypeCode(channelID int64, brandCode, modelCode string) (int64, error)
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

// CountByTypeCode 统计指定终端类型的终端数量
func (r *GormTerminalRepository) CountByTypeCode(channelID int64, brandCode, modelCode string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Terminal{}).
		Where("channel_id = ? AND brand_code = ? AND model_code = ?", channelID, brandCode, modelCode).
		Count(&count).Error
	return count, err
}

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

// TerminalFilterParams 终端列表筛选参数
type TerminalFilterParams struct {
	OwnerAgentID int64    // 所属代理商ID（必填）
	ChannelID    *int64   // 通道ID（可选）
	BrandCode    string   // 品牌编码（可选）
	ModelCode    string   // 型号编码（可选）
	StatusGroup  string   // 状态分组（可选）：all/unstock/stocked/unbound/inactive/active
	Keyword      string   // 搜索关键词（可选）：终端SN或商户号
	Limit        int
	Offset       int
}

// FindByOwnerWithFilter 带筛选条件查询终端列表
func (r *GormTerminalRepository) FindByOwnerWithFilter(params TerminalFilterParams) ([]*models.Terminal, int64, error) {
	var terminals []*models.Terminal
	var total int64

	query := r.db.Model(&models.Terminal{}).Where("owner_agent_id = ?", params.OwnerAgentID)

	// 通道筛选
	if params.ChannelID != nil && *params.ChannelID > 0 {
		query = query.Where("channel_id = ?", *params.ChannelID)
	}

	// 品牌筛选
	if params.BrandCode != "" {
		query = query.Where("brand_code = ?", params.BrandCode)
	}

	// 型号筛选
	if params.ModelCode != "" {
		query = query.Where("model_code = ?", params.ModelCode)
	}

	// 状态分组筛选
	switch params.StatusGroup {
	case "unstock":
		// 未出库: Status=1
		query = query.Where("status = ?", models.TerminalStatusPending)
	case "stocked":
		// 已出库: Status=2
		query = query.Where("status = ?", models.TerminalStatusAllocated)
	case "unbound":
		// 未绑定: Status=2 且 MerchantID=null
		query = query.Where("status = ? AND merchant_id IS NULL", models.TerminalStatusAllocated)
	case "inactive":
		// 未激活: Status=3 且 ActivatedAt=null
		query = query.Where("status = ? AND activated_at IS NULL", models.TerminalStatusBound)
	case "active":
		// 已激活: Status=4
		query = query.Where("status = ?", models.TerminalStatusActivated)
	case "all", "":
		// 全部，不添加状态筛选
	}

	// 关键词搜索
	if params.Keyword != "" {
		query = query.Where("terminal_sn LIKE ? OR merchant_no LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&terminals).Error
	return terminals, total, err
}

// TerminalTypeInfo 终端类型信息
type TerminalTypeInfo struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelCode string `json:"channel_code"`
	BrandCode   string `json:"brand_code"`
	ModelCode   string `json:"model_code"`
	Count       int64  `json:"count"`
}

// GetTerminalTypes 获取代理商拥有的终端类型列表
func (r *GormTerminalRepository) GetTerminalTypes(ownerAgentID int64) ([]TerminalTypeInfo, error) {
	var types []TerminalTypeInfo
	err := r.db.Model(&models.Terminal{}).
		Select("channel_id, channel_code, brand_code, model_code, COUNT(*) as count").
		Where("owner_agent_id = ?", ownerAgentID).
		Group("channel_id, channel_code, brand_code, model_code").
		Order("channel_id, brand_code, model_code").
		Find(&types).Error
	return types, err
}

// StatusGroupCount 状态分组统计
type StatusGroupCount struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}

// GetStatusGroupCounts 获取状态分组统计
func (r *GormTerminalRepository) GetStatusGroupCounts(ownerAgentID int64, channelID *int64, brandCode, modelCode string) ([]StatusGroupCount, error) {
	var counts []StatusGroupCount

	baseQuery := r.db.Model(&models.Terminal{}).Where("owner_agent_id = ?", ownerAgentID)
	if channelID != nil && *channelID > 0 {
		baseQuery = baseQuery.Where("channel_id = ?", *channelID)
	}
	if brandCode != "" {
		baseQuery = baseQuery.Where("brand_code = ?", brandCode)
	}
	if modelCode != "" {
		baseQuery = baseQuery.Where("model_code = ?", modelCode)
	}

	// 全部
	var totalCount int64
	baseQuery.Count(&totalCount)
	counts = append(counts, StatusGroupCount{Key: "all", Label: "全部", Count: totalCount})

	// 未出库
	var unstockCount int64
	baseQuery.Where("status = ?", models.TerminalStatusPending).Count(&unstockCount)
	counts = append(counts, StatusGroupCount{Key: "unstock", Label: "未出库", Count: unstockCount})

	// 已出库
	var stockedCount int64
	r.db.Model(&models.Terminal{}).Where("owner_agent_id = ? AND status = ?", ownerAgentID, models.TerminalStatusAllocated).Count(&stockedCount)
	counts = append(counts, StatusGroupCount{Key: "stocked", Label: "已出库", Count: stockedCount})

	// 未绑定
	var unboundCount int64
	r.db.Model(&models.Terminal{}).Where("owner_agent_id = ? AND status = ? AND merchant_id IS NULL", ownerAgentID, models.TerminalStatusAllocated).Count(&unboundCount)
	counts = append(counts, StatusGroupCount{Key: "unbound", Label: "未绑定", Count: unboundCount})

	// 未激活
	var inactiveCount int64
	r.db.Model(&models.Terminal{}).Where("owner_agent_id = ? AND status = ? AND activated_at IS NULL", ownerAgentID, models.TerminalStatusBound).Count(&inactiveCount)
	counts = append(counts, StatusGroupCount{Key: "inactive", Label: "未激活", Count: inactiveCount})

	// 已激活
	var activeCount int64
	r.db.Model(&models.Terminal{}).Where("owner_agent_id = ? AND status = ?", ownerAgentID, models.TerminalStatusActivated).Count(&activeCount)
	counts = append(counts, StatusGroupCount{Key: "active", Label: "已激活", Count: activeCount})

	return counts, nil
}

// TerminalFlowLog 终端流动记录
type TerminalFlowLog struct {
	ID            int64      `json:"id"`
	LogType       string     `json:"log_type"`        // distribute/recall/bind/unbind/activate
	LogTypeName   string     `json:"log_type_name"`   // 下发/回拨/绑定/解绑/激活
	FromAgentID   *int64     `json:"from_agent_id"`
	FromAgentName string     `json:"from_agent_name"`
	ToAgentID     *int64     `json:"to_agent_id"`
	ToAgentName   string     `json:"to_agent_name"`
	MerchantNo    string     `json:"merchant_no"`
	Status        int16      `json:"status"`
	StatusName    string     `json:"status_name"`
	Remark        string     `json:"remark"`
	CreatedAt     time.Time  `json:"created_at"`
	ConfirmedAt   *time.Time `json:"confirmed_at"`
}

// GetTerminalFlowLogs 获取终端流动记录
func (r *GormTerminalRepository) GetTerminalFlowLogs(terminalSN string, logType string, limit, offset int) ([]TerminalFlowLog, int64, error) {
	var logs []TerminalFlowLog
	var total int64

	// 构建UNION ALL查询
	// 1. 下发记录
	distributeQuery := r.db.Model(&models.TerminalDistribute{}).
		Select(`id, 'distribute' as log_type, '下发' as log_type_name,
			from_agent_id, to_agent_id, '' as merchant_no, status, remark,
			created_at, confirmed_at`).
		Where("terminal_sn = ?", terminalSN)

	// 2. 回拨记录
	recallQuery := r.db.Model(&models.TerminalRecall{}).
		Select(`id, 'recall' as log_type, '回拨' as log_type_name,
			from_agent_id, to_agent_id, '' as merchant_no, status, remark,
			created_at, confirmed_at`).
		Where("terminal_sn = ?", terminalSN)

	// 根据日志类型筛选
	switch logType {
	case "distribute":
		err := distributeQuery.Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = distributeQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
		return logs, total, err
	case "recall":
		err := recallQuery.Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = recallQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
		return logs, total, err
	default:
		// 查询全部类型，使用子查询合并
		// 先分别查询再合并
		var distributeLogs []TerminalFlowLog
		var recallLogs []TerminalFlowLog

		r.db.Model(&models.TerminalDistribute{}).
			Select(`id, 'distribute' as log_type, from_agent_id, to_agent_id, '' as merchant_no, status, remark, created_at, confirmed_at`).
			Where("terminal_sn = ?", terminalSN).
			Find(&distributeLogs)

		r.db.Model(&models.TerminalRecall{}).
			Select(`id, 'recall' as log_type, from_agent_id, to_agent_id, '' as merchant_no, status, remark, created_at, confirmed_at`).
			Where("terminal_sn = ?", terminalSN).
			Find(&recallLogs)

		// 为下发记录设置类型名称和状态名称
		for i := range distributeLogs {
			distributeLogs[i].LogType = "distribute"
			distributeLogs[i].LogTypeName = "下发"
			distributeLogs[i].StatusName = getDistributeStatusName(distributeLogs[i].Status)
		}

		// 为回拨记录设置类型名称和状态名称
		for i := range recallLogs {
			recallLogs[i].LogType = "recall"
			recallLogs[i].LogTypeName = "回拨"
			recallLogs[i].StatusName = getRecallStatusName(recallLogs[i].Status)
		}

		// 合并并按时间排序
		logs = append(distributeLogs, recallLogs...)
		total = int64(len(logs))

		// 简单排序（按创建时间倒序）
		for i := 0; i < len(logs)-1; i++ {
			for j := i + 1; j < len(logs); j++ {
				if logs[i].CreatedAt.Before(logs[j].CreatedAt) {
					logs[i], logs[j] = logs[j], logs[i]
				}
			}
		}

		// 分页
		start := offset
		if start > len(logs) {
			start = len(logs)
		}
		end := start + limit
		if end > len(logs) {
			end = len(logs)
		}
		logs = logs[start:end]

		return logs, total, nil
	}
}

// getDistributeStatusName 获取下发状态名称
func getDistributeStatusName(status int16) string {
	switch status {
	case models.TerminalDistributeStatusPending:
		return "待确认"
	case models.TerminalDistributeStatusConfirmed:
		return "已确认"
	case models.TerminalDistributeStatusRejected:
		return "已拒绝"
	case models.TerminalDistributeStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// getRecallStatusName 获取回拨状态名称
func getRecallStatusName(status int16) string {
	switch status {
	case models.TerminalRecallStatusPending:
		return "待确认"
	case models.TerminalRecallStatusConfirmed:
		return "已确认"
	case models.TerminalRecallStatusRejected:
		return "已拒绝"
	case models.TerminalRecallStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// GetChannelList 获取代理商拥有终端的通道列表
func (r *GormTerminalRepository) GetChannelList(ownerAgentID int64) ([]struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelCode string `json:"channel_code"`
}, error) {
	var channels []struct {
		ChannelID   int64  `json:"channel_id"`
		ChannelCode string `json:"channel_code"`
	}
	err := r.db.Model(&models.Terminal{}).
		Select("DISTINCT channel_id, channel_code").
		Where("owner_agent_id = ?", ownerAgentID).
		Find(&channels).Error
	return channels, err
}
