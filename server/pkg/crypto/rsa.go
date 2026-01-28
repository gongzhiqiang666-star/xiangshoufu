package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"sync"
	"time"
)

// RSAKeyPair RSA密钥对
type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	PublicPEM  string
	CreatedAt  time.Time
}

// RSAManager RSA密钥管理器
type RSAManager struct {
	keyPair    *RSAKeyPair
	mu         sync.RWMutex
	keyBits    int
	rotateTime time.Duration
}

var (
	defaultManager *RSAManager
	rsaOnce        sync.Once
)

// GetDefaultManager 获取默认RSA管理器（单例）
func GetDefaultManager() *RSAManager {
	rsaOnce.Do(func() {
		defaultManager = NewRSAManager(2048, 24*time.Hour)
	})
	return defaultManager
}

// NewRSAManager 创建RSA管理器
func NewRSAManager(keyBits int, rotateTime time.Duration) *RSAManager {
	return &RSAManager{
		keyBits:    keyBits,
		rotateTime: rotateTime,
	}
}

// GenerateKeyPair 生成新的密钥对
func (m *RSAManager) GenerateKeyPair() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	privateKey, err := rsa.GenerateKey(rand.Reader, m.keyBits)
	if err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	m.keyPair = &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		PublicPEM:  string(publicPEM),
		CreatedAt:  time.Now(),
	}

	return nil
}

// GetPublicKey 获取公钥（PEM格式）
func (m *RSAManager) GetPublicKey() (string, error) {
	m.mu.RLock()
	needGenerate := m.keyPair == nil || time.Since(m.keyPair.CreatedAt) > m.rotateTime
	m.mu.RUnlock()

	if needGenerate {
		if err := m.GenerateKeyPair(); err != nil {
			return "", err
		}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.keyPair == nil {
		return "", errors.New("密钥对未初始化")
	}

	return m.keyPair.PublicPEM, nil
}

// Decrypt 解密数据（Base64编码的密文）
func (m *RSAManager) Decrypt(encryptedBase64 string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.keyPair == nil {
		return "", errors.New("密钥对未初始化")
	}

	// Base64解码
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", errors.New("密文格式错误")
	}

	// RSA-OAEP解密
	decryptedBytes, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		m.keyPair.PrivateKey,
		encryptedBytes,
		nil,
	)
	if err != nil {
		return "", errors.New("解密失败")
	}

	return string(decryptedBytes), nil
}

// DecryptPassword 解密密码（便捷方法）
func DecryptPassword(encryptedPassword string) (string, error) {
	return GetDefaultManager().Decrypt(encryptedPassword)
}

// GetPublicKeyPEM 获取公钥PEM（便捷方法）
func GetPublicKeyPEM() (string, error) {
	return GetDefaultManager().GetPublicKey()
}
