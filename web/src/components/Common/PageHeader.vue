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
      <!-- 面包屑导航 -->
      <el-breadcrumb separator="/">
        <el-breadcrumb-item
          v-for="(item, index) in breadcrumbs"
          :key="index"
          :to="item.path"
        >
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="header-right">
      <slot name="extra"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'

interface Props {
  title?: string
  subTitle?: string
  showBack?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showBack: false,
})

const router = useRouter()
const appStore = useAppStore()

const breadcrumbs = computed(() => appStore.breadcrumbs)

function handleBack() {
  router.back()
}
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-md;
  padding: $spacing-sm $spacing-md;
  background: $bg-white;
  border-radius: $border-radius-sm;
  box-shadow: $shadow-sm;
  min-height: 40px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.header-right {
  display: flex;
  gap: $spacing-xs;
}

:deep(.el-breadcrumb) {
  font-size: 14px;

  .el-breadcrumb__item:last-child .el-breadcrumb__inner {
    font-weight: 600;
    color: $text-primary;
  }
}
</style>
