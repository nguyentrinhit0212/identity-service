package models

import (
	"time"

	"github.com/google/uuid"
)

type IPWhitelist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
	IPAddress string    `gorm:"type:inet;unique;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

func (IPWhitelist) TableName() string {
	return "ip_whitelists"
}
