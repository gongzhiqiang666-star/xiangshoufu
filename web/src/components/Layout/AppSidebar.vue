<template>
  <el-aside :width="sidebarWidth" class="app-sidebar">
    <!-- Logo区域 -->
    <div class="logo-container">
      <img src="/vite.svg" alt="Logo" class="logo" />
      <transition name="fade">
        <span v-show="!collapsed" class="title">管理系统</span>
      </transition>
    </div>

    <!-- 菜单 -->
    <el-scrollbar class="menu-scrollbar">
      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :unique-opened="true"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        router
      >
        <template v-for="item in menuList" :key="item.path">
          <!-- 有子菜单 -->
          <el-sub-menu v-if="item.children?.length" :index="item.path">
            <template #title>
              <el-icon>
                <component :is="item.icon" />
              </el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.path"
              :index="child.path"
            >
              <span>{{ child.title }}</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- 无子菜单 -->
          <el-menu-item v-else :index="item.path">
            <el-icon>
              <component :is="item.icon" />
            </el-icon>
            <template #title>
              <span>{{ item.title }}</span>
            </template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-scrollbar>
  </el-aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { MENU_LIST } from '@/constants/menu'

const appStore = useAppStore()

// 菜单列表
const menuList = MENU_LIST

// 侧边栏是否折叠
const collapsed = computed(() => appStore.sidebarCollapsed)

// 侧边栏宽度
const sidebarWidth = computed(() => (collapsed.value ? '64px' : '220px'))

// 当前激活菜单
const activeMenu = computed(() => appStore.activeMenu)
</script>

<style lang="scss" scoped>
.app-sidebar {
  height: 100vh;
  background-color: $sidebar-bg;
  transition: width $transition-normal;
  overflow: hidden;

  :deep(.el-menu) {
    border-right: none;
  }

  :deep(.el-menu--collapse) {
    width: 64px;
  }
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 $spacing-md;
  background-color: darken($sidebar-bg, 5%);
  overflow: hidden;

  .logo {
    width: 32px;
    height: 32px;
    flex-shrink: 0;
  }

  .title {
    margin-left: $spacing-sm;
    font-size: 16px;
    font-weight: 600;
    color: #ffffff;
    white-space: nowrap;
  }
}

.menu-scrollbar {
  height: calc(100vh - 60px);
}

// 菜单项样式
:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  &:hover {
    background-color: lighten($sidebar-bg, 5%) !important;
  }
}

:deep(.el-menu-item.is-active) {
  background-color: $primary-color !important;
  color: #ffffff !important;

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background-color: #ffffff;
  }
}

// fade动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity $transition-fast;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
