# PC端剩余功能开发计划

> 制定日期: 2026-01-19
> 更新日期: 2026-01-19
> 状态: ✅ 全部完成

---

## 一、已完成模块 ✅

| 模块 | 功能 | 状态 |
|------|------|------|
| 代扣管理 | Q6伙伴代扣、Q7每日扣款、Q8多钱包优先级 | ✅ 已完成 |
| 终端下发 | Q16跨级代扣链、Q29 APP/PC权限控制 | ✅ 已完成 |
| 流量费返现 | Q30三档返现、级差计算 | ✅ 已完成 |
| 回调处理 | 通道回调接收、解析、入库 | ✅ 已完成 |
| 分润计算 | 交易分润、级差分润 | ✅ 已完成 |
| 消息服务 | 系统消息、告警通知 | ✅ 已完成 |
| **认证授权** | JWT登录、Token刷新、权限中间件 | ✅ 已完成 |
| **代理商管理** | 详情、下级列表、团队树 | ✅ 已完成 |
| **钱包管理** | 余额查询、流水、提现 | ✅ 已完成 |
| **分润记录** | 分润列表、统计 | ✅ 已完成 |
| **交易查询** | 交易列表、统计、趋势 | ✅ 已完成 |
| **数据看板** | 概览、图表 | ✅ 已完成 |
| **消息中心** | 消息列表、已读标记 | ✅ 已完成 |
| **商户管理** | 商户列表、详情、统计、交易 | ✅ 已完成 |
| **终端管理** | 终端列表、详情、统计 | ✅ 已完成 |
| **政策模板** | 模板列表、详情、我的政策 | ✅ 已完成 |

---

## 二、待开发模块

### 第一优先级：核心基础模块（Day 1-2）

#### 1. 认证授权模块 (Auth)
```
功能点:
├── JWT登录认证
├── Token刷新机制
├── 登出处理
├── 认证中间件
└── 权限控制中间件

API端点:
├── POST /api/v1/auth/login        登录
├── POST /api/v1/auth/logout       登出
├── POST /api/v1/auth/refresh      刷新Token
└── GET  /api/v1/auth/profile      获取当前用户信息
```

#### 2. 代理商管理模块 (Agent)
```
功能点:
├── 代理商信息查询
├── 下级代理商列表
├── 代理商统计数据
├── 团队层级树
└── 邀请码/二维码管理

API端点:
├── GET  /api/v1/agents/:id           获取代理商详情
├── GET  /api/v1/agents/subordinates  获取下级代理商列表
├── GET  /api/v1/agents/team-tree     获取团队层级树
├── GET  /api/v1/agents/stats         获取代理商统计
├── PUT  /api/v1/agents/profile       更新代理商资料
└── GET  /api/v1/agents/invite-code   获取邀请码
```

### 第二优先级：资金模块（Day 2-3）

#### 3. 钱包管理模块 (Wallet)
```
功能点:
├── 多钱包余额查询
├── 钱包流水明细
├── 提现申请
├── 提现记录
└── 冻结/解冻记录

API端点:
├── GET  /api/v1/wallets                获取钱包列表
├── GET  /api/v1/wallets/:id/logs       获取钱包流水
├── POST /api/v1/wallets/:id/withdraw   申请提现
├── GET  /api/v1/withdrawals            获取提现记录
└── GET  /api/v1/wallets/summary        钱包汇总统计
```

#### 4. 分润记录模块 (Profit)
```
功能点:
├── 分润记录列表
├── 分润统计（日/周/月）
├── 分润类型筛选
└── 分润来源追溯

API端点:
├── GET  /api/v1/profits                获取分润记录
├── GET  /api/v1/profits/stats          分润统计
├── GET  /api/v1/profits/daily          每日分润汇总
└── GET  /api/v1/profits/:id            分润详情
```

### 第三优先级：业务查询模块（Day 3-4）

#### 5. 交易查询模块 (Transaction)
```
功能点:
├── 交易流水查询
├── 交易统计（金额/笔数）
├── 多维度筛选
├── 交易趋势图表
└── 交易导出

API端点:
├── GET  /api/v1/transactions           交易列表
├── GET  /api/v1/transactions/:id       交易详情
├── GET  /api/v1/transactions/stats     交易统计
├── GET  /api/v1/transactions/trend     交易趋势
└── GET  /api/v1/transactions/export    导出交易
```

#### 6. 商户管理模块 (Merchant)
```
功能点:
├── 商户列表查询
├── 商户详情
├── 商户统计
└── 商户交易汇总

API端点:
├── GET  /api/v1/merchants              商户列表
├── GET  /api/v1/merchants/:id          商户详情
├── GET  /api/v1/merchants/stats        商户统计
└── GET  /api/v1/merchants/:id/transactions  商户交易
```

#### 7. 终端管理模块增强 (Terminal)
```
功能点:
├── 终端列表查询
├── 终端详情
├── 终端统计（库存/已下发/已激活）
├── 终端入库（待明确）
└── 终端激活记录

API端点:
├── GET  /api/v1/terminals              终端列表
├── GET  /api/v1/terminals/:sn          终端详情
├── GET  /api/v1/terminals/stats        终端统计
├── POST /api/v1/terminals/import       终端入库（待明确）
└── GET  /api/v1/terminals/:sn/history  终端历史
```

### 第四优先级：运营模块（Day 4-5）

#### 8. 数据看板模块 (Dashboard)
```
功能点:
├── 今日数据概览
├── 交易趋势图
├── 分润趋势图
├── 下级排行榜
└── 待办事项提醒

API端点:
├── GET  /api/v1/dashboard/overview     数据概览
├── GET  /api/v1/dashboard/charts       图表数据
├── GET  /api/v1/dashboard/ranking      排行榜
└── GET  /api/v1/dashboard/todos        待办事项
```

#### 9. 政策模板模块 (Policy)
```
功能点:
├── 政策模板列表
├── 模板详情
├── 下级政策分配
└── 政策修改（待明确权限）

API端点:
├── GET  /api/v1/policies/templates           模板列表
├── GET  /api/v1/policies/templates/:id       模板详情
├── PUT  /api/v1/policies/subordinate/:id     修改下级政策（待明确）
└── GET  /api/v1/policies/my                  我的政策
```

#### 10. 消息中心模块 (Message)
```
功能点:
├── 消息列表
├── 未读消息数
├── 标记已读
├── 消息详情

API端点:
├── GET  /api/v1/messages               消息列表
├── GET  /api/v1/messages/unread-count  未读数量
├── PUT  /api/v1/messages/:id/read      标记已读
└── PUT  /api/v1/messages/read-all      全部已读
```

---

## 三、待明确事项处理

| 事项 | 优先级 | 处理方式 |
|------|--------|----------|
| 终端入库方式 | 低 | 预留接口，待产品确认后实现 |
| 代理商能否修改下级政策模板 | 中 | 实现接口，添加权限开关控制 |
| T+0秒到费用分配规则 | 中 | 在交易处理中预留扩展点 |

---

## 四、开发顺序

```
Day 1: 认证授权模块 + 代理商管理模块
       ├── auth_handler.go
       ├── auth_middleware.go
       ├── agent_handler.go
       └── agent_service.go

Day 2: 钱包管理模块 + 分润记录模块
       ├── wallet_handler.go
       ├── wallet_service.go
       ├── profit_handler.go
       └── profit_service.go (扩展)

Day 3: 交易查询模块 + 商户管理模块
       ├── transaction_handler.go
       ├── transaction_service.go
       ├── merchant_handler.go
       └── merchant_service.go

Day 4: 终端管理增强 + 数据看板
       ├── terminal_handler.go (扩展)
       ├── terminal_service.go
       ├── dashboard_handler.go
       └── dashboard_service.go

Day 5: 政策模板 + 消息中心 + 集成测试
       ├── policy_handler.go
       ├── message_handler.go
       └── 集成测试
```

---

## 五、技术要点

1. **认证**: JWT RS256签名，Token有效期2小时，Refresh Token 7天
2. **权限**: RBAC模型，代理商角色区分（普通/管理员）
3. **分页**: 统一分页格式，limit/offset + 游标分页
4. **缓存**: 热点数据Redis缓存（代理商信息、政策配置）
5. **导出**: 大数据量异步导出，生成文件后通知下载

---

## 六、文件清单

### 新增文件
```
internal/
├── handler/
│   ├── auth_handler.go
│   ├── agent_handler.go
│   ├── wallet_handler.go
│   ├── profit_handler.go
│   ├── transaction_handler.go
│   ├── merchant_handler.go
│   ├── terminal_handler.go
│   ├── dashboard_handler.go
│   ├── policy_handler.go
│   └── message_handler.go
├── service/
│   ├── auth_service.go
│   ├── agent_service.go
│   ├── wallet_service.go
│   ├── transaction_service.go
│   ├── merchant_service.go
│   ├── terminal_service.go
│   ├── dashboard_service.go
│   └── policy_service.go
├── middleware/
│   ├── auth_middleware.go
│   └── permission_middleware.go
├── models/
│   ├── user.go
│   └── withdrawal.go
└── repository/
    ├── user_repo.go
    ├── withdrawal_repo.go
    └── merchant_repo.go

migrations/
├── 009_create_users_table.sql
└── 010_create_withdrawals_table.sql
```
