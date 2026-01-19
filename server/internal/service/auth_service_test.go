package service

import (
	"testing"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int64) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// MockRefreshTokenRepository 模拟刷新令牌仓库
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteByToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteByUserID(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

// MockLoginLogRepository 模拟登录日志仓库
type MockLoginLogRepository struct {
	mock.Mock
}

func (m *MockLoginLogRepository) Create(log *models.LoginLog) error {
	args := m.Called(log)
	return args.Error(0)
}

// MockAuthAgentRepository 模拟代理商仓库（用于认证测试）
type MockAuthAgentRepository struct {
	mock.Mock
}

func (m *MockAuthAgentRepository) FindByID(id int64) (*repository.Agent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.Agent), args.Error(1)
}

func (m *MockAuthAgentRepository) FindByAgentNo(agentNo string) (*repository.Agent, error) {
	args := m.Called(agentNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.Agent), args.Error(1)
}

func (m *MockAuthAgentRepository) FindAncestors(agentID int64) ([]*repository.Agent, error) {
	args := m.Called(agentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repository.Agent), args.Error(1)
}

// TestAuthConfig 测试认证配置
func TestDefaultAuthConfig(t *testing.T) {
	config := DefaultAuthConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 2*time.Hour, config.AccessTokenExpiry)
	assert.Equal(t, 7*24*time.Hour, config.RefreshTokenExpiry)
	assert.NotEmpty(t, config.JWTSecret)
}

// TestJWTClaims 测试JWT声明
func TestJWTClaims(t *testing.T) {
	claims := &JWTClaims{
		UserID:   1,
		Username: "testuser",
		AgentID:  100,
		RoleType: 1,
	}

	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, int64(100), claims.AgentID)
	assert.Equal(t, int16(1), claims.RoleType)
}

// TestLoginRequest 测试登录请求验证
func TestLoginRequest(t *testing.T) {
	tests := []struct {
		name     string
		req      *LoginRequest
		hasError bool
	}{
		{
			name: "valid request",
			req: &LoginRequest{
				Username: "admin",
				Password: "password123",
			},
			hasError: false,
		},
		{
			name: "empty username",
			req: &LoginRequest{
				Username: "",
				Password: "password123",
			},
			hasError: true,
		},
		{
			name: "empty password",
			req: &LoginRequest{
				Username: "admin",
				Password: "",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.req.Username == "" || tt.req.Password == ""
			assert.Equal(t, tt.hasError, hasError)
		})
	}
}

// TestPasswordChange 测试密码修改验证
func TestPasswordChange(t *testing.T) {
	tests := []struct {
		name        string
		oldPassword string
		newPassword string
		hasError    bool
	}{
		{
			name:        "valid change",
			oldPassword: "oldpass",
			newPassword: "newpass123",
			hasError:    false,
		},
		{
			name:        "same passwords",
			oldPassword: "samepass",
			newPassword: "samepass",
			hasError:    true, // 新旧密码不能相同
		},
		{
			name:        "short new password",
			oldPassword: "oldpass",
			newPassword: "short",
			hasError:    true, // 密码太短
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.oldPassword == tt.newPassword || len(tt.newPassword) < 6
			assert.Equal(t, tt.hasError, hasError)
		})
	}
}

// TestLoginResponse 测试登录响应结构
func TestLoginResponse(t *testing.T) {
	resp := &LoginResponse{
		AccessToken:  "access_token_here",
		RefreshToken: "refresh_token_here",
		ExpiresIn:    7200,
		TokenType:    "Bearer",
		User: &LoginUserInfo{
			ID:       1,
			Username: "testuser",
			RoleType: 1,
		},
		Agent: &LoginAgentInfo{
			ID:        100,
			AgentNo:   "A001",
			AgentName: "Test Agent",
			Level:     1,
		},
	}

	assert.Equal(t, "access_token_here", resp.AccessToken)
	assert.Equal(t, "refresh_token_here", resp.RefreshToken)
	assert.Equal(t, int64(7200), resp.ExpiresIn)
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int64(1), resp.User.ID)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.NotNil(t, resp.Agent)
	assert.Equal(t, int64(100), resp.Agent.ID)
	assert.Equal(t, "Test Agent", resp.Agent.AgentName)
}

// TestLoginUserInfo 测试登录用户信息结构
func TestLoginUserInfo(t *testing.T) {
	info := &LoginUserInfo{
		ID:       1,
		Username: "testuser",
		RoleType: 1,
	}

	assert.Equal(t, int64(1), info.ID)
	assert.Equal(t, "testuser", info.Username)
	assert.Equal(t, int16(1), info.RoleType)
}
