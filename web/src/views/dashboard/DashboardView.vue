<template>
  <div class="dashboard-view" v-loading="loading">
    <!-- 范围切换Tab -->
    <el-tabs v-model="activeScope" class="scope-tabs" @tab-change="handleScopeChange">
      <el-tab-pane label="直营数据" name="direct" />
      <el-tab-pane label="团队数据" name="team" />
    </el-tabs>

    <!-- 统计卡片区域 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col
        v-for="(card, index) in statCards"
        :key="index"
        :xs="12"
        :sm="8"
        :lg="4"
      >
        <StatCard
          :title="card.title"
          :value="card.value"
          :prefix="card.prefix"
          :suffix="card.suffix"
          :trend="card.trend"
          :trend-label="card.trendLabel"
          :icon="card.icon"
          :color="card.color"
        />
      </el-col>
    </el-row>

    <!-- 图表区域第一行 -->
    <el-row :gutter="20" class="chart-section">
      <!-- 交易趋势图 -->
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card" v-loading="chartLoading">
          <template #header>
            <div class="card-header">
              <span class="title">交易趋势</span>
              <el-radio-group v-model="chartDays" size="small" @change="handleDaysChange">
                <el-radio-button :label="7">近7天</el-radio-button>
                <el-radio-button :label="15">近15天</el-radio-button>
                <el-radio-button :label="30">近30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <LineChart
            v-if="chartData"
            :dates="chartData.dates"
            :data="lineChartData"
            height="350px"
          />
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>

      <!-- 通道占比饼图 -->
      <el-col :xs="24" :lg="8">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">通道占比</span>
            </div>
          </template>
          <PieChart
            v-if="channelStats.length > 0"
            :data="channelChartData"
            height="350px"
          />
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域第二行 -->
    <el-row :gutter="20" class="chart-section">
      <!-- 商户类型分布 -->
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">商户类型分布</span>
            </div>
          </template>
          <BarChart
            v-if="merchantDistribution.length > 0"
            :data="merchantChartData"
            height="300px"
          />
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>

      <!-- 下级代理排名 -->
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card ranking-card">
          <template #header>
            <div class="card-header">
              <span class="title">下级代理排名 TOP10</span>
              <el-select v-model="rankingPeriod" size="small" @change="fetchAgentRanking">
                <el-option label="本日" value="day" />
                <el-option label="本周" value="week" />
                <el-option label="本月" value="month" />
              </el-select>
            </div>
          </template>
          <div v-if="agentRanking.length > 0" class="ranking-list">
            <div
              v-for="(agent, index) in agentRanking"
              :key="agent.agent_id"
              class="ranking-item"
            >
              <span class="rank" :class="{ 'top-3': index < 3 }">{{ index + 1 }}</span>
              <span class="name">{{ agent.agent_name }}</span>
              <span class="value">¥{{ formatAmount(agent.value) }}</span>
              <span class="change" :class="agent.change >= 0 ? 'up' : 'down'">
                {{ agent.change >= 0 ? '↑' : '↓' }}{{ Math.abs(agent.change_rate).toFixed(1) }}%
              </span>
            </div>
          </div>
          <el-empty v-else description="暂无排名数据" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 本月概览 + 终端统计 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card month-overview">
          <template #header>
            <div class="card-header">
              <span class="title">本月概览</span>
            </div>
          </template>
          <div v-if="overview" class="month-stats">
            <div class="stat-item">
              <div class="stat-label">本月交易额</div>
              <div class="stat-value">
                <span class="prefix">¥</span>
                {{ formatAmount(overview.month?.trans_amount_yuan || 0) }}
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-label">本月交易笔数</div>
              <div class="stat-value">
                {{ formatNumber(overview.month?.trans_count || 0) }}
                <span class="suffix">笔</span>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-label">本月分润</div>
              <div class="stat-value">
                <span class="prefix">¥</span>
                {{ formatAmount(overview.month?.profit_total_yuan || 0) }}
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-label">本月新增商户</div>
              <div class="stat-value">
                {{ formatNumber(overview.month?.merchant_new || 0) }}
                <span class="suffix">户</span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card class="chart-card terminal-stats">
          <template #header>
            <div class="card-header">
              <span class="title">终端统计</span>
            </div>
          </template>
          <div v-if="overview?.terminal" class="terminal-grid">
            <div class="terminal-item">
              <div class="terminal-value">{{ overview.terminal.total }}</div>
              <div class="terminal-label">总数</div>
            </div>
            <div class="terminal-item">
              <div class="terminal-value success">{{ overview.terminal.activated }}</div>
              <div class="terminal-label">已激活</div>
            </div>
            <div class="terminal-item">
              <div class="terminal-value primary">{{ overview.terminal.today_activated }}</div>
              <div class="terminal-label">今日激活</div>
            </div>
            <div class="terminal-item">
              <div class="terminal-value warning">{{ overview.terminal.month_activated }}</div>
              <div class="terminal-label">本月激活</div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作区域 -->
    <el-row :gutter="20" class="quick-actions">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="title">快捷操作</span>
            </div>
          </template>
          <div class="action-buttons">
            <el-button type="primary" @click="navigateTo('/agents/list')">
              <el-icon><User /></el-icon>
              代理管理
            </el-button>
            <el-button type="success" @click="navigateTo('/merchants/list')">
              <el-icon><Shop /></el-icon>
              商户管理
            </el-button>
            <el-button type="warning" @click="navigateTo('/transactions/list')">
              <el-icon><Money /></el-icon>
              交易记录
            </el-button>
            <el-button type="info" @click="navigateTo('/wallets/list')">
              <el-icon><Wallet /></el-icon>
              钱包管理
            </el-button>
            <el-button @click="navigateTo('/terminals/list')">
              <el-icon><Monitor /></el-icon>
              终端管理
            </el-button>
            <el-button @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新数据
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Shop, Money, Wallet, Refresh, Monitor } from '@element-plus/icons-vue'
import { useDashboardStore } from '@/stores/dashboard'
import { formatAmount, formatNumber } from '@/utils/format'
import StatCard from '@/components/Common/StatCard.vue'
import LineChart from '@/components/Charts/LineChart.vue'
import PieChart from '@/components/Charts/PieChart.vue'
import BarChart from '@/components/Charts/BarChart.vue'

const router = useRouter()
const dashboardStore = useDashboardStore()

// 范围切换
const activeScope = ref('direct')

// 图表天数
const chartDays = ref(7)

// 排名时间周期
const rankingPeriod = ref('month')

// 加载状态
const loading = computed(() => dashboardStore.loading)
const chartLoading = computed(() => dashboardStore.chartLoading)

// 统计卡片数据
const statCards = computed(() => dashboardStore.statCards)

// 概览数据
const overview = computed(() => dashboardStore.overview)

// 图表数据
const chartData = computed(() => dashboardStore.chartData)

// 通道统计
const channelStats = computed(() => dashboardStore.channelStats || [])

// 商户分布
const merchantDistribution = computed(() => dashboardStore.merchantDistribution || [])

// 代理商排名
const agentRanking = computed(() => dashboardStore.agentRanking || [])

// 折线图数据
const lineChartData = computed(() => {
  if (!chartData.value) return []
  return [
    {
      name: '交易额(元)',
      values: chartData.value.transaction_amounts.map((v: number) => v / 100),
      color: '#409eff',
    },
    {
      name: '分润(元)',
      values: chartData.value.profits.map((v: number) => v / 100),
      color: '#67c23a',
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

// 柱状图数据
const merchantChartData = computed(() => {
  return {
    categories: merchantDistribution.value.map((item: any) => item.type_name),
    values: merchantDistribution.value.map((item: any) => item.count),
  }
})

// 处理范围切换
function handleScopeChange(scope: string) {
  dashboardStore.setScope(scope)
  dashboardStore.refreshAll()
}

// 处理天数切换
function handleDaysChange(days: number) {
  dashboardStore.fetchChartData(days)
}

// 获取代理商排名
function fetchAgentRanking() {
  dashboardStore.fetchAgentRanking(rankingPeriod.value)
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
.dashboard-view {
  min-height: 100%;
}

.scope-tabs {
  margin-bottom: $spacing-lg;

  :deep(.el-tabs__header) {
    margin-bottom: 0;
  }
}

.stat-cards {
  margin-bottom: $spacing-lg;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.chart-section {
  margin-bottom: $spacing-lg;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.chart-card {
  height: 100%;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .title {
      font-size: 16px;
      font-weight: 500;
      color: $text-primary;
    }
  }
}

.month-overview {
  .month-stats {
    display: flex;
    flex-direction: column;
    gap: $spacing-lg;
  }

  .stat-item {
    padding: $spacing-md;
    background: $bg-color;
    border-radius: $border-radius-sm;
    transition: all $transition-normal;

    &:hover {
      background: lighten($bg-color, 2%);
      transform: translateX(4px);
    }
  }

  .stat-label {
    font-size: 13px;
    color: $text-secondary;
    margin-bottom: $spacing-xs;
  }

  .stat-value {
    font-size: 20px;
    font-weight: 600;
    color: $text-primary;

    .prefix,
    .suffix {
      font-size: 14px;
      font-weight: normal;
      color: $text-secondary;
    }
  }
}

.terminal-stats {
  .terminal-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: $spacing-lg;
  }

  .terminal-item {
    text-align: center;
    padding: $spacing-lg;
    background: $bg-color;
    border-radius: $border-radius-sm;
  }

  .terminal-value {
    font-size: 28px;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: $spacing-xs;

    &.success { color: $success-color; }
    &.primary { color: $primary-color; }
    &.warning { color: $warning-color; }
  }

  .terminal-label {
    font-size: 13px;
    color: $text-secondary;
  }
}

.ranking-card {
  .ranking-list {
    max-height: 300px;
    overflow-y: auto;
  }

  .ranking-item {
    display: flex;
    align-items: center;
    padding: $spacing-sm $spacing-md;
    border-bottom: 1px solid $border-color;

    &:last-child {
      border-bottom: none;
    }

    .rank {
      width: 24px;
      height: 24px;
      line-height: 24px;
      text-align: center;
      border-radius: 50%;
      background: $bg-color;
      font-size: 12px;
      font-weight: 600;
      margin-right: $spacing-md;

      &.top-3 {
        background: $warning-color;
        color: white;
      }
    }

    .name {
      flex: 1;
      font-size: 14px;
      color: $text-primary;
    }

    .value {
      font-size: 14px;
      font-weight: 600;
      color: $text-primary;
      margin-right: $spacing-md;
    }

    .change {
      font-size: 12px;

      &.up { color: $success-color; }
      &.down { color: $danger-color; }
    }
  }
}

.quick-actions {
  .action-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-md;

    .el-button {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
    }
  }
}

// 响应式调整
@media (max-width: 768px) {
  .chart-card {
    .card-header {
      flex-direction: column;
      align-items: flex-start;
      gap: $spacing-sm;
    }
  }
}
</style>
