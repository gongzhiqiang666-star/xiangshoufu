<template>
  <div class="page-header">
    <div class="header-left">
      <el-button
        v-if="showBack"
        :icon="ArrowLeft"
        text
        @click="handleBack"
      />
      <h2 class="title">{{ title }}</h2>
      <span v-if="subTitle" class="sub-title">{{ subTitle }}</span>
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
  title: string
  subTitle?: string
  showBack?: boolean
}

const props = withDefaults(defineProps<Props>(), {
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
  margin-bottom: $spacing-lg;
  padding: $spacing-md;
  background: $bg-white;
  border-radius: $border-radius-md;
  box-shadow: $shadow-sm;
}

.header-left {
  display: flex;
  align-items: center;
  gap: $spacing-sm;

  .title {
    font-size: 18px;
    font-weight: 600;
    color: $text-primary;
    margin: 0;
  }

  .sub-title {
    font-size: 14px;
    color: $text-secondary;
  }
}

.header-right {
  display: flex;
  gap: $spacing-sm;
}
</style>
