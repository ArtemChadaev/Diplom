package service

import (
	"context"
	"strconv"

	"github.com/ima/diplom-backend/internal/domain"
)

type orderService struct {
	repo         domain.OrderRepository
	batchRepo    domain.BatchRepository
	productRepo  domain.ProductRepository
	settingsRepo domain.SystemSettingsRepository
}

func NewOrderService(
	repo domain.OrderRepository,
	batchRepo domain.BatchRepository,
	productRepo domain.ProductRepository,
	settingsRepo domain.SystemSettingsRepository,
) domain.OrderService {
	return &orderService{
		repo:         repo,
		batchRepo:    batchRepo,
		productRepo:  productRepo,
		settingsRepo: settingsRepo,
	}
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
	if o.OrderType == "" {
		o.OrderType = domain.OrderTypeRegular
	}

	mosPercentStr, err := s.settingsRepo.Get(ctx, "mos_percent")
	if err != nil {
		return nil, err
	}
	mosPercent := 60
	if mosPercentStr != "" {
		if val, err := strconv.Atoi(mosPercentStr); err == nil {
			mosPercent = val
		}
	}

	var allocatedItems []domain.OrderItem

	for _, item := range o.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, domain.ErrProductNotFound
		}

		// а) Проверка VEN-категории
		if product.VenCategory == domain.VenV && o.OrderType != domain.OrderTypeCito {
			return nil, domain.ErrVCategoryReserveOnly
		}

		// б) Расчет MOS
		monthlyTurnover, err := s.repo.GetMonthlyTurnover(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		mosLimit := (monthlyTurnover * mosPercent) / 100

		// Текущий доступный остаток
		totalStock, err := s.batchRepo.GetTotalStock(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Проверка лимита MOS
		isMosBlocked := false
		if totalStock < mosLimit {
			if o.OrderType == domain.OrderTypeRegular {
				return nil, domain.ErrInsufficientStock
			}
			isMosBlocked = true
		}

		// в) Аллокация по FEFO
		batches, err := s.batchRepo.ListAvailableSorted(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		remainingToAllocate := item.Quantity
		var allocations []domain.BatchAllocation

		for _, b := range batches {
			if remainingToAllocate <= 0 {
				break
			}
			take := b.Quantity
			if take > remainingToAllocate {
				take = remainingToAllocate
			}
			allocations = append(allocations, domain.BatchAllocation{
				BatchID: b.ID,
				Qty:     take,
			})
			remainingToAllocate -= take
		}

		if remainingToAllocate > 0 {
			return nil, domain.ErrInsufficientStock
		}

		// Снижаем остатки в сериях и создаем отдельные OrderItem
		for _, alloc := range allocations {
			batch, err := s.batchRepo.GetByID(ctx, alloc.BatchID)
			if err != nil {
				return nil, err
			}
			batch.Quantity -= alloc.Qty
			if err := s.batchRepo.Update(ctx, batch); err != nil {
				return nil, err
			}

			allocatedItems = append(allocatedItems, domain.OrderItem{
				ProductID:  item.ProductID,
				Quantity:   alloc.Qty,
				BatchID:    &alloc.BatchID,
				MosBlocked: isMosBlocked,
			})
		}
	}

	o.Items = allocatedItems

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
