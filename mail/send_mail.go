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
	smtpHost := "mail.privateemail.com"
	smtpPort := "587"

	godotenv.Load()
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), smtpHost)

	htmlContent := `
	<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f6f6f6;
            margin: 0;
            padding: 0;
        }
        .email-container {
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
        .email-header {
            background-color: #232f3e;
            color: #ffffff;
            padding: 20px 0;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .email-header h1 {
            margin: 0;
            font-size: 24px;
        }
        .email-body {
            padding: 20px;
        }
        .email-body span {
            display: block;
            text-align: center;
            font-size: 16px;
            color: #333;
            line-height: 1.5;
        }
        .verification-code {
            font-size: 32px;
            font-weight: bold;
            color: #4f42d8;
            text-align: center;
            margin: 30px 0;
        }
        .cta-button {
            display: block;
            width: 100%;
            background-color: #4f42d8;
            color: #ffffff;
            padding: 15px;
            text-align: center;
            font-size: 18px;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
        }
        .cta-button:hover {
            background-color: #4f42d8;
        }
        .email-footer {
            font-size: 12px;
            color: #777;
            text-align: center;
            margin-top: 30px;
            border-top: 1px solid #ddd;
            padding-top: 20px;
        }
        .email-footer a {
            color: #555;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-header">
            <h1>Welcome to Wordlingo!</h1>
        </div>
        <div class="email-body">
            <span>Thank you for registering. Your verification code is:</span>
            <div class="verification-code">
                ` + verificationCode + `
            </div>
            <span>If you did not request this, please ignore this email.</span>
            <a href="https://www.wordlingo.me/verify" class="cta-button">Verify Your Account</a>
        </div>
        <div class="email-footer">
            <span>Wordlingo &copy; 2024</span>
            <span>Need help? Visit our <a href="https://www.wordlingo.me/help">help center</a>.</span>
        </div>
    </div>
</body>
</html>
`

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: Verification Code\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		htmlContent + "\r\n")

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
