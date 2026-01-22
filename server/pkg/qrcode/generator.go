package qrcode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// Config 二维码配置
type Config struct {
	Size       int    // 二维码尺寸（像素）
	BaseURL    string // 注册页面基础URL
	StaticDir  string // 静态文件存储目录
	PublicURL  string // 静态文件公开访问URL
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Size:      256,
		BaseURL:   "https://m.xiangshoufu.com/register",
		StaticDir: "./static/qrcode",
		PublicURL: "/static/qrcode",
	}
}

// Generator 二维码生成器
type Generator struct {
	config *Config
}

// NewGenerator 创建二维码生成器
func NewGenerator(config *Config) *Generator {
	if config == nil {
		config = DefaultConfig()
	}
	// 确保存储目录存在
	os.MkdirAll(config.StaticDir, 0755)
	return &Generator{config: config}
}

// GenerateInviteURL 生成邀请链接
func (g *Generator) GenerateInviteURL(inviteCode string) string {
	return fmt.Sprintf("%s?code=%s", g.config.BaseURL, inviteCode)
}

// GenerateQRCodeBase64 生成二维码并返回Base64编码
func (g *Generator) GenerateQRCodeBase64(content string) (string, error) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("创建二维码失败: %w", err)
	}

	// 生成PNG图片
	var buf bytes.Buffer
	err = png.Encode(&buf, qr.Image(g.config.Size))
	if err != nil {
		return "", fmt.Errorf("编码二维码图片失败: %w", err)
	}

	// 返回Base64编码
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + base64Str, nil
}

// GenerateQRCodeFile 生成二维码并保存到文件
func (g *Generator) GenerateQRCodeFile(content string, filename string) (string, error) {
	filePath := filepath.Join(g.config.StaticDir, filename)

	err := qrcode.WriteFile(content, qrcode.Medium, g.config.Size, filePath)
	if err != nil {
		return "", fmt.Errorf("保存二维码文件失败: %w", err)
	}

	// 返回公开访问URL
	return fmt.Sprintf("%s/%s", g.config.PublicURL, filename), nil
}

// GenerateInviteQRCode 生成邀请二维码（完整流程）
// 返回: inviteURL, qrCodeURL(或Base64), error
func (g *Generator) GenerateInviteQRCode(inviteCode string, agentID int64) (string, string, error) {
	// 生成邀请链接
	inviteURL := g.GenerateInviteURL(inviteCode)

	// 生成二维码文件
	filename := fmt.Sprintf("invite_%d.png", agentID)
	qrCodeURL, err := g.GenerateQRCodeFile(inviteURL, filename)
	if err != nil {
		// 如果文件保存失败，尝试返回Base64
		qrCodeBase64, err2 := g.GenerateQRCodeBase64(inviteURL)
		if err2 != nil {
			return inviteURL, "", fmt.Errorf("生成二维码失败: %w", err)
		}
		return inviteURL, qrCodeBase64, nil
	}

	return inviteURL, qrCodeURL, nil
}

// GenerateQRCodeBytes 生成二维码并返回PNG字节
func (g *Generator) GenerateQRCodeBytes(content string) ([]byte, error) {
	return qrcode.Encode(content, qrcode.Medium, g.config.Size)
}
