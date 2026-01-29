<template>
  <div class="deduction-list-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="计划类型">
        <el-select v-model="searchForm.plan_type" placeholder="请选择" clearable style="width: 120px">
          <el-option label="货款代扣" :value="1" />
          <el-option label="伙伴代扣" :value="2" />
          <el-option label="押金代扣" :value="3" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择" clearable style="width: 100px">
          <el-option label="进行中" :value="1" />
          <el-option label="已完成" :value="2" />
          <el-option label="已暂停" :value="3" />
          <el-option label="已取消" :value="4" />
        </el-select>
      </el-form-item>
      <el-form-item label="被扣方">
        <AgentSelect v-model="searchForm.deductee_id" placeholder="选择被扣款代理商" style="width: 150px" />
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
      <template #extra>
        <el-button type="primary" :icon="Plus" @click="handleCreate">发起代扣</el-button>
      </template>
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
      <el-table-column prop="plan_no" label="计划编号" width="180" />
      <el-table-column prop="plan_type" label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.plan_type)" size="small">
            {{ getTypeLabel(row.plan_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="deductor_name" label="扣款方" width="120" />
      <el-table-column prop="deductee_name" label="被扣方" width="120" />
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
      <el-table-column prop="period_progress" label="期数进度" width="100" align="center">
        <template #default="{ row }">
          {{ row.current_period }}/{{ row.total_periods }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170" />

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button
          v-if="row.status === 1"
          type="warning"
          link
          @click="handlePause(row)"
        >暂停</el-button>
        <el-button
          v-if="row.status === 3"
          type="success"
          link
          @click="handleResume(row)"
        >恢复</el-button>
        <el-button
          v-if="row.status === 1 || row.status === 3"
          type="danger"
          link
          @click="handleCancel(row)"
        >取消</el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import {
  getDeductionPlans,
  pauseDeductionPlan,
  resumeDeductionPlan,
  cancelDeductionPlan,
} from '@/api/deduction'
import { formatAmount } from '@/utils/format'
import type { DeductionPlan, DeductionPlanStatus, DeductionPlanType } from '@/types'
import {
  DEDUCTION_PLAN_STATUS_CONFIG,
  DEDUCTION_PLAN_TYPE_CONFIG,
} from '@/types/deduction'

const router = useRouter()

// 搜索表单
const searchForm = reactive({
  plan_type: undefined as DeductionPlanType | undefined,
  status: undefined as DeductionPlanStatus | undefined,
  deductee_id: undefined as number | undefined,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<DeductionPlan[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 类型标签
function getTypeTag(type: DeductionPlanType) {
  const config = DEDUCTION_PLAN_TYPE_CONFIG[type]
  return config?.color === '#409eff' ? 'primary' : config?.color === '#67c23a' ? 'success' : 'warning'
}

function getTypeLabel(type: DeductionPlanType) {
  return DEDUCTION_PLAN_TYPE_CONFIG[type]?.label || '未知'
}

// 状态标签
function getStatusTag(status: DeductionPlanStatus) {
  const config = DEDUCTION_PLAN_STATUS_CONFIG[status]
  return config?.type || 'info'
}

function getStatusLabel(status: DeductionPlanStatus) {
  return DEDUCTION_PLAN_STATUS_CONFIG[status]?.label || '未知'
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
    const res = await getDeductionPlans({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch deduction plans error:', error)
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

// 创建代扣
function handleCreate() {
  router.push('/deductions/create')
}

// 查看详情
function handleView(row: DeductionPlan) {
  router.push(`/deductions/${row.id}`)
}

// 暂停
async function handlePause(row: DeductionPlan) {
  try {
    await ElMessageBox.confirm(`确定要暂停代扣计划 ${row.plan_no} 吗？`, '暂停确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await pauseDeductionPlan(row.id)
    ElMessage.success('暂停成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Pause error:', error)
    }
  }
}

// 恢复
async function handleResume(row: DeductionPlan) {
  try {
    await ElMessageBox.confirm(`确定要恢复代扣计划 ${row.plan_no} 吗？`, '恢复确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info',
    })
    await resumeDeductionPlan(row.id)
    ElMessage.success('恢复成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Resume error:', error)
    }
  }
}

// 取消
async function handleCancel(row: DeductionPlan) {
  try {
    await ElMessageBox.confirm(
      `确定要取消代扣计划 ${row.plan_no} 吗？取消后不可恢复。`,
      '取消确认',
      {
        confirmButtonText: '确定取消',
        cancelButtonText: '返回',
        type: 'warning',
      }
    )
    await cancelDeductionPlan(row.id)
    ElMessage.success('已取消')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Cancel error:', error)
    }
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.deduction-list-view {
  padding: 0;
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
