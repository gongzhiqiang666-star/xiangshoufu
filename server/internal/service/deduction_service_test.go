package service

import (
	"testing"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// MockDeductionPlanRepository 模拟代扣计划仓库
type MockDeductionPlanRepository struct {
	plans      map[int64]*models.DeductionPlan
	nextID     int64
	createErr  error
	findErr    error
	updateErr  error
}

func NewMockDeductionPlanRepository() *MockDeductionPlanRepository {
	return &MockDeductionPlanRepository{
		plans:  make(map[int64]*models.DeductionPlan),
		nextID: 1,
	}
}

func (m *MockDeductionPlanRepository) Create(plan *models.DeductionPlan) error {
	if m.createErr != nil {
		return m.createErr
	}
	plan.ID = m.nextID
	m.nextID++
	m.plans[plan.ID] = plan
	return nil
}

func (m *MockDeductionPlanRepository) Update(plan *models.DeductionPlan) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.plans[plan.ID] = plan
	return nil
}

func (m *MockDeductionPlanRepository) FindByID(id int64) (*models.DeductionPlan, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	plan, ok := m.plans[id]
	if !ok {
		return nil, nil
	}
	return plan, nil
}

func (m *MockDeductionPlanRepository) FindByPlanNo(planNo string) (*models.DeductionPlan, error) {
	for _, plan := range m.plans {
		if plan.PlanNo == planNo {
			return plan, nil
		}
	}
	return nil, nil
}

func (m *MockDeductionPlanRepository) FindByDeductee(deducteeID int64, status []int16, limit, offset int) ([]*models.DeductionPlan, int64, error) {
	var result []*models.DeductionPlan
	for _, plan := range m.plans {
		if plan.DeducteeID == deducteeID {
			result = append(result, plan)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockDeductionPlanRepository) FindByDeductor(deductorID int64, status []int16, limit, offset int) ([]*models.DeductionPlan, int64, error) {
	var result []*models.DeductionPlan
	for _, plan := range m.plans {
		if plan.DeductorID == deductorID {
			result = append(result, plan)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockDeductionPlanRepository) FindActivePlans(limit int) ([]*models.DeductionPlan, error) {
	var result []*models.DeductionPlan
	for _, plan := range m.plans {
		if plan.Status == models.DeductionPlanStatusActive {
			result = append(result, plan)
		}
	}
	return result, nil
}

func (m *MockDeductionPlanRepository) UpdateStatus(id int64, status int16) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if plan, ok := m.plans[id]; ok {
		plan.Status = status
		if status == models.DeductionPlanStatusCompleted {
			now := time.Now()
			plan.CompletedAt = &now
		}
	}
	return nil
}

func (m *MockDeductionPlanRepository) UpdateDeductedAmount(id int64, amount int64, currentPeriod int) error {
	if plan, ok := m.plans[id]; ok {
		plan.DeductedAmount += amount
		plan.RemainingAmount -= amount
		plan.CurrentPeriod = currentPeriod
	}
	return nil
}

func (m *MockDeductionPlanRepository) List(offset, limit int, status, planType int16) ([]*models.DeductionPlan, int64, error) {
	var result []*models.DeductionPlan
	for _, plan := range m.plans {
		if status > 0 && plan.Status != status {
			continue
		}
		if planType > 0 && plan.PlanType != planType {
			continue
		}
		result = append(result, plan)
	}
	return result, int64(len(result)), nil
}

// MockDeductionRecordRepository 模拟代扣记录仓库
type MockDeductionRecordRepository struct {
	records   map[int64]*models.DeductionRecord
	nextID    int64
	createErr error
}

func NewMockDeductionRecordRepository() *MockDeductionRecordRepository {
	return &MockDeductionRecordRepository{
		records: make(map[int64]*models.DeductionRecord),
		nextID:  1,
	}
}

func (m *MockDeductionRecordRepository) Create(record *models.DeductionRecord) error {
	if m.createErr != nil {
		return m.createErr
	}
	record.ID = m.nextID
	m.nextID++
	m.records[record.ID] = record
	return nil
}

func (m *MockDeductionRecordRepository) BatchCreate(records []*models.DeductionRecord) error {
	for _, record := range records {
		if err := m.Create(record); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockDeductionRecordRepository) FindByID(id int64) (*models.DeductionRecord, error) {
	record, ok := m.records[id]
	if !ok {
		return nil, nil
	}
	return record, nil
}

func (m *MockDeductionRecordRepository) FindByPlanID(planID int64) ([]*models.DeductionRecord, error) {
	var result []*models.DeductionRecord
	for _, record := range m.records {
		if record.PlanID == planID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (m *MockDeductionRecordRepository) FindPendingRecords(scheduledBefore time.Time, limit int) ([]*models.DeductionRecord, error) {
	var result []*models.DeductionRecord
	for _, record := range m.records {
		if record.Status == models.DeductionRecordStatusPending && !record.ScheduledAt.After(scheduledBefore) {
			result = append(result, record)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockDeductionRecordRepository) UpdateStatus(id int64, status int16, actualAmount int64, walletDetails string, failReason string) error {
	if record, ok := m.records[id]; ok {
		record.Status = status
		record.ActualAmount = actualAmount
		record.WalletDetails = walletDetails
		record.FailReason = failReason
		now := time.Now()
		record.DeductedAt = &now
	}
	return nil
}

// MockDeductionChainRepository 模拟代扣链仓库
type MockDeductionChainRepository struct {
	chains map[int64]*models.DeductionChain
	nextID int64
}

func NewMockDeductionChainRepository() *MockDeductionChainRepository {
	return &MockDeductionChainRepository{
		chains: make(map[int64]*models.DeductionChain),
		nextID: 1,
	}
}

func (m *MockDeductionChainRepository) Create(chain *models.DeductionChain) error {
	chain.ID = m.nextID
	m.nextID++
	m.chains[chain.ID] = chain
	return nil
}

func (m *MockDeductionChainRepository) FindByID(id int64) (*models.DeductionChain, error) {
	chain, ok := m.chains[id]
	if !ok {
		return nil, nil
	}
	return chain, nil
}

func (m *MockDeductionChainRepository) FindByDistributeID(distributeID int64) (*models.DeductionChain, error) {
	for _, chain := range m.chains {
		if chain.DistributeID == distributeID {
			return chain, nil
		}
	}
	return nil, nil
}

func (m *MockDeductionChainRepository) UpdateStatus(id int64, status int16) error {
	if chain, ok := m.chains[id]; ok {
		chain.Status = status
	}
	return nil
}

// MockDeductionChainItemRepository 模拟代扣链节点仓库
type MockDeductionChainItemRepository struct {
	items  map[int64]*models.DeductionChainItem
	nextID int64
}

func NewMockDeductionChainItemRepository() *MockDeductionChainItemRepository {
	return &MockDeductionChainItemRepository{
		items:  make(map[int64]*models.DeductionChainItem),
		nextID: 1,
	}
}

func (m *MockDeductionChainItemRepository) Create(item *models.DeductionChainItem) error {
	item.ID = m.nextID
	m.nextID++
	m.items[item.ID] = item
	return nil
}

func (m *MockDeductionChainItemRepository) BatchCreate(items []*models.DeductionChainItem) error {
	for _, item := range items {
		if err := m.Create(item); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockDeductionChainItemRepository) FindByChainID(chainID int64) ([]*models.DeductionChainItem, error) {
	var result []*models.DeductionChainItem
	for _, item := range m.items {
		if item.ChainID == chainID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockDeductionChainItemRepository) UpdatePlanID(id int64, planID int64) error {
	if item, ok := m.items[id]; ok {
		item.PlanID = planID
	}
	return nil
}

func (m *MockDeductionChainItemRepository) UpdateStatus(id int64, status int16) error {
	if item, ok := m.items[id]; ok {
		item.Status = status
	}
	return nil
}

// MockWalletRepository 模拟钱包仓库
type MockWalletRepository struct {
	wallets map[int64]*repository.Wallet
	nextID  int64
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallets: make(map[int64]*repository.Wallet),
		nextID:  1,
	}
}

func (m *MockWalletRepository) AddWallet(agentID, channelID int64, walletType int16, balance int64) *repository.Wallet {
	wallet := &repository.Wallet{
		ID:         m.nextID,
		AgentID:    agentID,
		ChannelID:  channelID,
		WalletType: walletType,
		Balance:    balance,
	}
	m.wallets[wallet.ID] = wallet
	m.nextID++
	return wallet
}

func (m *MockWalletRepository) FindByAgentAndType(agentID int64, channelID int64, walletType int16) (*repository.Wallet, error) {
	for _, wallet := range m.wallets {
		if wallet.AgentID == agentID && wallet.ChannelID == channelID && wallet.WalletType == walletType {
			return wallet, nil
		}
	}
	return nil, nil
}

func (m *MockWalletRepository) UpdateBalance(id int64, amount int64) error {
	if wallet, ok := m.wallets[id]; ok {
		wallet.Balance += amount
		wallet.TotalIncome += amount
	}
	return nil
}

func (m *MockWalletRepository) BatchUpdateBalance(updates map[int64]int64) error {
	for id, amount := range updates {
		if wallet, ok := m.wallets[id]; ok {
			wallet.Balance += amount
			wallet.TotalIncome += amount
		}
	}
	return nil
}

// MockWalletLogRepository 模拟钱包流水仓库
type MockWalletLogRepository struct {
	logs   []*repository.WalletLog
	nextID int64
}

func NewMockWalletLogRepository() *MockWalletLogRepository {
	return &MockWalletLogRepository{
		logs:   make([]*repository.WalletLog, 0),
		nextID: 1,
	}
}

func (m *MockWalletLogRepository) Create(log *repository.WalletLog) error {
	log.ID = m.nextID
	m.nextID++
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockWalletLogRepository) BatchCreate(logs []*repository.WalletLog) error {
	for _, log := range logs {
		if err := m.Create(log); err != nil {
			return err
		}
	}
	return nil
}

// MockAgentRepository 模拟代理商仓库
type MockAgentRepository struct {
	agents map[int64]*repository.Agent
}

func NewMockAgentRepository() *MockAgentRepository {
	return &MockAgentRepository{
		agents: make(map[int64]*repository.Agent),
	}
}

func (m *MockAgentRepository) AddAgent(id int64, agentNo string, parentID int64, path string, level int) {
	m.agents[id] = &repository.Agent{
		ID:       id,
		AgentNo:  agentNo,
		ParentID: parentID,
		Path:     path,
		Level:    level,
		Status:   1,
	}
}

func (m *MockAgentRepository) FindByID(id int64) (*repository.Agent, error) {
	agent, ok := m.agents[id]
	if !ok {
		return nil, nil
	}
	return agent, nil
}

func (m *MockAgentRepository) FindByAgentNo(agentNo string) (*repository.Agent, error) {
	for _, agent := range m.agents {
		if agent.AgentNo == agentNo {
			return agent, nil
		}
	}
	return nil, nil
}

func (m *MockAgentRepository) FindAncestors(agentID int64) ([]*repository.Agent, error) {
	agent, ok := m.agents[agentID]
	if !ok {
		return nil, nil
	}
	// 简化实现：根据path解析上级
	var ancestors []*repository.Agent
	currentID := agent.ParentID
	for currentID != 0 {
		if parent, ok := m.agents[currentID]; ok {
			ancestors = append(ancestors, parent)
			currentID = parent.ParentID
		} else {
			break
		}
	}
	return ancestors, nil
}

// =============================================================================
// 测试用例
// =============================================================================

// TestCreateDeductionPlan_Success 测试成功创建代扣计划
func TestCreateDeductionPlan_Success(t *testing.T) {
	// 准备
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	// 添加测试代理商
	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	// 执行 - Q6测试：伙伴代扣（任意代理商之间）
	req := &CreateDeductionPlanRequest{
		DeductorID:   1,
		DeducteeID:   2,
		PlanType:     models.DeductionPlanTypePartner, // 伙伴代扣
		TotalAmount:  120000,                          // 1200元
		TotalPeriods: 12,                              // 12期
		RelatedType:  "partner_loan",
		Remark:       "伙伴代扣测试",
		CreatedBy:    1,
	}

	plan, err := service.CreateDeductionPlan(req)

	// 验证
	if err != nil {
		t.Fatalf("CreateDeductionPlan failed: %v", err)
	}

	if plan == nil {
		t.Fatal("Plan should not be nil")
	}

	if plan.ID == 0 {
		t.Error("Plan ID should be set")
	}

	if plan.TotalAmount != 120000 {
		t.Errorf("TotalAmount = %d, want 120000", plan.TotalAmount)
	}

	if plan.PeriodAmount != 10000 {
		t.Errorf("PeriodAmount = %d, want 10000", plan.PeriodAmount)
	}

	if plan.Status != models.DeductionPlanStatusActive {
		t.Errorf("Status = %d, want %d", plan.Status, models.DeductionPlanStatusActive)
	}

	// 验证生成的代扣记录数量
	records, _ := recordRepo.FindByPlanID(plan.ID)
	if len(records) != 12 {
		t.Errorf("Generated %d records, want 12", len(records))
	}
}

// TestCreateDeductionPlan_PartnerDeduction_AnyAgent 测试伙伴代扣不限层级关系（Q6）
func TestCreateDeductionPlan_PartnerDeduction_AnyAgent(t *testing.T) {
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	// 添加测试代理商（无层级关系）
	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 0, "/", 1)  // 与A001无层级关系
	agentRepo.AddAgent(3, "A003", 2, "/2/", 2)

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	testCases := []struct {
		name       string
		deductorID int64
		deducteeID int64
		wantErr    bool
	}{
		{"同级代理商之间代扣", 1, 2, false},
		{"非直属上下级代扣", 1, 3, false},
		{"自己代扣自己", 1, 1, false}, // 业务规则允许
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &CreateDeductionPlanRequest{
				DeductorID:   tc.deductorID,
				DeducteeID:   tc.deducteeID,
				PlanType:     models.DeductionPlanTypePartner,
				TotalAmount:  10000,
				TotalPeriods: 1,
				CreatedBy:    1,
			}

			_, err := service.CreateDeductionPlan(req)

			if tc.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestCreateDeductionPlan_InvalidAgent 测试代理商不存在
func TestCreateDeductionPlan_InvalidAgent(t *testing.T) {
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	// 被扣款方不存在
	req := &CreateDeductionPlanRequest{
		DeductorID:   1,
		DeducteeID:   999, // 不存在
		PlanType:     models.DeductionPlanTypePartner,
		TotalAmount:  10000,
		TotalPeriods: 1,
		CreatedBy:    1,
	}

	_, err := service.CreateDeductionPlan(req)
	if err == nil {
		t.Error("Expected error for invalid deductee")
	}
}

// TestPauseAndResumeDeductionPlan 测试暂停和恢复代扣计划
func TestPauseAndResumeDeductionPlan(t *testing.T) {
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	// 创建计划
	req := &CreateDeductionPlanRequest{
		DeductorID:   1,
		DeducteeID:   2,
		PlanType:     models.DeductionPlanTypeGoods,
		TotalAmount:  10000,
		TotalPeriods: 1,
		CreatedBy:    1,
	}
	plan, _ := service.CreateDeductionPlan(req)

	// 暂停
	err := service.PauseDeductionPlan(plan.ID)
	if err != nil {
		t.Fatalf("PauseDeductionPlan failed: %v", err)
	}

	pausedPlan, _ := planRepo.FindByID(plan.ID)
	if pausedPlan.Status != models.DeductionPlanStatusPaused {
		t.Errorf("Status = %d, want %d", pausedPlan.Status, models.DeductionPlanStatusPaused)
	}

	// 恢复
	err = service.ResumeDeductionPlan(plan.ID)
	if err != nil {
		t.Fatalf("ResumeDeductionPlan failed: %v", err)
	}

	resumedPlan, _ := planRepo.FindByID(plan.ID)
	if resumedPlan.Status != models.DeductionPlanStatusActive {
		t.Errorf("Status = %d, want %d", resumedPlan.Status, models.DeductionPlanStatusActive)
	}
}

// TestCancelDeductionPlan 测试取消代扣计划
func TestCancelDeductionPlan(t *testing.T) {
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	// 创建计划
	req := &CreateDeductionPlanRequest{
		DeductorID:   1,
		DeducteeID:   2,
		PlanType:     models.DeductionPlanTypeGoods,
		TotalAmount:  10000,
		TotalPeriods: 1,
		CreatedBy:    1,
	}
	plan, _ := service.CreateDeductionPlan(req)

	// 取消
	err := service.CancelDeductionPlan(plan.ID)
	if err != nil {
		t.Fatalf("CancelDeductionPlan failed: %v", err)
	}

	cancelledPlan, _ := planRepo.FindByID(plan.ID)
	if cancelledPlan.Status != models.DeductionPlanStatusCancelled {
		t.Errorf("Status = %d, want %d", cancelledPlan.Status, models.DeductionPlanStatusCancelled)
	}
}

// TestCreateDeductionChain_CrossLevel 测试跨级下发创建代扣链（Q16）
func TestCreateDeductionChain_CrossLevel(t *testing.T) {
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	// A -> B -> C 层级结构
	agentRepo.AddAgent(1, "A", 0, "/", 1)
	agentRepo.AddAgent(2, "B", 1, "/1/", 2)
	agentRepo.AddAgent(3, "C", 2, "/1/2/", 3)

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	// 创建跨级代扣链：A 下发给 C，需要生成 C→B→A 的代扣链
	req := &CreateDeductionChainRequest{
		DistributeID: 1,
		TerminalSN:   "SN123456",
		AgentPath:    []int64{1, 2, 3}, // A, B, C
		TotalAmount:  100000,           // 1000元
		TotalPeriods: 12,
		CreatedBy:    1,
	}

	chain, err := service.CreateDeductionChain(req)

	if err != nil {
		t.Fatalf("CreateDeductionChain failed: %v", err)
	}

	if chain == nil {
		t.Fatal("Chain should not be nil")
	}

	if chain.TotalLevels != 2 { // A→B, B→C 共2层
		t.Errorf("TotalLevels = %d, want 2", chain.TotalLevels)
	}

	// 验证生成的代扣链节点
	items, _ := chainItemRepo.FindByChainID(chain.ID)
	if len(items) != 2 {
		t.Errorf("Generated %d chain items, want 2", len(items))
	}

	// 验证生成的代扣计划
	plans, _, _ := planRepo.FindByDeductee(2, nil, 100, 0)
	if len(plans) == 0 {
		t.Error("Should have created deduction plan for agent B")
	}
}

// TestCreateDeductionChain_MinimumAgents 测试代扣链最少需要2个代理商
func TestCreateDeductionChain_MinimumAgents(t *testing.T) {
	planRepo := NewMockDeductionPlanRepository()
	recordRepo := NewMockDeductionRecordRepository()
	chainRepo := NewMockDeductionChainRepository()
	chainItemRepo := NewMockDeductionChainItemRepository()
	walletRepo := NewMockWalletRepository()
	walletLogRepo := NewMockWalletLogRepository()
	agentRepo := NewMockAgentRepository()

	service := NewDeductionService(
		planRepo, recordRepo, chainRepo, chainItemRepo,
		walletRepo, walletLogRepo, agentRepo,
	)

	req := &CreateDeductionChainRequest{
		DistributeID: 1,
		TerminalSN:   "SN123456",
		AgentPath:    []int64{1}, // 只有1个代理商
		TotalAmount:  100000,
		TotalPeriods: 12,
		CreatedBy:    1,
	}

	_, err := service.CreateDeductionChain(req)
	if err == nil {
		t.Error("Expected error for agent path with less than 2 agents")
	}
}

// TestGetWalletTypeName 测试获取钱包类型名称
func TestGetWalletTypeName(t *testing.T) {
	testCases := []struct {
		walletType int16
		want       string
	}{
		{1, "分润钱包"},
		{2, "服务费钱包"},
		{3, "奖励钱包"},
		{99, "未知钱包"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			got := getWalletTypeName(tc.walletType)
			if got != tc.want {
				t.Errorf("getWalletTypeName(%d) = %s, want %s", tc.walletType, got, tc.want)
			}
		})
	}
}
