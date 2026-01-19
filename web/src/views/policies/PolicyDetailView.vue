<template>
  <div class="policy-detail-view">
    <PageHeader title="政策模板详情" :sub-title="`模板ID: ${route.params.id}`">
      <template #extra>
        <el-button @click="handleBack">返回</el-button>
        <el-button type="primary" @click="handleEdit">编辑</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading">
      <!-- 基本信息 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <div class="header-tags">
              <el-tag v-if="detail?.is_default" type="warning" size="small">
                <el-icon><StarFilled /></el-icon> 默认模板
              </el-tag>
              <el-tag :type="detail?.status === 1 ? 'success' : 'danger'" size="small">
                {{ detail?.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </div>
          </div>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="模板ID">{{ detail?.id }}</el-descriptions-item>
          <el-descriptions-item label="模板名称">{{ detail?.name }}</el-descriptions-item>
          <el-descriptions-item label="所属通道">{{ detail?.channel_name }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ detail?.created_at }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ detail?.updated_at }}</el-descriptions-item>
          <el-descriptions-item label="创建人">{{ detail?.created_by || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 费率配置 -->
      <el-card class="detail-card">
        <template #header>费率配置</template>
        <el-row :gutter="20">
          <el-col :span="4">
            <div class="rate-item">
              <div class="rate-label">贷记卡费率</div>
              <div class="rate-value primary">{{ formatRate(detail?.credit_rate) }}%</div>
            </div>
          </el-col>
          <el-col :span="4">
            <div class="rate-item">
              <div class="rate-label">借记卡费率</div>
              <div class="rate-value success">{{ formatRate(detail?.debit_rate) }}%</div>
            </div>
          </el-col>
          <el-col :span="4">
            <div class="rate-item">
              <div class="rate-label">借记卡封顶</div>
              <div class="rate-value warning">¥{{ detail?.debit_cap || 0 }}</div>
            </div>
          </el-col>
          <el-col :span="4">
            <div class="rate-item">
              <div class="rate-label">云闪付费率</div>
              <div class="rate-value info">{{ formatRate(detail?.qrcode_rate) }}%</div>
            </div>
          </el-col>
          <el-col :span="4">
            <div class="rate-item">
              <div class="rate-label">扫码费率</div>
              <div class="rate-value">{{ formatRate(detail?.scan_rate) }}%</div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 分润规则 -->
      <el-card class="detail-card">
        <template #header>分润规则</template>
        <el-table :data="detail?.profit_rules || []" border stripe>
          <el-table-column prop="level" label="代理层级" width="120" align="center">
            <template #default="{ row }">
              <el-tag size="small">{{ getLevelLabel(row.level) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="profit_type" label="分润类型" width="120" align="center">
            <template #default="{ row }">
              {{ getProfitTypeLabel(row.profit_type) }}
            </template>
          </el-table-column>
          <el-table-column prop="profit_value" label="分润值" width="120" align="right">
            <template #default="{ row }">
              <span class="profit-value">{{ formatProfitValue(row) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="min_amount" label="最小交易额" width="120" align="right">
            <template #default="{ row }">
              {{ row.min_amount > 0 ? `¥${row.min_amount}` : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="max_amount" label="最大交易额" width="120" align="right">
            <template #default="{ row }">
              {{ row.max_amount > 0 ? `¥${row.max_amount}` : '不限' }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="150">
            <template #default="{ row }">
              {{ row.remark || '-' }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="!detail?.profit_rules?.length" description="暂无分润规则" />
      </el-card>

      <!-- 模板说明 -->
      <el-card v-if="detail?.description" class="detail-card">
        <template #header>模板说明</template>
        <div class="description-content">{{ detail?.description }}</div>
      </el-card>

      <!-- 使用情况 -->
      <el-card class="detail-card">
        <template #header>使用情况</template>
        <el-descriptions :column="4" border>
          <el-descriptions-item label="关联代理数">
            <span class="usage-count">{{ detail?.agent_count || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="关联商户数">
            <span class="usage-count">{{ detail?.merchant_count || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="本月交易笔数">
            <span class="usage-count">{{ detail?.month_transaction_count || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="本月交易金额">
            <span class="usage-count primary">¥{{ formatAmount(detail?.month_transaction_amount || 0) }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { StarFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getPolicyTemplate } from '@/api/policy'
import { formatAmount } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)

// 详情数据
const detail = ref<any>(null)

// 格式化费率
function formatRate(rate?: number) {
  if (rate === undefined || rate === null) return '0.00'
  return (rate * 100).toFixed(2)
}

// 获取层级标签
function getLevelLabel(level: number) {
  const labels: Record<number, string> = {
    1: '一级代理',
    2: '二级代理',
    3: '三级代理',
  }
  return labels[level] || `${level}级代理`
}

// 获取分润类型标签
function getProfitTypeLabel(type: string) {
  const labels: Record<string, string> = {
    fixed: '固定金额',
    rate_diff: '费率差',
    percentage: '交易比例',
  }
  return labels[type] || type
}

// 格式化分润值
function formatProfitValue(row: any) {
  switch (row.profit_type) {
    case 'fixed':
      return `¥${row.profit_value}`
    case 'rate_diff':
    case 'percentage':
      return `${row.profit_value}%`
    default:
      return row.profit_value
  }
}

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    detail.value = await getPolicyTemplate(Number(route.params.id))
  } catch (error) {
    console.error('Fetch policy template error:', error)
    ElMessage.error('获取模板详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/policies/list')
}

// 编辑
function handleEdit() {
  router.push(`/policies/${route.params.id}/edit`)
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.policy-detail-view {
  padding: 0;
}

.detail-card {
  margin-top: $spacing-md;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-tags {
      display: flex;
      gap: $spacing-sm;
    }
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
    font-size: 24px;
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

    &.info {
      color: $info-color;
    }
  }
}

.profit-value {
  color: $success-color;
  font-weight: 600;
}

.description-content {
  color: $text-secondary;
  line-height: 1.8;
  white-space: pre-wrap;
}

.usage-count {
  font-size: 18px;
  font-weight: 600;

  &.primary {
    color: $primary-color;
  }
}
</style>
