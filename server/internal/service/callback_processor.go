package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"xiangshoufu/internal/async"
	"xiangshoufu/internal/channel"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// CallbackProcessor 回调数据处理服务
type CallbackProcessor struct {
	factory         *channel.AdapterFactory
	callbackRepo    repository.RawCallbackRepository
	transactionRepo repository.TransactionRepository
	deviceFeeRepo   repository.DeviceFeeRepository
	rateChangeRepo  repository.RateChangeRepository
	profitService   *ProfitService
	queue           async.MessageQueue
}

// NewCallbackProcessor 创建回调处理服务
func NewCallbackProcessor(
	factory *channel.AdapterFactory,
	callbackRepo repository.RawCallbackRepository,
	transactionRepo repository.TransactionRepository,
	deviceFeeRepo repository.DeviceFeeRepository,
	rateChangeRepo repository.RateChangeRepository,
	profitService *ProfitService,
	queue async.MessageQueue,
) *CallbackProcessor {
	return &CallbackProcessor{
		factory:         factory,
		callbackRepo:    callbackRepo,
		transactionRepo: transactionRepo,
		deviceFeeRepo:   deviceFeeRepo,
		rateChangeRepo:  rateChangeRepo,
		profitService:   profitService,
		queue:           queue,
	}
}

// ProcessMessage 处理队列消息
func (p *CallbackProcessor) ProcessMessage(msgBytes []byte) error {
	var msg QueueMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return fmt.Errorf("unmarshal message failed: %w", err)
	}

	return p.ProcessCallback(msg.CallbackLogID, msg.ChannelCode, msg.ActionType, msg.RawBody)
}

// QueueMessage 队列消息
type QueueMessage struct {
	CallbackLogID int64  `json:"callback_log_id"`
	ChannelCode   string `json:"channel_code"`
	ActionType    string `json:"action_type"`
	RawBody       []byte `json:"raw_body"`
}

// ProcessCallback 处理单个回调
func (p *CallbackProcessor) ProcessCallback(logID int64, channelCode, actionType string, rawBody []byte) error {
	adapter, err := p.factory.GetAdapter(channelCode)
	if err != nil {
		return p.markFailed(logID, fmt.Sprintf("adapter not found: %s", channelCode))
	}

	var processErr error

	switch channel.ActionType(actionType) {
	case channel.ActionTransaction:
		processErr = p.processTransaction(adapter, rawBody, logID)
	case channel.ActionDeviceFee:
		processErr = p.processDeviceFee(adapter, rawBody)
	case channel.ActionRateChange:
		processErr = p.processRateChange(adapter, rawBody)
	case channel.ActionMerchantIncome:
		processErr = p.processMerchantIncome(adapter, rawBody)
	case channel.ActionTerminalBind:
		processErr = p.processTerminalBind(adapter, rawBody)
	default:
		processErr = fmt.Errorf("unknown action type: %s", actionType)
	}

	if processErr != nil {
		return p.markFailed(logID, processErr.Error())
	}

	return p.markSuccess(logID)
}

// processTransaction 处理交易回调
func (p *CallbackProcessor) processTransaction(adapter channel.ChannelAdapter, rawBody []byte, logID int64) error {
	// 1. 解析交易数据
	unified, err := adapter.ParseTransaction(rawBody)
	if err != nil {
		return fmt.Errorf("parse transaction failed: %w", err)
	}

	// 2. 转换为交易模型
	agentID, _ := strconv.ParseInt(unified.AgentID, 10, 64)
	extDataBytes, _ := json.Marshal(unified.ExtData)

	tx := &repository.Transaction{
		OrderNo:     unified.OrderNo,
		ChannelCode: unified.ChannelCode,
		TerminalSN:  unified.TerminalSN,
		AgentID:     agentID,
		TradeType:   1, // 消费
		PayType:     mapPayType(unified.CardType),
		CardType:    mapCardTypeToInt(unified.CardType),
		Amount:      unified.Amount,
		Rate:        unified.FeeRate,
		D0Fee:       unified.D0Fee,
		HighRate:    unified.HighRate,
		CardNo:      unified.CardNo,
		TradeTime:   unified.TransTime,
		ReceivedAt:  time.Now(),
		ExtData:     string(extDataBytes),
	}

	// 3. 检查是否已存在（幂等）
	existing, _ := p.transactionRepo.FindByOrderNo(tx.OrderNo)
	if existing != nil {
		log.Printf("[CallbackProcessor] Transaction already exists: %s", tx.OrderNo)
		return nil
	}

	// 4. 保存交易
	if err := p.transactionRepo.Create(tx); err != nil {
		return fmt.Errorf("save transaction failed: %w", err)
	}

	// 5. 发送到分润计算队列
	profitMsg := &ProfitMessage{
		TransactionID: tx.ID,
		OrderNo:       tx.OrderNo,
	}
	msgBytes, _ := json.Marshal(profitMsg)
	if err := p.queue.Publish(async.TopicProfitCalc, msgBytes); err != nil {
		log.Printf("[CallbackProcessor] Publish profit calc failed: %v", err)
		// 不返回错误，定时任务会兜底
	}

	return nil
}

// ProfitMessage 分润计算消息
type ProfitMessage struct {
	TransactionID int64  `json:"transaction_id"`
	OrderNo       string `json:"order_no"`
}

// processDeviceFee 处理流量费/服务费回调
func (p *CallbackProcessor) processDeviceFee(adapter channel.ChannelAdapter, rawBody []byte) error {
	unified, err := adapter.ParseDeviceFee(rawBody)
	if err != nil {
		return fmt.Errorf("parse device fee failed: %w", err)
	}

	agentID, _ := strconv.ParseInt(unified.AgentID, 10, 64)

	fee := &models.DeviceFee{
		ChannelCode:  unified.ChannelCode,
		TerminalSN:   unified.TerminalSN,
		MerchantNo:   unified.MerchantNo,
		AgentID:      agentID,
		OrderNo:      unified.OrderNo,
		FeeType:      int16(unified.FeeType),
		FeeAmount:    unified.FeeAmount,
		ChargingTime: unified.ChargingTime,
		ReceivedAt:   time.Now(),
		BrandCode:    unified.BrandCode,
	}

	// 检查是否已存在
	existing, _ := p.deviceFeeRepo.FindByOrderNo(fee.OrderNo)
	if existing != nil {
		return nil
	}

	return p.deviceFeeRepo.Create(fee)
}

// processRateChange 处理费率变更回调
func (p *CallbackProcessor) processRateChange(adapter channel.ChannelAdapter, rawBody []byte) error {
	unified, err := adapter.ParseRateChange(rawBody)
	if err != nil {
		return fmt.Errorf("parse rate change failed: %w", err)
	}

	agentID, _ := strconv.ParseInt(unified.AgentID, 10, 64)

	change := &models.RateChange{
		ChannelCode:          unified.ChannelCode,
		TerminalSN:           unified.TerminalSN,
		MerchantNo:           unified.MerchantNo,
		AgentID:              agentID,
		CreditRate:           unified.CreditRate,
		CreditExtraRate:      unified.CreditExtraRate,
		DebitRate:            unified.DebitRate,
		AlipayRate:           unified.AlipayRate,
		WechatRate:           unified.WechatRate,
		UnionpayRate:         unified.UnionpayRate,
		CreditAdditionRate:   unified.CreditAdditionRate,
		UnionpayAdditionRate: unified.UnionpayAdditionRate,
		AlipayAdditionRate:   unified.AlipayAdditionRate,
		WechatAdditionRate:   unified.WechatAdditionRate,
		ReceivedAt:           time.Now(),
		BrandCode:            unified.BrandCode,
	}

	return p.rateChangeRepo.Create(change)
}

// processMerchantIncome 处理商户入网回调
func (p *CallbackProcessor) processMerchantIncome(adapter channel.ChannelAdapter, rawBody []byte) error {
	unified, err := adapter.ParseMerchantIncome(rawBody)
	if err != nil {
		return fmt.Errorf("parse merchant income failed: %w", err)
	}

	// TODO: 更新商户表状态
	log.Printf("[CallbackProcessor] Merchant income: %s, status: %d", unified.MerchantNo, unified.ApproveStatus)
	return nil
}

// processTerminalBind 处理终端绑定回调
func (p *CallbackProcessor) processTerminalBind(adapter channel.ChannelAdapter, rawBody []byte) error {
	unified, err := adapter.ParseTerminalBind(rawBody)
	if err != nil {
		return fmt.Errorf("parse terminal bind failed: %w", err)
	}

	// TODO: 更新终端表状态
	log.Printf("[CallbackProcessor] Terminal bind: %s, status: %d", unified.TerminalSN, unified.BindStatus)
	return nil
}

// markSuccess 标记处理成功
func (p *CallbackProcessor) markSuccess(logID int64) error {
	return p.callbackRepo.UpdateStatus(logID, models.ProcessStatusSuccess, "")
}

// markFailed 标记处理失败
func (p *CallbackProcessor) markFailed(logID int64, errMsg string) error {
	p.callbackRepo.IncrementRetryCount(logID)
	return p.callbackRepo.UpdateStatus(logID, models.ProcessStatusFailed, errMsg)
}

// mapPayType 映射支付类型
func mapPayType(cardType channel.CardType) int16 {
	switch cardType {
	case channel.CardTypeDebit, channel.CardTypeCredit:
		return 1 // 刷卡
	case channel.CardTypeWechat:
		return 2 // 微信
	case channel.CardTypeAlipay:
		return 3 // 支付宝
	case channel.CardTypeUnionpay:
		return 4 // 云闪付
	default:
		return 1
	}
}

// mapCardTypeToInt 映射卡类型到整数
func mapCardTypeToInt(cardType channel.CardType) int16 {
	switch cardType {
	case channel.CardTypeDebit:
		return 1
	case channel.CardTypeCredit:
		return 2
	default:
		return 0
	}
}
