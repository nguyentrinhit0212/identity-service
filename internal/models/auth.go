package models

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	TenantID     uuid.UUID `json:"tenantId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
	LastUsedAt   time.Time `json:"lastUsedAt"`
	IPAddress    string    `json:"ipAddress"`
	UserAgent    string    `json:"userAgent"`
}

// OAuthState represents the state of an OAuth flow
type OAuthState struct {
	ID        uuid.UUID `json:"id"`
	State     string    `json:"state"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// OAuthToken represents an OAuth token
type OAuthToken struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	Provider     string    `json:"provider"`
	AccessToken  string    `json:"-"`
	TokenType    string    `json:"tokenType"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// OAuthProfile represents a user's OAuth profile
type OAuthProfile struct {
	ID        string            `json:"id"`
	Provider  string            `json:"provider"`
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	Picture   string            `json:"picture"`
	Raw       map[string]string `json:"raw"`
	CreatedAt time.Time         `json:"createdAt"`
}

// MFADevice represents a multi-factor authentication device
type MFADevice struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Secret    string    `json:"-"`
	Verified  bool      `json:"verified"`
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Token     string     `json:"-"`
	Used      bool       `json:"used"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

// EmailVerification represents an email verification request
type EmailVerification struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Email      string     `json:"email"`
	Token      string     `json:"-"`
	Verified   bool       `json:"verified"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// LoginCredentials represents a user's login credentials
type LoginCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type PKCEChallenge struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	CodeChallenge string    `json:"codeChallenge" gorm:"not null"`
	CodeVerifier  string    `json:"codeVerifier" gorm:"not null"`
	UserID        uuid.UUID `json:"userId" gorm:"type:uuid;not null"`
	TenantID      uuid.UUID `json:"tenantId" gorm:"type:uuid;not null"`
	ExpiresAt     time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt     time.Time `json:"createdAt" gorm:"not null"`
	Used          bool      `json:"used" gorm:"not null;default:false"`
}

type TokenExchangeRequest struct {
	Code         string `json:"code" binding:"required"`
	CodeVerifier string `json:"codeVerifier" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
