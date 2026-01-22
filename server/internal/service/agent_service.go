package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"xiangshoufu/internal/repository"
	"xiangshoufu/pkg/qrcode"
)

// AgentService 代理商服务
type AgentService struct {
	agentRepo       *repository.GormAgentRepository
	agentPolicyRepo *repository.GormAgentPolicyRepository
	walletRepo      *repository.GormWalletRepository
	transactionRepo *repository.GormTransactionRepository
	profitRepo      *repository.GormProfitRecordRepository
	qrCodeGenerator *qrcode.Generator
}

// NewAgentService 创建代理商服务
func NewAgentService(
	agentRepo *repository.GormAgentRepository,
	agentPolicyRepo *repository.GormAgentPolicyRepository,
	walletRepo *repository.GormWalletRepository,
	transactionRepo *repository.GormTransactionRepository,
	profitRepo *repository.GormProfitRecordRepository,
) *AgentService {
	return &AgentService{
		agentRepo:       agentRepo,
		agentPolicyRepo: agentPolicyRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		profitRepo:      profitRepo,
		qrCodeGenerator: qrcode.NewGenerator(nil), // 使用默认配置
	}
}

// AgentDetailResponse 代理商详情响应
type AgentDetailResponse struct {
	ID                   int64     `json:"id"`
	AgentNo              string    `json:"agent_no"`
	AgentName            string    `json:"agent_name"`
	ParentID             int64     `json:"parent_id"`
	ParentName           string    `json:"parent_name,omitempty"`
	Path                 string    `json:"path"`
	Level                int       `json:"level"`
	ContactName          string    `json:"contact_name"`
	ContactPhone         string    `json:"contact_phone"`
	IDCardNo             string    `json:"id_card_no"`
	InviteCode           string    `json:"invite_code"`
	QRCodeURL            string    `json:"qr_code_url"`
	// 结算卡信息
	BankName             string    `json:"bank_name"`
	BankAccount          string    `json:"bank_account"`
	BankCardNo           string    `json:"bank_card_no"`
	Status               int16     `json:"status"`
	StatusName           string    `json:"status_name"`
	DirectAgentCount     int       `json:"direct_agent_count"`
	DirectMerchantCount  int       `json:"direct_merchant_count"`
	TeamAgentCount       int       `json:"team_agent_count"`
	TeamMerchantCount    int       `json:"team_merchant_count"`
	RegisterTime         time.Time `json:"register_time"`
	CreatedAt            time.Time `json:"created_at"`
}

// GetAgentDetail 获取代理商详情
func (s *AgentService) GetAgentDetail(agentID int64) (*AgentDetailResponse, error) {
	agent, err := s.agentRepo.FindByIDFull(agentID)
	if err != nil {
		return nil, fmt.Errorf("查询代理商失败: %w", err)
	}
	if agent == nil {
		return nil, errors.New("代理商不存在")
	}

	resp := &AgentDetailResponse{
		ID:                  agent.ID,
		AgentNo:             agent.AgentNo,
		AgentName:           agent.AgentName,
		ParentID:            agent.ParentID,
		Path:                agent.Path,
		Level:               agent.Level,
		ContactName:         agent.ContactName,
		ContactPhone:        agent.ContactPhone,
		IDCardNo:            agent.IDCardNo,
		InviteCode:          agent.InviteCode,
		QRCodeURL:           agent.QRCodeURL,
		BankName:            agent.BankName,
		BankAccount:         agent.BankAccount,
		BankCardNo:          agent.BankCardNo,
		Status:              agent.Status,
		StatusName:          getAgentStatusName(agent.Status),
		DirectAgentCount:    agent.DirectAgentCount,
		DirectMerchantCount: agent.DirectMerchantCount,
		TeamAgentCount:      agent.TeamAgentCount,
		TeamMerchantCount:   agent.TeamMerchantCount,
		RegisterTime:        agent.RegisterTime,
		CreatedAt:           agent.CreatedAt,
	}

	// 获取上级代理商名称
	if agent.ParentID > 0 {
		parent, _ := s.agentRepo.FindByID(agent.ParentID)
		if parent != nil {
			resp.ParentName = parent.AgentName
		}
	}

	return resp, nil
}

// SubordinateListRequest 下级代理商列表请求
type SubordinateListRequest struct {
	AgentID  int64  `json:"-"`
	Keyword  string `form:"keyword"`    // 搜索关键词
	Status   *int16 `form:"status"`     // 状态筛选
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// SubordinateInfo 下级代理商信息
type SubordinateInfo struct {
	ID                  int64     `json:"id"`
	AgentNo             string    `json:"agent_no"`
	AgentName           string    `json:"agent_name"`
	Level               int       `json:"level"`
	ContactPhone        string    `json:"contact_phone"`
	Status              int16     `json:"status"`
	StatusName          string    `json:"status_name"`
	DirectAgentCount    int       `json:"direct_agent_count"`
	DirectMerchantCount int       `json:"direct_merchant_count"`
	TeamAgentCount      int       `json:"team_agent_count"`
	TeamMerchantCount   int       `json:"team_merchant_count"`
	RegisterTime        time.Time `json:"register_time"`
}

// GetSubordinateList 获取下级代理商列表
func (s *AgentService) GetSubordinateList(req *SubordinateListRequest) ([]*SubordinateInfo, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	agents, total, err := s.agentRepo.FindSubordinates(req.AgentID, req.Keyword, req.Status, req.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询下级代理商失败: %w", err)
	}

	list := make([]*SubordinateInfo, 0, len(agents))
	for _, agent := range agents {
		list = append(list, &SubordinateInfo{
			ID:                  agent.ID,
			AgentNo:             agent.AgentNo,
			AgentName:           agent.AgentName,
			Level:               agent.Level,
			ContactPhone:        agent.ContactPhone,
			Status:              agent.Status,
			StatusName:          getAgentStatusName(agent.Status),
			DirectAgentCount:    agent.DirectAgentCount,
			DirectMerchantCount: agent.DirectMerchantCount,
			TeamAgentCount:      agent.TeamAgentCount,
			TeamMerchantCount:   agent.TeamMerchantCount,
			RegisterTime:        agent.RegisterTime,
		})
	}

	return list, total, nil
}

// TeamTreeNode 团队层级树节点
type TeamTreeNode struct {
	ID           int64           `json:"id"`
	AgentNo      string          `json:"agent_no"`
	AgentName    string          `json:"agent_name"`
	Level        int             `json:"level"`
	Status       int16           `json:"status"`
	ChildCount   int             `json:"child_count"`
	Children     []*TeamTreeNode `json:"children,omitempty"`
}

// GetTeamTree 获取团队层级树
func (s *AgentService) GetTeamTree(agentID int64, maxDepth int) (*TeamTreeNode, error) {
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil || agent == nil {
		return nil, errors.New("代理商不存在")
	}

	root := &TeamTreeNode{
		ID:        agent.ID,
		AgentNo:   agent.AgentNo,
		AgentName: agent.AgentName,
		Level:     agent.Level,
		Status:    agent.Status,
	}

	if maxDepth > 0 {
		children, err := s.buildTreeChildren(agentID, 1, maxDepth)
		if err != nil {
			return nil, err
		}
		root.Children = children
		root.ChildCount = len(children)
	}

	return root, nil
}

// buildTreeChildren 递归构建子树
func (s *AgentService) buildTreeChildren(parentID int64, currentDepth, maxDepth int) ([]*TeamTreeNode, error) {
	if currentDepth >= maxDepth {
		return nil, nil
	}

	children, _, err := s.agentRepo.FindSubordinates(parentID, "", nil, 100, 0)
	if err != nil {
		return nil, err
	}

	nodes := make([]*TeamTreeNode, 0, len(children))
	for _, child := range children {
		node := &TeamTreeNode{
			ID:        child.ID,
			AgentNo:   child.AgentNo,
			AgentName: child.AgentName,
			Level:     child.Level,
			Status:    child.Status,
		}

		subChildren, _ := s.buildTreeChildren(child.ID, currentDepth+1, maxDepth)
		if len(subChildren) > 0 {
			node.Children = subChildren
			node.ChildCount = len(subChildren)
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AgentStatsResponse 代理商统计响应
type AgentStatsResponse struct {
	// 团队统计
	DirectAgentCount    int `json:"direct_agent_count"`
	DirectMerchantCount int `json:"direct_merchant_count"`
	TeamAgentCount      int `json:"team_agent_count"`
	TeamMerchantCount   int `json:"team_merchant_count"`

	// 今日数据
	TodayTransAmount   int64 `json:"today_trans_amount"`    // 今日交易额（分）
	TodayTransCount    int64 `json:"today_trans_count"`     // 今日交易笔数
	TodayProfitAmount  int64 `json:"today_profit_amount"`   // 今日分润（分）
	TodayNewAgents     int   `json:"today_new_agents"`      // 今日新增代理商
	TodayNewMerchants  int   `json:"today_new_merchants"`   // 今日新增商户

	// 本月数据
	MonthTransAmount   int64 `json:"month_trans_amount"`
	MonthTransCount    int64 `json:"month_trans_count"`
	MonthProfitAmount  int64 `json:"month_profit_amount"`
}

// GetAgentStats 获取代理商统计
func (s *AgentService) GetAgentStats(agentID int64) (*AgentStatsResponse, error) {
	agent, err := s.agentRepo.FindByIDFull(agentID)
	if err != nil || agent == nil {
		return nil, errors.New("代理商不存在")
	}

	stats := &AgentStatsResponse{
		DirectAgentCount:    agent.DirectAgentCount,
		DirectMerchantCount: agent.DirectMerchantCount,
		TeamAgentCount:      agent.TeamAgentCount,
		TeamMerchantCount:   agent.TeamMerchantCount,
	}

	// 获取今日交易统计
	todayStats, err := s.transactionRepo.GetAgentDailyStats(agentID, time.Now())
	if err == nil && todayStats != nil {
		stats.TodayTransAmount = todayStats.TotalAmount
		stats.TodayTransCount = todayStats.TotalCount
	}

	// 获取本月交易统计
	monthStats, err := s.transactionRepo.GetAgentMonthlyStats(agentID, time.Now())
	if err == nil && monthStats != nil {
		stats.MonthTransAmount = monthStats.TotalAmount
		stats.MonthTransCount = monthStats.TotalCount
	}

	// 获取今日分润统计
	todayProfitStats, err := s.profitRepo.GetAgentDailyProfitStats(agentID, time.Now())
	if err == nil && todayProfitStats != nil {
		stats.TodayProfitAmount = todayProfitStats.TotalAmount
	}

	// 获取本月分润统计
	monthProfitStats, err := s.profitRepo.GetAgentMonthlyProfitStats(agentID, time.Now())
	if err == nil && monthProfitStats != nil {
		stats.MonthProfitAmount = monthProfitStats.TotalAmount
	}

	return stats, nil
}

// UpdateAgentProfileRequest 更新代理商资料请求
type UpdateAgentProfileRequest struct {
	AgentID      int64  `json:"-"`
	AgentName    string `json:"agent_name"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	BankName     string `json:"bank_name"`
	BankAccount  string `json:"bank_account"`
	BankCardNo   string `json:"bank_card_no"`
}

// UpdateAgentProfile 更新代理商资料
func (s *AgentService) UpdateAgentProfile(req *UpdateAgentProfileRequest) error {
	agent, err := s.agentRepo.FindByIDFull(req.AgentID)
	if err != nil || agent == nil {
		return errors.New("代理商不存在")
	}

	// 更新可修改字段
	if req.AgentName != "" {
		agent.AgentName = req.AgentName
	}
	if req.ContactName != "" {
		agent.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		agent.ContactPhone = req.ContactPhone
	}
	if req.BankName != "" {
		agent.BankName = req.BankName
	}
	if req.BankAccount != "" {
		agent.BankAccount = req.BankAccount
	}
	if req.BankCardNo != "" {
		agent.BankCardNo = req.BankCardNo
	}

	if err := s.agentRepo.Update(agent); err != nil {
		return fmt.Errorf("更新失败: %w", err)
	}

	log.Printf("[AgentService] Updated agent profile: %d", agent.ID)

	return nil
}

// GetInviteCode 获取邀请码
func (s *AgentService) GetInviteCode(agentID int64) (string, string, error) {
	agent, err := s.agentRepo.FindByIDFull(agentID)
	if err != nil || agent == nil {
		return "", "", errors.New("代理商不存在")
	}

	// 如果没有邀请码，生成一个
	if agent.InviteCode == "" {
		inviteCode := generateInviteCode(agentID)
		agent.InviteCode = inviteCode
		s.agentRepo.Update(agent)
	}

	// 如果没有二维码URL，生成二维码
	if agent.QRCodeURL == "" {
		_, qrCodeURL, err := s.qrCodeGenerator.GenerateInviteQRCode(agent.InviteCode, agent.ID)
		if err == nil && qrCodeURL != "" {
			agent.QRCodeURL = qrCodeURL
			s.agentRepo.Update(agent)
		}
	}

	return agent.InviteCode, agent.QRCodeURL, nil
}

// GetAncestors 获取所有上级代理商
func (s *AgentService) GetAncestors(agentID int64) ([]*repository.Agent, error) {
	return s.agentRepo.FindAncestors(agentID)
}

// IsSubordinate 检查是否为下级代理商
func (s *AgentService) IsSubordinate(parentID, childID int64) (bool, error) {
	child, err := s.agentRepo.FindByID(childID)
	if err != nil || child == nil {
		return false, errors.New("代理商不存在")
	}

	// 检查path是否包含parentID
	parentPath := fmt.Sprintf("/%d/", parentID)
	return strings.Contains(child.Path, parentPath), nil
}

// getAgentStatusName 获取代理商状态名称
func getAgentStatusName(status int16) string {
	switch status {
	case 1:
		return "正常"
	case 2:
		return "禁用"
	case 3:
		return "待审核"
	default:
		return "未知"
	}
}

// generateInviteCode 生成邀请码
func generateInviteCode(agentID int64) string {
	// 简单的邀请码生成逻辑
	chars := "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[(agentID+int64(i)*7)%int64(len(chars))]
	}
	return string(code)
}

// CreateAgentRequest 创建代理商请求
type CreateAgentRequest struct {
	AgentName    string `json:"agent_name" binding:"required"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	IDCardNo     string `json:"id_card_no"`
	BankName     string `json:"bank_name"`
	BankAccount  string `json:"bank_account"`
	BankCardNo   string `json:"bank_card_no"`
	ParentID     int64  `json:"parent_id"` // 上级代理商ID，0表示顶级
}

// CreateAgent 创建代理商
func (s *AgentService) CreateAgent(req *CreateAgentRequest, operatorID int64) (*AgentDetailResponse, error) {
	// 验证手机号唯一性
	existing, _ := s.agentRepo.FindByPhone(req.ContactPhone)
	if existing != nil {
		return nil, errors.New("该手机号已被注册")
	}

	// 获取上级代理商信息（如果有）
	var parentAgent *repository.AgentFull
	var path string
	var level int = 1

	if req.ParentID > 0 {
		var err error
		parentAgent, err = s.agentRepo.FindByIDFull(req.ParentID)
		if err != nil || parentAgent == nil {
			return nil, errors.New("上级代理商不存在")
		}
		level = parentAgent.Level + 1
		// 限制最大层级为10
		if level > 10 {
			return nil, errors.New("代理商层级不能超过10级")
		}
	}

	// 生成代理商编号
	agentNo := s.generateAgentNo()

	// 创建代理商对象
	now := time.Now()
	agent := &repository.AgentFull{
		Agent: repository.Agent{
			AgentNo:     agentNo,
			AgentName:   req.AgentName,
			ParentID:    req.ParentID,
			Level:       level,
			DefaultRate: "0.0060", // 默认费率0.60%
			Status:      1,        // 默认正常状态
		},
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		IDCardNo:     req.IDCardNo,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankCardNo:   req.BankCardNo,
		RegisterTime: now,
		CreatedAt:    now,
	}

	// 创建代理商
	if err := s.agentRepo.Create(agent); err != nil {
		return nil, fmt.Errorf("创建代理商失败: %w", err)
	}

	// 构建物化路径
	if parentAgent != nil {
		path = fmt.Sprintf("%s%d/", parentAgent.Path, agent.ID)
	} else {
		path = fmt.Sprintf("/%d/", agent.ID)
	}
	agent.Path = path

	// 生成邀请码
	agent.InviteCode = generateInviteCode(agent.ID)

	// 更新路径和邀请码
	if err := s.agentRepo.Update(agent); err != nil {
		return nil, fmt.Errorf("更新代理商路径失败: %w", err)
	}

	// 更新上级代理商的统计计数
	if parentAgent != nil {
		s.agentRepo.IncrementDirectAgentCount(parentAgent.ID)
		// 更新所有祖先的团队代理商计数
		ancestors, _ := s.agentRepo.FindAncestors(agent.ID)
		for _, ancestor := range ancestors {
			s.agentRepo.IncrementTeamAgentCount(ancestor.ID)
		}
	}

	log.Printf("[AgentService] Created agent: %s (%d) by operator %d", agentNo, agent.ID, operatorID)

	return s.GetAgentDetail(agent.ID)
}

// CreateAgentWithStatus 创建代理商（指定状态，用于邀请码注册）
// status: 1正常 2禁用 3待审核
// password: 注册时设置的密码（用于创建登录账号）
func (s *AgentService) CreateAgentWithStatus(req *CreateAgentRequest, operatorID int64, status int16, password string) (*AgentDetailResponse, error) {
	// 验证手机号唯一性
	existing, _ := s.agentRepo.FindByPhone(req.ContactPhone)
	if existing != nil {
		return nil, errors.New("该手机号已被注册")
	}

	// 获取上级代理商信息（如果有）
	var parentAgent *repository.AgentFull
	var path string
	var level int = 1

	if req.ParentID > 0 {
		var err error
		parentAgent, err = s.agentRepo.FindByIDFull(req.ParentID)
		if err != nil || parentAgent == nil {
			return nil, errors.New("上级代理商不存在")
		}
		level = parentAgent.Level + 1
		// 限制最大层级为10
		if level > 10 {
			return nil, errors.New("代理商层级不能超过10级")
		}
	}

	// 生成代理商编号
	agentNo := s.generateAgentNo()

	// 创建代理商对象
	now := time.Now()
	agent := &repository.AgentFull{
		Agent: repository.Agent{
			AgentNo:     agentNo,
			AgentName:   req.AgentName,
			ParentID:    req.ParentID,
			Level:       level,
			DefaultRate: "0.0060", // 默认费率0.60%
			Status:      status,   // 使用指定的状态
		},
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		IDCardNo:     req.IDCardNo,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankCardNo:   req.BankCardNo,
		RegisterTime: now,
		CreatedAt:    now,
	}

	// 创建代理商
	if err := s.agentRepo.Create(agent); err != nil {
		return nil, fmt.Errorf("创建代理商失败: %w", err)
	}

	// 构建物化路径
	if parentAgent != nil {
		path = fmt.Sprintf("%s%d/", parentAgent.Path, agent.ID)
	} else {
		path = fmt.Sprintf("/%d/", agent.ID)
	}
	agent.Path = path

	// 生成邀请码
	agent.InviteCode = generateInviteCode(agent.ID)

	// 更新路径和邀请码
	if err := s.agentRepo.Update(agent); err != nil {
		return nil, fmt.Errorf("更新代理商路径失败: %w", err)
	}

	// 更新上级代理商的统计计数
	if parentAgent != nil {
		s.agentRepo.IncrementDirectAgentCount(parentAgent.ID)
		// 更新所有祖先的团队代理商计数
		ancestors, _ := s.agentRepo.FindAncestors(agent.ID)
		for _, ancestor := range ancestors {
			s.agentRepo.IncrementTeamAgentCount(ancestor.ID)
		}

		// 继承上级代理商的政策模板（状态设为待调整）
		s.inheritParentPolicies(parentAgent.ID, agent.ID)
	}

	log.Printf("[AgentService] Created agent via invite code: %s (%d) status=%d", agentNo, agent.ID, status)

	return s.GetAgentDetail(agent.ID)
}

// inheritParentPolicies 继承上级代理商的政策
// 复制上级的 AgentPolicy 记录到新代理商，费率需要上级调整后审核通过
func (s *AgentService) inheritParentPolicies(parentID, newAgentID int64) {
	// 获取上级代理商的所有政策
	parentPolicies, err := s.agentPolicyRepo.FindByAgentID(parentID)
	if err != nil || len(parentPolicies) == 0 {
		log.Printf("[AgentService] No policies to inherit from parent %d", parentID)
		return
	}

	// 复制政策到新代理商
	for _, policy := range parentPolicies {
		newPolicy := &repository.AgentPolicy{
			AgentID:    newAgentID,
			ChannelID:  policy.ChannelID,
			TemplateID: policy.TemplateID,
			CreditRate: policy.CreditRate, // 继承上级费率，上级需要调整后审核
			DebitRate:  policy.DebitRate,
		}
		// 创建新政策记录
		if err := s.agentPolicyRepo.Create(newPolicy); err != nil {
			log.Printf("[AgentService] Failed to inherit policy for channel %d: %v", policy.ChannelID, err)
		} else {
			log.Printf("[AgentService] Inherited policy from parent %d to agent %d for channel %d", parentID, newAgentID, policy.ChannelID)
		}
	}
}

// generateAgentNo 生成代理商编号
func (s *AgentService) generateAgentNo() string {
	// 格式: A + 年月日 + 4位序号，如 A202501200001
	now := time.Now()
	dateStr := now.Format("20060102")
	// 获取当日序号
	seq := s.agentRepo.GetDailyAgentSequence(dateStr)
	return fmt.Sprintf("A%s%04d", dateStr, seq)
}

// UpdateAgentStatusRequest 更新代理商状态请求
type UpdateAgentStatusRequest struct {
	AgentID int64 `json:"-"`
	Status  int16 `json:"status" binding:"required,oneof=1 2"` // 1正常 2禁用
}

// UpdateAgentStatus 更新代理商状态
func (s *AgentService) UpdateAgentStatus(req *UpdateAgentStatusRequest) error {
	agent, err := s.agentRepo.FindByIDFull(req.AgentID)
	if err != nil || agent == nil {
		return errors.New("代理商不存在")
	}

	if agent.Status == req.Status {
		return nil // 状态未变更
	}

	agent.Status = req.Status
	if err := s.agentRepo.Update(agent); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	statusName := getAgentStatusName(req.Status)
	log.Printf("[AgentService] Updated agent %d status to %s", req.AgentID, statusName)

	return nil
}

// RegisterByInviteCodeRequest 通过邀请码注册请求
type RegisterByInviteCodeRequest struct {
	InviteCode   string `json:"invite_code" binding:"required"`
	AgentName    string `json:"agent_name" binding:"required"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	IDCardNo     string `json:"id_card_no"`
	Password     string `json:"password" binding:"required,min=6"`
}

// RegisterByInviteCode 通过邀请码注册代理商
// 注册后状态为"待审核"(status=3)，需上级审核通过后才能使用
func (s *AgentService) RegisterByInviteCode(req *RegisterByInviteCodeRequest) (*AgentDetailResponse, error) {
	// 查找邀请码对应的代理商
	inviter, err := s.agentRepo.FindByInviteCode(req.InviteCode)
	if err != nil || inviter == nil {
		return nil, errors.New("邀请码无效")
	}

	if inviter.Status != 1 {
		return nil, errors.New("邀请人账号已被禁用")
	}

	// 使用创建代理商的逻辑（自注册模式，状态为待审核）
	createReq := &CreateAgentRequest{
		AgentName:    req.AgentName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		IDCardNo:     req.IDCardNo,
		ParentID:     inviter.ID,
	}

	return s.CreateAgentWithStatus(createReq, 0, 3, req.Password) // 0表示自注册, 3表示待审核状态
}

// SearchAgentsRequest 搜索代理商请求
type SearchAgentsRequest struct {
	Keyword  string `form:"keyword"`
	Status   *int16 `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// SearchAgents 全局搜索代理商
func (s *AgentService) SearchAgents(req *SearchAgentsRequest) ([]*SubordinateInfo, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize

	agents, total, err := s.agentRepo.SearchAgents(req.Keyword, req.Status, req.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索代理商失败: %w", err)
	}

	list := make([]*SubordinateInfo, 0, len(agents))
	for _, agent := range agents {
		list = append(list, &SubordinateInfo{
			ID:                  agent.ID,
			AgentNo:             agent.AgentNo,
			AgentName:           agent.AgentName,
			Level:               agent.Level,
			ContactPhone:        agent.ContactPhone,
			Status:              agent.Status,
			StatusName:          getAgentStatusName(agent.Status),
			DirectAgentCount:    agent.DirectAgentCount,
			DirectMerchantCount: agent.DirectMerchantCount,
			TeamAgentCount:      agent.TeamAgentCount,
			TeamMerchantCount:   agent.TeamMerchantCount,
			RegisterTime:        agent.RegisterTime,
		})
	}

	return list, total, nil
}

// SetCustomInviteCodeRequest 设置自定义邀请码请求
type SetCustomInviteCodeRequest struct {
	AgentID    int64  `json:"-"`
	InviteCode string `json:"invite_code" binding:"required,min=4,max=12,alphanum"`
}

// SetCustomInviteCode 设置自定义邀请码（靓号）
func (s *AgentService) SetCustomInviteCode(req *SetCustomInviteCodeRequest) error {
	// 验证邀请码格式
	if len(req.InviteCode) < 4 || len(req.InviteCode) > 12 {
		return fmt.Errorf("邀请码长度必须在4-12位之间")
	}

	// 检查邀请码是否只包含字母和数字
	for _, c := range req.InviteCode {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return fmt.Errorf("邀请码只能包含字母和数字")
		}
	}

	// 转换为大写以保持一致性
	inviteCode := strings.ToUpper(req.InviteCode)

	// 检查邀请码是否已被其他代理商使用
	existingAgent, err := s.agentRepo.FindByInviteCode(inviteCode)
	if err == nil && existingAgent != nil && existingAgent.ID != req.AgentID {
		return fmt.Errorf("该邀请码已被其他代理商使用")
	}

	// 更新代理商的邀请码
	agent, err := s.agentRepo.FindByIDFull(req.AgentID)
	if err != nil {
		return fmt.Errorf("代理商不存在: %w", err)
	}
	if agent == nil {
		return fmt.Errorf("代理商不存在")
	}

	agent.InviteCode = inviteCode
	if err := s.agentRepo.Update(agent); err != nil {
		return fmt.Errorf("更新邀请码失败: %w", err)
	}

	return nil
}
