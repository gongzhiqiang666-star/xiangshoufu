import type { Router, RouteLocationNormalized } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { getToken } from '@/utils/storage'

// 白名单路由（无需登录）
const whiteList = ['/login']

/**
 * 设置路由守卫
 */
export function setupRouterGuards(router: Router) {
  // 前置守卫
  router.beforeEach(async (to, _from, next) => {
    const userStore = useUserStore()
    const appStore = useAppStore()

    // 开始页面加载
    appStore.setPageLoading(true)

    // 设置页面标题
    const title = to.meta.title as string
    document.title = title ? `${title} - ${appStore.title}` : appStore.title

    // 获取Token
    const hasToken = getToken()

    if (hasToken) {
      if (to.path === '/login') {
        // 已登录访问登录页，重定向到首页
        next({ path: '/' })
      } else {
        // 检查是否有用户信息
        if (userStore.userInfo) {
          handleRouteChange(to, appStore)
          next()
        } else {
          // 尝试获取用户信息
          try {
            await userStore.fetchUserInfo()
            handleRouteChange(to, appStore)
            next()
          } catch (error) {
            // 获取用户信息失败，清除Token并跳转登录
            await userStore.logout()
            next(`/login?redirect=${to.path}`)
          }
        }
      }
    } else {
      // 未登录
      if (whiteList.includes(to.path)) {
        next()
      } else {
        // 保存目标路由，登录后跳转
        next(`/login?redirect=${to.path}`)
      }
    }
  })

  // 后置守卫
  router.afterEach(() => {
    const appStore = useAppStore()
    // 结束页面加载
    appStore.setPageLoading(false)
  })

  // 错误处理
  router.onError((error) => {
    console.error('Router error:', error)
  })
}

/**
 * 处理路由变化
 */
function handleRouteChange(to: RouteLocationNormalized, appStore: ReturnType<typeof useAppStore>) {
  // 设置激活菜单
  const matched = to.matched
  if (matched.length > 1) {
    // 子路由，使用父路由路径
    appStore.setActiveMenu(matched[1].path)
  } else {
    appStore.setActiveMenu(to.path)
  }

  // 设置面包屑
  const breadcrumb = to.meta.breadcrumb as { title: string; path?: string }[]
  if (breadcrumb) {
    appStore.setBreadcrumbs(breadcrumb)
  }
}
