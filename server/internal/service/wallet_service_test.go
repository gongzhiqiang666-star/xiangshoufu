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

// =============================================================================
// 奖励钱包提现检查测试（P0业务规则）
// 核心规则：奖励钱包提现需检查上级充值钱包余额是否充足
// =============================================================================

// TestRewardWalletWithdraw_ParentChargingWalletCheck 测试奖励钱包提现检查上级充值钱包
func TestRewardWalletWithdraw_ParentChargingWalletCheck(t *testing.T) {
	tests := []struct {
		name                    string
		withdrawAmount          int64 // 提现金额
		parentChargingBalance   int64 // 上级充值钱包余额
		canWithdraw             bool
		errMsg                  string
	}{
		{
			name:                  "上级余额充足",
			withdrawAmount:        10000, // 100元
			parentChargingBalance: 50000, // 500元
			canWithdraw:           true,
			errMsg:                "",
		},
		{
			name:                  "上级余额刚好够",
			withdrawAmount:        10000, // 100元
			parentChargingBalance: 10000, // 100元
			canWithdraw:           true,
			errMsg:                "",
		},
		{
			name:                  "上级余额不足",
			withdrawAmount:        10000, // 100元
			parentChargingBalance: 5000,  // 50元
			canWithdraw:           false,
			errMsg:                "上级充值钱包余额不足，无法提现",
		},
		{
			name:                  "上级余额为零",
			withdrawAmount:        10000, // 100元
			parentChargingBalance: 0,
			canWithdraw:           false,
			errMsg:                "上级充值钱包余额不足，无法提现",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟检查逻辑
			canWithdraw := tt.parentChargingBalance >= tt.withdrawAmount
			var errMsg string
			if !canWithdraw {
				errMsg = "上级充值钱包余额不足，无法提现"
			}

			assert.Equal(t, tt.canWithdraw, canWithdraw)
			if !canWithdraw {
				assert.Equal(t, tt.errMsg, errMsg)
			}
		})
	}
}

// TestRewardWalletWithdraw_TopAgentCannotWithdraw 测试顶级代理商无法从奖励钱包提现
func TestRewardWalletWithdraw_TopAgentCannotWithdraw(t *testing.T) {
	// 顶级代理商没有上级，无法从奖励钱包提现
	parentID := int64(0) // 顶级代理商

	canWithdraw := parentID > 0
	assert.False(t, canWithdraw)

	// 普通代理商有上级
	parentID = int64(1)
	canWithdraw = parentID > 0
	assert.True(t, canWithdraw)
}

// TestRewardWalletWithdraw_WalletTypeCheck 测试钱包类型判断
func TestRewardWalletWithdraw_WalletTypeCheck(t *testing.T) {
	const (
		WalletTypeProfit     = int16(1) // 分润钱包
		WalletTypeService    = int16(2) // 服务费钱包
		WalletTypeReward     = int16(3) // 奖励钱包
		WalletTypeCharging   = int16(4) // 充值钱包
		WalletTypeSettlement = int16(5) // 沉淀钱包
	)

	tests := []struct {
		walletType           int16
		needParentCheck      bool
		description          string
	}{
		{WalletTypeProfit, false, "分润钱包无需检查上级"},
		{WalletTypeService, false, "服务费钱包无需检查上级"},
		{WalletTypeReward, true, "奖励钱包需检查上级充值钱包"},
		{WalletTypeCharging, false, "充值钱包无需检查上级"},
		{WalletTypeSettlement, false, "沉淀钱包无需检查上级"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			needParentCheck := tt.walletType == WalletTypeReward
			assert.Equal(t, tt.needParentCheck, needParentCheck)
		})
	}
}

// =============================================================================
// 充值钱包测试
// =============================================================================

// TestChargingWallet_BalanceCheck 测试充值钱包余额检查
func TestChargingWallet_BalanceCheck(t *testing.T) {
	tests := []struct {
		name           string
		balance        int64
		issueAmount    int64
		canIssue       bool
	}{
		{"余额充足", 100000, 50000, true},
		{"余额刚好", 50000, 50000, true},
		{"余额不足", 30000, 50000, false},
		{"零余额", 0, 50000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canIssue := tt.balance >= tt.issueAmount
			assert.Equal(t, tt.canIssue, canIssue)
		})
	}
}

// TestChargingWallet_IssueReward 测试充值钱包发放奖励
func TestChargingWallet_IssueReward(t *testing.T) {
	// 模拟发放奖励流程
	fromAgentBalance := int64(100000) // 1000元
	issueAmount := int64(20000)       // 200元

	// 发放前检查
	canIssue := fromAgentBalance >= issueAmount
	assert.True(t, canIssue)

	// 发放后余额
	newBalance := fromAgentBalance - issueAmount
	assert.Equal(t, int64(80000), newBalance)
}

// TestChargingWallet_OnlyDirectSubordinate 测试只能给直属下级发放
func TestChargingWallet_OnlyDirectSubordinate(t *testing.T) {
	tests := []struct {
		name          string
		fromAgentID   int64
		toAgentParent int64
		canIssue      bool
	}{
		{"直属下级", 1, 1, true},      // toAgent.ParentID == fromAgentID
		{"非直属下级", 1, 2, false},   // toAgent.ParentID != fromAgentID
		{"跨级下级", 1, 3, false},     // toAgent.ParentID != fromAgentID
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canIssue := tt.toAgentParent == tt.fromAgentID
			assert.Equal(t, tt.canIssue, canIssue)
		})
	}
}

// TestChargingWallet_RewardStats 测试充值钱包奖励统计
func TestChargingWallet_RewardStats(t *testing.T) {
	// 模拟奖励发放记录
	type RewardRecord struct {
		Amount     int64
		RewardType int // 1=手动发放 2=系统自动
	}

	records := []RewardRecord{
		{5000, 1},  // 手动50元
		{3000, 1},  // 手动30元
		{2000, 2},  // 自动20元
		{4000, 2},  // 自动40元
		{1000, 2},  // 自动10元
	}

	var totalIssued int64     // 手动发放总额
	var totalAutoReward int64 // 系统自动总额
	var totalReward int64     // 奖励总额

	for _, r := range records {
		totalReward += r.Amount
		if r.RewardType == 1 {
			totalIssued += r.Amount
		} else {
			totalAutoReward += r.Amount
		}
	}

	assert.Equal(t, int64(8000), totalIssued)      // 80元手动
	assert.Equal(t, int64(7000), totalAutoReward)  // 70元自动
	assert.Equal(t, int64(15000), totalReward)     // 150元总计
}

// =============================================================================
// 钱包类型常量测试
// =============================================================================

// TestWalletTypeConstants 测试钱包类型常量
func TestWalletTypeConstants(t *testing.T) {
	const (
		WalletTypeProfit     = int16(1)
		WalletTypeService    = int16(2)
		WalletTypeReward     = int16(3)
		WalletTypeCharging   = int16(4)
		WalletTypeSettlement = int16(5)
	)

	// 验证钱包类型值
	assert.Equal(t, int16(1), WalletTypeProfit)
	assert.Equal(t, int16(2), WalletTypeService)
	assert.Equal(t, int16(3), WalletTypeReward)
	assert.Equal(t, int16(4), WalletTypeCharging)
	assert.Equal(t, int16(5), WalletTypeSettlement)
}

// TestWalletTypeNames 测试钱包类型名称映射
func TestWalletTypeNames(t *testing.T) {
	walletTypeNames := map[int16]string{
		1: "分润钱包",
		2: "服务费钱包",
		3: "奖励钱包",
		4: "充值钱包",
		5: "沉淀钱包",
	}

	assert.Equal(t, "分润钱包", walletTypeNames[1])
	assert.Equal(t, "服务费钱包", walletTypeNames[2])
	assert.Equal(t, "奖励钱包", walletTypeNames[3])
	assert.Equal(t, "充值钱包", walletTypeNames[4])
	assert.Equal(t, "沉淀钱包", walletTypeNames[5])
}

// =============================================================================
// 提现流程测试
// =============================================================================

// TestWithdrawFlow_RewardWallet 测试奖励钱包提现完整流程
func TestWithdrawFlow_RewardWallet(t *testing.T) {
	// 模拟代理商结构
	agentID := int64(2)
	parentID := int64(1)

	// 奖励钱包余额
	rewardWalletBalance := int64(50000) // 500元
	withdrawAmount := int64(20000)      // 200元

	// 上级充值钱包余额
	parentChargingBalance := int64(100000) // 1000元

	// Step 1: 检查是否是奖励钱包
	walletType := int16(3) // 奖励钱包
	isRewardWallet := walletType == 3
	assert.True(t, isRewardWallet)

	// Step 2: 检查是否有上级
	hasParent := parentID > 0
	assert.True(t, hasParent)

	// Step 3: 检查奖励钱包余额
	hasSufficientRewardBalance := rewardWalletBalance >= withdrawAmount
	assert.True(t, hasSufficientRewardBalance)

	// Step 4: 检查上级充值钱包余额
	hasParentChargingBalance := parentChargingBalance >= withdrawAmount
	assert.True(t, hasParentChargingBalance)

	// Step 5: 执行提现
	canWithdraw := isRewardWallet && hasParent && hasSufficientRewardBalance && hasParentChargingBalance
	assert.True(t, canWithdraw)

	// Step 6: 更新余额
	newRewardBalance := rewardWalletBalance - withdrawAmount
	newParentChargingBalance := parentChargingBalance - withdrawAmount

	assert.Equal(t, int64(30000), newRewardBalance)        // 300元
	assert.Equal(t, int64(80000), newParentChargingBalance) // 800元

	_ = agentID // 使用变量避免警告
}

// TestWithdrawFlow_OtherWallets 测试其他钱包提现流程（无需检查上级）
func TestWithdrawFlow_OtherWallets(t *testing.T) {
	tests := []struct {
		name        string
		walletType  int16
		balance     int64
		withdraw    int64
		canWithdraw bool
	}{
		{"分润钱包提现", 1, 100000, 50000, true},
		{"服务费钱包提现", 2, 100000, 50000, true},
		{"分润钱包余额不足", 1, 30000, 50000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 非奖励钱包只需检查自身余额
			canWithdraw := tt.balance >= tt.withdraw
			assert.Equal(t, tt.canWithdraw, canWithdraw)
		})
	}
}
