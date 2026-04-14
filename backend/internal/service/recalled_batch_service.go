package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type recalledBatchService struct {
	repo domain.RecalledBatchRepository
}

func NewRecalledBatchService(repo domain.RecalledBatchRepository) domain.RecalledBatchService {
	return &recalledBatchService{repo: repo}
}

func (s *recalledBatchService) ListRecalled(ctx context.Context, limit, offset int) ([]domain.RecalledBatch, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *recalledBatchService) CheckBatch(ctx context.Context, serial string) (bool, *domain.RecalledBatch, error) {
	batch, err := s.repo.GetBySerial(ctx, serial)
	if err != nil {
		return false, nil, err
	}
	if batch != nil {
		return true, batch, nil
	}
	return false, nil, nil
}

func (s *recalledBatchService) SyncRecalled(ctx context.Context, items []domain.RecalledBatch) error {
	for i := range items {
		if err := s.repo.Create(ctx, &items[i]); err != nil {
			return err
		}
	}
	return nil
}
