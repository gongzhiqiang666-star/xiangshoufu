package models

import (
	"time"
)

// LinkType 链接类型
type LinkType int

const (
	LinkTypeNone     LinkType = 0 // 无链接
	LinkTypeInternal LinkType = 1 // 内部页面
	LinkTypeExternal LinkType = 2 // 外部链接
)

// Banner 滚动图模型
type Banner struct {
	ID         int64      `json:"id" gorm:"primaryKey"`
	Title      string     `json:"title" gorm:"size:100;not null"`           // 标题
	ImageURL   string     `json:"image_url" gorm:"size:500;not null"`       // 图片URL
	LinkType   LinkType   `json:"link_type" gorm:"default:0"`               // 链接类型
	LinkURL    string     `json:"link_url" gorm:"size:500"`                 // 跳转链接
	SortOrder  int        `json:"sort_order" gorm:"default:0"`              // 排序（越大越靠前）
	Status     int        `json:"status" gorm:"default:1"`                  // 1启用 0禁用
	StartTime  *time.Time `json:"start_time"`                               // 开始展示时间
	EndTime    *time.Time `json:"end_time"`                                 // 结束展示时间
	ClickCount int64      `json:"click_count" gorm:"default:0"`             // 点击统计
	CreatedBy  int64      `json:"created_by"`                               // 创建人ID
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 表名
func (Banner) TableName() string {
	return "banners"
}

// IsActive 判断Banner是否有效
func (b *Banner) IsActive() bool {
	if b.Status != 1 {
		return false
	}
	now := time.Now()
	if b.StartTime != nil && now.Before(*b.StartTime) {
		return false
	}
	if b.EndTime != nil && now.After(*b.EndTime) {
		return false
	}
	return true
}

// BannerListRequest Banner列表请求
type BannerListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   *int `form:"status"`
}

// BannerCreateRequest Banner创建请求
type BannerCreateRequest struct {
	Title     string     `json:"title" binding:"required,max=100"`
	ImageURL  string     `json:"image_url" binding:"required,max=500"`
	LinkType  LinkType   `json:"link_type"`
	LinkURL   string     `json:"link_url" binding:"omitempty,max=500"`
	SortOrder int        `json:"sort_order"`
	Status    int        `json:"status"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// BannerUpdateRequest Banner更新请求
type BannerUpdateRequest struct {
	Title     string     `json:"title" binding:"required,max=100"`
	ImageURL  string     `json:"image_url" binding:"required,max=500"`
	LinkType  LinkType   `json:"link_type"`
	LinkURL   string     `json:"link_url" binding:"omitempty,max=500"`
	SortOrder int        `json:"sort_order"`
	Status    int        `json:"status"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// BannerStatusRequest 状态切换请求
type BannerStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1"`
}

// BannerSortRequest 批量排序请求
type BannerSortRequest struct {
	Items []BannerSortItem `json:"items" binding:"required,min=1"`
}

// BannerSortItem 排序项
type BannerSortItem struct {
	ID        int64 `json:"id" binding:"required"`
	SortOrder int   `json:"sort_order"`
}

// BannerResponse Banner响应
type BannerResponse struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	ImageURL   string     `json:"image_url"`
	LinkType   LinkType   `json:"link_type"`
	LinkURL    string     `json:"link_url"`
	SortOrder  int        `json:"sort_order"`
	Status     int        `json:"status"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	ClickCount int64      `json:"click_count"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应
func (b *Banner) ToResponse() *BannerResponse {
	return &BannerResponse{
		ID:         b.ID,
		Title:      b.Title,
		ImageURL:   b.ImageURL,
		LinkType:   b.LinkType,
		LinkURL:    b.LinkURL,
		SortOrder:  b.SortOrder,
		Status:     b.Status,
		StartTime:  b.StartTime,
		EndTime:    b.EndTime,
		ClickCount: b.ClickCount,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}

// AppBannerResponse APP端Banner响应（精简版）
type AppBannerResponse struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	ImageURL string   `json:"image_url"`
	LinkType LinkType `json:"link_type"`
	LinkURL  string   `json:"link_url"`
}

// ToAppResponse 转换为APP端响应
func (b *Banner) ToAppResponse() *AppBannerResponse {
	return &AppBannerResponse{
		ID:       b.ID,
		Title:    b.Title,
		ImageURL: b.ImageURL,
		LinkType: b.LinkType,
		LinkURL:  b.LinkURL,
	}
}
