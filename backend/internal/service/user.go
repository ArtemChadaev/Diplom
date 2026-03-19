package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

// IsLoginTaken — делегирует проверку занятости логина в репозиторий
func (s *userService) IsLoginTaken(ctx context.Context, login string) (bool, error) {
	return s.repo.IsLoginTaken(ctx, login)
}
