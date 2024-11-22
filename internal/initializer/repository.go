package initializer

import (
	"identity-service/db"
	"identity-service/internal/repositories"
)

type Repositories struct {
    UserRepo      repositories.UserRepository
    JwtTokenRepo  repositories.JwtTokenRepository
}

func InitRepositories() *Repositories {
    return &Repositories{
        UserRepo:     repositories.NewUserRepository(db.DB),
        JwtTokenRepo: repositories.NewJwtTokenRepository(db.DB),
    }
}