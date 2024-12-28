package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// DefaultKeySize is the default RSA key size in bits
	DefaultKeySize = 2048

	// KeyRotationInterval is the default interval for key rotation
	KeyRotationInterval = 24 * time.Hour
)

var (
	ErrInvalidKey  = errors.New("invalid key")
	ErrGeneration  = errors.New("failed to generate key")
	ErrKeyEncoding = errors.New("failed to encode key")
	ErrKeyDecoding = errors.New("failed to decode key")
)

// KeyPair represents an RSA key pair with metadata
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	KeyID      string // Unique identifier for the key pair
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

// KeyManager manages RSA key pairs with rotation
type KeyManager struct {
	currentKey       *KeyPair
	previousKey      *KeyPair
	rotationInterval time.Duration
	keySize          int
	mu               sync.RWMutex
}

// NewKeyManager creates a new KeyManager with the given parameters
func NewKeyManager(rotationInterval time.Duration, keySize int) (*KeyManager, error) {
	if rotationInterval == 0 {
		rotationInterval = KeyRotationInterval
	}
	if keySize == 0 {
		keySize = DefaultKeySize
	}

	manager := &KeyManager{
		rotationInterval: rotationInterval,
		keySize:          keySize,
	}

	// Generate initial key pair synchronously
	initialKey, err := manager.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate initial key pair: %v", err)
	}
	manager.currentKey = initialKey

	// Start key rotation goroutine
	go manager.startKeyRotation()

	return manager, nil
}

// GenerateKeyPair generates a new RSA key pair
func (m *KeyManager) GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, m.keySize)
	if err != nil {
		return nil, ErrGeneration
	}

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		KeyID:      generateKeyID(),
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(m.rotationInterval * 2),
	}, nil
}

// GetCurrentPrivateKey returns the current private key
func (m *KeyManager) GetCurrentPrivateKey() *rsa.PrivateKey {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentKey.PrivateKey
}

// GetCurrentPublicKey returns the current public key
func (m *KeyManager) GetCurrentPublicKey() *rsa.PublicKey {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentKey.PublicKey
}

// GetKeyPairByID returns the key pair matching the given ID
func (m *KeyManager) GetKeyPairByID(keyID string) *KeyPair {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentKey.KeyID == keyID {
		return m.currentKey
	}
	if m.previousKey != nil && m.previousKey.KeyID == keyID {
		return m.previousKey
	}
	return nil
}

// rotateKeys generates a new key pair and rotates the existing ones
func (m *KeyManager) rotateKeys() error {
	newKeyPair, err := m.GenerateKeyPair()
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.previousKey = m.currentKey
	m.currentKey = newKeyPair

	return nil
}

// startKeyRotation starts the key rotation goroutine
func (m *KeyManager) startKeyRotation() {
	ticker := time.NewTicker(m.rotationInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.rotateKeys(); err != nil {
			// TODO: Add proper error logging
			continue
		}
	}
}

// ExportPublicKeyPEM exports the current public key in PEM format
func (m *KeyManager) ExportPublicKeyPEM() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(m.currentKey.PublicKey)
	if err != nil {
		return nil, ErrKeyEncoding
	}

	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// ImportPublicKeyPEM imports a public key from PEM format
func ImportPublicKeyPEM(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrKeyDecoding
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrKeyDecoding
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKey
	}

	return rsaPub, nil
}

// generateKeyID generates a unique identifier for a key pair
func generateKeyID() string {
	var (
		b   = make([]byte, 16)
		n   int
		err error
	)
	n, err = rand.Read(b)
	if err != nil || n != 16 {
		// Fallback to timestamp-based ID if random generation fails
		return base64.RawURLEncoding.EncodeToString([]byte(time.Now().String()))
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// GetCurrentKeyID returns the current key's ID
func (m *KeyManager) GetCurrentKeyID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentKey.KeyID
}

// SignToken signs a token with the current key
func (m *KeyManager) SignToken(claims jwt.Claims) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentKey == nil {
		return "", ErrInvalidKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = m.currentKey.KeyID

	return token.SignedString(m.currentKey.PrivateKey)
}

// VerifyToken verifies a token and unmarshals its claims
func (m *KeyManager) VerifyToken(tokenString string, claims jwt.Claims) error {
	var err error
	parser := jwt.Parser{
		ValidMethods: []string{jwt.SigningMethodRS256.Name},
	}

	// First parse without validation to get the key ID
	var token *jwt.Token
	token, err = parser.Parse(tokenString, nil)
	if err != nil {
		return ErrInvalidKey
	}

	// Get key ID from token header
	var ok bool
	var keyID string
	if keyID, ok = token.Header["kid"].(string); !ok {
		return ErrInvalidKey
	}

	// Get the key pair for this key ID
	keyPair := m.GetKeyPairByID(keyID)
	if keyPair == nil {
		return ErrInvalidKey
	}

	// Parse and validate token with the correct public key
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		var isRSA bool
		if _, isRSA = token.Method.(*jwt.SigningMethodRSA); !isRSA {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return keyPair.PublicKey, nil
	})

	return err
}
