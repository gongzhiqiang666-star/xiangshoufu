<template>
  <div class="transaction-detail-view">
    <PageHeader title="交易详情" :sub-title="`交易单号: ${detail?.trade_no || ''}`">
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 交易状态 -->
      <el-card class="status-card">
        <div class="status-content">
          <div class="status-icon" :class="getStatusClass(detail?.status)">
            <el-icon v-if="detail?.status === 'success'"><CircleCheckFilled /></el-icon>
            <el-icon v-else-if="detail?.status === 'failed'"><CircleCloseFilled /></el-icon>
            <el-icon v-else><Loading /></el-icon>
          </div>
          <div class="status-info">
            <div class="status-text">{{ getStatusText(detail?.status) }}</div>
            <div class="status-amount">¥{{ formatAmount(detail?.amount) }}</div>
            <div class="status-time">{{ detail?.created_at }}</div>
          </div>
        </div>
      </el-card>

      <!-- 交易信息 -->
      <el-card class="detail-card">
        <template #header>交易信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="交易单号">{{ detail?.trade_no }}</el-descriptions-item>
          <el-descriptions-item label="通道单号">{{ detail?.channel_trade_no || '-' }}</el-descriptions-item>
          <el-descriptions-item label="交易类型">
            <el-tag size="small">{{ getTypeLabel(detail?.type) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="交易金额">
            <span class="amount">¥{{ formatAmount(detail?.amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="手续费">
            <span class="fee">¥{{ formatAmount(detail?.fee) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="结算金额">
            <span class="settle-amount">¥{{ formatAmount(detail?.settle_amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="费率">{{ detail?.rate ? `${(detail.rate * 100).toFixed(2)}%` : '-' }}</el-descriptions-item>
          <el-descriptions-item label="卡类型">{{ detail?.card_type || '-' }}</el-descriptions-item>
          <el-descriptions-item label="卡号">{{ detail?.card_no || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 商户信息 -->
      <el-card class="detail-card">
        <template #header>商户信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="商户名称">{{ detail?.merchant_name }}</el-descriptions-item>
          <el-descriptions-item label="商户编号">{{ detail?.merchant_no }}</el-descriptions-item>
          <el-descriptions-item label="所属通道">{{ detail?.channel_name }}</el-descriptions-item>
          <el-descriptions-item label="终端编号">{{ detail?.terminal_sn }}</el-descriptions-item>
          <el-descriptions-item label="所属代理">{{ detail?.agent_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="代理层级">{{ detail?.agent_level ? `${detail.agent_level}级` : '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 分润信息 -->
      <el-card class="detail-card">
        <template #header>分润信息</template>
        <el-table :data="detail?.profit_details || []" border stripe>
          <el-table-column prop="agent_name" label="代理商" width="150" />
          <el-table-column prop="level" label="层级" width="80" align="center">
            <template #default="{ row }">
              {{ row.level }}级
            </template>
          </el-table-column>
          <el-table-column prop="profit_type" label="分润类型" width="100" align="center" />
          <el-table-column prop="profit_rate" label="分润费率" width="100" align="right">
            <template #default="{ row }">
              {{ row.profit_rate ? `${(row.profit_rate * 100).toFixed(4)}%` : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="profit_amount" label="分润金额" width="120" align="right">
            <template #default="{ row }">
              <span class="profit-amount">¥{{ formatAmount(row.profit_amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="wallet_type" label="入账钱包" width="100" />
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'warning'" size="small">
                {{ row.status === 1 ? '已入账' : '待入账' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="入账时间" min-width="170" />
        </el-table>
        <el-empty v-if="!detail?.profit_details?.length" description="暂无分润记录" />
      </el-card>

      <!-- 回调信息 -->
      <el-card class="detail-card">
        <template #header>回调信息</template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="回调时间">{{ detail?.callback_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="回调状态">
            <el-tag :type="detail?.callback_status === 1 ? 'success' : 'info'" size="small">
              {{ detail?.callback_status === 1 ? '已回调' : '未回调' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="原始数据" :span="2">
            <div class="raw-data">
              <pre>{{ formatJson(detail?.raw_data) }}</pre>
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CircleCheckFilled, CircleCloseFilled, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getTransaction } from '@/api/transaction'
import { formatAmount } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)

// 详情数据
const detail = ref<any>(null)

// 获取状态样式
function getStatusClass(status?: string) {
  const classMap: Record<string, string> = {
    success: 'success',
    failed: 'failed',
    pending: 'pending',
  }
  return classMap[status || ''] || 'pending'
}

// 获取状态文本
function getStatusText(status?: string) {
  const textMap: Record<string, string> = {
    success: '交易成功',
    failed: '交易失败',
    pending: '处理中',
  }
  return textMap[status || ''] || '未知状态'
}

// 获取类型标签
function getTypeLabel(type?: string) {
  const labelMap: Record<string, string> = {
    consume: '消费',
    refund: '退款',
    void: '撤销',
  }
  return labelMap[type || ''] || type || '-'
}

// 格式化JSON
function formatJson(str?: string) {
  if (!str) return '-'
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    detail.value = await getTransaction(Number(route.params.id))
  } catch (error) {
    console.error('Fetch transaction error:', error)
    ElMessage.error('获取交易详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/transactions/list')
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.transaction-detail-view {
  padding: 0;
}

.status-card {
  margin-top: $spacing-md;

  .status-content {
    display: flex;
    align-items: center;
    gap: $spacing-lg;
    padding: $spacing-lg;
  }

  .status-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 40px;

    &.success {
      background: rgba(103, 194, 58, 0.1);
      color: $success-color;
    }

    &.failed {
      background: rgba(245, 108, 108, 0.1);
      color: $danger-color;
    }

    &.pending {
      background: rgba(230, 162, 60, 0.1);
      color: $warning-color;
    }
  }

  .status-info {
    .status-text {
      font-size: 18px;
      font-weight: 600;
      margin-bottom: $spacing-xs;
    }

    .status-amount {
      font-size: 32px;
      font-weight: 700;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }

    .status-time {
      font-size: 14px;
      color: $text-secondary;
    }
  }
}

.detail-card {
  margin-top: $spacing-md;
}

.amount {
  font-weight: 600;
  color: $text-primary;
}

.fee {
  color: $warning-color;
}

.settle-amount {
  color: $success-color;
  font-weight: 600;
}

.profit-amount {
  color: $success-color;
  font-weight: 600;
}

.raw-data {
  max-height: 200px;
  overflow: auto;
  background: $bg-color;
  border-radius: $border-radius-sm;
  padding: $spacing-sm;

  pre {
    margin: 0;
    font-size: 12px;
    white-space: pre-wrap;
    word-break: break-all;
  }
}
</style>
