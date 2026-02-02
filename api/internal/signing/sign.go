package signing

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"crypto/rsa"
)

type Signer struct {
	PrivateKey *rsa.PrivateKey
	Lifetime   time.Duration
}

func (s *Signer) Sign(data any) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(s.Lifetime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(s.PrivateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
