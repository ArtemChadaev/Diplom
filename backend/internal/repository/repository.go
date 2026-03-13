package repository

import (
	"github.com/ima/diplom-backend/internal/domain"
	"gorm.io/gorm"
)

// Repository — агрегатор всех репозиториев приложения
type Repository struct {
	User domain.UserRepository
}

// NewRepository создаёт новый слой репозиториев, используя GORM
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
