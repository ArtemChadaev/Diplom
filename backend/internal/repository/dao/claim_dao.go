package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type ClaimDAO struct {
	ID          string             `gorm:"column:id;primaryKey"`
	Title       string             `gorm:"column:title"`
	Description string             `gorm:"column:description"`
	InboundID   *string            `gorm:"column:inbound_id;index"`
	OrderID     *string            `gorm:"column:order_id;index"`
	Status      domain.ClaimStatus `gorm:"column:status"`
	CreatedBy   int                `gorm:"column:created_by"`
	CreatedAt   time.Time          `gorm:"column:created_at"`
	UpdatedAt   time.Time          `gorm:"column:updated_at"`
}

func (ClaimDAO) TableName() string {
	return "claims"
}

func (c ClaimDAO) ToDomain() domain.Claim {
	return domain.Claim{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		InboundID:   c.InboundID,
		OrderID:     c.OrderID,
		Status:      c.Status,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func FromClaimDomain(c domain.Claim) ClaimDAO {
	return ClaimDAO{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		InboundID:   c.InboundID,
		OrderID:     c.OrderID,
		Status:      c.Status,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
