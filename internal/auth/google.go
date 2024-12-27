package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"identity-service/config"
	"identity-service/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	Config *oauth2.Config
}

func init() {
	// Set up logging to write to both file and stdout
	logFile, err := os.OpenFile("oauth.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func NewGoogleProvider() *GoogleProvider {
	log.Println("=== Initializing Google Provider ===")
	oauthConfig := config.GoogleOAuth
	redirectURL := ""
	if len(oauthConfig.RedirectURIs) > 0 {
		redirectURL = oauthConfig.RedirectURIs[0]
	}
	log.Printf("Redirect URL: %s", redirectURL)
	log.Printf("Client ID present: %v", oauthConfig.ClientID != "")
	log.Printf("Client Secret present: %v", oauthConfig.ClientSecret != "")

	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		log.Fatal("ClientID or ClientSecret is missing in OAuth config")
	}

	provider := &GoogleProvider{
		Config: &oauth2.Config{
			ClientID:     oauthConfig.ClientID,
			ClientSecret: oauthConfig.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
	log.Println("Google Provider initialized successfully")
	return provider
}

func (g *GoogleProvider) GetAuthURL(state string) string {
	log.Printf("GetAuthURL called with state: %s", state)
	log.Printf("Config: ClientID=%v, RedirectURL=%v, Scopes=%v",
		g.Config.ClientID != "",
		g.Config.RedirectURL,
		g.Config.Scopes,
	)

	url := g.Config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
	)
	log.Printf("Generated OAuth URL: %s", url)
	return url
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
