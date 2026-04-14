package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type RecalledBatchResponse struct {
	ID           string    `json:"id"`
	SerialNumber string    `json:"serial_number"`
	ProductName  string    `json:"product_name"`
	RuNumber     string    `json:"ru_number"`
	RecallReason string    `json:"recall_reason"`
	IssuedAt     time.Time `json:"issued_at"`
	SyncedAt     time.Time `json:"synced_at"`
}

func ToRecalledBatchResponse(b domain.RecalledBatch) RecalledBatchResponse {
	return RecalledBatchResponse{
		ID:           b.ID,
		SerialNumber: b.SerialNumber,
		ProductName:  b.ProductName,
		RuNumber:     b.RuNumber,
		RecallReason: b.RecallReason,
		IssuedAt:     b.IssuedAt,
		SyncedAt:     b.SyncedAt,
	}
}

type RecalledCheckResponse struct {
	IsRecalled bool                   `json:"is_recalled"`
	Details    *RecalledBatchResponse `json:"details,omitempty"`
}
