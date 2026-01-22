package repository

import (
	"time"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// BannerRepository Banner仓库接口
type BannerRepository interface {
	Create(banner *models.Banner) error
	Update(banner *models.Banner) error
	Delete(id int64) error
	FindByID(id int64) (*models.Banner, error)
	FindAll(req *models.BannerListRequest) ([]*models.Banner, int64, error)
	FindActive() ([]*models.Banner, error)
	UpdateStatus(id int64, status int) error
	UpdateSortOrder(items []models.BannerSortItem) error
	IncrementClickCount(id int64) error
}

// GormBannerRepository GORM实现的Banner仓库
type GormBannerRepository struct {
	db *gorm.DB
}

// NewGormBannerRepository 创建Banner仓库实例
func NewGormBannerRepository(db *gorm.DB) *GormBannerRepository {
	return &GormBannerRepository{db: db}
}

// Create 创建Banner
func (r *GormBannerRepository) Create(banner *models.Banner) error {
	return r.db.Create(banner).Error
}

// Update 更新Banner
func (r *GormBannerRepository) Update(banner *models.Banner) error {
	return r.db.Save(banner).Error
}

// Delete 删除Banner
func (r *GormBannerRepository) Delete(id int64) error {
	return r.db.Delete(&models.Banner{}, id).Error
}

// FindByID 根据ID查找Banner
func (r *GormBannerRepository) FindByID(id int64) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// FindAll 查询Banner列表（管理端）
func (r *GormBannerRepository) FindAll(req *models.BannerListRequest) ([]*models.Banner, int64, error) {
	var banners []*models.Banner
	var total int64

	query := r.db.Model(&models.Banner{})

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
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
	err := query.Order("sort_order DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&banners).Error

	return banners, total, err
}

// FindActive 查询有效的Banner列表（APP端）
func (r *GormBannerRepository) FindActive() ([]*models.Banner, error) {
	var banners []*models.Banner
	now := time.Now()

	err := r.db.Where("status = ?", 1).
		Where("(start_time IS NULL OR start_time <= ?)", now).
		Where("(end_time IS NULL OR end_time >= ?)", now).
		Order("sort_order DESC, id DESC").
		Find(&banners).Error

	return banners, err
}

// UpdateStatus 更新状态
func (r *GormBannerRepository) UpdateStatus(id int64, status int) error {
	return r.db.Model(&models.Banner{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateSortOrder 批量更新排序
func (r *GormBannerRepository) UpdateSortOrder(items []models.BannerSortItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.Banner{}).
				Where("id = ?", item.ID).
				Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// IncrementClickCount 增加点击次数
func (r *GormBannerRepository) IncrementClickCount(id int64) error {
	return r.db.Model(&models.Banner{}).
		Where("id = ?", id).
		UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error
}
