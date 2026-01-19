<template>
  <div class="profit-list-view">
    <PageHeader title="分润管理" sub-title="分润明细">
      <template #extra>
        <el-button :icon="Download" @click="handleExport">导出Excel</el-button>
      </template>
    </PageHeader>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" />
      </el-form-item>
      <el-form-item label="分润类型">
        <el-select v-model="searchForm.profit_type" placeholder="请选择" clearable>
          <el-option label="交易分润" value="transaction" />
          <el-option label="押金返现" value="deposit_cashback" />
          <el-option label="流量返现" value="sim_cashback" />
          <el-option label="激活奖励" value="activation_reward" />
        </el-select>
      </el-form-item>
      <el-form-item label="代理商">
        <AgentSelect v-model="searchForm.agent_id" />
      </el-form-item>
      <el-form-item label="日期范围">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
        />
      </el-form-item>
    </SearchForm>

    <!-- 统计汇总 -->
    <el-card class="summary-card">
      <el-row :gutter="20">
        <el-col :span="4">
          <div class="summary-item">
            <span class="label">交易分润</span>
            <span class="value primary">¥{{ formatAmount(stats.transaction_profit) }}</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="summary-item">
            <span class="label">押金返现</span>
            <span class="value success">¥{{ formatAmount(stats.deposit_cashback) }}</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="summary-item">
            <span class="label">流量返现</span>
            <span class="value warning">¥{{ formatAmount(stats.sim_cashback) }}</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="summary-item">
            <span class="label">激活奖励</span>
            <span class="value danger">¥{{ formatAmount(stats.activation_reward) }}</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="summary-item total">
            <span class="label">总计</span>
            <span class="value">¥{{ formatAmount(stats.total) }}</span>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      :show-export="true"
      @refresh="fetchData"
      @export="handleExport"
    >
      <el-table-column prop="profit_no" label="分润编号" width="180" />
      <el-table-column prop="related_no" label="关联单号" width="180" />
      <el-table-column prop="agent_name" label="代理商" width="100" />
      <el-table-column prop="channel_name" label="通道" width="100" />
      <el-table-column prop="profit_type" label="分润类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.profit_type)" size="small">
            {{ getTypeLabel(row.profit_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="transaction_amount" label="交易金额" width="120" align="right">
        <template #default="{ row }">
          ¥{{ formatAmount(row.transaction_amount) }}
        </template>
      </el-table-column>
      <el-table-column prop="profit_amount" label="分润金额" width="120" align="right">
        <template #default="{ row }">
          <span class="profit-amount">¥{{ formatAmount(row.profit_amount) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="wallet_type" label="入账钱包" width="100" />
      <el-table-column prop="created_at" label="时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import { getProfits, getProfitStats } from '@/api/profit'
import { formatAmount } from '@/utils/format'
import type { Profit, ProfitType, ProfitStats } from '@/types'
import { PROFIT_TYPE_CONFIG } from '@/types/profit'

const router = useRouter()

// 统计数据
const stats = ref<ProfitStats>({
  transaction_profit: 0,
  deposit_cashback: 0,
  sim_cashback: 0,
  activation_reward: 0,
  total: 0,
})

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  profit_type: undefined as ProfitType | undefined,
  agent_id: undefined as number | undefined,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<Profit[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 类型配置
function getTypeTag(type: ProfitType) {
  const colorMap: Record<string, string> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
  }
  const config = PROFIT_TYPE_CONFIG[type]
  return colorMap[config?.color] || ''
}

function getTypeLabel(type: ProfitType) {
  return PROFIT_TYPE_CONFIG[type]?.label || type
}

// 处理日期变化
function handleDateChange(val: [string, string] | null) {
  if (val) {
    searchForm.start_date = val[0]
    searchForm.end_date = val[1]
  } else {
    searchForm.start_date = ''
    searchForm.end_date = ''
  }
}

// 获取统计数据
async function fetchStats() {
  try {
    stats.value = await getProfitStats({
      channel_id: searchForm.channel_id,
      start_date: searchForm.start_date,
      end_date: searchForm.end_date,
    })
  } catch (error) {
    console.error('Fetch profit stats error:', error)
  }
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getProfits({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch profits error:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  page.value = 1
  fetchData()
  fetchStats()
}

// 重置
function handleReset() {
  dateRange.value = null
  page.value = 1
  fetchData()
  fetchStats()
}

// 查看详情
function handleView(row: Profit) {
  router.push(`/profits/${row.id}`)
}

// 导出
function handleExport() {
  ElMessage.info('导出功能开发中...')
}

onMounted(() => {
  fetchStats()
  fetchData()
})
</script>

<style lang="scss" scoped>
.profit-list-view {
  padding: 0;
}

.summary-card {
  margin-bottom: $spacing-md;

  .summary-item {
    text-align: center;

    .label {
      display: block;
      font-size: 12px;
      color: $text-secondary;
      margin-bottom: $spacing-xs;
    }

    .value {
      font-size: 20px;
      font-weight: 600;

      &.primary {
        color: $primary-color;
      }

      &.success {
        color: $success-color;
      }

      &.warning {
        color: $warning-color;
      }

      &.danger {
        color: $danger-color;
      }
    }

    &.total {
      background: $bg-color;
      padding: $spacing-sm;
      border-radius: $border-radius-sm;

      .value {
        color: $text-primary;
      }
    }
  }
}

.profit-amount {
  color: $success-color;
  font-weight: 600;
}
</style>
