package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type ClaimResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	InboundID   *string            `json:"inbound_id,omitempty"`
	OrderID     *string            `json:"order_id,omitempty"`
	Status      domain.ClaimStatus `json:"status"`
	CreatedBy   int                `json:"created_by"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func ToClaimResponse(c domain.Claim) ClaimResponse {
	return ClaimResponse{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		InboundID:   c.InboundID,
		OrderID:     c.OrderID,
		Status:      c.Status,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type CreateClaimRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	InboundID   *string `json:"inbound_id" validate:"omitempty,uuid"`
	OrderID     *string `json:"order_id" validate:"omitempty,uuid"`
}

func (r CreateClaimRequest) ToDomain() domain.Claim {
	return domain.Claim{
		Title:       r.Title,
		Description: r.Description,
		InboundID:   r.InboundID,
		OrderID:     r.OrderID,
	}
}

type ClaimListResponse struct {
	Total  int             `json:"total"`
	Claims []ClaimResponse `json:"claims"`
}

type UpdateClaimStatusRequest struct {
	Status domain.ClaimStatus `json:"status" validate:"required"`
}
