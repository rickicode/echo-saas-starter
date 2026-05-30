package auth

import (
	"testing"
	"time"
)

const testSecret = "this-is-a-test-secret-key-32bytes!"

func TestGenerateAndValidateAccessToken(t *testing.T) {
	userID := "user-123"
	token, err := GenerateAccessToken(userID, testSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// Validate the token
	gotUserID, err := ValidateAccessToken(token, testSecret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if gotUserID != userID {
		t.Errorf("expected userID %q, got %q", userID, gotUserID)
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	userID := "user-456"
	token, err := GenerateAccessToken(userID, testSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	wrongSecret := "wrong-secret-key-that-is-32bytes!!"
	_, err = ValidateAccessToken(token, wrongSecret)
	if err == nil {
		t.Fatal("expected error when validating with wrong secret")
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	userID := "user-789"
	// Generate a token that expires immediately (negative duration)
	token, err := GenerateAccessToken(userID, testSecret, -1*time.Second)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = ValidateAccessToken(token, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestGenerateAccessToken_ShortSecret(t *testing.T) {
	_, err := GenerateAccessToken("user", "short", 15*time.Minute)
	if err == nil {
		t.Fatal("expected error for short secret")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token1 := GenerateRefreshToken()
	token2 := GenerateRefreshToken()

	if len(token1) != 64 {
		t.Errorf("expected refresh token length 64, got %d", len(token1))
	}
	if token1 == token2 {
		t.Error("expected unique refresh tokens")
	}
}
