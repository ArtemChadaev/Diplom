package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type inboundService struct {
	repo domain.InboundRepository
}

func NewInboundService(repo domain.InboundRepository) domain.InboundService {
	return &inboundService{repo: repo}
}

func (s *inboundService) ListInbounds(ctx context.Context, limit, offset int) ([]domain.InboundReceiving, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *inboundService) GetInbound(ctx context.Context, id string) (*domain.InboundReceiving, error) {
	inbound, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inbound == nil {
		return nil, domain.ErrInboundNotFound
	}
	return inbound, nil
}

func (s *inboundService) CreateInbound(ctx context.Context, callerRole domain.UserRole, i *domain.InboundReceiving) (*domain.InboundReceiving, error) {
	if callerRole == domain.RoleStorekeeper || callerRole == domain.RolePharmacist {
		return nil, domain.ErrInsufficientPerms
	}

	i.Status = domain.InboundStatusDraft
	if err := s.repo.Create(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *inboundService) UpdateInbound(ctx context.Context, callerRole domain.UserRole, i *domain.InboundReceiving) (*domain.InboundReceiving, error) {
	if callerRole == domain.RoleStorekeeper || callerRole == domain.RolePharmacist {
		return nil, domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, i.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrInboundNotFound
	}

	if existing.Status == domain.InboundStatusCompleted {
		return nil, domain.ErrConflict // Cannot update completed inbound
	}

	if err := s.repo.Update(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *inboundService) UpdateStatus(ctx context.Context, callerRole domain.UserRole, id string, status domain.InboundStatus) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrInboundNotFound
	}

	// Simple status flow validation
	switch status {
	case domain.InboundStatusReceived:
		// Any role except pharmacist? Standard is WarehouseManager or Storekeeper
		if callerRole == domain.RolePharmacist {
			return domain.ErrInsufficientPerms
		}
	case domain.InboundStatusCompleted:
		if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
			return domain.ErrInsufficientPerms
		}
	case domain.InboundStatusCancelled:
		if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
			return domain.ErrInsufficientPerms
		}
	}

	existing.Status = status
	return s.repo.Update(ctx, existing)
}

func (s *inboundService) DeleteInbound(ctx context.Context, callerRole domain.UserRole, id string) error {
	if callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrInboundNotFound
	}

	if existing.Status == domain.InboundStatusCompleted {
		return domain.ErrConflict // Cannot delete completed inbound
	}

	return s.repo.Delete(ctx, id)
}
