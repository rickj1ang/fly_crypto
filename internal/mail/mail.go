package mail

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func SendVerifyCode(email string, code string) error {
	subject := "fly_crypto_verify_code"
	content := fmt.Sprintf("Your code is %s input it in 6 min", code)

	return send(email, subject, content)
}

func SendNotify(message Message) error {
	subject := "fly_crypto_notify"
	content := fmt.Sprintf("Now the price of %s is touching %f please check", message.CoinSymbol, message.TargetPrice)

	return send(message.SendTo, subject, content)
}

// Send sends a verification code to the specified email address
func send(email string, subject string, content string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "rick0j1ang@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", content)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "rick0j1ang@gmail.com", "occe jexr cmkv tquz")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	err := d.DialAndSend(m)

	return err
}
