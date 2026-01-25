package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// ChannelService 通道服务
type ChannelService struct {
	channelRepo *repository.GormChannelRepository
}

// NewChannelService 创建通道服务
func NewChannelService(channelRepo *repository.GormChannelRepository) *ChannelService {
	return &ChannelService{
		channelRepo: channelRepo,
	}
}

// ChannelConfig 通道配置结构
type ChannelConfig struct {
	APIBaseURL string                       `json:"api_base_url"`
	PublicKey  string                       `json:"public_key"`
	RateTypes  []models.RateTypeDefinition  `json:"rate_types"`
}

// GetRateTypes 获取通道费率类型列表
func (s *ChannelService) GetRateTypes(channelID int64) ([]models.RateTypeDefinition, error) {
	channel, err := s.channelRepo.FindByID(channelID)
	if err != nil {
		return nil, fmt.Errorf("通道不存在: %d", channelID)
	}

	if channel.Config == "" {
		return []models.RateTypeDefinition{}, nil
	}

	var config ChannelConfig
	if err := json.Unmarshal([]byte(channel.Config), &config); err != nil {
		return nil, fmt.Errorf("解析通道配置失败: %w", err)
	}

	// 按 sort_order 排序
	sort.Slice(config.RateTypes, func(i, j int) bool {
		return config.RateTypes[i].SortOrder < config.RateTypes[j].SortOrder
	})

	return config.RateTypes, nil
}

// GetRateTypesByCode 根据通道编码获取费率类型列表
func (s *ChannelService) GetRateTypesByCode(channelCode string) ([]models.RateTypeDefinition, error) {
	channel, err := s.channelRepo.FindByCode(channelCode)
	if err != nil {
		return nil, fmt.Errorf("通道不存在: %s", channelCode)
	}

	return s.GetRateTypes(channel.ID)
}

// ValidateRateConfigs 验证费率配置是否在通道允许范围内
func (s *ChannelService) ValidateRateConfigs(channelID int64, rateConfigs models.RateConfigs) error {
	rateTypes, err := s.GetRateTypes(channelID)
	if err != nil {
		return err
	}

	// 构建费率类型映射
	rateTypeMap := make(map[string]models.RateTypeDefinition)
	for _, rt := range rateTypes {
		rateTypeMap[rt.Code] = rt
	}

	// 验证每个费率配置
	for code, config := range rateConfigs {
		rt, ok := rateTypeMap[code]
		if !ok {
			return fmt.Errorf("未知的费率类型: %s", code)
		}

		rate, err := strconv.ParseFloat(config.Rate, 64)
		if err != nil {
			return fmt.Errorf("%s费率格式错误: %s", rt.Name, config.Rate)
		}

		minRate, _ := strconv.ParseFloat(rt.MinRate, 64)
		maxRate, _ := strconv.ParseFloat(rt.MaxRate, 64)

		if rate < minRate {
			return fmt.Errorf("%s费率不能低于%.2f%%", rt.Name, minRate)
		}
		if rate > maxRate {
			return fmt.Errorf("%s费率不能高于%.2f%%", rt.Name, maxRate)
		}
	}

	return nil
}

// GetChannelByID 根据ID获取通道
func (s *ChannelService) GetChannelByID(channelID int64) (*models.Channel, error) {
	return s.channelRepo.FindByID(channelID)
}

// GetChannelByCode 根据编码获取通道
func (s *ChannelService) GetChannelByCode(channelCode string) (*models.Channel, error) {
	return s.channelRepo.FindByCode(channelCode)
}

// GetAllChannels 获取所有通道
func (s *ChannelService) GetAllChannels() ([]*models.Channel, error) {
	return s.channelRepo.FindAll()
}

// GetEnabledChannels 获取启用的通道
func (s *ChannelService) GetEnabledChannels() ([]*models.Channel, error) {
	return s.channelRepo.FindAllActive()
}
