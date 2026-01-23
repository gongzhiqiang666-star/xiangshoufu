package service

import (
	"context"
	"fmt"
	"log"

	"xiangshoufu/internal/channel"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// RateSyncService 费率同步服务
type RateSyncService struct {
	rateSyncRepo   repository.RateSyncLogRepository
	adapterFactory *channel.AdapterFactory
	messageService *MessageService
}

// NewRateSyncService 创建费率同步服务
func NewRateSyncService(
	rateSyncRepo repository.RateSyncLogRepository,
	adapterFactory *channel.AdapterFactory,
) *RateSyncService {
	return &RateSyncService{
		rateSyncRepo:   rateSyncRepo,
		adapterFactory: adapterFactory,
	}
}

// SetMessageService 设置消息服务（可选注入）
func (s *RateSyncService) SetMessageService(messageService *MessageService) {
	s.messageService = messageService
}

// RateUpdateParams 费率更新参数
type RateUpdateParams struct {
	MerchantID   int64
	MerchantNo   string
	TerminalSN   string
	ChannelCode  string
	AgentID      int64
	OldRates     *RateInfo
	NewRates     *RateInfo
}

// RateInfo 费率信息
type RateInfo struct {
	CreditRate   float64
	DebitRate    float64
	DebitCap     int64
	WechatRate   float64
	AlipayRate   float64
	UnionpayRate float64
}

// SyncResult 同步结果
type SyncResult struct {
	Success  bool
	LogID    int64
	TradeNo  string
	Message  string
}

// SyncRateToChannel 同步费率到通道（同步调用，立即返回结果）
func (s *RateSyncService) SyncRateToChannel(ctx context.Context, params *RateUpdateParams) (*SyncResult, error) {
	// 获取通道适配器
	adapter, err := s.adapterFactory.GetAdapter(params.ChannelCode)
	if err != nil {
		return nil, fmt.Errorf("获取通道适配器失败: %w", err)
	}

	// 检查通道是否支持费率实时更新
	if !adapter.SupportsRateUpdate() {
		return &SyncResult{
			Success: true,
			Message: "通道不支持费率实时更新，已跳过同步",
		}, nil
	}

	// 创建同步日志
	syncLog := &models.RateSyncLog{
		MerchantID:  params.MerchantID,
		MerchantNo:  params.MerchantNo,
		TerminalSN:  params.TerminalSN,
		ChannelCode: params.ChannelCode,
		AgentID:     params.AgentID,
		SyncStatus:  models.RateSyncStatusSyncing,
	}

	// 设置原费率
	if params.OldRates != nil {
		syncLog.OldCreditRate = &params.OldRates.CreditRate
		syncLog.OldDebitRate = &params.OldRates.DebitRate
		syncLog.OldDebitCap = &params.OldRates.DebitCap
		syncLog.OldWechatRate = &params.OldRates.WechatRate
		syncLog.OldAlipayRate = &params.OldRates.AlipayRate
		syncLog.OldUnionpayRate = &params.OldRates.UnionpayRate
	}

	// 设置新费率
	if params.NewRates != nil {
		syncLog.NewCreditRate = &params.NewRates.CreditRate
		syncLog.NewDebitRate = &params.NewRates.DebitRate
		syncLog.NewDebitCap = &params.NewRates.DebitCap
		syncLog.NewWechatRate = &params.NewRates.WechatRate
		syncLog.NewAlipayRate = &params.NewRates.AlipayRate
		syncLog.NewUnionpayRate = &params.NewRates.UnionpayRate
	}

	// 保存同步日志
	if err := s.rateSyncRepo.Create(ctx, syncLog); err != nil {
		log.Printf("[RateSyncService] 创建同步日志失败: %v", err)
	}

	// 构建请求
	req := &channel.RateUpdateRequest{
		MerchantNo: params.MerchantNo,
		TerminalSN: params.TerminalSN,
	}

	if params.NewRates != nil {
		req.CreditRate = params.NewRates.CreditRate
		req.DebitRate = params.NewRates.DebitRate
		req.DebitCap = params.NewRates.DebitCap
		req.WechatRate = params.NewRates.WechatRate
		req.AlipayRate = params.NewRates.AlipayRate
		req.UnionpayRate = params.NewRates.UnionpayRate
	}

	// 调用通道API
	resp, err := adapter.UpdateMerchantRate(req)
	if err != nil {
		// 记录失败日志
		syncLog.MarkFailed(err.Error())
		if updateErr := s.rateSyncRepo.Update(ctx, syncLog); updateErr != nil {
			log.Printf("[RateSyncService] 更新同步日志失败: %v", updateErr)
		}
		// 发送失败消息通知
		s.sendRateSyncNotification(params.AgentID, params.MerchantNo, false, err.Error())
		return &SyncResult{
			Success: false,
			LogID:   syncLog.ID,
			Message: fmt.Sprintf("通道同步失败: %v", err),
		}, nil
	}

	if !resp.Success {
		// 记录失败日志
		syncLog.MarkFailed(resp.Message)
		if updateErr := s.rateSyncRepo.Update(ctx, syncLog); updateErr != nil {
			log.Printf("[RateSyncService] 更新同步日志失败: %v", updateErr)
		}
		// 发送失败消息通知
		s.sendRateSyncNotification(params.AgentID, params.MerchantNo, false, resp.Message)
		return &SyncResult{
			Success: false,
			LogID:   syncLog.ID,
			Message: fmt.Sprintf("通道返回失败: %s", resp.Message),
		}, nil
	}

	// 同步成功
	syncLog.MarkSuccess(resp.TradeNo)
	if updateErr := s.rateSyncRepo.Update(ctx, syncLog); updateErr != nil {
		log.Printf("[RateSyncService] 更新同步日志失败: %v", updateErr)
	}
	// 发送成功消息通知
	s.sendRateSyncNotification(params.AgentID, params.MerchantNo, true, "")

	return &SyncResult{
		Success: true,
		LogID:   syncLog.ID,
		TradeNo: resp.TradeNo,
		Message: "费率同步成功",
	}, nil
}

// sendRateSyncNotification 发送费率同步结果通知
func (s *RateSyncService) sendRateSyncNotification(agentID int64, merchantNo string, success bool, errMsg string) {
	if s.messageService == nil {
		return
	}

	var title, content string
	var messageType int16

	if success {
		messageType = 8 // 交易通知类型
		title = "费率修改成功"
		content = fmt.Sprintf("商户[%s]费率已成功同步到支付通道", merchantNo)
	} else {
		messageType = 8
		title = "费率修改同步失败"
		content = fmt.Sprintf("商户[%s]费率同步到支付通道失败，请人工处理。失败原因：%s", merchantNo, errMsg)
	}

	msg := &NotificationMessage{
		AgentID:     agentID,
		MessageType: messageType,
		Title:       title,
		Content:     content,
		RelatedType: "rate_sync",
	}

	if err := s.messageService.SendNotification(msg); err != nil {
		log.Printf("[RateSyncService] 发送消息通知失败: %v", err)
	}
}

// GetSyncLogsByMerchant 获取商户的同步日志
func (s *RateSyncService) GetSyncLogsByMerchant(ctx context.Context, merchantID int64, page, pageSize int) ([]*models.RateSyncLog, int64, error) {
	return s.rateSyncRepo.GetByMerchantID(ctx, merchantID, page, pageSize)
}

// GetSyncLogByID 根据ID获取同步日志
func (s *RateSyncService) GetSyncLogByID(ctx context.Context, id int64) (*models.RateSyncLog, error) {
	return s.rateSyncRepo.GetByID(ctx, id)
}
