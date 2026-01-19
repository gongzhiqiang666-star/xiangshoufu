# 代理商分润管理系统 - 完整实现计划

## 项目概述

开发一款支付代理商分润管理系统，支持多级代理商体系、多通道对接、实时分润计算、多种奖励机制。

### 整体操作流程
```
开代理 → 设置政策（结算价）→ 下发机具 → 装机（流量卡-押金-费率）→ 产生交易 → 分润/押金奖励/流量费奖励/活动奖励
```

---

## 一、系统模块清单

### 核心业务模块

| 模块 | 功能说明 | 优先级 |
|------|----------|--------|
| **通道管理** | 对接8家第三方支付公司（如拉卡拉），每个通道独立配置 | P0 |
| **代理管理** | 代理商注册、树结构、政策模板分配 | P0 |
| **政策引擎** | 分润规则、奖励规则、返现规则配置 | P0 |
| **终端管理** | 机具下发/回拨、费率设置、押金/流量卡设置 | P0 |
| **商户管理** | 商户入网、费率修改、直营/团队商户 | P0 |
| **交易中心** | 接收通道交易流水、数据校验、触发分润 | P0 |
| **分润计算** | 4种分润类型计算（交易分润、奖励、押金返现、流量返现） | P0 |
| **钱包系统** | 正常钱包、充值钱包、沉淀钱包 | P0 |
| **提现管理** | 提现申请、审核、税筹通道打款 | P1 |
| **代扣管理** | 上级扣款、伙伴代扣 | P1 |
| **数据分析** | 交易查询、商户分类、代理排名 | P1 |
| **收益统计** | 收益明细、收益趋势 | P1 |
| **营销模块** | 营销海报、推广码、滚动图 | P2 |
| **消息通知** | 分润提醒、注册提醒、消费提醒 | P2 |

### 模块关系图
```
                    ┌─────────────┐
                    │   通道管理   │ ← 对接8家支付公司
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  代理管理    │    │  政策引擎    │    │  终端管理    │
│  (树结构)    │◄──►│ (分润/奖励)  │◄──►│ (机具/SN)   │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                  │
       └──────────────────┼──────────────────┘
                          ▼
                  ┌─────────────┐
                  │  商户管理    │
                  └──────┬──────┘
                         │
                         ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  交易中心    │───►│  分润计算    │───►│  钱包系统    │
│ (Kafka消费)  │    │ (4种类型)   │    │ (3种钱包)   │
└─────────────┘    └─────────────┘    └──────┬──────┘
                                             │
                                             ▼
                                      ┌─────────────┐
                                      │  提现管理    │
                                      │ (税筹通道)   │
                                      └─────────────┘
```

---

## 二、技术栈选型

| 层级 | 技术选择 | 说明 |
|------|----------|------|
| **移动端APP** | Flutter 3.x | 跨平台，iOS/Android |
| **管理后台** | Vue 3 + Element Plus | 企业级后台 |
| **后端框架** | Go 1.21 + Gin + GORM | 高并发处理 |
| **主数据库** | PostgreSQL 15 | 支持递归CTE查询 |
| **缓存** | Redis 7 Cluster | 代理商树缓存 |
| **消息队列** | Kafka | 交易流水削峰 |
| **API网关** | Kong | 限流、鉴权 |

---

## 三、核心数据库设计

### 3.1 通道相关表

```sql
-- 支付通道表
CREATE TABLE channels (
    id              BIGSERIAL PRIMARY KEY,
    channel_code    VARCHAR(32) NOT NULL UNIQUE,    -- 通道编码 如 LAKALA
    channel_name    VARCHAR(100) NOT NULL,          -- 通道名称

    -- 通道成本费率
    credit_rate     DECIMAL(10,4),                  -- 贷记卡成本
    debit_rate      DECIMAL(10,4),                  -- 借记卡成本
    debit_cap       DECIMAL(10,2),                  -- 借记卡封顶
    unionpay_rate   DECIMAL(10,4),                  -- 银联云闪付成本
    wechat_rate     DECIMAL(10,4),                  -- 微信扫码成本
    alipay_rate     DECIMAL(10,4),                  -- 支付宝扫码成本

    -- 接口配置
    api_url         VARCHAR(255),
    api_key         VARCHAR(255),
    api_secret      VARCHAR(255),

    status          SMALLINT DEFAULT 1,             -- 1启用 2禁用
    is_visible      BOOLEAN DEFAULT TRUE,           -- 是否APP可见
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 代理商-通道关联表
CREATE TABLE agent_channels (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    is_enabled      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, channel_id)
);
```

### 3.2 政策模板表

```sql
-- 政策模板主表
CREATE TABLE policy_templates (
    id              BIGSERIAL PRIMARY KEY,
    template_name   VARCHAR(100) NOT NULL,
    channel_id      BIGINT NOT NULL,                -- 所属通道
    is_default      BOOLEAN DEFAULT FALSE,          -- 是否默认模板

    -- 分润设置 (结算价)
    credit_rate     DECIMAL(10,4),                  -- 贷记卡结算价
    debit_rate      DECIMAL(10,4),                  -- 借记卡结算价
    debit_cap       DECIMAL(10,2),                  -- 借记卡封顶
    unionpay_rate   DECIMAL(10,4),                  -- 银联云闪付
    wechat_rate     DECIMAL(10,4),                  -- 微信扫码
    alipay_rate     DECIMAL(10,4),                  -- 支付宝扫码

    -- T+0秒到账设置（每笔固定金额）
    t0_fee_type     SMALLINT DEFAULT 0,             -- 0:不开通 1:加1元 2:加2元 3:加3元
    t0_fee_amount   DECIMAL(10,2) DEFAULT 0,        -- 秒到费用（元/笔）

    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 结算价阶梯调整表 (费率分阶段)
CREATE TABLE policy_rate_stages (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    stage_type      SMALLINT NOT NULL,              -- 1:按商户入网时间 2:按代理商入网时间
    day_start       INT NOT NULL,                   -- 开始天数
    day_end         INT,                            -- 结束天数 NULL表示无限
    rate_adjust     DECIMAL(10,4) NOT NULL,         -- 费率调整值
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 激活奖励规则表
CREATE TABLE policy_activation_rewards (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    day_start       INT NOT NULL,                   -- 激活后开始天数
    day_end         INT NOT NULL,                   -- 激活后结束天数
    trade_amount    DECIMAL(15,2) NOT NULL,         -- 达标交易额
    reward_amount   DECIMAL(10,2) NOT NULL,         -- 奖励金额
    effective_date  DATE,                           -- 生效日期(针对新商户)
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 押金返现规则表
CREATE TABLE policy_deposit_cashbacks (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    deposit_amount  DECIMAL(10,2) NOT NULL,         -- 押金金额 99/199/299
    cashback_amount DECIMAL(10,2) NOT NULL,         -- 返现金额
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 流量费返现规则表
CREATE TABLE policy_sim_cashbacks (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    sim_type        SMALLINT NOT NULL,              -- 1:首次 2:二次 3:N次
    sim_fee         DECIMAL(10,2) NOT NULL,         -- 流量费金额
    cashback_amount DECIMAL(10,2) NOT NULL,         -- 返现金额
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 3.3 代理商表 (增强版)

```sql
-- 代理商表
CREATE TABLE agents (
    id              BIGSERIAL PRIMARY KEY,
    agent_no        VARCHAR(32) NOT NULL UNIQUE,    -- 代理商编号
    agent_name      VARCHAR(100) NOT NULL,

    -- 树结构
    parent_id       BIGINT REFERENCES agents(id),
    path            VARCHAR(500) DEFAULT '',        -- 物化路径 /1/5/12/
    level           INT DEFAULT 1,

    -- 费率 (每个通道独立配置，这里是默认)
    default_rate    DECIMAL(10,4),

    -- 联系信息
    contact_name    VARCHAR(50),
    contact_phone   VARCHAR(20) NOT NULL,
    id_card_no      VARCHAR(18),                    -- 脱敏存储

    -- 结算信息
    bank_name       VARCHAR(100),
    bank_account    VARCHAR(30),
    bank_card_no    VARCHAR(25),

    -- 推广码
    invite_code     VARCHAR(20) UNIQUE,             -- 邀请码(可自定义靓号)
    qr_code_url     VARCHAR(255),                   -- 二维码URL

    -- 状态
    status          SMALLINT DEFAULT 1,
    register_time   TIMESTAMP DEFAULT NOW(),

    -- 统计
    direct_agent_count   INT DEFAULT 0,             -- 直属代理数
    direct_merchant_count INT DEFAULT 0,            -- 直营商户数
    team_agent_count     INT DEFAULT 0,             -- 团队代理数
    team_merchant_count  INT DEFAULT 0,             -- 团队商户数

    created_at      TIMESTAMP DEFAULT NOW()
);

-- 代理商-政策模板关联表
CREATE TABLE agent_policies (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    template_id     BIGINT NOT NULL,

    -- 可覆盖模板的费率
    credit_rate     DECIMAL(10,4),
    debit_rate      DECIMAL(10,4),

    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, channel_id)
);
```

### 3.4 终端(机具)表

```sql
-- 终端/机具表
CREATE TABLE terminals (
    id              BIGSERIAL PRIMARY KEY,
    sn              VARCHAR(50) NOT NULL UNIQUE,    -- 机具SN号
    terminal_no     VARCHAR(20),                    -- 终端号
    channel_id      BIGINT NOT NULL,                -- 所属通道

    -- 归属关系
    owner_agent_id  BIGINT,                         -- 当前持有代理商
    merchant_id     BIGINT,                         -- 绑定商户

    -- 机具设置
    credit_rate     DECIMAL(10,4),                  -- 贷记卡费率
    debit_rate      DECIMAL(10,4),                  -- 借记卡费率

    -- T+0秒到账设置
    t0_enabled      BOOLEAN DEFAULT FALSE,          -- 是否开通T+0
    t0_fee_type     SMALLINT DEFAULT 0,             -- 0:不开通 1:加1元 2:加2元 3:加3元
    t0_fee_amount   DECIMAL(10,2) DEFAULT 0,        -- 秒到费用（元/笔）

    -- 押金设置
    deposit_type    SMALLINT DEFAULT 0,             -- 0无押金 1:99 2:199 3:299
    deposit_amount  DECIMAL(10,2) DEFAULT 0,

    -- 流量卡设置
    sim_first_fee   DECIMAL(10,2),                  -- 首次扣费金额
    sim_next_fee    DECIMAL(10,2),                  -- 非首次扣费金额
    sim_interval    INT,                            -- 非首次间隔天数

    -- 状态
    status          SMALLINT DEFAULT 0,             -- 0库存 1已下发 2已绑定 3已激活

    -- 时间
    dispatch_time   TIMESTAMP,                      -- 下发时间
    bind_time       TIMESTAMP,                      -- 绑定时间
    activate_time   TIMESTAMP,                      -- 激活时间
    first_trade_time TIMESTAMP,                     -- 首刷时间

    created_at      TIMESTAMP DEFAULT NOW()
);

-- 终端流转记录表
CREATE TABLE terminal_logs (
    id              BIGSERIAL PRIMARY KEY,
    terminal_id     BIGINT NOT NULL,
    action          SMALLINT NOT NULL,              -- 1下发 2回拨 3绑定 4解绑 5激活
    from_agent_id   BIGINT,
    to_agent_id     BIGINT,
    operator_id     BIGINT NOT NULL,
    remark          VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 3.5 交易与分润表

```sql
-- 交易流水表 (按月分区)
CREATE TABLE transactions (
    id              BIGSERIAL,
    trade_no        VARCHAR(64) NOT NULL,           -- 通道交易号
    order_no        VARCHAR(64) NOT NULL UNIQUE,    -- 系统订单号
    channel_id      BIGINT NOT NULL,

    -- 关联
    terminal_sn     VARCHAR(50) NOT NULL,
    merchant_id     BIGINT NOT NULL,
    agent_id        BIGINT NOT NULL,                -- 直属代理商

    -- 交易信息
    trade_type      SMALLINT NOT NULL,              -- 1消费 2撤销 3退货 4服务费 5押金 6流量费
    pay_type        SMALLINT NOT NULL,              -- 1刷卡 2微信 3支付宝 4云闪付
    card_type       SMALLINT,                       -- 1借记卡 2贷记卡

    amount          DECIMAL(15,2) NOT NULL,         -- 交易金额
    fee             DECIMAL(15,2),                  -- 手续费
    rate            DECIMAL(10,4),                  -- 交易费率

    -- 分润状态
    profit_status   SMALLINT DEFAULT 0,             -- 0待计算 1已计算 2失败

    trade_time      TIMESTAMP NOT NULL,
    received_at     TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (id, trade_time)
) PARTITION BY RANGE (trade_time);

-- 分润明细表 (4种类型)
CREATE TABLE profit_records (
    id              BIGSERIAL PRIMARY KEY,
    transaction_id  BIGINT NOT NULL,
    order_no        VARCHAR(64) NOT NULL,

    agent_id        BIGINT NOT NULL,                -- 获得分润的代理商
    profit_type     SMALLINT NOT NULL,              -- 1交易分润 2激活奖励 3押金返现 4流量返现

    -- 计算明细
    trade_amount    DECIMAL(15,2) NOT NULL,
    self_rate       DECIMAL(10,4),
    lower_rate      DECIMAL(10,4),
    rate_diff       DECIMAL(10,4),
    profit_amount   DECIMAL(15,4) NOT NULL,

    -- 来源
    source_merchant_id BIGINT NOT NULL,
    source_agent_id    BIGINT NOT NULL,
    channel_id         BIGINT NOT NULL,

    -- 入账状态
    wallet_type     SMALLINT NOT NULL,              -- 1分润钱包 2服务费钱包 3奖励钱包
    wallet_status   SMALLINT DEFAULT 0,             -- 0待入账 1已入账

    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 3.6 钱包表 (3种钱包)

```sql
-- 钱包表 (每个代理商每个通道每种类型一个钱包)
CREATE TABLE wallets (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,              -- 1分润钱包 2服务费钱包 3奖励钱包

    balance         DECIMAL(15,4) DEFAULT 0,        -- 可用余额
    frozen_amount   DECIMAL(15,4) DEFAULT 0,        -- 冻结金额
    total_income    DECIMAL(15,4) DEFAULT 0,        -- 累计收入
    total_withdraw  DECIMAL(15,4) DEFAULT 0,        -- 累计提现

    -- 提现门槛
    withdraw_threshold DECIMAL(10,2) DEFAULT 100,

    version         INT DEFAULT 0,                  -- 乐观锁
    updated_at      TIMESTAMP DEFAULT NOW(),

    UNIQUE(agent_id, channel_id, wallet_type)
);

-- 充值钱包 (上级给下级的额外奖励)
CREATE TABLE recharge_wallets (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL UNIQUE,
    balance         DECIMAL(15,4) DEFAULT 0,
    total_recharge  DECIMAL(15,4) DEFAULT 0,        -- 累计充值
    total_paid      DECIMAL(15,4) DEFAULT 0,        -- 累计发放
    version         INT DEFAULT 0,
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 沉淀钱包 (可使用下级未提现的比例)
CREATE TABLE deposit_wallets (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL UNIQUE,
    available_ratio DECIMAL(5,2) DEFAULT 30,        -- 可使用比例
    available_amount DECIMAL(15,4) DEFAULT 0,       -- 可使用金额
    used_amount     DECIMAL(15,4) DEFAULT 0,        -- 已使用金额
    is_visible      BOOLEAN DEFAULT FALSE,          -- 是否可见
    version         INT DEFAULT 0,
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 钱包流水表
CREATE TABLE wallet_logs (
    id              BIGSERIAL PRIMARY KEY,
    wallet_id       BIGINT NOT NULL,
    agent_id        BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,

    log_type        SMALLINT NOT NULL,              -- 1分润入账 2提现冻结 3提现成功 4提现退回 5调账 6代扣
    amount          DECIMAL(15,4) NOT NULL,
    balance_before  DECIMAL(15,4) NOT NULL,
    balance_after   DECIMAL(15,4) NOT NULL,

    ref_type        VARCHAR(20),
    ref_id          BIGINT,
    remark          VARCHAR(500),

    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 3.7 其他表

```sql
-- 提现表
CREATE TABLE withdrawals (
    id              BIGSERIAL PRIMARY KEY,
    withdraw_no     VARCHAR(32) NOT NULL UNIQUE,
    agent_id        BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,
    channel_id      BIGINT NOT NULL,

    amount          DECIMAL(15,2) NOT NULL,
    tax_fee         DECIMAL(15,2) DEFAULT 0,        -- 税费 (如9%)
    service_fee     DECIMAL(10,2) DEFAULT 0,        -- 服务费 (如3元/笔)
    actual_amount   DECIMAL(15,2) NOT NULL,         -- 实际到账

    -- 税筹通道
    tax_channel     VARCHAR(50),

    status          SMALLINT DEFAULT 0,             -- 0待审核 1审核通过 2打款中 3完成 4拒绝

    created_at      TIMESTAMP DEFAULT NOW()
);

-- 代扣表
CREATE TABLE deductions (
    id              BIGSERIAL PRIMARY KEY,
    deduction_no    VARCHAR(32) NOT NULL UNIQUE,
    from_agent_id   BIGINT NOT NULL,                -- 被扣款方
    to_agent_id     BIGINT NOT NULL,                -- 发起方(上级/伙伴)

    deduction_type  SMALLINT NOT NULL,              -- 1上级扣款 2伙伴代扣
    amount          DECIMAL(15,2) NOT NULL,
    periods         INT DEFAULT 1,                  -- 期数
    wallet_source   SMALLINT NOT NULL,              -- 扣款来源钱包类型

    status          SMALLINT DEFAULT 0,             -- 0待确认 1已同意 2已拒绝 3进行中 4完成
    agreement_url   VARCHAR(255),                   -- 协议文件

    created_at      TIMESTAMP DEFAULT NOW()
);

-- 营销海报表
CREATE TABLE marketing_posters (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT,
    title           VARCHAR(100),
    image_url       VARCHAR(255) NOT NULL,
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 消息通知表
CREATE TABLE notifications (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT,                         -- NULL表示全体
    notify_type     SMALLINT NOT NULL,              -- 1分润 2注册 3消费 4节假日 5系统公告
    title           VARCHAR(100) NOT NULL,
    content         TEXT,
    is_read         BOOLEAN DEFAULT FALSE,
    expire_time     TIMESTAMP,                      -- 过期时间(3天)
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、分润计算核心逻辑

### 4.1 四种分润类型

| 类型 | 触发条件 | 计算方式 | 入账钱包 |
|------|----------|----------|----------|
| **交易分润** | 普通交易 | 费率差 × 交易额 | 分润钱包 |
| **激活奖励** | 达标交易额 | 固定金额(分阶段) | 奖励钱包 |
| **押金返现** | 押金扣取成功 | 固定金额 | 服务费钱包 |
| **流量返现** | 流量费扣取成功 | 固定金额(首次/N次) | 服务费钱包 |

### 4.2 分润计算流程

```
收到交易流水 (trade_type判断)
        │
        ├─── trade_type=1(消费) ──→ 计算【交易分润】
        │                              │
        │                              ├─ 查询商户费率
        │                              ├─ 向上追溯代理商链
        │                              ├─ 计算每级费率差分润
        │                              └─ 检查是否达标【激活奖励】
        │
        ├─── trade_type=5(押金) ──→ 计算【押金返现】
        │                              │
        │                              ├─ 查询押金返现规则
        │                              └─ 向上追溯分配返现
        │
        └─── trade_type=6(流量) ──→ 计算【流量返现】
                                       │
                                       ├─ 判断首次/N次
                                       ├─ 查询流量返现规则
                                       └─ 向上追溯分配返现
```

### 4.3 结算价阶梯计算

```go
// 根据商户入网时间或代理商入网时间计算实际结算价
func CalculateSettlementRate(baseRate decimal.Decimal, enrollDays int, stages []RateStage) decimal.Decimal {
    for _, stage := range stages {
        if enrollDays >= stage.DayStart && (stage.DayEnd == nil || enrollDays <= *stage.DayEnd) {
            return baseRate.Add(stage.RateAdjust)
        }
    }
    return baseRate
}
```

---

## 五、项目目录结构

```
profit-sharing-system/
├── backend/                          # 后端服务 (Go)
│   ├── cmd/
│   │   ├── api/                      # API服务
│   │   ├── consumer/                 # Kafka消费者
│   │   └── scheduler/                # 定时任务
│   ├── internal/
│   │   ├── app/
│   │   │   ├── channel/              # 通道管理
│   │   │   ├── agent/                # 代理商管理
│   │   │   ├── policy/               # 政策模板
│   │   │   ├── terminal/             # 终端管理
│   │   │   ├── merchant/             # 商户管理
│   │   │   ├── transaction/          # 交易处理
│   │   │   ├── profit/               # 分润计算 ⭐
│   │   │   ├── wallet/               # 钱包管理
│   │   │   ├── withdrawal/           # 提现管理
│   │   │   ├── deduction/            # 代扣管理
│   │   │   ├── marketing/            # 营销模块
│   │   │   ├── notification/         # 消息通知
│   │   │   ├── report/               # 数据分析
│   │   │   └── user/                 # 用户认证
│   │   └── domain/
│   └── pkg/
│
├── admin-web/                        # 管理后台 (Vue 3)
│
├── mobile-app/                       # 移动端 (Flutter)
│
└── docs/                             # 文档
```

---

## 六、开发阶段划分

### 第一阶段: 基础架构 (第1-3周)
- [ ] 项目框架搭建、数据库设计
- [ ] 通道管理模块
- [ ] 用户认证与权限

### 第二阶段: 核心业务 (第4-10周)
- [ ] 代理商管理（树结构、推广码）
- [ ] 政策模板引擎（4种规则配置）
- [ ] 终端管理（下发、回拨、设置）
- [ ] 商户管理
- [ ] 交易流水接收
- [ ] **分润计算（4种类型）** ⭐

### 第三阶段: 钱包与财务 (第11-14周)
- [ ] 钱包模块（3种钱包）
- [ ] 提现管理（税筹通道对接）
- [ ] 代扣管理
- [ ] 手动调账功能

### 第四阶段: 数据与营销 (第15-17周)
- [ ] 数据分析模块
- [ ] 收益统计
- [ ] 营销海报、推广码
- [ ] 消息通知

### 第五阶段: 前端开发 (第18-22周)
- [ ] Vue管理后台
- [ ] Flutter移动端APP
- [ ] 前后端联调

### 第六阶段: 测试上线 (第23-26周)
- [ ] 功能测试、性能测试
- [ ] 安全审计
- [ ] 部署上线

---

## 七、关键文件清单

| 文件 | 说明 |
|------|------|
| `internal/app/profit/calculator.go` | 分润计算核心 |
| `internal/app/profit/reward_checker.go` | 激活奖励检查 |
| `internal/app/policy/engine.go` | 政策模板引擎 |
| `internal/app/terminal/dispatch.go` | 机具下发/回拨 |
| `internal/app/wallet/multi_wallet.go` | 多钱包管理 |
| `migrations/*.sql` | 数据库迁移脚本 |

---

## 八、通道适配层设计

由于8家支付通道接口差异大，需要设计统一的适配层：

### 8.1 通道适配器接口

```go
// ChannelAdapter 通道适配器接口
type ChannelAdapter interface {
    // 基础信息
    GetChannelCode() string
    GetChannelName() string

    // 商户管理
    RegisterMerchant(req *MerchantRegisterReq) (*MerchantRegisterResp, error)
    UpdateMerchantRate(merchantNo string, rate decimal.Decimal) error

    // 终端管理
    BindTerminal(sn string, merchantNo string) error
    UnbindTerminal(sn string) error
    SetTerminalPolicy(sn string, policy *TerminalPolicy) error

    // 交易查询
    QueryTransaction(tradeNo string) (*TransactionDetail, error)

    // 回调处理
    ParseCallback(data []byte) (*CallbackData, error)
    VerifySign(data []byte, sign string) bool
}

// 每个通道实现该接口
type LakalaAdapter struct { ... }
type XxxPayAdapter struct { ... }
```

### 8.2 通道配置表补充

```sql
-- 通道扩展配置表
CREATE TABLE channel_configs (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    config_key      VARCHAR(100) NOT NULL,
    config_value    TEXT NOT NULL,
    description     VARCHAR(500),
    UNIQUE(channel_id, config_key)
);

-- 示例配置项
-- api_url, api_key, api_secret, callback_url
-- sign_type (MD5/RSA), charset, version
-- merchant_register_url, terminal_bind_url, rate_update_url
```

---

## 九、税筹通道与提现设计

### 9.1 税筹通道规则（按通道区分）

| 扣费类型 | 说明 | 计算公式 |
|----------|------|----------|
| **付款扣** | 充值时扣除税费 | 实际入账 = 充值金额 - 税费 |
| **出款扣** | 提现时扣除税费 | 实际到账 = 提现金额 - 税费 |
| **混合扣** | 税率+固定费用 | 实际到账 = 金额 × (1-税率) - 固定费用 |

### 9.2 税筹通道配置表

```sql
-- 税筹通道表
CREATE TABLE tax_channels (
    id              BIGSERIAL PRIMARY KEY,
    channel_code    VARCHAR(32) NOT NULL UNIQUE,
    channel_name    VARCHAR(100) NOT NULL,

    -- 扣费规则
    fee_type        SMALLINT NOT NULL,              -- 1付款扣 2出款扣
    tax_rate        DECIMAL(5,4) NOT NULL,          -- 税率 如0.09表示9%
    fixed_fee       DECIMAL(10,2) DEFAULT 0,        -- 固定费用 如3元/笔

    -- 接口配置
    api_url         VARCHAR(255),
    api_key         VARCHAR(255),

    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 通道-税筹通道关联 (不同支付通道可能走不同税筹通道)
CREATE TABLE channel_tax_mappings (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,                -- 支付通道
    tax_channel_id  BIGINT NOT NULL,                -- 税筹通道
    wallet_type     SMALLINT NOT NULL,              -- 钱包类型
    UNIQUE(channel_id, wallet_type)
);
```

### 9.3 提现金额计算逻辑

```go
// 计算实际到账金额 (确保代理商拿到的钱按平台规则结算)
func CalculateActualAmount(amount decimal.Decimal, taxChannel *TaxChannel) decimal.Decimal {
    // 税费 = 金额 × 税率
    taxFee := amount.Mul(taxChannel.TaxRate)

    // 固定费用
    fixedFee := taxChannel.FixedFee

    // 实际到账 = 金额 - 税费 - 固定费用
    actualAmount := amount.Sub(taxFee).Sub(fixedFee)

    return actualAmount
}

// 注意：无论付款扣还是出款扣，最终代理商拿到的钱必须一致
```

---

## 十、特殊钱包详细设计

### 10.1 充值钱包

**用途**：上级代理商给下级代理商发放额外奖励

**规则**：
1. 代理商需开启充值钱包功能才能使用
2. 奖励规则在政策模板中配置
3. 上级充值后，满足条件时自动发放给下级

```sql
-- 充值钱包表
CREATE TABLE recharge_wallets (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL UNIQUE,

    is_enabled      BOOLEAN DEFAULT FALSE,          -- 是否开启
    balance         DECIMAL(15,4) DEFAULT 0,        -- 余额
    total_recharge  DECIMAL(15,4) DEFAULT 0,        -- 累计充值
    total_paid      DECIMAL(15,4) DEFAULT 0,        -- 累计发放

    version         INT DEFAULT 0,
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 充值钱包奖励规则表 (在政策模板中配置)
CREATE TABLE policy_recharge_rewards (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,

    reward_type     SMALLINT NOT NULL,              -- 1激活奖励 2交易量奖励 3拉新奖励
    condition_type  SMALLINT NOT NULL,              -- 条件类型
    condition_value DECIMAL(15,2),                  -- 条件值
    reward_amount   DECIMAL(10,2) NOT NULL,         -- 奖励金额

    is_enabled      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 充值钱包流水表
CREATE TABLE recharge_wallet_logs (
    id              BIGSERIAL PRIMARY KEY,
    wallet_id       BIGINT NOT NULL,
    agent_id        BIGINT NOT NULL,

    log_type        SMALLINT NOT NULL,              -- 1充值 2发放奖励 3退回
    amount          DECIMAL(15,4) NOT NULL,
    balance_before  DECIMAL(15,4) NOT NULL,
    balance_after   DECIMAL(15,4) NOT NULL,

    -- 发放相关
    to_agent_id     BIGINT,                         -- 发放给谁
    reward_rule_id  BIGINT,                         -- 奖励规则

    remark          VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 10.2 沉淀钱包

**用途**：
1. 可使用下级未提现余额的一定比例（如30%）
2. 支持向平台借款（前置分润），线下签协议

```sql
-- 沉淀钱包表
CREATE TABLE deposit_wallets (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL UNIQUE,

    is_enabled      BOOLEAN DEFAULT FALSE,          -- 是否开启
    is_visible      BOOLEAN DEFAULT FALSE,          -- 是否可见

    -- 下级余额使用
    available_ratio DECIMAL(5,2) DEFAULT 30,        -- 可使用比例 30%
    subordinate_balance DECIMAL(15,4) DEFAULT 0,    -- 下级未提现总额
    available_amount DECIMAL(15,4) DEFAULT 0,       -- 可使用金额
    used_amount     DECIMAL(15,4) DEFAULT 0,        -- 已使用金额

    -- 平台借款
    loan_limit      DECIMAL(15,2) DEFAULT 0,        -- 借款额度
    loan_balance    DECIMAL(15,4) DEFAULT 0,        -- 借款余额
    loan_rate       DECIMAL(5,4) DEFAULT 0,         -- 借款利率

    version         INT DEFAULT 0,
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 平台借款记录表
CREATE TABLE platform_loans (
    id              BIGSERIAL PRIMARY KEY,
    loan_no         VARCHAR(32) NOT NULL UNIQUE,
    agent_id        BIGINT NOT NULL,

    loan_amount     DECIMAL(15,2) NOT NULL,         -- 借款金额
    loan_rate       DECIMAL(5,4) NOT NULL,          -- 借款利率
    interest_amount DECIMAL(15,2) DEFAULT 0,        -- 利息金额

    repaid_amount   DECIMAL(15,2) DEFAULT 0,        -- 已还金额
    remaining_amount DECIMAL(15,2) NOT NULL,        -- 剩余金额

    status          SMALLINT DEFAULT 1,             -- 1进行中 2已还清 3逾期

    -- 线下协议
    agreement_no    VARCHAR(64),                    -- 协议编号
    agreement_url   VARCHAR(255),                   -- 协议文件URL

    loan_date       DATE NOT NULL,
    due_date        DATE,                           -- 到期日

    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 沉淀钱包流水表
CREATE TABLE deposit_wallet_logs (
    id              BIGSERIAL PRIMARY KEY,
    wallet_id       BIGINT NOT NULL,
    agent_id        BIGINT NOT NULL,

    log_type        SMALLINT NOT NULL,              -- 1使用下级余额 2归还下级余额 3借款 4还款
    amount          DECIMAL(15,4) NOT NULL,

    ref_type        VARCHAR(20),                    -- loan/subordinate
    ref_id          BIGINT,

    remark          VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 10.3 沉淀钱包-下级余额计算

```go
// 定时任务：计算下级未提现余额
func CalculateSubordinateBalance(agentID int64) decimal.Decimal {
    // 获取所有直属下级
    subordinates := GetDirectSubordinates(agentID)

    totalBalance := decimal.Zero
    for _, sub := range subordinates {
        // 汇总下级所有钱包余额
        wallets := GetAgentWallets(sub.ID)
        for _, w := range wallets {
            totalBalance = totalBalance.Add(w.Balance)
        }
    }

    return totalBalance
}

// 可使用金额 = 下级未提现总额 × 可使用比例
func CalculateAvailableAmount(subordinateBalance decimal.Decimal, ratio decimal.Decimal) decimal.Decimal {
    return subordinateBalance.Mul(ratio).Div(decimal.NewFromInt(100))
}
```

### 10.4 特殊钱包注意事项

1. **充值钱包**：
   - 需先充值才能发放奖励
   - 奖励发放走政策模板规则，自动触发
   - 充值钱包余额不足时，奖励挂起待发放

2. **沉淀钱包**：
   - 使用下级余额有风险，若下级提现可能导致超支
   - 需设置预警机制，当沉淀使用接近上限时提醒
   - 平台借款需线下签署协议，系统只做记录
   - 借款通过未来分润自动扣除还款

---

## 十一、数据权限规则

### 11.1 APP端数据可见性

| 数据类型 | 可见范围 | 说明 |
|----------|----------|------|
| 团队数据 | 自己 + 直属子代理汇总 | 子代理再往下的数据看不到 |
| 交易数据 | 直营 + 团队（子代理） | 最近6个月 |
| 商户数据 | 直营商户 + 团队商户 | 团队商户为子代理的直营商户 |
| 收益数据 | 自己的收益明细 | 包括4种分润类型 |

### 11.2 PC后台数据可见性

| 角色 | 可见范围 | 说明 |
|------|----------|------|
| 平台管理员 | 全部数据 | 所有代理商、商户、交易 |
| 代理商 | 自己 + 所有下级 | 包括子代理的子代理（无限层级） |
| 财务人员 | 财务相关数据 | 提现审核、打款记录 |

---

## 十二、待完善事项 (遗留问题)

### 12.1 待确认业务规则

| 问题 | 文档位置 | 状态 |
|------|----------|------|
| 秒到费率的加1、加2、加3规则 | 第53行 | 待确认 |
| 服务商偷数据权限 | 第116行 | 待确认 |
| 退费如何处理分润 | 第20、100行 | 待完善 |

### 12.2 退费处理方案 (待定)

```
退费处理选项：
1. 反向扣除已发放分润 (需要考虑钱包余额不足情况)
2. 不做处理，退费不影响已发分润
3. 单独记录，人工处理
```

### 12.3 补充功能点

#### 1. 货款代扣功能
机具划拨给下级时，可设置货款代扣：
- 划拨时设置代扣金额
- 需要下级接收确认
- 从下级分润钱包中自动扣除

```sql
-- 货款代扣设置表
CREATE TABLE terminal_payment_deductions (
    id              BIGSERIAL PRIMARY KEY,
    terminal_id     BIGINT NOT NULL,
    from_agent_id   BIGINT NOT NULL,         -- 下级(被扣款方)
    to_agent_id     BIGINT NOT NULL,         -- 上级(收款方)
    amount          DECIMAL(10,2) NOT NULL,  -- 代扣金额
    deducted_amount DECIMAL(10,2) DEFAULT 0, -- 已扣金额
    status          SMALLINT DEFAULT 0,      -- 0待接收 1已接收 2已完成 3已拒绝
    created_at      TIMESTAMP DEFAULT NOW()
);
```

#### 2. 商户预警功能
客户登记预警提示（商户管理模块）：
- 30天无交易预警
- 交易量下降预警
- 费率即将调整提醒

```sql
-- 商户预警规则表
CREATE TABLE merchant_alert_rules (
    id              BIGSERIAL PRIMARY KEY,
    alert_type      SMALLINT NOT NULL,       -- 1:无交易 2:交易下降 3:费率调整
    threshold_days  INT,                      -- 天数阈值
    threshold_rate  DECIMAL(5,2),            -- 下降比例阈值
    is_enabled      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 商户预警记录表
CREATE TABLE merchant_alerts (
    id              BIGSERIAL PRIMARY KEY,
    merchant_id     BIGINT NOT NULL,
    agent_id        BIGINT NOT NULL,
    alert_type      SMALLINT NOT NULL,
    alert_message   VARCHAR(500),
    is_handled      BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

#### 3. PC后台跨级操作
终端管理中PC后台特殊权限：
- PC后台可以**跨级下发**机具（APP不可以）
- PC后台可以**跨级回拨**机具（APP不可以，只能一级一级回收）
- 跨级操作需记录日志，便于追溯

#### 4. PC端弹窗提醒
PC端登录后弹窗提醒功能：
- 登录后自动弹出最近重要更新内容
- 支持公告类型：系统更新、政策变更、重要通知等
- 可设置强制阅读（必须点确认才能关闭）

```sql
-- 系统公告表
CREATE TABLE system_announcements (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(200) NOT NULL,
    content         TEXT NOT NULL,
    announce_type   SMALLINT NOT NULL,       -- 1:系统更新 2:政策变更 3:重要通知
    is_force_read   BOOLEAN DEFAULT FALSE,   -- 是否强制阅读
    target_type     SMALLINT DEFAULT 0,      -- 0:全部 1:仅PC 2:仅APP
    start_time      TIMESTAMP,
    end_time        TIMESTAMP,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 公告阅读记录表
CREATE TABLE announcement_reads (
    id              BIGSERIAL PRIMARY KEY,
    announcement_id BIGINT NOT NULL,
    user_id         BIGINT NOT NULL,
    read_at         TIMESTAMP DEFAULT NOW(),
    UNIQUE(announcement_id, user_id)
);
```

#### 5. 政策继承规则
下级代理商注册时：
- **默认继承**上级的政策模板
- **费率需要调整**：每增加一级代理商，利率一般会加一点
- 也可以不加利润直接下放（为了推广，后续通过流量费返现等盈利）
- 取下原则：当没有费率差时，这级代理商不分润，继续向上追溯

#### 6. 通道费率范围管理
- 每个通道提供一个费率范围（如万分53-60）
- 给商户设置费率只能在这个范围内
- 通道表需增加费率上下限字段

```sql
-- 补充通道表字段
ALTER TABLE channels ADD COLUMN rate_min DECIMAL(10,4);  -- 费率下限
ALTER TABLE channels ADD COLUMN rate_max DECIMAL(10,4);  -- 费率上限
```

#### 7. 充值钱包贷款功能
充值钱包除了发放奖励外，还支持贷款功能：
- 代理商可以向下级提供贷款
- 收取利息，走线下协议
- 系统记录贷款信息

```sql
-- 充值钱包贷款记录表
CREATE TABLE recharge_wallet_loans (
    id              BIGSERIAL PRIMARY KEY,
    loan_no         VARCHAR(32) NOT NULL UNIQUE,
    from_agent_id   BIGINT NOT NULL,         -- 贷款方(上级)
    to_agent_id     BIGINT NOT NULL,         -- 借款方(下级)

    loan_amount     DECIMAL(15,2) NOT NULL,
    loan_rate       DECIMAL(5,4) NOT NULL,   -- 利率
    interest_amount DECIMAL(15,2) DEFAULT 0,

    repaid_amount   DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,

    status          SMALLINT DEFAULT 1,      -- 1进行中 2已还清
    agreement_url   VARCHAR(255),            -- 线下协议

    loan_date       DATE NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

#### 8. 营销海报分类
海报需要分类管理，APP端按分类展示：

```sql
-- 海报分类表
CREATE TABLE poster_categories (
    id              BIGSERIAL PRIMARY KEY,
    category_name   VARCHAR(50) NOT NULL,
    sort_order      INT DEFAULT 0,
    is_visible      BOOLEAN DEFAULT TRUE,    -- APP是否可见
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 修改海报表关联分类
ALTER TABLE marketing_posters ADD CONSTRAINT fk_poster_category
    FOREIGN KEY (category_id) REFERENCES poster_categories(id);
```

#### 9. 滚动图(轮播图)
首页轮播图管理：

```sql
-- 滚动图/轮播图表
CREATE TABLE banners (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(100),
    image_url       VARCHAR(255) NOT NULL,
    link_url        VARCHAR(255),            -- 点击跳转链接
    link_type       SMALLINT DEFAULT 0,      -- 0:无 1:内部页面 2:外部链接
    sort_order      INT DEFAULT 0,
    start_time      TIMESTAMP,
    end_time        TIMESTAMP,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

#### 10. 提现门槛差异化
不同通道、不同钱包类型有不同的提现门槛：

```sql
-- 提现门槛配置表
CREATE TABLE withdraw_thresholds (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,       -- 1分润 2服务费 3奖励
    threshold       DECIMAL(10,2) NOT NULL,  -- 提现门槛金额
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(channel_id, wallet_type)
);
```

#### 11. 费率实时修改
代理商可以实时调整商户费率：
- 调用通道接口实时生效
- 需记录费率变更历史
- 个别通道不支持实时修改，需结合接口文档

```sql
-- 商户费率变更记录表
CREATE TABLE merchant_rate_logs (
    id              BIGSERIAL PRIMARY KEY,
    merchant_id     BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    old_rate        DECIMAL(10,4) NOT NULL,
    new_rate        DECIMAL(10,4) NOT NULL,
    operator_id     BIGINT NOT NULL,
    sync_status     SMALLINT DEFAULT 0,      -- 0待同步 1已同步 2同步失败
    sync_time       TIMESTAMP,
    remark          VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

#### 12. 政策下发机具
机具切换或设置变更时，需要通过通道接口推送到机具：
- 流量卡设置、押金设置需下发
- 费率设置需下发
- 使用政策模板可下发默认政策

```go
// 通道适配器接口补充
type ChannelAdapter interface {
    // ... 其他方法

    // 下发机具政策
    PushTerminalPolicy(sn string, policy *TerminalPolicy) error

    // TerminalPolicy 包含: 费率、押金、流量卡设置
}
```

#### 13. 消息有效期机制
- 未读消息在首页显示红点
- 3天内未查看一直活跃
- 超过3天自动过期，不再显示红点
- 已读/过期消息可在消息列表中查看历史

#### 14. "我的信息"模块 (APP端)
代理商个人中心页面，包含：
- 姓名、手机号、服务商编号、入网时间
- 脱敏身份证号
- 结算卡信息 + **更改结算卡按钮**
- 费率成本显示（当前代理商的结算价）
- 唯一邀请码（可自定义靓号）

```sql
-- 结算卡变更记录表
CREATE TABLE agent_bank_card_logs (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    old_bank_name   VARCHAR(100),
    old_card_no     VARCHAR(30),           -- 脱敏存储
    new_bank_name   VARCHAR(100) NOT NULL,
    new_card_no     VARCHAR(30) NOT NULL,  -- 脱敏存储
    status          SMALLINT DEFAULT 0,    -- 0待审核 1已通过 2已拒绝
    auditor_id      BIGINT,
    audit_time      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

#### 15. 终端统计数据 (APP首页)
APP端终端管理模块显示统计数据：
- 终端总数、已激活台数、未激活台数
- 昨日激活数、今日激活数、本月激活数
- 库存区分：未绑定、已绑定、未激活、已激活

```sql
-- 代理商终端统计表 (定时任务更新)
CREATE TABLE agent_terminal_stats (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    stat_date       DATE NOT NULL,

    total_count     INT DEFAULT 0,         -- 终端总数
    activated_count INT DEFAULT 0,         -- 已激活数
    unactivated_count INT DEFAULT 0,       -- 未激活数
    unbound_count   INT DEFAULT 0,         -- 未绑定数
    bound_count     INT DEFAULT 0,         -- 已绑定数

    today_activated INT DEFAULT 0,         -- 今日激活
    month_activated INT DEFAULT 0,         -- 本月激活

    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, channel_id, stat_date)
);
```

#### 16. 商户互查功能
商户管理中支持三种方式互查：
- 通过**商户姓名**查询
- 通过**商户编号**查询
- 通过**机具编号**查询

需要建立相应索引优化查询性能。

#### 17. 数据分析详细需求

**交易查询**（第159-160行）：
- 分类：全部、直营、团队（代理商）
- 时间范围：最近6个月数据

**商户分类统计**（第162-163行）：
| 分类 | 月均交易额 |
|------|------------|
| 忠诚商户 | > 5万 |
| 优质商户 | 3万 ~ 5万 |
| 潜力商户 | 2万 ~ 3万 |
| 一般商户 | 1万 ~ 2万 |
| 低活跃商户 | < 1万 |
| 无交易 | 30天无交易 |

**代理数据分析**（第167行）：
- 各代理交易量排名
- 各代理交易笔数排名
- 各代理终端总数
- 各代理已激活终端数

**客户数据分析**（第168-170行）：
- 客户总交易量：刷卡占比、支付宝占比、微信占比
- 直营商户排名：按本月交易排名
- 商户详情：机具号、刷卡交易金额、扫码交易金额、近七天交易、近半年交易

```sql
-- 商户分类统计表 (按月统计)
CREATE TABLE merchant_category_stats (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    stat_month      VARCHAR(7) NOT NULL,   -- 格式: 2024-01

    loyal_count     INT DEFAULT 0,         -- 忠诚商户数 (>5万)
    quality_count   INT DEFAULT 0,         -- 优质商户数 (3-5万)
    potential_count INT DEFAULT 0,         -- 潜力商户数 (2-3万)
    normal_count    INT DEFAULT 0,         -- 一般商户数 (1-2万)
    low_count       INT DEFAULT 0,         -- 低活跃商户数 (<1万)
    inactive_count  INT DEFAULT 0,         -- 30天无交易

    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, stat_month)
);

-- 代理商业绩排名表 (定时计算)
CREATE TABLE agent_ranking (
    id              BIGSERIAL PRIMARY KEY,
    parent_agent_id BIGINT NOT NULL,       -- 上级代理商(排名范围)
    agent_id        BIGINT NOT NULL,
    stat_date       DATE NOT NULL,

    trade_amount_rank    INT,              -- 交易量排名
    trade_count_rank     INT,              -- 笔数排名
    terminal_count_rank  INT,              -- 终端数排名
    activated_count_rank INT,              -- 激活数排名

    trade_amount    DECIMAL(18,2) DEFAULT 0,
    trade_count     INT DEFAULT 0,
    terminal_count  INT DEFAULT 0,
    activated_count INT DEFAULT 0,

    UNIQUE(parent_agent_id, agent_id, stat_date)
);
```

#### 18. 收益统计详细需求（第175-177行）

**今日收益**：
- 今日收益总额
- 今日收益明细：交易收益、押金返现、流量返现

**收益趋势**：
- 近7天、近30天

**月收益查看**：
- 近六个月、近一年、近两年

```sql
-- 代理商收益日统计表
CREATE TABLE agent_profit_daily (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    stat_date       DATE NOT NULL,

    total_profit    DECIMAL(15,4) DEFAULT 0,    -- 总收益
    trade_profit    DECIMAL(15,4) DEFAULT 0,    -- 交易分润
    deposit_profit  DECIMAL(15,4) DEFAULT 0,    -- 押金返现
    sim_profit      DECIMAL(15,4) DEFAULT 0,    -- 流量返现
    reward_profit   DECIMAL(15,4) DEFAULT 0,    -- 激活奖励

    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, channel_id, stat_date)
);

-- 代理商收益月统计表
CREATE TABLE agent_profit_monthly (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    stat_month      VARCHAR(7) NOT NULL,        -- 格式: 2024-01

    total_profit    DECIMAL(18,4) DEFAULT 0,
    trade_profit    DECIMAL(18,4) DEFAULT 0,
    deposit_profit  DECIMAL(18,4) DEFAULT 0,
    sim_profit      DECIMAL(18,4) DEFAULT 0,
    reward_profit   DECIMAL(18,4) DEFAULT 0,

    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, channel_id, stat_month)
);
```

#### 19. 交易类型扩展字段
根据文档第74行，通道返回的交易数据包含交易类型：
- 服务费刷卡 → 涉及返现
- 流量费刷卡 → 涉及返现
- 押金收取成功 → 涉及返现
- 普通交易 → 涉及分润

交易表的trade_type字段已覆盖，确认映射关系正确。

#### 20. 商户详情字段补充（第154行）
商户详情页需要显示：
- 脱敏手机号、机具编号、激活时间
- 首次流量费时间、首次流量费金额
- 刷卡费率、扫码费率
- 交易额统计：累计、本月总额、本月贷记卡、本月借记卡、本月微信、本月支付宝

```sql
-- 商户交易统计表 (按月)
CREATE TABLE merchant_trade_stats (
    id              BIGSERIAL PRIMARY KEY,
    merchant_id     BIGINT NOT NULL,
    stat_month      VARCHAR(7) NOT NULL,

    total_amount    DECIMAL(18,2) DEFAULT 0,    -- 总交易额
    credit_amount   DECIMAL(18,2) DEFAULT 0,    -- 贷记卡交易额
    debit_amount    DECIMAL(18,2) DEFAULT 0,    -- 借记卡交易额
    wechat_amount   DECIMAL(18,2) DEFAULT 0,    -- 微信交易额
    alipay_amount   DECIMAL(18,2) DEFAULT 0,    -- 支付宝交易额
    unionpay_amount DECIMAL(18,2) DEFAULT 0,    -- 云闪付交易额

    trade_count     INT DEFAULT 0,              -- 交易笔数

    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(merchant_id, stat_month)
);

-- 商户表补充字段
ALTER TABLE merchants ADD COLUMN first_sim_time TIMESTAMP;      -- 首次流量费时间
ALTER TABLE merchants ADD COLUMN first_sim_amount DECIMAL(10,2);-- 首次流量费金额
ALTER TABLE merchants ADD COLUMN total_amount DECIMAL(18,2) DEFAULT 0; -- 累计交易额
```

---

#### 21. 代理拓展模块细节（第28-40行）

**主动注册**：
- 使用拓展码（二维码）展示
- 可保存相册
- 可复制链接给下级代理商

**被动注册**：
- 手动填写下级代理商信息
- 验证通过后成为代理商
- 下级根据二维码下载应用后登录使用

**机构树命名规则**：
- 总部
- 机构（一级代理商）= 字母 + 手机号
- 直属代理商（B端）/ 直营客户（C端）
- 下级代理商

#### 22. 营销海报大小控制（第45行）
- 海报图片需要控制大小，避免影响服务器性能
- 建议：图片压缩、缩略图生成、CDN分发

```sql
-- 海报表补充字段
ALTER TABLE marketing_posters ADD COLUMN file_size INT;           -- 文件大小(KB)
ALTER TABLE marketing_posters ADD COLUMN thumb_url VARCHAR(255);  -- 缩略图URL
ALTER TABLE marketing_posters ADD COLUMN max_size INT DEFAULT 2048; -- 最大允许大小(KB)
```

#### 23. 代理商注册属性（第106行）
代理商注册时需要填写的信息：
- 个人基本信息
- 身份证号（脱敏存储）
- 手机号
- 开户行信息
- 注册时间（系统自动记录）
- 给到的政策（默认继承上级）

---

## 十三、核对清单（业务逻辑文件逐行对照）

| 行号 | 功能点 | 计划覆盖情况 |
|------|--------|--------------|
| 1-21 | 分润逻辑、结算周期、层级追溯 | ✅ 已覆盖 |
| 22-25 | 通道概念、8家支付公司 | ✅ 已覆盖 |
| 28-40 | 代理拓展、主动/被动注册、机构树 | ✅ 已覆盖 |
| 42-45 | 营销海报、分类、大小控制 | ✅ 已覆盖 |
| 47-48 | 滚动图/轮播图 | ✅ 已覆盖 |
| 51-102 | 政策模板、4种分润类型、阶梯费率 | ✅ 已覆盖 |
| 109-118 | 服务商成本设置、调价配置 | ✅ 已覆盖(政策模板) |
| 120-138 | 终端管理、下发/回拨、批量操作 | ✅ 已覆盖 |
| 141-156 | 商户管理、费率修改、脱敏存储 | ✅ 已覆盖 |
| 159-173 | 数据分析、商户分类、排名 | ✅ 已覆盖 |
| 175-177 | 收益统计、趋势图 | ✅ 已覆盖 |
| 179-183 | 代扣管理、上级/伙伴代扣 | ✅ 已覆盖 |
| 185-205 | 钱包(3种)、提现、税筹 | ✅ 已覆盖 |
| 207-208 | 我的信息、结算卡变更 | ✅ 已覆盖 |
| 210-214 | 消息通知、3天有效期 | ✅ 已覆盖 |

### 遗留问题确认结果 ✅

| 问题 | 文档位置 | 状态 | 确认结果 |
|------|----------|------|----------|
| **T+0秒到费率加1/加2/加3** | 第53行 | ✅ 已确认 | 笔数费，每笔1-3元，可设为0 |
| **服务商偷数据权限** | 第116行 | ✅ 已确认 | 暂不开发，后续版本考虑 |
| **退费分润处理** | 第20、100行 | ✅ 已确认 | 分润回扣，标注退货 |

#### 24. T+0秒到费率规则（已确认）

**功能说明**：商户选择T+0实时结算时，额外收取的手续费

**确认规则**：
- **计费维度**：以每交易一笔为维度，按笔数收费
- **费用金额**：一笔结算费用为 **1到3元**（固定金额/笔，非百分比）
- **与税筹关系**：和税筹9%+3里的"3元/笔"是同一概念
- **可选设置**：也可以不加（设为0元/笔）

| 档位 | 笔数费 | 说明 |
|------|--------|------|
| 不加 | 0元/笔 | 不收取秒到费 |
| 加1 | 1元/笔 | 每笔交易加收1元 |
| 加2 | 2元/笔 | 每笔交易加收2元 |
| 加3 | 3元/笔 | 每笔交易加收3元 |

**计算示例**：
- 交易金额：10000元
- 基础费率：0.60%，手续费=60元
- 秒到加2：**2元/笔**
- 商户实际支付手续费：60 + 2 = **62元**

```sql
-- 政策模板补充秒到费率字段
ALTER TABLE policy_templates ADD COLUMN t0_fee_type SMALLINT DEFAULT 0;  -- 0:不开通 1:加1 2:加2 3:加3
ALTER TABLE policy_templates ADD COLUMN t0_fee_amount DECIMAL(10,2) DEFAULT 0;  -- 秒到费用（元/笔）

-- 终端表补充秒到设置
ALTER TABLE terminals ADD COLUMN t0_enabled BOOLEAN DEFAULT FALSE;       -- 是否开通T+0
ALTER TABLE terminals ADD COLUMN t0_fee_type SMALLINT DEFAULT 0;         -- 秒到费率档位

-- 交易表补充秒到标识
ALTER TABLE transactions ADD COLUMN is_t0 BOOLEAN DEFAULT FALSE;         -- 是否T+0结算
ALTER TABLE transactions ADD COLUMN t0_fee DECIMAL(10,2) DEFAULT 0;      -- 秒到手续费（元/笔）
```

**秒到费用分润**：
- 秒到手续费作为额外收入，可参与分润计算
- 分润规则：按笔数费的差价计算（上级设置减下级设置）

#### 25. 退费分润处理规则（已确认）

**确认规则**：
- **基本情况**：正常业务不会出现退费
- **异常处理**：如出现退费，进行 **分润回扣**
- **标识方式**：交易记录标注为"退货"类型

**实现方案**：
```sql
-- 交易类型增加退货标识
-- trade_type: 1=消费 2=撤销 3=退货 4=服务费 5=押金 6=流量费

-- 退货交易处理流程：
-- 1. 收到退货交易通知
-- 2. 查找原交易对应的分润记录
-- 3. 生成负数分润记录（回扣）
-- 4. 从各级代理商钱包中扣减对应金额
```

```go
// 退货分润回扣处理
func (c *ProfitCalculator) HandleRefund(refundTx *Transaction) error {
    // 1. 查找原交易
    originalTx := c.findOriginalTransaction(refundTx.OriginalTradeNo)

    // 2. 查找原交易的分润记录
    profitRecords := c.profitRepo.GetByTransactionID(originalTx.ID)

    // 3. 生成回扣记录（负数）
    for _, record := range profitRecords {
        refundProfit := &ProfitRecord{
            TransactionID: refundTx.ID,
            AgentID:       record.AgentID,
            ProfitType:    ProfitTypeRefund,  // 退货回扣
            ProfitAmount:  record.ProfitAmount.Neg(),  // 负数
            RefOrderNo:    originalTx.OrderNo,
        }
        c.profitRepo.Create(refundProfit)

        // 4. 从钱包扣减
        c.walletSvc.DeductBalance(
            record.AgentID,
            record.WalletType,
            record.ProfitAmount,
            "退货分润回扣",
            refundTx.OrderNo,
        )
    }
    return nil
}
```

#### 26. 服务商偷数据权限（暂不开发）

**确认规则**：
- **当前版本**：暂不开发此功能
- **功能说明**：根据一定条件设定，可以不跑特定设置下的交易数据
- **业务影响**：不跑数据 = 不结算分润

**后续版本设计预留**：
```sql
-- 数据过滤规则表（后续版本）
CREATE TABLE data_filter_rules (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT,
    rule_type       SMALLINT NOT NULL,       -- 1:金额区间 2:时间段 3:商户类型
    condition_json  JSONB NOT NULL,          -- 过滤条件
    is_enabled      BOOLEAN DEFAULT FALSE,
    effective_date  DATE,
    expire_date     DATE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 示例条件：
-- 金额区间: {"min_amount": 100, "max_amount": 1000}
-- 时间段: {"start_time": "00:00", "end_time": "06:00"}
```

---

## 十四、验证方案

1. **单元测试**: 分润计算4种类型的覆盖测试
2. **集成测试**: 完整交易链路测试
3. **压力测试**: 模拟高并发交易场景
4. **对账验证**: 定时对账任务，核对分润与钱包
5. **税筹验证**: 验证不同扣费方式下代理商实际到账金额一致

---

## 十五、数据库表汇总

### 核心业务表（共计45+张表）

| 模块 | 表名 | 说明 |
|------|------|------|
| **通道** | channels | 支付通道 |
| | channel_configs | 通道扩展配置 |
| | agent_channels | 代理商-通道关联 |
| | tax_channels | 税筹通道 |
| | channel_tax_mappings | 通道-税筹映射 |
| **政策** | policy_templates | 政策模板 |
| | policy_rate_stages | 费率阶梯 |
| | policy_activation_rewards | 激活奖励规则 |
| | policy_deposit_cashbacks | 押金返现规则 |
| | policy_sim_cashbacks | 流量返现规则 |
| | policy_recharge_rewards | 充值钱包奖励规则 |
| **代理商** | agents | 代理商主表 |
| | agent_policies | 代理商政策 |
| | agent_bank_card_logs | 结算卡变更记录 |
| | agent_terminal_stats | 终端统计 |
| | agent_profit_daily | 日收益统计 |
| | agent_profit_monthly | 月收益统计 |
| | agent_ranking | 业绩排名 |
| **终端** | terminals | 终端/机具 |
| | terminal_logs | 终端流转记录 |
| | terminal_payment_deductions | 货款代扣 |
| **商户** | merchants | 商户主表 |
| | merchant_rate_logs | 费率变更记录 |
| | merchant_trade_stats | 交易统计 |
| | merchant_category_stats | 分类统计 |
| | merchant_alert_rules | 预警规则 |
| | merchant_alerts | 预警记录 |
| **交易** | transactions | 交易流水(分区表) |
| **分润** | profit_records | 分润明细 |
| **钱包** | wallets | 正常钱包 |
| | wallet_logs | 钱包流水 |
| | recharge_wallets | 充值钱包 |
| | recharge_wallet_logs | 充值钱包流水 |
| | recharge_wallet_loans | 充值钱包贷款 |
| | deposit_wallets | 沉淀钱包 |
| | deposit_wallet_logs | 沉淀钱包流水 |
| | platform_loans | 平台借款 |
| | withdraw_thresholds | 提现门槛 |
| **提现** | withdrawals | 提现申请 |
| **代扣** | deductions | 代扣协议 |
| **营销** | poster_categories | 海报分类 |
| | marketing_posters | 营销海报 |
| | banners | 滚动图/轮播图 |
| **消息** | notifications | 消息通知 |
| **系统** | system_announcements | 系统公告 |
| | announcement_reads | 公告阅读记录 |

---

## 十六、关键文件清单

| 文件路径 | 说明 |
|----------|------|
| `internal/app/profit/calculator.go` | 分润计算核心引擎 |
| `internal/app/profit/reward_checker.go` | 激活奖励检查 |
| `internal/app/policy/engine.go` | 政策模板引擎 |
| `internal/app/policy/inherit.go` | 政策继承逻辑 |
| `internal/app/terminal/dispatch.go` | 机具下发/回拨 |
| `internal/app/wallet/multi_wallet.go` | 多钱包管理 |
| `internal/app/wallet/sedimentation.go` | 沉淀钱包逻辑 |
| `internal/adapters/channel/*.go` | 通道适配器(8家) |
| `internal/adapters/tax/*.go` | 税筹通道适配器 |
| `migrations/*.sql` | 数据库迁移脚本 |

---

## 十七、开发阶段计划

| 阶段 | 周期 | 内容 |
|------|------|------|
| **第一阶段** | 第1-3周 | 基础架构、通道管理、用户认证 |
| **第二阶段** | 第4-10周 | 代理商、政策模板、终端、商户、交易、分润计算 |
| **第三阶段** | 第11-14周 | 钱包模块、提现、代扣、手动调账 |
| **第四阶段** | 第15-17周 | 数据分析、收益统计、营销、消息 |
| **第五阶段** | 第18-22周 | Vue管理后台、Flutter APP |
| **第六阶段** | 第23-26周 | 测试、安全审计、部署上线 |

---

## 十八、基础设施与部署方案

### 18.1 云服务器租用方案

#### 推荐云服务商对比

| 云服务商 | 优势 | 适用场景 | 价格参考(月) |
|----------|------|----------|--------------|
| **阿里云（推荐）** | 国内市场占有率高、金融行业资质齐全、支付牌照认可 | 支付类系统首选 | ¥800-3000 |
| **腾讯云** | 微信生态对接好、CDN性能强 | 扫码支付、小程序 | ¥700-2500 |
| **华为云** | 政企合规性好、鸿蒙生态支持 | 政企客户、鸿蒙优化 | ¥800-2800 |

#### 基于日交易5000笔的精简配置方案

**业务量分析**：
| 指标 | 数值 | 说明 |
|------|------|------|
| 日交易笔数 | 5,000 | 当前业务量 |
| 日分润计算 | ~20,000条 | 每笔交易产生约4条分润记录 |
| 峰值QPS | ~5 | 假设高峰集中在2小时内 |
| 月数据增量 | ~60万条 | 交易+分润 |

**推荐方案一：最经济方案（约¥500-800/月）**

适用场景：开发测试、小规模上线

| 用途 | 规格 | 数量 | 月费用(阿里云) |
|------|------|------|----------------|
| 应用服务器 | 2核4G 轻量应用服务器 | 1台 | ¥99 (新人优惠) |
| 数据库 | 云数据库PostgreSQL 1核2G 基础版 | 1台 | ¥200 |
| Redis | 云Redis 1G 标准版 | 1台 | ¥130 |
| 对象存储 | OSS 50G | 1套 | ¥5 |
| **合计** | | | **约¥434/月** |

> 💡 新用户首年优惠可低至 ¥200-300/月

**推荐方案二：稳定运营方案（约¥1200-1500/月）**

适用场景：正式上线运营，保障稳定性

| 用途 | 规格 | 数量 | 月费用(阿里云) |
|------|------|------|----------------|
| 应用服务器 | ECS 2核4G (突发性能实例t6) | 2台 | ¥150×2 |
| 数据库 | RDS PostgreSQL 2核4G 高可用版 | 1套 | ¥500 |
| Redis | 云Redis 2G 标准版 | 1套 | ¥200 |
| Kafka | 消息队列(可选，初期可不用) | - | ¥0 |
| 对象存储 | OSS 100G | 1套 | ¥10 |
| CDN | 50G流量包 | 1套 | ¥20 |
| **合计** | | | **约¥1030/月** |

**推荐方案三：单机极简方案（约¥300-500/月）**

适用场景：MVP验证、预算紧张、日交易5000笔

| 用途 | 规格 | 说明 | 月费用 |
|------|------|------|--------|
| 单机部署 | ECS 4核8G | 应用+DB+Redis全部署 | ¥300-400 |

**单机极简方案 - 阿里云详细配置**：

| 配置项 | 推荐值 | 说明 |
|--------|--------|------|
| **实例规格** | ecs.t6-c1m2.large (2核4G) 或 ecs.c6.large (2核4G) | 突发性能实例更便宜 |
| **系统盘** | 40GB 高效云盘 | 装系统+应用 |
| **数据盘** | **100GB SSD云盘** | PostgreSQL+Redis+日志+OSS本地 |
| **带宽** | 5Mbps 固定带宽 | 支持API+APP访问 |
| **镜像** | CentOS 8 / Ubuntu 22.04 | 免费 |

**存储容量详细规划（100GB 数据盘）**：

```
100GB 数据盘分配：
├── /data/postgresql      50GB    # 数据库(可用5年+)
│   ├── 交易数据          ~5GB/年
│   ├── 分润数据          ~2GB/年
│   └── 索引+WAL          ~3GB/年
│
├── /data/redis           5GB     # Redis持久化
│
├── /data/uploads         20GB    # 文件存储(替代OSS)
│   ├── 海报图片          ~2GB
│   ├── 证件照片          ~5GB
│   └── 二维码            ~1GB
│
├── /data/logs            15GB    # 应用日志
│   └── 保留90天自动清理
│
└── /data/backup          10GB    # 本地备份
    └── 每日自动备份，保留7天
```

**单机部署架构图**：

```
┌─────────────────────────────────────────────────────────────────┐
│              单机极简方案架构 (ECS 2核4G + 100GB SSD)             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Nginx (反向代理)                        │   │
│  │                    端口: 80, 443                          │   │
│  └─────────────────────────┬────────────────────────────────┘   │
│                            │                                     │
│            ┌───────────────┴───────────────┐                    │
│            ▼                               ▼                    │
│  ┌─────────────────┐            ┌─────────────────┐            │
│  │   Go API 服务    │            │   Vue 静态资源   │            │
│  │   端口: 8080     │            │   (Nginx托管)   │            │
│  └────────┬────────┘            └─────────────────┘            │
│           │                                                      │
│           │  本地连接 (无网络延迟)                                │
│           │                                                      │
│  ┌────────┴─────────────────────────────────────────────┐       │
│  │                                                       │       │
│  │  ┌─────────────────┐      ┌─────────────────┐        │       │
│  │  │   PostgreSQL    │      │     Redis       │        │       │
│  │  │   端口: 5432    │      │   端口: 6379    │        │       │
│  │  │   50GB存储      │      │   1GB内存缓存   │        │       │
│  │  └─────────────────┘      └─────────────────┘        │       │
│  │                                                       │       │
│  │             /data 挂载点 (100GB SSD)                  │       │
│  └───────────────────────────────────────────────────────┘       │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  定时任务 (Cron)                                          │   │
│  │  ├── 每日凌晨: 数据库备份 → /data/backup                  │   │
│  │  ├── 每日凌晨: 日志清理 (保留90天)                         │   │
│  │  ├── 每小时: 分润计算任务                                  │   │
│  │  └── 每周: 备份上传OSS (灾备)                              │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**阿里云具体采购清单（方案三）**：

| 产品 | 规格 | 配置 | 月费用(原价) | 新人价 |
|------|------|------|-------------|--------|
| ECS服务器 | 2核4G 突发性能t6 | | ¥99 | ¥38 |
| 系统盘 | 高效云盘 | 40GB | ¥16 | 含在ECS |
| **数据盘** | **SSD云盘** | **100GB** | ¥80 | ¥30 |
| 公网IP | 固定带宽 | 5Mbps | ¥115 | ¥45 |
| **总计** | | | **¥310/月** | **¥113/月** |

> 💡 新用户首年优惠后约 **¥100-150/月**

**为什么选100GB数据盘**：

| 因素 | 说明 |
|------|------|
| **3年数据** | 交易+分润约15GB |
| **文件存储** | 替代OSS约10GB |
| **日志备份** | 约25GB |
| **余量空间** | 50GB扩展余量 |
| **性价比** | 100GB SSD约¥80/月，性价比最高 |

**后续扩展路径**：

```
单机100GB → 需要扩展时：
├── 方案A：数据盘扩容 100GB → 200GB（在线扩容，+¥80/月）
├── 方案B：分离数据库 → 购买RDS（+¥350/月）
└── 方案C：迁移到方案二（ECS+RDS分离部署）
```

#### 各方案对比

| 方案 | 月费用 | 高可用 | 适用阶段 | 支撑交易量 |
|------|--------|--------|----------|------------|
| **方案一** | ¥400-500 | ❌ | 开发测试 | 日5000笔 |
| **方案二** | ¥1000-1200 | ✅ | 正式运营 | 日1-2万笔 |
| **方案三** | ¥300-400 | ❌ | MVP验证 | 日3000笔 |

#### 省钱技巧

1. **新用户优惠**：阿里云/腾讯云新用户首年2-3折
2. **包年折扣**：包年比包月便宜30-50%
3. **预留实例**：长期使用可购买预留实例
4. **按量付费**：开发期间按量付费，正式运营再包年
5. **学生优惠**：有学生身份可享受超低价

#### 推荐采购策略

```
开发期（0-3月）：
├── 方案三（单机4核8G）≈ ¥300/月
├── 新用户优惠 2折起
└── 预算：¥100-200/月

上线初期（3-6月）：
├── 方案一（轻量服务器+云DB）≈ ¥500/月
├── 根据实际交易量调整
└── 预算：¥400-600/月

稳定运营（6月+）：
├── 方案二（ECS+RDS高可用）≈ ¥1000/月
├── 包年付费享折扣
└── 预算：¥800-1200/月
```

### 18.2 数据量增长应对方案（可扩展架构）

#### 当前业务量（日交易5000笔）

| 数据类型 | 单条大小 | 日增量 | 月增量 | 年增量 |
|----------|----------|--------|--------|--------|
| 交易流水 | ~500B | 5,000条 | 15万条 | 180万条 |
| 分润明细 | ~300B | 20,000条 | 60万条 | 720万条 |
| 钱包流水 | ~200B | 5,000条 | 15万条 | 180万条 |

**当前阶段存储需求**：约50GB/年（完全可控）

#### 扩展性设计原则

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     可扩展架构设计原则                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 【无状态设计】应用服务器无状态，随时可水平扩展                              │
│  2. 【分区预留】数据库表设计预留分区字段，后期无需改表                          │
│  3. 【接口解耦】模块间通过接口通信，便于拆分微服务                              │
│  4. 【配置化】关键参数配置化，不硬编码                                         │
│  5. 【监控预警】提前部署监控，容量预警                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 架构演进路线图

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                           架构演进路线图                                       │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  阶段一（当前）         阶段二（1-2万笔/日）    阶段三（5万+笔/日）            │
│  日交易5000笔          日交易1-2万笔          日交易5万+笔                   │
│                                                                               │
│  ┌─────────────┐       ┌─────────────┐       ┌─────────────┐                │
│  │   单机部署   │  ──▶  │  主从分离    │  ──▶  │  微服务集群  │                │
│  │   4核8G     │       │  读写分离    │       │  K8s部署    │                │
│  └─────────────┘       └─────────────┘       └─────────────┘                │
│                                                                               │
│  成本：¥300/月         成本：¥1000/月        成本：¥5000+/月                │
│  扩展：手动迁移         扩展：加从库          扩展：弹性伸缩                  │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

#### 预留扩展设计要点

**1. 数据库层面**

```sql
-- 表结构预留分区键（即使初期不分区，也预留字段）
CREATE TABLE transactions (
    id              BIGSERIAL,
    trade_no        VARCHAR(64) NOT NULL,
    channel_id      BIGINT NOT NULL,          -- 预留：可按通道分库
    trade_time      TIMESTAMP NOT NULL,       -- 预留：可按时间分区
    -- 其他字段...
    PRIMARY KEY (id)
);

-- 初期：不启用分区，当单表超过1000万条时再启用
-- 扩展方式：ALTER TABLE ... PARTITION BY RANGE (trade_time)

-- 索引预留（初期创建必要索引，后续按需添加）
CREATE INDEX idx_transactions_channel_time ON transactions(channel_id, trade_time);
CREATE INDEX idx_transactions_agent ON transactions(agent_id, trade_time);
```

**2. 应用层面**

```go
// 接口定义，便于后续替换实现
type TransactionRepository interface {
    Save(tx *Transaction) error
    FindByID(id int64) (*Transaction, error)
    FindByTimeRange(start, end time.Time) ([]*Transaction, error)
}

// 初期：单库实现
type PostgresTransactionRepo struct {
    db *gorm.DB
}

// 扩展：分库实现（后续平滑替换）
type ShardingTransactionRepo struct {
    shards []*gorm.DB
    router ShardRouter
}
```

**3. 缓存层面**

```go
// 缓存接口定义
type CacheService interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
}

// 初期：单节点Redis
// 扩展：Redis Cluster（接口不变，配置切换）
```

**4. 消息队列（可选，后续启用）**

```
初期架构（同步处理）：
API → 业务处理 → 数据库

扩展架构（异步处理）：
API → Kafka → 消费者 → 数据库

切换成本：
1. 引入Kafka（¥600/月）
2. 改造交易接收模块（约1天工作量）
3. 新增消费者服务
```

#### 扩展触发条件与方案

| 指标 | 阈值 | 扩展动作 | 成本增加 |
|------|------|----------|----------|
| 日交易量 | >1万笔 | 启用读写分离 | +¥500/月 |
| 单表数据 | >1000万条 | 启用表分区 | ¥0 |
| API响应 | >500ms | 加应用服务器 | +¥300/台 |
| 数据库CPU | >70% | 升级RDS配置 | +¥300/月 |
| 日交易量 | >5万笔 | 引入消息队列 | +¥600/月 |

#### 监控预警配置

```yaml
# 容量监控配置
alerts:
  - name: 数据库存储告警
    condition: db_storage_used > 80%
    action: 通知管理员扩容

  - name: 交易量突增告警
    condition: daily_transactions > 10000
    action: 评估是否需要架构升级

  - name: API响应慢告警
    condition: api_latency_p99 > 500ms
    action: 检查瓶颈，考虑水平扩展

  - name: 缓存命中率告警
    condition: cache_hit_rate < 80%
    action: 检查缓存策略
```

#### 总结：低成本起步 + 平滑扩展

| 阶段 | 日交易量 | 月成本 | 扩展方式 |
|------|----------|--------|----------|
| **起步** | 5000 | ¥300-500 | 单机部署 |
| **增长** | 1-2万 | ¥1000-1500 | 读写分离 + 缓存 |
| **爆发** | 5万+ | ¥3000-5000 | 微服务 + 消息队列 |
| **规模** | 10万+ | ¥10000+ | K8s + 分库分表 |

**核心思路**：
1. ✅ 代码架构预留扩展点（接口化设计）
2. ✅ 数据库表结构预留分区键
3. ✅ 初期用最低成本验证业务
4. ✅ 根据实际增长按需扩展
5. ✅ 监控预警，提前发现瓶颈

### 18.3 阿里云存储容量规划

#### 存储需求详细计算（日交易5000笔）

**数据库存储计算**：

| 表 | 单条大小 | 日增量 | 年增量 | 年存储 |
|------|----------|--------|--------|--------|
| transactions（交易） | 500B | 5,000条 | 180万条 | ~900MB |
| profit_records（分润） | 300B | 20,000条 | 720万条 | ~2.2GB |
| wallet_logs（钱包流水） | 200B | 5,000条 | 180万条 | ~360MB |
| agents（代理商） | 1KB | 100条 | 3.6万条 | ~36MB |
| merchants（商户） | 800B | 500条 | 18万条 | ~144MB |
| terminals（终端） | 600B | 200条 | 7.2万条 | ~43MB |
| 其他表 | - | - | - | ~500MB |
| **数据库总计** | | | | **~4.2GB/年** |

**索引存储（约为数据的30-50%）**：~2GB/年

**对象存储（OSS）计算**：

| 类型 | 单个大小 | 年增量 | 年存储 |
|------|----------|--------|--------|
| 营销海报 | 500KB | 200张 | ~100MB |
| 代理商证件 | 200KB | 5000份 | ~1GB |
| 二维码图片 | 50KB | 10000张 | ~500MB |
| 导出报表 | 100KB | 1000份 | ~100MB |
| **OSS总计** | | | **~2GB/年** |

**日志存储**：

| 类型 | 日增量 | 年存储 |
|------|--------|--------|
| 应用日志 | 50MB | ~18GB |
| 审计日志 | 10MB | ~3.6GB |
| **日志总计** | | **~22GB/年** |

#### 阿里云存储配置建议

**数据库存储（RDS PostgreSQL）**：

| 阶段 | 业务量 | 推荐存储 | 费用/月 | 可用年限 |
|------|--------|----------|---------|----------|
| **初期** | 日5000笔 | 20GB SSD | ¥0（含在RDS中） | 3年+ |
| **增长** | 日1万笔 | 50GB SSD | ¥50 | 5年+ |
| **稳定** | 日2万笔 | 100GB SSD | ¥100 | 5年+ |

> 💡 阿里云RDS基础版默认含20GB存储，初期完全够用

**推荐配置（基于日5000笔，考虑3年扩展）**：

```
┌─────────────────────────────────────────────────────────────────────┐
│                   阿里云存储配置推荐                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  【数据库 RDS PostgreSQL】                                           │
│  ├── 规格：2核4G 基础版（初期够用）                                   │
│  ├── 存储：50GB SSD（含3年数据增长余量）                              │
│  ├── 费用：约 ¥350/月（包年更优惠）                                   │
│  └── 扩展：支持在线扩容到1000GB                                      │
│                                                                      │
│  【对象存储 OSS】                                                     │
│  ├── 存储包：40GB 标准型                                             │
│  ├── 费用：¥9/年（资源包）                                           │
│  └── 扩展：按量自动扩展                                              │
│                                                                      │
│  【Redis 缓存】                                                       │
│  ├── 规格：1GB 标准版                                                │
│  ├── 费用：约 ¥130/月                                                │
│  └── 说明：缓存不占用持久存储                                         │
│                                                                      │
│  【日志服务 SLS】（可选）                                             │
│  ├── 存储：50GB                                                      │
│  ├── 费用：约 ¥30/月                                                 │
│  └── 替代方案：本地文件+定期清理                                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### 存储扩展策略

**阿里云RDS存储扩容**：

```
初始配置：50GB SSD
├── 支持在线扩容（不停机）
├── 扩容步长：5GB
├── 最大容量：6000GB（PostgreSQL）
└── 扩容费用：约 ¥1/GB/月

扩容时机建议：
├── 使用率达到 70% 时预警
├── 使用率达到 80% 时扩容
└── 每次扩容 50% 容量
```

**存储成本控制**：

| 策略 | 说明 | 节省比例 |
|------|------|----------|
| 历史数据归档 | 2年前数据转OSS归档存储 | 90% |
| 日志定期清理 | 保留90天，定期删除 | 70% |
| 图片压缩 | 海报/证件压缩存储 | 50% |
| 冷热分离 | 热数据SSD，冷数据HDD | 60% |

#### 阿里云完整配置清单（推荐）

**基于日交易5000笔，考虑3年扩展**：

| 产品 | 规格 | 存储 | 月费用 |
|------|------|------|--------|
| ECS服务器 | 2核4G 突发性能t6 | 40GB 系统盘 | ¥80 |
| RDS PostgreSQL | 2核4G 基础版 | **50GB SSD** | ¥350 |
| Redis | 1GB 标准版 | - | ¥130 |
| OSS对象存储 | 标准型 | **40GB资源包** | ¥9/年 |
| **总计** | | | **约¥570/月** |

> 💡 新用户首购优惠可低至 3-4折，约 ¥200/月

#### 存储容量预警设置

```sql
-- 数据库存储监控脚本
SELECT
    pg_database.datname AS database_name,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
ORDER BY pg_database_size(pg_database.datname) DESC;

-- 阿里云控制台设置告警规则
-- 1. 存储使用率 > 70%：发送预警通知
-- 2. 存储使用率 > 85%：自动扩容（需开启）
```

#### 总结

| 存储类型 | 推荐容量 | 可用年限 | 月费用 |
|----------|----------|----------|--------|
| 数据库 | 50GB SSD | 3-5年 | ¥0（含在RDS） |
| 对象存储 | 40GB | 10年+ | ¥9/年 |
| 系统盘 | 40GB | 永久 | 含在ECS |

**最终建议**：
- ✅ RDS选择 **50GB SSD**，足够3年使用
- ✅ OSS选择 **40GB资源包**，按年付费更划算
- ✅ 开启存储监控和自动扩容
- ✅ 定期归档历史数据到OSS低频存储

```
┌─────────────────────────────────────────────────────────────────┐
│                        数据分片策略                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐                                                │
│  │  交易表      │  按月分区 (PARTITION BY RANGE)                 │
│  │ transactions │  2024-01, 2024-02, ... 自动创建分区            │
│  └─────────────┘                                                │
│                                                                  │
│  ┌─────────────┐                                                │
│  │  分润明细表   │  按通道ID + 月份 分表                          │
│  │profit_records│  profit_records_channel1_202401               │
│  └─────────────┘                                                │
│                                                                  │
│  ┌─────────────┐                                                │
│  │  历史数据    │  超过2年的数据归档到冷存储                       │
│  │  归档策略    │  定期迁移到 OSS/S3 + 压缩存储                   │
│  └─────────────┘                                                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### PostgreSQL 分区表实现

```sql
-- 交易表按月自动分区
CREATE TABLE transactions (
    id              BIGSERIAL,
    trade_no        VARCHAR(64) NOT NULL,
    trade_time      TIMESTAMP NOT NULL,
    -- 其他字段...
    PRIMARY KEY (id, trade_time)
) PARTITION BY RANGE (trade_time);

-- 创建分区（可通过定时任务自动创建）
CREATE TABLE transactions_2024_01 PARTITION OF transactions
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE transactions_2024_02 PARTITION OF transactions
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- 自动分区维护脚本（定时任务每月执行）
CREATE OR REPLACE FUNCTION create_monthly_partition()
RETURNS void AS $$
DECLARE
    next_month DATE;
    partition_name TEXT;
BEGIN
    next_month := DATE_TRUNC('month', NOW()) + INTERVAL '1 month';
    partition_name := 'transactions_' || TO_CHAR(next_month, 'YYYY_MM');

    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF transactions
         FOR VALUES FROM (%L) TO (%L)',
        partition_name,
        next_month,
        next_month + INTERVAL '1 month'
    );
END;
$$ LANGUAGE plpgsql;
```

#### 读写分离架构

```
                    ┌───────────────┐
                    │   应用服务器   │
                    └───────┬───────┘
                            │
              ┌─────────────┴─────────────┐
              ▼                           ▼
    ┌─────────────────┐         ┌─────────────────┐
    │   主库 (写入)    │────────▶│   从库 (读取)    │
    │   PostgreSQL    │  同步    │   PostgreSQL    │
    │   RDS 高可用    │         │   只读实例       │
    └─────────────────┘         └─────────────────┘

    写操作：INSERT/UPDATE/DELETE → 主库
    读操作：SELECT (报表/查询) → 从库
```

#### 缓存策略

```go
// Redis 缓存层设计
type CacheStrategy struct {
    // 热点数据缓存
    AgentTree       time.Duration // 代理商树结构：24小时
    PolicyTemplate  time.Duration // 政策模板：1小时
    ChannelConfig   time.Duration // 通道配置：1小时

    // 实时数据缓存
    WalletBalance   time.Duration // 钱包余额：实时更新
    TodayProfit     time.Duration // 今日收益：5分钟

    // 统计数据缓存
    DailyStats      time.Duration // 日统计：10分钟
    MonthlyStats    time.Duration // 月统计：1小时
}
```

### 18.3 技术框架选型详解

#### 后端框架：Go + Gin + GORM

**选择理由**：
| 特性 | 说明 |
|------|------|
| **高并发** | Go协程轻量级，支持百万级并发 |
| **编译部署** | 单二进制文件，无依赖 |
| **Gin框架** | 高性能HTTP框架，路由性能极佳 |
| **GORM** | 成熟的ORM，支持PostgreSQL |
| **社区活跃** | 支付行业广泛使用 |

**项目结构**：
```
backend/
├── cmd/
│   ├── api/main.go           # API服务入口
│   ├── consumer/main.go      # Kafka消费者
│   └── scheduler/main.go     # 定时任务
├── internal/
│   ├── app/                  # 业务逻辑层
│   │   ├── agent/
│   │   ├── channel/
│   │   ├── merchant/
│   │   ├── policy/
│   │   ├── profit/           # 分润计算核心
│   │   ├── terminal/
│   │   ├── transaction/
│   │   └── wallet/
│   ├── domain/               # 领域模型
│   ├── repository/           # 数据访问层
│   └── service/              # 服务层
├── pkg/
│   ├── cache/                # Redis封装
│   ├── db/                   # 数据库连接
│   ├── kafka/                # 消息队列
│   ├── logger/               # 日志
│   └── utils/                # 工具函数
└── configs/                  # 配置文件
```

#### 前端框架

**APP端：Flutter 3.x**
- 一套代码，iOS/Android/鸿蒙三端运行
- Dart语言，热重载开发效率高
- 丰富的UI组件库

**管理后台：Vue 3 + Element Plus**
- 企业级后台标配
- TypeScript支持
- 响应式设计

### 18.4 华为鸿蒙系统兼容性

#### Flutter 对鸿蒙的支持

**当前状态（2024年）**：
| 项目 | 支持情况 | 说明 |
|------|----------|------|
| **Flutter OHOS** | ✅ 官方支持 | Flutter 3.22+ 原生支持鸿蒙 |
| **鸿蒙NEXT** | ✅ 兼容 | 需使用 flutter_ohos 插件 |
| **跨平台兼容** | ✅ 一套代码 | iOS/Android/HarmonyOS |

**鸿蒙适配方案**：

```yaml
# pubspec.yaml 配置
dependencies:
  flutter:
    sdk: flutter
  flutter_ohos: ^1.0.0  # 鸿蒙适配插件

# 编译命令
flutter build apk       # Android
flutter build ios       # iOS
flutter build hap       # 鸿蒙 (HAP包)
```

**鸿蒙特殊适配点**：

1. **推送通知**：需集成华为Push Kit
2. **支付SDK**：鸿蒙版支付SDK适配
3. **分享功能**：鸿蒙分享能力API

```dart
// 平台判断示例
import 'dart:io';

bool isHarmonyOS() {
  // 通过系统版本判断
  return Platform.operatingSystem == 'harmonyos';
}

// 条件编译
void initPushService() {
  if (isHarmonyOS()) {
    // 华为Push Kit
    HuaweiPushKit.init();
  } else if (Platform.isAndroid) {
    // FCM / 极光推送
    FirebaseMessaging.init();
  } else if (Platform.isIOS) {
    // APNs
    ApplePushNotification.init();
  }
}
```

#### 鸿蒙原生开发备选方案

如果未来需要深度适配鸿蒙特性：

| 方案 | 说明 | 工作量 |
|------|------|--------|
| **Flutter OHOS** | 推荐，一套代码三端 | 低 |
| **ArkTS原生** | 鸿蒙原生开发 | 高（需单独开发） |
| **跨平台+原生** | Flutter主体 + 鸿蒙原生插件 | 中 |

**推荐方案**：使用 Flutter 开发，可同时覆盖 iOS、Android、鸿蒙三端。

---

## 十九、分润计算引擎详细设计

### 19.1 分润计算核心架构

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              分润计算引擎架构                                            │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                          │
│  ┌─────────────────┐                                                                    │
│  │  交易数据入口    │  ← 通道推送/定时拉取                                               │
│  │  (Webhook/API)  │                                                                    │
│  └────────┬────────┘                                                                    │
│           │                                                                              │
│           ▼                                                                              │
│  ┌─────────────────┐                                                                    │
│  │  交易类型分发器   │  ← 根据 trade_type 分发到不同计算器                                │
│  │  TransactionRouter                                                                   │
│  └────────┬────────┘                                                                    │
│           │                                                                              │
│   ┌───────┴───────┬───────────────┬───────────────┐                                    │
│   ▼               ▼               ▼               ▼                                    │
│ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐                                   │
│ │交易分润   │ │激活奖励   │ │押金返现   │ │流量返现   │                                   │
│ │Calculator│ │Checker   │ │Calculator│ │Calculator│                                   │
│ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘                                   │
│      │            │            │            │                                          │
│      └────────────┴────────────┴────────────┘                                          │
│                         │                                                               │
│                         ▼                                                               │
│              ┌─────────────────┐                                                        │
│              │  代理商链追溯器   │  ← 向上追溯所有上级代理商                               │
│              │  AgentChainTracer                                                        │
│              └────────┬────────┘                                                        │
│                       │                                                                 │
│                       ▼                                                                 │
│              ┌─────────────────┐                                                        │
│              │  分润分配引擎    │  ← 计算每级代理商应得分润                               │
│              │  ProfitAllocator                                                         │
│              └────────┬────────┘                                                        │
│                       │                                                                 │
│                       ▼                                                                 │
│              ┌─────────────────┐                                                        │
│              │  钱包入账服务    │  ← 分润入对应钱包                                       │
│              │  WalletService                                                           │
│              └─────────────────┘                                                        │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### 19.2 四种分润类型详解

#### 类型一：交易分润（核心）

**触发条件**：普通消费交易（trade_type = 1）

**计算公式**：
```
分润金额 = 交易金额 × (商户费率 - 结算费率)

示例：
交易金额 = 10000元
商户费率 = 0.60%（万分之60）
结算费率 = 0.52%（万分之52）
分润金额 = 10000 × (0.0060 - 0.0052) = 10000 × 0.0008 = 8元
```

**多层级分润计算**：
```
商户费率：0.60%  ──────────────────────────────────────┐
                                                        │ 差价8元
直属代理结算价：0.52%  ───────────────────────┐         │
                                              │ 差价4元  ├── 直属代理得8元
二级代理结算价：0.48%  ────────────┐          │         │
                                   │ 差价6元   ├── 二级代理得4元
一级代理结算价：0.42%  ───┐        │          │         │
                          │ 差价4元 ├── 一级代理得6元    │
通道成本：0.38%  ─────────┴── 机构得4元       │         │
                                              │         │
验证：8 + 4 + 6 + 4 = 22元 = 10000 × (0.60% - 0.38%)  ✓
```

**Go代码实现**：

```go
// internal/app/profit/calculator.go

package profit

import (
    "github.com/shopspring/decimal"
)

// TransactionProfitCalculator 交易分润计算器
type TransactionProfitCalculator struct {
    agentRepo    AgentRepository
    policyRepo   PolicyRepository
    merchantRepo MerchantRepository
    profitRepo   ProfitRepository
    walletSvc    WalletService
}

// ProfitResult 分润计算结果
type ProfitResult struct {
    AgentID      int64           `json:"agent_id"`
    AgentName    string          `json:"agent_name"`
    Layer        int             `json:"layer"`        // 层级 1=直属
    SelfRate     decimal.Decimal `json:"self_rate"`    // 自己的结算价
    LowerRate    decimal.Decimal `json:"lower_rate"`   // 下级的结算价/商户费率
    RateDiff     decimal.Decimal `json:"rate_diff"`    // 费率差
    ProfitAmount decimal.Decimal `json:"profit_amount"` // 分润金额
    WalletType   int             `json:"wallet_type"`  // 入账钱包类型
}

// Calculate 计算交易分润
func (c *TransactionProfitCalculator) Calculate(tx *Transaction) ([]*ProfitResult, error) {
    results := make([]*ProfitResult, 0)

    // 1. 获取商户信息及费率
    merchant, err := c.merchantRepo.GetByID(tx.MerchantID)
    if err != nil {
        return nil, err
    }

    // 根据支付类型获取对应费率
    merchantRate := c.getMerchantRate(merchant, tx.PayType)

    // 2. 获取代理商链（从直属代理向上追溯）
    agentChain, err := c.agentRepo.GetAgentChain(tx.AgentID)
    if err != nil {
        return nil, err
    }

    // 3. 获取通道成本费率
    channel, _ := c.channelRepo.GetByID(tx.ChannelID)
    costRate := c.getChannelCostRate(channel, tx.PayType)

    // 4. 逐级计算分润
    prevRate := merchantRate // 上一级费率（初始为商户费率）

    for i, agent := range agentChain {
        // 获取该代理商的结算价（考虑阶梯费率）
        settlementRate := c.getSettlementRate(agent, tx.ChannelID, tx.PayType, merchant.RegisterTime)

        // 费率差 = 上一级费率 - 当前结算价
        rateDiff := prevRate.Sub(settlementRate)

        // 如果没有费率差，跳过（取下原则）
        if rateDiff.LessThanOrEqual(decimal.Zero) {
            prevRate = settlementRate
            continue
        }

        // 分润金额 = 交易金额 × 费率差
        profitAmount := tx.Amount.Mul(rateDiff)

        // 最低分润检查（可配置）
        if profitAmount.LessThan(decimal.NewFromFloat(0.01)) {
            prevRate = settlementRate
            continue
        }

        results = append(results, &ProfitResult{
            AgentID:      agent.ID,
            AgentName:    agent.AgentName,
            Layer:        i + 1,
            SelfRate:     settlementRate,
            LowerRate:    prevRate,
            RateDiff:     rateDiff,
            ProfitAmount: profitAmount,
            WalletType:   WalletTypeProfit, // 分润钱包
        })

        prevRate = settlementRate

        // 如果结算价已经等于成本价，停止追溯
        if settlementRate.LessThanOrEqual(costRate) {
            break
        }
    }

    return results, nil
}

// getSettlementRate 获取结算价（考虑阶梯费率）
func (c *TransactionProfitCalculator) getSettlementRate(
    agent *Agent,
    channelID int64,
    payType int,
    merchantRegisterTime time.Time,
) decimal.Decimal {
    // 1. 获取代理商对应通道的政策
    policy, _ := c.policyRepo.GetAgentPolicy(agent.ID, channelID)
    if policy == nil {
        return decimal.Zero
    }

    // 2. 获取基础结算价
    baseRate := c.getBaseRate(policy, payType)

    // 3. 检查阶梯费率调整
    stages, _ := c.policyRepo.GetRateStages(policy.TemplateID)
    if len(stages) == 0 {
        return baseRate
    }

    // 4. 计算天数（根据calc_base判断）
    var days int
    for _, stage := range stages {
        if stage.CalcBase == 1 {
            // 按商户入网时间计算
            days = int(time.Since(merchantRegisterTime).Hours() / 24)
        } else {
            // 按代理商入网时间计算
            days = int(time.Since(agent.RegisterTime).Hours() / 24)
        }

        // 匹配阶段
        if days >= stage.DayStart && (stage.DayEnd == -1 || days <= stage.DayEnd) {
            return baseRate.Add(stage.RateAdjust)
        }
    }

    return baseRate
}
```

#### 类型二：激活奖励

**触发条件**：商户在考核期内达成交易额目标

**计算流程**：
```
┌────────────────┐
│   终端激活      │  ← 首笔交易成功
│  (activate)    │
└───────┬────────┘
        │
        ▼
┌────────────────┐
│  创建奖励记录   │  ← status = 考核中
│  assess_start  │
│  assess_end    │
└───────┬────────┘
        │
        ▼
┌────────────────┐
│  定时任务检查   │  ← 每日凌晨运行
│  (Scheduler)   │
└───────┬────────┘
        │
   ┌────┴────┐
   ▼         ▼
┌─────┐   ┌─────┐
│达标  │   │未达标│
└──┬──┘   └──┬──┘
   │         │
   ▼         ▼
入账奖励   继续等待/过期
```

**Go代码实现**：

```go
// internal/app/profit/reward_checker.go

// ActivationRewardChecker 激活奖励检查器
type ActivationRewardChecker struct {
    rewardRepo     RewardRepository
    transactionRepo TransactionRepository
    walletSvc      WalletService
}

// CheckAndGrant 检查并发放激活奖励
func (c *ActivationRewardChecker) CheckAndGrant() error {
    // 1. 查询所有"考核中"的奖励记录
    pendingRewards, err := c.rewardRepo.GetPendingRewards()
    if err != nil {
        return err
    }

    for _, reward := range pendingRewards {
        // 2. 检查是否已过考核期
        if time.Now().After(reward.AssessEndDate) {
            // 标记为未达标
            c.rewardRepo.UpdateStatus(reward.ID, StatusNotAchieved)
            continue
        }

        // 3. 统计考核期内交易额
        actualAmount, err := c.transactionRepo.SumAmountByMerchant(
            reward.MerchantID,
            reward.AssessStartDate,
            reward.AssessEndDate,
        )
        if err != nil {
            continue
        }

        // 4. 判断是否达标
        if actualAmount.GreaterThanOrEqual(reward.RequireAmount) {
            // 更新为已达标
            c.rewardRepo.UpdateStatus(reward.ID, StatusAchieved)
            c.rewardRepo.UpdateActualAmount(reward.ID, actualAmount)

            // 5. 发放奖励到奖励钱包
            c.walletSvc.AddBalance(
                reward.AgentID,
                WalletTypeReward, // 奖励钱包
                reward.RewardAmount,
                "激活奖励",
                reward.RecordNo,
            )

            // 6. 如果有多层级奖励，继续发放上级
            if reward.LayerRewards != nil {
                c.grantLayerRewards(reward)
            }
        }
    }

    return nil
}

// grantLayerRewards 发放多层级奖励
func (c *ActivationRewardChecker) grantLayerRewards(reward *ActivationReward) {
    // 解析多层级奖励配置 [{layer:1,amount:50},{layer:2,amount:20}]
    var layers []LayerReward
    json.Unmarshal(reward.LayerRewards, &layers)

    // 获取代理商链
    agentChain, _ := c.agentRepo.GetAgentChain(reward.AgentID)

    for _, layer := range layers {
        if layer.Layer <= len(agentChain) {
            agent := agentChain[layer.Layer-1]
            c.walletSvc.AddBalance(
                agent.ID,
                WalletTypeReward,
                decimal.NewFromFloat(layer.Amount),
                fmt.Sprintf("激活奖励(第%d层)", layer.Layer),
                reward.RecordNo,
            )
        }
    }
}
```

#### 类型三：押金返现

**触发条件**：通道返回押金扣取成功（trade_type = 5）

**计算规则**：
```
押金金额    通道收取    返给代理商
99元    →   99元    →   0-99元（可配置）
199元   →   199元   →   0-199元（可配置）
299元   →   299元   →   0-299元（可配置）
```

**Go代码实现**：

```go
// internal/app/profit/deposit_cashback.go

// DepositCashbackCalculator 押金返现计算器
type DepositCashbackCalculator struct {
    policyRepo PolicyRepository
    walletSvc  WalletService
}

// Calculate 计算押金返现
func (c *DepositCashbackCalculator) Calculate(tx *Transaction) ([]*ProfitResult, error) {
    results := make([]*ProfitResult, 0)

    // 1. 获取终端信息
    terminal, _ := c.terminalRepo.GetBySN(tx.TerminalSN)
    if terminal == nil {
        return nil, errors.New("terminal not found")
    }

    // 2. 获取押金返现规则
    rule, _ := c.policyRepo.GetDepositCashbackRule(
        terminal.OwnerAgentID,
        tx.ChannelID,
        tx.Amount, // 押金金额
    )
    if rule == nil {
        return nil, nil // 无返现规则
    }

    // 3. 计算返现金额
    cashbackAmount := rule.CashbackAmount

    // 4. 创建返现记录
    results = append(results, &ProfitResult{
        AgentID:      terminal.OwnerAgentID,
        Layer:        1,
        ProfitAmount: cashbackAmount,
        WalletType:   WalletTypeService, // 服务费钱包
    })

    return results, nil
}
```

#### 类型四：流量费返现

**触发条件**：通道返回流量费扣取成功（trade_type = 6）

**计算规则**：
```
流量费    通道收取    首次返现    续费返现
79元   →   79元   →   69元     →   59元
89元   →   89元   →   79元     →   69元
99元   →   99元   →   89元     →   79元

多层级分配示例（首次79元流量费）：
├── 直属代理：69元
├── 二级代理：5元
└── 三级代理：3元（余下）
```

**Go代码实现**：

```go
// internal/app/profit/sim_cashback.go

// SimCashbackCalculator 流量费返现计算器
type SimCashbackCalculator struct {
    policyRepo PolicyRepository
    walletSvc  WalletService
}

// Calculate 计算流量费返现
func (c *SimCashbackCalculator) Calculate(tx *Transaction) ([]*ProfitResult, error) {
    results := make([]*ProfitResult, 0)

    // 1. 获取终端信息
    terminal, _ := c.terminalRepo.GetBySN(tx.TerminalSN)

    // 2. 判断是首次还是续费
    isFirstYear := c.isFirstYearSim(terminal)

    // 3. 获取流量费返现规则
    rules, _ := c.policyRepo.GetSimCashbackRules(
        terminal.OwnerAgentID,
        tx.ChannelID,
        isFirstYear,
    )

    // 4. 获取代理商链
    agentChain, _ := c.agentRepo.GetAgentChain(terminal.OwnerAgentID)

    // 5. 按层级分配返现
    for _, rule := range rules {
        if rule.Layer <= len(agentChain) {
            agent := agentChain[rule.Layer-1]
            results = append(results, &ProfitResult{
                AgentID:      agent.ID,
                Layer:        rule.Layer,
                ProfitAmount: rule.CashbackAmount,
                WalletType:   WalletTypeService, // 服务费钱包
            })
        }
    }

    return results, nil
}

// isFirstYearSim 判断是否首年流量费
func (c *SimCashbackCalculator) isFirstYearSim(terminal *Terminal) bool {
    if terminal.FirstSimTime.IsZero() {
        return true
    }
    // 首次流量费时间超过1年，则为续费
    return time.Since(terminal.FirstSimTime) < 365*24*time.Hour
}
```

### 19.3 代理商链追溯算法

#### 物化路径（Materialized Path）方案

```sql
-- 代理商表结构
CREATE TABLE agents (
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGINT,
    path        VARCHAR(500) DEFAULT '',  -- 物化路径 如 /1/5/12/
    level       INT DEFAULT 1
);

-- 示例数据
-- id=1, path='/', level=1        (总部)
-- id=5, path='/1/', level=2      (一级代理)
-- id=12, path='/1/5/', level=3   (二级代理)
-- id=25, path='/1/5/12/', level=4 (三级代理，直属代理)

-- 向上追溯查询（从直属代理到总部）
SELECT * FROM agents
WHERE '/1/5/12/25/' LIKE path || '%'
ORDER BY level DESC;
```

#### Go实现

```go
// internal/app/agent/chain_tracer.go

// AgentChainTracer 代理商链追溯器
type AgentChainTracer struct {
    db    *gorm.DB
    cache *redis.Client
}

// GetAgentChain 获取代理商链（向上追溯）
func (t *AgentChainTracer) GetAgentChain(agentID int64) ([]*Agent, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("agent:chain:%d", agentID)
    cached, err := t.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var chain []*Agent
        json.Unmarshal([]byte(cached), &chain)
        return chain, nil
    }

    // 2. 查询当前代理商
    var agent Agent
    if err := t.db.First(&agent, agentID).Error; err != nil {
        return nil, err
    }

    // 3. 使用路径向上追溯
    var chain []*Agent
    err = t.db.Raw(`
        SELECT * FROM agents
        WHERE ? LIKE path || '%'
        ORDER BY level DESC
    `, agent.Path+strconv.FormatInt(agentID, 10)+"/").Scan(&chain).Error

    if err != nil {
        return nil, err
    }

    // 4. 缓存结果（24小时）
    data, _ := json.Marshal(chain)
    t.cache.Set(ctx, cacheKey, data, 24*time.Hour)

    return chain, nil
}

// PostgreSQL 递归CTE方案（备选）
func (t *AgentChainTracer) GetAgentChainByCTE(agentID int64) ([]*Agent, error) {
    var chain []*Agent
    err := t.db.Raw(`
        WITH RECURSIVE agent_tree AS (
            -- 起点：当前代理商
            SELECT id, parent_id, agent_name, path, level,
                   credit_rate, debit_rate
            FROM agents
            WHERE id = ?

            UNION ALL

            -- 递归：向上追溯
            SELECT a.id, a.parent_id, a.agent_name, a.path, a.level,
                   a.credit_rate, a.debit_rate
            FROM agents a
            INNER JOIN agent_tree t ON a.id = t.parent_id
        )
        SELECT * FROM agent_tree ORDER BY level ASC
    `, agentID).Scan(&chain).Error

    return chain, err
}
```

### 19.4 分润入账与钱包更新

```go
// internal/app/wallet/service.go

// WalletService 钱包服务
type WalletService struct {
    db    *gorm.DB
    cache *redis.Client
}

// AddBalance 增加钱包余额（使用乐观锁）
func (s *WalletService) AddBalance(
    agentID int64,
    walletType int,
    amount decimal.Decimal,
    remark string,
    refNo string,
) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 获取钱包（加锁）
        var wallet Wallet
        if err := tx.Set("gorm:query_option", "FOR UPDATE").
            Where("agent_id = ? AND wallet_type = ?", agentID, walletType).
            First(&wallet).Error; err != nil {
            return err
        }

        // 2. 更新余额
        balanceBefore := wallet.Balance
        wallet.Balance = wallet.Balance.Add(amount)
        wallet.TotalIncome = wallet.TotalIncome.Add(amount)
        wallet.Version++

        if err := tx.Save(&wallet).Error; err != nil {
            return err
        }

        // 3. 记录流水
        log := WalletLog{
            WalletID:      wallet.ID,
            AgentID:       agentID,
            WalletType:    walletType,
            LogType:       LogTypeIncome,
            Amount:        amount,
            BalanceBefore: balanceBefore,
            BalanceAfter:  wallet.Balance,
            RefType:       "profit",
            RefNo:         refNo,
            Remark:        remark,
        }

        return tx.Create(&log).Error
    })
}
```

### 19.5 定时任务调度

```go
// internal/scheduler/profit_scheduler.go

// ProfitScheduler 分润调度器
type ProfitScheduler struct {
    profitCalculator *TransactionProfitCalculator
    rewardChecker    *ActivationRewardChecker
}

// Run 运行分润定时任务
func (s *ProfitScheduler) Run() {
    // 1. 每小时：处理待计算的交易分润
    c := cron.New()

    c.AddFunc("0 * * * *", func() {
        s.processTransactionProfits()
    })

    // 2. 每日凌晨2点：检查激活奖励
    c.AddFunc("0 2 * * *", func() {
        s.rewardChecker.CheckAndGrant()
    })

    // 3. 每日凌晨3点：更新统计数据
    c.AddFunc("0 3 * * *", func() {
        s.updateDailyStats()
    })

    c.Start()
}

// processTransactionProfits 批量处理交易分润
func (s *ProfitScheduler) processTransactionProfits() {
    // 查询待计算的交易
    transactions, _ := s.transactionRepo.GetPendingProfitTx(1000)

    for _, tx := range transactions {
        results, err := s.profitCalculator.Calculate(tx)
        if err != nil {
            log.Error("profit calculate failed", "tx_id", tx.ID, "error", err)
            continue
        }

        // 保存分润记录并入账
        for _, result := range results {
            s.saveProfitRecord(tx, result)
            s.walletSvc.AddBalance(
                result.AgentID,
                result.WalletType,
                result.ProfitAmount,
                "交易分润",
                tx.OrderNo,
            )
        }

        // 更新交易分润状态
        s.transactionRepo.UpdateProfitStatus(tx.ID, StatusCalculated)
    }
}
```

### 19.6 手动调账功能

```go
// internal/app/profit/adjustment.go

// ManualAdjustment 手动调账
type ManualAdjustment struct {
    ID           int64           `json:"id"`
    AdjustmentNo string          `json:"adjustment_no"`
    AgentID      int64           `json:"agent_id"`
    AdjustType   int             `json:"adjust_type"`   // 1=增加 2=扣减
    WalletType   int             `json:"wallet_type"`   // 1=分润 2=服务费 3=奖励
    Amount       decimal.Decimal `json:"amount"`
    Reason       string          `json:"reason"`
    Status       int             `json:"status"`        // 0=待审核 1=已通过 2=已拒绝
    ApplyUserID  int64           `json:"apply_user_id"`
    AuditUserID  int64           `json:"audit_user_id"`
}

// Execute 执行调账（审核通过后）
func (a *ManualAdjustment) Execute(walletSvc *WalletService) error {
    if a.AdjustType == 1 {
        // 增加余额
        return walletSvc.AddBalance(
            a.AgentID,
            a.WalletType,
            a.Amount,
            "手动调账:"+a.Reason,
            a.AdjustmentNo,
        )
    } else {
        // 扣减余额
        return walletSvc.DeductBalance(
            a.AgentID,
            a.WalletType,
            a.Amount,
            "手动调账:"+a.Reason,
            a.AdjustmentNo,
        )
    }
}
```

### 19.7 分润计算性能优化

#### 批量处理策略

```go
// 批量计算分润（提高性能）
func (c *TransactionProfitCalculator) BatchCalculate(txs []*Transaction) error {
    // 1. 预加载代理商链（减少数据库查询）
    agentIDs := make([]int64, len(txs))
    for i, tx := range txs {
        agentIDs[i] = tx.AgentID
    }
    agentChains := c.preloadAgentChains(agentIDs)

    // 2. 预加载政策模板
    policies := c.preloadPolicies(txs)

    // 3. 批量计算
    var profitRecords []*ProfitRecord
    for _, tx := range txs {
        results := c.calculateWithCache(tx, agentChains, policies)
        profitRecords = append(profitRecords, results...)
    }

    // 4. 批量插入分润记录
    return c.profitRepo.BatchCreate(profitRecords)
}
```

#### 缓存策略

```go
// Redis缓存结构
/*
1. 代理商链缓存
   Key: agent:chain:{agent_id}
   Value: JSON数组
   TTL: 24小时

2. 政策模板缓存
   Key: policy:template:{template_id}
   Value: JSON对象
   TTL: 1小时

3. 通道费率缓存
   Key: channel:rate:{channel_id}
   Value: JSON对象
   TTL: 1小时
*/

func (c *TransactionProfitCalculator) invalidateCache(agentID int64) {
    // 当代理商层级关系变化时，清除相关缓存
    c.cache.Del(ctx, fmt.Sprintf("agent:chain:%d", agentID))

    // 清除所有下级的缓存
    descendants := c.agentRepo.GetDescendants(agentID)
    for _, d := range descendants {
        c.cache.Del(ctx, fmt.Sprintf("agent:chain:%d", d.ID))
    }
}
```

### 19.8 分润对账机制

```go
// internal/app/profit/reconciliation.go

// ProfitReconciliation 分润对账
type ProfitReconciliation struct {
    transactionRepo TransactionRepository
    profitRepo      ProfitRepository
    walletRepo      WalletRepository
}

// DailyReconcile 每日对账
func (r *ProfitReconciliation) DailyReconcile(date time.Time) (*ReconcileResult, error) {
    result := &ReconcileResult{
        Date: date,
    }

    // 1. 统计当日交易总额
    txTotal, _ := r.transactionRepo.SumAmountByDate(date)
    result.TransactionTotal = txTotal

    // 2. 统计当日分润总额
    profitTotal, _ := r.profitRepo.SumAmountByDate(date)
    result.ProfitTotal = profitTotal

    // 3. 统计钱包变动总额
    walletTotal, _ := r.walletRepo.SumIncomeByDate(date)
    result.WalletTotal = walletTotal

    // 4. 检查差异
    if !profitTotal.Equal(walletTotal) {
        result.HasDiff = true
        result.DiffAmount = profitTotal.Sub(walletTotal)
        // 触发告警
        r.alertService.Send("分润对账异常", result)
    }

    return result, nil
}
```

---

## 二十、钱包系统详细设计

### 20.1 三种钱包体系

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              钱包系统架构                                                │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐    │
│  │                           正常钱包（按通道+类型）                                  │    │
│  ├─────────────────────────────────────────────────────────────────────────────────┤    │
│  │                                                                                  │    │
│  │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                        │    │
│  │   │  分润钱包    │    │  服务费钱包  │    │  奖励钱包   │                        │    │
│  │   │  (交易分润)  │    │ (押金+流量)  │    │ (激活奖励)  │                        │    │
│  │   │             │    │             │    │             │                        │    │
│  │   │ 提现门槛:   │    │ 提现门槛:   │    │ 提现门槛:   │                        │    │
│  │   │ 100元      │    │ 200元      │    │ 50元       │                        │    │
│  │   └─────────────┘    └─────────────┘    └─────────────┘                        │    │
│  │                                                                                  │    │
│  │   每个代理商 × 每个通道 = 3种钱包                                                │    │
│  │   例如：代理商A 在 拉卡拉通道 有3个钱包                                           │    │
│  │                                                                                  │    │
│  └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐    │
│  │                           特殊钱包（全局唯一）                                    │    │
│  ├─────────────────────────────────────────────────────────────────────────────────┤    │
│  │                                                                                  │    │
│  │   ┌─────────────────────────┐    ┌─────────────────────────┐                    │    │
│  │   │       充值钱包           │    │       沉淀钱包           │                    │    │
│  │   │   (上级给下级奖励)        │    │  (使用下级未提现比例)    │                    │    │
│  │   │                         │    │                         │                    │    │
│  │   │ • 上级充值后可发放       │    │ • 可使用下级余额30%     │                    │    │
│  │   │ • 满足条件自动发放       │    │ • 支持平台借款          │                    │    │
│  │   │ • 也可用于贷款          │    │ • 隐藏/显示可配置       │                    │    │
│  │   │ • 可隐藏               │    │ • 需承担风险            │                    │    │
│  │   └─────────────────────────┘    └─────────────────────────┘                    │    │
│  │                                                                                  │    │
│  │   每个代理商只有1个充值钱包 + 1个沉淀钱包                                         │    │
│  │                                                                                  │    │
│  └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### 20.2 钱包数据模型

```go
// internal/domain/wallet.go

// 钱包类型常量
const (
    WalletTypeProfit   = 1 // 分润钱包
    WalletTypeService  = 2 // 服务费钱包（押金+流量）
    WalletTypeReward   = 3 // 奖励钱包
    WalletTypeRecharge = 4 // 充值钱包
    WalletTypeDeposit  = 5 // 沉淀钱包
)

// Wallet 钱包模型（正常钱包）
type Wallet struct {
    ID                int64           `json:"id"`
    AgentID           int64           `json:"agent_id"`
    ChannelID         int64           `json:"channel_id"`
    WalletType        int             `json:"wallet_type"`
    Balance           decimal.Decimal `json:"balance"`           // 可用余额
    FrozenAmount      decimal.Decimal `json:"frozen_amount"`     // 冻结金额
    TotalIncome       decimal.Decimal `json:"total_income"`      // 累计收入
    TotalWithdraw     decimal.Decimal `json:"total_withdraw"`    // 累计提现
    WithdrawThreshold decimal.Decimal `json:"withdraw_threshold"` // 提现门槛
    Version           int             `json:"version"`           // 乐观锁
}

// RechargeWallet 充值钱包
type RechargeWallet struct {
    ID            int64           `json:"id"`
    AgentID       int64           `json:"agent_id"`
    IsEnabled     bool            `json:"is_enabled"`      // 是否开启
    IsVisible     bool            `json:"is_visible"`      // 是否可见
    Balance       decimal.Decimal `json:"balance"`         // 余额
    TotalRecharge decimal.Decimal `json:"total_recharge"`  // 累计充值
    TotalPaid     decimal.Decimal `json:"total_paid"`      // 累计发放
    CreditLimit   decimal.Decimal `json:"credit_limit"`    // 授信额度（贷款）
    UsedCredit    decimal.Decimal `json:"used_credit"`     // 已用授信
    Version       int             `json:"version"`
}

// DepositWallet 沉淀钱包
type DepositWallet struct {
    ID                  int64           `json:"id"`
    AgentID             int64           `json:"agent_id"`
    IsEnabled           bool            `json:"is_enabled"`
    IsVisible           bool            `json:"is_visible"`
    AvailableRatio      decimal.Decimal `json:"available_ratio"`      // 可使用比例（如30%）
    SubordinateBalance  decimal.Decimal `json:"subordinate_balance"`  // 下级未提现总额
    AvailableAmount     decimal.Decimal `json:"available_amount"`     // 可使用金额
    UsedAmount          decimal.Decimal `json:"used_amount"`          // 已使用金额
    LoanLimit           decimal.Decimal `json:"loan_limit"`           // 借款额度
    LoanBalance         decimal.Decimal `json:"loan_balance"`         // 借款余额
    Version             int             `json:"version"`
}
```

### 20.3 提现流程

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              提现完整流程                                                │
└─────────────────────────────────────────────────────────────────────────────────────────┘

    ┌──────────────┐
    │  发起提现申请  │
    │  (APP/后台)   │
    └───────┬──────┘
            │
            ▼
    ┌──────────────┐     ┌─────────────┐
    │  校验提现条件  │────▶│  校验失败    │───▶ 返回错误信息
    │  • 余额足够   │     │  • 余额不足  │
    │  • 达到门槛   │     │  • 未达门槛  │
    │  • 银行卡有效 │     └─────────────┘
    └───────┬──────┘
            │ 校验通过
            ▼
    ┌──────────────┐
    │   冻结金额    │ ← 从可用余额转到冻结金额
    │   创建提现单  │
    └───────┬──────┘
            │
            ▼
    ┌──────────────┐     ┌─────────────┐
    │   审核流程    │────▶│  审核拒绝    │───▶ 解冻金额，退回余额
    │   (人工审核)  │     └─────────────┘
    └───────┬──────┘
            │ 审核通过
            ▼
    ┌──────────────┐
    │  计算实际到账  │
    │  扣除税筹费用  │
    │  9% + 3元/笔  │
    └───────┬──────┘
            │
            ▼
    ┌──────────────┐
    │  调用税筹通道  │ ← 代付打款
    │  打款到银行卡  │
    └───────┬──────┘
            │
    ┌───────┴───────┐
    ▼               ▼
 ┌─────┐        ┌─────┐
 │成功  │        │失败  │───▶ 解冻金额，重试/人工处理
 └──┬──┘        └─────┘
    │
    ▼
 ┌──────────────┐
 │  扣减冻结金额  │
 │  更新提现单   │
 │  状态=完成    │
 └──────────────┘
```

### 20.4 提现代码实现

```go
// internal/app/wallet/withdraw_service.go

type WithdrawService struct {
    walletRepo   WalletRepository
    withdrawRepo WithdrawRepository
    taxChannel   TaxChannelAdapter
}

// ApplyWithdraw 申请提现
func (s *WithdrawService) ApplyWithdraw(req *WithdrawRequest) (*Withdrawal, error) {
    // 1. 获取钱包
    wallet, err := s.walletRepo.GetByAgentAndType(req.AgentID, req.ChannelID, req.WalletType)
    if err != nil {
        return nil, err
    }

    // 2. 校验余额
    if wallet.Balance.LessThan(req.Amount) {
        return nil, errors.New("余额不足")
    }

    // 3. 校验提现门槛
    threshold := s.getWithdrawThreshold(req.ChannelID, req.WalletType)
    if req.Amount.LessThan(threshold) {
        return nil, fmt.Errorf("提现金额需达到%.2f元", threshold.InexactFloat64())
    }

    // 4. 计算税筹费用
    taxConfig := s.getTaxChannelConfig(req.ChannelID)
    taxFee := req.Amount.Mul(taxConfig.TaxRate)           // 9%
    serviceFee := taxConfig.FixedFee                       // 3元
    actualAmount := req.Amount.Sub(taxFee).Sub(serviceFee) // 实际到账

    // 5. 冻结金额
    err = s.walletRepo.FreezeBalance(wallet.ID, req.Amount)
    if err != nil {
        return nil, err
    }

    // 6. 创建提现单
    withdrawal := &Withdrawal{
        WithdrawNo:   generateWithdrawNo(),
        AgentID:      req.AgentID,
        WalletID:     wallet.ID,
        WalletType:   req.WalletType,
        ChannelID:    req.ChannelID,
        Amount:       req.Amount,
        TaxFee:       taxFee.Add(serviceFee),
        ActualAmount: actualAmount,
        BankCardNo:   req.BankCardNo,
        BankName:     req.BankName,
        AccountName:  req.AccountName,
        Status:       WithdrawStatusPending,
    }

    return s.withdrawRepo.Create(withdrawal)
}

// ProcessWithdraw 处理提现（审核通过后）
func (s *WithdrawService) ProcessWithdraw(withdrawID int64) error {
    withdrawal, _ := s.withdrawRepo.GetByID(withdrawID)

    // 1. 调用税筹通道打款
    result, err := s.taxChannel.Transfer(&TransferRequest{
        OrderNo:     withdrawal.WithdrawNo,
        Amount:      withdrawal.ActualAmount,
        BankCardNo:  withdrawal.BankCardNo,
        BankName:    withdrawal.BankName,
        AccountName: withdrawal.AccountName,
    })

    if err != nil || !result.Success {
        // 打款失败，更新状态
        s.withdrawRepo.UpdateStatus(withdrawID, WithdrawStatusFailed)
        return err
    }

    // 2. 打款成功，扣减冻结金额
    s.walletRepo.DeductFrozenBalance(withdrawal.WalletID, withdrawal.Amount)

    // 3. 更新提现单状态
    s.withdrawRepo.UpdateStatus(withdrawID, WithdrawStatusSuccess)
    s.withdrawRepo.UpdatePayInfo(withdrawID, result.PayTransNo, time.Now())

    return nil
}
```

### 20.5 充值钱包业务逻辑

```go
// internal/app/wallet/recharge_wallet_service.go

type RechargeWalletService struct {
    rechargeRepo RechargeWalletRepository
    walletRepo   WalletRepository
    policyRepo   PolicyRepository
}

// Recharge 上级给充值钱包充值
func (s *RechargeWalletService) Recharge(agentID int64, amount decimal.Decimal) error {
    wallet, _ := s.rechargeRepo.GetByAgentID(agentID)
    if wallet == nil {
        // 首次充值，创建钱包
        wallet = &RechargeWallet{
            AgentID:   agentID,
            IsEnabled: true,
            IsVisible: true,
        }
        s.rechargeRepo.Create(wallet)
    }

    // 增加余额
    wallet.Balance = wallet.Balance.Add(amount)
    wallet.TotalRecharge = wallet.TotalRecharge.Add(amount)

    return s.rechargeRepo.Update(wallet)
}

// TransferToSubordinate 发放奖励给下级
func (s *RechargeWalletService) TransferToSubordinate(
    fromAgentID int64,
    toAgentID int64,
    amount decimal.Decimal,
    reason string,
) error {
    // 1. 检查充值钱包余额
    wallet, _ := s.rechargeRepo.GetByAgentID(fromAgentID)
    if wallet.Balance.LessThan(amount) {
        return errors.New("充值钱包余额不足")
    }

    // 2. 扣减上级充值钱包
    wallet.Balance = wallet.Balance.Sub(amount)
    wallet.TotalPaid = wallet.TotalPaid.Add(amount)
    s.rechargeRepo.Update(wallet)

    // 3. 增加下级奖励钱包
    // 下级收到的奖励入奖励钱包
    s.walletRepo.AddBalance(toAgentID, WalletTypeReward, amount, reason, "")

    // 4. 记录转账流水
    s.recordTransfer(fromAgentID, toAgentID, amount, reason)

    return nil
}

// AutoGrantRewards 自动发放奖励（定时任务）
func (s *RechargeWalletService) AutoGrantRewards() error {
    // 1. 查询所有待发放的奖励规则
    rules, _ := s.policyRepo.GetPendingRechargeRewards()

    for _, rule := range rules {
        // 2. 检查下级是否满足条件
        subordinates := s.getQualifiedSubordinates(rule)

        for _, sub := range subordinates {
            // 3. 发放奖励
            s.TransferToSubordinate(
                rule.AgentID,
                sub.ID,
                rule.RewardAmount,
                rule.RewardName,
            )
        }
    }

    return nil
}
```

### 20.6 沉淀钱包业务逻辑

```go
// internal/app/wallet/deposit_wallet_service.go

type DepositWalletService struct {
    depositRepo DepositWalletRepository
    walletRepo  WalletRepository
    agentRepo   AgentRepository
}

// CalculateSubordinateBalance 计算下级未提现余额（定时任务）
func (s *DepositWalletService) CalculateSubordinateBalance(agentID int64) error {
    // 1. 获取所有直属下级
    subordinates, _ := s.agentRepo.GetDirectSubordinates(agentID)

    totalBalance := decimal.Zero
    for _, sub := range subordinates {
        // 2. 汇总下级所有钱包余额
        wallets, _ := s.walletRepo.GetByAgentID(sub.ID)
        for _, w := range wallets {
            totalBalance = totalBalance.Add(w.Balance)
        }
    }

    // 3. 更新沉淀钱包可用金额
    depositWallet, _ := s.depositRepo.GetByAgentID(agentID)
    if depositWallet == nil {
        return nil
    }

    depositWallet.SubordinateBalance = totalBalance
    depositWallet.AvailableAmount = totalBalance.Mul(depositWallet.AvailableRatio).Div(decimal.NewFromInt(100))

    return s.depositRepo.Update(depositWallet)
}

// UseSubordinateBalance 使用下级余额
func (s *DepositWalletService) UseSubordinateBalance(agentID int64, amount decimal.Decimal) error {
    depositWallet, _ := s.depositRepo.GetByAgentID(agentID)

    // 1. 检查可用金额
    remaining := depositWallet.AvailableAmount.Sub(depositWallet.UsedAmount)
    if remaining.LessThan(amount) {
        return errors.New("沉淀钱包可用金额不足")
    }

    // 2. 增加已使用金额
    depositWallet.UsedAmount = depositWallet.UsedAmount.Add(amount)

    // 3. 记录流水
    s.recordDepositUse(agentID, amount)

    return s.depositRepo.Update(depositWallet)
}

// ApplyLoan 申请平台借款
func (s *DepositWalletService) ApplyLoan(req *LoanRequest) (*PlatformLoan, error) {
    depositWallet, _ := s.depositRepo.GetByAgentID(req.AgentID)

    // 1. 检查借款额度
    if req.Amount.GreaterThan(depositWallet.LoanLimit) {
        return nil, errors.New("超出借款额度")
    }

    // 2. 创建借款记录（需线下签协议）
    loan := &PlatformLoan{
        LoanNo:          generateLoanNo(),
        AgentID:         req.AgentID,
        LoanAmount:      req.Amount,
        LoanRate:        req.Rate,
        InterestAmount:  req.Amount.Mul(req.Rate),
        RemainingAmount: req.Amount,
        Status:          LoanStatusPending,
        AgreementNo:     req.AgreementNo,
    }

    return s.loanRepo.Create(loan)
}
```

### 20.7 提现门槛配置

```sql
-- 提现门槛配置表
CREATE TABLE withdraw_thresholds (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,
    threshold       DECIMAL(10,2) NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(channel_id, wallet_type)
);

-- 示例数据
INSERT INTO withdraw_thresholds (channel_id, wallet_type, threshold) VALUES
(1, 1, 100.00),   -- 拉卡拉-分润钱包-100元
(1, 2, 200.00),   -- 拉卡拉-服务费钱包-200元
(1, 3, 50.00),    -- 拉卡拉-奖励钱包-50元
(2, 1, 150.00),   -- 随行付-分润钱包-150元
(2, 2, 250.00);   -- 随行付-服务费钱包-250元
```

---

## 二十一、通道对接设计

### 21.1 通道适配器架构

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              通道适配器架构                                              │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                          │
│  ┌──────────────────────────────────────────────────────────────────────────────┐       │
│  │                          ChannelAdapter 接口                                  │       │
│  │  统一定义所有通道必须实现的方法                                                 │       │
│  └────────────────────────────────────┬─────────────────────────────────────────┘       │
│                                       │                                                  │
│          ┌────────────────────────────┼────────────────────────────┐                    │
│          │                            │                            │                    │
│          ▼                            ▼                            ▼                    │
│  ┌───────────────┐           ┌───────────────┐           ┌───────────────┐             │
│  │  LakalaAdapter │           │  SuixingfuAdapter │           │  XxxPayAdapter │             │
│  │    拉卡拉       │           │    随行付        │           │    其他通道     │             │
│  └───────────────┘           └───────────────┘           └───────────────┘             │
│          │                            │                            │                    │
│          ▼                            ▼                            ▼                    │
│  ┌───────────────┐           ┌───────────────┐           ┌───────────────┐             │
│  │  HTTP Client   │           │  HTTP Client   │           │  HTTP Client   │             │
│  │  签名/验签     │           │  签名/验签     │           │  签名/验签     │             │
│  │  加密/解密     │           │  加密/解密     │           │  加密/解密     │             │
│  └───────────────┘           └───────────────┘           └───────────────┘             │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### 21.2 通道适配器接口定义

```go
// internal/adapters/channel/interface.go

// ChannelAdapter 通道适配器接口
type ChannelAdapter interface {
    // 基础信息
    GetChannelCode() string
    GetChannelName() string

    // ============ 商户管理 ============
    // 商户入网
    RegisterMerchant(req *MerchantRegisterReq) (*MerchantRegisterResp, error)
    // 商户信息变更
    UpdateMerchant(merchantNo string, req *MerchantUpdateReq) error
    // 商户费率修改（实时生效）
    UpdateMerchantRate(merchantNo string, rateConfig *RateConfig) error
    // 商户状态查询
    QueryMerchantStatus(merchantNo string) (*MerchantStatus, error)

    // ============ 终端管理 ============
    // 终端绑定
    BindTerminal(sn string, merchantNo string) error
    // 终端解绑
    UnbindTerminal(sn string) error
    // 下发终端政策（费率/押金/流量卡）
    PushTerminalPolicy(sn string, policy *TerminalPolicy) error
    // 终端状态查询
    QueryTerminalStatus(sn string) (*TerminalStatus, error)

    // ============ 交易查询 ============
    // 单笔交易查询
    QueryTransaction(tradeNo string) (*TransactionDetail, error)
    // 批量交易查询（日期范围）
    QueryTransactionList(startDate, endDate time.Time, page int) (*TransactionListResp, error)
    // 交易对账文件下载
    DownloadReconciliation(date time.Time) ([]byte, error)

    // ============ 回调处理 ============
    // 解析回调数据
    ParseCallback(data []byte) (*CallbackData, error)
    // 验证签名
    VerifySign(data []byte, sign string) bool
    // 生成签名
    GenerateSign(data map[string]string) string
}

// RateConfig 费率配置
type RateConfig struct {
    CreditRate  decimal.Decimal `json:"credit_rate"`   // 贷记卡费率
    DebitRate   decimal.Decimal `json:"debit_rate"`    // 借记卡费率
    DebitCap    decimal.Decimal `json:"debit_cap"`     // 借记卡封顶
    CloudRate   decimal.Decimal `json:"cloud_rate"`    // 云闪付费率
    WechatRate  decimal.Decimal `json:"wechat_rate"`   // 微信费率
    AlipayRate  decimal.Decimal `json:"alipay_rate"`   // 支付宝费率
    T0FeeRate   decimal.Decimal `json:"t0_fee_rate"`   // 秒到费率
}

// TerminalPolicy 终端政策
type TerminalPolicy struct {
    RateConfig    *RateConfig     `json:"rate_config"`
    DepositAmount decimal.Decimal `json:"deposit_amount"`  // 押金金额
    SimFirstFee   decimal.Decimal `json:"sim_first_fee"`   // 首次流量费
    SimRenewFee   decimal.Decimal `json:"sim_renew_fee"`   // 续费金额
    SimInterval   int             `json:"sim_interval"`    // 扣费间隔天数
}

// CallbackData 回调数据
type CallbackData struct {
    TradeNo       string          `json:"trade_no"`
    OrderNo       string          `json:"order_no"`
    MerchantNo    string          `json:"merchant_no"`
    TerminalSN    string          `json:"terminal_sn"`
    TradeType     int             `json:"trade_type"`      // 1消费 2撤销 3退货 5押金 6流量
    PayType       int             `json:"pay_type"`        // 1刷卡 2微信 3支付宝 4云闪付
    CardType      int             `json:"card_type"`       // 1借记 2贷记
    Amount        decimal.Decimal `json:"amount"`
    Fee           decimal.Decimal `json:"fee"`
    Rate          decimal.Decimal `json:"rate"`
    TradeTime     time.Time       `json:"trade_time"`
    Status        int             `json:"status"`
    RawData       string          `json:"raw_data"`
}
```

### 21.3 拉卡拉通道实现示例

```go
// internal/adapters/channel/lakala/adapter.go

type LakalaAdapter struct {
    config     *LakalaConfig
    httpClient *http.Client
}

type LakalaConfig struct {
    ApiUrl     string
    MerchantNo string
    AppID      string
    PrivateKey string
    PublicKey  string
}

func NewLakalaAdapter(config *LakalaConfig) *LakalaAdapter {
    return &LakalaAdapter{
        config:     config,
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

func (a *LakalaAdapter) GetChannelCode() string {
    return "LAKALA"
}

func (a *LakalaAdapter) GetChannelName() string {
    return "拉卡拉"
}

// RegisterMerchant 商户入网
func (a *LakalaAdapter) RegisterMerchant(req *MerchantRegisterReq) (*MerchantRegisterResp, error) {
    // 1. 构建请求参数
    params := map[string]string{
        "merchant_name": req.MerchantName,
        "contact_phone": req.ContactPhone,
        "id_card_no":    req.IDCardNo,
        "bank_card_no":  req.BankCardNo,
        "mcc_code":      req.MCCCode,
        "timestamp":     time.Now().Format("20060102150405"),
    }

    // 2. 签名
    sign := a.GenerateSign(params)
    params["sign"] = sign

    // 3. 发送请求
    resp, err := a.post("/merchant/register", params)
    if err != nil {
        return nil, err
    }

    // 4. 解析响应
    var result struct {
        Code       string `json:"code"`
        Message    string `json:"message"`
        MerchantNo string `json:"merchant_no"`
    }
    json.Unmarshal(resp, &result)

    if result.Code != "0000" {
        return nil, errors.New(result.Message)
    }

    return &MerchantRegisterResp{
        MerchantNo: result.MerchantNo,
    }, nil
}

// UpdateMerchantRate 修改商户费率（实时生效）
func (a *LakalaAdapter) UpdateMerchantRate(merchantNo string, rateConfig *RateConfig) error {
    params := map[string]string{
        "merchant_no":  merchantNo,
        "credit_rate":  rateConfig.CreditRate.String(),
        "debit_rate":   rateConfig.DebitRate.String(),
        "debit_cap":    rateConfig.DebitCap.String(),
        "timestamp":    time.Now().Format("20060102150405"),
    }

    sign := a.GenerateSign(params)
    params["sign"] = sign

    resp, err := a.post("/merchant/rate/update", params)
    if err != nil {
        return err
    }

    var result struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    }
    json.Unmarshal(resp, &result)

    if result.Code != "0000" {
        return errors.New(result.Message)
    }

    return nil
}

// PushTerminalPolicy 下发终端政策
func (a *LakalaAdapter) PushTerminalPolicy(sn string, policy *TerminalPolicy) error {
    params := map[string]string{
        "terminal_sn":    sn,
        "credit_rate":    policy.RateConfig.CreditRate.String(),
        "debit_rate":     policy.RateConfig.DebitRate.String(),
        "deposit_amount": policy.DepositAmount.String(),
        "sim_first_fee":  policy.SimFirstFee.String(),
        "sim_renew_fee":  policy.SimRenewFee.String(),
        "sim_interval":   strconv.Itoa(policy.SimInterval),
        "timestamp":      time.Now().Format("20060102150405"),
    }

    sign := a.GenerateSign(params)
    params["sign"] = sign

    resp, err := a.post("/terminal/policy/push", params)
    if err != nil {
        return err
    }

    // 解析响应...
    return nil
}

// ParseCallback 解析回调
func (a *LakalaAdapter) ParseCallback(data []byte) (*CallbackData, error) {
    var raw map[string]interface{}
    json.Unmarshal(data, &raw)

    // 验证签名
    sign := raw["sign"].(string)
    if !a.VerifySign(data, sign) {
        return nil, errors.New("签名验证失败")
    }

    // 解析数据
    callback := &CallbackData{
        TradeNo:    raw["trade_no"].(string),
        OrderNo:    raw["order_no"].(string),
        MerchantNo: raw["merchant_no"].(string),
        TerminalSN: raw["terminal_sn"].(string),
        TradeType:  int(raw["trade_type"].(float64)),
        PayType:    int(raw["pay_type"].(float64)),
        RawData:    string(data),
    }

    amountStr := raw["amount"].(string)
    callback.Amount, _ = decimal.NewFromString(amountStr)

    tradeTimeStr := raw["trade_time"].(string)
    callback.TradeTime, _ = time.Parse("20060102150405", tradeTimeStr)

    return callback, nil
}

// GenerateSign RSA签名
func (a *LakalaAdapter) GenerateSign(params map[string]string) string {
    // 1. 参数排序拼接
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    var builder strings.Builder
    for _, k := range keys {
        builder.WriteString(k)
        builder.WriteString("=")
        builder.WriteString(params[k])
        builder.WriteString("&")
    }
    signStr := strings.TrimSuffix(builder.String(), "&")

    // 2. RSA签名
    hash := sha256.Sum256([]byte(signStr))
    privateKey, _ := parsePrivateKey(a.config.PrivateKey)
    signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])

    return base64.StdEncoding.EncodeToString(signature)
}

// post 发送POST请求
func (a *LakalaAdapter) post(path string, params map[string]string) ([]byte, error) {
    jsonData, _ := json.Marshal(params)
    req, _ := http.NewRequest("POST", a.config.ApiUrl+path, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    resp, err := a.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### 21.4 通道工厂模式

```go
// internal/adapters/channel/factory.go

type ChannelAdapterFactory struct {
    adapters map[string]ChannelAdapter
    configs  map[string]interface{}
}

func NewChannelAdapterFactory() *ChannelAdapterFactory {
    return &ChannelAdapterFactory{
        adapters: make(map[string]ChannelAdapter),
        configs:  make(map[string]interface{}),
    }
}

// Register 注册通道适配器
func (f *ChannelAdapterFactory) Register(code string, config interface{}) {
    f.configs[code] = config
}

// GetAdapter 获取通道适配器
func (f *ChannelAdapterFactory) GetAdapter(channelCode string) (ChannelAdapter, error) {
    // 1. 尝试从缓存获取
    if adapter, ok := f.adapters[channelCode]; ok {
        return adapter, nil
    }

    // 2. 根据编码创建适配器
    config := f.configs[channelCode]
    var adapter ChannelAdapter

    switch channelCode {
    case "LAKALA":
        adapter = NewLakalaAdapter(config.(*LakalaConfig))
    case "SUIXINGFU":
        adapter = NewSuixingfuAdapter(config.(*SuixingfuConfig))
    case "YEEPAY":
        adapter = NewYeepayAdapter(config.(*YeepayConfig))
    // ... 其他通道
    default:
        return nil, fmt.Errorf("unsupported channel: %s", channelCode)
    }

    f.adapters[channelCode] = adapter
    return adapter, nil
}

// 使用示例
func Example() {
    factory := NewChannelAdapterFactory()

    // 注册通道配置
    factory.Register("LAKALA", &LakalaConfig{
        ApiUrl:     "https://api.lakala.com",
        MerchantNo: "xxx",
        AppID:      "xxx",
    })

    // 获取适配器
    adapter, _ := factory.GetAdapter("LAKALA")

    // 使用适配器
    adapter.RegisterMerchant(&MerchantRegisterReq{...})
}
```

### 21.5 交易数据同步服务

```go
// internal/app/transaction/sync_service.go

type TransactionSyncService struct {
    channelFactory *ChannelAdapterFactory
    transactionRepo TransactionRepository
    profitService   *ProfitService
}

// SyncTransactions 同步交易数据
func (s *TransactionSyncService) SyncTransactions(channelID int64, date time.Time) error {
    // 1. 获取通道适配器
    channel, _ := s.channelRepo.GetByID(channelID)
    adapter, _ := s.channelFactory.GetAdapter(channel.ChannelCode)

    // 2. 查询交易列表
    page := 1
    for {
        resp, err := adapter.QueryTransactionList(date, date, page)
        if err != nil {
            return err
        }

        // 3. 批量保存交易
        for _, tx := range resp.Transactions {
            // 检查是否已存在
            exists, _ := s.transactionRepo.ExistsByTradeNo(tx.TradeNo)
            if exists {
                continue
            }

            // 保存交易
            transaction := &Transaction{
                TradeNo:     tx.TradeNo,
                ChannelID:   channelID,
                MerchantNo:  tx.MerchantNo,
                TerminalSN:  tx.TerminalSN,
                TradeType:   tx.TradeType,
                PayType:     tx.PayType,
                Amount:      tx.Amount,
                Fee:         tx.Fee,
                Rate:        tx.Rate,
                TradeTime:   tx.TradeTime,
                Status:      tx.Status,
                ProfitStatus: ProfitStatusPending,
            }
            s.transactionRepo.Create(transaction)
        }

        // 4. 判断是否还有下一页
        if page >= resp.TotalPages {
            break
        }
        page++
    }

    return nil
}

// HandleCallback 处理回调
func (s *TransactionSyncService) HandleCallback(channelCode string, data []byte) error {
    // 1. 获取适配器
    adapter, _ := s.channelFactory.GetAdapter(channelCode)

    // 2. 解析回调
    callback, err := adapter.ParseCallback(data)
    if err != nil {
        return err
    }

    // 3. 保存交易
    transaction := &Transaction{
        TradeNo:     callback.TradeNo,
        ChannelID:   s.getChannelID(channelCode),
        MerchantNo:  callback.MerchantNo,
        TerminalSN:  callback.TerminalSN,
        TradeType:   callback.TradeType,
        PayType:     callback.PayType,
        CardType:    callback.CardType,
        Amount:      callback.Amount,
        Fee:         callback.Fee,
        Rate:        callback.Rate,
        TradeTime:   callback.TradeTime,
        Status:      callback.Status,
        ProfitStatus: ProfitStatusPending,
    }

    // 4. 关联商户和代理商
    merchant, _ := s.merchantRepo.GetByChannelNo(callback.MerchantNo)
    if merchant != nil {
        transaction.MerchantID = merchant.ID
        transaction.AgentID = merchant.AgentID
    }

    s.transactionRepo.Create(transaction)

    // 5. 触发分润计算（异步）
    go s.profitService.TriggerCalculate(transaction.ID)

    return nil
}
```

### 21.6 通道配置管理

```sql
-- 通道扩展配置表
CREATE TABLE channel_configs (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    config_key      VARCHAR(100) NOT NULL,
    config_value    TEXT NOT NULL,
    description     VARCHAR(500),
    UNIQUE(channel_id, config_key)
);

-- 通道配置示例
INSERT INTO channel_configs (channel_id, config_key, config_value, description) VALUES
-- 拉卡拉配置
(1, 'api_url', 'https://api.lakala.com', 'API基础地址'),
(1, 'merchant_no', 'M123456789', '商户号'),
(1, 'app_id', 'APP123456', '应用ID'),
(1, 'sign_type', 'RSA2', '签名类型'),
(1, 'private_key', '-----BEGIN RSA PRIVATE KEY-----...', '私钥'),
(1, 'public_key', '-----BEGIN PUBLIC KEY-----...', '公钥'),
(1, 'callback_url', 'https://ourplatform.com/callback/lakala', '回调地址'),
(1, 'rate_update_api', '/merchant/rate/update', '费率修改接口'),
(1, 'rate_update_realtime', 'true', '费率是否实时生效'),

-- 随行付配置
(2, 'api_url', 'https://api.suixingpay.com', 'API基础地址'),
(2, 'sign_type', 'MD5', '签名类型');
```

---

## 二十二、政策模板引擎详细设计

### 22.1 政策模板体系

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              政策模板体系架构                                            │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐    │
│  │                           政策模板主表                                           │    │
│  │  policy_templates                                                               │    │
│  │  ├── 模板名称、所属通道、是否默认                                                │    │
│  │  └── 基础费率（贷记卡、借记卡、云闪付、微信、支付宝）                              │    │
│  └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                       │                                                  │
│          ┌────────────────────────────┼────────────────────────────┐                    │
│          ▼                            ▼                            ▼                    │
│  ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐               │
│  │  费率阶梯规则      │    │  激活奖励规则      │    │  返现规则          │               │
│  │  policy_rate_stages│    │  policy_activation │    │                   │               │
│  │                   │    │  _rewards          │    │  ┌─────────────┐  │               │
│  │  • 阶段天数范围   │    │                   │    │  │ 押金返现    │  │               │
│  │  • 费率调整值    │    │  • 考核天数范围   │    │  │ policy_     │  │               │
│  │  • 计算基准      │    │  • 达标交易额    │    │  │ deposit_    │  │               │
│  │    (商户/代理)   │    │  • 奖励金额      │    │  │ cashbacks   │  │               │
│  │                   │    │  • 多层级奖励    │    │  └─────────────┘  │               │
│  └───────────────────┘    │  • 生效时间      │    │  ┌─────────────┐  │               │
│                           └───────────────────┘    │  │ 流量返现    │  │               │
│                                                    │  │ policy_     │  │               │
│                                                    │  │ sim_        │  │               │
│                                                    │  │ cashbacks   │  │               │
│                                                    │  └─────────────┘  │               │
│                                                    └───────────────────┘               │
│                                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐    │
│  │                           代理商-政策关联表                                       │    │
│  │  agent_policies                                                                 │    │
│  │  ├── 代理商ID、通道ID、模板ID                                                    │    │
│  │  └── 可覆盖的自定义费率（继承模板 or 自定义覆盖）                                  │    │
│  └─────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

### 22.2 政策模板数据模型

```go
// internal/domain/policy.go

// PolicyTemplate 政策模板
type PolicyTemplate struct {
    ID           int64           `json:"id"`
    TemplateName string          `json:"template_name"`
    ChannelID    int64           `json:"channel_id"`
    IsDefault    bool            `json:"is_default"`

    // 基础费率（结算价）
    CreditRate   decimal.Decimal `json:"credit_rate"`    // 贷记卡
    DebitRate    decimal.Decimal `json:"debit_rate"`     // 借记卡
    DebitCap     decimal.Decimal `json:"debit_cap"`      // 借记卡封顶
    UnionpayRate decimal.Decimal `json:"unionpay_rate"`  // 银联云闪付
    WechatRate   decimal.Decimal `json:"wechat_rate"`    // 微信扫码
    AlipayRate   decimal.Decimal `json:"alipay_rate"`    // 支付宝扫码
    T0FeeType    int             `json:"t0_fee_type"`    // 秒到费率档位

    Status       int             `json:"status"`
    CreatedAt    time.Time       `json:"created_at"`
}

// PolicyRateStage 费率阶梯规则
type PolicyRateStage struct {
    ID          int64           `json:"id"`
    TemplateID  int64           `json:"template_id"`
    StageType   int             `json:"stage_type"`    // 1=按商户入网时间 2=按代理商入网时间
    DayStart    int             `json:"day_start"`     // 开始天数
    DayEnd      *int            `json:"day_end"`       // 结束天数 nil=无限
    RateAdjust  decimal.Decimal `json:"rate_adjust"`   // 费率调整值
    CreatedAt   time.Time       `json:"created_at"`
}

// PolicyActivationReward 激活奖励规则
type PolicyActivationReward struct {
    ID            int64           `json:"id"`
    TemplateID    int64           `json:"template_id"`
    DayStart      int             `json:"day_start"`      // 激活后开始天数
    DayEnd        int             `json:"day_end"`        // 激活后结束天数
    TradeAmount   decimal.Decimal `json:"trade_amount"`   // 达标交易额
    RewardAmount  decimal.Decimal `json:"reward_amount"`  // 奖励金额
    LayerRewards  string          `json:"layer_rewards"`  // 多层级奖励JSON
    EffectiveDate *time.Time      `json:"effective_date"` // 生效日期
    CreatedAt     time.Time       `json:"created_at"`
}

// PolicyDepositCashback 押金返现规则
type PolicyDepositCashback struct {
    ID              int64           `json:"id"`
    TemplateID      int64           `json:"template_id"`
    DepositAmount   decimal.Decimal `json:"deposit_amount"`   // 押金金额
    CashbackAmount  decimal.Decimal `json:"cashback_amount"`  // 返现金额
    CreatedAt       time.Time       `json:"created_at"`
}

// PolicySimCashback 流量费返现规则
type PolicySimCashback struct {
    ID              int64           `json:"id"`
    TemplateID      int64           `json:"template_id"`
    SimType         int             `json:"sim_type"`         // 1=首次 2=续费
    SimFee          decimal.Decimal `json:"sim_fee"`          // 流量费金额
    CashbackAmount  decimal.Decimal `json:"cashback_amount"`  // 返现金额
    LayerCashback   string          `json:"layer_cashback"`   // 多层级返现JSON
    CreatedAt       time.Time       `json:"created_at"`
}

// AgentPolicy 代理商政策配置
type AgentPolicy struct {
    ID          int64           `json:"id"`
    AgentID     int64           `json:"agent_id"`
    ChannelID   int64           `json:"channel_id"`
    TemplateID  int64           `json:"template_id"`

    // 自定义覆盖费率（nil表示继承模板）
    CreditRate  *decimal.Decimal `json:"credit_rate"`
    DebitRate   *decimal.Decimal `json:"debit_rate"`
    DebitCap    *decimal.Decimal `json:"debit_cap"`

    CreatedAt   time.Time       `json:"created_at"`
}
```

### 22.3 政策引擎核心服务

```go
// internal/app/policy/engine.go

// PolicyEngine 政策引擎
type PolicyEngine struct {
    templateRepo     PolicyTemplateRepository
    rateStageRepo    PolicyRateStageRepository
    rewardRepo       PolicyActivationRewardRepository
    depositRepo      PolicyDepositCashbackRepository
    simRepo          PolicySimCashbackRepository
    agentPolicyRepo  AgentPolicyRepository
    cache            *redis.Client
}

// GetEffectiveRate 获取生效的结算费率
func (e *PolicyEngine) GetEffectiveRate(
    agentID int64,
    channelID int64,
    payType string,
    merchantRegisterTime time.Time,
    agentRegisterTime time.Time,
) (decimal.Decimal, error) {
    // 1. 获取代理商政策
    agentPolicy, err := e.agentPolicyRepo.GetByAgentAndChannel(agentID, channelID)
    if err != nil {
        return decimal.Zero, err
    }

    // 2. 获取基础费率
    baseRate := e.getBaseRate(agentPolicy, payType)

    // 3. 获取阶梯调整
    template, _ := e.templateRepo.GetByID(agentPolicy.TemplateID)
    stages, _ := e.rateStageRepo.GetByTemplateID(template.ID)

    if len(stages) == 0 {
        return baseRate, nil
    }

    // 4. 计算生效的阶梯费率
    for _, stage := range stages {
        var days int
        if stage.StageType == 1 {
            // 按商户入网时间
            days = int(time.Since(merchantRegisterTime).Hours() / 24)
        } else {
            // 按代理商入网时间
            days = int(time.Since(agentRegisterTime).Hours() / 24)
        }

        // 匹配阶段
        if days >= stage.DayStart && (stage.DayEnd == nil || days <= *stage.DayEnd) {
            return baseRate.Add(stage.RateAdjust), nil
        }
    }

    return baseRate, nil
}

// getBaseRate 获取基础费率
func (e *PolicyEngine) getBaseRate(policy *AgentPolicy, payType string) decimal.Decimal {
    // 优先使用自定义覆盖费率
    template, _ := e.templateRepo.GetByID(policy.TemplateID)

    switch payType {
    case "CREDIT":
        if policy.CreditRate != nil {
            return *policy.CreditRate
        }
        return template.CreditRate
    case "DEBIT":
        if policy.DebitRate != nil {
            return *policy.DebitRate
        }
        return template.DebitRate
    case "UNIONPAY":
        return template.UnionpayRate
    case "WECHAT":
        return template.WechatRate
    case "ALIPAY":
        return template.AlipayRate
    default:
        return template.CreditRate
    }
}

// GetActivationRewardRules 获取激活奖励规则
func (e *PolicyEngine) GetActivationRewardRules(
    agentID int64,
    channelID int64,
    activateTime time.Time,
) ([]*PolicyActivationReward, error) {
    agentPolicy, _ := e.agentPolicyRepo.GetByAgentAndChannel(agentID, channelID)
    if agentPolicy == nil {
        return nil, nil
    }

    // 获取模板下的所有奖励规则
    rules, err := e.rewardRepo.GetByTemplateID(agentPolicy.TemplateID)
    if err != nil {
        return nil, err
    }

    // 过滤生效的规则
    var effectiveRules []*PolicyActivationReward
    for _, rule := range rules {
        // 检查生效日期（针对新商户）
        if rule.EffectiveDate != nil && activateTime.Before(*rule.EffectiveDate) {
            continue
        }
        effectiveRules = append(effectiveRules, rule)
    }

    return effectiveRules, nil
}

// GetDepositCashbackRule 获取押金返现规则
func (e *PolicyEngine) GetDepositCashbackRule(
    agentID int64,
    channelID int64,
    depositAmount decimal.Decimal,
) (*PolicyDepositCashback, error) {
    agentPolicy, _ := e.agentPolicyRepo.GetByAgentAndChannel(agentID, channelID)
    if agentPolicy == nil {
        return nil, nil
    }

    return e.depositRepo.GetByTemplateAndAmount(agentPolicy.TemplateID, depositAmount)
}

// GetSimCashbackRules 获取流量费返现规则
func (e *PolicyEngine) GetSimCashbackRules(
    agentID int64,
    channelID int64,
    isFirstYear bool,
) ([]*PolicySimCashback, error) {
    agentPolicy, _ := e.agentPolicyRepo.GetByAgentAndChannel(agentID, channelID)
    if agentPolicy == nil {
        return nil, nil
    }

    simType := 1 // 首次
    if !isFirstYear {
        simType = 2 // 续费
    }

    return e.simRepo.GetByTemplateAndType(agentPolicy.TemplateID, simType)
}
```

### 22.4 政策继承与覆盖

```go
// internal/app/policy/inherit.go

// PolicyInheritService 政策继承服务
type PolicyInheritService struct {
    agentRepo       AgentRepository
    agentPolicyRepo AgentPolicyRepository
    templateRepo    PolicyTemplateRepository
}

// InheritFromParent 从上级继承政策
func (s *PolicyInheritService) InheritFromParent(agentID int64, parentID int64) error {
    // 1. 获取上级的所有通道政策
    parentPolicies, err := s.agentPolicyRepo.GetByAgentID(parentID)
    if err != nil {
        return err
    }

    // 2. 为新代理商创建继承的政策
    for _, parentPolicy := range parentPolicies {
        // 默认继承上级的模板，但费率需要调整
        childPolicy := &AgentPolicy{
            AgentID:    agentID,
            ChannelID:  parentPolicy.ChannelID,
            TemplateID: parentPolicy.TemplateID,
            // 不设置自定义费率，使用模板默认值
            // 或者可以在这里设置比上级高一点的费率
        }

        s.agentPolicyRepo.Create(childPolicy)
    }

    return nil
}

// SetCustomRate 设置自定义费率覆盖
func (s *PolicyInheritService) SetCustomRate(
    agentID int64,
    channelID int64,
    rateType string,
    rate decimal.Decimal,
) error {
    policy, err := s.agentPolicyRepo.GetByAgentAndChannel(agentID, channelID)
    if err != nil {
        return err
    }

    // 验证费率范围
    template, _ := s.templateRepo.GetByID(policy.TemplateID)
    channel, _ := s.channelRepo.GetByID(channelID)

    // 费率必须在通道允许的范围内
    if rate.LessThan(channel.RateMin) || rate.GreaterThan(channel.RateMax) {
        return fmt.Errorf("费率必须在 %s - %s 范围内",
            channel.RateMin.String(), channel.RateMax.String())
    }

    // 费率必须 >= 上级的费率（保证有利润空间）
    parent, _ := s.agentRepo.GetByID(policy.AgentID)
    if parent.ParentID != nil {
        parentPolicy, _ := s.agentPolicyRepo.GetByAgentAndChannel(*parent.ParentID, channelID)
        parentRate := s.getRate(parentPolicy, rateType)
        if rate.LessThan(parentRate) {
            return errors.New("费率不能低于上级的费率")
        }
    }

    // 更新费率
    switch rateType {
    case "CREDIT":
        policy.CreditRate = &rate
    case "DEBIT":
        policy.DebitRate = &rate
    }

    return s.agentPolicyRepo.Update(policy)
}

// ApplyTemplate 应用政策模板
func (s *PolicyInheritService) ApplyTemplate(
    agentID int64,
    channelID int64,
    templateID int64,
) error {
    // 1. 验证模板存在
    template, err := s.templateRepo.GetByID(templateID)
    if err != nil {
        return err
    }

    // 2. 验证通道匹配
    if template.ChannelID != channelID {
        return errors.New("模板与通道不匹配")
    }

    // 3. 更新或创建代理商政策
    policy, _ := s.agentPolicyRepo.GetByAgentAndChannel(agentID, channelID)
    if policy == nil {
        policy = &AgentPolicy{
            AgentID:    agentID,
            ChannelID:  channelID,
            TemplateID: templateID,
        }
        return s.agentPolicyRepo.Create(policy)
    }

    policy.TemplateID = templateID
    // 清除自定义覆盖，使用新模板默认值
    policy.CreditRate = nil
    policy.DebitRate = nil
    policy.DebitCap = nil

    return s.agentPolicyRepo.Update(policy)
}
```

### 22.5 政策模板CRUD API

```go
// internal/app/policy/template_service.go

type PolicyTemplateService struct {
    templateRepo     PolicyTemplateRepository
    rateStageRepo    PolicyRateStageRepository
    rewardRepo       PolicyActivationRewardRepository
    depositRepo      PolicyDepositCashbackRepository
    simRepo          PolicySimCashbackRepository
}

// CreateTemplate 创建政策模板
func (s *PolicyTemplateService) CreateTemplate(req *CreateTemplateRequest) (*PolicyTemplate, error) {
    // 1. 验证通道存在
    channel, _ := s.channelRepo.GetByID(req.ChannelID)
    if channel == nil {
        return nil, errors.New("通道不存在")
    }

    // 2. 创建主表
    template := &PolicyTemplate{
        TemplateName: req.TemplateName,
        ChannelID:    req.ChannelID,
        IsDefault:    req.IsDefault,
        CreditRate:   req.CreditRate,
        DebitRate:    req.DebitRate,
        DebitCap:     req.DebitCap,
        UnionpayRate: req.UnionpayRate,
        WechatRate:   req.WechatRate,
        AlipayRate:   req.AlipayRate,
        T0FeeType:    req.T0FeeType,
        Status:       1,
    }

    if err := s.templateRepo.Create(template); err != nil {
        return nil, err
    }

    // 3. 创建费率阶梯规则
    for _, stage := range req.RateStages {
        rateStage := &PolicyRateStage{
            TemplateID: template.ID,
            StageType:  stage.StageType,
            DayStart:   stage.DayStart,
            DayEnd:     stage.DayEnd,
            RateAdjust: stage.RateAdjust,
        }
        s.rateStageRepo.Create(rateStage)
    }

    // 4. 创建激活奖励规则
    for _, reward := range req.ActivationRewards {
        activationReward := &PolicyActivationReward{
            TemplateID:    template.ID,
            DayStart:      reward.DayStart,
            DayEnd:        reward.DayEnd,
            TradeAmount:   reward.TradeAmount,
            RewardAmount:  reward.RewardAmount,
            LayerRewards:  reward.LayerRewards,
            EffectiveDate: reward.EffectiveDate,
        }
        s.rewardRepo.Create(activationReward)
    }

    // 5. 创建押金返现规则
    for _, deposit := range req.DepositCashbacks {
        depositCashback := &PolicyDepositCashback{
            TemplateID:     template.ID,
            DepositAmount:  deposit.DepositAmount,
            CashbackAmount: deposit.CashbackAmount,
        }
        s.depositRepo.Create(depositCashback)
    }

    // 6. 创建流量费返现规则
    for _, sim := range req.SimCashbacks {
        simCashback := &PolicySimCashback{
            TemplateID:     template.ID,
            SimType:        sim.SimType,
            SimFee:         sim.SimFee,
            CashbackAmount: sim.CashbackAmount,
            LayerCashback:  sim.LayerCashback,
        }
        s.simRepo.Create(simCashback)
    }

    // 7. 如果是默认模板，取消其他默认模板
    if req.IsDefault {
        s.templateRepo.ClearDefaultByChannel(req.ChannelID, template.ID)
    }

    return template, nil
}

// CopyTemplate 复制政策模板
func (s *PolicyTemplateService) CopyTemplate(templateID int64, newName string) (*PolicyTemplate, error) {
    // 1. 获取原模板
    original, _ := s.templateRepo.GetByID(templateID)
    if original == nil {
        return nil, errors.New("模板不存在")
    }

    // 2. 复制主表
    newTemplate := &PolicyTemplate{
        TemplateName: newName,
        ChannelID:    original.ChannelID,
        IsDefault:    false,
        CreditRate:   original.CreditRate,
        DebitRate:    original.DebitRate,
        DebitCap:     original.DebitCap,
        UnionpayRate: original.UnionpayRate,
        WechatRate:   original.WechatRate,
        AlipayRate:   original.AlipayRate,
        T0FeeType:    original.T0FeeType,
        Status:       1,
    }
    s.templateRepo.Create(newTemplate)

    // 3. 复制子规则
    s.copyRateStages(templateID, newTemplate.ID)
    s.copyActivationRewards(templateID, newTemplate.ID)
    s.copyDepositCashbacks(templateID, newTemplate.ID)
    s.copySimCashbacks(templateID, newTemplate.ID)

    return newTemplate, nil
}

// GetTemplateDetail 获取模板详情（含所有子规则）
func (s *PolicyTemplateService) GetTemplateDetail(templateID int64) (*TemplateDetailResponse, error) {
    template, _ := s.templateRepo.GetByID(templateID)
    if template == nil {
        return nil, errors.New("模板不存在")
    }

    rateStages, _ := s.rateStageRepo.GetByTemplateID(templateID)
    rewards, _ := s.rewardRepo.GetByTemplateID(templateID)
    deposits, _ := s.depositRepo.GetByTemplateID(templateID)
    sims, _ := s.simRepo.GetByTemplateID(templateID)

    return &TemplateDetailResponse{
        Template:           template,
        RateStages:         rateStages,
        ActivationRewards:  rewards,
        DepositCashbacks:   deposits,
        SimCashbacks:       sims,
    }, nil
}
```

### 22.6 费率阶梯示例配置

```sql
-- 政策模板示例
INSERT INTO policy_templates (template_name, channel_id, is_default, credit_rate, debit_rate, debit_cap)
VALUES ('拉卡拉默认模板', 1, true, 0.0051, 0.0051, 25.00);

-- 费率阶梯规则示例
-- 按商户入网时间：0-180天费率+0.04，181-360天费率+0.02，361天以上不调整
INSERT INTO policy_rate_stages (template_id, stage_type, day_start, day_end, rate_adjust) VALUES
(1, 1, 0, 180, 0.0004),     -- 商户入网0-180天，费率+0.04%
(1, 1, 181, 360, 0.0002),   -- 商户入网181-360天，费率+0.02%
(1, 1, 361, NULL, 0.0000);  -- 商户入网361天以上，不调整

-- 激活奖励规则示例
-- 30天内交易达10000，奖励50元
INSERT INTO policy_activation_rewards (template_id, day_start, day_end, trade_amount, reward_amount, layer_rewards)
VALUES (1, 0, 30, 10000.00, 50.00, '[{"layer":1,"amount":50},{"layer":2,"amount":20}]');

-- 60天内交易达30000，额外奖励100元
INSERT INTO policy_activation_rewards (template_id, day_start, day_end, trade_amount, reward_amount)
VALUES (1, 31, 60, 30000.00, 100.00);

-- 押金返现规则示例
INSERT INTO policy_deposit_cashbacks (template_id, deposit_amount, cashback_amount) VALUES
(1, 99.00, 89.00),    -- 99押金返89
(1, 199.00, 180.00),  -- 199押金返180
(1, 299.00, 270.00);  -- 299押金返270

-- 流量费返现规则示例
INSERT INTO policy_sim_cashbacks (template_id, sim_type, sim_fee, cashback_amount, layer_cashback) VALUES
(1, 1, 79.00, 69.00, '[{"layer":1,"amount":69},{"layer":2,"amount":5}]'),  -- 首次79返69（多层级）
(1, 2, 79.00, 59.00, NULL);  -- 续费79返59（仅直属）
```

### 22.7 政策变更审计日志

```sql
-- 政策变更日志表
CREATE TABLE policy_change_logs (
    id              BIGSERIAL PRIMARY KEY,
    change_type     VARCHAR(50) NOT NULL,       -- TEMPLATE_CREATE, RATE_UPDATE, REWARD_ADD, etc.
    template_id     BIGINT,
    agent_id        BIGINT,
    channel_id      BIGINT,
    before_value    JSONB,                      -- 变更前的值
    after_value     JSONB,                      -- 变更后的值
    operator_id     BIGINT NOT NULL,            -- 操作人
    operator_type   SMALLINT NOT NULL,          -- 1=管理员 2=代理商
    remark          VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_policy_change_template ON policy_change_logs(template_id);
CREATE INDEX idx_policy_change_agent ON policy_change_logs(agent_id);
CREATE INDEX idx_policy_change_time ON policy_change_logs(created_at);
```

```go
// 记录政策变更日志
func (s *PolicyTemplateService) logChange(
    changeType string,
    templateID *int64,
    agentID *int64,
    channelID *int64,
    beforeValue interface{},
    afterValue interface{},
    operatorID int64,
    operatorType int,
) {
    beforeJSON, _ := json.Marshal(beforeValue)
    afterJSON, _ := json.Marshal(afterValue)

    log := &PolicyChangeLog{
        ChangeType:   changeType,
        TemplateID:   templateID,
        AgentID:      agentID,
        ChannelID:    channelID,
        BeforeValue:  string(beforeJSON),
        AfterValue:   string(afterJSON),
        OperatorID:   operatorID,
        OperatorType: operatorType,
    }

    s.changeLogRepo.Create(log)
}
```

---

## 二十三、关键文件清单汇总

### 后端核心文件

| 模块 | 文件路径 | 功能说明 |
|------|----------|----------|
| **分润计算** | `internal/app/profit/calculator.go` | 交易分润计算器 |
| | `internal/app/profit/reward_checker.go` | 激活奖励检查器 |
| | `internal/app/profit/deposit_cashback.go` | 押金返现计算 |
| | `internal/app/profit/sim_cashback.go` | 流量费返现计算 |
| | `internal/app/profit/reconciliation.go` | 分润对账 |
| | `internal/app/profit/adjustment.go` | 手动调账 |
| **政策引擎** | `internal/app/policy/engine.go` | 政策引擎核心 |
| | `internal/app/policy/inherit.go` | 政策继承逻辑 |
| | `internal/app/policy/template_service.go` | 模板CRUD服务 |
| **代理商** | `internal/app/agent/service.go` | 代理商服务 |
| | `internal/app/agent/chain_tracer.go` | 代理商链追溯 |
| | `internal/app/agent/tree.go` | 机构树查询 |
| **钱包** | `internal/app/wallet/service.go` | 钱包服务 |
| | `internal/app/wallet/withdraw_service.go` | 提现服务 |
| | `internal/app/wallet/recharge_wallet_service.go` | 充值钱包 |
| | `internal/app/wallet/deposit_wallet_service.go` | 沉淀钱包 |
| **通道** | `internal/adapters/channel/interface.go` | 通道适配器接口 |
| | `internal/adapters/channel/factory.go` | 通道工厂 |
| | `internal/adapters/channel/lakala/adapter.go` | 拉卡拉适配器 |
| **交易** | `internal/app/transaction/sync_service.go` | 交易同步服务 |
| | `internal/app/transaction/callback_handler.go` | 回调处理 |
| **终端** | `internal/app/terminal/dispatch.go` | 机具下发/回拨 |
| | `internal/app/terminal/config.go` | 终端配置 |
| **定时任务** | `internal/scheduler/profit_scheduler.go` | 分润调度器 |
| | `internal/scheduler/reward_scheduler.go` | 奖励调度器 |
| | `internal/scheduler/stats_scheduler.go` | 统计调度器 |

### 数据库迁移文件

| 文件 | 说明 |
|------|------|
| `migrations/001_create_channels.sql` | 通道相关表 |
| `migrations/002_create_agents.sql` | 代理商相关表 |
| `migrations/003_create_policies.sql` | 政策模板相关表 |
| `migrations/004_create_terminals.sql` | 终端相关表 |
| `migrations/005_create_merchants.sql` | 商户相关表 |
| `migrations/006_create_transactions.sql` | 交易相关表 |
| `migrations/007_create_profits.sql` | 分润相关表 |
| `migrations/008_create_wallets.sql` | 钱包相关表 |
| `migrations/009_create_withdrawals.sql` | 提现相关表 |
| `migrations/010_create_deductions.sql` | 代扣相关表 |
| `migrations/011_create_messages.sql` | 消息相关表 |
| `migrations/012_create_marketing.sql` | 营销相关表 |
| `migrations/013_create_stats.sql` | 统计相关表 |

---

## 二十四、验证方案

### 24.1 单元测试

```go
// internal/app/profit/calculator_test.go

func TestTransactionProfitCalculator(t *testing.T) {
    calculator := NewTransactionProfitCalculator(...)

    t.Run("单级代理分润计算", func(t *testing.T) {
        tx := &Transaction{
            Amount:     decimal.NewFromInt(10000),
            MerchantID: 1,
            AgentID:    1,
        }

        results, err := calculator.Calculate(tx)
        assert.NoError(t, err)
        assert.Equal(t, 1, len(results))
        assert.Equal(t, decimal.NewFromFloat(8.0), results[0].ProfitAmount)
    })

    t.Run("多级代理分润计算", func(t *testing.T) {
        // ... 测试多层级
    })

    t.Run("阶梯费率分润计算", func(t *testing.T) {
        // ... 测试阶梯费率
    })

    t.Run("取下原则测试", func(t *testing.T) {
        // ... 测试费率相同时跳过
    })
}
```

### 24.2 集成测试

```go
// tests/integration/profit_flow_test.go

func TestCompleteProfitFlow(t *testing.T) {
    // 1. 创建测试数据
    channel := createTestChannel()
    template := createTestPolicyTemplate(channel.ID)
    agent := createTestAgent(template.ID)
    merchant := createTestMerchant(agent.ID)
    terminal := createTestTerminal(merchant.ID)

    // 2. 模拟交易
    tx := createTestTransaction(merchant.ID, terminal.ID, 10000)

    // 3. 触发分润计算
    profitService.Calculate(tx.ID)

    // 4. 验证分润结果
    profits, _ := profitRepo.GetByTransactionID(tx.ID)
    assert.NotEmpty(t, profits)

    // 5. 验证钱包余额
    wallet, _ := walletRepo.GetByAgentAndType(agent.ID, WalletTypeProfit)
    assert.Equal(t, decimal.NewFromFloat(8.0), wallet.Balance)
}
```

### 24.3 压力测试

```bash
# 使用 wrk 进行压力测试
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/transactions/callback

# 预期结果（日5000笔 = 0.06 TPS）
# 实际测试目标：100 TPS（留足余量）
```

### 24.4 对账验证

```sql
-- 每日对账SQL
SELECT
    t.settle_date,
    COUNT(*) as tx_count,
    SUM(t.amount) as tx_total,
    SUM(p.profit_amount) as profit_total,
    SUM(wl.amount) as wallet_total,
    CASE
        WHEN SUM(p.profit_amount) = SUM(wl.amount) THEN '✓ 一致'
        ELSE '✗ 差异: ' || (SUM(p.profit_amount) - SUM(wl.amount))::text
    END as check_result
FROM transactions t
LEFT JOIN profit_records p ON t.id = p.transaction_id
LEFT JOIN wallet_logs wl ON p.id = wl.ref_id AND wl.ref_type = 'profit'
WHERE t.settle_date = CURRENT_DATE - 1
GROUP BY t.settle_date;
```

---

**计划状态**: ✅ 完整设计完成（所有遗留问题已确认）

**设计文档包含**:
- 24个核心章节
- 45+张数据库表设计
- 完整Go代码实现示例
- 4种分润类型详细算法
- 3种钱包系统设计
- 8通道适配器架构
- 政策模板引擎
- 验证方案

**遗留问题确认结果** ✅:
| 问题 | 确认结果 |
|------|----------|
| T+0秒到费率加1/加2/加3 | ✅ 笔数费，每笔1-3元，可设为0 |
| 服务商偷数据权限 | ✅ 暂不开发，后续版本考虑 |
| 退费分润处理 | ✅ 分润回扣，标注退货 |
