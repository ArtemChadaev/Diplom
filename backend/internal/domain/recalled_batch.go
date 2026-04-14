package domain

import (
	"context"
	"time"
)

// RecalledBatch — сведения об изъятых из обращения сериях (из реестра).
type RecalledBatch struct {
	ID           string    `json:"id"`
	SerialNumber string    `json:"serial_number"`
	ProductName  string    `json:"product_name"`
	RuNumber     string    `json:"ru_number"`
	RecallReason string    `json:"recall_reason"`
	IssuedAt     time.Time `json:"issued_at"`
	SyncedAt     time.Time `json:"synced_at"`
}

// RecalledBatchRepository — интерфейс для работы с реестром изъятия.
type RecalledBatchRepository interface {
	List(ctx context.Context, limit, offset int) ([]RecalledBatch, int, error)
	GetBySerial(ctx context.Context, serial string) (*RecalledBatch, error)
	Create(ctx context.Context, b *RecalledBatch) error
}

// RecalledBatchService — бизнес-логика реестра изъятия.
type RecalledBatchService interface {
	ListRecalled(ctx context.Context, limit, offset int) ([]RecalledBatch, int, error)
	CheckBatch(ctx context.Context, serial string) (bool, *RecalledBatch, error)
	SyncRecalled(ctx context.Context, items []RecalledBatch) error
}
