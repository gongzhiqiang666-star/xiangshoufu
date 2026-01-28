<template>
  <el-card class="search-form-card">
    <el-form :model="modelValue" :inline="true" class="search-form" @submit.prevent="handleSearch">
      <slot></slot>

      <el-form-item class="form-buttons">
        <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
        <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        <slot name="extra"></slot>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { Search, Refresh } from '@element-plus/icons-vue'

interface Props {
  modelValue: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Record<string, any>): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

function handleSearch() {
  emit('search')
}

function handleReset() {
  // 重置表单值
  const resetValue: Record<string, any> = {}
  Object.keys(props.modelValue).forEach((key) => {
    resetValue[key] = undefined
  })
  emit('update:modelValue', resetValue)
  emit('reset')
}
</script>

<style lang="scss" scoped>
.search-form-card {
  margin-bottom: $spacing-sm;

  :deep(.el-card__body) {
    padding: $spacing-sm $spacing-md;
  }
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: $spacing-xs;

  :deep(.el-form-item) {
    margin-bottom: 0;
    margin-right: $spacing-sm;
  }

  :deep(.el-form-item__label) {
    padding-right: 6px;
  }

  .form-buttons {
    margin-left: auto;
  }
}
</style>
