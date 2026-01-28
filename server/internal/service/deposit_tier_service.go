package service

import (
	"errors"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// DepositTierService 押金档位服务
type DepositTierService struct {
	tierRepo repository.ChannelDepositTierRepository
}

// NewDepositTierService 创建押金档位服务
func NewDepositTierService(tierRepo repository.ChannelDepositTierRepository) *DepositTierService {
	return &DepositTierService{tierRepo: tierRepo}
}

// CreateDepositTierRequest 创建押金档位请求
type CreateDepositTierRequest struct {
	ChannelID     int64  `json:"channel_id" binding:"required"`
	BrandCode     string `json:"brand_code"`
	TierCode      string `json:"tier_code" binding:"required"`
	DepositAmount int64  `json:"deposit_amount" binding:"required"`
	TierName      string `json:"tier_name" binding:"required"`
	SortOrder     int    `json:"sort_order"`
}

// UpdateDepositTierRequest 更新押金档位请求
type UpdateDepositTierRequest struct {
	TierCode      string `json:"tier_code"`
	DepositAmount int64  `json:"deposit_amount"`
	TierName      string `json:"tier_name"`
	SortOrder     int    `json:"sort_order"`
	Status        *int16 `json:"status"`
}

// Create 创建押金档位
func (s *DepositTierService) Create(req *CreateDepositTierRequest) (*models.ChannelDepositTier, error) {
	// 检查是否已存在相同编码
	existing, _ := s.tierRepo.FindByTierCode(req.ChannelID, req.BrandCode, req.TierCode)
	if existing != nil {
		return nil, errors.New("档位编码已存在")
	}

	tier := &models.ChannelDepositTier{
		ChannelID:     req.ChannelID,
		BrandCode:     req.BrandCode,
		TierCode:      req.TierCode,
		DepositAmount: req.DepositAmount,
		TierName:      req.TierName,
		SortOrder:     req.SortOrder,
		Status:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.tierRepo.Create(tier); err != nil {
		return nil, err
	}

	return tier, nil
}

// Update 更新押金档位
func (s *DepositTierService) Update(id int64, req *UpdateDepositTierRequest) (*models.ChannelDepositTier, error) {
	tier, err := s.tierRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("押金档位不存在")
	}

	if req.TierCode != "" && req.TierCode != tier.TierCode {
		// 检查新编码是否冲突
		existing, _ := s.tierRepo.FindByTierCode(tier.ChannelID, tier.BrandCode, req.TierCode)
		if existing != nil && existing.ID != id {
			return nil, errors.New("档位编码已存在")
		}
		tier.TierCode = req.TierCode
	}

	if req.DepositAmount > 0 {
		tier.DepositAmount = req.DepositAmount
	}
	if req.TierName != "" {
		tier.TierName = req.TierName
	}
	if req.SortOrder > 0 {
		tier.SortOrder = req.SortOrder
	}
	if req.Status != nil {
		tier.Status = *req.Status
	}

	tier.UpdatedAt = time.Now()

	if err := s.tierRepo.Update(tier); err != nil {
		return nil, err
	}

	return tier, nil
}

// Delete 删除押金档位
func (s *DepositTierService) Delete(id int64) error {
	_, err := s.tierRepo.FindByID(id)
	if err != nil {
		return errors.New("押金档位不存在")
	}
	return s.tierRepo.Delete(id)
}

// GetByID 根据ID获取
func (s *DepositTierService) GetByID(id int64) (*models.ChannelDepositTier, error) {
	return s.tierRepo.FindByID(id)
}

// GetByChannelID 根据通道ID获取列表
func (s *DepositTierService) GetByChannelID(channelID int64) ([]*models.ChannelDepositTier, error) {
	return s.tierRepo.FindByChannelID(channelID)
}

// GetByChannelAndBrand 根据通道ID和品牌获取列表
func (s *DepositTierService) GetByChannelAndBrand(channelID int64, brandCode string) ([]*models.ChannelDepositTier, error) {
	return s.tierRepo.FindByChannelAndBrand(channelID, brandCode)
}

// GetDepositAmountByTierCode 根据档位编码获取押金金额
func (s *DepositTierService) GetDepositAmountByTierCode(channelID int64, brandCode string, tierCode string) (int64, error) {
	tier, err := s.tierRepo.FindByTierCode(channelID, brandCode, tierCode)
	if err != nil {
		return 0, err
	}
	return tier.DepositAmount, nil
}
