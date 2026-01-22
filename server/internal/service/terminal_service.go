package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// TerminalService 终端服务
type TerminalService struct {
	terminalRepo      *repository.GormTerminalRepository
	recallRepo        *repository.GormTerminalRecallRepository
	importRecordRepo  *repository.GormTerminalImportRecordRepository
	agentRepo         repository.AgentRepository
}

// NewTerminalService 创建终端服务
func NewTerminalService(
	terminalRepo *repository.GormTerminalRepository,
	recallRepo *repository.GormTerminalRecallRepository,
	importRecordRepo *repository.GormTerminalImportRecordRepository,
	agentRepo repository.AgentRepository,
) *TerminalService {
	return &TerminalService{
		terminalRepo:     terminalRepo,
		recallRepo:       recallRepo,
		importRecordRepo: importRecordRepo,
		agentRepo:        agentRepo,
	}
}

// ImportTerminalsRequest 终端入库请求
type ImportTerminalsRequest struct {
	ChannelID    int64    `json:"channel_id"`    // 通道ID
	ChannelCode  string   `json:"channel_code"`  // 通道编码
	BrandCode    string   `json:"brand_code"`    // 品牌编码
	ModelCode    string   `json:"model_code"`    // 型号编码
	SNList       []string `json:"sn_list"`       // SN列表
	OwnerAgentID int64    `json:"owner_agent_id"` // 入库代理商
	CreatedBy    int64    `json:"created_by"`    // 创建人
}

// ImportTerminalsResult 终端入库结果
type ImportTerminalsResult struct {
	ImportNo     string   `json:"import_no"`     // 批次号
	TotalCount   int      `json:"total_count"`   // 总数
	SuccessCount int      `json:"success_count"` // 成功数
	FailedCount  int      `json:"failed_count"`  // 失败数
	FailedSNs    []string `json:"failed_sns"`    // 失败的SN列表
	Errors       []string `json:"errors"`        // 错误信息
}

// ImportTerminals 批量入库终端
func (s *TerminalService) ImportTerminals(req *ImportTerminalsRequest) (*ImportTerminalsResult, error) {
	if len(req.SNList) == 0 {
		return nil, fmt.Errorf("SN列表不能为空")
	}

	if len(req.SNList) > 1000 {
		return nil, fmt.Errorf("单次入库不能超过1000台")
	}

	// 生成入库批次号
	importNo := fmt.Sprintf("IMP%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	// 去重和清理SN
	snSet := make(map[string]bool)
	cleanedSNs := make([]string, 0, len(req.SNList))
	for _, sn := range req.SNList {
		sn = strings.TrimSpace(sn)
		if sn != "" && !snSet[sn] {
			snSet[sn] = true
			cleanedSNs = append(cleanedSNs, sn)
		}
	}

	// 检查已存在的终端
	existingTerminals, err := s.terminalRepo.FindBySNs(cleanedSNs)
	if err != nil {
		return nil, fmt.Errorf("查询已存在终端失败: %w", err)
	}

	existingSNs := make(map[string]bool)
	for _, t := range existingTerminals {
		existingSNs[t.TerminalSN] = true
	}

	// 分离新SN和已存在SN
	var newSNs, failedSNs []string
	var errors []string
	for _, sn := range cleanedSNs {
		if existingSNs[sn] {
			failedSNs = append(failedSNs, sn)
			errors = append(errors, fmt.Sprintf("终端 %s 已存在", sn))
		} else {
			newSNs = append(newSNs, sn)
		}
	}

	// 批量创建终端
	now := time.Now()
	terminals := make([]*models.Terminal, 0, len(newSNs))
	for _, sn := range newSNs {
		terminals = append(terminals, &models.Terminal{
			TerminalSN:   sn,
			ChannelID:    req.ChannelID,
			ChannelCode:  req.ChannelCode,
			BrandCode:    req.BrandCode,
			ModelCode:    req.ModelCode,
			OwnerAgentID: req.OwnerAgentID,
			Status:       models.TerminalStatusPending,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	if len(terminals) > 0 {
		if err := s.terminalRepo.BatchCreate(terminals); err != nil {
			return nil, fmt.Errorf("批量创建终端失败: %w", err)
		}
	}

	// 记录入库批次
	failedSNsJSON, _ := json.Marshal(failedSNs)
	importRecord := &models.TerminalImportRecord{
		ImportNo:     importNo,
		ChannelID:    req.ChannelID,
		ChannelCode:  req.ChannelCode,
		BrandCode:    req.BrandCode,
		ModelCode:    req.ModelCode,
		TotalCount:   len(cleanedSNs),
		SuccessCount: len(newSNs),
		FailedCount:  len(failedSNs),
		FailedSNs:    string(failedSNsJSON),
		OwnerAgentID: req.OwnerAgentID,
		CreatedBy:    req.CreatedBy,
		CreatedAt:    now,
	}

	if err := s.importRecordRepo.Create(importRecord); err != nil {
		log.Printf("[TerminalService] Create import record failed: %v", err)
	}

	log.Printf("[TerminalService] Import terminals: batch=%s, total=%d, success=%d, failed=%d",
		importNo, len(cleanedSNs), len(newSNs), len(failedSNs))

	return &ImportTerminalsResult{
		ImportNo:     importNo,
		TotalCount:   len(cleanedSNs),
		SuccessCount: len(newSNs),
		FailedCount:  len(failedSNs),
		FailedSNs:    failedSNs,
		Errors:       errors,
	}, nil
}

// RecallTerminalRequest 终端回拨请求
type RecallTerminalRequest struct {
	FromAgentID int64  `json:"from_agent_id"` // 回拨方代理商ID（当前持有者）
	ToAgentID   int64  `json:"to_agent_id"`   // 接收方代理商ID（上级）
	TerminalSN  string `json:"terminal_sn"`   // 终端SN
	ChannelID   int64  `json:"channel_id"`    // 通道ID
	Source      int16  `json:"source"`        // 1:APP 2:PC
	Remark      string `json:"remark"`        // 备注
	CreatedBy   int64  `json:"created_by"`    // 创建人
}

// RecallTerminal 终端回拨
func (s *TerminalService) RecallTerminal(req *RecallTerminalRequest) (*models.TerminalRecall, error) {
	// 1. 验证终端是否存在且属于回拨方
	terminal, err := s.terminalRepo.FindBySN(req.TerminalSN)
	if err != nil || terminal == nil {
		return nil, fmt.Errorf("终端不存在: %s", req.TerminalSN)
	}

	if terminal.OwnerAgentID != req.FromAgentID {
		return nil, fmt.Errorf("终端不属于当前代理商")
	}

	// 只有未激活的终端才能回拨
	if terminal.Status == models.TerminalStatusActivated {
		return nil, fmt.Errorf("已激活的终端不能回拨")
	}

	// 2. 验证回拨方和接收方
	fromAgent, err := s.agentRepo.FindByID(req.FromAgentID)
	if err != nil || fromAgent == nil {
		return nil, fmt.Errorf("回拨方代理商不存在")
	}

	toAgent, err := s.agentRepo.FindByID(req.ToAgentID)
	if err != nil || toAgent == nil {
		return nil, fmt.Errorf("接收方代理商不存在")
	}

	// 3. 检查是否跨级回拨
	isCrossLevel, crossLevelPath, err := s.checkCrossLevelRecall(fromAgent, toAgent)
	if err != nil {
		return nil, fmt.Errorf("检查跨级关系失败: %w", err)
	}

	// 4. APP端不允许跨级回拨
	if req.Source == models.TerminalDistributeSourceApp && isCrossLevel {
		return nil, fmt.Errorf("APP端不支持跨级回拨，请使用PC端")
	}

	// 5. 生成回拨单号
	recallNo := fmt.Sprintf("TR%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	// 6. 创建回拨记录
	recall := &models.TerminalRecall{
		RecallNo:       recallNo,
		FromAgentID:    req.FromAgentID,
		ToAgentID:      req.ToAgentID,
		TerminalSN:     req.TerminalSN,
		ChannelID:      req.ChannelID,
		IsCrossLevel:   isCrossLevel,
		CrossLevelPath: crossLevelPath,
		Status:         models.TerminalRecallStatusPending,
		Source:         req.Source,
		Remark:         req.Remark,
		CreatedBy:      req.CreatedBy,
		CreatedAt:      time.Now(),
	}

	if err := s.recallRepo.Create(recall); err != nil {
		return nil, fmt.Errorf("创建回拨记录失败: %w", err)
	}

	log.Printf("[TerminalService] Created recall: %s, from: %d, to: %d, crossLevel: %v",
		recallNo, req.FromAgentID, req.ToAgentID, isCrossLevel)

	return recall, nil
}

// checkCrossLevelRecall 检查是否跨级回拨
func (s *TerminalService) checkCrossLevelRecall(fromAgent, toAgent *repository.Agent) (bool, string, error) {
	// 回拨是从下级回到上级
	// 检查toAgent是否是fromAgent的直属上级
	if fromAgent.ParentID == toAgent.ID {
		return false, "", nil // 非跨级
	}

	// 检查toAgent是否在fromAgent的上级链中
	// fromAgent.Path应该包含toAgent.ID
	if fromAgent.Path == "" {
		return false, "", fmt.Errorf("回拨方代理商路径为空")
	}

	toIDStr := fmt.Sprintf("/%d/", toAgent.ID)
	if !strings.Contains(fromAgent.Path, toIDStr) {
		return false, "", fmt.Errorf("接收方不是回拨方的上级")
	}

	// 是跨级回拨，提取路径
	// 从toAgent.ID到fromAgent.ID的路径
	idx := strings.Index(fromAgent.Path, toIDStr)
	if idx == -1 {
		return false, "", fmt.Errorf("无法解析跨级路径")
	}

	crossPath := fromAgent.Path[idx:]
	return true, crossPath, nil
}

// ConfirmRecall 确认接收终端回拨
func (s *TerminalService) ConfirmRecall(recallID int64, confirmedBy int64) error {
	// 1. 获取回拨记录
	recall, err := s.recallRepo.FindByID(recallID)
	if err != nil || recall == nil {
		return fmt.Errorf("回拨记录不存在")
	}

	if recall.Status != models.TerminalRecallStatusPending {
		return fmt.Errorf("回拨记录状态不允许确认")
	}

	// 2. 更新终端所有权
	terminal, err := s.terminalRepo.FindBySN(recall.TerminalSN)
	if err != nil || terminal == nil {
		return fmt.Errorf("终端不存在")
	}

	if err := s.terminalRepo.UpdateOwner(terminal.ID, recall.ToAgentID); err != nil {
		return fmt.Errorf("更新终端所有权失败: %w", err)
	}

	// 3. 更新回拨记录状态
	if err := s.recallRepo.UpdateStatus(recallID, models.TerminalRecallStatusConfirmed, &confirmedBy); err != nil {
		return fmt.Errorf("更新回拨记录状态失败: %w", err)
	}

	log.Printf("[TerminalService] Confirmed recall: %d", recallID)
	return nil
}

// RejectRecall 拒绝终端回拨
func (s *TerminalService) RejectRecall(recallID int64, confirmedBy int64) error {
	recall, err := s.recallRepo.FindByID(recallID)
	if err != nil || recall == nil {
		return fmt.Errorf("回拨记录不存在")
	}

	if recall.Status != models.TerminalRecallStatusPending {
		return fmt.Errorf("回拨记录状态不允许拒绝")
	}

	return s.recallRepo.UpdateStatus(recallID, models.TerminalRecallStatusRejected, &confirmedBy)
}

// CancelRecall 取消终端回拨（回拨方取消）
func (s *TerminalService) CancelRecall(recallID int64, cancelBy int64) error {
	recall, err := s.recallRepo.FindByID(recallID)
	if err != nil || recall == nil {
		return fmt.Errorf("回拨记录不存在")
	}

	if recall.Status != models.TerminalRecallStatusPending {
		return fmt.Errorf("回拨记录状态不允许取消")
	}

	if recall.FromAgentID != cancelBy {
		return fmt.Errorf("只有回拨方可以取消")
	}

	return s.recallRepo.UpdateStatus(recallID, models.TerminalRecallStatusCancelled, nil)
}

// GetRecallList 获取回拨列表
func (s *TerminalService) GetRecallList(agentID int64, direction string, status []int16, limit, offset int) ([]*models.TerminalRecall, int64, error) {
	if direction == "from" {
		return s.recallRepo.FindByFromAgent(agentID, status, limit, offset)
	}
	return s.recallRepo.FindByToAgent(agentID, status, limit, offset)
}

// TerminalStats 终端统计
type TerminalStats struct {
	Total              int64 `json:"total"`               // 总数
	PendingCount       int64 `json:"pending_count"`       // 待分配
	AllocatedCount     int64 `json:"allocated_count"`     // 已分配
	BoundCount         int64 `json:"bound_count"`         // 已绑定
	ActivatedCount     int64 `json:"activated_count"`     // 已激活
	UnboundCount       int64 `json:"unbound_count"`       // 未绑定 = 待分配 + 已分配
	YesterdayActivated int64 `json:"yesterday_activated"` // 昨日激活
	TodayActivated     int64 `json:"today_activated"`     // 今日激活
	MonthActivated     int64 `json:"month_activated"`     // 本月激活
}

// GetTerminalStats 获取终端统计
func (s *TerminalService) GetTerminalStats(agentID int64) (*TerminalStats, error) {
	stats := &TerminalStats{}

	// 各状态统计
	s.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusPending, &stats.PendingCount)
	s.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusAllocated, &stats.AllocatedCount)
	s.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusBound, &stats.BoundCount)
	s.terminalRepo.CountByOwnerAndStatus(agentID, models.TerminalStatusActivated, &stats.ActivatedCount)

	stats.Total = stats.PendingCount + stats.AllocatedCount + stats.BoundCount + stats.ActivatedCount
	stats.UnboundCount = stats.PendingCount + stats.AllocatedCount

	// 时间统计
	now := time.Now()
	loc := now.Location()

	// 今日开始时间
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	// 昨日时间范围
	yesterdayStart := todayStart.Add(-24 * time.Hour)
	yesterdayEnd := todayStart

	// 本月开始时间
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	// 统计激活数
	stats.TodayActivated, _ = s.terminalRepo.CountActivatedByDate(agentID, todayStart, todayEnd)
	stats.YesterdayActivated, _ = s.terminalRepo.CountActivatedByDate(agentID, yesterdayStart, yesterdayEnd)
	stats.MonthActivated, _ = s.terminalRepo.CountActivatedByDate(agentID, monthStart, todayEnd)

	return stats, nil
}

// BatchRecallRequest 批量回拨请求
type BatchRecallRequest struct {
	TerminalSNs []string `json:"terminal_sns"` // 终端SN列表
	ToAgentID   int64    `json:"to_agent_id"`  // 接收方代理商ID
	FromAgentID int64    `json:"from_agent_id"`// 回拨方代理商ID
	Source      int16    `json:"source"`       // 来源
	Remark      string   `json:"remark"`       // 备注
	CreatedBy   int64    `json:"created_by"`   // 创建人
}

// BatchRecallResult 批量回拨结果
type BatchRecallResult struct {
	TotalCount   int      `json:"total_count"`   // 总数
	SuccessCount int      `json:"success_count"` // 成功数
	FailedCount  int      `json:"failed_count"`  // 失败数
	Errors       []string `json:"errors"`        // 错误信息
}

// BatchRecallTerminals 批量回拨终端
func (s *TerminalService) BatchRecallTerminals(req *BatchRecallRequest) (*BatchRecallResult, error) {
	if len(req.TerminalSNs) == 0 {
		return nil, fmt.Errorf("终端列表不能为空")
	}

	result := &BatchRecallResult{
		TotalCount: len(req.TerminalSNs),
		Errors:     make([]string, 0),
	}

	for _, sn := range req.TerminalSNs {
		recallReq := &RecallTerminalRequest{
			FromAgentID: req.FromAgentID,
			ToAgentID:   req.ToAgentID,
			TerminalSN:  sn,
			Source:      req.Source,
			Remark:      req.Remark,
			CreatedBy:   req.CreatedBy,
		}

		_, err := s.RecallTerminal(recallReq)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: %v", sn, err))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// ========== 终端政策设置 ==========

// BatchSetRateRequest 批量设置费率请求
type BatchSetRateRequest struct {
	TerminalSNs  []string `json:"terminal_sns"`
	AgentID      int64    `json:"agent_id"`
	CreditRate   int      `json:"credit_rate"`
	DebitRate    int      `json:"debit_rate"`
	DebitCap     int      `json:"debit_cap"`
	UnionpayRate int      `json:"unionpay_rate"`
	WechatRate   int      `json:"wechat_rate"`
	AlipayRate   int      `json:"alipay_rate"`
	UpdatedBy    int64    `json:"updated_by"`
}

// BatchSetSimFeeRequest 批量设置SIM卡费用请求
type BatchSetSimFeeRequest struct {
	TerminalSNs        []string `json:"terminal_sns"`
	AgentID            int64    `json:"agent_id"`
	FirstSimFee        int      `json:"first_sim_fee"`
	NonFirstSimFee     int      `json:"non_first_sim_fee"`
	SimFeeIntervalDays int      `json:"sim_fee_interval_days"`
	UpdatedBy          int64    `json:"updated_by"`
}

// BatchSetDepositRequest 批量设置押金请求
type BatchSetDepositRequest struct {
	TerminalSNs   []string `json:"terminal_sns"`
	AgentID       int64    `json:"agent_id"`
	DepositAmount int      `json:"deposit_amount"`
	UpdatedBy     int64    `json:"updated_by"`
}

// BatchPolicyResult 批量设置政策结果
type BatchPolicyResult struct {
	TotalCount   int      `json:"total_count"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors"`
}

// GetTerminalPolicy 获取终端政策
func (s *TerminalService) GetTerminalPolicy(sn string, agentID int64) (*models.TerminalPolicy, error) {
	// 验证终端存在且属于该代理商
	terminal, err := s.terminalRepo.FindBySN(sn)
	if err != nil || terminal == nil {
		return nil, fmt.Errorf("终端不存在: %s", sn)
	}

	if terminal.OwnerAgentID != agentID {
		return nil, fmt.Errorf("终端不属于当前代理商")
	}

	// 查询政策
	policy, err := s.terminalRepo.FindPolicyBySN(sn)
	if err != nil {
		// 如果不存在，返回默认值
		return &models.TerminalPolicy{
			TerminalSN: sn,
			ChannelID:  terminal.ChannelID,
			AgentID:    agentID,
		}, nil
	}

	return policy, nil
}

// BatchSetRate 批量设置费率
func (s *TerminalService) BatchSetRate(req *BatchSetRateRequest) (*BatchPolicyResult, error) {
	if len(req.TerminalSNs) == 0 {
		return nil, fmt.Errorf("终端列表不能为空")
	}

	result := &BatchPolicyResult{
		TotalCount: len(req.TerminalSNs),
		Errors:     make([]string, 0),
	}

	now := time.Now()

	for _, sn := range req.TerminalSNs {
		// 验证终端存在且属于该代理商
		terminal, err := s.terminalRepo.FindBySN(sn)
		if err != nil || terminal == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 不存在", sn))
			continue
		}

		if terminal.OwnerAgentID != req.AgentID {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 不属于当前代理商", sn))
			continue
		}

		// 查询或创建政策
		policy, _ := s.terminalRepo.FindPolicyBySN(sn)
		if policy == nil {
			policy = &models.TerminalPolicy{
				TerminalSN: sn,
				ChannelID:  terminal.ChannelID,
				AgentID:    req.AgentID,
				CreatedBy:  req.UpdatedBy,
				CreatedAt:  now,
			}
		}

		// 更新费率
		policy.CreditRate = req.CreditRate
		policy.DebitRate = req.DebitRate
		policy.DebitCap = req.DebitCap
		policy.UnionpayRate = req.UnionpayRate
		policy.WechatRate = req.WechatRate
		policy.AlipayRate = req.AlipayRate
		policy.UpdatedBy = req.UpdatedBy
		policy.UpdatedAt = now
		policy.IsSynced = false // 标记需要同步

		if err := s.terminalRepo.SavePolicy(policy); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 保存失败 - %v", sn, err))
			continue
		}

		result.SuccessCount++
	}

	log.Printf("[TerminalService] BatchSetRate: total=%d, success=%d, failed=%d, credit_rate=%d",
		result.TotalCount, result.SuccessCount, result.FailedCount, req.CreditRate)

	return result, nil
}

// BatchSetSimFee 批量设置SIM卡费用
func (s *TerminalService) BatchSetSimFee(req *BatchSetSimFeeRequest) (*BatchPolicyResult, error) {
	if len(req.TerminalSNs) == 0 {
		return nil, fmt.Errorf("终端列表不能为空")
	}

	result := &BatchPolicyResult{
		TotalCount: len(req.TerminalSNs),
		Errors:     make([]string, 0),
	}

	now := time.Now()

	for _, sn := range req.TerminalSNs {
		// 验证终端存在且属于该代理商
		terminal, err := s.terminalRepo.FindBySN(sn)
		if err != nil || terminal == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 不存在", sn))
			continue
		}

		if terminal.OwnerAgentID != req.AgentID {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 不属于当前代理商", sn))
			continue
		}

		// 查询或创建政策
		policy, _ := s.terminalRepo.FindPolicyBySN(sn)
		if policy == nil {
			policy = &models.TerminalPolicy{
				TerminalSN: sn,
				ChannelID:  terminal.ChannelID,
				AgentID:    req.AgentID,
				CreatedBy:  req.UpdatedBy,
				CreatedAt:  now,
			}
		}

		// 更新SIM卡费用
		policy.FirstSimFee = req.FirstSimFee
		policy.NonFirstSimFee = req.NonFirstSimFee
		policy.SimFeeIntervalDays = req.SimFeeIntervalDays
		policy.UpdatedBy = req.UpdatedBy
		policy.UpdatedAt = now
		policy.IsSynced = false

		if err := s.terminalRepo.SavePolicy(policy); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 保存失败 - %v", sn, err))
			continue
		}

		result.SuccessCount++
	}

	log.Printf("[TerminalService] BatchSetSimFee: total=%d, success=%d, failed=%d",
		result.TotalCount, result.SuccessCount, result.FailedCount)

	return result, nil
}

// BatchSetDeposit 批量设置押金
func (s *TerminalService) BatchSetDeposit(req *BatchSetDepositRequest) (*BatchPolicyResult, error) {
	if len(req.TerminalSNs) == 0 {
		return nil, fmt.Errorf("终端列表不能为空")
	}

	result := &BatchPolicyResult{
		TotalCount: len(req.TerminalSNs),
		Errors:     make([]string, 0),
	}

	now := time.Now()

	for _, sn := range req.TerminalSNs {
		// 验证终端存在且属于该代理商
		terminal, err := s.terminalRepo.FindBySN(sn)
		if err != nil || terminal == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 不存在", sn))
			continue
		}

		if terminal.OwnerAgentID != req.AgentID {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 不属于当前代理商", sn))
			continue
		}

		// 查询或创建政策
		policy, _ := s.terminalRepo.FindPolicyBySN(sn)
		if policy == nil {
			policy = &models.TerminalPolicy{
				TerminalSN: sn,
				ChannelID:  terminal.ChannelID,
				AgentID:    req.AgentID,
				CreatedBy:  req.UpdatedBy,
				CreatedAt:  now,
			}
		}

		// 更新押金
		policy.DepositAmount = req.DepositAmount
		policy.UpdatedBy = req.UpdatedBy
		policy.UpdatedAt = now
		policy.IsSynced = false

		if err := s.terminalRepo.SavePolicy(policy); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("终端 %s: 保存失败 - %v", sn, err))
			continue
		}

		result.SuccessCount++
	}

	log.Printf("[TerminalService] BatchSetDeposit: total=%d, success=%d, failed=%d, amount=%d",
		result.TotalCount, result.SuccessCount, result.FailedCount, req.DepositAmount)

	return result, nil
}
