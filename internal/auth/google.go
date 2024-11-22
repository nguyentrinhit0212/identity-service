package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"identity-service/config"
	"identity-service/internal/models"

	"golang.org/x/oauth2"
)

type GoogleProvider struct {
	Config *oauth2.Config
}

func NewGoogleProvider() *GoogleProvider {
	oauthConfig := config.GoogleOAuth
	redirectURL := ""
	if len(oauthConfig.RedirectURIs) > 0 {
		redirectURL = oauthConfig.RedirectURIs[0]
	}
	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		log.Fatalf("ClientID or ClientSecret is missing in OAuth config")
	}
	if oauthConfig.AuthURI == "" || oauthConfig.TokenURI == "" {
		log.Fatalf("AuthURI or TokenURI is missing in OAuth config")
	}

	return &GoogleProvider{
		Config: &oauth2.Config{
			ClientID:     oauthConfig.ClientID,
			ClientSecret: oauthConfig.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  oauthConfig.AuthURI,
				TokenURL: oauthConfig.TokenURI,
			},
		},
	}
}

func (g *GoogleProvider) GetAuthURL(state string) string {
	return g.Config.AuthCodeURL(state)
}

func (g *GoogleProvider) ExchangeToken(ctx context.Context, code string) (string, error) {
	token, err := g.Config.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func (g *GoogleProvider) FetchUserInfo(token string) (*models.OAuthUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &models.OAuthUser{
		ID:            googleUser.ID,
		Email:         googleUser.Email,
		VerifiedEmail: googleUser.VerifiedEmail,
		Name:          googleUser.Name,
		Picture:       googleUser.Picture,
		Provider:      "google",
	}, nil
}

func (g *GoogleProvider) GetProviderName() string {
	return "google"
}
