package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type InventorySessionResponse struct {
	ID          string                 `json:"id"`
	ZoneID      string                 `json:"zone_id"`
	Status      domain.InventoryStatus `json:"status"`
	StartedBy   int                    `json:"started_by"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Items       []InventoryItemResponse `json:"items,omitempty"`
}

type InventoryItemResponse struct {
	ID                string `json:"id"`
	ProductID         string `json:"product_id"`
	BatchNumber       string `json:"batch_number"`
	SystemQuantity    int    `json:"system_qty"`
	PhysicalQuantity  int    `json:"physical_qty"`
	DiscrepancyReason string `json:"reason"`
}

func ToInventorySessionResponse(s domain.InventorySession) InventorySessionResponse {
	items := make([]InventoryItemResponse, len(s.Items))
	for i, item := range s.Items {
		items[i] = InventoryItemResponse{
			ID:                item.ID,
			ProductID:         item.ProductID,
			BatchNumber:       item.BatchNumber,
			SystemQuantity:    item.SystemQuantity,
			PhysicalQuantity:  item.PhysicalQuantity,
			DiscrepancyReason: item.DiscrepancyReason,
		}
	}

	return InventorySessionResponse{
		ID:          s.ID,
		ZoneID:      s.ZoneID,
		Status:      s.Status,
		StartedBy:   s.StartedBy,
		StartedAt:   s.StartedAt,
		CompletedAt: s.CompletedAt,
		Items:       items,
	}
}

type StartInventoryRequest struct {
	ZoneID string `json:"zone_id" validate:"required,uuid"`
}

type SubmitCountRequest struct {
	Items []CountItem `json:"items" validate:"required,dive"`
}

type CountItem struct {
	ProductID         string `json:"product_id" validate:"required,uuid"`
	BatchNumber       string `json:"batch_number" validate:"required"`
	SystemQuantity    int    `json:"system_qty"`
	PhysicalQuantity  int    `json:"physical_qty"`
	DiscrepancyReason string `json:"reason"`
}

func (r CountItem) ToDomain() domain.InventoryItem {
	return domain.InventoryItem{
		ProductID:        r.ProductID,
		BatchNumber:      r.BatchNumber,
		SystemQuantity:   r.SystemQuantity,
		PhysicalQuantity: r.PhysicalQuantity,
		DiscrepancyReason: r.DiscrepancyReason,
	}
}
