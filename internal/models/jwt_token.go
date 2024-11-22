package models

import (
	"time"

	"github.com/google/uuid"
)

type JWTToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
	TokenHash string    `gorm:"type:text;unique;not null"`
	IsRevoked bool      `gorm:"default:false"`
	IssuedAt  time.Time `gorm:"type:timestamp;not null"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null"`
	IPAddress string    `gorm:"type:varchar(45);not null"`
}

func (JWTToken) TableName() string {
	return "jwt_tokens"
}