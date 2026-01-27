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

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" :title="editForm.id ? '编辑结算价' : '新增结算价'" width="700px">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="代理商ID" required>
          <el-input v-model.number="editForm.agent_id" :disabled="!!editForm.id" />
        </el-form-item>
        <el-form-item label="通道ID" required>
          <el-input v-model.number="editForm.channel_id" :disabled="!!editForm.id" />
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

// 编辑对话框
const editDialogVisible = ref(false)
const editForm = reactive({
  id: 0,
  agent_id: 0,
  channel_id: 0,
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

// 新增
const handleCreate = () => {
  Object.assign(editForm, {
    id: 0,
    agent_id: 0,
    channel_id: 0,
    credit_rate: '',
    debit_rate: '',
    deposit_cashbacks: [],
    sim_first_cashback: 0,
    sim_second_cashback: 0,
    sim_third_plus_cashback: 0,
  })
  editDialogVisible.value = true
}

// 编辑
const handleEdit = async (row: SettlementPriceItem) => {
  try {
    const detail = await getSettlementPrice(row.id)
    Object.assign(editForm, {
      id: detail.id,
      agent_id: detail.agent_id,
      channel_id: detail.channel_id,
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
    if (editForm.id) {
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
    } else {
      await createSettlementPrice({
        agent_id: editForm.agent_id,
        channel_id: editForm.channel_id,
      })
      ElMessage.success('创建成功')
    }
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
