package repository

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// IsLoginTaken — проверяет, занят ли логин.
func (r *userRepository) IsLoginTaken(ctx context.Context, login string) (bool, error) {
	var count int64

	// Пример 1: Чистый SQL через GORM (db.Raw)
	/*
		err := r.db.WithContext(ctx).
			Raw(`SELECT count(1) FROM users WHERE login = ?`, login).
			Scan(&count).Error
	*/

	// Пример 2: Fluent GORM API (используется)
	err := r.db.WithContext(ctx).
		Model(&dao.UserDAO{}).
		Where("login = ?", login).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
