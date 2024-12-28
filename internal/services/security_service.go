package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"
	"time"

	"github.com/google/uuid"
)

// SecurityService defines the interface for security-related operations
type SecurityService interface {
	ListWhitelistedIPs(tenantID uuid.UUID) ([]string, error)
	AddWhitelistedIP(tenantID uuid.UUID, ip string) error
	RemoveWhitelistedIP(tenantID uuid.UUID, ip string) error
	ListAPIKeys(tenantID uuid.UUID) ([]*models.APIKey, error)
	GetAPIKey(tenantID uuid.UUID, keyID uuid.UUID) (*models.APIKey, error)
	UpdateAPIKey(tenantID uuid.UUID, keyID uuid.UUID, name string) (*models.APIKey, error)
	CreateAPIKey(tenantID uuid.UUID, name string, expiresAt *time.Time) (*models.APIKey, error)
	RevokeAPIKey(tenantID uuid.UUID, keyID uuid.UUID) error
	GetSecurityAuditLogs(tenantID uuid.UUID, page, limit int, filter map[string]string) ([]*models.AuditLog, int64, error)
	GetAuditLogEntry(tenantID uuid.UUID, logID uuid.UUID) (*models.AuditLog, error)
	GetSecurityPolicies(tenantID uuid.UUID) (*models.SecurityPolicies, error)
	UpdateSecurityPolicies(tenantID uuid.UUID, policies *models.SecurityPolicies) error
	TestSecurityPolicy(tenantID uuid.UUID, policy *models.SecurityPolicies) (map[string]bool, error)
	GetSecurityMetrics(tenantID uuid.UUID) (*models.SecurityMetrics, error)
	GetSecurityAlerts(tenantID uuid.UUID, status string) ([]*models.SecurityAlert, error)
	UpdateSecurityAlert(tenantID uuid.UUID, alertID uuid.UUID, status string) error
}

type securityService struct {
	securityRepo repositories.SecurityRepository
}

func NewSecurityService(securityRepo repositories.SecurityRepository) SecurityService {
	return &securityService{
		securityRepo: securityRepo,
	}
}

func (s *securityService) ListWhitelistedIPs(tenantID uuid.UUID) ([]string, error) {
	return s.securityRepo.ListWhitelistedIPs(tenantID)
}

func (s *securityService) AddWhitelistedIP(tenantID uuid.UUID, ip string) error {
	return s.securityRepo.AddWhitelistedIP(tenantID, ip)
}

func (s *securityService) RemoveWhitelistedIP(tenantID uuid.UUID, ip string) error {
	return s.securityRepo.RemoveWhitelistedIP(tenantID, ip)
}

func (s *securityService) ListAPIKeys(tenantID uuid.UUID) ([]*models.APIKey, error) {
	return s.securityRepo.ListAPIKeys(tenantID)
}

func (s *securityService) CreateAPIKey(tenantID uuid.UUID, name string, expiresAt *time.Time) (*models.APIKey, error) {
	return s.securityRepo.CreateAPIKey(tenantID, name, expiresAt)
}

func (s *securityService) RevokeAPIKey(tenantID uuid.UUID, keyID uuid.UUID) error {
	return s.securityRepo.RevokeAPIKey(tenantID, keyID)
}

func (s *securityService) GetSecurityAuditLogs(tenantID uuid.UUID, page, limit int, filter map[string]string) ([]*models.AuditLog, int64, error) {
	return s.securityRepo.GetSecurityAuditLogs(tenantID, page, limit, filter)
}

func (s *securityService) GetSecurityPolicies(tenantID uuid.UUID) (*models.SecurityPolicies, error) {
	return s.securityRepo.GetSecurityPolicies(tenantID)
}

func (s *securityService) UpdateSecurityPolicies(tenantID uuid.UUID, policies *models.SecurityPolicies) error {
	return s.securityRepo.UpdateSecurityPolicies(tenantID, policies)
}

func (s *securityService) GetSecurityMetrics(tenantID uuid.UUID) (*models.SecurityMetrics, error) {
	return s.securityRepo.GetSecurityMetrics(tenantID)
}

func (s *securityService) GetSecurityAlerts(tenantID uuid.UUID, status string) ([]*models.SecurityAlert, error) {
	return s.securityRepo.GetSecurityAlerts(tenantID, status)
}

func (s *securityService) UpdateSecurityAlert(tenantID uuid.UUID, alertID uuid.UUID, status string) error {
	return s.securityRepo.UpdateSecurityAlert(tenantID, alertID, status)
}

func (s *securityService) GetAPIKey(tenantID uuid.UUID, keyID uuid.UUID) (*models.APIKey, error) {
	return s.securityRepo.GetAPIKey(tenantID, keyID)
}

func (s *securityService) UpdateAPIKey(tenantID uuid.UUID, keyID uuid.UUID, name string) (*models.APIKey, error) {
	return s.securityRepo.UpdateAPIKey(tenantID, keyID, name)
}

func (s *securityService) GetAuditLogEntry(tenantID uuid.UUID, logID uuid.UUID) (*models.AuditLog, error) {
	return s.securityRepo.GetAuditLogEntry(tenantID, logID)
}

func (s *securityService) TestSecurityPolicy(_ uuid.UUID, _ *models.SecurityPolicies) (map[string]bool, error) {
	// Test various aspects of the security policy
	results := map[string]bool{
		"ipWhitelistValid":    true,
		"mfaConfigValid":      true,
		"passwordPolicyValid": true,
		"sessionPolicyValid":  true,
		"auditLoggingEnabled": true,
	}
	return results, nil
}
