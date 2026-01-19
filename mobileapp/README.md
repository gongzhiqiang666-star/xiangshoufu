# 享收付 APP - 代理商分润管理系统

## 项目简介

享收付是一款面向支付代理商的分润管理APP，支持iOS、Android和HarmonyOS三大平台。

## 技术栈

- **框架**: Flutter 3.x
- **语言**: Dart 3.x
- **UI组件库**: Bruno（贝壳找房）+ TDesign（腾讯）
- **状态管理**: Riverpod 2.0
- **路由**: go_router
- **网络**: Dio + Retrofit

## 目录结构

```
mobileapp/
├── lib/
│   ├── main.dart                 # 入口
│   ├── core/                     # 核心模块
│   │   ├── theme/                # 设计系统
│   │   ├── network/              # 网络层
│   │   ├── storage/              # 存储层
│   │   └── utils/                # 工具类
│   ├── shared/                   # 共享组件
│   ├── features/                 # 功能模块
│   │   ├── auth/                 # 认证
│   │   ├── home/                 # 首页
│   │   ├── terminal/             # 终端管理
│   │   ├── merchant/             # 商户管理
│   │   ├── wallet/               # 钱包
│   │   └── ...                   # 其他模块
│   └── router/                   # 路由
├── assets/                       # 资源
└── pubspec.yaml                  # 依赖
```

## 快速开始

### 1. 环境要求

- Flutter 3.16+
- Dart 3.2+
- Xcode 15+ (iOS)
- Android Studio / VS Code

### 2. 安装Flutter

```bash
# macOS - 使用Homebrew
brew install flutter

# 或从官网下载
# https://flutter.dev/docs/get-started/install
```

### 3. 检查环境

```bash
flutter doctor
```

### 4. 安装依赖

```bash
cd mobileapp
flutter pub get
```

### 5. 运行项目

```bash
# iOS模拟器
flutter run -d iPhone

# Android模拟器
flutter run -d android

# 指定设备
flutter devices  # 查看可用设备
flutter run -d <device_id>
```

### 6. 构建发布包

```bash
# Android APK
flutter build apk --release

# Android App Bundle
flutter build appbundle --release

# iOS
flutter build ios --release
```

## 设计系统

### 颜色

| 名称 | 色值 | 用途 |
|------|------|------|
| Primary | `#2563EB` | 主色 |
| Success | `#10B981` | 成功 |
| Warning | `#F59E0B` | 警告 |
| Danger | `#EF4444` | 错误 |

### 分润类型色

| 类型 | 色值 |
|------|------|
| 交易分润 | `#2563EB` |
| 押金返现 | `#10B981` |
| 流量返现 | `#F59E0B` |
| 激活奖励 | `#8B5CF6` |

## 核心功能

- [x] 首页数据总览
- [x] 终端管理（划拨/回拨）
- [x] 钱包管理
- [ ] 商户管理
- [ ] 收益统计
- [ ] 数据分析
- [ ] 代理拓展
- [ ] 消息通知
- [ ] 营销海报

## 相关文档

- [APP详细设计文档](../design/APP详细设计文档.md)
- [APP设计稿](../design/APP设计稿.md)
- [业务逻辑梳理](../design/业务逻辑梳理.md)

## 开发规范

### 代码风格

```bash
# 格式化代码
flutter format .

# 分析代码
flutter analyze
```

### 提交规范

```
feat: 新功能
fix: Bug修复
docs: 文档更新
style: 代码格式
refactor: 重构
test: 测试
chore: 构建/工具
```

## 常用命令

```bash
# 清理缓存
flutter clean

# 更新依赖
flutter pub upgrade

# 生成代码（如JSON序列化）
flutter pub run build_runner build

# 运行测试
flutter test

# 查看依赖树
flutter pub deps
```

## License

Private - All rights reserved
