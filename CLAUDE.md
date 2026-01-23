# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

收享付 (ShouXiangFu) - An agent profit-sharing management system that processes payment channel callbacks from multiple payment providers. The system handles transaction callbacks, calculates profit sharing across agent hierarchies, and manages wallets.

## 重要
每次代码改造后，需要将改造的内容同步给 docs/design/业务逻辑梳理.md 文件中， 不要同步详细设计，而是改造后的的业务流程，维护到对应模块下面

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

## 开发规范 - 测试驱动开发 (TDD)

### 核心原则

| 之前 | 之后 |
|------|------|
| 编译通过就部署 | 测试通过才部署 |
| 手动测试验证 | 一个命令自动验证 |
| 改代码担心破坏旧功能 | 测试会告诉你哪里坏了 |
| Bug在生产环境发现 | Bug在开发时就被发现 |

### 每次开发必须遵循的步骤

1. **开发前**：先写测试用例（或让Claude生成）
   - 正常情况测试
   - 边界情况测试
   - 错误处理测试

2. **开发中**：确保测试通过
   ```bash
   go test ./internal/service/... -v
   ```

3. **开发后**：运行全量测试
   ```bash
   go test ./... -v
   ```

4. **查看覆盖率**：
   ```bash
   go test ./internal/service/... -cover
   ```

### Claude Code 开发指令模板

每次让Claude开发功能时，使用以下模板：

```
请帮我实现[功能]，要求：
1. 先写单元测试，覆盖：正常流程、边界情况、错误处理
2. 再写实现代码
3. 确保 go test 通过
```

### 测试文件规范

| 源文件 | 测试文件 |
|--------|----------|
| `xxx_service.go` | `xxx_service_test.go` |
| `xxx_handler.go` | `xxx_handler_test.go` |
| `adapter.go` | `adapter_test.go` |

### Mock 规范

- 使用接口进行依赖注入，便于 mock
- Mock 实现放在 `_test.go` 文件中
- 命名规范：`Mock<Interface>` 如 `MockTransactionRepository`

---

## 前端测试规范 - PC端 (Vue 3 + Vitest)

### 测试命令

```bash
cd web

# 运行测试（监听模式）
npm run test

# 单次运行测试
npm run test:run

# 查看覆盖率报告
npm run test:coverage
```

### 测试文件规范

| 源文件 | 测试文件 |
|--------|----------|
| `src/utils/format.ts` | `src/utils/__tests__/format.test.ts` |
| `src/stores/user.ts` | `src/stores/__tests__/user.test.ts` |
| `src/components/StatCard.vue` | `src/components/__tests__/StatCard.test.ts` |
| `src/api/auth.ts` | `src/api/__tests__/auth.test.ts` |

### 测试分类与覆盖要求

| 测试类型 | 工具 | 占比 | 说明 |
|---------|------|------|------|
| 单元测试 | Vitest | 60% | 工具函数、纯逻辑 |
| Store测试 | Vitest + Pinia | 20% | 状态管理 |
| 组件测试 | @vue/test-utils | 20% | UI组件渲染与交互 |

### 必须测试的4种场景

```typescript
describe('formatAmount', () => {
  // ✅ 1. 正常流程 (Happy Path)
  it('should format 100 cents to "1.00"', () => {
    expect(formatAmount(100)).toBe('1.00')
  })

  // ✅ 2. 边界情况 (Edge Cases)
  it('should handle zero', () => {
    expect(formatAmount(0)).toBe('0.00')
  })

  // ✅ 3. 错误处理 (Error Handling)
  it('should handle negative values', () => {
    expect(formatAmount(-100)).toBe('-1.00')
  })

  // ✅ 4. 特殊输入 (Special Inputs)
  it('should handle large numbers', () => {
    expect(formatAmount(10000000)).toBe('100,000.00')
  })
})
```

### Store 测试模板

```typescript
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '@/stores/user'

describe('useUserStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with null user', () => {
    const store = useUserStore()
    expect(store.userInfo).toBeNull()
    expect(store.isLoggedIn).toBe(false)
  })

  it('should update state after login', async () => {
    const store = useUserStore()
    // Mock API and test login flow
  })
})
```

### 组件测试模板

```typescript
import { mount } from '@vue/test-utils'
import StatCard from '@/components/Common/StatCard.vue'

describe('StatCard.vue', () => {
  it('should render title and value', () => {
    const wrapper = mount(StatCard, {
      props: { title: '今日交易', value: '1,234' }
    })
    expect(wrapper.text()).toContain('今日交易')
    expect(wrapper.text()).toContain('1,234')
  })

  it('should emit click event', async () => {
    const wrapper = mount(StatCard, {
      props: { title: 'Test', value: '0' }
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })
})
```

---

## 前端测试规范 - APP端 (Flutter)

### 测试命令

```bash
cd mobileapp

# 运行所有测试
flutter test

# 运行指定目录测试
flutter test test/utils/

# 查看覆盖率
flutter test --coverage

# 生成覆盖率报告
genhtml coverage/lcov.info -o coverage/html
```

### 测试文件规范

| 源文件 | 测试文件 |
|--------|----------|
| `lib/utils/format.dart` | `test/utils/format_test.dart` |
| `lib/features/home/domain/home_model.dart` | `test/features/home/home_model_test.dart` |
| `lib/features/home/presentation/providers/home_provider.dart` | `test/providers/home_provider_test.dart` |
| `lib/widgets/stat_card.dart` | `test/widgets/stat_card_test.dart` |

### 测试分类

| 测试类型 | 工具 | 说明 |
|---------|------|------|
| 单元测试 | flutter_test | 工具函数、Model、纯逻辑 |
| Provider测试 | flutter_test + mocktail | 状态管理测试 |
| Widget测试 | flutter_test | UI组件渲染与交互 |

### 单元测试模板

```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/utils/format.dart';

void main() {
  group('formatAmount', () {
    // ✅ 正常流程
    test('should format 100 cents to "1.00"', () {
      expect(formatAmount(100), equals('1.00'));
    });

    // ✅ 边界情况
    test('should handle zero', () {
      expect(formatAmount(0), equals('0.00'));
    });

    // ✅ 错误处理
    test('should handle negative values', () {
      expect(formatAmount(-100), equals('-1.00'));
    });
  });
}
```

### Provider 测试模板

```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mocktail/mocktail.dart';

class MockHomeApi extends Mock implements HomeApi {}

void main() {
  late MockHomeApi mockApi;
  late ProviderContainer container;

  setUp(() {
    mockApi = MockHomeApi();
    container = ProviderContainer(overrides: [
      homeApiProvider.overrideWithValue(mockApi),
    ]);
  });

  tearDown(() {
    container.dispose();
  });

  test('should load dashboard data successfully', () async {
    when(() => mockApi.getDashboard()).thenAnswer(
      (_) async => DashboardData(todayAmount: 1000),
    );

    final result = await container.read(dashboardProvider.future);
    expect(result.todayAmount, equals(1000));
  });
}
```

### Widget 测试模板

```dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/widgets/stat_card.dart';

void main() {
  testWidgets('StatCard displays title and value', (tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: StatCard(title: '今日交易', value: '1,234'),
        ),
      ),
    );

    expect(find.text('今日交易'), findsOneWidget);
    expect(find.text('1,234'), findsOneWidget);
  });

  testWidgets('StatCard responds to tap', (tester) async {
    bool tapped = false;
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: StatCard(
            title: 'Test',
            value: '0',
            onTap: () => tapped = true,
          ),
        ),
      ),
    );

    await tester.tap(find.byType(StatCard));
    expect(tapped, isTrue);
  });
}
```

---

## Claude Code 前端开发指令模板

每次让Claude开发前端功能时，使用以下模板：

### PC端 (Vue 3)

```
请帮我实现[功能]，要求：
1. 先写单元测试，覆盖：正常流程、边界情况、错误处理
2. 再写实现代码
3. 确保 npm run test:run 通过
```

### APP端 (Flutter)

```
请帮我实现[功能]，要求：
1. 先写单元测试，覆盖：正常流程、边界情况、错误处理
2. 再写实现代码
3. 确保 flutter test 通过
```

---

## 全栈开发检查清单

每次功能开发完成后，必须执行以下检查：

```bash
# 1. 后端测试
cd server && go test ./... -v

# 2. PC端测试
cd web && npm run test:run

# 3. APP端测试
cd mobileapp && flutter test

# 4. 全部通过后才能提交
git add . && git commit -m "feat: xxx"
```

### 测试覆盖率目标

| 模块 | 最低覆盖率 |
|------|-----------|
| 后端 Service 层 | 80% |
| 后端 Handler 层 | 60% |
| PC端 工具函数 | 90% |
| PC端 Store | 70% |
| APP端 工具函数 | 90% |
| APP端 Provider | 70% |
