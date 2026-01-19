package async

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// MemoryQueue 内存队列实现
type MemoryQueue struct {
	queues      map[string]chan []byte
	handlers    map[string]func([]byte) error
	workerCount int
	bufferSize  int
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// MemoryQueueConfig 内存队列配置
type MemoryQueueConfig struct {
	WorkerCount int // 每个主题的工作协程数
	BufferSize  int // 队列缓冲区大小
}

// DefaultMemoryQueueConfig 默认配置
func DefaultMemoryQueueConfig() *MemoryQueueConfig {
	return &MemoryQueueConfig{
		WorkerCount: 10,   // 默认10个工作协程
		BufferSize:  1000, // 默认缓冲1000条消息
	}
}

// NewMemoryQueue 创建内存队列
func NewMemoryQueue(config *MemoryQueueConfig) *MemoryQueue {
	if config == nil {
		config = DefaultMemoryQueueConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MemoryQueue{
		queues:      make(map[string]chan []byte),
		handlers:    make(map[string]func([]byte) error),
		workerCount: config.WorkerCount,
		bufferSize:  config.BufferSize,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Publish 发布消息
func (q *MemoryQueue) Publish(topic string, msg []byte) error {
	q.mu.RLock()
	ch, exists := q.queues[topic]
	q.mu.RUnlock()

	if !exists {
		return fmt.Errorf("topic not found: %s", topic)
	}

	// 非阻塞发送，如果队列满了返回错误
	select {
	case ch <- msg:
		return nil
	case <-q.ctx.Done():
		return fmt.Errorf("queue is closed")
	default:
		return fmt.Errorf("queue is full for topic: %s", topic)
	}
}

// Subscribe 订阅主题
func (q *MemoryQueue) Subscribe(topic string, handler func([]byte) error) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 检查是否已订阅
	if _, exists := q.queues[topic]; exists {
		return fmt.Errorf("topic already subscribed: %s", topic)
	}

	// 创建队列通道
	ch := make(chan []byte, q.bufferSize)
	q.queues[topic] = ch
	q.handlers[topic] = handler

	// 启动工作协程
	for i := 0; i < q.workerCount; i++ {
		q.wg.Add(1)
		go q.worker(topic, ch, handler, i)
	}

	log.Printf("[MemoryQueue] Subscribed to topic: %s with %d workers", topic, q.workerCount)
	return nil
}

// worker 工作协程
func (q *MemoryQueue) worker(topic string, ch <-chan []byte, handler func([]byte) error, workerID int) {
	defer q.wg.Done()

	for {
		select {
		case <-q.ctx.Done():
			log.Printf("[MemoryQueue] Worker %d for topic %s stopped", workerID, topic)
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			q.processMessage(topic, msg, handler, workerID)
		}
	}
}

// processMessage 处理单条消息
func (q *MemoryQueue) processMessage(topic string, msg []byte, handler func([]byte) error, workerID int) {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[MemoryQueue] Worker %d panic in topic %s: %v", workerID, topic, r)
		}
	}()

	if err := handler(msg); err != nil {
		log.Printf("[MemoryQueue] Handler error in topic %s: %v (took %v)", topic, err, time.Since(startTime))
	}
}

// Close 关闭队列
func (q *MemoryQueue) Close() error {
	q.cancel()

	// 关闭所有通道
	q.mu.Lock()
	for topic, ch := range q.queues {
		close(ch)
		log.Printf("[MemoryQueue] Closed channel for topic: %s", topic)
	}
	q.mu.Unlock()

	// 等待所有工作协程退出
	q.wg.Wait()
	log.Printf("[MemoryQueue] All workers stopped")

	return nil
}

// GetQueueLength 获取队列长度（用于监控）
func (q *MemoryQueue) GetQueueLength(topic string) int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if ch, exists := q.queues[topic]; exists {
		return len(ch)
	}
	return 0
}

// GetAllQueueLengths 获取所有队列长度
func (q *MemoryQueue) GetAllQueueLengths() map[string]int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	result := make(map[string]int)
	for topic, ch := range q.queues {
		result[topic] = len(ch)
	}
	return result
}

// 确保实现了MessageQueue接口
var _ MessageQueue = (*MemoryQueue)(nil)
