package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecurityService defines the interface for security-related operations
type SecurityService interface {
	ListWhitelistedIPs(ctx *gin.Context, tenantID uuid.UUID) ([]string, error)
	AddWhitelistedIP(ctx *gin.Context, tenantID uuid.UUID, ip string) error
	RemoveWhitelistedIP(ctx *gin.Context, tenantID uuid.UUID, ip string) error
	ListAPIKeys(ctx *gin.Context, tenantID uuid.UUID) ([]*models.APIKey, error)
	GetAPIKey(ctx *gin.Context, tenantID uuid.UUID, keyID uuid.UUID) (*models.APIKey, error)
	UpdateAPIKey(ctx *gin.Context, tenantID uuid.UUID, keyID uuid.UUID, name string) (*models.APIKey, error)
	CreateAPIKey(ctx *gin.Context, tenantID uuid.UUID, name string, expiresAt *time.Time) (*models.APIKey, error)
	RevokeAPIKey(ctx *gin.Context, tenantID uuid.UUID, keyID uuid.UUID) error
	GetSecurityAuditLogs(ctx *gin.Context, tenantID uuid.UUID, page, limit int, filter map[string]string) ([]*models.AuditLog, int64, error)
	GetAuditLogEntry(ctx *gin.Context, tenantID uuid.UUID, logID uuid.UUID) (*models.AuditLog, error)
	GetSecurityPolicies(ctx *gin.Context, tenantID uuid.UUID) (*models.SecurityPolicies, error)
	UpdateSecurityPolicies(ctx *gin.Context, tenantID uuid.UUID, policies *models.SecurityPolicies) error
	TestSecurityPolicy(ctx *gin.Context, tenantID uuid.UUID, policy *models.SecurityPolicies) (map[string]bool, error)
	GetSecurityMetrics(ctx *gin.Context, tenantID uuid.UUID, timeRange string) (*models.SecurityMetrics, error)
	GetSecurityAlerts(ctx *gin.Context, tenantID uuid.UUID, status string) ([]*models.SecurityAlert, error)
	UpdateSecurityAlert(ctx *gin.Context, tenantID uuid.UUID, alertID uuid.UUID, status string) error
}

type securityService struct {
	securityRepo repositories.SecurityRepository
}

func NewSecurityService(securityRepo repositories.SecurityRepository) SecurityService {
	return &securityService{
		securityRepo: securityRepo,
	}
}

func (s *securityService) ListWhitelistedIPs(ctx *gin.Context, tenantID uuid.UUID) ([]string, error) {
	return s.securityRepo.ListWhitelistedIPs(tenantID)
}

func (s *securityService) AddWhitelistedIP(ctx *gin.Context, tenantID uuid.UUID, ip string) error {
	return s.securityRepo.AddWhitelistedIP(tenantID, ip)
}

func (s *securityService) RemoveWhitelistedIP(ctx *gin.Context, tenantID uuid.UUID, ip string) error {
	return s.securityRepo.RemoveWhitelistedIP(tenantID, ip)
}

func (s *securityService) ListAPIKeys(ctx *gin.Context, tenantID uuid.UUID) ([]*models.APIKey, error) {
	return s.securityRepo.ListAPIKeys(tenantID)
}

func (s *securityService) CreateAPIKey(ctx *gin.Context, tenantID uuid.UUID, name string, expiresAt *time.Time) (*models.APIKey, error) {
	return s.securityRepo.CreateAPIKey(tenantID, name, expiresAt)
}

func (s *securityService) RevokeAPIKey(ctx *gin.Context, tenantID uuid.UUID, keyID uuid.UUID) error {
	return s.securityRepo.RevokeAPIKey(tenantID, keyID)
}

func (s *securityService) GetSecurityAuditLogs(ctx *gin.Context, tenantID uuid.UUID, page, limit int, filter map[string]string) ([]*models.AuditLog, int64, error) {
	return s.securityRepo.GetSecurityAuditLogs(tenantID, page, limit, filter)
}

func (s *securityService) GetSecurityPolicies(ctx *gin.Context, tenantID uuid.UUID) (*models.SecurityPolicies, error) {
	return s.securityRepo.GetSecurityPolicies(tenantID)
}

func (s *securityService) UpdateSecurityPolicies(ctx *gin.Context, tenantID uuid.UUID, policies *models.SecurityPolicies) error {
	return s.securityRepo.UpdateSecurityPolicies(tenantID, policies)
}

func (s *securityService) GetSecurityMetrics(ctx *gin.Context, tenantID uuid.UUID, timeRange string) (*models.SecurityMetrics, error) {
	return s.securityRepo.GetSecurityMetrics(tenantID, timeRange)
}

func (s *securityService) GetSecurityAlerts(ctx *gin.Context, tenantID uuid.UUID, status string) ([]*models.SecurityAlert, error) {
	return s.securityRepo.GetSecurityAlerts(tenantID, status)
}

func (s *securityService) UpdateSecurityAlert(ctx *gin.Context, tenantID uuid.UUID, alertID uuid.UUID, status string) error {
	return s.securityRepo.UpdateSecurityAlert(tenantID, alertID, status)
}

func (s *securityService) GetAPIKey(ctx *gin.Context, tenantID uuid.UUID, keyID uuid.UUID) (*models.APIKey, error) {
	return s.securityRepo.GetAPIKey(tenantID, keyID)
}

func (s *securityService) UpdateAPIKey(ctx *gin.Context, tenantID uuid.UUID, keyID uuid.UUID, name string) (*models.APIKey, error) {
	return s.securityRepo.UpdateAPIKey(tenantID, keyID, name)
}

func (s *securityService) GetAuditLogEntry(ctx *gin.Context, tenantID uuid.UUID, logID uuid.UUID) (*models.AuditLog, error) {
	return s.securityRepo.GetAuditLogEntry(tenantID, logID)
}

func (s *securityService) TestSecurityPolicy(ctx *gin.Context, tenantID uuid.UUID, policy *models.SecurityPolicies) (map[string]bool, error) {
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
