// 海报相关类型定义

// 海报分类
export interface PosterCategory {
  id: number
  name: string
  sort_order: number
  status: number
  poster_count?: number
  created_at: string
  updated_at: string
}

// 海报分类创建请求
export interface PosterCategoryCreateRequest {
  name: string
  sort_order?: number
  status: number
}

// 海报分类更新请求
export interface PosterCategoryUpdateRequest {
  name: string
  sort_order?: number
  status: number
}

// 海报实体
export interface Poster {
  id: number
  title: string
  category_id: number
  category_name?: string
  image_url: string
  thumbnail_url: string
  description: string
  file_size: number
  width: number
  height: number
  sort_order: number
  status: number
  download_count: number
  share_count: number
  created_at: string
  updated_at: string
}

// 海报创建请求
export interface PosterCreateRequest {
  title: string
  category_id: number
  image_url: string
  thumbnail_url?: string
  description?: string
  file_size?: number
  width?: number
  height?: number
  sort_order?: number
  status: number
}

// 海报更新请求
export interface PosterUpdateRequest {
  title: string
  category_id: number
  image_url: string
  thumbnail_url?: string
  description?: string
  file_size?: number
  width?: number
  height?: number
  sort_order?: number
  status: number
}

// 海报列表请求
export interface PosterListRequest {
  page?: number
  page_size?: number
  category_id?: number
  status?: number
  keyword?: string
}

// 海报导入项
export interface PosterImportItem {
  title: string
  image_url: string
  thumbnail_url?: string
  description?: string
  file_size?: number
  width?: number
  height?: number
}

// 批量导入请求
export interface PosterBatchImportRequest {
  category_id: number
  items: PosterImportItem[]
}

// 分类列表请求
export interface PosterCategoryListRequest {
  status?: number
}

// 状态选项
export const posterStatusOptions = [
  { value: 1, label: '启用' },
  { value: 0, label: '禁用' }
]
