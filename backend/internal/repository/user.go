package repository

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// IsEmailTaken — проверяет, зарегистрирован ли уже такой e-mail
func (r *userRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
