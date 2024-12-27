package initializer

import (
	"identity-service/internal/auth/jwt"
	"identity-service/internal/services"
	"log"
	"time"
)

const (
	defaultTokenDuration     = 24 * time.Hour
	defaultKeySize           = 2048
	defaultKeyRotationPeriod = 24 * time.Hour
)

// Services holds all service instances
type Services struct {
	AuthService     services.AuthService
	UserService     services.UserService
	TenantService   services.TenantService
	SecurityService services.SecurityService
	PKCEService     services.PKCEService
	keyManager      *jwt.KeyManager
}

// InitServices initializes all services with their required repositories
func InitServices(repos *Repositories) *Services {
	userService := services.NewUserService(repos.UserRepo, repos.TenantRepo)

	// Initialize key manager with default settings
	keyManager, err := jwt.NewKeyManager(defaultKeyRotationPeriod, defaultKeySize)
	if err != nil {
		log.Fatalf("Failed to initialize key manager: %v", err)
	}

	return &Services{
		AuthService:     services.NewAuthService(userService, repos.SessionRepo, keyManager),
		UserService:     userService,
		TenantService:   services.NewTenantService(repos.TenantRepo),
		SecurityService: services.NewSecurityService(repos.SecurityRepo),
		PKCEService:     services.NewPKCEService(repos.PKCERepository),
		keyManager:      keyManager,
	}
}

// GetKeyManager returns the JWT key manager instance
func (s *Services) GetKeyManager() *jwt.KeyManager {
	return s.keyManager
}
