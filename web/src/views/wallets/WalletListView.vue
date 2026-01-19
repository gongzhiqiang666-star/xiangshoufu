<template>
  <div class="wallet-list-view">
    <PageHeader title="钱包管理" sub-title="钱包总览" />

    <!-- 钱包汇总 -->
    <el-row :gutter="20" class="wallet-summary">
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="wallet-card profit">
          <div class="wallet-icon">
            <el-icon><Wallet /></el-icon>
          </div>
          <div class="wallet-info">
            <div class="wallet-name">分润钱包</div>
            <div class="wallet-balance">¥{{ formatAmount(summary.profit_balance) }}</div>
            <div class="wallet-desc">交易分润收入</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="wallet-card service">
          <div class="wallet-icon">
            <el-icon><CreditCard /></el-icon>
          </div>
          <div class="wallet-info">
            <div class="wallet-name">服务费钱包</div>
            <div class="wallet-balance">¥{{ formatAmount(summary.service_balance) }}</div>
            <div class="wallet-desc">流量费+押金返现</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="wallet-card reward">
          <div class="wallet-icon">
            <el-icon><Present /></el-icon>
          </div>
          <div class="wallet-info">
            <div class="wallet-name">奖励钱包</div>
            <div class="wallet-balance">¥{{ formatAmount(summary.reward_balance) }}</div>
            <div class="wallet-desc">激活奖励收入</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="wallet-card total">
          <div class="wallet-icon">
            <el-icon><Money /></el-icon>
          </div>
          <div class="wallet-info">
            <div class="wallet-name">可用余额</div>
            <div class="wallet-balance">¥{{ formatAmount(summary.total_available) }}</div>
            <div class="wallet-desc">总余额: ¥{{ formatAmount(summary.total_balance) }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 钱包列表 -->
    <el-card class="wallet-list-card">
      <template #header>
        <div class="card-header">
          <span>钱包明细</span>
          <el-button type="primary" :icon="Download" @click="handleWithdraw">
            申请提现
          </el-button>
        </div>
      </template>

      <el-table :data="wallets" v-loading="loading" border stripe>
        <el-table-column prop="channel_name" label="通道" width="120" />
        <el-table-column prop="wallet_type" label="钱包类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getWalletTypeTag(row.wallet_type)" size="small">
              {{ getWalletTypeLabel(row.wallet_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="150" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.balance) }}
          </template>
        </el-table-column>
        <el-table-column prop="frozen" label="冻结金额" width="150" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.frozen) }}
          </template>
        </el-table-column>
        <el-table-column prop="available" label="可用金额" width="150" align="right">
          <template #default="{ row }">
            <span class="available-amount">¥{{ formatAmount(row.available) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_income" label="累计收入" width="150" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.total_income) }}
          </template>
        </el-table-column>
        <el-table-column prop="total_withdraw" label="累计提现" width="150" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.total_withdraw) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewLogs(row)">流水</el-button>
            <el-button
              type="success"
              link
              :disabled="row.available <= 0"
              @click="handleWithdrawSingle(row)"
            >
              提现
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 提现弹窗 -->
    <el-dialog v-model="withdrawDialogVisible" title="申请提现" width="500px">
      <el-form :model="withdrawForm" label-width="100px">
        <el-form-item label="选择钱包" required>
          <el-select v-model="withdrawForm.wallet_id" placeholder="请选择钱包">
            <el-option
              v-for="wallet in availableWallets"
              :key="wallet.id"
              :label="`${wallet.channel_name} - ${getWalletTypeLabel(wallet.wallet_type)} (可用: ¥${formatAmount(wallet.available)})`"
              :value="wallet.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="提现金额" required>
          <el-input-number
            v-model="withdrawForm.amount"
            :min="0"
            :max="maxWithdrawAmount"
            :precision="2"
            style="width: 200px"
          />
          <span class="form-tip">元</span>
        </el-form-item>
        <el-form-item label="可提金额">
          <span class="max-amount">¥{{ formatAmount(maxWithdrawAmount * 100) }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="withdrawDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitWithdraw">确认提现</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Wallet, CreditCard, Present, Money, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import { getWallets, getWalletSummary, applyWithdraw } from '@/api/wallet'
import { formatAmount } from '@/utils/format'
import type { Wallet as WalletType, WalletSummary, WalletType as WalletTypeEnum } from '@/types'
import { WALLET_TYPE_CONFIG } from '@/types/wallet'

const router = useRouter()

// 钱包汇总
const summary = ref<WalletSummary>({
  profit_balance: 0,
  service_balance: 0,
  reward_balance: 0,
  recharge_balance: 0,
  deposit_balance: 0,
  total_balance: 0,
  total_available: 0,
  total_frozen: 0,
})

// 钱包列表
const wallets = ref<WalletType[]>([])
const loading = ref(false)

// 提现弹窗
const withdrawDialogVisible = ref(false)
const withdrawForm = reactive({
  wallet_id: undefined as number | undefined,
  amount: 0,
})

// 可提现的钱包
const availableWallets = computed(() => wallets.value.filter((w) => w.available > 0))

// 最大提现金额
const maxWithdrawAmount = computed(() => {
  const wallet = wallets.value.find((w) => w.id === withdrawForm.wallet_id)
  return wallet ? wallet.available / 100 : 0
})

// 钱包类型配置
function getWalletTypeTag(type: WalletTypeEnum) {
  const colorMap: Record<string, string> = {
    '#409eff': 'primary',
    '#67c23a': 'success',
    '#e6a23c': 'warning',
    '#f56c6c': 'danger',
    '#909399': 'info',
  }
  const config = WALLET_TYPE_CONFIG[type]
  return colorMap[config?.color] || ''
}

function getWalletTypeLabel(type: WalletTypeEnum) {
  return WALLET_TYPE_CONFIG[type]?.label || type
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const [summaryData, walletsData] = await Promise.all([getWalletSummary(), getWallets()])
    summary.value = summaryData
    wallets.value = walletsData
  } catch (error) {
    console.error('Fetch wallets error:', error)
  } finally {
    loading.value = false
  }
}

// 查看流水
function handleViewLogs(wallet: WalletType) {
  router.push(`/wallets/${wallet.id}/logs`)
}

// 申请提现
function handleWithdraw() {
  withdrawForm.wallet_id = undefined
  withdrawForm.amount = 0
  withdrawDialogVisible.value = true
}

// 单个钱包提现
function handleWithdrawSingle(wallet: WalletType) {
  withdrawForm.wallet_id = wallet.id
  withdrawForm.amount = 0
  withdrawDialogVisible.value = true
}

// 提交提现
async function handleSubmitWithdraw() {
  if (!withdrawForm.wallet_id) {
    ElMessage.warning('请选择钱包')
    return
  }
  if (withdrawForm.amount <= 0) {
    ElMessage.warning('请输入提现金额')
    return
  }
  if (withdrawForm.amount > maxWithdrawAmount.value) {
    ElMessage.warning('提现金额超过可用余额')
    return
  }

  try {
    await applyWithdraw({
      wallet_id: withdrawForm.wallet_id,
      amount: Math.round(withdrawForm.amount * 100), // 转换为分
    })
    ElMessage.success('提现申请已提交')
    withdrawDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Apply withdraw error:', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.wallet-list-view {
  padding: 0;
}

.wallet-summary {
  margin-bottom: $spacing-lg;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.wallet-card {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-lg;
  background: $bg-white;
  border-radius: $border-radius-md;
  box-shadow: $shadow-sm;
  transition: all $transition-normal;

  &:hover {
    box-shadow: $shadow-md;
    transform: translateY(-2px);
  }

  .wallet-icon {
    width: 60px;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    font-size: 28px;
    color: #ffffff;
  }

  &.profit .wallet-icon {
    background: linear-gradient(135deg, #409eff, #66b1ff);
  }

  &.service .wallet-icon {
    background: linear-gradient(135deg, #67c23a, #85ce61);
  }

  &.reward .wallet-icon {
    background: linear-gradient(135deg, #e6a23c, #ebb563);
  }

  &.total .wallet-icon {
    background: linear-gradient(135deg, #f56c6c, #f78989);
  }

  .wallet-info {
    flex: 1;

    .wallet-name {
      font-size: 14px;
      color: $text-secondary;
      margin-bottom: $spacing-xs;
    }

    .wallet-balance {
      font-size: 24px;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }

    .wallet-desc {
      font-size: 12px;
      color: $text-placeholder;
    }
  }
}

.wallet-list-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.available-amount {
  color: $success-color;
  font-weight: 600;
}

.form-tip {
  margin-left: $spacing-sm;
  color: $text-secondary;
}

.max-amount {
  font-size: 18px;
  font-weight: 600;
  color: $primary-color;
}
</style>
