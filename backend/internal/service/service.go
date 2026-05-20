package service

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/mailer"
	"github.com/ima/diplom-backend/internal/repository"
)

// Service — агрегатор всех сервисов приложения
type Service struct {
	Auth            domain.AuthService
	Token           domain.TokenService
	EmployeeProfile domain.EmployeeProfileService
	Reference       domain.ReferenceService
	Product         domain.ProductService
	Supplier        domain.SupplierService
	Zone            domain.ZoneService
	Inbound         domain.InboundService
	EnvironmentLog  domain.EnvironmentLogService
	Order           domain.OrderService
	Inventory       domain.InventoryService
	Claim           domain.ClaimService
	Settings        domain.SystemSettingsService
	Batch           domain.BatchService
	RecalledBatch   domain.RecalledBatchService
}

func NewService(
	repos *repository.Repository,
	jwtSecret string,
	googleClientID string,
	otpHMACSecret string,
	m mailer.Mailer,
) *Service {
	// Standard TTL configuration: Access Token 15m, Refresh Token 15d
	tokenSvc := NewTokenService(jwtSecret, 15*time.Minute, 15*24*time.Hour)

	return &Service{
		Auth:            NewAuthService(repos.User, repos.Session, repos.OTP, tokenSvc, 15*24*time.Hour, googleClientID, m, otpHMACSecret),
		Token:           tokenSvc,
		EmployeeProfile: NewEmployeeProfileService(repos.EmployeeProfile),
		Reference:       NewReferenceService(repos.Reference),
		Product:         NewProductService(repos.Product, repos.Order, repos.Batch),
		Supplier:        NewSupplierService(repos.Supplier),
		Zone:            NewZoneService(repos.Zone),
		Inbound:         NewInboundService(repos.Inbound, repos.EmployeeProfile),
		EnvironmentLog:  NewEnvironmentLogService(repos.EnvironmentLog, repos.Zone, repos.EmployeeProfile),
		Order:           NewOrderService(repos.Order, repos.Batch, repos.Product, repos.Settings),
		Inventory:       NewInventoryService(repos.Inventory, repos.EmployeeProfile, repos.Product),
		Claim:           NewClaimService(repos.Claim, repos.Batch),
		Settings:        NewSystemSettingsService(repos.Settings),
		Batch:           NewBatchService(repos.Batch),
		RecalledBatch:   NewRecalledBatchService(repos.RecalledBatch),
	}
}

