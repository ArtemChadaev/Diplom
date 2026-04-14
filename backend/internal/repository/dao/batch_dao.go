package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type BatchDAO struct {
	ID               string             `gorm:"column:id;primaryKey"`
	ProductID        string             `gorm:"column:product_id;index"`
	ZoneID           *string            `gorm:"column:zone_id;index"`
	SerialNumber     string             `gorm:"column:serial_number"`
	ManufactureDate  time.Time          `gorm:"column:manufacture_date"`
	ExpiryDate       time.Time          `gorm:"column:expiry_date"`
	Quantity         int                `gorm:"column:quantity"`
	Status           domain.BatchStatus `gorm:"column:status"`
	UpdatedAt        time.Time          `gorm:"column:updated_at"`
}

func (BatchDAO) TableName() string {
	return "batches"
}

func (b BatchDAO) ToDomain() domain.Batch {
	return domain.Batch{
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

func FromBatchDomain(b domain.Batch) BatchDAO {
	return BatchDAO{
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
