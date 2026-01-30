package service

import (
	"fmt"
	"log"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"

	"gorm.io/gorm"
)

// WalletAdjustmentService 钱包调账服务
type WalletAdjustmentService struct {
	adjustmentRepo *repository.GormWalletAdjustmentRepository
	walletRepo     *repository.GormWalletRepository
	walletLogRepo  *repository.GormWalletLogRepository
	agentRepo      *repository.GormAgentRepository
}

// NewWalletAdjustmentService 创建钱包调账服务
func NewWalletAdjustmentService(
	adjustmentRepo *repository.GormWalletAdjustmentRepository,
	walletRepo *repository.GormWalletRepository,
	walletLogRepo *repository.GormWalletLogRepository,
	agentRepo *repository.GormAgentRepository,
) *WalletAdjustmentService {
	return &WalletAdjustmentService{
		adjustmentRepo: adjustmentRepo,
		walletRepo:     walletRepo,
		walletLogRepo:  walletLogRepo,
		agentRepo:      agentRepo,
	}
}

// 调账流水类型
const (
	WalletLogTypeAdjustmentIn  int16 = 11 // 调账充入
	WalletLogTypeAdjustmentOut int16 = 12 // 调账扣减
)

// CreateAdjustmentRequest 创建调账请求
type CreateAdjustmentRequest struct {
	AgentID      int64  `json:"agent_id" binding:"required"`
	WalletType   int16  `json:"wallet_type" binding:"required"` // 1分润 2服务费 3奖励 4充值 5沉淀
	ChannelID    int64  `json:"channel_id"`                     // 0表示不区分通道
	Amount       int64  `json:"amount" binding:"required"`      // 正数充入，负数扣减
	Reason       string `json:"reason" binding:"required"`
	OperatorID   int64  `json:"operator_id" binding:"required"`
	OperatorName string `json:"operator_name"`
}

// CreateAdjustment 创建调账
func (s *WalletAdjustmentService) CreateAdjustment(req *CreateAdjustmentRequest) (*AdjustmentInfo, error) {
	// 1. 验证代理商存在
	agent, err := s.agentRepo.FindByIDFull(req.AgentID)
	if err != nil || agent == nil {
		return nil, fmt.Errorf("代理商不存在")
	}

	// 2. 获取钱包
	wallet, err := s.walletRepo.FindByAgentAndType(req.AgentID, req.ChannelID, req.WalletType)
	if err != nil || wallet == nil {
		return nil, fmt.Errorf("钱包不存在，请确认代理商和钱包类型")
	}

	// 3. 扣减时检查余额
	if req.Amount < 0 {
		availableBalance := wallet.Balance - wallet.FrozenAmount
		if availableBalance < -req.Amount {
			return nil, fmt.Errorf("可用余额不足，当前可用余额%.2f元", float64(availableBalance)/100)
		}
	}

	// 4. 生成调账单号
	adjustmentNo := fmt.Sprintf("ADJ%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)

	now := time.Now()
	balanceBefore := wallet.Balance
	balanceAfter := wallet.Balance + req.Amount

	// 5. 使用事务执行调账
	db := s.adjustmentRepo.GetDB()
	var adjustment *models.WalletAdjustment
	var walletLogID int64

	err = db.Transaction(func(tx *gorm.DB) error {
		// 5.1 创建调账记录
		adjustment = &models.WalletAdjustment{
			AdjustmentNo:  adjustmentNo,
			AgentID:       req.AgentID,
			WalletID:      wallet.ID,
			WalletType:    req.WalletType,
			ChannelID:     req.ChannelID,
			Amount:        req.Amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Reason:        req.Reason,
			OperatorID:    req.OperatorID,
			OperatorName:  req.OperatorName,
			Status:        models.AdjustmentStatusEffective,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := tx.Create(adjustment).Error; err != nil {
			return fmt.Errorf("创建调账记录失败: %w", err)
		}

		// 5.2 更新钱包余额
		if err := tx.Model(&repository.Wallet{}).
			Where("id = ?", wallet.ID).
			Updates(map[string]interface{}{
				"balance":      gorm.Expr("balance + ?", req.Amount),
				"total_income": gorm.Expr("total_income + ?", req.Amount),
				"version":      gorm.Expr("version + 1"),
			}).Error; err != nil {
			return fmt.Errorf("更新钱包余额失败: %w", err)
		}

		// 5.3 创建钱包流水
		logType := WalletLogTypeAdjustmentIn
		if req.Amount < 0 {
			logType = WalletLogTypeAdjustmentOut
		}

		walletLog := &repository.WalletLog{
			WalletID:      wallet.ID,
			AgentID:       req.AgentID,
			WalletType:    req.WalletType,
			LogType:       logType,
			Amount:        req.Amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			RefType:       "wallet_adjustment",
			RefID:         adjustment.ID,
			Remark:        fmt.Sprintf("手动调账: %s", req.Reason),
			CreatedAt:     now,
		}

		if err := tx.Create(walletLog).Error; err != nil {
			return fmt.Errorf("创建钱包流水失败: %w", err)
		}
		walletLogID = walletLog.ID

		// 5.4 更新调账记录的流水ID
		if err := tx.Model(&models.WalletAdjustment{}).
			Where("id = ?", adjustment.ID).
			Update("wallet_log_id", walletLogID).Error; err != nil {
			return fmt.Errorf("更新调账记录失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[WalletAdjustmentService] Created adjustment: %s, agent=%d, wallet_type=%d, amount=%d",
		adjustmentNo, req.AgentID, req.WalletType, req.Amount)

	return s.toAdjustmentInfo(adjustment, agent.AgentName), nil
}

// AdjustmentInfo 调账信息
type AdjustmentInfo struct {
	ID               int64      `json:"id"`
	AdjustmentNo     string     `json:"adjustment_no"`
	AgentID          int64      `json:"agent_id"`
	AgentName        string     `json:"agent_name"`
	WalletID         int64      `json:"wallet_id"`
	WalletType       int16      `json:"wallet_type"`
	WalletTypeName   string     `json:"wallet_type_name"`
	ChannelID        int64      `json:"channel_id"`
	Amount           int64      `json:"amount"`
	AmountYuan       float64    `json:"amount_yuan"`
	AdjustmentType   string     `json:"adjustment_type"` // 充入/扣减
	BalanceBefore    int64      `json:"balance_before"`
	BalanceBeforeYuan float64   `json:"balance_before_yuan"`
	BalanceAfter     int64      `json:"balance_after"`
	BalanceAfterYuan float64    `json:"balance_after_yuan"`
	Reason           string     `json:"reason"`
	OperatorID       int64      `json:"operator_id"`
	OperatorName     string     `json:"operator_name"`
	Status           int16      `json:"status"`
	StatusName       string     `json:"status_name"`
	CreatedAt        time.Time  `json:"created_at"`
}

func (s *WalletAdjustmentService) toAdjustmentInfo(adj *models.WalletAdjustment, agentName string) *AdjustmentInfo {
	return &AdjustmentInfo{
		ID:                adj.ID,
		AdjustmentNo:      adj.AdjustmentNo,
		AgentID:           adj.AgentID,
		AgentName:         agentName,
		WalletID:          adj.WalletID,
		WalletType:        adj.WalletType,
		WalletTypeName:    models.GetWalletTypeName(adj.WalletType),
		ChannelID:         adj.ChannelID,
		Amount:            adj.Amount,
		AmountYuan:        float64(adj.Amount) / 100,
		AdjustmentType:    models.GetAdjustmentTypeName(adj.Amount),
		BalanceBefore:     adj.BalanceBefore,
		BalanceBeforeYuan: float64(adj.BalanceBefore) / 100,
		BalanceAfter:      adj.BalanceAfter,
		BalanceAfterYuan:  float64(adj.BalanceAfter) / 100,
		Reason:            adj.Reason,
		OperatorID:        adj.OperatorID,
		OperatorName:      adj.OperatorName,
		Status:            adj.Status,
		StatusName:        models.GetAdjustmentStatusName(adj.Status),
		CreatedAt:         adj.CreatedAt,
	}
}

// AdjustmentListParams 调账列表查询参数
type AdjustmentListParams struct {
	AgentID    int64
	WalletType *int16
	ChannelID  *int64
	StartTime  *time.Time
	EndTime    *time.Time
	Page       int
	PageSize   int
}

// GetAdjustmentList 获取调账列表
func (s *WalletAdjustmentService) GetAdjustmentList(params *AdjustmentListParams) ([]*AdjustmentInfo, int64, error) {
	offset := (params.Page - 1) * params.PageSize

	repoParams := &repository.WalletAdjustmentQueryParams{
		AgentID:    params.AgentID,
		WalletType: params.WalletType,
		ChannelID:  params.ChannelID,
		StartTime:  params.StartTime,
		EndTime:    params.EndTime,
		Limit:      params.PageSize,
		Offset:     offset,
	}

	adjustments, total, err := s.adjustmentRepo.List(repoParams)
	if err != nil {
		return nil, 0, fmt.Errorf("查询调账列表失败: %w", err)
	}

	// 获取代理商名称映射
	agentNames := make(map[int64]string)
	for _, adj := range adjustments {
		if _, ok := agentNames[adj.AgentID]; !ok {
			if agent, _ := s.agentRepo.FindByIDFull(adj.AgentID); agent != nil {
				agentNames[adj.AgentID] = agent.AgentName
			}
		}
	}

	list := make([]*AdjustmentInfo, 0, len(adjustments))
	for _, adj := range adjustments {
		list = append(list, s.toAdjustmentInfo(adj, agentNames[adj.AgentID]))
	}

	return list, total, nil
}

// GetAdjustmentDetail 获取调账详情
func (s *WalletAdjustmentService) GetAdjustmentDetail(id int64) (*AdjustmentInfo, error) {
	adjustment, err := s.adjustmentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("查询调账记录失败: %w", err)
	}
	if adjustment == nil {
		return nil, fmt.Errorf("调账记录不存在")
	}

	agentName := ""
	if agent, _ := s.agentRepo.FindByIDFull(adjustment.AgentID); agent != nil {
		agentName = agent.AgentName
	}

	return s.toAdjustmentInfo(adjustment, agentName), nil
}
