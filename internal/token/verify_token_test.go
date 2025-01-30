package token

import (
	"testing"
)

func TestGenerateVerificationCode(t *testing.T) {
	// Test case 1: Basic code generation
	code1 := GenerateVerificationCode()

	// Verify code length is 6
	if len(code1) != 6 {
		t.Errorf("GenerateVerificationCode() returned code of length %d, want 6", len(code1))
	}

	// Verify code contains only digits
	for _, c := range code1 {
		if c < '0' || c > '9' {
			t.Errorf("GenerateVerificationCode() returned code containing non-digit character: %c", c)
		}
	}

	// Test case 2: Uniqueness check
	code2 := GenerateVerificationCode()

	// Verify codes are different (this might rarely fail due to random chance)
	if code1 == code2 {
		t.Log("Warning: Generated identical codes. This might happen rarely due to random chance")
	}

	// Verify code2 also meets the requirements
	if len(code2) != 6 {
		t.Errorf("GenerateVerificationCode() returned code of length %d, want 6", len(code2))
	}
	for _, c := range code2 {
		if c < '0' || c > '9' {
			t.Errorf("GenerateVerificationCode() returned code containing non-digit character: %c", c)
		}
	}
}
