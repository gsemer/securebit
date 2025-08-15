package utils

import (
	"securebit/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignedToken(user domain.User, duration time.Time, secretKey string) (string, error) {
	claims := domain.Claims{
		Username: user.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(duration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", domain.ErrTokenSigningFailed
	}

	return signedToken, nil
}
