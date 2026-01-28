package channel

import (
	"fmt"
	"sync"
)

// 通道编码常量
const (
	ChannelCodeHengxintong = "HENGXINTONG" // 恒信通
	ChannelCodeLakala      = "LAKALA"      // 拉卡拉
	ChannelCodeYeahka      = "YEAHKA"      // 乐刷
	ChannelCodeSuixingfu   = "SUIXINGFU"   // 随行付
	ChannelCodeLianlian    = "LIANLIAN"    // 连连支付
	ChannelCodeSandpay     = "SANDPAY"     // 杉德支付
	ChannelCodeFuiou       = "FUIOU"       // 富友支付
	ChannelCodeHeepay      = "HEEPAY"      // 汇付天下
	ChannelCodeLiandong    = "LIANDONG"    // 联动
)

// AdapterFactory 适配器工厂
type AdapterFactory struct {
	adapters map[string]ChannelAdapter
	mu       sync.RWMutex
}

var (
	factory     *AdapterFactory
	factoryOnce sync.Once
)

// GetFactory 获取适配器工厂单例
func GetFactory() *AdapterFactory {
	factoryOnce.Do(func() {
		factory = &AdapterFactory{
			adapters: make(map[string]ChannelAdapter),
		}
	})
	return factory
}

// Register 注册适配器
func (f *AdapterFactory) Register(adapter ChannelAdapter) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.adapters[adapter.GetChannelCode()] = adapter
}

// GetAdapter 根据通道编码获取适配器
func (f *AdapterFactory) GetAdapter(channelCode string) (ChannelAdapter, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	adapter, exists := f.adapters[channelCode]
	if !exists {
		return nil, fmt.Errorf("adapter not found for channel: %s", channelCode)
	}
	return adapter, nil
}

// GetAllAdapters 获取所有已注册的适配器
func (f *AdapterFactory) GetAllAdapters() []ChannelAdapter {
	f.mu.RLock()
	defer f.mu.RUnlock()

	adapters := make([]ChannelAdapter, 0, len(f.adapters))
	for _, adapter := range f.adapters {
		adapters = append(adapters, adapter)
	}
	return adapters
}

// GetSupportedChannels 获取所有支持的通道编码
func (f *AdapterFactory) GetSupportedChannels() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	channels := make([]string, 0, len(f.adapters))
	for code := range f.adapters {
		channels = append(channels, code)
	}
	return channels
}

// HasAdapter 检查是否存在指定通道的适配器
func (f *AdapterFactory) HasAdapter(channelCode string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, exists := f.adapters[channelCode]
	return exists
}

// ChannelConfig 通道配置
type ChannelConfig struct {
	ChannelCode string `json:"channel_code"` // 通道编码
	ChannelName string `json:"channel_name"` // 通道名称
	PublicKey   string `json:"public_key"`   // RSA公钥（用于验签）
	PrivateKey  string `json:"private_key"`  // RSA私钥（用于签名，如有需要）
	APIKey      string `json:"api_key"`      // API密钥
	APISecret   string `json:"api_secret"`   // API密钥
	APIBaseURL  string `json:"api_base_url"` // API基础URL（用于费率更新等主动调用）
	CallbackURL string `json:"callback_url"` // 回调URL
	Enabled     bool   `json:"enabled"`      // 是否启用
}

// ConfigurableAdapter 可配置的适配器接口（可选实现）
type ConfigurableAdapter interface {
	ChannelAdapter
	// Configure 配置适配器
	Configure(config *ChannelConfig) error
}
