package models

import (
	"time"
)

// UploadedFile 上传文件记录模型
type UploadedFile struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	OriginalName string    `json:"original_name" gorm:"size:255;not null"` // 原始文件名
	StoredName   string    `json:"stored_name" gorm:"size:255;not null"`   // 存储文件名
	FilePath     string    `json:"file_path" gorm:"size:500;not null"`     // 文件路径
	FileURL      string    `json:"file_url" gorm:"size:500;not null"`      // 访问URL
	FileSize     int64     `json:"file_size" gorm:"default:0"`             // 文件大小（字节）
	MimeType     string    `json:"mime_type" gorm:"size:100"`              // MIME类型
	Width        int       `json:"width"`                                  // 图片宽度
	Height       int       `json:"height"`                                 // 图片高度
	Module       string    `json:"module" gorm:"size:50"`                  // 所属模块 banner/poster
	RefID        int64     `json:"ref_id"`                                 // 关联ID
	UploadedBy   int64     `json:"uploaded_by"`                            // 上传人ID
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 表名
func (UploadedFile) TableName() string {
	return "uploaded_files"
}

// UploadImageRequest 图片上传请求
type UploadImageRequest struct {
	Module string `form:"module" binding:"omitempty,oneof=banner poster"` // 模块
}

// UploadImageResponse 图片上传响应
type UploadImageResponse struct {
	ID           int64  `json:"id"`
	OriginalName string `json:"original_name"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"` // 缩略图URL（仅海报）
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	MimeType     string `json:"mime_type"`
}

// ToResponse 转换为响应
func (f *UploadedFile) ToResponse() *UploadImageResponse {
	return &UploadImageResponse{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		FileURL:      f.FileURL,
		FileSize:     f.FileSize,
		Width:        f.Width,
		Height:       f.Height,
		MimeType:     f.MimeType,
	}
}

// ImageConfig 图片配置
type ImageConfig struct {
	MaxSize       int64    // 最大文件大小（字节）
	MaxWidth      int      // 最大宽度
	MaxHeight     int      // 最大高度
	Quality       int      // 压缩质量 1-100
	AllowedTypes  []string // 允许的MIME类型
	GenerateThumb bool     // 是否生成缩略图
	ThumbWidth    int      // 缩略图宽度
}

// BannerImageConfig Banner图片配置
var BannerImageConfig = ImageConfig{
	MaxSize:       2 * 1024 * 1024, // 2MB
	MaxWidth:      1500,
	MaxHeight:     800,
	Quality:       85,
	AllowedTypes:  []string{"image/jpeg", "image/png", "image/webp"},
	GenerateThumb: false,
}

// PosterImageConfig 海报图片配置
var PosterImageConfig = ImageConfig{
	MaxSize:       5 * 1024 * 1024, // 5MB
	MaxWidth:      2000,
	MaxHeight:     3000,
	Quality:       85,
	AllowedTypes:  []string{"image/jpeg", "image/png", "image/webp"},
	GenerateThumb: true,
	ThumbWidth:    300,
}
