package utils

import (
	"securebit/domain"
	"strings"
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

func Validate(token string, prefix string) (string, error) {
	if !strings.HasPrefix(token, prefix) {
		return "", domain.ErrInvalidTokenFormat
	}
	splitToken := strings.Split(token, " ")
	return splitToken[1], nil
}
