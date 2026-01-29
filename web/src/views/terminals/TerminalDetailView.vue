<template>
  <div class="terminal-detail-view">
    <PageHeader title="终端详情" :sub-title="`终端SN: ${detail?.sn || ''}`">
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 终端状态卡片 -->
      <el-card class="status-card">
        <div class="terminal-header">
          <div class="terminal-icon" :class="getStatusClass(detail?.status)">
            <el-icon><Monitor /></el-icon>
          </div>
          <div class="terminal-info">
            <div class="terminal-sn">{{ detail?.sn }}</div>
            <div class="terminal-model">{{ detail?.model }} | {{ detail?.brand }}</div>
            <div class="terminal-tags">
              <el-tag :type="detail?.status === 1 ? 'success' : 'danger'" size="small">
                {{ detail?.status === 1 ? '正常' : '停用' }}
              </el-tag>
              <el-tag :type="detail?.is_activated ? 'success' : 'info'" size="small">
                {{ detail?.is_activated ? '已激活' : '未激活' }}
              </el-tag>
              <el-tag type="info" size="small">{{ detail?.channel_name }}</el-tag>
            </div>
          </div>
          <div class="terminal-stats">
            <div class="stat-item">
              <div class="stat-value">¥{{ formatAmount(detail?.total_amount) }}</div>
              <div class="stat-label">累计交易</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ detail?.total_count || 0 }}</div>
              <div class="stat-label">交易笔数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ detail?.days_active || 0 }}</div>
              <div class="stat-label">活跃天数</div>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 基本信息 -->
      <el-card class="detail-card">
        <template #header>基本信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="终端ID">{{ detail?.id }}</el-descriptions-item>
          <el-descriptions-item label="终端SN">{{ detail?.sn }}</el-descriptions-item>
          <el-descriptions-item label="所属通道">{{ detail?.channel_name }}</el-descriptions-item>
          <el-descriptions-item label="终端品牌">{{ detail?.brand || '-' }}</el-descriptions-item>
          <el-descriptions-item label="终端型号">{{ detail?.model || '-' }}</el-descriptions-item>
          <el-descriptions-item label="SIM卡号">{{ detail?.sim_no || '-' }}</el-descriptions-item>
          <el-descriptions-item label="入库时间">{{ detail?.created_at }}</el-descriptions-item>
          <el-descriptions-item label="绑定时间">{{ detail?.bind_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="激活时间">{{ detail?.activated_at || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 绑定信息 -->
      <el-card class="detail-card">
        <template #header>绑定信息</template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="绑定商户">
            <el-link v-if="detail?.merchant_id" type="primary" @click="handleViewMerchant">
              {{ detail?.merchant_name }}
            </el-link>
            <span v-else class="text-placeholder">未绑定</span>
          </el-descriptions-item>
          <el-descriptions-item label="商户编号">{{ detail?.merchant_no || '-' }}</el-descriptions-item>
          <el-descriptions-item label="所属代理">{{ detail?.agent_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="代理层级">{{ detail?.agent_level ? `${detail.agent_level}级` : '-' }}</el-descriptions-item>
          <el-descriptions-item label="政策模板">{{ detail?.policy_name || '默认' }}</el-descriptions-item>
          <el-descriptions-item label="绑定状态">
            <el-tag :type="detail?.merchant_id ? 'success' : 'info'" size="small">
              {{ detail?.merchant_id ? '已绑定' : '未绑定' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 费用信息 -->
      <el-card class="detail-card">
        <template #header>费用信息</template>
        <el-descriptions :column="4" border>
          <el-descriptions-item label="押金金额">
            <span class="amount">¥{{ formatAmount(detail?.deposit_amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="押金状态">
            <el-tag :type="getDepositStatusType(detail?.deposit_status)" size="small">
              {{ getDepositStatusText(detail?.deposit_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="流量费">
            <span class="amount">¥{{ formatAmount(detail?.sim_fee) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="流量费状态">
            <el-tag :type="detail?.sim_fee_paid ? 'success' : 'warning'" size="small">
              {{ detail?.sim_fee_paid ? '已缴纳' : '待缴纳' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 近期交易 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>近期交易</span>
            <el-button type="primary" link @click="handleViewAllTransactions">
              查看全部
            </el-button>
          </div>
        </template>
        <el-table :data="recentTransactions" border stripe max-height="300">
          <el-table-column prop="trade_no" label="交易单号" width="200" />
          <el-table-column prop="amount" label="交易金额" width="120" align="right">
            <template #default="{ row }">
              ¥{{ formatAmount(row.amount) }}
            </template>
          </el-table-column>
          <el-table-column prop="card_type" label="卡类型" width="80" />
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
                {{ row.status === 'success' ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="交易时间" min-width="170" />
        </el-table>
        <el-empty v-if="!recentTransactions.length" description="暂无交易记录" />
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Monitor } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getTerminal, getTerminalTransactions } from '@/api/terminal'
import { formatAmount } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)

// 详情数据
const detail = ref<any>(null)
const recentTransactions = ref<any[]>([])

// 获取状态样式
function getStatusClass(status?: number) {
  return status === 1 ? 'active' : 'inactive'
}

// 获取押金状态类型
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getDepositStatusType(status?: string): TagType {
  const typeMap: Record<string, TagType> = {
    paid: 'success',
    refunded: 'info',
    pending: 'warning',
  }
  return typeMap[status || ''] || 'info'
}

// 获取押金状态文本
function getDepositStatusText(status?: string) {
  const textMap: Record<string, string> = {
    paid: '已收取',
    refunded: '已退还',
    pending: '待收取',
  }
  return textMap[status || ''] || '未知'
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    const [terminalData, transactionsData] = await Promise.all([
      getTerminal(Number(route.params.id)),
      getTerminalTransactions(Number(route.params.id), { page: 1, page_size: 5 }),
    ])
    detail.value = terminalData
    recentTransactions.value = transactionsData?.list || []
  } catch (error) {
    console.error('Fetch terminal detail error:', error)
    ElMessage.error('获取终端详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/terminals/list')
}

// 查看商户
function handleViewMerchant() {
  if (detail.value?.merchant_id) {
    router.push(`/merchants/${detail.value.merchant_id}`)
  }
}

// 查看全部交易
function handleViewAllTransactions() {
  router.push({
    path: '/transactions/list',
    query: { terminal_sn: detail.value?.sn },
  })
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.terminal-detail-view {
  padding: 0;
}

.status-card {
  margin-top: $spacing-md;

  .terminal-header {
    display: flex;
    align-items: center;
    gap: $spacing-xl;
  }

  .terminal-icon {
    width: 80px;
    height: 80px;
    border-radius: $border-radius-md;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36px;
    color: #fff;

    &.active {
      background: linear-gradient(135deg, $success-color, lighten($success-color, 15%));
    }

    &.inactive {
      background: linear-gradient(135deg, $danger-color, lighten($danger-color, 15%));
    }
  }

  .terminal-info {
    flex: 1;

    .terminal-sn {
      font-size: 24px;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }

    .terminal-model {
      font-size: 14px;
      color: $text-secondary;
      margin-bottom: $spacing-sm;
    }

    .terminal-tags {
      display: flex;
      gap: $spacing-sm;
    }
  }

  .terminal-stats {
    display: flex;
    gap: $spacing-xl;

    .stat-item {
      text-align: center;
      padding: 0 $spacing-lg;
      border-left: 1px solid $border-color;

      &:first-child {
        border-left: none;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: $text-primary;
      }

      .stat-label {
        font-size: 12px;
        color: $text-secondary;
        margin-top: $spacing-xs;
      }
    }
  }
}

.detail-card {
  margin-top: $spacing-md;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.amount {
  font-weight: 600;
  color: $primary-color;
}

.text-placeholder {
  color: $text-placeholder;
}
</style>
