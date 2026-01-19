package service

import (
	"testing"
	"time"

	"xiangshoufu/internal/models"
)

// MockTerminalRepository 模拟终端仓库
type MockTerminalRepository struct {
	terminals map[int64]*models.Terminal
	nextID    int64
}

func NewMockTerminalRepository() *MockTerminalRepository {
	return &MockTerminalRepository{
		terminals: make(map[int64]*models.Terminal),
		nextID:    1,
	}
}

func (m *MockTerminalRepository) AddTerminal(sn string, ownerAgentID int64, channelID int64, status int16) *models.Terminal {
	terminal := &models.Terminal{
		ID:           m.nextID,
		TerminalSN:   sn,
		OwnerAgentID: ownerAgentID,
		ChannelID:    channelID,
		Status:       status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	m.terminals[terminal.ID] = terminal
	m.nextID++
	return terminal
}

func (m *MockTerminalRepository) Create(terminal *models.Terminal) error {
	terminal.ID = m.nextID
	m.nextID++
	m.terminals[terminal.ID] = terminal
	return nil
}

func (m *MockTerminalRepository) Update(terminal *models.Terminal) error {
	m.terminals[terminal.ID] = terminal
	return nil
}

func (m *MockTerminalRepository) FindByID(id int64) (*models.Terminal, error) {
	terminal, ok := m.terminals[id]
	if !ok {
		return nil, nil
	}
	return terminal, nil
}

func (m *MockTerminalRepository) FindBySN(terminalSN string) (*models.Terminal, error) {
	for _, terminal := range m.terminals {
		if terminal.TerminalSN == terminalSN {
			return terminal, nil
		}
	}
	return nil, nil
}

func (m *MockTerminalRepository) FindByOwner(ownerAgentID int64, status []int16, limit, offset int) ([]*models.Terminal, int64, error) {
	var result []*models.Terminal
	for _, terminal := range m.terminals {
		if terminal.OwnerAgentID == ownerAgentID {
			result = append(result, terminal)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockTerminalRepository) UpdateOwner(id int64, newOwnerID int64) error {
	if terminal, ok := m.terminals[id]; ok {
		terminal.OwnerAgentID = newOwnerID
		terminal.Status = models.TerminalStatusAllocated
	}
	return nil
}

func (m *MockTerminalRepository) UpdateStatus(id int64, status int16) error {
	if terminal, ok := m.terminals[id]; ok {
		terminal.Status = status
	}
	return nil
}

func (m *MockTerminalRepository) UpdateSimFeeCount(id int64, count int) error {
	if terminal, ok := m.terminals[id]; ok {
		terminal.SimFeeCount = count
		now := time.Now()
		terminal.LastSimFeeAt = &now
	}
	return nil
}

// MockTerminalDistributeRepository 模拟终端下发仓库
type MockTerminalDistributeRepository struct {
	distributes map[int64]*models.TerminalDistribute
	nextID      int64
}

func NewMockTerminalDistributeRepository() *MockTerminalDistributeRepository {
	return &MockTerminalDistributeRepository{
		distributes: make(map[int64]*models.TerminalDistribute),
		nextID:      1,
	}
}

func (m *MockTerminalDistributeRepository) Create(distribute *models.TerminalDistribute) error {
	distribute.ID = m.nextID
	m.nextID++
	m.distributes[distribute.ID] = distribute
	return nil
}

func (m *MockTerminalDistributeRepository) Update(distribute *models.TerminalDistribute) error {
	m.distributes[distribute.ID] = distribute
	return nil
}

func (m *MockTerminalDistributeRepository) FindByID(id int64) (*models.TerminalDistribute, error) {
	distribute, ok := m.distributes[id]
	if !ok {
		return nil, nil
	}
	return distribute, nil
}

func (m *MockTerminalDistributeRepository) FindByDistributeNo(distributeNo string) (*models.TerminalDistribute, error) {
	for _, distribute := range m.distributes {
		if distribute.DistributeNo == distributeNo {
			return distribute, nil
		}
	}
	return nil, nil
}

func (m *MockTerminalDistributeRepository) FindByFromAgent(fromAgentID int64, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error) {
	var result []*models.TerminalDistribute
	for _, distribute := range m.distributes {
		if distribute.FromAgentID == fromAgentID {
			result = append(result, distribute)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockTerminalDistributeRepository) FindByToAgent(toAgentID int64, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error) {
	var result []*models.TerminalDistribute
	for _, distribute := range m.distributes {
		if distribute.ToAgentID == toAgentID {
			result = append(result, distribute)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockTerminalDistributeRepository) UpdateStatus(id int64, status int16, confirmedBy *int64) error {
	if distribute, ok := m.distributes[id]; ok {
		distribute.Status = status
		distribute.ConfirmedBy = confirmedBy
		now := time.Now()
		distribute.ConfirmedAt = &now
	}
	return nil
}

// =============================================================================
// 测试用例
// =============================================================================

// TestDistributeTerminal_Success_DirectLevel 测试直属下级下发成功
func TestDistributeTerminal_Success_DirectLevel(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	// 设置代理商层级：A -> B（直属）
	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	// 添加终端
	terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil, // deductionService可以为nil，因为这里只测试非分期场景
	)

	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     2,
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000, // 1000元
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourceApp, // APP端
		CreatedBy:     1,
	}

	distribute, err := service.DistributeTerminal(req)

	if err != nil {
		t.Fatalf("DistributeTerminal failed: %v", err)
	}

	if distribute == nil {
		t.Fatal("Distribute should not be nil")
	}

	if distribute.IsCrossLevel {
		t.Error("Should not be cross level for direct subordinate")
	}

	if distribute.Status != models.TerminalDistributeStatusPending {
		t.Errorf("Status = %d, want %d", distribute.Status, models.TerminalDistributeStatusPending)
	}
}

// TestDistributeTerminal_AppCrossLevel_Denied 测试APP端跨级下发被拒绝（Q29）
func TestDistributeTerminal_AppCrossLevel_Denied(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	// 设置代理商层级：A -> B -> C
	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)
	agentRepo.AddAgent(3, "A003", 2, "/1/2/", 3)

	// 添加终端
	terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// APP端尝试跨级下发 A -> C（应该被拒绝）
	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     3, // 跨过B，直接到C
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourceApp, // APP端
		CreatedBy:     1,
	}

	_, err := service.DistributeTerminal(req)

	if err == nil {
		t.Error("Expected error for APP cross-level distribute")
	}

	if err.Error() != "APP端不支持跨级下发，请使用PC端" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestDistributeTerminal_PCCrossLevel_Allowed 测试PC端跨级下发允许（Q29）
func TestDistributeTerminal_PCCrossLevel_Allowed(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	// 设置代理商层级：A -> B -> C
	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)
	agentRepo.AddAgent(3, "A003", 2, "/1/2/", 3)

	// 添加终端
	terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// PC端跨级下发 A -> C（应该允许）
	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     3, // 跨过B，直接到C
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourcePC, // PC端
		CreatedBy:     1,
	}

	distribute, err := service.DistributeTerminal(req)

	if err != nil {
		t.Fatalf("DistributeTerminal failed: %v", err)
	}

	if !distribute.IsCrossLevel {
		t.Error("Should be marked as cross level")
	}

	if distribute.CrossLevelPath == "" {
		t.Error("CrossLevelPath should not be empty")
	}
}

// TestDistributeTerminal_TerminalNotFound 测试终端不存在
func TestDistributeTerminal_TerminalNotFound(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     2,
		TerminalSN:    "NOTEXIST",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourcePC,
		CreatedBy:     1,
	}

	_, err := service.DistributeTerminal(req)

	if err == nil {
		t.Error("Expected error for non-existent terminal")
	}
}

// TestDistributeTerminal_TerminalNotOwned 测试终端不属于下发方
func TestDistributeTerminal_TerminalNotOwned(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)
	agentRepo.AddAgent(3, "A003", 1, "/1/", 2)

	// 终端属于代理商3
	terminalRepo.AddTerminal("SN123456", 3, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// 代理商1尝试下发不属于自己的终端
	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     2,
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourcePC,
		CreatedBy:     1,
	}

	_, err := service.DistributeTerminal(req)

	if err == nil {
		t.Error("Expected error for terminal not owned by distributor")
	}
}

// TestRejectDistribute 测试拒绝下发
func TestRejectDistribute(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// 创建下发
	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     2,
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourcePC,
		CreatedBy:     1,
	}
	distribute, _ := service.DistributeTerminal(req)

	// 拒绝
	err := service.RejectDistribute(distribute.ID, 2)
	if err != nil {
		t.Fatalf("RejectDistribute failed: %v", err)
	}

	rejected, _ := distributeRepo.FindByID(distribute.ID)
	if rejected.Status != models.TerminalDistributeStatusRejected {
		t.Errorf("Status = %d, want %d", rejected.Status, models.TerminalDistributeStatusRejected)
	}
}

// TestCancelDistribute 测试取消下发
func TestCancelDistribute(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// 创建下发
	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     2,
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourcePC,
		CreatedBy:     1,
	}
	distribute, _ := service.DistributeTerminal(req)

	// 下发方取消
	err := service.CancelDistribute(distribute.ID, 1)
	if err != nil {
		t.Fatalf("CancelDistribute failed: %v", err)
	}

	cancelled, _ := distributeRepo.FindByID(distribute.ID)
	if cancelled.Status != models.TerminalDistributeStatusCancelled {
		t.Errorf("Status = %d, want %d", cancelled.Status, models.TerminalDistributeStatusCancelled)
	}
}

// TestCancelDistribute_NotOwner 测试非下发方不能取消
func TestCancelDistribute_NotOwner(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	terminalRepo.AddTerminal("SN123456", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// 创建下发
	req := &DistributeTerminalRequest{
		FromAgentID:   1,
		ToAgentID:     2,
		TerminalSN:    "SN123456",
		ChannelID:     1,
		GoodsPrice:    100000,
		DeductionType: models.TerminalDistributeDeductionOneTime,
		Source:        models.TerminalDistributeSourcePC,
		CreatedBy:     1,
	}
	distribute, _ := service.DistributeTerminal(req)

	// 接收方尝试取消（应该失败）
	err := service.CancelDistribute(distribute.ID, 2)
	if err == nil {
		t.Error("Expected error when non-owner tries to cancel")
	}
}

// TestGetDistributeList 测试获取下发列表
func TestGetDistributeList(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	agentRepo.AddAgent(1, "A001", 0, "/", 1)
	agentRepo.AddAgent(2, "A002", 1, "/1/", 2)

	terminalRepo.AddTerminal("SN001", 1, 1, models.TerminalStatusPending)
	terminalRepo.AddTerminal("SN002", 1, 1, models.TerminalStatusPending)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	// 创建两个下发
	for _, sn := range []string{"SN001", "SN002"} {
		req := &DistributeTerminalRequest{
			FromAgentID:   1,
			ToAgentID:     2,
			TerminalSN:    sn,
			ChannelID:     1,
			GoodsPrice:    100000,
			DeductionType: models.TerminalDistributeDeductionOneTime,
			Source:        models.TerminalDistributeSourcePC,
			CreatedBy:     1,
		}
		service.DistributeTerminal(req)
	}

	// 查询下发方的列表
	list, total, err := service.GetDistributeList(1, "from", nil, 10, 0)
	if err != nil {
		t.Fatalf("GetDistributeList failed: %v", err)
	}

	if total != 2 {
		t.Errorf("Total = %d, want 2", total)
	}

	if len(list) != 2 {
		t.Errorf("List length = %d, want 2", len(list))
	}

	// 查询接收方的列表
	list, total, err = service.GetDistributeList(2, "to", nil, 10, 0)
	if err != nil {
		t.Fatalf("GetDistributeList failed: %v", err)
	}

	if total != 2 {
		t.Errorf("Total = %d, want 2", total)
	}
}

// TestCheckCrossLevel 测试跨级检查辅助函数
func TestCheckCrossLevel(t *testing.T) {
	terminalRepo := NewMockTerminalRepository()
	distributeRepo := NewMockTerminalDistributeRepository()
	agentRepo := NewMockAgentRepository()

	// A -> B -> C
	agentRepo.AddAgent(1, "A", 0, "/", 1)
	agentRepo.AddAgent(2, "B", 1, "/1/", 2)
	agentRepo.AddAgent(3, "C", 2, "/1/2/", 3)

	service := NewTerminalDistributeService(
		terminalRepo,
		distributeRepo,
		agentRepo,
		nil,
	)

	testCases := []struct {
		name           string
		fromID         int64
		toID           int64
		wantCrossLevel bool
	}{
		{"Direct: A->B", 1, 2, false},
		{"Direct: B->C", 2, 3, false},
		{"Cross: A->C", 1, 3, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fromAgent, _ := agentRepo.FindByID(tc.fromID)
			toAgent, _ := agentRepo.FindByID(tc.toID)

			isCross, _, err := service.checkCrossLevel(fromAgent, toAgent)
			if err != nil {
				t.Fatalf("checkCrossLevel failed: %v", err)
			}

			if isCross != tc.wantCrossLevel {
				t.Errorf("isCrossLevel = %v, want %v", isCross, tc.wantCrossLevel)
			}
		})
	}
}
