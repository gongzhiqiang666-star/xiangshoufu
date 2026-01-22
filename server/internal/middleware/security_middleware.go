// Package middleware 安全中间件
// 提供XSS防护、CSRF防护、安全头部、SQL注入防护等功能
// 满足三级等保安全要求
package middleware

import (
	"html"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== 安全头部中间件 ====================

// SecurityHeadersMiddleware 添加安全响应头
// 符合三级等保和OWASP安全标准
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")

		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 启用XSS过滤器
		c.Header("X-XSS-Protection", "1; mode=block")

		// 强制HTTPS（生产环境启用）
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// 缓存控制（敏感接口禁止缓存）
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// ==================== XSS防护中间件 ====================

// XSSProtectionMiddleware XSS防护中间件
// 对请求参数进行HTML转义处理
func XSSProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对URL查询参数进行转义
		for key, values := range c.Request.URL.Query() {
			for i, v := range values {
				values[i] = html.EscapeString(v)
			}
			c.Request.URL.RawQuery = strings.Replace(
				c.Request.URL.RawQuery,
				key+"="+values[0],
				key+"="+html.EscapeString(values[0]),
				1,
			)
		}

		c.Next()
	}
}

// SanitizeInput 输入净化函数
// 移除潜在的XSS攻击代码
func SanitizeInput(input string) string {
	// 移除script标签
	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)
	input = scriptPattern.ReplaceAllString(input, "")

	// 移除事件处理属性
	eventPattern := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	input = eventPattern.ReplaceAllString(input, "")

	// 移除javascript:协议
	jsPattern := regexp.MustCompile(`(?i)javascript:`)
	input = jsPattern.ReplaceAllString(input, "")

	// HTML实体转义
	input = html.EscapeString(input)

	return input
}

// ==================== SQL注入防护 ====================

// SQLInjectionPatterns SQL注入检测模式
var SQLInjectionPatterns = []string{
	`(?i)(\s|^)(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|TRUNCATE)(\s|$)`,
	`(?i)(\s|^)(OR|AND)\s+[\d\w]+\s*=\s*[\d\w]+`,
	`(?i)'(\s|%20)*OR(\s|%20)*'`,
	`(?i)--`,
	`(?i);(\s|%20)*DROP`,
	`(?i)(\s|^)EXEC(\s|%20|\()+`,
	`(?i)WAITFOR(\s|%20)+DELAY`,
	`(?i)BENCHMARK\s*\(`,
	`(?i)SLEEP\s*\(`,
}

var sqlPatterns []*regexp.Regexp
var sqlPatternOnce sync.Once

func getSQLPatterns() []*regexp.Regexp {
	sqlPatternOnce.Do(func() {
		sqlPatterns = make([]*regexp.Regexp, 0, len(SQLInjectionPatterns))
		for _, pattern := range SQLInjectionPatterns {
			if re, err := regexp.Compile(pattern); err == nil {
				sqlPatterns = append(sqlPatterns, re)
			}
		}
	})
	return sqlPatterns
}

// DetectSQLInjection 检测SQL注入
func DetectSQLInjection(input string) bool {
	patterns := getSQLPatterns()
	for _, pattern := range patterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// SQLInjectionMiddleware SQL注入防护中间件
func SQLInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查URL参数
		for _, values := range c.Request.URL.Query() {
			for _, v := range values {
				if DetectSQLInjection(v) {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    400,
						"message": "请求包含非法字符",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// ==================== CSRF防护 ====================

// CSRFConfig CSRF配置
type CSRFConfig struct {
	TokenLength   int
	CookieName    string
	HeaderName    string
	CookieMaxAge  int
	CookieSecure  bool
	CookieHTTPOnly bool
	ExemptPaths   []string // 豁免路径
	ExemptMethods []string // 豁免方法
}

// DefaultCSRFConfig 默认CSRF配置
func DefaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		TokenLength:   32,
		CookieName:    "csrf_token",
		HeaderName:    "X-CSRF-Token",
		CookieMaxAge:  3600,
		CookieSecure:  false, // 生产环境设为true
		CookieHTTPOnly: true,
		ExemptPaths:   []string{"/api/v1/auth/login", "/api/v1/auth/register", "/callback/"},
		ExemptMethods: []string{"GET", "HEAD", "OPTIONS"},
	}
}

// CSRFMiddleware CSRF防护中间件
func CSRFMiddleware(config *CSRFConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCSRFConfig()
	}

	return func(c *gin.Context) {
		// 检查是否豁免
		for _, path := range config.ExemptPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		for _, method := range config.ExemptMethods {
			if c.Request.Method == method {
				c.Next()
				return
			}
		}

		// 验证CSRF Token
		cookieToken, err := c.Cookie(config.CookieName)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "缺少CSRF令牌",
			})
			c.Abort()
			return
		}

		headerToken := c.GetHeader(config.HeaderName)
		if headerToken == "" {
			headerToken = c.PostForm("_csrf")
		}

		if cookieToken != headerToken {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "CSRF令牌验证失败",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==================== IP黑名单 ====================

// IPBlacklist IP黑名单管理
type IPBlacklist struct {
	blacklist map[string]time.Time // IP -> 过期时间
	mu        sync.RWMutex
}

// NewIPBlacklist 创建IP黑名单
func NewIPBlacklist() *IPBlacklist {
	bl := &IPBlacklist{
		blacklist: make(map[string]time.Time),
	}
	// 启动清理协程
	go bl.cleanup()
	return bl
}

// Add 添加IP到黑名单
func (b *IPBlacklist) Add(ip string, duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.blacklist[ip] = time.Now().Add(duration)
}

// Remove 从黑名单移除IP
func (b *IPBlacklist) Remove(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.blacklist, ip)
}

// IsBlocked 检查IP是否被封禁
func (b *IPBlacklist) IsBlocked(ip string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	expiry, exists := b.blacklist[ip]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		return false
	}

	return true
}

// cleanup 定期清理过期的黑名单记录
func (b *IPBlacklist) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		now := time.Now()
		for ip, expiry := range b.blacklist {
			if now.After(expiry) {
				delete(b.blacklist, ip)
			}
		}
		b.mu.Unlock()
	}
}

// IPBlacklistMiddleware IP黑名单中间件
func IPBlacklistMiddleware(blacklist *IPBlacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if blacklist.IsBlocked(ip) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "您的IP已被封禁",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==================== 登录失败锁定 ====================

// LoginAttemptTracker 登录尝试追踪器
type LoginAttemptTracker struct {
	attempts  map[string]*loginAttempt // key: username 或 IP
	mu        sync.RWMutex
	maxFails  int           // 最大失败次数
	lockTime  time.Duration // 锁定时间
	blacklist *IPBlacklist  // IP黑名单引用
}

type loginAttempt struct {
	failCount  int
	lastFail   time.Time
	lockedUtil time.Time
}

// NewLoginAttemptTracker 创建登录尝试追踪器
func NewLoginAttemptTracker(maxFails int, lockTime time.Duration, blacklist *IPBlacklist) *LoginAttemptTracker {
	tracker := &LoginAttemptTracker{
		attempts:  make(map[string]*loginAttempt),
		maxFails:  maxFails,
		lockTime:  lockTime,
		blacklist: blacklist,
	}
	// 启动清理协程
	go tracker.cleanup()
	return tracker
}

// RecordFailure 记录登录失败
func (t *LoginAttemptTracker) RecordFailure(key string, ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	attempt, exists := t.attempts[key]
	if !exists {
		attempt = &loginAttempt{}
		t.attempts[key] = attempt
	}

	// 如果距离上次失败超过锁定时间，重置计数
	if time.Since(attempt.lastFail) > t.lockTime {
		attempt.failCount = 0
	}

	attempt.failCount++
	attempt.lastFail = time.Now()

	// 达到最大失败次数，锁定账户
	if attempt.failCount >= t.maxFails {
		attempt.lockedUtil = time.Now().Add(t.lockTime)

		// 如果失败次数过多，将IP加入黑名单
		if attempt.failCount >= t.maxFails*2 && t.blacklist != nil {
			t.blacklist.Add(ip, t.lockTime*2)
		}

		return true // 已锁定
	}

	return false
}

// IsLocked 检查是否被锁定
func (t *LoginAttemptTracker) IsLocked(key string) (bool, time.Duration) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	attempt, exists := t.attempts[key]
	if !exists {
		return false, 0
	}

	if time.Now().Before(attempt.lockedUtil) {
		remaining := time.Until(attempt.lockedUtil)
		return true, remaining
	}

	return false, 0
}

// ResetAttempts 重置登录尝试（登录成功后调用）
func (t *LoginAttemptTracker) ResetAttempts(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.attempts, key)
}

// GetFailCount 获取失败次数
func (t *LoginAttemptTracker) GetFailCount(key string) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	attempt, exists := t.attempts[key]
	if !exists {
		return 0
	}
	return attempt.failCount
}

// cleanup 定期清理过期记录
func (t *LoginAttemptTracker) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		t.mu.Lock()
		now := time.Now()
		for key, attempt := range t.attempts {
			// 清理超过2倍锁定时间未活动的记录
			if now.Sub(attempt.lastFail) > t.lockTime*2 {
				delete(t.attempts, key)
			}
		}
		t.mu.Unlock()
	}
}

// ==================== 密码强度校验 ====================

// PasswordStrength 密码强度等级
type PasswordStrength int

const (
	PasswordWeak     PasswordStrength = 1
	PasswordMedium   PasswordStrength = 2
	PasswordStrong   PasswordStrength = 3
	PasswordVeryStrong PasswordStrength = 4
)

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	MinLength      int  // 最小长度
	MaxLength      int  // 最大长度
	RequireUpper   bool // 需要大写字母
	RequireLower   bool // 需要小写字母
	RequireDigit   bool // 需要数字
	RequireSpecial bool // 需要特殊字符
	MinStrength    PasswordStrength // 最低强度要求
}

// DefaultPasswordPolicy 默认密码策略（满足三级等保要求）
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      8,
		MaxLength:      32,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false, // 可选
		MinStrength:    PasswordMedium,
	}
}

// ValidatePassword 验证密码是否符合策略
func ValidatePassword(password string, policy *PasswordPolicy) (bool, string) {
	if policy == nil {
		policy = DefaultPasswordPolicy()
	}

	// 长度检查
	if len(password) < policy.MinLength {
		return false, "密码长度不能少于" + string(rune('0'+policy.MinLength)) + "位"
	}
	if len(password) > policy.MaxLength {
		return false, "密码长度不能超过" + string(rune('0'+policy.MaxLength)) + "位"
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, c):
			hasSpecial = true
		}
	}

	if policy.RequireUpper && !hasUpper {
		return false, "密码必须包含大写字母"
	}
	if policy.RequireLower && !hasLower {
		return false, "密码必须包含小写字母"
	}
	if policy.RequireDigit && !hasDigit {
		return false, "密码必须包含数字"
	}
	if policy.RequireSpecial && !hasSpecial {
		return false, "密码必须包含特殊字符"
	}

	// 检查密码强度
	strength := CalculatePasswordStrength(password)
	if strength < policy.MinStrength {
		return false, "密码强度不足，请使用更复杂的密码"
	}

	return true, ""
}

// CalculatePasswordStrength 计算密码强度
func CalculatePasswordStrength(password string) PasswordStrength {
	score := 0

	// 长度得分
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	// 字符类型得分
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, c):
			hasSpecial = true
		}
	}

	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score += 2
	}

	// 根据得分返回强度等级
	switch {
	case score >= 8:
		return PasswordVeryStrong
	case score >= 5:
		return PasswordStrong
	case score >= 3:
		return PasswordMedium
	default:
		return PasswordWeak
	}
}

// ==================== 敏感日志脱敏 ====================

// SensitiveFields 敏感字段列表
var SensitiveFields = []string{
	"password", "passwd", "pwd",
	"token", "access_token", "refresh_token",
	"secret", "api_key", "apikey",
	"id_card", "idcard", "id_number",
	"bank_card", "bankcard", "card_no",
	"phone", "mobile", "tel",
	"email",
}

// MaskSensitiveData 脱敏敏感数据
func MaskSensitiveData(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		lowerKey := strings.ToLower(key)
		isSensitive := false

		for _, field := range SensitiveFields {
			if strings.Contains(lowerKey, field) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			if str, ok := value.(string); ok {
				result[key] = maskString(str)
			} else {
				result[key] = "***"
			}
		} else {
			result[key] = value
		}
	}

	return result
}

// maskString 对字符串进行脱敏
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	if len(s) <= 8 {
		return s[:2] + "****"
	}
	return s[:3] + "****" + s[len(s)-3:]
}

// ==================== 请求签名验证 ====================

// RequestSignatureConfig 请求签名配置
type RequestSignatureConfig struct {
	HeaderName     string        // 签名头名称
	TimestampHeader string       // 时间戳头名称
	NonceHeader    string        // 随机数头名称
	MaxTimeDiff    time.Duration // 最大时间差
	Secret         string        // 签名密钥
}

// DefaultRequestSignatureConfig 默认签名配置
func DefaultRequestSignatureConfig() *RequestSignatureConfig {
	return &RequestSignatureConfig{
		HeaderName:      "X-Signature",
		TimestampHeader: "X-Timestamp",
		NonceHeader:     "X-Nonce",
		MaxTimeDiff:     5 * time.Minute,
		Secret:          "xiangshoufu-api-secret",
	}
}
