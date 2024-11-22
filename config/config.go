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

// LoadOAuthConfig loads OAuth2 configuration from a JSON file
func LoadOAuthConfig() {
	// Mở file cấu hình OAuth
	file, err := os.Open("config/google_oauth.json")
	if err != nil {
		log.Fatalf("Failed to open OAuth config file: %v", err)
	}
	defer file.Close()

	// Kiểm tra xem file đã mở thành công
	log.Println("OAuth config file opened successfully")

	var data struct {
		Web OAuthConfig `json:"web"`
	}

	// Giải mã cấu hình JSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		log.Fatalf("Failed to decode OAuth config: %v", err)
	}

	// Gán cấu hình cho GoogleOAuth
	GoogleOAuth = data.Web

	// Kiểm tra cấu hình đã được load chính xác
	log.Printf("Google OAuth Config: %+v", GoogleOAuth)

	log.Println("Google OAuth2 configuration loaded successfully!")
}
