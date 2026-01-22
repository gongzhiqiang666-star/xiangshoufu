package repository

import (
	"xiangshoufu/internal/models"

	"gorm.io/gorm"
)

// UploadedFileRepository 上传文件仓库接口
type UploadedFileRepository interface {
	Create(file *models.UploadedFile) error
	FindByID(id int64) (*models.UploadedFile, error)
	FindByModule(module string, refID int64) ([]*models.UploadedFile, error)
	UpdateRefID(id int64, refID int64) error
	Delete(id int64) error
}

// GormUploadedFileRepository GORM实现的上传文件仓库
type GormUploadedFileRepository struct {
	db *gorm.DB
}

// NewGormUploadedFileRepository 创建上传文件仓库实例
func NewGormUploadedFileRepository(db *gorm.DB) *GormUploadedFileRepository {
	return &GormUploadedFileRepository{db: db}
}

// Create 创建上传记录
func (r *GormUploadedFileRepository) Create(file *models.UploadedFile) error {
	return r.db.Create(file).Error
}

// FindByID 根据ID查找上传记录
func (r *GormUploadedFileRepository) FindByID(id int64) (*models.UploadedFile, error) {
	var file models.UploadedFile
	err := r.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByModule 根据模块和关联ID查找上传记录
func (r *GormUploadedFileRepository) FindByModule(module string, refID int64) ([]*models.UploadedFile, error) {
	var files []*models.UploadedFile
	err := r.db.Where("module = ? AND ref_id = ?", module, refID).Find(&files).Error
	return files, err
}

// UpdateRefID 更新关联ID
func (r *GormUploadedFileRepository) UpdateRefID(id int64, refID int64) error {
	return r.db.Model(&models.UploadedFile{}).
		Where("id = ?", id).
		Update("ref_id", refID).Error
}

// Delete 删除上传记录
func (r *GormUploadedFileRepository) Delete(id int64) error {
	return r.db.Delete(&models.UploadedFile{}, id).Error
}
