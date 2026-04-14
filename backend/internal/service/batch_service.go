package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type batchService struct {
	repo domain.BatchRepository
}

func NewBatchService(repo domain.BatchRepository) domain.BatchService {
	return &batchService{repo: repo}
}

func (s *batchService) ListBatches(ctx context.Context, filter domain.BatchFilter) ([]domain.Batch, int, error) {
	return s.repo.List(ctx, filter)
}

func (s *batchService) GetBatch(ctx context.Context, id string) (*domain.Batch, error) {
	batch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, domain.ErrBatchNotFound
	}
	return batch, nil
}

func (s *batchService) UpdateStatus(ctx context.Context, callerRole domain.UserRole, id string, status domain.BatchStatus) error {
	// Only Manager or Admin can manually change batch status
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return domain.ErrInsufficientPerms
	}

	batch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if batch == nil {
		return domain.ErrBatchNotFound
	}

	batch.Status = status
	return s.repo.Update(ctx, batch)
}

func (s *batchService) TransferBatch(ctx context.Context, callerRole domain.UserRole, id string, targetZoneID string) error {
	// Storekeeper, WarehouseManager and Admin can move goods; Pharmacist cannot
	if callerRole == domain.RolePharmacist {
		return domain.ErrInsufficientPerms
	}
	batch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if batch == nil {
		return domain.ErrBatchNotFound
	}

	batch.ZoneID = &targetZoneID
	return s.repo.Update(ctx, batch)
}
