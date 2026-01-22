// Banner 滚动图类型定义

// 链接类型
export enum LinkType {
  None = 0,      // 无链接
  Internal = 1,  // 内部页面
  External = 2   // 外部链接
}

// Banner实体
export interface Banner {
  id: number
  title: string
  image_url: string
  link_type: LinkType
  link_url: string
  sort_order: number
  status: number
  start_time: string | null
  end_time: string | null
  click_count: number
  created_at: string
  updated_at: string
}

// Banner创建请求
export interface BannerCreateRequest {
  title: string
  image_url: string
  link_type: LinkType
  link_url?: string
  sort_order?: number
  status: number
  start_time?: string | null
  end_time?: string | null
}

// Banner更新请求
export interface BannerUpdateRequest {
  title: string
  image_url: string
  link_type: LinkType
  link_url?: string
  sort_order?: number
  status: number
  start_time?: string | null
  end_time?: string | null
}

// Banner列表请求
export interface BannerListRequest {
  page?: number
  page_size?: number
  status?: number
}

// Banner状态更新请求
export interface BannerStatusRequest {
  status: number
}

// Banner排序项
export interface BannerSortItem {
  id: number
  sort_order: number
}

// Banner排序请求
export interface BannerSortRequest {
  items: BannerSortItem[]
}

// 链接类型选项
export const linkTypeOptions = [
  { value: LinkType.None, label: '无链接' },
  { value: LinkType.Internal, label: '内部页面' },
  { value: LinkType.External, label: '外部链接' }
]

// 状态选项
export const statusOptions = [
  { value: 1, label: '启用' },
  { value: 0, label: '禁用' }
]
