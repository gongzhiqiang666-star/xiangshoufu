package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// WalletSplitService 钱包拆分配置服务
type WalletSplitService struct {
	splitConfigRepo   *repository.GormWalletSplitConfigRepository
	thresholdRepo     *repository.GormPolicyWithdrawThresholdRepository
	agentRepo         *repository.GormAgentRepository
	walletRepo        *repository.GormWalletRepository
	agentPolicyRepo   *repository.GormAgentPolicyRepository
	withdrawRepo      *repository.GormWithdrawRepository
}

// NewWalletSplitService 创建钱包拆分配置服务
func NewWalletSplitService(
	splitConfigRepo *repository.GormWalletSplitConfigRepository,
	thresholdRepo *repository.GormPolicyWithdrawThresholdRepository,
	agentRepo *repository.GormAgentRepository,
	walletRepo *repository.GormWalletRepository,
	agentPolicyRepo *repository.GormAgentPolicyRepository,
	withdrawRepo *repository.GormWithdrawRepository,
) *WalletSplitService {
	return &WalletSplitService{
		splitConfigRepo:   splitConfigRepo,
		thresholdRepo:     thresholdRepo,
		agentRepo:         agentRepo,
		walletRepo:        walletRepo,
		agentPolicyRepo:   agentPolicyRepo,
		withdrawRepo:      withdrawRepo,
	}
}

// ============================================================
// 拆分配置管理
// ============================================================

// GetSplitConfig 获取代理商钱包拆分配置
func (s *WalletSplitService) GetSplitConfig(agentID int64) (*models.AgentWalletSplitConfig, error) {
	config, err := s.splitConfigRepo.FindByAgentID(agentID)
	if err != nil {
		return nil, fmt.Errorf("查询拆分配置失败: %w", err)
	}

	// 如果没有配置，返回默认值
	if config == nil {
		return &models.AgentWalletSplitConfig{
			AgentID:        agentID,
			SplitByChannel: false,
		}, nil
	}

	return config, nil
}

// IsSplitByChannel 检查代理商是否按通道拆分
// 规则：如果代理商自己开启了拆分，或者任意上级开启了拆分，则返回true
func (s *WalletSplitService) IsSplitByChannel(agentID int64) (bool, error) {
	// 1. 先检查自己的配置
	config, err := s.splitConfigRepo.FindByAgentID(agentID)
	if err != nil {
		return false, fmt.Errorf("查询拆分配置失败: %w", err)
	}
	if config != nil && config.SplitByChannel {
		return true, nil
	}

	// 2. 检查上级链路是否有开启拆分的
	ancestors, err := s.agentRepo.FindAncestors(agentID)
	if err != nil {
		log.Printf("[WalletSplitService] 查询上级代理商失败: %v", err)
		return false, nil // 查询失败时默认不拆分
	}

	for _, ancestor := range ancestors {
		ancestorConfig, err := s.splitConfigRepo.FindByAgentID(ancestor.ID)
		if err != nil {
			continue
		}
		if ancestorConfig != nil && ancestorConfig.SplitByChannel {
			return true, nil // 上级开启了拆分，自动继承
		}
	}

	return false, nil
}

// SetSplitConfigRequest 设置拆分配置请求
type SetSplitConfigRequest struct {
	AgentID        int64 `json:"agent_id" binding:"required"`
	SplitByChannel bool  `json:"split_by_channel"`
	ConfiguredBy   int64 `json:"-"` // 配置人ID
}

// SetSplitConfig 设置代理商钱包拆分配置
// 权限：仅管理员或直接上级可设置
// 规则：
// 1. 一旦开启不可关闭
// 2. 有待审核/待打款的提现时不能开启
func (s *WalletSplitService) SetSplitConfig(req *SetSplitConfigRequest) error {
	// 获取现有配置
	existingConfig, err := s.splitConfigRepo.FindByAgentID(req.AgentID)
	if err != nil {
		return fmt.Errorf("查询现有配置失败: %w", err)
	}

	// 规则1：已开启的不能关闭
	if existingConfig != nil && existingConfig.SplitByChannel && !req.SplitByChannel {
		return errors.New("钱包拆分开关一旦开启不可关闭")
	}

	// 如果要开启拆分
	if req.SplitByChannel {
		// 规则2：检查是否有待处理的提现
		hasPending, err := s.hasPendingWithdraw(req.AgentID)
		if err != nil {
			return fmt.Errorf("检查待处理提现失败: %w", err)
		}
		if hasPending {
			return errors.New("存在待审核或待打款的提现申请，请等待所有提现完成后再开启拆分")
		}
	}

	// 创建或更新配置
	now := time.Now()
	if existingConfig == nil {
		config := &models.AgentWalletSplitConfig{
			AgentID:        req.AgentID,
			SplitByChannel: req.SplitByChannel,
			ConfiguredBy:   &req.ConfiguredBy,
			ConfiguredAt:   &now,
		}
		if err := s.splitConfigRepo.Create(config); err != nil {
			return fmt.Errorf("创建拆分配置失败: %w", err)
		}
	} else {
		existingConfig.SplitByChannel = req.SplitByChannel
		existingConfig.ConfiguredBy = &req.ConfiguredBy
		existingConfig.ConfiguredAt = &now
		if err := s.splitConfigRepo.Update(existingConfig); err != nil {
			return fmt.Errorf("更新拆分配置失败: %w", err)
		}
	}

	log.Printf("[WalletSplitService] 设置拆分配置: AgentID=%d, SplitByChannel=%v, ConfiguredBy=%d",
		req.AgentID, req.SplitByChannel, req.ConfiguredBy)

	return nil
}

// hasPendingWithdraw 检查是否有待处理的提现
func (s *WalletSplitService) hasPendingWithdraw(agentID int64) (bool, error) {
	if s.withdrawRepo == nil {
		return false, nil // 如果没有注入提现仓库，跳过检查
	}

	// 查询待审核(0)和待打款(1)的提现
	count, err := s.withdrawRepo.CountPendingByAgentID(agentID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ============================================================
// 提现门槛配置管理
// ============================================================

// GetWithdrawThresholds 获取政策模版的提现门槛配置
func (s *WalletSplitService) GetWithdrawThresholds(templateID int64) ([]*models.PolicyWithdrawThreshold, error) {
	thresholds, err := s.thresholdRepo.FindByTemplateID(templateID)
	if err != nil {
		return nil, fmt.Errorf("查询提现门槛失败: %w", err)
	}
	return thresholds, nil
}

// GetWithdrawThreshold 获取特定钱包类型和通道的提现门槛
// 如果未配置特定通道门槛，返回通用门槛(channel_id=0)
// 如果都未配置，返回默认值100元
func (s *WalletSplitService) GetWithdrawThreshold(templateID int64, walletType int16, channelID int64) (int64, error) {
	// 先查特定通道的门槛
	if channelID > 0 {
		threshold, err := s.thresholdRepo.FindByTemplateWalletAndChannel(templateID, walletType, channelID)
		if err != nil {
			return 0, fmt.Errorf("查询提现门槛失败: %w", err)
		}
		if threshold != nil {
			return threshold.ThresholdAmount, nil
		}
	}

	// 查通用门槛(channel_id=0)
	threshold, err := s.thresholdRepo.FindByTemplateWalletAndChannel(templateID, walletType, 0)
	if err != nil {
		return 0, fmt.Errorf("查询通用提现门槛失败: %w", err)
	}
	if threshold != nil {
		return threshold.ThresholdAmount, nil
	}

	// 返回默认值
	return getDefaultThreshold(walletType), nil
}

// getDefaultThreshold 获取默认提现门槛
func getDefaultThreshold(walletType int16) int64 {
	switch walletType {
	case models.WalletTypeProfit:
		return 10000 // 分润钱包默认100元
	case models.WalletTypeService:
		return 5000 // 服务费钱包默认50元
	case models.WalletTypeReward:
		return 10000 // 奖励钱包默认100元
	default:
		return 10000
	}
}

// SetWithdrawThresholdRequest 设置提现门槛请求
type SetWithdrawThresholdRequest struct {
	TemplateID      int64 `json:"template_id" binding:"required"`
	WalletType      int16 `json:"wallet_type" binding:"required,min=1,max=3"`
	ChannelID       int64 `json:"channel_id"` // 0表示通用门槛
	ThresholdAmount int64 `json:"threshold_amount" binding:"required,min=100"` // 最少1元
}

// SetWithdrawThreshold 设置提现门槛
func (s *WalletSplitService) SetWithdrawThreshold(req *SetWithdrawThresholdRequest) error {
	// 检查是否已存在
	existing, err := s.thresholdRepo.FindByTemplateWalletAndChannel(req.TemplateID, req.WalletType, req.ChannelID)
	if err != nil {
		return fmt.Errorf("查询现有门槛失败: %w", err)
	}

	if existing != nil {
		// 更新
		existing.ThresholdAmount = req.ThresholdAmount
		existing.UpdatedAt = time.Now()
		if err := s.thresholdRepo.Update(existing); err != nil {
			return fmt.Errorf("更新提现门槛失败: %w", err)
		}
	} else {
		// 创建
		threshold := &models.PolicyWithdrawThreshold{
			TemplateID:      req.TemplateID,
			WalletType:      req.WalletType,
			ChannelID:       req.ChannelID,
			ThresholdAmount: req.ThresholdAmount,
		}
		if err := s.thresholdRepo.Create(threshold); err != nil {
			return fmt.Errorf("创建提现门槛失败: %w", err)
		}
	}

	log.Printf("[WalletSplitService] 设置提现门槛: TemplateID=%d, WalletType=%d, ChannelID=%d, Amount=%d",
		req.TemplateID, req.WalletType, req.ChannelID, req.ThresholdAmount)

	return nil
}

// BatchSetWithdrawThresholds 批量设置提现门槛
func (s *WalletSplitService) BatchSetWithdrawThresholds(templateID int64, thresholds []*SetWithdrawThresholdRequest) error {
	for _, req := range thresholds {
		req.TemplateID = templateID
		if err := s.SetWithdrawThreshold(req); err != nil {
			return err
		}
	}
	return nil
}

// DeleteWithdrawThreshold 删除提现门槛配置
func (s *WalletSplitService) DeleteWithdrawThreshold(templateID int64, walletType int16, channelID int64) error {
	existing, err := s.thresholdRepo.FindByTemplateWalletAndChannel(templateID, walletType, channelID)
	if err != nil {
		return fmt.Errorf("查询提现门槛失败: %w", err)
	}
	if existing == nil {
		return nil // 不存在，视为成功
	}

	return s.thresholdRepo.DeleteByTemplateAndChannel(templateID, channelID)
}

// ============================================================
// 代理商提现门槛获取（通过政策模版）
// ============================================================

// GetAgentWithdrawThreshold 获取代理商特定钱包类型的提现门槛
// 通过代理商政策关联的模版获取门槛配置
func (s *WalletSplitService) GetAgentWithdrawThreshold(agentID int64, walletType int16, channelID int64) (int64, error) {
	// 获取代理商政策（先尝试获取特定通道的政策）
	var templateID int64

	if channelID > 0 {
		policy, err := s.agentPolicyRepo.FindByAgentAndChannel(agentID, channelID)
		if err == nil && policy != nil {
			templateID = policy.TemplateID
		}
	}

	// 如果没有特定通道政策，尝试获取默认政策
	if templateID == 0 {
		policies, err := s.agentPolicyRepo.FindByAgentID(agentID)
		if err != nil || len(policies) == 0 {
			// 没有政策，返回默认门槛
			return getDefaultThreshold(walletType), nil
		}
		templateID = policies[0].TemplateID
	}

	return s.GetWithdrawThreshold(templateID, walletType, channelID)
}
