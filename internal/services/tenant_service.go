package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"

	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// TenantService defines the interface for tenant-related operations
type TenantService interface {
	ListTenants(page, limit int, search string, filter map[string]string) ([]*models.Tenant, int64, error)
	CreateTenant(tenant *models.Tenant) (*models.Tenant, error)
	GetTenant(id uuid.UUID) (*models.Tenant, error)
	GetTenantByID(id uuid.UUID) (*models.Tenant, error)
	UpdateTenant(id uuid.UUID, updates *models.TenantUpdate) (*models.Tenant, error)
	DeleteTenant(id uuid.UUID) error
	GetTenantSettings(id uuid.UUID) (*models.TenantSettings, error)
	UpdateTenantSettings(id uuid.UUID, settings *models.TenantSettings) (*models.TenantSettings, error)
	GetTenantMembers(id uuid.UUID, page, limit int, search string, filter map[string]string) ([]*models.User, int64, error)
	GetTenantFeatures(id uuid.UUID) (*models.TenantFeatures, error)
	UpdateTenantFeatures(id uuid.UUID, features *models.TenantFeatures) (*models.TenantFeatures, error)
	SwitchTenant(userID, tenantID uuid.UUID) error
	UpgradeTenant(id uuid.UUID, upgrade *models.TenantUpgrade) error
	CreateTenantInvite(id uuid.UUID, invite *models.TenantInvite) (*models.TenantInvite, error)
	DeleteTenantInvite(tenantID, inviteID uuid.UUID) error
}

type tenantService struct {
	tenantRepo repositories.TenantRepository
}

func NewTenantService(tenantRepo repositories.TenantRepository) TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
	}
}

func (s *tenantService) ListTenants(page, limit int, search string, filter map[string]string) ([]*models.Tenant, int64, error) {
	return s.tenantRepo.ListTenants(page, limit, search, filter)
}

func (s *tenantService) CreateTenant(tenant *models.Tenant) (*models.Tenant, error) {
	if err := s.tenantRepo.CreateTenant(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *tenantService) GetTenant(id uuid.UUID) (*models.Tenant, error) {
	return s.tenantRepo.GetTenantByID(id)
}

func (s *tenantService) UpdateTenant(id uuid.UUID, updates *models.TenantUpdate) (*models.Tenant, error) {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.Name != nil {
		tenant.Name = *updates.Name
	}
	if updates.Type != nil {
		tenant.Type = *updates.Type
	}
	if updates.Domain != nil {
		tenant.Domain = *updates.Domain
	}
	if updates.DomainVerified != nil {
		tenant.DomainVerified = *updates.DomainVerified
	}
	if updates.OwnerID != nil {
		tenant.OwnerID = updates.OwnerID
	}
	if updates.MaxUsers != nil {
		tenant.MaxUsers = updates.MaxUsers
	}
	if updates.AuthProviders != nil {
		tenant.AuthProviders = *updates.AuthProviders
	}
	if updates.Features != nil {
		tenant.Features = *updates.Features
	}
	if updates.Settings != nil {
		tenant.Settings = *updates.Settings
	}
	if updates.SubscriptionStatus != nil {
		tenant.SubscriptionStatus = *updates.SubscriptionStatus
	}
	if updates.SubscriptionPlan != nil {
		tenant.SubscriptionPlan = *updates.SubscriptionPlan
	}
	if updates.SubscriptionExpiresAt != nil {
		tenant.SubscriptionExpiresAt = updates.SubscriptionExpiresAt
	}

	if err := s.tenantRepo.UpdateTenant(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *tenantService) DeleteTenant(id uuid.UUID) error {
	return s.tenantRepo.DeleteTenant(id)
}

func (s *tenantService) GetTenantSettings(id uuid.UUID) (*models.TenantSettings, error) {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, err
	}
	var settings models.TenantSettings
	if err := json.Unmarshal(tenant.Settings, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *tenantService) UpdateTenantSettings(id uuid.UUID, settings *models.TenantSettings) (*models.TenantSettings, error) {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, err
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	tenant.Settings = settingsJSON
	if err := s.tenantRepo.UpdateTenant(tenant); err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *tenantService) GetTenantMembers(id uuid.UUID, page, limit int, search string, filter map[string]string) ([]*models.User, int64, error) {
	accesses, total, err := s.tenantRepo.GetTenantMembers(id, page, limit, search, filter)
	if err != nil {
		return nil, 0, err
	}

	users := make([]*models.User, len(accesses))
	for i, access := range accesses {
		users[i] = &access.User
	}
	return users, total, nil
}

func (s *tenantService) GetTenantFeatures(id uuid.UUID) (*models.TenantFeatures, error) {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, err
	}
	var features models.TenantFeatures
	if err := json.Unmarshal(tenant.Features, &features); err != nil {
		return nil, err
	}
	return &features, nil
}

func (s *tenantService) UpdateTenantFeatures(id uuid.UUID, features *models.TenantFeatures) (*models.TenantFeatures, error) {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return nil, err
	}

	featuresJSON, err := json.Marshal(features)
	if err != nil {
		return nil, err
	}

	tenant.Features = featuresJSON
	if err := s.tenantRepo.UpdateTenant(tenant); err != nil {
		return nil, err
	}

	return features, nil
}

func (s *tenantService) GetTenantByID(id uuid.UUID) (*models.Tenant, error) {
	return s.tenantRepo.GetTenantByID(id)
}

func (s *tenantService) SwitchTenant(userID, tenantID uuid.UUID) error {
	// Verify user has access to tenant
	access, err := s.tenantRepo.GetUserTenantAccess(userID, tenantID)
	if err != nil || access == nil {
		return fmt.Errorf("user does not have access to tenant")
	}
	return nil
}

func (s *tenantService) UpgradeTenant(id uuid.UUID, upgrade *models.TenantUpgrade) error {
	tenant, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return err
	}

	tenant.SubscriptionPlan = upgrade.Plan
	tenant.SubscriptionStatus = "active"
	tenant.SubscriptionExpiresAt = upgrade.ExpiresAt

	return s.tenantRepo.UpdateTenant(tenant)
}

func (s *tenantService) CreateTenantInvite(tenantID uuid.UUID, invite *models.TenantInvite) (*models.TenantInvite, error) {
	return s.tenantRepo.CreateTenantInvite(tenantID, invite)
}

func (s *tenantService) DeleteTenantInvite(tenantID, inviteID uuid.UUID) error {
	return s.tenantRepo.DeleteTenantInvite(tenantID, inviteID)
}
