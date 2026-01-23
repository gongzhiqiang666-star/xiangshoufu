<template>
  <div class="settlement-wallet-view">
    <PageHeader title="沉淀钱包" sub-title="沉淀钱包管理" />

    <!-- 钱包汇总 -->
    <el-card class="summary-card" v-loading="summaryLoading">
      <template #header>
        <div class="card-header">
          <span>沉淀钱包概览</span>
          <el-tag v-if="summary.settlement_ratio > 0" type="info">
            沉淀比例: {{ summary.settlement_ratio }}%
          </el-tag>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="stat-item">
            <div class="stat-label">下级未提现总额</div>
            <div class="stat-value primary">¥{{ formatAmount(summary.subordinate_total_balance) }}</div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="stat-item">
            <div class="stat-label">可用沉淀额度</div>
            <div class="stat-value success">¥{{ formatAmount(summary.available_amount) }}</div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="stat-item">
            <div class="stat-label">已使用沉淀额</div>
            <div class="stat-value warning">¥{{ formatAmount(summary.used_amount) }}</div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="stat-item">
            <div class="stat-label">待归还金额</div>
            <div class="stat-value danger">¥{{ formatAmount(summary.pending_return_amount) }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 操作区 -->
    <el-card class="action-card">
      <template #header>
        <div class="card-header">
          <span>快捷操作</span>
        </div>
      </template>

      <el-space wrap>
        <el-button type="primary" :icon="Download" @click="handleUseSettlement">
          使用沉淀款
        </el-button>
        <el-button type="success" :icon="Upload" @click="handleReturnSettlement">
          归还沉淀款
        </el-button>
      </el-space>

      <el-alert
        v-if="remainingAmount > 0"
        type="info"
        :closable="false"
        class="remaining-tip"
      >
        当前可用沉淀额度: ¥{{ formatAmount(remainingAmount) }}
      </el-alert>
    </el-card>

    <!-- 下级余额明细 -->
    <el-card class="subordinate-card">
      <template #header>
        <div class="card-header">
          <span>下级余额明细</span>
          <el-button type="primary" link :icon="Refresh" @click="fetchSubordinates">刷新</el-button>
        </div>
      </template>

      <el-table :data="subordinates" v-loading="subordinatesLoading" border stripe>
        <el-table-column prop="agent_id" label="代理ID" width="100" />
        <el-table-column prop="agent_name" label="代理名称" min-width="150" />
        <el-table-column prop="available_balance_yuan" label="可用余额" width="150" align="right">
          <template #default="{ row }">
            <span class="balance-amount">¥{{ row.available_balance_yuan.toFixed(2) }}</span>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="subordinates.length === 0 && !subordinatesLoading" description="暂无下级余额数据" />
    </el-card>

    <!-- 使用/归还记录 -->
    <el-card class="usage-card">
      <template #header>
        <div class="card-header">
          <span>使用记录</span>
          <el-radio-group v-model="usageType" @change="fetchUsages">
            <el-radio-button :value="undefined">全部</el-radio-button>
            <el-radio-button :value="1">使用</el-radio-button>
            <el-radio-button :value="2">归还</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="usages" v-loading="usagesLoading" border stripe>
        <el-table-column prop="usage_no" label="单号" width="200" />
        <el-table-column prop="usage_type_name" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="row.usage_type === 1 ? 'primary' : 'success'" size="small">
              {{ row.usage_type_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount_yuan" label="金额" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.usage_type === 1 ? 'amount-use' : 'amount-return'">
              {{ row.usage_type === 1 ? '-' : '+' }}¥{{ row.amount_yuan.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status_name" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'warning'" size="small">
              {{ row.status_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="returned_at" label="归还时间" width="180" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      </el-table>

      <el-pagination
        v-if="usageTotal > 0"
        class="pagination"
        background
        layout="total, sizes, prev, pager, next"
        :total="usageTotal"
        v-model:current-page="usagePage"
        v-model:page-size="usagePageSize"
        :page-sizes="[10, 20, 50]"
        @change="fetchUsages"
      />
    </el-card>

    <!-- 使用沉淀款弹窗 -->
    <el-dialog v-model="useDialogVisible" title="使用沉淀款" width="500px">
      <el-form :model="useForm" label-width="100px">
        <el-form-item label="可用额度">
          <span class="available-amount">¥{{ formatAmount(remainingAmount) }}</span>
        </el-form-item>
        <el-form-item label="使用金额" required>
          <el-input-number
            v-model="useForm.amount"
            :min="1"
            :max="remainingAmount / 100"
            :precision="2"
            style="width: 200px"
          />
          <span class="form-tip">元 (最少1元)</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="useForm.remark" type="textarea" rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="useDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitUse">确认使用</el-button>
      </template>
    </el-dialog>

    <!-- 归还沉淀款弹窗 -->
    <el-dialog v-model="returnDialogVisible" title="归还沉淀款" width="500px">
      <el-form :model="returnForm" label-width="100px">
        <el-form-item label="待归还总额">
          <span class="pending-amount">¥{{ formatAmount(summary.pending_return_amount) }}</span>
        </el-form-item>
        <el-form-item label="归还金额" required>
          <el-input-number
            v-model="returnForm.amount"
            :min="0.01"
            :precision="2"
            style="width: 200px"
          />
          <span class="form-tip">元</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="returnForm.remark" type="textarea" rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="returnDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitReturn">确认归还</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Download, Upload, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import {
  getSettlementWalletSummary,
  getSubordinateBalances,
  getSettlementUsageList,
  useSettlement,
  returnSettlement,
} from '@/api/settlementWallet'
import { formatAmount } from '@/utils/format'
import type {
  SettlementWalletSummary,
  SubordinateBalance,
  SettlementWalletUsage,
  SettlementUsageType,
} from '@/types'

// 汇总
const summary = ref<SettlementWalletSummary>({
  subordinate_total_balance: 0,
  subordinate_total_balance_yuan: 0,
  settlement_ratio: 0,
  available_amount: 0,
  available_amount_yuan: 0,
  used_amount: 0,
  used_amount_yuan: 0,
  pending_return_amount: 0,
  pending_return_amount_yuan: 0,
})
const summaryLoading = ref(false)

// 剩余可用额度
const remainingAmount = computed(() => {
  return Math.max(0, summary.value.available_amount - summary.value.used_amount)
})

// 下级余额
const subordinates = ref<SubordinateBalance[]>([])
const subordinatesLoading = ref(false)

// 使用记录
const usages = ref<SettlementWalletUsage[]>([])
const usagesLoading = ref(false)
const usageType = ref<SettlementUsageType | undefined>(undefined)
const usagePage = ref(1)
const usagePageSize = ref(10)
const usageTotal = ref(0)

// 使用弹窗
const useDialogVisible = ref(false)
const useForm = reactive({
  amount: 100,
  remark: '',
})

// 归还弹窗
const returnDialogVisible = ref(false)
const returnForm = reactive({
  amount: 100,
  remark: '',
})

// 获取汇总
async function fetchSummary() {
  summaryLoading.value = true
  try {
    summary.value = await getSettlementWalletSummary()
  } catch (error) {
    console.error('Fetch summary error:', error)
  } finally {
    summaryLoading.value = false
  }
}

// 获取下级余额
async function fetchSubordinates() {
  subordinatesLoading.value = true
  try {
    const res = await getSubordinateBalances()
    subordinates.value = res.list || []
  } catch (error) {
    console.error('Fetch subordinates error:', error)
  } finally {
    subordinatesLoading.value = false
  }
}

// 获取使用记录
async function fetchUsages() {
  usagesLoading.value = true
  try {
    const res = await getSettlementUsageList({
      usage_type: usageType.value,
      page: usagePage.value,
      page_size: usagePageSize.value,
    })
    usages.value = res.list || []
    usageTotal.value = res.total
  } catch (error) {
    console.error('Fetch usages error:', error)
  } finally {
    usagesLoading.value = false
  }
}

// 使用沉淀款
function handleUseSettlement() {
  if (remainingAmount.value <= 0) {
    ElMessage.warning('当前没有可用的沉淀额度')
    return
  }
  useForm.amount = Math.min(100, remainingAmount.value / 100)
  useForm.remark = ''
  useDialogVisible.value = true
}

// 提交使用
async function handleSubmitUse() {
  if (useForm.amount <= 0) {
    ElMessage.warning('请输入使用金额')
    return
  }
  if (useForm.amount * 100 > remainingAmount.value) {
    ElMessage.warning('使用金额超过可用额度')
    return
  }

  try {
    await useSettlement({
      amount: Math.round(useForm.amount * 100),
      remark: useForm.remark,
    })
    ElMessage.success('使用成功')
    useDialogVisible.value = false
    fetchSummary()
    fetchUsages()
  } catch (error) {
    console.error('Submit use error:', error)
  }
}

// 归还沉淀款
function handleReturnSettlement() {
  returnForm.amount = Math.min(100, summary.value.used_amount / 100)
  returnForm.remark = ''
  returnDialogVisible.value = true
}

// 提交归还
async function handleSubmitReturn() {
  if (returnForm.amount <= 0) {
    ElMessage.warning('请输入归还金额')
    return
  }

  try {
    await returnSettlement({
      amount: Math.round(returnForm.amount * 100),
      remark: returnForm.remark,
    })
    ElMessage.success('归还成功')
    returnDialogVisible.value = false
    fetchSummary()
    fetchUsages()
  } catch (error) {
    console.error('Submit return error:', error)
  }
}

onMounted(() => {
  fetchSummary()
  fetchSubordinates()
  fetchUsages()
})
</script>

<style lang="scss" scoped>
.settlement-wallet-view {
  padding: 0;
}

.summary-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.stat-item {
  text-align: center;
  padding: $spacing-md;

  .stat-label {
    font-size: 14px;
    color: $text-secondary;
    margin-bottom: $spacing-sm;
  }

  .stat-value {
    font-size: 24px;
    font-weight: 600;

    &.primary {
      color: $primary-color;
    }

    &.success {
      color: $success-color;
    }

    &.warning {
      color: $warning-color;
    }

    &.danger {
      color: $danger-color;
    }
  }
}

.action-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .remaining-tip {
    margin-top: $spacing-md;
  }
}

.subordinate-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .balance-amount {
    font-weight: 600;
    color: $primary-color;
  }
}

.usage-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .amount-use {
    color: $danger-color;
    font-weight: 600;
  }

  .amount-return {
    color: $success-color;
    font-weight: 600;
  }
}

.pagination {
  margin-top: $spacing-md;
  justify-content: flex-end;
}

.available-amount {
  font-size: 18px;
  font-weight: 600;
  color: $success-color;
}

.pending-amount {
  font-size: 18px;
  font-weight: 600;
  color: $warning-color;
}

.form-tip {
  margin-left: $spacing-sm;
  color: $text-secondary;
}
</style>
