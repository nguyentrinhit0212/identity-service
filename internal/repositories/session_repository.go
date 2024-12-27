package repositories

import (
	"identity-service/internal/models"

	"github.com/google/uuid"
)

type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSession(id uuid.UUID) (*models.Session, error)
	UpdateSession(session *models.Session) error
	DeleteSession(id uuid.UUID) error
	ListUserSessions(userID uuid.UUID) ([]*models.Session, error)
}

type sessionRepository struct {
	db GormDB
}

func NewSessionRepository(db GormDB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) GetSession(id uuid.UUID) (*models.Session, error) {
	var session models.Session
	if err := r.db.First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) UpdateSession(session *models.Session) error {
	return r.db.Save(session).Error
}

func (r *sessionRepository) DeleteSession(id uuid.UUID) error {
	return r.db.Delete(&models.Session{}, "id = ?", id).Error
}

func (r *sessionRepository) ListUserSessions(userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	if err := r.db.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
