package service

import (
	"fmt"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// PosterService 海报服务接口
type PosterService interface {
	// 分类管理
	CreateCategory(req *models.PosterCategoryCreateRequest) (*models.PosterCategory, error)
	UpdateCategory(id int64, req *models.PosterCategoryUpdateRequest) (*models.PosterCategory, error)
	DeleteCategory(id int64) error
	GetCategoryByID(id int64) (*models.PosterCategory, error)
	GetCategories(req *models.PosterCategoryListRequest) ([]*models.PosterCategory, error)

	// 海报管理
	Create(req *models.PosterCreateRequest, createdBy int64) (*models.Poster, error)
	Update(id int64, req *models.PosterUpdateRequest) (*models.Poster, error)
	Delete(id int64) error
	GetByID(id int64) (*models.Poster, error)
	GetList(req *models.PosterListRequest) ([]*models.Poster, int64, error)
	BatchImport(req *models.PosterBatchImportRequest, createdBy int64) (int, error)

	// APP端接口
	GetActiveCategories() ([]*models.AppPosterCategoryResponse, error)
	GetActivePosters(categoryID *int64, page, pageSize int) ([]*models.Poster, int64, error)
	RecordDownload(id int64) error
	RecordShare(id int64) error
}

// posterService 海报服务实现
type posterService struct {
	posterRepo   repository.PosterRepository
	categoryRepo repository.PosterCategoryRepository
}

// NewPosterService 创建海报服务实例
func NewPosterService(posterRepo repository.PosterRepository, categoryRepo repository.PosterCategoryRepository) PosterService {
	return &posterService{
		posterRepo:   posterRepo,
		categoryRepo: categoryRepo,
	}
}

// ========== 分类管理 ==========

// CreateCategory 创建分类
func (s *posterService) CreateCategory(req *models.PosterCategoryCreateRequest) (*models.PosterCategory, error) {
	category := &models.PosterCategory{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}

	return category, nil
}

// UpdateCategory 更新分类
func (s *posterService) UpdateCategory(id int64, req *models.PosterCategoryUpdateRequest) (*models.PosterCategory, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("分类不存在: %w", err)
	}

	category.Name = req.Name
	category.SortOrder = req.SortOrder
	category.Status = req.Status

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("更新分类失败: %w", err)
	}

	return category, nil
}

// DeleteCategory 删除分类
func (s *posterService) DeleteCategory(id int64) error {
	// 检查是否存在
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("分类不存在: %w", err)
	}

	// 检查分类下是否有海报
	count, err := s.posterRepo.CountByCategory(id)
	if err != nil {
		return fmt.Errorf("检查分类海报失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("该分类下还有%d张海报，请先删除或移动海报", count)
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}

	return nil
}

// GetCategoryByID 根据ID获取分类
func (s *posterService) GetCategoryByID(id int64) (*models.PosterCategory, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("分类不存在: %w", err)
	}
	return category, nil
}

// GetCategories 获取分类列表
func (s *posterService) GetCategories(req *models.PosterCategoryListRequest) ([]*models.PosterCategory, error) {
	return s.categoryRepo.FindAll(req)
}

// ========== 海报管理 ==========

// Create 创建海报
func (s *posterService) Create(req *models.PosterCreateRequest, createdBy int64) (*models.Poster, error) {
	// 检查分类是否存在
	exists, err := s.categoryRepo.ExistsByID(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("检查分类失败: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("分类不存在")
	}

	poster := &models.Poster{
		Title:        req.Title,
		CategoryID:   req.CategoryID,
		ImageURL:     req.ImageURL,
		ThumbnailURL: req.ThumbnailURL,
		Description:  req.Description,
		FileSize:     req.FileSize,
		Width:        req.Width,
		Height:       req.Height,
		SortOrder:    req.SortOrder,
		Status:       req.Status,
		CreatedBy:    createdBy,
	}

	if err := s.posterRepo.Create(poster); err != nil {
		return nil, fmt.Errorf("创建海报失败: %w", err)
	}

	return poster, nil
}

// Update 更新海报
func (s *posterService) Update(id int64, req *models.PosterUpdateRequest) (*models.Poster, error) {
	poster, err := s.posterRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("海报不存在: %w", err)
	}

	// 检查分类是否存在
	if req.CategoryID != poster.CategoryID {
		exists, err := s.categoryRepo.ExistsByID(req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("检查分类失败: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("分类不存在")
		}
	}

	poster.Title = req.Title
	poster.CategoryID = req.CategoryID
	poster.ImageURL = req.ImageURL
	poster.ThumbnailURL = req.ThumbnailURL
	poster.Description = req.Description
	poster.FileSize = req.FileSize
	poster.Width = req.Width
	poster.Height = req.Height
	poster.SortOrder = req.SortOrder
	poster.Status = req.Status

	if err := s.posterRepo.Update(poster); err != nil {
		return nil, fmt.Errorf("更新海报失败: %w", err)
	}

	return poster, nil
}

// Delete 删除海报
func (s *posterService) Delete(id int64) error {
	// 检查是否存在
	_, err := s.posterRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("海报不存在: %w", err)
	}

	if err := s.posterRepo.Delete(id); err != nil {
		return fmt.Errorf("删除海报失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取海报
func (s *posterService) GetByID(id int64) (*models.Poster, error) {
	poster, err := s.posterRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("海报不存在: %w", err)
	}
	return poster, nil
}

// GetList 获取海报列表
func (s *posterService) GetList(req *models.PosterListRequest) ([]*models.Poster, int64, error) {
	return s.posterRepo.FindAll(req)
}

// BatchImport 批量导入海报
func (s *posterService) BatchImport(req *models.PosterBatchImportRequest, createdBy int64) (int, error) {
	// 检查分类是否存在
	exists, err := s.categoryRepo.ExistsByID(req.CategoryID)
	if err != nil {
		return 0, fmt.Errorf("检查分类失败: %w", err)
	}
	if !exists {
		return 0, fmt.Errorf("分类不存在")
	}

	posters := make([]*models.Poster, 0, len(req.Items))
	for _, item := range req.Items {
		poster := &models.Poster{
			Title:        item.Title,
			CategoryID:   req.CategoryID,
			ImageURL:     item.ImageURL,
			ThumbnailURL: item.ThumbnailURL,
			Description:  item.Description,
			FileSize:     item.FileSize,
			Width:        item.Width,
			Height:       item.Height,
			Status:       1,
			CreatedBy:    createdBy,
		}
		posters = append(posters, poster)
	}

	if err := s.posterRepo.CreateBatch(posters); err != nil {
		return 0, fmt.Errorf("批量导入失败: %w", err)
	}

	return len(posters), nil
}

// ========== APP端接口 ==========

// GetActiveCategories 获取有效的分类列表（APP端，含海报数量）
func (s *posterService) GetActiveCategories() ([]*models.AppPosterCategoryResponse, error) {
	categories, err := s.categoryRepo.FindActive()
	if err != nil {
		return nil, err
	}

	// 获取各分类的海报数量
	counts, err := s.categoryRepo.GetPosterCounts()
	if err != nil {
		return nil, err
	}

	result := make([]*models.AppPosterCategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, &models.AppPosterCategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			PosterCount: counts[c.ID],
		})
	}

	return result, nil
}

// GetActivePosters 获取有效的海报列表（APP端）
func (s *posterService) GetActivePosters(categoryID *int64, page, pageSize int) ([]*models.Poster, int64, error) {
	return s.posterRepo.FindActive(categoryID, page, pageSize)
}

// RecordDownload 记录下载
func (s *posterService) RecordDownload(id int64) error {
	return s.posterRepo.IncrementDownloadCount(id)
}

// RecordShare 记录分享
func (s *posterService) RecordShare(id int64) error {
	return s.posterRepo.IncrementShareCount(id)
}
