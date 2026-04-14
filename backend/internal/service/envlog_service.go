package service

import (
	"context"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type environmentLogService struct {
	repo domain.EnvironmentLogRepository
}

func NewEnvironmentLogService(repo domain.EnvironmentLogRepository) domain.EnvironmentLogService {
	return &environmentLogService{repo: repo}
}

func (s *environmentLogService) ListLogs(ctx context.Context, zoneID string, limit, offset int) ([]domain.EnvironmentLog, int, error) {
	return s.repo.List(ctx, zoneID, limit, offset)
}

func (s *environmentLogService) RecordLogs(ctx context.Context, userID int, logs []domain.EnvironmentLog) error {
	for i := range logs {
		logs[i].RecordedBy = userID
		if logs[i].RecordedAt.IsZero() {
			logs[i].RecordedAt = time.Now()
		}
		if err := s.repo.Create(ctx, &logs[i]); err != nil {
			return err
		}
	}
	return nil
}
