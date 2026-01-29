import { get } from './request'

// 费率同步日志
export interface RateSyncLog {
  id: number
  merchant_id: number
  merchant_no: string
  terminal_sn: string
  channel_code: string
  old_credit_rate: number
  old_debit_rate: number
  new_credit_rate: number
  new_debit_rate: number
  sync_status: number
  sync_status_name: string
  channel_trade_no: string
  error_message: string
  created_at: string
  synced_at: string
}

// 查询参数
export interface RateSyncLogQuery {
  merchant_id?: number
  sync_status?: number
  page?: number
  page_size?: number
}

// 费率同步日志列表响应
export interface RateSyncLogListResponse {
  items: RateSyncLog[]
  total: number
  page: number
  page_size: number
}

// 获取费率同步日志列表
export function getRateSyncLogs(params: RateSyncLogQuery): Promise<RateSyncLogListResponse> {
  return get<RateSyncLogListResponse>('/v1/rate-sync/logs', params)
}

// 获取费率同步日志详情
export function getRateSyncLogDetail(id: number): Promise<RateSyncLog> {
  return get<RateSyncLog>(`/v1/rate-sync/logs/${id}`)
}
