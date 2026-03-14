package service

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository"
)

// Service — агрегатор всех сервисов приложения
type Service struct {
	User  domain.UserService
	Auth  domain.AuthService
	Token domain.TokenService
}

func NewService(repos *repository.Repository, jwtSecret string) *Service {
	// Standard TTL configuration: Access Token 15m, Refresh Token 15d
	tokenSvc := NewTokenService(jwtSecret, 15*time.Minute, 15*24*time.Hour)

	return &Service{
		User:  NewUserService(repos.User),
		Auth:  NewAuthService(repos.User, repos.Session, tokenSvc, 15*24*time.Hour),
		Token: tokenSvc,
	}
}
