<template>
  <el-breadcrumb separator="/" class="app-breadcrumb">
    <el-breadcrumb-item
      v-for="(item, index) in breadcrumbs"
      :key="index"
      :to="item.path"
    >
      {{ item.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const breadcrumbs = computed(() => appStore.breadcrumbs)
</script>

<style lang="scss" scoped>
.app-breadcrumb {
  margin-bottom: $spacing-md;
  padding: $spacing-md;
  background-color: $bg-white;
  border-radius: $border-radius-sm;
  box-shadow: $shadow-sm;
}
</style>
