package domain

import (
	"context"
	"time"
)

type InboundStatus string

const (
	InboundStatusDraft     InboundStatus = "draft"      // Plan created
	InboundStatusReceived  InboundStatus = "received"   // Physical arrival, verification in progress
	InboundStatusCompleted InboundStatus = "completed"  // Items put away, stock updated
	InboundStatusCancelled InboundStatus = "cancelled"
)

// InboundReceiving — приходная накладная / поступление.
type InboundReceiving struct {
	ID                string           `json:"id"`
	InvoiceNumber     string           `json:"invoice_number"`
	InvoiceDate       time.Time        `json:"invoice_date"`
	SupplierID        string           `json:"supplier_id"`
	Status            InboundStatus    `json:"status"`
	TotalAmount       float64          `json:"total_amount"`
	VATAmount         float64          `json:"vat_amount"`
	Notes             string           `json:"notes"`
	ReceivedBy        int              `json:"received_by"` // UserID
	Items             []InboundItem    `json:"items,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// InboundItem — позиция в приходной накладной.
type InboundItem struct {
	ID                string    `json:"id"`
	InboundID         string    `json:"inbound_id"`
	ProductID         string    `json:"product_id"`
	BatchNumber       string    `json:"batch_number"`
	ExpirationDate    time.Time `json:"expiration_date"`
	Quantity          int       `json:"quantity"`
	PriceNetto        float64   `json:"price_netto"`
	VATRate           float64   `json:"vat_rate"`
	PriceBrutto       float64   `json:"price_brutto"`
	CertificateNumber string    `json:"cert_number"`
	ZoneID            string    `json:"zone_id"` // Target zone for storage
}

// InboundRepository — интерфейс для работы с поступлениями.
type InboundRepository interface {
	List(ctx context.Context, limit, offset int) ([]InboundReceiving, int, error)
	GetByID(ctx context.Context, id string) (*InboundReceiving, error)
	Create(ctx context.Context, i *InboundReceiving) error
	Update(ctx context.Context, i *InboundReceiving) error
	Delete(ctx context.Context, id string) error
	
	// Items
	AddItems(ctx context.Context, inboundID string, items []InboundItem) error
	UpdateItem(ctx context.Context, item *InboundItem) error
	RemoveItem(ctx context.Context, id string) error
}

// InboundService — бизнес-логика поступлений.
type InboundService interface {
	ListInbounds(ctx context.Context, limit, offset int) ([]InboundReceiving, int, error)
	GetInbound(ctx context.Context, id string) (*InboundReceiving, error)
	CreateInbound(ctx context.Context, callerRole UserRole, i *InboundReceiving) (*InboundReceiving, error)
	UpdateInbound(ctx context.Context, callerRole UserRole, i *InboundReceiving) (*InboundReceiving, error)
	UpdateStatus(ctx context.Context, callerRole UserRole, id string, status InboundStatus) error
	DeleteInbound(ctx context.Context, callerRole UserRole, id string) error
}
