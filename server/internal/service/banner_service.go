package service

import (
	"fmt"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// BannerService Banner服务接口
type BannerService interface {
	// 管理端接口
	Create(req *models.BannerCreateRequest, createdBy int64) (*models.Banner, error)
	Update(id int64, req *models.BannerUpdateRequest) (*models.Banner, error)
	Delete(id int64) error
	GetByID(id int64) (*models.Banner, error)
	GetList(req *models.BannerListRequest) ([]*models.Banner, int64, error)
	UpdateStatus(id int64, status int) error
	UpdateSortOrder(req *models.BannerSortRequest) error

	// APP端接口
	GetActiveBanners() ([]*models.Banner, error)
	RecordClick(id int64) error
}

// bannerService Banner服务实现
type bannerService struct {
	repo repository.BannerRepository
}

// NewBannerService 创建Banner服务实例
func NewBannerService(repo repository.BannerRepository) BannerService {
	return &bannerService{repo: repo}
}

// Create 创建Banner
func (s *bannerService) Create(req *models.BannerCreateRequest, createdBy int64) (*models.Banner, error) {
	banner := &models.Banner{
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		LinkType:  req.LinkType,
		LinkURL:   req.LinkURL,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(banner); err != nil {
		return nil, fmt.Errorf("创建Banner失败: %w", err)
	}

	return banner, nil
}

// Update 更新Banner
func (s *bannerService) Update(id int64, req *models.BannerUpdateRequest) (*models.Banner, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("Banner不存在: %w", err)
	}

	banner.Title = req.Title
	banner.ImageURL = req.ImageURL
	banner.LinkType = req.LinkType
	banner.LinkURL = req.LinkURL
	banner.SortOrder = req.SortOrder
	banner.Status = req.Status
	banner.StartTime = req.StartTime
	banner.EndTime = req.EndTime

	if err := s.repo.Update(banner); err != nil {
		return nil, fmt.Errorf("更新Banner失败: %w", err)
	}

	return banner, nil
}

// Delete 删除Banner
func (s *bannerService) Delete(id int64) error {
	// 检查是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("Banner不存在: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除Banner失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取Banner
func (s *bannerService) GetByID(id int64) (*models.Banner, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("Banner不存在: %w", err)
	}
	return banner, nil
}

// GetList 获取Banner列表
func (s *bannerService) GetList(req *models.BannerListRequest) ([]*models.Banner, int64, error) {
	return s.repo.FindAll(req)
}

// UpdateStatus 更新状态
func (s *bannerService) UpdateStatus(id int64, status int) error {
	// 检查是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("Banner不存在: %w", err)
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	return nil
}

// UpdateSortOrder 批量更新排序
func (s *bannerService) UpdateSortOrder(req *models.BannerSortRequest) error {
	if err := s.repo.UpdateSortOrder(req.Items); err != nil {
		return fmt.Errorf("更新排序失败: %w", err)
	}
	return nil
}

// GetActiveBanners 获取有效的Banner列表（APP端）
func (s *bannerService) GetActiveBanners() ([]*models.Banner, error) {
	return s.repo.FindActive()
}

// RecordClick 记录点击
func (s *bannerService) RecordClick(id int64) error {
	return s.repo.IncrementClickCount(id)
}
