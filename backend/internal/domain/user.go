package domain

import (
	"context"
	"time"
)

// User — чистая доменная модель пользователя.
// Не содержит тегов gorm, json или validate — только бизнес-поля.
type User struct {
	ID           string
	Login        string
	PasswordHash string
	Role         string
	IsBlocked    bool
	CreatedAt    time.Time
}

// UserRepository — интерфейс для работы с хранилищем пользователей.
type UserRepository interface {
	IsLoginTaken(ctx context.Context, login string) (bool, error)
}

// UserService — интерфейс бизнес-логики пользователей.
type UserService interface {
	IsLoginTaken(ctx context.Context, login string) (bool, error)
}
