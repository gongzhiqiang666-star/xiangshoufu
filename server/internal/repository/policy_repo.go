package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// ============================================================
// 接口定义
// ============================================================

// DepositCashbackPolicyRepository 押金返现政策仓库接口
type DepositCashbackPolicyRepository interface {
	Create(policy *models.DepositCashbackPolicy) error
	Update(policy *models.DepositCashbackPolicy) error
	FindByID(id int64) (*models.DepositCashbackPolicy, error)
	FindByTemplateID(templateID int64) ([]*models.DepositCashbackPolicy, error)
	FindByTemplateAndAmount(templateID int64, depositAmount int64) (*models.DepositCashbackPolicy, error)
	DeleteByTemplateID(templateID int64) error
}

// DepositCashbackRecordRepository 押金返现记录仓库接口
type DepositCashbackRecordRepository interface {
	Create(record *models.DepositCashbackRecord) error
	BatchCreate(records []*models.DepositCashbackRecord) error
	FindByTerminalID(terminalID int64) ([]*models.DepositCashbackRecord, error)
	FindPendingRecords(limit int) ([]*models.DepositCashbackRecord, error)
	UpdateWalletStatus(id int64, status int16) error
}

// ActivationRewardPolicyRepository 激活奖励政策仓库接口
type ActivationRewardPolicyRepository interface {
	Create(policy *models.ActivationRewardPolicy) error
	Update(policy *models.ActivationRewardPolicy) error
	FindByID(id int64) (*models.ActivationRewardPolicy, error)
	FindByTemplateID(templateID int64) ([]*models.ActivationRewardPolicy, error)
	FindActiveByChannelID(channelID int64) ([]*models.ActivationRewardPolicy, error)
	DeleteByTemplateID(templateID int64) error
}

// ActivationRewardRecordRepository 激活奖励记录仓库接口
type ActivationRewardRecordRepository interface {
	Create(record *models.ActivationRewardRecord) error
	BatchCreate(records []*models.ActivationRewardRecord) error
	FindByPolicyAndTerminal(policyID, terminalID int64, checkDate time.Time) (*models.ActivationRewardRecord, error)
	FindPendingRecords(limit int) ([]*models.ActivationRewardRecord, error)
	UpdateWalletStatus(id int64, status int16) error
}

// RateStagePolicyRepository 费率阶梯政策仓库接口
type RateStagePolicyRepository interface {
	Create(policy *models.RateStagePolicy) error
	Update(policy *models.RateStagePolicy) error
	FindByID(id int64) (*models.RateStagePolicy, error)
	FindByTemplateID(templateID int64) ([]*models.RateStagePolicy, error)
	FindActiveByChannelAndApplyTo(channelID int64, applyTo int16) ([]*models.RateStagePolicy, error)
	DeleteByTemplateID(templateID int64) error
}

// AgentDepositCashbackPolicyRepository 代理商押金返现政策仓库接口
type AgentDepositCashbackPolicyRepository interface {
	Create(policy *models.AgentDepositCashbackPolicy) error
	Update(policy *models.AgentDepositCashbackPolicy) error
	Upsert(policy *models.AgentDepositCashbackPolicy) error
	FindByAgentAndChannel(agentID, channelID int64) ([]*models.AgentDepositCashbackPolicy, error)
	FindByAgentChannelAndAmount(agentID, channelID, depositAmount int64) (*models.AgentDepositCashbackPolicy, error)
	DeleteByAgentAndChannel(agentID, channelID int64) error
}

// AgentSimCashbackPolicyRepository 代理商流量卡返现政策仓库接口
type AgentSimCashbackPolicyRepository interface {
	Create(policy *models.AgentSimCashbackPolicy) error
	Update(policy *models.AgentSimCashbackPolicy) error
	Upsert(policy *models.AgentSimCashbackPolicy) error
	FindByAgentAndChannel(agentID, channelID int64) (*models.AgentSimCashbackPolicy, error)
	DeleteByAgentAndChannel(agentID, channelID int64) error
}

// AgentActivationRewardPolicyRepository 代理商激活奖励政策仓库接口
type AgentActivationRewardPolicyRepository interface {
	Create(policy *models.AgentActivationRewardPolicy) error
	Update(policy *models.AgentActivationRewardPolicy) error
	FindByAgentAndChannel(agentID, channelID int64) ([]*models.AgentActivationRewardPolicy, error)
	DeleteByAgentAndChannel(agentID, channelID int64) error
}

// PolicyTemplateRepository 政策模板仓库接口
type PolicyTemplateRepository interface {
	Create(template *models.PolicyTemplateComplete) error
	Update(template *models.PolicyTemplateComplete) error
	FindByID(id int64) (*models.PolicyTemplateComplete, error)
	FindByChannelID(channelID int64, status []int16, limit, offset int) ([]*models.PolicyTemplateComplete, int64, error)
	FindDefaultByChannelID(channelID int64) (*models.PolicyTemplateComplete, error)
	UpdateStatus(id int64, status int16) error
}

// ============================================================
// 押金返现政策仓库实现
// ============================================================

// GormDepositCashbackPolicyRepository GORM实现
type GormDepositCashbackPolicyRepository struct {
	db *gorm.DB
}

// NewGormDepositCashbackPolicyRepository 创建仓库
func NewGormDepositCashbackPolicyRepository(db *gorm.DB) *GormDepositCashbackPolicyRepository {
	return &GormDepositCashbackPolicyRepository{db: db}
}

func (r *GormDepositCashbackPolicyRepository) Create(policy *models.DepositCashbackPolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormDepositCashbackPolicyRepository) Update(policy *models.DepositCashbackPolicy) error {
	policy.UpdatedAt = time.Now()
	return r.db.Save(policy).Error
}

func (r *GormDepositCashbackPolicyRepository) FindByID(id int64) (*models.DepositCashbackPolicy, error) {
	var policy models.DepositCashbackPolicy
	if err := r.db.First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormDepositCashbackPolicyRepository) FindByTemplateID(templateID int64) ([]*models.DepositCashbackPolicy, error) {
	var policies []*models.DepositCashbackPolicy
	err := r.db.Where("template_id = ? AND status = 1", templateID).
		Order("deposit_amount ASC").
		Find(&policies).Error
	return policies, err
}

func (r *GormDepositCashbackPolicyRepository) FindByTemplateAndAmount(templateID int64, depositAmount int64) (*models.DepositCashbackPolicy, error) {
	var policy models.DepositCashbackPolicy
	err := r.db.Where("template_id = ? AND deposit_amount = ? AND status = 1", templateID, depositAmount).
		First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormDepositCashbackPolicyRepository) DeleteByTemplateID(templateID int64) error {
	return r.db.Where("template_id = ?", templateID).Delete(&models.DepositCashbackPolicy{}).Error
}

var _ DepositCashbackPolicyRepository = (*GormDepositCashbackPolicyRepository)(nil)

// ============================================================
// 押金返现记录仓库实现
// ============================================================

// GormDepositCashbackRecordRepository GORM实现
type GormDepositCashbackRecordRepository struct {
	db *gorm.DB
}

// NewGormDepositCashbackRecordRepository 创建仓库
func NewGormDepositCashbackRecordRepository(db *gorm.DB) *GormDepositCashbackRecordRepository {
	return &GormDepositCashbackRecordRepository{db: db}
}

func (r *GormDepositCashbackRecordRepository) Create(record *models.DepositCashbackRecord) error {
	return r.db.Create(record).Error
}

func (r *GormDepositCashbackRecordRepository) BatchCreate(records []*models.DepositCashbackRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.Create(&records).Error
}

func (r *GormDepositCashbackRecordRepository) FindByTerminalID(terminalID int64) ([]*models.DepositCashbackRecord, error) {
	var records []*models.DepositCashbackRecord
	err := r.db.Where("terminal_id = ?", terminalID).
		Order("created_at DESC").
		Find(&records).Error
	return records, err
}

func (r *GormDepositCashbackRecordRepository) FindPendingRecords(limit int) ([]*models.DepositCashbackRecord, error) {
	var records []*models.DepositCashbackRecord
	err := r.db.Where("wallet_status = 0").
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

func (r *GormDepositCashbackRecordRepository) UpdateWalletStatus(id int64, status int16) error {
	now := time.Now()
	return r.db.Model(&models.DepositCashbackRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"wallet_status": status,
			"processed_at":  &now,
		}).Error
}

var _ DepositCashbackRecordRepository = (*GormDepositCashbackRecordRepository)(nil)

// ============================================================
// 激活奖励政策仓库实现
// ============================================================

// GormActivationRewardPolicyRepository GORM实现
type GormActivationRewardPolicyRepository struct {
	db *gorm.DB
}

// NewGormActivationRewardPolicyRepository 创建仓库
func NewGormActivationRewardPolicyRepository(db *gorm.DB) *GormActivationRewardPolicyRepository {
	return &GormActivationRewardPolicyRepository{db: db}
}

func (r *GormActivationRewardPolicyRepository) Create(policy *models.ActivationRewardPolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormActivationRewardPolicyRepository) Update(policy *models.ActivationRewardPolicy) error {
	policy.UpdatedAt = time.Now()
	return r.db.Save(policy).Error
}

func (r *GormActivationRewardPolicyRepository) FindByID(id int64) (*models.ActivationRewardPolicy, error) {
	var policy models.ActivationRewardPolicy
	if err := r.db.First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormActivationRewardPolicyRepository) FindByTemplateID(templateID int64) ([]*models.ActivationRewardPolicy, error) {
	var policies []*models.ActivationRewardPolicy
	err := r.db.Where("template_id = ? AND status = 1", templateID).
		Order("priority DESC, min_register_days ASC").
		Find(&policies).Error
	return policies, err
}

func (r *GormActivationRewardPolicyRepository) FindActiveByChannelID(channelID int64) ([]*models.ActivationRewardPolicy, error) {
	var policies []*models.ActivationRewardPolicy
	err := r.db.Where("channel_id = ? AND status = 1", channelID).
		Order("priority DESC").
		Find(&policies).Error
	return policies, err
}

func (r *GormActivationRewardPolicyRepository) DeleteByTemplateID(templateID int64) error {
	return r.db.Where("template_id = ?", templateID).Delete(&models.ActivationRewardPolicy{}).Error
}

var _ ActivationRewardPolicyRepository = (*GormActivationRewardPolicyRepository)(nil)

// ============================================================
// 激活奖励记录仓库实现
// ============================================================

// GormActivationRewardRecordRepository GORM实现
type GormActivationRewardRecordRepository struct {
	db *gorm.DB
}

// NewGormActivationRewardRecordRepository 创建仓库
func NewGormActivationRewardRecordRepository(db *gorm.DB) *GormActivationRewardRecordRepository {
	return &GormActivationRewardRecordRepository{db: db}
}

func (r *GormActivationRewardRecordRepository) Create(record *models.ActivationRewardRecord) error {
	return r.db.Create(record).Error
}

func (r *GormActivationRewardRecordRepository) BatchCreate(records []*models.ActivationRewardRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.Create(&records).Error
}

func (r *GormActivationRewardRecordRepository) FindByPolicyAndTerminal(policyID, terminalID int64, checkDate time.Time) (*models.ActivationRewardRecord, error) {
	var record models.ActivationRewardRecord
	err := r.db.Where("policy_id = ? AND terminal_id = ? AND check_date = ?", policyID, terminalID, checkDate).
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *GormActivationRewardRecordRepository) FindPendingRecords(limit int) ([]*models.ActivationRewardRecord, error) {
	var records []*models.ActivationRewardRecord
	err := r.db.Where("wallet_status = 0").
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

func (r *GormActivationRewardRecordRepository) UpdateWalletStatus(id int64, status int16) error {
	now := time.Now()
	return r.db.Model(&models.ActivationRewardRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"wallet_status": status,
			"processed_at":  &now,
		}).Error
}

var _ ActivationRewardRecordRepository = (*GormActivationRewardRecordRepository)(nil)

// ============================================================
// 费率阶梯政策仓库实现
// ============================================================

// GormRateStagePolicyRepository GORM实现
type GormRateStagePolicyRepository struct {
	db *gorm.DB
}

// NewGormRateStagePolicyRepository 创建仓库
func NewGormRateStagePolicyRepository(db *gorm.DB) *GormRateStagePolicyRepository {
	return &GormRateStagePolicyRepository{db: db}
}

func (r *GormRateStagePolicyRepository) Create(policy *models.RateStagePolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormRateStagePolicyRepository) Update(policy *models.RateStagePolicy) error {
	policy.UpdatedAt = time.Now()
	return r.db.Save(policy).Error
}

func (r *GormRateStagePolicyRepository) FindByID(id int64) (*models.RateStagePolicy, error) {
	var policy models.RateStagePolicy
	if err := r.db.First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormRateStagePolicyRepository) FindByTemplateID(templateID int64) ([]*models.RateStagePolicy, error) {
	var policies []*models.RateStagePolicy
	err := r.db.Where("template_id = ? AND status = 1", templateID).
		Order("priority DESC, min_days ASC").
		Find(&policies).Error
	return policies, err
}

func (r *GormRateStagePolicyRepository) FindActiveByChannelAndApplyTo(channelID int64, applyTo int16) ([]*models.RateStagePolicy, error) {
	var policies []*models.RateStagePolicy
	err := r.db.Where("channel_id = ? AND apply_to = ? AND status = 1", channelID, applyTo).
		Order("priority DESC, min_days ASC").
		Find(&policies).Error
	return policies, err
}

// FindApplicablePolicy 查找适用的费率阶梯政策（根据通道、应用对象和入网天数）
func (r *GormRateStagePolicyRepository) FindApplicablePolicy(channelID int64, applyTo int16, days int) (*models.RateStagePolicy, error) {
	var policy models.RateStagePolicy
	err := r.db.Where("channel_id = ? AND apply_to = ? AND status = 1 AND min_days <= ? AND (max_days >= ? OR max_days = -1)",
		channelID, applyTo, days, days).
		Order("priority DESC, min_days DESC").
		First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormRateStagePolicyRepository) DeleteByTemplateID(templateID int64) error {
	return r.db.Where("template_id = ?", templateID).Delete(&models.RateStagePolicy{}).Error
}

var _ RateStagePolicyRepository = (*GormRateStagePolicyRepository)(nil)

// ============================================================
// 代理商押金返现政策仓库实现
// ============================================================

// GormAgentDepositCashbackPolicyRepository GORM实现
type GormAgentDepositCashbackPolicyRepository struct {
	db *gorm.DB
}

// NewGormAgentDepositCashbackPolicyRepository 创建仓库
func NewGormAgentDepositCashbackPolicyRepository(db *gorm.DB) *GormAgentDepositCashbackPolicyRepository {
	return &GormAgentDepositCashbackPolicyRepository{db: db}
}

func (r *GormAgentDepositCashbackPolicyRepository) Create(policy *models.AgentDepositCashbackPolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormAgentDepositCashbackPolicyRepository) Update(policy *models.AgentDepositCashbackPolicy) error {
	policy.UpdatedAt = time.Now()
	return r.db.Save(policy).Error
}

func (r *GormAgentDepositCashbackPolicyRepository) Upsert(policy *models.AgentDepositCashbackPolicy) error {
	return r.db.Where("agent_id = ? AND channel_id = ? AND deposit_amount = ?", policy.AgentID, policy.ChannelID, policy.DepositAmount).
		Assign(map[string]interface{}{
			"cashback_amount": policy.CashbackAmount,
			"status":          policy.Status,
			"updated_at":      time.Now(),
		}).FirstOrCreate(policy).Error
}

func (r *GormAgentDepositCashbackPolicyRepository) FindByAgentAndChannel(agentID, channelID int64) ([]*models.AgentDepositCashbackPolicy, error) {
	var policies []*models.AgentDepositCashbackPolicy
	err := r.db.Where("agent_id = ? AND channel_id = ? AND status = 1", agentID, channelID).
		Order("deposit_amount ASC").
		Find(&policies).Error
	return policies, err
}

func (r *GormAgentDepositCashbackPolicyRepository) FindByAgentChannelAndAmount(agentID, channelID, depositAmount int64) (*models.AgentDepositCashbackPolicy, error) {
	var policy models.AgentDepositCashbackPolicy
	err := r.db.Where("agent_id = ? AND channel_id = ? AND deposit_amount = ? AND status = 1", agentID, channelID, depositAmount).
		First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormAgentDepositCashbackPolicyRepository) DeleteByAgentAndChannel(agentID, channelID int64) error {
	return r.db.Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Delete(&models.AgentDepositCashbackPolicy{}).Error
}

var _ AgentDepositCashbackPolicyRepository = (*GormAgentDepositCashbackPolicyRepository)(nil)

// ============================================================
// 代理商流量卡返现政策仓库实现
// ============================================================

// GormAgentSimCashbackPolicyRepository GORM实现
type GormAgentSimCashbackPolicyRepository struct {
	db *gorm.DB
}

// NewGormAgentSimCashbackPolicyRepository 创建仓库
func NewGormAgentSimCashbackPolicyRepository(db *gorm.DB) *GormAgentSimCashbackPolicyRepository {
	return &GormAgentSimCashbackPolicyRepository{db: db}
}

func (r *GormAgentSimCashbackPolicyRepository) Create(policy *models.AgentSimCashbackPolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormAgentSimCashbackPolicyRepository) Update(policy *models.AgentSimCashbackPolicy) error {
	policy.UpdatedAt = time.Now()
	return r.db.Save(policy).Error
}

func (r *GormAgentSimCashbackPolicyRepository) Upsert(policy *models.AgentSimCashbackPolicy) error {
	return r.db.Where("agent_id = ? AND channel_id = ? AND (brand_code = ? OR (brand_code IS NULL AND ? IS NULL))",
		policy.AgentID, policy.ChannelID, policy.BrandCode, policy.BrandCode).
		Assign(map[string]interface{}{
			"first_time_cashback":  policy.FirstTimeCashback,
			"second_time_cashback": policy.SecondTimeCashback,
			"third_plus_cashback":  policy.ThirdPlusCashback,
			"status":               policy.Status,
			"updated_at":           time.Now(),
		}).FirstOrCreate(policy).Error
}

func (r *GormAgentSimCashbackPolicyRepository) FindByAgentAndChannel(agentID, channelID int64) (*models.AgentSimCashbackPolicy, error) {
	var policy models.AgentSimCashbackPolicy
	err := r.db.Where("agent_id = ? AND channel_id = ? AND status = 1", agentID, channelID).
		First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *GormAgentSimCashbackPolicyRepository) DeleteByAgentAndChannel(agentID, channelID int64) error {
	return r.db.Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Delete(&models.AgentSimCashbackPolicy{}).Error
}

var _ AgentSimCashbackPolicyRepository = (*GormAgentSimCashbackPolicyRepository)(nil)

// ============================================================
// 代理商激活奖励政策仓库实现
// ============================================================

// GormAgentActivationRewardPolicyRepository GORM实现
type GormAgentActivationRewardPolicyRepository struct {
	db *gorm.DB
}

// NewGormAgentActivationRewardPolicyRepository 创建仓库
func NewGormAgentActivationRewardPolicyRepository(db *gorm.DB) *GormAgentActivationRewardPolicyRepository {
	return &GormAgentActivationRewardPolicyRepository{db: db}
}

func (r *GormAgentActivationRewardPolicyRepository) Create(policy *models.AgentActivationRewardPolicy) error {
	return r.db.Create(policy).Error
}

func (r *GormAgentActivationRewardPolicyRepository) Update(policy *models.AgentActivationRewardPolicy) error {
	policy.UpdatedAt = time.Now()
	return r.db.Save(policy).Error
}

func (r *GormAgentActivationRewardPolicyRepository) FindByAgentAndChannel(agentID, channelID int64) ([]*models.AgentActivationRewardPolicy, error) {
	var policies []*models.AgentActivationRewardPolicy
	err := r.db.Where("agent_id = ? AND channel_id = ? AND status = 1", agentID, channelID).
		Order("priority DESC, min_register_days ASC").
		Find(&policies).Error
	return policies, err
}

func (r *GormAgentActivationRewardPolicyRepository) DeleteByAgentAndChannel(agentID, channelID int64) error {
	return r.db.Where("agent_id = ? AND channel_id = ?", agentID, channelID).
		Delete(&models.AgentActivationRewardPolicy{}).Error
}

var _ AgentActivationRewardPolicyRepository = (*GormAgentActivationRewardPolicyRepository)(nil)

// ============================================================
// 政策模板仓库实现
// ============================================================

// GormPolicyTemplateRepository GORM实现
type GormPolicyTemplateRepository struct {
	db *gorm.DB
}

// NewGormPolicyTemplateRepository 创建仓库
func NewGormPolicyTemplateRepository(db *gorm.DB) *GormPolicyTemplateRepository {
	return &GormPolicyTemplateRepository{db: db}
}

func (r *GormPolicyTemplateRepository) Create(template *models.PolicyTemplateComplete) error {
	return r.db.Create(template).Error
}

func (r *GormPolicyTemplateRepository) Update(template *models.PolicyTemplateComplete) error {
	template.UpdatedAt = time.Now()
	return r.db.Save(template).Error
}

func (r *GormPolicyTemplateRepository) FindByID(id int64) (*models.PolicyTemplateComplete, error) {
	var template models.PolicyTemplateComplete
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *GormPolicyTemplateRepository) FindByChannelID(channelID int64, status []int16, limit, offset int) ([]*models.PolicyTemplateComplete, int64, error) {
	var templates []*models.PolicyTemplateComplete
	var total int64

	query := r.db.Model(&models.PolicyTemplateComplete{}).Where("channel_id = ?", channelID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("is_default DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (r *GormPolicyTemplateRepository) FindDefaultByChannelID(channelID int64) (*models.PolicyTemplateComplete, error) {
	var template models.PolicyTemplateComplete
	err := r.db.Where("channel_id = ? AND is_default = true AND status = 1", channelID).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *GormPolicyTemplateRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.PolicyTemplateComplete{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

var _ PolicyTemplateRepository = (*GormPolicyTemplateRepository)(nil)

// ============================================================
// SimCashbackPolicy 仓库扩展方法（补充到terminal_repo.go中的实现）
// ============================================================

// FindByTemplateID 根据模板ID查询流量卡返现政策
func (r *GormSimCashbackPolicyRepository) FindByTemplateID(templateID int64) ([]*models.SimCashbackPolicy, error) {
	var policies []*models.SimCashbackPolicy
	err := r.db.Where("template_id = ? AND status = 1", templateID).
		Find(&policies).Error
	return policies, err
}

// DeleteByTemplateID 根据模板ID删除流量卡返现政策
func (r *GormSimCashbackPolicyRepository) DeleteByTemplateID(templateID int64) error {
	return r.db.Where("template_id = ?", templateID).Delete(&models.SimCashbackPolicy{}).Error
}

// ============================================================
// ChannelRepository 通道仓库
// ============================================================

// ChannelRepository 通道仓库接口
type ChannelRepository interface {
	Create(channel *models.Channel) error
	Update(channel *models.Channel) error
	FindByID(id int64) (*models.Channel, error)
	FindByCode(code string) (*models.Channel, error)
	FindAll() ([]*models.Channel, error)
	FindAllActive() ([]*models.Channel, error)
	UpdateStatus(id int64, status int16) error
}

// GormChannelRepository GORM实现
type GormChannelRepository struct {
	db *gorm.DB
}

// NewGormChannelRepository 创建通道仓库
func NewGormChannelRepository(db *gorm.DB) *GormChannelRepository {
	return &GormChannelRepository{db: db}
}

func (r *GormChannelRepository) Create(channel *models.Channel) error {
	return r.db.Create(channel).Error
}

func (r *GormChannelRepository) Update(channel *models.Channel) error {
	return r.db.Save(channel).Error
}

func (r *GormChannelRepository) FindByID(id int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.Where("id = ?", id).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *GormChannelRepository) FindByCode(code string) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.Where("channel_code = ?", code).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *GormChannelRepository) FindAll() ([]*models.Channel, error) {
	var channels []*models.Channel
	err := r.db.Order("priority DESC").Find(&channels).Error
	return channels, err
}

func (r *GormChannelRepository) FindAllActive() ([]*models.Channel, error) {
	var channels []*models.Channel
	err := r.db.Where("status = 1").Order("priority DESC").Find(&channels).Error
	return channels, err
}

func (r *GormChannelRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.Channel{}).Where("id = ?", id).
		Update("status", status).Error
}

var _ ChannelRepository = (*GormChannelRepository)(nil)
