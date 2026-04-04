package service

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/mailer"
	"github.com/ima/diplom-backend/internal/repository"
)

// Service — агрегатор всех сервисов приложения
type Service struct {
	Auth            domain.AuthService
	Token           domain.TokenService
	EmployeeProfile domain.EmployeeProfileService
}

func NewService(
	repos *repository.Repository,
	jwtSecret string,
	googleClientID string,
	otpHMACSecret string,
	m mailer.Mailer,
) *Service {
	// Standard TTL configuration: Access Token 15m, Refresh Token 15d
	tokenSvc := NewTokenService(jwtSecret, 15*time.Minute, 15*24*time.Hour)

	return &Service{
		Auth:            NewAuthService(repos.User, repos.Session, repos.OTP, tokenSvc, 15*24*time.Hour, googleClientID, m, otpHMACSecret),
		Token:           tokenSvc,
		EmployeeProfile: NewEmployeeProfileService(repos.EmployeeProfile),
	}
}

