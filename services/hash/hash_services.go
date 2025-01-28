// services/hash_services.go
package hash_services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 parametreleri (Güvenlik ve performans dengesi için ayarlayın)
const (
	argon2Time    = 1      // CPU/maliyet faktörü (iterasyon sayısı)
	argon2Memory  = 64 * 1024 // 64MB bellek kullanımı
	argon2Threads = 4      // Paralel thread sayısı
	argon2KeyLen  = 32     // Üretilecek hash uzunluğu (32 byte = 256 bit)
	saltLength    = 16     // Salt uzunluğu (16 byte = 128 bit)
)

// GenerateSecureSalt - Kriptografik olarak güvenli rastgele salt üretir
func GenerateSecureSalt() (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt üretilemedi: %v", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword - Şifreyi Argon2 ile hash'ler ve hash+salt'ı birleştirir
func HashPassword(password string) (string, error) {
	// Salt üret
	saltBase64, err := GenerateSecureSalt()
	if err != nil {
		return "", err
	}
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return "", fmt.Errorf("salt decode hatası: %v", err)
	}

	// Argon2ID ile hash'le
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Hash ve salt'ı birleştir (format: hash:salt)
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s:%s", encodedHash, saltBase64), nil
}

// VerifyPassword - Hash'i doğrular
func VerifyPassword(password, encodedHashWithSalt string) bool {
	// Hash ve salt'ı ayır
	parts := strings.Split(encodedHashWithSalt, ":")
	if len(parts) != 2 {
		return false
	}

	encodedHash := parts[0]
	saltBase64 := parts[1]

	// Decode işlemleri
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false
	}

	storedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}

	// Girilen şifreyi hash'le
	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Timing attack korumalı karşılaştırma
	return subtle.ConstantTimeCompare(storedHash, newHash) == 1
}