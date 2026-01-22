import type { RouteRecordRaw } from 'vue-router'

// 公开路由（无需登录）
export const publicRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: {
      title: '登录',
      hidden: true,
    },
  },
]

// 需要认证的路由
export const privateRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/components/Layout/AppLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: {
          title: '仪表盘',
          icon: 'Odometer',
          breadcrumb: [{ title: '首页' }, { title: '仪表盘' }],
        },
      },
      // 代理管理
      {
        path: 'agents',
        name: 'Agents',
        redirect: '/agents/list',
        meta: {
          title: '代理管理',
          icon: 'User',
        },
        children: [
          {
            path: 'list',
            name: 'AgentList',
            component: () => import('@/views/agents/AgentListView.vue'),
            meta: {
              title: '代理列表',
              breadcrumb: [{ title: '首页' }, { title: '代理管理' }, { title: '代理列表' }],
            },
          },
          {
            path: 'create',
            name: 'AgentCreate',
            component: () => import('@/views/agents/AgentFormView.vue'),
            meta: {
              title: '新增代理',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '代理管理' }, { title: '新增代理' }],
            },
          },
          {
            path: ':id',
            name: 'AgentDetail',
            component: () => import('@/views/agents/AgentDetailView.vue'),
            meta: {
              title: '代理详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '代理管理' }, { title: '代理详情' }],
            },
          },
          {
            path: ':id/edit',
            name: 'AgentEdit',
            component: () => import('@/views/agents/AgentFormView.vue'),
            meta: {
              title: '编辑代理',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '代理管理' }, { title: '编辑代理' }],
            },
          },
        ],
      },
      // 商户管理
      {
        path: 'merchants',
        name: 'Merchants',
        redirect: '/merchants/list',
        meta: {
          title: '商户管理',
          icon: 'Shop',
        },
        children: [
          {
            path: 'list',
            name: 'MerchantList',
            component: () => import('@/views/merchants/MerchantListView.vue'),
            meta: {
              title: '商户列表',
              breadcrumb: [{ title: '首页' }, { title: '商户管理' }, { title: '商户列表' }],
            },
          },
          {
            path: ':id',
            name: 'MerchantDetail',
            component: () => import('@/views/merchants/MerchantDetailView.vue'),
            meta: {
              title: '商户详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '商户管理' }, { title: '商户详情' }],
            },
          },
        ],
      },
      // 终端管理
      {
        path: 'terminals',
        name: 'Terminals',
        redirect: '/terminals/list',
        meta: {
          title: '终端管理',
          icon: 'Monitor',
        },
        children: [
          {
            path: 'list',
            name: 'TerminalList',
            component: () => import('@/views/terminals/TerminalListView.vue'),
            meta: {
              title: '终端列表',
              breadcrumb: [{ title: '首页' }, { title: '终端管理' }, { title: '终端列表' }],
            },
          },
          {
            path: ':id',
            name: 'TerminalDetail',
            component: () => import('@/views/terminals/TerminalDetailView.vue'),
            meta: {
              title: '终端详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '终端管理' }, { title: '终端详情' }],
            },
          },
        ],
      },
      // 交易管理
      {
        path: 'transactions',
        name: 'Transactions',
        redirect: '/transactions/list',
        meta: {
          title: '交易管理',
          icon: 'Money',
        },
        children: [
          {
            path: 'list',
            name: 'TransactionList',
            component: () => import('@/views/transactions/TransactionListView.vue'),
            meta: {
              title: '交易记录',
              breadcrumb: [{ title: '首页' }, { title: '交易管理' }, { title: '交易记录' }],
            },
          },
          {
            path: ':id',
            name: 'TransactionDetail',
            component: () => import('@/views/transactions/TransactionDetailView.vue'),
            meta: {
              title: '交易详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '交易管理' }, { title: '交易详情' }],
            },
          },
        ],
      },
      // 分润管理
      {
        path: 'profits',
        name: 'Profits',
        redirect: '/profits/list',
        meta: {
          title: '分润管理',
          icon: 'TrendCharts',
        },
        children: [
          {
            path: 'list',
            name: 'ProfitList',
            component: () => import('@/views/profits/ProfitListView.vue'),
            meta: {
              title: '分润明细',
              breadcrumb: [{ title: '首页' }, { title: '分润管理' }, { title: '分润明细' }],
            },
          },
          {
            path: ':id',
            name: 'ProfitDetail',
            component: () => import('@/views/profits/ProfitDetailView.vue'),
            meta: {
              title: '分润详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '分润管理' }, { title: '分润详情' }],
            },
          },
        ],
      },
      // 代扣管理
      {
        path: 'deductions',
        name: 'Deductions',
        redirect: '/deductions/list',
        meta: {
          title: '代扣管理',
          icon: 'CreditCard',
        },
        children: [
          {
            path: 'list',
            name: 'DeductionList',
            component: () => import('@/views/deductions/DeductionListView.vue'),
            meta: {
              title: '代扣列表',
              breadcrumb: [{ title: '首页' }, { title: '代扣管理' }, { title: '代扣列表' }],
            },
          },
          {
            path: 'create',
            name: 'DeductionCreate',
            component: () => import('@/views/deductions/DeductionCreateView.vue'),
            meta: {
              title: '发起代扣',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '代扣管理' }, { title: '发起代扣' }],
            },
          },
          {
            path: ':id',
            name: 'DeductionDetail',
            component: () => import('@/views/deductions/DeductionDetailView.vue'),
            meta: {
              title: '代扣详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '代扣管理' }, { title: '代扣详情' }],
            },
          },
        ],
      },
      // 货款代扣
      {
        path: 'goods-deductions',
        name: 'GoodsDeductions',
        redirect: '/goods-deductions/list',
        meta: {
          title: '货款代扣',
          icon: 'ShoppingCart',
        },
        children: [
          {
            path: 'list',
            name: 'GoodsDeductionList',
            component: () => import('@/views/goods-deductions/GoodsDeductionListView.vue'),
            meta: {
              title: '货款代扣列表',
              breadcrumb: [{ title: '首页' }, { title: '货款代扣' }, { title: '货款代扣列表' }],
            },
          },
          {
            path: ':id',
            name: 'GoodsDeductionDetail',
            component: () => import('@/views/goods-deductions/GoodsDeductionDetailView.vue'),
            meta: {
              title: '货款代扣详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '货款代扣' }, { title: '货款代扣详情' }],
            },
          },
        ],
      },
      // 钱包管理
      {
        path: 'wallets',
        name: 'Wallets',
        redirect: '/wallets/list',
        meta: {
          title: '钱包管理',
          icon: 'Wallet',
        },
        children: [
          {
            path: 'list',
            name: 'WalletList',
            component: () => import('@/views/wallets/WalletListView.vue'),
            meta: {
              title: '钱包总览',
              breadcrumb: [{ title: '首页' }, { title: '钱包管理' }, { title: '钱包总览' }],
            },
          },
          {
            path: 'charging',
            name: 'ChargingWallet',
            component: () => import('@/views/wallets/ChargingWalletView.vue'),
            meta: {
              title: '充值钱包',
              breadcrumb: [{ title: '首页' }, { title: '钱包管理' }, { title: '充值钱包' }],
            },
          },
          {
            path: 'settlement',
            name: 'SettlementWallet',
            component: () => import('@/views/wallets/SettlementWalletView.vue'),
            meta: {
              title: '沉淀钱包',
              breadcrumb: [{ title: '首页' }, { title: '钱包管理' }, { title: '沉淀钱包' }],
            },
          },
          {
            path: 'tax-channels',
            name: 'TaxChannels',
            component: () => import('@/views/wallets/TaxChannelView.vue'),
            meta: {
              title: '税筹通道',
              breadcrumb: [{ title: '首页' }, { title: '钱包管理' }, { title: '税筹通道' }],
            },
          },
          {
            path: ':id/logs',
            name: 'WalletLogs',
            component: () => import('@/views/wallets/WalletLogsView.vue'),
            meta: {
              title: '钱包流水',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '钱包管理' }, { title: '钱包流水' }],
            },
          },
        ],
      },
      // 政策管理
      {
        path: 'policies',
        name: 'Policies',
        redirect: '/policies/list',
        meta: {
          title: '政策管理',
          icon: 'Document',
        },
        children: [
          {
            path: 'list',
            name: 'PolicyList',
            component: () => import('@/views/policies/PolicyListView.vue'),
            meta: {
              title: '政策模板',
              breadcrumb: [{ title: '首页' }, { title: '政策管理' }, { title: '政策模板' }],
            },
          },
          {
            path: 'create',
            name: 'PolicyCreate',
            component: () => import('@/views/policies/PolicyFormView.vue'),
            meta: {
              title: '新建模板',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '政策管理' }, { title: '新建模板' }],
            },
          },
          {
            path: ':id',
            name: 'PolicyDetail',
            component: () => import('@/views/policies/PolicyDetailView.vue'),
            meta: {
              title: '模板详情',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '政策管理' }, { title: '模板详情' }],
            },
          },
          {
            path: ':id/edit',
            name: 'PolicyEdit',
            component: () => import('@/views/policies/PolicyFormView.vue'),
            meta: {
              title: '编辑模板',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '政策管理' }, { title: '编辑模板' }],
            },
          },
        ],
      },
      // 营销管理
      {
        path: 'marketing',
        name: 'Marketing',
        redirect: '/marketing/banners',
        meta: {
          title: '营销管理',
          icon: 'Picture',
        },
        children: [
          {
            path: 'banners',
            name: 'BannerList',
            component: () => import('@/views/marketing/BannerListView.vue'),
            meta: {
              title: '滚动图管理',
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '滚动图管理' }],
            },
          },
          {
            path: 'banners/create',
            name: 'BannerCreate',
            component: () => import('@/views/marketing/BannerFormView.vue'),
            meta: {
              title: '新增滚动图',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '新增滚动图' }],
            },
          },
          {
            path: 'banners/:id/edit',
            name: 'BannerEdit',
            component: () => import('@/views/marketing/BannerFormView.vue'),
            meta: {
              title: '编辑滚动图',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '编辑滚动图' }],
            },
          },
          {
            path: 'posters',
            name: 'PosterList',
            component: () => import('@/views/marketing/PosterListView.vue'),
            meta: {
              title: '海报管理',
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '海报管理' }],
            },
          },
          {
            path: 'posters/create',
            name: 'PosterCreate',
            component: () => import('@/views/marketing/PosterFormView.vue'),
            meta: {
              title: '新增海报',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '新增海报' }],
            },
          },
          {
            path: 'posters/:id/edit',
            name: 'PosterEdit',
            component: () => import('@/views/marketing/PosterFormView.vue'),
            meta: {
              title: '编辑海报',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '编辑海报' }],
            },
          },
          {
            path: 'poster-categories',
            name: 'PosterCategories',
            component: () => import('@/views/marketing/PosterCategoryView.vue'),
            meta: {
              title: '海报分类',
              breadcrumb: [{ title: '首页' }, { title: '营销管理' }, { title: '海报分类' }],
            },
          },
        ],
      },
      // 系统管理
      {
        path: 'system',
        name: 'System',
        redirect: '/system/users',
        meta: {
          title: '系统管理',
          icon: 'Setting',
        },
        children: [
          {
            path: 'users',
            name: 'UserList',
            component: () => import('@/views/system/UserListView.vue'),
            meta: {
              title: '用户管理',
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '用户管理' }],
            },
          },
          {
            path: 'logs',
            name: 'LogList',
            component: () => import('@/views/system/LogListView.vue'),
            meta: {
              title: '操作日志',
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '操作日志' }],
            },
          },
          {
            path: 'messages',
            name: 'MessageManagement',
            component: () => import('@/views/messages/MessageListView.vue'),
            meta: {
              title: '消息管理',
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '消息管理' }],
            },
          },
          {
            path: 'messages/send',
            name: 'SendMessage',
            component: () => import('@/views/messages/SendMessageView.vue'),
            meta: {
              title: '发送消息',
              hidden: true,
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '消息管理' }, { title: '发送消息' }],
            },
          },
          {
            path: 'jobs',
            name: 'JobList',
            component: () => import('@/views/system/JobListView.vue'),
            meta: {
              title: '定时任务',
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '定时任务' }],
            },
          },
          {
            path: 'job-logs',
            name: 'JobLogList',
            component: () => import('@/views/system/JobLogListView.vue'),
            meta: {
              title: '执行日志',
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '执行日志' }],
            },
          },
          {
            path: 'alert-configs',
            name: 'AlertConfigList',
            component: () => import('@/views/system/AlertConfigView.vue'),
            meta: {
              title: '告警配置',
              breadcrumb: [{ title: '首页' }, { title: '系统管理' }, { title: '告警配置' }],
            },
          },
        ],
      },
      // 个人中心
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/ProfileView.vue'),
        meta: {
          title: '个人中心',
          hidden: true,
          breadcrumb: [{ title: '首页' }, { title: '个人中心' }],
        },
      },
      // 修改密码
      {
        path: 'change-password',
        name: 'ChangePassword',
        component: () => import('@/views/profile/ChangePasswordView.vue'),
        meta: {
          title: '修改密码',
          hidden: true,
          breadcrumb: [{ title: '首页' }, { title: '修改密码' }],
        },
      },
      // 消息中心
      {
        path: 'messages',
        name: 'Messages',
        component: () => import('@/views/profile/MessagesView.vue'),
        meta: {
          title: '消息中心',
          hidden: true,
          breadcrumb: [{ title: '首页' }, { title: '消息中心' }],
        },
      },
    ],
  },
]

// 404路由
export const notFoundRoute: RouteRecordRaw = {
  path: '/:pathMatch(.*)*',
  name: 'NotFound',
  redirect: '/dashboard',
}

// 所有路由
export const routes: RouteRecordRaw[] = [...publicRoutes, ...privateRoutes, notFoundRoute]
