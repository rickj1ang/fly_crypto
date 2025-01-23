package token

import (
	"errors"
	"crypto/rand"
	"fmt"

	"github.com/rick/fly_crypto/internal/data"
)

// GenerateVerificationCode generates a random 6-digit code
func GenerateVerificationCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = byte(rand.Intn(10) + '0')
	}
	return string(code)
}

// ValidateVerificationCode validates if the provided code matches the stored code for the given email
func ValidateVerificationCode(code, email string, app *data.App) error {
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
	storedCode, err := app.GetVerifyCodeByEmail(fmt.Sprintf("verify:%s", email))
	if err != nil {
		return errors.New("invalid or expired verification code")
	}

	// Compare the codes
	if storedCode != code {
		return errors.New("invalid verification code")
	}

	return nil
}