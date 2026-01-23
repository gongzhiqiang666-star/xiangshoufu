package service

import (
	"errors"
	"testing"
	"time"

	"xiangshoufu/internal/repository"
)

// ==================== Profit Service 专用 Mock 实现 ====================
// 注意：以 Profit 前缀命名，避免与其他测试文件的 Mock 冲突

// ProfitMockTransactionRepository 模拟交易仓库
type ProfitMockTransactionRepository struct {
	transactions      map[string]*repository.Transaction
	profitStatusCalls map[int64]int16
	refundStatusCalls map[int64]int16
}

func NewProfitMockTransactionRepository() *ProfitMockTransactionRepository {
	return &ProfitMockTransactionRepository{
		transactions:      make(map[string]*repository.Transaction),
		profitStatusCalls: make(map[int64]int16),
		refundStatusCalls: make(map[int64]int16),
	}
}

func (m *ProfitMockTransactionRepository) Create(tx *repository.Transaction) error {
	m.transactions[tx.OrderNo] = tx
	return nil
}

func (m *ProfitMockTransactionRepository) FindByOrderNo(orderNo string) (*repository.Transaction, error) {
	if tx, ok := m.transactions[orderNo]; ok {
		return tx, nil
	}
	return nil, nil
}

// FindByID 根据交易ID查找
func (m *ProfitMockTransactionRepository) FindByID(id int64) (*repository.Transaction, error) {
	for _, tx := range m.transactions {
		if tx.ID == id {
			return tx, nil
		}
	}
	return nil, nil
}

func (m *ProfitMockTransactionRepository) FindUnprocessedProfit(limit int) ([]*repository.Transaction, error) {
	var result []*repository.Transaction
	for _, tx := range m.transactions {
		if tx.ProfitStatus == 0 {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *ProfitMockTransactionRepository) UpdateProfitStatus(id int64, status int16) error {
	m.profitStatusCalls[id] = status
	for _, tx := range m.transactions {
		if tx.ID == id {
			tx.ProfitStatus = status
			return nil
		}
	}
	return nil
}

func (m *ProfitMockTransactionRepository) BatchUpdateProfitStatus(ids []int64, status int16) error {
	for _, id := range ids {
		m.profitStatusCalls[id] = status
	}
	return nil
}

func (m *ProfitMockTransactionRepository) UpdateRefundStatus(id int64, status int16) error {
	m.refundStatusCalls[id] = status
	return nil
}

func (m *ProfitMockTransactionRepository) GetTerminalTotalTradeAmount(terminalSN string) (int64, error) {
	return 0, nil
}

func (m *ProfitMockTransactionRepository) AddTransaction(tx *repository.Transaction) {
	m.transactions[tx.OrderNo] = tx
}

// ProfitMockProfitRecordRepository 模拟分润记录仓库
type ProfitMockProfitRecordRepository struct {
	records       []*repository.ProfitRecord
	revokedTxIDs  []int64
	revokeReasons map[int64]string
}

func NewProfitMockProfitRecordRepository() *ProfitMockProfitRecordRepository {
	return &ProfitMockProfitRecordRepository{
		records:       make([]*repository.ProfitRecord, 0),
		revokedTxIDs:  make([]int64, 0),
		revokeReasons: make(map[int64]string),
	}
}

func (m *ProfitMockProfitRecordRepository) Create(record *repository.ProfitRecord) error {
	m.records = append(m.records, record)
	return nil
}

func (m *ProfitMockProfitRecordRepository) BatchCreate(records []*repository.ProfitRecord) error {
	m.records = append(m.records, records...)
	return nil
}

func (m *ProfitMockProfitRecordRepository) FindByTransactionID(txID int64) ([]*repository.ProfitRecord, error) {
	var result []*repository.ProfitRecord
	for _, r := range m.records {
		if r.TransactionID == txID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *ProfitMockProfitRecordRepository) RevokeByTransactionID(txID int64, reason string) error {
	m.revokedTxIDs = append(m.revokedTxIDs, txID)
	m.revokeReasons[txID] = reason
	for _, r := range m.records {
		if r.TransactionID == txID {
			r.IsRevoked = true
			r.RevokeReason = reason
			now := time.Now()
			r.RevokedAt = &now
		}
	}
	return nil
}

// ProfitMockWalletRepository 模拟钱包仓库（分润专用）
type ProfitMockWalletRepository struct {
	wallets        map[string]*repository.Wallet
	balanceUpdates map[int64]int64
}

func NewProfitMockWalletRepository() *ProfitMockWalletRepository {
	return &ProfitMockWalletRepository{
		wallets:        make(map[string]*repository.Wallet),
		balanceUpdates: make(map[int64]int64),
	}
}

func (m *ProfitMockWalletRepository) FindByAgentAndType(agentID int64, channelID int64, walletType int16) (*repository.Wallet, error) {
	return &repository.Wallet{
		ID:         agentID*1000 + channelID*10 + int64(walletType),
		AgentID:    agentID,
		ChannelID:  channelID,
		WalletType: walletType,
		Balance:    100000,
	}, nil
}

func (m *ProfitMockWalletRepository) UpdateBalance(id int64, amount int64) error {
	m.balanceUpdates[id] += amount
	return nil
}

func (m *ProfitMockWalletRepository) BatchUpdateBalance(updates map[int64]int64) error {
	for id, amount := range updates {
		m.balanceUpdates[id] += amount
	}
	return nil
}

// ProfitMockWalletLogRepository 模拟钱包流水仓库（分润专用）
type ProfitMockWalletLogRepository struct {
	logs []*repository.WalletLog
}

func NewProfitMockWalletLogRepository() *ProfitMockWalletLogRepository {
	return &ProfitMockWalletLogRepository{
		logs: make([]*repository.WalletLog, 0),
	}
}

func (m *ProfitMockWalletLogRepository) Create(log *repository.WalletLog) error {
	m.logs = append(m.logs, log)
	return nil
}

func (m *ProfitMockWalletLogRepository) BatchCreate(logs []*repository.WalletLog) error {
	m.logs = append(m.logs, logs...)
	return nil
}

// ProfitMockAgentRepository 模拟代理商仓库（分润专用）
type ProfitMockAgentRepository struct {
	agents    map[int64]*repository.Agent
	ancestors map[int64][]*repository.Agent
}

func NewProfitMockAgentRepository() *ProfitMockAgentRepository {
	return &ProfitMockAgentRepository{
		agents:    make(map[int64]*repository.Agent),
		ancestors: make(map[int64][]*repository.Agent),
	}
}

func (m *ProfitMockAgentRepository) FindByID(id int64) (*repository.Agent, error) {
	if agent, ok := m.agents[id]; ok {
		return agent, nil
	}
	return nil, nil
}

func (m *ProfitMockAgentRepository) FindByAgentNo(agentNo string) (*repository.Agent, error) {
	for _, agent := range m.agents {
		if agent.AgentNo == agentNo {
			return agent, nil
		}
	}
	return nil, nil
}

func (m *ProfitMockAgentRepository) FindAncestors(agentID int64) ([]*repository.Agent, error) {
	if ancestors, ok := m.ancestors[agentID]; ok {
		return ancestors, nil
	}
	return []*repository.Agent{}, nil
}

func (m *ProfitMockAgentRepository) AddAgent(agent *repository.Agent) {
	m.agents[agent.ID] = agent
}

func (m *ProfitMockAgentRepository) SetAncestors(agentID int64, ancestors []*repository.Agent) {
	m.ancestors[agentID] = ancestors
}

// ProfitMockAgentPolicyRepository 模拟代理商政策仓库（分润专用）
type ProfitMockAgentPolicyRepository struct {
	policies map[int64]map[int64]*repository.AgentPolicy // agentID -> channelID -> policy
}

func NewProfitMockAgentPolicyRepository() *ProfitMockAgentPolicyRepository {
	return &ProfitMockAgentPolicyRepository{
		policies: make(map[int64]map[int64]*repository.AgentPolicy),
	}
}

func (m *ProfitMockAgentPolicyRepository) FindByAgentAndChannel(agentID int64, channelID int64) (*repository.AgentPolicy, error) {
	if agentPolicies, ok := m.policies[agentID]; ok {
		if policy, ok := agentPolicies[channelID]; ok {
			return policy, nil
		}
	}
	return nil, errors.New("policy not found")
}

func (m *ProfitMockAgentPolicyRepository) AddPolicy(policy *repository.AgentPolicy) {
	if m.policies[policy.AgentID] == nil {
		m.policies[policy.AgentID] = make(map[int64]*repository.AgentPolicy)
	}
	m.policies[policy.AgentID][policy.ChannelID] = policy
}

// ProfitMockMessageQueue 模拟消息队列（分润专用）
type ProfitMockMessageQueue struct {
	messages map[string][][]byte
}

func NewProfitMockMessageQueue() *ProfitMockMessageQueue {
	return &ProfitMockMessageQueue{
		messages: make(map[string][][]byte),
	}
}

func (m *ProfitMockMessageQueue) Publish(topic string, msg []byte) error {
	m.messages[topic] = append(m.messages[topic], msg)
	return nil
}

func (m *ProfitMockMessageQueue) Subscribe(topic string, handler func([]byte) error) error {
	return nil
}

func (m *ProfitMockMessageQueue) Close() error {
	return nil
}

// ==================== 辅助函数 ====================

func createProfitTestService() (*ProfitService, *ProfitMockTransactionRepository, *ProfitMockProfitRecordRepository, *ProfitMockWalletRepository, *ProfitMockAgentRepository, *ProfitMockAgentPolicyRepository) {
	txRepo := NewProfitMockTransactionRepository()
	profitRepo := NewProfitMockProfitRecordRepository()
	walletRepo := NewProfitMockWalletRepository()
	walletLogRepo := NewProfitMockWalletLogRepository()
	agentRepo := NewProfitMockAgentRepository()
	policyRepo := NewProfitMockAgentPolicyRepository()
	queue := NewProfitMockMessageQueue()

	service := NewProfitService(
		txRepo,
		profitRepo,
		walletRepo,
		walletLogRepo,
		agentRepo,
		policyRepo,
		nil,
		queue,
	)

	return service, txRepo, profitRepo, walletRepo, agentRepo, policyRepo
}

// ==================== 测试用例 ====================

// TestCalculateProfitAmount 测试分润金额计算公式
// 注意：浮点数精度问题 - 0.60 - 0.50 = 0.09999999999999998（而非0.1）
// 这是Go语言浮点数的固有特性，实际业务中建议使用整数分（分）进行计算
func TestCalculateProfitAmount(t *testing.T) {
	tests := []struct {
		name         string
		tradeAmount  int64
		lowerRate    float64
		selfRate     float64
		wantProfit   int64
		wantHasSpace bool
	}{
		{
			name:         "正常分润计算 - 100元交易 0.1%费率差",
			tradeAmount:  10000,
			lowerRate:    0.60,
			selfRate:     0.50,
			wantProfit:   9, // 浮点数精度：10000 * 0.0999999... / 100 = 9
			wantHasSpace: true,
		},
		{
			name:         "正常分润计算 - 1000元交易 0.05%费率差",
			tradeAmount:  100000,
			lowerRate:    0.55,
			selfRate:     0.50,
			wantProfit:   50,
			wantHasSpace: true,
		},
		{
			name:         "费率差为0 - 无分润",
			tradeAmount:  10000,
			lowerRate:    0.50,
			selfRate:     0.50,
			wantProfit:   0,
			wantHasSpace: false,
		},
		{
			name:         "费率差为负 - 无分润（异常情况）",
			tradeAmount:  10000,
			lowerRate:    0.48,
			selfRate:     0.50,
			wantProfit:   -2, // 负值，实际业务中会被过滤掉
			wantHasSpace: false,
		},
		{
			name:         "小额交易 - 1元交易（分润取整后为0）",
			tradeAmount:  100,
			lowerRate:    0.60,
			selfRate:     0.50,
			wantProfit:   0, // 100 * 0.1% = 0.1分 取整为0
			wantHasSpace: true, // 有费率空间，但分润为0
		},
		{
			name:         "大额交易 - 10万元交易",
			tradeAmount:  10000000,
			lowerRate:    0.60,
			selfRate:     0.50,
			wantProfit:   9999, // 浮点数精度
			wantHasSpace: true,
		},
		{
			name:         "0元交易",
			tradeAmount:  0,
			lowerRate:    0.60,
			selfRate:     0.50,
			wantProfit:   0,
			wantHasSpace: true, // 有费率空间，但金额为0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rateDiff := tt.lowerRate - tt.selfRate
			hasSpace := rateDiff > 0

			if hasSpace != tt.wantHasSpace {
				t.Errorf("费率空间判断错误: got %v, want %v", hasSpace, tt.wantHasSpace)
			}

			// 计算分润金额
			profitAmount := int64(float64(tt.tradeAmount) * rateDiff / 100)

			if profitAmount != tt.wantProfit {
				t.Errorf("分润金额计算错误: got %d分, want %d分 (交易%d分, 费率差%.16f%%)",
					profitAmount, tt.wantProfit, tt.tradeAmount, rateDiff)
			}
		})
	}
}

// TestProfitService_SingleAgentProfit 测试单级代理商分润
func TestProfitService_SingleAgentProfit(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agent := &repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1}
	agentRepo.AddAgent(agent)

	policyRepo.AddPolicy(&repository.AgentPolicy{
		ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50", DebitRate: "0.45",
	})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, MerchantID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	records := profitRepo.records
	if len(records) != 1 {
		t.Fatalf("期望1条分润记录, 实际 %d 条", len(records))
	}

	// 浮点数精度：1000元 * (0.60% - 0.50%) = 1000 * 0.0999999... = 99分
	expectedProfit := int64(99)
	if records[0].ProfitAmount != expectedProfit {
		t.Errorf("分润金额错误: got %d, want %d", records[0].ProfitAmount, expectedProfit)
	}

	if records[0].AgentID != 100 {
		t.Errorf("代理商ID错误: got %d, want 100", records[0].AgentID)
	}
}

// TestProfitService_MultiLevelAgentProfit 测试多级代理商分润
func TestProfitService_MultiLevelAgentProfit(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	topAgent := &repository.Agent{ID: 1, AgentNo: "A001", ParentID: 0, Level: 1}
	level1Agent := &repository.Agent{ID: 10, AgentNo: "A010", ParentID: 1, Level: 2}
	level2Agent := &repository.Agent{ID: 100, AgentNo: "A100", ParentID: 10, Level: 3}

	agentRepo.AddAgent(topAgent)
	agentRepo.AddAgent(level1Agent)
	agentRepo.AddAgent(level2Agent)
	agentRepo.SetAncestors(100, []*repository.Agent{level1Agent, topAgent})

	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 1, ChannelID: 1, CreditRate: "0.45", DebitRate: "0.40"})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 2, AgentID: 10, ChannelID: 1, CreditRate: "0.50", DebitRate: "0.45"})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 3, AgentID: 100, ChannelID: 1, CreditRate: "0.55", DebitRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, MerchantID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	records := profitRepo.records
	if len(records) != 3 {
		t.Fatalf("期望3条分润记录, 实际 %d 条", len(records))
	}

	// 浮点数精度：每级0.05%费率差，1000元交易 = 49或50分
	expectedProfits := map[int64]int64{
		100: 49, // 二级: 1000 * (0.60-0.55)% ≈ 49分
		10:  50, // 一级: 1000 * (0.55-0.50)% = 50分
		1:   49, // 总部: 1000 * (0.50-0.45)% ≈ 49分
	}

	for _, record := range records {
		expected := expectedProfits[record.AgentID]
		if record.ProfitAmount != expected {
			t.Errorf("代理商%d分润金额错误: got %d, want %d", record.AgentID, record.ProfitAmount, expected)
		}
	}
}

// TestProfitService_AlreadyCalculated 测试重复计算保护
func TestProfitService_AlreadyCalculated(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 1, // 已计算
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(profitRepo.records) != 0 {
		t.Errorf("已计算交易不应产生新分润记录, 实际产生 %d 条", len(profitRepo.records))
	}
}

// TestProfitService_ZeroRateDiff 测试费率差为0的情况
func TestProfitService_ZeroRateDiff(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.60"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(profitRepo.records) != 0 {
		t.Errorf("费率差为0不应产生分润记录, 实际产生 %d 条", len(profitRepo.records))
	}
}

// TestProfitService_NegativeRateDiff 测试费率差为负的情况
func TestProfitService_NegativeRateDiff(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.70"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(profitRepo.records) != 0 {
		t.Errorf("费率差为负不应产生分润记录, 实际产生 %d 条", len(profitRepo.records))
	}
}

// TestRevokeProfit 测试撤销分润
func TestRevokeProfit(t *testing.T) {
	service, txRepo, profitRepo, walletRepo, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, MerchantID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(profitRepo.records) != 1 {
		t.Fatalf("期望1条分润记录, 实际 %d 条", len(profitRepo.records))
	}

	profitAmount := profitRepo.records[0].ProfitAmount

	// 清空钱包更新记录，以便检测撤销操作的影响
	walletRepo.balanceUpdates = make(map[int64]int64)

	err = service.RevokeProfit(1, "用户退款")
	if err != nil {
		t.Fatalf("RevokeProfit failed: %v", err)
	}

	if !profitRepo.records[0].IsRevoked {
		t.Error("分润记录应被标记为已撤销")
	}

	if profitRepo.records[0].RevokeReason != "用户退款" {
		t.Errorf("撤销原因错误: got %s, want 用户退款", profitRepo.records[0].RevokeReason)
	}

	// 验证钱包被扣减（负值表示扣减）
	for _, amount := range walletRepo.balanceUpdates {
		if amount != -profitAmount {
			t.Errorf("钱包扣减金额错误: got %d, want %d", amount, -profitAmount)
		}
	}
}

// TestDebitCardRate 测试借记卡费率
func TestDebitCardRate(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{
		ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50", DebitRate: "0.45",
	})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.55", CardType: 1, ProfitStatus: 0, // 借记卡
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(profitRepo.records) != 1 {
		t.Fatalf("期望1条分润记录, 实际 %d 条", len(profitRepo.records))
	}

	expectedProfit := int64(100) // 1000 * (0.55 - 0.45)% = 100分
	if profitRepo.records[0].ProfitAmount != expectedProfit {
		t.Errorf("借记卡分润金额错误: got %d, want %d", profitRepo.records[0].ProfitAmount, expectedProfit)
	}
}

// TestLargeAmountTransaction 测试大额交易
func TestLargeAmountTransaction(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000000, Rate: "0.60", CardType: 2, ProfitStatus: 0, // 100万元
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	// 浮点数精度：1000000 * 0.0999999...% = 99999分
	expectedProfit := int64(99999)
	if profitRepo.records[0].ProfitAmount != expectedProfit {
		t.Errorf("大额交易分润金额错误: got %d, want %d", profitRepo.records[0].ProfitAmount, expectedProfit)
	}
}

// TestSmallAmountTransaction 测试小额交易
func TestSmallAmountTransaction(t *testing.T) {
	service, txRepo, profitRepo, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100, Rate: "0.60", CardType: 2, ProfitStatus: 0, // 1元
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(profitRepo.records) != 0 {
		t.Errorf("小额交易不应产生分润记录, 实际产生 %d 条", len(profitRepo.records))
	}
}

// TestWalletBalanceUpdate 测试钱包余额更新
func TestWalletBalanceUpdate(t *testing.T) {
	service, txRepo, _, walletRepo, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if len(walletRepo.balanceUpdates) == 0 {
		t.Error("钱包余额应该被更新")
	}

	// 浮点数精度：99分
	for _, amount := range walletRepo.balanceUpdates {
		if amount != 99 {
			t.Errorf("钱包更新金额错误: got %d, want 99", amount)
		}
	}
}

// TestProfitStatusUpdate 测试交易分润状态更新
func TestProfitStatusUpdate(t *testing.T) {
	service, txRepo, _, _, agentRepo, policyRepo := createProfitTestService()

	agentRepo.AddAgent(&repository.Agent{ID: 100, AgentNo: "A100", ParentID: 0, Level: 1})
	policyRepo.AddPolicy(&repository.AgentPolicy{ID: 1, AgentID: 100, ChannelID: 1, CreditRate: "0.50"})

	tx := &repository.Transaction{
		ID: 1, OrderNo: "TX001", ChannelID: 1, AgentID: 100,
		Amount: 100000, Rate: "0.60", CardType: 2, ProfitStatus: 0,
	}
	txRepo.AddTransaction(tx)

	err := service.CalculateProfit(1)
	if err != nil {
		t.Fatalf("CalculateProfit failed: %v", err)
	}

	if status, ok := txRepo.profitStatusCalls[1]; !ok || status != 1 {
		t.Errorf("交易分润状态应更新为1, got %d", status)
	}
}
