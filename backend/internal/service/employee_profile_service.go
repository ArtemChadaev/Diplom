package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

type employeeProfileService struct {
	repo domain.EmployeeProfileRepository
}

func NewEmployeeProfileService(repo domain.EmployeeProfileRepository) domain.EmployeeProfileService {
	return &employeeProfileService{repo: repo}
}

func requireAdmin(callerRole domain.UserRole) error {
	if callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return nil
}

func (s *employeeProfileService) GetProfile(
	ctx context.Context, callerID int, callerRole domain.UserRole, targetUserID int,
) (*domain.EmployeeProfile, error) {
	if err := requireAdmin(callerRole); err != nil {
		return nil, err
	}
	return s.repo.FindByUserID(ctx, targetUserID)
}

func (s *employeeProfileService) UpdateProfile(
	ctx context.Context, callerID int, callerRole domain.UserRole,
	targetUserID int, input domain.UpdateEmployeeProfileInput,
) (*domain.EmployeeProfile, error) {
	if err := requireAdmin(callerRole); err != nil {
		return nil, err
	}

	profile, err := s.repo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, int(profile.ID), input)
	if err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("admin updated employee profile",
		"admin_id", callerID,
		"target_user_id", targetUserID,
		"profile_id", int(profile.ID),
	)

	return updated, nil
}

func (s *employeeProfileService) ListProfiles(
	ctx context.Context, callerID int, callerRole domain.UserRole, limit, offset int,
) ([]domain.EmployeeProfile, error) {
	if err := requireAdmin(callerRole); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.List(ctx, limit, offset)
}
