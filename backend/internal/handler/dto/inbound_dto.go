package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type InboundResponse struct {
	ID            string               `json:"id"`
	InvoiceNumber string               `json:"invoice_number"`
	InvoiceDate   time.Time            `json:"invoice_date"`
	SupplierID    string               `json:"supplier_id"`
	Status        domain.InboundStatus `json:"status"`
	TotalAmount   float64              `json:"total_amount"`
	VATAmount     float64              `json:"vat_amount"`
	Notes         string               `json:"notes"`
	ReceivedBy    int                  `json:"received_by"`
	Items         []InboundItemResponse `json:"items,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

type InboundItemResponse struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	BatchNumber       string    `json:"batch_number"`
	ExpirationDate    time.Time `json:"expiration_date"`
	Quantity          int       `json:"quantity"`
	PriceNetto        float64   `json:"price_netto"`
	VATRate           float64   `json:"vat_rate"`
	PriceBrutto       float64   `json:"price_brutto"`
	CertificateNumber string    `json:"cert_number"`
	ZoneID            string    `json:"zone_id"`
}

func ToInboundResponse(i domain.InboundReceiving) InboundResponse {
	items := make([]InboundItemResponse, len(i.Items))
	for idx, item := range i.Items {
		items[idx] = InboundItemResponse{
			ID:                item.ID,
			ProductID:         item.ProductID,
			BatchNumber:       item.BatchNumber,
			ExpirationDate:    item.ExpirationDate,
			Quantity:          item.Quantity,
			PriceNetto:        item.PriceNetto,
			VATRate:           item.VATRate,
			PriceBrutto:       item.PriceBrutto,
			CertificateNumber: item.CertificateNumber,
			ZoneID:            item.ZoneID,
		}
	}

	return InboundResponse{
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

type CreateInboundRequest struct {
	InvoiceNumber string               `json:"invoice_number" validate:"required"`
	InvoiceDate   time.Time            `json:"invoice_date" validate:"required"`
	SupplierID    string               `json:"supplier_id" validate:"required,uuid"`
	Notes         string               `json:"notes"`
	Items         []CreateInboundItem  `json:"items" validate:"required,dive"`
}

type CreateInboundItem struct {
	ProductID         string    `json:"product_id" validate:"required,uuid"`
	BatchNumber       string    `json:"batch_number" validate:"required"`
	ExpirationDate    time.Time `json:"expiration_date" validate:"required"`
	Quantity          int       `json:"quantity" validate:"required,min=1"`
	PriceNetto        float64   `json:"price_netto" validate:"required"`
	VATRate           float64   `json:"vat_rate"`
	PriceBrutto       float64   `json:"price_brutto"`
	CertificateNumber string    `json:"cert_number"`
	ZoneID            string    `json:"zone_id" validate:"required,uuid"`
}

func (r CreateInboundRequest) ToDomain() domain.InboundReceiving {
	items := make([]domain.InboundItem, len(r.Items))
	var total, vatTotal float64
	for i, item := range r.Items {
		items[i] = domain.InboundItem{
			ProductID:         item.ProductID,
			BatchNumber:       item.BatchNumber,
			ExpirationDate:    item.ExpirationDate,
			Quantity:          item.Quantity,
			PriceNetto:        item.PriceNetto,
			VATRate:           item.VATRate,
			PriceBrutto:       item.PriceBrutto,
			CertificateNumber: item.CertificateNumber,
			ZoneID:            item.ZoneID,
		}
		total += item.PriceBrutto * float64(item.Quantity)
		vatTotal += (item.PriceBrutto - item.PriceNetto) * float64(item.Quantity)
	}

	return domain.InboundReceiving{
		InvoiceNumber: r.InvoiceNumber,
		InvoiceDate:   r.InvoiceDate,
		SupplierID:    r.SupplierID,
		Notes:         r.Notes,
		Items:         items,
		TotalAmount:   total,
		VATAmount:     vatTotal,
	}
}

type InboundListResponse struct {
	Total    int               `json:"total"`
	Inbounds []InboundResponse `json:"inbounds"`
}

type UpdateInboundStatusRequest struct {
	Status domain.InboundStatus `json:"status" validate:"required,oneof=draft received completed cancelled"`
}
