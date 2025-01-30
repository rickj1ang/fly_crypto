package mail

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

// Send sends a verification code to the specified email address
func Send(email string, code string) error {
	// TODO: Implement email sending logic
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "rick0j1ang@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "fly_crypto_verify_code")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", code)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "rick0j1ang@gmail.com", "occe jexr cmkv tquz")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	err := d.DialAndSend(m)

	return err
}
