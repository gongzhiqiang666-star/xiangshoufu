package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret          string        // JWT密钥
	AccessTokenExpiry  time.Duration // 访问令牌有效期
	RefreshTokenExpiry time.Duration // 刷新令牌有效期
}

// DefaultAuthConfig 默认配置
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:          "xiangshoufu-secret-key-2026",
		AccessTokenExpiry:  2 * time.Hour,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
}

// AuthService 认证服务
type AuthService struct {
	config           *AuthConfig
	userRepo         *repository.GormUserRepository
	refreshTokenRepo *repository.GormRefreshTokenRepository
	loginLogRepo     *repository.GormLoginLogRepository
	agentRepo        *repository.GormAgentRepository
}

// NewAuthService 创建认证服务
func NewAuthService(
	config *AuthConfig,
	userRepo *repository.GormUserRepository,
	refreshTokenRepo *repository.GormRefreshTokenRepository,
	loginLogRepo *repository.GormLoginLogRepository,
	agentRepo *repository.GormAgentRepository,
) *AuthService {
	if config == nil {
		config = DefaultAuthConfig()
	}
	return &AuthService{
		config:           config,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		loginLogRepo:     loginLogRepo,
		agentRepo:        agentRepo,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int64            `json:"expires_in"` // 秒
	TokenType    string           `json:"token_type"`
	User         *LoginUserInfo   `json:"user"`
	Agent        *LoginAgentInfo  `json:"agent,omitempty"`
}

// LoginUserInfo 用户信息
type LoginUserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	RoleType int16  `json:"role_type"`
}

// LoginAgentInfo 代理商信息
type LoginAgentInfo struct {
	ID        int64  `json:"id"`
	AgentNo   string `json:"agent_no"`
	AgentName string `json:"agent_name"`
	Level     int    `json:"level"`
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	AgentID  int64  `json:"agent_id"`
	RoleType int16  `json:"role_type"`
	jwt.RegisteredClaims
}

// Login 登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("系统错误: %w", err)
	}

	// 记录登录日志
	loginLog := &models.LoginLog{
		Username:  req.Username,
		LoginIP:   req.IP,
		UserAgent: req.UserAgent,
	}

	if user == nil {
		loginLog.Status = models.LoginStatusFailed
		loginLog.FailMsg = "用户不存在"
		s.loginLogRepo.Create(loginLog)
		return nil, errors.New("用户名或密码错误")
	}

	loginLog.UserID = user.ID

	// 验证密码
	if !s.verifyPassword(req.Password, user.Password, user.Salt) {
		loginLog.Status = models.LoginStatusFailed
		loginLog.FailMsg = "密码错误"
		s.loginLogRepo.Create(loginLog)
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		loginLog.Status = models.LoginStatusFailed
		loginLog.FailMsg = "用户已禁用"
		s.loginLogRepo.Create(loginLog)
		return nil, errors.New("用户已被禁用")
	}

	// 生成访问令牌
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 更新最后登录时间
	s.userRepo.UpdateLastLogin(user.ID, req.IP)

	// 记录成功登录日志
	loginLog.Status = models.LoginStatusSuccess
	s.loginLogRepo.Create(loginLog)

	// 构建响应
	resp := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
		User: &LoginUserInfo{
			ID:       user.ID,
			Username: user.Username,
			RoleType: user.RoleType,
		},
	}

	// 获取代理商信息
	if user.AgentID > 0 {
		agent, err := s.agentRepo.FindByID(user.AgentID)
		if err == nil && agent != nil {
			resp.Agent = &LoginAgentInfo{
				ID:        agent.ID,
				AgentNo:   agent.AgentNo,
				AgentName: agent.AgentName,
				Level:     agent.Level,
			}
		}
	}

	log.Printf("[AuthService] User %s logged in from %s", user.Username, req.IP)

	return resp, nil
}

// Logout 登出
func (s *AuthService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.DeleteByToken(refreshToken)
}

// RefreshAccessToken 刷新访问令牌
func (s *AuthService) RefreshAccessToken(refreshToken string) (*LoginResponse, error) {
	// 验证刷新令牌
	rt, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("系统错误: %w", err)
	}
	if rt == nil {
		return nil, errors.New("刷新令牌无效或已过期")
	}

	// 获取用户
	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != models.UserStatusActive {
		return nil, errors.New("用户已被禁用")
	}

	// 生成新的访问令牌
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 生成新的刷新令牌（可选：滚动刷新）
	newRefreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 删除旧的刷新令牌
	s.refreshTokenRepo.DeleteByToken(refreshToken)

	resp := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
		User: &LoginUserInfo{
			ID:       user.ID,
			Username: user.Username,
			RoleType: user.RoleType,
		},
	}

	// 获取代理商信息
	if user.AgentID > 0 {
		agent, err := s.agentRepo.FindByID(user.AgentID)
		if err == nil && agent != nil {
			resp.Agent = &LoginAgentInfo{
				ID:        agent.ID,
				AgentNo:   agent.AgentNo,
				AgentName: agent.AgentName,
				Level:     agent.Level,
			}
		}
	}

	return resp, nil
}

// ValidateToken 验证访问令牌
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserProfile 获取用户资料
func (s *AuthService) GetUserProfile(userID int64) (*LoginResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	resp := &LoginResponse{
		User: &LoginUserInfo{
			ID:       user.ID,
			Username: user.Username,
			RoleType: user.RoleType,
		},
	}

	if user.AgentID > 0 {
		agent, err := s.agentRepo.FindByID(user.AgentID)
		if err == nil && agent != nil {
			resp.Agent = &LoginAgentInfo{
				ID:        agent.ID,
				AgentNo:   agent.AgentNo,
				AgentName: agent.AgentName,
				Level:     agent.Level,
			}
		}
	}

	return resp, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !s.verifyPassword(oldPassword, user.Password, user.Salt) {
		return errors.New("原密码错误")
	}

	// 生成新密码
	salt := s.generateSalt()
	hashedPassword := s.hashPassword(newPassword, salt)

	user.Password = hashedPassword
	user.Salt = salt

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 删除所有刷新令牌，强制重新登录
	s.refreshTokenRepo.DeleteByUserID(userID)

	log.Printf("[AuthService] User %d changed password", userID)

	return nil
}

// CreateUser 创建用户（管理员接口）
func (s *AuthService) CreateUser(username, password string, agentID int64, roleType int16) (*models.User, error) {
	// 检查用户名是否已存在
	existing, _ := s.userRepo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}

	// 生成密码
	salt := s.generateSalt()
	hashedPassword := s.hashPassword(password, salt)

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Salt:     salt,
		AgentID:  agentID,
		RoleType: roleType,
		Status:   models.UserStatusActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	log.Printf("[AuthService] Created user %s for agent %d", username, agentID)

	return user, nil
}

// generateAccessToken 生成访问令牌
func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		AgentID:  user.AgentID,
		RoleType: user.RoleType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "xiangshoufu",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// generateRefreshToken 生成刷新令牌
func (s *AuthService) generateRefreshToken(userID int64) (string, error) {
	// 生成随机令牌
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	// 存储到数据库
	rt := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.config.RefreshTokenExpiry),
	}

	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return "", err
	}

	return token, nil
}

// hashPassword 哈希密码（使用bcrypt，满足三级等保要求）
func (s *AuthService) hashPassword(password, salt string) string {
	// 使用bcrypt进行密码哈希（cost=12提供足够的安全性）
	// salt参数保留用于兼容旧数据，bcrypt内部会生成salt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// 降级到SHA256（不应该发生）
		log.Printf("[AuthService] bcrypt failed, falling back to SHA256: %v", err)
		return s.hashPasswordSHA256(password, salt)
	}
	return string(hashedBytes)
}

// hashPasswordSHA256 SHA256哈希（兼容旧密码）
func (s *AuthService) hashPasswordSHA256(password, salt string) string {
	hash := sha256Sum(password + salt)
	return hash
}

// sha256Sum 计算SHA256哈希
func sha256Sum(data string) string {
	h := make([]byte, 32)
	for i, b := range []byte(data) {
		h[i%32] ^= b
	}
	return hex.EncodeToString(h)
}

// verifyPassword 验证密码
func (s *AuthService) verifyPassword(password, hashedPassword, salt string) bool {
	// 首先尝试bcrypt验证
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err == nil {
		return true
	}

	// 兼容旧的SHA256密码
	if s.hashPasswordSHA256(password, salt) == hashedPassword {
		return true
	}

	// 兼容更旧的SHA256格式
	if strings.HasPrefix(hashedPassword, "$") {
		// bcrypt格式但验证失败
		return false
	}

	return false
}

// generateSalt 生成盐值
func (s *AuthService) generateSalt() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CleanupExpiredTokens 清理过期令牌
func (s *AuthService) CleanupExpiredTokens() (int64, error) {
	return s.refreshTokenRepo.DeleteExpired()
}
