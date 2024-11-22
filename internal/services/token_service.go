package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"
)

type TokenService interface {
    GenerateToken(userID string) (string, error)
    ValidateToken(token string) (*models.JWTToken, error)
    RevokeToken(tokenID string) error
}

type tokenService struct {
    tokenRepo repositories.JwtTokenRepository
}

func NewTokenService(tokenRepo repositories.JwtTokenRepository) TokenService {
    return &tokenService{
        tokenRepo: tokenRepo,
    }
}

func (s *tokenService) GenerateToken(userID string) (string, error) {
    // Logic to generate token
    return "", nil
}

func (s *tokenService) ValidateToken(token string) (*models.JWTToken, error) {
    // Logic to validate token
    return nil, nil
}

func (s *tokenService) RevokeToken(tokenID string) error {
    return s.tokenRepo.DeleteToken(tokenID)
}