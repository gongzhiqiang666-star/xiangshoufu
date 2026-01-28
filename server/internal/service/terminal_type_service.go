package service

import (
	"context"
	"errors"
	"fmt"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// TerminalTypeService 终端类型服务
type TerminalTypeService struct {
	repo        repository.TerminalTypeRepository
	channelRepo *repository.GormChannelRepository
}

// NewTerminalTypeService 创建终端类型服务
func NewTerminalTypeService(repo repository.TerminalTypeRepository, channelRepo *repository.GormChannelRepository) *TerminalTypeService {
	return &TerminalTypeService{
		repo:        repo,
		channelRepo: channelRepo,
	}
}

// CreateTerminalTypeRequest 创建终端类型请求
type CreateTerminalTypeRequest struct {
	ChannelID   int64  `json:"channel_id" binding:"required"`
	BrandCode   string `json:"brand_code" binding:"required"`
	BrandName   string `json:"brand_name" binding:"required"`
	ModelCode   string `json:"model_code" binding:"required"`
	ModelName   string `json:"model_name"`
	Description string `json:"description"`
}

// UpdateTerminalTypeRequest 更新终端类型请求
type UpdateTerminalTypeRequest struct {
	BrandCode   string `json:"brand_code"`
	BrandName   string `json:"brand_name"`
	ModelCode   string `json:"model_code"`
	ModelName   string `json:"model_name"`
	Description string `json:"description"`
}

// TerminalTypeListRequest 终端类型列表请求
type TerminalTypeListRequest struct {
	ChannelID int64  `form:"channel_id"`
	Status    *int16 `form:"status"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// TerminalTypeResponse 终端类型响应
type TerminalTypeResponse struct {
	ID          int64  `json:"id"`
	ChannelID   int64  `json:"channel_id"`
	ChannelCode string `json:"channel_code"`
	ChannelName string `json:"channel_name"`
	BrandCode   string `json:"brand_code"`
	BrandName   string `json:"brand_name"`
	ModelCode   string `json:"model_code"`
	ModelName   string `json:"model_name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Status      int16  `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Create 创建终端类型
func (s *TerminalTypeService) Create(ctx context.Context, req *CreateTerminalTypeRequest) (*models.TerminalType, error) {
	// 验证通道是否存在
	channel, err := s.channelRepo.FindByID(req.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("查询通道失败: %w", err)
	}
	if channel == nil {
		return nil, errors.New("通道不存在")
	}

	// 检查是否重复
	existing, err := s.repo.FindByChannelAndCodes(ctx, req.ChannelID, req.BrandCode, req.ModelCode)
	if err != nil {
		return nil, fmt.Errorf("检查重复失败: %w", err)
	}
	if existing != nil {
		return nil, errors.New("该通道下已存在相同品牌和型号的终端类型")
	}

	terminalType := &models.TerminalType{
		ChannelID:   req.ChannelID,
		ChannelCode: channel.ChannelCode,
		BrandCode:   req.BrandCode,
		BrandName:   req.BrandName,
		ModelCode:   req.ModelCode,
		ModelName:   req.ModelName,
		Description: req.Description,
		Status:      models.TerminalTypeStatusEnabled,
	}

	if err := s.repo.Create(ctx, terminalType); err != nil {
		return nil, fmt.Errorf("创建终端类型失败: %w", err)
	}

	return terminalType, nil
}

// Update 更新终端类型
func (s *TerminalTypeService) Update(ctx context.Context, id int64, req *UpdateTerminalTypeRequest) (*models.TerminalType, error) {
	terminalType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询终端类型失败: %w", err)
	}
	if terminalType == nil {
		return nil, errors.New("终端类型不存在")
	}

	// 如果修改了编码，检查是否重复
	if (req.BrandCode != "" && req.BrandCode != terminalType.BrandCode) ||
		(req.ModelCode != "" && req.ModelCode != terminalType.ModelCode) {
		brandCode := terminalType.BrandCode
		modelCode := terminalType.ModelCode
		if req.BrandCode != "" {
			brandCode = req.BrandCode
		}
		if req.ModelCode != "" {
			modelCode = req.ModelCode
		}

		existing, err := s.repo.FindByChannelAndCodes(ctx, terminalType.ChannelID, brandCode, modelCode)
		if err != nil {
			return nil, fmt.Errorf("检查重复失败: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("该通道下已存在相同品牌和型号的终端类型")
		}
	}

	// 更新字段
	if req.BrandCode != "" {
		terminalType.BrandCode = req.BrandCode
	}
	if req.BrandName != "" {
		terminalType.BrandName = req.BrandName
	}
	if req.ModelCode != "" {
		terminalType.ModelCode = req.ModelCode
	}
	if req.ModelName != "" {
		terminalType.ModelName = req.ModelName
	}
	if req.Description != "" {
		terminalType.Description = req.Description
	}

	if err := s.repo.Update(ctx, terminalType); err != nil {
		return nil, fmt.Errorf("更新终端类型失败: %w", err)
	}

	return terminalType, nil
}

// GetByID 根据ID获取终端类型
func (s *TerminalTypeService) GetByID(ctx context.Context, id int64) (*models.TerminalType, error) {
	terminalType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询终端类型失败: %w", err)
	}
	if terminalType == nil {
		return nil, errors.New("终端类型不存在")
	}
	return terminalType, nil
}

// List 获取终端类型列表
func (s *TerminalTypeService) List(ctx context.Context, req *TerminalTypeListRequest) ([]models.TerminalType, int64, error) {
	filter := repository.TerminalTypeFilter{
		ChannelID: req.ChannelID,
		Status:    req.Status,
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	return s.repo.List(ctx, filter)
}

// ListByChannelID 根据通道ID获取终端类型列表（用于下拉选择）
func (s *TerminalTypeService) ListByChannelID(ctx context.Context, channelID int64) ([]models.TerminalType, error) {
	return s.repo.ListByChannelID(ctx, channelID, true)
}

// UpdateStatus 更新状态
func (s *TerminalTypeService) UpdateStatus(ctx context.Context, id int64, status int16) error {
	terminalType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("查询终端类型失败: %w", err)
	}
	if terminalType == nil {
		return errors.New("终端类型不存在")
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	return nil
}

// Delete 删除终端类型
func (s *TerminalTypeService) Delete(ctx context.Context, id int64) error {
	terminalType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("查询终端类型失败: %w", err)
	}
	if terminalType == nil {
		return errors.New("终端类型不存在")
	}

	// TODO: 检查是否有关联的终端，如果有则不允许删除

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除终端类型失败: %w", err)
	}

	return nil
}

// ToResponse 转换为响应
func (s *TerminalTypeService) ToResponse(t *models.TerminalType) *TerminalTypeResponse {
	resp := &TerminalTypeResponse{
		ID:          t.ID,
		ChannelID:   t.ChannelID,
		ChannelCode: t.ChannelCode,
		BrandCode:   t.BrandCode,
		BrandName:   t.BrandName,
		ModelCode:   t.ModelCode,
		ModelName:   t.ModelName,
		FullName:    t.FullName(),
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if t.Channel != nil {
		resp.ChannelName = t.Channel.ChannelName
	}

	return resp
}

// ToResponseList 转换为响应列表
func (s *TerminalTypeService) ToResponseList(list []models.TerminalType) []*TerminalTypeResponse {
	result := make([]*TerminalTypeResponse, len(list))
	for i, t := range list {
		result[i] = s.ToResponse(&t)
	}
	return result
}
