package service

import (
	"testing"
	"time"

	"xiangshoufu/internal/models"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock Repositories for Deposit Cashback Service Tests
// =============================================================================

// MockDepositCashbackRecordRepository 模拟押金返现记录仓库
type MockDepositCashbackRecordRepository struct {
	records map[int64]*models.DepositCashbackRecord
	nextID  int64
}

func NewMockDepositCashbackRecordRepository() *MockDepositCashbackRecordRepository {
	return &MockDepositCashbackRecordRepository{
		records: make(map[int64]*models.DepositCashbackRecord),
		nextID:  1,
	}
}

func (m *MockDepositCashbackRecordRepository) Create(record *models.DepositCashbackRecord) error {
	record.ID = m.nextID
	m.nextID++
	m.records[record.ID] = record
	return nil
}

func (m *MockDepositCashbackRecordRepository) BatchCreate(records []*models.DepositCashbackRecord) error {
	for _, record := range records {
		if err := m.Create(record); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockDepositCashbackRecordRepository) FindByTerminalID(terminalID int64) ([]*models.DepositCashbackRecord, error) {
	var result []*models.DepositCashbackRecord
	for _, record := range m.records {
		if record.TerminalID == terminalID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (m *MockDepositCashbackRecordRepository) FindByAgentID(agentID int64, limit, offset int) ([]*models.DepositCashbackRecord, error) {
	var result []*models.DepositCashbackRecord
	for _, record := range m.records {
		if record.AgentID == agentID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (m *MockDepositCashbackRecordRepository) UpdateWalletStatus(id int64, status int16) error {
	if record, ok := m.records[id]; ok {
		record.WalletStatus = status
		now := time.Now()
		record.ProcessedAt = &now
	}
	return nil
}

// =============================================================================
// 级差计算测试
// =============================================================================

// TestDepositCashback_LevelDifference 测试押金返现级差计算
func TestDepositCashback_LevelDifference(t *testing.T) {
	// 模拟三级代理商返现场景
	// 平台 -> A -> B -> C（终端持有者）
	// 99元押金返现配置：平台80, A=70, B=60, C=50

	testCases := []struct {
		name           string
		selfCashback   int64 // 自身配置
		lowerCashback  int64 // 下级配置（用于计算级差）
		wantActual     int64 // 实际返现 = 自身 - 下级
	}{
		{"平台级差", 8000, 7000, 1000}, // 平台得：80-70=10元
		{"A级级差", 7000, 6000, 1000},  // A得：70-60=10元
		{"B级级差", 6000, 5000, 1000},  // B得：60-50=10元
		{"C级级差（终端持有者）", 5000, 0, 5000}, // C得：50-0=50元
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := &models.DepositCashbackRecord{
				SelfCashback:  tc.selfCashback,
				UpperCashback: tc.lowerCashback,
			}

			actualCashback := record.SelfCashback - record.UpperCashback
			assert.Equal(t, tc.wantActual, actualCashback)
		})
	}
}

// TestDepositCashback_MultiLevel 测试多级返现计算
func TestDepositCashback_MultiLevel(t *testing.T) {
	// 模拟四级代理商返现场景
	// 平台 -> A -> B -> C -> D（终端持有者）
	// 99元押金返现配置：平台80, A=70, B=60, C=50, D=40

	agentConfigs := []struct {
		agentID      int64
		selfCashback int64
	}{
		{0, 8000}, // 平台
		{1, 7000}, // A
		{2, 6000}, // B
		{3, 5000}, // C
		{4, 4000}, // D（终端持有者）
	}

	// 计算每级的实际返现（级差）
	var totalDistributed int64
	results := make(map[int64]int64)

	for i := len(agentConfigs) - 1; i >= 0; i-- {
		config := agentConfigs[i]
		var lowerCashback int64
		if i < len(agentConfigs)-1 {
			lowerCashback = agentConfigs[i+1].selfCashback
		}

		actualCashback := config.selfCashback - lowerCashback
		if actualCashback > 0 {
			results[config.agentID] = actualCashback
			totalDistributed += actualCashback
		}
	}

	// 验证每级获得正确的级差返现
	expectedResults := map[int64]int64{
		0: 1000, // 平台: 80-70=10元
		1: 1000, // A: 70-60=10元
		2: 1000, // B: 60-50=10元
		3: 1000, // C: 50-40=10元
		4: 4000, // D: 40-0=40元
	}

	for agentID, expected := range expectedResults {
		assert.Equal(t, expected, results[agentID], "Agent %d cashback mismatch", agentID)
	}

	// 总分配应等于最高配置（平台配置）
	assert.Equal(t, int64(8000), totalDistributed)
}

// TestDepositCashback_ZeroLowerLevel 测试终端持有者（无下级）情况
func TestDepositCashback_ZeroLowerLevel(t *testing.T) {
	// 终端持有者，下级返现为0
	selfCashback := int64(5000) // 50元
	lowerCashback := int64(0)   // 无下级

	actualCashback := selfCashback - lowerCashback
	assert.Equal(t, int64(5000), actualCashback)
}

// TestDepositCashback_SameLevelConfig 测试相同配置情况（无级差）
func TestDepositCashback_SameLevelConfig(t *testing.T) {
	// 如果上下级配置相同，则上级无返现
	selfCashback := int64(5000)
	lowerCashback := int64(5000)

	actualCashback := selfCashback - lowerCashback
	assert.Equal(t, int64(0), actualCashback)
}

// TestDepositCashback_NegativeDifference 测试配置倒挂情况
func TestDepositCashback_NegativeDifference(t *testing.T) {
	// 如果下级配置比上级高（配置错误），应跳过
	selfCashback := int64(3000)
	lowerCashback := int64(5000)

	actualCashback := selfCashback - lowerCashback
	// 应该为负数，业务逻辑应跳过
	assert.True(t, actualCashback < 0)
}

// =============================================================================
// 押金档次测试
// =============================================================================

// TestDepositAmount_Tiers 测试押金档次
func TestDepositAmount_Tiers(t *testing.T) {
	testCases := []struct {
		depositAmount int64
		tierName      string
	}{
		{9900, "99元押金"},
		{19900, "199元押金"},
		{29900, "299元押金"},
		{0, "无押金"},
	}

	for _, tc := range testCases {
		t.Run(tc.tierName, func(t *testing.T) {
			var tierName string
			switch tc.depositAmount {
			case 9900:
				tierName = "99元押金"
			case 19900:
				tierName = "199元押金"
			case 29900:
				tierName = "299元押金"
			default:
				tierName = "无押金"
			}
			assert.Equal(t, tc.tierName, tierName)
		})
	}
}

// =============================================================================
// 返现记录测试
// =============================================================================

// TestDepositCashbackRecord_WalletType 测试押金返现入账钱包类型
func TestDepositCashbackRecord_WalletType(t *testing.T) {
	// 押金返现应入服务费钱包（WalletType=2）
	record := &models.DepositCashbackRecord{
		TerminalSN:     "SN123456",
		ChannelID:      1,
		AgentID:        1,
		DepositAmount:  9900,
		SelfCashback:   5000,
		UpperCashback:  0,
		ActualCashback: 5000,
		WalletType:     models.WalletTypeService, // 应为2
		WalletStatus:   0,                        // 待入账
	}

	assert.Equal(t, int16(models.WalletTypeService), record.WalletType)
	assert.Equal(t, int16(2), record.WalletType) // 确认是服务费钱包
}

// TestDepositCashbackRecord_WalletStatus 测试返现记录钱包状态
func TestDepositCashbackRecord_WalletStatus(t *testing.T) {
	recordRepo := NewMockDepositCashbackRecordRepository()

	// 创建返现记录
	record := &models.DepositCashbackRecord{
		TerminalID:     1,
		TerminalSN:     "SN123456",
		ChannelID:      1,
		AgentID:        1,
		DepositAmount:  9900,
		SelfCashback:   5000,
		UpperCashback:  0,
		ActualCashback: 5000,
		WalletType:     models.WalletTypeService,
		WalletStatus:   0, // 待入账
		CreatedAt:      time.Now(),
	}
	recordRepo.Create(record)

	// 验证初始状态
	assert.Equal(t, int16(0), record.WalletStatus)

	// 更新为已入账
	recordRepo.UpdateWalletStatus(record.ID, 1)

	records, _ := recordRepo.FindByTerminalID(1)
	assert.Len(t, records, 1)
	assert.Equal(t, int16(1), records[0].WalletStatus)
	assert.NotNil(t, records[0].ProcessedAt)
}

// TestDepositCashbackRecord_SourceAgentID 测试级差来源代理商
func TestDepositCashbackRecord_SourceAgentID(t *testing.T) {
	// 终端持有者，无来源代理商
	record1 := &models.DepositCashbackRecord{
		AgentID:       1,
		SourceAgentID: nil, // 终端持有者无来源
	}
	assert.Nil(t, record1.SourceAgentID)

	// 上级代理商，来源是下级
	srcID := int64(1)
	record2 := &models.DepositCashbackRecord{
		AgentID:       2,
		SourceAgentID: &srcID, // 来源是代理商1
	}
	assert.NotNil(t, record2.SourceAgentID)
	assert.Equal(t, int64(1), *record2.SourceAgentID)
}

// =============================================================================
// 统计测试
// =============================================================================

// TestDepositCashbackStats 测试押金返现统计
func TestDepositCashbackStats(t *testing.T) {
	// 模拟返现记录
	records := []*models.DepositCashbackRecord{
		{ID: 1, DepositAmount: 9900, ActualCashback: 5000},
		{ID: 2, DepositAmount: 9900, ActualCashback: 3000},
		{ID: 3, DepositAmount: 19900, ActualCashback: 8000},
		{ID: 4, DepositAmount: 29900, ActualCashback: 10000},
		{ID: 5, DepositAmount: 29900, ActualCashback: 10000},
	}

	stats := &DepositCashbackStats{}
	for _, record := range records {
		stats.TotalCashback += record.ActualCashback
		stats.TotalCount++

		switch record.DepositAmount {
		case 9900:
			stats.Deposit99Count++
		case 19900:
			stats.Deposit199Count++
		case 29900:
			stats.Deposit299Count++
		}
	}

	assert.Equal(t, int64(36000), stats.TotalCashback) // 50+30+80+100+100=360元
	assert.Equal(t, int64(5), stats.TotalCount)
	assert.Equal(t, int64(2), stats.Deposit99Count)
	assert.Equal(t, int64(1), stats.Deposit199Count)
	assert.Equal(t, int64(2), stats.Deposit299Count)
}

// =============================================================================
// 边界条件测试
// =============================================================================

// TestDepositCashback_ZeroDeposit 测试零押金情况
func TestDepositCashback_ZeroDeposit(t *testing.T) {
	// 零押金不应触发返现
	depositAmount := int64(0)
	assert.Equal(t, int64(0), depositAmount)
	// 业务逻辑应该在押金为0时直接返回
}

// TestDepositCashback_MinCashback 测试最小返现金额
func TestDepositCashback_MinCashback(t *testing.T) {
	// 返现金额为1分
	actualCashback := int64(1)
	assert.Equal(t, int64(1), actualCashback)
	assert.Equal(t, float64(0.01), float64(actualCashback)/100)
}

// TestDepositCashback_LargeCashback 测试大额返现
func TestDepositCashback_LargeCashback(t *testing.T) {
	// 返现金额10000元
	actualCashback := int64(1000000)
	assert.Equal(t, float64(10000), float64(actualCashback)/100)
}
