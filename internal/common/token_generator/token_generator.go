package token_generator

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Username string
}

const (
	signingKey = "" // Your token
	tokenTTL   = time.Hour
)

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(tokenTTL)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		username,
	})

	return token.SignedString([]byte(signingKey))
}
