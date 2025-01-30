package token

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate creates a new random token for the given email
func Generate() (string, error) {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Convert to hex string
	return hex.EncodeToString(bytes), nil
}
