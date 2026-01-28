<template>
  <div class="high-rate-config-form">
    <el-alert
      v-if="showHint"
      type="info"
      :closable="false"
      style="margin-bottom: 16px"
    >
      高调费率配置：初始值继承自政策模版，可根据级差规则单独调整。下级费率 ≥ 上级费率。
    </el-alert>

    <div v-if="rateTypes.length === 0" class="empty-tip">
      <el-empty description="请先选择通道以获取费率类型" :image-size="60" />
    </div>

    <el-form v-else label-width="140px" :model="formData">
      <el-form-item
        v-for="rateType in rateTypes"
        :key="rateType.code"
        :label="`${rateType.name}高调费率`"
      >
        <div class="rate-input-wrapper">
          <el-input-number
            v-model="formData[rateType.code]"
            :min="getMinRate(rateType.code)"
            :max="100"
            :precision="2"
            :step="0.01"
            :disabled="disabled"
            placeholder="0.00"
            controls-position="right"
            style="width: 160px"
          />
          <span class="unit">%</span>
          <span v-if="parentRates[rateType.code]" class="hint">
            上级: {{ parentRates[rateType.code] }}%
          </span>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface RateType {
  code: string
  name: string
}

interface HighRateConfigs {
  [key: string]: { rate: string }
}

const props = withDefaults(defineProps<{
  modelValue: HighRateConfigs
  rateTypes: RateType[]
  parentRates?: Record<string, string>
  disabled?: boolean
  showHint?: boolean
}>(), {
  parentRates: () => ({}),
  disabled: false,
  showHint: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: HighRateConfigs): void
}>()

// 内部表单数据（数值形式）
const formData = ref<Record<string, number>>({})

// 初始化表单数据
const initFormData = () => {
  const data: Record<string, number> = {}
  for (const rateType of props.rateTypes) {
    const config = props.modelValue[rateType.code]
    data[rateType.code] = config ? parseFloat(config.rate) || 0 : 0
  }
  formData.value = data
}

// 监听 modelValue 变化
watch(() => props.modelValue, initFormData, { immediate: true, deep: true })

// 监听 rateTypes 变化
watch(() => props.rateTypes, initFormData, { immediate: true })

// 监听表单数据变化，向上同步
watch(formData, (newData) => {
  const configs: HighRateConfigs = {}
  for (const [code, rate] of Object.entries(newData)) {
    configs[code] = { rate: rate.toFixed(2) }
  }
  emit('update:modelValue', configs)
}, { deep: true })

// 获取最小费率（上级费率）
const getMinRate = (code: string): number => {
  const parentRate = props.parentRates[code]
  return parentRate ? parseFloat(parentRate) || 0 : 0
}
</script>

<style scoped>
.high-rate-config-form {
  padding: 8px 0;
}

.rate-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.unit {
  color: #606266;
  font-size: 14px;
}

.hint {
  color: #909399;
  font-size: 12px;
  margin-left: 8px;
}

.empty-tip {
  padding: 20px 0;
}
</style>
