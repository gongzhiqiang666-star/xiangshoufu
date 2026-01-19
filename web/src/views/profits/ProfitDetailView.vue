<template>
  <div class="profit-detail-view">
    <PageHeader title="分润详情" :sub-title="`分润编号: ${detail?.profit_no || ''}`">
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 分润状态 -->
      <el-card class="status-card">
        <div class="profit-header">
          <div class="profit-icon">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="profit-info">
            <div class="profit-amount">¥{{ formatAmount(detail?.profit_amount) }}</div>
            <div class="profit-type">
              <el-tag :type="getTypeTag(detail?.profit_type)" size="small">
                {{ getTypeLabel(detail?.profit_type) }}
              </el-tag>
            </div>
            <div class="profit-time">{{ detail?.created_at }}</div>
          </div>
          <div class="profit-status">
            <el-tag :type="detail?.status === 1 ? 'success' : 'warning'" size="large">
              {{ detail?.status === 1 ? '已入账' : '待入账' }}
            </el-tag>
          </div>
        </div>
      </el-card>

      <!-- 分润信息 -->
      <el-card class="detail-card">
        <template #header>分润信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="分润编号">{{ detail?.profit_no }}</el-descriptions-item>
          <el-descriptions-item label="分润类型">
            <el-tag :type="getTypeTag(detail?.profit_type)" size="small">
              {{ getTypeLabel(detail?.profit_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="分润金额">
            <span class="profit-value">¥{{ formatAmount(detail?.profit_amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="入账钱包">{{ detail?.wallet_type_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="入账状态">
            <el-tag :type="detail?.status === 1 ? 'success' : 'warning'" size="small">
              {{ detail?.status === 1 ? '已入账' : '待入账' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="入账时间">{{ detail?.settled_at || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 关联交易 -->
      <el-card v-if="detail?.transaction" class="detail-card">
        <template #header>关联交易</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="交易单号">
            <el-link type="primary" @click="handleViewTransaction">
              {{ detail?.related_no }}
            </el-link>
          </el-descriptions-item>
          <el-descriptions-item label="交易金额">
            ¥{{ formatAmount(detail?.transaction_amount) }}
          </el-descriptions-item>
          <el-descriptions-item label="交易状态">
            <el-tag :type="detail?.transaction?.status === 'success' ? 'success' : 'danger'" size="small">
              {{ detail?.transaction?.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="商户名称">{{ detail?.transaction?.merchant_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="终端SN">{{ detail?.transaction?.terminal_sn || '-' }}</el-descriptions-item>
          <el-descriptions-item label="交易时间">{{ detail?.transaction?.created_at || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 代理商信息 -->
      <el-card class="detail-card">
        <template #header>代理商信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="代理商名称">
            <el-link type="primary" @click="handleViewAgent">
              {{ detail?.agent_name }}
            </el-link>
          </el-descriptions-item>
          <el-descriptions-item label="代理层级">{{ detail?.agent_level }}级</el-descriptions-item>
          <el-descriptions-item label="所属通道">{{ detail?.channel_name }}</el-descriptions-item>
          <el-descriptions-item label="分润费率">
            {{ detail?.profit_rate ? `${(detail.profit_rate * 100).toFixed(4)}%` : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="政策模板">{{ detail?.policy_name || '默认' }}</el-descriptions-item>
          <el-descriptions-item label="分润方式">{{ getProfitModeLabel(detail?.profit_mode) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 分润链路 -->
      <el-card v-if="detail?.profit_chain?.length" class="detail-card">
        <template #header>分润链路</template>
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in detail.profit_chain"
            :key="index"
            :type="item.is_current ? 'primary' : 'info'"
            :hollow="!item.is_current"
          >
            <div class="chain-item">
              <div class="chain-header">
                <span class="chain-name">{{ item.agent_name }}</span>
                <el-tag size="small" type="info">{{ item.level }}级代理</el-tag>
                <span v-if="item.is_current" class="current-tag">当前</span>
              </div>
              <div class="chain-detail">
                <span>分润金额: <b class="profit-value">¥{{ formatAmount(item.profit_amount) }}</b></span>
                <span>费率: {{ (item.rate * 100).toFixed(4) }}%</span>
              </div>
            </div>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { TrendCharts } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getProfit } from '@/api/profit'
import { formatAmount } from '@/utils/format'
import type { ProfitType } from '@/types'
import { PROFIT_TYPE_CONFIG } from '@/types/profit'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)

// 详情数据
const detail = ref<any>(null)

// 获取类型标签
function getTypeTag(type?: ProfitType) {
  if (!type) return ''
  const colorMap: Record<string, string> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
  }
  const config = PROFIT_TYPE_CONFIG[type]
  return colorMap[config?.color] || ''
}

// 获取类型名称
function getTypeLabel(type?: ProfitType) {
  if (!type) return '-'
  return PROFIT_TYPE_CONFIG[type]?.label || type
}

// 获取分润方式名称
function getProfitModeLabel(mode?: string) {
  const modeMap: Record<string, string> = {
    rate_diff: '费率差',
    fixed: '固定金额',
    percentage: '交易比例',
  }
  return modeMap[mode || ''] || mode || '-'
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    detail.value = await getProfit(Number(route.params.id))
  } catch (error) {
    console.error('Fetch profit detail error:', error)
    ElMessage.error('获取分润详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/profits/list')
}

// 查看交易
function handleViewTransaction() {
  if (detail.value?.transaction_id) {
    router.push(`/transactions/${detail.value.transaction_id}`)
  }
}

// 查看代理商
function handleViewAgent() {
  if (detail.value?.agent_id) {
    router.push(`/agents/${detail.value.agent_id}`)
  }
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.profit-detail-view {
  padding: 0;
}

.status-card {
  margin-top: $spacing-md;

  .profit-header {
    display: flex;
    align-items: center;
    gap: $spacing-xl;
    padding: $spacing-lg;
  }

  .profit-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: linear-gradient(135deg, $success-color, lighten($success-color, 15%));
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36px;
    color: #fff;
  }

  .profit-info {
    flex: 1;

    .profit-amount {
      font-size: 32px;
      font-weight: 700;
      color: $success-color;
      margin-bottom: $spacing-xs;
    }

    .profit-type {
      margin-bottom: $spacing-sm;
    }

    .profit-time {
      font-size: 14px;
      color: $text-secondary;
    }
  }

  .profit-status {
    padding: 0 $spacing-xl;
  }
}

.detail-card {
  margin-top: $spacing-md;
}

.profit-value {
  color: $success-color;
  font-weight: 600;
}

.chain-item {
  .chain-header {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    margin-bottom: $spacing-xs;

    .chain-name {
      font-weight: 600;
    }

    .current-tag {
      color: $primary-color;
      font-size: 12px;
    }
  }

  .chain-detail {
    display: flex;
    gap: $spacing-lg;
    font-size: 13px;
    color: $text-secondary;
  }
}
</style>
