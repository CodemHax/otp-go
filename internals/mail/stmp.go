package mail

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendOTPEmail(to, otp string) error {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || user == "" || pass == "" || from == "" {
		return fmt.Errorf("missing SMTP configuration environment variables")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP Code")

	m.SetBody(
		"text/plain",
		fmt.Sprintf(
			"Your verification code is %s.\n\nThis code expires in 5 minutes.",
			otp,
		),
	)

	d := gomail.NewDialer(host, port, user, pass)

	return d.DialAndSend(m)
}
