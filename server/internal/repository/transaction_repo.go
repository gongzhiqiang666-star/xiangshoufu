package repository

import (
	"time"

	"gorm.io/gorm"
)

// GormTransactionRepository GORM实现的交易仓库
type GormTransactionRepository struct {
	db *gorm.DB
}

// NewGormTransactionRepository 创建仓库
func NewGormTransactionRepository(db *gorm.DB) *GormTransactionRepository {
	return &GormTransactionRepository{db: db}
}

// Create 创建交易
func (r *GormTransactionRepository) Create(tx *Transaction) error {
	return r.db.Create(tx).Error
}

// FindByOrderNo 根据订单号查找
func (r *GormTransactionRepository) FindByOrderNo(orderNo string) (*Transaction, error) {
	var tx Transaction
	err := r.db.Where("order_no = ?", orderNo).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// FindUnprocessedProfit 查找未计算分润的交易
func (r *GormTransactionRepository) FindUnprocessedProfit(limit int) ([]*Transaction, error) {
	var txs []*Transaction
	// 查找5分钟前创建且未计算分润的交易
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	err := r.db.Where("profit_status = ? AND received_at < ?", 0, fiveMinutesAgo).
		Order("received_at ASC").
		Limit(limit).
		Find(&txs).Error
	return txs, err
}

// UpdateProfitStatus 更新分润状态
func (r *GormTransactionRepository) UpdateProfitStatus(id int64, status int16) error {
	return r.db.Model(&Transaction{}).
		Where("id = ?", id).
		Update("profit_status", status).Error
}

// BatchUpdateProfitStatus 批量更新分润状态
func (r *GormTransactionRepository) BatchUpdateProfitStatus(ids []int64, status int16) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&Transaction{}).
		Where("id IN ?", ids).
		Update("profit_status", status).Error
}

// UpdateRefundStatus 更新退款状态
func (r *GormTransactionRepository) UpdateRefundStatus(id int64, status int16) error {
	return r.db.Model(&Transaction{}).
		Where("id = ?", id).
		Update("refund_status", status).Error
}

// 确保实现了接口
var _ TransactionRepository = (*GormTransactionRepository)(nil)

// GormProfitRecordRepository GORM实现的分润记录仓库
type GormProfitRecordRepository struct {
	db *gorm.DB
}

// NewGormProfitRecordRepository 创建仓库
func NewGormProfitRecordRepository(db *gorm.DB) *GormProfitRecordRepository {
	return &GormProfitRecordRepository{db: db}
}

// Create 创建分润记录
func (r *GormProfitRecordRepository) Create(record *ProfitRecord) error {
	return r.db.Create(record).Error
}

// BatchCreate 批量创建分润记录
func (r *GormProfitRecordRepository) BatchCreate(records []*ProfitRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.CreateInBatches(records, 100).Error
}

// FindByTransactionID 根据交易ID查找分润记录
func (r *GormProfitRecordRepository) FindByTransactionID(txID int64) ([]*ProfitRecord, error) {
	var records []*ProfitRecord
	err := r.db.Where("transaction_id = ?", txID).Find(&records).Error
	return records, err
}

// RevokeByTransactionID 撤销交易相关的分润记录
func (r *GormProfitRecordRepository) RevokeByTransactionID(txID int64, reason string) error {
	now := time.Now()
	return r.db.Model(&ProfitRecord{}).
		Where("transaction_id = ? AND is_revoked = ?", txID, false).
		Updates(map[string]interface{}{
			"is_revoked":    true,
			"revoked_at":    &now,
			"revoke_reason": reason,
		}).Error
}

// 确保实现了接口
var _ ProfitRecordRepository = (*GormProfitRecordRepository)(nil)

// TransactionStats 交易统计
type TransactionStats struct {
	TotalAmount int64 `json:"total_amount"`
	TotalCount  int64 `json:"total_count"`
	TotalFee    int64 `json:"total_fee"`
}

// GetAgentDailyStats 获取代理商今日交易统计
func (r *GormTransactionRepository) GetAgentDailyStats(agentID int64, date time.Time) (*TransactionStats, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var stats TransactionStats
	err := r.db.Model(&Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count, COALESCE(SUM(fee), 0) as total_fee").
		Where("agent_id = ? AND trade_time >= ? AND trade_time < ? AND trade_type = 1", agentID, startOfDay, endOfDay).
		Scan(&stats).Error

	return &stats, err
}

// GetAgentMonthlyStats 获取代理商本月交易统计
func (r *GormTransactionRepository) GetAgentMonthlyStats(agentID int64, date time.Time) (*TransactionStats, error) {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var stats TransactionStats
	err := r.db.Model(&Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count, COALESCE(SUM(fee), 0) as total_fee").
		Where("agent_id = ? AND trade_time >= ? AND trade_time < ? AND trade_type = 1", agentID, startOfMonth, endOfMonth).
		Scan(&stats).Error

	return &stats, err
}

// FindByAgentID 分页查询代理商交易
func (r *GormTransactionRepository) FindByAgentID(agentID int64, startTime, endTime *time.Time, tradeType *int16, limit, offset int) ([]*Transaction, int64, error) {
	query := r.db.Model(&Transaction{}).Where("agent_id = ?", agentID)

	if startTime != nil {
		query = query.Where("trade_time >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("trade_time < ?", *endTime)
	}
	if tradeType != nil {
		query = query.Where("trade_type = ?", *tradeType)
	}

	var total int64
	query.Count(&total)

	var transactions []*Transaction
	err := query.Order("trade_time DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

// GetTransactionTrend 获取交易趋势（按天）
func (r *GormTransactionRepository) GetTransactionTrend(agentID int64, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&Transaction{}).
		Select("DATE(trade_time) as date, SUM(amount) as amount, COUNT(*) as count").
		Where("agent_id = ? AND trade_time >= ? AND trade_time < ? AND trade_type = 1", agentID, startDate, endDate).
		Group("DATE(trade_time)").
		Order("date ASC").
		Find(&results).Error
	return results, err
}

// ProfitStats 分润统计
type ProfitStats struct {
	TotalAmount int64 `json:"total_amount"`
	TotalCount  int64 `json:"total_count"`
}

// GetAgentDailyProfitStats 获取代理商今日分润统计
func (r *GormProfitRecordRepository) GetAgentDailyProfitStats(agentID int64, date time.Time) (*ProfitStats, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var stats ProfitStats
	err := r.db.Model(&ProfitRecord{}).
		Select("COALESCE(SUM(profit_amount), 0) as total_amount, COUNT(*) as total_count").
		Where("agent_id = ? AND created_at >= ? AND created_at < ? AND is_revoked = false", agentID, startOfDay, endOfDay).
		Scan(&stats).Error

	return &stats, err
}

// GetAgentMonthlyProfitStats 获取代理商本月分润统计
func (r *GormProfitRecordRepository) GetAgentMonthlyProfitStats(agentID int64, date time.Time) (*ProfitStats, error) {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var stats ProfitStats
	err := r.db.Model(&ProfitRecord{}).
		Select("COALESCE(SUM(profit_amount), 0) as total_amount, COUNT(*) as total_count").
		Where("agent_id = ? AND created_at >= ? AND created_at < ? AND is_revoked = false", agentID, startOfMonth, endOfMonth).
		Scan(&stats).Error

	return &stats, err
}

// FindByMerchantID 分页查询商户交易
func (r *GormTransactionRepository) FindByMerchantID(merchantID int64, startTime, endTime *time.Time, limit, offset int) ([]*Transaction, int64, error) {
	query := r.db.Model(&Transaction{}).Where("merchant_id = ?", merchantID)

	if startTime != nil {
		query = query.Where("trade_time >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("trade_time < ?", *endTime)
	}

	var total int64
	query.Count(&total)

	var transactions []*Transaction
	err := query.Order("trade_time DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

// GetTerminalTotalTradeAmount 获取终端累计交易量（分）
// 用于激活奖励检查
func (r *GormTransactionRepository) GetTerminalTotalTradeAmount(terminalSN string) (int64, error) {
	var totalAmount int64
	err := r.db.Model(&Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("terminal_sn = ? AND trade_type = 1", terminalSN). // trade_type=1 表示消费
		Scan(&totalAmount).Error
	return totalAmount, err
}

// GetTerminalTradeAmountBetween 获取终端指定时间范围内的交易量（分）
func (r *GormTransactionRepository) GetTerminalTradeAmountBetween(terminalSN string, startTime, endTime time.Time) (int64, error) {
	var totalAmount int64
	err := r.db.Model(&Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("terminal_sn = ? AND trade_type = 1 AND trade_time >= ? AND trade_time < ?", terminalSN, startTime, endTime).
		Scan(&totalAmount).Error
	return totalAmount, err
}

// FindByAgentID 分页查询代理商分润记录
func (r *GormProfitRecordRepository) FindByAgentID(agentID int64, profitType *int16, startTime, endTime *time.Time, limit, offset int) ([]*ProfitRecord, int64, error) {
	query := r.db.Model(&ProfitRecord{}).Where("agent_id = ? AND is_revoked = false", agentID)

	if profitType != nil {
		query = query.Where("profit_type = ?", *profitType)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at < ?", *endTime)
	}

	var total int64
	query.Count(&total)

	var records []*ProfitRecord
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error
	return records, total, err
}
