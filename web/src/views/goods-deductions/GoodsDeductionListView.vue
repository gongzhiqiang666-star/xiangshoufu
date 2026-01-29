<template>
  <div class="goods-deduction-list-view">
    <!-- Tab切换 -->
    <div class="tab-wrapper">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="我发起的" name="sent" />
        <el-tab-pane label="我接收的" name="received" />
      </el-tabs>
    </div>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择" clearable style="width: 100px">
          <el-option label="待接收" :value="1" />
          <el-option label="进行中" :value="2" />
          <el-option label="已完成" :value="3" />
          <el-option label="已拒绝" :value="4" />
        </el-select>
      </el-form-item>
      <el-form-item label="扣款来源">
        <el-select v-model="searchForm.deduction_source" placeholder="请选择" clearable style="width: 120px">
          <el-option label="分润钱包" :value="1" />
          <el-option label="服务费钱包" :value="2" />
          <el-option label="分润+服务费" :value="3" />
        </el-select>
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
      <el-table-column prop="deduction_no" label="代扣编号" width="180" />
      <el-table-column :label="activeTab === 'sent' ? '接收方' : '发起方'" width="120">
        <template #default="{ row }">
          {{ activeTab === 'sent' ? row.to_agent_name : row.from_agent_name }}
        </template>
      </el-table-column>
      <el-table-column prop="terminal_count" label="终端数量" width="90" align="center" />
      <el-table-column prop="total_amount" label="总金额" width="120" align="right">
        <template #default="{ row }">
          ¥{{ formatAmount(row.total_amount) }}
        </template>
      </el-table-column>
      <el-table-column prop="deducted_amount" label="已扣金额" width="120" align="right">
        <template #default="{ row }">
          <span class="success-text">¥{{ formatAmount(row.deducted_amount) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="remaining_amount" label="剩余金额" width="120" align="right">
        <template #default="{ row }">
          <span class="danger-text">¥{{ formatAmount(row.remaining_amount) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="progress" label="进度" width="100" align="center">
        <template #default="{ row }">
          <el-progress
            :percentage="row.progress"
            :stroke-width="6"
            :show-text="false"
            :status="row.progress >= 100 ? 'success' : ''"
          />
          <span class="progress-text">{{ row.progress.toFixed(1) }}%</span>
        </template>
      </el-table-column>
      <el-table-column prop="source_name" label="扣款来源" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getSourceTag(row.deduction_source)" size="small">
            {{ row.source_name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)" size="small">
            {{ row.status_name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <template v-if="activeTab === 'received' && row.status === 1">
          <el-button type="success" link @click="handleAccept(row)">接收</el-button>
          <el-button type="danger" link @click="handleReject(row)">拒绝</el-button>
        </template>
      </template>
    </ProTable>

    <!-- 拒绝弹窗 -->
    <el-dialog
      v-model="rejectDialogVisible"
      title="拒绝货款代扣"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form ref="rejectFormRef" :model="rejectForm" :rules="rejectRules">
        <el-form-item label="拒绝原因" prop="reason">
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="rejecting" @click="confirmReject">
          确认拒绝
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import {
  getSentGoodsDeductions,
  getReceivedGoodsDeductions,
  acceptGoodsDeduction,
  rejectGoodsDeduction,
} from '@/api/goodsDeduction'
import { formatAmount } from '@/utils/format'
import type {
  GoodsDeduction,
  GoodsDeductionStatus,
  DeductionSource,
} from '@/types'
import { GOODS_DEDUCTION_STATUS_CONFIG, DEDUCTION_SOURCE_CONFIG } from '@/types/deduction'

const router = useRouter()

// Tab状态
const activeTab = ref<'sent' | 'received'>('sent')

// 搜索表单
const searchForm = reactive({
  status: undefined as GoodsDeductionStatus | undefined,
  deduction_source: undefined as DeductionSource | undefined,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<GoodsDeduction[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 拒绝弹窗
const rejectDialogVisible = ref(false)
const rejectFormRef = ref<FormInstance>()
const rejecting = ref(false)
const currentRejectId = ref<number>(0)
const rejectForm = reactive({
  reason: '',
})
const rejectRules: FormRules = {
  reason: [{ required: true, message: '请输入拒绝原因', trigger: 'blur' }],
}

// 状态标签
function getStatusTag(status: GoodsDeductionStatus) {
  return GOODS_DEDUCTION_STATUS_CONFIG[status]?.type || 'info'
}

// 来源标签
function getSourceTag(source: DeductionSource) {
  const config = DEDUCTION_SOURCE_CONFIG[source]
  return config?.color === '#409eff' ? 'primary' : config?.color === '#67c23a' ? 'success' : 'warning'
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

// Tab切换
function handleTabChange() {
  page.value = 1
  fetchData()
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const fetchFn = activeTab.value === 'sent' ? getSentGoodsDeductions : getReceivedGoodsDeductions
    const res = await fetchFn({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch goods deductions error:', error)
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
function handleView(row: GoodsDeduction) {
  router.push(`/goods-deductions/${row.id}`)
}

// 接收
async function handleAccept(row: GoodsDeduction) {
  try {
    await ElMessageBox.confirm(
      `确定要接收货款代扣 ${row.deduction_no} 吗？接收后将开始自动扣款。`,
      '接收确认',
      {
        confirmButtonText: '确定接收',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    await acceptGoodsDeduction(row.id)
    ElMessage.success('接收成功，代扣已开始')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Accept error:', error)
    }
  }
}

// 打开拒绝弹窗
function handleReject(row: GoodsDeduction) {
  currentRejectId.value = row.id
  rejectForm.reason = ''
  rejectDialogVisible.value = true
}

// 确认拒绝
async function confirmReject() {
  if (!rejectFormRef.value) return

  try {
    await rejectFormRef.value.validate()
    rejecting.value = true
    await rejectGoodsDeduction(currentRejectId.value, { reason: rejectForm.reason })
    ElMessage.success('已拒绝')
    rejectDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Reject error:', error)
  } finally {
    rejecting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.goods-deduction-list-view {
  padding: 0;
}

.tab-wrapper {
  margin-bottom: $spacing-md;
  background: #fff;
  padding: 0 16px;
  border-radius: 4px;
  border: 1px solid #ebeef5;

  :deep(.el-tabs__header) {
    margin-bottom: 0;
  }
}

.progress-text {
  display: block;
  font-size: 12px;
  color: $text-secondary;
  text-align: center;
  margin-top: 2px;
}

.success-text {
  color: $success-color;
  font-weight: 600;
}

.danger-text {
  color: $danger-color;
  font-weight: 600;
}
</style>
