package auth

import "errors"

// Common authentication errors
var (
	ErrInvalidToken = errors.New("invalid token")
)
