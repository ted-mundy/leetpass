package signing_test

import (
	"testing"

	"crypto/rand"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ted-mundy/leetpass-api/internal/signing"
)

func TestSign(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	signer := &signing.Signer{
		PrivateKey: key,
	}

	data := "test-data"
	token, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if len(token) == 0 {
		t.Fatalf("Expected non-empty token")
	}

	// ensure that the token has been signed with the correct key
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return &key.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatalf("Expected token to be valid")
	}

	// ensure the content is correct
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("Failed to get claims from token")
	}

	if string(claims["data"].(string)) != string(data) {
		t.Fatalf("Expected data %s, got %s", data, claims["data"])
	}
}
