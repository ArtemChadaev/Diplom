package repository

import (
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// Repository — агрегатор всех репозиториев приложения
type Repository struct {
	User domain.UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
