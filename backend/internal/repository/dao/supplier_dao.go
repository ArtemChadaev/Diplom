package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type SupplierDAO struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name"`
	INN         string    `gorm:"column:inn;uniqueIndex"`
	KPP         string    `gorm:"column:kpp"`
	ContactName string    `gorm:"column:contact_name"`
	Phone       string    `gorm:"column:phone"`
	Email       string    `gorm:"column:email"`
	Address     string    `gorm:"column:address"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (SupplierDAO) TableName() string {
	return "suppliers"
}

func (s SupplierDAO) ToDomain() domain.Supplier {
	return domain.Supplier{
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

func FromSupplierDomain(s domain.Supplier) SupplierDAO {
	return SupplierDAO{
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
