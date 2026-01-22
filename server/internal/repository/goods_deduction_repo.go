package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// GoodsDeductionRepository 货款代扣仓库接口
type GoodsDeductionRepository interface {
	Create(deduction *models.GoodsDeduction) error
	Update(deduction *models.GoodsDeduction) error
	FindByID(id int64) (*models.GoodsDeduction, error)
	FindByDeductionNo(deductionNo string) (*models.GoodsDeduction, error)
	FindByFromAgent(agentID int64, status []int16, limit, offset int) ([]*models.GoodsDeduction, int64, error)
	FindByToAgent(agentID int64, status []int16, limit, offset int) ([]*models.GoodsDeduction, int64, error)
	FindActiveByToAgent(agentID int64) ([]*models.GoodsDeduction, error)
	FindByDistributeID(distributeID int64) (*models.GoodsDeduction, error)
	UpdateStatus(id int64, status int16) error
	UpdateDeductedAmount(id int64, deductAmount int64) error
	GetSummaryByFromAgent(agentID int64) (*models.GoodsDeductionSummary, error)
	GetSummaryByToAgent(agentID int64) (*models.GoodsDeductionSummary, error)
}

// GoodsDeductionDetailRepository 货款代扣明细仓库接口
type GoodsDeductionDetailRepository interface {
	Create(detail *models.GoodsDeductionDetail) error
	BatchCreate(details []*models.GoodsDeductionDetail) error
	FindByID(id int64) (*models.GoodsDeductionDetail, error)
	FindByDeductionID(deductionID int64, limit, offset int) ([]*models.GoodsDeductionDetail, int64, error)
	GetTotalDeductedByDeductionID(deductionID int64) (int64, error)
}

// GoodsDeductionTerminalRepository 货款代扣终端仓库接口
type GoodsDeductionTerminalRepository interface {
	Create(terminal *models.GoodsDeductionTerminal) error
	BatchCreate(terminals []*models.GoodsDeductionTerminal) error
	FindByDeductionID(deductionID int64) ([]*models.GoodsDeductionTerminal, error)
	FindByTerminalID(terminalID int64) ([]*models.GoodsDeductionTerminal, error)
}

// GoodsDeductionNotificationRepository 货款代扣通知仓库接口
type GoodsDeductionNotificationRepository interface {
	Create(notification *models.GoodsDeductionNotification) error
	FindByAgentID(agentID int64, isRead *bool, limit, offset int) ([]*models.GoodsDeductionNotification, int64, error)
	FindUnreadCount(agentID int64) (int64, error)
	MarkAsRead(id int64) error
	MarkAllAsRead(agentID int64) error
}

// ===============================
// GormGoodsDeductionRepository GORM实现
// ===============================

type GormGoodsDeductionRepository struct {
	db *gorm.DB
}

func NewGormGoodsDeductionRepository(db *gorm.DB) *GormGoodsDeductionRepository {
	return &GormGoodsDeductionRepository{db: db}
}

func (r *GormGoodsDeductionRepository) Create(deduction *models.GoodsDeduction) error {
	return r.db.Create(deduction).Error
}

func (r *GormGoodsDeductionRepository) Update(deduction *models.GoodsDeduction) error {
	return r.db.Save(deduction).Error
}

func (r *GormGoodsDeductionRepository) FindByID(id int64) (*models.GoodsDeduction, error) {
	var deduction models.GoodsDeduction
	err := r.db.First(&deduction, id).Error
	if err != nil {
		return nil, err
	}
	return &deduction, nil
}

func (r *GormGoodsDeductionRepository) FindByDeductionNo(deductionNo string) (*models.GoodsDeduction, error) {
	var deduction models.GoodsDeduction
	err := r.db.Where("deduction_no = ?", deductionNo).First(&deduction).Error
	if err != nil {
		return nil, err
	}
	return &deduction, nil
}

func (r *GormGoodsDeductionRepository) FindByFromAgent(agentID int64, status []int16, limit, offset int) ([]*models.GoodsDeduction, int64, error) {
	var deductions []*models.GoodsDeduction
	var total int64

	query := r.db.Model(&models.GoodsDeduction{}).Where("from_agent_id = ?", agentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&deductions).Error
	return deductions, total, err
}

func (r *GormGoodsDeductionRepository) FindByToAgent(agentID int64, status []int16, limit, offset int) ([]*models.GoodsDeduction, int64, error) {
	var deductions []*models.GoodsDeduction
	var total int64

	query := r.db.Model(&models.GoodsDeduction{}).Where("to_agent_id = ?", agentID)
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&deductions).Error
	return deductions, total, err
}

// FindActiveByToAgent 查找代理商所有进行中的货款代扣（用于实时扣款触发）
func (r *GormGoodsDeductionRepository) FindActiveByToAgent(agentID int64) ([]*models.GoodsDeduction, error) {
	var deductions []*models.GoodsDeduction
	err := r.db.Where("to_agent_id = ?", agentID).
		Where("status = ?", models.GoodsDeductionStatusInProgress).
		Where("remaining_amount > 0").
		Order("created_at ASC"). // 按创建时间排序，先创建的先扣
		Find(&deductions).Error
	return deductions, err
}

func (r *GormGoodsDeductionRepository) FindByDistributeID(distributeID int64) (*models.GoodsDeduction, error) {
	var deduction models.GoodsDeduction
	err := r.db.Where("distribute_id = ?", distributeID).First(&deduction).Error
	if err != nil {
		return nil, err
	}
	return &deduction, nil
}

func (r *GormGoodsDeductionRepository) UpdateStatus(id int64, status int16) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// 根据状态更新相应时间字段
	now := time.Now()
	switch status {
	case models.GoodsDeductionStatusInProgress:
		updates["accepted_at"] = &now
	case models.GoodsDeductionStatusCompleted:
		updates["completed_at"] = &now
	}

	return r.db.Model(&models.GoodsDeduction{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateDeductedAmount 更新已扣金额（原子操作）
func (r *GormGoodsDeductionRepository) UpdateDeductedAmount(id int64, deductAmount int64) error {
	return r.db.Model(&models.GoodsDeduction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deducted_amount":  gorm.Expr("deducted_amount + ?", deductAmount),
			"remaining_amount": gorm.Expr("remaining_amount - ?", deductAmount),
			"updated_at":       time.Now(),
		}).Error
}

// GetSummaryByFromAgent 获取发起方的货款代扣统计
func (r *GormGoodsDeductionRepository) GetSummaryByFromAgent(agentID int64) (*models.GoodsDeductionSummary, error) {
	var summary models.GoodsDeductionSummary

	// 总数和各状态数量
	r.db.Model(&models.GoodsDeduction{}).Where("from_agent_id = ?", agentID).Count(&summary.TotalCount)
	r.db.Model(&models.GoodsDeduction{}).Where("from_agent_id = ? AND status = ?", agentID, models.GoodsDeductionStatusPendingAccept).Count(&summary.PendingCount)
	r.db.Model(&models.GoodsDeduction{}).Where("from_agent_id = ? AND status = ?", agentID, models.GoodsDeductionStatusInProgress).Count(&summary.InProgressCount)
	r.db.Model(&models.GoodsDeduction{}).Where("from_agent_id = ? AND status = ?", agentID, models.GoodsDeductionStatusCompleted).Count(&summary.CompletedCount)

	// 金额统计
	var amounts struct {
		TotalAmount     int64
		DeductedAmount  int64
		RemainingAmount int64
	}
	r.db.Model(&models.GoodsDeduction{}).
		Where("from_agent_id = ?", agentID).
		Where("status IN ?", []int16{models.GoodsDeductionStatusInProgress, models.GoodsDeductionStatusCompleted}).
		Select("COALESCE(SUM(total_amount), 0) as total_amount, COALESCE(SUM(deducted_amount), 0) as deducted_amount, COALESCE(SUM(remaining_amount), 0) as remaining_amount").
		Scan(&amounts)

	summary.TotalAmount = amounts.TotalAmount
	summary.DeductedAmount = amounts.DeductedAmount
	summary.RemainingAmount = amounts.RemainingAmount

	return &summary, nil
}

// GetSummaryByToAgent 获取接收方的货款代扣统计
func (r *GormGoodsDeductionRepository) GetSummaryByToAgent(agentID int64) (*models.GoodsDeductionSummary, error) {
	var summary models.GoodsDeductionSummary

	// 总数和各状态数量
	r.db.Model(&models.GoodsDeduction{}).Where("to_agent_id = ?", agentID).Count(&summary.TotalCount)
	r.db.Model(&models.GoodsDeduction{}).Where("to_agent_id = ? AND status = ?", agentID, models.GoodsDeductionStatusPendingAccept).Count(&summary.PendingCount)
	r.db.Model(&models.GoodsDeduction{}).Where("to_agent_id = ? AND status = ?", agentID, models.GoodsDeductionStatusInProgress).Count(&summary.InProgressCount)
	r.db.Model(&models.GoodsDeduction{}).Where("to_agent_id = ? AND status = ?", agentID, models.GoodsDeductionStatusCompleted).Count(&summary.CompletedCount)

	// 金额统计
	var amounts struct {
		TotalAmount     int64
		DeductedAmount  int64
		RemainingAmount int64
	}
	r.db.Model(&models.GoodsDeduction{}).
		Where("to_agent_id = ?", agentID).
		Where("status IN ?", []int16{models.GoodsDeductionStatusInProgress, models.GoodsDeductionStatusCompleted}).
		Select("COALESCE(SUM(total_amount), 0) as total_amount, COALESCE(SUM(deducted_amount), 0) as deducted_amount, COALESCE(SUM(remaining_amount), 0) as remaining_amount").
		Scan(&amounts)

	summary.TotalAmount = amounts.TotalAmount
	summary.DeductedAmount = amounts.DeductedAmount
	summary.RemainingAmount = amounts.RemainingAmount

	return &summary, nil
}

var _ GoodsDeductionRepository = (*GormGoodsDeductionRepository)(nil)

// ===============================
// GormGoodsDeductionDetailRepository GORM实现
// ===============================

type GormGoodsDeductionDetailRepository struct {
	db *gorm.DB
}

func NewGormGoodsDeductionDetailRepository(db *gorm.DB) *GormGoodsDeductionDetailRepository {
	return &GormGoodsDeductionDetailRepository{db: db}
}

func (r *GormGoodsDeductionDetailRepository) Create(detail *models.GoodsDeductionDetail) error {
	return r.db.Create(detail).Error
}

func (r *GormGoodsDeductionDetailRepository) BatchCreate(details []*models.GoodsDeductionDetail) error {
	if len(details) == 0 {
		return nil
	}
	return r.db.CreateInBatches(details, 100).Error
}

func (r *GormGoodsDeductionDetailRepository) FindByID(id int64) (*models.GoodsDeductionDetail, error) {
	var detail models.GoodsDeductionDetail
	err := r.db.First(&detail, id).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *GormGoodsDeductionDetailRepository) FindByDeductionID(deductionID int64, limit, offset int) ([]*models.GoodsDeductionDetail, int64, error) {
	var details []*models.GoodsDeductionDetail
	var total int64

	query := r.db.Model(&models.GoodsDeductionDetail{}).Where("deduction_id = ?", deductionID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&details).Error
	return details, total, err
}

func (r *GormGoodsDeductionDetailRepository) GetTotalDeductedByDeductionID(deductionID int64) (int64, error) {
	var total int64
	err := r.db.Model(&models.GoodsDeductionDetail{}).
		Where("deduction_id = ?", deductionID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

var _ GoodsDeductionDetailRepository = (*GormGoodsDeductionDetailRepository)(nil)

// ===============================
// GormGoodsDeductionTerminalRepository GORM实现
// ===============================

type GormGoodsDeductionTerminalRepository struct {
	db *gorm.DB
}

func NewGormGoodsDeductionTerminalRepository(db *gorm.DB) *GormGoodsDeductionTerminalRepository {
	return &GormGoodsDeductionTerminalRepository{db: db}
}

func (r *GormGoodsDeductionTerminalRepository) Create(terminal *models.GoodsDeductionTerminal) error {
	return r.db.Create(terminal).Error
}

func (r *GormGoodsDeductionTerminalRepository) BatchCreate(terminals []*models.GoodsDeductionTerminal) error {
	if len(terminals) == 0 {
		return nil
	}
	return r.db.CreateInBatches(terminals, 100).Error
}

func (r *GormGoodsDeductionTerminalRepository) FindByDeductionID(deductionID int64) ([]*models.GoodsDeductionTerminal, error) {
	var terminals []*models.GoodsDeductionTerminal
	err := r.db.Where("deduction_id = ?", deductionID).Find(&terminals).Error
	return terminals, err
}

func (r *GormGoodsDeductionTerminalRepository) FindByTerminalID(terminalID int64) ([]*models.GoodsDeductionTerminal, error) {
	var terminals []*models.GoodsDeductionTerminal
	err := r.db.Where("terminal_id = ?", terminalID).Find(&terminals).Error
	return terminals, err
}

var _ GoodsDeductionTerminalRepository = (*GormGoodsDeductionTerminalRepository)(nil)

// ===============================
// GormGoodsDeductionNotificationRepository GORM实现
// ===============================

type GormGoodsDeductionNotificationRepository struct {
	db *gorm.DB
}

func NewGormGoodsDeductionNotificationRepository(db *gorm.DB) *GormGoodsDeductionNotificationRepository {
	return &GormGoodsDeductionNotificationRepository{db: db}
}

func (r *GormGoodsDeductionNotificationRepository) Create(notification *models.GoodsDeductionNotification) error {
	return r.db.Create(notification).Error
}

func (r *GormGoodsDeductionNotificationRepository) FindByAgentID(agentID int64, isRead *bool, limit, offset int) ([]*models.GoodsDeductionNotification, int64, error) {
	var notifications []*models.GoodsDeductionNotification
	var total int64

	query := r.db.Model(&models.GoodsDeductionNotification{}).Where("agent_id = ?", agentID)
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, total, err
}

func (r *GormGoodsDeductionNotificationRepository) FindUnreadCount(agentID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.GoodsDeductionNotification{}).
		Where("agent_id = ? AND is_read = ?", agentID, false).
		Count(&count).Error
	return count, err
}

func (r *GormGoodsDeductionNotificationRepository) MarkAsRead(id int64) error {
	now := time.Now()
	return r.db.Model(&models.GoodsDeductionNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *GormGoodsDeductionNotificationRepository) MarkAllAsRead(agentID int64) error {
	now := time.Now()
	return r.db.Model(&models.GoodsDeductionNotification{}).
		Where("agent_id = ? AND is_read = ?", agentID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

var _ GoodsDeductionNotificationRepository = (*GormGoodsDeductionNotificationRepository)(nil)
