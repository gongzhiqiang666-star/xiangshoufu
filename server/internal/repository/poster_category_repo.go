package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// PosterCategoryRepository 海报分类仓库接口
type PosterCategoryRepository interface {
	Create(category *models.PosterCategory) error
	Update(category *models.PosterCategory) error
	Delete(id int64) error
	FindByID(id int64) (*models.PosterCategory, error)
	FindAll(req *models.PosterCategoryListRequest) ([]*models.PosterCategory, error)
	FindActive() ([]*models.PosterCategory, error)
	GetPosterCounts() (map[int64]int64, error)
	ExistsByID(id int64) (bool, error)
}

// GormPosterCategoryRepository GORM实现的海报分类仓库
type GormPosterCategoryRepository struct {
	db *gorm.DB
}

// NewGormPosterCategoryRepository 创建海报分类仓库实例
func NewGormPosterCategoryRepository(db *gorm.DB) *GormPosterCategoryRepository {
	return &GormPosterCategoryRepository{db: db}
}

// Create 创建分类
func (r *GormPosterCategoryRepository) Create(category *models.PosterCategory) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *GormPosterCategoryRepository) Update(category *models.PosterCategory) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *GormPosterCategoryRepository) Delete(id int64) error {
	return r.db.Delete(&models.PosterCategory{}, id).Error
}

// FindByID 根据ID查找分类
func (r *GormPosterCategoryRepository) FindByID(id int64) (*models.PosterCategory, error) {
	var category models.PosterCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAll 查询分类列表（管理端）
func (r *GormPosterCategoryRepository) FindAll(req *models.PosterCategoryListRequest) ([]*models.PosterCategory, error) {
	var categories []*models.PosterCategory

	query := r.db.Model(&models.PosterCategory{})

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	err := query.Order("sort_order DESC, id DESC").Find(&categories).Error
	return categories, err
}

// FindActive 查询有效的分类列表（APP端）
func (r *GormPosterCategoryRepository) FindActive() ([]*models.PosterCategory, error) {
	var categories []*models.PosterCategory

	err := r.db.Where("status = ?", 1).
		Order("sort_order DESC, id DESC").
		Find(&categories).Error

	return categories, err
}

// GetPosterCounts 获取各分类的海报数量
func (r *GormPosterCategoryRepository) GetPosterCounts() (map[int64]int64, error) {
	type Result struct {
		CategoryID int64
		Count      int64
	}
	var results []Result

	err := r.db.Model(&models.Poster{}).
		Select("category_id, COUNT(*) as count").
		Where("status = ?", 1).
		Group("category_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[int64]int64)
	for _, r := range results {
		counts[r.CategoryID] = r.Count
	}
	return counts, nil
}

// ExistsByID 检查分类是否存在
func (r *GormPosterCategoryRepository) ExistsByID(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.PosterCategory{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
