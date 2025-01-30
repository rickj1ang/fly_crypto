package token

import (
	"encoding/hex"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Test case 1: Basic token generation
	token1, err := Generate()
	if err != nil {
		t.Errorf("Generate() returned unexpected error: %v", err)
	}

	// Verify token length (32 bytes = 64 hex characters)
	if len(token1) != 64 {
		t.Errorf("Generate() returned token of length %d, want 64", len(token1))
	}

	// Verify token is valid hex string
	_, err = hex.DecodeString(token1)
	if err != nil {
		t.Errorf("Generate() returned invalid hex string: %v", err)
	}

	// Test case 2: Uniqueness check
	token2, err := Generate()
	if err != nil {
		t.Errorf("Generate() returned unexpected error: %v", err)
	}

	// Verify tokens are different
	if token1 == token2 {
		t.Error("Generate() returned identical tokens, expected unique tokens")
	}
}
