package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/pkg/crypto"
)

// MerchantService 商户服务
type MerchantService struct {
	merchantRepo    *repository.GormMerchantRepository
	agentRepo       *repository.GormAgentRepository
	transactionRepo *repository.GormTransactionRepository
	terminalRepo    *repository.GormTerminalRepository
	rateSyncService *RateSyncService
}

// NewMerchantService 创建商户服务
func NewMerchantService(
	merchantRepo *repository.GormMerchantRepository,
	agentRepo *repository.GormAgentRepository,
	transactionRepo *repository.GormTransactionRepository,
	terminalRepo *repository.GormTerminalRepository,
) *MerchantService {
	return &MerchantService{
		merchantRepo:    merchantRepo,
		agentRepo:       agentRepo,
		transactionRepo: transactionRepo,
		terminalRepo:    terminalRepo,
	}
}

// SetRateSyncService 设置费率同步服务（可选注入，避免循环依赖）
func (s *MerchantService) SetRateSyncService(rateSyncService *RateSyncService) {
	s.rateSyncService = rateSyncService
}

// ==================== 请求/响应结构体 ====================

// CreateMerchantRequest 创建商户请求
type CreateMerchantRequest struct {
	MerchantNo   string `json:"merchant_no" binding:"required"`
	MerchantName string `json:"merchant_name" binding:"required"`
	AgentID      int64  `json:"agent_id" binding:"required"`
	ChannelID    int64  `json:"channel_id" binding:"required"`
	TerminalSN   string `json:"terminal_sn"`
	LegalName    string `json:"legal_name"`
	LegalIDCard  string `json:"legal_id_card"`
	MCC          string `json:"mcc"`
	CreditRate   string `json:"credit_rate"`
	DebitRate    string `json:"debit_rate"`
	IsDirect     bool   `json:"is_direct"`
}

// UpdateMerchantRequest 更新商户请求
type UpdateMerchantRequest struct {
	MerchantName string `json:"merchant_name"`
	LegalName    string `json:"legal_name"`
	LegalIDCard  string `json:"legal_id_card"`
	MCC          string `json:"mcc"`
	TerminalSN   string `json:"terminal_sn"`
}

// UpdateRateRequest 更新费率请求
type UpdateRateRequest struct {
	CreditRate float64 `json:"credit_rate" binding:"required"`
	DebitRate  float64 `json:"debit_rate" binding:"required"`
}

// RegisterMerchantRequest 商户登记请求
type RegisterMerchantRequest struct {
	Phone  string `json:"phone" binding:"required"`
	Remark string `json:"remark"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status int16 `json:"status" binding:"required,oneof=1 2"`
}

// MerchantDetailResponse 商户详情响应
type MerchantDetailResponse struct {
	ID              int64      `json:"id"`
	MerchantNo      string     `json:"merchant_no"`
	MerchantName    string     `json:"merchant_name"`
	AgentID         int64      `json:"agent_id"`
	AgentName       string     `json:"agent_name"`
	AgentLevel      int        `json:"agent_level"`
	ChannelID       int64      `json:"channel_id"`
	ChannelName     string     `json:"channel_name"`
	TerminalSN      string     `json:"terminal_sn"`
	Status          int16      `json:"status"`
	StatusName      string     `json:"status_name"`
	ApproveStatus   int16      `json:"approve_status"`
	ApproveStatusName string   `json:"approve_status_name"`
	LegalName       string     `json:"legal_name"`
	LegalIDCard     string     `json:"legal_id_card"`
	MCC             string     `json:"mcc"`
	CreditRate      string     `json:"credit_rate"`
	DebitRate       string     `json:"debit_rate"`
	MerchantType    string     `json:"merchant_type"`
	MerchantTypeName string    `json:"merchant_type_name"`
	IsDirect        bool       `json:"is_direct"`
	OwnerType       string     `json:"owner_type"`
	ActivatedAt     *time.Time `json:"activated_at"`
	RegisteredPhone string     `json:"registered_phone"`
	RegisterRemark  string     `json:"register_remark"`
	// 统计数据
	MonthAmount     int64      `json:"month_amount"`
	MonthCount      int64      `json:"month_count"`
	TerminalCount   int        `json:"terminal_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ==================== 服务方法 ====================

// CreateMerchant 创建商户
func (s *MerchantService) CreateMerchant(req *CreateMerchantRequest) (*models.Merchant, error) {
	// 验证商户编号唯一性
	existing, _ := s.merchantRepo.FindByMerchantNo(req.MerchantNo)
	if existing != nil {
		return nil, errors.New("商户编号已存在")
	}

	// 验证代理商存在
	agent, err := s.agentRepo.FindByID(req.AgentID)
	if err != nil || agent == nil {
		return nil, errors.New("代理商不存在")
	}

	// 加密身份证号（三级等保要求）
	encryptedIDCard := req.LegalIDCard
	if req.LegalIDCard != "" {
		encryptedIDCard, err = crypto.EncryptIDCard(req.LegalIDCard)
		if err != nil {
			log.Printf("[MerchantService] Failed to encrypt ID card: %v", err)
			return nil, fmt.Errorf("身份证加密失败: %w", err)
		}
	}

	merchant := &models.Merchant{
		MerchantNo:   req.MerchantNo,
		MerchantName: req.MerchantName,
		AgentID:      req.AgentID,
		ChannelID:    req.ChannelID,
		TerminalSN:   req.TerminalSN,
		LegalName:    req.LegalName,
		LegalIDCard:  encryptedIDCard,
		MCC:          req.MCC,
		CreditRate:   req.CreditRate,
		DebitRate:    req.DebitRate,
		IsDirect:     req.IsDirect,
		MerchantType: models.MerchantTypeNormal,
		Status:       models.MerchantStatusActive,
		ApproveStatus: models.MerchantApproveStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.merchantRepo.Create(merchant); err != nil {
		return nil, fmt.Errorf("创建商户失败: %w", err)
	}

	return merchant, nil
}

// UpdateMerchant 更新商户信息
func (s *MerchantService) UpdateMerchant(id int64, req *UpdateMerchantRequest) error {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("商户不存在: %w", err)
	}

	if req.MerchantName != "" {
		merchant.MerchantName = req.MerchantName
	}
	if req.LegalName != "" {
		merchant.LegalName = req.LegalName
	}
	if req.LegalIDCard != "" {
		merchant.LegalIDCard = req.LegalIDCard
	}
	if req.MCC != "" {
		merchant.MCC = req.MCC
	}
	if req.TerminalSN != "" {
		merchant.TerminalSN = req.TerminalSN
	}
	merchant.UpdatedAt = time.Now()

	return s.merchantRepo.Update(merchant)
}

// RateUpdateResult 费率更新结果
type RateUpdateResult struct {
	Success     bool   `json:"success"`      // 是否成功
	SyncSuccess bool   `json:"sync_success"` // 通道同步是否成功
	SyncMessage string `json:"sync_message"` // 通道同步消息
}

// UpdateRate 修改费率（校验费率范围，并同步到通道）
// 返回更新结果，包含本地更新和通道同步状态
func (s *MerchantService) UpdateRate(id int64, agentID int64, creditRate, debitRate float64) (*RateUpdateResult, error) {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("商户不存在: %w", err)
	}

	// 验证商户归属
	if merchant.AgentID != agentID {
		return nil, errors.New("无权修改此商户费率")
	}

	// 验证费率范围 (0-1之间，通常为0.001-0.01即0.1%-1%)
	if creditRate < 0 || creditRate > 0.1 {
		return nil, errors.New("贷记卡费率范围无效")
	}
	if debitRate < 0 || debitRate > 0.1 {
		return nil, errors.New("借记卡费率范围无效")
	}

	// 校验费率不能低于上级代理商费率
	agent, err := s.agentRepo.FindByIDFull(agentID)
	if err != nil || agent == nil {
		return nil, errors.New("代理商不存在")
	}

	// 获取上级代理商的政策费率（如果有上级）
	if agent.ParentID > 0 {
		parentPolicy, err := s.getParentPolicyRate(agent.ParentID, merchant.ChannelID)
		if err == nil && parentPolicy != nil {
			parentCreditRate := parseRate(parentPolicy.CreditRate)
			parentDebitRate := parseRate(parentPolicy.DebitRate)

			if creditRate < parentCreditRate {
				return nil, fmt.Errorf("贷记卡费率不能低于上级费率 %.4f", parentCreditRate)
			}
			if debitRate < parentDebitRate {
				return nil, fmt.Errorf("借记卡费率不能低于上级费率 %.4f", parentDebitRate)
			}
		}
	}

	// 保存原费率用于日志记录
	oldCreditRate := parseRate(merchant.CreditRate)
	oldDebitRate := parseRate(merchant.DebitRate)

	creditRateStr := strconv.FormatFloat(creditRate, 'f', 4, 64)
	debitRateStr := strconv.FormatFloat(debitRate, 'f', 4, 64)

	result := &RateUpdateResult{
		Success:     true,
		SyncSuccess: true,
		SyncMessage: "费率更新成功",
	}

	// 【重要】先同步到支付通道，通道成功后才能修改本地费率
	// 原则：费率修改必须和通道联动，先调用通道，通道返回成功才能修改自己的费率
	if s.rateSyncService != nil {
		syncResult := s.syncRateToChannel(merchant, agentID, oldCreditRate, oldDebitRate, creditRate, debitRate)
		result.SyncSuccess = syncResult.Success
		result.SyncMessage = syncResult.Message

		// 通道同步失败，不更新本地数据库
		if !syncResult.Success {
			log.Printf("[MerchantService] 通道费率同步失败，本地费率不更新: merchantID=%d, err=%s", id, syncResult.Message)
			return result, nil
		}
	}

	// 通道同步成功（或无需同步），更新本地数据库
	if err := s.merchantRepo.UpdateRate(id, creditRateStr, debitRateStr); err != nil {
		return nil, fmt.Errorf("通道同步成功但本地更新失败: %w", err)
	}

	return result, nil
}

// syncRateToChannel 同步费率到通道（同步调用）
func (s *MerchantService) syncRateToChannel(merchant *models.Merchant, agentID int64, oldCreditRate, oldDebitRate, newCreditRate, newDebitRate float64) *SyncResult {
	// 获取通道编码
	channelCode := s.getChannelCode(merchant.ChannelID)
	if channelCode == "" {
		log.Printf("[MerchantService] 获取通道编码失败，跳过同步: merchantID=%d", merchant.ID)
		return &SyncResult{
			Success: true,
			Message: "通道编码未配置，跳过同步",
		}
	}

	params := &RateUpdateParams{
		MerchantID:  merchant.ID,
		MerchantNo:  merchant.MerchantNo,
		TerminalSN:  merchant.TerminalSN,
		ChannelCode: channelCode,
		AgentID:     agentID,
		OldRates: &RateInfo{
			CreditRate: oldCreditRate,
			DebitRate:  oldDebitRate,
		},
		NewRates: &RateInfo{
			CreditRate: newCreditRate,
			DebitRate:  newDebitRate,
		},
	}

	result, err := s.rateSyncService.SyncRateToChannel(nil, params)
	if err != nil {
		log.Printf("[MerchantService] 费率同步失败: merchantID=%d, error=%v", merchant.ID, err)
		return &SyncResult{
			Success: false,
			Message: fmt.Sprintf("同步失败: %v", err),
		}
	}

	if result.Success {
		log.Printf("[MerchantService] 费率同步成功: merchantID=%d, tradeNo=%s", merchant.ID, result.TradeNo)
	} else {
		log.Printf("[MerchantService] 费率同步失败: merchantID=%d, message=%s", merchant.ID, result.Message)
	}

	return result
}

// getChannelCode 根据通道ID获取通道编码
func (s *MerchantService) getChannelCode(channelID int64) string {
	query := `SELECT channel_code FROM channels WHERE id = ? LIMIT 1`
	var channelCode string
	if err := s.merchantRepo.GetDB().Raw(query, channelID).Scan(&channelCode).Error; err != nil {
		return ""
	}
	return channelCode
}

// RegisterMerchant 商户登记
func (s *MerchantService) RegisterMerchant(id int64, agentID int64, phone, remark string) error {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("商户不存在: %w", err)
	}

	// 验证商户归属
	if merchant.AgentID != agentID {
		return errors.New("无权登记此商户")
	}

	// 对手机号进行加密存储（三级等保要求）
	encryptedPhone, err := crypto.EncryptPhone(phone)
	if err != nil {
		log.Printf("[MerchantService] Failed to encrypt phone: %v", err)
		return fmt.Errorf("手机号加密失败: %w", err)
	}

	return s.merchantRepo.Register(id, encryptedPhone, remark)
}

// UpdateStatus 更新商户状态
func (s *MerchantService) UpdateStatus(id int64, agentID int64, status int16) error {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("商户不存在: %w", err)
	}

	// 验证商户归属
	if merchant.AgentID != agentID {
		return errors.New("无权修改此商户状态")
	}

	// 验证状态值
	if status != models.MerchantStatusActive && status != models.MerchantStatusDisabled {
		return errors.New("无效的状态值")
	}

	return s.merchantRepo.UpdateStatus(id, status)
}

// DeleteMerchant 删除商户
func (s *MerchantService) DeleteMerchant(id int64, agentID int64) error {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("商户不存在: %w", err)
	}

	// 验证商户归属
	if merchant.AgentID != agentID {
		return errors.New("无权删除此商户")
	}

	return s.merchantRepo.Delete(id)
}

// CalculateMerchantType 计算商户类型（5档分类）
// 分类逻辑：
// 1. 先判断60天无交易 → churned（流失）
// 2. 再判断30天无交易 → warning（预警）
// 3. 有交易按金额分档：≥5万→quality，3-5万→medium，<3万→normal
func (s *MerchantService) CalculateMerchantType(merchantID int64) (string, error) {
	endTime := time.Now()

	// 先检查60天内是否有交易
	startTime60 := endTime.AddDate(0, 0, -60)
	stats60, err := s.merchantRepo.GetMerchantTransStats(merchantID, &startTime60, &endTime)
	if err != nil {
		return "", fmt.Errorf("获取60天交易统计失败: %w", err)
	}

	// 60天无交易 → 流失
	if stats60.TotalCount == 0 {
		if err := s.merchantRepo.UpdateMerchantType(merchantID, models.MerchantTypeChurned); err != nil {
			return "", fmt.Errorf("更新商户类型失败: %w", err)
		}
		return models.MerchantTypeChurned, nil
	}

	// 检查30天内是否有交易
	startTime30 := endTime.AddDate(0, 0, -30)
	stats30, err := s.merchantRepo.GetMerchantTransStats(merchantID, &startTime30, &endTime)
	if err != nil {
		return "", fmt.Errorf("获取30天交易统计失败: %w", err)
	}

	// 30天无交易 → 预警
	if stats30.TotalCount == 0 {
		if err := s.merchantRepo.UpdateMerchantType(merchantID, models.MerchantTypeWarning); err != nil {
			return "", fmt.Errorf("更新商户类型失败: %w", err)
		}
		return models.MerchantTypeWarning, nil
	}

	// 有交易，按月均交易额分档（分转元）
	monthAmount := float64(stats30.TotalAmount) / 100

	var merchantType string
	switch {
	case monthAmount >= 50000:
		merchantType = models.MerchantTypeQuality // 优质: ≥5万
	case monthAmount >= 30000:
		merchantType = models.MerchantTypeMedium // 中等: 3-5万
	default:
		merchantType = models.MerchantTypeNormal // 普通: <3万
	}

	// 更新商户类型
	if err := s.merchantRepo.UpdateMerchantType(merchantID, merchantType); err != nil {
		return "", fmt.Errorf("更新商户类型失败: %w", err)
	}

	return merchantType, nil
}

// GetMerchantDetail 获取商户详情（包含关联信息）
func (s *MerchantService) GetMerchantDetail(id int64, agentID int64) (*MerchantDetailResponse, error) {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("商户不存在: %w", err)
	}

	// 验证商户归属
	if merchant.AgentID != agentID {
		return nil, errors.New("无权查看此商户")
	}

	// 获取代理商信息
	agent, _ := s.agentRepo.FindByID(merchant.AgentID)
	agentName := ""
	agentLevel := 0
	if agent != nil {
		agentName = agent.AgentName
		agentLevel = agent.Level
	}

	// 获取月度交易统计
	endTime := time.Now()
	startTime := time.Date(endTime.Year(), endTime.Month(), 1, 0, 0, 0, 0, endTime.Location())
	transStats, _ := s.merchantRepo.GetMerchantTransStats(id, &startTime, &endTime)

	monthAmount := int64(0)
	monthCount := int64(0)
	if transStats != nil {
		monthAmount = transStats.TotalAmount
		monthCount = transStats.TotalCount
	}

	// 获取终端数量
	// TODO: 添加终端仓储的FindByMerchantID方法后启用
	terminalCount := 0

	resp := &MerchantDetailResponse{
		ID:                merchant.ID,
		MerchantNo:        merchant.MerchantNo,
		MerchantName:      merchant.MerchantName,
		AgentID:           merchant.AgentID,
		AgentName:         agentName,
		AgentLevel:        agentLevel,
		ChannelID:         merchant.ChannelID,
		TerminalSN:        merchant.TerminalSN,
		Status:            merchant.Status,
		StatusName:        getMerchantStatusName(merchant.Status),
		ApproveStatus:     merchant.ApproveStatus,
		ApproveStatusName: getMerchantApproveStatusName(merchant.ApproveStatus),
		LegalName:         merchant.LegalName,
		LegalIDCard:       maskIDCard(merchant.LegalIDCard),
		MCC:               merchant.MCC,
		CreditRate:        merchant.CreditRate,
		DebitRate:         merchant.DebitRate,
		MerchantType:      merchant.MerchantType,
		MerchantTypeName:  getMerchantTypeName(merchant.MerchantType),
		IsDirect:          merchant.IsDirect,
		OwnerType:         getOwnerType(merchant.IsDirect),
		ActivatedAt:       merchant.ActivatedAt,
		RegisteredPhone:   maskPhone(merchant.RegisteredPhone),
		RegisterRemark:    merchant.RegisterRemark,
		MonthAmount:       monthAmount,
		MonthCount:        monthCount,
		TerminalCount:     terminalCount,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}

	return resp, nil
}

// GetExtendedStats 获取扩展统计
func (s *MerchantService) GetExtendedStats(agentID int64) (*repository.ExtendedMerchantStats, error) {
	return s.merchantRepo.GetExtendedStats(agentID)
}

// ==================== 辅助函数 ====================

func getMerchantStatusName(status int16) string {
	switch status {
	case models.MerchantStatusActive:
		return "正常"
	case models.MerchantStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

func getMerchantApproveStatusName(status int16) string {
	switch status {
	case models.MerchantApproveStatusPending:
		return "待审核"
	case models.MerchantApproveStatusApproved:
		return "已通过"
	case models.MerchantApproveStatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}

func getMerchantTypeName(merchantType string) string {
	switch merchantType {
	case models.MerchantTypeQuality:
		return "优质商户"
	case models.MerchantTypeMedium:
		return "中等商户"
	case models.MerchantTypeNormal:
		return "普通商户"
	case models.MerchantTypeWarning:
		return "预警商户"
	case models.MerchantTypeChurned:
		return "流失商户"
	default:
		return "未知"
	}
}

func getOwnerType(isDirect bool) string {
	if isDirect {
		return "direct"
	}
	return "team"
}

func maskIDCard(idCard string) string {
	// 先尝试解密（如果是加密的数据）
	if crypto.IsEncrypted(idCard) {
		decrypted, err := crypto.DecryptIDCard(idCard)
		if err == nil {
			idCard = decrypted
		}
	}
	return crypto.MaskIDCard(idCard)
}

func maskPhone(phone string) string {
	// 先尝试解密（如果是加密的数据）
	if crypto.IsEncrypted(phone) {
		decrypted, err := crypto.DecryptPhone(phone)
		if err == nil {
			phone = decrypted
		}
	}
	return crypto.MaskPhone(phone)
}

// getParentPolicyRate 获取上级代理商的政策费率
func (s *MerchantService) getParentPolicyRate(parentID int64, channelID int64) (*parentPolicyRate, error) {
	// 查询上级代理商的政策
	query := `
		SELECT credit_rate, debit_rate
		FROM agent_policies
		WHERE agent_id = ? AND channel_id = ?
		LIMIT 1
	`
	var result parentPolicyRate
	if err := s.merchantRepo.GetDB().Raw(query, parentID, channelID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// parentPolicyRate 上级政策费率
type parentPolicyRate struct {
	CreditRate string `json:"credit_rate"`
	DebitRate  string `json:"debit_rate"`
}

// parseRate 解析费率字符串为浮点数
func parseRate(rateStr string) float64 {
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return 0
	}
	return rate
}
