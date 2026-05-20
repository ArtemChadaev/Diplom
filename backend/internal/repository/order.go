package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) List(ctx context.Context, limit, offset int) ([]domain.Order, int, error) {
	var daos []dao.OrderDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.OrderDAO{})
	
	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	if err := model.Order("priority DESC, created_at ASC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.Order, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var d dao.OrderDAO
	if err := r.db.WithContext(ctx).Preload("Items").First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *orderRepository) Create(ctx context.Context, o *domain.Order) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	for idx := range o.Items {
		if o.Items[idx].ID == "" {
			o.Items[idx].ID = uuid.NewString()
		}
		o.Items[idx].OrderID = o.ID
	}
	d := dao.FromOrderDomain(*o)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *orderRepository) Update(ctx context.Context, o *domain.Order) error {
	d := dao.FromOrderDomain(*o)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&d).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&dao.OrderItemDAO{}, "order_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&dao.OrderDAO{}, "id = ?", id).Error
	})
}

func (r *orderRepository) GetMonthlyTurnover(ctx context.Context, productID string) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.product_id = ? AND orders.created_at >= NOW() - INTERVAL '30 days' AND orders.status != ?", productID, domain.OrderStatusCancelled).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Row().Scan(&total)
	if err != nil {
		return 0, err
	}
	return int(total), nil
}
