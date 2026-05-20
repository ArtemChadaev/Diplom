package domain

import (
	"context"
	"time"
)

type InventoryStatus string

const (
	InventoryStatusDraft     InventoryStatus = "draft"
	InventoryStatusActive    InventoryStatus = "active"
	InventoryStatusCompleted InventoryStatus = "completed"
)

// InventorySession — сессия инвентаризации (проверка остатков).
type InventorySession struct {
	ID          string          `json:"id"`
	ZoneID      string          `json:"zone_id"` // Session usually limited to one zone
	Status      InventoryStatus `json:"status"`
	StartedBy   int             `json:"started_by"`
	StartedAt   time.Time       `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Items       []InventoryItem `json:"items,omitempty"`
}

// InventoryItem — запись о пересчёте конкретного товара/партии.
type InventoryItem struct {
	ID                 string    `json:"id"`
	SessionID          string    `json:"session_id"`
	ProductID          string    `json:"product_id"`
	BatchNumber        string    `json:"batch_number"`
	SystemQuantity     int       `json:"system_qty"` // Qty in DB before session
	PhysicalQuantity   int       `json:"physical_qty"` // Qty counted by storekeeper
	DiscrepancyReason  string    `json:"reason"`
}

// InventoryRepository — интерфейс для работы с инвентаризацией.
type InventoryRepository interface {
	ListSessions(ctx context.Context, limit, offset int) ([]InventorySession, int, error)
	GetSessionByID(ctx context.Context, id string) (*InventorySession, error)
	CreateSession(ctx context.Context, s *InventorySession) error
	UpdateSession(ctx context.Context, s *InventorySession) error
	
	// Count items
	AddCount(ctx context.Context, item *InventoryItem) error
}

// NettingLine — строка зачёта излишков и недостач (пересортица).
type NettingLine struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	ATCGroup    string `json:"atc_group"` // АТХ группа 3-го уровня (первые 3 символа)
	Delta       int    `json:"delta"`     // Разница (physical_qty - system_qty)
}

// InventoryService — бизнес-логика инвентаризации.
type InventoryService interface {
	ListSessions(ctx context.Context, limit, offset int) ([]InventorySession, int, error)
	GetSession(ctx context.Context, id string) (*InventorySession, error)
	StartSession(ctx context.Context, userID int, zoneID string) (*InventorySession, error)
	FinishSession(ctx context.Context, id string) error
	SubmitCount(ctx context.Context, sessionID string, items []InventoryItem) error
	CalculateNetting(ctx context.Context, sessionID string) ([]NettingLine, error)
}
