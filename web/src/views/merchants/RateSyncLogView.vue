<template>
  <div class="rate-sync-log-view">
    <PageHeader title="费率同步日志" sub-title="查看商户费率修改同步记录">
      <template #extra>
        <el-button type="primary" @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </template>
    </PageHeader>

    <!-- 搜索筛选 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="商户编号">
          <el-input
            v-model="queryParams.merchant_no"
            placeholder="请输入商户编号"
            clearable
            style="width: 180px"
          />
        </el-form-item>
        <el-form-item label="同步状态">
          <el-select v-model="queryParams.sync_status" placeholder="全部" clearable style="width: 120px">
            <el-option label="待同步" :value="0" />
            <el-option label="同步中" :value="1" />
            <el-option label="同步成功" :value="2" />
            <el-option label="同步失败" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="merchant_no" label="商户编号" width="150" />
        <el-table-column prop="terminal_sn" label="终端SN" width="140" />
        <el-table-column prop="channel_code" label="通道" width="120" />
        <el-table-column label="原费率" width="140">
          <template #default="{ row }">
            <div>贷记: {{ formatRatePercent(row.old_credit_rate) }}</div>
            <div>借记: {{ formatRatePercent(row.old_debit_rate) }}</div>
          </template>
        </el-table-column>
        <el-table-column label="新费率" width="140">
          <template #default="{ row }">
            <div class="new-rate">贷记: {{ formatRatePercent(row.new_credit_rate) }}</div>
            <div class="new-rate">借记: {{ formatRatePercent(row.new_debit_rate) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="sync_status_name" label="同步状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getSyncStatusType(row.sync_status)" size="small">
              {{ row.sync_status_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error_message" label="失败原因" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column prop="synced_at" label="同步时间" width="170" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleViewDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="同步日志详情" width="600px">
      <el-descriptions :column="2" border v-if="currentLog">
        <el-descriptions-item label="ID">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item label="商户编号">{{ currentLog.merchant_no }}</el-descriptions-item>
        <el-descriptions-item label="终端SN">{{ currentLog.terminal_sn || '-' }}</el-descriptions-item>
        <el-descriptions-item label="通道">{{ currentLog.channel_code }}</el-descriptions-item>
        <el-descriptions-item label="原贷记卡费率">{{ formatRatePercent(currentLog.old_credit_rate) }}</el-descriptions-item>
        <el-descriptions-item label="新贷记卡费率">{{ formatRatePercent(currentLog.new_credit_rate) }}</el-descriptions-item>
        <el-descriptions-item label="原借记卡费率">{{ formatRatePercent(currentLog.old_debit_rate) }}</el-descriptions-item>
        <el-descriptions-item label="新借记卡费率">{{ formatRatePercent(currentLog.new_debit_rate) }}</el-descriptions-item>
        <el-descriptions-item label="同步状态">
          <el-tag :type="getSyncStatusType(currentLog.sync_status)" size="small">
            {{ currentLog.sync_status_name }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="通道流水号">{{ currentLog.channel_trade_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item label="同步时间" :span="2">{{ currentLog.synced_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="失败原因" :span="2" v-if="currentLog.error_message">
          <span class="error-message">{{ currentLog.error_message }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getRateSyncLogs, type RateSyncLog } from '@/api/rateSync'

// 查询参数
const queryParams = reactive({
  merchant_no: '',
  sync_status: undefined as number | undefined,
  page: 1,
  page_size: 20,
})

// 数据
const loading = ref(false)
const tableData = ref<RateSyncLog[]>([])
const total = ref(0)

// 详情弹窗
const detailVisible = ref(false)
const currentLog = ref<RateSyncLog | null>(null)

// 格式化费率为百分比
function formatRatePercent(rate: number | undefined): string {
  if (rate === undefined || rate === null) return '-'
  return `${(rate * 100).toFixed(2)}%`
}

// 获取同步状态类型
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getSyncStatusType(status: number): TagType {
  switch (status) {
    case 0: return 'info'
    case 1: return 'warning'
    case 2: return 'success'
    case 3: return 'danger'
    default: return 'info'
  }
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getRateSyncLogs({
      ...queryParams,
      sync_status: queryParams.sync_status,
    })
    tableData.value = res.items || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Fetch rate sync logs error:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  queryParams.page = 1
  fetchData()
}

// 重置
function handleReset() {
  queryParams.merchant_no = ''
  queryParams.sync_status = undefined
  queryParams.page = 1
  fetchData()
}

// 刷新
function handleRefresh() {
  fetchData()
}

// 查看详情
function handleViewDetail(row: RateSyncLog) {
  currentLog.value = row
  detailVisible.value = true
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.rate-sync-log-view {
  padding: 0;
}

.filter-card {
  margin-top: 16px;
}

.table-card {
  margin-top: 16px;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.new-rate {
  color: #409eff;
  font-weight: 500;
}

.error-message {
  color: #f56c6c;
}
</style>
