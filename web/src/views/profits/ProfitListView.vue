<template>
  <div class="profit-list-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" style="width: 150px" />
      </el-form-item>
      <el-form-item label="分润类型">
        <el-select v-model="searchForm.profit_type" placeholder="请选择" clearable style="width: 120px">
          <el-option label="交易分润" value="transaction" />
          <el-option label="押金返现" value="deposit_cashback" />
          <el-option label="流量返现" value="sim_cashback" />
          <el-option label="激活奖励" value="activation_reward" />
        </el-select>
      </el-form-item>
      <el-form-item label="代理商">
        <AgentSelect v-model="searchForm.agent_id" style="width: 150px" />
      </el-form-item>
      <el-form-item label="日期范围">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始"
          end-placeholder="结束"
          value-format="YYYY-MM-DD"
          style="width: 240px"
          @change="handleDateChange"
        />
      </el-form-item>
    </SearchForm>

    <!-- 表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      @refresh="fetchData"
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
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import { getProfits } from '@/api/profit'
import { formatAmount } from '@/utils/format'
import type { Profit, ProfitType } from '@/types'
import { PROFIT_TYPE_CONFIG } from '@/types/profit'

const router = useRouter()

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
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getTypeTag(type: ProfitType): TagType {
  const colorMap: Record<string, TagType> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
  }
  const config = PROFIT_TYPE_CONFIG[type]
  return colorMap[config?.color] || 'info'
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
}

// 重置
function handleReset() {
  dateRange.value = null
  page.value = 1
  fetchData()
}

// 查看详情
function handleView(row: Profit) {
  router.push(`/profits/${row.id}`)
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.profit-list-view {
  padding: 0;
}

.profit-amount {
  color: $success-color;
  font-weight: 600;
}
</style>
