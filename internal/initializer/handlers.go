package initializer

import (
	"identity-service/internal/auth"
	"identity-service/internal/handlers"
	"identity-service/internal/services"
	"log"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	AuthHandler     *handlers.AuthHandler
	OAuthHandler    *handlers.OAuthHandler
	UserHandler     *handlers.UserHandler
	TenantHandler   *handlers.TenantHandler
	SecurityHandler *handlers.SecurityHandler
}

// InitHandlers initializes all handlers with their required services
func InitHandlers(s *Services) *Handlers {
	// Initialize OAuth providers
	providers := make(map[string]services.OAuthProvider)

	// Initialize Google OAuth provider
	googleProvider := auth.NewGoogleProvider()
	providers["google"] = googleProvider
	log.Printf("Initialized Google OAuth provider")

	return &Handlers{
		AuthHandler:     handlers.NewAuthHandler(s.AuthService),
		OAuthHandler:    handlers.NewOAuthHandler(providers, s.UserService, s.AuthService),
		UserHandler:     handlers.NewUserHandler(s.UserService),
		TenantHandler:   handlers.NewTenantHandler(s.TenantService),
		SecurityHandler: handlers.NewSecurityHandler(s.SecurityService),
	}
}
