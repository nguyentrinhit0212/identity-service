package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserTenantAccess struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null" json:"tenantId"`
	Roles       pq.StringArray `gorm:"type:text[]" json:"roles"`
	Permissions pq.StringArray `gorm:"type:text[]" json:"permissions"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	Tenant      Tenant         `gorm:"foreignKey:TenantID" json:"tenant"`
}

func (UserTenantAccess) TableName() string {
	return "user_tenant_access"
}
