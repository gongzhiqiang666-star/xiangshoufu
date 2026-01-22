package service

import (
	"log"
	"strconv"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// RateStagingService 费率阶梯服务
// 根据商户/代理商入网时间自动调整费率
type RateStagingService struct {
	rateStagePolicyRepo *repository.GormRateStagePolicyRepository
	merchantRepo        *repository.GormMerchantRepository
	agentRepo           repository.AgentRepository
}

// NewRateStagingService 创建费率阶梯服务
func NewRateStagingService(
	rateStagePolicyRepo *repository.GormRateStagePolicyRepository,
	merchantRepo *repository.GormMerchantRepository,
	agentRepo repository.AgentRepository,
) *RateStagingService {
	return &RateStagingService{
		rateStagePolicyRepo: rateStagePolicyRepo,
		merchantRepo:        merchantRepo,
		agentRepo:           agentRepo,
	}
}

// RateAdjustment 费率调整结果
type RateAdjustment struct {
	OriginalRate   float64 `json:"original_rate"`    // 原始费率
	AdjustedRate   float64 `json:"adjusted_rate"`    // 调整后费率
	RateDelta      float64 `json:"rate_delta"`       // 费率调整值
	PolicyID       int64   `json:"policy_id"`        // 应用的政策ID
	PolicyName     string  `json:"policy_name"`      // 政策名称
	RegisterDays   int     `json:"register_days"`    // 入网天数
	ApplyTo        int16   `json:"apply_to"`         // 应用对象 1-商户 2-代理商
}

// GetMerchantRateAdjustment 获取商户费率调整
// 根据商户入网时间计算费率调整值
func (s *RateStagingService) GetMerchantRateAdjustment(
	merchantID int64,
	channelID int64,
	baseRate float64,
	cardType int16, // 1-借记卡 2-贷记卡
) (*RateAdjustment, error) {
	// 获取商户信息
	merchant, err := s.merchantRepo.FindByID(merchantID)
	if err != nil {
		return nil, err
	}

	// 计算入网天数
	registerDays := s.calculateDaysFromPtr(merchant.ActivatedAt)

	// 查找适用的费率阶梯政策
	policy, err := s.rateStagePolicyRepo.FindApplicablePolicy(
		channelID,
		models.RateStageApplyToMerchant, // 商户
		registerDays,
	)
	if err != nil {
		// 没有找到适用的政策，返回原始费率
		return &RateAdjustment{
			OriginalRate: baseRate,
			AdjustedRate: baseRate,
			RateDelta:    0,
			RegisterDays: registerDays,
			ApplyTo:      models.RateStageApplyToMerchant,
		}, nil
	}

	// 获取对应卡类型的费率调整值
	rateDelta := s.getRateDelta(policy, cardType)

	// 计算调整后费率
	adjustedRate := baseRate + rateDelta

	log.Printf("[RateStagingService] Merchant %d rate adjustment: base=%.4f, delta=%.4f, adjusted=%.4f, days=%d",
		merchantID, baseRate, rateDelta, adjustedRate, registerDays)

	return &RateAdjustment{
		OriginalRate: baseRate,
		AdjustedRate: adjustedRate,
		RateDelta:    rateDelta,
		PolicyID:     policy.ID,
		PolicyName:   policy.StageName,
		RegisterDays: registerDays,
		ApplyTo:      models.RateStageApplyToMerchant,
	}, nil
}

// GetAgentRateAdjustment 获取代理商费率调整
// 根据代理商入网时间计算费率调整值
func (s *RateStagingService) GetAgentRateAdjustment(
	agentID int64,
	channelID int64,
	baseRate float64,
	cardType int16,
) (*RateAdjustment, error) {
	// 获取代理商信息
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil {
		return nil, err
	}

	// 计算入网天数（使用代理商注册时间）
	// 注意：这里假设 Agent 模型有 CreatedAt 字段
	// 如果需要使用其他字段，需要扩展 AgentRepository 接口
	registerDays := 0 // 默认0天，需要从 AgentFull 获取

	// 查找适用的费率阶梯政策
	policy, err := s.rateStagePolicyRepo.FindApplicablePolicy(
		channelID,
		models.RateStageApplyToAgent, // 代理商
		registerDays,
	)
	if err != nil {
		// 没有找到适用的政策，返回原始费率
		return &RateAdjustment{
			OriginalRate: baseRate,
			AdjustedRate: baseRate,
			RateDelta:    0,
			RegisterDays: registerDays,
			ApplyTo:      models.RateStageApplyToAgent,
		}, nil
	}

	// 获取对应卡类型的费率调整值
	rateDelta := s.getRateDelta(policy, cardType)

	// 计算调整后费率
	adjustedRate := baseRate + rateDelta

	log.Printf("[RateStagingService] Agent %d (%s) rate adjustment: base=%.4f, delta=%.4f, adjusted=%.4f, days=%d",
		agentID, agent.AgentNo, baseRate, rateDelta, adjustedRate, registerDays)

	return &RateAdjustment{
		OriginalRate: baseRate,
		AdjustedRate: adjustedRate,
		RateDelta:    rateDelta,
		PolicyID:     policy.ID,
		PolicyName:   policy.StageName,
		RegisterDays: registerDays,
		ApplyTo:      models.RateStageApplyToAgent,
	}, nil
}

// calculateDays 计算从指定时间到现在的天数
func (s *RateStagingService) calculateDays(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	duration := time.Since(t)
	return int(duration.Hours() / 24)
}

// calculateDaysFromPtr 计算从指定时间指针到现在的天数
func (s *RateStagingService) calculateDaysFromPtr(t *time.Time) int {
	if t == nil || t.IsZero() {
		return 0
	}
	duration := time.Since(*t)
	return int(duration.Hours() / 24)
}

// getRateDelta 根据卡类型获取费率调整值
func (s *RateStagingService) getRateDelta(policy *models.RateStagePolicy, cardType int16) float64 {
	var deltaStr string
	switch cardType {
	case 1: // 借记卡
		deltaStr = policy.DebitRateDelta
	case 2: // 贷记卡
		deltaStr = policy.CreditRateDelta
	case 3: // 云闪付
		deltaStr = policy.UnionpayRateDelta
	case 4: // 微信
		deltaStr = policy.WechatRateDelta
	case 5: // 支付宝
		deltaStr = policy.AlipayRateDelta
	default:
		deltaStr = policy.CreditRateDelta
	}

	delta, _ := strconv.ParseFloat(deltaStr, 64)
	return delta
}

// ApplyRateStaging 应用费率阶梯（综合商户和代理商调整）
// 返回最终应使用的费率
func (s *RateStagingService) ApplyRateStaging(
	merchantID int64,
	agentID int64,
	channelID int64,
	baseRate float64,
	cardType int16,
) (float64, error) {
	finalRate := baseRate

	// 1. 先应用商户入网时间的费率调整
	merchantAdj, err := s.GetMerchantRateAdjustment(merchantID, channelID, baseRate, cardType)
	if err == nil && merchantAdj.RateDelta != 0 {
		finalRate = merchantAdj.AdjustedRate
	}

	// 2. 再应用代理商入网时间的费率调整（在商户调整基础上）
	agentAdj, err := s.GetAgentRateAdjustment(agentID, channelID, finalRate, cardType)
	if err == nil && agentAdj.RateDelta != 0 {
		finalRate = agentAdj.AdjustedRate
	}

	return finalRate, nil
}

// GetApplicableStages 获取所有适用的费率阶梯配置
// 用于前端展示当前生效的费率阶梯
func (s *RateStagingService) GetApplicableStages(channelID int64, applyTo int16) ([]*models.RateStagePolicy, error) {
	return s.rateStagePolicyRepo.FindActiveByChannelAndApplyTo(channelID, applyTo)
}

// PreviewMerchantRateStaging 预览商户费率阶梯
// 返回所有阶梯的费率调整预览
type RateStagingPreview struct {
	StageName   string  `json:"stage_name"`
	MinDays     int     `json:"min_days"`
	MaxDays     int     `json:"max_days"`
	CreditDelta float64 `json:"credit_delta"`
	DebitDelta  float64 `json:"debit_delta"`
}

func (s *RateStagingService) PreviewRateStaging(channelID int64, applyTo int16) ([]RateStagingPreview, error) {
	policies, err := s.rateStagePolicyRepo.FindActiveByChannelAndApplyTo(channelID, applyTo)
	if err != nil {
		return nil, err
	}

	previews := make([]RateStagingPreview, len(policies))
	for i, p := range policies {
		creditDelta, _ := strconv.ParseFloat(p.CreditRateDelta, 64)
		debitDelta, _ := strconv.ParseFloat(p.DebitRateDelta, 64)
		previews[i] = RateStagingPreview{
			StageName:   p.StageName,
			MinDays:     p.MinDays,
			MaxDays:     p.MaxDays,
			CreditDelta: creditDelta,
			DebitDelta:  debitDelta,
		}
	}

	return previews, nil
}
