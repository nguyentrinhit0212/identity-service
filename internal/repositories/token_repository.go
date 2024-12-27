package repositories

import (
	"identity-service/internal/models"
)

type JwtTokenRepository interface {
	CreateToken(token *models.JWTToken) error
	FindTokenByID(id string) (*models.JWTToken, error)
	FindTokenByHash(hash string) (*models.JWTToken, error)
	DeleteToken(id string) error
}

type jwtTokenRepository struct {
	db GormDB
}

func NewJwtTokenRepository(db GormDB) JwtTokenRepository {
	return &jwtTokenRepository{
		db: db,
	}
}

func (r *jwtTokenRepository) CreateToken(token *models.JWTToken) error {
	return r.db.Create(token).Error
}

func (r *jwtTokenRepository) FindTokenByID(id string) (*models.JWTToken, error) {
	var token models.JWTToken
	err := r.db.Where("id = ?", id).First(&token).Error
	return &token, err
}

func (r *jwtTokenRepository) FindTokenByHash(hash string) (*models.JWTToken, error) {
	var token models.JWTToken
	err := r.db.Where("token_hash = ?", hash).First(&token).Error
	return &token, err
}

func (r *jwtTokenRepository) DeleteToken(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.JWTToken{}).Error
}
