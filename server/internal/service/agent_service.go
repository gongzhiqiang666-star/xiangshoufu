package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"xiangshoufu/internal/repository"
)

// AgentService 代理商服务
type AgentService struct {
	agentRepo       *repository.GormAgentRepository
	agentPolicyRepo *repository.GormAgentPolicyRepository
	walletRepo      *repository.GormWalletRepository
	transactionRepo *repository.GormTransactionRepository
}

// NewAgentService 创建代理商服务
func NewAgentService(
	agentRepo *repository.GormAgentRepository,
	agentPolicyRepo *repository.GormAgentPolicyRepository,
	walletRepo *repository.GormWalletRepository,
	transactionRepo *repository.GormTransactionRepository,
) *AgentService {
	return &AgentService{
		agentRepo:       agentRepo,
		agentPolicyRepo: agentPolicyRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
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
	InviteCode           string    `json:"invite_code"`
	QRCodeURL            string    `json:"qr_code_url"`
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
		InviteCode:          agent.InviteCode,
		QRCodeURL:           agent.QRCodeURL,
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

	// TODO: 获取分润统计（需要扩展repository）

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
