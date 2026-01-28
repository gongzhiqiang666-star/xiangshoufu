<template>
  <div class="settlement-price-list">
    <div class="page-header">
      <h2>结算价管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增结算价
      </el-button>
    </div>

    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="代理商ID">
          <el-input v-model.number="searchForm.agent_id" placeholder="请输入代理商ID" clearable />
        </el-form-item>
        <el-form-item label="通道ID">
          <el-input v-model.number="searchForm.channel_id" placeholder="请输入通道ID" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="agent_id" label="代理商ID" width="100" />
        <el-table-column prop="agent_name" label="代理商名称" min-width="120" />
        <el-table-column prop="channel_id" label="通道ID" width="80" />
        <el-table-column prop="channel_name" label="通道名称" min-width="100" />
        <el-table-column label="费率配置" min-width="200">
          <template #default="{ row }">
            <div v-if="row.rate_configs && Object.keys(row.rate_configs).length">
              <el-tag v-for="(config, key) in row.rate_configs" :key="key" size="small" class="rate-tag">
                {{ key }}: {{ config.rate }}%
              </el-tag>
            </div>
            <span v-else class="text-muted">未配置</span>
          </template>
        </el-table-column>
        <el-table-column label="押金返现" min-width="150">
          <template #default="{ row }">
            <div v-if="row.deposit_cashbacks && row.deposit_cashbacks.length">
              <el-tag v-for="(dc, idx) in row.deposit_cashbacks" :key="idx" size="small" type="success" class="deposit-tag">
                ¥{{ (dc.deposit_amount / 100).toFixed(0) }} → ¥{{ (dc.cashback_amount / 100).toFixed(0) }}
              </el-tag>
            </div>
            <span v-else class="text-muted">未配置</span>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="70" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link size="small" @click="handleViewLogs(row)">调价记录</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 新增对话框：选择代理商+通道+模板 -->
    <el-dialog v-model="createDialogVisible" title="新增结算价" width="800px" destroy-on-close>
      <el-form :model="createForm" label-width="100px">
        <!-- Step 1: 选择代理商 -->
        <el-form-item label="代理商" required>
          <el-select
            v-model="createForm.agent_id"
            filterable
            remote
            reserve-keyword
            placeholder="请输入代理商名称搜索"
            :remote-method="handleAgentSearch"
            :loading="agentSearchLoading"
            style="width: 100%"
          >
            <el-option
              v-for="agent in agentOptions"
              :key="agent.id"
              :label="`${agent.name} (ID: ${agent.id})`"
              :value="agent.id"
            />
          </el-select>
        </el-form-item>

        <!-- Step 2: 选择通道 -->
        <el-form-item label="通道" required>
          <el-select
            v-model="createForm.channel_id"
            placeholder="请选择通道"
            style="width: 100%"
            @change="handleChannelChange"
          >
            <el-option
              v-for="channel in channelList"
              :key="channel.id"
              :label="channel.channel_name"
              :value="channel.id"
            />
          </el-select>
        </el-form-item>

        <!-- Step 3: 选择政策模板（通道选择后显示） -->
        <el-form-item v-if="createForm.channel_id" label="政策模板" required>
          <el-select
            v-model="createForm.template_id"
            placeholder="请选择政策模板"
            :loading="templateLoading"
            style="width: 100%"
            @change="handleTemplateChange"
          >
            <el-option
              v-for="tpl in templateList"
              :key="tpl.id"
              :label="`${tpl.name}${tpl.is_default ? ' (默认)' : ''}`"
              :value="tpl.id"
            />
          </el-select>
        </el-form-item>

        <!-- Step 4: 模板配置预览（模板选择后显示） -->
        <template v-if="templateDetail">
          <el-divider content-position="left">模板配置预览（只读）</el-divider>

          <!-- 费率配置 -->
          <el-descriptions title="费率配置" :column="3" border size="small">
            <el-descriptions-item
              v-for="(config, key) in templateDetail.rate_configs"
              :key="key"
              :label="key"
            >
              {{ config.rate }}%
            </el-descriptions-item>
          </el-descriptions>

          <!-- 押金返现配置 -->
          <el-descriptions
            v-if="templateDetail.deposit_cashbacks && templateDetail.deposit_cashbacks.length"
            title="押金返现"
            :column="2"
            border
            size="small"
            style="margin-top: 16px"
          >
            <el-descriptions-item
              v-for="(item, idx) in templateDetail.deposit_cashbacks"
              :key="idx"
              :label="`¥${(item.deposit_amount/100).toFixed(0)}押金`"
            >
              返现 ¥{{ (item.cashback_amount/100).toFixed(0) }}
            </el-descriptions-item>
          </el-descriptions>

          <!-- 流量费返现 -->
          <el-descriptions
            v-if="templateDetail.sim_cashbacks && templateDetail.sim_cashbacks.length"
            title="流量费返现"
            :column="3"
            border
            size="small"
            style="margin-top: 16px"
          >
            <el-descriptions-item
              v-for="(item, idx) in templateDetail.sim_cashbacks"
              :key="idx"
              :label="getSimCashbackLabel(item.type)"
            >
              ¥{{ (item.cashback_amount/100).toFixed(0) }}
            </el-descriptions-item>
          </el-descriptions>

          <!-- 激活奖励 -->
          <el-descriptions
            v-if="templateDetail.activation_rewards && templateDetail.activation_rewards.length"
            title="激活奖励"
            :column="1"
            border
            size="small"
            style="margin-top: 16px"
          >
            <el-descriptions-item
              v-for="(reward, idx) in templateDetail.activation_rewards"
              :key="idx"
              :label="`${reward.start_day}-${reward.end_day}天`"
            >
              达标 ¥{{ (reward.target_amount/100).toFixed(0) }} → 奖励 ¥{{ (reward.reward_amount/100).toFixed(0) }}
            </el-descriptions-item>
          </el-descriptions>

          <el-alert type="info" :closable="false" style="margin-top: 16px">
            创建后可在编辑页面调整以上配置
          </el-alert>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="handleConfirmCreate"
          :loading="creating"
          :disabled="!createForm.agent_id || !createForm.channel_id || !createForm.template_id"
        >
          确认导入
        </el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑结算价" width="700px">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="代理商">
          <el-input :value="`${editForm.agent_name} (ID: ${editForm.agent_id})`" disabled />
        </el-form-item>
        <el-form-item label="通道">
          <el-input :value="`${editForm.channel_name} (ID: ${editForm.channel_id})`" disabled />
        </el-form-item>
        <el-divider content-position="left">费率配置</el-divider>
        <el-form-item label="贷记卡费率">
          <el-input v-model="editForm.credit_rate" placeholder="如: 0.60">
            <template #append>%</template>
          </el-input>
        </el-form-item>
        <el-form-item label="借记卡费率">
          <el-input v-model="editForm.debit_rate" placeholder="如: 0.50">
            <template #append>%</template>
          </el-input>
        </el-form-item>
        <el-divider content-position="left">押金返现配置</el-divider>
        <div v-for="(dc, idx) in editForm.deposit_cashbacks" :key="idx" class="deposit-item">
          <el-row :gutter="10">
            <el-col :span="10">
              <el-form-item label="押金金额(分)">
                <el-input-number v-model="dc.deposit_amount" :min="0" />
              </el-form-item>
            </el-col>
            <el-col :span="10">
              <el-form-item label="返现金额(分)">
                <el-input-number v-model="dc.cashback_amount" :min="0" />
              </el-form-item>
            </el-col>
            <el-col :span="4">
              <el-button type="danger" :icon="Delete" circle @click="removeDeposit(idx)" />
            </el-col>
          </el-row>
        </div>
        <el-button type="primary" link @click="addDeposit">+ 添加押金返现配置</el-button>
        <el-divider content-position="left">流量费返现配置</el-divider>
        <el-form-item label="首次返现(分)">
          <el-input-number v-model="editForm.sim_first_cashback" :min="0" />
        </el-form-item>
        <el-form-item label="第2次返现(分)">
          <el-input-number v-model="editForm.sim_second_cashback" :min="0" />
        </el-form-item>
        <el-form-item label="第3次+返现(分)">
          <el-input-number v-model="editForm.sim_third_plus_cashback" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <!-- 调价记录对话框 -->
    <el-dialog v-model="logDialogVisible" title="调价记录" width="900px">
      <el-table :data="changeLogs" v-loading="logsLoading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="change_type_name" label="变更类型" width="120" />
        <el-table-column prop="change_summary" label="变更摘要" min-width="150" />
        <el-table-column prop="operator_name" label="操作人" width="100" />
        <el-table-column prop="source" label="来源" width="80" />
        <el-table-column prop="created_at" label="操作时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import {
  getSettlementPrices,
  getSettlementPrice,
  createSettlementPrice,
  updateSettlementPriceRate,
  updateSettlementPriceDeposit,
  updateSettlementPriceSim,
  getSettlementPriceChangeLogs,
  type SettlementPriceItem,
  type PriceChangeLog,
  type DepositCashbackItem,
} from '@/api/settlementPrice'
import { searchAgentList } from '@/api/agent'
import { getChannelList } from '@/api/channel'
import { getPolicyTemplates, getPolicyTemplateDetail } from '@/api/policy'
import type { PolicyTemplate, PolicyTemplateDetail, Channel, Agent } from '@/types'

// 搜索表单
const searchForm = reactive({
  agent_id: undefined as number | undefined,
  channel_id: undefined as number | undefined,
  status: undefined as number | undefined,
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

// 数据
const tableData = ref<SettlementPriceItem[]>([])
const loading = ref(false)
const saving = ref(false)

// ============= 新增相关状态 =============
const createDialogVisible = ref(false)
const creating = ref(false)
const channelList = ref<Channel[]>([])
const agentOptions = ref<Agent[]>([])
const agentSearchLoading = ref(false)
const templateList = ref<PolicyTemplate[]>([])
const templateLoading = ref(false)
const templateDetail = ref<PolicyTemplateDetail | null>(null)

const createForm = reactive({
  agent_id: null as number | null,
  channel_id: null as number | null,
  template_id: null as number | null,
})

// ============= 编辑相关状态 =============
const editDialogVisible = ref(false)
const editForm = reactive({
  id: 0,
  agent_id: 0,
  agent_name: '',
  channel_id: 0,
  channel_name: '',
  credit_rate: '',
  debit_rate: '',
  deposit_cashbacks: [] as DepositCashbackItem[],
  sim_first_cashback: 0,
  sim_second_cashback: 0,
  sim_third_plus_cashback: 0,
})

// 调价记录对话框
const logDialogVisible = ref(false)
const changeLogs = ref<PriceChangeLog[]>([])
const logsLoading = ref(false)

// 格式化日期时间
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 获取流量返现类型标签
const getSimCashbackLabel = (type: string) => {
  const labels: Record<string, string> = {
    'first': '首次',
    'second': '第2次',
    'renewal': '第3次+',
  }
  return labels[type] || type
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const resp = await getSettlementPrices({
      agent_id: searchForm.agent_id,
      channel_id: searchForm.channel_id,
      status: searchForm.status,
      page: pagination.page,
      page_size: pagination.pageSize,
    })
    tableData.value = resp.list || []
    pagination.total = resp.total
  } catch (e) {
    console.error('加载结算价列表失败', e)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  searchForm.agent_id = undefined
  searchForm.channel_id = undefined
  searchForm.status = undefined
  handleSearch()
}

// 分页
const handleSizeChange = () => {
  pagination.page = 1
  loadData()
}
const handleCurrentChange = () => {
  loadData()
}

// ============= 新增功能 =============

// 打开新增对话框
const handleCreate = async () => {
  // 重置表单
  createForm.agent_id = null
  createForm.channel_id = null
  createForm.template_id = null
  templateDetail.value = null
  templateList.value = []
  agentOptions.value = []

  // 加载通道列表
  try {
    channelList.value = await getChannelList()
  } catch (e) {
    ElMessage.error('加载通道列表失败')
    return
  }
  createDialogVisible.value = true
}

// 搜索代理商
const handleAgentSearch = async (keyword: string) => {
  if (!keyword || keyword.length < 2) {
    agentOptions.value = []
    return
  }
  agentSearchLoading.value = true
  try {
    const res = await searchAgentList({ keyword, page_size: 20 })
    agentOptions.value = res.list || []
  } catch (e) {
    console.error('搜索代理商失败', e)
  } finally {
    agentSearchLoading.value = false
  }
}

// 通道变更
const handleChannelChange = async (channelId: number) => {
  createForm.template_id = null
  templateDetail.value = null

  if (!channelId) {
    templateList.value = []
    return
  }

  templateLoading.value = true
  try {
    const res = await getPolicyTemplates({ channel_id: channelId })
    templateList.value = res.list || []
  } catch (e) {
    ElMessage.error('获取政策模板列表失败')
  } finally {
    templateLoading.value = false
  }
}

// 模板变更
const handleTemplateChange = async (templateId: number) => {
  if (!templateId) {
    templateDetail.value = null
    return
  }
  try {
    templateDetail.value = await getPolicyTemplateDetail(templateId)
  } catch (e) {
    ElMessage.error('获取模板详情失败')
  }
}

// 确认创建
const handleConfirmCreate = async () => {
  if (!createForm.agent_id || !createForm.channel_id || !createForm.template_id) {
    ElMessage.warning('请完整填写信息')
    return
  }

  creating.value = true
  try {
    await createSettlementPrice({
      agent_id: createForm.agent_id,
      channel_id: createForm.channel_id,
      template_id: createForm.template_id,
    })
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e?.message || '创建失败')
  } finally {
    creating.value = false
  }
}

// ============= 编辑功能 =============

// 编辑
const handleEdit = async (row: SettlementPriceItem) => {
  try {
    const detail = await getSettlementPrice(row.id)
    Object.assign(editForm, {
      id: detail.id,
      agent_id: detail.agent_id,
      agent_name: row.agent_name || '',
      channel_id: detail.channel_id,
      channel_name: row.channel_name || '',
      credit_rate: detail.credit_rate || '',
      debit_rate: detail.debit_rate || '',
      deposit_cashbacks: detail.deposit_cashbacks || [],
      sim_first_cashback: detail.sim_first_cashback,
      sim_second_cashback: detail.sim_second_cashback,
      sim_third_plus_cashback: detail.sim_third_plus_cashback,
    })
    editDialogVisible.value = true
  } catch (e) {
    ElMessage.error('获取结算价详情失败')
  }
}

// 保存
const handleSave = async () => {
  saving.value = true
  try {
    // 更新费率
    if (editForm.credit_rate || editForm.debit_rate) {
      await updateSettlementPriceRate(editForm.id, {
        credit_rate: editForm.credit_rate || undefined,
        debit_rate: editForm.debit_rate || undefined,
      })
    }
    // 更新押金返现
    if (editForm.deposit_cashbacks.length) {
      await updateSettlementPriceDeposit(editForm.id, {
        deposit_cashbacks: editForm.deposit_cashbacks,
      })
    }
    // 更新流量费返现
    await updateSettlementPriceSim(editForm.id, {
      sim_first_cashback: editForm.sim_first_cashback,
      sim_second_cashback: editForm.sim_second_cashback,
      sim_third_plus_cashback: editForm.sim_third_plus_cashback,
    })
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 添加押金返现配置
const addDeposit = () => {
  editForm.deposit_cashbacks.push({ deposit_amount: 0, cashback_amount: 0 })
}

// 删除押金返现配置
const removeDeposit = (idx: number) => {
  editForm.deposit_cashbacks.splice(idx, 1)
}

// 查看调价记录
const handleViewLogs = async (row: SettlementPriceItem) => {
  logDialogVisible.value = true
  logsLoading.value = true
  try {
    const resp = await getSettlementPriceChangeLogs(row.id)
    changeLogs.value = resp.list || []
  } catch (e) {
    ElMessage.error('获取调价记录失败')
  } finally {
    logsLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.settlement-price-list {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.rate-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.deposit-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.text-muted {
  color: #909399;
}

.deposit-item {
  margin-bottom: 10px;
}
</style>
