import request from './request'

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

// 获取费率同步日志列表
export function getRateSyncLogs(params: RateSyncLogQuery) {
  return request.get<{
    items: RateSyncLog[]
    total: number
    page: number
    page_size: number
  }>('/api/v1/rate-sync/logs', { params })
}

// 获取费率同步日志详情
export function getRateSyncLogDetail(id: number) {
  return request.get<RateSyncLog>(`/api/v1/rate-sync/logs/${id}`)
}
