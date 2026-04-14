package repository

import (
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

// Repository — агрегатор всех репозиториев приложения
type Repository struct {
	User            domain.UserRepository
	Session         domain.SessionRepository
	EmployeeProfile domain.EmployeeProfileRepository
	OTP             domain.OTPRepository
	Reference       domain.ReferenceRepository
	Product         domain.ProductRepository
	Supplier        domain.SupplierRepository
	Zone            domain.ZoneRepository
	Inbound         domain.InboundRepository
	EnvironmentLog  domain.EnvironmentLogRepository
	Order           domain.OrderRepository
	Inventory       domain.InventoryRepository
	Claim           domain.ClaimRepository
	Settings        domain.SystemSettingsRepository
	Batch           domain.BatchRepository
	RecalledBatch   domain.RecalledBatchRepository
}

// NewRepository создаёт новый слой репозиториев, используя GORM и Valkey
func NewRepository(db *gorm.DB, valkeyClient valkey.Client) *Repository {
	return &Repository{
		User:            NewUserRepository(db),
		Session:         NewSessionRepository(db),
		EmployeeProfile: NewEmployeeProfileRepository(db),
		OTP:             NewOTPValkeyRepository(valkeyClient),
		Reference:       NewReferenceRepository(db),
		Product:         NewProductRepository(db),
		Supplier:        NewSupplierRepository(db),
		Zone:            NewZoneRepository(db),
		Inbound:         NewInboundRepository(db),
		EnvironmentLog:  NewEnvironmentLogRepository(db),
		Order:           NewOrderRepository(db),
		Inventory:       NewInventoryRepository(db),
		Claim:           NewClaimRepository(db),
		Settings:        NewSystemSettingsRepository(db),
		Batch:           NewBatchRepository(db),
		RecalledBatch:   NewRecalledBatchRepository(db),
	}
}

