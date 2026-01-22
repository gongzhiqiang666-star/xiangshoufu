package service

import (
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// AgentChannelService 代理商通道服务
type AgentChannelService struct {
	agentChannelRepo *repository.GormAgentChannelRepository
	channelRepo      *repository.GormChannelRepository
	agentRepo        repository.AgentRepository
}

// NewAgentChannelService 创建代理商通道服务
func NewAgentChannelService(
	agentChannelRepo *repository.GormAgentChannelRepository,
	channelRepo *repository.GormChannelRepository,
	agentRepo repository.AgentRepository,
) *AgentChannelService {
	return &AgentChannelService{
		agentChannelRepo: agentChannelRepo,
		channelRepo:      channelRepo,
		agentRepo:        agentRepo,
	}
}

// GetAgentChannels 获取代理商的通道配置列表
func (s *AgentChannelService) GetAgentChannels(agentID int64) ([]*models.AgentChannelWithInfo, error) {
	return s.agentChannelRepo.FindByAgentIDWithInfo(agentID)
}

// GetEnabledChannels 获取代理商已启用的通道列表（用于APP端显示）
func (s *AgentChannelService) GetEnabledChannels(agentID int64) ([]*models.AgentChannelWithInfo, error) {
	// 获取所有配置
	channels, err := s.agentChannelRepo.FindByAgentIDWithInfo(agentID)
	if err != nil {
		return nil, err
	}

	// 过滤出已启用且可见的通道
	var result []*models.AgentChannelWithInfo
	for _, ch := range channels {
		if ch.IsEnabled && ch.IsVisible {
			result = append(result, ch)
		}
	}

	return result, nil
}

// EnableChannel 启用代理商通道
func (s *AgentChannelService) EnableChannel(agentID, channelID int64, operatorID int64) error {
	// 检查通道是否存在
	channel, err := s.channelRepo.FindByID(channelID)
	if err != nil {
		return fmt.Errorf("通道不存在: %w", err)
	}
	if channel.Status != 1 {
		return fmt.Errorf("通道已禁用")
	}

	// 检查是否已有配置
	existing, _ := s.agentChannelRepo.FindByAgentAndChannel(agentID, channelID)
	if existing != nil {
		// 更新启用状态
		return s.agentChannelRepo.Enable(agentID, channelID, operatorID)
	}

	// 创建新配置
	now := time.Now()
	ac := &models.AgentChannel{
		AgentID:   agentID,
		ChannelID: channelID,
		IsEnabled: true,
		IsVisible: true,
		EnabledAt: &now,
		EnabledBy: &operatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.agentChannelRepo.Create(ac)
}

// DisableChannel 禁用代理商通道
func (s *AgentChannelService) DisableChannel(agentID, channelID int64, operatorID int64) error {
	// 检查是否已有配置
	existing, _ := s.agentChannelRepo.FindByAgentAndChannel(agentID, channelID)
	if existing == nil {
		return fmt.Errorf("代理商未配置该通道")
	}

	return s.agentChannelRepo.Disable(agentID, channelID, operatorID)
}

// SetChannelVisibility 设置通道可见性
func (s *AgentChannelService) SetChannelVisibility(agentID, channelID int64, isVisible bool) error {
	return s.agentChannelRepo.SetVisibility(agentID, channelID, isVisible)
}

// BatchEnableChannels 批量启用通道
func (s *AgentChannelService) BatchEnableChannels(agentID int64, channelIDs []int64, operatorID int64) error {
	if len(channelIDs) == 0 {
		return nil
	}

	// 确保所有通道配置都存在
	for _, channelID := range channelIDs {
		existing, _ := s.agentChannelRepo.FindByAgentAndChannel(agentID, channelID)
		if existing == nil {
			// 创建配置
			now := time.Now()
			ac := &models.AgentChannel{
				AgentID:   agentID,
				ChannelID: channelID,
				IsEnabled: true,
				IsVisible: true,
				EnabledAt: &now,
				EnabledBy: &operatorID,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := s.agentChannelRepo.Create(ac); err != nil {
				log.Printf("[AgentChannelService] Create channel config failed: %v", err)
			}
		}
	}

	return s.agentChannelRepo.BatchEnable(agentID, channelIDs, operatorID)
}

// BatchDisableChannels 批量禁用通道
func (s *AgentChannelService) BatchDisableChannels(agentID int64, channelIDs []int64, operatorID int64) error {
	return s.agentChannelRepo.BatchDisable(agentID, channelIDs, operatorID)
}

// InitAgentChannels 初始化代理商通道配置
// 在创建代理商时调用，为其配置所有活跃通道
func (s *AgentChannelService) InitAgentChannels(agentID int64, operatorID int64) error {
	// 获取所有活跃通道
	channels, err := s.channelRepo.FindAllActive()
	if err != nil {
		return fmt.Errorf("获取通道列表失败: %w", err)
	}

	if len(channels) == 0 {
		return nil
	}

	// 收集通道ID
	channelIDs := make([]int64, len(channels))
	for i, ch := range channels {
		channelIDs[i] = ch.ID
	}

	// 初始化配置
	return s.agentChannelRepo.InitAgentChannels(agentID, channelIDs, operatorID)
}

// SyncAgentChannelsFromParent 从上级代理商同步通道配置
// 下级代理商的通道配置不能超过上级的范围
func (s *AgentChannelService) SyncAgentChannelsFromParent(agentID int64, parentID int64, operatorID int64) error {
	// 获取上级的通道配置
	parentChannels, err := s.agentChannelRepo.FindEnabledByAgentID(parentID)
	if err != nil {
		return fmt.Errorf("获取上级通道配置失败: %w", err)
	}

	if len(parentChannels) == 0 {
		log.Printf("[AgentChannelService] Parent %d has no enabled channels", parentID)
		return nil
	}

	// 收集上级启用的通道ID
	channelIDs := make([]int64, len(parentChannels))
	for i, ch := range parentChannels {
		channelIDs[i] = ch.ChannelID
	}

	// 为下级创建相同的通道配置（默认启用）
	return s.agentChannelRepo.InitAgentChannels(agentID, channelIDs, operatorID)
}

// GetChannelStatus 获取代理商某个通道的启用状态
func (s *AgentChannelService) GetChannelStatus(agentID, channelID int64) (bool, error) {
	ac, err := s.agentChannelRepo.FindByAgentAndChannel(agentID, channelID)
	if err != nil {
		return false, err
	}
	return ac.IsEnabled, nil
}

// AgentChannelStats 代理商通道统计
type AgentChannelStats struct {
	TotalChannels   int `json:"total_channels"`   // 总通道数
	EnabledChannels int `json:"enabled_channels"` // 已启用通道数
	VisibleChannels int `json:"visible_channels"` // 可见通道数
}

// GetAgentChannelStats 获取代理商通道统计
func (s *AgentChannelService) GetAgentChannelStats(agentID int64) (*AgentChannelStats, error) {
	channels, err := s.agentChannelRepo.FindByAgentID(agentID)
	if err != nil {
		return nil, err
	}

	stats := &AgentChannelStats{
		TotalChannels: len(channels),
	}

	for _, ch := range channels {
		if ch.IsEnabled {
			stats.EnabledChannels++
		}
		if ch.IsVisible {
			stats.VisibleChannels++
		}
	}

	return stats, nil
}
