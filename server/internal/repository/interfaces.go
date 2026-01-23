package repository

import (
	"time"

	"xiangshoufu/internal/models"
)

// RawCallbackRepository 原始回调日志仓库接口
type RawCallbackRepository interface {
	// Create 创建回调日志
	Create(log *models.RawCallbackLog) error

	// Update 更新回调日志
	Update(log *models.RawCallbackLog) error

	// FindByID 根据ID查找
	FindByID(id int64) (*models.RawCallbackLog, error)

	// FindByIdempotentKey 根据幂等键查找
	FindByIdempotentKey(key string) (*models.RawCallbackLog, error)

	// FindPendingLogs 查找待处理的日志
	FindPendingLogs(limit int) ([]*models.RawCallbackLog, error)

	// FindFailedLogs 查找失败的日志（用于重试）
	FindFailedLogs(maxRetry int, limit int) ([]*models.RawCallbackLog, error)

	// UpdateStatus 更新处理状态
	UpdateStatus(id int64, status int16, errorMsg string) error

	// IncrementRetryCount 增加重试次数
	IncrementRetryCount(id int64) error
}

// DeviceFeeRepository 流量费仓库接口
type DeviceFeeRepository interface {
	Create(fee *models.DeviceFee) error
	Update(fee *models.DeviceFee) error
	FindByOrderNo(orderNo string) (*models.DeviceFee, error)
	FindPendingCashback(limit int) ([]*models.DeviceFee, error)
	UpdateCashbackStatus(id int64, status int16, amount int64) error
}

// RateChangeRepository 费率变更仓库接口
type RateChangeRepository interface {
	Create(change *models.RateChange) error
	FindPendingSync(limit int) ([]*models.RateChange, error)
	UpdateSyncStatus(id int64, status int16) error
}

// MessageRepository 消息仓库接口
type MessageRepository interface {
	Create(msg *models.Message) error
	BatchCreate(msgs []*models.Message) error
	FindByAgentID(agentID int64, limit, offset int) ([]*models.Message, error)
	FindUnreadByAgentID(agentID int64) ([]*models.Message, error)
	MarkAsRead(id int64) error
	MarkAllAsRead(agentID int64) error
	DeleteExpired() (int64, error) // 删除过期消息

	// 扩展方法 - 按类型筛选
	FindByAgentIDAndTypes(agentID int64, types []int16, limit, offset int) ([]*models.Message, error)
	CountByAgentIDAndTypes(agentID int64, types []int16) (int64, error)

	// 扩展方法 - 分类统计
	GetStatsByAgentID(agentID int64) (*MessageStats, error)

	// 管理端方法
	FindAll(limit, offset int) ([]*models.Message, error)
	CountAll() (int64, error)
	FindByID(id int64) (*models.Message, error)
	Delete(id int64) error
	FindByAgentIDs(agentIDs []int64, limit, offset int) ([]*models.Message, error)
}

// MessageStats 消息统计
type MessageStats struct {
	Total        int64 `json:"total"`         // 总消息数
	UnreadTotal  int64 `json:"unread_total"`  // 未读总数
	ProfitCount  int64 `json:"profit_count"`  // 分润类消息数
	RegisterCount int64 `json:"register_count"` // 注册类消息数
	ConsumptionCount int64 `json:"consumption_count"` // 消费类消息数
	SystemCount  int64 `json:"system_count"`  // 系统类消息数
}

// TransactionRepository 交易仓库接口
type TransactionRepository interface {
	Create(tx *Transaction) error
	FindByID(id int64) (*Transaction, error)
	FindByOrderNo(orderNo string) (*Transaction, error)
	FindUnprocessedProfit(limit int) ([]*Transaction, error)
	UpdateProfitStatus(id int64, status int16) error
	BatchUpdateProfitStatus(ids []int64, status int16) error
	UpdateRefundStatus(id int64, status int16) error
	// 激活奖励相关
	GetTerminalTotalTradeAmount(terminalSN string) (int64, error)
}

// Transaction 交易模型（简化版，实际应从现有模型导入）
type Transaction struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	TradeNo      string    `json:"trade_no" gorm:"size:64"`
	OrderNo      string    `json:"order_no" gorm:"size:64;uniqueIndex"`
	ChannelID    int64     `json:"channel_id"`
	ChannelCode  string    `json:"channel_code" gorm:"size:32"`
	TerminalSN   string    `json:"terminal_sn" gorm:"size:50"`
	MerchantID   int64     `json:"merchant_id"`
	AgentID      int64     `json:"agent_id"`
	TradeType    int16     `json:"trade_type"`                     // 1消费 2撤销 3退货
	PayType      int16     `json:"pay_type"`                       // 1刷卡 2微信 3支付宝 4云闪付
	CardType     int16     `json:"card_type"`                      // 1借记卡 2贷记卡
	Amount       int64     `json:"amount"`                         // 分
	Fee          int64     `json:"fee"`                            // 手续费（分）
	Rate         string    `json:"rate"`                           // 费率
	D0Fee        int64     `json:"d0_fee"`                         // D0手续费
	HighRate     string    `json:"high_rate"`                      // 调价费率
	CardNo       string    `json:"card_no"`                        // 脱敏卡号
	ProfitStatus int16     `json:"profit_status" gorm:"default:0"` // 0待计算 1已计算 2失败
	RefundStatus int16     `json:"refund_status" gorm:"default:0"` // 0正常 1已退款
	TradeTime    time.Time `json:"trade_time"`
	ReceivedAt   time.Time `json:"received_at" gorm:"default:now()"`
	ExtData      string    `json:"ext_data" gorm:"type:jsonb"`
}

// ProfitRecordRepository 分润记录仓库接口
type ProfitRecordRepository interface {
	Create(record *ProfitRecord) error
	BatchCreate(records []*ProfitRecord) error
	FindByTransactionID(txID int64) ([]*ProfitRecord, error)
	RevokeByTransactionID(txID int64, reason string) error
}

// ProfitRecord 分润记录模型
type ProfitRecord struct {
	ID               int64      `json:"id" gorm:"primaryKey"`
	TransactionID    int64      `json:"transaction_id" gorm:"not null;index"`
	OrderNo          string     `json:"order_no" gorm:"size:64"`
	AgentID          int64      `json:"agent_id" gorm:"not null;index"`
	ProfitType       int16      `json:"profit_type" gorm:"not null"` // 1交易分润 2激活奖励 3押金返现 4流量返现
	TradeAmount      int64      `json:"trade_amount"`                // 交易金额（分）
	SelfRate         string     `json:"self_rate"`                   // 自身费率
	LowerRate        string     `json:"lower_rate"`                  // 下级费率
	RateDiff         string     `json:"rate_diff"`                   // 费率差
	ProfitAmount     int64      `json:"profit_amount"`               // 分润金额（分）
	SourceMerchantID int64      `json:"source_merchant_id"`
	SourceAgentID    int64      `json:"source_agent_id"`
	ChannelID        int64      `json:"channel_id"`
	WalletType       int16      `json:"wallet_type"`                    // 1分润钱包 2服务费钱包 3奖励钱包
	WalletStatus     int16      `json:"wallet_status" gorm:"default:0"` // 0待入账 1已入账
	IsRevoked        bool       `json:"is_revoked" gorm:"default:false"`
	RevokedAt        *time.Time `json:"revoked_at"`
	RevokeReason     string     `json:"revoke_reason"`
	CreatedAt        time.Time  `json:"created_at" gorm:"default:now()"`
}

// WalletRepository 钱包仓库接口
type WalletRepository interface {
	FindByAgentAndType(agentID int64, channelID int64, walletType int16) (*Wallet, error)
	UpdateBalance(id int64, amount int64) error
	BatchUpdateBalance(updates map[int64]int64) error
}

// Wallet 钱包模型
type Wallet struct {
	ID                int64     `json:"id" gorm:"primaryKey"`
	AgentID           int64     `json:"agent_id" gorm:"not null"`
	ChannelID         int64     `json:"channel_id" gorm:"not null"`
	WalletType        int16     `json:"wallet_type" gorm:"not null"`
	Balance           int64     `json:"balance" gorm:"default:0"` // 分
	FrozenAmount      int64     `json:"frozen_amount" gorm:"default:0"`
	TotalIncome       int64     `json:"total_income" gorm:"default:0"`
	TotalWithdraw     int64     `json:"total_withdraw" gorm:"default:0"`
	WithdrawThreshold int64     `json:"withdraw_threshold" gorm:"default:10000"` // 默认100元
	Version           int       `json:"version" gorm:"default:0"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"default:now()"`
}

// WalletLogRepository 钱包流水仓库接口
type WalletLogRepository interface {
	Create(log *WalletLog) error
	BatchCreate(logs []*WalletLog) error
}

// WalletLog 钱包流水
type WalletLog struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	WalletID      int64     `json:"wallet_id" gorm:"not null"`
	AgentID       int64     `json:"agent_id" gorm:"not null"`
	WalletType    int16     `json:"wallet_type" gorm:"not null"`
	LogType       int16     `json:"log_type" gorm:"not null"` // 1分润入账 2提现冻结 3提现成功 4提现退回 5调账 6代扣
	Amount        int64     `json:"amount"`                   // 分（可为负）
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	RefType       string    `json:"ref_type"`
	RefID         int64     `json:"ref_id"`
	Remark        string    `json:"remark"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:now()"`
}

// AgentRepository 代理商仓库接口
type AgentRepository interface {
	FindByID(id int64) (*Agent, error)
	FindByAgentNo(agentNo string) (*Agent, error)
	FindAncestors(agentID int64) ([]*Agent, error) // 查找所有上级代理商
}

// Agent 代理商模型
type Agent struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	AgentNo     string `json:"agent_no" gorm:"size:32;uniqueIndex"`
	AgentName   string `json:"agent_name" gorm:"size:100"`
	ParentID    int64  `json:"parent_id"`
	Path        string `json:"path" gorm:"size:500"` // 物化路径 /1/5/12/
	Level       int    `json:"level" gorm:"default:1"`
	DefaultRate string `json:"default_rate"`
	Status      int16  `json:"status" gorm:"default:1"`
}

// AgentPolicyRepository 代理商政策仓库接口
type AgentPolicyRepository interface {
	FindByAgentAndChannel(agentID int64, channelID int64) (*AgentPolicy, error)
}

// AgentPolicy 代理商政策
type AgentPolicy struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
	AgentID    int64  `json:"agent_id" gorm:"not null"`
	ChannelID  int64  `json:"channel_id" gorm:"not null"`
	TemplateID int64  `json:"template_id" gorm:"not null"`
	CreditRate string `json:"credit_rate"`
	DebitRate  string `json:"debit_rate"`
}

// MerchantRepository 商户仓库接口
type MerchantRepository interface {
	FindByID(id int64) (*models.Merchant, error)
	FindByMerchantNo(merchantNo string) (*models.Merchant, error)
	UpdateApproveStatus(id int64, status int16) error
}
