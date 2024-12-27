package services

import (
	"encoding/json"
	"errors"

	"identity-service/internal/models"
	"identity-service/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
	ErrInvalidTemplate  = errors.New("invalid template format")
)

type TenantTemplate struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        models.TenantType      `json:"type"`
	Features    map[string]bool        `json:"features"`
	Settings    map[string]interface{} `json:"settings"`
	MaxUsers    *int                   `json:"maxUsers"`
	IsDefault   bool                   `json:"isDefault"`
}

type TenantTemplateService interface {
	GetTemplateByType(tenantType models.TenantType) (*TenantTemplate, error)
	GetTemplateByID(templateID uuid.UUID) (*TenantTemplate, error)
	ListTemplates() ([]*TenantTemplate, error)
	CreateTemplate(template *TenantTemplate) error
	UpdateTemplate(template *TenantTemplate) error
	DeleteTemplate(templateID uuid.UUID) error
	ApplyTemplate(tenant *models.Tenant, template *TenantTemplate) error
}

type tenantTemplateService struct {
	tenantService TenantService
}

func NewTenantTemplateService(tenantService TenantService) TenantTemplateService {
	return &tenantTemplateService{
		tenantService: tenantService,
	}
}

func (s *tenantTemplateService) GetTemplateByType(tenantType models.TenantType) (*TenantTemplate, error) {
	// Get default template for the tenant type
	templates, err := s.ListTemplates()
	if err != nil {
		return nil, err
	}

	for _, tmpl := range templates {
		if tmpl.Type == tenantType && tmpl.IsDefault {
			return tmpl, nil
		}
	}

	// If no default template found, create one
	return s.createDefaultTemplate(tenantType)
}

func (s *tenantTemplateService) GetTemplateByID(templateID uuid.UUID) (*TenantTemplate, error) {
	templates, err := s.ListTemplates()
	if err != nil {
		return nil, err
	}

	for _, tmpl := range templates {
		if tmpl.ID == templateID {
			return tmpl, nil
		}
	}

	return nil, ErrTemplateNotFound
}

func (s *tenantTemplateService) ListTemplates() ([]*TenantTemplate, error) {
	// TODO: Implement template storage and retrieval
	// For now, return default templates
	return []*TenantTemplate{
		s.getPersonalTemplate(),
		s.getTeamTemplate(),
		s.getEnterpriseTemplate(),
	}, nil
}

func (s *tenantTemplateService) CreateTemplate(template *TenantTemplate) error {
	// TODO: Implement template storage
	return nil
}

func (s *tenantTemplateService) UpdateTemplate(template *TenantTemplate) error {
	// TODO: Implement template storage
	return nil
}

func (s *tenantTemplateService) DeleteTemplate(templateID uuid.UUID) error {
	// TODO: Implement template storage
	return nil
}

func (s *tenantTemplateService) ApplyTemplate(tenant *models.Tenant, template *TenantTemplate) error {
	if template == nil {
		return ErrTemplateNotFound
	}

	// Apply template features
	featuresJSON, err := json.Marshal(template.Features)
	if err != nil {
		return err
	}
	tenant.Features = featuresJSON

	// Apply template settings
	settingsJSON, err := json.Marshal(template.Settings)
	if err != nil {
		return err
	}
	tenant.Settings = settingsJSON

	// Apply max users
	tenant.MaxUsers = template.MaxUsers

	// Update tenant
	_, err = s.tenantService.UpdateTenant(nil, tenant.ID, &models.TenantUpdate{})
	return err
}

func (s *tenantTemplateService) createDefaultTemplate(tenantType models.TenantType) (*TenantTemplate, error) {
	switch tenantType {
	case models.PersonalTenant:
		return s.getPersonalTemplate(), nil
	case models.TeamTenant:
		return s.getTeamTemplate(), nil
	case models.EnterpriseTenant:
		return s.getEnterpriseTemplate(), nil
	default:
		return nil, ErrInvalidTemplate
	}
}

func (s *tenantTemplateService) getPersonalTemplate() *TenantTemplate {
	return &TenantTemplate{
		ID:          uuid.New(),
		Name:        "Personal Workspace",
		Description: "Default template for personal workspaces",
		Type:        models.PersonalTenant,
		Features: map[string]bool{
			"personalFeatures":  true,
			"personalDashboard": true,
			"personalStorage":   true,
			"basicIntegrations": true,
		},
		Settings: map[string]interface{}{
			"storageLimit": "5GB",
			"apiLimit":     1000,
		},
		MaxUsers:  utils.Ptr(1),
		IsDefault: true,
	}
}

func (s *tenantTemplateService) getTeamTemplate() *TenantTemplate {
	return &TenantTemplate{
		ID:          uuid.New(),
		Name:        "Team Workspace",
		Description: "Default template for team workspaces",
		Type:        models.TeamTenant,
		Features: map[string]bool{
			"teamFeatures":         true,
			"collaborationTools":   true,
			"teamDashboard":        true,
			"sharedStorage":        true,
			"advancedIntegrations": true,
			"basicAnalytics":       true,
			"inviteSystem":         true,
		},
		Settings: map[string]interface{}{
			"storageLimit": "50GB",
			"apiLimit":     10000,
		},
		MaxUsers:  utils.Ptr(10),
		IsDefault: true,
	}
}

func (s *tenantTemplateService) getEnterpriseTemplate() *TenantTemplate {
	return &TenantTemplate{
		ID:          uuid.New(),
		Name:        "Enterprise Workspace",
		Description: "Default template for enterprise workspaces",
		Type:        models.EnterpriseTenant,
		Features: map[string]bool{
			"enterpriseFeatures": true,
			"sso":                true,
			"audit":              true,
			"advancedSecurity":   true,
			"customBranding":     true,
			"apiAccess":          true,
			"prioritySupport":    true,
		},
		Settings: map[string]interface{}{
			"storageLimit": "Unlimited",
			"apiLimit":     100000,
		},
		MaxUsers:  nil, // Unlimited
		IsDefault: true,
	}
}
