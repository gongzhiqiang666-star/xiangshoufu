package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	c := NewCrypto("test-key-32-bytes-long-string!!")

	tests := []struct {
		name      string
		plaintext string
	}{
		{"手机号", "13800138000"},
		{"身份证", "110101199001011234"},
		{"空字符串", ""},
		{"短字符串", "abc"},
		{"中文", "张三"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := c.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			decrypted, err := c.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("Decrypt mismatch: got %s, want %s", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptPhone(t *testing.T) {
	phone := "13800138000"

	encrypted, err := EncryptPhone(phone)
	if err != nil {
		t.Fatalf("EncryptPhone failed: %v", err)
	}

	if encrypted == phone {
		t.Error("Encrypted phone should be different from original")
	}

	decrypted, err := DecryptPhone(encrypted)
	if err != nil {
		t.Fatalf("DecryptPhone failed: %v", err)
	}

	if decrypted != phone {
		t.Errorf("DecryptPhone mismatch: got %s, want %s", decrypted, phone)
	}
}

func TestEncryptIDCard(t *testing.T) {
	idCard := "110101199001011234"

	encrypted, err := EncryptIDCard(idCard)
	if err != nil {
		t.Fatalf("EncryptIDCard failed: %v", err)
	}

	if encrypted == idCard {
		t.Error("Encrypted ID card should be different from original")
	}

	decrypted, err := DecryptIDCard(encrypted)
	if err != nil {
		t.Fatalf("DecryptIDCard failed: %v", err)
	}

	if decrypted != idCard {
		t.Errorf("DecryptIDCard mismatch: got %s, want %s", decrypted, idCard)
	}
}

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		phone    string
		expected string
	}{
		{"13800138000", "138****8000"},
		{"123456", "123456"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskPhone(tt.phone)
		if result != tt.expected {
			t.Errorf("MaskPhone(%s) = %s, want %s", tt.phone, result, tt.expected)
		}
	}
}

func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		idCard   string
		expected string
	}{
		{"110101199001011234", "1101**********1234"},
		{"1234567", "1234567"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskIDCard(tt.idCard)
		if result != tt.expected {
			t.Errorf("MaskIDCard(%s) = %s, want %s", tt.idCard, result, tt.expected)
		}
	}
}

func TestIsEncrypted(t *testing.T) {
	c := NewCrypto("test-key")

	// 未加密的手机号
	if IsEncrypted("13800138000") {
		t.Error("Plain phone should not be detected as encrypted")
	}

	// 加密后的手机号
	encrypted, _ := c.Encrypt("13800138000")
	if !IsEncrypted(encrypted) {
		t.Error("Encrypted string should be detected as encrypted")
	}

	// 空字符串
	if IsEncrypted("") {
		t.Error("Empty string should not be detected as encrypted")
	}
}

func TestDifferentEncryptions(t *testing.T) {
	c := NewCrypto("test-key")
	plaintext := "13800138000"

	// 同一明文多次加密应该产生不同密文（因为使用随机nonce）
	encrypted1, _ := c.Encrypt(plaintext)
	encrypted2, _ := c.Encrypt(plaintext)

	if encrypted1 == encrypted2 {
		t.Error("Different encryptions of same plaintext should produce different ciphertexts")
	}

	// 但解密结果应该相同
	decrypted1, _ := c.Decrypt(encrypted1)
	decrypted2, _ := c.Decrypt(encrypted2)

	if decrypted1 != decrypted2 {
		t.Error("Decrypted values should be the same")
	}
}
