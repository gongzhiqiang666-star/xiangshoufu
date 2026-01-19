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

// GormPolicyTemplateRepository 政策模板仓储
type GormPolicyTemplateRepository struct {
	db *gorm.DB
}

// NewGormPolicyTemplateRepository 创建政策模板仓储
func NewGormPolicyTemplateRepository(db *gorm.DB) *GormPolicyTemplateRepository {
	return &GormPolicyTemplateRepository{db: db}
}

// FindAll 查询所有模板
func (r *GormPolicyTemplateRepository) FindAll(channelID *int64, status *int16) ([]*models.PolicyTemplate, error) {
	var templates []*models.PolicyTemplate

	query := r.db.Model(&models.PolicyTemplate{})

	if channelID != nil {
		query = query.Where("channel_id = ?", *channelID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("channel_id, is_default DESC, created_at DESC").Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

// FindByID 根据ID查询模板
func (r *GormPolicyTemplateRepository) FindByID(id int64) (*models.PolicyTemplate, error) {
	var template models.PolicyTemplate
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// FindDefaultByChannel 查询通道默认模板
func (r *GormPolicyTemplateRepository) FindDefaultByChannel(channelID int64) (*models.PolicyTemplate, error) {
	var template models.PolicyTemplate
	if err := r.db.Where("channel_id = ? AND is_default = true", channelID).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}
