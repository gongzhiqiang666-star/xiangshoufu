<template>
  <div class="transaction-list-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" style="width: 150px" />
      </el-form-item>
      <el-form-item label="交易类型">
        <el-select v-model="searchForm.transaction_type" placeholder="请选择" clearable style="width: 120px">
          <el-option label="贷记卡" value="credit" />
          <el-option label="借记卡" value="debit" />
          <el-option label="微信" value="wechat" />
          <el-option label="支付宝" value="alipay" />
          <el-option label="云闪付" value="unionpay_qr" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择" clearable style="width: 100px">
          <el-option label="成功" value="success" />
          <el-option label="失败" value="failed" />
          <el-option label="处理中" value="pending" />
          <el-option label="已退款" value="refunded" />
        </el-select>
      </el-form-item>
      <el-form-item label="机具号">
        <el-input v-model="searchForm.terminal_sn" placeholder="请输入机具号" clearable style="width: 150px" />
      </el-form-item>
      <el-form-item label="交易时间">
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
      :show-export="true"
      @refresh="fetchData"
      @export="handleExport"
    >
      <el-table-column prop="transaction_no" label="交易号" width="180" />
      <el-table-column prop="channel_name" label="通道" width="100" />
      <el-table-column prop="merchant_name" label="商户" width="100" />
      <el-table-column prop="terminal_sn" label="机具号" width="130" />
      <el-table-column prop="transaction_type" label="交易类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.transaction_type)" size="small">
            {{ getTypeLabel(row.transaction_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="交易金额" width="120" align="right">
        <template #default="{ row }">
          ¥{{ formatAmount(row.amount) }}
        </template>
      </el-table-column>
      <el-table-column prop="fee" label="手续费" width="100" align="right">
        <template #default="{ row }">
          ¥{{ formatAmount(row.fee) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="transaction_time" label="交易时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import { getTransactions } from '@/api/transaction'
import { formatAmount } from '@/utils/format'
import type { Transaction, TransactionType, TransactionStatus } from '@/types'
import { TRANSACTION_TYPE_CONFIG } from '@/types/transaction'

const router = useRouter()

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  transaction_type: undefined as TransactionType | undefined,
  status: undefined as TransactionStatus | undefined,
  terminal_sn: '',
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<Transaction[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 类型配置
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getTypeTag(type: TransactionType): TagType {
  const colorMap: Record<string, TagType> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#07c160': 'success',
    '#1677ff': 'primary',
    '#e60012': 'danger',
  }
  const config = TRANSACTION_TYPE_CONFIG[type]
  return colorMap[config?.color] || 'info'
}

function getTypeLabel(type: TransactionType) {
  return TRANSACTION_TYPE_CONFIG[type]?.label || type
}

function getStatusTag(status: TransactionStatus): TagType {
  const map: Record<TransactionStatus, TagType> = {
    success: 'success',
    failed: 'danger',
    pending: 'warning',
    refunded: 'info',
  }
  return map[status] || 'info'
}

function getStatusLabel(status: TransactionStatus) {
  const map: Record<TransactionStatus, string> = {
    success: '成功',
    failed: '失败',
    pending: '处理中',
    refunded: '已退款',
  }
  return map[status] || status
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
    const res = await getTransactions({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch transactions error:', error)
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
function handleView(row: Transaction) {
  router.push(`/transactions/${row.id}`)
}

// 导出
function handleExport() {
  ElMessage.info('导出功能开发中...')
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.transaction-list-view {
  padding: 0;
}

.summary-card {
  margin-bottom: $spacing-md;

  .summary-item {
    display: flex;
    align-items: center;
    gap: $spacing-sm;

    .label {
      color: $text-secondary;
    }

    .value {
      font-size: 18px;
      font-weight: 600;
      color: $primary-color;
    }
  }
}
</style>
