package service

import (
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository"
)

// Service — агрегатор всех сервисов приложения
type Service struct {
	domain.UserService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repos.User),
	}
}
