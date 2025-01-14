// services/hash_services.go
package hash_services

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)


func GenerateSalt(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = charset[rand.Intn(len(charset))]
	}
	return string(salt)
}

func HashPasswordWithSalt(password, salt string) (string, error) {
	passwordWithSalt := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func VerifyPassword(password, salt, hashedPassword string) error {
	passwordWithSalt := password + salt
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordWithSalt))
}

