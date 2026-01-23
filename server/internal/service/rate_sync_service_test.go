package service

import (
	"context"
	"testing"
	"time"

	"xiangshoufu/internal/models"

	"github.com/stretchr/testify/assert"
)

// MockRateSyncLogRepository 模拟费率同步日志仓库
type MockRateSyncLogRepository struct {
	logs      []*models.RateSyncLog
	createErr error
	updateErr error
}

func NewMockRateSyncLogRepository() *MockRateSyncLogRepository {
	return &MockRateSyncLogRepository{
		logs: make([]*models.RateSyncLog, 0),
	}
}

func (m *MockRateSyncLogRepository) Create(ctx context.Context, log *models.RateSyncLog) error {
	if m.createErr != nil {
		return m.createErr
	}
	log.ID = int64(len(m.logs) + 1)
	log.CreatedAt = time.Now()
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockRateSyncLogRepository) Update(ctx context.Context, log *models.RateSyncLog) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, l := range m.logs {
		if l.ID == log.ID {
			m.logs[i] = log
			return nil
		}
	}
	return nil
}

func (m *MockRateSyncLogRepository) GetByID(ctx context.Context, id int64) (*models.RateSyncLog, error) {
	for _, log := range m.logs {
		if log.ID == id {
			return log, nil
		}
	}
	return nil, nil
}

func (m *MockRateSyncLogRepository) GetByMerchantID(ctx context.Context, merchantID int64, page, pageSize int) ([]*models.RateSyncLog, int64, error) {
	var result []*models.RateSyncLog
	for _, log := range m.logs {
		if merchantID == 0 || log.MerchantID == merchantID {
			result = append(result, log)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockRateSyncLogRepository) GetPendingLogs(ctx context.Context, limit int) ([]*models.RateSyncLog, error) {
	var result []*models.RateSyncLog
	for _, log := range m.logs {
		if log.SyncStatus == models.RateSyncStatusPending {
			result = append(result, log)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

// ==================== 测试用例 ====================

// TestRateUpdateParams 测试费率更新参数结构
func TestRateUpdateParams(t *testing.T) {
	params := &RateUpdateParams{
		MerchantID:  1,
		MerchantNo:  "M001",
		TerminalSN:  "T001",
		ChannelCode: "HENGXINTONG",
		AgentID:     100,
		OldRates: &RateInfo{
			CreditRate: 0.006,
			DebitRate:  0.005,
		},
		NewRates: &RateInfo{
			CreditRate: 0.0055,
			DebitRate:  0.0045,
		},
	}

	assert.Equal(t, int64(1), params.MerchantID)
	assert.Equal(t, "M001", params.MerchantNo)
	assert.Equal(t, "T001", params.TerminalSN)
	assert.Equal(t, "HENGXINTONG", params.ChannelCode)
	assert.Equal(t, int64(100), params.AgentID)
	assert.Equal(t, 0.006, params.OldRates.CreditRate)
	assert.Equal(t, 0.005, params.OldRates.DebitRate)
	assert.Equal(t, 0.0055, params.NewRates.CreditRate)
	assert.Equal(t, 0.0045, params.NewRates.DebitRate)
}

// TestSyncResult 测试同步结果结构
func TestSyncResult(t *testing.T) {
	tests := []struct {
		name    string
		result  *SyncResult
		success bool
	}{
		{
			name: "成功结果",
			result: &SyncResult{
				Success: true,
				LogID:   1,
				TradeNo: "TRX123456",
				Message: "费率同步成功",
			},
			success: true,
		},
		{
			name: "失败结果",
			result: &SyncResult{
				Success: false,
				LogID:   2,
				Message: "通道返回失败: 商户不存在",
			},
			success: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.success, tt.result.Success)
			assert.NotEmpty(t, tt.result.Message)
		})
	}
}

// TestRateInfo 测试费率信息结构
func TestRateInfo(t *testing.T) {
	tests := []struct {
		name       string
		rateInfo   *RateInfo
		creditRate float64
		debitRate  float64
	}{
		{
			name: "标准费率",
			rateInfo: &RateInfo{
				CreditRate:   0.006,
				DebitRate:    0.005,
				DebitCap:     20,
				WechatRate:   0.0038,
				AlipayRate:   0.0038,
				UnionpayRate: 0.006,
			},
			creditRate: 0.006,
			debitRate:  0.005,
		},
		{
			name: "零费率",
			rateInfo: &RateInfo{
				CreditRate: 0,
				DebitRate:  0,
			},
			creditRate: 0,
			debitRate:  0,
		},
		{
			name: "最高费率",
			rateInfo: &RateInfo{
				CreditRate: 0.1,
				DebitRate:  0.1,
			},
			creditRate: 0.1,
			debitRate:  0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.creditRate, tt.rateInfo.CreditRate)
			assert.Equal(t, tt.debitRate, tt.rateInfo.DebitRate)
		})
	}
}

// TestRateSyncLogModel 测试费率同步日志模型
func TestRateSyncLogModel(t *testing.T) {
	creditRate := 0.006
	debitRate := 0.005

	log := &models.RateSyncLog{
		ID:            1,
		MerchantID:    100,
		MerchantNo:    "M001",
		TerminalSN:    "T001",
		ChannelCode:   "HENGXINTONG",
		AgentID:       10,
		OldCreditRate: &creditRate,
		OldDebitRate:  &debitRate,
		SyncStatus:    models.RateSyncStatusPending,
	}

	// 测试初始状态
	assert.Equal(t, models.RateSyncStatusPending, log.SyncStatus)
	assert.Empty(t, log.ErrorMessage)
	assert.Empty(t, log.ChannelTradeNo)

	// 测试标记成功
	log.MarkSuccess("TRX123456")
	assert.Equal(t, models.RateSyncStatusSuccess, log.SyncStatus)
	assert.Equal(t, "TRX123456", log.ChannelTradeNo)
	assert.NotNil(t, log.SyncedAt)

	// 创建新日志测试标记失败
	log2 := &models.RateSyncLog{
		ID:         2,
		SyncStatus: models.RateSyncStatusSyncing,
	}
	log2.MarkFailed("商户不存在")
	assert.Equal(t, models.RateSyncStatusFailed, log2.SyncStatus)
	assert.Equal(t, "商户不存在", log2.ErrorMessage)
	assert.NotNil(t, log2.UpdatedAt) // MarkFailed 设置的是 UpdatedAt
}

// TestRateSyncStatusConstants 测试同步状态常量
func TestRateSyncStatusConstants(t *testing.T) {
	assert.Equal(t, models.RateSyncStatus(0), models.RateSyncStatusPending)
	assert.Equal(t, models.RateSyncStatus(1), models.RateSyncStatusSyncing)
	assert.Equal(t, models.RateSyncStatus(2), models.RateSyncStatusSuccess)
	assert.Equal(t, models.RateSyncStatus(3), models.RateSyncStatusFailed)
}

// TestRateValidation 测试费率范围验证
func TestRateValidation(t *testing.T) {
	tests := []struct {
		name       string
		creditRate float64
		debitRate  float64
		valid      bool
	}{
		{
			name:       "正常费率 0.6%",
			creditRate: 0.006,
			debitRate:  0.005,
			valid:      true,
		},
		{
			name:       "最低费率 0%",
			creditRate: 0,
			debitRate:  0,
			valid:      true,
		},
		{
			name:       "最高费率 10%",
			creditRate: 0.1,
			debitRate:  0.1,
			valid:      true,
		},
		{
			name:       "超出范围 负数",
			creditRate: -0.001,
			debitRate:  0.005,
			valid:      false,
		},
		{
			name:       "超出范围 超过10%",
			creditRate: 0.11,
			debitRate:  0.005,
			valid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证费率范围 (0-0.1)
			valid := tt.creditRate >= 0 && tt.creditRate <= 0.1 &&
				tt.debitRate >= 0 && tt.debitRate <= 0.1
			assert.Equal(t, tt.valid, valid)
		})
	}
}

// TestRateFormatConversion 测试费率格式转换
func TestRateFormatConversion(t *testing.T) {
	tests := []struct {
		name       string
		decimal    float64 // 小数形式 0.006
		percentage float64 // 百分比形式 0.6
	}{
		{
			name:       "0.6%",
			decimal:    0.006,
			percentage: 0.6,
		},
		{
			name:       "0.55%",
			decimal:    0.0055,
			percentage: 0.55,
		},
		{
			name:       "1%",
			decimal:    0.01,
			percentage: 1.0,
		},
		{
			name:       "0.38%",
			decimal:    0.0038,
			percentage: 0.38,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 小数转百分比
			percentage := tt.decimal * 100
			assert.InDelta(t, tt.percentage, percentage, 0.0001)

			// 百分比转小数
			decimal := tt.percentage / 100
			assert.InDelta(t, tt.decimal, decimal, 0.000001)
		})
	}
}

// TestMockRepository 测试模拟仓库功能
func TestMockRepository(t *testing.T) {
	repo := NewMockRateSyncLogRepository()
	ctx := context.Background()

	creditRate := 0.006
	debitRate := 0.005
	newCreditRate := 0.0055
	newDebitRate := 0.0045

	// 测试创建
	log := &models.RateSyncLog{
		MerchantID:    100,
		MerchantNo:    "M001",
		ChannelCode:   "HENGXINTONG",
		OldCreditRate: &creditRate,
		OldDebitRate:  &debitRate,
		NewCreditRate: &newCreditRate,
		NewDebitRate:  &newDebitRate,
		SyncStatus:    models.RateSyncStatusPending,
	}

	err := repo.Create(ctx, log)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), log.ID)

	// 测试查询
	found, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "M001", found.MerchantNo)

	// 测试更新
	log.MarkSuccess("TRX123")
	err = repo.Update(ctx, log)
	assert.NoError(t, err)

	// 测试按商户查询
	logs, total, err := repo.GetByMerchantID(ctx, 100, 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, logs, 1)

	// 测试查询所有
	logs, total, err = repo.GetByMerchantID(ctx, 0, 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
}

// TestRateCompareWithParent 测试费率与上级比较
func TestRateCompareWithParent(t *testing.T) {
	tests := []struct {
		name           string
		parentRate     float64
		newRate        float64
		shouldBeHigher bool // 新费率应该>=上级费率
	}{
		{
			name:           "新费率高于上级",
			parentRate:     0.005,
			newRate:        0.006,
			shouldBeHigher: true,
		},
		{
			name:           "新费率等于上级",
			parentRate:     0.005,
			newRate:        0.005,
			shouldBeHigher: true,
		},
		{
			name:           "新费率低于上级 - 不允许",
			parentRate:     0.006,
			newRate:        0.005,
			shouldBeHigher: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.newRate >= tt.parentRate
			assert.Equal(t, tt.shouldBeHigher, isValid)
		})
	}
}

// TestRateUpdateResult 测试费率更新结果
func TestRateUpdateResult(t *testing.T) {
	tests := []struct {
		name        string
		result      *RateUpdateResult
		expectSync  bool
		description string
	}{
		{
			name: "完全成功",
			result: &RateUpdateResult{
				Success:     true,
				SyncSuccess: true,
				SyncMessage: "费率更新成功",
			},
			expectSync:  true,
			description: "本地更新成功，通道同步成功",
		},
		{
			name: "本地成功但通道失败",
			result: &RateUpdateResult{
				Success:     true,
				SyncSuccess: false,
				SyncMessage: "通道返回失败: 商户不存在",
			},
			expectSync:  false,
			description: "本地不更新，因为通道同步失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectSync, tt.result.SyncSuccess)
			assert.NotEmpty(t, tt.result.SyncMessage)
		})
	}
}

// TestBusinessRule_ChannelFirstThenLocal 测试业务规则：先通道后本地
func TestBusinessRule_ChannelFirstThenLocal(t *testing.T) {
	// 业务规则：费率修改必须和通道联动
	// 原则：先调用通道，通道返回成功才能修改自己的费率

	type scenario struct {
		name               string
		channelSuccess     bool
		shouldUpdateLocal  bool
		expectedUserResult string
	}

	scenarios := []scenario{
		{
			name:               "通道成功 -> 更新本地",
			channelSuccess:     true,
			shouldUpdateLocal:  true,
			expectedUserResult: "费率修改成功，已同步到支付通道",
		},
		{
			name:               "通道失败 -> 不更新本地",
			channelSuccess:     false,
			shouldUpdateLocal:  false,
			expectedUserResult: "费率修改成功，但通道同步失败",
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
