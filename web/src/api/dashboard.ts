import { get } from './request'
import type { DashboardOverview, ChartData } from '@/types'

/**
 * 获取仪表盘概览数据
 */
export function getDashboardOverview(): Promise<DashboardOverview> {
  return get<DashboardOverview>('/v1/dashboard/overview')
}

/**
 * 获取图表数据
 * @param days 天数（默认7天）
 */
export function getChartData(days = 7): Promise<ChartData> {
  return get<ChartData>('/v1/dashboard/charts', { days })
}

/**
 * 获取今日实时数据
 */
export function getTodayRealtime(): Promise<{
  transaction_amount: number
  transaction_count: number
  profit: number
}> {
  return get('/v1/dashboard/realtime')
}
