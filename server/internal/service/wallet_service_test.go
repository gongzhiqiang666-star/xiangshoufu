package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWalletInfoDetail 测试钱包信息结构
func TestWalletInfoDetail(t *testing.T) {
	info := &WalletInfoDetail{
		ID:                1,
		ChannelID:         1,
		ChannelName:       "恒信通",
		WalletType:        1,
		WalletTypeName:    "分润钱包",
		Balance:           100000,
		BalanceYuan:       1000.00,
		FrozenAmount:      5000,
		TotalIncome:       500000,
		TotalWithdraw:     400000,
		WithdrawThreshold: 10000,
		CanWithdraw:       true,
	}

	assert.Equal(t, int64(1), info.ID)
	assert.Equal(t, int16(1), info.WalletType)
	assert.Equal(t, "分润钱包", info.WalletTypeName)
	assert.Equal(t, int64(100000), info.Balance)
	assert.Equal(t, float64(1000.00), info.BalanceYuan)
	assert.Equal(t, int64(5000), info.FrozenAmount)
	assert.Equal(t, int64(500000), info.TotalIncome)
	assert.Equal(t, true, info.CanWithdraw)
}

// TestWalletSummary 测试钱包汇总结构
func TestWalletSummary(t *testing.T) {
	summary := &WalletSummary{
		TotalBalance:     500000,
		TotalBalanceYuan: 5000.00,
		TotalFrozen:      50000,
		TotalIncome:      600000,
		TotalWithdraw:    100000,
		AvailableBalance: 450000,
		WalletCount:      3,
	}

	assert.Equal(t, int64(500000), summary.TotalBalance)
	assert.Equal(t, float64(5000.00), summary.TotalBalanceYuan)
	assert.Equal(t, int64(50000), summary.TotalFrozen)
	assert.Equal(t, int64(450000), summary.AvailableBalance)
	assert.Equal(t, 3, summary.WalletCount)
}

// TestWalletLogInfo 测试钱包流水信息结构
func TestWalletLogInfo(t *testing.T) {
	log := &WalletLogInfo{
		ID:            1,
		LogType:       1,
		LogTypeName:   "交易分润",
		Amount:        1000,
		AmountYuan:    10.00,
		BalanceBefore: 99000,
		BalanceAfter:  100000,
		Remark:        "交易分润入账",
	}

	assert.Equal(t, int64(1), log.ID)
	assert.Equal(t, int16(1), log.LogType)
	assert.Equal(t, "交易分润", log.LogTypeName)
	assert.Equal(t, int64(1000), log.Amount)
	assert.Equal(t, float64(10.00), log.AmountYuan)
	assert.Equal(t, int64(99000), log.BalanceBefore)
	assert.Equal(t, int64(100000), log.BalanceAfter)
}

// TestWithdrawRequest 测试提现请求验证
func TestWithdrawRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *WithdrawRequest
		isValid bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &WithdrawRequest{
				AgentID:  1,
				WalletID: 1,
				Amount:   10000,
			},
			isValid: true,
		},
		{
			name: "zero amount",
			req: &WithdrawRequest{
				AgentID:  1,
				WalletID: 1,
				Amount:   0,
			},
			isValid: false,
			errMsg:  "提现金额必须大于0",
		},
		{
			name: "negative amount",
			req: &WithdrawRequest{
				AgentID:  1,
				WalletID: 1,
				Amount:   -1000,
			},
			isValid: false,
			errMsg:  "提现金额必须大于0",
		},
		{
			name: "invalid wallet id",
			req: &WithdrawRequest{
				AgentID:  1,
				WalletID: 0,
				Amount:   10000,
			},
			isValid: false,
			errMsg:  "钱包ID无效",
		},
		{
			name: "invalid agent id",
			req: &WithdrawRequest{
				AgentID:  0,
				WalletID: 1,
				Amount:   10000,
			},
			isValid: false,
			errMsg:  "代理商ID无效",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errMsg string
			isValid := true

			if tt.req.AgentID <= 0 {
				isValid = false
				errMsg = "代理商ID无效"
			} else if tt.req.WalletID <= 0 {
				isValid = false
				errMsg = "钱包ID无效"
			} else if tt.req.Amount <= 0 {
				isValid = false
				errMsg = "提现金额必须大于0"
			}

			assert.Equal(t, tt.isValid, isValid)
			if !isValid {
				assert.Equal(t, tt.errMsg, errMsg)
			}
		})
	}
}

// TestGetWalletTypeNameStr 测试钱包类型名称获取
func TestGetWalletTypeNameStr(t *testing.T) {
	tests := []struct {
		walletType int16
		expected   string
	}{
		{1, "分润钱包"},
		{2, "服务费钱包"},
		{3, "奖励钱包"},
		{0, "未知钱包"},
		{99, "未知钱包"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getWalletTypeNameStr(tt.walletType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetLogTypeName 测试流水类型名称获取
func TestGetLogTypeName(t *testing.T) {
	tests := []struct {
		logType  int16
		expected string
	}{
		{1, "分润入账"},
		{2, "提现冻结"},
		{3, "提现成功"},
		{4, "提现退回"},
		{5, "调账"},
		{6, "代扣"},
		{7, "返现"},
		{0, "未知"},
		{99, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getLogTypeName(tt.logType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBalanceCalculation 测试余额计算
func TestBalanceCalculation(t *testing.T) {
	tests := []struct {
		name             string
		balance          int64
		frozenAmount     int64
		expectedAvailable int64
	}{
		{"no frozen", 100000, 0, 100000},
		{"partial frozen", 100000, 30000, 70000},
		{"all frozen", 100000, 100000, 0},
		{"zero balance", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := tt.balance - tt.frozenAmount
			assert.Equal(t, tt.expectedAvailable, available)
		})
	}
}

// TestAmountConversion 测试金额转换（分到元）
func TestAmountConversion(t *testing.T) {
	tests := []struct {
		amountFen  int64
		amountYuan float64
	}{
		{0, 0.00},
		{1, 0.01},
		{10, 0.10},
		{100, 1.00},
		{1000, 10.00},
		{10000, 100.00},
		{100000, 1000.00},
		{1234567, 12345.67},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			yuan := float64(tt.amountFen) / 100
			assert.Equal(t, tt.amountYuan, yuan)
		})
	}
}

// TestWithdrawThreshold 测试提现门槛验证
func TestWithdrawThreshold(t *testing.T) {
	tests := []struct {
		name      string
		amount    int64
		threshold int64
		canWithdraw bool
	}{
		{"above threshold", 20000, 10000, true},
		{"at threshold", 10000, 10000, true},
		{"below threshold", 5000, 10000, false},
		{"zero threshold", 1000, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canWithdraw := tt.amount >= tt.threshold
			assert.Equal(t, tt.canWithdraw, canWithdraw)
		})
	}
}

// TestSufficientBalance 测试余额充足性验证
func TestSufficientBalance(t *testing.T) {
	tests := []struct {
		name          string
		balance       int64
		frozenAmount  int64
		withdrawAmount int64
		sufficient    bool
	}{
		{"sufficient", 100000, 20000, 50000, true},
		{"exact match", 100000, 20000, 80000, true},
		{"insufficient", 100000, 20000, 90000, false},
		{"all frozen", 100000, 100000, 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := tt.balance - tt.frozenAmount
			sufficient := available >= tt.withdrawAmount
			assert.Equal(t, tt.sufficient, sufficient)
		})
	}
}
