<template>
  <div class="dashboard-view" v-loading="loading">
    <!-- 统计卡片区域 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col
        v-for="(card, index) in statCards"
        :key="index"
        :xs="24"
        :sm="12"
        :lg="6"
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

    <!-- 图表区域 -->
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

      <!-- 本月数据概览 -->
      <el-col :xs="24" :lg="8">
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
                {{ formatAmount(overview.month_transaction_amount) }}
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-label">本月交易笔数</div>
              <div class="stat-value">
                {{ formatNumber(overview.month_transaction_count) }}
                <span class="suffix">笔</span>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-label">本月分润</div>
              <div class="stat-value">
                <span class="prefix">¥</span>
                {{ formatAmount(overview.month_profit) }}
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-label">本月新增商户</div>
              <div class="stat-value">
                {{ formatNumber(overview.month_new_merchants) }}
                <span class="suffix">户</span>
              </div>
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
import { User, Shop, Money, Wallet, Refresh } from '@element-plus/icons-vue'
import { useDashboardStore } from '@/stores/dashboard'
import { formatAmount, formatNumber } from '@/utils/format'
import StatCard from '@/components/Common/StatCard.vue'
import LineChart from '@/components/Charts/LineChart.vue'

const router = useRouter()
const dashboardStore = useDashboardStore()

// 图表天数
const chartDays = ref(7)

// 加载状态
const loading = computed(() => dashboardStore.loading)
const chartLoading = computed(() => dashboardStore.chartLoading)

// 统计卡片数据
const statCards = computed(() => dashboardStore.statCards)

// 概览数据
const overview = computed(() => dashboardStore.overview)

// 图表数据
const chartData = computed(() => dashboardStore.chartData)

// 折线图数据
const lineChartData = computed(() => {
  if (!chartData.value) return []
  return [
    {
      name: '交易额(元)',
      values: chartData.value.transaction_amounts.map((v) => v / 100),
      color: '#409eff',
    },
    {
      name: '分润(元)',
      values: chartData.value.profits.map((v) => v / 100),
      color: '#67c23a',
    },
  ]
})

// 处理天数切换
function handleDaysChange(days: number) {
  dashboardStore.fetchChartData(days)
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
