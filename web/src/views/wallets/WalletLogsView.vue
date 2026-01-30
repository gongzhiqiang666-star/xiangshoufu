<template>
  <div class="wallet-logs-view">
    <!-- 钱包信息 -->
    <el-card class="wallet-info-card">
      <el-descriptions :column="5" border>
        <el-descriptions-item label="钱包类型">
          <el-tag :type="getWalletTypeTag(wallet?.wallet_type)" size="small">
            {{ getWalletTypeLabel(wallet?.wallet_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="所属通道">{{ wallet?.channel_name }}</el-descriptions-item>
        <el-descriptions-item label="余额">
          <span class="balance">¥{{ formatAmount(wallet?.balance) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="冻结金额">
          <span class="frozen">¥{{ formatAmount(wallet?.frozen) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="可用金额">
          <span class="available">¥{{ formatAmount(wallet?.available) }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="流水类型">
        <el-select v-model="searchForm.log_type" placeholder="请选择类型" clearable style="width: 120px">
          <el-option label="收入" value="income" />
          <el-option label="支出" value="expense" />
          <el-option label="冻结" value="freeze" />
          <el-option label="解冻" value="unfreeze" />
        </el-select>
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

    <!-- 表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      @refresh="fetchData"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="log_no" label="流水号" width="200" />
      <el-table-column prop="log_type" label="类型" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="getLogTypeTag(row.log_type)" size="small">
            {{ getLogTypeLabel(row.log_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="120" align="right">
        <template #default="{ row }">
          <span :class="['log-amount', row.log_type === 'income' ? 'income' : 'expense']">
            {{ row.log_type === 'income' ? '+' : '-' }}¥{{ formatAmount(Math.abs(row.amount)) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="balance_before" label="变动前余额" width="120" align="right">
        <template #default="{ row }">
          ¥{{ formatAmount(row.balance_before) }}
        </template>
      </el-table-column>
      <el-table-column prop="balance_after" label="变动后余额" width="120" align="right">
        <template #default="{ row }">
          ¥{{ formatAmount(row.balance_after) }}
        </template>
      </el-table-column>
      <el-table-column prop="related_type" label="关联类型" width="100" />
      <el-table-column prop="related_no" label="关联单号" width="200" />
      <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="170" />
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import { getWallet, getWalletLogs } from '@/api/wallet'
import { formatAmount } from '@/utils/format'
import type { WalletType } from '@/types'
import { WALLET_TYPE_CONFIG } from '@/types/wallet'

const route = useRoute()

// 钱包信息
const wallet = ref<any>(null)

// 搜索表单
const searchForm = reactive({
  log_type: undefined as string | undefined,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 获取钱包类型标签
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getWalletTypeTag(type?: WalletType): TagType {
  if (!type) return 'info'
  const colorMap: Record<string, TagType> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
    '#909399': 'info',
  }
  const config = WALLET_TYPE_CONFIG[type]
  return colorMap[config?.color] || 'info'
}

// 获取钱包类型名称
function getWalletTypeLabel(type?: WalletType) {
  if (!type) return '-'
  return WALLET_TYPE_CONFIG[type]?.label || type
}

// 获取流水类型标签
function getLogTypeTag(type: string): TagType {
  const tagMap: Record<string, TagType> = {
    income: 'success',
    expense: 'danger',
    freeze: 'warning',
    unfreeze: 'info',
  }
  return tagMap[type] || 'info'
}

// 获取流水类型名称
function getLogTypeLabel(type: string) {
  const labelMap: Record<string, string> = {
    income: '收入',
    expense: '支出',
    freeze: '冻结',
    unfreeze: '解冻',
  }
  return labelMap[type] || type
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

// 获取钱包信息
async function fetchWallet() {
  try {
    wallet.value = await getWallet(Number(route.params.id))
  } catch (error) {
    console.error('Fetch wallet error:', error)
  }
}

// 获取流水数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getWalletLogs(Number(route.params.id), {
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch wallet logs error:', error)
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

onMounted(() => {
  fetchWallet()
  fetchData()
})
</script>

<style lang="scss" scoped>
.wallet-logs-view {
  padding: 0;
}

.wallet-info-card {
  margin-bottom: $spacing-md;

  .balance {
    font-weight: 600;
    color: $text-primary;
  }

  .frozen {
    color: $warning-color;
  }

  .available {
    color: $success-color;
    font-weight: 600;
  }
}

.log-amount {
  font-weight: 600;

  &.income {
    color: $success-color;
  }

  &.expense {
    color: $danger-color;
  }
}
</style>
