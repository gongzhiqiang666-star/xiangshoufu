import { get } from './request'

/**
 * 获取仪表盘概览数据
 * @param scope 统计范围 (direct=直营, team=团队)
 */
export function getDashboardOverview(scope = 'direct'): Promise<any> {
  return get('/v1/dashboard/overview', { scope })
}

/**
 * 获取图表数据
 * @param days 天数（默认7天）
 * @param scope 统计范围
 */
export function getChartData(days = 7, scope = 'direct'): Promise<any> {
  return get('/v1/dashboard/charts', { days, scope })
}

/**
 * 获取通道统计
 * @param scope 统计范围
 * @param period 时间范围 (day/week/month)
 */
export function getChannelStats(scope = 'direct', period = 'month'): Promise<any> {
  return get('/v1/dashboard/channel-stats', { scope, period })
}

/**
 * 获取商户类型分布
 * @param scope 统计范围
 */
export function getMerchantDistribution(scope = 'direct'): Promise<any> {
  return get('/v1/dashboard/merchant-distribution', { scope })
}

/**
 * 获取最近交易列表
 * @param limit 数量限制
 */
export function getRecentTransactions(limit = 10): Promise<any> {
  return get('/v1/dashboard/recent-transactions', { limit })
}

/**
 * 获取代理商排名
 * @param period 时间范围 (day/week/month)
 * @param rankBy 排名维度 (trans_amount/profit/terminal)
 * @param limit 数量限制
 */
export function getAgentRanking(period = 'month', rankBy = 'trans_amount', limit = 10): Promise<any> {
  return get('/v1/analytics/agent-ranking', { period, rank_by: rankBy, limit })
}

/**
 * 获取商户排名
 * @param merchantType 商户类型筛选
 * @param scope 统计范围
 * @param limit 数量限制
 */
export function getMerchantRanking(merchantType = 'all', scope = 'direct', limit = 20): Promise<any> {
  return get('/v1/analytics/merchant-ranking', { merchant_type: merchantType, scope, limit })
}

/**
 * 获取分析汇总数据
 * @param period 时间范围
 * @param scope 统计范围
 */
export function getAnalyticsSummary(period = 'month', scope = 'direct'): Promise<any> {
  return get('/v1/analytics/summary', { period, scope })
}
