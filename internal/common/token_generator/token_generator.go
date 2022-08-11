package token_generator

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Username string
}

const tokenTTL = time.Hour

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(tokenTTL)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		username,
	})

	return token.SignedString([]byte(os.Getenv("SIGNINKEY")))
}

func GenerateRefreshJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 48)},
	})

	return token.SignedString([]byte(os.Getenv("SIGNINKEY")))
}
