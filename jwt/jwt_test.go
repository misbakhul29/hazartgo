package jwt

import (
	"testing"
	"time"
)

func TestJWT_SignAndVerify(t *testing.T) {
	jwtManager := New("super-secret-key-123")

	claims := MapClaims{
		"sub":   "user_123",
		"name":  "Budi Santoso",
		"roles": []string{"admin", "editor"},
	}

	token, err := jwtManager.Sign(claims, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	parsedClaims, err := jwtManager.Verify(token)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}

	if parsedClaims["sub"] != "user_123" {
		t.Errorf("expected sub user_123, got %v", parsedClaims["sub"])
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	jwtManager := New("super-secret-key-123")

	claims := MapClaims{"sub": "user_123"}
	token, err := jwtManager.Sign(claims, -1*time.Minute)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = jwtManager.Verify(token)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}
