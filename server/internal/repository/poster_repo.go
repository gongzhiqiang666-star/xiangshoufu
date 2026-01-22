package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// PosterRepository 海报仓库接口
type PosterRepository interface {
	Create(poster *models.Poster) error
	CreateBatch(posters []*models.Poster) error
	Update(poster *models.Poster) error
	Delete(id int64) error
	FindByID(id int64) (*models.Poster, error)
	FindAll(req *models.PosterListRequest) ([]*models.Poster, int64, error)
	FindActive(categoryID *int64, page, pageSize int) ([]*models.Poster, int64, error)
	IncrementDownloadCount(id int64) error
	IncrementShareCount(id int64) error
	CountByCategory(categoryID int64) (int64, error)
}

// GormPosterRepository GORM实现的海报仓库
type GormPosterRepository struct {
	db *gorm.DB
}

// NewGormPosterRepository 创建海报仓库实例
func NewGormPosterRepository(db *gorm.DB) *GormPosterRepository {
	return &GormPosterRepository{db: db}
}

// Create 创建海报
func (r *GormPosterRepository) Create(poster *models.Poster) error {
	return r.db.Create(poster).Error
}

// CreateBatch 批量创建海报
func (r *GormPosterRepository) CreateBatch(posters []*models.Poster) error {
	return r.db.CreateInBatches(posters, 100).Error
}

// Update 更新海报
func (r *GormPosterRepository) Update(poster *models.Poster) error {
	return r.db.Save(poster).Error
}

// Delete 删除海报
func (r *GormPosterRepository) Delete(id int64) error {
	return r.db.Delete(&models.Poster{}, id).Error
}

// FindByID 根据ID查找海报
func (r *GormPosterRepository) FindByID(id int64) (*models.Poster, error) {
	var poster models.Poster
	err := r.db.Preload("Category").First(&poster, id).Error
	if err != nil {
		return nil, err
	}
	return &poster, nil
}

// FindAll 查询海报列表（管理端）
func (r *GormPosterRepository) FindAll(req *models.PosterListRequest) ([]*models.Poster, int64, error) {
	var posters []*models.Poster
	var total int64

	query := r.db.Model(&models.Poster{})

	// 分类筛选
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 排序并查询
	err := query.Preload("Category").
		Order("sort_order DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posters).Error

	return posters, total, err
}

// FindActive 查询有效的海报列表（APP端）
func (r *GormPosterRepository) FindActive(categoryID *int64, page, pageSize int) ([]*models.Poster, int64, error) {
	var posters []*models.Poster
	var total int64

	query := r.db.Model(&models.Poster{}).Where("status = ?", 1)

	// 分类筛选
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 排序并查询
	err := query.Order("sort_order DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posters).Error

	return posters, total, err
}

// IncrementDownloadCount 增加下载次数
func (r *GormPosterRepository) IncrementDownloadCount(id int64) error {
	return r.db.Model(&models.Poster{}).
		Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

// IncrementShareCount 增加分享次数
func (r *GormPosterRepository) IncrementShareCount(id int64) error {
	return r.db.Model(&models.Poster{}).
		Where("id = ?", id).
		UpdateColumn("share_count", gorm.Expr("share_count + 1")).Error
}

// CountByCategory 统计分类下的海报数量
func (r *GormPosterRepository) CountByCategory(categoryID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Poster{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
