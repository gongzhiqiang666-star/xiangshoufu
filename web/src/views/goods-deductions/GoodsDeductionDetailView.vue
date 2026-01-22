<template>
  <div class="goods-deduction-detail-view">
    <PageHeader title="货款代扣详情" :sub-title="`代扣编号: ${detail?.deduction_no || ''}`">
      <template #extra>
        <el-button @click="handleBack">返回列表</el-button>
        <template v-if="isReceived && detail?.status === 1">
          <el-button type="success" @click="handleAccept">接收</el-button>
          <el-button type="danger" @click="handleReject">拒绝</el-button>
        </template>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 状态卡片 -->
      <el-card class="status-card">
        <div class="status-content">
          <div class="status-icon" :class="getStatusClass(detail?.status)">
            <el-icon v-if="detail?.status === 1"><Clock /></el-icon>
            <el-icon v-else-if="detail?.status === 2"><Loading /></el-icon>
            <el-icon v-else-if="detail?.status === 3"><CircleCheckFilled /></el-icon>
            <el-icon v-else><CircleCloseFilled /></el-icon>
          </div>
          <div class="status-info">
            <div class="status-text">{{ detail?.status_name }}</div>
            <div class="status-desc">
              <span v-if="detail?.status === 1">等待对方接收确认</span>
              <span v-else-if="detail?.status === 2">正在自动扣款中</span>
              <span v-else-if="detail?.status === 3">代扣已完成</span>
              <span v-else>对方已拒绝此代扣</span>
            </div>
          </div>
          <div class="progress-section">
            <el-progress
              type="dashboard"
              :percentage="detail?.progress || 0"
              :width="120"
              :stroke-width="10"
              :status="detail?.progress >= 100 ? 'success' : ''"
            >
              <template #default>
                <div class="progress-content">
                  <div class="percent">{{ (detail?.progress || 0).toFixed(1) }}%</div>
                  <div class="label">扣款进度</div>
                </div>
              </template>
            </el-progress>
          </div>
        </div>
      </el-card>

      <!-- 金额信息 -->
      <el-card class="amount-card">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="amount-item">
              <span class="label">代扣总额</span>
              <span class="value">¥{{ formatAmount(detail?.total_amount) }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="amount-item success">
              <span class="label">已扣金额</span>
              <span class="value">¥{{ formatAmount(detail?.deducted_amount) }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="amount-item danger">
              <span class="label">剩余待扣</span>
              <span class="value">¥{{ formatAmount(detail?.remaining_amount) }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="amount-item">
              <span class="label">终端数量</span>
              <span class="value">{{ detail?.terminal_count }} 台</span>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 基本信息 -->
      <el-card class="detail-card">
        <template #header>基本信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="代扣编号">{{ detail?.deduction_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusTag(detail?.status)" size="small">
              {{ detail?.status_name }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="扣款来源">
            <el-tag :type="getSourceTag(detail?.deduction_source)" size="small">
              {{ detail?.source_name }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="发起方">{{ detail?.from_agent_name }}</el-descriptions-item>
          <el-descriptions-item label="接收方">{{ detail?.to_agent_name }}</el-descriptions-item>
          <el-descriptions-item label="终端单价">¥{{ formatAmount(detail?.unit_price) }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ detail?.created_at }}</el-descriptions-item>
          <el-descriptions-item v-if="detail?.accepted_at" label="接收时间">
            {{ detail?.accepted_at }}
          </el-descriptions-item>
          <el-descriptions-item v-if="detail?.completed_at" label="完成时间">
            {{ detail?.completed_at }}
          </el-descriptions-item>
          <el-descriptions-item v-if="detail?.remark" label="备注" :span="3">
            {{ detail?.remark }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 关联终端 -->
      <el-card v-if="terminals.length" class="detail-card">
        <template #header>
          <div class="card-header">
            <span>关联终端</span>
            <el-tag type="info" size="small">共 {{ terminals.length }} 台</el-tag>
          </div>
        </template>
        <el-table :data="terminals" border stripe max-height="300">
          <el-table-column type="index" label="序号" width="60" align="center" />
          <el-table-column prop="terminal_sn" label="终端SN" width="200" />
          <el-table-column prop="unit_price" label="单价" width="120" align="right">
            <template #default="{ row }">
              ¥{{ formatAmount(row.unit_price) }}
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="添加时间" />
        </el-table>
      </el-card>

      <!-- 扣款明细 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>扣款明细</span>
            <el-tag type="info" size="small">共 {{ deductionDetails.length }} 条</el-tag>
          </div>
        </template>
        <el-timeline v-if="deductionDetails.length">
          <el-timeline-item
            v-for="item in deductionDetails"
            :key="item.id"
            :timestamp="item.created_at"
            placement="top"
            :type="'success'"
          >
            <el-card shadow="hover" class="timeline-card">
              <div class="timeline-content">
                <div class="amount">
                  <span class="label">扣款金额</span>
                  <span class="value">¥{{ formatAmount(item.amount) }}</span>
                </div>
                <div class="info">
                  <span class="wallet">{{ item.wallet_type_name }}</span>
                  <span v-if="item.channel_name" class="channel">{{ item.channel_name }}</span>
                </div>
                <div class="balance">
                  <span>余额: ¥{{ formatAmount(item.wallet_balance_before) }} → ¥{{ formatAmount(item.wallet_balance_after) }}</span>
                </div>
                <div class="cumulative">
                  <span>累计已扣: ¥{{ formatAmount(item.cumulative_deducted) }}</span>
                  <span class="remaining">剩余: ¥{{ formatAmount(item.remaining_after) }}</span>
                </div>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-else description="暂无扣款记录" />
      </el-card>
    </div>

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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Clock,
  Loading,
  CircleCheckFilled,
  CircleCloseFilled,
} from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import {
  getGoodsDeductionDetail,
  getGoodsDeductionDetails,
  acceptGoodsDeduction,
  rejectGoodsDeduction,
} from '@/api/goodsDeduction'
import { formatAmount } from '@/utils/format'
import type {
  GoodsDeductionDetail,
  GoodsDeductionDetailRecord,
  GoodsDeductionTerminal,
  GoodsDeductionStatus,
  DeductionSource,
} from '@/types'
import { GOODS_DEDUCTION_STATUS_CONFIG, DEDUCTION_SOURCE_CONFIG } from '@/types/deduction'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<GoodsDeductionDetail | null>(null)
const terminals = ref<GoodsDeductionTerminal[]>([])
const deductionDetails = ref<GoodsDeductionDetailRecord[]>([])

// 是否为接收方视角（判断逻辑：通过URL参数或后端返回标记）
const isReceived = computed(() => {
  // 这里可以根据实际业务判断，比如对比当前登录用户ID与to_agent_id
  return route.query.type === 'received'
})

// 拒绝弹窗
const rejectDialogVisible = ref(false)
const rejectFormRef = ref<FormInstance>()
const rejecting = ref(false)
const rejectForm = reactive({
  reason: '',
})
const rejectRules: FormRules = {
  reason: [{ required: true, message: '请输入拒绝原因', trigger: 'blur' }],
}

// 状态样式
function getStatusClass(status?: GoodsDeductionStatus) {
  const map: Record<number, string> = {
    1: 'pending',
    2: 'active',
    3: 'completed',
    4: 'rejected',
  }
  return status ? map[status] : ''
}

function getStatusTag(status?: GoodsDeductionStatus) {
  return status ? GOODS_DEDUCTION_STATUS_CONFIG[status]?.type : 'info'
}

// 来源标签
function getSourceTag(source?: DeductionSource) {
  if (!source) return 'info'
  const config = DEDUCTION_SOURCE_CONFIG[source]
  return config?.color === '#409eff' ? 'primary' : config?.color === '#67c23a' ? 'success' : 'warning'
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    detail.value = await getGoodsDeductionDetail(id)
    terminals.value = detail.value.terminals || []

    // 获取扣款明细
    const detailRes = await getGoodsDeductionDetails(id, { page_size: 100 })
    deductionDetails.value = detailRes.list
  } catch (error) {
    console.error('Fetch detail error:', error)
    ElMessage.error('获取详情失败')
  } finally {
    loading.value = false
  }
}

// 返回
function handleBack() {
  router.push('/goods-deductions/list')
}

// 接收
async function handleAccept() {
  if (!detail.value) return
  try {
    await ElMessageBox.confirm(
      '确定要接收此货款代扣吗？接收后将开始自动扣款。',
      '接收确认',
      {
        confirmButtonText: '确定接收',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    await acceptGoodsDeduction(detail.value.id)
    ElMessage.success('接收成功，代扣已开始')
    fetchDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Accept error:', error)
    }
  }
}

// 打开拒绝弹窗
function handleReject() {
  rejectForm.reason = ''
  rejectDialogVisible.value = true
}

// 确认拒绝
async function confirmReject() {
  if (!rejectFormRef.value || !detail.value) return

  try {
    await rejectFormRef.value.validate()
    rejecting.value = true
    await rejectGoodsDeduction(detail.value.id, { reason: rejectForm.reason })
    ElMessage.success('已拒绝')
    rejectDialogVisible.value = false
    fetchDetail()
  } catch (error) {
    console.error('Reject error:', error)
  } finally {
    rejecting.value = false
  }
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.goods-deduction-detail-view {
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

      &.pending {
        background: rgba($warning-color, 0.1);
        color: $warning-color;
      }

      &.active {
        background: rgba($primary-color, 0.1);
        color: $primary-color;
      }

      &.completed {
        background: rgba($success-color, 0.1);
        color: $success-color;
      }

      &.rejected {
        background: rgba($danger-color, 0.1);
        color: $danger-color;
      }
    }

    .status-info {
      flex: 1;

      .status-text {
        font-size: 20px;
        font-weight: 600;
        margin-bottom: $spacing-xs;
      }

      .status-desc {
        color: $text-secondary;
      }
    }

    .progress-section {
      .progress-content {
        text-align: center;

        .percent {
          font-size: 20px;
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

.timeline-card {
  .timeline-content {
    .amount {
      margin-bottom: $spacing-sm;

      .label {
        color: $text-secondary;
        margin-right: $spacing-sm;
      }

      .value {
        font-size: 18px;
        font-weight: 600;
        color: $success-color;
      }
    }

    .info {
      margin-bottom: $spacing-xs;

      .wallet {
        color: $primary-color;
        margin-right: $spacing-md;
      }

      .channel {
        color: $text-secondary;
      }
    }

    .balance,
    .cumulative {
      font-size: 12px;
      color: $text-secondary;
      margin-bottom: $spacing-xs;

      .remaining {
        margin-left: $spacing-md;
        color: $danger-color;
      }
    }
  }
}
</style>
