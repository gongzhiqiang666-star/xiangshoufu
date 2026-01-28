<template>
  <div class="d0-extra-config-form">
    <el-alert
      v-if="showHint"
      type="info"
      :closable="false"
      style="margin-bottom: 16px"
    >
      P+0加价配置：初始值继承自政策模版，可根据级差规则单独调整。下级加价 ≤ 上级给自己的配置金额。
    </el-alert>

    <div v-if="rateTypes.length === 0" class="empty-tip">
      <el-empty description="请先选择通道以获取费率类型" :image-size="60" />
    </div>

    <el-form v-else label-width="140px" :model="formData">
      <el-form-item
        v-for="rateType in rateTypes"
        :key="rateType.code"
        :label="`${rateType.name}P+0加价`"
      >
        <div class="extra-input-wrapper">
          <el-input-number
            v-model="formData[rateType.code]"
            :min="0"
            :max="getMaxExtra(rateType.code)"
            :step="10"
            :disabled="disabled"
            placeholder="0"
            controls-position="right"
            style="width: 160px"
          />
          <span class="unit">分</span>
          <span class="yuan-hint">({{ formatYuan(formData[rateType.code]) }}元)</span>
          <span v-if="parentExtras[rateType.code] !== undefined" class="hint">
            上级: {{ parentExtras[rateType.code] }}分
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

interface D0ExtraConfigs {
  [key: string]: { extra_fee: number }
}

const props = withDefaults(defineProps<{
  modelValue: D0ExtraConfigs
  rateTypes: RateType[]
  parentExtras?: Record<string, number>
  disabled?: boolean
  showHint?: boolean
}>(), {
  parentExtras: () => ({}),
  disabled: false,
  showHint: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: D0ExtraConfigs): void
}>()

// 内部表单数据（数值形式）
const formData = ref<Record<string, number>>({})

// 初始化表单数据
const initFormData = () => {
  const data: Record<string, number> = {}
  for (const rateType of props.rateTypes) {
    const config = props.modelValue[rateType.code]
    data[rateType.code] = config ? config.extra_fee || 0 : 0
  }
  formData.value = data
}

// 监听 modelValue 变化
watch(() => props.modelValue, initFormData, { immediate: true, deep: true })

// 监听 rateTypes 变化
watch(() => props.rateTypes, initFormData, { immediate: true })

// 监听表单数据变化，向上同步
watch(formData, (newData) => {
  const configs: D0ExtraConfigs = {}
  for (const [code, extraFee] of Object.entries(newData)) {
    configs[code] = { extra_fee: extraFee }
  }
  emit('update:modelValue', configs)
}, { deep: true })

// 获取最大加价（上级配置金额）
const getMaxExtra = (code: string): number => {
  const parentExtra = props.parentExtras[code]
  return parentExtra !== undefined ? parentExtra : 99999
}

// 格式化为元
const formatYuan = (fen: number): string => {
  if (!fen) return '0.00'
  return (fen / 100).toFixed(2)
}
</script>

<style scoped>
.d0-extra-config-form {
  padding: 8px 0;
}

.extra-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.unit {
  color: #606266;
  font-size: 14px;
}

.yuan-hint {
  color: #67c23a;
  font-size: 12px;
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
