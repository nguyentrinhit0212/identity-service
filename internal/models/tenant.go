package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenantType string

const (
	PersonalTenant   TenantType = "personal"
	TeamTenant       TenantType = "team"
	EnterpriseTenant TenantType = "enterprise"
)

type AuthProvider struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	DisplayName string          `json:"displayName"`
	Config      json.RawMessage `json:"config"`
}

// AuthProviders is a custom type for handling JSON serialization of []AuthProvider
type AuthProviders []AuthProvider

// Scan implements the sql.Scanner interface for AuthProviders
func (ap *AuthProviders) Scan(value interface{}) error {
	if value == nil {
		*ap = make([]AuthProvider, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), ap)
	}
	return json.Unmarshal(bytes, ap)
}

// Value implements the driver.Valuer interface for AuthProviders
func (ap AuthProviders) Value() (driver.Value, error) {
	if ap == nil {
		return json.Marshal([]AuthProvider{})
	}
	return json.Marshal(ap)
}

type Tenant struct {
	ID                    uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Slug                  string          `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Name                  string          `gorm:"type:varchar(255);not null" json:"name"`
	Type                  TenantType      `gorm:"type:tenant_type;not null;default:'team'" json:"type"`
	Domain                string          `gorm:"type:varchar(255)" json:"domain,omitempty"`
	DomainVerified        bool            `gorm:"type:boolean;default:false" json:"domainVerified"`
	OwnerID               *uuid.UUID      `gorm:"type:uuid;references:users(id)" json:"ownerId,omitempty"`
	MaxUsers              *int            `gorm:"type:integer" json:"maxUsers,omitempty"`
	AuthProviders         AuthProviders   `gorm:"type:jsonb" json:"authProviders"`
	Features              json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"features,omitempty"`
	Settings              json.RawMessage `gorm:"type:jsonb" json:"settings,omitempty"`
	SubscriptionStatus    string          `gorm:"type:varchar(50)" json:"subscriptionStatus,omitempty"`
	SubscriptionPlan      string          `gorm:"type:varchar(50)" json:"subscriptionPlan,omitempty"`
	SubscriptionExpiresAt *time.Time      `gorm:"type:timestamp" json:"subscriptionExpiresAt,omitempty"`
	UsageStats            json.RawMessage `gorm:"type:jsonb" json:"usageStats,omitempty"`
	CreatedAt             time.Time       `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt             time.Time       `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}

// TenantUpdate represents the fields that can be updated in a tenant
type TenantUpdate struct {
	Name                  *string          `json:"name,omitempty"`
	Type                  *TenantType      `json:"type,omitempty"`
	Domain                *string          `json:"domain,omitempty"`
	DomainVerified        *bool            `json:"domainVerified,omitempty"`
	OwnerID               *uuid.UUID       `json:"ownerId,omitempty"`
	MaxUsers              *int             `json:"maxUsers,omitempty"`
	AuthProviders         *[]AuthProvider  `json:"authProviders,omitempty"`
	Features              *json.RawMessage `json:"features,omitempty"`
	Settings              *json.RawMessage `json:"settings,omitempty"`
	SubscriptionStatus    *string          `json:"subscriptionStatus,omitempty"`
	SubscriptionPlan      *string          `json:"subscriptionPlan,omitempty"`
	SubscriptionExpiresAt *time.Time       `json:"subscriptionExpiresAt,omitempty"`
}

// TenantUpgrade represents a tenant subscription upgrade request
type TenantUpgrade struct {
	Plan      string     `json:"plan"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// TenantInvite represents an invitation to join a tenant
type TenantInvite struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID  uuid.UUID  `gorm:"type:uuid;not null" json:"tenantId"`
	Email     string     `gorm:"type:varchar(255);not null" json:"email"`
	Role      string     `gorm:"type:varchar(50);not null" json:"role"`
	Status    string     `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	ExpiresAt *time.Time `gorm:"type:timestamp" json:"expiresAt,omitempty"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"updatedAt"`
}

type TenantFeatures struct {
	EnabledFeatures []string        `json:"enabledFeatures"`
	FeatureSettings json.RawMessage `json:"featureSettings,omitempty"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type TenantSettings struct {
	GeneralSettings  json.RawMessage `json:"generalSettings,omitempty"`
	SecuritySettings json.RawMessage `json:"securitySettings,omitempty"`
	CustomSettings   json.RawMessage `json:"customSettings,omitempty"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// Helper methods for feature checking
func (t *Tenant) HasFeatures() map[string]bool {
	var features map[string]bool
	if err := json.Unmarshal(t.Features, &features); err != nil {
		return make(map[string]bool)
	}
	return features
}

func (t *Tenant) UpdateFeatures(features map[string]bool) error {
	featuresJSON, err := json.Marshal(features)
	if err != nil {
		return err
	}
	t.Features = featuresJSON
	return nil
}

// Helper methods for subscription
func (t *Tenant) IsSubscriptionActive() bool {
	if t.Type == PersonalTenant {
		return true // Personal tenants are always active
	}
	if t.SubscriptionExpiresAt == nil {
		return false
	}
	return t.SubscriptionExpiresAt.After(time.Now())
}

// Helper methods for usage stats
func (t *Tenant) GetUsageStats() map[string]interface{} {
	var stats map[string]interface{}
	if err := json.Unmarshal(t.UsageStats, &stats); err != nil {
		return make(map[string]interface{})
	}
	return stats
}

func (t *Tenant) UpdateUsageStats(stats map[string]interface{}) error {
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	t.UsageStats = statsJSON
	return nil
}

// Helper methods for settings
func (t *Tenant) GetSettings() map[string]interface{} {
	var settings map[string]interface{}
	if err := json.Unmarshal(t.Settings, &settings); err != nil {
		return make(map[string]interface{})
	}
	return settings
}

func (t *Tenant) UpdateSettings(settings map[string]interface{}) error {
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	t.Settings = settingsJSON
	return nil
}
