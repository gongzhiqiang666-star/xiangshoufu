package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// ChargingWalletRepository 充值钱包仓库接口
type ChargingWalletRepository interface {
	// 配置相关
	GetConfig(agentID int64) (*models.AgentWalletConfig, error)
	SaveConfig(config *models.AgentWalletConfig) error

	// 充值记录相关
	CreateDeposit(deposit *models.ChargingWalletDeposit) error
	GetDeposit(id int64) (*models.ChargingWalletDeposit, error)
	GetDepositByNo(depositNo string) (*models.ChargingWalletDeposit, error)
	UpdateDeposit(deposit *models.ChargingWalletDeposit) error
	GetDepositsByAgent(agentID int64, status *int16, limit, offset int) ([]*models.ChargingWalletDeposit, int64, error)
	GetPendingDeposits(limit, offset int) ([]*models.ChargingWalletDeposit, int64, error)

	// 奖励发放相关
	CreateReward(reward *models.ChargingWalletReward) error
	GetReward(id int64) (*models.ChargingWalletReward, error)
	GetRewardByNo(rewardNo string) (*models.ChargingWalletReward, error)
	UpdateReward(reward *models.ChargingWalletReward) error
	GetRewardsByFromAgent(fromAgentID int64, limit, offset int) ([]*models.ChargingWalletReward, int64, error)
	GetRewardsByToAgent(toAgentID int64, limit, offset int) ([]*models.ChargingWalletReward, int64, error)
}

// GormChargingWalletRepository GORM实现
type GormChargingWalletRepository struct {
	db *gorm.DB
}

// NewGormChargingWalletRepository 创建仓库
func NewGormChargingWalletRepository(db *gorm.DB) *GormChargingWalletRepository {
	return &GormChargingWalletRepository{db: db}
}

// ========== 配置相关 ==========

// GetConfig 获取代理商钱包配置
func (r *GormChargingWalletRepository) GetConfig(agentID int64) (*models.AgentWalletConfig, error) {
	var config models.AgentWalletConfig
	err := r.db.Where("agent_id = ?", agentID).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

// SaveConfig 保存代理商钱包配置
func (r *GormChargingWalletRepository) SaveConfig(config *models.AgentWalletConfig) error {
	if config.ID == 0 {
		return r.db.Create(config).Error
	}
	return r.db.Save(config).Error
}

// ========== 充值记录相关 ==========

// CreateDeposit 创建充值记录
func (r *GormChargingWalletRepository) CreateDeposit(deposit *models.ChargingWalletDeposit) error {
	return r.db.Create(deposit).Error
}

// GetDeposit 获取充值记录
func (r *GormChargingWalletRepository) GetDeposit(id int64) (*models.ChargingWalletDeposit, error) {
	var deposit models.ChargingWalletDeposit
	err := r.db.First(&deposit, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &deposit, err
}

// GetDepositByNo 根据单号获取充值记录
func (r *GormChargingWalletRepository) GetDepositByNo(depositNo string) (*models.ChargingWalletDeposit, error) {
	var deposit models.ChargingWalletDeposit
	err := r.db.Where("deposit_no = ?", depositNo).First(&deposit).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &deposit, err
}

// UpdateDeposit 更新充值记录
func (r *GormChargingWalletRepository) UpdateDeposit(deposit *models.ChargingWalletDeposit) error {
	deposit.UpdatedAt = time.Now()
	return r.db.Save(deposit).Error
}

// GetDepositsByAgent 获取代理商充值记录
func (r *GormChargingWalletRepository) GetDepositsByAgent(agentID int64, status *int16, limit, offset int) ([]*models.ChargingWalletDeposit, int64, error) {
	query := r.db.Model(&models.ChargingWalletDeposit{}).Where("agent_id = ?", agentID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	query.Count(&total)

	var deposits []*models.ChargingWalletDeposit
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&deposits).Error
	return deposits, total, err
}

// GetPendingDeposits 获取待确认的充值记录
func (r *GormChargingWalletRepository) GetPendingDeposits(limit, offset int) ([]*models.ChargingWalletDeposit, int64, error) {
	query := r.db.Model(&models.ChargingWalletDeposit{}).Where("status = ?", models.ChargingDepositStatusPending)

	var total int64
	query.Count(&total)

	var deposits []*models.ChargingWalletDeposit
	err := query.Order("created_at ASC").Limit(limit).Offset(offset).Find(&deposits).Error
	return deposits, total, err
}

// ========== 奖励发放相关 ==========

// CreateReward 创建奖励记录
func (r *GormChargingWalletRepository) CreateReward(reward *models.ChargingWalletReward) error {
	return r.db.Create(reward).Error
}

// GetReward 获取奖励记录
func (r *GormChargingWalletRepository) GetReward(id int64) (*models.ChargingWalletReward, error) {
	var reward models.ChargingWalletReward
	err := r.db.First(&reward, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &reward, err
}

// GetRewardByNo 根据单号获取奖励记录
func (r *GormChargingWalletRepository) GetRewardByNo(rewardNo string) (*models.ChargingWalletReward, error) {
	var reward models.ChargingWalletReward
	err := r.db.Where("reward_no = ?", rewardNo).First(&reward).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &reward, err
}

// UpdateReward 更新奖励记录
func (r *GormChargingWalletRepository) UpdateReward(reward *models.ChargingWalletReward) error {
	return r.db.Save(reward).Error
}

// GetRewardsByFromAgent 获取代理商发放的奖励记录
func (r *GormChargingWalletRepository) GetRewardsByFromAgent(fromAgentID int64, limit, offset int) ([]*models.ChargingWalletReward, int64, error) {
	query := r.db.Model(&models.ChargingWalletReward{}).Where("from_agent_id = ?", fromAgentID)

	var total int64
	query.Count(&total)

	var rewards []*models.ChargingWalletReward
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rewards).Error
	return rewards, total, err
}

// GetRewardsByToAgent 获取代理商收到的奖励记录
func (r *GormChargingWalletRepository) GetRewardsByToAgent(toAgentID int64, limit, offset int) ([]*models.ChargingWalletReward, int64, error) {
	query := r.db.Model(&models.ChargingWalletReward{}).Where("to_agent_id = ?", toAgentID)

	var total int64
	query.Count(&total)

	var rewards []*models.ChargingWalletReward
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rewards).Error
	return rewards, total, err
}

// GetTotalRewardsIssuedByAgent 获取代理商发放的奖励总额
func (r *GormChargingWalletRepository) GetTotalRewardsIssuedByAgent(agentID int64) (int64, error) {
	var total int64
	err := r.db.Model(&models.ChargingWalletReward{}).
		Where("from_agent_id = ? AND status = ?", agentID, models.ChargingRewardStatusIssued).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

var _ ChargingWalletRepository = (*GormChargingWalletRepository)(nil)
