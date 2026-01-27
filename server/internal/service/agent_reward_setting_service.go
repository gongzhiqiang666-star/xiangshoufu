package service

import (
	"encoding/json"
	"fmt"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// AgentRewardSettingService 代理商奖励配置服务
type AgentRewardSettingService struct {
	repo          repository.AgentRewardSettingRepository
	changeLogRepo repository.PriceChangeLogRepository
}

// NewAgentRewardSettingService 创建代理商奖励配置服务
func NewAgentRewardSettingService(
	repo repository.AgentRewardSettingRepository,
	changeLogRepo repository.PriceChangeLogRepository,
) *AgentRewardSettingService {
	return &AgentRewardSettingService{
		repo:          repo,
		changeLogRepo: changeLogRepo,
	}
}

// CreateFromTemplate 从模板创建奖励配置
func (s *AgentRewardSettingService) CreateFromTemplate(
	agentID int64,
	templateID *int64,
	rewardAmount int64,
	activationRewards models.ActivationRewards,
	operatorID int64,
	operatorName string,
	source string,
) (*models.AgentRewardSetting, error) {
	// 检查是否已存在
	existing, err := s.repo.GetByAgentID(agentID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("代理商奖励配置已存在")
	}

	// 创建奖励配置
	setting := &models.AgentRewardSetting{
		AgentID:           agentID,
		TemplateID:        templateID,
		RewardAmount:      rewardAmount,
		ActivationRewards: activationRewards,
		Version:           1,
		Status:            1,
		EffectiveAt:       time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         &operatorID,
	}

	err = s.repo.Create(setting)
	if err != nil {
		return nil, fmt.Errorf("创建奖励配置失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLog(setting, nil, models.ChangeTypeInit, operatorID, operatorName, source, "", "初始化奖励配置")

	return setting, nil
}

// UpdateRewardAmount 更新奖励金额
func (s *AgentRewardSettingService) UpdateRewardAmount(
	id int64,
	rewardAmount int64,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.AgentRewardSetting, error) {
	setting, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("奖励配置不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(setting)

	// 更新奖励金额
	setting.RewardAmount = rewardAmount
	setting.Version++
	setting.UpdatedBy = &operatorID

	err = s.repo.Update(setting)
	if err != nil {
		return nil, fmt.Errorf("更新奖励金额失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(setting, snapshotBefore, models.ChangeTypeActivation, operatorID, operatorName, source, ipAddress, "奖励金额", "奖励金额调整")

	return setting, nil
}

// UpdateActivationReward 更新激活奖励
func (s *AgentRewardSettingService) UpdateActivationReward(
	id int64,
	req *models.UpdateActivationRewardRequest,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.AgentRewardSetting, error) {
	setting, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("奖励配置不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(setting)

	// 更新激活奖励
	setting.ActivationRewards = req.ActivationRewards
	setting.Version++
	setting.UpdatedBy = &operatorID

	err = s.repo.Update(setting)
	if err != nil {
		return nil, fmt.Errorf("更新激活奖励失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(setting, snapshotBefore, models.ChangeTypeActivation, operatorID, operatorName, source, ipAddress, "激活奖励", "激活奖励调整")

	return setting, nil
}

// GetByID 根据ID获取奖励配置
func (s *AgentRewardSettingService) GetByID(id int64) (*models.AgentRewardSetting, error) {
	return s.repo.GetByID(id)
}

// GetByAgentID 根据代理商ID获取奖励配置
func (s *AgentRewardSettingService) GetByAgentID(agentID int64) (*models.AgentRewardSetting, error) {
	return s.repo.GetByAgentID(agentID)
}

// List 获取奖励配置列表
func (s *AgentRewardSettingService) List(page, pageSize int) ([]models.AgentRewardSetting, int64, error) {
	return s.repo.List(page, pageSize)
}

// GetAgentActivationReward 获取代理商激活奖励配置（供奖励计算使用）
func (s *AgentRewardSettingService) GetAgentActivationReward(agentID int64, registerDays int, tradeAmount int64) (*models.ActivationRewardItem, error) {
	setting, err := s.repo.GetByAgentID(agentID)
	if err != nil {
		return nil, fmt.Errorf("获取奖励配置失败: %w", err)
	}

	// 查找匹配的激活奖励
	for _, reward := range setting.ActivationRewards {
		if registerDays >= reward.MinRegisterDays && registerDays <= reward.MaxRegisterDays {
			if tradeAmount >= reward.TargetAmount {
				return &reward, nil
			}
		}
	}

	return nil, nil
}

// GetAgentRewardAmount 获取代理商奖励金额（差额分配模式）
func (s *AgentRewardSettingService) GetAgentRewardAmount(agentID int64) (int64, error) {
	setting, err := s.repo.GetByAgentID(agentID)
	if err != nil {
		return 0, fmt.Errorf("获取奖励配置失败: %w", err)
	}

	return setting.RewardAmount, nil
}

// createSnapshot 创建快照
func (s *AgentRewardSettingService) createSnapshot(setting *models.AgentRewardSetting) models.JSONMap {
	snapshot := models.JSONMap{
		"id":                 setting.ID,
		"agent_id":           setting.AgentID,
		"template_id":        setting.TemplateID,
		"reward_amount":      setting.RewardAmount,
		"activation_rewards": setting.ActivationRewards,
		"version":            setting.Version,
	}
	return snapshot
}

// createChangeLog 创建调价日志
func (s *AgentRewardSettingService) createChangeLog(
	setting *models.AgentRewardSetting,
	snapshotBefore models.JSONMap,
	changeType models.ChangeType,
	operatorID int64,
	operatorName string,
	source string,
	fieldName string,
	summary string,
) {
	snapshotAfter := s.createSnapshot(setting)

	log := &models.PriceChangeLog{
		AgentID:         setting.AgentID,
		RewardSettingID: &setting.ID,
		ChangeType:      changeType,
		ConfigType:      models.ConfigTypeReward,
		FieldName:       fieldName,
		ChangeSummary:   summary,
		SnapshotBefore:  snapshotBefore,
		SnapshotAfter:   snapshotAfter,
		OperatorType:    models.OperatorTypeAdmin,
		OperatorID:      operatorID,
		OperatorName:    operatorName,
		Source:          source,
		CreatedAt:       time.Now(),
	}

	s.changeLogRepo.Create(log)
}

// createChangeLogWithSnapshot 创建带快照的调价日志
func (s *AgentRewardSettingService) createChangeLogWithSnapshot(
	setting *models.AgentRewardSetting,
	snapshotBefore models.JSONMap,
	changeType models.ChangeType,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
	fieldName string,
	summary string,
) {
	snapshotAfter := s.createSnapshot(setting)

	// 计算变更内容
	oldValue, _ := json.Marshal(snapshotBefore)
	newValue, _ := json.Marshal(snapshotAfter)

	log := &models.PriceChangeLog{
		AgentID:         setting.AgentID,
		RewardSettingID: &setting.ID,
		ChangeType:      changeType,
		ConfigType:      models.ConfigTypeReward,
		FieldName:       fieldName,
		OldValue:        string(oldValue),
		NewValue:        string(newValue),
		ChangeSummary:   summary,
		SnapshotBefore:  snapshotBefore,
		SnapshotAfter:   snapshotAfter,
		OperatorType:    models.OperatorTypeAdmin,
		OperatorID:      operatorID,
		OperatorName:    operatorName,
		Source:          source,
		IPAddress:       ipAddress,
		CreatedAt:       time.Now(),
	}

	s.changeLogRepo.Create(log)
}
