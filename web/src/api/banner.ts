import request from '@/utils/request'
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
  return request<{ data: Banner[]; total: number }>({
    url: '/api/v1/admin/banners',
    method: 'get',
    params
  })
}

// 获取Banner详情
export function getBannerDetail(id: number) {
  return request<Banner>({
    url: `/api/v1/admin/banners/${id}`,
    method: 'get'
  })
}

// 创建Banner
export function createBanner(data: BannerCreateRequest) {
  return request<Banner>({
    url: '/api/v1/admin/banners',
    method: 'post',
    data
  })
}

// 更新Banner
export function updateBanner(id: number, data: BannerUpdateRequest) {
  return request<Banner>({
    url: `/api/v1/admin/banners/${id}`,
    method: 'put',
    data
  })
}

// 删除Banner
export function deleteBanner(id: number) {
  return request({
    url: `/api/v1/admin/banners/${id}`,
    method: 'delete'
  })
}

// 更新Banner状态
export function updateBannerStatus(id: number, data: BannerStatusRequest) {
  return request({
    url: `/api/v1/admin/banners/${id}/status`,
    method: 'put',
    data
  })
}

// 批量更新排序
export function updateBannerSort(data: BannerSortRequest) {
  return request({
    url: '/api/v1/admin/banners/sort',
    method: 'put',
    data
  })
}
