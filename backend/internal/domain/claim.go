package domain

import (
	"context"
	"time"
)

type ClaimStatus string

const (
	ClaimStatusNew        ClaimStatus = "new"
	ClaimStatusUnderReview ClaimStatus = "review"
	ClaimStatusResolved    ClaimStatus = "resolved"
	ClaimStatusRejected    ClaimStatus = "rejected"
)

// Claim — претензия или дефект.
type Claim struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	InboundID   *string     `json:"inbound_id,omitempty"` // Related inbound
	OrderID     *string     `json:"order_id,omitempty"`   // Related order
	Status      ClaimStatus `json:"status"`
	CreatedBy   int         `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ClaimRepository — интерфейс для работы с претензиями.
type ClaimRepository interface {
	List(ctx context.Context, limit, offset int) ([]Claim, int, error)
	GetByID(ctx context.Context, id string) (*Claim, error)
	Create(ctx context.Context, c *Claim) error
	Update(ctx context.Context, c *Claim) error
}

// ClaimService — бизнес-логика претензий.
type ClaimService interface {
	ListClaims(ctx context.Context, limit, offset int) ([]Claim, int, error)
	GetClaim(ctx context.Context, id string) (*Claim, error)
	CreateClaim(ctx context.Context, c *Claim) (*Claim, error)
	UpdateStatus(ctx context.Context, callerRole UserRole, id string, status ClaimStatus) error
}
