package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
	"xiangshoufu/internal/service"
)

// SystemHandler 系统管理处理器
type SystemHandler struct {
	userRepo     *repository.GormUserRepository
	auditLogRepo *repository.GormAuditLogRepository
	authService  *service.AuthService
}

// NewSystemHandler 创建系统管理处理器
func NewSystemHandler(
	userRepo *repository.GormUserRepository,
	auditLogRepo *repository.GormAuditLogRepository,
) *SystemHandler {
	return &SystemHandler{
		userRepo:     userRepo,
		auditLogRepo: auditLogRepo,
	}
}

// SetAuthService 设置认证服务
func (h *SystemHandler) SetAuthService(authService *service.AuthService) {
	h.authService = authService
}

// ========== 用户管理 ==========

// UserListRequest 用户列表请求
type UserListRequest struct {
	Username string `form:"username"`
	RoleType *int16 `form:"role_type"`
	Status   *int16 `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	AgentID     int64      `json:"agent_id"`
	RoleType    int16      `json:"role_type"`
	RoleTypeName string    `json:"role_type_name"`
	Status      int16      `json:"status"`
	StatusName  string     `json:"status_name"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ListUsers 用户列表
// @Summary 获取用户列表
// @Tags 系统管理
// @Produce json
// @Param username query string false "用户名"
// @Param role_type query int false "角色类型"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/system/users [get]
func (h *SystemHandler) ListUsers(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	params := repository.UserQueryParams{
		Username: req.Username,
		RoleType: req.RoleType,
		Status:   req.Status,
		Limit:    req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
	}

	users, total, err := h.userRepo.FindByParams(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	var list []UserListResponse
	for _, user := range users {
		item := UserListResponse{
			ID:          user.ID,
			Username:    user.Username,
			AgentID:     user.AgentID,
			RoleType:    user.RoleType,
			RoleTypeName: getRoleTypeName(user.RoleType),
			Status:      user.Status,
			StatusName:  getUserStatusName(user.Status),
			LastLoginAt: user.LastLoginAt,
			LastLoginIP: user.LastLoginIP,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":  list,
			"total": total,
			"page":  req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Tags 系统管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/system/users/{id} [get]
func (h *SystemHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := h.userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": UserListResponse{
			ID:          user.ID,
			Username:    user.Username,
			AgentID:     user.AgentID,
			RoleType:    user.RoleType,
			RoleTypeName: getRoleTypeName(user.RoleType),
			Status:      user.Status,
			StatusName:  getUserStatusName(user.Status),
			LastLoginAt: user.LastLoginAt,
			LastLoginIP: user.LastLoginIP,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	})
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status int16 `json:"status" binding:"required,oneof=1 2"`
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param body body UpdateUserStatusRequest true "状态"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/system/users/{id}/status [put]
func (h *SystemHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 检查用户是否存在
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新状态
	if err := h.userRepo.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// ========== 操作日志 ==========

// LogListRequest 日志列表请求
type LogListRequest struct {
	LogType   *int16 `form:"log_type"`
	LogLevel  *int16 `form:"log_level"`
	UserID    *int64 `form:"user_id"`
	AgentID   *int64 `form:"agent_id"`
	Keyword   string `form:"keyword"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// LogListResponse 日志列表响应
type LogListResponse struct {
	ID            int64     `json:"id"`
	LogType       int16     `json:"log_type"`
	LogTypeName   string    `json:"log_type_name"`
	LogLevel      int16     `json:"log_level"`
	LogLevelName  string    `json:"log_level_name"`
	UserID        int64     `json:"user_id"`
	Username      string    `json:"username"`
	AgentID       int64     `json:"agent_id"`
	AgentName     string    `json:"agent_name"`
	TargetType    string    `json:"target_type"`
	TargetID      int64     `json:"target_id"`
	TargetName    string    `json:"target_name"`
	Action        string    `json:"action"`
	Description   string    `json:"description"`
	IP            string    `json:"ip"`
	Result        int16     `json:"result"`
	ResultName    string    `json:"result_name"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListLogs 操作日志列表
// @Summary 获取操作日志列表
// @Tags 系统管理
// @Produce json
// @Param log_type query int false "日志类型"
// @Param log_level query int false "日志级别"
// @Param user_id query int false "用户ID"
// @Param agent_id query int false "代理商ID"
// @Param keyword query string false "关键词"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/system/logs [get]
func (h *SystemHandler) ListLogs(c *gin.Context) {
	var req LogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	params := repository.AuditLogQueryParams{
		Keyword: req.Keyword,
		Limit:   req.PageSize,
		Offset:  (req.Page - 1) * req.PageSize,
	}

	if req.LogType != nil {
		logType := models.AuditLogType(*req.LogType)
		params.LogType = &logType
	}
	if req.LogLevel != nil {
		logLevel := models.AuditLogLevel(*req.LogLevel)
		params.LogLevel = &logLevel
	}
	if req.UserID != nil {
		params.UserID = req.UserID
	}
	if req.AgentID != nil {
		params.AgentID = req.AgentID
	}
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			params.StartTime = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endOfDay := t.Add(24*time.Hour - time.Second)
			params.EndTime = &endOfDay
		}
	}

	logs, total, err := h.auditLogRepo.FindByParams(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询日志失败"})
		return
	}

	var list []LogListResponse
	for _, log := range logs {
		item := LogListResponse{
			ID:           log.ID,
			LogType:      int16(log.LogType),
			LogTypeName:  models.GetLogTypeName(log.LogType),
			LogLevel:     int16(log.LogLevel),
			LogLevelName: models.GetLogLevelName(log.LogLevel),
			UserID:       log.UserID,
			Username:     log.Username,
			AgentID:      log.AgentID,
			AgentName:    log.AgentName,
			TargetType:   log.TargetType,
			TargetID:     log.TargetID,
			TargetName:   log.TargetName,
			Action:       log.Action,
			Description:  log.Description,
			IP:           log.IP,
			Result:       log.Result,
			ResultName:   getResultName(log.Result),
			CreatedAt:    log.CreatedAt,
		}
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetLog 获取日志详情
// @Summary 获取日志详情
// @Tags 系统管理
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/system/logs/{id} [get]
func (h *SystemHandler) GetLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"})
		return
	}

	log, err := h.auditLogRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"id":             log.ID,
			"log_type":       log.LogType,
			"log_type_name":  models.GetLogTypeName(log.LogType),
			"log_level":      log.LogLevel,
			"log_level_name": models.GetLogLevelName(log.LogLevel),
			"user_id":        log.UserID,
			"username":       log.Username,
			"agent_id":       log.AgentID,
			"agent_name":     log.AgentName,
			"target_type":    log.TargetType,
			"target_id":      log.TargetID,
			"target_name":    log.TargetName,
			"action":         log.Action,
			"description":    log.Description,
			"old_value":      log.OldValue,
			"new_value":      log.NewValue,
			"ip":             log.IP,
			"user_agent":     log.UserAgent,
			"request_path":   log.RequestPath,
			"request_method": log.RequestMethod,
			"result":         log.Result,
			"result_name":    getResultName(log.Result),
			"error_msg":      log.ErrorMsg,
			"created_at":     log.CreatedAt,
		},
	})
}

// ========== 辅助函数 ==========

func getRoleTypeName(roleType int16) string {
	switch roleType {
	case models.UserRoleTypeNormal:
		return "普通用户"
	case models.UserRoleTypeAdmin:
		return "管理员"
	default:
		return "未知"
	}
}

func getUserStatusName(status int16) string {
	switch status {
	case models.UserStatusActive:
		return "正常"
	case models.UserStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

func getResultName(result int16) string {
	switch result {
	case 1:
		return "成功"
	case 2:
		return "失败"
	default:
		return "未知"
	}
}

// RegisterSystemRoutes 注册系统管理路由
func RegisterSystemRoutes(rg *gin.RouterGroup, h *SystemHandler, authService *service.AuthService) {
	systemGroup := rg.Group("/system")
	systemGroup.Use(middleware.AuthMiddleware(authService))
	{
		// 用户管理
		systemGroup.GET("/users", h.ListUsers)
		systemGroup.GET("/users/:id", h.GetUser)
		systemGroup.PUT("/users/:id/status", h.UpdateUserStatus)

		// 操作日志
		systemGroup.GET("/logs", h.ListLogs)
		systemGroup.GET("/logs/:id", h.GetLog)
	}
}
