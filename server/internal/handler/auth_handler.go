package handler

import (
	"net/http"

	"xiangshoufu/internal/middleware"
	"xiangshoufu/internal/service"
	"xiangshoufu/pkg/crypto"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService  *service.AuthService
	agentService *service.AgentService
	auditService *service.AuditService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SetAgentService 设置代理商服务（延迟注入，避免循环依赖）
func (h *AuthHandler) SetAgentService(agentService *service.AgentService) {
	h.agentService = agentService
}

// SetAuditService 设置审计服务
func (h *AuthHandler) SetAuditService(auditService *service.AuditService) {
	h.auditService = auditService
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Encrypted bool   `json:"encrypted"` // 密码是否已加密
}

// Login 登录
// @Summary 用户登录
// @Description 使用用户名密码登录，返回访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 如果密码已加密，先解密
	password := req.Password
	if req.Encrypted {
		decrypted, err := crypto.DecryptPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "密码解密失败，请刷新页面重试",
			})
			return
		}
		password = decrypted
	}

	loginReq := &service.LoginRequest{
		Username:  req.Username,
		Password:  password,
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	resp, err := h.authService.Login(loginReq)
	if err != nil {
		// 记录登录失败审计日志
		if h.auditService != nil {
			auditCtx := &service.AuditContext{
				Username:      req.Username,
				IP:            c.ClientIP(),
				UserAgent:     c.GetHeader("User-Agent"),
				RequestPath:   c.Request.URL.Path,
				RequestMethod: c.Request.Method,
			}
			h.auditService.LogLogin(auditCtx, false, err.Error())
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	// 记录登录成功审计日志
	if h.auditService != nil {
		agentID := int64(0)
		if resp.Agent != nil {
			agentID = resp.Agent.ID
		}
		auditCtx := &service.AuditContext{
			UserID:        resp.User.ID,
			Username:      resp.User.Username,
			AgentID:       agentID,
			IP:            c.ClientIP(),
			UserAgent:     c.GetHeader("User-Agent"),
			RequestPath:   c.Request.URL.Path,
			RequestMethod: c.Request.Method,
		}
		h.auditService.LogLogin(auditCtx, true, "")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    resp,
	})
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Logout 登出
// @Summary 用户登出
// @Description 登出并使刷新令牌失效
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LogoutRequest false "登出请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	c.ShouldBindJSON(&req)

	if req.RefreshToken != "" {
		h.authService.Logout(req.RefreshToken)
	}

	// 记录登出审计日志
	if h.auditService != nil {
		auditCtx := service.NewAuditContextFromGin(c)
		h.auditService.LogLogout(auditCtx)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登出成功",
	})
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "刷新请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	resp, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "刷新成功",
		"data":    resp,
	})
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	resp, err := h.authService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    resp,
	})
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		// 记录密码修改失败审计日志
		if h.auditService != nil {
			auditCtx := service.NewAuditContextFromGin(c)
			h.auditService.LogPasswordChange(auditCtx, false, err.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 记录密码修改成功审计日志
	if h.auditService != nil {
		auditCtx := service.NewAuditContextFromGin(c)
		h.auditService.LogPasswordChange(auditCtx, true, "")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码修改成功，请重新登录",
	})
}

// RegisterRequest 公开注册请求
type RegisterRequest struct {
	InviteCode   string `json:"invite_code" binding:"required"`
	AgentName    string `json:"agent_name" binding:"required"`
	ContactName  string `json:"contact_name" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	IDCardNo     string `json:"id_card_no"`
	Password     string `json:"password" binding:"required,min=6"`
}

// Register 通过邀请码注册代理商
// @Summary 通过邀请码注册
// @Description 使用邀请码自助注册成为代理商，注册后状态为待审核
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	if h.agentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 调用代理商服务进行注册
	serviceReq := &service.RegisterByInviteCodeRequest{
		InviteCode:   req.InviteCode,
		AgentName:    req.AgentName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		IDCardNo:     req.IDCardNo,
		Password:     req.Password,
	}

	agent, err := h.agentService.RegisterByInviteCode(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功，请等待上级代理商审核",
		"data":    agent,
	})
}

// GetPublicKey 获取RSA公钥
// @Summary 获取RSA公钥
// @Description 获取用于密码加密的RSA公钥
// @Tags 认证
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/public-key [get]
func (h *AuthHandler) GetPublicKey(c *gin.Context) {
	publicKey, err := crypto.GetPublicKeyPEM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取公钥失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"public_key": publicKey,
		},
	})
}

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(r *gin.RouterGroup, h *AuthHandler, authService *service.AuthService) {
	auth := r.Group("/auth")
	{
		// 公开接口（无需登录）
		auth.GET("/public-key", h.GetPublicKey) // RSA公钥获取
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/register", h.Register) // 邀请码注册接口

		// 需要登录的接口
		authRequired := auth.Group("")
		authRequired.Use(middleware.AuthMiddleware(authService))
		{
			authRequired.POST("/logout", h.Logout)
			authRequired.GET("/profile", h.GetProfile)
			authRequired.PUT("/password", h.ChangePassword)
		}
	}
}
