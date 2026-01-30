package service

import (
	"time"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// GlobalWithdrawThresholdService 全局提现门槛服务
type GlobalWithdrawThresholdService struct {
	repo repository.GlobalWithdrawThresholdRepository
}

// NewGlobalWithdrawThresholdService 创建服务实例
func NewGlobalWithdrawThresholdService(repo repository.GlobalWithdrawThresholdRepository) *GlobalWithdrawThresholdService {
	return &GlobalWithdrawThresholdService{repo: repo}
}

// ThresholdConfig 门槛配置响应结构
type ThresholdConfig struct {
	ID              int64  `json:"id"`
	WalletType      int16  `json:"wallet_type"`
	WalletTypeName  string `json:"wallet_type_name"`
	ChannelID       int64  `json:"channel_id"`
	ChannelName     string `json:"channel_name,omitempty"`
	ThresholdAmount int64  `json:"threshold_amount"` // 分
}

// ThresholdListResponse 门槛列表响应
type ThresholdListResponse struct {
	GeneralThresholds []ThresholdConfig `json:"general_thresholds"` // 通用门槛(channel_id=0)
	ChannelThresholds []ThresholdConfig `json:"channel_thresholds"` // 按通道门槛
}

// GetAllThresholds 获取所有门槛配置
func (s *GlobalWithdrawThresholdService) GetAllThresholds() (*ThresholdListResponse, error) {
	thresholds, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	response := &ThresholdListResponse{
		GeneralThresholds: make([]ThresholdConfig, 0),
		ChannelThresholds: make([]ThresholdConfig, 0),
	}

	for _, t := range thresholds {
		config := ThresholdConfig{
			ID:              t.ID,
			WalletType:      t.WalletType,
			WalletTypeName:  models.WalletTypeName(t.WalletType),
			ChannelID:       t.ChannelID,
			ThresholdAmount: t.ThresholdAmount,
		}

		if t.ChannelID == 0 {
			response.GeneralThresholds = append(response.GeneralThresholds, config)
		} else {
			response.ChannelThresholds = append(response.ChannelThresholds, config)
		}
	}

	return response, nil
}

// SetGeneralThresholdRequest 设置通用门槛请求
type SetGeneralThresholdRequest struct {
	ProfitThreshold     int64 `json:"profit_threshold"`      // 分润钱包门槛（分）
	ServiceFeeThreshold int64 `json:"service_fee_threshold"` // 服务费钱包门槛（分）
	RewardThreshold     int64 `json:"reward_threshold"`      // 奖励钱包门槛（分）
}

// SetGeneralThresholds 设置通用门槛
func (s *GlobalWithdrawThresholdService) SetGeneralThresholds(req *SetGeneralThresholdRequest) error {
	now := time.Now()
	thresholds := []models.GlobalWithdrawThreshold{
		{
			WalletType:      models.WalletTypeProfit,
			ChannelID:       0,
			ThresholdAmount: req.ProfitThreshold,
			UpdatedAt:       now,
		},
		{
			WalletType:      models.WalletTypeServiceFee,
			ChannelID:       0,
			ThresholdAmount: req.ServiceFeeThreshold,
			UpdatedAt:       now,
		},
		{
			WalletType:      models.WalletTypeReward,
			ChannelID:       0,
			ThresholdAmount: req.RewardThreshold,
			UpdatedAt:       now,
		},
	}

	return s.repo.UpsertBatch(thresholds)
}

// SetChannelThresholdRequest 设置通道门槛请求
type SetChannelThresholdRequest struct {
	ChannelID           int64 `json:"channel_id" binding:"required"`
	ProfitThreshold     int64 `json:"profit_threshold"`      // 分润钱包门槛（分），0表示使用通用门槛
	ServiceFeeThreshold int64 `json:"service_fee_threshold"` // 服务费钱包门槛（分），0表示使用通用门槛
	RewardThreshold     int64 `json:"reward_threshold"`      // 奖励钱包门槛（分），0表示使用通用门槛
}

// SetChannelThresholds 设置通道门槛
func (s *GlobalWithdrawThresholdService) SetChannelThresholds(req *SetChannelThresholdRequest) error {
	now := time.Now()
	var thresholds []models.GlobalWithdrawThreshold

	// 只有非零值才设置通道门槛
	if req.ProfitThreshold > 0 {
		thresholds = append(thresholds, models.GlobalWithdrawThreshold{
			WalletType:      models.WalletTypeProfit,
			ChannelID:       req.ChannelID,
			ThresholdAmount: req.ProfitThreshold,
			UpdatedAt:       now,
		})
	}

	if req.ServiceFeeThreshold > 0 {
		thresholds = append(thresholds, models.GlobalWithdrawThreshold{
			WalletType:      models.WalletTypeServiceFee,
			ChannelID:       req.ChannelID,
			ThresholdAmount: req.ServiceFeeThreshold,
			UpdatedAt:       now,
		})
	}

	if req.RewardThreshold > 0 {
		thresholds = append(thresholds, models.GlobalWithdrawThreshold{
			WalletType:      models.WalletTypeReward,
			ChannelID:       req.ChannelID,
			ThresholdAmount: req.RewardThreshold,
			UpdatedAt:       now,
		})
	}

	if len(thresholds) == 0 {
		// 如果所有门槛都是0，删除该通道的所有门槛配置
		return s.repo.DeleteByChannel(req.ChannelID)
	}

	return s.repo.UpsertBatch(thresholds)
}

// DeleteChannelThreshold 删除通道门槛配置
func (s *GlobalWithdrawThresholdService) DeleteChannelThreshold(channelID int64) error {
	return s.repo.DeleteByChannel(channelID)
}

// GetWithdrawThreshold 获取提现门槛金额
// 优先级：特定通道门槛 > 通用门槛 > 系统默认值
func (s *GlobalWithdrawThresholdService) GetWithdrawThreshold(walletType int16, channelID int64) int64 {
	threshold, err := s.repo.GetThreshold(walletType, channelID)
	if err == nil && threshold != nil {
		return threshold.ThresholdAmount
	}

	// 返回系统默认值
	return s.getDefaultThreshold(walletType)
}

// getDefaultThreshold 获取默认门槛
func (s *GlobalWithdrawThresholdService) getDefaultThreshold(walletType int16) int64 {
	switch walletType {
	case models.WalletTypeProfit:
		return 10000 // 分润钱包默认100元
	case models.WalletTypeServiceFee:
		return 5000 // 服务费钱包默认50元
	case models.WalletTypeReward:
		return 10000 // 奖励钱包默认100元
	default:
		return 10000 // 默认100元
	}
}
