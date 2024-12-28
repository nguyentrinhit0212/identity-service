package services

import (
	"encoding/json"
	"time"

	"identity-service/internal/models"

	"github.com/google/uuid"
)

type UsageStats struct {
	APICallsCount     int64     `json:"apiCallsCount"`
	StorageUsageBytes int64     `json:"storageUsageBytes"`
	LastUpdated       time.Time `json:"lastUpdated"`
	Limits            Limits    `json:"limits"`
}

type Limits struct {
	MaxAPICalls     int64  `json:"maxApiCalls"`
	MaxStorageBytes int64  `json:"maxStorageBytes"`
	MaxUsers        int    `json:"maxUsers"`
	StorageUnit     string `json:"storageUnit"`
}

// UsageService defines the interface for usage tracking operations
type UsageService interface {
	TrackAPICall(tenantID uuid.UUID)
	GetStorageUsage(tenantID uuid.UUID) (int64, error)
	CheckUsageLimits(tenantID uuid.UUID) map[string]bool
	GetUsageStats(tenantID uuid.UUID) (*UsageStats, error)
}

type usageService struct {
	tenantService TenantService
}

func NewUsageService(tenantService TenantService) UsageService {
	return &usageService{
		tenantService: tenantService,
	}
}

func (s *usageService) TrackAPICall(tenantID uuid.UUID) {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return
	}

	var stats UsageStats
	if err := json.Unmarshal(tenant.Settings, &stats); err != nil {
		stats = UsageStats{LastUpdated: time.Now()}
	}

	stats.APICallsCount++
	stats.LastUpdated = time.Now()

	rawJSON, err := json.Marshal(stats)
	if err != nil {
		return
	}

	update := &models.TenantUpdate{
		Settings: (*json.RawMessage)(&rawJSON),
	}
	_, err = s.tenantService.UpdateTenant(tenantID, update)
	if err != nil {
		return
	}
}

func (s *usageService) GetStorageUsage(tenantID uuid.UUID) (int64, error) {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return 0, err
	}

	var stats UsageStats
	if err := json.Unmarshal(tenant.Settings, &stats); err != nil {
		return 0, err
	}

	return stats.StorageUsageBytes, nil
}

func (s *usageService) CheckUsageLimits(tenantID uuid.UUID) map[string]bool {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return map[string]bool{"error": true}
	}

	stats, err := s.GetUsageStats(tenantID)
	if err != nil {
		return map[string]bool{"error": true}
	}

	limits := s.getLimitsForTenant(tenant.Type)

	return map[string]bool{
		"withinAPILimit":     stats.APICallsCount < limits.MaxAPICalls,
		"withinStorageLimit": stats.StorageUsageBytes < limits.MaxStorageBytes,
		"canDowngrade":       stats.StorageUsageBytes < (5 * 1024 * 1024 * 1024),
	}
}

func (s *usageService) GetUsageStats(tenantID uuid.UUID) (*UsageStats, error) {
	tenant, err := s.tenantService.GetTenantByID(tenantID)
	if err != nil {
		return nil, err
	}

	var stats UsageStats
	if err := json.Unmarshal(tenant.Settings, &stats); err != nil {
		stats = UsageStats{
			LastUpdated: time.Now(),
			Limits:      s.getLimitsForTenant(tenant.Type),
		}
	}

	return &stats, nil
}

func (s *usageService) getLimitsForTenant(tenantType models.TenantType) Limits {
	switch tenantType {
	case models.PersonalTenant:
		return Limits{
			MaxAPICalls:     1000,
			MaxStorageBytes: 5 * 1024 * 1024 * 1024,
			MaxUsers:        1,
			StorageUnit:     "GB",
		}
	case models.TeamTenant:
		return Limits{
			MaxAPICalls:     10000,
			MaxStorageBytes: 50 * 1024 * 1024 * 1024,
			MaxUsers:        10,
			StorageUnit:     "GB",
		}
	case models.EnterpriseTenant:
		return Limits{
			MaxAPICalls:     100000,
			MaxStorageBytes: -1,
			MaxUsers:        -1,
			StorageUnit:     "TB",
		}
	default:
		return Limits{
			MaxAPICalls:     1000,
			MaxStorageBytes: 1 * 1024 * 1024 * 1024,
			MaxUsers:        1,
			StorageUnit:     "GB",
		}
	}
}
