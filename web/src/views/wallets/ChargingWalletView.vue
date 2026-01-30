<template>
  <div class="charging-wallet-view">
    <!-- 钱包配置状态 -->
    <el-card class="config-card" v-loading="configLoading">
      <template #header>
        <div class="card-header">
          <span>钱包配置</span>
          <el-button v-if="!config?.charging_wallet_enabled" type="primary" @click="handleEnableWallet">
            开通充值钱包
          </el-button>
        </div>
      </template>

      <el-descriptions :column="3" border v-if="config">
        <el-descriptions-item label="状态">
          <el-tag :type="config.charging_wallet_enabled ? 'success' : 'info'">
            {{ config.charging_wallet_enabled ? '已开通' : '未开通' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="充值限额">
          ¥{{ formatAmount(config.charging_wallet_limit) }}
        </el-descriptions-item>
        <el-descriptions-item label="开通时间">
          {{ config.enabled_at || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-empty v-else description="暂无配置信息" />
    </el-card>

    <!-- 钱包余额汇总 -->
    <el-row :gutter="20" class="summary-row" v-if="config?.charging_wallet_enabled">
      <el-col :xs="24" :sm="12" :lg="8">
        <div class="summary-card balance">
          <div class="summary-icon">
            <el-icon><Wallet /></el-icon>
          </div>
          <div class="summary-info">
            <div class="summary-label">当前余额</div>
            <div class="summary-value">¥{{ formatAmount(summary.balance) }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="8">
        <div class="summary-card issued">
          <div class="summary-icon">
            <el-icon><Present /></el-icon>
          </div>
          <div class="summary-info">
            <div class="summary-label">累计发放</div>
            <div class="summary-value">¥{{ formatAmount(summary.total_issued) }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 功能操作区 -->
    <el-card class="action-card" v-if="config?.charging_wallet_enabled">
      <template #header>
        <div class="card-header">
          <span>快捷操作</span>
        </div>
      </template>

      <el-space wrap>
        <el-button type="primary" :icon="Plus" @click="handleDeposit">申请充值</el-button>
        <el-button type="success" :icon="Present" @click="handleIssueReward">发放奖励</el-button>
      </el-space>
    </el-card>

    <!-- 充值记录 -->
    <el-card class="list-card" v-if="config?.charging_wallet_enabled">
      <template #header>
        <div class="card-header">
          <span>充值记录</span>
          <el-radio-group v-model="depositStatus" @change="fetchDeposits">
            <el-radio-button :value="undefined">全部</el-radio-button>
            <el-radio-button :value="0">待确认</el-radio-button>
            <el-radio-button :value="1">已确认</el-radio-button>
            <el-radio-button :value="2">已拒绝</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="deposits" v-loading="depositsLoading" border stripe>
        <el-table-column prop="deposit_no" label="充值单号" width="180" />
        <el-table-column prop="amount_yuan" label="充值金额" width="120" align="right">
          <template #default="{ row }">
            ¥{{ row.amount_yuan.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="payment_method_name" label="付款方式" width="100" />
        <el-table-column prop="status_name" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getDepositStatusType(row.status)">{{ row.status_name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column prop="confirmed_at" label="确认时间" width="180" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      </el-table>

      <el-pagination
        v-if="depositTotal > 0"
        class="pagination"
        background
        layout="total, sizes, prev, pager, next"
        :total="depositTotal"
        v-model:current-page="depositPage"
        v-model:page-size="depositPageSize"
        :page-sizes="[10, 20, 50]"
        @change="fetchDeposits"
      />
    </el-card>

    <!-- 奖励记录 -->
    <el-card class="list-card" v-if="config?.charging_wallet_enabled">
      <template #header>
        <div class="card-header">
          <span>奖励发放记录</span>
          <el-radio-group v-model="rewardDirection" @change="fetchRewards">
            <el-radio-button value="from">我发放的</el-radio-button>
            <el-radio-button value="to">我收到的</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="rewards" v-loading="rewardsLoading" border stripe>
        <el-table-column prop="reward_no" label="奖励单号" width="180" />
        <el-table-column prop="from_agent_name" label="发放人" width="120" v-if="rewardDirection === 'to'" />
        <el-table-column prop="to_agent_name" label="接收人" width="120" v-if="rewardDirection === 'from'" />
        <el-table-column prop="amount_yuan" label="奖励金额" width="120" align="right">
          <template #default="{ row }">
            ¥{{ row.amount_yuan.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="reward_type_name" label="奖励类型" width="100" />
        <el-table-column prop="status_name" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status_name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发放时间" width="180" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      </el-table>

      <el-pagination
        v-if="rewardTotal > 0"
        class="pagination"
        background
        layout="total, sizes, prev, pager, next"
        :total="rewardTotal"
        v-model:current-page="rewardPage"
        v-model:page-size="rewardPageSize"
        :page-sizes="[10, 20, 50]"
        @change="fetchRewards"
      />
    </el-card>

    <!-- 充值弹窗 -->
    <el-dialog v-model="depositDialogVisible" title="申请充值" width="500px">
      <el-form :model="depositForm" label-width="100px">
        <el-form-item label="充值金额" required>
          <el-input-number v-model="depositForm.amount" :min="1" :precision="2" style="width: 200px" />
          <span class="form-tip">元</span>
        </el-form-item>
        <el-form-item label="付款方式" required>
          <el-select v-model="depositForm.payment_method" placeholder="请选择" style="width: 100%">
            <el-option label="银行转账" :value="1" />
            <el-option label="微信转账" :value="2" />
            <el-option label="支付宝转账" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="支付凭证">
          <el-input v-model="depositForm.payment_ref" placeholder="转账流水号或截图" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="depositForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="depositDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitDeposit">确认提交</el-button>
      </template>
    </el-dialog>

    <!-- 发放奖励弹窗 -->
    <el-dialog v-model="rewardDialogVisible" title="发放奖励" width="500px">
      <el-form :model="rewardForm" label-width="100px">
        <el-form-item label="接收代理" required>
          <AgentSelect v-model="rewardForm.to_agent_id" placeholder="选择下级代理商" />
        </el-form-item>
        <el-form-item label="奖励金额" required>
          <el-input-number v-model="rewardForm.amount" :min="0.01" :precision="2" style="width: 200px" />
          <span class="form-tip">元</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="rewardForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rewardDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitReward">确认发放</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Wallet, Present, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import {
  getMyWalletConfig,
  getChargingWalletSummary,
  getChargingDepositList,
  getChargingRewardList,
  createChargingDeposit,
  issueChargingReward,
} from '@/api/chargingWallet'
import { formatAmount } from '@/utils/format'
import type {
  AgentWalletConfig,
  ChargingWalletSummary,
  ChargingWalletDeposit,
  ChargingWalletReward,
} from '@/types'

// 配置
const config = ref<AgentWalletConfig | null>(null)
const configLoading = ref(false)

// 汇总
const summary = ref<ChargingWalletSummary>({
  balance: 0,
  balance_yuan: 0,
  total_issued: 0,
  total_issued_yuan: 0,
})

// 充值记录
const deposits = ref<ChargingWalletDeposit[]>([])
const depositsLoading = ref(false)
const depositStatus = ref<number | undefined>(undefined)
const depositPage = ref(1)
const depositPageSize = ref(10)
const depositTotal = ref(0)

// 奖励记录
const rewards = ref<ChargingWalletReward[]>([])
const rewardsLoading = ref(false)
const rewardDirection = ref<'from' | 'to'>('from')
const rewardPage = ref(1)
const rewardPageSize = ref(10)
const rewardTotal = ref(0)

// 充值弹窗
const depositDialogVisible = ref(false)
const depositForm = reactive({
  amount: 100,
  payment_method: 1,
  payment_ref: '',
  remark: '',
})

// 奖励弹窗
const rewardDialogVisible = ref(false)
const rewardForm = reactive({
  to_agent_id: undefined as number | undefined,
  amount: 10,
  remark: '',
})

// 获取配置
async function fetchConfig() {
  configLoading.value = true
  try {
    config.value = await getMyWalletConfig()
  } catch (error) {
    console.error('Fetch config error:', error)
  } finally {
    configLoading.value = false
  }
}

// 获取汇总
async function fetchSummary() {
  try {
    summary.value = await getChargingWalletSummary()
  } catch (error) {
    console.error('Fetch summary error:', error)
  }
}

// 获取充值记录
async function fetchDeposits() {
  depositsLoading.value = true
  try {
    const res = await getChargingDepositList({
      status: depositStatus.value,
      page: depositPage.value,
      page_size: depositPageSize.value,
    })
    deposits.value = res.list || []
    depositTotal.value = res.total
  } catch (error) {
    console.error('Fetch deposits error:', error)
  } finally {
    depositsLoading.value = false
  }
}

// 获取奖励记录
async function fetchRewards() {
  rewardsLoading.value = true
  try {
    const res = await getChargingRewardList({
      direction: rewardDirection.value,
      page: rewardPage.value,
      page_size: rewardPageSize.value,
    })
    rewards.value = res.list || []
    rewardTotal.value = res.total
  } catch (error) {
    console.error('Fetch rewards error:', error)
  } finally {
    rewardsLoading.value = false
  }
}

// 充值状态类型
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'
function getDepositStatusType(status: number): TagType {
  const typeMap: Record<number, TagType> = {
    0: 'warning',
    1: 'success',
    2: 'danger',
  }
  return typeMap[status] || 'info'
}

// 开通钱包
function handleEnableWallet() {
  ElMessage.info('请联系管理员开通充值钱包')
}

// 申请充值
function handleDeposit() {
  depositForm.amount = 100
  depositForm.payment_method = 1
  depositForm.payment_ref = ''
  depositForm.remark = ''
  depositDialogVisible.value = true
}

// 提交充值
async function handleSubmitDeposit() {
  if (depositForm.amount <= 0) {
    ElMessage.warning('请输入充值金额')
    return
  }

  try {
    await createChargingDeposit({
      amount: Math.round(depositForm.amount * 100),
      payment_method: depositForm.payment_method,
      payment_ref: depositForm.payment_ref,
      remark: depositForm.remark,
    })
    ElMessage.success('充值申请已提交')
    depositDialogVisible.value = false
    fetchDeposits()
  } catch (error) {
    console.error('Submit deposit error:', error)
  }
}

// 发放奖励
function handleIssueReward() {
  rewardForm.to_agent_id = undefined
  rewardForm.amount = 10
  rewardForm.remark = ''
  rewardDialogVisible.value = true
}

// 提交奖励
async function handleSubmitReward() {
  if (!rewardForm.to_agent_id) {
    ElMessage.warning('请选择接收代理')
    return
  }
  if (rewardForm.amount <= 0) {
    ElMessage.warning('请输入奖励金额')
    return
  }

  try {
    await issueChargingReward({
      to_agent_id: rewardForm.to_agent_id,
      amount: Math.round(rewardForm.amount * 100),
      remark: rewardForm.remark,
    })
    ElMessage.success('奖励发放成功')
    rewardDialogVisible.value = false
    fetchRewards()
    fetchSummary()
  } catch (error) {
    console.error('Submit reward error:', error)
  }
}

onMounted(() => {
  fetchConfig()
  fetchSummary()
  fetchDeposits()
  fetchRewards()
})
</script>

<style lang="scss" scoped>
.charging-wallet-view {
  padding: 0;
}

.config-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.summary-row {
  margin-bottom: $spacing-lg;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.summary-card {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-lg;
  background: $bg-white;
  border-radius: $border-radius-md;
  box-shadow: $shadow-sm;

  .summary-icon {
    width: 50px;
    height: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    font-size: 24px;
    color: #ffffff;
  }

  &.balance .summary-icon {
    background: linear-gradient(135deg, #409eff, #66b1ff);
  }

  &.issued .summary-icon {
    background: linear-gradient(135deg, #67c23a, #85ce61);
  }

  .summary-info {
    .summary-label {
      font-size: 14px;
      color: $text-secondary;
      margin-bottom: $spacing-xs;
    }

    .summary-value {
      font-size: 24px;
      font-weight: 600;
      color: $text-primary;
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
}

.list-card {
  margin-bottom: $spacing-lg;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.pagination {
  margin-top: $spacing-md;
  justify-content: flex-end;
}

.form-tip {
  margin-left: $spacing-sm;
  color: $text-secondary;
}
</style>
