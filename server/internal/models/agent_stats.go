package models

import "time"

// AgentDailyStats 代理商每日统计汇总
type AgentDailyStats struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	AgentID   int64     `json:"agent_id" gorm:"index"`
	StatDate  time.Time `json:"stat_date" gorm:"type:date;index"`
	Scope     string    `json:"scope" gorm:"size:10"` // direct=直营, team=团队

	// 交易统计
	TransAmount int64 `json:"trans_amount"` // 交易金额(分)
	TransCount  int   `json:"trans_count"`  // 交易笔数

	// 分润统计(按类型分)
	ProfitTrade   int64 `json:"profit_trade"`   // 交易分润(分)
	ProfitDeposit int64 `json:"profit_deposit"` // 押金返现(分)
	ProfitSim     int64 `json:"profit_sim"`     // 流量返现(分)
	ProfitReward  int64 `json:"profit_reward"`  // 激活奖励(分)
	ProfitTotal   int64 `json:"profit_total"`   // 总分润(分)

	// 商户与终端统计
	MerchantCount        int `json:"merchant_count"`         // 商户数量
	MerchantNew          int `json:"merchant_new"`           // 新增商户数
	TerminalTotal        int `json:"terminal_total"`         // 终端总数
	TerminalActivated    int `json:"terminal_activated"`     // 已激活终端数
	TerminalNewActivated int `json:"terminal_new_activated"` // 当日新激活数

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (AgentDailyStats) TableName() string {
	return "agent_daily_stats"
}

// AgentMonthlyStats 代理商每月统计汇总
type AgentMonthlyStats struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	AgentID   int64  `json:"agent_id" gorm:"index"`
	StatMonth string `json:"stat_month" gorm:"size:7;index"` // '2026-01' 格式
	Scope     string `json:"scope" gorm:"size:10"`

	// 交易统计
	TransAmount int64 `json:"trans_amount"`
	TransCount  int   `json:"trans_count"`

	// 分润统计
	ProfitTrade   int64 `json:"profit_trade"`
	ProfitDeposit int64 `json:"profit_deposit"`
	ProfitSim     int64 `json:"profit_sim"`
	ProfitReward  int64 `json:"profit_reward"`
	ProfitTotal   int64 `json:"profit_total"`

	// 商户与终端统计
	MerchantCount     int `json:"merchant_count"`
	MerchantNew       int `json:"merchant_new"`
	TerminalTotal     int `json:"terminal_total"`
	TerminalActivated int `json:"terminal_activated"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (AgentMonthlyStats) TableName() string {
	return "agent_monthly_stats"
}

// 统计范围常量
const (
	StatScopeDirect = "direct" // 直营
	StatScopeTeam   = "team"   // 团队
)

// 分润类型常量(用于profit_records表的profit_type字段)
const (
	ProfitTypeTrade   int16 = 1 // 交易分润
	ProfitTypeDeposit int16 = 2 // 押金返现
	ProfitTypeSim     int16 = 3 // 流量返现
	ProfitTypeReward  int16 = 4 // 激活奖励
)

// OverviewData 首页概览数据结构
type OverviewData struct {
	Today     *DayStats      `json:"today"`
	Yesterday *DayStats      `json:"yesterday"`
	Week      *PeriodStats   `json:"week"`
	Month     *PeriodStats   `json:"month"`
	Team      *TeamStats     `json:"team"`
	Terminal  *TerminalStats `json:"terminal"`
	Wallet    *WalletStats   `json:"wallet"`
}

// DayStats 每日统计
type DayStats struct {
	TransAmount   int64   `json:"trans_amount"`
	TransAmountYuan float64 `json:"trans_amount_yuan"`
	TransCount    int     `json:"trans_count"`
	ProfitTotal   int64   `json:"profit_total"`
	ProfitTotalYuan float64 `json:"profit_total_yuan"`
	ProfitTrade   int64   `json:"profit_trade"`
	ProfitDeposit int64   `json:"profit_deposit"`
	ProfitSim     int64   `json:"profit_sim"`
	ProfitReward  int64   `json:"profit_reward"`
}

// PeriodStats 时间段统计
type PeriodStats struct {
	TransAmount   int64   `json:"trans_amount"`
	TransAmountYuan float64 `json:"trans_amount_yuan"`
	TransCount    int     `json:"trans_count"`
	ProfitTotal   int64   `json:"profit_total"`
	ProfitTotalYuan float64 `json:"profit_total_yuan"`
	MerchantNew   int     `json:"merchant_new"`
}

// TeamStats 团队统计
type TeamStats struct {
	DirectAgentCount    int `json:"direct_agent_count"`
	DirectMerchantCount int `json:"direct_merchant_count"`
	TeamAgentCount      int `json:"team_agent_count"`
	TeamMerchantCount   int `json:"team_merchant_count"`
}

// TerminalStats 终端统计
type TerminalStats struct {
	Total          int `json:"total"`
	Activated      int `json:"activated"`
	TodayActivated int `json:"today_activated"`
	MonthActivated int `json:"month_activated"`
}

// WalletStats 钱包统计
type WalletStats struct {
	TotalBalance     int64   `json:"total_balance"`
	TotalBalanceYuan float64 `json:"total_balance_yuan"`
}

// ChannelStats 通道统计
type ChannelStats struct {
	ChannelID   int64   `json:"channel_id"`
	ChannelCode string  `json:"channel_code"`
	ChannelName string  `json:"channel_name"`
	TransAmount int64   `json:"trans_amount"`
	TransCount  int     `json:"trans_count"`
	Percentage  float64 `json:"percentage"` // 占比百分比
}

// MerchantDistribution 商户类型分布
type MerchantDistribution struct {
	MerchantType string `json:"merchant_type"`
	TypeName     string `json:"type_name"`
	Count        int    `json:"count"`
	Percentage   float64 `json:"percentage"`
}

// AgentRanking 代理商排名
type AgentRanking struct {
	Rank        int     `json:"rank"`
	AgentID     int64   `json:"agent_id"`
	AgentName   string  `json:"agent_name"`
	AgentNo     string  `json:"agent_no"`
	Value       int64   `json:"value"`       // 排名值(金额或数量)
	ValueYuan   float64 `json:"value_yuan"`  // 金额(元)
	Change      int64   `json:"change"`      // 较上期变化
	ChangeRate  float64 `json:"change_rate"` // 变化率(%)
}

// MerchantRanking 商户排名
type MerchantRanking struct {
	Rank           int     `json:"rank"`
	MerchantID     int64   `json:"merchant_id"`
	MerchantName   string  `json:"merchant_name"`   // 脱敏显示
	MerchantType   string  `json:"merchant_type"`
	TotalAmount    int64   `json:"total_amount"`    // 累计交易额
	MonthAmount    int64   `json:"month_amount"`    // 本月交易额
	MonthAmountYuan float64 `json:"month_amount_yuan"`
}

// TrendPoint 趋势图数据点
type TrendPoint struct {
	Date        string  `json:"date"`
	TransAmount int64   `json:"trans_amount"`
	TransAmountYuan float64 `json:"trans_amount_yuan"`
	TransCount  int     `json:"trans_count"`
	ProfitTotal int64   `json:"profit_total"`
	ProfitTotalYuan float64 `json:"profit_total_yuan"`
}

// RecentTransaction 最近交易
type RecentTransaction struct {
	ID           int64     `json:"id"`
	MerchantName string    `json:"merchant_name"`
	PayType      int16     `json:"pay_type"`
	PayTypeName  string    `json:"pay_type_name"`
	Amount       int64     `json:"amount"`
	AmountYuan   float64   `json:"amount_yuan"`
	TradeTime    time.Time `json:"trade_time"`
	TimeAgo      string    `json:"time_ago"` // 如 "10分钟前"
}
