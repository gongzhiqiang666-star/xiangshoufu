package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// ============================================================
// 奖励政策模版仓库接口
// ============================================================

// RewardPolicyTemplateRepository 奖励政策模版仓库接口
type RewardPolicyTemplateRepository interface {
	Create(template *models.RewardPolicyTemplate) error
	Update(template *models.RewardPolicyTemplate) error
	FindByID(id int64) (*models.RewardPolicyTemplate, error)
	FindAll(enabled *bool, limit, offset int) ([]*models.RewardPolicyTemplate, int64, error)
	UpdateEnabled(id int64, enabled bool) error
	Delete(id int64) error
}

// RewardStageRepository 奖励阶段仓库接口
type RewardStageRepository interface {
	Create(stage *models.RewardStage) error
	BatchCreate(stages []*models.RewardStage) error
	FindByTemplateID(templateID int64) ([]*models.RewardStage, error)
	DeleteByTemplateID(templateID int64) error
}

// AgentRewardRateRepository 代理商奖励比例仓库接口
type AgentRewardRateRepository interface {
	Create(rate *models.AgentRewardRate) error
	Update(rate *models.AgentRewardRate) error
	Upsert(rate *models.AgentRewardRate) error
	FindByAgentID(agentID int64) (*models.AgentRewardRate, error)
	FindByAgentIDs(agentIDs []int64) ([]*models.AgentRewardRate, error)
	Delete(agentID int64) error
}

// TerminalRewardProgressRepository 终端奖励进度仓库接口
type TerminalRewardProgressRepository interface {
	Create(progress *models.TerminalRewardProgress) error
	Update(progress *models.TerminalRewardProgress) error
	FindByID(id int64) (*models.TerminalRewardProgress, error)
	FindByTerminalSN(terminalSN string) ([]*models.TerminalRewardProgress, error)
	FindActiveByTerminalSN(terminalSN string) (*models.TerminalRewardProgress, error)
	FindActiveByBindTime(beforeTime time.Time, limit, offset int) ([]*models.TerminalRewardProgress, error)
	UpdateStatus(id int64, status models.RewardProgressStatus) error
	Terminate(id int64) error
}

// TerminalStageRewardRepository 终端阶段奖励仓库接口
type TerminalStageRewardRepository interface {
	Create(reward *models.TerminalStageReward) error
	BatchCreate(rewards []*models.TerminalStageReward) error
	Update(reward *models.TerminalStageReward) error
	FindByID(id int64) (*models.TerminalStageReward, error)
	FindByProgressID(progressID int64) ([]*models.TerminalStageReward, error)
	FindPendingByStageEnd(beforeTime time.Time, limit int) ([]*models.TerminalStageReward, error)
	UpdateStatus(id int64, status models.StageRewardStatus) error
	UpdateActualValue(id int64, actualValue int64) error
}

// RewardDistributionRepository 奖励发放记录仓库接口
type RewardDistributionRepository interface {
	Create(distribution *models.RewardDistribution) error
	BatchCreate(distributions []*models.RewardDistribution) error
	FindByStageRewardID(stageRewardID int64) ([]*models.RewardDistribution, error)
	FindByAgentID(agentID int64, limit, offset int) ([]*models.RewardDistribution, int64, error)
	FindPendingWallet(limit int) ([]*models.RewardDistribution, error)
	UpdateWalletStatus(id int64, walletRecordID int64, status int16) error
}

// RewardOverflowLogRepository 奖励溢出日志仓库接口
type RewardOverflowLogRepository interface {
	Create(log *models.RewardOverflowLog) error
	FindUnresolved(limit, offset int) ([]*models.RewardOverflowLog, int64, error)
	Resolve(id int64, resolvedBy string) error
}

// ============================================================
// 奖励政策模版仓库实现
// ============================================================

// GormRewardPolicyTemplateRepository GORM实现
type GormRewardPolicyTemplateRepository struct {
	db *gorm.DB
}

// NewGormRewardPolicyTemplateRepository 创建仓库
func NewGormRewardPolicyTemplateRepository(db *gorm.DB) *GormRewardPolicyTemplateRepository {
	return &GormRewardPolicyTemplateRepository{db: db}
}

func (r *GormRewardPolicyTemplateRepository) Create(template *models.RewardPolicyTemplate) error {
	return r.db.Create(template).Error
}

func (r *GormRewardPolicyTemplateRepository) Update(template *models.RewardPolicyTemplate) error {
	template.UpdatedAt = time.Now()
	return r.db.Save(template).Error
}

func (r *GormRewardPolicyTemplateRepository) FindByID(id int64) (*models.RewardPolicyTemplate, error) {
	var template models.RewardPolicyTemplate
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *GormRewardPolicyTemplateRepository) FindAll(enabled *bool, limit, offset int) ([]*models.RewardPolicyTemplate, int64, error) {
	var templates []*models.RewardPolicyTemplate
	var total int64

	query := r.db.Model(&models.RewardPolicyTemplate{})
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (r *GormRewardPolicyTemplateRepository) UpdateEnabled(id int64, enabled bool) error {
	return r.db.Model(&models.RewardPolicyTemplate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"enabled":    enabled,
			"updated_at": time.Now(),
		}).Error
}

func (r *GormRewardPolicyTemplateRepository) Delete(id int64) error {
	return r.db.Delete(&models.RewardPolicyTemplate{}, id).Error
}

var _ RewardPolicyTemplateRepository = (*GormRewardPolicyTemplateRepository)(nil)

// ============================================================
// 奖励阶段仓库实现
// ============================================================

// GormRewardStageRepository GORM实现
type GormRewardStageRepository struct {
	db *gorm.DB
}

// NewGormRewardStageRepository 创建仓库
func NewGormRewardStageRepository(db *gorm.DB) *GormRewardStageRepository {
	return &GormRewardStageRepository{db: db}
}

func (r *GormRewardStageRepository) Create(stage *models.RewardStage) error {
	return r.db.Create(stage).Error
}

func (r *GormRewardStageRepository) BatchCreate(stages []*models.RewardStage) error {
	if len(stages) == 0 {
		return nil
	}
	return r.db.Create(&stages).Error
}

func (r *GormRewardStageRepository) FindByTemplateID(templateID int64) ([]*models.RewardStage, error) {
	var stages []*models.RewardStage
	err := r.db.Where("template_id = ?", templateID).
		Order("stage_order ASC").
		Find(&stages).Error
	return stages, err
}

func (r *GormRewardStageRepository) DeleteByTemplateID(templateID int64) error {
	return r.db.Where("template_id = ?", templateID).Delete(&models.RewardStage{}).Error
}

var _ RewardStageRepository = (*GormRewardStageRepository)(nil)

// ============================================================
// 代理商奖励比例仓库实现
// ============================================================

// GormAgentRewardRateRepository GORM实现
type GormAgentRewardRateRepository struct {
	db *gorm.DB
}

// NewGormAgentRewardRateRepository 创建仓库
func NewGormAgentRewardRateRepository(db *gorm.DB) *GormAgentRewardRateRepository {
	return &GormAgentRewardRateRepository{db: db}
}

func (r *GormAgentRewardRateRepository) Create(rate *models.AgentRewardRate) error {
	return r.db.Create(rate).Error
}

func (r *GormAgentRewardRateRepository) Update(rate *models.AgentRewardRate) error {
	rate.UpdatedAt = time.Now()
	return r.db.Save(rate).Error
}

func (r *GormAgentRewardRateRepository) Upsert(rate *models.AgentRewardRate) error {
	return r.db.Where("agent_id = ?", rate.AgentID).
		Assign(map[string]interface{}{
			"reward_rate": rate.RewardRate,
			"updated_at":  time.Now(),
		}).FirstOrCreate(rate).Error
}

func (r *GormAgentRewardRateRepository) FindByAgentID(agentID int64) (*models.AgentRewardRate, error) {
	var rate models.AgentRewardRate
	if err := r.db.Where("agent_id = ?", agentID).First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *GormAgentRewardRateRepository) FindByAgentIDs(agentIDs []int64) ([]*models.AgentRewardRate, error) {
	if len(agentIDs) == 0 {
		return nil, nil
	}
	var rates []*models.AgentRewardRate
	err := r.db.Where("agent_id IN ?", agentIDs).Find(&rates).Error
	return rates, err
}

func (r *GormAgentRewardRateRepository) Delete(agentID int64) error {
	return r.db.Where("agent_id = ?", agentID).Delete(&models.AgentRewardRate{}).Error
}

var _ AgentRewardRateRepository = (*GormAgentRewardRateRepository)(nil)

// ============================================================
// 终端奖励进度仓库实现
// ============================================================

// GormTerminalRewardProgressRepository GORM实现
type GormTerminalRewardProgressRepository struct {
	db *gorm.DB
}

// NewGormTerminalRewardProgressRepository 创建仓库
func NewGormTerminalRewardProgressRepository(db *gorm.DB) *GormTerminalRewardProgressRepository {
	return &GormTerminalRewardProgressRepository{db: db}
}

func (r *GormTerminalRewardProgressRepository) Create(progress *models.TerminalRewardProgress) error {
	return r.db.Create(progress).Error
}

func (r *GormTerminalRewardProgressRepository) Update(progress *models.TerminalRewardProgress) error {
	progress.UpdatedAt = time.Now()
	return r.db.Save(progress).Error
}

func (r *GormTerminalRewardProgressRepository) FindByID(id int64) (*models.TerminalRewardProgress, error) {
	var progress models.TerminalRewardProgress
	if err := r.db.First(&progress, id).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *GormTerminalRewardProgressRepository) FindByTerminalSN(terminalSN string) ([]*models.TerminalRewardProgress, error) {
	var progresses []*models.TerminalRewardProgress
	err := r.db.Where("terminal_sn = ?", terminalSN).
		Order("created_at DESC").
		Find(&progresses).Error
	return progresses, err
}

func (r *GormTerminalRewardProgressRepository) FindActiveByTerminalSN(terminalSN string) (*models.TerminalRewardProgress, error) {
	var progress models.TerminalRewardProgress
	err := r.db.Where("terminal_sn = ? AND status = ?", terminalSN, models.RewardProgressStatusActive).
		First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *GormTerminalRewardProgressRepository) FindActiveByBindTime(beforeTime time.Time, limit, offset int) ([]*models.TerminalRewardProgress, error) {
	var progresses []*models.TerminalRewardProgress
	err := r.db.Where("status = ? AND bind_time <= ?", models.RewardProgressStatusActive, beforeTime).
		Order("bind_time ASC").
		Limit(limit).Offset(offset).
		Find(&progresses).Error
	return progresses, err
}

func (r *GormTerminalRewardProgressRepository) UpdateStatus(id int64, status models.RewardProgressStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if status == models.RewardProgressStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return r.db.Model(&models.TerminalRewardProgress{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *GormTerminalRewardProgressRepository) Terminate(id int64) error {
	now := time.Now()
	return r.db.Model(&models.TerminalRewardProgress{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        models.RewardProgressStatusTerminated,
			"terminated_at": &now,
			"updated_at":    now,
		}).Error
}

var _ TerminalRewardProgressRepository = (*GormTerminalRewardProgressRepository)(nil)

// ============================================================
// 终端阶段奖励仓库实现
// ============================================================

// GormTerminalStageRewardRepository GORM实现
type GormTerminalStageRewardRepository struct {
	db *gorm.DB
}

// NewGormTerminalStageRewardRepository 创建仓库
func NewGormTerminalStageRewardRepository(db *gorm.DB) *GormTerminalStageRewardRepository {
	return &GormTerminalStageRewardRepository{db: db}
}

func (r *GormTerminalStageRewardRepository) Create(reward *models.TerminalStageReward) error {
	return r.db.Create(reward).Error
}

func (r *GormTerminalStageRewardRepository) BatchCreate(rewards []*models.TerminalStageReward) error {
	if len(rewards) == 0 {
		return nil
	}
	return r.db.Create(&rewards).Error
}

func (r *GormTerminalStageRewardRepository) Update(reward *models.TerminalStageReward) error {
	reward.UpdatedAt = time.Now()
	return r.db.Save(reward).Error
}

func (r *GormTerminalStageRewardRepository) FindByID(id int64) (*models.TerminalStageReward, error) {
	var reward models.TerminalStageReward
	if err := r.db.First(&reward, id).Error; err != nil {
		return nil, err
	}
	return &reward, nil
}

func (r *GormTerminalStageRewardRepository) FindByProgressID(progressID int64) ([]*models.TerminalStageReward, error) {
	var rewards []*models.TerminalStageReward
	err := r.db.Where("progress_id = ?", progressID).
		Order("stage_order ASC").
		Find(&rewards).Error
	return rewards, err
}

func (r *GormTerminalStageRewardRepository) FindPendingByStageEnd(beforeTime time.Time, limit int) ([]*models.TerminalStageReward, error) {
	var rewards []*models.TerminalStageReward
	err := r.db.Where("status = ? AND stage_end < ?", models.StageRewardStatusPending, beforeTime).
		Order("stage_end ASC").
		Limit(limit).
		Find(&rewards).Error
	return rewards, err
}

func (r *GormTerminalStageRewardRepository) UpdateStatus(id int64, status models.StageRewardStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if status == models.StageRewardStatusSettled {
		now := time.Now()
		updates["settled_at"] = &now
	}
	return r.db.Model(&models.TerminalStageReward{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *GormTerminalStageRewardRepository) UpdateActualValue(id int64, actualValue int64) error {
	return r.db.Model(&models.TerminalStageReward{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"actual_value": actualValue,
			"updated_at":   time.Now(),
		}).Error
}

var _ TerminalStageRewardRepository = (*GormTerminalStageRewardRepository)(nil)

// ============================================================
// 奖励发放记录仓库实现
// ============================================================

// GormRewardDistributionRepository GORM实现
type GormRewardDistributionRepository struct {
	db *gorm.DB
}

// NewGormRewardDistributionRepository 创建仓库
func NewGormRewardDistributionRepository(db *gorm.DB) *GormRewardDistributionRepository {
	return &GormRewardDistributionRepository{db: db}
}

func (r *GormRewardDistributionRepository) Create(distribution *models.RewardDistribution) error {
	return r.db.Create(distribution).Error
}

func (r *GormRewardDistributionRepository) BatchCreate(distributions []*models.RewardDistribution) error {
	if len(distributions) == 0 {
		return nil
	}
	return r.db.Create(&distributions).Error
}

func (r *GormRewardDistributionRepository) FindByStageRewardID(stageRewardID int64) ([]*models.RewardDistribution, error) {
	var distributions []*models.RewardDistribution
	err := r.db.Where("stage_reward_id = ?", stageRewardID).
		Order("agent_level ASC").
		Find(&distributions).Error
	return distributions, err
}

func (r *GormRewardDistributionRepository) FindByAgentID(agentID int64, limit, offset int) ([]*models.RewardDistribution, int64, error) {
	var distributions []*models.RewardDistribution
	var total int64

	query := r.db.Model(&models.RewardDistribution{}).Where("agent_id = ?", agentID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&distributions).Error; err != nil {
		return nil, 0, err
	}

	return distributions, total, nil
}

func (r *GormRewardDistributionRepository) FindPendingWallet(limit int) ([]*models.RewardDistribution, error) {
	var distributions []*models.RewardDistribution
	err := r.db.Where("wallet_status = 0").
		Order("created_at ASC").
		Limit(limit).
		Find(&distributions).Error
	return distributions, err
}

func (r *GormRewardDistributionRepository) UpdateWalletStatus(id int64, walletRecordID int64, status int16) error {
	return r.db.Model(&models.RewardDistribution{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"wallet_record_id": walletRecordID,
			"wallet_status":    status,
		}).Error
}

var _ RewardDistributionRepository = (*GormRewardDistributionRepository)(nil)

// ============================================================
// 奖励溢出日志仓库实现
// ============================================================

// GormRewardOverflowLogRepository GORM实现
type GormRewardOverflowLogRepository struct {
	db *gorm.DB
}

// NewGormRewardOverflowLogRepository 创建仓库
func NewGormRewardOverflowLogRepository(db *gorm.DB) *GormRewardOverflowLogRepository {
	return &GormRewardOverflowLogRepository{db: db}
}

func (r *GormRewardOverflowLogRepository) Create(log *models.RewardOverflowLog) error {
	return r.db.Create(log).Error
}

func (r *GormRewardOverflowLogRepository) FindUnresolved(limit, offset int) ([]*models.RewardOverflowLog, int64, error) {
	var logs []*models.RewardOverflowLog
	var total int64

	query := r.db.Model(&models.RewardOverflowLog{}).Where("resolved = false")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *GormRewardOverflowLogRepository) Resolve(id int64, resolvedBy string) error {
	now := time.Now()
	return r.db.Model(&models.RewardOverflowLog{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"resolved":    true,
			"resolved_at": &now,
			"resolved_by": resolvedBy,
		}).Error
}

var _ RewardOverflowLogRepository = (*GormRewardOverflowLogRepository)(nil)
