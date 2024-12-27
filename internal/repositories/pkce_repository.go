package repositories

import (
	"identity-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type PKCERepository interface {
	CreateChallenge(challenge *models.PKCEChallenge) error
	GetChallenge(id uuid.UUID) (*models.PKCEChallenge, error)
	MarkChallengeAsUsed(id uuid.UUID) error
}

type pkceRepository struct {
	db GormDB
}

func NewPKCERepository(db GormDB) PKCERepository {
	return &pkceRepository{
		db: db,
	}
}

func (r *pkceRepository) CreateChallenge(challenge *models.PKCEChallenge) error {
	return r.db.Create(challenge).Error
}

func (r *pkceRepository) GetChallenge(id uuid.UUID) (*models.PKCEChallenge, error) {
	var challenge models.PKCEChallenge
	if err := r.db.Where("id = ? AND used = ? AND expires_at > ?", id, false, time.Now()).First(&challenge).Error; err != nil {
		return nil, err
	}
	return &challenge, nil
}

func (r *pkceRepository) MarkChallengeAsUsed(id uuid.UUID) error {
	return r.db.Model(&models.PKCEChallenge{}).Where("id = ?", id).Update("used", true).Error
}
