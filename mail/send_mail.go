package mail

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func SendMail(toEmail string, verificationCode string) error {
	// Namecheap SMTP ayarları
	smtpHost := "mail.privateemail.com"  // SMTP sunucu
	smtpPort := "587"                    // TLS için port

	// .env dosyasından kullanıcı adı ve şifreyi yükle
	godotenv.Load()
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), smtpHost)
	
	log.Println("SMTP Auth: ", auth);
	// E-posta içeriği
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: Verification Code\r\n" +
		"\r\n" +
		"Your verification code is: " + verificationCode + "\r\n")

	// E-posta gönderme
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "info@wordlingo.me", []string{toEmail}, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func GenerateVerificationCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000)) // 000000 - 999999 arası
}
