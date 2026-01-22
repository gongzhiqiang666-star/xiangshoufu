// Package crypto 提供敏感信息加密解密功能
// 用于满足三级等保对敏感数据加密存储的要求
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	// 全局加密器实例
	globalCrypto *Crypto
	once         sync.Once
)

// Crypto 加密解密器
type Crypto struct {
	key []byte
}

// GetDefaultCrypto 获取全局加密器实例
func GetDefaultCrypto() *Crypto {
	once.Do(func() {
		key := os.Getenv("CRYPTO_KEY")
		if key == "" {
			// 默认密钥（生产环境必须通过环境变量设置）
			key = "xiangshoufu-2026-crypto-key-32b"
		}
		globalCrypto = NewCrypto(key)
	})
	return globalCrypto
}

// NewCrypto 创建加密器
// key 必须是 16, 24, 或 32 字节长度（对应 AES-128, AES-192, AES-256）
func NewCrypto(key string) *Crypto {
	keyBytes := []byte(key)
	// 确保密钥长度为32字节（AES-256）
	if len(keyBytes) < 32 {
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}
	return &Crypto{key: keyBytes}
}

// Encrypt 加密字符串
func (c *Crypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	// GCM模式提供认证加密
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func (c *Crypto) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptPhone 加密手机号
func EncryptPhone(phone string) (string, error) {
	return GetDefaultCrypto().Encrypt(phone)
}

// DecryptPhone 解密手机号
func DecryptPhone(encrypted string) (string, error) {
	return GetDefaultCrypto().Decrypt(encrypted)
}

// EncryptIDCard 加密身份证号
func EncryptIDCard(idCard string) (string, error) {
	return GetDefaultCrypto().Encrypt(idCard)
}

// DecryptIDCard 解密身份证号
func DecryptIDCard(encrypted string) (string, error) {
	return GetDefaultCrypto().Decrypt(encrypted)
}

// MaskPhone 手机号脱敏显示（前3后4）
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskIDCard 身份证号脱敏显示（前4后4）
func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:4] + "**********" + idCard[len(idCard)-4:]
}

// IsEncrypted 判断字符串是否已加密（通过base64特征判断）
func IsEncrypted(s string) bool {
	if s == "" {
		return false
	}
	// 加密后的字符串是base64编码，长度会变长
	// 手机号11位加密后约44位，身份证18位加密后约60位
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil && len(s) > 20
}
