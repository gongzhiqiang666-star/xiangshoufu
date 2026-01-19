<template>
  <el-header class="app-header">
    <!-- 左侧：折叠按钮 -->
    <div class="header-left">
      <el-icon class="collapse-btn" @click="toggleSidebar">
        <Fold v-if="!collapsed" />
        <Expand v-else />
      </el-icon>
    </div>

    <!-- 右侧：用户信息等 -->
    <div class="header-right">
      <!-- 消息通知 -->
      <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="message-badge">
        <el-icon class="header-icon" @click="showMessages">
          <Bell />
        </el-icon>
      </el-badge>

      <!-- 全屏切换 -->
      <el-icon class="header-icon" @click="toggleFullscreen">
        <FullScreen />
      </el-icon>

      <!-- 用户下拉菜单 -->
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32" class="user-avatar">
            {{ userInitial }}
          </el-avatar>
          <span class="user-name">{{ realName }}</span>
          <el-icon class="arrow-icon">
            <ArrowDown />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="password">
              <el-icon><Lock /></el-icon>
              修改密码
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </el-header>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { getUnreadCount } from '@/api/message'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

// 侧边栏折叠状态
const collapsed = computed(() => appStore.sidebarCollapsed)

// 用户信息
const realName = computed(() => userStore.realName || userStore.username)
const userInitial = computed(() => realName.value.charAt(0).toUpperCase())

// 未读消息数
const unreadCount = ref(0)

// 获取未读消息数
async function fetchUnreadCount() {
  try {
    const data = await getUnreadCount()
    unreadCount.value = data.total
  } catch (error) {
    console.error('Failed to fetch unread count:', error)
  }
}

// 切换侧边栏
function toggleSidebar() {
  appStore.toggleSidebar()
}

// 显示消息
function showMessages() {
  router.push('/messages')
}

// 切换全屏
function toggleFullscreen() {
  if (document.fullscreenElement) {
    document.exitFullscreen()
  } else {
    document.documentElement.requestFullscreen()
  }
}

// 处理下拉菜单命令
async function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'password':
      router.push('/change-password')
      break
    case 'logout':
      await handleLogout()
      break
  }
}

// 退出登录
async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await userStore.logout()
    router.push('/login')
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  fetchUnreadCount()
})
</script>

<style lang="scss" scoped>
.app-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 $spacing-md;
  background-color: $header-bg;
  box-shadow: $shadow-sm;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: $text-regular;
  transition: color $transition-fast;

  &:hover {
    color: $primary-color;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: $spacing-lg;
}

.header-icon {
  font-size: 18px;
  cursor: pointer;
  color: $text-regular;
  transition: color $transition-fast;

  &:hover {
    color: $primary-color;
  }
}

.message-badge {
  :deep(.el-badge__content) {
    top: 4px;
    right: 4px;
  }
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: $spacing-xs $spacing-sm;
  border-radius: $border-radius-sm;
  transition: background-color $transition-fast;

  &:hover {
    background-color: $bg-color;
  }
}

.user-avatar {
  background-color: $primary-color;
  color: #ffffff;
  font-weight: 600;
}

.user-name {
  margin-left: $spacing-sm;
  font-size: 14px;
  color: $text-primary;
}

.arrow-icon {
  margin-left: $spacing-xs;
  font-size: 12px;
  color: $text-secondary;
}
</style>
