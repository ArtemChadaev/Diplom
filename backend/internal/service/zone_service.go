package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type zoneService struct {
	repo domain.ZoneRepository
}

func NewZoneService(repo domain.ZoneRepository) domain.ZoneService {
	return &zoneService{repo: repo}
}

func (s *zoneService) ListZones(ctx context.Context) ([]domain.Zone, error) {
	return s.repo.List(ctx)
}

func (s *zoneService) GetZone(ctx context.Context, id string) (*domain.Zone, error) {
	z, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if z == nil {
		return nil, domain.ErrZoneNotFound
	}
	return z, nil
}

func (s *zoneService) CreateZone(ctx context.Context, callerRole domain.UserRole, z *domain.Zone) (*domain.Zone, error) {
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return nil, domain.ErrInsufficientPerms
	}

	if err := s.repo.Create(ctx, z); err != nil {
		return nil, err
	}
	return z, nil
}

func (s *zoneService) UpdateZone(ctx context.Context, callerRole domain.UserRole, z *domain.Zone) (*domain.Zone, error) {
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return nil, domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, z.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrZoneNotFound
	}

	if err := s.repo.Update(ctx, z); err != nil {
		return nil, err
	}
	return z, nil
}

func (s *zoneService) DeleteZone(ctx context.Context, callerRole domain.UserRole, id string) error {
	if callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrZoneNotFound
	}

	return s.repo.Delete(ctx, id)
}
