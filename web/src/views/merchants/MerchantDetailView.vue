<template>
  <div class="merchant-detail-view">
    <PageHeader title="商户详情" :sub-title="`商户编号: ${detail?.merchant_no || ''}`">
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 商户状态 -->
      <el-card class="status-card">
        <div class="merchant-header">
          <div class="merchant-avatar">
            <el-icon><Shop /></el-icon>
          </div>
          <div class="merchant-info">
            <div class="merchant-name">{{ detail?.name }}</div>
            <div class="merchant-no">商户编号: {{ detail?.merchant_no }}</div>
            <div class="merchant-tags">
              <el-tag :type="detail?.status === 1 ? 'success' : 'danger'" size="small">
                {{ detail?.status === 1 ? '正常' : '禁用' }}
              </el-tag>
              <el-tag type="info" size="small">{{ detail?.channel_name }}</el-tag>
            </div>
          </div>
          <div class="merchant-stats">
            <div class="stat-item">
              <div class="stat-value">{{ detail?.terminal_count || 0 }}</div>
              <div class="stat-label">终端数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">¥{{ formatAmount(detail?.month_amount) }}</div>
              <div class="stat-label">本月交易</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ detail?.month_count || 0 }}</div>
              <div class="stat-label">本月笔数</div>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Tab切换 -->
      <el-card class="detail-card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="基本信息" name="info">
            <!-- 基本信息 -->
            <el-descriptions :column="3" border>
              <el-descriptions-item label="商户ID">{{ detail?.id }}</el-descriptions-item>
              <el-descriptions-item label="商户名称">{{ detail?.name }}</el-descriptions-item>
              <el-descriptions-item label="商户编号">{{ detail?.merchant_no }}</el-descriptions-item>
              <el-descriptions-item label="所属通道">{{ detail?.channel_name }}</el-descriptions-item>
              <el-descriptions-item label="所属代理">{{ detail?.agent_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="代理层级">{{ detail?.agent_level ? `${detail.agent_level}级` : '-' }}</el-descriptions-item>
              <el-descriptions-item label="联系人">{{ detail?.contact_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="联系电话">{{ detail?.contact_phone || '-' }}</el-descriptions-item>
              <el-descriptions-item label="经营地址">{{ detail?.address || '-' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detail?.created_at }}</el-descriptions-item>
              <el-descriptions-item label="更新时间">{{ detail?.updated_at }}</el-descriptions-item>
              <el-descriptions-item label="激活时间">{{ detail?.activated_at || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane label="费率信息" name="rate">
            <div class="rate-section">
              <div class="rate-header">
                <span>当前费率</span>
                <el-button type="primary" size="small" @click="handleEditRate">修改费率</el-button>
              </div>
              <el-row :gutter="20" style="margin-top: 16px;">
                <el-col :span="6">
                  <div class="rate-item">
                    <div class="rate-label">贷记卡费率</div>
                    <div class="rate-value primary">{{ formatRate(detail?.credit_rate) }}%</div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="rate-item">
                    <div class="rate-label">借记卡费率</div>
                    <div class="rate-value success">{{ formatRate(detail?.debit_rate) }}%</div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="rate-item">
                    <div class="rate-label">借记卡封顶</div>
                    <div class="rate-value warning">¥{{ detail?.debit_cap || 0 }}</div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="rate-item">
                    <div class="rate-label">政策模板</div>
                    <div class="rate-value">{{ detail?.policy_name || '默认' }}</div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>

          <el-tab-pane label="费率同步记录" name="sync-logs">
            <el-table :data="syncLogs" v-loading="syncLogsLoading" border stripe max-height="400">
              <el-table-column prop="id" label="ID" width="80" />
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
              <el-table-column prop="created_at" label="操作时间" width="170" />
            </el-table>
            <el-empty v-if="!syncLogsLoading && !syncLogs.length" description="暂无费率同步记录" />
          </el-tab-pane>

          <el-tab-pane label="关联终端" name="terminals">
            <el-table :data="terminals" border stripe max-height="300">
              <el-table-column prop="sn" label="终端SN" width="180" />
              <el-table-column prop="model" label="终端型号" width="120" />
              <el-table-column prop="status" label="状态" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                    {{ row.status === 1 ? '正常' : '停用' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="bind_at" label="绑定时间" width="170" />
              <el-table-column prop="last_trade_at" label="最后交易" min-width="170" />
            </el-table>
            <el-empty v-if="!terminals.length" description="暂无关联终端" />
          </el-tab-pane>

          <el-tab-pane label="近期交易" name="transactions">
            <div class="transactions-header">
              <el-button type="primary" link @click="handleViewAllTransactions">
                查看全部
              </el-button>
            </div>
            <el-table :data="recentTransactions" border stripe max-height="300">
              <el-table-column prop="trade_no" label="交易单号" width="200" />
              <el-table-column prop="amount" label="交易金额" width="120" align="right">
                <template #default="{ row }">
                  ¥{{ formatAmount(row.amount) }}
                </template>
              </el-table-column>
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
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <!-- 费率修改弹窗 -->
    <RateEditDialog
      v-model="rateDialogVisible"
      :merchant="rateDialogMerchant"
      @success="handleRateEditSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Shop } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import RateEditDialog from '@/components/Merchant/RateEditDialog.vue'
import { getMerchant, getMerchantTerminals, getMerchantTransactions } from '@/api/merchant'
import { getRateSyncLogs, type RateSyncLog } from '@/api/rateSync'
import { formatAmount } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)
const syncLogsLoading = ref(false)

// Tab
const activeTab = ref('info')

// 详情数据
const detail = ref<any>(null)
const terminals = ref<any[]>([])
const recentTransactions = ref<any[]>([])
const syncLogs = ref<RateSyncLog[]>([])

// 费率修改弹窗
const rateDialogVisible = ref(false)
const rateDialogMerchant = computed(() => {
  if (!detail.value) return null
  return {
    id: detail.value.id,
    merchant_no: detail.value.merchant_no,
    merchant_name: detail.value.name,
    credit_rate: detail.value.credit_rate,
    debit_rate: detail.value.debit_rate,
  }
})

// 格式化费率
function formatRate(rate?: number | string) {
  if (rate === undefined || rate === null) return '0.00'
  const num = typeof rate === 'string' ? parseFloat(rate) : rate
  return (num * 100).toFixed(2)
}

// 格式化费率为百分比
function formatRatePercent(rate: number | undefined): string {
  if (rate === undefined || rate === null) return '-'
  return `${(rate * 100).toFixed(2)}%`
}

// 获取同步状态类型
function getSyncStatusType(status: number): string {
  switch (status) {
    case 0: return 'info'
    case 1: return 'warning'
    case 2: return 'success'
    case 3: return 'danger'
    default: return 'info'
  }
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    const [merchantData, terminalsData, transactionsData] = await Promise.all([
      getMerchant(Number(route.params.id)),
      getMerchantTerminals(Number(route.params.id)).catch(() => []),
      getMerchantTransactions(Number(route.params.id), { page: 1, page_size: 5 }),
    ])
    detail.value = merchantData
    terminals.value = terminalsData || []
    recentTransactions.value = transactionsData?.list || []
  } catch (error) {
    console.error('Fetch merchant detail error:', error)
    ElMessage.error('获取商户详情失败')
  } finally {
    loading.value = false
  }
}

// 获取费率同步日志
async function fetchSyncLogs() {
  if (!detail.value) return
  syncLogsLoading.value = true
  try {
    const res = await getRateSyncLogs({
      merchant_id: detail.value.id,
      page: 1,
      page_size: 20,
    })
    syncLogs.value = res.items || []
  } catch (error) {
    console.error('Fetch sync logs error:', error)
  } finally {
    syncLogsLoading.value = false
  }
}

// 监听tab切换，懒加载同步日志
watch(activeTab, (tab) => {
  if (tab === 'sync-logs' && syncLogs.value.length === 0) {
    fetchSyncLogs()
  }
})

// 返回列表
function handleBack() {
  router.push('/merchants/list')
}

// 修改费率
function handleEditRate() {
  rateDialogVisible.value = true
}

// 费率修改成功后刷新数据
function handleRateEditSuccess() {
  fetchDetail()
  // 如果在同步日志tab，刷新同步日志
  if (activeTab.value === 'sync-logs') {
    fetchSyncLogs()
  }
}

// 查看全部交易
function handleViewAllTransactions() {
  router.push({
    path: '/transactions/list',
    query: { merchant_id: route.params.id as string },
  })
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.merchant-detail-view {
  padding: 0;
}

.status-card {
  margin-top: $spacing-md;

  .merchant-header {
    display: flex;
    align-items: center;
    gap: $spacing-xl;
  }

  .merchant-avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: linear-gradient(135deg, $primary-color, lighten($primary-color, 15%));
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36px;
    color: #fff;
  }

  .merchant-info {
    flex: 1;

    .merchant-name {
      font-size: 24px;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }

    .merchant-no {
      font-size: 14px;
      color: $text-secondary;
      margin-bottom: $spacing-sm;
    }

    .merchant-tags {
      display: flex;
      gap: $spacing-sm;
    }
  }

  .merchant-stats {
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

.rate-section {
  .rate-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    color: $text-primary;
  }
}

.rate-item {
  text-align: center;
  padding: $spacing-md;
  background: $bg-color;
  border-radius: $border-radius-sm;

  .rate-label {
    font-size: 12px;
    color: $text-secondary;
    margin-bottom: $spacing-xs;
  }

  .rate-value {
    font-size: 20px;
    font-weight: 600;
    color: $text-primary;

    &.primary {
      color: $primary-color;
    }

    &.success {
      color: $success-color;
    }

    &.warning {
      color: $warning-color;
    }
  }
}

.new-rate {
  color: #409eff;
  font-weight: 500;
}

.transactions-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}
</style>
