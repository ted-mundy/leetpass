package signing

import (
	"github.com/golang-jwt/jwt/v5"

	"crypto/rsa"
)

type Signer struct {
	PrivateKey *rsa.PrivateKey
}

func (s *Signer) Sign(data any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"data": data,
	})

	signedToken, err := token.SignedString(s.PrivateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
