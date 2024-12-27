package services

import (
	"encoding/json"
	"identity-service/internal/models"

	"github.com/google/uuid"
)

type FeatureService interface {
	IsFeatureEnabled(tenantID uuid.UUID, feature string) bool
	GetFeatureList(tenantType models.TenantType) []string
	RequiresUpgrade(feature string, currentType models.TenantType) (bool, models.TenantType)
}

type featureService struct {
	tenantService TenantService
}

func NewFeatureService(tenantService TenantService) FeatureService {
	return &featureService{
		tenantService: tenantService,
	}
}

func (s *featureService) IsFeatureEnabled(tenantID uuid.UUID, feature string) bool {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return false
	}

	var features map[string]bool
	if err := json.Unmarshal(tenant.Features, &features); err != nil {
		return false
	}

	return features[feature]
}

func (s *featureService) GetFeatureList(tenantType models.TenantType) []string {
	switch tenantType {
	case models.PersonalTenant:
		return []string{
			"personalFeatures",
			"personalDashboard",
			"personalStorage",
			"basicIntegrations",
		}
	case models.TeamTenant:
		return []string{
			"teamFeatures",
			"collaborationTools",
			"teamDashboard",
			"sharedStorage",
			"advancedIntegrations",
			"basicAnalytics",
			"inviteSystem",
		}
	case models.EnterpriseTenant:
		return []string{
			"enterpriseFeatures",
			"sso",
			"audit",
			"advancedSecurity",
			"enterpriseDashboard",
			"unlimitedStorage",
			"customIntegrations",
			"advancedAnalytics",
			"bulkInviteSystem",
			"customBranding",
			"apiAccess",
			"prioritySupport",
		}
	default:
		return []string{}
	}
}

func (s *featureService) RequiresUpgrade(feature string, currentType models.TenantType) (bool, models.TenantType) {
	featureToTypeMap := map[string]models.TenantType{
		// Team features
		"collaborationTools":   models.TeamTenant,
		"teamDashboard":        models.TeamTenant,
		"sharedStorage":        models.TeamTenant,
		"advancedIntegrations": models.TeamTenant,
		"basicAnalytics":       models.TeamTenant,
		"inviteSystem":         models.TeamTenant,

		// Enterprise features
		"sso":              models.EnterpriseTenant,
		"audit":            models.EnterpriseTenant,
		"advancedSecurity": models.EnterpriseTenant,
		"customBranding":   models.EnterpriseTenant,
		"apiAccess":        models.EnterpriseTenant,
		"prioritySupport":  models.EnterpriseTenant,
	}

	requiredType, exists := featureToTypeMap[feature]
	if !exists {
		return false, currentType
	}

	// Check if current type is less than required type
	switch currentType {
	case models.PersonalTenant:
		return true, requiredType
	case models.TeamTenant:
		return requiredType == models.EnterpriseTenant, models.EnterpriseTenant
	case models.EnterpriseTenant:
		return false, models.EnterpriseTenant
	default:
		return true, requiredType
	}
}
