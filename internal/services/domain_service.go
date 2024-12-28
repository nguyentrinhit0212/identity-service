package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"identity-service/internal/models"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DomainService handles domain verification and management
type DomainService interface {
	InitiateDomainVerification(tenantID uuid.UUID, domain string) (*DomainVerification, error)
	CheckVerificationStatus(tenantID uuid.UUID, domain string) (bool, error)
	VerifyDomain(tenantID uuid.UUID, domain string, method string) error
	GetVerifiedDomains(tenantID uuid.UUID) ([]string, error)
	RemoveDomain(tenantID uuid.UUID, domain string) error
}

// DomainVerification represents a domain verification request
type DomainVerification struct {
	Domain    string    `json:"domain"`
	Token     string    `json:"token"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type domainService struct {
	tenantService TenantService
}

var ErrVerificationFailed = errors.New("domain verification failed")

func NewDomainService(tenantService TenantService) DomainService {
	return &domainService{
		tenantService: tenantService,
	}
}

func (s *domainService) InitiateDomainVerification(tenantID uuid.UUID, domain string) (*DomainVerification, error) {
	// Validate domain format
	if !isValidDomain(domain) {
		return nil, errors.New("invalid domain format")
	}

	// Generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Create verification record
	verification := &DomainVerification{
		Domain:    domain,
		Token:     token,
		Method:    "dns",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Store verification details in tenant settings
	_, err = s.tenantService.UpdateTenant(tenantID, &models.TenantUpdate{
		Domain: &domain,
	})
	if err != nil {
		return nil, err
	}

	return verification, nil
}

func (s *domainService) CheckVerificationStatus(tenantID uuid.UUID, domain string) (bool, error) {
	tenant, err := s.tenantService.GetTenant(tenantID)
	if err != nil {
		return false, err
	}

	if tenant.Domain != domain {
		return false, errors.New("domain mismatch")
	}

	return tenant.DomainVerified, nil
}

func (s *domainService) VerifyDomain(tenantID uuid.UUID, domain string, method string) error {
	tenant, err := s.tenantService.GetTenant(tenantID)
	if err != nil {
		return err
	}

	if tenant.Domain != domain {
		return errors.New("domain mismatch")
	}

	var verified bool
	switch method {
	case "dns":
		verified = s.verifyDNSRecord(domain)
	case "file":
		verified = s.verifyFileToken(domain)
	default:
		return errors.New("invalid verification method")
	}

	if !verified {
		return ErrVerificationFailed
	}

	// Update tenant with verified domain
	_, err = s.tenantService.UpdateTenant(tenantID, &models.TenantUpdate{
		DomainVerified: &verified,
	})
	return err
}

func (s *domainService) GetVerifiedDomains(tenantID uuid.UUID) ([]string, error) {
	tenant, err := s.tenantService.GetTenant(tenantID)
	if err != nil {
		return nil, err
	}

	if tenant.DomainVerified && tenant.Domain != "" {
		return []string{tenant.Domain}, nil
	}

	return []string{}, nil
}

func (s *domainService) RemoveDomain(tenantID uuid.UUID, domain string) error {
	tenant, err := s.tenantService.GetTenant(tenantID)
	if err != nil {
		return err
	}

	if tenant.Domain != domain {
		return errors.New("domain mismatch")
	}

	emptyDomain := ""
	verified := false
	_, err = s.tenantService.UpdateTenant(tenantID, &models.TenantUpdate{
		Domain:         &emptyDomain,
		DomainVerified: &verified,
	})
	return err
}

// Helper functions

func isValidDomain(domain string) bool {
	// Basic domain validation
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}
	return true
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *domainService) verifyDNSRecord(domain string) bool {
	// Implement DNS record verification logic
	// This is a placeholder implementation
	records, err := net.LookupTXT(domain)
	if err != nil {
		return false
	}

	// Check if any of the TXT records match our verification token
	for _, record := range records {
		if strings.HasPrefix(record, "identity-verification=") {
			return true
		}
	}

	return false
}

func (s *domainService) verifyFileToken(_ string) bool {
	// Implement file-based verification logic
	// This is a placeholder implementation
	return false
}
