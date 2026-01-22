package jobs

import (
	"time"
)

// 指数退避重试策略
// 默认重试间隔: 1分钟 -> 5分钟 -> 15分钟
var DefaultRetryIntervals = []time.Duration{
	1 * time.Minute,
	5 * time.Minute,
	15 * time.Minute,
}

// RetryStrategy 重试策略接口
type RetryStrategy interface {
	// GetRetryInterval 获取第n次重试的等待时间
	GetRetryInterval(retryCount int) time.Duration
	// ShouldRetry 判断是否应该继续重试
	ShouldRetry(retryCount, maxRetries int) bool
}

// ExponentialBackoffStrategy 指数退避策略
type ExponentialBackoffStrategy struct {
	intervals []time.Duration
}

// NewExponentialBackoffStrategy 创建指数退避策略
func NewExponentialBackoffStrategy(intervals []time.Duration) *ExponentialBackoffStrategy {
	if len(intervals) == 0 {
		intervals = DefaultRetryIntervals
	}
	return &ExponentialBackoffStrategy{intervals: intervals}
}

// GetRetryInterval 获取重试间隔
func (s *ExponentialBackoffStrategy) GetRetryInterval(retryCount int) time.Duration {
	if retryCount <= 0 {
		return s.intervals[0]
	}
	idx := retryCount - 1
	if idx >= len(s.intervals) {
		idx = len(s.intervals) - 1
	}
	return s.intervals[idx]
}

// ShouldRetry 判断是否应该继续重试
func (s *ExponentialBackoffStrategy) ShouldRetry(retryCount, maxRetries int) bool {
	return retryCount < maxRetries
}

// FixedIntervalStrategy 固定间隔策略
type FixedIntervalStrategy struct {
	interval time.Duration
}

// NewFixedIntervalStrategy 创建固定间隔策略
func NewFixedIntervalStrategy(interval time.Duration) *FixedIntervalStrategy {
	return &FixedIntervalStrategy{interval: interval}
}

// GetRetryInterval 获取重试间隔
func (s *FixedIntervalStrategy) GetRetryInterval(retryCount int) time.Duration {
	return s.interval
}

// ShouldRetry 判断是否应该继续重试
func (s *FixedIntervalStrategy) ShouldRetry(retryCount, maxRetries int) bool {
	return retryCount < maxRetries
}

// LinearBackoffStrategy 线性退避策略
type LinearBackoffStrategy struct {
	baseInterval time.Duration
	increment    time.Duration
	maxInterval  time.Duration
}

// NewLinearBackoffStrategy 创建线性退避策略
func NewLinearBackoffStrategy(baseInterval, increment, maxInterval time.Duration) *LinearBackoffStrategy {
	return &LinearBackoffStrategy{
		baseInterval: baseInterval,
		increment:    increment,
		maxInterval:  maxInterval,
	}
}

// GetRetryInterval 获取重试间隔
func (s *LinearBackoffStrategy) GetRetryInterval(retryCount int) time.Duration {
	interval := s.baseInterval + time.Duration(retryCount)*s.increment
	if interval > s.maxInterval {
		interval = s.maxInterval
	}
	return interval
}

// ShouldRetry 判断是否应该继续重试
func (s *LinearBackoffStrategy) ShouldRetry(retryCount, maxRetries int) bool {
	return retryCount < maxRetries
}

// CreateRetryIntervalsFromConfig 从配置创建重试间隔
// baseInterval: 初始间隔（秒）
// maxRetries: 最大重试次数
func CreateRetryIntervalsFromConfig(baseIntervalSeconds int, maxRetries int) []time.Duration {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	if baseIntervalSeconds <= 0 {
		baseIntervalSeconds = 60
	}

	intervals := make([]time.Duration, maxRetries)
	base := time.Duration(baseIntervalSeconds) * time.Second

	for i := 0; i < maxRetries; i++ {
		// 指数增长: base * 2^i，但设置上限为15分钟
		multiplier := 1 << i // 2^i
		interval := base * time.Duration(multiplier)
		if interval > 15*time.Minute {
			interval = 15 * time.Minute
		}
		intervals[i] = interval
	}

	return intervals
}
