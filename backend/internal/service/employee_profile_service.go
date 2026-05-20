package service

import (
	"context"
	"regexp"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

var employeeCodeRegex = regexp.MustCompile("^[A-Z]{2}-[0-9]{3}$")

// CheckGDPValid verifies if the given employee profile has a valid, non-expired GDP training.
func CheckGDPValid(profile *domain.EmployeeProfile) error {
	if profile == nil || len(profile.GDPTrainingHistory) == 0 {
		return domain.ErrGDPTrainingRequired
	}

	var latestRecord *domain.GDPTrainingRecord
	var latestTime time.Time

	for i := range profile.GDPTrainingHistory {
		rec := &profile.GDPTrainingHistory[i]
		t, err := time.Parse("2006-01-02", rec.Date)
		if err != nil {
			continue
		}
		if latestRecord == nil || t.After(latestTime) {
			latestRecord = rec
			latestTime = t
		}
	}

	if latestRecord == nil {
		latestRecord = &profile.GDPTrainingHistory[len(profile.GDPTrainingHistory)-1]
		var err error
		latestTime, err = time.Parse("2006-01-02", latestRecord.Date)
		if err != nil {
			if latestRecord.Result != "pass" {
				return domain.ErrGDPTrainingRequired
			}
			return nil
		}
	}

	if latestRecord.Result != "pass" {
		return domain.ErrGDPTrainingRequired
	}

	if time.Since(latestTime) > 365*24*time.Hour {
		return domain.ErrGDPTrainingExpired
	}

	return nil
}

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

func (s *employeeProfileService) GetSelfProfile(ctx context.Context, userID int) (*domain.EmployeeProfile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *employeeProfileService) CreateProfile(
	ctx context.Context, callerID int, callerRole domain.UserRole, input domain.CreateEmployeeProfileInput,
) (*domain.EmployeeProfile, error) {
	if err := requireAdmin(callerRole); err != nil {
		return nil, err
	}

	if !employeeCodeRegex.MatchString(input.EmployeeCode) {
		return nil, domain.ErrInvalidEmployeeCode
	}

	created, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("admin created employee profile",
		"admin_id", callerID,
		"target_user_id", int(input.UserID),
		"profile_id", int(created.ID),
	)

	return created, nil
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
