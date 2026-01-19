package async

// MessageQueue 消息队列接口
// 当前使用内存队列实现，后续可无缝切换到Kafka
type MessageQueue interface {
	// Publish 发布消息到指定主题
	Publish(topic string, msg []byte) error

	// Subscribe 订阅指定主题，handler处理接收到的消息
	Subscribe(topic string, handler func([]byte) error) error

	// Close 关闭队列
	Close() error
}

// 预定义主题
const (
	TopicRawCallback  = "raw_callback" // 原始回调数据处理
	TopicProfitCalc   = "profit_calc"  // 分润计算
	TopicNotification = "notification" // 消息通知
)
