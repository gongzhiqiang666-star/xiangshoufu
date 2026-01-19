# PC端设计文档补充计划

## 一、背景信息

### 当前状态
- **后端**: Go语言，已实现回调处理、分润计算、消息通知
- **PC端**: 设计文档已完成，代码目录为空
- **APP端**: 设计文档已完成，代码目录为空

### 决策确认
- ✅ **共用一套后端** - PC端和APP端共用Go后端
- ✅ **共用账号体系** - 同一个代理商账号可以同时登录PC端和APP端

---

## 二、需要补充的设计文档内容

### 2.1 后端API架构设计

#### API路由结构
```
/api/v1/
├── auth/              # 认证（公开）
├── common/            # 公共API（PC/APP共用）
│   ├── dashboard/     # 仪表盘
│   ├── agents/        # 代理商
│   ├── terminals/     # 终端
│   ├── merchants/     # 商户
│   ├── transactions/  # 交易
│   ├── profits/       # 分润
│   ├── wallets/       # 钱包
│   ├── withdrawals/   # 提现
│   ├── messages/      # 消息
│   └── marketing/     # 营销
├── pc/                # PC端专属
│   ├── channels/      # 通道管理
│   ├── policies/      # 政策模板管理
│   ├── terminals/cross/  # 跨级终端操作
│   ├── adjustments/   # 手动调账
│   ├── withdrawals/audit/  # 提现审核
│   ├── reports/       # 数据报表
│   └── system/        # 系统设置
└── app/               # APP端专属
    ├── home/          # 首页聚合
    ├── terminals/direct/  # 直属终端操作
    ├── cargo-deduction/   # 货款代扣
    └── qrcode/        # 推广二维码
```

---

### 2.2 认证授权设计

#### JWT Token方案
| 项目 | PC端 | APP端 |
|------|------|-------|
| Access Token有效期 | 2小时 | 24小时 |
| Refresh Token有效期 | 7天 | 30天 |
| 存储位置 | HttpOnly Cookie + 内存 | Keychain/Keystore |

#### Token Claims结构
```go
type Claims struct {
    UserID      int64      // 用户ID
    AgentID     int64      // 代理商ID
    AgentNo     string     // 代理商编号
    RoleID      int64      // 角色ID
    RoleCode    string     // 角色编码
    TokenType   string     // access/refresh
    ClientType  string     // pc/app
    DeviceID    string     // 设备ID
    Permissions []string   // 权限列表
}
```

#### 多端登录策略
- **策略**: 同端互踢（PC踢PC，APP踢APP）
- **PC端最大设备数**: 1
- **APP端最大设备数**: 2

---

### 2.3 RBAC权限设计

#### 角色定义
| 角色编码 | 角色名称 | 适用端 |
|----------|----------|--------|
| super_admin | 超级管理员 | PC |
| operator | 运营管理员 | PC |
| finance | 财务人员 | PC |
| customer | 客服人员 | PC |
| agent_level1 | 一级代理商 | PC/APP |
| agent_level2 | 二级代理商 | PC/APP |
| agent_level3 | 三级代理商 | PC/APP |

#### 权限矩阵（关键权限）
| 权限 | 超管 | 运营 | 财务 | 客服 | 代理商 |
|------|------|------|------|------|--------|
| channel:edit | ✅ | ❌ | ❌ | ❌ | ❌ |
| policy:edit | ✅ | ✅ | ❌ | ❌ | ❌ |
| terminal:cross_level | ✅ | ✅ | ❌ | ❌ | ❌ |
| profit:adjust | ✅ | ❌ | ✅ | ❌ | ❌ |
| withdraw:audit | ✅ | ❌ | ✅ | ❌ | ❌ |
| system:edit | ✅ | ❌ | ❌ | ❌ | ❌ |

#### 数据权限
| 角色 | 数据范围 |
|------|----------|
| 超管/运营/财务/客服 | 全部数据 |
| 代理商 | 自己及下级（物化路径判断） |

---

### 2.4 安全设计

#### Token存储
- **PC端Access Token**: 内存变量（防XSS）
- **PC端Refresh Token**: HttpOnly Cookie（Secure=true, SameSite=Strict）
- **APP端**: iOS Keychain / Android EncryptedSharedPreferences

#### XSS/CSRF防护
| 防护措施 | 实现方式 |
|----------|----------|
| XSS | CSP头、X-XSS-Protection头、Vue自动转义 |
| CSRF | Origin/Referer验证、SameSite Cookie |
| 点击劫持 | X-Frame-Options: DENY |

#### 敏感数据脱敏规则
| 字段类型 | 脱敏规则 | 示例 |
|----------|----------|------|
| 手机号 | 前3后4 | 138****8888 |
| 身份证 | 前3后4 | 110***********1234 |
| 银行卡 | 后4位 | ****5678 |
| 姓名 | 中间* | 张*明 |

---

### 2.5 错误处理规范

#### 错误码规范（5位数字）
```
AB CCC
│  └── 错误编号 (001-999)
└───── 模块编号 (10-99)
```

| 范围 | 模块 |
|------|------|
| 10xxx | 通用错误 |
| 40xxx | 认证授权 |
| 41xxx | 代理商 |
| 42xxx | 终端 |
| 43xxx | 商户 |
| 44xxx | 交易 |
| 45xxx | 分润 |
| 46xxx | 钱包 |
| 47xxx | 政策 |

#### 统一响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "request_id": "xxx",
  "timestamp": 1234567890
}
```

---

### 2.6 性能优化

#### 分页规范
- 默认页大小: 20
- 最大页大小: 100
- 大数据量使用游标分页

#### 缓存策略
| 数据类型 | 缓存时间 | 失效策略 |
|----------|----------|----------|
| 用户信息 | 30分钟 | 更新时清除 |
| 代理商信息 | 30分钟 | 更新时清除 |
| 政策模板 | 2小时 | 更新时清除 |
| 费率信息 | 2小时 | 政策变更时清除 |
| 仪表盘统计 | 5分钟 | 定时刷新 |

---

## 三、实施步骤

### Step 1: 更新PC端设计文档
**文件**: `design/PC端管理功能详细设计.md`

添加以下章节：
1. 后端API架构设计
2. 认证授权设计（JWT方案、多端登录策略）
3. RBAC权限设计（角色、权限、数据权限）
4. 安全设计（Token存储、XSS/CSRF防护、数据脱敏）
5. 错误处理规范（错误码、响应格式）
6. 性能优化（分页、缓存策略）

### Step 2: 创建后端API设计文档
**文件**: `design/后端API设计.md`

包含：
1. API路由完整定义
2. 请求/响应数据结构
3. 认证中间件说明
4. 权限中间件说明

### Step 3: 创建安全设计文档
**文件**: `design/安全设计规范.md`

包含：
1. 认证流程图
2. Token刷新流程
3. 安全配置清单
4. 敏感数据处理规范

---

## 四、关键文件

| 文件 | 操作 |
|------|------|
| `design/PC端管理功能详细设计.md` | 修改 - 补充缺失章节 |
| `design/后端API设计.md` | 新建 - API完整定义 |
| `design/安全设计规范.md` | 新建 - 安全规范 |

---

## 五、验证方式

1. 检查设计文档完整性 - 所有模块都有对应设计
2. 检查API定义一致性 - PC端和APP端API无冲突
3. 检查权限覆盖 - 所有敏感操作都有权限控制
4. 检查安全措施 - Token、XSS、CSRF等都有明确方案
