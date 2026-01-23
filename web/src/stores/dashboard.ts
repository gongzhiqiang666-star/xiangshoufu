import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { StatCardData } from '@/types'
import {
  getDashboardOverview,
  getChartData,
  getChannelStats,
  getMerchantDistribution,
  getAgentRanking,
} from '@/api/dashboard'
import { formatAmount, formatNumber, calculateTrend } from '@/utils/format'

// 概览数据类型
interface OverviewData {
  today: {
    trans_amount: number
    trans_amount_yuan: number
    trans_count: number
    profit_total: number
    profit_total_yuan: number
    profit_trade: number
    profit_deposit: number
    profit_sim: number
    profit_reward: number
  }
  yesterday: {
    trans_amount: number
    trans_amount_yuan: number
    trans_count: number
    profit_total: number
    profit_total_yuan: number
  }
  month: {
    trans_amount: number
    trans_amount_yuan: number
    trans_count: number
    profit_total: number
    profit_total_yuan: number
    merchant_new: number
  }
  team: {
    direct_agent_count: number
    direct_merchant_count: number
    team_agent_count: number
    team_merchant_count: number
  }
  terminal: {
    total: number
    activated: number
    today_activated: number
    month_activated: number
  }
  wallet: {
    total_balance: number
    total_balance_yuan: number
  }
}

// 图表数据类型
interface ChartData {
  dates: string[]
  transaction_amounts: number[]
  profits: number[]
}

// 通道统计
interface ChannelStatsItem {
  channel_id: number
  channel_code: string
  channel_name: string
  trans_amount: number
  trans_count: number
  percentage: number
}

// 商户分布
interface MerchantDistributionItem {
  merchant_type: string
  type_name: string
  count: number
  percentage: number
}

// 代理商排名
interface AgentRankingItem {
  rank: number
  agent_id: number
  agent_name: string
  agent_no: string
  value: number
  value_yuan: number
  change: number
  change_rate: number
}

export const useDashboardStore = defineStore('dashboard', () => {
  // 状态
  const scope = ref('direct') // 'direct' | 'team'
  const overview = ref<OverviewData | null>(null)
  const chartData = ref<ChartData | null>(null)
  const channelStats = ref<ChannelStatsItem[]>([])
  const merchantDistribution = ref<MerchantDistributionItem[]>([])
  const agentRanking = ref<AgentRankingItem[]>([])
  const loading = ref(false)
  const chartLoading = ref(false)

  // 统计卡片数据
  const statCards = ref<StatCardData[]>([])

  // 设置统计范围
  function setScope(newScope: string) {
    scope.value = newScope
  }

  // 获取概览数据
  async function fetchOverview() {
    loading.value = true
    try {
      const data = await getDashboardOverview(scope.value)
      overview.value = data

      // 生成统计卡片数据(6个卡片)
      statCards.value = [
        {
          title: '今日交易额',
          value: formatAmount(data.today?.trans_amount_yuan || 0),
          prefix: '¥',
          trend: calculateTrend(data.today?.trans_amount || 0, data.yesterday?.trans_amount || 0),
          trendLabel: '较昨日',
          icon: 'Money',
          color: '#409eff',
        },
        {
          title: '今日收益',
          value: formatAmount(data.today?.profit_total_yuan || 0),
          prefix: '¥',
          trend: calculateTrend(data.today?.profit_total || 0, data.yesterday?.profit_total || 0),
          trendLabel: '较昨日',
          icon: 'TrendCharts',
          color: '#67c23a',
        },
        {
          title: '交易分润',
          value: formatAmount((data.today?.profit_trade || 0) / 100),
          prefix: '¥',
          icon: 'Connection',
          color: '#409eff',
        },
        {
          title: '押金返现',
          value: formatAmount((data.today?.profit_deposit || 0) / 100),
          prefix: '¥',
          icon: 'Coin',
          color: '#e6a23c',
        },
        {
          title: '流量返现',
          value: formatAmount((data.today?.profit_sim || 0) / 100),
          prefix: '¥',
          icon: 'Connection',
          color: '#909399',
        },
        {
          title: '激活奖励',
          value: formatAmount((data.today?.profit_reward || 0) / 100),
          prefix: '¥',
          icon: 'Present',
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
      const data = await getChartData(days, scope.value)
      // 转换趋势数据格式
      if (data.trans_trend) {
        chartData.value = {
          dates: data.trans_trend.map((t: any) => t.date),
          transaction_amounts: data.trans_trend.map((t: any) => t.trans_amount),
          profits: data.trans_trend.map((t: any) => t.profit_total),
        }
      }
      return data
    } finally {
      chartLoading.value = false
    }
  }

  // 获取通道统计
  async function fetchChannelStats() {
    try {
      const data = await getChannelStats(scope.value)
      channelStats.value = data.channel_stats || []
      return data
    } catch (e) {
      console.error('获取通道统计失败:', e)
      channelStats.value = []
    }
  }

  // 获取商户分布
  async function fetchMerchantDistribution() {
    try {
      const data = await getMerchantDistribution(scope.value)
      merchantDistribution.value = data.distribution || []
      return data
    } catch (e) {
      console.error('获取商户分布失败:', e)
      merchantDistribution.value = []
    }
  }

  // 获取代理商排名
  async function fetchAgentRanking(period = 'month') {
    try {
      const data = await getAgentRanking(period)
      agentRanking.value = data.ranking || []
      return data
    } catch (e) {
      console.error('获取代理商排名失败:', e)
      agentRanking.value = []
    }
  }

  // 刷新所有数据
  async function refreshAll() {
    await Promise.all([
      fetchOverview(),
      fetchChartData(),
      fetchChannelStats(),
      fetchMerchantDistribution(),
      fetchAgentRanking(),
    ])
  }

  return {
    // 状态
    scope,
    overview,
    chartData,
    channelStats,
    merchantDistribution,
    agentRanking,
    loading,
    chartLoading,
    statCards,
    // 方法
    setScope,
    fetchOverview,
    fetchChartData,
    fetchChannelStats,
    fetchMerchantDistribution,
    fetchAgentRanking,
    refreshAll,
  }
})
