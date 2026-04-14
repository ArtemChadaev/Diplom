package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type supplierService struct {
	repo domain.SupplierRepository
}

func NewSupplierService(repo domain.SupplierRepository) domain.SupplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) ListSuppliers(ctx context.Context, limit, offset int) ([]domain.Supplier, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *supplierService) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	sup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sup == nil {
		return nil, domain.ErrSupplierNotFound
	}
	return sup, nil
}

func (s *supplierService) CreateSupplier(ctx context.Context, callerRole domain.UserRole, sup *domain.Supplier) (*domain.Supplier, error) {
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return nil, domain.ErrInsufficientPerms
	}

	// Check INN uniqueness
	existing, err := s.repo.GetByINN(ctx, sup.INN)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrConflict
	}

	if err := s.repo.Create(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *supplierService) UpdateSupplier(ctx context.Context, callerRole domain.UserRole, sup *domain.Supplier) (*domain.Supplier, error) {
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return nil, domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, sup.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrSupplierNotFound
	}

	if err := s.repo.Update(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *supplierService) DeleteSupplier(ctx context.Context, callerRole domain.UserRole, id string) error {
	if callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrSupplierNotFound
	}

	return s.repo.Delete(ctx, id)
}
