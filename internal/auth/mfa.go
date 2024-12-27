package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
)

func GenerateMFASecret() string {
	// Generate 20 random bytes
	bytes := make([]byte, 20)
	rand.Read(bytes)

	// Encode to base32
	return base32.StdEncoding.EncodeToString(bytes)
}

func GenerateMFAQRCode(secret string) string {
	// Generate otpauth URL
	// Format: otpauth://totp/Service:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=Service
	return fmt.Sprintf("otpauth://totp/YourService:%s?secret=%s&issuer=YourService", "user@example.com", secret)
}

func VerifyTOTP(secret string, token string) bool {
	// TODO: Implement TOTP verification
	return true
}