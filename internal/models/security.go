package models

import (
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID   uuid.UUID  `gorm:"type:uuid;not null" json:"tenantId"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name"`
	Key        string     `gorm:"type:varchar(255);not null" json:"key"`
	ExpiresAt  *time.Time `gorm:"type:timestamp" json:"expiresAt,omitempty"`
	LastUsedAt *time.Time `gorm:"type:timestamp" json:"lastUsedAt,omitempty"`
	CreatedAt  time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	RevokedAt  *time.Time `gorm:"type:timestamp" json:"revokedAt,omitempty"`
}

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null" json:"tenantId"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"userId,omitempty"`
	Action    string    `gorm:"type:varchar(255);not null" json:"action"`
	Resource  string    `gorm:"type:varchar(255);not null" json:"resource"`
	Details   string    `gorm:"type:text" json:"details,omitempty"`
	IP        string    `gorm:"type:varchar(45)" json:"ip,omitempty"`
	UserAgent string    `gorm:"type:varchar(255)" json:"userAgent,omitempty"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

type SecurityPolicies struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID              uuid.UUID `gorm:"type:uuid;not null" json:"tenantId"`
	RequireMFA            bool      `gorm:"type:boolean;default:false" json:"requireMFA"`
	PasswordMinLength     int       `gorm:"type:integer;default:8" json:"passwordMinLength"`
	PasswordRequireUpper  bool      `gorm:"type:boolean;default:true" json:"passwordRequireUpper"`
	PasswordRequireLower  bool      `gorm:"type:boolean;default:true" json:"passwordRequireLower"`
	PasswordRequireNumber bool      `gorm:"type:boolean;default:true" json:"passwordRequireNumber"`
	PasswordRequireSymbol bool      `gorm:"type:boolean;default:true" json:"passwordRequireSymbol"`
	SessionTimeout        int       `gorm:"type:integer;default:3600" json:"sessionTimeout"` // in seconds
	MaxLoginAttempts      int       `gorm:"type:integer;default:5" json:"maxLoginAttempts"`
	LockoutDuration       int       `gorm:"type:integer;default:300" json:"lockoutDuration"` // in seconds
	UpdatedAt             time.Time `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}

type SecurityMetrics struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID      uuid.UUID `gorm:"type:timestamp;not null" json:"tenantId"`
	LoginAttempts int       `gorm:"type:integer;default:0" json:"loginAttempts"`
	FailedLogins  int       `gorm:"type:integer;default:0" json:"failedLogins"`
	MFAUsage      int       `gorm:"type:integer;default:0" json:"mfaUsage"`
	APIKeyUsage   int       `gorm:"type:integer;default:0" json:"apiKeyUsage"`
	SuspiciousIPs []string  `gorm:"type:text[]" json:"suspiciousIPs"`
	BlockedIPs    []string  `gorm:"type:text[]" json:"blockedIPs"`
	LastUpdated   time.Time `gorm:"type:timestamp;default:current_timestamp"`
}

type SecurityAlert struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID    uuid.UUID  `gorm:"type:uuid;not null" json:"tenantId"`
	Type        string     `gorm:"type:varchar(50);not null" json:"type"`
	Severity    string     `gorm:"type:varchar(20);not null" json:"severity"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Status      string     `gorm:"type:varchar(20);not null;default:'open'" json:"status"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
	ResolvedAt  *time.Time `gorm:"type:timestamp" json:"resolvedAt,omitempty"`
}

type SecuritySettings struct {
	MFAEnabled bool      `json:"mfaEnabled"`
	LastLogin  time.Time `json:"lastLogin"`
}
