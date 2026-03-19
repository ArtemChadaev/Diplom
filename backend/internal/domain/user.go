package domain

import (
	"context"
	"time"
)

type UserRole string
type UserStatus string

const (
	RoleAdmin      UserRole = "admin"
	RoleEmployee   UserRole = "employee"
	RoleUnverified UserRole = "unverified"

	StatusUnverified UserStatus = "unverified"
	StatusActive     UserStatus = "active"
	StatusBlocked    UserStatus = "blocked"
)

// User — чистая доменная модель пользователя.
// Не содержит тегов gorm, json или validate — только бизнес-поля.
type User struct {
	ID           int
	Login        string
	Email        *string // nullable
	GoogleID     *string // nullable
	TelegramID   *int64  // nullable
	PasswordHash *string // nullable (social-only users have no password)
	Role         UserRole
	Status       UserStatus
	IsBlocked    bool
	CreatedAt    time.Time
}

// UserProfile — flattened view joining users + employee_profiles
type UserProfile struct {
	User
	EmployeeCode string
	FullName     string
	Position     string
	Department   string
}

// UserRepository — интерфейс для работы с хранилищем пользователей.
type UserRepository interface {
	// Identity resolution
	FindByID(ctx context.Context, id int) (*User, error)
	FindByLogin(ctx context.Context, login string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*User, error)
	FindByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	IsLoginTaken(ctx context.Context, login string) (bool, error)

	// Mutation
	Create(ctx context.Context, u *User) (*User, error)
	UpdateRole(ctx context.Context, userID int, role UserRole) error
	UpdateStatus(ctx context.Context, userID int, status UserStatus) error
	LinkGoogle(ctx context.Context, userID int, googleID string) error
	LinkTelegram(ctx context.Context, userID int, telegramID int64) error
	SetPasswordHash(ctx context.Context, userID int, hash string) error

	// Profile
	FindProfileByUserID(ctx context.Context, userID int) (*UserProfile, error)
}
