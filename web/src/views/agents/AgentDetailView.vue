<template>
  <div class="agent-detail-view">
    <PageHeader :title="agentDetail?.name || '代理商详情'" show-back>
      <template #extra>
        <el-button type="primary" @click="handleEdit">编辑资料</el-button>
      </template>
    </PageHeader>

    <div v-loading="loading" class="detail-content">
      <!-- 基本信息 -->
      <el-card class="info-card">
        <template #header>
          <span>基本信息</span>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="代理商编号">
            {{ agentDetail?.agent_code }}
          </el-descriptions-item>
          <el-descriptions-item label="姓名">
            {{ agentDetail?.name }}
          </el-descriptions-item>
          <el-descriptions-item label="手机号">
            {{ agentDetail?.phone }}
          </el-descriptions-item>
          <el-descriptions-item label="身份证号">
            {{ agentDetail?.id_card_no }}
          </el-descriptions-item>
          <el-descriptions-item label="上级代理">
            {{ agentDetail?.parent_name || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="层级">
            <el-tag size="small">{{ agentDetail?.level }}级</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="所属通道">
            {{ agentDetail?.channel_name }}
          </el-descriptions-item>
          <el-descriptions-item label="邀请码">
            <el-tag type="info">{{ agentDetail?.invite_code }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ agentDetail?.created_at }}
          </el-descriptions-item>
          <el-descriptions-item label="结算银行">
            {{ agentDetail?.bank_name }}
          </el-descriptions-item>
          <el-descriptions-item label="银行卡号" :span="2">
            {{ agentDetail?.bank_card_no }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 统计数据 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">{{ agentDetail?.direct_agent_count || 0 }}</div>
            <div class="stat-label">直属代理</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">{{ agentDetail?.team_agent_count || 0 }}</div>
            <div class="stat-label">团队代理</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">{{ agentDetail?.direct_merchant_count || 0 }}</div>
            <div class="stat-label">直营商户</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">{{ agentDetail?.team_merchant_count || 0 }}</div>
            <div class="stat-label">团队商户</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">{{ agentDetail?.terminal_total || 0 }}</div>
            <div class="stat-label">终端总数</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">{{ agentDetail?.terminal_activated || 0 }}</div>
            <div class="stat-label">已激活</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">¥{{ formatAmount(agentDetail?.month_transaction_amount) }}</div>
            <div class="stat-label">本月交易</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-value">¥{{ formatAmount(agentDetail?.total_profit) }}</div>
            <div class="stat-label">累计分润</div>
          </div>
        </el-col>
      </el-row>

      <!-- Tab区域 -->
      <el-card class="tabs-card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="政策模板" name="policy">
            <div class="tab-content">
              <p class="empty-tip">暂无政策模板数据</p>
            </div>
          </el-tab-pane>
          <el-tab-pane label="钱包信息" name="wallet">
            <div class="tab-content">
              <p class="empty-tip">暂无钱包数据</p>
            </div>
          </el-tab-pane>
          <el-tab-pane label="下级代理" name="subordinates">
            <div class="tab-content">
              <p class="empty-tip">暂无下级代理数据</p>
            </div>
          </el-tab-pane>
          <el-tab-pane label="商户列表" name="merchants">
            <div class="tab-content">
              <p class="empty-tip">暂无商户数据</p>
            </div>
          </el-tab-pane>
          <el-tab-pane label="终端列表" name="terminals">
            <div class="tab-content">
              <p class="empty-tip">暂无终端数据</p>
            </div>
          </el-tab-pane>
          <el-tab-pane label="交易记录" name="transactions">
            <div class="tab-content">
              <p class="empty-tip">暂无交易数据</p>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getAgentDetail } from '@/api/agent'
import { formatAmount } from '@/utils/format'
import type { AgentDetail } from '@/types'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const agentDetail = ref<AgentDetail | null>(null)
const activeTab = ref('policy')

// 获取代理商ID
const agentId = Number(route.params.id)

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    agentDetail.value = await getAgentDetail(agentId)
  } catch (error) {
    console.error('Fetch agent detail error:', error)
  } finally {
    loading.value = false
  }
}

// 编辑
function handleEdit() {
  router.push(`/agents/${agentId}/edit`)
}

onMounted(() => {
  fetchDetail()
})
</script>

<style lang="scss" scoped>
.agent-detail-view {
  padding: 0;
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.info-card {
  :deep(.el-descriptions__label) {
    width: 120px;
  }
}

.stats-row {
  .el-col {
    margin-bottom: $spacing-md;
  }
}

.stat-card {
  background: $bg-white;
  border-radius: $border-radius-md;
  padding: $spacing-lg;
  text-align: center;
  box-shadow: $shadow-sm;
  transition: all $transition-normal;

  &:hover {
    box-shadow: $shadow-md;
    transform: translateY(-2px);
  }

  .stat-value {
    font-size: 24px;
    font-weight: 600;
    color: $primary-color;
    margin-bottom: $spacing-xs;
  }

  .stat-label {
    font-size: 14px;
    color: $text-secondary;
  }
}

.tabs-card {
  .tab-content {
    min-height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .empty-tip {
    color: $text-secondary;
  }
}
</style>
