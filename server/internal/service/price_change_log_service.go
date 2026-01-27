package service

import (
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// PriceChangeLogService 调价记录服务
type PriceChangeLogService struct {
	repo repository.PriceChangeLogRepository
}

// NewPriceChangeLogService 创建调价记录服务
func NewPriceChangeLogService(repo repository.PriceChangeLogRepository) *PriceChangeLogService {
	return &PriceChangeLogService{repo: repo}
}

// GetByID 根据ID获取调价记录
func (s *PriceChangeLogService) GetByID(id int64) (*models.PriceChangeLog, error) {
	return s.repo.GetByID(id)
}

// List 获取调价记录列表
func (s *PriceChangeLogService) List(req *models.PriceChangeLogListRequest) (*models.PriceChangeLogListResponse, error) {
	logs, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]models.PriceChangeLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, models.PriceChangeLogItem{
			ID:             l.ID,
			AgentID:        l.AgentID,
			ChannelID:      l.ChannelID,
			ChangeType:     l.ChangeType,
			ChangeTypeName: models.ChangeTypeName(l.ChangeType),
			ConfigType:     l.ConfigType,
			ConfigTypeName: models.ConfigTypeName(l.ConfigType),
			FieldName:      l.FieldName,
			OldValue:       l.OldValue,
			NewValue:       l.NewValue,
			ChangeSummary:  l.ChangeSummary,
			OperatorName:   l.OperatorName,
			Source:         l.Source,
			CreatedAt:      l.CreatedAt,
		})
	}

	return &models.PriceChangeLogListResponse{
		List:  items,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// ListByAgent 按代理商获取调价记录
func (s *PriceChangeLogService) ListByAgent(agentID int64, page, pageSize int) (*models.PriceChangeLogListResponse, error) {
	logs, total, err := s.repo.ListByAgent(agentID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]models.PriceChangeLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, models.PriceChangeLogItem{
			ID:             l.ID,
			AgentID:        l.AgentID,
			ChannelID:      l.ChannelID,
			ChangeType:     l.ChangeType,
			ChangeTypeName: models.ChangeTypeName(l.ChangeType),
			ConfigType:     l.ConfigType,
			ConfigTypeName: models.ConfigTypeName(l.ConfigType),
			FieldName:      l.FieldName,
			OldValue:       l.OldValue,
			NewValue:       l.NewValue,
			ChangeSummary:  l.ChangeSummary,
			OperatorName:   l.OperatorName,
			Source:         l.Source,
			CreatedAt:      l.CreatedAt,
		})
	}

	return &models.PriceChangeLogListResponse{
		List:  items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// ListBySettlementPrice 按结算价获取调价记录
func (s *PriceChangeLogService) ListBySettlementPrice(settlementPriceID int64, page, pageSize int) (*models.PriceChangeLogListResponse, error) {
	logs, total, err := s.repo.ListBySettlementPrice(settlementPriceID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]models.PriceChangeLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, models.PriceChangeLogItem{
			ID:             l.ID,
			AgentID:        l.AgentID,
			ChannelID:      l.ChannelID,
			ChangeType:     l.ChangeType,
			ChangeTypeName: models.ChangeTypeName(l.ChangeType),
			ConfigType:     l.ConfigType,
			ConfigTypeName: models.ConfigTypeName(l.ConfigType),
			FieldName:      l.FieldName,
			OldValue:       l.OldValue,
			NewValue:       l.NewValue,
			ChangeSummary:  l.ChangeSummary,
			OperatorName:   l.OperatorName,
			Source:         l.Source,
			CreatedAt:      l.CreatedAt,
		})
	}

	return &models.PriceChangeLogListResponse{
		List:  items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// ListByRewardSetting 按奖励配置获取调价记录
func (s *PriceChangeLogService) ListByRewardSetting(rewardSettingID int64, page, pageSize int) (*models.PriceChangeLogListResponse, error) {
	logs, total, err := s.repo.ListByRewardSetting(rewardSettingID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]models.PriceChangeLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, models.PriceChangeLogItem{
			ID:             l.ID,
			AgentID:        l.AgentID,
			ChannelID:      l.ChannelID,
			ChangeType:     l.ChangeType,
			ChangeTypeName: models.ChangeTypeName(l.ChangeType),
			ConfigType:     l.ConfigType,
			ConfigTypeName: models.ConfigTypeName(l.ConfigType),
			FieldName:      l.FieldName,
			OldValue:       l.OldValue,
			NewValue:       l.NewValue,
			ChangeSummary:  l.ChangeSummary,
			OperatorName:   l.OperatorName,
			Source:         l.Source,
			CreatedAt:      l.CreatedAt,
		})
	}

	return &models.PriceChangeLogListResponse{
		List:  items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}
