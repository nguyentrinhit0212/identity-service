package config

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// 1. Load environment variables from .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 2. Get the base64 encoded string from the environment
	encodedString := os.Getenv("GOOGLE_OAUTH")

	// 3. Decode the base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		log.Fatalf("Error decoding base64 string: %v", err)
	}

	// 4. Unmarshal the JSON data into a struct with the "web" property
	var data struct {
		Web OAuthConfig `json:"web"`
	}
	if err := json.Unmarshal(decodedBytes, &data); err != nil {
		log.Fatalf("Failed to unmarshal OAuth config: %v", err)
	}

	// 5. Assign the nested "web" data to GoogleOAuth
	GoogleOAuth = data.Web

	log.Println("Google OAuth2 configuration loaded successfully!")
}
