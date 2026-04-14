package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type systemSettingsService struct {
	repo domain.SystemSettingsRepository
}

func NewSystemSettingsService(repo domain.SystemSettingsRepository) domain.SystemSettingsService {
	return &systemSettingsService{repo: repo}
}

func (s *systemSettingsService) GetSetting(ctx context.Context, key string) (string, error) {
	return s.repo.Get(ctx, key)
}

func (s *systemSettingsService) UpdateSetting(ctx context.Context, callerRole domain.UserRole, key, value string) error {
	// Only Admin can change system settings
	if callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return s.repo.Set(ctx, key, value)
}

func (s *systemSettingsService) ListSettings(ctx context.Context) ([]domain.SystemSetting, error) {
	return s.repo.List(ctx)
}
