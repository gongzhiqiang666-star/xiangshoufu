package service

import (
	"encoding/json"
	"log"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"

	"github.com/gin-gonic/gin"
)

// AuditService 审计服务
type AuditService struct {
	auditRepo *repository.GormAuditLogRepository
}

// NewAuditService 创建审计服务
func NewAuditService(auditRepo *repository.GormAuditLogRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

// AuditContext 审计上下文
type AuditContext struct {
	UserID      int64
	Username    string
	AgentID     int64
	AgentName   string
	IP          string
	UserAgent   string
	RequestPath string
	RequestMethod string
}

// NewAuditContextFromGin 从Gin上下文创建审计上下文
func NewAuditContextFromGin(c *gin.Context) *AuditContext {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	agentID, _ := c.Get("agent_id")

	ctx := &AuditContext{
		IP:            c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		RequestPath:   c.Request.URL.Path,
		RequestMethod: c.Request.Method,
	}

	if uid, ok := userID.(int64); ok {
		ctx.UserID = uid
	}
	if uname, ok := username.(string); ok {
		ctx.Username = uname
	}
	if aid, ok := agentID.(int64); ok {
		ctx.AgentID = aid
	}

	return ctx
}

// LogLogin 记录登录
func (s *AuditService) LogLogin(ctx *AuditContext, success bool, failMsg string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeLogin,
		LogLevel:      models.AuditLogLevelInfo,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		Action:        "login",
		Description:   "用户登录",
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	if !success {
		auditLog.Result = 2
		auditLog.ErrorMsg = failMsg
		auditLog.LogLevel = models.AuditLogLevelWarning
	}

	s.saveLog(auditLog)
}

// LogLogout 记录登出
func (s *AuditService) LogLogout(ctx *AuditContext) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeLogout,
		LogLevel:      models.AuditLogLevelInfo,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		Action:        "logout",
		Description:   "用户登出",
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogPasswordChange 记录密码修改
func (s *AuditService) LogPasswordChange(ctx *AuditContext, success bool, failMsg string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypePasswordChange,
		LogLevel:      models.AuditLogLevelCritical,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		Action:        "change_password",
		Description:   "修改密码",
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	if !success {
		auditLog.Result = 2
		auditLog.ErrorMsg = failMsg
	}

	s.saveLog(auditLog)
}

// LogDataExport 记录数据导出
func (s *AuditService) LogDataExport(ctx *AuditContext, targetType string, description string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeDataExport,
		LogLevel:      models.AuditLogLevelWarning,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    targetType,
		Action:        "export",
		Description:   description,
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogWithdraw 记录提现操作
func (s *AuditService) LogWithdraw(ctx *AuditContext, withdrawID int64, amount int64, success bool, failMsg string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeWithdraw,
		LogLevel:      models.AuditLogLevelCritical,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    "withdraw",
		TargetID:      withdrawID,
		Action:        "withdraw",
		Description:   "提现申请",
		NewValue:      toJSON(map[string]interface{}{"amount": amount}),
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	if !success {
		auditLog.Result = 2
		auditLog.ErrorMsg = failMsg
	}

	s.saveLog(auditLog)
}

// LogAgentCreate 记录创建代理商
func (s *AuditService) LogAgentCreate(ctx *AuditContext, newAgentID int64, newAgentName string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeAgentCreate,
		LogLevel:      models.AuditLogLevelInfo,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    "agent",
		TargetID:      newAgentID,
		TargetName:    newAgentName,
		Action:        "create_agent",
		Description:   "创建代理商: " + newAgentName,
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogAgentStatusChange 记录代理商状态变更
func (s *AuditService) LogAgentStatusChange(ctx *AuditContext, targetAgentID int64, targetAgentName string, oldStatus, newStatus int16) {
	logType := models.AuditLogTypeAgentDisable
	action := "disable_agent"
	description := "禁用代理商"

	if newStatus == 1 {
		action = "enable_agent"
		description = "启用代理商"
	}

	auditLog := &models.AuditLog{
		LogType:       logType,
		LogLevel:      models.AuditLogLevelWarning,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    "agent",
		TargetID:      targetAgentID,
		TargetName:    targetAgentName,
		Action:        action,
		Description:   description + ": " + targetAgentName,
		OldValue:      toJSON(map[string]interface{}{"status": oldStatus}),
		NewValue:      toJSON(map[string]interface{}{"status": newStatus}),
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogRateChange 记录费率变更
func (s *AuditService) LogRateChange(ctx *AuditContext, targetType string, targetID int64, targetName string, oldRate, newRate string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeRateChange,
		LogLevel:      models.AuditLogLevelWarning,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    targetType,
		TargetID:      targetID,
		TargetName:    targetName,
		Action:        "change_rate",
		Description:   "费率变更",
		OldValue:      oldRate,
		NewValue:      newRate,
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogTerminalOperation 记录终端操作
func (s *AuditService) LogTerminalOperation(ctx *AuditContext, action string, terminalSN string, description string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeTerminalOp,
		LogLevel:      models.AuditLogLevelInfo,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    "terminal",
		TargetName:    terminalSN,
		Action:        action,
		Description:   description,
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogDeduction 记录代扣操作
func (s *AuditService) LogDeduction(ctx *AuditContext, deductionID int64, amount int64, action string, description string) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeDeduction,
		LogLevel:      models.AuditLogLevelWarning,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    "deduction",
		TargetID:      deductionID,
		Action:        action,
		Description:   description,
		NewValue:      toJSON(map[string]interface{}{"amount": amount}),
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogReward 记录奖励发放
func (s *AuditService) LogReward(ctx *AuditContext, targetAgentID int64, targetAgentName string, amount int64) {
	auditLog := &models.AuditLog{
		LogType:       models.AuditLogTypeReward,
		LogLevel:      models.AuditLogLevelInfo,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    "agent",
		TargetID:      targetAgentID,
		TargetName:    targetAgentName,
		Action:        "issue_reward",
		Description:   "发放奖励",
		NewValue:      toJSON(map[string]interface{}{"amount": amount}),
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	s.saveLog(auditLog)
}

// LogGeneric 记录通用操作
func (s *AuditService) LogGeneric(ctx *AuditContext, logType models.AuditLogType, level models.AuditLogLevel,
	targetType string, targetID int64, targetName string, action string, description string,
	oldValue, newValue interface{}, success bool, errMsg string) {

	auditLog := &models.AuditLog{
		LogType:       logType,
		LogLevel:      level,
		UserID:        ctx.UserID,
		Username:      ctx.Username,
		AgentID:       ctx.AgentID,
		TargetType:    targetType,
		TargetID:      targetID,
		TargetName:    targetName,
		Action:        action,
		Description:   description,
		OldValue:      toJSON(oldValue),
		NewValue:      toJSON(newValue),
		IP:            ctx.IP,
		UserAgent:     ctx.UserAgent,
		RequestPath:   ctx.RequestPath,
		RequestMethod: ctx.RequestMethod,
		Result:        1,
	}

	if !success {
		auditLog.Result = 2
		auditLog.ErrorMsg = errMsg
	}

	s.saveLog(auditLog)
}

// saveLog 保存日志（异步）
func (s *AuditService) saveLog(auditLog *models.AuditLog) {
	go func() {
		if err := s.auditRepo.Create(auditLog); err != nil {
			log.Printf("[AuditService] Failed to save audit log: %v", err)
		}
	}()
}

// QueryLogs 查询审计日志
func (s *AuditService) QueryLogs(params repository.AuditLogQueryParams) ([]*models.AuditLog, int64, error) {
	return s.auditRepo.FindByParams(params)
}

// toJSON 转换为JSON字符串
func toJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}
