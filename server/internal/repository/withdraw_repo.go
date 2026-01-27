package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
)

// WithdrawRepository 提现记录仓储接口
type WithdrawRepository interface {
	Create(record *models.WithdrawRecord) error
	FindByID(id int64) (*models.WithdrawRecord, error)
	FindByNo(withdrawNo string) (*models.WithdrawRecord, error)
	FindByAgentID(agentID int64, status *int16, limit, offset int) ([]*models.WithdrawRecord, int64, error)
	FindPending(limit, offset int) ([]*models.WithdrawRecord, int64, error)
	Update(record *models.WithdrawRecord) error
	UpdateStatus(id int64, status int16) error
	GetStatsByAgent(agentID int64) (*models.WithdrawStats, error)
}

// GormWithdrawRepository GORM实现的提现记录仓储
type GormWithdrawRepository struct {
	db *gorm.DB
}

// NewGormWithdrawRepository 创建提现记录仓储
func NewGormWithdrawRepository(db *gorm.DB) *GormWithdrawRepository {
	return &GormWithdrawRepository{db: db}
}

// Create 创建提现记录
func (r *GormWithdrawRepository) Create(record *models.WithdrawRecord) error {
	return r.db.Create(record).Error
}

// FindByID 根据ID查询
func (r *GormWithdrawRepository) FindByID(id int64) (*models.WithdrawRecord, error) {
	var record models.WithdrawRecord
	err := r.db.Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByNo 根据提现单号查询
func (r *GormWithdrawRepository) FindByNo(withdrawNo string) (*models.WithdrawRecord, error) {
	var record models.WithdrawRecord
	err := r.db.Where("withdraw_no = ?", withdrawNo).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByAgentID 根据代理商ID查询提现记录
func (r *GormWithdrawRepository) FindByAgentID(agentID int64, status *int16, limit, offset int) ([]*models.WithdrawRecord, int64, error) {
	var records []*models.WithdrawRecord
	var total int64

	query := r.db.Model(&models.WithdrawRecord{}).Where("agent_id = ?", agentID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// FindPending 查询待审核的提现记录
func (r *GormWithdrawRepository) FindPending(limit, offset int) ([]*models.WithdrawRecord, int64, error) {
	var records []*models.WithdrawRecord
	var total int64

	query := r.db.Model(&models.WithdrawRecord{}).Where("status = ?", models.WithdrawStatusPending)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at ASC").Limit(limit).Offset(offset).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// Update 更新提现记录
func (r *GormWithdrawRepository) Update(record *models.WithdrawRecord) error {
	record.UpdatedAt = time.Now()
	return r.db.Save(record).Error
}

// UpdateStatus 更新状态
func (r *GormWithdrawRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.WithdrawRecord{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// GetStatsByAgent 获取代理商提现统计
func (r *GormWithdrawRepository) GetStatsByAgent(agentID int64) (*models.WithdrawStats, error) {
	stats := &models.WithdrawStats{}

	// 总计
	var result struct {
		Count  int64
		Amount int64
	}
	err := r.db.Model(&models.WithdrawRecord{}).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("agent_id = ?", agentID).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询总计失败: %w", err)
	}
	stats.TotalCount = result.Count
	stats.TotalAmount = result.Amount

	// 待审核
	err = r.db.Model(&models.WithdrawRecord{}).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("agent_id = ? AND status = ?", agentID, models.WithdrawStatusPending).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询待审核失败: %w", err)
	}
	stats.PendingCount = result.Count
	stats.PendingAmount = result.Amount

	// 已打款
	err = r.db.Model(&models.WithdrawRecord{}).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("agent_id = ? AND status = ?", agentID, models.WithdrawStatusPaid).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询已打款失败: %w", err)
	}
	stats.PaidCount = result.Count
	stats.PaidAmount = result.Amount

	// 已拒绝
	err = r.db.Model(&models.WithdrawRecord{}).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("agent_id = ? AND status = ?", agentID, models.WithdrawStatusRejected).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询已拒绝失败: %w", err)
	}
	stats.RejectedCount = result.Count
	stats.RejectedAmount = result.Amount

	return stats, nil
}

// GenerateWithdrawNo 生成提现单号
func GenerateWithdrawNo() string {
	return fmt.Sprintf("WD%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000000)
}

// CountPendingByAgentID 统计代理商待处理的提现数量（待审核+待打款）
func (r *GormWithdrawRepository) CountPendingByAgentID(agentID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.WithdrawRecord{}).
		Where("agent_id = ? AND status IN (?, ?)", agentID, models.WithdrawStatusPending, models.WithdrawStatusApproved).
		Count(&count).Error
	return count, err
}
