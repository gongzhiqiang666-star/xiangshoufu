package models

import (
	"time"
)

// PosterCategory 海报分类模型
type PosterCategory struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null"`  // 分类名称
	SortOrder int       `json:"sort_order" gorm:"default:0"`   // 排序
	Status    int       `json:"status" gorm:"default:1"`       // 1启用 0禁用
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 表名
func (PosterCategory) TableName() string {
	return "poster_categories"
}

// Poster 营销海报模型
type Poster struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title" gorm:"size:100;not null"`       // 标题
	CategoryID    int64     `json:"category_id" gorm:"not null"`          // 分类ID
	ImageURL      string    `json:"image_url" gorm:"size:500;not null"`   // 原图URL
	ThumbnailURL  string    `json:"thumbnail_url" gorm:"size:500"`        // 缩略图URL
	Description   string    `json:"description" gorm:"type:text"`         // 描述
	FileSize      int64     `json:"file_size" gorm:"default:0"`           // 文件大小（字节）
	Width         int       `json:"width" gorm:"default:0"`               // 图片宽度
	Height        int       `json:"height" gorm:"default:0"`              // 图片高度
	SortOrder     int       `json:"sort_order" gorm:"default:0"`          // 排序
	Status        int       `json:"status" gorm:"default:1"`              // 1启用 0禁用
	DownloadCount int64     `json:"download_count" gorm:"default:0"`      // 下载次数
	ShareCount    int64     `json:"share_count" gorm:"default:0"`         // 分享次数
	CreatedBy     int64     `json:"created_by"`                           // 创建人ID
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联
	Category *PosterCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName 表名
func (Poster) TableName() string {
	return "posters"
}

// ========== 海报分类相关请求/响应 ==========

// PosterCategoryListRequest 分类列表请求
type PosterCategoryListRequest struct {
	Status *int `form:"status"`
}

// PosterCategoryCreateRequest 分类创建请求
type PosterCategoryCreateRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	SortOrder int    `json:"sort_order"`
	Status    int    `json:"status"`
}

// PosterCategoryUpdateRequest 分类更新请求
type PosterCategoryUpdateRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	SortOrder int    `json:"sort_order"`
	Status    int    `json:"status"`
}

// PosterCategoryResponse 分类响应
type PosterCategoryResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	SortOrder   int       `json:"sort_order"`
	Status      int       `json:"status"`
	PosterCount int64     `json:"poster_count,omitempty"` // 海报数量
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse 转换为响应
func (c *PosterCategory) ToResponse() *PosterCategoryResponse {
	return &PosterCategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		SortOrder: c.SortOrder,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// AppPosterCategoryResponse APP端分类响应
type AppPosterCategoryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PosterCount int64  `json:"poster_count"`
}

// ========== 海报相关请求/响应 ==========

// PosterListRequest 海报列表请求
type PosterListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	CategoryID *int64 `form:"category_id"`
	Status     *int   `form:"status"`
	Keyword    string `form:"keyword"`
}

// PosterCreateRequest 海报创建请求
type PosterCreateRequest struct {
	Title        string `json:"title" binding:"required,max=100"`
	CategoryID   int64  `json:"category_id" binding:"required"`
	ImageURL     string `json:"image_url" binding:"required,max=500"`
	ThumbnailURL string `json:"thumbnail_url" binding:"omitempty,max=500"`
	Description  string `json:"description"`
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	SortOrder    int    `json:"sort_order"`
	Status       int    `json:"status"`
}

// PosterUpdateRequest 海报更新请求
type PosterUpdateRequest struct {
	Title        string `json:"title" binding:"required,max=100"`
	CategoryID   int64  `json:"category_id" binding:"required"`
	ImageURL     string `json:"image_url" binding:"required,max=500"`
	ThumbnailURL string `json:"thumbnail_url" binding:"omitempty,max=500"`
	Description  string `json:"description"`
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	SortOrder    int    `json:"sort_order"`
	Status       int    `json:"status"`
}

// PosterBatchImportRequest 批量导入请求
type PosterBatchImportRequest struct {
	CategoryID int64               `json:"category_id" binding:"required"`
	Items      []PosterImportItem  `json:"items" binding:"required,min=1"`
}

// PosterImportItem 导入项
type PosterImportItem struct {
	Title        string `json:"title" binding:"required,max=100"`
	ImageURL     string `json:"image_url" binding:"required,max=500"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// PosterResponse 海报响应
type PosterResponse struct {
	ID            int64                   `json:"id"`
	Title         string                  `json:"title"`
	CategoryID    int64                   `json:"category_id"`
	CategoryName  string                  `json:"category_name,omitempty"`
	ImageURL      string                  `json:"image_url"`
	ThumbnailURL  string                  `json:"thumbnail_url"`
	Description   string                  `json:"description"`
	FileSize      int64                   `json:"file_size"`
	Width         int                     `json:"width"`
	Height        int                     `json:"height"`
	SortOrder     int                     `json:"sort_order"`
	Status        int                     `json:"status"`
	DownloadCount int64                   `json:"download_count"`
	ShareCount    int64                   `json:"share_count"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// ToResponse 转换为响应
func (p *Poster) ToResponse() *PosterResponse {
	resp := &PosterResponse{
		ID:            p.ID,
		Title:         p.Title,
		CategoryID:    p.CategoryID,
		ImageURL:      p.ImageURL,
		ThumbnailURL:  p.ThumbnailURL,
		Description:   p.Description,
		FileSize:      p.FileSize,
		Width:         p.Width,
		Height:        p.Height,
		SortOrder:     p.SortOrder,
		Status:        p.Status,
		DownloadCount: p.DownloadCount,
		ShareCount:    p.ShareCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
	if p.Category != nil {
		resp.CategoryName = p.Category.Name
	}
	return resp
}

// AppPosterResponse APP端海报响应
type AppPosterResponse struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	CategoryID   int64  `json:"category_id"`
	ImageURL     string `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// ToAppResponse 转换为APP端响应
func (p *Poster) ToAppResponse() *AppPosterResponse {
	return &AppPosterResponse{
		ID:           p.ID,
		Title:        p.Title,
		CategoryID:   p.CategoryID,
		ImageURL:     p.ImageURL,
		ThumbnailURL: p.ThumbnailURL,
		Description:  p.Description,
		Width:        p.Width,
		Height:       p.Height,
	}
}
