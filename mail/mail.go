package mail

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func sendMail(toEmail string, verificationCode string) error {
	smtpHost := "email-smtp.eu-north-1.amazonaws.com"
	smtpPort := "587"

	godotenv.Load()
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), smtpHost)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: Verification Code\r\n" +
		"\r\n" +
		"Your verification code is: " + verificationCode + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "yaar.unylmz@outlook.com", []string{toEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}



func generateVerificationCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000)) // 000000 - 999999 arası
}