package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

const (
	// SecretKeyLength is the length of the random secret key in bytes (32 bytes = 256 bits)
	SecretKeyLength = 32

	// DefaultRotationInterval is the default interval for key rotation
	DefaultRotationInterval = 24 * time.Hour
)

var (
	ErrInvalidKeyLength = errors.New("invalid key length")
	ErrKeyGeneration    = errors.New("failed to generate key")
)

// SecretKey represents a JWT secret key with metadata
type SecretKey struct {
	Key       []byte
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SecretManager manages JWT secret keys with rotation
type SecretManager struct {
	currentKey       *SecretKey
	previousKey      *SecretKey
	rotationInterval time.Duration
	mu               sync.RWMutex
}

// NewSecretManager creates a new SecretManager with the given rotation interval
func NewSecretManager(rotationInterval time.Duration) (*SecretManager, error) {
	if rotationInterval == 0 {
		rotationInterval = DefaultRotationInterval
	}

	manager := &SecretManager{
		rotationInterval: rotationInterval,
	}

	// Generate initial key
	if err := manager.rotateKey(); err != nil {
		return nil, err
	}

	// Start key rotation goroutine
	go manager.startKeyRotation()

	return manager, nil
}

// GenerateKey generates a cryptographically secure random key
func GenerateKey() ([]byte, error) {
	key := make([]byte, SecretKeyLength)
	if _, err := rand.Read(key); err != nil {
		return nil, ErrKeyGeneration
	}
	return key, nil
}

// GetCurrentKey returns the current secret key
func (m *SecretManager) GetCurrentKey() []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentKey.Key
}

// GetKey returns either the current or previous key based on the creation time
func (m *SecretManager) GetKey(keyCreationTime time.Time) []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.previousKey != nil && !keyCreationTime.Before(m.previousKey.CreatedAt) && keyCreationTime.Before(m.currentKey.CreatedAt) {
		return m.previousKey.Key
	}
	return m.currentKey.Key
}

// rotateKey generates a new key and rotates the existing ones
func (m *SecretManager) rotateKey() error {
	newKey, err := GenerateKey()
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.previousKey = m.currentKey
	m.currentKey = &SecretKey{
		Key:       newKey,
		CreatedAt: now,
		ExpiresAt: now.Add(m.rotationInterval * 2), // Keep keys valid for 2x rotation interval
	}

	return nil
}

// startKeyRotation starts the key rotation goroutine
func (m *SecretManager) startKeyRotation() {
	ticker := time.NewTicker(m.rotationInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.rotateKey(); err != nil {
			// Log the error but continue running
			// TODO: Add proper error logging
			continue
		}
	}
}

// EncodeKey encodes a binary key to a base64 string
func EncodeKey(key []byte) string {
	return base64.RawURLEncoding.EncodeToString(key)
}

// DecodeKey decodes a base64 string to a binary key
func DecodeKey(encodedKey string) ([]byte, error) {
	key, err := base64.RawURLEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}
	if len(key) != SecretKeyLength {
		return nil, ErrInvalidKeyLength
	}
	return key, nil
}
