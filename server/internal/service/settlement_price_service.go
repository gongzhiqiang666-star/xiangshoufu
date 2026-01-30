package service

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// SettlementPriceService 结算价服务
type SettlementPriceService struct {
	repo              repository.SettlementPriceRepository
	changeLogRepo     repository.PriceChangeLogRepository
	agentRepo         repository.AgentRepository
	channelConfigRepo *repository.GormChannelConfigRepository
	db                *gorm.DB
}

// NewSettlementPriceService 创建结算价服务
func NewSettlementPriceService(
	repo repository.SettlementPriceRepository,
	changeLogRepo repository.PriceChangeLogRepository,
	agentRepo repository.AgentRepository,
	channelConfigRepo *repository.GormChannelConfigRepository,
	db *gorm.DB,
) *SettlementPriceService {
	return &SettlementPriceService{
		repo:              repo,
		changeLogRepo:     changeLogRepo,
		agentRepo:         agentRepo,
		channelConfigRepo: channelConfigRepo,
		db:                db,
	}
}

// CreateFromTemplate 从模板创建结算价
func (s *SettlementPriceService) CreateFromTemplate(
	agentID, channelID int64,
	templateID *int64,
	brandCode string,
	template *models.PolicyTemplateComplete,
	operatorID int64,
	operatorName string,
	source string,
) (*models.SettlementPrice, error) {
	// 检查是否已存在
	existing, err := s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("结算价已存在")
	}

	// 创建结算价
	price := &models.SettlementPrice{
		AgentID:    agentID,
		ChannelID:  channelID,
		TemplateID: templateID,
		BrandCode:  brandCode,
		Version:    1,
		Status:     1,
		EffectiveAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CreatedBy:  &operatorID,
	}

	// 如果有模板，从模板复制配置
	if template != nil {
		price.RateConfigs = template.RateConfigs
		price.CreditRate = &template.CreditRate
		price.DebitRate = &template.DebitRate
		price.DebitCap = &template.DebitCap
		price.UnionpayRate = &template.UnionpayRate
		price.WechatRate = &template.WechatRate
		price.AlipayRate = &template.AlipayRate

		// 复制押金返现配置
		if len(template.DepositCashbackPolicies) > 0 {
			depositCashbacks := make(models.DepositCashbacks, 0, len(template.DepositCashbackPolicies))
			for _, dp := range template.DepositCashbackPolicies {
				depositCashbacks = append(depositCashbacks, models.DepositCashbackItem{
					DepositAmount:  dp.DepositAmount,
					CashbackAmount: dp.CashbackAmount,
				})
			}
			price.DepositCashbacks = depositCashbacks
		}

		// 复制流量卡返现配置
		if len(template.SimCashbackPolicies) > 0 {
			// 取第一个配置
			simPolicy := template.SimCashbackPolicies[0]
			price.SimFirstCashback = simPolicy.FirstTimeCashback
			price.SimSecondCashback = simPolicy.SecondTimeCashback
			price.SimThirdPlusCashback = simPolicy.ThirdPlusCashback
		}
	}

	// 校验结算价配置是否符合约束
	simCashbacks := models.SimCashbacks{
		{TierOrder: 1, CashbackAmount: price.SimFirstCashback},
		{TierOrder: 2, CashbackAmount: price.SimSecondCashback},
		{TierOrder: 3, CashbackAmount: price.SimThirdPlusCashback},
	}
	if err := s.ValidateSettlementPriceConstraints(agentID, channelID, brandCode, price.RateConfigs, price.DepositCashbacks, simCashbacks); err != nil {
		return nil, fmt.Errorf("结算价约束校验失败: %w", err)
	}

	// 保存结算价
	err = s.repo.Create(price)
	if err != nil {
		return nil, fmt.Errorf("创建结算价失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLog(price, nil, models.ChangeTypeInit, operatorID, operatorName, source, "", "初始化结算价")

	return price, nil
}

// UpdateRate 更新费率
func (s *SettlementPriceService) UpdateRate(
	id int64,
	req *models.UpdateRateRequest,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.SettlementPrice, error) {
	price, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("结算价不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(price)

	// 更新费率
	if req.RateConfigs != nil {
		price.RateConfigs = req.RateConfigs
	}
	if req.CreditRate != nil {
		price.CreditRate = req.CreditRate
	}
	if req.DebitRate != nil {
		price.DebitRate = req.DebitRate
	}
	if req.DebitCap != nil {
		price.DebitCap = req.DebitCap
	}
	if req.UnionpayRate != nil {
		price.UnionpayRate = req.UnionpayRate
	}
	if req.WechatRate != nil {
		price.WechatRate = req.WechatRate
	}
	if req.AlipayRate != nil {
		price.AlipayRate = req.AlipayRate
	}

	// 校验费率配置是否符合约束
	simCashbacks := models.SimCashbacks{
		{TierOrder: 1, CashbackAmount: price.SimFirstCashback},
		{TierOrder: 2, CashbackAmount: price.SimSecondCashback},
		{TierOrder: 3, CashbackAmount: price.SimThirdPlusCashback},
	}
	if err := s.ValidateSettlementPriceConstraints(price.AgentID, price.ChannelID, price.BrandCode, price.RateConfigs, price.DepositCashbacks, simCashbacks); err != nil {
		return nil, fmt.Errorf("费率约束校验失败: %w", err)
	}

	// 版本号+1
	price.Version++
	price.UpdatedBy = &operatorID

	err = s.repo.Update(price)
	if err != nil {
		return nil, fmt.Errorf("更新费率失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(price, snapshotBefore, models.ChangeTypeRate, operatorID, operatorName, source, ipAddress, "费率", "费率调整")

	return price, nil
}

// UpdateDepositCashback 更新押金返现
func (s *SettlementPriceService) UpdateDepositCashback(
	id int64,
	req *models.UpdateDepositCashbackRequest,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.SettlementPrice, error) {
	price, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("结算价不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(price)

	// 校验押金返现配置是否符合约束
	simCashbacks := models.SimCashbacks{
		{TierOrder: 1, CashbackAmount: price.SimFirstCashback},
		{TierOrder: 2, CashbackAmount: price.SimSecondCashback},
		{TierOrder: 3, CashbackAmount: price.SimThirdPlusCashback},
	}
	if err := s.ValidateSettlementPriceConstraints(price.AgentID, price.ChannelID, price.BrandCode, price.RateConfigs, req.DepositCashbacks, simCashbacks); err != nil {
		return nil, fmt.Errorf("押金返现约束校验失败: %w", err)
	}

	// 更新押金返现
	price.DepositCashbacks = req.DepositCashbacks
	price.Version++
	price.UpdatedBy = &operatorID

	err = s.repo.Update(price)
	if err != nil {
		return nil, fmt.Errorf("更新押金返现失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(price, snapshotBefore, models.ChangeTypeDeposit, operatorID, operatorName, source, ipAddress, "押金返现", "押金返现调整")

	return price, nil
}

// UpdateSimCashback 更新流量费返现
func (s *SettlementPriceService) UpdateSimCashback(
	id int64,
	req *models.UpdateSimCashbackRequest,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.SettlementPrice, error) {
	price, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("结算价不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(price)

	// 校验流量费返现配置是否符合约束
	simCashbacks := models.SimCashbacks{
		{TierOrder: 1, CashbackAmount: req.SimFirstCashback},
		{TierOrder: 2, CashbackAmount: req.SimSecondCashback},
		{TierOrder: 3, CashbackAmount: req.SimThirdPlusCashback},
	}
	if err := s.ValidateSettlementPriceConstraints(price.AgentID, price.ChannelID, price.BrandCode, price.RateConfigs, price.DepositCashbacks, simCashbacks); err != nil {
		return nil, fmt.Errorf("流量费返现约束校验失败: %w", err)
	}

	// 更新流量费返现
	price.SimFirstCashback = req.SimFirstCashback
	price.SimSecondCashback = req.SimSecondCashback
	price.SimThirdPlusCashback = req.SimThirdPlusCashback
	price.Version++
	price.UpdatedBy = &operatorID

	err = s.repo.Update(price)
	if err != nil {
		return nil, fmt.Errorf("更新流量费返现失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(price, snapshotBefore, models.ChangeTypeSim, operatorID, operatorName, source, ipAddress, "流量费返现", "流量费返现调整")

	return price, nil
}

// GetByID 根据ID获取结算价
func (s *SettlementPriceService) GetByID(id int64) (*models.SettlementPrice, error) {
	return s.repo.GetByID(id)
}

// GetByAgentAndChannel 根据代理商ID和通道ID获取结算价
func (s *SettlementPriceService) GetByAgentAndChannel(agentID, channelID int64, brandCode string) (*models.SettlementPrice, error) {
	return s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
}

// List 获取结算价列表
func (s *SettlementPriceService) List(req *models.SettlementPriceListRequest) (*models.SettlementPriceListResponse, error) {
	prices, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]models.SettlementPriceItem, 0, len(prices))
	for _, p := range prices {
		items = append(items, models.SettlementPriceItem{
			ID:                   p.ID,
			AgentID:              p.AgentID,
			ChannelID:            p.ChannelID,
			BrandCode:            p.BrandCode,
			RateConfigs:          p.RateConfigs,
			DepositCashbacks:     p.DepositCashbacks,
			SimFirstCashback:     p.SimFirstCashback,
			SimSecondCashback:    p.SimSecondCashback,
			SimThirdPlusCashback: p.SimThirdPlusCashback,
			Version:              p.Version,
			Status:               p.Status,
			EffectiveAt:          p.EffectiveAt,
			UpdatedAt:            p.UpdatedAt,
		})
	}

	return &models.SettlementPriceListResponse{
		List:  items,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// GetAgentRate 获取代理商费率（供分润计算使用）
func (s *SettlementPriceService) GetAgentRate(agentID, channelID int64, brandCode string, rateType string) (string, error) {
	price, err := s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
	if err != nil {
		return "", fmt.Errorf("获取结算价失败: %w", err)
	}

	// 优先使用动态费率配置
	if price.RateConfigs != nil {
		if rateConfig, ok := price.RateConfigs[rateType]; ok {
			return rateConfig.Rate, nil
		}
	}

	// 降级使用旧字段
	switch rateType {
	case "credit":
		if price.CreditRate != nil {
			return *price.CreditRate, nil
		}
	case "debit":
		if price.DebitRate != nil {
			return *price.DebitRate, nil
		}
	case "unionpay":
		if price.UnionpayRate != nil {
			return *price.UnionpayRate, nil
		}
	case "wechat":
		if price.WechatRate != nil {
			return *price.WechatRate, nil
		}
	case "alipay":
		if price.AlipayRate != nil {
			return *price.AlipayRate, nil
		}
	}

	return "0", nil
}

// GetAgentDepositCashback 获取代理商押金返现配置
func (s *SettlementPriceService) GetAgentDepositCashback(agentID, channelID int64, brandCode string, depositAmount int64) (int64, error) {
	price, err := s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
	if err != nil {
		return 0, fmt.Errorf("获取结算价失败: %w", err)
	}

	for _, dc := range price.DepositCashbacks {
		if dc.DepositAmount == depositAmount {
			return dc.CashbackAmount, nil
		}
	}

	return 0, nil
}

// GetAgentSimCashback 获取代理商流量费返现配置
func (s *SettlementPriceService) GetAgentSimCashback(agentID, channelID int64, brandCode string) (firstCashback, secondCashback, thirdPlusCashback int64, err error) {
	price, err := s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("获取结算价失败: %w", err)
	}

	return price.SimFirstCashback, price.SimSecondCashback, price.SimThirdPlusCashback, nil
}

// GetAgentHighRate 获取代理商高调费率配置
func (s *SettlementPriceService) GetAgentHighRate(agentID, channelID int64, brandCode string, rateType string) (string, error) {
	price, err := s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
	if err != nil {
		return "0", fmt.Errorf("获取结算价失败: %w", err)
	}

	if price.HighRateConfigs != nil {
		if config, ok := price.HighRateConfigs[rateType]; ok {
			return config.Rate, nil
		}
	}

	return "0", nil
}

// GetAgentD0Extra 获取代理商P+0加价配置
func (s *SettlementPriceService) GetAgentD0Extra(agentID, channelID int64, brandCode string, rateType string) (int64, error) {
	price, err := s.repo.GetByAgentAndChannel(agentID, channelID, brandCode)
	if err != nil {
		return 0, fmt.Errorf("获取结算价失败: %w", err)
	}

	if price.D0ExtraConfigs != nil {
		if config, ok := price.D0ExtraConfigs[rateType]; ok {
			return config.ExtraFee, nil
		}
	}

	return 0, nil
}

// UpdateHighRateRequest 更新高调费率请求
type UpdateHighRateRequest struct {
	HighRateConfigs models.HighRateConfigs `json:"high_rate_configs" binding:"required"`
}

// UpdateD0ExtraRequest 更新P+0加价请求
type UpdateD0ExtraRequest struct {
	D0ExtraConfigs models.D0ExtraConfigs `json:"d0_extra_configs" binding:"required"`
}

// UpdateHighRate 更新高调费率配置
func (s *SettlementPriceService) UpdateHighRate(
	id int64,
	req *UpdateHighRateRequest,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.SettlementPrice, error) {
	price, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("结算价不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(price)

	// 更新高调费率配置
	price.HighRateConfigs = req.HighRateConfigs
	price.Version++
	price.UpdatedBy = &operatorID

	err = s.repo.Update(price)
	if err != nil {
		return nil, fmt.Errorf("更新高调费率失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(price, snapshotBefore, models.ChangeTypeRate, operatorID, operatorName, source, ipAddress, "高调费率", "高调费率调整")

	return price, nil
}

// UpdateD0Extra 更新P+0加价配置
func (s *SettlementPriceService) UpdateD0Extra(
	id int64,
	req *UpdateD0ExtraRequest,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
) (*models.SettlementPrice, error) {
	price, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("结算价不存在: %w", err)
	}

	// 保存变更前的快照
	snapshotBefore := s.createSnapshot(price)

	// 更新P+0加价配置
	price.D0ExtraConfigs = req.D0ExtraConfigs
	price.Version++
	price.UpdatedBy = &operatorID

	err = s.repo.Update(price)
	if err != nil {
		return nil, fmt.Errorf("更新P+0加价失败: %w", err)
	}

	// 记录调价日志
	s.createChangeLogWithSnapshot(price, snapshotBefore, models.ChangeTypeRate, operatorID, operatorName, source, ipAddress, "P+0加价", "P+0加价调整")

	return price, nil
}

// createSnapshot 创建快照
func (s *SettlementPriceService) createSnapshot(price *models.SettlementPrice) models.JSONMap {
	snapshot := models.JSONMap{
		"id":                     price.ID,
		"agent_id":               price.AgentID,
		"channel_id":             price.ChannelID,
		"brand_code":             price.BrandCode,
		"rate_configs":           price.RateConfigs,
		"deposit_cashbacks":      price.DepositCashbacks,
		"sim_first_cashback":     price.SimFirstCashback,
		"sim_second_cashback":    price.SimSecondCashback,
		"sim_third_plus_cashback": price.SimThirdPlusCashback,
		"version":                price.Version,
	}
	return snapshot
}

// createChangeLog 创建调价日志
func (s *SettlementPriceService) createChangeLog(
	price *models.SettlementPrice,
	snapshotBefore models.JSONMap,
	changeType models.ChangeType,
	operatorID int64,
	operatorName string,
	source string,
	fieldName string,
	summary string,
) {
	snapshotAfter := s.createSnapshot(price)

	log := &models.PriceChangeLog{
		AgentID:           price.AgentID,
		ChannelID:         &price.ChannelID,
		SettlementPriceID: &price.ID,
		ChangeType:        changeType,
		ConfigType:        models.ConfigTypeSettlement,
		FieldName:         fieldName,
		ChangeSummary:     summary,
		SnapshotBefore:    snapshotBefore,
		SnapshotAfter:     snapshotAfter,
		OperatorType:      models.OperatorTypeAdmin,
		OperatorID:        operatorID,
		OperatorName:      operatorName,
		Source:            source,
		CreatedAt:         time.Now(),
	}

	s.changeLogRepo.Create(log)
}

// createChangeLogWithSnapshot 创建带快照的调价日志
func (s *SettlementPriceService) createChangeLogWithSnapshot(
	price *models.SettlementPrice,
	snapshotBefore models.JSONMap,
	changeType models.ChangeType,
	operatorID int64,
	operatorName string,
	source string,
	ipAddress string,
	fieldName string,
	summary string,
) {
	snapshotAfter := s.createSnapshot(price)

	// 计算变更内容
	oldValue, _ := json.Marshal(snapshotBefore)
	newValue, _ := json.Marshal(snapshotAfter)

	log := &models.PriceChangeLog{
		AgentID:           price.AgentID,
		ChannelID:         &price.ChannelID,
		SettlementPriceID: &price.ID,
		ChangeType:        changeType,
		ConfigType:        models.ConfigTypeSettlement,
		FieldName:         fieldName,
		OldValue:          string(oldValue),
		NewValue:          string(newValue),
		ChangeSummary:     summary,
		SnapshotBefore:    snapshotBefore,
		SnapshotAfter:     snapshotAfter,
		OperatorType:      models.OperatorTypeAdmin,
		OperatorID:        operatorID,
		OperatorName:      operatorName,
		Source:            source,
		IPAddress:         ipAddress,
		CreatedAt:         time.Now(),
	}

	s.changeLogRepo.Create(log)
}

// ============================================================
// 上下级约束校验
// ============================================================

// ValidateSettlementPriceConstraints 校验结算价是否符合上下级约束
func (s *SettlementPriceService) ValidateSettlementPriceConstraints(
	agentID, channelID int64,
	brandCode string,
	rateConfigs models.RateConfigs,
	depositCashbacks models.DepositCashbacks,
	simCashbacks models.SimCashbacks,
) error {
	if s.agentRepo == nil {
		return nil // 没有注入代理商仓库时跳过校验
	}

	// 1. 获取代理商信息
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil || agent == nil {
		return nil // 代理商不存在时跳过校验
	}

	// 2. 如果有上级，校验不能超过上级配置
	if agent.ParentID > 0 {
		upperPrice, err := s.repo.GetByAgentAndChannel(agent.ParentID, channelID, brandCode)
		if err == nil && upperPrice != nil {
			// 校验费率：下级费率 >= 上级费率
			for rateCode, rateConfig := range rateConfigs {
				if upperRateConfig, ok := upperPrice.RateConfigs[rateCode]; ok {
					if err := models.ValidateSettlementRate(rateConfig.Rate, upperRateConfig.Rate, "0", "100"); err != nil {
						return fmt.Errorf("%s: %w", getRateCodeName(rateCode), err)
					}
				}
			}

			// 校验押金返现：下级返现 <= 上级返现
			for _, dc := range depositCashbacks {
				for _, upperDc := range upperPrice.DepositCashbacks {
					if dc.DepositAmount == upperDc.DepositAmount {
						if dc.CashbackAmount > upperDc.CashbackAmount {
							return fmt.Errorf("押金%d元的返现%d元不能超过上级配置%d元",
								dc.DepositAmount/100, dc.CashbackAmount/100, upperDc.CashbackAmount/100)
						}
					}
				}
			}

			// 校验流量费返现：下级每档返现 <= 上级对应档返现
			for _, tier := range simCashbacks {
				upperAmount := upperPrice.SimCashbacks.GetCashbackAmount(tier.TierOrder)
				if upperAmount > 0 && tier.CashbackAmount > upperAmount {
					return fmt.Errorf("流量费第%d次返现%d元不能超过上级配置%d元",
						tier.TierOrder, tier.CashbackAmount/100, upperAmount/100)
				}
			}
		}
	}

	// 3. 校验不超过通道上限
	if s.channelConfigRepo != nil {
		channelConfig, err := s.channelConfigRepo.GetFullConfig(nil, channelID)
		if err == nil && channelConfig != nil {
			// 校验费率不超过通道上限
			for rateCode, rateConfig := range rateConfigs {
				for _, channelRate := range channelConfig.RateConfigs {
					if channelRate.RateCode == rateCode {
						if err := models.ValidateRateRange(rateConfig.Rate, channelRate.MinRate, channelRate.MaxRate); err != nil {
							return fmt.Errorf("%s费率校验失败: %w", channelRate.RateName, err)
						}
					}
				}
			}

			// 校验押金返现不超过通道上限
			for _, dc := range depositCashbacks {
				for _, tier := range channelConfig.DepositTiers {
					if tier.DepositAmount == dc.DepositAmount && tier.MaxCashbackAmount > 0 {
						if dc.CashbackAmount > tier.MaxCashbackAmount {
							return fmt.Errorf("押金%d元的返现%d元超过通道上限%d元",
								dc.DepositAmount/100, dc.CashbackAmount/100, tier.MaxCashbackAmount/100)
						}
					}
				}
			}

			// 校验流量费返现不超过通道上限
			for _, sc := range simCashbacks {
				for _, tier := range channelConfig.SimCashbackTiers {
					if tier.TierOrder == sc.TierOrder {
						if sc.CashbackAmount > tier.MaxCashbackAmount {
							return fmt.Errorf("流量费第%d次返现%d元超过通道上限%d元",
								sc.TierOrder, sc.CashbackAmount/100, tier.MaxCashbackAmount/100)
						}
					}
				}
			}
		}
	}

	return nil
}

// getRateCodeName 获取费率编码的中文名称
func getRateCodeName(code string) string {
	names := map[string]string{
		"CREDIT":    "贷记卡",
		"DEBIT":     "借记卡",
		"WECHAT":    "微信",
		"ALIPAY":    "支付宝",
		"UNIONPAY":  "云闪付",
		"CLOUD_PAY": "云闪付",
	}
	if name, ok := names[code]; ok {
		return name
	}
	return code
}
