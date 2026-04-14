package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type RecalledBatchDAO struct {
	ID           string    `gorm:"column:id;primaryKey"`
	SerialNumber string    `gorm:"column:serial_number;index"`
	ProductName  string    `gorm:"column:product_name"`
	RuNumber     string    `gorm:"column:ru_number"`
	RecallReason string    `gorm:"column:recall_reason"`
	IssuedAt     time.Time `gorm:"column:issued_at"`
	SyncedAt     time.Time `gorm:"column:synced_at"`
}

func (RecalledBatchDAO) TableName() string {
	return "recalled_batches"
}

func (b RecalledBatchDAO) ToDomain() domain.RecalledBatch {
	return domain.RecalledBatch{
		ID:           b.ID,
		SerialNumber: b.SerialNumber,
		ProductName:  b.ProductName,
		RuNumber:     b.RuNumber,
		RecallReason: b.RecallReason,
		IssuedAt:     b.IssuedAt,
		SyncedAt:     b.SyncedAt,
	}
}

func FromRecalledBatchDomain(b domain.RecalledBatch) RecalledBatchDAO {
	return RecalledBatchDAO{
		ID:           b.ID,
		SerialNumber: b.SerialNumber,
		ProductName:  b.ProductName,
		RuNumber:     b.RuNumber,
		RecallReason: b.RecallReason,
		IssuedAt:     b.IssuedAt,
		SyncedAt:     b.SyncedAt,
	}
}
