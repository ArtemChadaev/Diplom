package domain

import (
	"context"
	"time"
)

// User — доменная сущность пользователя
type User struct {
	ID           int       `db:"id"`
	FullName     string    `db:"full_name"`     // ФИО
	PasswordHash string    `db:"password_hash"` // Хеш пароля
	Email        string    `db:"email"`         // E-mail (уникальный)
	Position     string    `db:"position"`      // Должность
	Role         string    `db:"role"`          // Полномочия / роль
	CreatedAt    time.Time `db:"created_at"`    // Дата создания
	IsActive     bool      `db:"is_active"`     // true — работает, false — уволен
}

// UserRepository — интерфейс для работы с хранилищем пользователей
type UserRepository interface {
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

// UserService — интерфейс бизнес-логики пользователей
type UserService interface {
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}
