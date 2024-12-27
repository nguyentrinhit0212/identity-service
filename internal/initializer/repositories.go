package initializer

import (
	"identity-service/db"
	"identity-service/internal/repositories"
)

// Repositories contains all data access repositories
type Repositories struct {
	UserRepo       repositories.UserRepository
	TenantRepo     repositories.TenantRepository
	SecurityRepo   repositories.SecurityRepository
	SessionRepo    repositories.SessionRepository
	PKCERepository repositories.PKCERepository
}

// InitRepositories initializes all repositories with database connections
func InitRepositories() *Repositories {
	database := repositories.WrapDB(db.GetDB())
	return &Repositories{
		UserRepo:       repositories.NewUserRepository(database),
		TenantRepo:     repositories.NewTenantRepository(database),
		SecurityRepo:   repositories.NewSecurityRepository(database),
		SessionRepo:    repositories.NewSessionRepository(database),
		PKCERepository: repositories.NewPKCERepository(database),
	}
}
