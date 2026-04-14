package domain

import (
	"context"
	"time"
)

// Supplier — поставщик продукции.
type Supplier struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`        // наименование организации
	INN         string    `json:"inn"`         // ИНН (уникальный)
	KPP         string    `json:"kpp"`         // КПП
	ContactName string    `json:"contact_name"` // ФИО контактного лица
	Phone       string    `json:"phone"`       // телефон
	Email       string    `json:"email"`       // email
	Address     string    `json:"address"`     // юридический адрес
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SupplierRepository — интерфейс для работы с поставщиками.
type SupplierRepository interface {
	List(ctx context.Context, limit, offset int) ([]Supplier, int, error)
	GetByID(ctx context.Context, id string) (*Supplier, error)
	GetByINN(ctx context.Context, inn string) (*Supplier, error)
	Create(ctx context.Context, s *Supplier) error
	Update(ctx context.Context, s *Supplier) error
	Delete(ctx context.Context, id string) error
}

// SupplierService — бизнес-логика поставщиков.
type SupplierService interface {
	ListSuppliers(ctx context.Context, limit, offset int) ([]Supplier, int, error)
	GetSupplier(ctx context.Context, id string) (*Supplier, error)
	CreateSupplier(ctx context.Context, callerRole UserRole, s *Supplier) (*Supplier, error)
	UpdateSupplier(ctx context.Context, callerRole UserRole, s *Supplier) (*Supplier, error)
	DeleteSupplier(ctx context.Context, callerRole UserRole, id string) error
}
