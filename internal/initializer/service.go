package initializer

import (
	"identity-service/internal/services"
)

type Services struct {
    UserService  services.UserService
    TokenService services.TokenService
}

func InitServices(repos *Repositories) *Services {
    return &Services{
        UserService:  services.NewUserService(repos.UserRepo),
        TokenService: services.NewTokenService(repos.JwtTokenRepo),
    }
}