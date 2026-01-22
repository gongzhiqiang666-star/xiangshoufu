import request from '@/utils/request'
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
  return request<PosterCategory[]>({
    url: '/api/v1/admin/poster-categories',
    method: 'get',
    params
  })
}

// 创建分类
export function createPosterCategory(data: PosterCategoryCreateRequest) {
  return request<PosterCategory>({
    url: '/api/v1/admin/poster-categories',
    method: 'post',
    data
  })
}

// 更新分类
export function updatePosterCategory(id: number, data: PosterCategoryUpdateRequest) {
  return request<PosterCategory>({
    url: `/api/v1/admin/poster-categories/${id}`,
    method: 'put',
    data
  })
}

// 删除分类
export function deletePosterCategory(id: number) {
  return request({
    url: `/api/v1/admin/poster-categories/${id}`,
    method: 'delete'
  })
}

// ========== 海报管理 ==========

// 获取海报列表
export function getPosterList(params: PosterListRequest) {
  return request<{ data: Poster[]; total: number }>({
    url: '/api/v1/admin/posters',
    method: 'get',
    params
  })
}

// 获取海报详情
export function getPosterDetail(id: number) {
  return request<Poster>({
    url: `/api/v1/admin/posters/${id}`,
    method: 'get'
  })
}

// 创建海报
export function createPoster(data: PosterCreateRequest) {
  return request<Poster>({
    url: '/api/v1/admin/posters',
    method: 'post',
    data
  })
}

// 更新海报
export function updatePoster(id: number, data: PosterUpdateRequest) {
  return request<Poster>({
    url: `/api/v1/admin/posters/${id}`,
    method: 'put',
    data
  })
}

// 删除海报
export function deletePoster(id: number) {
  return request({
    url: `/api/v1/admin/posters/${id}`,
    method: 'delete'
  })
}

// 批量导入海报
export function batchImportPosters(data: PosterBatchImportRequest) {
  return request<{ imported_count: number }>({
    url: '/api/v1/admin/posters/batch-import',
    method: 'post',
    data
  })
}
