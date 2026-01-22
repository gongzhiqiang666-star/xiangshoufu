package service

import (
	"fmt"
	"log"
	"strings"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// TerminalDistributeService 终端下发服务
// 业务规则：
// - Q16: 跨级下发时系统自动按层级生成A→B→C的货款代扣链
// - Q29: APP不能跨级，PC可以跨级（保留整个层级关系）
type TerminalDistributeService struct {
	terminalRepo          repository.TerminalRepository
	distributeRepo        repository.TerminalDistributeRepository
	agentRepo             repository.AgentRepository
	deductionService      *DeductionService
	goodsDeductionService *GoodsDeductionService // 货款代扣服务
}

// NewTerminalDistributeService 创建终端下发服务
func NewTerminalDistributeService(
	terminalRepo repository.TerminalRepository,
	distributeRepo repository.TerminalDistributeRepository,
	agentRepo repository.AgentRepository,
	deductionService *DeductionService,
) *TerminalDistributeService {
	return &TerminalDistributeService{
		terminalRepo:     terminalRepo,
		distributeRepo:   distributeRepo,
		agentRepo:        agentRepo,
		deductionService: deductionService,
	}
}

// SetGoodsDeductionService 设置货款代扣服务（延迟注入，避免循环依赖）
func (s *TerminalDistributeService) SetGoodsDeductionService(gds *GoodsDeductionService) {
	s.goodsDeductionService = gds
}

// DistributeTerminalRequest 终端下发请求
type DistributeTerminalRequest struct {
	FromAgentID      int64  `json:"from_agent_id"`      // 下发方代理商ID
	ToAgentID        int64  `json:"to_agent_id"`        // 接收方代理商ID
	TerminalSN       string `json:"terminal_sn"`        // 终端SN
	ChannelID        int64  `json:"channel_id"`         // 通道ID
	GoodsPrice       int64  `json:"goods_price"`        // 货款金额（分）
	DeductionType    int16  `json:"deduction_type"`     // 1:一次性付款 2:分期代扣 3:货款代扣（实时）
	DeductionPeriods int    `json:"deduction_periods"`  // 分期期数（分期代扣时必填）
	DeductionSource  int16  `json:"deduction_source"`   // 货款代扣来源: 1=分润 2=服务费 3=两者（货款代扣时必填）
	Source           int16  `json:"source"`             // 1:APP 2:PC
	Remark           string `json:"remark"`             // 备注
	CreatedBy        int64  `json:"created_by"`         // 创建人
}

// 代扣类型常量
const (
	DistributeDeductionTypeOnce        int16 = 1 // 一次性付款
	DistributeDeductionTypeInstallment int16 = 2 // 分期代扣
	DistributeDeductionTypeRealtime    int16 = 3 // 货款代扣（实时扣款）
)

// DistributeTerminal 终端下发
func (s *TerminalDistributeService) DistributeTerminal(req *DistributeTerminalRequest) (*models.TerminalDistribute, error) {
	// 1. 验证终端是否存在且属于下发方
	terminal, err := s.terminalRepo.FindBySN(req.TerminalSN)
	if err != nil || terminal == nil {
		return nil, fmt.Errorf("终端不存在: %s", req.TerminalSN)
	}

	if terminal.OwnerAgentID != req.FromAgentID {
		return nil, fmt.Errorf("终端不属于当前代理商")
	}

	if terminal.Status != models.TerminalStatusPending && terminal.Status != models.TerminalStatusAllocated {
		return nil, fmt.Errorf("终端状态不允许下发")
	}

	// 2. 验证下发方和接收方
	fromAgent, err := s.agentRepo.FindByID(req.FromAgentID)
	if err != nil || fromAgent == nil {
		return nil, fmt.Errorf("下发方代理商不存在")
	}

	toAgent, err := s.agentRepo.FindByID(req.ToAgentID)
	if err != nil || toAgent == nil {
		return nil, fmt.Errorf("接收方代理商不存在")
	}

	// 3. 检查是否跨级下发
	isCrossLevel, crossLevelPath, err := s.checkCrossLevel(fromAgent, toAgent)
	if err != nil {
		return nil, fmt.Errorf("检查跨级关系失败: %w", err)
	}

	// 4. APP端不允许跨级下发（Q29规则）
	if req.Source == models.TerminalDistributeSourceApp && isCrossLevel {
		return nil, fmt.Errorf("APP端不支持跨级下发，请使用PC端")
	}

	// 5. 生成下发单号
	distributeNo := fmt.Sprintf("TD%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	// 6. 创建下发记录
	distribute := &models.TerminalDistribute{
		DistributeNo:   distributeNo,
		FromAgentID:    req.FromAgentID,
		ToAgentID:      req.ToAgentID,
		TerminalSN:     req.TerminalSN,
		ChannelID:      req.ChannelID,
		IsCrossLevel:   isCrossLevel,
		CrossLevelPath: crossLevelPath,
		GoodsPrice:     req.GoodsPrice,
		DeductionType:  req.DeductionType,
		Status:         models.TerminalDistributeStatusPending,
		Source:         req.Source,
		Remark:         req.Remark,
		CreatedBy:      req.CreatedBy,
		CreatedAt:      time.Now(),
	}

	if err := s.distributeRepo.Create(distribute); err != nil {
		return nil, fmt.Errorf("创建下发记录失败: %w", err)
	}

	log.Printf("[TerminalDistributeService] Created distribute: %s, from: %d, to: %d, crossLevel: %v",
		distributeNo, req.FromAgentID, req.ToAgentID, isCrossLevel)

	return distribute, nil
}

// checkCrossLevel 检查是否跨级下发
func (s *TerminalDistributeService) checkCrossLevel(fromAgent, toAgent *repository.Agent) (bool, string, error) {
	// 检查toAgent是否是fromAgent的直属下级
	if toAgent.ParentID == fromAgent.ID {
		return false, "", nil // 非跨级
	}

	// 检查toAgent是否在fromAgent的下级链中
	// toAgent.Path应该包含fromAgent.ID
	if toAgent.Path == "" {
		return false, "", fmt.Errorf("接收方代理商路径为空")
	}

	fromIDStr := fmt.Sprintf("/%d/", fromAgent.ID)
	if !strings.Contains(toAgent.Path, fromIDStr) {
		return false, "", fmt.Errorf("接收方不是下发方的下级")
	}

	// 是跨级下发，提取路径
	// 路径格式: /1/5/12/，toAgent.ID是12
	// 从fromAgent.ID开始的子路径
	idx := strings.Index(toAgent.Path, fromIDStr)
	if idx == -1 {
		return false, "", fmt.Errorf("无法解析跨级路径")
	}

	crossPath := toAgent.Path[idx:]
	return true, crossPath, nil
}

// ConfirmDistribute 确认接收终端下发
func (s *TerminalDistributeService) ConfirmDistribute(distributeID int64, confirmedBy int64) error {
	// 1. 获取下发记录
	distribute, err := s.distributeRepo.FindByID(distributeID)
	if err != nil || distribute == nil {
		return fmt.Errorf("下发记录不存在")
	}

	if distribute.Status != models.TerminalDistributeStatusPending {
		return fmt.Errorf("下发记录状态不允许确认")
	}

	// 2. 更新终端所有权
	terminal, err := s.terminalRepo.FindBySN(distribute.TerminalSN)
	if err != nil || terminal == nil {
		return fmt.Errorf("终端不存在")
	}

	if err := s.terminalRepo.UpdateOwner(terminal.ID, distribute.ToAgentID); err != nil {
		return fmt.Errorf("更新终端所有权失败: %w", err)
	}

	// 3. 处理货款代扣
	if distribute.GoodsPrice > 0 {
		switch distribute.DeductionType {
		case DistributeDeductionTypeRealtime:
			// 货款代扣（实时扣款）- 使用新的GoodsDeductionService
			if err := s.createGoodsDeductionForDistribute(distribute, terminal, confirmedBy); err != nil {
				log.Printf("[TerminalDistributeService] Create goods deduction failed: %v", err)
				// 货款代扣创建失败不阻塞终端下发
			}
		case DistributeDeductionTypeInstallment:
			// 分期代扣 - 使用原有的DeductionService
			if distribute.IsCrossLevel {
				// 跨级下发：生成代扣链（Q16规则）
				if err := s.createDeductionChainForCrossLevel(distribute, confirmedBy); err != nil {
					log.Printf("[TerminalDistributeService] Create deduction chain failed: %v", err)
				}
			} else {
				// 非跨级下发：直接生成代扣计划
				if err := s.createDeductionPlanForDistribute(distribute, confirmedBy); err != nil {
					log.Printf("[TerminalDistributeService] Create deduction plan failed: %v", err)
				}
			}
		}
	}

	// 4. 更新下发记录状态
	if err := s.distributeRepo.UpdateStatus(distributeID, models.TerminalDistributeStatusConfirmed, &confirmedBy); err != nil {
		return fmt.Errorf("更新下发记录状态失败: %w", err)
	}

	log.Printf("[TerminalDistributeService] Confirmed distribute: %d", distributeID)
	return nil
}

// createDeductionChainForCrossLevel 为跨级下发创建代扣链
func (s *TerminalDistributeService) createDeductionChainForCrossLevel(distribute *models.TerminalDistribute, createdBy int64) error {
	// 获取代理商路径
	agentPath, err := s.deductionService.GetAgentPathBetween(distribute.FromAgentID, distribute.ToAgentID)
	if err != nil {
		return fmt.Errorf("获取代理商路径失败: %w", err)
	}

	// 默认分12期
	periods := 12

	// 创建代扣链
	chainReq := &CreateDeductionChainRequest{
		DistributeID: distribute.ID,
		TerminalSN:   distribute.TerminalSN,
		AgentPath:    agentPath,
		TotalAmount:  distribute.GoodsPrice,
		TotalPeriods: periods,
		CreatedBy:    createdBy,
	}

	chain, err := s.deductionService.CreateDeductionChain(chainReq)
	if err != nil {
		return fmt.Errorf("创建代扣链失败: %w", err)
	}

	// 更新下发记录的代扣链ID
	distribute.ChainID = &chain.ID
	return s.distributeRepo.Update(distribute)
}

// createDeductionPlanForDistribute 为非跨级下发创建代扣计划
func (s *TerminalDistributeService) createDeductionPlanForDistribute(distribute *models.TerminalDistribute, createdBy int64) error {
	// 默认分12期
	periods := 12

	planReq := &CreateDeductionPlanRequest{
		DeductorID:   distribute.FromAgentID,
		DeducteeID:   distribute.ToAgentID,
		PlanType:     models.DeductionPlanTypeGoods,
		TotalAmount:  distribute.GoodsPrice,
		TotalPeriods: periods,
		RelatedType:  "terminal_distribute",
		RelatedID:    distribute.ID,
		Remark:       fmt.Sprintf("终端下发货款代扣 - %s", distribute.TerminalSN),
		CreatedBy:    createdBy,
	}

	plan, err := s.deductionService.CreateDeductionPlan(planReq)
	if err != nil {
		return fmt.Errorf("创建代扣计划失败: %w", err)
	}

	// 更新下发记录的代扣计划ID
	distribute.DeductionPlanID = &plan.ID
	return s.distributeRepo.Update(distribute)
}

// createGoodsDeductionForDistribute 为终端下发创建货款代扣（实时扣款）
func (s *TerminalDistributeService) createGoodsDeductionForDistribute(distribute *models.TerminalDistribute, terminal *models.Terminal, createdBy int64) error {
	if s.goodsDeductionService == nil {
		return fmt.Errorf("货款代扣服务未初始化")
	}

	// 获取扣款来源，默认两者都扣
	deductionSource := int16(3) // 默认: 分润+服务费
	if distribute.DeductionSource > 0 {
		deductionSource = distribute.DeductionSource
	}

	// 创建货款代扣请求
	req := &models.CreateGoodsDeductionRequest{
		ToAgentID:       distribute.ToAgentID,
		UnitPrice:       distribute.GoodsPrice,
		DeductionSource: deductionSource,
		Terminals: []models.CreateGoodsDeductionTerminal{
			{
				TerminalID: terminal.ID,
				TerminalSN: terminal.TerminalSN,
				UnitPrice:  distribute.GoodsPrice,
			},
		},
		Remark:       fmt.Sprintf("终端划拨货款代扣 - %s", distribute.TerminalSN),
		DistributeID: &distribute.ID,
	}

	deduction, err := s.goodsDeductionService.CreateGoodsDeduction(req, distribute.FromAgentID, createdBy)
	if err != nil {
		return fmt.Errorf("创建货款代扣失败: %w", err)
	}

	// 更新下发记录的货款代扣ID
	distribute.GoodsDeductionID = &deduction.ID
	if err := s.distributeRepo.Update(distribute); err != nil {
		log.Printf("[TerminalDistributeService] Update distribute goods_deduction_id failed: %v", err)
	}

	log.Printf("[TerminalDistributeService] Created goods deduction %d for distribute %d", deduction.ID, distribute.ID)
	return nil
}

// RejectDistribute 拒绝终端下发
func (s *TerminalDistributeService) RejectDistribute(distributeID int64, confirmedBy int64) error {
	distribute, err := s.distributeRepo.FindByID(distributeID)
	if err != nil || distribute == nil {
		return fmt.Errorf("下发记录不存在")
	}

	if distribute.Status != models.TerminalDistributeStatusPending {
		return fmt.Errorf("下发记录状态不允许拒绝")
	}

	return s.distributeRepo.UpdateStatus(distributeID, models.TerminalDistributeStatusRejected, &confirmedBy)
}

// CancelDistribute 取消终端下发（下发方取消）
func (s *TerminalDistributeService) CancelDistribute(distributeID int64, cancelBy int64) error {
	distribute, err := s.distributeRepo.FindByID(distributeID)
	if err != nil || distribute == nil {
		return fmt.Errorf("下发记录不存在")
	}

	if distribute.Status != models.TerminalDistributeStatusPending {
		return fmt.Errorf("下发记录状态不允许取消")
	}

	if distribute.FromAgentID != cancelBy {
		return fmt.Errorf("只有下发方可以取消")
	}

	return s.distributeRepo.UpdateStatus(distributeID, models.TerminalDistributeStatusCancelled, nil)
}

// GetDistributeList 获取下发列表
func (s *TerminalDistributeService) GetDistributeList(agentID int64, direction string, status []int16, limit, offset int) ([]*models.TerminalDistribute, int64, error) {
	if direction == "from" {
		return s.distributeRepo.FindByFromAgent(agentID, status, limit, offset)
	}
	return s.distributeRepo.FindByToAgent(agentID, status, limit, offset)
}
