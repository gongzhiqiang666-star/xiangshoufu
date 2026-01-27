package repository

import (
	"fmt"
	"strings"
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// GormWalletRepository GORM实现的钱包仓库
type GormWalletRepository struct {
	db *gorm.DB
}

// NewGormWalletRepository 创建仓库
func NewGormWalletRepository(db *gorm.DB) *GormWalletRepository {
	return &GormWalletRepository{db: db}
}

// FindByAgentAndType 根据代理商和类型查找钱包
func (r *GormWalletRepository) FindByAgentAndType(agentID int64, channelID int64, walletType int16) (*Wallet, error) {
	var wallet Wallet
	err := r.db.Where("agent_id = ? AND channel_id = ? AND wallet_type = ?", agentID, channelID, walletType).
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// UpdateBalance 更新钱包余额（使用乐观锁）
func (r *GormWalletRepository) UpdateBalance(id int64, amount int64) error {
	return r.db.Model(&Wallet{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"balance":      gorm.Expr("balance + ?", amount),
			"total_income": gorm.Expr("total_income + ?", amount),
			"version":      gorm.Expr("version + 1"),
		}).Error
}

// BatchUpdateBalance 批量更新钱包余额
func (r *GormWalletRepository) BatchUpdateBalance(updates map[int64]int64) error {
	if len(updates) == 0 {
		return nil
	}

	// 使用事务批量更新
	return r.db.Transaction(func(tx *gorm.DB) error {
		for walletID, amount := range updates {
			err := tx.Model(&Wallet{}).
				Where("id = ?", walletID).
				Updates(map[string]interface{}{
					"balance":      gorm.Expr("balance + ?", amount),
					"total_income": gorm.Expr("total_income + ?", amount),
					"version":      gorm.Expr("version + 1"),
				}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// 确保实现了接口
var _ WalletRepository = (*GormWalletRepository)(nil)

// UpdateFrozenAmount 更新冻结金额（正数增加，负数减少）
func (r *GormWalletRepository) UpdateFrozenAmount(id int64, amount int64) error {
	return r.db.Model(&Wallet{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"frozen_amount": gorm.Expr("frozen_amount + ?", amount),
			"version":       gorm.Expr("version + 1"),
		}).Error
}

// GormWalletLogRepository GORM实现的钱包流水仓库
type GormWalletLogRepository struct {
	db *gorm.DB
}

// NewGormWalletLogRepository 创建仓库
func NewGormWalletLogRepository(db *gorm.DB) *GormWalletLogRepository {
	return &GormWalletLogRepository{db: db}
}

// Create 创建流水记录
func (r *GormWalletLogRepository) Create(log *WalletLog) error {
	return r.db.Create(log).Error
}

// BatchCreate 批量创建流水记录
func (r *GormWalletLogRepository) BatchCreate(logs []*WalletLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.CreateInBatches(logs, 100).Error
}

// 确保实现了接口
var _ WalletLogRepository = (*GormWalletLogRepository)(nil)

// FindByAgentID 查询代理商的所有钱包
func (r *GormWalletRepository) FindByAgentID(agentID int64) ([]*Wallet, error) {
	var wallets []*Wallet
	err := r.db.Where("agent_id = ?", agentID).Order("channel_id, wallet_type").Find(&wallets).Error
	return wallets, err
}

// FindByID 根据ID查询钱包
func (r *GormWalletRepository) FindByID(id int64) (*Wallet, error) {
	var wallet Wallet
	err := r.db.First(&wallet, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &wallet, err
}

// FreezeBalance 冻结余额
func (r *GormWalletRepository) FreezeBalance(walletID int64, amount int64) error {
	return r.db.Model(&Wallet{}).
		Where("id = ? AND balance - frozen_amount >= ?", walletID, amount).
		Updates(map[string]interface{}{
			"frozen_amount": gorm.Expr("frozen_amount + ?", amount),
			"version":       gorm.Expr("version + 1"),
		}).Error
}

// UnfreezeBalance 解冻余额
func (r *GormWalletRepository) UnfreezeBalance(walletID int64, amount int64) error {
	return r.db.Model(&Wallet{}).
		Where("id = ? AND frozen_amount >= ?", walletID, amount).
		Updates(map[string]interface{}{
			"frozen_amount": gorm.Expr("frozen_amount - ?", amount),
			"version":       gorm.Expr("version + 1"),
		}).Error
}

// DeductFrozenBalance 扣除冻结余额（提现成功）
func (r *GormWalletRepository) DeductFrozenBalance(walletID int64, amount int64) error {
	return r.db.Model(&Wallet{}).
		Where("id = ? AND frozen_amount >= ? AND balance >= ?", walletID, amount, amount).
		Updates(map[string]interface{}{
			"balance":        gorm.Expr("balance - ?", amount),
			"frozen_amount":  gorm.Expr("frozen_amount - ?", amount),
			"total_withdraw": gorm.Expr("total_withdraw + ?", amount),
			"version":        gorm.Expr("version + 1"),
		}).Error
}

// FindByAgentID 查询代理商钱包流水
func (r *GormWalletLogRepository) FindByAgentID(agentID int64, walletID int64, logType *int16, startTime, endTime *time.Time, limit, offset int) ([]*WalletLog, int64, error) {
	query := r.db.Model(&WalletLog{}).Where("agent_id = ?", agentID)

	if walletID > 0 {
		query = query.Where("wallet_id = ?", walletID)
	}
	if logType != nil {
		query = query.Where("log_type = ?", *logType)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at < ?", *endTime)
	}

	var total int64
	query.Count(&total)

	var logs []*WalletLog
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

// GormAgentRepository GORM实现的代理商仓库
type GormAgentRepository struct {
	db *gorm.DB
}

// NewGormAgentRepository 创建仓库
func NewGormAgentRepository(db *gorm.DB) *GormAgentRepository {
	return &GormAgentRepository{db: db}
}

// FindByID 根据ID查找代理商
func (r *GormAgentRepository) FindByID(id int64) (*Agent, error) {
	var agent Agent
	err := r.db.First(&agent, id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// FindByAgentNo 根据编号查找代理商
func (r *GormAgentRepository) FindByAgentNo(agentNo string) (*Agent, error) {
	var agent Agent
	err := r.db.Where("agent_no = ?", agentNo).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// FindAncestors 查找所有上级代理商（通过物化路径）
func (r *GormAgentRepository) FindAncestors(agentID int64) ([]*Agent, error) {
	// 先获取当前代理商
	agent, err := r.FindByID(agentID)
	if err != nil {
		return nil, err
	}

	// 解析物化路径获取所有上级ID
	// 路径格式: /1/5/12/
	if agent.Path == "" || agent.Path == "/" {
		return []*Agent{}, nil
	}

	// 去除首尾斜杠，分割
	pathStr := strings.Trim(agent.Path, "/")
	if pathStr == "" {
		return []*Agent{}, nil
	}

	// 查询所有上级（不包括自己）
	var ancestors []*Agent
	err = r.db.Raw(`
		SELECT * FROM agents
		WHERE id IN (
			SELECT unnest(string_to_array(?, '/'))::bigint
		)
		ORDER BY level ASC
	`, pathStr).Scan(&ancestors).Error

	return ancestors, err
}

// GetAllAgentIDs 获取所有活跃代理商的ID（用于消息广播）
func (r *GormAgentRepository) GetAllAgentIDs() ([]int64, error) {
	var ids []int64
	err := r.db.Model(&Agent{}).
		Where("status = ?", 1). // 只获取活跃状态的代理商
		Pluck("id", &ids).Error
	return ids, err
}

// GetAgentIDsByLevel 获取指定层级的所有活跃代理商ID（用于消息广播）
func (r *GormAgentRepository) GetAgentIDsByLevel(level int) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&Agent{}).
		Where("level = ? AND status = ?", level, 1). // 指定层级且活跃
		Pluck("id", &ids).Error
	return ids, err
}

// 确保实现了接口
var _ AgentRepository = (*GormAgentRepository)(nil)

// GormAgentPolicyRepository GORM实现的代理商政策仓库
type GormAgentPolicyRepository struct {
	db *gorm.DB
}

// NewGormAgentPolicyRepository 创建仓库
func NewGormAgentPolicyRepository(db *gorm.DB) *GormAgentPolicyRepository {
	return &GormAgentPolicyRepository{db: db}
}

// FindByAgentAndChannel 根据代理商和通道查找政策
func (r *GormAgentPolicyRepository) FindByAgentAndChannel(agentID int64, channelID int64) (*AgentPolicy, error) {
	var policy AgentPolicy
	err := r.db.Where("agent_id = ? AND channel_id = ?", agentID, channelID).First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// FindByAgentID 根据代理商查找所有政策
func (r *GormAgentPolicyRepository) FindByAgentID(agentID int64) ([]*AgentPolicy, error) {
	var policies []*AgentPolicy
	err := r.db.Where("agent_id = ?", agentID).Order("channel_id").Find(&policies).Error
	return policies, err
}

// Create 创建代理商政策
func (r *GormAgentPolicyRepository) Create(policy *AgentPolicy) error {
	return r.db.Create(policy).Error
}

// 确保实现了接口
var _ AgentPolicyRepository = (*GormAgentPolicyRepository)(nil)

// Agent 完整的代理商模型（扩展）
type AgentFull struct {
	Agent
	ContactName         string    `json:"contact_name" gorm:"size:50"`
	ContactPhone        string    `json:"contact_phone" gorm:"size:20"`
	IDCardNo            string    `json:"id_card_no" gorm:"size:18"`
	BankName            string    `json:"bank_name" gorm:"size:100"`
	BankAccount         string    `json:"bank_account" gorm:"size:30"`
	BankCardNo          string    `json:"bank_card_no" gorm:"size:25"`
	InviteCode          string    `json:"invite_code" gorm:"size:20"`
	QRCodeURL           string    `json:"qr_code_url" gorm:"size:255"`
	DirectAgentCount    int       `json:"direct_agent_count" gorm:"default:0"`
	DirectMerchantCount int       `json:"direct_merchant_count" gorm:"default:0"`
	TeamAgentCount      int       `json:"team_agent_count" gorm:"default:0"`
	TeamMerchantCount   int       `json:"team_merchant_count" gorm:"default:0"`
	RegisterTime        time.Time `json:"register_time" gorm:"default:now()"`
	CreatedAt           time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName 指定表名
func (AgentFull) TableName() string {
	return "agents"
}

// FindSubordinates 查找直属下级代理商
func (r *GormAgentRepository) FindSubordinates(parentID int64, keyword string, status *int16, limit, offset int) ([]*AgentFull, int64, error) {
	query := r.db.Model(&AgentFull{}).Where("parent_id = ?", parentID)

	if keyword != "" {
		query = query.Where("agent_no LIKE ? OR agent_name LIKE ? OR contact_phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	query.Count(&total)

	var agents []*AgentFull
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&agents).Error
	return agents, total, err
}

// FindByIDFull 根据ID查找完整代理商信息
func (r *GormAgentRepository) FindByIDFull(id int64) (*AgentFull, error) {
	var agent AgentFull
	err := r.db.First(&agent, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &agent, err
}

// Update 更新代理商
func (r *GormAgentRepository) Update(agent *AgentFull) error {
	return r.db.Save(agent).Error
}

// FindAllByParentPath 查找所有下级（通过路径匹配）
func (r *GormAgentRepository) FindAllByParentPath(parentID int64, limit, offset int) ([]*AgentFull, int64, error) {
	pathPattern := fmt.Sprintf("%%/%d/%%", parentID)

	query := r.db.Model(&AgentFull{}).Where("path LIKE ?", pathPattern)

	var total int64
	query.Count(&total)

	var agents []*AgentFull
	err := query.Order("level ASC, created_at DESC").Limit(limit).Offset(offset).Find(&agents).Error
	return agents, total, err
}

// Create 创建代理商
func (r *GormAgentRepository) Create(agent *AgentFull) error {
	return r.db.Create(agent).Error
}

// FindByPhone 根据手机号查找代理商
func (r *GormAgentRepository) FindByPhone(phone string) (*AgentFull, error) {
	var agent AgentFull
	err := r.db.Where("contact_phone = ?", phone).First(&agent).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &agent, err
}

// FindByInviteCode 根据邀请码查找代理商
func (r *GormAgentRepository) FindByInviteCode(inviteCode string) (*AgentFull, error) {
	var agent AgentFull
	err := r.db.Where("invite_code = ?", inviteCode).First(&agent).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &agent, err
}

// GetDailyAgentSequence 获取当日代理商序号
func (r *GormAgentRepository) GetDailyAgentSequence(dateStr string) int {
	var count int64
	pattern := "A" + dateStr + "%"
	r.db.Model(&AgentFull{}).Where("agent_no LIKE ?", pattern).Count(&count)
	return int(count) + 1
}

// IncrementDirectAgentCount 增加直属代理商计数
func (r *GormAgentRepository) IncrementDirectAgentCount(agentID int64) error {
	return r.db.Model(&AgentFull{}).
		Where("id = ?", agentID).
		Update("direct_agent_count", gorm.Expr("direct_agent_count + 1")).Error
}

// IncrementTeamAgentCount 增加团队代理商计数
func (r *GormAgentRepository) IncrementTeamAgentCount(agentID int64) error {
	return r.db.Model(&AgentFull{}).
		Where("id = ?", agentID).
		Update("team_agent_count", gorm.Expr("team_agent_count + 1")).Error
}

// SearchAgents 全局搜索代理商
func (r *GormAgentRepository) SearchAgents(keyword string, status *int16, limit, offset int) ([]*AgentFull, int64, error) {
	query := r.db.Model(&AgentFull{})

	if keyword != "" {
		query = query.Where("agent_no LIKE ? OR agent_name LIKE ? OR contact_phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	query.Count(&total)

	var agents []*AgentFull
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&agents).Error
	return agents, total, err
}

// ============================================================
// 钱包拆分配置仓库
// ============================================================

// GormWalletSplitConfigRepository GORM实现的钱包拆分配置仓库
type GormWalletSplitConfigRepository struct {
	db *gorm.DB
}

// NewGormWalletSplitConfigRepository 创建仓库
func NewGormWalletSplitConfigRepository(db *gorm.DB) *GormWalletSplitConfigRepository {
	return &GormWalletSplitConfigRepository{db: db}
}

// FindByAgentID 根据代理商ID查找拆分配置
func (r *GormWalletSplitConfigRepository) FindByAgentID(agentID int64) (*models.AgentWalletSplitConfig, error) {
	var config models.AgentWalletSplitConfig
	err := r.db.Where("agent_id = ?", agentID).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

// Create 创建拆分配置
func (r *GormWalletSplitConfigRepository) Create(config *models.AgentWalletSplitConfig) error {
	return r.db.Create(config).Error
}

// Update 更新拆分配置
func (r *GormWalletSplitConfigRepository) Update(config *models.AgentWalletSplitConfig) error {
	return r.db.Save(config).Error
}

// Upsert 创建或更新拆分配置
func (r *GormWalletSplitConfigRepository) Upsert(config *models.AgentWalletSplitConfig) error {
	return r.db.Save(config).Error
}

// ============================================================
// 政策提现门槛配置仓库
// ============================================================

// GormPolicyWithdrawThresholdRepository GORM实现的政策提现门槛仓库
type GormPolicyWithdrawThresholdRepository struct {
	db *gorm.DB
}

// NewGormPolicyWithdrawThresholdRepository 创建仓库
func NewGormPolicyWithdrawThresholdRepository(db *gorm.DB) *GormPolicyWithdrawThresholdRepository {
	return &GormPolicyWithdrawThresholdRepository{db: db}
}

// FindByTemplateID 根据政策模版ID查找所有门槛配置
func (r *GormPolicyWithdrawThresholdRepository) FindByTemplateID(templateID int64) ([]*models.PolicyWithdrawThreshold, error) {
	var thresholds []*models.PolicyWithdrawThreshold
	err := r.db.Where("template_id = ?", templateID).Order("wallet_type, channel_id").Find(&thresholds).Error
	return thresholds, err
}

// FindByTemplateAndWalletType 根据政策模版和钱包类型查找门槛配置
func (r *GormPolicyWithdrawThresholdRepository) FindByTemplateAndWalletType(templateID int64, walletType int16) ([]*models.PolicyWithdrawThreshold, error) {
	var thresholds []*models.PolicyWithdrawThreshold
	err := r.db.Where("template_id = ? AND wallet_type = ?", templateID, walletType).Order("channel_id").Find(&thresholds).Error
	return thresholds, err
}

// FindByTemplateWalletAndChannel 根据政策模版、钱包类型和通道查找门槛配置
func (r *GormPolicyWithdrawThresholdRepository) FindByTemplateWalletAndChannel(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
	var threshold models.PolicyWithdrawThreshold
	err := r.db.Where("template_id = ? AND wallet_type = ? AND channel_id = ?", templateID, walletType, channelID).First(&threshold).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &threshold, err
}

// Create 创建门槛配置
func (r *GormPolicyWithdrawThresholdRepository) Create(threshold *models.PolicyWithdrawThreshold) error {
	return r.db.Create(threshold).Error
}

// Update 更新门槛配置
func (r *GormPolicyWithdrawThresholdRepository) Update(threshold *models.PolicyWithdrawThreshold) error {
	return r.db.Save(threshold).Error
}

// Upsert 创建或更新门槛配置
func (r *GormPolicyWithdrawThresholdRepository) Upsert(threshold *models.PolicyWithdrawThreshold) error {
	return r.db.Save(threshold).Error
}

// BatchUpsert 批量创建或更新门槛配置
func (r *GormPolicyWithdrawThresholdRepository) BatchUpsert(thresholds []*models.PolicyWithdrawThreshold) error {
	if len(thresholds) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, threshold := range thresholds {
			if err := tx.Save(threshold).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteByTemplateID 删除政策模版的所有门槛配置
func (r *GormPolicyWithdrawThresholdRepository) DeleteByTemplateID(templateID int64) error {
	return r.db.Where("template_id = ?", templateID).Delete(&models.PolicyWithdrawThreshold{}).Error
}

// DeleteByTemplateAndChannel 删除政策模版指定通道的门槛配置
func (r *GormPolicyWithdrawThresholdRepository) DeleteByTemplateAndChannel(templateID int64, channelID int64) error {
	return r.db.Where("template_id = ? AND channel_id = ?", templateID, channelID).Delete(&models.PolicyWithdrawThreshold{}).Error
}
