package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

// GenerateAccessToken creates a PASETO v4 local token with userID as subject.
func GenerateAccessToken(userID string, secret string, expiry time.Duration) (string, error) {
	token := paseto.NewToken()
	token.SetSubject(userID)
	token.SetIssuedAt(time.Now().UTC())
	token.SetExpiration(time.Now().UTC().Add(expiry))
	token.SetNotBefore(time.Now().UTC())

	key, err := symmetricKeyFromSecret(secret)
	if err != nil {
		return "", fmt.Errorf("failed to create symmetric key: %w", err)
	}

	return token.V4Encrypt(key, nil), nil
}

// ValidateAccessToken validates a PASETO v4 local token and returns the userID.
func ValidateAccessToken(tokenString string, secret string) (string, error) {
	key, err := symmetricKeyFromSecret(secret)
	if err != nil {
		return "", fmt.Errorf("failed to create symmetric key: %w", err)
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now().UTC()))

	token, err := parser.ParseV4Local(key, tokenString, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	userID, err := token.GetSubject()
	if err != nil {
		return "", fmt.Errorf("failed to get subject: %w", err)
	}

	return userID, nil
}

// GenerateRefreshToken creates a random 64-character hex string.
func GenerateRefreshToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func symmetricKeyFromSecret(secret string) (paseto.V4SymmetricKey, error) {
	secretBytes := []byte(secret)
	if len(secretBytes) < 32 {
		return paseto.V4SymmetricKey{}, fmt.Errorf("secret must be at least 32 bytes")
	}
	var keyBytes [32]byte
	copy(keyBytes[:], secretBytes[:32])
	return paseto.V4SymmetricKeyFromBytes(keyBytes[:])
}
