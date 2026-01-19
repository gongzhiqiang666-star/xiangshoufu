// 菜单配置
export interface MenuItem {
  path: string
  title: string
  icon: string
  children?: MenuItem[]
}

export const MENU_LIST: MenuItem[] = [
  {
    path: '/dashboard',
    title: '仪表盘',
    icon: 'Odometer',
  },
  {
    path: '/agents',
    title: '代理管理',
    icon: 'User',
    children: [
      { path: '/agents/list', title: '代理列表', icon: '' },
    ],
  },
  {
    path: '/merchants',
    title: '商户管理',
    icon: 'Shop',
    children: [
      { path: '/merchants/list', title: '商户列表', icon: '' },
    ],
  },
  {
    path: '/terminals',
    title: '终端管理',
    icon: 'Monitor',
    children: [
      { path: '/terminals/list', title: '终端列表', icon: '' },
    ],
  },
  {
    path: '/transactions',
    title: '交易管理',
    icon: 'Money',
    children: [
      { path: '/transactions/list', title: '交易记录', icon: '' },
    ],
  },
  {
    path: '/profits',
    title: '分润管理',
    icon: 'TrendCharts',
    children: [
      { path: '/profits/list', title: '分润明细', icon: '' },
    ],
  },
  {
    path: '/wallets',
    title: '钱包管理',
    icon: 'Wallet',
    children: [
      { path: '/wallets/list', title: '钱包总览', icon: '' },
    ],
  },
  {
    path: '/policies',
    title: '政策管理',
    icon: 'Document',
    children: [
      { path: '/policies/list', title: '政策模板', icon: '' },
    ],
  },
  {
    path: '/system',
    title: '系统管理',
    icon: 'Setting',
    children: [
      { path: '/system/users', title: '用户管理', icon: '' },
      { path: '/system/logs', title: '操作日志', icon: '' },
    ],
  },
]

// 用户角色名称映射
export const ROLE_NAMES: Record<string, string> = {
  admin: '管理员',
  finance: '财务',
  operation: '运营',
  readonly: '只读用户',
}

// 状态映射
export const STATUS_MAP: Record<number, { text: string; type: string }> = {
  1: { text: '正常', type: 'success' },
  0: { text: '禁用', type: 'danger' },
}
