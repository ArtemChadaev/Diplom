package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type SupplierResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	INN         string    `json:"inn"`
	KPP         string    `json:"kpp"`
	ContactName string    `json:"contact_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToSupplierResponse(s domain.Supplier) SupplierResponse {
	return SupplierResponse{
		ID:          s.ID,
		Name:        s.Name,
		INN:         s.INN,
		KPP:         s.KPP,
		ContactName: s.ContactName,
		Phone:       s.Phone,
		Email:       s.Email,
		Address:     s.Address,
		IsActive:    s.IsActive,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

type CreateSupplierRequest struct {
	Name        string `json:"name" validate:"required"`
	INN         string `json:"inn" validate:"required,len=10|len=12"`
	KPP         string `json:"kpp" validate:"required"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	Address     string `json:"address"`
}

func (r CreateSupplierRequest) ToDomain() domain.Supplier {
	return domain.Supplier{
		Name:        r.Name,
		INN:         r.INN,
		KPP:         r.KPP,
		ContactName: r.ContactName,
		Phone:       r.Phone,
		Email:       r.Email,
		Address:     r.Address,
		IsActive:    true,
	}
}

type UpdateSupplierRequest struct {
	Name        *string `json:"name"`
	KPP         *string `json:"kpp"`
	ContactName *string `json:"contact_name"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
	IsActive    *bool   `json:"is_active"`
}

func (r UpdateSupplierRequest) ApplyTo(s *domain.Supplier) {
	if r.Name != nil { s.Name = *r.Name }
	if r.KPP != nil { s.KPP = *r.KPP }
	if r.ContactName != nil { s.ContactName = *r.ContactName }
	if r.Phone != nil { s.Phone = *r.Phone }
	if r.Email != nil { s.Email = *r.Email }
	if r.Address != nil { s.Address = *r.Address }
	if r.IsActive != nil { s.IsActive = *r.IsActive }
}

type SupplierListResponse struct {
	Total     int                `json:"total"`
	Suppliers []SupplierResponse `json:"suppliers"`
}
