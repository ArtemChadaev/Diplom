package domain

import (
	"context"
	"time"
)

// UserRole — roles defined in the ERP spec.
type UserRole string

const (
	RoleAdmin            UserRole = "admin"
	RoleQP               UserRole = "qp"               // Authorised Person / QP
	RoleWarehouseManager UserRole = "warehouse_manager"
	RoleStorekeeper      UserRole = "storekeeper"
	RolePharmacist       UserRole = "pharmacist"
)

// User — clean domain model for a system user.
// No gorm/json/validate tags — only business fields.
type User struct {
	ID          int
	Email       string    // primary identity (unique, not null)
	GoogleID    *string   // nullable — set when linked via Google OAuth
	TelegramID  *int64    // nullable — set when linked via Telegram
	Role        UserRole
	NsPvAccess  bool      // access to narcotic/psychotropic substances (НС/ПВ)
	UkepBound   bool      // qualified electronic signature linked
	IsBlocked   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserProfile — flattened view joining users + employee_profiles.
type UserProfile struct {
	User
	EmployeeCode string
	FullName     string
	Position     string
	Department   string
}

// UserListFilter — filter for searching and listing users.
type UserListFilter struct {
	Query string
	Role  UserRole
	Page  int
	Limit int
}

// UserRepository — storage interface for users.
type UserRepository interface {
	// Identity resolution
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*User, error)
	FindByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)

	// Mutation
	Create(ctx context.Context, u *User) (*User, error)
	UpdateRole(ctx context.Context, userID int, role UserRole) error
	LinkGoogle(ctx context.Context, userID int, googleID string) error
	LinkTelegram(ctx context.Context, userID int, telegramID int64) error
	SetNsPvAccess(ctx context.Context, userID int, access bool) error
	SetBlocked(ctx context.Context, userID int, blocked bool) error

	// Profile
	FindProfileByUserID(ctx context.Context, userID int) (*UserProfile, error)
	List(ctx context.Context, filter UserListFilter) ([]*UserProfile, int, error)
}
