<template>
  <el-card class="search-form-card">
    <el-form :model="modelValue" :inline="true" class="search-form" @submit.prevent="handleSearch">
      <!-- 主要条件区域 -->
      <div class="form-fields">
        <!-- 常用条件（始终显示） -->
        <div class="form-row primary-fields">
          <slot></slot>
        </div>

        <!-- 更多条件（可折叠） -->
        <el-collapse-transition>
          <div v-show="expanded" class="form-row collapse-fields">
            <slot name="collapse"></slot>
          </div>
        </el-collapse-transition>
      </div>

      <!-- 按钮区域（固定右侧） -->
      <div class="form-actions">
        <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
        <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        <el-button
          v-if="hasCollapseSlot"
          :icon="expanded ? ArrowUp : ArrowDown"
          link
          type="primary"
          @click="toggleExpand"
        >
          {{ expanded ? '收起' : '展开' }}
        </el-button>
        <slot name="extra"></slot>
      </div>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, useSlots, computed } from 'vue'
import { Search, Refresh, ArrowUp, ArrowDown } from '@element-plus/icons-vue'

interface Props {
  modelValue: Record<string, any>
  /** 默认是否展开 */
  defaultExpanded?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  defaultExpanded: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: Record<string, any>): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

const slots = useSlots()

// 是否有折叠插槽内容
const hasCollapseSlot = computed(() => !!slots.collapse)

// 展开状态
const expanded = ref(props.defaultExpanded)

function toggleExpand() {
  expanded.value = !expanded.value
}

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

// 暴露方法给父组件
defineExpose({
  expand: () => { expanded.value = true },
  collapse: () => { expanded.value = false },
  toggle: toggleExpand,
})
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
  align-items: flex-start;
  gap: $spacing-sm;

  :deep(.el-form-item) {
    margin-bottom: 0;
    margin-right: $spacing-sm;
  }

  :deep(.el-form-item__label) {
    padding-right: 6px;
  }
}

.form-fields {
  flex: 1;
  min-width: 0;
}

.form-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: $spacing-xs;

  &.collapse-fields {
    margin-top: $spacing-sm;
    padding-top: $spacing-sm;
    border-top: 1px dashed $border-color;
  }
}

.form-actions {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  flex-shrink: 0;
  padding-top: 2px; // 与表单项对齐

  :deep(.el-button + .el-button) {
    margin-left: 0;
  }
}
</style>
