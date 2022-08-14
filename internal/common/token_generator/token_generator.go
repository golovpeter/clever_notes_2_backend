package token_generator

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Username string
}

const (
	tokenTTL        = time.Hour
	refreshTokenTTL = time.Hour * 720
)

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
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(refreshTokenTTL)},
	})

	return token.SignedString([]byte(os.Getenv("SIGNINKEY")))
}

func ValidateToken(inputToken string) error {
	_, err := jwt.Parse(inputToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNINKEY")), nil
	})

	if err != nil {
		return err
	}

	return nil
}
