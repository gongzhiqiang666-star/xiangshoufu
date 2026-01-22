package models

import (
	"time"
)

// RawCallbackLog 原始回调日志
type RawCallbackLog struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	ChannelCode   string     `json:"channel_code" gorm:"size:32;not null"`
	ActionType    string     `json:"action_type" gorm:"size:64;not null"`
	RawRequest    string     `json:"raw_request" gorm:"type:jsonb;not null"`
	SignVerified  bool       `json:"sign_verified" gorm:"default:false"`
	ProcessStatus int16      `json:"process_status" gorm:"default:0"` // 0:待处理 1:成功 2:失败
	ErrorMessage  string     `json:"error_message" gorm:"type:text"`
	RetryCount    int16      `json:"retry_count" gorm:"default:0"`
	IdempotentKey string     `json:"idempotent_key" gorm:"size:128;not null"`
	ClientIP      string     `json:"client_ip" gorm:"size:45"`
	ReceivedAt    time.Time  `json:"received_at" gorm:"default:now()"`
	ProcessedAt   *time.Time `json:"processed_at"`
	CreatedDate   time.Time  `json:"created_date" gorm:"default:CURRENT_DATE"`
}

func (RawCallbackLog) TableName() string {
	return "raw_callback_logs"
}

// ProcessStatus 常量
const (
	ProcessStatusPending = 0 // 待处理
	ProcessStatusSuccess = 1 // 处理成功
	ProcessStatusFailed  = 2 // 处理失败
)

// DeviceFee 流量费/服务费记录
type DeviceFee struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	ChannelID      int64     `json:"channel_id" gorm:"not null"`
	ChannelCode    string    `json:"channel_code" gorm:"size:32;not null"`
	TerminalSN     string    `json:"terminal_sn" gorm:"size:50;not null"`
	MerchantNo     string    `json:"merchant_no" gorm:"size:64"`
	AgentID        int64     `json:"agent_id"`
	OrderNo        string    `json:"order_no" gorm:"size:64;not null;uniqueIndex"`
	FeeType        int16     `json:"fee_type" gorm:"not null"`         // 1:服务费 2:流量费
	FeeAmount      int64     `json:"fee_amount" gorm:"not null"`       // 分
	CashbackStatus int16     `json:"cashback_status" gorm:"default:0"` // 0:待计算 1:已返现 2:不返现
	CashbackAmount int64     `json:"cashback_amount" gorm:"default:0"`
	ChargingTime   time.Time `json:"charging_time" gorm:"not null"`
	ReceivedAt     time.Time `json:"received_at" gorm:"default:now()"`
	BrandCode      string    `json:"brand_code" gorm:"size:32"`
	ExtData        string    `json:"ext_data" gorm:"type:jsonb"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:now()"`
}

func (DeviceFee) TableName() string {
	return "device_fees"
}

// RateChange 费率变更记录
type RateChange struct {
	ID                   int64     `json:"id" gorm:"primaryKey"`
	ChannelID            int64     `json:"channel_id" gorm:"not null"`
	ChannelCode          string    `json:"channel_code" gorm:"size:32;not null"`
	TerminalSN           string    `json:"terminal_sn" gorm:"size:50;not null"`
	MerchantNo           string    `json:"merchant_no" gorm:"size:64;not null"`
	AgentID              int64     `json:"agent_id"`
	CreditRate           string    `json:"credit_rate" gorm:"type:decimal(10,4)"`
	CreditExtraRate      string    `json:"credit_extra_rate" gorm:"type:decimal(10,4)"`
	DebitRate            string    `json:"debit_rate" gorm:"type:decimal(10,4)"`
	AlipayRate           string    `json:"alipay_rate" gorm:"type:decimal(10,4)"`
	WechatRate           string    `json:"wechat_rate" gorm:"type:decimal(10,4)"`
	UnionpayRate         string    `json:"unionpay_rate" gorm:"type:decimal(10,4)"`
	CreditAdditionRate   string    `json:"credit_addition_rate" gorm:"type:decimal(10,4)"`
	UnionpayAdditionRate string    `json:"unionpay_addition_rate" gorm:"type:decimal(10,4)"`
	AlipayAdditionRate   string    `json:"alipay_addition_rate" gorm:"type:decimal(10,4)"`
	WechatAdditionRate   string    `json:"wechat_addition_rate" gorm:"type:decimal(10,4)"`
	SyncStatus           int16     `json:"sync_status" gorm:"default:0"` // 0:待同步 1:已同步
	ReceivedAt           time.Time `json:"received_at" gorm:"default:now()"`
	BrandCode            string    `json:"brand_code" gorm:"size:32"`
	ExtData              string    `json:"ext_data" gorm:"type:jsonb"`
	CreatedAt            time.Time `json:"created_at" gorm:"default:now()"`
}

func (RateChange) TableName() string {
	return "rate_changes"
}

// Message 消息通知
type Message struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	AgentID     int64      `json:"agent_id" gorm:"not null;index"`
	MessageType int16      `json:"message_type" gorm:"not null"` // 1:分润 2:激活奖励 3:押金返现 4:流量返现 5:退款撤销 6:系统公告
	Title       string     `json:"title" gorm:"size:64;not null"`
	Content     string     `json:"content" gorm:"type:text"`
	IsRead      bool       `json:"is_read" gorm:"default:false"`
	IsPushed    bool       `json:"is_pushed" gorm:"default:false"`
	RelatedID   int64      `json:"related_id"`
	RelatedType string     `json:"related_type" gorm:"size:32"`
	ExpireAt    *time.Time `json:"expire_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:now()"`
}

func (Message) TableName() string {
	return "messages"
}

// MessageType 消息类型常量
const (
	MessageTypeProfit       = 1 // 交易分润
	MessageTypeActivation   = 2 // 激活奖励
	MessageTypeDeposit      = 3 // 押金返现
	MessageTypeSimCashback  = 4 // 流量返现
	MessageTypeRefund       = 5 // 退款撤销
	MessageTypeAnnouncement = 6 // 系统公告
	MessageTypeNewAgent     = 7 // 新代理注册
	MessageTypeTransaction  = 8 // 交易通知
)

// MessageCategory APP端消息分类
const (
	MessageCategoryAll         = "all"         // 全部
	MessageCategoryProfit      = "profit"      // 分润（类型1,2,3,4）
	MessageCategoryRegister    = "register"    // 注册（类型7）
	MessageCategoryConsumption = "consumption" // 消费（类型8）
	MessageCategorySystem      = "system"      // 系统（类型5,6）
)

// GetMessageTypesByCategory 根据分类获取消息类型列表
func GetMessageTypesByCategory(category string) []int16 {
	switch category {
	case MessageCategoryProfit:
		return []int16{MessageTypeProfit, MessageTypeActivation, MessageTypeDeposit, MessageTypeSimCashback}
	case MessageCategoryRegister:
		return []int16{MessageTypeNewAgent}
	case MessageCategoryConsumption:
		return []int16{MessageTypeTransaction}
	case MessageCategorySystem:
		return []int16{MessageTypeRefund, MessageTypeAnnouncement}
	default:
		return nil // 全部类型
	}
}

// GetMessageTypeName 获取消息类型名称
func GetMessageTypeName(messageType int16) string {
	switch messageType {
	case MessageTypeProfit:
		return "交易分润"
	case MessageTypeActivation:
		return "激活奖励"
	case MessageTypeDeposit:
		return "押金返现"
	case MessageTypeSimCashback:
		return "流量返现"
	case MessageTypeRefund:
		return "退款撤销"
	case MessageTypeAnnouncement:
		return "系统公告"
	case MessageTypeNewAgent:
		return "新代理注册"
	case MessageTypeTransaction:
		return "交易通知"
	default:
		return "未知类型"
	}
}
