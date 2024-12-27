package utils

import (
	"errors"
	"time"

	jwtmanager "identity-service/internal/auth/jwt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

var keyManager *jwtmanager.KeyManager

// SetKeyManager sets the JWT key manager instance
func SetKeyManager(km *jwtmanager.KeyManager) {
	keyManager = km
}

func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	if keyManager == nil {
		return "", errors.New("JWT manager not initialized")
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "identity-service",
		},
	}

	return keyManager.SignToken(claims)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	if keyManager == nil {
		return nil, errors.New("JWT manager not initialized")
	}

	claims := &Claims{}
	if err := keyManager.VerifyToken(tokenString, claims); err != nil {
		return nil, err
	}

	return claims, nil
}
