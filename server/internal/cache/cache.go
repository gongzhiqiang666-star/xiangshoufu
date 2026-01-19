package cache

import "time"

// Cache 缓存接口
// 当前使用本地缓存实现，后续可无缝切换到Redis
type Cache interface {
	// Get 获取缓存值
	Get(key string) (interface{}, bool)

	// Set 设置缓存值
	Set(key string, value interface{}, ttl time.Duration)

	// Delete 删除缓存
	Delete(key string)

	// Exists 检查key是否存在
	Exists(key string) bool

	// SetNX 如果key不存在则设置（用于分布式锁/幂等检查）
	SetNX(key string, value interface{}, ttl time.Duration) bool

	// Close 关闭缓存
	Close() error
}
