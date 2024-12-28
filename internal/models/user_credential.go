package models

import (
	"time"

	"github.com/google/uuid"
)

type UserCredential struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
	PasswordHash string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}

func (UserCredential) TableName() string {
	return "user_credentials"
}
