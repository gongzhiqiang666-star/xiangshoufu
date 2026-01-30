import { get, put, del } from './request'

// 门槛配置项
export interface ThresholdConfig {
  id: number
  wallet_type: number
  wallet_type_name: string
  channel_id: number
  channel_name?: string
  threshold_amount: number // 分
}

// 门槛列表响应
export interface ThresholdListResponse {
  general_thresholds: ThresholdConfig[]
  channel_thresholds: ThresholdConfig[]
}

// 设置通用门槛请求
export interface SetGeneralThresholdRequest {
  profit_threshold: number      // 分润钱包门槛（分）
  service_fee_threshold: number // 服务费钱包门槛（分）
  reward_threshold: number      // 奖励钱包门槛（分）
}

// 设置通道门槛请求
export interface SetChannelThresholdRequest {
  channel_id: number
  profit_threshold: number      // 分润钱包门槛（分），0表示使用通用门槛
  service_fee_threshold: number // 服务费钱包门槛（分），0表示使用通用门槛
  reward_threshold: number      // 奖励钱包门槛（分），0表示使用通用门槛
}

// 获取所有门槛配置
export function getWithdrawThresholds(): Promise<ThresholdListResponse> {
  return get<ThresholdListResponse>('/v1/withdraw-thresholds')
}

// 设置通用门槛
export function setGeneralThresholds(data: SetGeneralThresholdRequest): Promise<void> {
  return put('/v1/withdraw-thresholds/general', data)
}

// 设置通道门槛
export function setChannelThresholds(data: SetChannelThresholdRequest): Promise<void> {
  return put('/v1/withdraw-thresholds/channel', data)
}

// 删除通道门槛
export function deleteChannelThreshold(channelId: number): Promise<void> {
  return del(`/v1/withdraw-thresholds/channel/${channelId}`)
}
