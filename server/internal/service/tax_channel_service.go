package service

import (
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// TaxChannelService 税筹通道服务
type TaxChannelService struct {
	taxChannelRepo *repository.GormTaxChannelRepository
}

// NewTaxChannelService 创建税筹通道服务
func NewTaxChannelService(taxChannelRepo *repository.GormTaxChannelRepository) *TaxChannelService {
	return &TaxChannelService{
		taxChannelRepo: taxChannelRepo,
	}
}

// ========== 税筹通道管理 ==========

// CreateTaxChannelRequest 创建税筹通道请求
type CreateTaxChannelRequest struct {
	ChannelCode string  `json:"channel_code" binding:"required"`
	ChannelName string  `json:"channel_name" binding:"required"`
	FeeType     int16   `json:"fee_type" binding:"required,oneof=1 2"` // 1=付款扣 2=出款扣
	TaxRate     float64 `json:"tax_rate" binding:"required,min=0,max=1"`
	FixedFee    int64   `json:"fixed_fee"`
	ApiURL      string  `json:"api_url"`
	ApiKey      string  `json:"api_key"`
	ApiSecret   string  `json:"api_secret"`
	Remark      string  `json:"remark"`
}

// CreateTaxChannel 创建税筹通道
func (s *TaxChannelService) CreateTaxChannel(req *CreateTaxChannelRequest) (*models.TaxChannel, error) {
	// 检查编码是否已存在
	existing, err := s.taxChannelRepo.GetByCode(req.ChannelCode)
	if err != nil {
		return nil, fmt.Errorf("检查通道编码失败: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("通道编码已存在: %s", req.ChannelCode)
	}

	now := time.Now()
	taxChannel := &models.TaxChannel{
		ChannelCode: req.ChannelCode,
		ChannelName: req.ChannelName,
		FeeType:     req.FeeType,
		TaxRate:     req.TaxRate,
		FixedFee:    req.FixedFee,
		ApiURL:      req.ApiURL,
		ApiKey:      req.ApiKey,
		ApiSecret:   req.ApiSecret,
		Status:      models.TaxChannelStatusEnabled,
		Remark:      req.Remark,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taxChannelRepo.Create(taxChannel); err != nil {
		return nil, fmt.Errorf("创建税筹通道失败: %w", err)
	}

	log.Printf("[TaxChannelService] Created tax channel: %s (%s)", req.ChannelCode, req.ChannelName)
	return taxChannel, nil
}

// UpdateTaxChannelRequest 更新税筹通道请求
type UpdateTaxChannelRequest struct {
	ID          int64   `json:"id" binding:"required"`
	ChannelName string  `json:"channel_name"`
	FeeType     int16   `json:"fee_type"`
	TaxRate     float64 `json:"tax_rate"`
	FixedFee    int64   `json:"fixed_fee"`
	ApiURL      string  `json:"api_url"`
	ApiKey      string  `json:"api_key"`
	ApiSecret   string  `json:"api_secret"`
	Status      *int16  `json:"status"`
	Remark      string  `json:"remark"`
}

// UpdateTaxChannel 更新税筹通道
func (s *TaxChannelService) UpdateTaxChannel(req *UpdateTaxChannelRequest) (*models.TaxChannel, error) {
	taxChannel, err := s.taxChannelRepo.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("获取税筹通道失败: %w", err)
	}
	if taxChannel == nil {
		return nil, fmt.Errorf("税筹通道不存在")
	}

	// 更新字段
	if req.ChannelName != "" {
		taxChannel.ChannelName = req.ChannelName
	}
	if req.FeeType > 0 {
		taxChannel.FeeType = req.FeeType
	}
	if req.TaxRate > 0 {
		taxChannel.TaxRate = req.TaxRate
	}
	if req.FixedFee >= 0 {
		taxChannel.FixedFee = req.FixedFee
	}
	if req.ApiURL != "" {
		taxChannel.ApiURL = req.ApiURL
	}
	if req.ApiKey != "" {
		taxChannel.ApiKey = req.ApiKey
	}
	if req.ApiSecret != "" {
		taxChannel.ApiSecret = req.ApiSecret
	}
	if req.Status != nil {
		taxChannel.Status = *req.Status
	}
	if req.Remark != "" {
		taxChannel.Remark = req.Remark
	}
	taxChannel.UpdatedAt = time.Now()

	if err := s.taxChannelRepo.Update(taxChannel); err != nil {
		return nil, fmt.Errorf("更新税筹通道失败: %w", err)
	}

	log.Printf("[TaxChannelService] Updated tax channel: %d", req.ID)
	return taxChannel, nil
}

// GetTaxChannel 获取税筹通道详情
func (s *TaxChannelService) GetTaxChannel(id int64) (*TaxChannelInfo, error) {
	taxChannel, err := s.taxChannelRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取税筹通道失败: %w", err)
	}
	if taxChannel == nil {
		return nil, fmt.Errorf("税筹通道不存在")
	}

	return s.toTaxChannelInfo(taxChannel), nil
}

// GetTaxChannelList 获取税筹通道列表
func (s *TaxChannelService) GetTaxChannelList(status *int16) ([]*TaxChannelInfo, error) {
	taxChannels, err := s.taxChannelRepo.GetAll(status)
	if err != nil {
		return nil, fmt.Errorf("获取税筹通道列表失败: %w", err)
	}

	list := make([]*TaxChannelInfo, 0, len(taxChannels))
	for _, tc := range taxChannels {
		list = append(list, s.toTaxChannelInfo(tc))
	}
	return list, nil
}

// DeleteTaxChannel 删除税筹通道
func (s *TaxChannelService) DeleteTaxChannel(id int64) error {
	taxChannel, err := s.taxChannelRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取税筹通道失败: %w", err)
	}
	if taxChannel == nil {
		return fmt.Errorf("税筹通道不存在")
	}

	// 检查是否有关联的映射
	// TODO: 检查是否有关联的channel_tax_mappings

	if err := s.taxChannelRepo.Delete(id); err != nil {
		return fmt.Errorf("删除税筹通道失败: %w", err)
	}

	log.Printf("[TaxChannelService] Deleted tax channel: %d", id)
	return nil
}

// ========== 通道-税筹通道映射 ==========

// SetChannelTaxMappingRequest 设置通道税筹映射请求
type SetChannelTaxMappingRequest struct {
	ChannelID    int64 `json:"channel_id" binding:"required"`
	WalletType   int16 `json:"wallet_type" binding:"required,min=1,max=5"`
	TaxChannelID int64 `json:"tax_channel_id" binding:"required"`
}

// SetChannelTaxMapping 设置通道税筹映射
func (s *TaxChannelService) SetChannelTaxMapping(req *SetChannelTaxMappingRequest) error {
	// 验证税筹通道是否存在
	taxChannel, err := s.taxChannelRepo.GetByID(req.TaxChannelID)
	if err != nil {
		return fmt.Errorf("获取税筹通道失败: %w", err)
	}
	if taxChannel == nil {
		return fmt.Errorf("税筹通道不存在")
	}

	// 查找是否已有映射
	existing, err := s.taxChannelRepo.GetMappingByChannelAndWallet(req.ChannelID, req.WalletType)
	if err != nil {
		return fmt.Errorf("查询映射失败: %w", err)
	}

	now := time.Now()
	if existing != nil {
		// 更新现有映射
		existing.TaxChannelID = req.TaxChannelID
		existing.UpdatedAt = now
		if err := s.taxChannelRepo.UpdateMapping(existing); err != nil {
			return fmt.Errorf("更新映射失败: %w", err)
		}
	} else {
		// 创建新映射
		mapping := &models.ChannelTaxMapping{
			ChannelID:    req.ChannelID,
			WalletType:   req.WalletType,
			TaxChannelID: req.TaxChannelID,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.taxChannelRepo.CreateMapping(mapping); err != nil {
			return fmt.Errorf("创建映射失败: %w", err)
		}
	}

	log.Printf("[TaxChannelService] Set tax mapping: channel=%d, wallet=%d, taxChannel=%d",
		req.ChannelID, req.WalletType, req.TaxChannelID)
	return nil
}

// GetChannelTaxMappings 获取通道的税筹映射
func (s *TaxChannelService) GetChannelTaxMappings(channelID int64) ([]*ChannelTaxMappingInfo, error) {
	mappings, err := s.taxChannelRepo.GetMappingsByChannel(channelID)
	if err != nil {
		return nil, fmt.Errorf("获取映射失败: %w", err)
	}

	list := make([]*ChannelTaxMappingInfo, 0, len(mappings))
	for _, m := range mappings {
		taxChannel, _ := s.taxChannelRepo.GetByID(m.TaxChannelID)
		info := &ChannelTaxMappingInfo{
			ID:           m.ID,
			ChannelID:    m.ChannelID,
			WalletType:   m.WalletType,
			WalletTypeName: getTaxWalletTypeName(m.WalletType),
			TaxChannelID: m.TaxChannelID,
			CreatedAt:    m.CreatedAt,
		}
		if taxChannel != nil {
			info.TaxChannelName = taxChannel.ChannelName
			info.TaxRate = taxChannel.TaxRate
			info.FixedFee = taxChannel.FixedFee
		}
		list = append(list, info)
	}
	return list, nil
}

// DeleteChannelTaxMapping 删除通道税筹映射
func (s *TaxChannelService) DeleteChannelTaxMapping(id int64) error {
	if err := s.taxChannelRepo.DeleteMapping(id); err != nil {
		return fmt.Errorf("删除映射失败: %w", err)
	}
	log.Printf("[TaxChannelService] Deleted tax mapping: %d", id)
	return nil
}

// ========== 税费计算 ==========

// TaxCalculationResult 税费计算结果
type TaxCalculationResult struct {
	OriginalAmount int64   `json:"original_amount"`      // 原金额(分)
	TaxRate        float64 `json:"tax_rate"`             // 税率
	TaxFee         int64   `json:"tax_fee"`              // 税费(分)
	FixedFee       int64   `json:"fixed_fee"`            // 固定费用(分)
	TotalFee       int64   `json:"total_fee"`            // 总费用(分)
	ActualAmount   int64   `json:"actual_amount"`        // 实际到账(分)
	TaxChannelID   int64   `json:"tax_channel_id"`       // 税筹通道ID
	TaxChannelName string  `json:"tax_channel_name"`     // 税筹通道名称
}

// CalculateWithdrawalTax 计算提现税费
func (s *TaxChannelService) CalculateWithdrawalTax(channelID int64, walletType int16, amount int64) (*TaxCalculationResult, error) {
	// 获取对应的税筹通道
	taxChannel, err := s.taxChannelRepo.GetTaxChannelForWithdrawal(channelID, walletType)
	if err != nil {
		return nil, fmt.Errorf("获取税筹通道失败: %w", err)
	}

	result := &TaxCalculationResult{
		OriginalAmount: amount,
	}

	if taxChannel == nil {
		// 没有配置税筹通道，不扣税
		result.ActualAmount = amount
		return result, nil
	}

	// 计算税费
	taxFee := int64(float64(amount) * taxChannel.TaxRate)
	fixedFee := taxChannel.FixedFee
	totalFee := taxFee + fixedFee
	actualAmount := amount - totalFee

	// 实际到账不能为负
	if actualAmount < 0 {
		actualAmount = 0
		totalFee = amount
	}

	result.TaxRate = taxChannel.TaxRate
	result.TaxFee = taxFee
	result.FixedFee = fixedFee
	result.TotalFee = totalFee
	result.ActualAmount = actualAmount
	result.TaxChannelID = taxChannel.ID
	result.TaxChannelName = taxChannel.ChannelName

	return result, nil
}

// ========== 数据转换 ==========

// TaxChannelInfo 税筹通道信息
type TaxChannelInfo struct {
	ID             int64     `json:"id"`
	ChannelCode    string    `json:"channel_code"`
	ChannelName    string    `json:"channel_name"`
	FeeType        int16     `json:"fee_type"`
	FeeTypeName    string    `json:"fee_type_name"`
	TaxRate        float64   `json:"tax_rate"`
	TaxRatePercent float64   `json:"tax_rate_percent"` // 百分比显示
	FixedFee       int64     `json:"fixed_fee"`
	FixedFeeYuan   float64   `json:"fixed_fee_yuan"`
	Status         int16     `json:"status"`
	StatusName     string    `json:"status_name"`
	Remark         string    `json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s *TaxChannelService) toTaxChannelInfo(tc *models.TaxChannel) *TaxChannelInfo {
	return &TaxChannelInfo{
		ID:             tc.ID,
		ChannelCode:    tc.ChannelCode,
		ChannelName:    tc.ChannelName,
		FeeType:        tc.FeeType,
		FeeTypeName:    models.GetFeeTypeName(tc.FeeType),
		TaxRate:        tc.TaxRate,
		TaxRatePercent: tc.TaxRate * 100,
		FixedFee:       tc.FixedFee,
		FixedFeeYuan:   float64(tc.FixedFee) / 100,
		Status:         tc.Status,
		StatusName:     models.GetTaxChannelStatusName(tc.Status),
		Remark:         tc.Remark,
		CreatedAt:      tc.CreatedAt,
		UpdatedAt:      tc.UpdatedAt,
	}
}

// ChannelTaxMappingInfo 通道税筹映射信息
type ChannelTaxMappingInfo struct {
	ID             int64     `json:"id"`
	ChannelID      int64     `json:"channel_id"`
	WalletType     int16     `json:"wallet_type"`
	WalletTypeName string    `json:"wallet_type_name"`
	TaxChannelID   int64     `json:"tax_channel_id"`
	TaxChannelName string    `json:"tax_channel_name"`
	TaxRate        float64   `json:"tax_rate"`
	FixedFee       int64     `json:"fixed_fee"`
	CreatedAt      time.Time `json:"created_at"`
}

func getTaxWalletTypeName(walletType int16) string {
	switch walletType {
	case models.WalletTypeProfit:
		return "分润钱包"
	case models.WalletTypeServiceFee:
		return "服务费钱包"
	case models.WalletTypeReward:
		return "奖励钱包"
	case models.WalletTypeCharging:
		return "充值钱包"
	case models.WalletTypeSettlement:
		return "沉淀钱包"
	default:
		return "未知"
	}
}
