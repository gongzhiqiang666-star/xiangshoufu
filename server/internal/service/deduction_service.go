package service

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// DeductionService 代扣服务
// 业务规则：
// - Q6: 伙伴代扣 - 任意代理商之间都可以发起（不限层级关系）
// - Q7: 每日扣款 - 每天固定时间检查余额并扣款
// - Q8: 多钱包扣款优先级 - 按余额从多到少扣
type DeductionService struct {
	planRepo      repository.DeductionPlanRepository
	recordRepo    repository.DeductionRecordRepository
	chainRepo     repository.DeductionChainRepository
	chainItemRepo repository.DeductionChainItemRepository
	walletRepo    repository.WalletRepository
	walletLogRepo repository.WalletLogRepository
	agentRepo     repository.AgentRepository
}

// NewDeductionService 创建代扣服务
func NewDeductionService(
	planRepo repository.DeductionPlanRepository,
	recordRepo repository.DeductionRecordRepository,
	chainRepo repository.DeductionChainRepository,
	chainItemRepo repository.DeductionChainItemRepository,
	walletRepo repository.WalletRepository,
	walletLogRepo repository.WalletLogRepository,
	agentRepo repository.AgentRepository,
) *DeductionService {
	return &DeductionService{
		planRepo:      planRepo,
		recordRepo:    recordRepo,
		chainRepo:     chainRepo,
		chainItemRepo: chainItemRepo,
		walletRepo:    walletRepo,
		walletLogRepo: walletLogRepo,
		agentRepo:     agentRepo,
	}
}

// CreateDeductionPlanRequest 创建代扣计划请求
type CreateDeductionPlanRequest struct {
	DeductorID   int64  `json:"deductor_id"`   // 扣款方代理商ID
	DeducteeID   int64  `json:"deductee_id"`   // 被扣款方代理商ID
	PlanType     int16  `json:"plan_type"`     // 计划类型
	TotalAmount  int64  `json:"total_amount"`  // 总金额（分）
	TotalPeriods int    `json:"total_periods"` // 总期数
	RelatedType  string `json:"related_type"`  // 关联类型
	RelatedID    int64  `json:"related_id"`    // 关联ID
	Remark       string `json:"remark"`        // 备注
	CreatedBy    int64  `json:"created_by"`    // 创建人
}

// CreateDeductionPlan 创建代扣计划
// 支持伙伴代扣（任意代理商之间，不限层级关系）
func (s *DeductionService) CreateDeductionPlan(req *CreateDeductionPlanRequest) (*models.DeductionPlan, error) {
	// 验证扣款方和被扣款方是否存在
	deductor, err := s.agentRepo.FindByID(req.DeductorID)
	if err != nil || deductor == nil {
		return nil, fmt.Errorf("扣款方代理商不存在: %d", req.DeductorID)
	}

	deductee, err := s.agentRepo.FindByID(req.DeducteeID)
	if err != nil || deductee == nil {
		return nil, fmt.Errorf("被扣款方代理商不存在: %d", req.DeducteeID)
	}

	// 注意：伙伴代扣不限层级关系，任意代理商之间都可以发起

	// 计算每期金额
	periodAmount := req.TotalAmount / int64(req.TotalPeriods)
	if periodAmount <= 0 {
		return nil, fmt.Errorf("每期金额必须大于0")
	}

	// 生成计划编号
	planNo := fmt.Sprintf("DP%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	plan := &models.DeductionPlan{
		PlanNo:          planNo,
		DeductorID:      req.DeductorID,
		DeducteeID:      req.DeducteeID,
		PlanType:        req.PlanType,
		TotalAmount:     req.TotalAmount,
		DeductedAmount:  0,
		RemainingAmount: req.TotalAmount,
		TotalPeriods:    req.TotalPeriods,
		CurrentPeriod:   0,
		PeriodAmount:    periodAmount,
		Status:          models.DeductionPlanStatusActive,
		RelatedType:     req.RelatedType,
		RelatedID:       req.RelatedID,
		Remark:          req.Remark,
		CreatedBy:       req.CreatedBy,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, fmt.Errorf("创建代扣计划失败: %w", err)
	}

	// 生成代扣记录（按期数）
	if err := s.generateDeductionRecords(plan); err != nil {
		return nil, fmt.Errorf("生成代扣记录失败: %w", err)
	}

	log.Printf("[DeductionService] Created deduction plan: %s, deductor: %d, deductee: %d, amount: %d",
		planNo, req.DeductorID, req.DeducteeID, req.TotalAmount)

	return plan, nil
}

// generateDeductionRecords 生成代扣记录
func (s *DeductionService) generateDeductionRecords(plan *models.DeductionPlan) error {
	records := make([]*models.DeductionRecord, 0, plan.TotalPeriods)

	// 从明天开始，每天一期
	baseTime := time.Now().AddDate(0, 0, 1)
	baseTime = time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), 8, 0, 0, 0, baseTime.Location()) // 每天8点扣款

	remainingAmount := plan.TotalAmount
	for i := 1; i <= plan.TotalPeriods; i++ {
		// 最后一期扣除剩余金额
		amount := plan.PeriodAmount
		if i == plan.TotalPeriods {
			amount = remainingAmount
		}
		remainingAmount -= amount

		record := &models.DeductionRecord{
			PlanID:      plan.ID,
			PlanNo:      plan.PlanNo,
			DeductorID:  plan.DeductorID,
			DeducteeID:  plan.DeducteeID,
			PeriodNum:   i,
			Amount:      amount,
			Status:      models.DeductionRecordStatusPending,
			ScheduledAt: baseTime.AddDate(0, 0, i-1), // 每天扣一期
			CreatedAt:   time.Now(),
		}
		records = append(records, record)
	}

	return s.recordRepo.BatchCreate(records)
}

// WalletInfo 钱包信息（用于多钱包扣款排序）
type WalletInfo struct {
	Wallet     *repository.Wallet
	ChannelID  int64
	WalletType int16
}

// ExecuteDailyDeduction 执行每日扣款
// 业务规则：每天固定时间检查余额并扣款
func (s *DeductionService) ExecuteDailyDeduction() error {
	log.Printf("[DeductionService] Starting daily deduction job...")

	// 获取今天待扣款的记录
	now := time.Now()
	pendingRecords, err := s.recordRepo.FindPendingRecords(now, 1000)
	if err != nil {
		return fmt.Errorf("查询待扣款记录失败: %w", err)
	}

	log.Printf("[DeductionService] Found %d pending deduction records", len(pendingRecords))

	successCount := 0
	failCount := 0

	for _, record := range pendingRecords {
		if err := s.executeDeduction(record); err != nil {
			log.Printf("[DeductionService] Deduction failed for record %d: %v", record.ID, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("[DeductionService] Daily deduction completed: success=%d, fail=%d", successCount, failCount)
	return nil
}

// executeDeduction 执行单条代扣
// 业务规则：多钱包扣款按余额从多到少扣
func (s *DeductionService) executeDeduction(record *models.DeductionRecord) error {
	// 获取代扣计划
	plan, err := s.planRepo.FindByID(record.PlanID)
	if err != nil || plan == nil {
		return fmt.Errorf("代扣计划不存在: %d", record.PlanID)
	}

	// 检查计划状态
	if plan.Status != models.DeductionPlanStatusActive {
		return fmt.Errorf("代扣计划状态异常: %d", plan.Status)
	}

	// 获取被扣款方的所有钱包
	wallets, err := s.getAgentWallets(record.DeducteeID)
	if err != nil {
		return fmt.Errorf("获取钱包失败: %w", err)
	}

	if len(wallets) == 0 {
		// 标记为失败
		s.recordRepo.UpdateStatus(record.ID, models.DeductionRecordStatusFailed, 0, "", "无可用钱包")
		return fmt.Errorf("无可用钱包")
	}

	// 按余额从多到少排序（Q8规则）
	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].Balance > wallets[j].Balance
	})

	// 执行多钱包扣款
	totalDeducted, walletDetails, err := s.deductFromMultipleWallets(wallets, record.Amount, record.ID)
	if err != nil {
		s.recordRepo.UpdateStatus(record.ID, models.DeductionRecordStatusFailed, totalDeducted, walletDetails, err.Error())
		return err
	}

	// 判断扣款结果
	var status int16
	var failReason string
	if totalDeducted >= record.Amount {
		status = models.DeductionRecordStatusSuccess
	} else if totalDeducted > 0 {
		status = models.DeductionRecordStatusPartialSuccess
		failReason = fmt.Sprintf("部分成功，应扣%d分，实扣%d分", record.Amount, totalDeducted)
	} else {
		status = models.DeductionRecordStatusFailed
		failReason = "余额不足"
	}

	// 更新代扣记录状态
	if err := s.recordRepo.UpdateStatus(record.ID, status, totalDeducted, walletDetails, failReason); err != nil {
		return fmt.Errorf("更新代扣记录状态失败: %w", err)
	}

	// 更新代扣计划进度
	if totalDeducted > 0 {
		if err := s.planRepo.UpdateDeductedAmount(plan.ID, totalDeducted, record.PeriodNum); err != nil {
			return fmt.Errorf("更新代扣计划进度失败: %w", err)
		}

		// 检查是否已完成
		if record.PeriodNum >= plan.TotalPeriods {
			s.planRepo.UpdateStatus(plan.ID, models.DeductionPlanStatusCompleted)
		}
	}

	// 扣款成功后，将金额转入扣款方钱包
	if totalDeducted > 0 {
		if err := s.transferToDeductor(record.DeductorID, totalDeducted, record.ID); err != nil {
			log.Printf("[DeductionService] Transfer to deductor failed: %v", err)
		}
	}

	log.Printf("[DeductionService] Deduction executed: record=%d, amount=%d, deducted=%d, status=%d",
		record.ID, record.Amount, totalDeducted, status)

	return nil
}

// getAgentWallets 获取代理商的所有钱包
func (s *DeductionService) getAgentWallets(agentID int64) ([]*repository.Wallet, error) {
	// 获取所有通道的所有类型钱包
	// 这里简化处理，实际应该从数据库查询所有钱包
	wallets := make([]*repository.Wallet, 0)

	// 钱包类型：1分润钱包 2服务费钱包 3奖励钱包
	walletTypes := []int16{1, 2, 3}
	// 假设通道ID为1（实际应该查询所有通道）
	channelIDs := []int64{1}

	for _, channelID := range channelIDs {
		for _, walletType := range walletTypes {
			wallet, err := s.walletRepo.FindByAgentAndType(agentID, channelID, walletType)
			if err == nil && wallet != nil && wallet.Balance > 0 {
				wallets = append(wallets, wallet)
			}
		}
	}

	return wallets, nil
}

// deductFromMultipleWallets 从多个钱包扣款（按余额从多到少）
func (s *DeductionService) deductFromMultipleWallets(wallets []*repository.Wallet, targetAmount int64, recordID int64) (int64, string, error) {
	var totalDeducted int64
	walletDetails := make([]models.WalletDeductDetail, 0)
	remainingAmount := targetAmount

	for _, wallet := range wallets {
		if remainingAmount <= 0 {
			break
		}

		// 计算本钱包可扣金额
		deductAmount := wallet.Balance
		if deductAmount > remainingAmount {
			deductAmount = remainingAmount
		}

		if deductAmount <= 0 {
			continue
		}

		// 扣减钱包余额
		if err := s.walletRepo.UpdateBalance(wallet.ID, -deductAmount); err != nil {
			log.Printf("[DeductionService] Deduct from wallet %d failed: %v", wallet.ID, err)
			continue
		}

		// 记录钱包流水
		walletLog := &repository.WalletLog{
			WalletID:      wallet.ID,
			AgentID:       wallet.AgentID,
			WalletType:    wallet.WalletType,
			LogType:       6, // 代扣
			Amount:        -deductAmount,
			BalanceBefore: wallet.Balance,
			BalanceAfter:  wallet.Balance - deductAmount,
			RefType:       "deduction_record",
			RefID:         recordID,
			Remark:        "代扣扣款",
			CreatedAt:     time.Now(),
		}
		s.walletLogRepo.Create(walletLog)

		// 记录扣款明细
		detail := models.WalletDeductDetail{
			WalletID:      wallet.ID,
			WalletType:    wallet.WalletType,
			WalletName:    getWalletTypeName(wallet.WalletType),
			BalanceBefore: wallet.Balance,
			DeductAmount:  deductAmount,
			BalanceAfter:  wallet.Balance - deductAmount,
		}
		walletDetails = append(walletDetails, detail)

		totalDeducted += deductAmount
		remainingAmount -= deductAmount
	}

	// 转换为JSON
	detailsJSON, _ := json.Marshal(walletDetails)

	return totalDeducted, string(detailsJSON), nil
}

// transferToDeductor 将扣款金额转入扣款方钱包
func (s *DeductionService) transferToDeductor(deductorID int64, amount int64, recordID int64) error {
	// 转入扣款方的分润钱包（默认通道1，钱包类型1）
	wallet, err := s.walletRepo.FindByAgentAndType(deductorID, 1, 1)
	if err != nil {
		return fmt.Errorf("扣款方钱包不存在: %w", err)
	}

	// 增加余额
	if err := s.walletRepo.UpdateBalance(wallet.ID, amount); err != nil {
		return fmt.Errorf("增加扣款方余额失败: %w", err)
	}

	// 记录钱包流水
	walletLog := &repository.WalletLog{
		WalletID:      wallet.ID,
		AgentID:       deductorID,
		WalletType:    wallet.WalletType,
		LogType:       6, // 代扣
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance + amount,
		RefType:       "deduction_record",
		RefID:         recordID,
		Remark:        "代扣收款",
		CreatedAt:     time.Now(),
	}
	return s.walletLogRepo.Create(walletLog)
}

// getWalletTypeName 获取钱包类型名称
func getWalletTypeName(walletType int16) string {
	switch walletType {
	case 1:
		return "分润钱包"
	case 2:
		return "服务费钱包"
	case 3:
		return "奖励钱包"
	default:
		return "未知钱包"
	}
}

// PauseDeductionPlan 暂停代扣计划
func (s *DeductionService) PauseDeductionPlan(planID int64) error {
	plan, err := s.planRepo.FindByID(planID)
	if err != nil || plan == nil {
		return fmt.Errorf("代扣计划不存在: %d", planID)
	}

	if plan.Status != models.DeductionPlanStatusActive {
		return fmt.Errorf("计划状态不允许暂停")
	}

	return s.planRepo.UpdateStatus(planID, models.DeductionPlanStatusPaused)
}

// ResumeDeductionPlan 恢复代扣计划
func (s *DeductionService) ResumeDeductionPlan(planID int64) error {
	plan, err := s.planRepo.FindByID(planID)
	if err != nil || plan == nil {
		return fmt.Errorf("代扣计划不存在: %d", planID)
	}

	if plan.Status != models.DeductionPlanStatusPaused {
		return fmt.Errorf("计划状态不允许恢复")
	}

	return s.planRepo.UpdateStatus(planID, models.DeductionPlanStatusActive)
}

// CancelDeductionPlan 取消代扣计划
func (s *DeductionService) CancelDeductionPlan(planID int64) error {
	plan, err := s.planRepo.FindByID(planID)
	if err != nil || plan == nil {
		return fmt.Errorf("代扣计划不存在: %d", planID)
	}

	if plan.Status == models.DeductionPlanStatusCompleted {
		return fmt.Errorf("已完成的计划不能取消")
	}

	return s.planRepo.UpdateStatus(planID, models.DeductionPlanStatusCancelled)
}

// CreateDeductionChainRequest 创建代扣链请求
type CreateDeductionChainRequest struct {
	DistributeID int64   `json:"distribute_id"` // 终端下发记录ID
	TerminalSN   string  `json:"terminal_sn"`   // 终端SN
	AgentPath    []int64 `json:"agent_path"`    // 代理商路径 [A, B, C]（从上到下）
	TotalAmount  int64   `json:"total_amount"`  // 总金额
	TotalPeriods int     `json:"total_periods"` // 总期数
	CreatedBy    int64   `json:"created_by"`    // 创建人
}

// CreateDeductionChain 创建代扣链（用于跨级下发）
// 业务规则Q16：跨级下发时系统自动按层级生成A→B→C的货款代扣链
func (s *DeductionService) CreateDeductionChain(req *CreateDeductionChainRequest) (*models.DeductionChain, error) {
	if len(req.AgentPath) < 2 {
		return nil, fmt.Errorf("代扣链至少需要2个代理商")
	}

	// 生成代扣链编号
	chainNo := fmt.Sprintf("DC%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	chain := &models.DeductionChain{
		ChainNo:      chainNo,
		DistributeID: req.DistributeID,
		TerminalSN:   req.TerminalSN,
		TotalLevels:  len(req.AgentPath) - 1,
		TotalAmount:  req.TotalAmount,
		Status:       1, // 进行中
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.chainRepo.Create(chain); err != nil {
		return nil, fmt.Errorf("创建代扣链失败: %w", err)
	}

	// 生成代扣链节点（A→B, B→C, ...）
	items := make([]*models.DeductionChainItem, 0)
	for i := 0; i < len(req.AgentPath)-1; i++ {
		item := &models.DeductionChainItem{
			ChainID:     chain.ID,
			ChainNo:     chainNo,
			Level:       i + 1,
			FromAgentID: req.AgentPath[i+1], // 下级代理商扣款
			ToAgentID:   req.AgentPath[i],   // 上级代理商收款
			Amount:      req.TotalAmount,
			Status:      0, // 待处理
		}
		items = append(items, item)
	}

	if err := s.chainItemRepo.BatchCreate(items); err != nil {
		return nil, fmt.Errorf("创建代扣链节点失败: %w", err)
	}

	// 为每个节点创建代扣计划
	for _, item := range items {
		planReq := &CreateDeductionPlanRequest{
			DeductorID:   item.ToAgentID,
			DeducteeID:   item.FromAgentID,
			PlanType:     models.DeductionPlanTypeGoods,
			TotalAmount:  item.Amount,
			TotalPeriods: req.TotalPeriods,
			RelatedType:  "deduction_chain_item",
			RelatedID:    item.ID,
			Remark:       fmt.Sprintf("跨级下发货款代扣 - 终端%s", req.TerminalSN),
			CreatedBy:    req.CreatedBy,
		}

		plan, err := s.CreateDeductionPlan(planReq)
		if err != nil {
			log.Printf("[DeductionService] Create deduction plan for chain item %d failed: %v", item.ID, err)
			continue
		}

		// 更新节点关联的计划ID
		s.chainItemRepo.UpdatePlanID(item.ID, plan.ID)
		s.chainItemRepo.UpdateStatus(item.ID, 1) // 已生成计划
	}

	log.Printf("[DeductionService] Created deduction chain: %s, levels: %d, amount: %d",
		chainNo, chain.TotalLevels, req.TotalAmount)

	return chain, nil
}

// GetAgentPathBetween 获取两个代理商之间的路径
func (s *DeductionService) GetAgentPathBetween(fromAgentID, toAgentID int64) ([]int64, error) {
	// 获取目标代理商（下级）
	toAgent, err := s.agentRepo.FindByID(toAgentID)
	if err != nil || toAgent == nil {
		return nil, fmt.Errorf("目标代理商不存在: %d", toAgentID)
	}

	// 解析物化路径
	// 路径格式: /1/5/12/
	if toAgent.Path == "" {
		return nil, fmt.Errorf("代理商路径为空")
	}

	pathStr := strings.Trim(toAgent.Path, "/")
	if pathStr == "" {
		return nil, fmt.Errorf("代理商路径为空")
	}

	// 分割路径
	pathParts := strings.Split(pathStr, "/")
	agentPath := make([]int64, 0)

	// 找到fromAgentID的位置
	foundFrom := false
	for _, part := range pathParts {
		var id int64
		fmt.Sscanf(part, "%d", &id)

		if id == fromAgentID {
			foundFrom = true
		}

		if foundFrom {
			agentPath = append(agentPath, id)
		}
	}

	// 添加目标代理商
	agentPath = append(agentPath, toAgentID)

	if !foundFrom {
		return nil, fmt.Errorf("下发方不在目标代理商的上级链中")
	}

	return agentPath, nil
}
