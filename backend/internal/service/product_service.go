package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{repo: repo}
}

func (s *productService) ListProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	return s.repo.List(ctx, filter)
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, domain.ErrProductNotFound
	}
	return p, nil
}

func (s *productService) CreateProduct(ctx context.Context, callerRole domain.UserRole, p *domain.Product) (*domain.Product, error) {
	// Only Admin or WarehouseManager can create products
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return nil, domain.ErrInsufficientPerms
	}

	// Check SKU uniqueness
	existing, err := s.repo.GetBySKU(ctx, p.SKU)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrConflict
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) UpdateProduct(ctx context.Context, callerRole domain.UserRole, p *domain.Product) (*domain.Product, error) {
	if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
		return nil, domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrProductNotFound
	}

	// Update logic: preserve non-updatable fields if any. 
	// For now, allow updating everything except ID.
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) DeleteProduct(ctx context.Context, callerRole domain.UserRole, id string) error {
	if callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrProductNotFound
	}

	return s.repo.Delete(ctx, id)
}
