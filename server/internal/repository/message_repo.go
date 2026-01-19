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

// 确保实现了接口
var _ MessageRepository = (*GormMessageRepository)(nil)
