package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenSigningFailed = errors.New("failed to sign token")
	ErrInvalidTokenFormat = errors.New("invalid token format: expected 'Bearer <token>'")
	ErrExpiredToken       = errors.New("invalid or expired token")
)
