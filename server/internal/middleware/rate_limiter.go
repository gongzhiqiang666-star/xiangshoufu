package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 单机限流器（令牌桶算法）
type RateLimiter struct {
	rate       float64   // 每秒产生的令牌数
	capacity   int64     // 桶容量
	tokens     float64   // 当前令牌数
	lastUpdate time.Time // 上次更新时间
	mu         sync.Mutex
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate float64, capacity int64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     float64(capacity),
		lastUpdate: time.Now(),
	}
}

// Allow 判断是否允许通过
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastUpdate).Seconds()

	// 添加新令牌
	r.tokens += elapsed * r.rate
	if r.tokens > float64(r.capacity) {
		r.tokens = float64(r.capacity)
	}
	r.lastUpdate = now

	// 检查是否有可用令牌
	if r.tokens >= 1 {
		r.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    -1,
				"message": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter 基于IP的限流器
type IPRateLimiter struct {
	limiters map[string]*RateLimiter
	rate     float64
	capacity int64
	mu       sync.RWMutex
}

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter(rate float64, capacity int64) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		rate:     rate,
		capacity: capacity,
	}
}

// GetLimiter 获取指定IP的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	i.mu.RLock()
	limiter, exists := i.limiters[ip]
	i.mu.RUnlock()

	if exists {
		return limiter
	}

	// 创建新的限流器
	i.mu.Lock()
	defer i.mu.Unlock()

	// 双重检查
	if limiter, exists = i.limiters[ip]; exists {
		return limiter
	}

	limiter = NewRateLimiter(i.rate, i.capacity)
	i.limiters[ip] = limiter

	return limiter
}

// IPRateLimitMiddleware IP限流中间件
func IPRateLimitMiddleware(ipLimiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := ipLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    -1,
				"message": "too many requests from your IP",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件（防止panic）
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    -1,
					"message": "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// 只记录非200状态或慢请求
		if status != http.StatusOK || latency > 100*time.Millisecond {
			// log.Printf("[HTTP] %s %s %d %v", c.Request.Method, c.Request.URL.Path, status, latency)
		}
	}
}
