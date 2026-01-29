<template>
  <div class="agent-detail-view">
    <PageHeader :title="agentDetail?.agent_name || '代理商详情'" show-back>
      <template #extra>
        <el-button-group>
          <el-button
            :type="agentDetail?.status === 1 ? 'warning' : 'success'"
            @click="handleToggleStatus"
          >
            {{ agentDetail?.status === 1 ? '禁用' : '启用' }}
          </el-button>
          <el-button type="primary" @click="handleEdit">编辑资料</el-button>
        </el-button-group>
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
            {{ agentDetail?.agent_no }}
          </el-descriptions-item>
          <el-descriptions-item label="代理商名称">
            {{ agentDetail?.agent_name }}
          </el-descriptions-item>
          <el-descriptions-item label="联系电话">
            {{ agentDetail?.contact_phone }}
          </el-descriptions-item>
          <el-descriptions-item label="联系人">
            {{ agentDetail?.contact_name }}
          </el-descriptions-item>
          <el-descriptions-item label="上级代理">
            {{ agentDetail?.parent_name || '无（顶级代理）' }}
          </el-descriptions-item>
          <el-descriptions-item label="层级">
            <el-tag size="small">{{ agentDetail?.level }}级</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="agentDetail?.status === 1 ? 'success' : 'danger'">
              {{ agentDetail?.status_name }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="邀请码">
            <el-tag type="info">{{ agentDetail?.invite_code }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ formatDate(agentDetail?.register_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="结算银行">
            {{ agentDetail?.bank_name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="银行卡号" :span="2">
            {{ agentDetail?.bank_card_no || '-' }}
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
      </el-row>

      <!-- Tab区域 -->
      <el-card class="tabs-card">
        <el-tabs v-model="activeTab" @tab-change="handleTabChange">
          <!-- 政策模板Tab -->
          <el-tab-pane label="政策模板" name="policy">
            <div class="tab-content">
              <div class="tab-toolbar">
                <el-button type="primary" @click="showPolicyAssignDialog = true">
                  分配政策
                </el-button>
              </div>
              <el-table v-if="policies.length > 0" :data="policies" style="width: 100%">
                <el-table-column prop="channel_name" label="通道名称" width="150" />
                <el-table-column prop="template_name" label="政策模板" width="180" />
                <el-table-column prop="credit_rate" label="贷记卡费率" width="120">
                  <template #default="{ row }">
                    {{ (row.credit_rate * 100).toFixed(2) }}%
                  </template>
                </el-table-column>
                <el-table-column prop="debit_rate" label="借记卡费率" width="120">
                  <template #default="{ row }">
                    {{ (row.debit_rate * 100).toFixed(2) }}%
                  </template>
                </el-table-column>
                <el-table-column prop="debit_cap" label="借记卡封顶" width="120">
                  <template #default="{ row }">
                    {{ row.debit_cap ? `${row.debit_cap}元` : '不封顶' }}
                  </template>
                </el-table-column>
                <el-table-column prop="created_at" label="生效时间" min-width="160">
                  <template #default="{ row }">
                    {{ formatDate(row.created_at) }}
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无政策数据" />
            </div>
          </el-tab-pane>

          <!-- 钱包信息Tab -->
          <el-tab-pane label="钱包信息" name="wallet">
            <div class="tab-content">
              <el-table v-if="wallets.length > 0" :data="wallets" style="width: 100%">
                <el-table-column prop="channel_name" label="通道名称" width="150" />
                <el-table-column prop="wallet_type" label="钱包类型" width="120">
                  <template #default="{ row }">
                    <el-tag size="small" :type="getWalletTagType(row.wallet_type)">
                      {{ getWalletTypeName(row.wallet_type) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="balance" label="余额" width="140">
                  <template #default="{ row }">
                    <span class="money">¥{{ formatAmount(row.balance) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="frozen" label="冻结" width="140">
                  <template #default="{ row }">
                    <span class="frozen">¥{{ formatAmount(row.frozen) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="available" label="可用" width="140">
                  <template #default="{ row }">
                    <span class="available">¥{{ formatAmount(row.available) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="total_income" label="累计收入" min-width="140">
                  <template #default="{ row }">
                    ¥{{ formatAmount(row.total_income) }}
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无钱包数据" />
            </div>
          </el-tab-pane>

          <!-- 下级代理Tab -->
          <el-tab-pane label="下级代理" name="subordinates">
            <div class="tab-content">
              <div class="tab-toolbar">
                <el-input
                  v-model="subordinatesSearch"
                  placeholder="搜索代理商名称/编号"
                  style="width: 250px"
                  clearable
                  @keyup.enter="fetchSubordinates"
                >
                  <template #append>
                    <el-button @click="fetchSubordinates">搜索</el-button>
                  </template>
                </el-input>
              </div>
              <el-table v-if="subordinates.length > 0" :data="subordinates" style="width: 100%">
                <el-table-column prop="agent_no" label="代理商编号" width="150" />
                <el-table-column prop="agent_name" label="代理商名称" width="150" />
                <el-table-column prop="contact_phone" label="联系电话" width="140" />
                <el-table-column prop="level" label="层级" width="80">
                  <template #default="{ row }">
                    {{ row.level }}级
                  </template>
                </el-table-column>
                <el-table-column prop="status" label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                      {{ row.status_name }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="direct_agent_count" label="直属代理" width="100" />
                <el-table-column prop="direct_merchant_count" label="直营商户" width="100" />
                <el-table-column prop="register_time" label="注册时间" min-width="160">
                  <template #default="{ row }">
                    {{ formatDate(row.register_time) }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="100" fixed="right">
                  <template #default="{ row }">
                    <el-button type="primary" link @click="viewAgent(row.id)">查看</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无下级代理" />
              <el-pagination
                v-if="subordinatesTotal > 10"
                v-model:current-page="subordinatesPage"
                :page-size="10"
                :total="subordinatesTotal"
                layout="prev, pager, next"
                @current-change="fetchSubordinates"
              />
            </div>
          </el-tab-pane>

          <!-- 商户列表Tab -->
          <el-tab-pane label="商户列表" name="merchants">
            <div class="tab-content">
              <div class="tab-toolbar">
                <el-input
                  v-model="merchantsSearch"
                  placeholder="搜索商户名称/编号"
                  style="width: 250px"
                  clearable
                  @keyup.enter="fetchMerchants"
                >
                  <template #append>
                    <el-button @click="fetchMerchants">搜索</el-button>
                  </template>
                </el-input>
              </div>
              <el-table v-if="merchants.length > 0" :data="merchants" style="width: 100%">
                <el-table-column prop="merchant_no" label="商户编号" width="150" />
                <el-table-column prop="merchant_name" label="商户名称" width="180" />
                <el-table-column prop="contact_phone" label="联系电话" width="140" />
                <el-table-column prop="terminal_count" label="终端数" width="80" />
                <el-table-column prop="status" label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                      {{ row.status === 1 ? '正常' : '禁用' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" min-width="160">
                  <template #default="{ row }">
                    {{ formatDate(row.created_at) }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="100" fixed="right">
                  <template #default="{ row }">
                    <el-button type="primary" link @click="viewMerchant(row.id)">查看</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无商户数据" />
              <el-pagination
                v-if="merchantsTotal > 10"
                v-model:current-page="merchantsPage"
                :page-size="10"
                :total="merchantsTotal"
                layout="prev, pager, next"
                @current-change="fetchMerchants"
              />
            </div>
          </el-tab-pane>

          <!-- 终端列表Tab -->
          <el-tab-pane label="终端列表" name="terminals">
            <div class="tab-content">
              <div class="tab-toolbar">
                <el-input
                  v-model="terminalsSearch"
                  placeholder="搜索终端SN/商户名称"
                  style="width: 250px"
                  clearable
                  @keyup.enter="fetchTerminals"
                >
                  <template #append>
                    <el-button @click="fetchTerminals">搜索</el-button>
                  </template>
                </el-input>
              </div>
              <el-table v-if="terminals.length > 0" :data="terminals" style="width: 100%">
                <el-table-column prop="terminal_sn" label="终端SN" width="150" />
                <el-table-column prop="merchant_name" label="商户名称" width="180" />
                <el-table-column prop="device_model" label="设备型号" width="120" />
                <el-table-column prop="status" label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="getTerminalStatusType(row.status)" size="small">
                      {{ getTerminalStatusName(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="activated_at" label="激活时间" min-width="160">
                  <template #default="{ row }">
                    {{ row.activated_at ? formatDate(row.activated_at) : '-' }}
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无终端数据" />
              <el-pagination
                v-if="terminalsTotal > 10"
                v-model:current-page="terminalsPage"
                :page-size="10"
                :total="terminalsTotal"
                layout="prev, pager, next"
                @current-change="fetchTerminals"
              />
            </div>
          </el-tab-pane>

          <!-- 交易记录Tab -->
          <el-tab-pane label="交易记录" name="transactions">
            <div class="tab-content">
              <div class="tab-toolbar">
                <el-date-picker
                  v-model="transactionDateRange"
                  type="daterange"
                  range-separator="至"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  value-format="YYYY-MM-DD"
                  @change="fetchTransactions"
                />
              </div>
              <el-table v-if="transactions.length > 0" :data="transactions" style="width: 100%">
                <el-table-column prop="order_no" label="订单号" width="180" />
                <el-table-column prop="merchant_name" label="商户名称" width="150" />
                <el-table-column prop="amount" label="交易金额" width="120">
                  <template #default="{ row }">
                    <span class="money">¥{{ formatAmount(row.amount) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="fee" label="手续费" width="100">
                  <template #default="{ row }">
                    ¥{{ formatAmount(row.fee) }}
                  </template>
                </el-table-column>
                <el-table-column prop="pay_type" label="支付方式" width="100">
                  <template #default="{ row }">
                    {{ getPayTypeName(row.pay_type) }}
                  </template>
                </el-table-column>
                <el-table-column prop="trade_time" label="交易时间" min-width="160">
                  <template #default="{ row }">
                    {{ formatDate(row.trade_time) }}
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无交易记录" />
              <el-pagination
                v-if="transactionsTotal > 10"
                v-model:current-page="transactionsPage"
                :page-size="10"
                :total="transactionsTotal"
                layout="prev, pager, next"
                @current-change="fetchTransactions"
              />
            </div>
          </el-tab-pane>

          <!-- 通道管理Tab -->
          <el-tab-pane label="通道管理" name="channels">
            <div v-loading="channelsLoading" class="tab-content">
              <div class="tab-toolbar">
                <el-button type="primary" @click="handleInitChannels">
                  初始化通道配置
                </el-button>
                <span class="toolbar-tip">
                  提示：初始化将为代理商配置所有可用的支付通道
                </span>
              </div>
              <el-table v-if="channels.length > 0" :data="channels" style="width: 100%">
                <el-table-column prop="channel_code" label="通道编码" width="140" />
                <el-table-column prop="channel_name" label="通道名称" width="160" />
                <el-table-column prop="is_enabled" label="启用状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.is_enabled ? 'success' : 'danger'" size="small">
                      {{ row.is_enabled ? '已启用' : '已禁用' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="is_visible" label="APP可见" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.is_visible ? 'success' : 'info'" size="small">
                      {{ row.is_visible ? '可见' : '隐藏' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="enabled_at" label="启用时间" width="160">
                  <template #default="{ row }">
                    {{ row.enabled_at ? formatDate(row.enabled_at) : '-' }}
                  </template>
                </el-table-column>
                <el-table-column prop="created_at" label="配置时间" min-width="160">
                  <template #default="{ row }">
                    {{ formatDate(row.created_at) }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="180" fixed="right">
                  <template #default="{ row }">
                    <el-button
                      :type="row.is_enabled ? 'warning' : 'success'"
                      link
                      @click="handleToggleChannelEnabled(row)"
                    >
                      {{ row.is_enabled ? '禁用' : '启用' }}
                    </el-button>
                    <el-button
                      :type="row.is_visible ? 'info' : 'primary'"
                      link
                      @click="handleToggleChannelVisibility(row)"
                    >
                      {{ row.is_visible ? '隐藏' : '显示' }}
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无通道配置，请点击「初始化通道配置」按钮" />
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <!-- 政策分配弹窗 -->
    <PolicyAssignDialog
      v-model="showPolicyAssignDialog"
      :agent-id="agentId"
      :agent-name="agentDetail?.agent_name || ''"
      @success="handlePolicyAssignSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import PolicyAssignDialog from '@/components/Policy/PolicyAssignDialog.vue'
import {
  getAgentDetail,
  getAgentPolicyList,
  getAgentWallets,
  getAgentSubordinates,
  getAgentMerchants,
  getAgentTerminals,
  getAgentTransactions,
  updateAgentStatus,
} from '@/api/agent'
import {
  getAgentChannels,
  enableChannel,
  disableChannel,
  setChannelVisibility,
  initAgentChannels,
} from '@/api/agent-channel'
import { formatAmount, formatDate } from '@/utils/format'
import type { AgentDetail, AgentPolicy, Wallet, Agent, Merchant, Terminal, Transaction, AgentChannel } from '@/types'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const agentDetail = ref<AgentDetail | null>(null)
const activeTab = ref('policy')
const showPolicyAssignDialog = ref(false)

// 获取代理商ID
const agentId = Number(route.params.id)

// Tab数据
const policies = ref<AgentPolicy[]>([])
const wallets = ref<Wallet[]>([])
const subordinates = ref<Agent[]>([])
const subordinatesSearch = ref('')
const subordinatesPage = ref(1)
const subordinatesTotal = ref(0)
const merchants = ref<Merchant[]>([])
const merchantsSearch = ref('')
const merchantsPage = ref(1)
const merchantsTotal = ref(0)
const terminals = ref<Terminal[]>([])
const terminalsSearch = ref('')
const terminalsPage = ref(1)
const terminalsTotal = ref(0)
const transactions = ref<Transaction[]>([])
const transactionDateRange = ref<string[]>([])
const transactionsPage = ref(1)
const transactionsTotal = ref(0)
const channels = ref<AgentChannel[]>([])
const channelsLoading = ref(false)

// 获取详情
async function fetchDetail() {
  loading.value = true
  try {
    agentDetail.value = await getAgentDetail(agentId) as AgentDetail
    // 默认加载政策数据
    await fetchPolicies()
  } catch (error) {
    console.error('Fetch agent detail error:', error)
    ElMessage.error('获取代理商详情失败')
  } finally {
    loading.value = false
  }
}

// 获取政策
async function fetchPolicies() {
  try {
    policies.value = await getAgentPolicyList(agentId)
  } catch (error) {
    console.error('Fetch policies error:', error)
  }
}

// 获取钱包
async function fetchWallets() {
  try {
    wallets.value = await getAgentWallets(agentId)
  } catch (error) {
    console.error('Fetch wallets error:', error)
  }
}

// 获取下级代理
async function fetchSubordinates() {
  try {
    const res = await getAgentSubordinates(agentId, {
      keyword: subordinatesSearch.value,
      page: subordinatesPage.value,
      page_size: 10,
    })
    subordinates.value = res.list as Agent[]
    subordinatesTotal.value = res.total
  } catch (error) {
    console.error('Fetch subordinates error:', error)
  }
}

// 获取商户
async function fetchMerchants() {
  try {
    const res = await getAgentMerchants(agentId, {
      keyword: merchantsSearch.value,
      page: merchantsPage.value,
      page_size: 10,
    })
    merchants.value = res.list as Merchant[]
    merchantsTotal.value = res.total
  } catch (error) {
    console.error('Fetch merchants error:', error)
  }
}

// 获取终端
async function fetchTerminals() {
  try {
    const res = await getAgentTerminals(agentId, {
      keyword: terminalsSearch.value,
      page: terminalsPage.value,
      page_size: 10,
    })
    terminals.value = res.list as Terminal[]
    terminalsTotal.value = res.total
  } catch (error) {
    console.error('Fetch terminals error:', error)
  }
}

// 获取交易
async function fetchTransactions() {
  try {
    const params: { page: number; page_size: number; start_date?: string; end_date?: string } = {
      page: transactionsPage.value,
      page_size: 10,
    }
    if (transactionDateRange.value?.length === 2) {
      params.start_date = transactionDateRange.value[0]
      params.end_date = transactionDateRange.value[1]
    }
    const res = await getAgentTransactions(agentId, params)
    transactions.value = res.list as Transaction[]
    transactionsTotal.value = res.total
  } catch (error) {
    console.error('Fetch transactions error:', error)
  }
}

// 获取通道配置
async function fetchChannels() {
  channelsLoading.value = true
  try {
    channels.value = await getAgentChannels(agentId)
  } catch (error) {
    console.error('Fetch channels error:', error)
  } finally {
    channelsLoading.value = false
  }
}

// 切换通道启用状态
async function handleToggleChannelEnabled(channel: AgentChannel) {
  const action = channel.is_enabled ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}「${channel.channel_name}」通道吗？`, '确认操作', {
      type: 'warning',
    })
    if (channel.is_enabled) {
      await disableChannel(agentId, channel.channel_id)
    } else {
      await enableChannel(agentId, channel.channel_id)
    }
    ElMessage.success(`${action}成功`)
    fetchChannels()
  } catch (error: unknown) {
    if ((error as string) !== 'cancel') {
      ElMessage.error(`${action}失败`)
    }
  }
}

// 切换通道可见性
async function handleToggleChannelVisibility(channel: AgentChannel) {
  const action = channel.is_visible ? '隐藏' : '显示'
  try {
    await setChannelVisibility(agentId, channel.channel_id, !channel.is_visible)
    ElMessage.success(`${action}成功`)
    fetchChannels()
  } catch (error) {
    ElMessage.error(`${action}失败`)
  }
}

// 初始化通道配置
async function handleInitChannels() {
  try {
    await ElMessageBox.confirm('确定要初始化通道配置吗？这将为代理商配置所有可用通道。', '确认操作', {
      type: 'info',
    })
    await initAgentChannels(agentId)
    ElMessage.success('初始化成功')
    fetchChannels()
  } catch (error: unknown) {
    if ((error as string) !== 'cancel') {
      ElMessage.error('初始化失败')
    }
  }
}

// Tab切换
function handleTabChange(tab: string | number) {
  switch (tab) {
    case 'policy':
      if (policies.value.length === 0) fetchPolicies()
      break
    case 'wallet':
      if (wallets.value.length === 0) fetchWallets()
      break
    case 'subordinates':
      if (subordinates.value.length === 0) fetchSubordinates()
      break
    case 'merchants':
      if (merchants.value.length === 0) fetchMerchants()
      break
    case 'terminals':
      if (terminals.value.length === 0) fetchTerminals()
      break
    case 'transactions':
      if (transactions.value.length === 0) fetchTransactions()
      break
    case 'channels':
      if (channels.value.length === 0) fetchChannels()
      break
  }
}

// 编辑
function handleEdit() {
  router.push(`/agents/${agentId}/edit`)
}

// 切换状态
async function handleToggleStatus() {
  if (!agentDetail.value) return
  const newStatus = agentDetail.value.status === 1 ? 2 : 1
  const action = newStatus === 2 ? '禁用' : '启用'

  try {
    await ElMessageBox.confirm(`确定要${action}该代理商吗？`, '确认操作', {
      type: 'warning',
    })
    await updateAgentStatus(agentId, newStatus)
    ElMessage.success(`${action}成功`)
    fetchDetail()
  } catch (error: unknown) {
    if ((error as string) !== 'cancel') {
      ElMessage.error(`${action}失败`)
    }
  }
}

// 查看代理商
function viewAgent(id: number) {
  router.push(`/agents/${id}`)
}

// 查看商户
function viewMerchant(id: number) {
  router.push(`/merchants/${id}`)
}

// 政策分配成功回调
function handlePolicyAssignSuccess() {
  fetchPolicies()
}

// 辅助函数
function getWalletTypeName(type: string): string {
  const map: Record<string, string> = {
    profit: '分润钱包',
    service: '服务费钱包',
    reward: '奖励钱包',
    recharge: '充值钱包',
    deposit: '沉淀钱包',
  }
  return map[type] || type
}

function getWalletTagType(type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    profit: 'primary',
    service: 'success',
    reward: 'warning',
    recharge: 'danger',
    deposit: 'info',
  }
  return map[type] || 'info'
}

function getTerminalStatusName(status: number): string {
  const map: Record<number, string> = {
    0: '未激活',
    1: '已激活',
    2: '已禁用',
    3: '已解绑',
  }
  return map[status] || '未知'
}

function getTerminalStatusType(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'info',
    1: 'success',
    2: 'danger',
    3: 'warning',
  }
  return map[status] || 'info'
}

function getPayTypeName(type: number): string {
  const map: Record<number, string> = {
    1: '刷卡',
    2: '微信',
    3: '支付宝',
    4: '云闪付',
  }
  return map[type] || '其他'
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
  }

  .tab-toolbar {
    margin-bottom: $spacing-md;
    display: flex;
    gap: $spacing-md;
    align-items: center;

    .toolbar-tip {
      font-size: 12px;
      color: $text-secondary;
    }
  }

  .el-pagination {
    margin-top: $spacing-md;
    justify-content: flex-end;
  }
}

.money {
  color: $success-color;
  font-weight: 600;
}

.frozen {
  color: $warning-color;
}

.available {
  color: $primary-color;
  font-weight: 600;
}
</style>
