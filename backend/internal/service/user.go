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

// IsEmailTaken — делегирует проверку занятости email в репозиторий
func (s *userService) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	return s.repo.IsEmailTaken(ctx, email)
}
