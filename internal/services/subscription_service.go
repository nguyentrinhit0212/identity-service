package services

import (
	"encoding/json"
	"errors"
	"time"

	"identity-service/internal/models"
	"identity-service/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrInvalidPlan          = errors.New("invalid subscription plan")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrSubscriptionExpired  = errors.New("subscription expired")
	ErrDowngradeNotAllowed  = errors.New("downgrade not allowed with current usage")
)

type SubscriptionService interface {
	GetSubscriptionStatus(tenantID uuid.UUID) (string, error)
	UpdateSubscription(tenantID uuid.UUID, plan string) error
	HandleSubscriptionWebhook(event interface{}) error
	CheckSubscriptionValid(tenantID uuid.UUID) error
}

type subscriptionService struct {
	tenantService TenantService
	usageService  UsageService
}

func NewSubscriptionService(tenantService TenantService, usageService UsageService) SubscriptionService {
	return &subscriptionService{
		tenantService: tenantService,
		usageService:  usageService,
	}
}

func (s *subscriptionService) GetSubscriptionStatus(tenantID uuid.UUID) (string, error) {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return "", err
	}

	if tenant.Type == models.PersonalTenant {
		return "free", nil
	}

	return tenant.SubscriptionStatus, nil
}

func (s *subscriptionService) UpdateSubscription(tenantID uuid.UUID, plan string) error {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return err
	}

	// Validate plan change
	if err := s.validatePlanChange(tenant, plan); err != nil {
		return err
	}

	// Update subscription details
	tenant.SubscriptionStatus = plan
	tenant.SubscriptionExpiresAt = utils.Ptr(time.Now().AddDate(1, 0, 0)) // 1 year from now

	// Update features based on plan
	features := s.getFeaturesByPlan(plan)
	featuresJSON, err := json.Marshal(features)
	if err != nil {
		return err
	}
	rawJSON := json.RawMessage(featuresJSON)
	update := &models.TenantUpdate{
		Features: &rawJSON,
	}
	_, err = s.tenantService.UpdateTenant(nil, tenant.ID, update)
	return err
}

func (s *subscriptionService) HandleSubscriptionWebhook(event interface{}) error {
	// TODO: Implement webhook handling for your payment provider
	// This would handle events like:
	// - subscription.created
	// - subscription.updated
	// - subscription.deleted
	// - payment.succeeded
	// - payment.failed
	return nil
}

func (s *subscriptionService) CheckSubscriptionValid(tenantID uuid.UUID) error {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return err
	}

	if tenant.Type == models.PersonalTenant {
		return nil // Personal tenants are always valid
	}

	if tenant.SubscriptionExpiresAt.Before(time.Now()) {
		return ErrSubscriptionExpired
	}

	return nil
}

func (s *subscriptionService) validatePlanChange(tenant *models.Tenant, newPlan string) error {
	// Check if downgrading
	if tenant.SubscriptionStatus == "enterprise" && newPlan == "team" {
		// Check if current usage allows downgrade
		usageLimits := s.usageService.CheckUsageLimits(tenant.ID)
		if !usageLimits["canDowngrade"] {
			return ErrDowngradeNotAllowed
		}
	}

	// Validate plan name
	switch newPlan {
	case "team", "enterprise":
		return nil
	default:
		return ErrInvalidPlan
	}
}

func (s *subscriptionService) getFeaturesByPlan(plan string) map[string]bool {
	switch plan {
	case "team":
		return map[string]bool{
			"teamFeatures":         true,
			"collaborationTools":   true,
			"teamDashboard":        true,
			"sharedStorage":        true,
			"advancedIntegrations": true,
			"basicAnalytics":       true,
			"inviteSystem":         true,
		}
	case "enterprise":
		return map[string]bool{
			"enterpriseFeatures": true,
			"sso":                true,
			"audit":              true,
			"advancedSecurity":   true,
			"customBranding":     true,
			"apiAccess":          true,
			"prioritySupport":    true,
		}
	default:
		return map[string]bool{}
	}
}
