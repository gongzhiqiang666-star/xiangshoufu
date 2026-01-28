# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

享收付 (XiangShouFu) - An agent profit-sharing management system that processes payment channel callbacks from multiple payment providers. The system handles transaction callbacks, calculates profit sharing across agent hierarchies, and manages wallets.

---

## ⚠️ 重要行为规则（必读）

### 1. "开发完成"的定义

必须满足以下**全部条件**才能声明"开发完成"：

| 条件 | 说明 |
|------|------|
| ✅ 编译通过 | `go build` / `npm run build` / `flutter analyze` |
| ✅ 测试通过 | `go test` / `npm run test:run` / `flutter test` |
| ✅ 集成配置完成 | 路由注册、依赖添加、定时任务启动等 |
| ✅ 验证命令已实际执行 | 必须有真实的命令输出结果 |

### 2. 禁止假设性完成

| ❌ 禁止 | ✅ 正确 |
|--------|--------|
| "代码写完了，应该可以用了" | "代码写完了，我已验证：编译通过、测试通过、路由已注册" |
| "理论上没问题" | "实际运行验证通过" |
| "按照规范写的，应该OK" | "已执行验证命令，输出结果如下..." |

### 3. 每次开发结束必须输出验证报告

```
## ✅ 验证报告
- 后端编译: ✅ 通过
- 后端测试: ✅ 通过 (X passed)
- PC端编译: ✅ 通过（如涉及）
- APP分析: ✅ 通过（如涉及）
- 路由注册: ✅ 已检查
- 定时任务: ✅ 已注册（如涉及）
- 数据库迁移: ✅ 已创建（如涉及）
```

---

## 🚨 开发完成前必须执行的验证（强制）

### 声明"开发完成"之前，必须执行以下所有验证：

#### 1. 后端验证（必须）

```bash
cd server && go build ./...           # 编译必须通过
cd server && go test ./... -v         # 测试必须通过
```

#### 2. PC端验证（如涉及前端改动）

```bash
cd web && npm run build               # 编译必须通过
cd web && npm run test:run            # 测试必须通过（如有）
```

#### 3. APP端验证（如涉及APP改动）

```bash
cd mobileapp && flutter pub get       # 依赖必须安装成功
cd mobileapp && flutter analyze       # 静态分析必须通过
cd mobileapp && flutter test          # 测试必须通过
```

#### 4. 集成配置检查清单

| 检查项 | 验证方法 |
|--------|---------|
| 新Handler是否注册路由 | 搜索 `main.go` 或 `routes.go` 中的路由注册 |
| 新依赖是否添加 | 检查 `go.mod` / `package.json` / `pubspec.yaml` |
| 新定时任务是否注册 | 搜索 `setupScheduler` 或 `jobs/` 目录 |
| 数据库迁移是否创建 | 检查 `migrations/` 目录是否有新文件 |
| 前端路由是否配置 | 检查 `router/routes.ts` 或路由配置文件 |
| 环境变量是否文档化 | 检查是否需要更新环境变量说明 |

#### 5. 运行时验证（重要！）

> ⚠️ **编译通过 ≠ 代码正确**
>
> 以下问题**只能通过实际运行**才能发现，编译和测试都无法捕获：

| 问题类型 | 示例 | 为什么编译不报错 |
|----------|------|------------------|
| API路径错误 | `/api/api/users` 双重前缀 | 字符串拼接语法正确 |
| import路径错误 | `@/api/user` vs `@/api/users` | 路径语法正确，文件不存在时才报错 |
| Vue prop类型错误 | `:visible="dialogVisible"` vs `v-model:visible` | 模板语法正确 |
| 运行时依赖缺失 | Provider未导出、常量未定义 | 类型定义存在但值不存在 |

**前端改动必须执行：**

```bash
# 1. 启动后端
cd server && go run cmd/server/main.go

# 2. 启动前端（另一个终端）
cd web && npm run dev

# 3. 打开浏览器验证
# - 打开 DevTools → Network 面板
# - 执行相关操作
# - 检查：请求路径是否正确、响应状态是否200
# - 打开 Console 面板，确认无红色错误
```

**APP改动必须执行：**

```bash
# 1. 启动后端
cd server && go run cmd/server/main.go

# 2. 运行APP
cd mobileapp && flutter run

# 3. 验证
# - 执行相关功能
# - 检查日志输出无错误
# - 网络请求正确返回
```

---

### 🔍 常见问题自查表

> 基于历史Bug总结，每次提交前必须检查：

#### 前端（Vue/TypeScript）

| 检查项 | 常见错误 | 正确做法 |
|--------|----------|----------|
| API路径 | `${API_BASE}/api/xxx` 导致双重 `/api` | 检查 `API_BASE` 是否已包含 `/api` |
| import路径 | 文件名单复数不一致 `user` vs `users` | 使用IDE跳转验证路径存在 |
| Element Plus对话框 | `:visible` 属性 | 使用 `v-model="visible"` |
| 路由常量 | 使用了未定义的路由名 | 检查 `router/routes.ts` 中是否定义 |
| API函数导出 | 函数定义了但未export | 检查 `export { functionName }` |

#### APP（Flutter）

| 检查项 | 常见错误 | 正确做法 |
|--------|----------|----------|
| Provider导出 | Provider定义了但未在 `providers.dart` 导出 | 检查barrel文件导出 |
| 路由注册 | 页面创建了但未注册路由 | 检查 `router.dart` |
| API路径 | 硬编码路径与后端不一致 | 使用常量定义，对照Swagger |

#### 后端（Go）

| 检查项 | 常见错误 | 正确做法 |
|--------|----------|----------|
| 路由注册 | Handler写了但没注册 | 检查 `main.go` 路由组 |
| 依赖注入 | Service未注入到Handler | 检查构造函数参数 |
| 数据库字段 | Model字段与表结构不匹配 | 对照migration文件 |

---

### ❌ 绝对禁止的行为

1. **不允许**说"开发完成"但没有执行上述验证命令
2. **不允许**说"应该能通过"但没有实际运行验证
3. **不允许**遗漏任何一项必要的检查
4. **不允许**假设配置已完成而不去实际检查
5. **不允许**只做编译验证而跳过运行时验证

---

## 重要
每次代码改造后，需要将改造的内容涉及的业务逻辑调整写给 docs/业务逻辑梳理.md ，切记只写业务流程，不要将接口设计等内容写入该文件中。
每次代码提交时都要做测试，测试参考 测试规范.md 文件，完成测试后才可以提交代码。

---

## Project Structure

```
xiangshoufu/
├── server/                      # 后端服务 (Go)
│   ├── cmd/server/main.go       # 服务入口
│   ├── internal/                # 内部模块
│   │   ├── handler/             # HTTP处理器
│   │   ├── service/             # 业务逻辑
│   │   ├── repository/          # 数据仓库
│   │   ├── channel/             # 支付通道适配器
│   │   ├── middleware/          # 中间件
│   │   ├── models/              # 数据模型
│   │   ├── async/               # 异步处理
│   │   ├── cache/               # 缓存层
│   │   └── jobs/                # 定时任务
│   ├── pkg/                     # 公共包
│   ├── migrations/              # 数据库迁移
│   ├── scripts/                 # 脚本工具
│   ├── swagger/                 # Swagger API文档
│   ├── bin/                     # 编译产物
│   ├── go.mod
│   └── go.sum
├── web/                         # PC端前台 (Vue 3)
│   ├── src/
│   └── package.json
├── mobileapp/                   # 移动端APP (Flutter)
│   ├── lib/
│   └── pubspec.yaml
├── docs/                        # 项目文档
│   ├── api/                     # API接口文档
│   ├── plans/                   # 开发计划
│   ├── 测试规范.md              # 测试规范文档
│   ├── 业务逻辑梳理.md          # 业务逻辑文档
│   ├── PC端管理功能详细设计.md  # PC端设计文档
│   └── APP设计稿.md             # APP设计文档
├── CLAUDE.md                    # Claude指引
└── README.md                    # 项目说明
```

## Tech Stack

- **后端**: Go 1.24, Gin, GORM, PostgreSQL 15+
- **PC端**: Vue 3, TypeScript, Element Plus, Pinia, Vite
- **移动端**: Flutter, Dart

## Common Commands

### 后端 (server/)

```bash
cd server

# Run the server
go run cmd/server/main.go

# Build binary
go build -o bin/server cmd/server/main.go

# Run all tests
go test ./...

# Run specific module tests
go test ./internal/channel/hengxintong/...

# Run single test function
go test ./internal/service/... -run TestWalletService -v

# Run tests with coverage
go test ./internal/service/... -cover

# Format code
go fmt ./...

# Database migrations
psql -d xiangshoufu -f migrations/000_create_core_tables.sql
```

### PC端 (web/)

```bash
cd web

# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Run all tests
npm run test:run

# Run tests in watch mode
npm run test

# Run single test file
npm run test:run src/utils/__tests__/format.test.ts

# Run tests with coverage
npm run test:coverage
```

### 移动端 (mobileapp/)

```bash
cd mobileapp

# Get dependencies
flutter pub get

# Run app
flutter run

# Run all tests
flutter test

# Run single test file
flutter test test/utils/format_test.dart

# Run tests with coverage
flutter test --coverage

# Static analysis
flutter analyze
```

## Access URLs

| 服务 | 地址 |
|------|------|
| PC端管理系统 | http://localhost:5173 |
| 后端API | http://localhost:8080 |
| Swagger文档 | http://localhost:8080/swagger/index.html |

## Default Account

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | 123456 | 管理员 |

## Architecture

### Clean Architecture Layers

1. **Handler Layer** (`server/internal/handler/`) - HTTP request handling
2. **Service Layer** (`server/internal/service/`) - Business logic
3. **Repository Layer** (`server/internal/repository/`) - Data persistence
4. **Channel Adapter Layer** (`server/internal/channel/`) - Payment provider integrations

### Key Design Patterns

- **Factory Pattern**: `AdapterFactory` creates channel adapters by channel code
- **Adapter Pattern**: `ChannelAdapter` interface normalizes different payment provider APIs
- **Pub/Sub Pattern**: `MessageQueue` for async processing

## Adding a New Payment Channel

1. Create directory: `server/internal/channel/<channel_name>/`
2. Create files:
   - `adapter.go` - Implement `ChannelAdapter` interface
   - `models.go` - Channel-specific request/response models
   - `adapter_test.go` - Unit tests
3. Register in `server/cmd/server/main.go` via `factory.Register()`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://apple@localhost:5432/xiangshoufu?sslmode=disable` |
| `SERVER_PORT` | HTTP port | `8080` |
| `HENGXINTONG_PUBLIC_KEY` | RSA public key for signature verification | - |
| `ALERT_WEBHOOK_URL` | WeChat/DingTalk webhook for alerts | - |
| `SWAGGER_ENABLED` | Enable Swagger UI | `true` |

## Supported Payment Channels

| Code | Name | Status |
|------|------|--------|
| `HENGXINTONG` | 恒信通 | Implemented |
| `LAKALA` | 拉卡拉 | Pending |
| `YEAHKA` | 乐刷 | Pending |
| `SUIXINGFU` | 随行付 | Pending |
| `LIANLIAN` | 连连支付 | Pending |
| `SANDPAY` | 杉德支付 | Pending |
| `FUIOU` | 富友支付 | Pending |
| `HEEPAY` | 汇付天下 | Pending |

## Code Conventions

- Chinese comments for business logic documentation
- Interface-first design with `New<Type>()` constructors
- Table-driven tests with `t.Run()` subtests
- Error wrapping: `fmt.Errorf("message: %w", err)`
- Repository naming: `Gorm<Entity>Repository`
