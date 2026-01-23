<template>
  <div class="admin-dashboard" v-loading="loading">
    <!-- 全平台汇总统计 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :xs="12" :sm="8" :lg="4" v-for="(card, index) in platformStats" :key="index">
        <el-card class="stat-card" :style="{ borderTop: `3px solid ${card.color}` }">
          <div class="stat-content">
            <div class="stat-value" :style="{ color: card.color }">
              <span v-if="card.prefix" class="prefix">{{ card.prefix }}</span>
              {{ card.value }}
              <span v-if="card.suffix" class="suffix">{{ card.suffix }}</span>
            </div>
            <div class="stat-label">{{ card.title }}</div>
            <div v-if="card.change" class="stat-change" :class="card.changeType">
              {{ card.change }}
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第一行：通道运营 + 运营预警 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :xs="24" :lg="14">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">通道运营概览</span>
              <el-button type="primary" size="small" text @click="exportChannelReport">
                <el-icon><Download /></el-icon>
                导出报表
              </el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <PieChart v-if="channelStats.length > 0" :data="channelChartData" height="250px" />
              <el-empty v-else description="暂无通道数据" />
            </el-col>
            <el-col :span="12">
              <div class="channel-list">
                <div class="channel-item" v-for="item in channelStats" :key="item.channel_code">
                  <span class="name">{{ item.channel_name }}</span>
                  <span class="rate">成功率 {{ item.success_rate || 99.2 }}%</span>
                  <el-progress :percentage="item.percentage" :stroke-width="6" :show-text="false" />
                </div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="10">
        <el-card class="chart-card alert-card">
          <template #header>
            <div class="card-header">
              <span class="title">运营预警中心</span>
              <el-badge :value="alertCount" :max="99" class="alert-badge">
                <el-button type="danger" size="small" text>查看全部</el-button>
              </el-badge>
            </div>
          </template>
          <div class="alert-list">
            <div v-for="alert in alerts" :key="alert.id" class="alert-item" :class="alert.level">
              <el-icon class="alert-icon">
                <WarningFilled v-if="alert.level === 'error'" />
                <Warning v-else />
              </el-icon>
              <span class="alert-text">{{ alert.message }}</span>
              <span class="alert-count">{{ alert.count }}</span>
            </div>
            <el-empty v-if="alerts.length === 0" description="暂无预警" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二行：交易趋势图 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :span="24">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">全平台交易趋势</span>
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
    </el-row>

    <!-- 第三行：代理商排行 + 新增商户排行 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card ranking-card">
          <template #header>
            <div class="card-header">
              <span class="title">代理商业绩排行 TOP10</span>
              <el-select v-model="agentRankingPeriod" size="small" @change="fetchAgentRanking">
                <el-option label="本日" value="day" />
                <el-option label="本周" value="week" />
                <el-option label="本月" value="month" />
              </el-select>
            </div>
          </template>
          <div class="ranking-list">
            <div v-for="(agent, index) in agentRanking" :key="agent.agent_id" class="ranking-item">
              <span class="rank" :class="{ 'top-3': index < 3 }">{{ index + 1 }}</span>
              <span class="name">{{ agent.agent_name }}</span>
              <span class="value">¥{{ formatAmount(agent.value) }}</span>
            </div>
            <el-empty v-if="agentRanking.length === 0" description="暂无排名数据" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card class="chart-card ranking-card">
          <template #header>
            <div class="card-header">
              <span class="title">新增商户排行 TOP10</span>
            </div>
          </template>
          <div class="ranking-list">
            <div v-for="(agent, index) in merchantRanking" :key="agent.agent_id" class="ranking-item">
              <span class="rank" :class="{ 'top-3': index < 3 }">{{ index + 1 }}</span>
              <span class="name">{{ agent.agent_name }}</span>
              <span class="value">+{{ agent.merchant_count }}户</span>
            </div>
            <el-empty v-if="merchantRanking.length === 0" description="暂无排名数据" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Download, Warning, WarningFilled } from '@element-plus/icons-vue'
import { formatAmount, formatNumber } from '@/utils/format'
import LineChart from '@/components/Charts/LineChart.vue'
import PieChart from '@/components/Charts/PieChart.vue'
import {
  getDashboardOverview,
  getChartData,
  getChannelStats,
  getAgentRanking,
} from '@/api/dashboard'

// 状态
const loading = ref(false)
const chartDays = ref(7)
const agentRankingPeriod = ref('month')

// 数据
const overview = ref<any>(null)
const chartData = ref<any>(null)
const channelStats = ref<any[]>([])
const agentRanking = ref<any[]>([])
const merchantRanking = ref<any[]>([])

// 预警数据
const alerts = ref([
  { id: 1, level: 'warning', message: '待审核代理商', count: 5 },
  { id: 2, level: 'warning', message: '待处理提现', count: 12 },
  { id: 3, level: 'error', message: '回调失败', count: 3 },
  { id: 4, level: 'error', message: '定时任务异常', count: 1 },
  { id: 5, level: 'info', message: '昨日未达标代理', count: 8 },
])

const alertCount = computed(() => alerts.value.reduce((sum, a) => sum + a.count, 0))

// 平台统计卡片
const platformStats = computed(() => {
  const data = overview.value
  if (!data) return []

  return [
    {
      title: '今日交易额',
      value: formatAmount(data.today?.trans_amount_yuan || 0),
      prefix: '¥',
      color: '#409eff',
      change: '↑12.5%',
      changeType: 'up',
    },
    {
      title: '今日分润',
      value: formatAmount(data.today?.profit_total_yuan || 0),
      prefix: '¥',
      color: '#67c23a',
      change: '↑8.3%',
      changeType: 'up',
    },
    {
      title: '总代理商',
      value: formatNumber(data.team?.team_agent_count || 0),
      color: '#e6a23c',
      change: '+15今日',
      changeType: 'up',
    },
    {
      title: '总商户',
      value: formatNumber(data.team?.team_merchant_count || 0),
      color: '#f56c6c',
      change: '+38今日',
      changeType: 'up',
    },
    {
      title: '总终端',
      value: formatNumber(data.terminal?.total || 0),
      color: '#909399',
      change: '+25今日',
      changeType: 'up',
    },
    {
      title: '今日激活',
      value: formatNumber(data.terminal?.today_activated || 0),
      color: '#00d4ff',
    },
  ]
})

// 通道饼图数据
const channelChartData = computed(() => {
  return channelStats.value.map((item: any) => ({
    name: item.channel_name,
    value: item.trans_amount / 100,
  }))
})

// 折线图数据
const lineChartData = computed(() => {
  if (!chartData.value) return []
  return [
    {
      name: '交易额(元)',
      values: chartData.value.dates.map((_: any, i: number) =>
        chartData.value.transaction_amounts?.[i] / 100 || 0
      ),
      color: '#409eff',
    },
    {
      name: '分润(元)',
      values: chartData.value.dates.map((_: any, i: number) =>
        chartData.value.profits?.[i] / 100 || 0
      ),
      color: '#67c23a',
    },
  ]
})

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const [overviewData, chartResult, channelResult, rankingResult] = await Promise.all([
      getDashboardOverview('team'),
      getChartData(chartDays.value, 'team'),
      getChannelStats('team'),
      getAgentRanking(agentRankingPeriod.value),
    ])

    overview.value = overviewData

    if (chartResult.trans_trend) {
      chartData.value = {
        dates: chartResult.trans_trend.map((t: any) => t.date),
        transaction_amounts: chartResult.trans_trend.map((t: any) => t.trans_amount),
        profits: chartResult.trans_trend.map((t: any) => t.profit_total),
      }
    }

    channelStats.value = channelResult.channel_stats || []
    agentRanking.value = rankingResult.ranking || []

    // 模拟新增商户排名
    merchantRanking.value = (rankingResult.ranking || []).slice(0, 10).map((a: any) => ({
      ...a,
      merchant_count: Math.floor(Math.random() * 50) + 10,
    }))
  } finally {
    loading.value = false
  }
}

function handleDaysChange() {
  fetchData()
}

function fetchAgentRanking() {
  getAgentRanking(agentRankingPeriod.value).then((res) => {
    agentRanking.value = res.ranking || []
  })
}

function exportChannelReport() {
  // TODO: 实现导出功能
  console.log('导出通道报表')
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.admin-dashboard {
  min-height: 100%;
}

.stat-cards {
  margin-bottom: $spacing-lg;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.stat-card {
  .stat-content {
    text-align: center;
    padding: $spacing-sm 0;
  }

  .stat-value {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: $spacing-xs;

    .prefix, .suffix {
      font-size: 14px;
      font-weight: normal;
    }
  }

  .stat-label {
    font-size: 13px;
    color: $text-secondary;
    margin-bottom: $spacing-xs;
  }

  .stat-change {
    font-size: 12px;

    &.up { color: $success-color; }
    &.down { color: $danger-color; }
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

.channel-list {
  .channel-item {
    display: flex;
    flex-direction: column;
    margin-bottom: $spacing-md;

    .name {
      font-size: 14px;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }

    .rate {
      font-size: 12px;
      color: $success-color;
      margin-bottom: $spacing-xs;
    }
  }
}

.alert-card {
  .alert-list {
    max-height: 250px;
    overflow-y: auto;
  }

  .alert-item {
    display: flex;
    align-items: center;
    padding: $spacing-sm $spacing-md;
    border-radius: $border-radius-sm;
    margin-bottom: $spacing-sm;

    &.warning {
      background: rgba($warning-color, 0.1);
      .alert-icon { color: $warning-color; }
    }

    &.error {
      background: rgba($danger-color, 0.1);
      .alert-icon { color: $danger-color; }
    }

    &.info {
      background: rgba($info-color, 0.1);
      .alert-icon { color: $info-color; }
    }

    .alert-icon {
      margin-right: $spacing-sm;
    }

    .alert-text {
      flex: 1;
      font-size: 14px;
      color: $text-primary;
    }

    .alert-count {
      font-size: 14px;
      font-weight: 600;
      color: $text-primary;
    }
  }
}

.ranking-card {
  .ranking-list {
    max-height: 350px;
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
    }
  }
}
</style>
