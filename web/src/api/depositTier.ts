import { get, post, put, del } from './request'

// ============================================================
// 类型定义
// ============================================================

/** 押金档位 */
export interface DepositTier {
  id: number
  channel_id: number
  brand_code: string
  tier_code: string
  deposit_amount: number
  tier_name: string
  sort_order: number
  status: number
  created_at: string
  updated_at: string
}

/** 创建押金档位请求 */
export interface CreateDepositTierRequest {
  channel_id: number
  brand_code?: string
  tier_code: string
  deposit_amount: number
  tier_name: string
  sort_order?: number
}

/** 更新押金档位请求 */
export interface UpdateDepositTierRequest {
  tier_code?: string
  deposit_amount?: number
  tier_name?: string
  sort_order?: number
  status?: number
}

// ============================================================
// API 函数
// ============================================================

/** 获取押金档位列表（按通道） */
export function getDepositTiersByChannel(channelId: number, brandCode?: string) {
  const params: Record<string, any> = {}
  if (brandCode) {
    params.brand_code = brandCode
  }
  return get<DepositTier[]>(`/v1/channels/${channelId}/deposit-tiers`, params)
}

/** 获取押金档位列表 */
export function getDepositTiers(channelId: number, brandCode?: string) {
  const params: Record<string, any> = { channel_id: channelId }
  if (brandCode) {
    params.brand_code = brandCode
  }
  return get<DepositTier[]>('/v1/deposit-tiers', params)
}

/** 获取押金档位详情 */
export function getDepositTier(id: number) {
  return get<DepositTier>(`/v1/deposit-tiers/${id}`)
}

/** 创建押金档位 */
export function createDepositTier(data: CreateDepositTierRequest) {
  return post<DepositTier>('/v1/deposit-tiers', data)
}

/** 更新押金档位 */
export function updateDepositTier(id: number, data: UpdateDepositTierRequest) {
  return put<DepositTier>(`/v1/deposit-tiers/${id}`, data)
}

/** 删除押金档位 */
export function deleteDepositTier(id: number) {
  return del(`/v1/deposit-tiers/${id}`)
}
