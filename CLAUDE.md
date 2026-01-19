# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

收享付 (ShouXiangFu) - An agent profit-sharing management system that processes payment channel callbacks from multiple payment providers. The system handles transaction callbacks, calculates profit sharing across agent hierarchies, and manages wallets.

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
│   ├── design/                  # 设计文档
│   └── plans/                   # 开发计划
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

# Run specific channel adapter tests
go test ./internal/channel/hengxintong/...

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
```

### 移动端 (mobileapp/)

```bash
cd mobileapp

# Get dependencies
flutter pub get

# Run app
flutter run
```

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
