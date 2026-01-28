package repository

import (
	"context"
	"errors"

	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// TerminalTypeRepository 终端类型仓库接口
type TerminalTypeRepository interface {
	Create(ctx context.Context, terminalType *models.TerminalType) error
	Update(ctx context.Context, terminalType *models.TerminalType) error
	FindByID(ctx context.Context, id int64) (*models.TerminalType, error)
	FindByChannelAndCodes(ctx context.Context, channelID int64, brandCode, modelCode string) (*models.TerminalType, error)
	List(ctx context.Context, filter TerminalTypeFilter) ([]models.TerminalType, int64, error)
	ListByChannelID(ctx context.Context, channelID int64, onlyEnabled bool) ([]models.TerminalType, error)
	UpdateStatus(ctx context.Context, id int64, status int16) error
	Delete(ctx context.Context, id int64) error
}

// TerminalTypeFilter 终端类型筛选条件
type TerminalTypeFilter struct {
	ChannelID   int64
	ChannelCode string
	BrandCode   string
	ModelCode   string
	Status      *int16
	Keyword     string // 搜索关键词（品牌名称、型号名称）
	Page        int
	PageSize    int
}

// GormTerminalTypeRepository GORM实现
type GormTerminalTypeRepository struct {
	db *gorm.DB
}

// NewGormTerminalTypeRepository 创建仓库实例
func NewGormTerminalTypeRepository(db *gorm.DB) *GormTerminalTypeRepository {
	return &GormTerminalTypeRepository{db: db}
}

// Create 创建终端类型
func (r *GormTerminalTypeRepository) Create(ctx context.Context, terminalType *models.TerminalType) error {
	return r.db.WithContext(ctx).Create(terminalType).Error
}

// Update 更新终端类型
func (r *GormTerminalTypeRepository) Update(ctx context.Context, terminalType *models.TerminalType) error {
	return r.db.WithContext(ctx).Save(terminalType).Error
}

// FindByID 根据ID查找
func (r *GormTerminalTypeRepository) FindByID(ctx context.Context, id int64) (*models.TerminalType, error) {
	var terminalType models.TerminalType
	err := r.db.WithContext(ctx).Preload("Channel").First(&terminalType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &terminalType, nil
}

// FindByChannelAndCodes 根据通道和编码查找
func (r *GormTerminalTypeRepository) FindByChannelAndCodes(ctx context.Context, channelID int64, brandCode, modelCode string) (*models.TerminalType, error) {
	var terminalType models.TerminalType
	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND brand_code = ? AND model_code = ?", channelID, brandCode, modelCode).
		First(&terminalType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &terminalType, nil
}

// List 列表查询
func (r *GormTerminalTypeRepository) List(ctx context.Context, filter TerminalTypeFilter) ([]models.TerminalType, int64, error) {
	var terminalTypes []models.TerminalType
	var total int64

	query := r.db.WithContext(ctx).Model(&models.TerminalType{})

	// 筛选条件
	if filter.ChannelID > 0 {
		query = query.Where("channel_id = ?", filter.ChannelID)
	}
	if filter.ChannelCode != "" {
		query = query.Where("channel_code = ?", filter.ChannelCode)
	}
	if filter.BrandCode != "" {
		query = query.Where("brand_code = ?", filter.BrandCode)
	}
	if filter.ModelCode != "" {
		query = query.Where("model_code = ?", filter.ModelCode)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("brand_name ILIKE ? OR model_name ILIKE ? OR model_code ILIKE ?", keyword, keyword, keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 查询
	err := query.Preload("Channel").Order("id DESC").Find(&terminalTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return terminalTypes, total, nil
}

// ListByChannelID 根据通道ID列出终端类型
func (r *GormTerminalTypeRepository) ListByChannelID(ctx context.Context, channelID int64, onlyEnabled bool) ([]models.TerminalType, error) {
	var terminalTypes []models.TerminalType

	query := r.db.WithContext(ctx).Where("channel_id = ?", channelID)
	if onlyEnabled {
		query = query.Where("status = ?", models.TerminalTypeStatusEnabled)
	}

	err := query.Order("brand_name ASC, model_code ASC").Find(&terminalTypes).Error
	if err != nil {
		return nil, err
	}

	return terminalTypes, nil
}

// UpdateStatus 更新状态
func (r *GormTerminalTypeRepository) UpdateStatus(ctx context.Context, id int64, status int16) error {
	return r.db.WithContext(ctx).Model(&models.TerminalType{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除终端类型
func (r *GormTerminalTypeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.TerminalType{}, id).Error
}
