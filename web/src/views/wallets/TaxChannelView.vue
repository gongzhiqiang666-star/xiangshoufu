<template>
  <div class="tax-channel-view">
    <PageHeader title="税筹通道" sub-title="税筹通道管理" />

    <!-- 税筹通道列表 -->
    <el-card class="list-card">
      <template #header>
        <div class="card-header">
          <span>税筹通道列表</span>
          <el-button type="primary" :icon="Plus" @click="handleCreate">
            新增通道
          </el-button>
        </div>
      </template>

      <el-table :data="taxChannels" v-loading="loading" border stripe>
        <el-table-column prop="channel_code" label="通道编码" width="120" />
        <el-table-column prop="channel_name" label="通道名称" width="150" />
        <el-table-column prop="fee_type_name" label="扣费类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.fee_type === 1 ? 'warning' : 'info'" size="small">
              {{ row.fee_type_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="tax_rate_percent" label="税率" width="100" align="right">
          <template #default="{ row }">
            {{ row.tax_rate_percent.toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column prop="fixed_fee_yuan" label="固定费用" width="120" align="right">
          <template #default="{ row }">
            ¥{{ row.fixed_fee_yuan.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status_name" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link @click="handleToggleStatus(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 税费计算器 -->
    <el-card class="calculator-card">
      <template #header>
        <div class="card-header">
          <span>税费计算器</span>
        </div>
      </template>

      <el-form :inline="true" :model="calcForm">
        <el-form-item label="支付通道">
          <el-input-number v-model="calcForm.channel_id" :min="1" placeholder="通道ID" style="width: 120px" />
        </el-form-item>
        <el-form-item label="钱包类型">
          <el-select v-model="calcForm.wallet_type" placeholder="选择钱包类型" style="width: 120px">
            <el-option label="分润钱包" :value="1" />
            <el-option label="服务费钱包" :value="2" />
            <el-option label="奖励钱包" :value="3" />
            <el-option label="充值钱包" :value="4" />
            <el-option label="沉淀钱包" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="calcForm.amount" :min="0.01" :precision="2" style="width: 150px" />
          <span class="form-tip">元</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleCalculate">计算税费</el-button>
        </el-form-item>
      </el-form>

      <el-descriptions v-if="calcResult" :column="3" border class="calc-result">
        <el-descriptions-item label="原金额">¥{{ (calcResult.original_amount / 100).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="税率">{{ (calcResult.tax_rate * 100).toFixed(2) }}%</el-descriptions-item>
        <el-descriptions-item label="税费">¥{{ (calcResult.tax_fee / 100).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="固定费用">¥{{ (calcResult.fixed_fee / 100).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="总费用">
          <span class="fee-amount">¥{{ (calcResult.total_fee / 100).toFixed(2) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="实际到账">
          <span class="actual-amount">¥{{ (calcResult.actual_amount / 100).toFixed(2) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="税筹通道" :span="3">
          {{ calcResult.tax_channel_name || '未配置' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑税筹通道' : '新增税筹通道'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="通道编码" required v-if="!isEdit">
          <el-input v-model="form.channel_code" placeholder="如: TAXCH01" />
        </el-form-item>
        <el-form-item label="通道名称" required>
          <el-input v-model="form.channel_name" placeholder="税筹通道名称" />
        </el-form-item>
        <el-form-item label="扣费类型" required>
          <el-radio-group v-model="form.fee_type">
            <el-radio :value="1">付款扣(充值时扣)</el-radio>
            <el-radio :value="2">出款扣(提现时扣)</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="税率" required>
          <el-input-number
            v-model="form.tax_rate"
            :min="0"
            :max="100"
            :precision="2"
            style="width: 150px"
          />
          <span class="form-tip">% (如9表示9%)</span>
        </el-form-item>
        <el-form-item label="固定费用">
          <el-input-number
            v-model="form.fixed_fee"
            :min="0"
            :precision="2"
            style="width: 150px"
          />
          <span class="form-tip">元/笔</span>
        </el-form-item>
        <el-form-item label="API地址">
          <el-input v-model="form.api_url" placeholder="可选" />
        </el-form-item>
        <el-form-item label="API密钥">
          <el-input v-model="form.api_key" placeholder="可选" type="password" show-password />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import {
  getTaxChannelList,
  createTaxChannel,
  updateTaxChannel,
  deleteTaxChannel,
  calculateWithdrawalTax,
} from '@/api/taxChannel'
import type { TaxChannel, TaxCalculationResult } from '@/types'

// 列表
const taxChannels = ref<TaxChannel[]>([])
const loading = ref(false)

// 弹窗
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const form = reactive({
  channel_code: '',
  channel_name: '',
  fee_type: 2 as 1 | 2,
  tax_rate: 9,
  fixed_fee: 0,
  api_url: '',
  api_key: '',
  api_secret: '',
  remark: '',
})

// 计算器
const calcForm = reactive({
  channel_id: 1,
  wallet_type: 1,
  amount: 100,
})
const calcResult = ref<TaxCalculationResult | null>(null)

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await getTaxChannelList()
    taxChannels.value = res.list || []
  } catch (error) {
    console.error('Fetch list error:', error)
  } finally {
    loading.value = false
  }
}

// 新增
function handleCreate() {
  isEdit.value = false
  editId.value = null
  form.channel_code = ''
  form.channel_name = ''
  form.fee_type = 2
  form.tax_rate = 9
  form.fixed_fee = 0
  form.api_url = ''
  form.api_key = ''
  form.api_secret = ''
  form.remark = ''
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: TaxChannel) {
  isEdit.value = true
  editId.value = row.id
  form.channel_code = row.channel_code
  form.channel_name = row.channel_name
  form.fee_type = row.fee_type
  form.tax_rate = row.tax_rate_percent
  form.fixed_fee = row.fixed_fee_yuan
  form.api_url = ''
  form.api_key = ''
  form.api_secret = ''
  form.remark = row.remark
  dialogVisible.value = true
}

// 提交
async function handleSubmit() {
  if (!form.channel_name) {
    ElMessage.warning('请输入通道名称')
    return
  }
  if (form.tax_rate < 0 || form.tax_rate > 100) {
    ElMessage.warning('税率必须在0-100之间')
    return
  }

  try {
    if (isEdit.value && editId.value) {
      await updateTaxChannel(editId.value, {
        channel_name: form.channel_name,
        fee_type: form.fee_type,
        tax_rate: form.tax_rate / 100,
        fixed_fee: Math.round(form.fixed_fee * 100),
        api_url: form.api_url || undefined,
        api_key: form.api_key || undefined,
        api_secret: form.api_secret || undefined,
        remark: form.remark || undefined,
      })
      ElMessage.success('更新成功')
    } else {
      if (!form.channel_code) {
        ElMessage.warning('请输入通道编码')
        return
      }
      await createTaxChannel({
        channel_code: form.channel_code,
        channel_name: form.channel_name,
        fee_type: form.fee_type,
        tax_rate: form.tax_rate / 100,
        fixed_fee: Math.round(form.fixed_fee * 100),
        api_url: form.api_url || undefined,
        api_key: form.api_key || undefined,
        api_secret: form.api_secret || undefined,
        remark: form.remark || undefined,
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('Submit error:', error)
  }
}

// 切换状态
async function handleToggleStatus(row: TaxChannel) {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'

  try {
    await ElMessageBox.confirm(`确定${action}税筹通道 "${row.channel_name}" 吗?`, '提示', {
      type: 'warning',
    })

    await updateTaxChannel(row.id, {
      status: newStatus as 0 | 1,
    })
    ElMessage.success(`${action}成功`)
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Toggle status error:', error)
    }
  }
}

// 删除
async function handleDelete(row: TaxChannel) {
  try {
    await ElMessageBox.confirm(`确定删除税筹通道 "${row.channel_name}" 吗?`, '警告', {
      type: 'warning',
    })

    await deleteTaxChannel(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
    }
  }
}

// 计算税费
async function handleCalculate() {
  if (calcForm.amount <= 0) {
    ElMessage.warning('请输入金额')
    return
  }

  try {
    calcResult.value = await calculateWithdrawalTax({
      channel_id: calcForm.channel_id,
      wallet_type: calcForm.wallet_type,
      amount: Math.round(calcForm.amount * 100),
    })
  } catch (error) {
    console.error('Calculate error:', error)
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style lang="scss" scoped>
.tax-channel-view {
  padding: 0;
}

.list-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.calculator-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .calc-result {
    margin-top: $spacing-md;
  }

  .fee-amount {
    color: $warning-color;
    font-weight: 600;
  }

  .actual-amount {
    color: $success-color;
    font-weight: 600;
    font-size: 16px;
  }
}

.form-tip {
  margin-left: $spacing-sm;
  color: $text-secondary;
}
</style>
