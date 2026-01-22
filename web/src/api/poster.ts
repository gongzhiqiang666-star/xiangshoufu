import { get, post, put, del } from './request'
import type {
  Poster,
  PosterCategory,
  PosterCreateRequest,
  PosterUpdateRequest,
  PosterListRequest,
  PosterCategoryCreateRequest,
  PosterCategoryUpdateRequest,
  PosterCategoryListRequest,
  PosterBatchImportRequest
} from '@/types/poster'

// ========== 分类管理 ==========

// 获取分类列表
export function getPosterCategories(params?: PosterCategoryListRequest) {
  return get<PosterCategory[]>('/api/v1/admin/poster-categories', params)
}

// 创建分类
export function createPosterCategory(data: PosterCategoryCreateRequest) {
  return post<PosterCategory>('/api/v1/admin/poster-categories', data)
}

// 更新分类
export function updatePosterCategory(id: number, data: PosterCategoryUpdateRequest) {
  return put<PosterCategory>(`/api/v1/admin/poster-categories/${id}`, data)
}

// 删除分类
export function deletePosterCategory(id: number) {
  return del(`/api/v1/admin/poster-categories/${id}`)
}

// ========== 海报管理 ==========

// 获取海报列表
export function getPosterList(params: PosterListRequest) {
  return get<{ data: Poster[]; total: number }>('/api/v1/admin/posters', params)
}

// 获取海报详情
export function getPosterDetail(id: number) {
  return get<Poster>(`/api/v1/admin/posters/${id}`)
}

// 创建海报
export function createPoster(data: PosterCreateRequest) {
  return post<Poster>('/api/v1/admin/posters', data)
}

// 更新海报
export function updatePoster(id: number, data: PosterUpdateRequest) {
  return put<Poster>(`/api/v1/admin/posters/${id}`, data)
}

// 删除海报
export function deletePoster(id: number) {
  return del(`/api/v1/admin/posters/${id}`)
}

// 批量导入海报
export function batchImportPosters(data: PosterBatchImportRequest) {
  return post<{ imported_count: number }>('/api/v1/admin/posters/batch-import', data)
}
