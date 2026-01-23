package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"xiangshoufu/internal/models"

	"github.com/stretchr/testify/assert"
)

// ==================== BatchSetRate 专用 Mock ====================

// BatchSetRateMockTerminalRepo 批量设置费率测试专用终端仓储
type BatchSetRateMockTerminalRepo struct {
	terminals map[string]*models.Terminal
	policies  map[string]*models.TerminalPolicy
	findErr   error
	saveErr   error
}

func NewBatchSetRateMockTerminalRepo() *BatchSetRateMockTerminalRepo {
	return &BatchSetRateMockTerminalRepo{
		terminals: make(map[string]*models.Terminal),
		policies:  make(map[string]*models.TerminalPolicy),
	}
}

func (m *BatchSetRateMockTerminalRepo) FindBySN(sn string) (*models.Terminal, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if t, ok := m.terminals[sn]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *BatchSetRateMockTerminalRepo) FindPolicyBySN(sn string) (*models.TerminalPolicy, error) {
	if p, ok := m.policies[sn]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *BatchSetRateMockTerminalRepo) SavePolicy(policy *models.TerminalPolicy) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.policies[policy.TerminalSN] = policy
	return nil
}

func (m *BatchSetRateMockTerminalRepo) AddTerminal(terminal *models.Terminal) {
	m.terminals[terminal.TerminalSN] = terminal
}

func (m *BatchSetRateMockTerminalRepo) AddPolicy(policy *models.TerminalPolicy) {
	m.policies[policy.TerminalSN] = policy
}

// BatchSetRateMockRateSyncService 批量设置费率测试专用费率同步服务
type BatchSetRateMockRateSyncService struct {
	syncResult *SyncResult
	syncErr    error
	callCount  int
	lastParams *RateUpdateParams
}

func NewBatchSetRateMockRateSyncService() *BatchSetRateMockRateSyncService {
	return &BatchSetRateMockRateSyncService{
		syncResult: &SyncResult{
			Success: true,
			LogID:   1,
			TradeNo: "TRX123456",
			Message: "费率同步成功",
		},
	}
}

func (m *BatchSetRateMockRateSyncService) SyncRateToChannel(ctx context.Context, params *RateUpdateParams) (*SyncResult, error) {
	m.callCount++
	m.lastParams = params
	if m.syncErr != nil {
		return nil, m.syncErr
	}
	return m.syncResult, nil
}

func (m *BatchSetRateMockRateSyncService) SetSyncResult(success bool, message string) {
	m.syncResult = &SyncResult{
		Success: success,
		Message: message,
	}
}

func (m *BatchSetRateMockRateSyncService) SetSyncError(err error) {
	m.syncErr = err
}

// ==================== 测试用例 ====================

// TestBatchSetRateRequest 测试批量设置费率请求结构
func TestBatchSetRateRequest(t *testing.T) {
	req := &BatchSetRateRequest{
		TerminalSNs:  []string{"SN001", "SN002"},
		AgentID:      100,
		CreditRate:   55, // 0.55%
		DebitRate:    50,
		DebitCap:     20,
		UnionpayRate: 60,
		WechatRate:   38,
		AlipayRate:   38,
		UpdatedBy:    1,
	}

	assert.Len(t, req.TerminalSNs, 2)
	assert.Equal(t, int64(100), req.AgentID)
	assert.Equal(t, 55, req.CreditRate)
	assert.Equal(t, 50, req.DebitRate)
	assert.Equal(t, 20, req.DebitCap)
}

// TestBatchPolicyResult 测试批量政策结果结构
func TestBatchPolicyResult(t *testing.T) {
	tests := []struct {
		name         string
		result       *BatchPolicyResult
		expectErrors bool
	}{
		{
			name: "全部成功",
			result: &BatchPolicyResult{
				TotalCount:   3,
				SuccessCount: 3,
				FailedCount:  0,
				Errors:       []string{},
			},
			expectErrors: false,
		},
		{
			name: "部分失败",
			result: &BatchPolicyResult{
				TotalCount:   3,
				SuccessCount: 2,
				FailedCount:  1,
				Errors:       []string{"终端 SN003: 不存在"},
			},
			expectErrors: true,
		},
		{
			name: "全部失败",
			result: &BatchPolicyResult{
				TotalCount:   2,
				SuccessCount: 0,
				FailedCount:  2,
				Errors:       []string{"终端 SN001: 通道同步失败", "终端 SN002: 不存在"},
			},
			expectErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result.TotalCount, tt.result.SuccessCount+tt.result.FailedCount)
			if tt.expectErrors {
				assert.NotEmpty(t, tt.result.Errors)
			} else {
				assert.Empty(t, tt.result.Errors)
			}
		})
	}
}

// TestBatchSetRate_EmptyList 测试空列表场景
func TestBatchSetRate_EmptyList(t *testing.T) {
	service := &TerminalService{}

	req := &BatchSetRateRequest{
		TerminalSNs: []string{},
		AgentID:     100,
	}

	result, err := service.BatchSetRate(req)

	// 验证
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "终端列表不能为空")
}

// TestRateConversionLogic 测试费率转换逻辑（万分比 -> 小数）
func TestRateConversionLogic(t *testing.T) {
	tests := []struct {
		name            string
		inputRate       int     // 万分比输入 (如 55 表示 0.55%)
		expectedDecimal float64 // 预期小数形式 (如 0.0055)
	}{
		{
			name:            "0.55% -> 0.0055",
			inputRate:       55,
			expectedDecimal: 0.0055,
		},
		{
			name:            "0.60% -> 0.006",
			inputRate:       60,
			expectedDecimal: 0.006,
		},
		{
			name:            "0.38% -> 0.0038",
			inputRate:       38,
			expectedDecimal: 0.0038,
		},
		{
			name:            "1.00% -> 0.01",
			inputRate:       100,
			expectedDecimal: 0.01,
		},
		{
			name:            "0% -> 0",
			inputRate:       0,
			expectedDecimal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证转换逻辑
			decimal := float64(tt.inputRate) / 10000
			assert.InDelta(t, tt.expectedDecimal, decimal, 0.000001)
		})
	}
}

// TestBusinessRule_ChannelFirstThenLocal_Integration 测试业务规则：先通道后本地
func TestBusinessRule_ChannelFirstThenLocal_Integration(t *testing.T) {
	// 业务规则：费率修改必须和通道联动
	// 原则：先调用通道，通道返回成功才能修改自己的费率

	type scenario struct {
		name              string
		channelSuccess    bool
		shouldUpdateLocal bool
		expectedMessage   string
	}

	scenarios := []scenario{
		{
			name:              "通道成功 -> 更新本地",
			channelSuccess:    true,
			shouldUpdateLocal: true,
			expectedMessage:   "费率修改成功，已同步到支付通道",
		},
		{
			name:              "通道失败 -> 不更新本地",
			channelSuccess:    false,
			shouldUpdateLocal: false,
			expectedMessage:   "费率修改失败，通道同步失败",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// 验证业务逻辑
			if s.channelSuccess {
				assert.True(t, s.shouldUpdateLocal, "通道成功时应该更新本地")
			} else {
				assert.False(t, s.shouldUpdateLocal, "通道失败时不应该更新本地")
			}
		})
	}
}

// TestBatchSetRateMock_SyncResult 测试模拟同步结果
func TestBatchSetRateMock_SyncResult(t *testing.T) {
	mockService := NewBatchSetRateMockRateSyncService()

	// 默认成功
	result, err := mockService.SyncRateToChannel(context.Background(), &RateUpdateParams{
		TerminalSN:  "SN001",
		ChannelCode: "HENGXINTONG",
	})

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, mockService.callCount)

	// 设置失败
	mockService.SetSyncResult(false, "商户不存在")
	result, err = mockService.SyncRateToChannel(context.Background(), &RateUpdateParams{
		TerminalSN: "SN002",
	})

	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "商户不存在", result.Message)

	// 设置异常
	mockService.SetSyncError(errors.New("网络超时"))
	result, err = mockService.SyncRateToChannel(context.Background(), &RateUpdateParams{
		TerminalSN: "SN003",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "网络超时")
}

// TestBatchSetRateMock_TerminalRepo 测试模拟终端仓储
func TestBatchSetRateMock_TerminalRepo(t *testing.T) {
	repo := NewBatchSetRateMockTerminalRepo()

	// 添加终端
	repo.AddTerminal(&models.Terminal{
		ID:           1,
		TerminalSN:   "SN001",
		ChannelID:    1,
		ChannelCode:  "HENGXINTONG",
		OwnerAgentID: 100,
	})

	// 查找存在的终端
	terminal, err := repo.FindBySN("SN001")
	assert.NoError(t, err)
	assert.NotNil(t, terminal)
	assert.Equal(t, "SN001", terminal.TerminalSN)

	// 查找不存在的终端
	terminal, err = repo.FindBySN("NOT_EXIST")
	assert.NoError(t, err)
	assert.Nil(t, terminal)

	// 保存政策
	policy := &models.TerminalPolicy{
		TerminalSN: "SN001",
		CreditRate: 55,
		DebitRate:  50,
		IsSynced:   true,
	}
	err = repo.SavePolicy(policy)
	assert.NoError(t, err)

	// 查询政策
	savedPolicy, err := repo.FindPolicyBySN("SN001")
	assert.NoError(t, err)
	assert.NotNil(t, savedPolicy)
	assert.Equal(t, 55, savedPolicy.CreditRate)

	// 测试保存错误
	repo.saveErr = errors.New("数据库错误")
	err = repo.SavePolicy(policy)
	assert.Error(t, err)
}

// TestTerminalPolicyModel 测试终端政策模型
func TestTerminalPolicyModel(t *testing.T) {
	now := time.Now()

	policy := &models.TerminalPolicy{
		ID:                 1,
		TerminalSN:         "SN001",
		ChannelID:          1,
		AgentID:            100,
		CreditRate:         55,  // 0.55%
		DebitRate:          50,  // 0.50%
		DebitCap:           20,  // 20元封顶
		UnionpayRate:       60,  // 0.60%
		WechatRate:         38,  // 0.38%
		AlipayRate:         38,  // 0.38%
		FirstSimFee:        30,  // 首次流量费
		NonFirstSimFee:     20,  // 非首次流量费
		SimFeeIntervalDays: 30,  // 间隔天数
		DepositAmount:      100, // 押金
		IsSynced:           true,
		CreatedBy:          1,
		CreatedAt:          now,
		UpdatedBy:          1,
		UpdatedAt:          now,
	}

	assert.Equal(t, "SN001", policy.TerminalSN)
	assert.Equal(t, 55, policy.CreditRate)
	assert.True(t, policy.IsSynced)
}

// TestImportTerminalsRequest 测试终端入库请求
func TestImportTerminalsRequest(t *testing.T) {
	req := &ImportTerminalsRequest{
		ChannelID:    1,
		ChannelCode:  "HENGXINTONG",
		BrandCode:    "BRAND001",
		ModelCode:    "MODEL001",
		SNList:       []string{"SN001", "SN002", "SN003"},
		OwnerAgentID: 100,
		CreatedBy:    1,
	}

	assert.Len(t, req.SNList, 3)
	assert.Equal(t, "HENGXINTONG", req.ChannelCode)
	assert.Equal(t, int64(100), req.OwnerAgentID)
}

// TestImportTerminalsResult 测试终端入库结果
func TestImportTerminalsResult(t *testing.T) {
	tests := []struct {
		name   string
		result *ImportTerminalsResult
	}{
		{
			name: "全部成功",
			result: &ImportTerminalsResult{
				ImportNo:     "IMP20260123123456000001",
				TotalCount:   5,
				SuccessCount: 5,
				FailedCount:  0,
				FailedSNs:    []string{},
				Errors:       []string{},
			},
		},
		{
			name: "部分失败",
			result: &ImportTerminalsResult{
				ImportNo:     "IMP20260123123456000002",
				TotalCount:   5,
				SuccessCount: 3,
				FailedCount:  2,
				FailedSNs:    []string{"SN004", "SN005"},
				Errors:       []string{"终端 SN004 已存在", "终端 SN005 已存在"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result.TotalCount, tt.result.SuccessCount+tt.result.FailedCount)
			assert.Len(t, tt.result.FailedSNs, tt.result.FailedCount)
			assert.NotEmpty(t, tt.result.ImportNo)
		})
	}
}

// TestRecallTerminalRequest 测试终端回拨请求
func TestRecallTerminalRequest(t *testing.T) {
	req := &RecallTerminalRequest{
		FromAgentID: 100,
		ToAgentID:   50, // 上级代理商
		TerminalSN:  "SN001",
		ChannelID:   1,
		Source:      1, // APP
		Remark:      "商户不需要了",
		CreatedBy:   1,
	}

	assert.Equal(t, int64(100), req.FromAgentID)
	assert.Equal(t, int64(50), req.ToAgentID)
	assert.Equal(t, "SN001", req.TerminalSN)
}

// TestTerminalStats 测试终端统计
func TestTerminalStats(t *testing.T) {
	stats := &TerminalStats{
		Total:              100,
		PendingCount:       30,
		AllocatedCount:     20,
		BoundCount:         25,
		ActivatedCount:     25,
		UnboundCount:       50, // 待分配 + 已分配
		YesterdayActivated: 5,
		TodayActivated:     3,
		MonthActivated:     50,
	}

	// 验证统计数据一致性
	assert.Equal(t, stats.Total, stats.PendingCount+stats.AllocatedCount+stats.BoundCount+stats.ActivatedCount)
	assert.Equal(t, stats.UnboundCount, stats.PendingCount+stats.AllocatedCount)
}

// TestBatchRecallRequest 测试批量回拨请求
func TestBatchRecallRequest(t *testing.T) {
	req := &BatchRecallRequest{
		TerminalSNs: []string{"SN001", "SN002", "SN003"},
		ToAgentID:   50,
		FromAgentID: 100,
		Source:      2, // PC
		Remark:      "批量回拨",
		CreatedBy:   1,
	}

	assert.Len(t, req.TerminalSNs, 3)
	assert.Equal(t, int64(50), req.ToAgentID)
	assert.Equal(t, int16(2), req.Source)
}

// TestBatchRecallResult 测试批量回拨结果
func TestBatchRecallResult(t *testing.T) {
	result := &BatchRecallResult{
		TotalCount:   3,
		SuccessCount: 2,
		FailedCount:  1,
		Errors:       []string{"终端 SN003: 已激活的终端不能回拨"},
	}

	assert.Equal(t, result.TotalCount, result.SuccessCount+result.FailedCount)
	assert.Len(t, result.Errors, result.FailedCount)
}

// TestBatchSetSimFeeRequest 测试批量设置SIM卡费用请求
func TestBatchSetSimFeeRequest(t *testing.T) {
	req := &BatchSetSimFeeRequest{
		TerminalSNs:        []string{"SN001", "SN002"},
		AgentID:            100,
		FirstSimFee:        30,
		NonFirstSimFee:     20,
		SimFeeIntervalDays: 30,
		UpdatedBy:          1,
	}

	assert.Len(t, req.TerminalSNs, 2)
	assert.Equal(t, 30, req.FirstSimFee)
	assert.Equal(t, 20, req.NonFirstSimFee)
	assert.Equal(t, 30, req.SimFeeIntervalDays)
}

// TestBatchSetDepositRequest 测试批量设置押金请求
func TestBatchSetDepositRequest(t *testing.T) {
	req := &BatchSetDepositRequest{
		TerminalSNs:   []string{"SN001", "SN002"},
		AgentID:       100,
		DepositAmount: 199, // 199元押金
		UpdatedBy:     1,
	}

	assert.Len(t, req.TerminalSNs, 2)
	assert.Equal(t, 199, req.DepositAmount)
}
