package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"

	"github.com/google/uuid"
)

type PKCEService interface {
	CreateChallenge(challenge *models.PKCEChallenge) error
	GetChallenge(id uuid.UUID) (*models.PKCEChallenge, error)
	MarkChallengeAsUsed(id uuid.UUID) error
}

type pkceService struct {
	repo repositories.PKCERepository
}

func NewPKCEService(repo repositories.PKCERepository) PKCEService {
	return &pkceService{
		repo: repo,
	}
}

func (s *pkceService) CreateChallenge(challenge *models.PKCEChallenge) error {
	return s.repo.CreateChallenge(challenge)
}

func (s *pkceService) GetChallenge(id uuid.UUID) (*models.PKCEChallenge, error) {
	return s.repo.GetChallenge(id)
}

func (s *pkceService) MarkChallengeAsUsed(id uuid.UUID) error {
	return s.repo.MarkChallengeAsUsed(id)
}
