# APP详细设计文档 - 享收付代理商分润管理系统

## 文档信息

| 项目 | 内容 |
|-----|------|
| **项目名称** | 享收付 - 代理商分润管理系统APP |
| **版本** | v1.0.0 |
| **最后更新** | 2025-01-19 |
| **开发框架** | Flutter 3.x |
| **UI组件库** | Bruno（贝壳找房）+ TDesign（腾讯）|

---

## 一、技术架构

### 1.1 技术栈

| 层级 | 技术选型 | 说明 |
|------|----------|------|
| **开发框架** | Flutter 3.x | 跨平台：iOS/Android/鸿蒙 |
| **开发语言** | Dart 3.x | Flutter官方语言 |
| **状态管理** | Riverpod 2.0 | 类型安全、简洁 |
| **路由管理** | go_router | 官方推荐声明式路由 |
| **UI组件库** | Bruno + TDesign | 中国本土化组件 |
| **网络请求** | Dio + Retrofit | 强大的HTTP客户端 |
| **本地存储** | SharedPreferences + SQLite | 配置+结构化数据 |
| **图表** | fl_chart | Flutter原生图表 |
| **推送** | 极光推送 | 国内推送服务 |

### 1.2 项目结构

```
mobileapp/
├── lib/
│   ├── main.dart                      # 入口文件
│   │
│   ├── core/                          # 核心模块
│   │   ├── theme/                     # 主题设计系统
│   │   │   ├── app_colors.dart        # 颜色规范
│   │   │   ├── app_typography.dart    # 字体规范
│   │   │   ├── app_spacing.dart       # 间距规范
│   │   │   ├── app_decorations.dart   # 圆角/阴影规范
│   │   │   └── app_theme.dart         # 主题配置
│   │   ├── network/                   # 网络层
│   │   │   ├── api_client.dart        # API客户端
│   │   │   ├── api_endpoints.dart     # 接口地址
│   │   │   └── interceptors/          # 拦截器
│   │   ├── storage/                   # 本地存储
│   │   ├── utils/                     # 工具类
│   │   │   ├── format_utils.dart      # 格式化工具
│   │   │   ├── validators.dart        # 表单验证
│   │   │   └── platform_utils.dart    # 平台判断
│   │   └── constants/                 # 常量定义
│   │
│   ├── shared/                        # 共享组件
│   │   ├── widgets/                   # 通用组件
│   │   │   ├── buttons/               # 按钮组件
│   │   │   ├── cards/                 # 卡片组件
│   │   │   ├── inputs/                # 输入组件
│   │   │   ├── dialogs/               # 弹窗组件
│   │   │   ├── charts/                # 图表组件
│   │   │   ├── loading/               # 加载组件
│   │   │   ├── tags/                  # 标签组件
│   │   │   ├── empty/                 # 空状态组件
│   │   │   └── main_scaffold.dart     # 主框架
│   │   └── extensions/                # 扩展方法
│   │
│   ├── features/                      # 功能模块
│   │   ├── auth/                      # 认证模块
│   │   ├── home/                      # 首页
│   │   ├── agent/                     # 代理拓展
│   │   ├── terminal/                  # 终端管理
│   │   ├── cargo_deduction/           # 货款代扣
│   │   ├── merchant/                  # 商户管理
│   │   ├── data_analysis/             # 数据分析
│   │   ├── profit/                    # 收益统计
│   │   ├── wallet/                    # 钱包
│   │   ├── deduction/                 # 代扣管理
│   │   ├── marketing/                 # 营销海报
│   │   ├── message/                   # 消息通知
│   │   └── profile/                   # 我的信息
│   │
│   └── router/                        # 路由配置
│       └── app_router.dart
│
├── assets/                            # 资源文件
│   ├── images/                        # 图片
│   ├── icons/                         # 图标
│   └── fonts/                         # 字体
│
├── pubspec.yaml                       # 依赖配置
└── README.md
```

---

## 二、设计系统

### 2.1 颜色规范

#### 主色系
| 名称 | 色值 | 用途 |
|------|------|------|
| Primary | `#2563EB` | 品牌主色、主要按钮 |
| Primary Light | `#60A5FA` | 次要强调 |
| Primary Dark | `#1D4ED8` | 按压状态 |

#### 功能色
| 名称 | 色值 | 用途 |
|------|------|------|
| Success | `#10B981` | 成功、收益增长 |
| Warning | `#F59E0B` | 警告、提醒 |
| Danger | `#EF4444` | 错误、亏损 |
| Info | `#3B82F6` | 信息提示 |

#### 分润类型专属色
| 类型 | 色值 | 图标 |
|------|------|------|
| 交易分润 | `#2563EB` | swap_horiz |
| 押金返现 | `#10B981` | monetization_on |
| 流量返现 | `#F59E0B` | signal_cellular_alt |
| 激活奖励 | `#8B5CF6` | card_giftcard |

#### 支付类型色
| 类型 | 色值 |
|------|------|
| 微信支付 | `#07C160` |
| 支付宝 | `#1677FF` |
| 银联 | `#E60012` |

### 2.2 字体规范

| 类型 | 字号 | 字重 | 用途 |
|------|------|------|------|
| H1 | 24px | Bold | 页面主标题 |
| H2 | 20px | SemiBold | 区块标题 |
| H3 | 18px | SemiBold | 卡片标题 |
| Body1 | 16px | Regular | 主要正文 |
| Body2 | 14px | Regular | 次要正文 |
| Caption | 12px | Regular | 辅助说明 |
| Amount Large | 32px | Bold | 大额金额 |
| Amount Medium | 24px | Bold | 中额金额 |
| Amount Small | 18px | SemiBold | 列表金额 |

### 2.3 间距规范

| 名称 | 数值 | 用途 |
|------|------|------|
| XS | 4px | 最小间距 |
| SM | 8px | 小间距 |
| MD | 16px | 常用间距 |
| LG | 24px | 大间距 |
| XL | 32px | 超大间距 |
| Page Padding | 16px | 页面边距 |
| Card Gap | 12px | 卡片间距 |

### 2.4 圆角规范

| 类型 | 数值 |
|------|------|
| 标签 | 4px |
| 按钮/输入框 | 8px |
| 卡片 | 12px |
| 底部弹窗 | 16px |

---

## 三、页面功能设计

### 3.1 首页（Home）

**页面路径**: `/`

**功能说明**:
- 展示今日收益总览
- 分润明细（交易分润、押金返现、流量返现、激活奖励）
- 快捷入口（8个功能入口）
- 最近交易记录
- 轮播图/公告

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/dashboard/today` | GET | 今日数据统计 |
| `/api/v1/dashboard/profit-detail` | GET | 分润明细 |
| `/api/v1/transactions/recent` | GET | 最近交易 |
| `/api/v1/banners` | GET | 轮播图 |

**Bruno组件使用**:
- `BrnEnhanceNumberCard` - 今日收益卡片
- `BrnBannerWidget` - 轮播图

---

### 3.2 终端管理（Terminal）

**页面路径**: `/terminal`

**功能说明**:
- 终端统计（总数、已激活、未激活、今日激活）
- 终端列表（支持筛选：全部/已激活/未激活/库存）
- 终端划拨（仅限直属下级）
- 终端回拨（仅限直属下级）
- 批量设置费率/流量卡/押金

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/terminals/stats` | GET | 终端统计 |
| `/api/v1/terminals` | GET | 终端列表 |
| `/api/v1/terminals/transfer` | POST | 终端划拨 |
| `/api/v1/terminals/recall` | POST | 终端回拨 |
| `/api/v1/terminals/batch-settings` | POST | 批量设置 |

**业务规则**:
1. APP端仅支持划拨给直属下级，不可跨级
2. 已激活终端不可回拨
3. 划拨时可设置货款代扣

---

### 3.3 货款代扣（Cargo Deduction）

**页面路径**: `/cargo-deduction`

**功能说明**:
- 待接收：显示上级发起的代扣请求，可接收/拒绝
- 进行中：显示正在执行的代扣，显示进度
- 已完成：历史代扣记录

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/cargo-deductions` | GET | 代扣列表 |
| `/api/v1/cargo-deductions/{id}/accept` | POST | 接收代扣 |
| `/api/v1/cargo-deductions/{id}/reject` | POST | 拒绝代扣 |

**业务规则**:
1. 代扣需要下级接收后才生效
2. 可选择从分润钱包/服务费钱包/奖励钱包扣除
3. 扣款进度实时更新

---

### 3.4 商户管理（Merchant）

**页面路径**: `/merchant`

**功能说明**:
- 直营商户：直接拓展的商户
- 团队商户：下级代理拓展的商户
- 商户搜索（按名称/编号/机具号）
- 商户详情（基本信息、交易统计、费率设置）
- 费率修改

**商户分类**:
| 类型 | 条件 |
|------|------|
| 忠诚商户 | 月均 > 5万 |
| 优质商户 | 3万 ≤ 月均 < 5万 |
| 潜力商户 | 2万 ≤ 月均 < 3万 |
| 一般商户 | 1万 ≤ 月均 < 2万 |
| 低活跃商户 | 0 < 月均 < 1万 |
| 30天无交易 | 30天内无交易 |

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/merchants` | GET | 商户列表 |
| `/api/v1/merchants/{id}` | GET | 商户详情 |
| `/api/v1/merchants/{id}/rate` | PUT | 修改费率 |
| `/api/v1/merchants/{id}/transactions` | GET | 交易记录 |

---

### 3.5 钱包（Wallet）

**页面路径**: `/wallet`

**功能说明**:
- 总资产展示
- 分类钱包（分润钱包、服务费钱包、奖励钱包）
- 按通道筛选
- 申请提现
- 钱包流水
- 提现记录

**钱包类型**:
| 类型 | 说明 |
|------|------|
| 分润钱包 | 交易分润收入 |
| 服务费钱包 | 流量费+押金返现 |
| 奖励钱包 | 激活奖励等 |
| 充值钱包 | 给下级的额外奖励（可选） |
| 沉淀钱包 | 下级未提现资金（可选） |

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/wallets` | GET | 钱包列表 |
| `/api/v1/wallets/{id}/withdraw` | POST | 申请提现 |
| `/api/v1/wallets/{id}/flows` | GET | 钱包流水 |
| `/api/v1/withdrawals` | GET | 提现记录 |

**提现规则**:
1. 每个钱包有独立的提现门槛
2. 提现需扣除税筹费用（如9%+3元/笔）
3. 通过代付通道打款

---

### 3.6 收益统计（Profit）

**页面路径**: `/profit`

**功能说明**:
- 今日收益及趋势
- 收益明细（按分润类型）
- 收益趋势图（7天/30天）
- 月收益统计（6月/1年/2年）

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/profits/today` | GET | 今日收益 |
| `/api/v1/profits/trend` | GET | 收益趋势 |
| `/api/v1/profits/monthly` | GET | 月度收益 |

---

### 3.7 代理拓展（Agent）

**页面路径**: `/agent`

**功能说明**:
- 推广二维码展示
- 复制邀请链接
- 保存二维码到相册
- 手动添加代理
- 直属代理列表
- 团队统计

**数据接口**:
| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/agent/qrcode` | GET | 获取推广二维码 |
| `/api/v1/agent/invite-link` | GET | 获取邀请链接 |
| `/api/v1/agent/subordinates` | GET | 下级代理列表 |
| `/api/v1/agent/add` | POST | 手动添加代理 |

---

### 3.8 消息通知（Message）

**页面路径**: `/message`

**功能说明**:
- 消息分类（全部/分润/注册/消费/系统）
- 消息列表（按时间分组）
- 未读标记
- 消息3天自动过期

**消息类型**:
| 类型 | 图标 | 颜色 |
|------|------|------|
| 分润到账 | 💰 | 绿色 |
| 新代理注册 | 👤 | 蓝色 |
| 交易通知 | 💳 | 蓝色 |
| 系统公告 | 📢 | 橙色 |

---

### 3.9 营销海报（Marketing）

**页面路径**: `/marketing`

**功能说明**:
- 海报分类筛选
- 海报瀑布流展示
- 保存到相册

---

### 3.10 我的信息（Profile）

**页面路径**: `/profile`

**功能说明**:
- 个人信息展示（脱敏）
- 结算卡管理
- 费率成本查看
- 邀请码自定义
- 退出登录

---

## 四、通用组件

### 4.1 统计卡片（StatCard）

```dart
StatCard(
  title: '今日收益',
  value: '¥1,234.56',
  icon: Icons.trending_up,
  iconColor: AppColors.success,
  trend: '+12.5%',
  onTap: () {},
)
```

### 4.2 钱包卡片（WalletCard）

```dart
WalletCard(
  name: '分润钱包',
  channel: '拉卡拉',
  balance: 5680.00,
  threshold: 100,
  canWithdraw: true,
  gradientColors: AppColors.walletProfitGradient,
  onWithdraw: () {},
)
```

### 4.3 交易列表项（TransactionItem）

```dart
TransactionItem(
  merchantName: '张三商店',
  amount: 1500.00,
  time: '10:30',
  type: 'credit', // credit/debit/wechat/alipay
  onTap: () {},
)
```

### 4.4 分润类型标签（ProfitTypeTag）

```dart
ProfitTypeTag(type: ProfitType.trade)
// 交易分润 - 蓝色
// 押金返现 - 绿色
// 流量返现 - 橙色
// 激活奖励 - 紫色
```

### 4.5 空状态（EmptyState）

```dart
EmptyState(
  icon: Icons.inbox_outlined,
  title: '暂无数据',
  description: '暂时没有相关记录',
  buttonText: '刷新',
  onButtonTap: () {},
)
```

---

## 五、开发计划

### Phase 1: 基础框架（1周）
- [x] 项目初始化
- [x] 设计系统实现
- [x] 路由配置
- [x] 主框架搭建
- [ ] 网络层封装
- [ ] 本地存储封装

### Phase 2: 认证模块（0.5周）
- [ ] 登录页面
- [ ] 验证码功能
- [ ] Token管理
- [ ] 自动登录

### Phase 3: 首页模块（1周）
- [x] 首页布局
- [ ] 数据接口对接
- [ ] 轮播图
- [ ] 下拉刷新

### Phase 4: 终端管理（1.5周）
- [x] 终端列表
- [x] 终端划拨页面
- [ ] 终端回拨
- [ ] 批量设置
- [ ] 货款代扣

### Phase 5: 商户管理（1周）
- [ ] 商户列表
- [ ] 商户详情
- [ ] 费率修改
- [ ] 交易记录

### Phase 6: 钱包模块（1周）
- [x] 钱包列表
- [ ] 申请提现
- [ ] 钱包流水
- [ ] 提现记录

### Phase 7: 其他模块（2周）
- [ ] 收益统计
- [ ] 数据分析
- [ ] 代理拓展
- [ ] 消息通知
- [ ] 营销海报
- [ ] 我的信息

### Phase 8: 联调测试（2周）
- [ ] 接口联调
- [ ] 功能测试
- [ ] 性能优化
- [ ] Bug修复

---

## 六、注意事项

### 6.1 金额处理
- 后端金额使用**分**为单位存储
- 前端展示时转换为**元**
- 使用`FormatUtils.formatCents()`统一处理

### 6.2 数据脱敏
- 手机号：138****8888
- 身份证：110***********34
- 银行卡：**** **** **** 5678
- 使用`FormatUtils.maskXXX()`方法

### 6.3 敏感数据存储
- Token使用加密存储
- 用户敏感信息不缓存本地

### 6.4 平台适配
- iOS: APNs推送
- Android: 极光推送
- HarmonyOS: 华为Push Kit

---

*文档版本: v1.0*
*最后更新: 2025-01-19*
