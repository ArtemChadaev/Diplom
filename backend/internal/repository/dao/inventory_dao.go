package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type InventorySessionDAO struct {
	ID          string                 `gorm:"column:id;primaryKey"`
	ZoneID      *string                `gorm:"column:zone_id;index"`
	Status      domain.InventoryStatus `gorm:"column:status"`
	StartedBy   int                    `gorm:"column:started_by"`
	StartedAt   time.Time              `gorm:"column:started_at"`
	CompletedAt *time.Time             `gorm:"column:completed_at"`
	Items       []InventoryItemDAO     `gorm:"foreignKey:SessionID"`
}

func (InventorySessionDAO) TableName() string {
	return "inventory_sessions"
}

func (s InventorySessionDAO) ToDomain() domain.InventorySession {
	items := make([]domain.InventoryItem, len(s.Items))
	for i, item := range s.Items {
		items[i] = item.ToDomain()
	}

	var zoneID string
	if s.ZoneID != nil {
		zoneID = *s.ZoneID
	}

	return domain.InventorySession{
		ID:          s.ID,
		ZoneID:      zoneID,
		Status:      s.Status,
		StartedBy:   s.StartedBy,
		StartedAt:   s.StartedAt,
		CompletedAt: s.CompletedAt,
		Items:       items,
	}
}

type InventoryItemDAO struct {
	ID               string    `gorm:"column:id;primaryKey"`
	SessionID        string    `gorm:"column:session_id;index"`
	ProductID        string    `gorm:"column:product_id"`
	BatchNumber      string    `gorm:"column:batch_number"`
	SystemQuantity   int       `gorm:"column:system_qty"`
	PhysicalQuantity int       `gorm:"column:physical_qty"`
	DiscrepancyReason string    `gorm:"column:reason"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (InventoryItemDAO) TableName() string {
	return "inventory_items"
}

func (i InventoryItemDAO) ToDomain() domain.InventoryItem {
	return domain.InventoryItem{
		ID:               i.ID,
		SessionID:        i.SessionID,
		ProductID:        i.ProductID,
		BatchNumber:      i.BatchNumber,
		SystemQuantity:   i.SystemQuantity,
		PhysicalQuantity: i.PhysicalQuantity,
		DiscrepancyReason: i.DiscrepancyReason,
	}
}

func FromInventorySessionDomain(s domain.InventorySession) InventorySessionDAO {
	items := make([]InventoryItemDAO, len(s.Items))
	for i, item := range s.Items {
		items[i] = FromInventoryItemDomain(item)
	}

	var zoneID *string
	if s.ZoneID != "" {
		zoneID = &s.ZoneID
	}

	return InventorySessionDAO{
		ID:          s.ID,
		ZoneID:      zoneID,
		Status:      s.Status,
		StartedBy:   s.StartedBy,
		StartedAt:   s.StartedAt,
		CompletedAt: s.CompletedAt,
		Items:       items,
	}
}

func FromInventoryItemDomain(i domain.InventoryItem) InventoryItemDAO {
	return InventoryItemDAO{
		ID:               i.ID,
		SessionID:        i.SessionID,
		ProductID:        i.ProductID,
		BatchNumber:      i.BatchNumber,
		SystemQuantity:   i.SystemQuantity,
		PhysicalQuantity: i.PhysicalQuantity,
		DiscrepancyReason: i.DiscrepancyReason,
	}
}
