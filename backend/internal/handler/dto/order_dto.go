package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type OrderResponse struct {
	ID           string             `json:"id"`
	OrderNumber  string             `json:"order_number"`
	CustomerName string             `json:"customer_name"`
	Status       domain.OrderStatus `json:"status"`
	Priority     int                `json:"priority"`
	Items        []OrderItemResponse `json:"items,omitempty"`
	CreatedBy    int                `json:"created_by"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	PickedQty int    `json:"picked_qty"`
}

func ToOrderResponse(o domain.Order) OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			PickedQty: item.PickedQty,
		}
	}

	return OrderResponse{
		ID:           o.ID,
		OrderNumber:  o.OrderNumber,
		CustomerName: o.CustomerName,
		Status:       o.Status,
		Priority:     o.Priority,
		Items:        items,
		CreatedBy:    o.CreatedBy,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

type CreateOrderRequest struct {
	CustomerName string           `json:"customer_name" validate:"required"`
	Priority     int              `json:"priority" validate:"min=1,max=3"`
	Items        []CreateOrderItem `json:"items" validate:"required,dive"`
}

type CreateOrderItem struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

func (r CreateOrderRequest) ToDomain() domain.Order {
	items := make([]domain.OrderItem, len(r.Items))
	for i, item := range r.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return domain.Order{
		CustomerName: r.CustomerName,
		Priority:     r.Priority,
		Items:        items,
	}
}

type OrderListResponse struct {
	Total  int             `json:"total"`
	Orders []OrderResponse `json:"orders"`
}

type UpdateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status" validate:"required"`
}
