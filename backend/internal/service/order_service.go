package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type orderService struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) domain.OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) ListOrders(ctx context.Context, limit, offset int) ([]domain.Order, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) CreateOrder(ctx context.Context, o *domain.Order) (*domain.Order, error) {
	o.Status = domain.OrderStatusNew
	if o.Priority == 0 {
		o.Priority = 1
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, callerRole domain.UserRole, id string, status domain.OrderStatus) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrOrderNotFound
	}

	// Status flow logic
	// Only Storekeeper or Manager can move to Assembling/Assembled
	// Only Manager or Admin can Ship
	switch status {
	case domain.OrderStatusAssembling, domain.OrderStatusAssembled:
		if callerRole == domain.RolePharmacist {
			return domain.ErrInsufficientPerms
		}
	case domain.OrderStatusShipped:
		if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
			return domain.ErrInsufficientPerms
		}
	}

	existing.Status = status
	return s.repo.Update(ctx, existing)
}
