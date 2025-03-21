package mail

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/joho/godotenv"
)


var dialer *mail.Dialer

func InitMailer() {
	godotenv.Load()

	dialer = mail.NewDialer(
		"mail.privateemail.com", // SMTP Host
		587,                     // SMTP Port
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS
}


func SendMail(toEmail, verificationCode string) error {
	if dialer == nil {
		InitMailer()
	}

  m := mail.NewMessage()
	m.SetHeader("From", "info@wordlingo.me")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
<html>
  <head>
    <title>Wordlingo Verification</title>
  </head>
  <body
    style="
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: Arial, sans-serif;
      background-color: #f6f6f6;
      margin: 0;
      padding: 0;
    "
  >
    <div
      style="
        width: 100%;
        max-width: 600px;
        margin: 0 auto;
        background-color: #ffffff;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      "
    >
      <div
        style="
          background-color: #4f42d8;
          width: 100%;

          color: #ffffff;
          padding: 20px 0;
          text-align: center;
          border-radius: 8px 8px 0 0;
        "
      >
        <h1 style="margin: 0; font-size: 24px">Welcome to Wordlingo!</h1>
      </div>
      <div style="padding: 20px">
        <p
          style="
            text-align: center;
            font-size: 16px;
            color: #333;
            line-height: 1.5;
          "
        >
        Your verification code is:
        </p>
        <div
          style="
            font-size: 32px;
            font-weight: bold;
            color: #4f42d8;
            text-align: center;
            margin: 10px 0;
          "
        >
          ` + verificationCode + `
          
        </div>
        <p
          style="
            text-align: center;
            font-size: 16px;
            color: #333;
            line-height: 1.5;
          "
        >
        Please use this code within the next 10 minutes to verify your email address.
        </p>
        <a
          href="https://www.wordlingo.me/verify"
          style="
            display: block;
            width: 90%;
            background-color: #232f3e;
            color: #ffffff;
            padding: 15px;
            text-align: center;
            font-size: 18px;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
            margin-left: 15px;
          "
        >
          Verify Your Account
        </a>
      </div>
      <div
        style="
          font-size: 12px;
          color: #777;
          text-align: center;
          margin-top: 10px;
          border-top: 1px solid #ddd;
          padding-top: 20px;
        "
      >
        <p>Wordlingo &copy; 2024</p>
        <p>
          Need help? Visit our
          <a
            href="https://www.wordlingo.me/help"
            style="color: #555; text-decoration: none"
            >help center</a
          >
        </p>
      </div>
    </div>
  </body>
</html>
`, verificationCode))


err := dialer.DialAndSend(m)
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
