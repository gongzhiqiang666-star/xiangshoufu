<template>
  <div class="adjustment-list-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="代理商">
        <AgentSelect v-model="searchForm.agent_id" placeholder="请选择代理商" style="width: 180px" />
      </el-form-item>
      <el-form-item label="钱包类型">
        <el-select v-model="searchForm.wallet_type" placeholder="请选择" clearable style="width: 120px">
          <el-option label="分润钱包" :value="1" />
          <el-option label="服务费钱包" :value="2" />
          <el-option label="奖励钱包" :value="3" />
          <el-option label="充值钱包" :value="4" />
          <el-option label="沉淀钱包" :value="5" />
        </el-select>
      </el-form-item>
      <el-form-item label="日期范围">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
        />
      </el-form-item>

      <template #extra>
        <el-button type="primary" :icon="Plus" @click="showCreateDialog">新增调账</el-button>
      </template>
    </SearchForm>

    <!-- 表格 -->
    <ProTable
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      @refresh="fetchData"
    >
      <el-table-column prop="adjustment_no" label="调账单号" width="180" />
      <el-table-column prop="agent_name" label="代理商" min-width="120" show-overflow-tooltip />
      <el-table-column prop="wallet_type_name" label="钱包类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getWalletTypeTag(row.wallet_type)" size="small">
            {{ row.wallet_type_name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="adjustment_type" label="类型" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.amount >= 0 ? 'success' : 'danger'" size="small">
            {{ row.adjustment_type }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount_yuan" label="金额" width="120" align="right">
        <template #default="{ row }">
          <span :class="['amount', row.amount >= 0 ? 'income' : 'expense']">
            {{ row.amount >= 0 ? '+' : '' }}¥{{ row.amount_yuan.toFixed(2) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="balance_before_yuan" label="调账前" width="110" align="right">
        <template #default="{ row }">
          ¥{{ row.balance_before_yuan.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="balance_after_yuan" label="调账后" width="110" align="right">
        <template #default="{ row }">
          ¥{{ row.balance_after_yuan.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="原因" min-width="150" show-overflow-tooltip />
      <el-table-column prop="operator_name" label="操作人" width="100" />
      <el-table-column prop="status_name" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status_name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="170" />
    </ProTable>

    <!-- 新增调账对话框 -->
    <el-dialog v-model="dialogVisible" title="新增调账" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="代理商" prop="agent_id">
          <AgentSelect v-model="form.agent_id" placeholder="请选择代理商" style="width: 100%" />
        </el-form-item>
        <el-form-item label="钱包类型" prop="wallet_type">
          <el-select v-model="form.wallet_type" placeholder="请选择钱包类型" style="width: 100%">
            <el-option label="分润钱包" :value="1" />
            <el-option label="服务费钱包" :value="2" />
            <el-option label="奖励钱包" :value="3" />
            <el-option label="充值钱包" :value="4" />
            <el-option label="沉淀钱包" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="调账类型" prop="adjustment_type">
          <el-radio-group v-model="adjustmentType">
            <el-radio value="add">充入</el-radio>
            <el-radio value="deduct">扣减</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="金额(元)" prop="amount_yuan">
          <el-input-number
            v-model="form.amount_yuan"
            :min="0.01"
            :precision="2"
            :step="100"
            placeholder="请输入金额"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="调账原因" prop="reason">
          <el-input
            v-model="form.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入调账原因"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import { getWalletAdjustments, createWalletAdjustment } from '@/api/wallet'
import type { WalletAdjustment, CreateAdjustmentParams } from '@/types'

// 搜索表单
const searchForm = reactive({
  agent_id: undefined as number | undefined,
  wallet_type: undefined as number | undefined,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// 表格数据
const tableData = ref<WalletAdjustment[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// 对话框
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const submitting = ref(false)
const adjustmentType = ref<'add' | 'deduct'>('add')

const form = reactive({
  agent_id: undefined as number | undefined,
  wallet_type: undefined as number | undefined,
  amount_yuan: undefined as number | undefined,
  reason: '',
})

const formRules: FormRules = {
  agent_id: [{ required: true, message: '请选择代理商', trigger: 'change' }],
  wallet_type: [{ required: true, message: '请选择钱包类型', trigger: 'change' }],
  amount_yuan: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  reason: [{ required: true, message: '请输入调账原因', trigger: 'blur' }],
}

// 获取钱包类型标签
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getWalletTypeTag(type: number): TagType {
  const tagMap: Record<number, TagType> = {
    1: 'primary', // 分润
    2: 'success', // 服务费
    3: 'warning', // 奖励
    4: 'danger',  // 充值
    5: 'info',    // 沉淀
  }
  return tagMap[type] || 'info'
}

// 处理日期变化
function handleDateChange(val: [string, string] | null) {
  if (val) {
    searchForm.start_date = val[0]
    searchForm.end_date = val[1]
  } else {
    searchForm.start_date = ''
    searchForm.end_date = ''
  }
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getWalletAdjustments({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list || []
    total.value = res.total
  } catch (error) {
    console.error('Fetch adjustments error:', error)
  } finally {
    loading.value = false
  }
}

// 监听分页变化
watch([page, pageSize], () => {
  fetchData()
})

// 搜索
function handleSearch() {
  page.value = 1
  fetchData()
}

// 重置
function handleReset() {
  dateRange.value = null
  page.value = 1
  fetchData()
}

// 显示新增对话框
function showCreateDialog() {
  form.agent_id = undefined
  form.wallet_type = undefined
  form.amount_yuan = undefined
  form.reason = ''
  adjustmentType.value = 'add'
  dialogVisible.value = true
}

// 提交
async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  if (!form.agent_id || !form.wallet_type || !form.amount_yuan) {
    return
  }

  // 计算金额（元转分，扣减为负数）
  let amountFen = Math.round(form.amount_yuan * 100)
  if (adjustmentType.value === 'deduct') {
    amountFen = -amountFen
  }

  const params: CreateAdjustmentParams = {
    agent_id: form.agent_id,
    wallet_type: form.wallet_type,
    channel_id: 0,
    amount: amountFen,
    reason: form.reason,
  }

  submitting.value = true
  try {
    await createWalletAdjustment(params)
    ElMessage.success('调账成功')
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '调账失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.adjustment-list-view {
  padding: 0;
}

.amount {
  font-weight: 600;

  &.income {
    color: $success-color;
  }

  &.expense {
    color: $danger-color;
  }
}
</style>
