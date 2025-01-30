package token

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/rickj1ang/fly_crypto/internal/app"
)

// GenerateVerificationCode generates a random 6-digit code
func GenerateVerificationCode() string {
	// Generate 2 random bytes (16 bits) which is enough for a 6-digit number
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		// In case of error, return a default code
		return "000000"
	}
	// Convert bytes to uint16 and ensure it's between 100000 and 999999
	num := (uint16(b[0]) << 8) | uint16(b[1])
	code := int(num)%900000 + 100000
	return fmt.Sprintf("%06d", code)
}

// ValidateVerificationCode validates if the provided code matches the stored code for the given email
func ValidateVerificationCode(app app.App, code, email string) error {
	// First validate the format
	if len(code) != 6 {
		return errors.New("verification code must be 6 digits")
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			return errors.New("verification code must contain only digits")
		}
	}

	// Get the stored code from Redis
	realEmail, err := app.Data.GetEmailByVerificationCode(fmt.Sprintf("verify:%s", code))
	if err != nil {
		return errors.New("invalid or expired verification code")
	}

	// Compare the codes
	if realEmail != email {
		return errors.New("invalid verification code")
	}

	return nil
}
