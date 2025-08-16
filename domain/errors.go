package domain

import "errors"

var (
	ErrTokenSigningFailed = errors.New("failed to sign token")
)
