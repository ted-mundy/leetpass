package signing_test

import (
	"testing"
	"time"

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
		Lifetime:   1 * time.Hour,
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

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("Failed to get exp from claims")
	}

	expTime := time.Unix(int64(exp), 0)
	if time.Until(expTime) > signer.Lifetime || time.Until(expTime) < signer.Lifetime-time.Minute { // allow 1 minute of leeway
		t.Fatalf("Expected exp to be within lifetime, got %v", time.Until(expTime))
	}
}
