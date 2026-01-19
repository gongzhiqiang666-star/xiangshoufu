package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://apple@localhost:5432/xiangshoufu?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// 计算正确的密码哈希
	password := "admin123"
	salt := "default_salt_12345678"
	hash := sha256.Sum256([]byte(password + salt))
	hashedPassword := hex.EncodeToString(hash[:])

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Salt: %s\n", salt)
	fmt.Printf("Hash: %s\n", hashedPassword)

	// 更新管理员密码
	result := db.Exec("UPDATE users SET password = ? WHERE username = ?", hashedPassword, "admin")
	if result.Error != nil {
		log.Fatalf("Failed to update: %v", result.Error)
	}

	fmt.Printf("Updated %d rows\n", result.RowsAffected)

	if result.RowsAffected == 0 {
		// 如果没有 admin 用户，创建一个
		fmt.Println("No admin user found, creating one...")
		result = db.Exec(`INSERT INTO users (username, password, salt, agent_id, role_type, status)
			VALUES (?, ?, ?, NULL, 2, 1)`, "admin", hashedPassword, salt)
		if result.Error != nil {
			log.Fatalf("Failed to insert: %v", result.Error)
		}
		fmt.Printf("Created admin user\n")
	}

	fmt.Println("Done!")
}
