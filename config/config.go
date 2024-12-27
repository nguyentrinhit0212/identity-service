package config

import (
	"encoding/json"
	"log"
	"os"
)

type OAuthConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
	RedirectURIs []string `json:"redirect_uris"`
}

var GoogleOAuth OAuthConfig

func LoadOAuthConfig() {
	file, err := os.Open("google_oauth.json")
	if err != nil {
		log.Fatalf("Failed to open OAuth config file: %v", err)
	}
	defer file.Close()

	log.Println("OAuth config file opened successfully")

	var data struct {
		Web OAuthConfig `json:"web"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		log.Fatalf("Failed to decode OAuth config: %v", err)
	}

	GoogleOAuth = data.Web
	log.Printf("Google OAuth Config: %+v", GoogleOAuth)

	log.Println("Google OAuth2 configuration loaded successfully!")
}
