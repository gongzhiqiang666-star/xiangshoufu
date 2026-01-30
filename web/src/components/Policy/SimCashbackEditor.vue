<template>
  <div class="sim-cashback-editor">
    <el-alert type="warning" :closable="false" class="mb-4">
      <template #title>
        <el-icon><InfoFilled /></el-icon>
        配置流量卡返现：商户缴纳流量费后，按次数返现给代理商
      </template>
    </el-alert>

    <el-form label-position="top" v-loading="loading">
      <el-row :gutter="16" v-if="simTiers.length > 0">
        <el-col :span="8" v-for="tier in simTiers" :key="tier.tierOrder">
          <el-form-item :label="tier.tierName">
            <el-input-number
              v-model="tier.cashback"
              :min="0"
              :max="tier.maxCashback ?? 99"
              :precision="2"
              :step="1"
              controls-position="right"
              @change="emitChange"
            />
            <div class="limit-hint" v-if="tier.maxCashback">
              最高 ¥{{ tier.maxCashback }}
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-empty v-else-if="!loading" description="该通道未配置流量费档位" />
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'
import { getChannelSimCashbackTiers } from '@/api/channel'

interface SimCashbackConfig {
  first_time_cashback: number
  second_time_cashback: number
  third_plus_cashback: number
}

interface SimTier {
  tierOrder: number
  tierName: string
  cashback: number
  maxCashback: number
  isLastTier: boolean
}

interface Limits {
  firstCashback?: number
  secondCashback?: number
  thirdPlusCashback?: number
}

const props = defineProps<{
  modelValue?: SimCashbackConfig
  channelId?: number
  limits?: Limits
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: SimCashbackConfig): void
}>()

const loading = ref(false)
const simTiers = ref<SimTier[]>([])

// 从通道配置加载流量费档位
async function loadSimTiers() {
  if (!props.channelId) {
    // 没有channelId时使用默认档位（兼容旧逻辑）
    simTiers.value = [
      { tierOrder: 1, tierName: '首次返现', cashback: 0, maxCashback: props.limits?.firstCashback ?? 99, isLastTier: false },
      { tierOrder: 2, tierName: '二次返现', cashback: 0, maxCashback: props.limits?.secondCashback ?? 99, isLastTier: false },
      { tierOrder: 3, tierName: '后续返现（第三次及以后）', cashback: 0, maxCashback: props.limits?.thirdPlusCashback ?? 99, isLastTier: true },
    ]
    syncFromModelValue()
    return
  }

  loading.value = true
  try {
    const tiers = await getChannelSimCashbackTiers(props.channelId)
    if (tiers?.length) {
      simTiers.value = tiers.map((tier: any) => ({
        tierOrder: tier.tier_order,
        tierName: tier.tier_name,
        cashback: 0,
        maxCashback: tier.max_cashback_amount / 100,
        isLastTier: tier.is_last_tier,
      }))
    } else {
      simTiers.value = []
    }
    syncFromModelValue()
  } catch (error) {
    console.error('加载流量费档位失败:', error)
    // 降级使用默认档位
    simTiers.value = [
      { tierOrder: 1, tierName: '首次返现', cashback: 0, maxCashback: props.limits?.firstCashback ?? 99, isLastTier: false },
      { tierOrder: 2, tierName: '二次返现', cashback: 0, maxCashback: props.limits?.secondCashback ?? 99, isLastTier: false },
      { tierOrder: 3, tierName: '后续返现（第三次及以后）', cashback: 0, maxCashback: props.limits?.thirdPlusCashback ?? 99, isLastTier: true },
    ]
    syncFromModelValue()
  } finally {
    loading.value = false
  }
}

function syncFromModelValue() {
  if (props.modelValue) {
    const tier1 = simTiers.value.find(t => t.tierOrder === 1)
    const tier2 = simTiers.value.find(t => t.tierOrder === 2)
    const tier3 = simTiers.value.find(t => t.tierOrder === 3)
    if (tier1) tier1.cashback = props.modelValue.first_time_cashback / 100
    if (tier2) tier2.cashback = props.modelValue.second_time_cashback / 100
    if (tier3) tier3.cashback = props.modelValue.third_plus_cashback / 100
  }
}

onMounted(() => {
  loadSimTiers()
})

watch(() => props.channelId, () => {
  loadSimTiers()
})

watch(() => props.modelValue, () => {
  syncFromModelValue()
}, { deep: true })

function emitChange() {
  const tier1 = simTiers.value.find(t => t.tierOrder === 1)
  const tier2 = simTiers.value.find(t => t.tierOrder === 2)
  const tier3 = simTiers.value.find(t => t.tierOrder === 3)

  emit('update:modelValue', {
    first_time_cashback: (tier1?.cashback ?? 0) * 100,
    second_time_cashback: (tier2?.cashback ?? 0) * 100,
    third_plus_cashback: (tier3?.cashback ?? 0) * 100,
  })
}
</script>

<style scoped>
.sim-cashback-editor {
  padding: 16px;
  background: #fffbf0;
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
