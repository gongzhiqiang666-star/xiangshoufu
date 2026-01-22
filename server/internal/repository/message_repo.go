package repository

import (
	"time"

	"gorm.io/gorm"

	"xiangshoufu/internal/models"
)

// GormMessageRepository GORM实现的消息仓库
type GormMessageRepository struct {
	db *gorm.DB
}

// NewGormMessageRepository 创建仓库
func NewGormMessageRepository(db *gorm.DB) *GormMessageRepository {
	return &GormMessageRepository{db: db}
}

// Create 创建消息
func (r *GormMessageRepository) Create(msg *models.Message) error {
	return r.db.Create(msg).Error
}

// BatchCreate 批量创建消息
func (r *GormMessageRepository) BatchCreate(msgs []*models.Message) error {
	if len(msgs) == 0 {
		return nil
	}
	return r.db.CreateInBatches(msgs, 100).Error
}

// FindByAgentID 根据代理商ID查找消息
func (r *GormMessageRepository) FindByAgentID(agentID int64, limit, offset int) ([]*models.Message, error) {
	var msgs []*models.Message
	err := r.db.Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	return msgs, err
}

// FindUnreadByAgentID 查找未读消息
func (r *GormMessageRepository) FindUnreadByAgentID(agentID int64) ([]*models.Message, error) {
	var msgs []*models.Message
	err := r.db.Where("agent_id = ? AND is_read = ?", agentID, false).
		Order("created_at DESC").
		Find(&msgs).Error
	return msgs, err
}

// MarkAsRead 标记消息已读
func (r *GormMessageRepository) MarkAsRead(id int64) error {
	return r.db.Model(&models.Message{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记所有消息已读
func (r *GormMessageRepository) MarkAllAsRead(agentID int64) error {
	return r.db.Model(&models.Message{}).
		Where("agent_id = ? AND is_read = ?", agentID, false).
		Update("is_read", true).Error
}

// DeleteExpired 删除过期消息
func (r *GormMessageRepository) DeleteExpired() (int64, error) {
	result := r.db.Where("expire_at IS NOT NULL AND expire_at < ?", time.Now()).
		Delete(&models.Message{})
	return result.RowsAffected, result.Error
}

// FindByAgentIDAndTypes 根据代理商ID和消息类型列表查找消息
func (r *GormMessageRepository) FindByAgentIDAndTypes(agentID int64, types []int16, limit, offset int) ([]*models.Message, error) {
	var msgs []*models.Message
	query := r.db.Where("agent_id = ?", agentID)
	if len(types) > 0 {
		query = query.Where("message_type IN ?", types)
	}
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	return msgs, err
}

// CountByAgentIDAndTypes 统计代理商指定类型的消息数量
func (r *GormMessageRepository) CountByAgentIDAndTypes(agentID int64, types []int16) (int64, error) {
	var count int64
	query := r.db.Model(&models.Message{}).Where("agent_id = ?", agentID)
	if len(types) > 0 {
		query = query.Where("message_type IN ?", types)
	}
	err := query.Count(&count).Error
	return count, err
}

// GetStatsByAgentID 获取代理商消息分类统计
func (r *GormMessageRepository) GetStatsByAgentID(agentID int64) (*MessageStats, error) {
	stats := &MessageStats{}

	// 总消息数
	if err := r.db.Model(&models.Message{}).Where("agent_id = ?", agentID).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 未读总数
	if err := r.db.Model(&models.Message{}).Where("agent_id = ? AND is_read = ?", agentID, false).Count(&stats.UnreadTotal).Error; err != nil {
		return nil, err
	}

	// 分润类消息数（类型1,2,3,4）
	profitTypes := []int16{models.MessageTypeProfit, models.MessageTypeActivation, models.MessageTypeDeposit, models.MessageTypeSimCashback}
	if err := r.db.Model(&models.Message{}).Where("agent_id = ? AND message_type IN ?", agentID, profitTypes).Count(&stats.ProfitCount).Error; err != nil {
		return nil, err
	}

	// 注册类消息数（类型7）
	registerTypes := []int16{models.MessageTypeNewAgent}
	if err := r.db.Model(&models.Message{}).Where("agent_id = ? AND message_type IN ?", agentID, registerTypes).Count(&stats.RegisterCount).Error; err != nil {
		return nil, err
	}

	// 消费类消息数（类型8）
	consumptionTypes := []int16{models.MessageTypeTransaction}
	if err := r.db.Model(&models.Message{}).Where("agent_id = ? AND message_type IN ?", agentID, consumptionTypes).Count(&stats.ConsumptionCount).Error; err != nil {
		return nil, err
	}

	// 系统类消息数（类型5,6）
	systemTypes := []int16{models.MessageTypeRefund, models.MessageTypeAnnouncement}
	if err := r.db.Model(&models.Message{}).Where("agent_id = ? AND message_type IN ?", agentID, systemTypes).Count(&stats.SystemCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// FindAll 查找所有消息（管理端）
func (r *GormMessageRepository) FindAll(limit, offset int) ([]*models.Message, error) {
	var msgs []*models.Message
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	return msgs, err
}

// CountAll 统计所有消息数量
func (r *GormMessageRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).Count(&count).Error
	return count, err
}

// FindByID 根据ID查找消息
func (r *GormMessageRepository) FindByID(id int64) (*models.Message, error) {
	var msg models.Message
	err := r.db.Where("id = ?", id).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Delete 删除消息
func (r *GormMessageRepository) Delete(id int64) error {
	return r.db.Delete(&models.Message{}, id).Error
}

// FindByAgentIDs 根据代理商ID列表查找消息
func (r *GormMessageRepository) FindByAgentIDs(agentIDs []int64, limit, offset int) ([]*models.Message, error) {
	var msgs []*models.Message
	query := r.db.Order("created_at DESC")
	if len(agentIDs) > 0 {
		query = query.Where("agent_id IN ?", agentIDs)
	}
	err := query.Limit(limit).Offset(offset).Find(&msgs).Error
	return msgs, err
}

// 确保实现了接口
var _ MessageRepository = (*GormMessageRepository)(nil)
