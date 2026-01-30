<template>
  <div class="deposit-cashback-editor">
    <el-alert type="info" :closable="false" class="mb-4">
      <template #title>
        <el-icon><InfoFilled /></el-icon>
        配置押金返现：商户缴纳押金后，按此配置返现给代理商
      </template>
    </el-alert>

    <el-form label-position="top" v-loading="loading">
      <el-row :gutter="16" v-if="depositItems.length > 0">
        <el-col :span="6" v-for="item in depositItems" :key="item.deposit">
          <el-form-item :label="`押金 ¥${item.deposit} 返现`">
            <el-input-number
              v-model="item.cashback"
              :min="0"
              :max="item.maxCashback ?? 999"
              :precision="2"
              :step="1"
              controls-position="right"
              @change="emitChange"
            />
            <div class="limit-hint" v-if="item.maxCashback">
              最高 ¥{{ item.maxCashback }}
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-empty v-else-if="!loading" description="该通道未配置押金档位" />
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'
import { getChannelDepositTiers } from '@/api/channel'

interface DepositCashbackItem {
  deposit_amount: number
  cashback_amount: number
}

interface DepositItem {
  deposit: number
  cashback: number
  maxCashback: number
}

const props = defineProps<{
  modelValue: DepositCashbackItem[]
  channelId?: number
  limits?: Record<number, number>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: DepositCashbackItem[]): void
}>()

const loading = ref(false)
const depositItems = ref<DepositItem[]>([])

// 从通道配置加载押金档位
async function loadDepositTiers() {
  if (!props.channelId) {
    // 没有channelId时使用默认档位（兼容旧逻辑）
    depositItems.value = [
      { deposit: 0, cashback: 0, maxCashback: props.limits?.[0] ?? 999 },
      { deposit: 99, cashback: 0, maxCashback: props.limits?.[99] ?? 999 },
      { deposit: 199, cashback: 0, maxCashback: props.limits?.[199] ?? 999 },
      { deposit: 299, cashback: 0, maxCashback: props.limits?.[299] ?? 999 },
    ]
    syncFromModelValue()
    return
  }

  loading.value = true
  try {
    const tiers = await getChannelDepositTiers(props.channelId)
    if (tiers?.length) {
      depositItems.value = tiers.map((tier: any) => ({
        deposit: tier.deposit_amount / 100,
        cashback: 0,
        maxCashback: tier.max_cashback_amount / 100,
      }))
    } else {
      depositItems.value = []
    }
    syncFromModelValue()
  } catch (error) {
    console.error('加载押金档位失败:', error)
    // 降级使用默认档位
    depositItems.value = [
      { deposit: 0, cashback: 0, maxCashback: props.limits?.[0] ?? 999 },
      { deposit: 99, cashback: 0, maxCashback: props.limits?.[99] ?? 999 },
      { deposit: 199, cashback: 0, maxCashback: props.limits?.[199] ?? 999 },
      { deposit: 299, cashback: 0, maxCashback: props.limits?.[299] ?? 999 },
    ]
    syncFromModelValue()
  } finally {
    loading.value = false
  }
}

function syncFromModelValue() {
  if (props.modelValue?.length) {
    for (const item of props.modelValue) {
      const deposit = item.deposit_amount / 100
      const target = depositItems.value.find(d => d.deposit === deposit)
      if (target) {
        target.cashback = item.cashback_amount / 100
      }
    }
  }
}

onMounted(() => {
  loadDepositTiers()
})

watch(() => props.channelId, () => {
  loadDepositTiers()
})

watch(() => props.modelValue, () => {
  syncFromModelValue()
}, { deep: true })

function emitChange() {
  const result: DepositCashbackItem[] = depositItems.value.map(item => ({
    deposit_amount: item.deposit * 100,
    cashback_amount: item.cashback * 100,
  }))
  emit('update:modelValue', result)
}
</script>

<style scoped>
.deposit-cashback-editor {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}
.limit-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.mb-4 {
  margin-bottom: 16px;
}
</style>
