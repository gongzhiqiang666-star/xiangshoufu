package service

import (
	"context"
	"fmt"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// ChannelConfigService 通道配置服务接口
type ChannelConfigService interface {
	// 费率配置
	GetRateConfigs(ctx context.Context, channelID int64) ([]models.ChannelRateConfigResponse, error)
	CreateRateConfig(ctx context.Context, channelID int64, req *models.CreateChannelRateConfigRequest) (*models.ChannelRateConfig, error)
	UpdateRateConfig(ctx context.Context, channelID int64, configID int64, req *models.UpdateChannelRateConfigRequest) error
	DeleteRateConfig(ctx context.Context, channelID int64, configID int64) error

	// 押金档位
	GetDepositTiers(ctx context.Context, channelID int64) ([]models.ChannelDepositTierResponse, error)
	UpdateDepositTier(ctx context.Context, channelID int64, tierID int64, req *models.UpdateChannelDepositTierRequest) error

	// 流量费返现档位
	GetSimCashbackTiers(ctx context.Context, channelID int64) ([]models.ChannelSimCashbackTierResponse, error)
	BatchSetSimCashbackTiers(ctx context.Context, channelID int64, req *models.BatchSetSimCashbackTiersRequest) error

	// 通道完整配置
	GetFullConfig(ctx context.Context, channelID int64) (*models.ChannelFullConfig, error)

	// 校验方法
	ValidateRateForTemplate(ctx context.Context, channelID int64, rateCode string, rate string) error
	ValidateRateForSettlement(ctx context.Context, channelID int64, rateCode string, rate string, upperRate string) error
	ValidateCashbackForTemplate(ctx context.Context, channelID int64, depositAmount int64, cashback int64) error
	ValidateSimCashbackForTemplate(ctx context.Context, channelID int64, tierOrder int, cashback int64) error

	// 高调费率和P+0加价校验
	ValidateHighRateAndD0Extra(ctx context.Context, channelID int64, highRateConfigs models.HighRateConfigs, d0ExtraConfigs models.D0ExtraConfigs) error

	// 配置变更影响检查
	CheckRateConfigChangeImpact(ctx context.Context, channelID int64, configID int64, req *models.CheckRateConfigChangeImpactRequest) (*models.ConfigChangeImpact, error)
}

// channelConfigService 通道配置服务实现
type channelConfigService struct {
	repo repository.ChannelConfigRepository
}

// NewChannelConfigService 创建通道配置服务实例
func NewChannelConfigService(repo repository.ChannelConfigRepository) ChannelConfigService {
	return &channelConfigService{repo: repo}
}

// ============================================================
// 费率配置方法
// ============================================================

// GetRateConfigs 获取费率配置列表
func (s *channelConfigService) GetRateConfigs(ctx context.Context, channelID int64) ([]models.ChannelRateConfigResponse, error) {
	configs, err := s.repo.GetRateConfigs(ctx, channelID)
	if err != nil {
		return nil, err
	}

	result := make([]models.ChannelRateConfigResponse, len(configs))
	for i, c := range configs {
		result[i] = models.ChannelRateConfigResponse{
			ID:          c.ID,
			ChannelID:   c.ChannelID,
			RateCode:    c.RateCode,
			RateName:    c.RateName,
			MinRate:     c.MinRate,
			MaxRate:     c.MaxRate,
			DefaultRate: c.DefaultRate,
			MaxHighRate: c.MaxHighRate,
			MaxD0Extra:  c.MaxD0Extra,
			SortOrder:   c.SortOrder,
			Status:      c.Status,
		}
	}
	return result, nil
}

// CreateRateConfig 创建费率配置
func (s *channelConfigService) CreateRateConfig(ctx context.Context, channelID int64, req *models.CreateChannelRateConfigRequest) (*models.ChannelRateConfig, error) {
	// 检查是否已存在相同编码
	existing, _ := s.repo.GetRateConfigByCode(ctx, channelID, req.RateCode)
	if existing != nil {
		return nil, fmt.Errorf("费率编码 %s 已存在", req.RateCode)
	}

	// 校验费率范围
	minRate := models.ParseRateToFloat(req.MinRate)
	maxRate := models.ParseRateToFloat(req.MaxRate)
	if minRate > maxRate {
		return nil, fmt.Errorf("最低费率不能大于最高费率")
	}

	config := &models.ChannelRateConfig{
		ChannelID:   channelID,
		RateCode:    req.RateCode,
		RateName:    req.RateName,
		MinRate:     req.MinRate,
		MaxRate:     req.MaxRate,
		DefaultRate: req.DefaultRate,
		SortOrder:   req.SortOrder,
		Status:      1,
	}

	if err := s.repo.CreateRateConfig(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

// UpdateRateConfig 更新费率配置
func (s *channelConfigService) UpdateRateConfig(ctx context.Context, channelID int64, configID int64, req *models.UpdateChannelRateConfigRequest) error {
	configs, err := s.repo.GetRateConfigs(ctx, channelID)
	if err != nil {
		return err
	}

	var config *models.ChannelRateConfig
	for i := range configs {
		if configs[i].ID == configID {
			config = &configs[i]
			break
		}
	}
	if config == nil {
		return fmt.Errorf("费率配置不存在")
	}

	if req.RateName != "" {
		config.RateName = req.RateName
	}
	if req.MinRate != "" {
		config.MinRate = req.MinRate
	}
	if req.MaxRate != "" {
		config.MaxRate = req.MaxRate
	}
	if req.DefaultRate != "" {
		config.DefaultRate = req.DefaultRate
	}
	if req.MaxHighRate != nil {
		config.MaxHighRate = req.MaxHighRate
	}
	if req.MaxD0Extra != nil {
		config.MaxD0Extra = req.MaxD0Extra
	}
	if req.SortOrder != 0 {
		config.SortOrder = req.SortOrder
	}
	if req.Status != nil {
		config.Status = *req.Status
	}

	return s.repo.UpdateRateConfig(ctx, config)
}

// DeleteRateConfig 删除费率配置
func (s *channelConfigService) DeleteRateConfig(ctx context.Context, channelID int64, configID int64) error {
	return s.repo.DeleteRateConfig(ctx, configID)
}

// ============================================================
// 押金档位方法
// ============================================================

// GetDepositTiers 获取押金档位列表
func (s *channelConfigService) GetDepositTiers(ctx context.Context, channelID int64) ([]models.ChannelDepositTierResponse, error) {
	tiers, err := s.repo.GetDepositTiers(ctx, channelID)
	if err != nil {
		return nil, err
	}

	result := make([]models.ChannelDepositTierResponse, len(tiers))
	for i, t := range tiers {
		result[i] = models.ChannelDepositTierResponse{
			ID:                t.ID,
			ChannelID:         t.ChannelID,
			BrandCode:         t.BrandCode,
			TierCode:          t.TierCode,
			DepositAmount:     t.DepositAmount,
			TierName:          t.TierName,
			MaxCashbackAmount: t.MaxCashbackAmount,
			DefaultCashback:   t.DefaultCashback,
			SortOrder:         t.SortOrder,
			Status:            t.Status,
		}
	}
	return result, nil
}

// UpdateDepositTier 更新押金档位
func (s *channelConfigService) UpdateDepositTier(ctx context.Context, channelID int64, tierID int64, req *models.UpdateChannelDepositTierRequest) error {
	tier, err := s.repo.GetDepositTierByID(ctx, tierID)
	if err != nil {
		return err
	}

	if tier.ChannelID != channelID {
		return fmt.Errorf("押金档位不属于该通道")
	}

	tier.MaxCashbackAmount = req.MaxCashbackAmount
	tier.DefaultCashback = req.DefaultCashback
	if req.Status != nil {
		tier.Status = *req.Status
	}

	return s.repo.UpdateDepositTier(ctx, tier)
}

// ============================================================
// 流量费返现档位方法
// ============================================================

// GetSimCashbackTiers 获取流量费返现档位列表
func (s *channelConfigService) GetSimCashbackTiers(ctx context.Context, channelID int64) ([]models.ChannelSimCashbackTierResponse, error) {
	tiers, err := s.repo.GetSimCashbackTiers(ctx, channelID, "")
	if err != nil {
		return nil, err
	}

	result := make([]models.ChannelSimCashbackTierResponse, len(tiers))
	for i, t := range tiers {
		result[i] = models.ChannelSimCashbackTierResponse{
			ID:                t.ID,
			ChannelID:         t.ChannelID,
			BrandCode:         t.BrandCode,
			TierOrder:         t.TierOrder,
			TierName:          t.TierName,
			IsLastTier:        t.IsLastTier,
			MaxCashbackAmount: t.MaxCashbackAmount,
			DefaultCashback:   t.DefaultCashback,
			SimFeeAmount:      t.SimFeeAmount,
			Status:            t.Status,
		}
	}
	return result, nil
}

// BatchSetSimCashbackTiers 批量设置流量费返现档位
func (s *channelConfigService) BatchSetSimCashbackTiers(ctx context.Context, channelID int64, req *models.BatchSetSimCashbackTiersRequest) error {
	// 校验档位数据
	hasLastTier := false
	for _, t := range req.Tiers {
		if t.IsLastTier {
			if hasLastTier {
				return fmt.Errorf("只能有一个最后档位")
			}
			hasLastTier = true
		}
		if t.DefaultCashback > t.MaxCashbackAmount {
			return fmt.Errorf("档位 %d 的默认返现不能超过返现上限", t.TierOrder)
		}
	}

	// 转换为模型
	tiers := make([]models.ChannelSimCashbackTier, len(req.Tiers))
	for i, t := range req.Tiers {
		tiers[i] = models.ChannelSimCashbackTier{
			ChannelID:         channelID,
			BrandCode:         "",
			TierOrder:         t.TierOrder,
			TierName:          t.TierName,
			IsLastTier:        t.IsLastTier,
			MaxCashbackAmount: t.MaxCashbackAmount,
			DefaultCashback:   t.DefaultCashback,
			SimFeeAmount:      t.SimFeeAmount,
			Status:            1,
		}
	}

	return s.repo.BatchSetSimCashbackTiers(ctx, channelID, "", tiers)
}

// ============================================================
// 通道完整配置方法
// ============================================================

// GetFullConfig 获取通道完整配置
func (s *channelConfigService) GetFullConfig(ctx context.Context, channelID int64) (*models.ChannelFullConfig, error) {
	return s.repo.GetFullConfig(ctx, channelID)
}

// ============================================================
// 校验方法
// ============================================================

// ValidateRateForTemplate 校验政策模板费率是否在通道允许范围内
func (s *channelConfigService) ValidateRateForTemplate(ctx context.Context, channelID int64, rateCode string, rate string) error {
	config, err := s.repo.GetRateConfigByCode(ctx, channelID, rateCode)
	if err != nil {
		return fmt.Errorf("获取通道费率配置失败: %w", err)
	}

	return models.ValidateRateRange(rate, config.MinRate, config.MaxRate)
}

// ValidateRateForSettlement 校验结算价费率（必须 >= 上级费率 且在通道范围内）
func (s *channelConfigService) ValidateRateForSettlement(ctx context.Context, channelID int64, rateCode string, rate string, upperRate string) error {
	config, err := s.repo.GetRateConfigByCode(ctx, channelID, rateCode)
	if err != nil {
		return fmt.Errorf("获取通道费率配置失败: %w", err)
	}

	return models.ValidateSettlementRate(rate, upperRate, config.MinRate, config.MaxRate)
}

// ValidateCashbackForTemplate 校验政策模板押金返现是否在通道允许范围内
func (s *channelConfigService) ValidateCashbackForTemplate(ctx context.Context, channelID int64, depositAmount int64, cashback int64) error {
	tiers, err := s.repo.GetDepositTiers(ctx, channelID)
	if err != nil {
		return fmt.Errorf("获取通道押金档位失败: %w", err)
	}

	for _, tier := range tiers {
		if tier.DepositAmount == depositAmount {
			return models.ValidateCashbackAmount(cashback, tier.MaxCashbackAmount)
		}
	}

	return fmt.Errorf("未找到押金金额 %d 对应的档位配置", depositAmount)
}

// ValidateSimCashbackForTemplate 校验政策模板流量费返现是否在通道允许范围内
func (s *channelConfigService) ValidateSimCashbackForTemplate(ctx context.Context, channelID int64, tierOrder int, cashback int64) error {
	tiers, err := s.repo.GetSimCashbackTiers(ctx, channelID, "")
	if err != nil {
		return fmt.Errorf("获取通道流量费返现档位失败: %w", err)
	}

	// 查找对应档位
	for _, tier := range tiers {
		if tier.TierOrder == tierOrder {
			return models.ValidateCashbackAmount(cashback, tier.MaxCashbackAmount)
		}
	}

	// 超出配置档位，使用最后一档
	for _, tier := range tiers {
		if tier.IsLastTier {
			return models.ValidateCashbackAmount(cashback, tier.MaxCashbackAmount)
		}
	}

	return fmt.Errorf("未找到档位 %d 对应的流量费返现配置", tierOrder)
}

// ValidateHighRateAndD0Extra 校验高调费率和P+0加价是否在通道上限内
func (s *channelConfigService) ValidateHighRateAndD0Extra(ctx context.Context, channelID int64, highRateConfigs models.HighRateConfigs, d0ExtraConfigs models.D0ExtraConfigs) error {
	configs, err := s.repo.GetRateConfigs(ctx, channelID)
	if err != nil {
		return fmt.Errorf("获取通道费率配置失败: %w", err)
	}

	// 构建费率配置映射
	configMap := make(map[string]*models.ChannelRateConfig)
	for i := range configs {
		configMap[configs[i].RateCode] = &configs[i]
	}

	// 校验高调费率
	for rateCode, hrConfig := range highRateConfigs {
		channelConfig, ok := configMap[rateCode]
		if !ok {
			continue // 未找到对应通道配置，跳过
		}
		if err := models.ValidateHighRate(hrConfig.Rate, channelConfig.MaxHighRate, channelConfig.RateName); err != nil {
			return err
		}
	}

	// 校验P+0加价
	for rateCode, d0Config := range d0ExtraConfigs {
		channelConfig, ok := configMap[rateCode]
		if !ok {
			continue // 未找到对应通道配置，跳过
		}
		if err := models.ValidateD0Extra(d0Config.ExtraFee, channelConfig.MaxD0Extra, channelConfig.RateName); err != nil {
			return err
		}
	}

	return nil
}

// CheckRateConfigChangeImpact 检查费率配置变更影响
func (s *channelConfigService) CheckRateConfigChangeImpact(ctx context.Context, channelID int64, configID int64, req *models.CheckRateConfigChangeImpactRequest) (*models.ConfigChangeImpact, error) {
	// 获取当前费率配置
	config, err := s.repo.GetRateConfigByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("获取费率配置失败: %w", err)
	}

	if config.ChannelID != channelID {
		return nil, fmt.Errorf("费率配置不属于该通道")
	}

	impact := &models.ConfigChangeImpact{
		AffectedTemplates:   make([]models.AffectedTemplate, 0),
		AffectedSettlements: make([]models.AffectedSettlement, 0),
	}

	// 检查政策模版
	templates, err := s.repo.GetPolicyTemplatesByChannel(ctx, channelID)
	if err == nil {
		for _, t := range templates {
			if rateVal, ok := t.RateConfigs[config.RateCode]; ok {
				rateFloat := models.ParseRateToFloat(rateVal.Rate)
				newMinFloat := models.ParseRateToFloat(req.NewMinRate)
				newMaxFloat := models.ParseRateToFloat(req.NewMaxRate)

				var issue string
				if rateFloat < newMinFloat {
					issue = fmt.Sprintf("费率%s低于新下限%s", rateVal.Rate, req.NewMinRate)
				} else if rateFloat > newMaxFloat {
					issue = fmt.Sprintf("费率%s超过新上限%s", rateVal.Rate, req.NewMaxRate)
				}

				if issue != "" {
					impact.AffectedTemplates = append(impact.AffectedTemplates, models.AffectedTemplate{
						TemplateID:   t.ID,
						TemplateName: t.TemplateName,
						Issue:        issue,
					})
				}
			}
		}
	}

	// 检查结算价
	settlements, err := s.repo.GetSettlementPricesByChannel(ctx, channelID)
	if err == nil {
		agentMap := make(map[int64]bool)
		for _, sp := range settlements {
			if rateVal, ok := sp.RateConfigs[config.RateCode]; ok {
				rateFloat := models.ParseRateToFloat(rateVal.Rate)
				newMinFloat := models.ParseRateToFloat(req.NewMinRate)
				newMaxFloat := models.ParseRateToFloat(req.NewMaxRate)

				var issue string
				if rateFloat < newMinFloat {
					issue = fmt.Sprintf("费率%s低于新下限%s", rateVal.Rate, req.NewMinRate)
				} else if rateFloat > newMaxFloat {
					issue = fmt.Sprintf("费率%s超过新上限%s", rateVal.Rate, req.NewMaxRate)
				}

				// 检查高调费率
				if req.NewMaxHighRate != nil {
					if hrConfig, ok := sp.HighRateConfigs[config.RateCode]; ok {
						if err := models.ValidateHighRate(hrConfig.Rate, req.NewMaxHighRate, config.RateName); err != nil {
							if issue != "" {
								issue += "; "
							}
							issue += err.Error()
						}
					}
				}

				// 检查P+0加价
				if req.NewMaxD0Extra != nil {
					if d0Config, ok := sp.D0ExtraConfigs[config.RateCode]; ok {
						if err := models.ValidateD0Extra(d0Config.ExtraFee, req.NewMaxD0Extra, config.RateName); err != nil {
							if issue != "" {
								issue += "; "
							}
							issue += err.Error()
						}
					}
				}

				if issue != "" {
					impact.AffectedSettlements = append(impact.AffectedSettlements, models.AffectedSettlement{
						SettlementID: sp.ID,
						AgentID:      sp.AgentID,
						AgentName:    "", // 需要额外查询代理商名称
						Issue:        issue,
					})
					agentMap[sp.AgentID] = true
				}
			}
		}
		impact.TotalAffectedAgents = len(agentMap)
	}

	return impact, nil
}
