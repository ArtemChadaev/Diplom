package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type productService struct {
	repo      domain.ProductRepository
	orderRepo domain.OrderRepository
	batchRepo domain.BatchRepository
}

func NewProductService(
	repo domain.ProductRepository,
	orderRepo domain.OrderRepository,
	batchRepo domain.BatchRepository,
) domain.ProductService {
	return &productService{
		repo:      repo,
		orderRepo: orderRepo,
		batchRepo: batchRepo,
	}
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

func (s *productService) CheckReorderPoint(ctx context.Context, productID string) (*domain.ROPResult, error) {
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, domain.ErrProductNotFound
	}

	turnover, err := s.orderRepo.GetMonthlyTurnover(ctx, productID)
	if err != nil {
		return nil, err
	}

	currentStock, err := s.batchRepo.GetTotalStock(ctx, productID)
	if err != nil {
		return nil, err
	}

	dailyUsage := float64(turnover) / 30.0
	rop := int(dailyUsage*float64(product.LeadTimeDays)) + product.SafetyStockQty
	needsReorder := currentStock <= rop

	reorderQty := 0
	if needsReorder {
		reorderQty = product.MaxStockQty - currentStock
		if reorderQty < 0 {
			reorderQty = 0
		}
	}

	return &domain.ROPResult{
		ProductID:       product.ID,
		SKU:             product.SKU,
		Name:            product.Name,
		CurrentStock:    currentStock,
		SafetyStock:     product.SafetyStockQty,
		MaxStock:        product.MaxStockQty,
		MonthlyTurnover: turnover,
		DailyUsage:      dailyUsage,
		ROP:             rop,
		NeedsReorder:    needsReorder,
		ReorderQty:      reorderQty,
	}, nil
}

func (s *productService) RunReorderCheckAll(ctx context.Context) ([]domain.ROPResult, error) {
	products, _, err := s.repo.List(ctx, domain.ProductFilter{Limit: 1000})
	if err != nil {
		return nil, err
	}

	results := make([]domain.ROPResult, 0, len(products))
	for _, p := range products {
		res, err := s.CheckReorderPoint(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		results = append(results, *res)
	}

	return results, nil
}
