package services

import (
	"context"
	"fmt"
	"identity-service/internal/auth"
	"identity-service/internal/models"
	"log"
)

// OAuthProvider defines the interface for OAuth providers
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeToken(ctx context.Context, code string) (string, error)
	FetchUserInfo(token string) (*models.OAuthUser, error)
	GetProviderName() string
}

// OAuthService handles multiple OAuth providers
type OAuthService struct {
	providers map[string]OAuthProvider
}

// NewOAuthService creates a new OAuth service with registered providers
func NewOAuthService() *OAuthService {
	service := &OAuthService{
		providers: make(map[string]OAuthProvider),
	}

	// Register providers
	service.RegisterProvider("google", auth.NewGoogleProvider())
	// Add more providers here
	// service.RegisterProvider("github", auth.NewGithubProvider())
	// service.RegisterProvider("microsoft", auth.NewMicrosoftProvider())

	return service
}

// RegisterProvider adds a new OAuth provider to the service
func (s *OAuthService) RegisterProvider(name string, provider OAuthProvider) {
	log.Printf("Registering OAuth provider: %s", name)
	s.providers[name] = provider
}

// GetProvider returns a provider by name
func (s *OAuthService) GetProvider(name string) (OAuthProvider, error) {
	provider, exists := s.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// GetAuthURL returns the auth URL for the specified provider
func (s *OAuthService) GetAuthURL(provider string, state string) (string, error) {
	p, err := s.GetProvider(provider)
	if err != nil {
		return "", err
	}

	log.Printf("Getting auth URL for provider %s with state: %s", provider, state)
	url := p.GetAuthURL(state)
	log.Printf("Generated URL: %s", url)
	return url, nil
}

// ExchangeToken exchanges the auth code for an access token
func (s *OAuthService) ExchangeToken(ctx context.Context, provider string, code string) (string, error) {
	p, err := s.GetProvider(provider)
	if err != nil {
		return "", err
	}
	return p.ExchangeToken(ctx, code)
}

// FetchUserInfo fetches the user information using the access token
func (s *OAuthService) FetchUserInfo(provider string, token string) (*models.OAuthUser, error) {
	p, err := s.GetProvider(provider)
	if err != nil {
		return nil, err
	}
	return p.FetchUserInfo(token)
}

// ListProviders returns a list of registered provider names
func (s *OAuthService) ListProviders() []string {
	providers := make([]string, 0, len(s.providers))
	for name := range s.providers {
		providers = append(providers, name)
	}
	return providers
}
