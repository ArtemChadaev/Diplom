package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

var employeeCodeRegex = regexp.MustCompile("^[A-Z]{2}-[0-9]{3}$")

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

	if input.EmployeeCode != nil {
		if !employeeCodeRegex.MatchString(*input.EmployeeCode) {
			return nil, domain.ErrInvalidEmployeeCode
		}
	}

	profile, err := s.repo.FindByUserID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeProfileNotFound) {
			// If not found, CREATE it!
			newProfile := &domain.EmployeeProfile{
				UserID: uint(targetUserID),
			}
			if input.EmployeeCode != nil {
				newProfile.EmployeeCode = *input.EmployeeCode
			}
			if input.FullName != nil {
				newProfile.FullName = *input.FullName
			}
			if input.CorporateEmail != nil {
				newProfile.CorporateEmail = *input.CorporateEmail
			}
			if input.Phone != nil {
				newProfile.Phone = *input.Phone
			}
			if input.Position != nil {
				newProfile.Position = *input.Position
			}
			if input.Department != nil {
				newProfile.Department = *input.Department
			}
			if input.BirthDate != nil {
				newProfile.BirthDate = *input.BirthDate
			}
			if input.AvatarURL != nil {
				newProfile.AvatarURL = *input.AvatarURL
			}
			if input.HireDate != nil {
				newProfile.HireDate = *input.HireDate
			}
			if input.DismissalDate != nil {
				newProfile.DismissalDate = input.DismissalDate
			}
			if input.MedicalBookScanURL != nil {
				newProfile.MedicalBookScanURL = *input.MedicalBookScanURL
			}
			if input.SpecialZoneAccess != nil {
				newProfile.SpecialZoneAccess = *input.SpecialZoneAccess
			}
			if input.GDPTrainingHistory != nil {
				var gdp []domain.GDPTrainingRecord
				_ = json.Unmarshal(input.GDPTrainingHistory, &gdp)
				newProfile.GDPTrainingHistory = gdp
			}

			created, err := s.repo.Create(ctx, newProfile)
			if err != nil {
				return nil, err
			}

			logger.FromContext(ctx).Info("admin created employee profile",
				"admin_id", callerID,
				"target_user_id", targetUserID,
				"profile_id", int(created.ID),
			)
			return created, nil
		}
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

func (s *employeeProfileService) PatchSelfProfile(
	ctx context.Context, userID int, input domain.UpdateEmployeeProfileInput,
) (*domain.EmployeeProfile, error) {
	profile, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, int(profile.ID), input)
	if err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("user patched own profile",
		"user_id", userID,
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
