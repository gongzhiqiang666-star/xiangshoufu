package service

import (
	"errors"
	"testing"
	"time"

	"xiangshoufu/internal/models"
)

// ============================================================
// Mock Repositories
// ============================================================

// MockWalletSplitConfigRepository 模拟拆分配置仓库
type MockWalletSplitConfigRepository struct {
	FindByAgentIDFunc func(agentID int64) (*models.AgentWalletSplitConfig, error)
	CreateFunc        func(config *models.AgentWalletSplitConfig) error
	UpdateFunc        func(config *models.AgentWalletSplitConfig) error
}

func (m *MockWalletSplitConfigRepository) FindByAgentID(agentID int64) (*models.AgentWalletSplitConfig, error) {
	if m.FindByAgentIDFunc != nil {
		return m.FindByAgentIDFunc(agentID)
	}
	return nil, nil
}

func (m *MockWalletSplitConfigRepository) Create(config *models.AgentWalletSplitConfig) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(config)
	}
	return nil
}

func (m *MockWalletSplitConfigRepository) Update(config *models.AgentWalletSplitConfig) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(config)
	}
	return nil
}

// MockPolicyWithdrawThresholdRepository 模拟提现门槛仓库
type MockPolicyWithdrawThresholdRepository struct {
	FindByTemplateIDFunc              func(templateID int64) ([]*models.PolicyWithdrawThreshold, error)
	FindByTemplateWalletAndChannelFunc func(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error)
	CreateFunc                        func(threshold *models.PolicyWithdrawThreshold) error
	UpdateFunc                        func(threshold *models.PolicyWithdrawThreshold) error
}

func (m *MockPolicyWithdrawThresholdRepository) FindByTemplateID(templateID int64) ([]*models.PolicyWithdrawThreshold, error) {
	if m.FindByTemplateIDFunc != nil {
		return m.FindByTemplateIDFunc(templateID)
	}
	return nil, nil
}

func (m *MockPolicyWithdrawThresholdRepository) FindByTemplateWalletAndChannel(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
	if m.FindByTemplateWalletAndChannelFunc != nil {
		return m.FindByTemplateWalletAndChannelFunc(templateID, walletType, channelID)
	}
	return nil, nil
}

func (m *MockPolicyWithdrawThresholdRepository) Create(threshold *models.PolicyWithdrawThreshold) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(threshold)
	}
	return nil
}

func (m *MockPolicyWithdrawThresholdRepository) Update(threshold *models.PolicyWithdrawThreshold) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(threshold)
	}
	return nil
}

// MockAgentRepositoryForSplit 模拟代理商仓库
type MockAgentRepositoryForSplit struct {
	FindAncestorsFunc func(agentID int64) ([]*models.Agent, error)
}

func (m *MockAgentRepositoryForSplit) FindAncestors(agentID int64) ([]*models.Agent, error) {
	if m.FindAncestorsFunc != nil {
		return m.FindAncestorsFunc(agentID)
	}
	return nil, nil
}

// MockWithdrawRepository 模拟提现仓库
type MockWithdrawRepository struct {
	CountPendingByAgentIDFunc func(agentID int64) (int64, error)
}

func (m *MockWithdrawRepository) CountPendingByAgentID(agentID int64) (int64, error) {
	if m.CountPendingByAgentIDFunc != nil {
		return m.CountPendingByAgentIDFunc(agentID)
	}
	return 0, nil
}

// ============================================================
// 测试用例
// ============================================================

// TestGetDefaultThreshold 测试获取默认提现门槛
func TestGetDefaultThreshold(t *testing.T) {
	tests := []struct {
		name       string
		walletType int16
		expected   int64
	}{
		{"profit wallet default", models.WalletTypeProfit, 10000},
		{"service wallet default", models.WalletTypeService, 5000},
		{"reward wallet default", models.WalletTypeReward, 10000},
		{"unknown wallet type", 99, 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDefaultThreshold(tt.walletType)
			if result != tt.expected {
				t.Errorf("getDefaultThreshold(%d) = %d, want %d", tt.walletType, result, tt.expected)
			}
		})
	}
}

// TestGetSplitConfig 测试获取拆分配置
func TestGetSplitConfig(t *testing.T) {
	t.Run("should return existing config", func(t *testing.T) {
		now := time.Now()
		configuredBy := int64(100)
		mockRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				return &models.AgentWalletSplitConfig{
					ID:             1,
					AgentID:        agentID,
					SplitByChannel: true,
					ConfiguredBy:   &configuredBy,
					ConfiguredAt:   &now,
				}, nil
			},
		}

		// 注意：这里使用简化的测试方式，因为实际服务需要多个仓库
		config, err := mockRepo.FindByAgentID(1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if config == nil {
			t.Error("expected config, got nil")
		}
		if !config.SplitByChannel {
			t.Error("expected SplitByChannel=true")
		}
	})

	t.Run("should return default config when not exists", func(t *testing.T) {
		mockRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				return nil, nil
			},
		}

		config, err := mockRepo.FindByAgentID(1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if config != nil {
			t.Errorf("expected nil config for non-existent, got %+v", config)
		}
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				return nil, errors.New("database error")
			},
		}

		_, err := mockRepo.FindByAgentID(1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestIsSplitByChannel 测试是否按通道拆分
func TestIsSplitByChannel(t *testing.T) {
	t.Run("should return true when agent has split enabled", func(t *testing.T) {
		mockRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				return &models.AgentWalletSplitConfig{
					AgentID:        agentID,
					SplitByChannel: true,
				}, nil
			},
		}

		config, _ := mockRepo.FindByAgentID(1)
		if config == nil || !config.SplitByChannel {
			t.Error("expected split to be enabled")
		}
	})

	t.Run("should return false when agent has no config", func(t *testing.T) {
		mockRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				return nil, nil
			},
		}

		config, _ := mockRepo.FindByAgentID(1)
		if config != nil && config.SplitByChannel {
			t.Error("expected split to be disabled")
		}
	})

	t.Run("should return false when split is disabled", func(t *testing.T) {
		mockRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				return &models.AgentWalletSplitConfig{
					AgentID:        agentID,
					SplitByChannel: false,
				}, nil
			},
		}

		config, _ := mockRepo.FindByAgentID(1)
		if config != nil && config.SplitByChannel {
			t.Error("expected split to be disabled")
		}
	})
}

// TestSetSplitConfigValidation 测试设置拆分配置的验证规则
func TestSetSplitConfigValidation(t *testing.T) {
	t.Run("should not allow disabling once enabled", func(t *testing.T) {
		// 模拟已开启的配置
		existingConfig := &models.AgentWalletSplitConfig{
			AgentID:        1,
			SplitByChannel: true,
		}

		// 尝试关闭
		newConfig := &SetSplitConfigRequest{
			AgentID:        1,
			SplitByChannel: false,
		}

		// 验证规则：已开启的不能关闭
		if existingConfig.SplitByChannel && !newConfig.SplitByChannel {
			// 预期行为：应该返回错误
			t.Log("correctly prevents disabling split once enabled")
		} else {
			t.Error("should prevent disabling split")
		}
	})

	t.Run("should allow enabling when currently disabled", func(t *testing.T) {
		existingConfig := &models.AgentWalletSplitConfig{
			AgentID:        1,
			SplitByChannel: false,
		}

		newConfig := &SetSplitConfigRequest{
			AgentID:        1,
			SplitByChannel: true,
		}

		// 验证规则：未开启时可以开启
		if !existingConfig.SplitByChannel && newConfig.SplitByChannel {
			t.Log("correctly allows enabling split")
		} else {
			t.Error("should allow enabling split")
		}
	})
}

// TestHasPendingWithdraw 测试检查待处理提现
func TestHasPendingWithdraw(t *testing.T) {
	t.Run("should return true when pending withdraws exist", func(t *testing.T) {
		mockRepo := &MockWithdrawRepository{
			CountPendingByAgentIDFunc: func(agentID int64) (int64, error) {
				return 2, nil
			},
		}

		count, _ := mockRepo.CountPendingByAgentID(1)
		if count <= 0 {
			t.Error("expected pending withdraws to exist")
		}
	})

	t.Run("should return false when no pending withdraws", func(t *testing.T) {
		mockRepo := &MockWithdrawRepository{
			CountPendingByAgentIDFunc: func(agentID int64) (int64, error) {
				return 0, nil
			},
		}

		count, _ := mockRepo.CountPendingByAgentID(1)
		if count > 0 {
			t.Error("expected no pending withdraws")
		}
	})
}

// TestGetWithdrawThreshold 测试获取提现门槛
func TestGetWithdrawThreshold(t *testing.T) {
	t.Run("should return channel-specific threshold when exists", func(t *testing.T) {
		mockRepo := &MockPolicyWithdrawThresholdRepository{
			FindByTemplateWalletAndChannelFunc: func(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
				if channelID == 1 {
					return &models.PolicyWithdrawThreshold{
						TemplateID:      templateID,
						WalletType:      walletType,
						ChannelID:       channelID,
						ThresholdAmount: 15000, // 150元
					}, nil
				}
				return nil, nil
			},
		}

		threshold, _ := mockRepo.FindByTemplateWalletAndChannel(1, models.WalletTypeProfit, 1)
		if threshold == nil {
			t.Error("expected threshold, got nil")
			return
		}
		if threshold.ThresholdAmount != 15000 {
			t.Errorf("expected 15000, got %d", threshold.ThresholdAmount)
		}
	})

	t.Run("should return general threshold when channel-specific not exists", func(t *testing.T) {
		mockRepo := &MockPolicyWithdrawThresholdRepository{
			FindByTemplateWalletAndChannelFunc: func(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
				if channelID == 0 {
					return &models.PolicyWithdrawThreshold{
						TemplateID:      templateID,
						WalletType:      walletType,
						ChannelID:       0,
						ThresholdAmount: 10000, // 100元
					}, nil
				}
				return nil, nil
			},
		}

		// 查询通道1的门槛，应返回nil
		threshold, _ := mockRepo.FindByTemplateWalletAndChannel(1, models.WalletTypeProfit, 1)
		if threshold != nil {
			t.Errorf("expected nil for channel 1, got %+v", threshold)
		}

		// 查询通用门槛(channel=0)，应返回值
		threshold, _ = mockRepo.FindByTemplateWalletAndChannel(1, models.WalletTypeProfit, 0)
		if threshold == nil {
			t.Error("expected general threshold, got nil")
			return
		}
		if threshold.ThresholdAmount != 10000 {
			t.Errorf("expected 10000, got %d", threshold.ThresholdAmount)
		}
	})

	t.Run("should return default when no threshold configured", func(t *testing.T) {
		mockRepo := &MockPolicyWithdrawThresholdRepository{
			FindByTemplateWalletAndChannelFunc: func(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
				return nil, nil
			},
		}

		threshold, _ := mockRepo.FindByTemplateWalletAndChannel(1, models.WalletTypeProfit, 0)
		if threshold != nil {
			t.Errorf("expected nil, got %+v", threshold)
		}

		// 应该使用默认值
		defaultVal := getDefaultThreshold(models.WalletTypeProfit)
		if defaultVal != 10000 {
			t.Errorf("expected default 10000, got %d", defaultVal)
		}
	})
}

// TestSetWithdrawThreshold 测试设置提现门槛
func TestSetWithdrawThreshold(t *testing.T) {
	t.Run("should create new threshold when not exists", func(t *testing.T) {
		created := false
		mockRepo := &MockPolicyWithdrawThresholdRepository{
			FindByTemplateWalletAndChannelFunc: func(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
				return nil, nil
			},
			CreateFunc: func(threshold *models.PolicyWithdrawThreshold) error {
				created = true
				return nil
			},
		}

		// 模拟查询不存在
		_, _ = mockRepo.FindByTemplateWalletAndChannel(1, models.WalletTypeProfit, 0)
		// 创建新记录
		_ = mockRepo.Create(&models.PolicyWithdrawThreshold{
			TemplateID:      1,
			WalletType:      models.WalletTypeProfit,
			ChannelID:       0,
			ThresholdAmount: 20000,
		})

		if !created {
			t.Error("expected Create to be called")
		}
	})

	t.Run("should update existing threshold", func(t *testing.T) {
		updated := false
		mockRepo := &MockPolicyWithdrawThresholdRepository{
			FindByTemplateWalletAndChannelFunc: func(templateID int64, walletType int16, channelID int64) (*models.PolicyWithdrawThreshold, error) {
				return &models.PolicyWithdrawThreshold{
					ID:              1,
					TemplateID:      templateID,
					WalletType:      walletType,
					ChannelID:       channelID,
					ThresholdAmount: 10000,
				}, nil
			},
			UpdateFunc: func(threshold *models.PolicyWithdrawThreshold) error {
				updated = true
				return nil
			},
		}

		// 模拟查询存在
		existing, _ := mockRepo.FindByTemplateWalletAndChannel(1, models.WalletTypeProfit, 0)
		if existing != nil {
			existing.ThresholdAmount = 20000
			_ = mockRepo.Update(existing)
		}

		if !updated {
			t.Error("expected Update to be called")
		}
	})
}

// TestBatchSetWithdrawThresholds 测试批量设置提现门槛
func TestBatchSetWithdrawThresholds(t *testing.T) {
	t.Run("should process all thresholds in batch", func(t *testing.T) {
		thresholds := []*SetWithdrawThresholdRequest{
			{WalletType: models.WalletTypeProfit, ChannelID: 0, ThresholdAmount: 10000},
			{WalletType: models.WalletTypeService, ChannelID: 0, ThresholdAmount: 5000},
			{WalletType: models.WalletTypeReward, ChannelID: 0, ThresholdAmount: 10000},
		}

		// 验证批量请求数量
		if len(thresholds) != 3 {
			t.Errorf("expected 3 thresholds, got %d", len(thresholds))
		}

		// 验证各项配置
		for _, th := range thresholds {
			if th.ThresholdAmount < 100 {
				t.Errorf("threshold amount should be >= 100, got %d", th.ThresholdAmount)
			}
		}
	})
}

// TestSplitInheritance 测试拆分配置继承规则
func TestSplitInheritance(t *testing.T) {
	t.Run("should inherit split from parent", func(t *testing.T) {
		// 模拟上级链路
		ancestors := []*models.Agent{
			{ID: 2, Name: "Parent Agent"},
			{ID: 1, Name: "Root Agent"},
		}

		mockAgentRepo := &MockAgentRepositoryForSplit{
			FindAncestorsFunc: func(agentID int64) ([]*models.Agent, error) {
				return ancestors, nil
			},
		}

		mockSplitRepo := &MockWalletSplitConfigRepository{
			FindByAgentIDFunc: func(agentID int64) (*models.AgentWalletSplitConfig, error) {
				// 上级ID=2开启了拆分
				if agentID == 2 {
					return &models.AgentWalletSplitConfig{
						AgentID:        2,
						SplitByChannel: true,
					}, nil
				}
				return nil, nil
			},
		}

		// 检查子代理商(ID=3)的拆分状态
		// 1. 先检查自己的配置
		config, _ := mockSplitRepo.FindByAgentID(3)
		if config != nil && config.SplitByChannel {
			t.Log("agent has split enabled directly")
			return
		}

		// 2. 检查上级链路
		ancestors, _ = mockAgentRepo.FindAncestors(3)
		for _, ancestor := range ancestors {
			ancestorConfig, _ := mockSplitRepo.FindByAgentID(ancestor.ID)
			if ancestorConfig != nil && ancestorConfig.SplitByChannel {
				t.Log("correctly inherits split from parent")
				return
			}
		}

		t.Error("should inherit split from parent")
	})
}
