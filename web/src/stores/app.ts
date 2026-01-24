import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const sidebarCollapsed = ref(false)

  // 当前激活的菜单
  const activeMenu = ref('')

  // 面包屑导航
  const breadcrumbs = ref<{ title: string; path?: string }[]>([])

  // 页面加载状态
  const pageLoading = ref(false)

  // 系统标题
  const title = computed(() => import.meta.env.VITE_APP_TITLE || '享收付管理系统')

  // 切换侧边栏
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // 设置侧边栏状态
  function setSidebarCollapsed(collapsed: boolean) {
    sidebarCollapsed.value = collapsed
  }

  // 设置激活菜单
  function setActiveMenu(menu: string) {
    activeMenu.value = menu
  }

  // 设置面包屑
  function setBreadcrumbs(items: { title: string; path?: string }[]) {
    breadcrumbs.value = items
  }

  // 设置页面加载状态
  function setPageLoading(loading: boolean) {
    pageLoading.value = loading
  }

  return {
    // 状态
    sidebarCollapsed,
    activeMenu,
    breadcrumbs,
    pageLoading,
    // 计算属性
    title,
    // 方法
    toggleSidebar,
    setSidebarCollapsed,
    setActiveMenu,
    setBreadcrumbs,
    setPageLoading,
  }
})
