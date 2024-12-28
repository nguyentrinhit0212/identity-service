package services

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	jwtmanager "identity-service/internal/auth/jwt"
	"identity-service/internal/models"
	"identity-service/internal/repositories"
)

type TokenService interface {
	GenerateToken(userID string, tenantID string) (string, error)
	ValidateToken(token string) (*models.JWTToken, error)
	RevokeToken(tokenID string) error
}

type tokenService struct {
	tokenRepo  repositories.JwtTokenRepository
	keyManager *jwtmanager.KeyManager
}

func NewTokenService(tokenRepo repositories.JwtTokenRepository, keyManager *jwtmanager.KeyManager) TokenService {
	return &tokenService{
		tokenRepo:  tokenRepo,
		keyManager: keyManager,
	}
}

func (s *tokenService) GenerateToken(userID string, tenantID string) (string, error) {
	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}

	// Parse tenant ID
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return "", err
	}

	// Create claims
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		Audience:  jwt.ClaimStrings{tenantID},
	}

	// Create token with RS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.keyManager.GetCurrentKeyID()

	// Sign token
	tokenString, err := token.SignedString(s.keyManager.GetCurrentPrivateKey())
	if err != nil {
		return "", err
	}

	// Hash token for storage
	hash := sha256.Sum256([]byte(tokenString))
	hashString := hex.EncodeToString(hash[:])

	// Store token in database
	jwtToken := &models.JWTToken{
		UserID:    uid,
		TokenHash: hashString,
		IssuedAt:  now,
		ExpiresAt: now.Add(24 * time.Hour),
		IPAddress: "0.0.0.0", // Should be passed from the request context
		TenantID:  tid,
	}

	if err := s.tokenRepo.CreateToken(jwtToken); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *tokenService) ValidateToken(tokenString string) (*models.JWTToken, error) {
	// First parse without validation to get the key ID
	var (
		err         error
		parsedToken *jwt.Token
		ok          bool
		keyID       string
	)

	// Initial parse to get the key ID
	parsedToken, err = jwt.Parse(tokenString, nil)
	if err != nil {
		return nil, jwt.ErrSignatureInvalid
	}

	// Get key ID from token header
	if keyID, ok = parsedToken.Header["kid"].(string); !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	// Get the key pair for this key ID
	keyPair := s.keyManager.GetKeyPairByID(keyID)
	if keyPair == nil {
		return nil, jwt.ErrSignatureInvalid
	}

	// Parse and validate token
	validatedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		var isRSA bool
		if _, isRSA = token.Method.(*jwt.SigningMethodRSA); !isRSA {
			return nil, jwt.ErrSignatureInvalid
		}
		return keyPair.PublicKey, nil
	})

	if err != nil || !validatedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// Hash token for lookup
	hash := sha256.Sum256([]byte(tokenString))
	hashString := hex.EncodeToString(hash[:])

	// Get token from database
	jwtToken, err := s.tokenRepo.FindTokenByHash(hashString)
	if err != nil {
		return nil, err
	}

	// Check if token is revoked
	if jwtToken.IsRevoked {
		return nil, jwt.ErrSignatureInvalid
	}

	// Check if token is expired
	if time.Now().After(jwtToken.ExpiresAt) {
		return nil, jwt.ErrTokenExpired
	}

	return jwtToken, nil
}

func (s *tokenService) RevokeToken(tokenID string) error {
	return s.tokenRepo.DeleteToken(tokenID)
}
