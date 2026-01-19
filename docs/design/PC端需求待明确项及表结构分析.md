# PC端需求待明确项及表结构分析

## 一、需求不明确点汇总

### 1.1 业务规则不明确

| 序号 | 模块 | 待明确项 | 影响范围 | 建议 |
|------|------|----------|----------|------|
| 1 | **分润计算** | 分润计算的精度规则（四舍五入/向下取整/银行家舍入） | 分润金额准确性 | 需确认精度规则，建议使用分作为最小单位 |
| 2 | **分润计算** | 分润计算时机：实时计算还是T+1结算？ | 系统架构设计 | 文档说"实时计算"但需确认 |
| 3 | **费率阶梯** | 阶梯费率边界处理（如第180天算第一阶段还是第二阶段） | 分润金额 | 建议采用左闭右开区间 [0,180) |
| 4 | **激活奖励** | 激活奖励的"多层级"具体指几级？无限级还是有限制？ | 奖励发放 | 需确认多层级分润规则 |
| 5 | **激活奖励** | 达标交易额的统计口径：累计还是单笔？ | 奖励计算 | 需确认统计规则 |
| 6 | **押金返现** | 押金返现是实时返还还是T+N返还？ | 返现时效 | 需确认时效要求 |
| 7 | **流量费返现** | 续费次数如何判断？同一商户多台机具如何计算？ | 返现金额 | 需确认按商户还是按机具统计 |
| 8 | **沉淀钱包** | 使用下级余额后，下级提现失败的具体处理流程 | 资金安全 | 需明确回充机制和时限 |
| 9 | **沉淀钱包** | 沉淀钱包比例是按单个下级还是所有下级总额？ | 可用金额计算 | 需确认计算口径 |
| 10 | **充值钱包贷款** | 贷款利率计算方式（日息/月息/年息） | 利息计算 | 需确认计息规则 |
| 11 | **充值钱包贷款** | 逾期处理机制和违约金规则 | 风控 | 需明确逾期规则 |
| 12 | **代扣管理** | 多钱包扣款的优先级顺序 | 扣款逻辑 | 需确认扣款顺序 |
| 13 | **代扣管理** | 代扣期数的每期扣款时间点（每日/每周/每月） | 扣款周期 | 需确认扣款频率 |
| 14 | **终端下发** | 跨级下发时，中间层级代理商的货款代扣如何处理 | 货款结算 | 需确认跨级财务关系 |
| 15 | **提现** | 提现手续费由谁承担（代理商/平台） | 费用归属 | 需确认费用规则 |
| 16 | **提现** | 单日/单笔提现限额是否有限制 | 风控规则 | 需确认限额 |
| 17 | **商户类型** | 月均交易额计算时，新入网不足6个月的商户如何处理 | 分类准确性 | 需确认计算规则 |
| 18 | **偷数据权限** | "偷数据"的具体含义和业务场景 | 功能设计 | 需详细说明业务场景 |
| 19 | **税筹通道** | 税筹费率（9%+3元）是固定还是可配置 | 费率配置 | 需确认配置需求 |
| 20 | **T+0秒到** | 秒到费用（加1/加2/加3）的收取对象和分配规则 | 费用分配 | 需确认分成规则 |

### 1.2 界面交互不明确

| 序号 | 页面 | 待明确项 | 建议 |
|------|------|----------|------|
| 1 | 代理商列表 | 是否支持批量操作（批量禁用/批量分配模板） | 建议支持批量操作 |
| 2 | 终端入库 | 入库方式：Excel导入/手动输入/API同步？ | 需确认入库方式 |
| 3 | 终端入库 | 入库时是否需要校验SN号格式和通道归属 | 建议增加校验 |
| 4 | 商户费率修改 | 修改后是否需要同步到通道方？同步失败如何处理？ | 需确认同步机制 |
| 5 | 手动调账 | 调账审批流程：单级审批还是多级审批？ | 需确认审批流程 |
| 6 | 手动调账 | 调账金额是否有上限限制 | 建议设置限额 |
| 7 | 提现审核 | 批量审核时部分失败如何处理 | 需确认异常处理 |
| 8 | 机构树 | 代理商层级深度是否有限制 | 需确认最大层级 |
| 9 | 政策模板 | 模板删除时如果有代理商正在使用如何处理 | 建议禁止删除或软删除 |
| 10 | 伙伴代扣 | "伙伴"的定义：同级代理商？还是任意代理商？ | 需确认关系定义 |

### 1.3 数据权限不明确

| 序号 | 场景 | 待明确项 | 建议 |
|------|------|----------|------|
| 1 | 代理商登录PC端 | 代理商能看到几级下级的数据？ | 需确认数据范围 |
| 2 | 代理商登录PC端 | 代理商能否修改下级代理商的政策模板？ | 需确认权限范围 |
| 3 | 运营管理员 | 运营管理员能管理哪些通道的数据？ | 需确认通道权限 |
| 4 | 财务人员 | 财务能查看所有代理商的钱包还是按通道隔离？ | 需确认数据隔离 |
| 5 | 跨级操作 | PC端跨级下发/回拨的权限控制规则 | 需确认权限规则 |

---

## 二、表结构缺失分析

### 2.1 现有表结构（已实现）

| 表名 | 说明 | 状态 |
|------|------|------|
| channels | 支付通道表 | ✅ 已有 |
| agents | 代理商表 | ✅ 已有 |
| agent_policies | 代理商政策表 | ✅ 已有 |
| merchants | 商户表 | ✅ 已有 |
| terminals | 终端表 | ✅ 已有 |
| transactions | 交易流水表 | ✅ 已有 |
| profit_records | 分润明细表 | ✅ 已有 |
| wallets | 钱包表 | ✅ 已有 |
| wallet_logs | 钱包流水表 | ✅ 已有 |
| policy_templates | 政策模板表 | ✅ 已有 |
| device_fees | 流量费/服务费记录表 | ✅ 已有 |
| rate_changes | 费率变更记录表 | ✅ 已有 |
| messages | 消息通知表 | ✅ 已有 |
| raw_callback_logs | 原始回调日志表 | ✅ 已有 |

### 2.2 缺失的表结构

#### 2.2.1 政策模板相关表（缺失）

```sql
-- 费率阶梯配置表
CREATE TABLE IF NOT EXISTS policy_rate_stages (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,              -- 关联政策模板
    stage_type      SMALLINT NOT NULL,            -- 1:按商户入网时间 2:按代理商入网时间
    start_day       INT NOT NULL,                 -- 开始天数
    end_day         INT,                          -- 结束天数（NULL表示无穷）
    rate_adjustment DECIMAL(10,4) NOT NULL,       -- 费率调整值
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 激活奖励配置表
CREATE TABLE IF NOT EXISTS policy_rewards (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,              -- 关联政策模板
    start_day       INT NOT NULL,                 -- 开始天数
    end_day         INT NOT NULL,                 -- 结束天数
    target_amount   BIGINT NOT NULL,              -- 达标交易额（分）
    reward_amount   BIGINT NOT NULL,              -- 奖励金额（分）
    is_multi_level  BOOLEAN DEFAULT FALSE,        -- 是否多层级分润
    effective_date  DATE,                         -- 生效日期
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 押金返现配置表
CREATE TABLE IF NOT EXISTS policy_deposit_cashbacks (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    deposit_amount  BIGINT NOT NULL,              -- 押金金额（分）
    cashback_amount BIGINT NOT NULL,              -- 返现金额（分）
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 流量费返现配置表
CREATE TABLE IF NOT EXISTS policy_sim_cashbacks (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    fee_type        SMALLINT NOT NULL,            -- 1:首次 2:第2次 3:第3次及以后
    fee_amount      BIGINT NOT NULL,              -- 流量费金额（分）
    cashback_amount BIGINT NOT NULL,              -- 返现金额（分）
    is_multi_level  BOOLEAN DEFAULT FALSE,        -- 是否多层级分润
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.2 终端管理相关表（缺失）

```sql
-- 终端流转记录表
CREATE TABLE IF NOT EXISTS terminal_transfers (
    id              BIGSERIAL PRIMARY KEY,
    terminal_id     BIGINT NOT NULL,
    terminal_sn     VARCHAR(50) NOT NULL,
    transfer_type   SMALLINT NOT NULL,            -- 1:入库 2:下发 3:回拨
    from_agent_id   BIGINT,                       -- 来源代理商
    to_agent_id     BIGINT,                       -- 目标代理商
    operator_id     BIGINT,                       -- 操作人ID
    operator_type   SMALLINT,                     -- 1:管理员 2:代理商
    remark          VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 货款代扣记录表
CREATE TABLE IF NOT EXISTS cargo_deductions (
    id              BIGSERIAL PRIMARY KEY,
    transfer_id     BIGINT NOT NULL,              -- 关联终端流转记录
    agent_id        BIGINT NOT NULL,              -- 被扣款代理商
    terminal_count  INT NOT NULL,                 -- 终端数量
    unit_price      BIGINT NOT NULL,              -- 单价（分）
    total_amount    BIGINT NOT NULL,              -- 总金额（分）
    deducted_amount BIGINT DEFAULT 0,             -- 已扣金额（分）
    status          SMALLINT DEFAULT 0,           -- 0:待扣款 1:扣款中 2:已完成
    wallet_sources  JSONB,                        -- 扣款钱包来源 [1,2,3]
    accepted_at     TIMESTAMPTZ,                  -- 下级接受时间
    completed_at    TIMESTAMPTZ,                  -- 完成时间
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.3 代扣管理相关表（缺失）

```sql
-- 代扣记录表（上级扣款/伙伴代扣）
CREATE TABLE IF NOT EXISTS deductions (
    id              BIGSERIAL PRIMARY KEY,
    deduction_type  SMALLINT NOT NULL,            -- 1:上级扣款 2:伙伴代扣
    initiator_id    BIGINT NOT NULL,              -- 发起人（代理商ID）
    target_id       BIGINT NOT NULL,              -- 被扣款人（代理商ID）
    total_amount    BIGINT NOT NULL,              -- 总金额（分）
    deducted_amount BIGINT DEFAULT 0,             -- 已扣金额（分）
    periods         INT NOT NULL,                 -- 总期数（0表示一次性）
    current_period  INT DEFAULT 0,                -- 当前期数
    wallet_sources  JSONB,                        -- 扣款钱包来源
    agreement_type  SMALLINT,                     -- 1:系统协议 2:线下协议
    agreement_no    VARCHAR(64),                  -- 协议编号
    status          SMALLINT DEFAULT 0,           -- 0:待确认 1:进行中 2:已完成 3:已终止
    confirmed_at    TIMESTAMPTZ,                  -- 确认时间
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 代扣执行记录表
CREATE TABLE IF NOT EXISTS deduction_logs (
    id              BIGSERIAL PRIMARY KEY,
    deduction_id    BIGINT NOT NULL,
    period_no       INT NOT NULL,                 -- 第几期
    plan_amount     BIGINT NOT NULL,              -- 计划扣款金额
    actual_amount   BIGINT NOT NULL,              -- 实际扣款金额
    wallet_type     SMALLINT NOT NULL,            -- 从哪个钱包扣
    status          SMALLINT DEFAULT 0,           -- 0:待扣 1:成功 2:余额不足
    executed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.4 提现管理相关表（缺失）

```sql
-- 提现申请表
CREATE TABLE IF NOT EXISTS withdrawals (
    id              BIGSERIAL PRIMARY KEY,
    withdrawal_no   VARCHAR(32) NOT NULL UNIQUE,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,
    amount          BIGINT NOT NULL,              -- 申请金额（分）
    fee             BIGINT DEFAULT 0,             -- 手续费（分）
    tax             BIGINT DEFAULT 0,             -- 税费（分）
    actual_amount   BIGINT NOT NULL,              -- 实际到账金额（分）
    bank_name       VARCHAR(100),
    bank_account    VARCHAR(30),
    bank_card_no    VARCHAR(25),
    status          SMALLINT DEFAULT 0,           -- 0:待审核 1:审核通过 2:审核拒绝 3:打款中 4:打款成功 5:打款失败
    audit_user_id   BIGINT,
    audit_time      TIMESTAMPTZ,
    audit_remark    VARCHAR(255),
    pay_time        TIMESTAMPTZ,
    pay_remark      VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 提现门槛配置表
CREATE TABLE IF NOT EXISTS withdrawal_thresholds (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,
    threshold       BIGINT NOT NULL,              -- 提现门槛（分）
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(channel_id, wallet_type)
);
```

#### 2.2.5 商户扩展表（缺失）

```sql
-- 商户登记信息表（存储代理商登记的完整手机号）
CREATE TABLE IF NOT EXISTS merchant_registrations (
    id              BIGSERIAL PRIMARY KEY,
    merchant_id     BIGINT NOT NULL UNIQUE,
    agent_id        BIGINT NOT NULL,              -- 登记的代理商
    full_phone      VARCHAR(255),                 -- 加密存储的完整手机号
    remark          TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 商户统计表（用于商户类型分类）
CREATE TABLE IF NOT EXISTS merchant_statistics (
    id              BIGSERIAL PRIMARY KEY,
    merchant_id     BIGINT NOT NULL UNIQUE,
    merchant_type   SMALLINT,                     -- 1:忠诚 2:优质 3:潜力 4:一般 5:低活跃 6:30天无交易
    monthly_avg_amount BIGINT DEFAULT 0,          -- 月均交易额（分）
    total_amount    BIGINT DEFAULT 0,             -- 累计交易额（分）
    month_amount    BIGINT DEFAULT 0,             -- 本月交易额（分）
    last_trade_time TIMESTAMPTZ,                  -- 最后交易时间
    calculated_at   TIMESTAMPTZ,                  -- 计算时间
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.6 通道配置相关表（缺失）

```sql
-- 通道成本费率表
CREATE TABLE IF NOT EXISTS channel_costs (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL UNIQUE,
    credit_rate     DECIMAL(10,4),                -- 贷记卡成本费率
    debit_rate      DECIMAL(10,4),                -- 借记卡成本费率
    debit_cap       DECIMAL(10,2),                -- 借记卡封顶
    unionpay_rate   DECIMAL(10,4),                -- 云闪付成本费率
    wechat_rate     DECIMAL(10,4),                -- 微信成本费率
    alipay_rate     DECIMAL(10,4),                -- 支付宝成本费率
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 通道调价阶梯表
CREATE TABLE IF NOT EXISTS channel_rate_stages (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    stage_type      SMALLINT NOT NULL,            -- 1:按商户入网时间 2:按代理商入网时间
    start_day       INT NOT NULL,
    end_day         INT,
    rate_adjustment DECIMAL(10,4) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 通道可见性配置表
CREATE TABLE IF NOT EXISTS channel_visibility (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    agent_id        BIGINT,                       -- NULL表示全局配置
    is_visible      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(channel_id, agent_id)
);

-- 偷数据权限配置表
CREATE TABLE IF NOT EXISTS data_access_permissions (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,
    is_enabled      BOOLEAN DEFAULT FALSE,
    min_amount      BIGINT,                       -- 最小金额（分）
    max_amount      BIGINT,                       -- 最大金额（分）
    start_time      TIME,                         -- 生效开始时间
    end_time        TIME,                         -- 生效结束时间
    agent_levels    JSONB,                        -- 适用代理商层级 [1,2]
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.7 营销管理相关表（缺失）

```sql
-- 海报分类表
CREATE TABLE IF NOT EXISTS poster_categories (
    id              BIGSERIAL PRIMARY KEY,
    category_name   VARCHAR(50) NOT NULL,
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 海报表
CREATE TABLE IF NOT EXISTS posters (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT NOT NULL,
    title           VARCHAR(100) NOT NULL,
    image_url       VARCHAR(255) NOT NULL,
    thumbnail_url   VARCHAR(255),
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 滚动图表
CREATE TABLE IF NOT EXISTS banners (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(100),
    image_url       VARCHAR(255) NOT NULL,
    link_url        VARCHAR(255),
    link_type       SMALLINT,                     -- 1:内部跳转 2:外部链接 3:无跳转
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    start_time      TIMESTAMPTZ,
    end_time        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.8 充值钱包/沉淀钱包相关表（缺失）

```sql
-- 沉淀钱包配置表
CREATE TABLE IF NOT EXISTS deposit_wallet_configs (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT,                       -- NULL表示全局配置
    is_visible      BOOLEAN DEFAULT TRUE,         -- 是否对代理商可见
    usage_ratio     DECIMAL(5,2) DEFAULT 30,      -- 可使用比例（%）
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 充值钱包贷款表
CREATE TABLE IF NOT EXISTS wallet_loans (
    id              BIGSERIAL PRIMARY KEY,
    loan_no         VARCHAR(32) NOT NULL UNIQUE,
    agent_id        BIGINT NOT NULL,
    loan_amount     BIGINT NOT NULL,              -- 贷款金额（分）
    interest_rate   DECIMAL(10,4) NOT NULL,       -- 月利率（%）
    total_interest  BIGINT DEFAULT 0,             -- 应付利息（分）
    repaid_principal BIGINT DEFAULT 0,            -- 已还本金（分）
    repaid_interest BIGINT DEFAULT 0,             -- 已还利息（分）
    agreement_no    VARCHAR(64),                  -- 线下协议编号
    status          SMALLINT DEFAULT 0,           -- 0:还款中 1:已结清 2:逾期
    loan_date       DATE NOT NULL,
    due_date        DATE,
    settled_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 贷款还款记录表
CREATE TABLE IF NOT EXISTS loan_repayments (
    id              BIGSERIAL PRIMARY KEY,
    loan_id         BIGINT NOT NULL,
    repay_amount    BIGINT NOT NULL,              -- 还款金额（分）
    principal       BIGINT NOT NULL,              -- 本金部分（分）
    interest        BIGINT NOT NULL,              -- 利息部分（分）
    repay_type      SMALLINT NOT NULL,            -- 1:自动扣款 2:手动还款
    wallet_type     SMALLINT,                     -- 从哪个钱包扣
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.9 代理商扩展表（缺失）

```sql
-- 代理商邀请码申请表
CREATE TABLE IF NOT EXISTS invite_code_applications (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    current_code    VARCHAR(20),                  -- 当前邀请码
    applied_code    VARCHAR(20) NOT NULL,         -- 申请的靓号
    status          SMALLINT DEFAULT 0,           -- 0:待审核 1:通过 2:拒绝
    audit_user_id   BIGINT,
    audit_time      TIMESTAMPTZ,
    audit_remark    VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 代理商可用通道表
CREATE TABLE IF NOT EXISTS agent_channels (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    is_enabled      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(agent_id, channel_id)
);
```

#### 2.2.10 系统管理相关表（缺失）

```sql
-- 系统用户表（后台管理员）
CREATE TABLE IF NOT EXISTS sys_users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    real_name       VARCHAR(50),
    phone           VARCHAR(20),
    email           VARCHAR(100),
    role_id         BIGINT,
    status          SMALLINT DEFAULT 1,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 角色表
CREATE TABLE IF NOT EXISTS sys_roles (
    id              BIGSERIAL PRIMARY KEY,
    role_name       VARCHAR(50) NOT NULL,
    role_code       VARCHAR(50) NOT NULL UNIQUE,
    description     VARCHAR(255),
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 权限表
CREATE TABLE IF NOT EXISTS sys_permissions (
    id              BIGSERIAL PRIMARY KEY,
    permission_name VARCHAR(50) NOT NULL,
    permission_code VARCHAR(100) NOT NULL UNIQUE,
    parent_id       BIGINT,
    menu_type       SMALLINT,                     -- 1:目录 2:菜单 3:按钮
    path            VARCHAR(255),
    component       VARCHAR(255),
    icon            VARCHAR(50),
    sort_order      INT DEFAULT 0,
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS sys_role_permissions (
    id              BIGSERIAL PRIMARY KEY,
    role_id         BIGINT NOT NULL,
    permission_id   BIGINT NOT NULL,
    UNIQUE(role_id, permission_id)
);

-- 操作日志表
CREATE TABLE IF NOT EXISTS sys_operation_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT,
    user_type       SMALLINT,                     -- 1:管理员 2:代理商
    module          VARCHAR(50),
    action          VARCHAR(50),
    method          VARCHAR(10),
    url             VARCHAR(255),
    params          JSONB,
    ip              VARCHAR(50),
    user_agent      VARCHAR(500),
    result          SMALLINT,                     -- 1:成功 0:失败
    error_msg       TEXT,
    duration        INT,                          -- 执行时长(ms)
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 系统配置表
CREATE TABLE IF NOT EXISTS sys_configs (
    id              BIGSERIAL PRIMARY KEY,
    config_key      VARCHAR(100) NOT NULL UNIQUE,
    config_value    TEXT,
    config_type     VARCHAR(20),                  -- string/number/boolean/json
    description     VARCHAR(255),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

#### 2.2.11 手动调账相关表（缺失）

```sql
-- 调账申请表
CREATE TABLE IF NOT EXISTS adjustments (
    id              BIGSERIAL PRIMARY KEY,
    adjustment_no   VARCHAR(32) NOT NULL UNIQUE,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT,
    wallet_type     SMALLINT NOT NULL,
    adjust_type     SMALLINT NOT NULL,            -- 1:增加 2:扣减
    amount          BIGINT NOT NULL,              -- 调账金额（分）
    reason          VARCHAR(500) NOT NULL,
    attachment_url  VARCHAR(255),                 -- 附件
    status          SMALLINT DEFAULT 0,           -- 0:待审核 1:通过 2:拒绝
    apply_user_id   BIGINT NOT NULL,              -- 申请人
    audit_user_id   BIGINT,                       -- 审核人
    audit_time      TIMESTAMPTZ,
    audit_remark    VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 三、agents表字段缺失

当前 agents 表需要补充的字段：

```sql
-- 补充代理商表字段
ALTER TABLE agents ADD COLUMN IF NOT EXISTS custom_invite_code VARCHAR(20);  -- 自定义靓号
ALTER TABLE agents ADD COLUMN IF NOT EXISTS id_card_encrypted VARCHAR(255);  -- 加密身份证
ALTER TABLE agents ADD COLUMN IF NOT EXISTS phone_encrypted VARCHAR(255);    -- 加密手机号
```

---

## 四、merchants表字段缺失

当前 merchants 表需要补充的字段：

```sql
-- 补充商户表字段
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS merchant_type SMALLINT;        -- 商户类型
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS is_direct BOOLEAN DEFAULT TRUE; -- 是否直营
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS activate_time TIMESTAMPTZ;      -- 激活时间
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS first_sim_time TIMESTAMPTZ;     -- 首次流量费时间
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS first_sim_amount BIGINT;        -- 首次流量费金额
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS sim_count INT DEFAULT 0;        -- 续费次数
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS scan_rate DECIMAL(10,4);        -- 扫码费率
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS t0_fee SMALLINT DEFAULT 0;      -- T+0秒到费用(元)
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS debit_cap DECIMAL(10,2);        -- 借记卡封顶
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS deposit_amount BIGINT DEFAULT 0; -- 押金金额
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS network_days INT DEFAULT 0;     -- 入网天数
```

---

## 五、terminals表字段缺失

当前 terminals 表需要补充的字段：

```sql
-- 补充终端表字段
ALTER TABLE terminals ADD COLUMN IF NOT EXISTS scan_rate DECIMAL(10,4);       -- 扫码费率
ALTER TABLE terminals ADD COLUMN IF NOT EXISTS t0_fee SMALLINT DEFAULT 0;     -- T+0费用(元): 0/1/2/3
ALTER TABLE terminals ADD COLUMN IF NOT EXISTS sim_interval_days INT DEFAULT 360; -- SIM续费间隔天数
```

---

## 六、优先级建议

### P0 - 核心功能必需（必须先完成）

1. 政策模板相关表（费率阶梯、激活奖励、押金返现、流量返现）
2. 提现管理相关表（提现申请、提现门槛）
3. 终端流转记录表
4. 系统用户、角色、权限表

### P1 - 重要功能（第二优先级）

5. 代扣管理相关表
6. 手动调账表
7. 商户扩展表（登记信息、统计表）
8. 货款代扣记录表

### P2 - 增强功能（第三优先级）

9. 通道配置相关表（成本、调价、可见性）
10. 充值钱包/沉淀钱包相关表
11. 营销管理相关表
12. 代理商扩展表

---

## 七、待确认事项汇总

请与产品/业务方确认以下事项：

1. **分润计算精度**：使用什么舍入规则？
2. **激活奖励多层级**：具体指几级分润？
3. **伙伴代扣**："伙伴"的具体定义？
4. **偷数据权限**：具体业务场景说明？
5. **沉淀钱包**：使用比例是按单个下级还是总额？
6. **贷款利率**：按日/月/年计息？
7. **代扣周期**：每期扣款的时间点？
8. **提现限额**：单日/单笔是否有限制？
9. **税筹费率**：是否可配置？
10. **代理商层级**：是否有最大层级限制？

---

*文档版本: v1.0*
*创建时间: 2026-01-18*
*状态: 待评审*
