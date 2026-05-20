package domain

import (
	"context"
	"time"
)

type OrderStatus string
type OrderType string

const (
	OrderStatusNew        OrderStatus = "new"         // Order created
	OrderStatusAssembling OrderStatus = "assembling"  // Storekeeper picking items
	OrderStatusAssembled  OrderStatus = "assembled"   // Ready for shipment
	OrderStatusShipped    OrderStatus = "shipped"     // Left warehouse
	OrderStatusCancelled  OrderStatus = "cancelled"

	OrderTypeRegular      OrderType   = "regular"
	OrderTypeCito         OrderType   = "cito"
)

// Order — заказ на отгрузку.
type Order struct {
	ID             string      `json:"id"`
	OrderNumber    string      `json:"order_number"`
	CustomerName   string      `json:"customer_name"`
	Status         OrderStatus `json:"status"`
	OrderType      OrderType   `json:"order_type"`
	Priority       int         `json:"priority"` // 1-normal, 2-high, 3-urgent
	Items          []OrderItem `json:"items,omitempty"`
	CreatedBy      int         `json:"created_by"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// OrderItem — позиция в заказе.
type OrderItem struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	PickedQty  int     `json:"picked_qty"` // Only filled during assembly
	BatchID    *string `json:"batch_id"`
	MosBlocked bool    `json:"mos_blocked"`
}

// BatchAllocation — серия и выделенное количество.
type BatchAllocation struct {
	BatchID string `json:"batch_id"`
	Qty     int    `json:"qty"`
}

// OrderRepository — интерфейс для работы с заказами.
type OrderRepository interface {
	List(ctx context.Context, limit, offset int) ([]Order, int, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, o *Order) error
	Update(ctx context.Context, o *Order) error
	Delete(ctx context.Context, id string) error
	GetMonthlyTurnover(ctx context.Context, productID string) (int, error)
}

// OrderService — бизнес-логика заказов.
type OrderService interface {
	ListOrders(ctx context.Context, limit, offset int) ([]Order, int, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	CreateOrder(ctx context.Context, o *Order) (*Order, error)
	UpdateStatus(ctx context.Context, callerRole UserRole, id string, status OrderStatus) error
}
