<template>
  <div class="page-header">
    <div class="header-left">
      <el-button
        v-if="showBack"
        :icon="ArrowLeft"
        text
        size="small"
        @click="handleBack"
      />
      <div class="title-section">
        <span class="page-title">{{ title }}</span>
        <span v-if="subTitle" class="page-subtitle">{{ subTitle }}</span>
      </div>
    </div>
    <div class="header-right">
      <slot name="extra"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

interface Props {
  title?: string
  subTitle?: string
  showBack?: boolean
}

withDefaults(defineProps<Props>(), {
  showBack: false,
})

const router = useRouter()

function handleBack() {
  router.back()
}
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-sm;
  padding: $spacing-xs $spacing-sm;
  background: $bg-white;
  border-radius: $border-radius-sm;
  box-shadow: $shadow-sm;
  min-height: 36px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
}

.title-section {
  display: flex;
  align-items: baseline;
  gap: $spacing-xs;
}

.page-title {
  font-size: 14px;
  font-weight: 600;
  color: $text-primary;
}

.page-subtitle {
  font-size: 12px;
  color: $text-secondary;
}

.header-right {
  display: flex;
  gap: $spacing-xs;
}
</style>
