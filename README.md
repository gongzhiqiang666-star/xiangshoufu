# 收享付

代理商分润管理系统 - 支持多级代理商体系、多通道对接、实时分润计算。

## 项目结构

```
xiangshoufu/
├── server/                      # 后端服务 (Go)
│   ├── cmd/server/main.go       # 服务入口
│   ├── internal/                # 内部模块
│   │   ├── handler/             # HTTP处理器
│   │   ├── service/             # 业务逻辑层
│   │   ├── repository/          # 数据仓库层
│   │   ├── channel/             # 支付通道适配器
│   │   ├── middleware/          # 中间件
│   │   ├── models/              # 数据模型
│   │   ├── async/               # 异步处理
│   │   ├── cache/               # 缓存层
│   │   └── jobs/                # 定时任务
│   ├── pkg/                     # 公共包
│   ├── migrations/              # 数据库迁移脚本
│   ├── scripts/                 # 工具脚本
│   ├── swagger/                 # Swagger API文档
│   └── bin/                     # 编译产物
├── web/                         # PC端管理系统 (Vue 3)
│   ├── src/
│   │   ├── api/                 # API接口
│   │   ├── components/          # 组件
│   │   ├── views/               # 页面
│   │   ├── stores/              # 状态管理
│   │   └── router/              # 路由
│   └── package.json
├── mobileapp/                   # 移动端APP (Flutter)
│   ├── lib/
│   └── pubspec.yaml
├── docs/                        # 项目文档
│   ├── api/                     # API接口文档
│   ├── design/                  # 设计文档
│   └── plans/                   # 开发计划
├── CLAUDE.md
└── README.md
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 15+
- Node.js 18+
- Flutter 3.x (移动端开发)

### 2. 数据库迁移

```bash
cd server

# 按顺序执行迁移脚本
psql -d xiangshoufu -f migrations/000_create_core_tables.sql
psql -d xiangshoufu -f migrations/001_create_raw_callback_logs.sql
psql -d xiangshoufu -f migrations/002_create_device_fees.sql
psql -d xiangshoufu -f migrations/003_create_rate_changes.sql
psql -d xiangshoufu -f migrations/004_create_messages.sql
psql -d xiangshoufu -f migrations/005_alter_transactions_and_profits.sql
```

### 3. 配置环境变量

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/xiangshoufu"
export HENGXINTONG_PUBLIC_KEY="..."  # 恒信通RSA公钥
export ALERT_WEBHOOK_URL="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"
```

### 4. 启动后端服务

```bash
cd server
go run cmd/server/main.go
```

### 5. 启动PC端前台

```bash
cd web
npm install
npm run dev
```

### 6. 访问系统

- PC端管理系统: http://localhost:5173
- 后端API: http://localhost:8080
- Swagger文档: http://localhost:8080/swagger/index.html

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |

## 支持的支付通道

| 通道编码 | 通道名称 | 状态 |
|---------|---------|------|
| HENGXINTONG | 恒信通 | ✅ 已实现 |
| LAKALA | 拉卡拉 | 🚧 待实现 |
| YEAHKA | 乐刷 | 🚧 待实现 |
| SUIXINGFU | 随行付 | 🚧 待实现 |
| LIANLIAN | 连连支付 | 🚧 待实现 |
| SANDPAY | 杉德支付 | 🚧 待实现 |
| FUIOU | 富友支付 | 🚧 待实现 |
| HEEPAY | 汇付天下 | 🚧 待实现 |

## 添加新通道

1. 在 `server/internal/channel/` 下创建新通道目录
2. 实现 `ChannelAdapter` 接口
3. 在 `server/cmd/server/main.go` 中注册适配器

```go
// 示例：添加拉卡拉适配器
package lakala

type Adapter struct {
    // ...
}

func (a *Adapter) GetChannelCode() string {
    return "LAKALA"
}

func (a *Adapter) VerifySign(rawBody []byte) (bool, error) {
    // 实现签名验证
}

func (a *Adapter) ParseTransaction(rawBody []byte) (*channel.UnifiedTransaction, error) {
    // 实现交易解析
}
```

## 升级路径

| 当前（单机） | 升级后（分布式） | 触发条件 |
|-------------|-----------------|----------|
| MemoryQueue | → Kafka | QPS > 2000 |
| LocalCache | → Redis | 多实例部署 |
| 单机PostgreSQL | → 主从复制 | 写入QPS > 500 |

## 监控告警

| 指标 | 阈值 | 告警方式 |
|------|------|----------|
| 通道成功率 | < 95% | 企业微信 |
| 队列积压 | > 500条 | 企业微信 |
| 平均延迟 | > 500ms | 企业微信 |
| 签名失败 | > 10次/分钟 | 企业微信 |

## 文档

- [收享付实施方案](docs/design/收享付实施方案.md)
- [业务逻辑梳理](docs/design/业务逻辑梳理.md)
- [PC端管理功能详细设计](docs/design/PC端管理功能详细设计.md)
- [APP详细设计文档](docs/design/APP详细设计文档.md)
- [恒信通API文档](docs/api/恒信通-20251222-API推送接口文档.md)

## License

Private - All rights reserved
