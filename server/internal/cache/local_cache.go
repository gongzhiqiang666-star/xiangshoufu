package cache

import (
	"sync"
	"time"
)

// LocalCache 本地内存缓存实现
type LocalCache struct {
	items           map[string]*cacheItem
	mu              sync.RWMutex
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// cacheItem 缓存项
type cacheItem struct {
	value     interface{}
	expireAt  time.Time
	hasExpiry bool
}

// LocalCacheConfig 本地缓存配置
type LocalCacheConfig struct {
	CleanupInterval time.Duration // 清理过期项的间隔
}

// DefaultLocalCacheConfig 默认配置
func DefaultLocalCacheConfig() *LocalCacheConfig {
	return &LocalCacheConfig{
		CleanupInterval: 5 * time.Minute,
	}
}

// NewLocalCache 创建本地缓存
func NewLocalCache(config *LocalCacheConfig) *LocalCache {
	if config == nil {
		config = DefaultLocalCacheConfig()
	}

	c := &LocalCache{
		items:           make(map[string]*cacheItem),
		cleanupInterval: config.CleanupInterval,
		stopCleanup:     make(chan struct{}),
	}

	// 启动后台清理协程
	go c.cleanupLoop()

	return c
}

// Get 获取缓存值
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if item.hasExpiry && time.Now().After(item.expireAt) {
		return nil, false
	}

	return item.value, true
}

// Set 设置缓存值
func (c *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := &cacheItem{
		value: value,
	}

	if ttl > 0 {
		item.expireAt = time.Now().Add(ttl)
		item.hasExpiry = true
	}

	c.items[key] = item
}

// Delete 删除缓存
func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Exists 检查key是否存在
func (c *LocalCache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	// 检查是否过期
	if item.hasExpiry && time.Now().After(item.expireAt) {
		return false
	}

	return true
}

// SetNX 如果key不存在则设置（原子操作）
func (c *LocalCache) SetNX(key string, value interface{}, ttl time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否已存在且未过期
	if item, exists := c.items[key]; exists {
		if !item.hasExpiry || time.Now().Before(item.expireAt) {
			return false
		}
	}

	// 设置新值
	item := &cacheItem{
		value: value,
	}

	if ttl > 0 {
		item.expireAt = time.Now().Add(ttl)
		item.hasExpiry = true
	}

	c.items[key] = item
	return true
}

// Close 关闭缓存
func (c *LocalCache) Close() error {
	close(c.stopCleanup)
	return nil
}

// cleanupLoop 定期清理过期项
func (c *LocalCache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanup 清理过期项
func (c *LocalCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if item.hasExpiry && now.After(item.expireAt) {
			delete(c.items, key)
		}
	}
}

// Size 获取缓存大小（用于监控）
func (c *LocalCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// 确保实现了Cache接口
var _ Cache = (*LocalCache)(nil)
