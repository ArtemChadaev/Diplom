package domain

import (
	"context"
	"time"
)

type BatchStatus string

const (
	BatchStatusQuarantine BatchStatus = "quarantine"
	BatchStatusAvailable  BatchStatus = "available"
	BatchStatusRejected   BatchStatus = "rejected"
	BatchStatusBlocked    BatchStatus = "blocked"
)

// Batch — конкретная серия товара на складе.
type Batch struct {
	ID               string      `json:"id"`
	ProductID        string      `json:"product_id"`
	ZoneID           *string     `json:"zone_id,omitempty"`
	SerialNumber     string      `json:"serial_number"`
	ManufactureDate  time.Time   `json:"manufacture_date"`
	ExpiryDate       time.Time   `json:"expiry_date"`
	Quantity         int         `json:"quantity"`
	Status           BatchStatus `json:"status"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

// BatchRepository — интерфейс для работы с остатками по сериям.
type BatchRepository interface {
	List(ctx context.Context, filter BatchFilter) ([]Batch, int, error)
	GetByID(ctx context.Context, id string) (*Batch, error)
	Update(ctx context.Context, b *Batch) error
	Delete(ctx context.Context, id string) error
	BlockAllByProductID(ctx context.Context, productID string) error
	ListAvailableSorted(ctx context.Context, productID string) ([]Batch, error)
	GetTotalStock(ctx context.Context, productID string) (int, error)
}

type BatchFilter struct {
	ProductID string
	ZoneID    string
	Status    BatchStatus
	Limit     int
	Offset    int
}

// BatchService — бизнес-логика управления сериями.
type BatchService interface {
	ListBatches(ctx context.Context, filter BatchFilter) ([]Batch, int, error)
	GetBatch(ctx context.Context, id string) (*Batch, error)
	UpdateStatus(ctx context.Context, callerRole UserRole, id string, status BatchStatus) error
	TransferBatch(ctx context.Context, callerRole UserRole, id string, targetZoneID string) error
}
