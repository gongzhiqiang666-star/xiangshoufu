import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DashboardOverview, ChartData, StatCardData } from '@/types'
import { getDashboardOverview, getChartData } from '@/api/dashboard'
import { formatAmount, formatNumber, calculateTrend } from '@/utils/format'

export const useDashboardStore = defineStore('dashboard', () => {
  // 状态
  const overview = ref<DashboardOverview | null>(null)
  const chartData = ref<ChartData | null>(null)
  const loading = ref(false)
  const chartLoading = ref(false)

  // 统计卡片数据
  const statCards = ref<StatCardData[]>([])

  // 获取概览数据
  async function fetchOverview() {
    loading.value = true
    try {
      const data = await getDashboardOverview()
      overview.value = data

      // 生成统计卡片数据
      statCards.value = [
        {
          title: '今日交易额',
          value: formatAmount(data.today_transaction_amount),
          prefix: '¥',
          trend: calculateTrend(data.today_transaction_amount, data.yesterday_transaction_amount),
          trendLabel: '较昨日',
          icon: 'Money',
          color: '#409eff',
        },
        {
          title: '今日交易笔数',
          value: formatNumber(data.today_transaction_count),
          suffix: '笔',
          trend: calculateTrend(data.today_transaction_count, data.yesterday_transaction_count),
          trendLabel: '较昨日',
          icon: 'Document',
          color: '#67c23a',
        },
        {
          title: '今日分润',
          value: formatAmount(data.today_profit),
          prefix: '¥',
          trend: calculateTrend(data.today_profit, data.yesterday_profit),
          trendLabel: '较昨日',
          icon: 'TrendCharts',
          color: '#e6a23c',
        },
        {
          title: '今日新增商户',
          value: formatNumber(data.today_new_merchants),
          suffix: '户',
          trend: calculateTrend(data.today_new_merchants, data.yesterday_new_merchants),
          trendLabel: '较昨日',
          icon: 'Shop',
          color: '#f56c6c',
        },
      ]

      return data
    } finally {
      loading.value = false
    }
  }

  // 获取图表数据
  async function fetchChartData(days = 7) {
    chartLoading.value = true
    try {
      const data = await getChartData(days)
      chartData.value = data
      return data
    } finally {
      chartLoading.value = false
    }
  }

  // 刷新所有数据
  async function refreshAll() {
    await Promise.all([fetchOverview(), fetchChartData()])
  }

  return {
    // 状态
    overview,
    chartData,
    loading,
    chartLoading,
    statCards,
    // 方法
    fetchOverview,
    fetchChartData,
    refreshAll,
  }
})
