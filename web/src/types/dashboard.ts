// 仪表盘概览数据
export interface DashboardOverview {
  // 今日数据
  today_transaction_amount: number
  today_transaction_count: number
  today_profit: number
  today_new_merchants: number

  // 昨日数据（用于计算同比）
  yesterday_transaction_amount: number
  yesterday_transaction_count: number
  yesterday_profit: number
  yesterday_new_merchants: number

  // 本月数据
  month_transaction_amount: number
  month_transaction_count: number
  month_profit: number
  month_new_merchants: number
}

// 图表数据
export interface ChartData {
  dates: string[]
  transaction_amounts: number[]
  transaction_counts: number[]
  profits: number[]
}

// 统计卡片数据
export interface StatCardData {
  title: string
  value: number | string
  prefix?: string
  suffix?: string
  trend?: number
  trendLabel?: string
  icon?: string
  color?: string
}
