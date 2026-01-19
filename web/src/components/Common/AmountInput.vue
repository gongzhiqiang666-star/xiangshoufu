<template>
  <el-input
    v-model="displayValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :size="size"
    class="amount-input"
    @input="handleInput"
    @blur="handleBlur"
  >
    <template #prefix>
      <span class="prefix">¥</span>
    </template>
    <template #suffix>
      <el-select
        v-if="showUnitSwitch"
        v-model="unit"
        class="unit-select"
        @change="handleUnitChange"
      >
        <el-option label="元" value="yuan" />
        <el-option label="分" value="fen" />
      </el-select>
      <span v-else class="suffix">{{ unitLabel }}</span>
    </template>
  </el-input>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface Props {
  modelValue: number | undefined
  placeholder?: string
  disabled?: boolean
  size?: 'small' | 'default' | 'large'
  valueUnit?: 'yuan' | 'fen' // 值的单位（默认分）
  showUnitSwitch?: boolean // 是否显示单位切换
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请输入金额',
  disabled: false,
  size: 'default',
  valueUnit: 'fen',
  showUnitSwitch: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | undefined): void
}>()

const unit = ref<'yuan' | 'fen'>(props.valueUnit)
const displayValue = ref('')

const unitLabel = computed(() => (unit.value === 'yuan' ? '元' : '分'))

// 格式化显示值
function formatDisplayValue(value: number | undefined) {
  if (value === undefined || value === null) {
    displayValue.value = ''
    return
  }

  if (unit.value === 'yuan') {
    // 值是分，显示元
    if (props.valueUnit === 'fen') {
      displayValue.value = (value / 100).toFixed(2)
    } else {
      displayValue.value = value.toFixed(2)
    }
  } else {
    // 显示分
    if (props.valueUnit === 'fen') {
      displayValue.value = value.toString()
    } else {
      displayValue.value = (value * 100).toString()
    }
  }
}

// 处理输入
function handleInput(value: string) {
  displayValue.value = value.replace(/[^\d.]/g, '')
}

// 处理失焦
function handleBlur() {
  const num = parseFloat(displayValue.value)
  if (isNaN(num)) {
    emit('update:modelValue', undefined)
    displayValue.value = ''
    return
  }

  let result: number
  if (unit.value === 'yuan') {
    // 输入的是元
    if (props.valueUnit === 'fen') {
      result = Math.round(num * 100)
    } else {
      result = num
    }
  } else {
    // 输入的是分
    if (props.valueUnit === 'fen') {
      result = Math.round(num)
    } else {
      result = num / 100
    }
  }

  emit('update:modelValue', result)
}

// 处理单位切换
function handleUnitChange() {
  formatDisplayValue(props.modelValue)
}

// 监听值变化
watch(
  () => props.modelValue,
  (val) => {
    formatDisplayValue(val)
  },
  { immediate: true }
)
</script>

<style lang="scss" scoped>
.amount-input {
  .prefix {
    color: $text-secondary;
  }

  .suffix {
    color: $text-secondary;
    font-size: 12px;
  }

  .unit-select {
    width: 60px;

    :deep(.el-input__wrapper) {
      box-shadow: none !important;
    }
  }
}
</style>
