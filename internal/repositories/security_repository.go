package repositories

import (
	"identity-service/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SecurityRepository interface {
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
	GetSecurityMetrics(tenantID uuid.UUID) (*models.SecurityMetrics, error)
	GetSecurityAlerts(tenantID uuid.UUID, status string) ([]*models.SecurityAlert, error)
	UpdateSecurityAlert(tenantID uuid.UUID, alertID uuid.UUID, status string) error
}

type securityRepository struct {
	db GormDB
}

func NewSecurityRepository(db GormDB) SecurityRepository {
	return &securityRepository{
		db: db,
	}
}

func (r *securityRepository) ListWhitelistedIPs(tenantID uuid.UUID) ([]string, error) {
	var ips []string
	err := r.db.Model(&models.SecurityPolicies{}).
		Where("tenant_id = ?", tenantID).
		Pluck("whitelisted_ips", &ips).Error
	return ips, err
}

func (r *securityRepository) AddWhitelistedIP(tenantID uuid.UUID, ip string) error {
	return r.db.Model(&models.SecurityPolicies{}).
		Where("tenant_id = ?", tenantID).
		Update("whitelisted_ips", gorm.Expr("array_append(whitelisted_ips, ?)", ip)).Error
}

func (r *securityRepository) RemoveWhitelistedIP(tenantID uuid.UUID, ip string) error {
	return r.db.Model(&models.SecurityPolicies{}).
		Where("tenant_id = ?", tenantID).
		Update("whitelisted_ips", gorm.Expr("array_remove(whitelisted_ips, ?)", ip)).Error
}

func (r *securityRepository) ListAPIKeys(tenantID uuid.UUID) ([]*models.APIKey, error) {
	var keys []*models.APIKey
	err := r.db.Where("tenant_id = ? AND revoked_at IS NULL", tenantID).Find(&keys).Error
	return keys, err
}

func (r *securityRepository) CreateAPIKey(tenantID uuid.UUID, name string, expiresAt *time.Time) (*models.APIKey, error) {
	key := &models.APIKey{
		TenantID:  tenantID,
		Name:      name,
		Key:       generateAPIKey(),
		ExpiresAt: expiresAt,
	}
	err := r.db.Create(key).Error
	return key, err
}

func (r *securityRepository) RevokeAPIKey(tenantID uuid.UUID, keyID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.APIKey{}).
		Where("tenant_id = ? AND id = ?", tenantID, keyID).
		Update("revoked_at", &now).Error
}

func (r *securityRepository) GetSecurityAuditLogs(tenantID uuid.UUID, page, limit int, filter map[string]string) ([]*models.AuditLog, int64, error) {
	var logs []*models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{}).Where("tenant_id = ?", tenantID)

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&logs).Error
	return logs, total, err
}

func (r *securityRepository) GetSecurityPolicies(tenantID uuid.UUID) (*models.SecurityPolicies, error) {
	var policies models.SecurityPolicies
	err := r.db.First(&policies, "tenant_id = ?", tenantID).Error
	return &policies, err
}

func (r *securityRepository) UpdateSecurityPolicies(tenantID uuid.UUID, policies *models.SecurityPolicies) error {
	policies.TenantID = tenantID
	return r.db.Save(policies).Error
}

func (r *securityRepository) GetSecurityMetrics(tenantID uuid.UUID) (*models.SecurityMetrics, error) {
	var metrics models.SecurityMetrics
	err := r.db.First(&metrics, "tenant_id = ?", tenantID).Error
	return &metrics, err
}

func (r *securityRepository) GetSecurityAlerts(tenantID uuid.UUID, status string) ([]*models.SecurityAlert, error) {
	var alerts []*models.SecurityAlert
	query := r.db.Where("tenant_id = ?", tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&alerts).Error
	return alerts, err
}

func (r *securityRepository) UpdateSecurityAlert(tenantID uuid.UUID, alertID uuid.UUID, status string) error {
	return r.db.Model(&models.SecurityAlert{}).
		Where("tenant_id = ? AND id = ?", tenantID, alertID).
		Update("status", status).Error
}

func (r *securityRepository) GetAPIKey(tenantID uuid.UUID, keyID uuid.UUID) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, keyID).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *securityRepository) UpdateAPIKey(tenantID uuid.UUID, keyID uuid.UUID, name string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.Model(&apiKey).
		Where("tenant_id = ? AND id = ?", tenantID, keyID).
		Update("name", name).
		First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *securityRepository) GetAuditLogEntry(tenantID uuid.UUID, logID uuid.UUID) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, logID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// Helper function to generate API key
func generateAPIKey() string {
	// This is a placeholder implementation
	// In a real application, you would use a secure method to generate API keys
	return uuid.New().String()
}
