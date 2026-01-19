package middleware

import (
	"net/http"
	"strings"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请先登录",
			})
			c.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证格式",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证Token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证已过期，请重新登录",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("agent_id", claims.AgentID)
		c.Set("role_type", claims.RoleType)
		c.Set("claims", claims)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleType, exists := c.Get("role_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请先登录",
			})
			c.Abort()
			return
		}

		if roleType.(int16) != models.UserRoleTypeAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制登录，但如果有token则解析）
func OptionalAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("agent_id", claims.AgentID)
			c.Set("role_type", claims.RoleType)
			c.Set("claims", claims)
		}

		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if id, ok := userID.(int64); ok {
		return id
	}
	return 0
}

// GetCurrentAgentID 从上下文获取当前代理商ID
func GetCurrentAgentID(c *gin.Context) int64 {
	agentID, exists := c.Get("agent_id")
	if !exists {
		return 0
	}
	if id, ok := agentID.(int64); ok {
		return id
	}
	return 0
}

// GetCurrentUsername 从上下文获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	if name, ok := username.(string); ok {
		return name
	}
	return ""
}

// GetCurrentRoleType 从上下文获取当前角色类型
func GetCurrentRoleType(c *gin.Context) int16 {
	roleType, exists := c.Get("role_type")
	if !exists {
		return 0
	}
	if rt, ok := roleType.(int16); ok {
		return rt
	}
	return 0
}

// IsAdmin 检查是否为管理员
func IsAdmin(c *gin.Context) bool {
	return GetCurrentRoleType(c) == models.UserRoleTypeAdmin
}
