package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email       string          `gorm:"type:varchar(255);unique;not null" json:"email"`
	Name        string          `gorm:"type:varchar(255);not null" json:"name"`
	Status      string          `gorm:"type:varchar(50);not null;default:'active'" json:"status"`
	Role        string          `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	Settings    json.RawMessage `gorm:"type:jsonb" json:"settings,omitempty"`
	MFAEnabled  bool            `json:"mfaEnabled" gorm:"default:false"`
	LastLoginAt time.Time       `json:"lastLoginAt"`
	CreatedAt   time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}

type UserUpdate struct {
	Email    *string          `json:"email,omitempty"`
	Name     *string          `json:"name,omitempty"`
	Status   *string          `json:"status,omitempty"`
	Role     *string          `json:"role,omitempty"`
	Settings *json.RawMessage `json:"settings,omitempty"`
}

type UserProfile struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null" json:"userId"`
	Avatar      string          `gorm:"type:varchar(255)" json:"avatar,omitempty"`
	Bio         string          `gorm:"type:text" json:"bio,omitempty"`
	Location    string          `gorm:"type:varchar(255)" json:"location,omitempty"`
	PhoneNumber string          `gorm:"type:varchar(50)" json:"phoneNumber,omitempty"`
	Preferences json.RawMessage `gorm:"type:jsonb" json:"preferences,omitempty"`
	CreatedAt   time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
	User        User            `gorm:"foreignKey:UserID" json:"user"`
}

func (User) TableName() string {
	return "users"
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
