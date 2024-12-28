package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// GenerateRandomState generates a random state string for OAuth
func GenerateRandomState() string {
	return GenerateRandomString(32)
}
