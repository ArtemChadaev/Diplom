package domain

import (
	"context"
	"time"
)

// EmployeeProfile — чистая доменная модель профиля сотрудника.
// Связана с таблицей employee_profiles.
type EmployeeProfile struct {
	ID               uint
	UserID           uint
	EmployeeCode     string
	FullName         string
	CorporateEmail   string
	Phone            string
	TelegramHandle   string
	EmergencyContact string
	Position         string
	Department       string
	BirthDate        time.Time
	AvatarURL        string
	HireDate         time.Time
	DismissalDate    *time.Time
}

// UpdateEmployeeProfileInput carries only the fields the admin wants to change.
// Pointer semantics: nil means "don't touch this field".
type UpdateEmployeeProfileInput struct {
	FullName         *string    `json:"full_name"`
	CorporateEmail   *string    `json:"corporate_email"`
	Phone            *string    `json:"phone"`
	TelegramHandle   *string    `json:"telegram_handle"`
	EmergencyContact *string    `json:"emergency_contact"`
	Position         *string    `json:"position"`
	Department       *string    `json:"department"`
	BirthDate        *time.Time `json:"birth_date"`
	AvatarURL        *string    `json:"avatar_url"`
	HireDate         *time.Time `json:"hire_date"`
	DismissalDate    *time.Time `json:"dismissal_date"`
}

// EmployeeProfileRepository — persistence interface for employee profiles.
type EmployeeProfileRepository interface {
	FindByUserID(ctx context.Context, userID int) (*EmployeeProfile, error)
	FindByID(ctx context.Context, id int) (*EmployeeProfile, error)
	Update(ctx context.Context, id int, input UpdateEmployeeProfileInput) (*EmployeeProfile, error)
	List(ctx context.Context, limit, offset int) ([]EmployeeProfile, error)
}

// EmployeeProfileService — business logic interface.
type EmployeeProfileService interface {
	GetProfile(ctx context.Context, callerID int, callerRole UserRole, targetUserID int) (*EmployeeProfile, error)
	UpdateProfile(ctx context.Context, callerID int, callerRole UserRole, targetUserID int, input UpdateEmployeeProfileInput) (*EmployeeProfile, error)
	ListProfiles(ctx context.Context, callerID int, callerRole UserRole, limit, offset int) ([]EmployeeProfile, error)
}

