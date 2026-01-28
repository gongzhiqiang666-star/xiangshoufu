<template>
  <div class="dashboard-v2" v-loading="loading">
    <!-- 顶部收益概览区 - 突出赚钱主题 -->
    <div class="hero-section">
      <div class="hero-background"></div>
      <div class="hero-content">
        <div class="greeting">
          <h1>{{ greeting }}，{{ userName }}</h1>
          <p class="date">{{ currentDate }}</p>
        </div>

        <!-- 核心收益数据 -->
        <div class="profit-showcase">
          <div class="main-profit">
            <div class="profit-label">
              <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
              </svg>
              今日收益
            </div>
            <div class="profit-value">
              <span class="currency">¥</span>
              <span class="amount">{{ formatLargeAmount(todayProfit) }}</span>
            </div>
            <div class="profit-trend" :class="todayTrend >= 0 ? 'up' : 'down'">
              <svg v-if="todayTrend >= 0" class="trend-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
                <polyline points="17 6 23 6 23 12"/>
              </svg>
              <svg v-else class="trend-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/>
                <polyline points="17 18 23 18 23 12"/>
              </svg>
              {{ todayTrend >= 0 ? '+' : '' }}{{ todayTrend.toFixed(1) }}% 较昨日
            </div>
          </div>

          <div class="secondary-profits">
            <div class="profit-card">
              <div class="card-icon transaction">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
                  <line x1="1" y1="10" x2="23" y2="10"/>
                </svg>
              </div>
              <div class="card-info">
                <span class="card-label">交易分润</span>
                <span class="card-value">¥{{ formatAmount(transProfit) }}</span>
              </div>
            </div>
            <div class="profit-card">
              <div class="card-icon deposit">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                </svg>
              </div>
              <div class="card-info">
                <span class="card-label">押金返现</span>
                <span class="card-value">¥{{ formatAmount(depositCashback) }}</span>
              </div>
            </div>
            <div class="profit-card">
              <div class="card-icon sim">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
                  <line x1="12" y1="18" x2="12.01" y2="18"/>
                </svg>
              </div>
              <div class="card-info">
                <span class="card-label">流量返现</span>
                <span class="card-value">¥{{ formatAmount(simCashback) }}</span>
              </div>
            </div>
            <div class="profit-card">
              <div class="card-icon reward">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="8" r="7"/>
                  <polyline points="8.21 13.89 7 23 12 20 17 23 15.79 13.88"/>
                </svg>
              </div>
              <div class="card-info">
                <span class="card-label">激活奖励</span>
                <span class="card-value">¥{{ formatAmount(activationReward) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 数据范围切换 -->
    <div class="scope-switcher">
      <button
        :class="['scope-btn', { active: activeScope === 'direct' }]"
        @click="handleScopeChange('direct')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
        </svg>
        直营数据
      </button>
      <button
        :class="['scope-btn', { active: activeScope === 'team' }]"
        @click="handleScopeChange('team')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        团队数据
      </button>
    </div>

    <!-- 核心指标卡片 -->
    <div class="metrics-grid">
      <div class="metric-card gold">
        <div class="metric-header">
          <span class="metric-title">今日交易额</span>
          <div class="metric-badge" :class="transAmountTrend >= 0 ? 'up' : 'down'">
            {{ transAmountTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(transAmountTrend).toFixed(1) }}%
          </div>
        </div>
        <div class="metric-value">
          <span class="prefix">¥</span>{{ formatLargeAmount(todayTransAmount) }}
        </div>
        <div class="metric-footer">
          <span>昨日: ¥{{ formatAmount(yesterdayTransAmount) }}</span>
        </div>
        <div class="metric-glow"></div>
      </div>

      <div class="metric-card purple">
        <div class="metric-header">
          <span class="metric-title">本月累计收益</span>
          <div class="metric-badge up">
            目标 {{ monthProfitProgress.toFixed(0) }}%
          </div>
        </div>
        <div class="metric-value">
          <span class="prefix">¥</span>{{ formatLargeAmount(monthProfit) }}
        </div>
        <div class="metric-progress">
          <div class="progress-bar" :style="{ width: monthProfitProgress + '%' }"></div>
        </div>
        <div class="metric-footer">
          <span>目标: ¥{{ formatAmount(monthProfitTarget) }}</span>
        </div>
        <div class="metric-glow"></div>
      </div>

      <div class="metric-card green">
        <div class="metric-header">
          <span class="metric-title">今日交易笔数</span>
          <div class="metric-badge" :class="transCountTrend >= 0 ? 'up' : 'down'">
            {{ transCountTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(transCountTrend).toFixed(1) }}%
          </div>
        </div>
        <div class="metric-value">
          {{ formatNumber(todayTransCount) }}<span class="suffix">笔</span>
        </div>
        <div class="metric-footer">
          <span>昨日: {{ formatNumber(yesterdayTransCount) }} 笔</span>
        </div>
        <div class="metric-glow"></div>
      </div>

      <div class="metric-card blue">
        <div class="metric-header">
          <span class="metric-title">活跃商户</span>
          <div class="metric-badge up">
            +{{ newMerchantsToday }} 今日
          </div>
        </div>
        <div class="metric-value">
          {{ formatNumber(activeMerchants) }}<span class="suffix">户</span>
        </div>
        <div class="metric-footer">
          <span>总商户: {{ formatNumber(totalMerchants) }} 户</span>
        </div>
        <div class="metric-glow"></div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-section">
      <!-- 收益趋势图 -->
      <div class="chart-container main-chart">
        <div class="chart-header">
          <h3>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
            </svg>
            收益趋势
          </h3>
          <div class="chart-controls">
            <button
              v-for="days in [7, 15, 30]"
              :key="days"
              :class="['control-btn', { active: chartDays === days }]"
              @click="handleDaysChange(days)"
            >
              近{{ days }}天
            </button>
          </div>
        </div>
        <div class="chart-body" v-loading="chartLoading">
          <LineChart
            v-if="chartData"
            :dates="chartData.dates"
            :data="lineChartData"
            height="320px"
          />
          <div v-else class="chart-empty">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M3 3v18h18"/>
              <path d="M18 9l-5 5-4-4-6 6"/>
            </svg>
            <span>暂无数据</span>
          </div>
        </div>
      </div>

      <!-- 通道收益占比 -->
      <div class="chart-container side-chart">
        <div class="chart-header">
          <h3>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21.21 15.89A10 10 0 1 1 8 2.83"/>
              <path d="M22 12A10 10 0 0 0 12 2v10z"/>
            </svg>
            通道收益占比
          </h3>
        </div>
        <div class="chart-body">
          <PieChart
            v-if="channelStats.length > 0"
            :data="channelChartData"
            height="280px"
          />
          <div v-else class="chart-empty">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 2a10 10 0 0 1 10 10"/>
            </svg>
            <span>暂无数据</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 下级代理排行榜 -->
    <div class="leaderboard-section">
      <div class="leaderboard-container">
        <div class="leaderboard-header">
          <h3>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
            </svg>
            下级代理收益排行榜
          </h3>
          <div class="leaderboard-controls">
            <button
              v-for="period in periods"
              :key="period.value"
              :class="['control-btn', { active: rankingPeriod === period.value }]"
              @click="handlePeriodChange(period.value)"
            >
              {{ period.label }}
            </button>
          </div>
        </div>
        <div class="leaderboard-body">
          <div v-if="agentRanking.length > 0" class="ranking-list">
            <div
              v-for="(agent, index) in agentRanking"
              :key="agent.agent_id"
              :class="['ranking-item', { 'top-three': index < 3 }]"
            >
              <div class="rank-badge" :class="`rank-${index + 1}`">
                <span v-if="index < 3">
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/>
                  </svg>
                </span>
                <span v-else>{{ index + 1 }}</span>
              </div>
              <div class="agent-info">
                <span class="agent-name">{{ agent.agent_name }}</span>
                <span class="agent-level">{{ agent.level || 'Lv1' }}</span>
              </div>
              <div class="agent-value">
                <span class="value">¥{{ formatAmount(agent.value) }}</span>
                <span class="change" :class="agent.change >= 0 ? 'up' : 'down'">
                  {{ agent.change >= 0 ? '+' : '' }}{{ agent.change_rate?.toFixed(1) || 0 }}%
                </span>
              </div>
            </div>
          </div>
          <div v-else class="leaderboard-empty">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
            <span>暂无排名数据</span>
          </div>
        </div>
      </div>

      <!-- 终端激活统计 -->
      <div class="terminal-container">
        <div class="terminal-header">
          <h3>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
              <line x1="8" y1="21" x2="16" y2="21"/>
              <line x1="12" y1="17" x2="12" y2="21"/>
            </svg>
            终端激活
          </h3>
        </div>
        <div class="terminal-body">
          <div class="terminal-stats">
            <div class="terminal-stat">
              <div class="stat-circle total">
                <svg viewBox="0 0 36 36">
                  <circle cx="18" cy="18" r="16" fill="none" stroke="currentColor" stroke-width="3" opacity="0.2"/>
                  <circle cx="18" cy="18" r="16" fill="none" stroke="currentColor" stroke-width="3"
                    :stroke-dasharray="`${activationRate} 100`"
                    stroke-linecap="round"
                    transform="rotate(-90 18 18)"/>
                </svg>
                <div class="stat-value">{{ overview?.terminal?.total || 0 }}</div>
              </div>
              <span class="stat-label">终端总数</span>
            </div>
            <div class="terminal-stat">
              <div class="stat-circle activated">
                <svg viewBox="0 0 36 36">
                  <circle cx="18" cy="18" r="16" fill="none" stroke="currentColor" stroke-width="3" opacity="0.2"/>
                  <circle cx="18" cy="18" r="16" fill="none" stroke="currentColor" stroke-width="3"
                    stroke-dasharray="100 100"
                    stroke-linecap="round"
                    transform="rotate(-90 18 18)"/>
                </svg>
                <div class="stat-value">{{ overview?.terminal?.activated || 0 }}</div>
              </div>
              <span class="stat-label">已激活</span>
            </div>
            <div class="terminal-stat highlight">
              <div class="stat-circle today">
                <div class="stat-value">{{ overview?.terminal?.today_activated || 0 }}</div>
              </div>
              <span class="stat-label">今日激活</span>
            </div>
            <div class="terminal-stat">
              <div class="stat-circle month">
                <div class="stat-value">{{ overview?.terminal?.month_activated || 0 }}</div>
              </div>
              <span class="stat-label">本月激活</span>
            </div>
          </div>
          <div class="activation-rate">
            <span class="rate-label">激活率</span>
            <span class="rate-value">{{ activationRate.toFixed(1) }}%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="quick-actions">
      <button class="action-btn primary" @click="navigateTo('/agents/list')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        代理管理
      </button>
      <button class="action-btn success" @click="navigateTo('/merchants/list')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
          <polyline points="9 22 9 12 15 12 15 22"/>
        </svg>
        商户管理
      </button>
      <button class="action-btn warning" @click="navigateTo('/transactions/list')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="1" x2="12" y2="23"/>
          <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
        </svg>
        交易记录
      </button>
      <button class="action-btn info" @click="navigateTo('/wallets/list')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
          <line x1="1" y1="10" x2="23" y2="10"/>
        </svg>
        钱包管理
      </button>
      <button class="action-btn default" @click="navigateTo('/terminals/list')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
          <line x1="8" y1="21" x2="16" y2="21"/>
          <line x1="12" y1="17" x2="12" y2="21"/>
        </svg>
        终端管理
      </button>
      <button class="action-btn refresh" @click="handleRefresh">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23 4 23 10 17 10"/>
          <polyline points="1 20 1 14 7 14"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
        </svg>
        刷新数据
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboardStore } from '@/stores/dashboard'
import { useUserStore } from '@/stores/user'
import { formatAmount, formatNumber } from '@/utils/format'
import LineChart from '@/components/Charts/LineChart.vue'
import PieChart from '@/components/Charts/PieChart.vue'

const router = useRouter()
const dashboardStore = useDashboardStore()
const userStore = useUserStore()

// 用户信息
const userName = computed(() => userStore.userInfo?.real_name || userStore.agentInfo?.agent_name || '用户')

// 问候语
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了'
})

// 当前日期
const currentDate = computed(() => {
  const now = new Date()
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return `${now.getFullYear()}年${now.getMonth() + 1}月${now.getDate()}日 ${weekdays[now.getDay()]}`
})

// 范围切换
const activeScope = ref('direct')

// 图表天数
const chartDays = ref(7)

// 排名时间周期
const rankingPeriod = ref('month')
const periods = [
  { label: '今日', value: 'day' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
]

// 加载状态
const loading = computed(() => dashboardStore.loading)
const chartLoading = computed(() => dashboardStore.chartLoading)

// 概览数据
const overview = computed(() => dashboardStore.overview)

// 收益数据
const todayProfit = computed(() => overview.value?.today?.profit_total || 0)
const todayTrend = computed(() => overview.value?.today?.profit_change_rate || 0)
const transProfit = computed(() => overview.value?.today?.trans_profit || 0)
const depositCashback = computed(() => overview.value?.today?.deposit_cashback || 0)
const simCashback = computed(() => overview.value?.today?.sim_cashback || 0)
const activationReward = computed(() => overview.value?.today?.activation_reward || 0)

// 交易数据
const todayTransAmount = computed(() => overview.value?.today?.trans_amount || 0)
const yesterdayTransAmount = computed(() => overview.value?.yesterday?.trans_amount || 0)
const transAmountTrend = computed(() => overview.value?.today?.trans_amount_change_rate || 0)

const todayTransCount = computed(() => overview.value?.today?.trans_count || 0)
const yesterdayTransCount = computed(() => overview.value?.yesterday?.trans_count || 0)
const transCountTrend = computed(() => overview.value?.today?.trans_count_change_rate || 0)

// 月度数据
const monthProfit = computed(() => overview.value?.month?.profit_total || 0)
const monthProfitTarget = computed(() => overview.value?.month?.profit_target || 100000)
const monthProfitProgress = computed(() => {
  if (!monthProfitTarget.value) return 0
  return Math.min((monthProfit.value / monthProfitTarget.value) * 100, 100)
})

// 商户数据
const activeMerchants = computed(() => overview.value?.merchant?.active || 0)
const totalMerchants = computed(() => overview.value?.merchant?.total || 0)
const newMerchantsToday = computed(() => overview.value?.merchant?.today_new || 0)

// 激活率
const activationRate = computed(() => {
  const total = overview.value?.terminal?.total || 0
  const activated = overview.value?.terminal?.activated || 0
  if (!total) return 0
  return (activated / total) * 100
})

// 图表数据
const chartData = computed(() => dashboardStore.chartData)
const channelStats = computed(() => dashboardStore.channelStats || [])
const agentRanking = computed(() => dashboardStore.agentRanking || [])

// 折线图数据
const lineChartData = computed(() => {
  if (!chartData.value) return []
  return [
    {
      name: '交易额(元)',
      values: chartData.value.transaction_amounts?.map((v: number) => v / 100) || [],
      color: '#F59E0B',
    },
    {
      name: '收益(元)',
      values: chartData.value.profits?.map((v: number) => v / 100) || [],
      color: '#10B981',
    },
  ]
})

// 饼图数据
const channelChartData = computed(() => {
  return channelStats.value.map((item: any) => ({
    name: item.channel_name,
    value: item.trans_amount / 100,
    percentage: item.percentage,
  }))
})

// 格式化大金额
function formatLargeAmount(amount: number): string {
  const yuan = amount / 100
  if (yuan >= 100000000) {
    return (yuan / 100000000).toFixed(2) + '亿'
  }
  if (yuan >= 10000) {
    return (yuan / 10000).toFixed(2) + '万'
  }
  return yuan.toFixed(2)
}

// 处理范围切换
function handleScopeChange(scope: string) {
  activeScope.value = scope
  dashboardStore.setScope(scope)
  dashboardStore.refreshAll()
}

// 处理天数切换
function handleDaysChange(days: number) {
  chartDays.value = days
  dashboardStore.fetchChartData(days)
}

// 处理排名周期切换
function handlePeriodChange(period: string) {
  rankingPeriod.value = period
  dashboardStore.fetchAgentRanking(period)
}

// 导航
function navigateTo(path: string) {
  router.push(path)
}

// 刷新数据
function handleRefresh() {
  dashboardStore.refreshAll()
}

// 页面加载时获取数据
onMounted(() => {
  dashboardStore.refreshAll()
})
</script>

<style lang="scss" scoped>
// 设计系统变量
$gold-primary: #F59E0B;
$gold-secondary: #FBBF24;
$purple-accent: #8B5CF6;
$green-success: #10B981;
$blue-info: #3B82F6;
$red-danger: #EF4444;
$bg-dark: #0F172A;
$bg-card: #1E293B;
$bg-hover: #334155;
$text-primary: #F8FAFC;
$text-secondary: #94A3B8;
$text-muted: #64748B;

.dashboard-v2 {
  min-height: 100vh;
  background: linear-gradient(135deg, $bg-dark 0%, #1a1f35 100%);
  color: $text-primary;
  padding-bottom: 40px;
}

// Hero Section
.hero-section {
  position: relative;
  padding: 40px 24px 60px;
  overflow: hidden;

  .hero-background {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background:
      radial-gradient(ellipse at 20% 20%, rgba($gold-primary, 0.15) 0%, transparent 50%),
      radial-gradient(ellipse at 80% 80%, rgba($purple-accent, 0.1) 0%, transparent 50%);
    pointer-events: none;
  }

  .hero-content {
    position: relative;
    max-width: 1400px;
    margin: 0 auto;
  }

  .greeting {
    margin-bottom: 32px;

    h1 {
      font-size: 28px;
      font-weight: 600;
      margin: 0 0 8px;
      background: linear-gradient(135deg, $text-primary 0%, $gold-secondary 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .date {
      color: $text-secondary;
      font-size: 14px;
      margin: 0;
    }
  }
}

// Profit Showcase
.profit-showcase {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;

  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
  }

  .main-profit {
    background: linear-gradient(135deg, rgba($gold-primary, 0.2) 0%, rgba($gold-secondary, 0.1) 100%);
    border: 1px solid rgba($gold-primary, 0.3);
    border-radius: 20px;
    padding: 32px;
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: -50%;
      right: -20%;
      width: 200px;
      height: 200px;
      background: radial-gradient(circle, rgba($gold-primary, 0.3) 0%, transparent 70%);
      pointer-events: none;
    }

    .profit-label {
      display: flex;
      align-items: center;
      gap: 8px;
      color: $gold-secondary;
      font-size: 16px;
      margin-bottom: 16px;

      .icon {
        width: 20px;
        height: 20px;
      }
    }

    .profit-value {
      margin-bottom: 16px;

      .currency {
        font-size: 24px;
        color: $gold-secondary;
        margin-right: 4px;
      }

      .amount {
        font-size: 48px;
        font-weight: 700;
        background: linear-gradient(135deg, $gold-primary 0%, $gold-secondary 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
    }

    .profit-trend {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;

      .trend-icon {
        width: 16px;
        height: 16px;
      }

      &.up {
        color: $green-success;
      }

      &.down {
        color: $red-danger;
      }
    }
  }

  .secondary-profits {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;

    .profit-card {
      background: $bg-card;
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 16px;
      padding: 20px;
      display: flex;
      align-items: center;
      gap: 16px;
      transition: all 0.3s ease;
      cursor: pointer;

      &:hover {
        background: $bg-hover;
        transform: translateY(-2px);
        border-color: rgba(255, 255, 255, 0.2);
      }

      .card-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;

        svg {
          width: 24px;
          height: 24px;
        }

        &.transaction {
          background: rgba($gold-primary, 0.2);
          color: $gold-primary;
        }

        &.deposit {
          background: rgba($purple-accent, 0.2);
          color: $purple-accent;
        }

        &.sim {
          background: rgba($blue-info, 0.2);
          color: $blue-info;
        }

        &.reward {
          background: rgba($green-success, 0.2);
          color: $green-success;
        }
      }

      .card-info {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .card-label {
          font-size: 13px;
          color: $text-secondary;
        }

        .card-value {
          font-size: 18px;
          font-weight: 600;
          color: $text-primary;
        }
      }
    }
  }
}

// Scope Switcher
.scope-switcher {
  max-width: 1400px;
  margin: -30px auto 24px;
  padding: 0 24px;
  display: flex;
  gap: 12px;
  position: relative;
  z-index: 1;

  .scope-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 24px;
    background: $bg-card;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    color: $text-secondary;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;

    svg {
      width: 18px;
      height: 18px;
    }

    &:hover {
      background: $bg-hover;
      color: $text-primary;
    }

    &.active {
      background: linear-gradient(135deg, rgba($gold-primary, 0.2) 0%, rgba($purple-accent, 0.1) 100%);
      border-color: $gold-primary;
      color: $gold-primary;
    }
  }
}

// Metrics Grid
.metrics-grid {
  max-width: 1400px;
  margin: 0 auto 24px;
  padding: 0 24px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;

  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }

  .metric-card {
    background: $bg-card;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 16px;
    padding: 20px;
    position: relative;
    overflow: hidden;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-4px);

      .metric-glow {
        opacity: 1;
      }
    }

    .metric-glow {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 2px;
      opacity: 0;
      transition: opacity 0.3s ease;
    }

    &.gold {
      .metric-glow { background: linear-gradient(90deg, $gold-primary, $gold-secondary); }
      .metric-value { color: $gold-primary; }
    }

    &.purple {
      .metric-glow { background: linear-gradient(90deg, $purple-accent, #A78BFA); }
      .metric-value { color: $purple-accent; }
    }

    &.green {
      .metric-glow { background: linear-gradient(90deg, $green-success, #34D399); }
      .metric-value { color: $green-success; }
    }

    &.blue {
      .metric-glow { background: linear-gradient(90deg, $blue-info, #60A5FA); }
      .metric-value { color: $blue-info; }
    }

    .metric-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .metric-title {
        font-size: 13px;
        color: $text-secondary;
      }

      .metric-badge {
        font-size: 12px;
        padding: 4px 8px;
        border-radius: 6px;

        &.up {
          background: rgba($green-success, 0.2);
          color: $green-success;
        }

        &.down {
          background: rgba($red-danger, 0.2);
          color: $red-danger;
        }
      }
    }

    .metric-value {
      font-size: 28px;
      font-weight: 700;
      margin-bottom: 12px;

      .prefix, .suffix {
        font-size: 14px;
        font-weight: 400;
        color: $text-secondary;
      }
    }

    .metric-progress {
      height: 4px;
      background: rgba(255, 255, 255, 0.1);
      border-radius: 2px;
      margin-bottom: 12px;
      overflow: hidden;

      .progress-bar {
        height: 100%;
        background: linear-gradient(90deg, $purple-accent, #A78BFA);
        border-radius: 2px;
        transition: width 0.6s ease;
      }
    }

    .metric-footer {
      font-size: 12px;
      color: $text-muted;
    }
  }
}

// Charts Section
.charts-section {
  max-width: 1400px;
  margin: 0 auto 24px;
  padding: 0 24px;
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;

  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
  }

  .chart-container {
    background: $bg-card;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 16px;
    overflow: hidden;

    .chart-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.05);

      h3 {
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 16px;
        font-weight: 500;
        margin: 0;
        color: $text-primary;

        svg {
          width: 20px;
          height: 20px;
          color: $gold-primary;
        }
      }

      .chart-controls {
        display: flex;
        gap: 8px;

        .control-btn {
          padding: 6px 14px;
          background: transparent;
          border: 1px solid rgba(255, 255, 255, 0.1);
          border-radius: 8px;
          color: $text-secondary;
          font-size: 13px;
          cursor: pointer;
          transition: all 0.2s ease;

          &:hover {
            background: $bg-hover;
            color: $text-primary;
          }

          &.active {
            background: rgba($gold-primary, 0.2);
            border-color: $gold-primary;
            color: $gold-primary;
          }
        }
      }
    }

    .chart-body {
      padding: 20px;
      min-height: 320px;
    }

    .chart-empty {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 280px;
      color: $text-muted;

      svg {
        width: 48px;
        height: 48px;
        margin-bottom: 12px;
        opacity: 0.5;
      }
    }
  }
}

// Leaderboard Section
.leaderboard-section {
  max-width: 1400px;
  margin: 0 auto 24px;
  padding: 0 24px;
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 16px;

  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
  }

  .leaderboard-container,
  .terminal-container {
    background: $bg-card;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 16px;
    overflow: hidden;
  }

  .leaderboard-header,
  .terminal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);

    h3 {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 16px;
      font-weight: 500;
      margin: 0;
      color: $text-primary;

      svg {
        width: 20px;
        height: 20px;
        color: $gold-primary;
      }
    }

    .leaderboard-controls {
      display: flex;
      gap: 8px;

      .control-btn {
        padding: 6px 14px;
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        color: $text-secondary;
        font-size: 13px;
        cursor: pointer;
        transition: all 0.2s ease;

        &:hover {
          background: $bg-hover;
          color: $text-primary;
        }

        &.active {
          background: rgba($gold-primary, 0.2);
          border-color: $gold-primary;
          color: $gold-primary;
        }
      }
    }
  }

  .leaderboard-body {
    padding: 16px;
    max-height: 380px;
    overflow-y: auto;

    .ranking-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 14px 16px;
      border-radius: 12px;
      margin-bottom: 8px;
      transition: all 0.2s ease;
      cursor: pointer;

      &:last-child {
        margin-bottom: 0;
      }

      &:hover {
        background: $bg-hover;
      }

      &.top-three {
        background: linear-gradient(135deg, rgba($gold-primary, 0.1) 0%, rgba($gold-secondary, 0.05) 100%);

        &:hover {
          background: linear-gradient(135deg, rgba($gold-primary, 0.15) 0%, rgba($gold-secondary, 0.1) 100%);
        }
      }

      .rank-badge {
        width: 32px;
        height: 32px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: 600;
        background: $bg-hover;
        color: $text-secondary;

        svg {
          width: 18px;
          height: 18px;
        }

        &.rank-1 {
          background: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);
          color: #000;
        }

        &.rank-2 {
          background: linear-gradient(135deg, #C0C0C0 0%, #A0A0A0 100%);
          color: #000;
        }

        &.rank-3 {
          background: linear-gradient(135deg, #CD7F32 0%, #8B4513 100%);
          color: #fff;
        }
      }

      .agent-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 2px;

        .agent-name {
          font-size: 14px;
          font-weight: 500;
          color: $text-primary;
        }

        .agent-level {
          font-size: 12px;
          color: $text-muted;
        }
      }

      .agent-value {
        text-align: right;

        .value {
          display: block;
          font-size: 15px;
          font-weight: 600;
          color: $gold-primary;
        }

        .change {
          font-size: 12px;

          &.up { color: $green-success; }
          &.down { color: $red-danger; }
        }
      }
    }
  }

  .leaderboard-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    color: $text-muted;

    svg {
      width: 48px;
      height: 48px;
      margin-bottom: 12px;
      opacity: 0.5;
    }
  }

  // Terminal Stats
  .terminal-body {
    padding: 24px;

    .terminal-stats {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 20px;
      margin-bottom: 24px;

      .terminal-stat {
        text-align: center;

        .stat-circle {
          width: 80px;
          height: 80px;
          margin: 0 auto 12px;
          position: relative;
          display: flex;
          align-items: center;
          justify-content: center;

          svg {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
          }

          .stat-value {
            font-size: 22px;
            font-weight: 600;
            color: $text-primary;
          }

          &.total { color: $blue-info; }
          &.activated { color: $green-success; }
          &.today {
            background: linear-gradient(135deg, rgba($gold-primary, 0.2) 0%, rgba($gold-secondary, 0.1) 100%);
            border-radius: 50%;

            .stat-value { color: $gold-primary; }
          }
          &.month {
            background: rgba($purple-accent, 0.15);
            border-radius: 50%;

            .stat-value { color: $purple-accent; }
          }
        }

        .stat-label {
          font-size: 13px;
          color: $text-secondary;
        }

        &.highlight .stat-label {
          color: $gold-primary;
        }
      }
    }

    .activation-rate {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px 20px;
      background: linear-gradient(135deg, rgba($green-success, 0.15) 0%, rgba($green-success, 0.05) 100%);
      border-radius: 12px;

      .rate-label {
        font-size: 14px;
        color: $text-secondary;
      }

      .rate-value {
        font-size: 24px;
        font-weight: 700;
        color: $green-success;
      }
    }
  }
}

// Quick Actions
.quick-actions {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;

  .action-btn {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 14px 24px;
    border: none;
    border-radius: 12px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;

    svg {
      width: 20px;
      height: 20px;
    }

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
    }

    &.primary {
      background: linear-gradient(135deg, $gold-primary 0%, $gold-secondary 100%);
      color: #000;
    }

    &.success {
      background: linear-gradient(135deg, $green-success 0%, #34D399 100%);
      color: #fff;
    }

    &.warning {
      background: linear-gradient(135deg, $purple-accent 0%, #A78BFA 100%);
      color: #fff;
    }

    &.info {
      background: linear-gradient(135deg, $blue-info 0%, #60A5FA 100%);
      color: #fff;
    }

    &.default {
      background: $bg-card;
      border: 1px solid rgba(255, 255, 255, 0.1);
      color: $text-primary;

      &:hover {
        background: $bg-hover;
      }
    }

    &.refresh {
      background: transparent;
      border: 1px solid rgba(255, 255, 255, 0.2);
      color: $text-secondary;

      &:hover {
        background: $bg-card;
        color: $text-primary;
      }
    }
  }
}

// 滚动条样式
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;

  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
}
</style>
