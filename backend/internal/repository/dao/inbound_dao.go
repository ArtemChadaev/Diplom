package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type InboundReceivingDAO struct {
	ID            string               `gorm:"column:id;primaryKey"`
	InvoiceNumber string               `gorm:"column:invoice_number"`
	InvoiceDate   time.Time            `gorm:"column:invoice_date"`
	SupplierID    string               `gorm:"column:supplier_id"`
	Status        domain.InboundStatus `gorm:"column:status"`
	TotalAmount   float64              `gorm:"column:total_amount"`
	VATAmount     float64              `gorm:"column:vat_amount"`
	Notes         string               `gorm:"column:notes"`
	ReceivedBy    int                  `gorm:"column:received_by"`
	CreatedAt     time.Time            `gorm:"column:created_at"`
	UpdatedAt     time.Time            `gorm:"column:updated_at"`
	Items         []InboundItemDAO     `gorm:"foreignKey:InboundID"`
}

func (InboundReceivingDAO) TableName() string {
	return "inbound_receivings"
}

func (i InboundReceivingDAO) ToDomain() domain.InboundReceiving {
	items := make([]domain.InboundItem, len(i.Items))
	for idx, item := range i.Items {
		items[idx] = item.ToDomain()
	}

	return domain.InboundReceiving{
		ID:            i.ID,
		InvoiceNumber: i.InvoiceNumber,
		InvoiceDate:   i.InvoiceDate,
		SupplierID:    i.SupplierID,
		Status:        i.Status,
		TotalAmount:   i.TotalAmount,
		VATAmount:     i.VATAmount,
		Notes:         i.Notes,
		ReceivedBy:    i.ReceivedBy,
		Items:         items,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

type InboundItemDAO struct {
	ID                string    `gorm:"column:id;primaryKey"`
	InboundID         string    `gorm:"column:inbound_id;index"`
	ProductID         string    `gorm:"column:product_id"`
	BatchNumber       string    `gorm:"column:batch_number"`
	ExpirationDate    time.Time `gorm:"column:expiration_date"`
	Quantity          int       `gorm:"column:quantity"`
	PriceNetto        float64   `gorm:"column:price_netto"`
	VATRate           float64   `gorm:"column:vat_rate"`
	PriceBrutto       float64   `gorm:"column:price_brutto"`
	CertificateNumber string    `gorm:"column:cert_number"`
	ZoneID            string    `gorm:"column:zone_id"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (InboundItemDAO) TableName() string {
	return "inbound_items"
}

func (i InboundItemDAO) ToDomain() domain.InboundItem {
	return domain.InboundItem{
		ID:                i.ID,
		InboundID:         i.InboundID,
		ProductID:         i.ProductID,
		BatchNumber:       i.BatchNumber,
		ExpirationDate:    i.ExpirationDate,
		Quantity:          i.Quantity,
		PriceNetto:        i.PriceNetto,
		VATRate:           i.VATRate,
		PriceBrutto:       i.PriceBrutto,
		CertificateNumber: i.CertificateNumber,
		ZoneID:            i.ZoneID,
	}
}

func FromInboundDomain(i domain.InboundReceiving) InboundReceivingDAO {
	items := make([]InboundItemDAO, len(i.Items))
	for idx, item := range i.Items {
		items[idx] = FromInboundItemDomain(item)
	}

	return InboundReceivingDAO{
		ID:            i.ID,
		InvoiceNumber: i.InvoiceNumber,
		InvoiceDate:   i.InvoiceDate,
		SupplierID:    i.SupplierID,
		Status:        i.Status,
		TotalAmount:   i.TotalAmount,
		VATAmount:     i.VATAmount,
		Notes:         i.Notes,
		ReceivedBy:    i.ReceivedBy,
		Items:         items,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

func FromInboundItemDomain(i domain.InboundItem) InboundItemDAO {
	return InboundItemDAO{
		ID:                i.ID,
		InboundID:         i.InboundID,
		ProductID:         i.ProductID,
		BatchNumber:       i.BatchNumber,
		ExpirationDate:    i.ExpirationDate,
		Quantity:          i.Quantity,
		PriceNetto:        i.PriceNetto,
		VATRate:           i.VATRate,
		PriceBrutto:       i.PriceBrutto,
		CertificateNumber: i.CertificateNumber,
		ZoneID:            i.ZoneID,
	}
}
