package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// PolicyService 政策管理服务
// 处理政策模板CRUD、代理商政策分配和调整
type PolicyService struct {
	templateRepo           *repository.GormPolicyTemplateRepository
	depositPolicyRepo      *repository.GormDepositCashbackPolicyRepository
	simPolicyRepo          *repository.GormSimCashbackPolicyRepository
	rewardPolicyRepo       *repository.GormActivationRewardPolicyRepository
	rateStagePolicyRepo    *repository.GormRateStagePolicyRepository
	agentPolicyRepo        repository.AgentPolicyRepository
	agentDepositRepo       *repository.GormAgentDepositCashbackPolicyRepository
	agentSimRepo           *repository.GormAgentSimCashbackPolicyRepository
	agentRewardRepo        *repository.GormAgentActivationRewardPolicyRepository
	agentRepo              repository.AgentRepository
	channelConfigRepo      *repository.GormChannelConfigRepository
	db                     *gorm.DB
}

// NewPolicyService 创建政策服务
func NewPolicyService(
	templateRepo *repository.GormPolicyTemplateRepository,
	depositPolicyRepo *repository.GormDepositCashbackPolicyRepository,
	simPolicyRepo *repository.GormSimCashbackPolicyRepository,
	rewardPolicyRepo *repository.GormActivationRewardPolicyRepository,
	rateStagePolicyRepo *repository.GormRateStagePolicyRepository,
	agentPolicyRepo repository.AgentPolicyRepository,
	agentDepositRepo *repository.GormAgentDepositCashbackPolicyRepository,
	agentSimRepo *repository.GormAgentSimCashbackPolicyRepository,
	agentRewardRepo *repository.GormAgentActivationRewardPolicyRepository,
	agentRepo repository.AgentRepository,
	channelConfigRepo *repository.GormChannelConfigRepository,
	db *gorm.DB,
) *PolicyService {
	return &PolicyService{
		templateRepo:        templateRepo,
		depositPolicyRepo:   depositPolicyRepo,
		simPolicyRepo:       simPolicyRepo,
		rewardPolicyRepo:    rewardPolicyRepo,
		rateStagePolicyRepo: rateStagePolicyRepo,
		agentPolicyRepo:     agentPolicyRepo,
		agentDepositRepo:    agentDepositRepo,
		agentSimRepo:        agentSimRepo,
		agentRewardRepo:     agentRewardRepo,
		agentRepo:           agentRepo,
		channelConfigRepo:   channelConfigRepo,
		db:                  db,
	}
}

// ============================================================
// 政策模板请求/响应结构
// ============================================================

// CreatePolicyTemplateRequest 创建政策模板请求
type CreatePolicyTemplateRequest struct {
	TemplateName string `json:"template_name" binding:"required"`
	ChannelID    int64  `json:"channel_id" binding:"required"`
	IsDefault    bool   `json:"is_default"`

	// 动态费率配置（新版）
	RateConfigs models.RateConfigs `json:"rate_configs"`

	// 成本（费率）- 旧字段，保留兼容
	CreditRate   string `json:"credit_rate"`
	DebitRate    string `json:"debit_rate"`
	DebitCap     string `json:"debit_cap"`
	UnionpayRate string `json:"unionpay_rate"`
	WechatRate   string `json:"wechat_rate"`
	AlipayRate   string `json:"alipay_rate"`

	// 押金返现配置
	DepositCashbacks []DepositCashbackInput `json:"deposit_cashbacks"`

	// 流量卡返现配置
	SimCashback *SimCashbackInput `json:"sim_cashback"`

	// 激活奖励配置
	ActivationRewards []ActivationRewardInput `json:"activation_rewards"`

	// 费率阶梯配置
	RateStages []RateStageInput `json:"rate_stages"`
}

// DepositCashbackInput 押金返现输入
type DepositCashbackInput struct {
	DepositAmount  int64 `json:"deposit_amount"`  // 押金金额（分）
	CashbackAmount int64 `json:"cashback_amount"` // 返现金额（分）
}

// SimCashbackInput 流量卡返现输入
type SimCashbackInput struct {
	FirstTimeCashback  int64 `json:"first_time_cashback"`
	SecondTimeCashback int64 `json:"second_time_cashback"`
	ThirdPlusCashback  int64 `json:"third_plus_cashback"`
	SimFeeAmount       int64 `json:"sim_fee_amount"`
}

// ActivationRewardInput 激活奖励输入
type ActivationRewardInput struct {
	RewardName      string `json:"reward_name"`
	MinRegisterDays int    `json:"min_register_days"`
	MaxRegisterDays int    `json:"max_register_days"`
	TargetAmount    int64  `json:"target_amount"`
	RewardAmount    int64  `json:"reward_amount"`
	Priority        int    `json:"priority"`
}

// RateStageInput 费率阶梯输入
type RateStageInput struct {
	StageName string `json:"stage_name"`
	ApplyTo   int16  `json:"apply_to"` // 1-商户 2-代理商
	MinDays   int    `json:"min_days"`
	MaxDays   int    `json:"max_days"`

	// 动态费率阶梯调整值（新版）
	RateDeltas models.RateDeltas `json:"rate_deltas"`

	// 旧字段，保留兼容
	CreditRateDelta   string `json:"credit_rate_delta"`
	DebitRateDelta    string `json:"debit_rate_delta"`
	UnionpayRateDelta string `json:"unionpay_rate_delta"`
	WechatRateDelta   string `json:"wechat_rate_delta"`
	AlipayRateDelta   string `json:"alipay_rate_delta"`
	Priority          int    `json:"priority"`
}

// PolicyTemplateResponse 政策模板响应
type PolicyTemplateResponse struct {
	ID           int64  `json:"id"`
	TemplateName string `json:"template_name"`
	ChannelID    int64  `json:"channel_id"`
	IsDefault    bool   `json:"is_default"`
	Status       int16  `json:"status"`
	CreatedAt    string `json:"created_at"`

	// 动态费率配置（新版）
	RateConfigs models.RateConfigs `json:"rate_configs"`

	// 成本（费率）- 旧字段，保留兼容
	CreditRate   string `json:"credit_rate"`
	DebitRate    string `json:"debit_rate"`
	DebitCap     string `json:"debit_cap"`
	UnionpayRate string `json:"unionpay_rate"`
	WechatRate   string `json:"wechat_rate"`
	AlipayRate   string `json:"alipay_rate"`

	// 押金返现配置
	DepositCashbacks []DepositCashbackInput `json:"deposit_cashbacks"`

	// 流量卡返现配置
	SimCashback *SimCashbackInput `json:"sim_cashback"`

	// 激活奖励配置
	ActivationRewards []ActivationRewardInput `json:"activation_rewards"`

	// 费率阶梯配置
	RateStages []RateStageInput `json:"rate_stages"`
}

// ============================================================
// 政策模板CRUD
// ============================================================

// CreatePolicyTemplate 创建政策模板
func (s *PolicyService) CreatePolicyTemplate(req *CreatePolicyTemplateRequest) (*PolicyTemplateResponse, error) {
	if req.TemplateName == "" {
		return nil, errors.New("模板名称不能为空")
	}

	// 校验模板配置是否符合通道约束
	ctx := context.Background()
	if err := s.validateTemplateAgainstChannel(ctx, req.ChannelID, req); err != nil {
		return nil, fmt.Errorf("通道约束校验失败: %w", err)
	}

	// 创建主模板
	template := &models.PolicyTemplateComplete{
		TemplateName: req.TemplateName,
		ChannelID:    req.ChannelID,
		IsDefault:    req.IsDefault,
		RateConfigs:  req.RateConfigs,
		CreditRate:   req.CreditRate,
		DebitRate:    req.DebitRate,
		DebitCap:     req.DebitCap,
		UnionpayRate: req.UnionpayRate,
		WechatRate:   req.WechatRate,
		AlipayRate:   req.AlipayRate,
		Status:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.templateRepo.Create(template); err != nil {
		return nil, fmt.Errorf("创建政策模板失败: %w", err)
	}

	// 创建押金返现政策
	if err := s.createDepositCashbackPolicies(template.ID, req.ChannelID, req.DepositCashbacks); err != nil {
		log.Printf("[PolicyService] Create deposit cashback policies failed: %v", err)
	}

	// 创建流量卡返现政策
	if req.SimCashback != nil {
		if err := s.createSimCashbackPolicy(template.ID, req.ChannelID, req.SimCashback); err != nil {
			log.Printf("[PolicyService] Create sim cashback policy failed: %v", err)
		}
	}

	// 创建激活奖励政策
	if err := s.createActivationRewardPolicies(template.ID, req.ChannelID, req.ActivationRewards); err != nil {
		log.Printf("[PolicyService] Create activation reward policies failed: %v", err)
	}

	// 创建费率阶梯政策
	if err := s.createRateStagePolicies(template.ID, req.ChannelID, req.RateStages); err != nil {
		log.Printf("[PolicyService] Create rate stage policies failed: %v", err)
	}

	log.Printf("[PolicyService] Created policy template: id=%d, name=%s", template.ID, template.TemplateName)

	return s.GetPolicyTemplateDetail(template.ID)
}

// UpdatePolicyTemplate 更新政策模板
func (s *PolicyService) UpdatePolicyTemplate(id int64, req *CreatePolicyTemplateRequest) (*PolicyTemplateResponse, error) {
	template, err := s.templateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("政策模板不存在: %d", id)
	}

	// 校验模板配置是否符合通道约束
	ctx := context.Background()
	if err := s.validateTemplateAgainstChannel(ctx, req.ChannelID, req); err != nil {
		return nil, fmt.Errorf("通道约束校验失败: %w", err)
	}

	// 使用事务保护更新操作
	if s.db != nil {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			// 更新主模板
			template.TemplateName = req.TemplateName
			template.IsDefault = req.IsDefault
			template.RateConfigs = req.RateConfigs
			template.CreditRate = req.CreditRate
			template.DebitRate = req.DebitRate
			template.DebitCap = req.DebitCap
			template.UnionpayRate = req.UnionpayRate
			template.WechatRate = req.WechatRate
			template.AlipayRate = req.AlipayRate
			template.UpdatedAt = time.Now()

			if err := tx.Save(template).Error; err != nil {
				return fmt.Errorf("更新政策模板失败: %w", err)
			}

			// 删除旧的关联政策
			if err := tx.Where("template_id = ?", id).Delete(&models.DepositCashbackPolicy{}).Error; err != nil {
				return fmt.Errorf("删除押金返现政策失败: %w", err)
			}
			if err := tx.Where("template_id = ?", id).Delete(&models.SimCashbackPolicy{}).Error; err != nil {
				return fmt.Errorf("删除流量费返现政策失败: %w", err)
			}
			if err := tx.Where("template_id = ?", id).Delete(&models.ActivationRewardPolicy{}).Error; err != nil {
				return fmt.Errorf("删除激活奖励政策失败: %w", err)
			}
			if err := tx.Where("template_id = ?", id).Delete(&models.RateStagePolicy{}).Error; err != nil {
				return fmt.Errorf("删除费率阶梯政策失败: %w", err)
			}

			// 重新创建关联政策
			for _, input := range req.DepositCashbacks {
				policy := &models.DepositCashbackPolicy{
					TemplateID:     id,
					ChannelID:      req.ChannelID,
					DepositAmount:  input.DepositAmount,
					CashbackAmount: input.CashbackAmount,
					Status:         1,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				if err := tx.Create(policy).Error; err != nil {
					return fmt.Errorf("创建押金返现政策失败: %w", err)
				}
			}

			if req.SimCashback != nil {
				policy := &models.SimCashbackPolicy{
					TemplateID:         id,
					ChannelID:          req.ChannelID,
					FirstTimeCashback:  req.SimCashback.FirstTimeCashback,
					SecondTimeCashback: req.SimCashback.SecondTimeCashback,
					ThirdPlusCashback:  req.SimCashback.ThirdPlusCashback,
					SimFeeAmount:       req.SimCashback.SimFeeAmount,
					Status:             1,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				}
				// 同时设置新版N档格式
				policy.SetFromOldFields()
				if err := tx.Create(policy).Error; err != nil {
					return fmt.Errorf("创建流量费返现政策失败: %w", err)
				}
			}

			for _, input := range req.ActivationRewards {
				policy := &models.ActivationRewardPolicy{
					TemplateID:      id,
					ChannelID:       req.ChannelID,
					RewardName:      input.RewardName,
					MinRegisterDays: input.MinRegisterDays,
					MaxRegisterDays: input.MaxRegisterDays,
					TargetAmount:    input.TargetAmount,
					RewardAmount:    input.RewardAmount,
					Priority:        input.Priority,
					Status:          1,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				if err := tx.Create(policy).Error; err != nil {
					return fmt.Errorf("创建激活奖励政策失败: %w", err)
				}
			}

			for _, input := range req.RateStages {
				policy := &models.RateStagePolicy{
					TemplateID:        id,
					ChannelID:         req.ChannelID,
					StageName:         input.StageName,
					ApplyTo:           input.ApplyTo,
					MinDays:           input.MinDays,
					MaxDays:           input.MaxDays,
					RateDeltas:        input.RateDeltas,
					CreditRateDelta:   input.CreditRateDelta,
					DebitRateDelta:    input.DebitRateDelta,
					UnionpayRateDelta: input.UnionpayRateDelta,
					WechatRateDelta:   input.WechatRateDelta,
					AlipayRateDelta:   input.AlipayRateDelta,
					Priority:          input.Priority,
					Status:            1,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}
				if err := tx.Create(policy).Error; err != nil {
					return fmt.Errorf("创建费率阶梯政策失败: %w", err)
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// 无事务支持时的降级处理（保持原有逻辑）
		template.TemplateName = req.TemplateName
		template.IsDefault = req.IsDefault
		template.RateConfigs = req.RateConfigs
		template.CreditRate = req.CreditRate
		template.DebitRate = req.DebitRate
		template.DebitCap = req.DebitCap
		template.UnionpayRate = req.UnionpayRate
		template.WechatRate = req.WechatRate
		template.AlipayRate = req.AlipayRate
		template.UpdatedAt = time.Now()

		if err := s.templateRepo.Update(template); err != nil {
			return nil, fmt.Errorf("更新政策模板失败: %w", err)
		}

		// 删除并重新创建关联政策
		s.depositPolicyRepo.DeleteByTemplateID(id)
		s.simPolicyRepo.DeleteByTemplateID(id)
		s.rewardPolicyRepo.DeleteByTemplateID(id)
		s.rateStagePolicyRepo.DeleteByTemplateID(id)

		s.createDepositCashbackPolicies(id, req.ChannelID, req.DepositCashbacks)
		if req.SimCashback != nil {
			s.createSimCashbackPolicy(id, req.ChannelID, req.SimCashback)
		}
		s.createActivationRewardPolicies(id, req.ChannelID, req.ActivationRewards)
		s.createRateStagePolicies(id, req.ChannelID, req.RateStages)
	}

	log.Printf("[PolicyService] Updated policy template: id=%d", id)

	return s.GetPolicyTemplateDetail(id)
}

// GetPolicyTemplateDetail 获取政策模板详情
func (s *PolicyService) GetPolicyTemplateDetail(id int64) (*PolicyTemplateResponse, error) {
	template, err := s.templateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("政策模板不存在: %d", id)
	}

	resp := &PolicyTemplateResponse{
		ID:           template.ID,
		TemplateName: template.TemplateName,
		ChannelID:    template.ChannelID,
		IsDefault:    template.IsDefault,
		Status:       template.Status,
		CreatedAt:    template.CreatedAt.Format("2006-01-02 15:04:05"),
		RateConfigs:  template.RateConfigs,
		CreditRate:   template.CreditRate,
		DebitRate:    template.DebitRate,
		DebitCap:     template.DebitCap,
		UnionpayRate: template.UnionpayRate,
		WechatRate:   template.WechatRate,
		AlipayRate:   template.AlipayRate,
	}

	// 加载押金返现配置
	depositPolicies, _ := s.depositPolicyRepo.FindByTemplateID(id)
	resp.DepositCashbacks = make([]DepositCashbackInput, len(depositPolicies))
	for i, p := range depositPolicies {
		resp.DepositCashbacks[i] = DepositCashbackInput{
			DepositAmount:  p.DepositAmount,
			CashbackAmount: p.CashbackAmount,
		}
	}

	// 加载流量卡返现配置
	simPolicies, _ := s.simPolicyRepo.FindByTemplateID(id)
	if len(simPolicies) > 0 {
		resp.SimCashback = &SimCashbackInput{
			FirstTimeCashback:  simPolicies[0].FirstTimeCashback,
			SecondTimeCashback: simPolicies[0].SecondTimeCashback,
			ThirdPlusCashback:  simPolicies[0].ThirdPlusCashback,
			SimFeeAmount:       simPolicies[0].SimFeeAmount,
		}
	}

	// 加载激活奖励配置
	rewardPolicies, _ := s.rewardPolicyRepo.FindByTemplateID(id)
	resp.ActivationRewards = make([]ActivationRewardInput, len(rewardPolicies))
	for i, p := range rewardPolicies {
		resp.ActivationRewards[i] = ActivationRewardInput{
			RewardName:      p.RewardName,
			MinRegisterDays: p.MinRegisterDays,
			MaxRegisterDays: p.MaxRegisterDays,
			TargetAmount:    p.TargetAmount,
			RewardAmount:    p.RewardAmount,
			Priority:        p.Priority,
		}
	}

	// 加载费率阶梯配置
	rateStagePolicies, _ := s.rateStagePolicyRepo.FindByTemplateID(id)
	resp.RateStages = make([]RateStageInput, len(rateStagePolicies))
	for i, p := range rateStagePolicies {
		resp.RateStages[i] = RateStageInput{
			StageName:         p.StageName,
			ApplyTo:           p.ApplyTo,
			MinDays:           p.MinDays,
			MaxDays:           p.MaxDays,
			RateDeltas:        p.RateDeltas,
			CreditRateDelta:   p.CreditRateDelta,
			DebitRateDelta:    p.DebitRateDelta,
			UnionpayRateDelta: p.UnionpayRateDelta,
			WechatRateDelta:   p.WechatRateDelta,
			AlipayRateDelta:   p.AlipayRateDelta,
			Priority:          p.Priority,
		}
	}

	return resp, nil
}

// GetPolicyTemplateList 获取政策模板列表
func (s *PolicyService) GetPolicyTemplateList(channelID int64, page, pageSize int) ([]*PolicyTemplateResponse, int64, error) {
	offset := (page - 1) * pageSize
	templates, total, err := s.templateRepo.FindByChannelID(channelID, []int16{1}, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*PolicyTemplateResponse, len(templates))
	for i, t := range templates {
		list[i] = &PolicyTemplateResponse{
			ID:           t.ID,
			TemplateName: t.TemplateName,
			ChannelID:    t.ChannelID,
			IsDefault:    t.IsDefault,
			Status:       t.Status,
			CreatedAt:    t.CreatedAt.Format("2006-01-02 15:04:05"),
			CreditRate:   t.CreditRate,
			DebitRate:    t.DebitRate,
			DebitCap:     t.DebitCap,
			UnionpayRate: t.UnionpayRate,
			WechatRate:   t.WechatRate,
			AlipayRate:   t.AlipayRate,
		}
	}

	return list, total, nil
}

// ============================================================
// 代理商政策分配
// ============================================================

// AssignAgentPolicyRequest 分配代理商政策请求
type AssignAgentPolicyRequest struct {
	AgentID    int64 `json:"agent_id" binding:"required"`
	ChannelID  int64 `json:"channel_id" binding:"required"`
	TemplateID int64 `json:"template_id" binding:"required"`

	// 可选覆盖值（在模板基础上调整）
	CreditRate   *string `json:"credit_rate"`
	DebitRate    *string `json:"debit_rate"`
	DebitCap     *string `json:"debit_cap"`
	UnionpayRate *string `json:"unionpay_rate"`
	WechatRate   *string `json:"wechat_rate"`
	AlipayRate   *string `json:"alipay_rate"`

	// 押金返现覆盖
	DepositCashbacks []DepositCashbackInput `json:"deposit_cashbacks"`

	// 流量卡返现覆盖
	SimCashback *SimCashbackInput `json:"sim_cashback"`

	// 激活奖励覆盖
	ActivationRewards []ActivationRewardInput `json:"activation_rewards"`
}

// AssignAgentPolicy 分配政策给代理商
func (s *PolicyService) AssignAgentPolicy(req *AssignAgentPolicyRequest, operatorID int64) error {
	// 验证代理商存在
	agent, err := s.agentRepo.FindByID(req.AgentID)
	if err != nil || agent == nil {
		return fmt.Errorf("代理商不存在: %d", req.AgentID)
	}

	// 验证模板存在
	template, err := s.templateRepo.FindByID(req.TemplateID)
	if err != nil || template == nil {
		return fmt.Errorf("政策模板不存在: %d", req.TemplateID)
	}

	// 保存押金返现政策
	if len(req.DepositCashbacks) > 0 {
		s.agentDepositRepo.DeleteByAgentAndChannel(req.AgentID, req.ChannelID)
		for _, dc := range req.DepositCashbacks {
			policy := &models.AgentDepositCashbackPolicy{
				AgentID:        req.AgentID,
				ChannelID:      req.ChannelID,
				DepositAmount:  dc.DepositAmount,
				CashbackAmount: dc.CashbackAmount,
				Status:         1,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			s.agentDepositRepo.Create(policy)
		}
	}

	// 保存流量卡返现政策
	if req.SimCashback != nil {
		policy := &models.AgentSimCashbackPolicy{
			AgentID:            req.AgentID,
			ChannelID:          req.ChannelID,
			FirstTimeCashback:  req.SimCashback.FirstTimeCashback,
			SecondTimeCashback: req.SimCashback.SecondTimeCashback,
			ThirdPlusCashback:  req.SimCashback.ThirdPlusCashback,
			Status:             1,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		s.agentSimRepo.Upsert(policy)
	}

	// 保存激活奖励政策
	if len(req.ActivationRewards) > 0 {
		s.agentRewardRepo.DeleteByAgentAndChannel(req.AgentID, req.ChannelID)
		for _, ar := range req.ActivationRewards {
			policy := &models.AgentActivationRewardPolicy{
				AgentID:         req.AgentID,
				ChannelID:       req.ChannelID,
				RewardName:      ar.RewardName,
				MinRegisterDays: ar.MinRegisterDays,
				MaxRegisterDays: ar.MaxRegisterDays,
				TargetAmount:    ar.TargetAmount,
				RewardAmount:    ar.RewardAmount,
				Priority:        ar.Priority,
				Status:          1,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			s.agentRewardRepo.Create(policy)
		}
	}

	log.Printf("[PolicyService] Assigned policy to agent: agent=%d, template=%d, operator=%d",
		req.AgentID, req.TemplateID, operatorID)

	return nil
}

// ============================================================
// 下级代理商政策调整（APP端用）
// ============================================================

// UpdateSubordinatePolicyRequest 更新下级代理商政策请求
type UpdateSubordinatePolicyRequest struct {
	SubordinateID int64 `json:"subordinate_id" binding:"required"` // 下级代理商ID
	ChannelID     int64 `json:"channel_id" binding:"required"`

	// 费率配置
	CreditRate   string `json:"credit_rate"`
	DebitRate    string `json:"debit_rate"`
	DebitCap     string `json:"debit_cap"`
	UnionpayRate string `json:"unionpay_rate"`
	WechatRate   string `json:"wechat_rate"`
	AlipayRate   string `json:"alipay_rate"`

	// 押金返现配置
	DepositCashbacks []DepositCashbackInput `json:"deposit_cashbacks"`

	// 流量卡返现配置
	SimCashback *SimCashbackInput `json:"sim_cashback"`

	// 激活奖励配置
	ActivationRewards []ActivationRewardInput `json:"activation_rewards"`
}

// UpdateSubordinatePolicy 更新下级代理商政策（当前代理商调整下级的政策）
func (s *PolicyService) UpdateSubordinatePolicy(operatorID int64, req *UpdateSubordinatePolicyRequest) error {
	// 验证操作者存在
	operator, err := s.agentRepo.FindByID(operatorID)
	if err != nil || operator == nil {
		return fmt.Errorf("操作者不存在: %d", operatorID)
	}

	// 验证下级代理商存在且是操作者的直属下级
	subordinate, err := s.agentRepo.FindByID(req.SubordinateID)
	if err != nil || subordinate == nil {
		return fmt.Errorf("下级代理商不存在: %d", req.SubordinateID)
	}

	if subordinate.ParentID != operatorID {
		return errors.New("只能调整直属下级的政策")
	}

	// 获取操作者的政策配置（用于验证范围限制）
	operatorPolicy, err := s.GetAgentPolicy(operatorID, req.ChannelID)
	if err != nil {
		return fmt.Errorf("获取操作者政策失败: %w", err)
	}

	// 验证政策范围（下级不能超过上级）
	if err := s.validatePolicyRange(operatorPolicy, req); err != nil {
		return err
	}

	// 保存押金返现政策
	if len(req.DepositCashbacks) > 0 {
		s.agentDepositRepo.DeleteByAgentAndChannel(req.SubordinateID, req.ChannelID)
		for _, dc := range req.DepositCashbacks {
			policy := &models.AgentDepositCashbackPolicy{
				AgentID:        req.SubordinateID,
				ChannelID:      req.ChannelID,
				DepositAmount:  dc.DepositAmount,
				CashbackAmount: dc.CashbackAmount,
				Status:         1,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			s.agentDepositRepo.Create(policy)
		}
	}

	// 保存流量卡返现政策
	if req.SimCashback != nil {
		policy := &models.AgentSimCashbackPolicy{
			AgentID:            req.SubordinateID,
			ChannelID:          req.ChannelID,
			FirstTimeCashback:  req.SimCashback.FirstTimeCashback,
			SecondTimeCashback: req.SimCashback.SecondTimeCashback,
			ThirdPlusCashback:  req.SimCashback.ThirdPlusCashback,
			Status:             1,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		s.agentSimRepo.Upsert(policy)
	}

	// 保存激活奖励政策
	if len(req.ActivationRewards) > 0 {
		s.agentRewardRepo.DeleteByAgentAndChannel(req.SubordinateID, req.ChannelID)
		for _, ar := range req.ActivationRewards {
			policy := &models.AgentActivationRewardPolicy{
				AgentID:         req.SubordinateID,
				ChannelID:       req.ChannelID,
				RewardName:      ar.RewardName,
				MinRegisterDays: ar.MinRegisterDays,
				MaxRegisterDays: ar.MaxRegisterDays,
				TargetAmount:    ar.TargetAmount,
				RewardAmount:    ar.RewardAmount,
				Priority:        ar.Priority,
				Status:          1,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			s.agentRewardRepo.Create(policy)
		}
	}

	log.Printf("[PolicyService] Updated subordinate policy: operator=%d, subordinate=%d, channel=%d",
		operatorID, req.SubordinateID, req.ChannelID)

	return nil
}

// GetAgentPolicy 获取代理商完整政策配置
func (s *PolicyService) GetAgentPolicy(agentID, channelID int64) (*models.AgentPolicyComplete, error) {
	policy := &models.AgentPolicyComplete{
		AgentID:   agentID,
		ChannelID: channelID,
	}

	// 获取基础政策
	agentPolicy, err := s.agentPolicyRepo.FindByAgentAndChannel(agentID, channelID)
	if err == nil && agentPolicy != nil {
		policy.CreditRate = agentPolicy.CreditRate
		policy.DebitRate = agentPolicy.DebitRate
	}

	// 获取押金返现配置
	depositPolicies, _ := s.agentDepositRepo.FindByAgentAndChannel(agentID, channelID)
	policy.DepositCashbacks = make([]models.DepositCashbackConfig, len(depositPolicies))
	for i, p := range depositPolicies {
		policy.DepositCashbacks[i] = models.DepositCashbackConfig{
			DepositAmount:  p.DepositAmount,
			CashbackAmount: p.CashbackAmount,
		}
	}

	// 获取流量卡返现配置
	simPolicy, err := s.agentSimRepo.FindByAgentAndChannel(agentID, channelID)
	if err == nil && simPolicy != nil {
		policy.SimCashback = &models.SimCashbackConfig{
			FirstTimeCashback:  simPolicy.FirstTimeCashback,
			SecondTimeCashback: simPolicy.SecondTimeCashback,
			ThirdPlusCashback:  simPolicy.ThirdPlusCashback,
		}
	}

	// 获取激活奖励配置
	rewardPolicies, _ := s.agentRewardRepo.FindByAgentAndChannel(agentID, channelID)
	policy.ActivationRewards = make([]models.ActivationRewardConfig, len(rewardPolicies))
	for i, p := range rewardPolicies {
		policy.ActivationRewards[i] = models.ActivationRewardConfig{
			RewardName:      p.RewardName,
			MinRegisterDays: p.MinRegisterDays,
			MaxRegisterDays: p.MaxRegisterDays,
			TargetAmount:    p.TargetAmount,
			RewardAmount:    p.RewardAmount,
		}
	}

	return policy, nil
}

// GetPolicyLimits 获取政策限制（用于前端显示可调整范围）
func (s *PolicyService) GetPolicyLimits(agentID, channelID int64) (*models.PolicyLimits, error) {
	policy, err := s.GetAgentPolicy(agentID, channelID)
	if err != nil {
		return nil, err
	}

	limits := &models.PolicyLimits{
		MinCreditRate:   policy.CreditRate,
		MinDebitRate:    policy.DebitRate,
		MinUnionpayRate: policy.UnionpayRate,
		MinWechatRate:   policy.WechatRate,
		MinAlipayRate:   policy.AlipayRate,
	}

	// 押金返现限制
	limits.MaxDepositCashbacks = policy.DepositCashbacks

	// 流量卡返现限制
	limits.MaxSimCashback = policy.SimCashback

	// 激活奖励限制
	limits.MaxActivationRewards = policy.ActivationRewards

	return limits, nil
}

// ============================================================
// 辅助方法
// ============================================================

func (s *PolicyService) createDepositCashbackPolicies(templateID, channelID int64, inputs []DepositCashbackInput) error {
	for _, input := range inputs {
		policy := &models.DepositCashbackPolicy{
			TemplateID:     templateID,
			ChannelID:      channelID,
			DepositAmount:  input.DepositAmount,
			CashbackAmount: input.CashbackAmount,
			Status:         1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.depositPolicyRepo.Create(policy); err != nil {
			return err
		}
	}
	return nil
}

func (s *PolicyService) createSimCashbackPolicy(templateID, channelID int64, input *SimCashbackInput) error {
	policy := &models.SimCashbackPolicy{
		TemplateID:         templateID,
		ChannelID:          channelID,
		FirstTimeCashback:  input.FirstTimeCashback,
		SecondTimeCashback: input.SecondTimeCashback,
		ThirdPlusCashback:  input.ThirdPlusCashback,
		SimFeeAmount:       input.SimFeeAmount,
		Status:             1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	return s.simPolicyRepo.Create(policy)
}

func (s *PolicyService) createActivationRewardPolicies(templateID, channelID int64, inputs []ActivationRewardInput) error {
	for _, input := range inputs {
		policy := &models.ActivationRewardPolicy{
			TemplateID:      templateID,
			ChannelID:       channelID,
			RewardName:      input.RewardName,
			MinRegisterDays: input.MinRegisterDays,
			MaxRegisterDays: input.MaxRegisterDays,
			TargetAmount:    input.TargetAmount,
			RewardAmount:    input.RewardAmount,
			Priority:        input.Priority,
			Status:          1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := s.rewardPolicyRepo.Create(policy); err != nil {
			return err
		}
	}
	return nil
}

func (s *PolicyService) createRateStagePolicies(templateID, channelID int64, inputs []RateStageInput) error {
	for _, input := range inputs {
		policy := &models.RateStagePolicy{
			TemplateID:        templateID,
			ChannelID:         channelID,
			StageName:         input.StageName,
			ApplyTo:           input.ApplyTo,
			MinDays:           input.MinDays,
			MaxDays:           input.MaxDays,
			RateDeltas:        input.RateDeltas,
			CreditRateDelta:   input.CreditRateDelta,
			DebitRateDelta:    input.DebitRateDelta,
			UnionpayRateDelta: input.UnionpayRateDelta,
			WechatRateDelta:   input.WechatRateDelta,
			AlipayRateDelta:   input.AlipayRateDelta,
			Priority:          input.Priority,
			Status:            1,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		if err := s.rateStagePolicyRepo.Create(policy); err != nil {
			return err
		}
	}
	return nil
}

func (s *PolicyService) validatePolicyRange(operatorPolicy *models.AgentPolicyComplete, req *UpdateSubordinatePolicyRequest) error {
	// 验证押金返现不超过上级
	for _, dc := range req.DepositCashbacks {
		for _, opDc := range operatorPolicy.DepositCashbacks {
			if dc.DepositAmount == opDc.DepositAmount && dc.CashbackAmount > opDc.CashbackAmount {
				return fmt.Errorf("押金%d元的返现金额不能超过上级配置", dc.DepositAmount/100)
			}
		}
	}

	// 验证流量卡返现不超过上级
	if req.SimCashback != nil && operatorPolicy.SimCashback != nil {
		if req.SimCashback.FirstTimeCashback > operatorPolicy.SimCashback.FirstTimeCashback {
			return errors.New("首次流量费返现不能超过上级配置")
		}
		if req.SimCashback.SecondTimeCashback > operatorPolicy.SimCashback.SecondTimeCashback {
			return errors.New("第2次流量费返现不能超过上级配置")
		}
		if req.SimCashback.ThirdPlusCashback > operatorPolicy.SimCashback.ThirdPlusCashback {
			return errors.New("第3次及以后流量费返现不能超过上级配置")
		}
	}

	// 验证激活奖励不超过上级
	for _, ar := range req.ActivationRewards {
		for _, opAr := range operatorPolicy.ActivationRewards {
			if ar.MinRegisterDays == opAr.MinRegisterDays &&
				ar.MaxRegisterDays == opAr.MaxRegisterDays &&
				ar.TargetAmount == opAr.TargetAmount &&
				ar.RewardAmount > opAr.RewardAmount {
				return fmt.Errorf("激活奖励「%s」的金额不能超过上级配置", ar.RewardName)
			}
		}
	}

	return nil
}

// ============================================================
// 通道约束校验
// ============================================================

// validateTemplateAgainstChannel 校验模板配置是否符合通道约束
func (s *PolicyService) validateTemplateAgainstChannel(ctx context.Context, channelID int64, req *CreatePolicyTemplateRequest) error {
	if s.channelConfigRepo == nil {
		// 如果没有注入通道配置仓库，跳过校验
		return nil
	}

	// 1. 获取通道完整配置
	channelConfig, err := s.channelConfigRepo.GetFullConfig(ctx, channelID)
	if err != nil {
		log.Printf("[PolicyService] Failed to get channel config: %v", err)
		return nil // 获取失败时不阻断流程，仅记录日志
	}

	// 2. 校验费率配置
	for rateCode, rateConfig := range req.RateConfigs {
		channelRateConfig := findChannelRateConfig(channelConfig.RateConfigs, rateCode)
		if channelRateConfig != nil {
			if err := models.ValidateRateRange(rateConfig.Rate, channelRateConfig.MinRate, channelRateConfig.MaxRate); err != nil {
				return fmt.Errorf("%s费率校验失败: %w", channelRateConfig.RateName, err)
			}
		}
	}

	// 3. 校验押金返现配置
	for _, dc := range req.DepositCashbacks {
		depositTier := findChannelDepositTier(channelConfig.DepositTiers, dc.DepositAmount)
		if depositTier != nil && depositTier.MaxCashbackAmount > 0 {
			if dc.CashbackAmount > depositTier.MaxCashbackAmount {
				return fmt.Errorf("押金%d元的返现%d元超过通道上限%d元",
					dc.DepositAmount/100, dc.CashbackAmount/100, depositTier.MaxCashbackAmount/100)
			}
		}
	}

	// 4. 校验流量费返现配置
	if req.SimCashback != nil {
		// 校验首次返现
		tier1 := findChannelSimCashbackTier(channelConfig.SimCashbackTiers, 1)
		if tier1 != nil && req.SimCashback.FirstTimeCashback > tier1.MaxCashbackAmount {
			return fmt.Errorf("首次流量费返现%d元超过通道上限%d元",
				req.SimCashback.FirstTimeCashback/100, tier1.MaxCashbackAmount/100)
		}
		// 校验第2次返现
		tier2 := findChannelSimCashbackTier(channelConfig.SimCashbackTiers, 2)
		if tier2 != nil && req.SimCashback.SecondTimeCashback > tier2.MaxCashbackAmount {
			return fmt.Errorf("第2次流量费返现%d元超过通道上限%d元",
				req.SimCashback.SecondTimeCashback/100, tier2.MaxCashbackAmount/100)
		}
		// 校验第3次及以后返现
		tier3 := findChannelSimCashbackTier(channelConfig.SimCashbackTiers, 3)
		if tier3 != nil && req.SimCashback.ThirdPlusCashback > tier3.MaxCashbackAmount {
			return fmt.Errorf("第3次及以后流量费返现%d元超过通道上限%d元",
				req.SimCashback.ThirdPlusCashback/100, tier3.MaxCashbackAmount/100)
		}
	}

	return nil
}

// findChannelRateConfig 查找通道费率配置
func findChannelRateConfig(configs []models.ChannelRateConfig, rateCode string) *models.ChannelRateConfig {
	for i := range configs {
		if configs[i].RateCode == rateCode {
			return &configs[i]
		}
	}
	return nil
}

// findChannelDepositTier 查找通道押金档位
func findChannelDepositTier(tiers []models.ChannelDepositTier, depositAmount int64) *models.ChannelDepositTier {
	for i := range tiers {
		if tiers[i].DepositAmount == depositAmount {
			return &tiers[i]
		}
	}
	return nil
}

// findChannelSimCashbackTier 查找通道流量费返现档位
func findChannelSimCashbackTier(tiers []models.ChannelSimCashbackTier, tierOrder int) *models.ChannelSimCashbackTier {
	for i := range tiers {
		if tiers[i].TierOrder == tierOrder {
			return &tiers[i]
		}
	}
	return nil
}
