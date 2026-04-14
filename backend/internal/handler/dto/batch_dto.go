package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type BatchResponse struct {
	ID               string             `json:"id"`
	ProductID        string             `json:"product_id"`
	ZoneID           *string            `json:"zone_id,omitempty"`
	SerialNumber     string             `json:"serial_number"`
	ManufactureDate  time.Time          `json:"manufacture_date"`
	ExpiryDate       time.Time          `json:"expiry_date"`
	Quantity         int                `json:"quantity"`
	Status           domain.BatchStatus `json:"status"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func ToBatchResponse(b domain.Batch) BatchResponse {
	return BatchResponse{
		ID:               b.ID,
		ProductID:        b.ProductID,
		ZoneID:           b.ZoneID,
		SerialNumber:     b.SerialNumber,
		ManufactureDate:  b.ManufactureDate,
		ExpiryDate:       b.ExpiryDate,
		Quantity:         b.Quantity,
		Status:           b.Status,
		UpdatedAt:        b.UpdatedAt,
	}
}

type BatchListResponse struct {
	Total   int             `json:"total"`
	Batches []BatchResponse `json:"batches"`
}

type UpdateBatchStatusRequest struct {
	Status domain.BatchStatus `json:"status" validate:"required"`
}

type TransferBatchRequest struct {
	TargetZoneID string `json:"target_zone_id" validate:"required,uuid"`
}
