<template>
  <el-container class="app-layout">
    <!-- 侧边栏 -->
    <AppSidebar />

    <el-container class="main-container">
      <!-- 顶部栏 -->
      <AppHeader />

      <!-- 主内容区 -->
      <el-main class="main-content">
        <!-- 页面内容 -->
        <div class="page-content">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
</script>

<style lang="scss" scoped>
.app-layout {
  height: 100vh;
  width: 100%;
}

.main-container {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: margin-left $transition-normal;
}

.main-content {
  flex: 1;
  padding: $spacing-sm $spacing-md;
  background-color: $bg-color;
  overflow-y: auto;
}

.page-content {
  min-height: calc(100vh - 100px);
}

// 页面切换动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
