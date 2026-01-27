package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// GormMerchantRepository 商户仓储
type GormMerchantRepository struct {
	db *gorm.DB
}

// NewGormMerchantRepository 创建商户仓储
func NewGormMerchantRepository(db *gorm.DB) *GormMerchantRepository {
	return &GormMerchantRepository{db: db}
}

// MerchantStats 商户统计
type MerchantStats struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	PendingCount  int64 `json:"pending_count"`
	DisabledCount int64 `json:"disabled_count"`
}

// FindByAgentID 根据代理商ID查询商户列表
func (r *GormMerchantRepository) FindByAgentID(agentID int64, keyword string, status *int16, limit, offset int) ([]*models.Merchant, int64, error) {
	var merchants []*models.Merchant
	var total int64

	query := r.db.Model(&models.Merchant{}).Where("agent_id = ?", agentID)

	if keyword != "" {
		query = query.Where("merchant_no LIKE ? OR merchant_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

// FindByID 根据ID查询商户
func (r *GormMerchantRepository) FindByID(id int64) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.First(&merchant, id).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}

// FindByMerchantNo 根据商户号查询
func (r *GormMerchantRepository) FindByMerchantNo(merchantNo string) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.Where("merchant_no = ?", merchantNo).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}

// GetAgentMerchantStats 获取代理商商户统计
func (r *GormMerchantRepository) GetAgentMerchantStats(agentID int64) (*MerchantStats, error) {
	stats := &MerchantStats{}

	// 总数
	r.db.Model(&models.Merchant{}).Where("agent_id = ?", agentID).Count(&stats.TotalCount)

	// 正常
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND status = 1", agentID).Count(&stats.ActiveCount)

	// 待审核
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND approve_status = 1", agentID).Count(&stats.PendingCount)

	// 禁用
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND status = 2", agentID).Count(&stats.DisabledCount)

	return stats, nil
}

// GetMerchantTransStats 获取商户交易统计
type MerchantTransStats struct {
	TotalAmount int64 `json:"total_amount"`
	TotalCount  int64 `json:"total_count"`
	TotalFee    int64 `json:"total_fee"`
}

func (r *GormMerchantRepository) GetMerchantTransStats(merchantID int64, startTime, endTime *time.Time) (*MerchantTransStats, error) {
	stats := &MerchantTransStats{}

	query := r.db.Table("transactions").
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count, COALESCE(SUM(fee), 0) as total_fee").
		Where("merchant_id = ?", merchantID)

	if startTime != nil {
		query = query.Where("trade_time >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("trade_time < ?", *endTime)
	}

	if err := query.Scan(stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// ==================== 商户仓储新增方法 ====================

// MerchantQueryParams 商户查询参数
type MerchantQueryParams struct {
	AgentID      int64
	Keyword      string
	Status       *int16
	MerchantType string
	IsDirect     *bool
	ChannelID    *int64
	Limit        int
	Offset       int
}

// FindByParams 根据参数查询商户列表（增强版）
func (r *GormMerchantRepository) FindByParams(params MerchantQueryParams) ([]*models.Merchant, int64, error) {
	var merchants []*models.Merchant
	var total int64

	query := r.db.Model(&models.Merchant{}).Where("agent_id = ?", params.AgentID)

	if params.Keyword != "" {
		query = query.Where("merchant_no LIKE ? OR merchant_name LIKE ? OR terminal_sn LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.MerchantType != "" {
		query = query.Where("merchant_type = ?", params.MerchantType)
	}
	if params.IsDirect != nil {
		query = query.Where("is_direct = ?", *params.IsDirect)
	}
	if params.ChannelID != nil {
		query = query.Where("channel_id = ?", *params.ChannelID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

// Create 创建商户
func (r *GormMerchantRepository) Create(merchant *models.Merchant) error {
	return r.db.Create(merchant).Error
}

// Update 更新商户
func (r *GormMerchantRepository) Update(merchant *models.Merchant) error {
	return r.db.Save(merchant).Error
}

// Delete 删除商户（软删除场景下可改为状态变更）
func (r *GormMerchantRepository) Delete(id int64) error {
	return r.db.Delete(&models.Merchant{}, id).Error
}

// UpdateStatus 更新商户状态
func (r *GormMerchantRepository) UpdateStatus(id int64, status int16) error {
	return r.db.Model(&models.Merchant{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// UpdateApproveStatus 更新商户审核状态
func (r *GormMerchantRepository) UpdateApproveStatus(id int64, status int16) error {
	return r.db.Model(&models.Merchant{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"approve_status": status,
			"updated_at":     time.Now(),
		}).Error
}

// UpdateRate 更新商户费率
func (r *GormMerchantRepository) UpdateRate(id int64, creditRate, debitRate string) error {
	return r.db.Model(&models.Merchant{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"credit_rate": creditRate,
			"debit_rate":  debitRate,
			"updated_at":  time.Now(),
		}).Error
}

// Register 商户登记（更新登记手机号和备注）
func (r *GormMerchantRepository) Register(id int64, phone, remark string) error {
	return r.db.Model(&models.Merchant{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"registered_phone": phone,
			"register_remark":  remark,
			"updated_at":       time.Now(),
		}).Error
}

// UpdateMerchantType 更新商户类型
func (r *GormMerchantRepository) UpdateMerchantType(id int64, merchantType string) error {
	return r.db.Model(&models.Merchant{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"merchant_type": merchantType,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateActivatedAt 更新激活时间
func (r *GormMerchantRepository) UpdateActivatedAt(id int64, activatedAt time.Time) error {
	return r.db.Model(&models.Merchant{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"activated_at": activatedAt,
			"updated_at":   time.Now(),
		}).Error
}

// FindByTerminalSN 根据终端SN查询商户
func (r *GormMerchantRepository) FindByTerminalSN(terminalSN string) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.Where("terminal_sn = ?", terminalSN).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}

// GetExtendedStats 获取扩展统计（包含直营/团队、商户类型分布 - 5档分类）
type ExtendedMerchantStats struct {
	TotalCount    int64 `json:"total_count"`
	ActiveCount   int64 `json:"active_count"`
	PendingCount  int64 `json:"pending_count"`
	DisabledCount int64 `json:"disabled_count"`
	DirectCount   int64 `json:"direct_count"`
	TeamCount     int64 `json:"team_count"`
	TodayNewCount int64 `json:"today_new_count"`
	// 商户类型分布（5档）
	QualityCount int64 `json:"quality_count"` // 优质
	MediumCount  int64 `json:"medium_count"`  // 中等
	NormalCount  int64 `json:"normal_count"`  // 普通
	WarningCount int64 `json:"warning_count"` // 预警
	ChurnedCount int64 `json:"churned_count"` // 流失
}

func (r *GormMerchantRepository) GetExtendedStats(agentID int64) (*ExtendedMerchantStats, error) {
	stats := &ExtendedMerchantStats{}

	baseQuery := r.db.Model(&models.Merchant{}).Where("agent_id = ?", agentID)

	// 基础统计
	baseQuery.Count(&stats.TotalCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND status = 1", agentID).Count(&stats.ActiveCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND approve_status = 1", agentID).Count(&stats.PendingCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND status = 2", agentID).Count(&stats.DisabledCount)

	// 直营/团队统计
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND is_direct = true", agentID).Count(&stats.DirectCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND is_direct = false", agentID).Count(&stats.TeamCount)

	// 今日新增
	today := time.Now().Truncate(24 * time.Hour)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND created_at >= ?", agentID, today).Count(&stats.TodayNewCount)

	// 商户类型分布（5档）
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND merchant_type = ?", agentID, models.MerchantTypeQuality).Count(&stats.QualityCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND merchant_type = ?", agentID, models.MerchantTypeMedium).Count(&stats.MediumCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND merchant_type = ?", agentID, models.MerchantTypeNormal).Count(&stats.NormalCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND merchant_type = ?", agentID, models.MerchantTypeWarning).Count(&stats.WarningCount)
	r.db.Model(&models.Merchant{}).Where("agent_id = ? AND merchant_type = ?", agentID, models.MerchantTypeChurned).Count(&stats.ChurnedCount)

	return stats, nil
}

// ExportMerchants 导出商户数据（不分页）
func (r *GormMerchantRepository) ExportMerchants(params MerchantQueryParams) ([]*models.Merchant, error) {
	var merchants []*models.Merchant

	query := r.db.Model(&models.Merchant{}).Where("agent_id = ?", params.AgentID)

	if params.Keyword != "" {
		query = query.Where("merchant_no LIKE ? OR merchant_name LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.MerchantType != "" {
		query = query.Where("merchant_type = ?", params.MerchantType)
	}
	if params.IsDirect != nil {
		query = query.Where("is_direct = ?", *params.IsDirect)
	}
	if params.ChannelID != nil {
		query = query.Where("channel_id = ?", *params.ChannelID)
	}

	if err := query.Order("created_at DESC").Find(&merchants).Error; err != nil {
		return nil, err
	}

	return merchants, nil
}

// FindAllMerchantIDs 获取所有活跃商户的ID（用于批量计算商户类型）
func (r *GormMerchantRepository) FindAllMerchantIDs(limit, offset int) ([]int64, error) {
	var ids []int64
	if err := r.db.Model(&models.Merchant{}).
		Where("status = ?", models.MerchantStatusActive).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// CountAllActiveMerchants 统计所有活跃商户数量
func (r *GormMerchantRepository) CountAllActiveMerchants() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Merchant{}).
		Where("status = ?", models.MerchantStatusActive).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetDB 获取数据库连接（用于原生SQL操作）
func (r *GormMerchantRepository) GetDB() *gorm.DB {
	return r.db
}
