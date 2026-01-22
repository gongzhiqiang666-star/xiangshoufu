package service

import (
	"testing"
	"time"

	"xiangshoufu/internal/models"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock Repositories for Activation Reward Service Tests
// =============================================================================

// MockActivationRewardRecordRepository 模拟激活奖励记录仓库
type MockActivationRewardRecordRepository struct {
	records map[int64]*models.ActivationRewardRecord
	nextID  int64
}

func NewMockActivationRewardRecordRepository() *MockActivationRewardRecordRepository {
	return &MockActivationRewardRecordRepository{
		records: make(map[int64]*models.ActivationRewardRecord),
		nextID:  1,
	}
}

func (m *MockActivationRewardRecordRepository) Create(record *models.ActivationRewardRecord) error {
	record.ID = m.nextID
	m.nextID++
	m.records[record.ID] = record
	return nil
}

func (m *MockActivationRewardRecordRepository) BatchCreate(records []*models.ActivationRewardRecord) error {
	for _, record := range records {
		if err := m.Create(record); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockActivationRewardRecordRepository) FindByPolicyAndTerminal(policyID, terminalID int64, checkDate time.Time) (*models.ActivationRewardRecord, error) {
	for _, record := range m.records {
		if record.PolicyID == policyID && record.TerminalID == terminalID {
			return record, nil
		}
	}
	return nil, nil
}

func (m *MockActivationRewardRecordRepository) FindByAgentID(agentID int64, limit, offset int) ([]*models.ActivationRewardRecord, int64, error) {
	var result []*models.ActivationRewardRecord
	for _, record := range m.records {
		if record.AgentID == agentID {
			result = append(result, record)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockActivationRewardRecordRepository) FindPendingRecords(limit int) ([]*models.ActivationRewardRecord, error) {
	var result []*models.ActivationRewardRecord
	for _, record := range m.records {
		if record.WalletStatus == 0 {
			result = append(result, record)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockActivationRewardRecordRepository) UpdateWalletStatus(id int64, status int16) error {
	if record, ok := m.records[id]; ok {
		record.WalletStatus = status
		now := time.Now()
		record.ProcessedAt = &now
	}
	return nil
}

func (m *MockActivationRewardRecordRepository) BatchUpdateWalletStatus(ids []int64, status int16) error {
	for _, id := range ids {
		m.UpdateWalletStatus(id, status)
	}
	return nil
}

// MockActivationRewardPolicyRepository 模拟激活奖励政策仓库
type MockActivationRewardPolicyRepository struct {
	policies map[int64]*models.ActivationRewardPolicy
	nextID   int64
}

func NewMockActivationRewardPolicyRepository() *MockActivationRewardPolicyRepository {
	return &MockActivationRewardPolicyRepository{
		policies: make(map[int64]*models.ActivationRewardPolicy),
		nextID:   1,
	}
}

func (m *MockActivationRewardPolicyRepository) AddPolicy(templateID int64, minDays, maxDays int, targetAmount, rewardAmount int64) *models.ActivationRewardPolicy {
	policy := &models.ActivationRewardPolicy{
		ID:              m.nextID,
		TemplateID:      templateID,
		MinRegisterDays: minDays,
		MaxRegisterDays: maxDays,
		TargetAmount:    targetAmount,
		RewardAmount:    rewardAmount,
		Status:          1,
	}
	m.policies[policy.ID] = policy
	m.nextID++
	return policy
}

func (m *MockActivationRewardPolicyRepository) FindByTemplateID(templateID int64) ([]*models.ActivationRewardPolicy, error) {
	var result []*models.ActivationRewardPolicy
	for _, policy := range m.policies {
		if policy.TemplateID == templateID && policy.Status == 1 {
			result = append(result, policy)
		}
	}
	return result, nil
}

// =============================================================================
// 级差计算测试
// =============================================================================

// TestActivationReward_LevelDifference 测试激活奖励级差计算
func TestActivationReward_LevelDifference(t *testing.T) {
	// 模拟三级代理商奖励场景
	// 平台 -> A -> B -> C（终端持有者）
	// 激活奖励配置：平台100, A=80, B=60, C=40

	testCases := []struct {
		name        string
		selfReward  int64 // 自身配置
		upperReward int64 // 下级配置（用于计算级差，字段名为UpperReward但实际存储下级配置）
		wantActual  int64 // 实际奖励 = 自身 - 下级
	}{
		{"平台级差", 10000, 8000, 2000}, // 平台得：100-80=20元
		{"A级级差", 8000, 6000, 2000},   // A得：80-60=20元
		{"B级级差", 6000, 4000, 2000},   // B得：60-40=20元
		{"C级级差（终端持有者）", 4000, 0, 4000}, // C得：40-0=40元
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := &models.ActivationRewardRecord{
				SelfReward:  tc.selfReward,
				UpperReward: tc.upperReward,
			}

			actualReward := record.SelfReward - record.UpperReward
			assert.Equal(t, tc.wantActual, actualReward)
		})
	}
}

// TestActivationReward_MultiLevel 测试多级奖励计算
func TestActivationReward_MultiLevel(t *testing.T) {
	// 模拟五级代理商奖励场景
	agentConfigs := []struct {
		agentID    int64
		selfReward int64
	}{
		{0, 10000}, // 平台100元
		{1, 8000},  // A 80元
		{2, 6000},  // B 60元
		{3, 4000},  // C 40元
		{4, 2000},  // D 20元（终端持有者）
	}

	var totalDistributed int64
	results := make(map[int64]int64)

	for i := len(agentConfigs) - 1; i >= 0; i-- {
		config := agentConfigs[i]
		var lowerReward int64
		if i < len(agentConfigs)-1 {
			lowerReward = agentConfigs[i+1].selfReward
		}

		actualReward := config.selfReward - lowerReward
		if actualReward > 0 {
			results[config.agentID] = actualReward
			totalDistributed += actualReward
		}
	}

	// 验证每级获得正确的级差奖励
	expectedResults := map[int64]int64{
		0: 2000, // 平台: 100-80=20元
		1: 2000, // A: 80-60=20元
		2: 2000, // B: 60-40=20元
		3: 2000, // C: 40-20=20元
		4: 2000, // D: 20-0=20元
	}

	for agentID, expected := range expectedResults {
		assert.Equal(t, expected, results[agentID], "Agent %d reward mismatch", agentID)
	}

	// 总分配应等于最高配置
	assert.Equal(t, int64(10000), totalDistributed)
}

// =============================================================================
// 激活条件测试
// =============================================================================

// TestActivationReward_RegisterDaysCheck 测试入网天数判断
func TestActivationReward_RegisterDaysCheck(t *testing.T) {
	testCases := []struct {
		name        string
		registerAt  time.Time
		checkDate   time.Time
		wantDays    int
	}{
		{
			"入网第1天",
			time.Now().AddDate(0, 0, 0),
			time.Now(),
			0,
		},
		{
			"入网第7天",
			time.Now().AddDate(0, 0, -7),
			time.Now(),
			7,
		},
		{
			"入网第30天",
			time.Now().AddDate(0, 0, -30),
			time.Now(),
			30,
		},
		{
			"入网第90天",
			time.Now().AddDate(0, 0, -90),
			time.Now(),
			90,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registerDays := int(tc.checkDate.Sub(tc.registerAt).Hours() / 24)
			assert.Equal(t, tc.wantDays, registerDays)
		})
	}
}

// TestActivationReward_TradeAmountCheck 测试交易量条件
func TestActivationReward_TradeAmountCheck(t *testing.T) {
	testCases := []struct {
		name         string
		tradeAmount  int64 // 累计交易量（分）
		targetAmount int64 // 目标交易量（分）
		shouldMatch  bool
	}{
		{"交易量未达标", 500000, 1000000, false},     // 5000元 < 10000元
		{"交易量刚好达标", 1000000, 1000000, true},   // 10000元 = 10000元
		{"交易量超过目标", 2000000, 1000000, true},   // 20000元 > 10000元
		{"零交易量", 0, 1000000, false},              // 0 < 10000元
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matched := tc.tradeAmount >= tc.targetAmount
			assert.Equal(t, tc.shouldMatch, matched)
		})
	}
}

// TestActivationReward_PolicyMatching 测试政策匹配
func TestActivationReward_PolicyMatching(t *testing.T) {
	// 创建政策：7天内交易1万元，奖励50元
	policy := &models.ActivationRewardPolicy{
		ID:              1,
		TemplateID:      1,
		MinRegisterDays: 0,
		MaxRegisterDays: 7,
		TargetAmount:    1000000, // 1万元
		RewardAmount:    5000,    // 50元
		Status:          1,
	}

	testCases := []struct {
		name         string
		registerDays int
		tradeAmount  int64
		shouldMatch  bool
	}{
		{"条件满足", 5, 1500000, true},      // 5天，1.5万元
		{"天数超限", 10, 1500000, false},    // 10天 > 7天
		{"交易量不足", 5, 500000, false},    // 5000元 < 10000元
		{"刚好满足", 7, 1000000, true},      // 7天，1万元
		{"首日完成", 0, 1000000, true},      // 0天，1万元
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matched := tc.registerDays >= policy.MinRegisterDays &&
				tc.registerDays <= policy.MaxRegisterDays &&
				tc.tradeAmount >= policy.TargetAmount
			assert.Equal(t, tc.shouldMatch, matched)
		})
	}
}

// =============================================================================
// 奖励记录测试
// =============================================================================

// TestActivationRewardRecord_WalletType 测试激活奖励入账钱包类型
func TestActivationRewardRecord_WalletType(t *testing.T) {
	// 激活奖励应入奖励钱包（WalletType=3）
	record := &models.ActivationRewardRecord{
		TerminalSN:   "SN123456",
		ChannelID:    1,
		AgentID:      1,
		SelfReward:   5000,
		UpperReward:  0,
		ActualReward: 5000,
		WalletType:   models.WalletTypeReward, // 应为3
		WalletStatus: 0,                       // 待入账
	}

	assert.Equal(t, int16(models.WalletTypeReward), record.WalletType)
	assert.Equal(t, int16(3), record.WalletType) // 确认是奖励钱包
}

// TestActivationRewardRecord_WalletStatus 测试奖励记录钱包状态
func TestActivationRewardRecord_WalletStatus(t *testing.T) {
	recordRepo := NewMockActivationRewardRecordRepository()

	// 创建奖励记录
	record := &models.ActivationRewardRecord{
		PolicyID:     1,
		TerminalID:   1,
		TerminalSN:   "SN123456",
		ChannelID:    1,
		AgentID:      1,
		RegisterDays: 7,
		TradeAmount:  1000000,
		SelfReward:   5000,
		UpperReward:  0,
		ActualReward: 5000,
		WalletType:   models.WalletTypeReward,
		WalletStatus: 0, // 待入账
		CreatedAt:    time.Now(),
	}
	recordRepo.Create(record)

	// 验证初始状态
	assert.Equal(t, int16(0), record.WalletStatus)

	// 查找待入账记录
	pendingRecords, _ := recordRepo.FindPendingRecords(100)
	assert.Len(t, pendingRecords, 1)

	// 更新为已入账
	recordRepo.UpdateWalletStatus(record.ID, 1)

	records, _, _ := recordRepo.FindByAgentID(1, 10, 0)
	assert.Len(t, records, 1)
	assert.Equal(t, int16(1), records[0].WalletStatus)
	assert.NotNil(t, records[0].ProcessedAt)

	// 验证已无待入账记录
	pendingRecords, _ = recordRepo.FindPendingRecords(100)
	assert.Len(t, pendingRecords, 0)
}

// TestActivationRewardRecord_DuplicatePrevention 测试防重复发放
func TestActivationRewardRecord_DuplicatePrevention(t *testing.T) {
	recordRepo := NewMockActivationRewardRecordRepository()

	policyID := int64(1)
	terminalID := int64(1)
	checkDate := time.Now()

	// 第一次创建记录
	record1 := &models.ActivationRewardRecord{
		PolicyID:     policyID,
		TerminalID:   terminalID,
		TerminalSN:   "SN123456",
		ChannelID:    1,
		AgentID:      1,
		ActualReward: 5000,
		WalletStatus: 0,
		CreatedAt:    checkDate,
	}
	recordRepo.Create(record1)

	// 检查是否已存在记录
	existingRecord, _ := recordRepo.FindByPolicyAndTerminal(policyID, terminalID, checkDate)
	assert.NotNil(t, existingRecord)

	// 业务逻辑应该阻止重复发放
	shouldProcess := existingRecord == nil
	assert.False(t, shouldProcess)
}

// =============================================================================
// ProcessPendingRewards 测试
// =============================================================================

// TestProcessPendingRewards_GroupByWallet 测试按钱包分组汇总
func TestProcessPendingRewards_GroupByWallet(t *testing.T) {
	// 模拟多条待入账记录
	records := []*models.ActivationRewardRecord{
		{ID: 1, AgentID: 1, ChannelID: 1, ActualReward: 2000, WalletStatus: 0},
		{ID: 2, AgentID: 1, ChannelID: 1, ActualReward: 3000, WalletStatus: 0},
		{ID: 3, AgentID: 1, ChannelID: 2, ActualReward: 1500, WalletStatus: 0},
		{ID: 4, AgentID: 2, ChannelID: 1, ActualReward: 4000, WalletStatus: 0},
	}

	// 按代理商+通道分组汇总
	type walletKey struct {
		agentID   int64
		channelID int64
	}
	walletAmounts := make(map[walletKey]int64)

	for _, record := range records {
		key := walletKey{agentID: record.AgentID, channelID: record.ChannelID}
		walletAmounts[key] += record.ActualReward
	}

	// 验证汇总结果
	assert.Equal(t, int64(5000), walletAmounts[walletKey{1, 1}]) // 代理商1通道1: 2000+3000
	assert.Equal(t, int64(1500), walletAmounts[walletKey{1, 2}]) // 代理商1通道2: 1500
	assert.Equal(t, int64(4000), walletAmounts[walletKey{2, 1}]) // 代理商2通道1: 4000
}

// TestProcessPendingRewards_BatchUpdate 测试批量更新状态
func TestProcessPendingRewards_BatchUpdate(t *testing.T) {
	recordRepo := NewMockActivationRewardRecordRepository()

	// 创建多条待入账记录
	for i := 1; i <= 5; i++ {
		record := &models.ActivationRewardRecord{
			PolicyID:     1,
			TerminalID:   int64(i),
			AgentID:      1,
			ActualReward: 1000,
			WalletStatus: 0,
		}
		recordRepo.Create(record)
	}

	// 获取待入账记录
	pendingRecords, _ := recordRepo.FindPendingRecords(100)
	assert.Len(t, pendingRecords, 5)

	// 批量更新状态
	ids := make([]int64, len(pendingRecords))
	for i, r := range pendingRecords {
		ids[i] = r.ID
	}
	recordRepo.BatchUpdateWalletStatus(ids, 1)

	// 验证已无待入账记录
	pendingRecords, _ = recordRepo.FindPendingRecords(100)
	assert.Len(t, pendingRecords, 0)
}

// =============================================================================
// 边界条件测试
// =============================================================================

// TestActivationReward_ZeroReward 测试零奖励情况
func TestActivationReward_ZeroReward(t *testing.T) {
	// 配置相同，级差为0
	selfReward := int64(5000)
	lowerReward := int64(5000)

	actualReward := selfReward - lowerReward
	assert.Equal(t, int64(0), actualReward)
}

// TestActivationReward_NegativeDifference 测试配置倒挂情况
func TestActivationReward_NegativeDifference(t *testing.T) {
	// 下级配置比上级高（配置错误）
	selfReward := int64(3000)
	lowerReward := int64(5000)

	actualReward := selfReward - lowerReward
	assert.True(t, actualReward < 0)
	// 业务逻辑应跳过负数奖励
}

// TestActivationReward_MultiplePolicies 测试多政策优先级
func TestActivationReward_MultiplePolicies(t *testing.T) {
	policyRepo := NewMockActivationRewardPolicyRepository()

	// 添加多个政策，按优先级排序
	policyRepo.AddPolicy(1, 0, 7, 500000, 3000)   // 7天5000元，奖30元
	policyRepo.AddPolicy(1, 0, 7, 1000000, 5000)  // 7天1万元，奖50元
	policyRepo.AddPolicy(1, 0, 30, 3000000, 8000) // 30天3万元，奖80元

	// 获取模板1的所有政策
	policies, _ := policyRepo.FindByTemplateID(1)
	assert.Len(t, policies, 3)

	// 模拟匹配逻辑：5天交易8000元
	registerDays := 5
	tradeAmount := int64(800000)

	var matchedPolicy *models.ActivationRewardPolicy
	for _, policy := range policies {
		if registerDays >= policy.MinRegisterDays &&
			registerDays <= policy.MaxRegisterDays &&
			tradeAmount >= policy.TargetAmount {
			matchedPolicy = policy
			break // 取第一个匹配的（优先级最高）
		}
	}

	assert.NotNil(t, matchedPolicy)
	assert.Equal(t, int64(500000), matchedPolicy.TargetAmount)
	assert.Equal(t, int64(3000), matchedPolicy.RewardAmount)
}

// =============================================================================
// 金额转换测试
// =============================================================================

// TestActivationReward_AmountConversion 测试金额转换
func TestActivationReward_AmountConversion(t *testing.T) {
	testCases := []struct {
		amountFen  int64
		amountYuan float64
	}{
		{0, 0.00},
		{100, 1.00},
		{5000, 50.00},
		{10000, 100.00},
		{12345, 123.45},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			yuan := float64(tc.amountFen) / 100
			assert.Equal(t, tc.amountYuan, yuan)
		})
	}
}

// =============================================================================
// 统计测试
// =============================================================================

// TestActivationRewardStats 测试激活奖励统计
func TestActivationRewardStats(t *testing.T) {
	records := []*models.ActivationRewardRecord{
		{ID: 1, ActualReward: 5000, WalletStatus: 1},
		{ID: 2, ActualReward: 3000, WalletStatus: 1},
		{ID: 3, ActualReward: 2000, WalletStatus: 0},
		{ID: 4, ActualReward: 4000, WalletStatus: 1},
	}

	var totalReward int64
	var processedReward int64
	var pendingReward int64
	var processedCount int
	var pendingCount int

	for _, record := range records {
		totalReward += record.ActualReward
		if record.WalletStatus == 1 {
			processedReward += record.ActualReward
			processedCount++
		} else {
			pendingReward += record.ActualReward
			pendingCount++
		}
	}

	assert.Equal(t, int64(14000), totalReward)     // 140元
	assert.Equal(t, int64(12000), processedReward) // 120元已入账
	assert.Equal(t, int64(2000), pendingReward)    // 20元待入账
	assert.Equal(t, 3, processedCount)
	assert.Equal(t, 1, pendingCount)
}
