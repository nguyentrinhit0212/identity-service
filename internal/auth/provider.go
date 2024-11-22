package auth

import (
	"context"
	"identity-service/internal/models"
)

type OAuthProviderInterface interface {
	GetAuthURL(state string) string
	ExchangeToken(ctx context.Context, code string) (string, error)
	FetchUserInfo(token string) (*models.OAuthUser, error)
	GetProviderName() string
}
