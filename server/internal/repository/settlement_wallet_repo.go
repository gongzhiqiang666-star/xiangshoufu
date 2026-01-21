package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// SettlementWalletRepository 沉淀钱包仓库接口
type SettlementWalletRepository interface {
	// 使用记录相关
	CreateUsage(usage *models.SettlementWalletUsage) error
	GetUsage(id int64) (*models.SettlementWalletUsage, error)
	GetUsageByNo(usageNo string) (*models.SettlementWalletUsage, error)
	UpdateUsage(usage *models.SettlementWalletUsage) error
	GetUsagesByAgent(agentID int64, usageType *int16, limit, offset int) ([]*models.SettlementWalletUsage, int64, error)
	GetPendingReturn(agentID int64) ([]*models.SettlementWalletUsage, error)
}

// GormSettlementWalletRepository GORM实现
type GormSettlementWalletRepository struct {
	db *gorm.DB
}

// NewGormSettlementWalletRepository 创建仓库
func NewGormSettlementWalletRepository(db *gorm.DB) *GormSettlementWalletRepository {
	return &GormSettlementWalletRepository{db: db}
}

// CreateUsage 创建使用记录
func (r *GormSettlementWalletRepository) CreateUsage(usage *models.SettlementWalletUsage) error {
	return r.db.Create(usage).Error
}

// GetUsage 获取使用记录
func (r *GormSettlementWalletRepository) GetUsage(id int64) (*models.SettlementWalletUsage, error) {
	var usage models.SettlementWalletUsage
	err := r.db.First(&usage, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &usage, err
}

// GetUsageByNo 根据单号获取使用记录
func (r *GormSettlementWalletRepository) GetUsageByNo(usageNo string) (*models.SettlementWalletUsage, error) {
	var usage models.SettlementWalletUsage
	err := r.db.Where("usage_no = ?", usageNo).First(&usage).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &usage, err
}

// UpdateUsage 更新使用记录
func (r *GormSettlementWalletRepository) UpdateUsage(usage *models.SettlementWalletUsage) error {
	return r.db.Save(usage).Error
}

// GetUsagesByAgent 获取代理商使用记录
func (r *GormSettlementWalletRepository) GetUsagesByAgent(agentID int64, usageType *int16, limit, offset int) ([]*models.SettlementWalletUsage, int64, error) {
	query := r.db.Model(&models.SettlementWalletUsage{}).Where("agent_id = ?", agentID)
	if usageType != nil {
		query = query.Where("usage_type = ?", *usageType)
	}

	var total int64
	query.Count(&total)

	var usages []*models.SettlementWalletUsage
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&usages).Error
	return usages, total, err
}

// GetPendingReturn 获取待归还的使用记录
func (r *GormSettlementWalletRepository) GetPendingReturn(agentID int64) ([]*models.SettlementWalletUsage, error) {
	var usages []*models.SettlementWalletUsage
	err := r.db.Where("agent_id = ? AND status = ? AND usage_type = ?",
		agentID, models.SettlementUsageStatusToReturn, models.SettlementUsageTypeUse).
		Order("created_at ASC").
		Find(&usages).Error
	return usages, err
}

// GetSubordinateUnwithdrawBalance 获取下级未提现余额汇总
func (r *GormSettlementWalletRepository) GetSubordinateUnwithdrawBalance(parentAgentID int64) (int64, error) {
	// 查询所有直属下级的钱包余额汇总
	var total int64
	err := r.db.Raw(`
		SELECT COALESCE(SUM(w.balance - w.frozen_amount), 0)
		FROM wallets w
		JOIN agents a ON w.agent_id = a.id
		WHERE a.parent_id = ? AND w.wallet_type IN (1, 2, 3)
	`, parentAgentID).Scan(&total).Error
	return total, err
}

// GetSubordinateBalanceDetails 获取下级余额明细
func (r *GormSettlementWalletRepository) GetSubordinateBalanceDetails(parentAgentID int64) ([]SubordinateBalance, error) {
	var results []SubordinateBalance
	err := r.db.Raw(`
		SELECT
			a.id as agent_id,
			a.agent_name,
			COALESCE(SUM(w.balance - w.frozen_amount), 0) as available_balance
		FROM agents a
		LEFT JOIN wallets w ON a.id = w.agent_id AND w.wallet_type IN (1, 2, 3)
		WHERE a.parent_id = ?
		GROUP BY a.id, a.agent_name
		HAVING COALESCE(SUM(w.balance - w.frozen_amount), 0) > 0
		ORDER BY available_balance DESC
	`, parentAgentID).Scan(&results).Error
	return results, err
}

// SubordinateBalance 下级余额
type SubordinateBalance struct {
	AgentID          int64  `json:"agent_id"`
	AgentName        string `json:"agent_name"`
	AvailableBalance int64  `json:"available_balance"`
}

var _ SettlementWalletRepository = (*GormSettlementWalletRepository)(nil)
