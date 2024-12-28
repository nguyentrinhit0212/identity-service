package models

import (
	"time"

	"github.com/google/uuid"
)

type FailedLogin struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"type:varchar(255);not null"`
	IPAddress    string    `gorm:"type:varchar(45);not null"`
	FailedAt     time.Time `gorm:"type:timestamp;default:current_timestamp"`
	AttemptCount int       `gorm:"default:1"`
	Reason       string    `gorm:"type:varchar(255)"`
}

func (FailedLogin) TableName() string {
	return "failed_logins"
}
