package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type OrderDAO struct {
	ID           string             `gorm:"column:id;primaryKey"`
	OrderNumber  string             `gorm:"column:order_number;uniqueIndex"`
	CustomerName string             `gorm:"column:customer_name"`
	Status       domain.OrderStatus `gorm:"column:status"`
	OrderType    domain.OrderType   `gorm:"column:order_type"`
	Priority     int                `gorm:"column:priority"`
	CreatedBy    int                `gorm:"column:created_by"`
	CreatedAt    time.Time          `gorm:"column:created_at"`
	UpdatedAt    time.Time          `gorm:"column:updated_at"`
	Items        []OrderItemDAO     `gorm:"foreignKey:OrderID"`
}

func (OrderDAO) TableName() string {
	return "orders"
}

func (o OrderDAO) ToDomain() domain.Order {
	items := make([]domain.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = item.ToDomain()
	}

	return domain.Order{
		ID:           o.ID,
		OrderNumber:  o.OrderNumber,
		CustomerName: o.CustomerName,
		Status:       o.Status,
		OrderType:    o.OrderType,
		Priority:     o.Priority,
		CreatedBy:    o.CreatedBy,
		Items:        items,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

type OrderItemDAO struct {
	ID         string    `gorm:"column:id;primaryKey"`
	OrderID    string    `gorm:"column:order_id;index"`
	ProductID  string    `gorm:"column:product_id"`
	Quantity   int       `gorm:"column:quantity"`
	PickedQty  int       `gorm:"column:picked_qty"`
	BatchID    *string   `gorm:"column:batch_id;index"`
	MosBlocked bool      `gorm:"column:mos_blocked"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (OrderItemDAO) TableName() string {
	return "order_items"
}

func (i OrderItemDAO) ToDomain() domain.OrderItem {
	return domain.OrderItem{
		ID:         i.ID,
		OrderID:    i.OrderID,
		ProductID:  i.ProductID,
		Quantity:   i.Quantity,
		PickedQty:  i.PickedQty,
		BatchID:    i.BatchID,
		MosBlocked: i.MosBlocked,
	}
}

func FromOrderDomain(o domain.Order) OrderDAO {
	items := make([]OrderItemDAO, len(o.Items))
	for i, item := range o.Items {
		items[i] = FromOrderItemDomain(item)
	}

	return OrderDAO{
		ID:           o.ID,
		OrderNumber:  o.OrderNumber,
		CustomerName: o.CustomerName,
		Status:       o.Status,
		OrderType:    o.OrderType,
		Priority:     o.Priority,
		CreatedBy:    o.CreatedBy,
		Items:        items,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

func FromOrderItemDomain(i domain.OrderItem) OrderItemDAO {
	return OrderItemDAO{
		ID:         i.ID,
		OrderID:    i.OrderID,
		ProductID:  i.ProductID,
		Quantity:   i.Quantity,
		PickedQty:  i.PickedQty,
		BatchID:    i.BatchID,
		MosBlocked: i.MosBlocked,
	}
}
