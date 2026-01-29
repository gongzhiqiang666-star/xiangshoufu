<template>
  <div class="deduction-detail-view">
    <PageHeader title="代扣详情" :sub-title="`计划编号: ${detail?.plan_no || ''}`">
      <template #extra>
        <el-button @click="handleBack">返回列表</el-button>
        <el-button
          v-if="detail?.status === 1"
          type="warning"
          @click="handlePause"
        >暂停</el-button>
        <el-button
          v-if="detail?.status === 3"
          type="success"
          @click="handleResume"
        >恢复</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 状态卡片 -->
      <el-card class="status-card">
        <div class="status-content">
          <div class="status-icon" :class="getStatusClass(detail?.status)">
            <el-icon v-if="detail?.status === 1"><Loading /></el-icon>
            <el-icon v-else-if="detail?.status === 2"><CircleCheckFilled /></el-icon>
            <el-icon v-else-if="detail?.status === 3"><VideoPause /></el-icon>
            <el-icon v-else><CircleCloseFilled /></el-icon>
          </div>
          <div class="status-info">
            <div class="status-text">{{ getStatusLabel(detail?.status) }}</div>
            <div class="status-amount">
              <span class="label">代扣总额</span>
              <span class="value">¥{{ formatAmount(detail?.total_amount || 0) }}</span>
            </div>
          </div>
          <div class="progress-info">
            <el-progress
              type="dashboard"
              :percentage="progressPercent"
              :width="100"
              :stroke-width="8"
            >
              <template #default>
                <div class="progress-text">
                  <div class="period">{{ detail?.current_period || 0 }}/{{ detail?.total_periods || 0 }}</div>
                  <div class="label">期</div>
                </div>
              </template>
            </el-progress>
          </div>
        </div>
      </el-card>

      <!-- 金额信息 -->
      <el-card class="amount-card">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="amount-item">
              <span class="label">代扣总额</span>
              <span class="value">¥{{ formatAmount(detail?.total_amount || 0) }}</span>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="amount-item success">
              <span class="label">已扣金额</span>
              <span class="value">¥{{ formatAmount(detail?.deducted_amount || 0) }}</span>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="amount-item danger">
              <span class="label">剩余待扣</span>
              <span class="value">¥{{ formatAmount(detail?.remaining_amount || 0) }}</span>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 基本信息 -->
      <el-card class="detail-card">
        <template #header>基本信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="计划编号">{{ detail?.plan_no }}</el-descriptions-item>
          <el-descriptions-item label="计划类型">
            <el-tag :type="getTypeTag(detail?.plan_type)" size="small">
              {{ getTypeLabel(detail?.plan_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusTag(detail?.status)" size="small">
              {{ getStatusLabel(detail?.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="扣款方">{{ detail?.deductor_name }}</el-descriptions-item>
          <el-descriptions-item label="被扣款方">{{ detail?.deductee_name }}</el-descriptions-item>
          <el-descriptions-item label="每期金额">¥{{ formatAmount(detail?.period_amount || 0) }}</el-descriptions-item>
          <el-descriptions-item label="总期数">{{ detail?.total_periods }} 期</el-descriptions-item>
          <el-descriptions-item label="当前期数">第 {{ detail?.current_period }} 期</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ detail?.created_at }}</el-descriptions-item>
          <el-descriptions-item v-if="detail?.completed_at" label="完成时间">
            {{ detail?.completed_at }}
          </el-descriptions-item>
          <el-descriptions-item v-if="detail?.remark" label="备注" :span="2">
            {{ detail?.remark }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 扣款记录 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>扣款记录</span>
            <el-tag type="info" size="small">共 {{ records.length }} 条</el-tag>
          </div>
        </template>
        <el-table :data="records" border stripe>
          <el-table-column prop="period_num" label="期数" width="80" align="center">
            <template #default="{ row }">
              第 {{ row.period_num }} 期
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="应扣金额" width="120" align="right">
            <template #default="{ row }">
              ¥{{ formatAmount(row.amount) }}
            </template>
          </el-table-column>
          <el-table-column prop="actual_amount" label="实扣金额" width="120" align="right">
            <template #default="{ row }">
              <span :class="{ 'success-text': row.actual_amount > 0 }">
                ¥{{ formatAmount(row.actual_amount) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getRecordStatusTag(row.status)" size="small">
                {{ getRecordStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="scheduled_at" label="计划扣款时间" width="170" />
          <el-table-column prop="deducted_at" label="实际扣款时间" width="170">
            <template #default="{ row }">
              {{ row.deducted_at || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="fail_reason" label="失败原因" min-width="150">
            <template #default="{ row }">
              <span v-if="row.fail_reason" class="danger-text">{{ row.fail_reason }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  CircleCheckFilled,
  CircleCloseFilled,
  Loading,
  VideoPause,
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import {
  getDeductionPlanDetail,
  getDeductionRecords,
  pauseDeductionPlan,
  resumeDeductionPlan,
} from '@/api/deduction'
import { formatAmount } from '@/utils/format'
import type {
  DeductionPlanDetail,
  DeductionRecord,
  DeductionPlanStatus,
  DeductionPlanType,
  DeductionRecordStatus,
} from '@/types'
import {
  DEDUCTION_PLAN_STATUS_CONFIG,
  DEDUCTION_PLAN_TYPE_CONFIG,
  DEDUCTION_RECORD_STATUS_CONFIG,
} from '@/types/deduction'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<DeductionPlanDetail | null>(null)
const records = ref<DeductionRecord[]>([])

// 进度百分比
const progressPercent = computed(() => {
  if (!detail.value || !detail.value.total_amount) return 0
  return Math.round((detail.value.deducted_amount / detail.value.total_amount) * 100)
})

// 状态样式
function getStatusClass(status?: DeductionPlanStatus) {
  const map: Record<number, string> = {
    1: 'active',
    2: 'completed',
    3: 'paused',
    4: 'cancelled',
  }
  return status ? map[status] : ''
}

function getStatusLabel(status?: DeductionPlanStatus) {
  return status ? DEDUCTION_PLAN_STATUS_CONFIG[status]?.label : '未知'
}

function getStatusTag(status?: DeductionPlanStatus) {
  return status ? DEDUCTION_PLAN_STATUS_CONFIG[status]?.type : 'info'
}

// 类型标签
function getTypeTag(type?: DeductionPlanType) {
  if (!type) return 'info'
  const config = DEDUCTION_PLAN_TYPE_CONFIG[type]
  return config?.color === '#409eff' ? 'primary' : config?.color === '#67c23a' ? 'success' : 'warning'
}

function getTypeLabel(type?: DeductionPlanType) {
  return type ? DEDUCTION_PLAN_TYPE_CONFIG[type]?.label : '未知'
}

// 记录状态
function getRecordStatusTag(status: DeductionRecordStatus) {
  return DEDUCTION_RECORD_STATUS_CONFIG[status]?.type || 'info'
}

function getRecordStatusLabel(status: DeductionRecordStatus) {
  return DEDUCTION_RECORD_STATUS_CONFIG[status]?.label || '未知'
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    detail.value = await getDeductionPlanDetail(id)
    // 获取扣款记录
    const recordRes = await getDeductionRecords(id, { page_size: 100 })
    records.value = recordRes.list
  } catch (error) {
    console.error('Fetch detail error:', error)
    ElMessage.error('获取详情失败')
  } finally {
    loading.value = false
  }
}

// 返回
function handleBack() {
  router.push('/deductions/list')
}

// 暂停
async function handlePause() {
  if (!detail.value) return
  try {
    await ElMessageBox.confirm('确定要暂停此代扣计划吗？', '暂停确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await pauseDeductionPlan(detail.value.id)
    ElMessage.success('暂停成功')
    fetchDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Pause error:', error)
    }
  }
}

// 恢复
async function handleResume() {
  if (!detail.value) return
  try {
    await ElMessageBox.confirm('确定要恢复此代扣计划吗？', '恢复确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info',
    })
    await resumeDeductionPlan(detail.value.id)
    ElMessage.success('恢复成功')
    fetchDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Resume error:', error)
    }
  }
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.deduction-detail-view {
  padding: 0;
}

.status-card {
  margin-bottom: $spacing-md;

  .status-content {
    display: flex;
    align-items: center;
    gap: $spacing-xl;

    .status-icon {
      width: 64px;
      height: 64px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 32px;

      &.active {
        background: rgba($primary-color, 0.1);
        color: $primary-color;
      }

      &.completed {
        background: rgba($success-color, 0.1);
        color: $success-color;
      }

      &.paused {
        background: rgba($warning-color, 0.1);
        color: $warning-color;
      }

      &.cancelled {
        background: rgba($text-secondary, 0.1);
        color: $text-secondary;
      }
    }

    .status-info {
      flex: 1;

      .status-text {
        font-size: 18px;
        font-weight: 600;
        margin-bottom: $spacing-xs;
      }

      .status-amount {
        .label {
          color: $text-secondary;
          margin-right: $spacing-sm;
        }

        .value {
          font-size: 24px;
          font-weight: 600;
          color: $primary-color;
        }
      }
    }

    .progress-info {
      .progress-text {
        text-align: center;

        .period {
          font-size: 18px;
          font-weight: 600;
          color: $text-primary;
        }

        .label {
          font-size: 12px;
          color: $text-secondary;
        }
      }
    }
  }
}

.amount-card {
  margin-bottom: $spacing-md;

  .amount-item {
    text-align: center;
    padding: $spacing-md;
    background: $bg-color;
    border-radius: $border-radius-md;

    .label {
      display: block;
      font-size: 14px;
      color: $text-secondary;
      margin-bottom: $spacing-xs;
    }

    .value {
      font-size: 24px;
      font-weight: 600;
      color: $text-primary;
    }

    &.success .value {
      color: $success-color;
    }

    &.danger .value {
      color: $danger-color;
    }
  }
}

.detail-card {
  margin-bottom: $spacing-md;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.success-text {
  color: $success-color;
  font-weight: 600;
}

.danger-text {
  color: $danger-color;
}
</style>
