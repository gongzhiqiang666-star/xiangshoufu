import { get, post, put, del } from './request'
import type {
  Banner,
  BannerCreateRequest,
  BannerUpdateRequest,
  BannerListRequest,
  BannerStatusRequest,
  BannerSortRequest
} from '@/types/banner'

// 获取Banner列表（管理端）
export function getBannerList(params: BannerListRequest) {
  return get<{ data: Banner[]; total: number }>('/v1/admin/banners', params)
}

// 获取Banner详情
export function getBannerDetail(id: number) {
  return get<Banner>(`/v1/admin/banners/${id}`)
}

// 创建Banner
export function createBanner(data: BannerCreateRequest) {
  return post<Banner>('/v1/admin/banners', data)
}

// 更新Banner
export function updateBanner(id: number, data: BannerUpdateRequest) {
  return put<Banner>(`/v1/admin/banners/${id}`, data)
}

// 删除Banner
export function deleteBanner(id: number) {
  return del(`/v1/admin/banners/${id}`)
}

// 更新Banner状态
export function updateBannerStatus(id: number, data: BannerStatusRequest) {
  return put(`/v1/admin/banners/${id}/status`, data)
}

// 批量更新排序
export function updateBannerSort(data: BannerSortRequest) {
  return put('/v1/admin/banners/sort', data)
}

// APP端 - 获取有效Banner列表
export function getActiveBanners() {
  return get<Banner[]>('/v1/banners')
}

// APP端 - 记录点击
export function recordBannerClick(id: number) {
  return post(`/v1/banners/${id}/click`)
}
