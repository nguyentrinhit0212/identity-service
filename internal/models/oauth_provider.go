package models

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
	Provider       string    `gorm:"type:varchar(50);not null"`
	ProviderUserID string    `gorm:"type:varchar(255);unique;not null"`
	CreatedAt      time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

func (OAuthProvider) TableName() string {
	return "oauth_providers"
}