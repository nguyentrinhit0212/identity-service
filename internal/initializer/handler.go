package initializer

import (
	"identity-service/internal/auth"
)

type Handlers struct {
    AuthHandler *auth.OAuthHandler
}

func InitHandlers(services *Services) *Handlers {
    providers := map[string]auth.OAuthProviderInterface{
        "google": auth.NewGoogleProvider(),
    }
    
	return &Handlers{
        AuthHandler: auth.NewOAuthHandler(providers, services.UserService),
    }
}