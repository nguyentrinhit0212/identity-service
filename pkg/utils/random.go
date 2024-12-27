package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomState generates a random state string for OAuth
func GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
