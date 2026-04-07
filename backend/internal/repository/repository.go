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
}

// NewRepository создаёт новый слой репозиториев, используя GORM и Valkey
func NewRepository(db *gorm.DB, valkeyClient valkey.Client) *Repository {
	return &Repository{
		User:            NewUserRepository(db),
		Session:         NewSessionRepository(db),
		EmployeeProfile: NewEmployeeProfileRepository(db),
		OTP:             NewOTPValkeyRepository(valkeyClient),
	}
}

