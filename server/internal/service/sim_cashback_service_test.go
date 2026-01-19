package service

import (
	"testing"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// MockSimCashbackPolicyRepository 模拟流量费返现政策仓库
type MockSimCashbackPolicyRepository struct {
	policies map[int64]*models.SimCashbackPolicy
	nextID   int64
}

func NewMockSimCashbackPolicyRepository() *MockSimCashbackPolicyRepository {
	return &MockSimCashbackPolicyRepository{
		policies: make(map[int64]*models.SimCashbackPolicy),
		nextID:   1,
	}
}

func (m *MockSimCashbackPolicyRepository) AddPolicy(templateID, channelID int64, brandCode string, first, second, thirdPlus, simFeeAmount int64) *models.SimCashbackPolicy {
	policy := &models.SimCashbackPolicy{
		ID:                 m.nextID,
		TemplateID:         templateID,
		ChannelID:          channelID,
		BrandCode:          brandCode,
		FirstTimeCashback:  first,
		SecondTimeCashback: second,
		ThirdPlusCashback:  thirdPlus,
		SimFeeAmount:       simFeeAmount,
		Status:             1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	m.policies[policy.ID] = policy
	m.nextID++
	return policy
}

func (m *MockSimCashbackPolicyRepository) Create(policy *models.SimCashbackPolicy) error {
	policy.ID = m.nextID
	m.nextID++
	m.policies[policy.ID] = policy
	return nil
}

func (m *MockSimCashbackPolicyRepository) Update(policy *models.SimCashbackPolicy) error {
	m.policies[policy.ID] = policy
	return nil
}

func (m *MockSimCashbackPolicyRepository) FindByID(id int64) (*models.SimCashbackPolicy, error) {
	policy, ok := m.policies[id]
	if !ok {
		return nil, nil
	}
	return policy, nil
}

func (m *MockSimCashbackPolicyRepository) FindByTemplateAndChannel(templateID, channelID int64, brandCode string) (*models.SimCashbackPolicy, error) {
	for _, policy := range m.policies {
		if policy.TemplateID == templateID && policy.ChannelID == channelID && policy.Status == 1 {
			if brandCode == "" || policy.BrandCode == brandCode {
				return policy, nil
			}
		}
	}
	return nil, nil
}

func (m *MockSimCashbackPolicyRepository) FindByChannel(channelID int64) ([]*models.SimCashbackPolicy, error) {
	var result []*models.SimCashbackPolicy
	for _, policy := range m.policies {
		if policy.ChannelID == channelID && policy.Status == 1 {
			result = append(result, policy)
		}
	}
	return result, nil
}

// MockSimCashbackRecordRepository 模拟流量费返现记录仓库
type MockSimCashbackRecordRepository struct {
	records map[int64]*models.SimCashbackRecord
	nextID  int64
}

func NewMockSimCashbackRecordRepository() *MockSimCashbackRecordRepository {
	return &MockSimCashbackRecordRepository{
		records: make(map[int64]*models.SimCashbackRecord),
		nextID:  1,
	}
}

func (m *MockSimCashbackRecordRepository) Create(record *models.SimCashbackRecord) error {
	record.ID = m.nextID
	m.nextID++
	m.records[record.ID] = record
	return nil
}

func (m *MockSimCashbackRecordRepository) BatchCreate(records []*models.SimCashbackRecord) error {
	for _, record := range records {
		if err := m.Create(record); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockSimCashbackRecordRepository) FindByDeviceFeeID(deviceFeeID int64) ([]*models.SimCashbackRecord, error) {
	var result []*models.SimCashbackRecord
	for _, record := range m.records {
		if record.DeviceFeeID == deviceFeeID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (m *MockSimCashbackRecordRepository) FindByAgent(agentID int64, limit, offset int) ([]*models.SimCashbackRecord, int64, error) {
	var result []*models.SimCashbackRecord
	for _, record := range m.records {
		if record.AgentID == agentID {
			result = append(result, record)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockSimCashbackRecordRepository) UpdateWalletStatus(id int64, status int16) error {
	if record, ok := m.records[id]; ok {
		record.WalletStatus = status
		now := time.Now()
		record.ProcessedAt = &now
	}
	return nil
}

// MockDeviceFeeRepository 模拟流量费仓库
type MockDeviceFeeRepository struct {
	fees   map[int64]*models.DeviceFee
	nextID int64
}

func NewMockDeviceFeeRepository() *MockDeviceFeeRepository {
	return &MockDeviceFeeRepository{
		fees:   make(map[int64]*models.DeviceFee),
		nextID: 1,
	}
}

func (m *MockDeviceFeeRepository) AddDeviceFee(terminalSN string, channelID int64, feeAmount int64) *models.DeviceFee {
	fee := &models.DeviceFee{
		ID:         m.nextID,
		TerminalSN: terminalSN,
		ChannelID:  channelID,
		FeeAmount:  feeAmount,
		FeeType:    2, // 流量费
		CreatedAt:  time.Now(),
	}
	m.fees[fee.ID] = fee
	m.nextID++
	return fee
}

func (m *MockDeviceFeeRepository) Create(fee *models.DeviceFee) error {
	fee.ID = m.nextID
	m.nextID++
	m.fees[fee.ID] = fee
	return nil
}

func (m *MockDeviceFeeRepository) Update(fee *models.DeviceFee) error {
	m.fees[fee.ID] = fee
	return nil
}

func (m *MockDeviceFeeRepository) FindByOrderNo(orderNo string) (*models.DeviceFee, error) {
	for _, fee := range m.fees {
		if fee.OrderNo == orderNo {
			return fee, nil
		}
	}
	return nil, nil
}

func (m *MockDeviceFeeRepository) FindPendingCashback(limit int) ([]*models.DeviceFee, error) {
	var result []*models.DeviceFee
	for _, fee := range m.fees {
		if fee.CashbackStatus == 0 {
			result = append(result, fee)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockDeviceFeeRepository) UpdateCashbackStatus(id int64, status int16, amount int64) error {
	if fee, ok := m.fees[id]; ok {
		fee.CashbackStatus = status
		fee.CashbackAmount = amount
	}
	return nil
}

// MockAgentPolicyRepository 模拟代理商政策仓库
type MockAgentPolicyRepository struct {
	policies map[string]*repository.AgentPolicy
}

func NewMockAgentPolicyRepository() *MockAgentPolicyRepository {
	return &MockAgentPolicyRepository{
		policies: make(map[string]*repository.AgentPolicy),
	}
}

func (m *MockAgentPolicyRepository) AddPolicy(agentID, channelID, templateID int64) {
	key := m.makeKey(agentID, channelID)
	m.policies[key] = &repository.AgentPolicy{
		AgentID:    agentID,
		ChannelID:  channelID,
		TemplateID: templateID,
	}
}

func (m *MockAgentPolicyRepository) makeKey(agentID, channelID int64) string {
	return string(rune(agentID)) + "-" + string(rune(channelID))
}

func (m *MockAgentPolicyRepository) FindByAgentAndChannel(agentID int64, channelID int64) (*repository.AgentPolicy, error) {
	for _, policy := range m.policies {
		if policy.AgentID == agentID && policy.ChannelID == channelID {
			return policy, nil
		}
	}
	return nil, nil
}

// =============================================================================
// 测试用例
// =============================================================================

// TestGetCashbackTier 测试返现档次计算（Q30：三档）
func TestGetCashbackTier(t *testing.T) {
	testCases := []struct {
		simFeeCount int
		wantTier    int16
		wantName    string
	}{
		{1, models.SimCashbackTierFirst, "首次"},
		{2, models.SimCashbackTierSecond, "第2次"},
		{3, models.SimCashbackTierThirdPlus, "第3次及以后"},
		{4, models.SimCashbackTierThirdPlus, "第3次及以后"},
		{10, models.SimCashbackTierThirdPlus, "第3次及以后"},
		{100, models.SimCashbackTierThirdPlus, "第3次及以后"},
	}

	for _, tc := range testCases {
		t.Run(tc.wantName, func(t *testing.T) {
			tier := models.GetCashbackTier(tc.simFeeCount)
			if tier != tc.wantTier {
				t.Errorf("GetCashbackTier(%d) = %d, want %d", tc.simFeeCount, tier, tc.wantTier)
			}
		})
	}
}

// TestGetTierName 测试档次名称
func TestGetTierName(t *testing.T) {
	testCases := []struct {
		tier     int16
		wantName string
	}{
		{models.SimCashbackTierFirst, "首次"},
		{models.SimCashbackTierSecond, "第2次"},
		{models.SimCashbackTierThirdPlus, "第3次及以后"},
		{99, ""}, // 未知档次
	}

	for _, tc := range testCases {
		t.Run(tc.wantName, func(t *testing.T) {
			name := getTierName(tc.tier)
			if name != tc.wantName {
				t.Errorf("getTierName(%d) = %s, want %s", tc.tier, name, tc.wantName)
			}
		})
	}
}

// TestSimCashbackPolicy_ThreeTiers 测试三档返现配置
func TestSimCashbackPolicy_ThreeTiers(t *testing.T) {
	// 创建三档返现政策
	policy := &models.SimCashbackPolicy{
		ID:                 1,
		TemplateID:         1,
		ChannelID:          1,
		BrandCode:          "HENGXINTONG",
		FirstTimeCashback:  5000,  // 首次50元
		SecondTimeCashback: 3000,  // 第2次30元
		ThirdPlusCashback:  2000,  // 第3次及以后20元
		SimFeeAmount:       9900,  // 流量费99元
		Status:             1,
	}

	// 验证各档次返现金额
	tiers := []struct {
		tier       int16
		wantAmount int64
	}{
		{models.SimCashbackTierFirst, 5000},
		{models.SimCashbackTierSecond, 3000},
		{models.SimCashbackTierThirdPlus, 2000},
	}

	for _, tc := range tiers {
		var amount int64
		switch tc.tier {
		case models.SimCashbackTierFirst:
			amount = policy.FirstTimeCashback
		case models.SimCashbackTierSecond:
			amount = policy.SecondTimeCashback
		case models.SimCashbackTierThirdPlus:
			amount = policy.ThirdPlusCashback
		}

		if amount != tc.wantAmount {
			t.Errorf("Tier %d amount = %d, want %d", tc.tier, amount, tc.wantAmount)
		}
	}
}

// TestSimCashbackRecord_LevelDifference 测试级差计算
func TestSimCashbackRecord_LevelDifference(t *testing.T) {
	// 模拟三级代理商返现配置
	// 平台: 首次60元
	// A级: 首次50元
	// B级: 首次40元
	// C级: 首次30元

	testCases := []struct {
		name           string
		selfCashback   int64 // 自身配置
		upperCashback  int64 // 下级配置（用于计算级差）
		wantActual     int64 // 实际返现 = 自身 - 下级
	}{
		{"平台级差", 6000, 5000, 1000}, // 平台得：60-50=10元
		{"A级级差", 5000, 4000, 1000},  // A得：50-40=10元
		{"B级级差", 4000, 3000, 1000},  // B得：40-30=10元
		{"C级级差（终端持有者）", 3000, 0, 3000}, // C得：30-0=30元
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := &models.SimCashbackRecord{
				SelfCashback:  tc.selfCashback,
				UpperCashback: tc.upperCashback,
			}

			actualCashback := record.SelfCashback - record.UpperCashback
			if actualCashback != tc.wantActual {
				t.Errorf("ActualCashback = %d, want %d", actualCashback, tc.wantActual)
			}
		})
	}
}

// TestCashbackStats 测试返现统计
func TestCashbackStats(t *testing.T) {
	// 模拟返现记录
	records := []*models.SimCashbackRecord{
		{ID: 1, CashbackTier: models.SimCashbackTierFirst, ActualCashback: 5000},
		{ID: 2, CashbackTier: models.SimCashbackTierFirst, ActualCashback: 5000},
		{ID: 3, CashbackTier: models.SimCashbackTierSecond, ActualCashback: 3000},
		{ID: 4, CashbackTier: models.SimCashbackTierThirdPlus, ActualCashback: 2000},
		{ID: 5, CashbackTier: models.SimCashbackTierThirdPlus, ActualCashback: 2000},
		{ID: 6, CashbackTier: models.SimCashbackTierThirdPlus, ActualCashback: 2000},
	}

	// 计算统计
	stats := &CashbackStats{}
	for _, record := range records {
		stats.TotalCashback += record.ActualCashback
		stats.TotalCount++

		switch record.CashbackTier {
		case models.SimCashbackTierFirst:
			stats.FirstTierCashback += record.ActualCashback
		case models.SimCashbackTierSecond:
			stats.SecondTierCashback += record.ActualCashback
		case models.SimCashbackTierThirdPlus:
			stats.ThirdTierCashback += record.ActualCashback
		}
	}

	// 验证
	if stats.TotalCashback != 19000 {
		t.Errorf("TotalCashback = %d, want 19000", stats.TotalCashback)
	}

	if stats.TotalCount != 6 {
		t.Errorf("TotalCount = %d, want 6", stats.TotalCount)
	}

	if stats.FirstTierCashback != 10000 {
		t.Errorf("FirstTierCashback = %d, want 10000", stats.FirstTierCashback)
	}

	if stats.SecondTierCashback != 3000 {
		t.Errorf("SecondTierCashback = %d, want 3000", stats.SecondTierCashback)
	}

	if stats.ThirdTierCashback != 6000 {
		t.Errorf("ThirdTierCashback = %d, want 6000", stats.ThirdTierCashback)
	}
}

// TestSimCashbackRecord_WalletStatus 测试返现记录钱包状态
func TestSimCashbackRecord_WalletStatus(t *testing.T) {
	recordRepo := NewMockSimCashbackRecordRepository()

	// 创建返现记录
	record := &models.SimCashbackRecord{
		DeviceFeeID:    1,
		TerminalSN:     "SN123456",
		ChannelID:      1,
		AgentID:        1,
		SimFeeCount:    1,
		SimFeeAmount:   9900,
		CashbackTier:   models.SimCashbackTierFirst,
		SelfCashback:   5000,
		UpperCashback:  0,
		ActualCashback: 5000,
		WalletType:     4,
		WalletStatus:   0, // 待入账
		CreatedAt:      time.Now(),
	}
	recordRepo.Create(record)

	// 验证初始状态
	if record.WalletStatus != 0 {
		t.Errorf("Initial WalletStatus = %d, want 0", record.WalletStatus)
	}

	// 更新为已入账
	recordRepo.UpdateWalletStatus(record.ID, 1)

	updated, _, _ := recordRepo.FindByAgent(1, 10, 0)
	if len(updated) == 0 {
		t.Fatal("Should find the record")
	}

	if updated[0].WalletStatus != 1 {
		t.Errorf("Updated WalletStatus = %d, want 1", updated[0].WalletStatus)
	}

	if updated[0].ProcessedAt == nil {
		t.Error("ProcessedAt should be set")
	}
}

// TestSimCashbackPolicyRepository 测试政策仓库
func TestSimCashbackPolicyRepository(t *testing.T) {
	repo := NewMockSimCashbackPolicyRepository()

	// 添加多个模板的政策
	repo.AddPolicy(1, 1, "HENGXINTONG", 5000, 3000, 2000, 9900) // 模板1
	repo.AddPolicy(2, 1, "HENGXINTONG", 4000, 2500, 1500, 9900) // 模板2
	repo.AddPolicy(3, 1, "HENGXINTONG", 3000, 2000, 1000, 9900) // 模板3

	// 查找模板1的政策
	policy, err := repo.FindByTemplateAndChannel(1, 1, "HENGXINTONG")
	if err != nil {
		t.Fatalf("FindByTemplateAndChannel failed: %v", err)
	}

	if policy == nil {
		t.Fatal("Policy should not be nil")
	}

	if policy.FirstTimeCashback != 5000 {
		t.Errorf("FirstTimeCashback = %d, want 5000", policy.FirstTimeCashback)
	}

	// 查找模板2的政策
	policy2, _ := repo.FindByTemplateAndChannel(2, 1, "HENGXINTONG")
	if policy2.FirstTimeCashback != 4000 {
		t.Errorf("Template 2 FirstTimeCashback = %d, want 4000", policy2.FirstTimeCashback)
	}
}

// TestTerminalSimFeeCount 测试终端流量费次数更新
func TestTerminalSimFeeCount(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()

	// 创建终端
	terminal := terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusActivated)

	if terminal.SimFeeCount != 0 {
		t.Errorf("Initial SimFeeCount = %d, want 0", terminal.SimFeeCount)
	}

	// 更新流量费次数
	terminalRepo.UpdateSimFeeCount(terminal.ID, 1)

	updated, _ := terminalRepo.FindBySN("SN123456")
	if updated.SimFeeCount != 1 {
		t.Errorf("Updated SimFeeCount = %d, want 1", updated.SimFeeCount)
	}

	if updated.LastSimFeeAt == nil {
		t.Error("LastSimFeeAt should be set")
	}

	// 继续更新
	terminalRepo.UpdateSimFeeCount(terminal.ID, 2)
	updated, _ = terminalRepo.FindBySN("SN123456")
	if updated.SimFeeCount != 2 {
		t.Errorf("Second update SimFeeCount = %d, want 2", updated.SimFeeCount)
	}
}

// TestMultiLevelCashback 测试多级返现计算
func TestMultiLevelCashback(t *testing.T) {
	// 模拟三级代理商返现场景
	// 平台 -> A -> B -> C（终端持有者）
	// 首次返现配置：平台60, A=50, B=40, C=30

	agentConfigs := []struct {
		agentID      int64
		selfCashback int64
	}{
		{0, 6000}, // 平台
		{1, 5000}, // A
		{2, 4000}, // B
		{3, 3000}, // C（终端持有者）
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

	// 验证
	expectedResults := map[int64]int64{
		0: 1000, // 平台: 60-50=10元
		1: 1000, // A: 50-40=10元
		2: 1000, // B: 40-30=10元
		3: 3000, // C: 30-0=30元
	}

	for agentID, expected := range expectedResults {
		if results[agentID] != expected {
			t.Errorf("Agent %d cashback = %d, want %d", agentID, results[agentID], expected)
		}
	}

	// 总分配应等于最高配置
	if totalDistributed != 6000 {
		t.Errorf("TotalDistributed = %d, want 6000", totalDistributed)
	}
}
