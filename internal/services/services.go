package services

// Services holds all service instances
type Services struct {
	AuthService     AuthService
	UserService     UserService
	TenantService   TenantService
	SecurityService SecurityService
}
