package service

import (
	"context"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type environmentLogService struct {
	repo        domain.EnvironmentLogRepository
	zoneRepo    domain.ZoneRepository
	profileRepo domain.EmployeeProfileRepository
}

func NewEnvironmentLogService(
	repo domain.EnvironmentLogRepository,
	zoneRepo domain.ZoneRepository,
	profileRepo domain.EmployeeProfileRepository,
) domain.EnvironmentLogService {
	return &environmentLogService{
		repo:        repo,
		zoneRepo:    zoneRepo,
		profileRepo: profileRepo,
	}
}

func (s *environmentLogService) ListLogs(ctx context.Context, zoneID string, limit, offset int) ([]domain.EnvironmentLog, int, error) {
	return s.repo.List(ctx, zoneID, limit, offset)
}

func (s *environmentLogService) RecordLogs(ctx context.Context, userID int, logs []domain.EnvironmentLog) error {
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return domain.ErrGDPTrainingRequired
	}
	if err := CheckGDPValid(profile); err != nil {
		return err
	}

	for i := range logs {
		logs[i].RecordedBy = userID
		if logs[i].RecordedAt.IsZero() {
			logs[i].RecordedAt = time.Now()
		}

		exists, err := s.repo.ExistsByZoneShiftDate(ctx, logs[i].ZoneID, logs[i].Shift, logs[i].RecordedAt)
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrEnvLogDuplicateShift
		}

		// Получаем информацию о зоне для проверки температурного диапазона
		zone, err := s.zoneRepo.GetByID(ctx, logs[i].ZoneID)
		if err != nil {
			return err
		}
		if zone == nil {
			return domain.ErrZoneNotFound
		}

		if logs[i].Temperature < zone.TemperatureMin || logs[i].Temperature > zone.TemperatureMax {
			return domain.ErrZoneTempViolation
		}

		if err := s.repo.Create(ctx, &logs[i]); err != nil {
			return err
		}
	}
	return nil
}
