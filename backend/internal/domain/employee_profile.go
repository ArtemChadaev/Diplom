package domain

import (
	"context"
	"encoding/json"
	"time"
)

// GDPTrainingRecord — one entry in the GDP training history JSONB array.
type GDPTrainingRecord struct {
	Date           string `json:"date"`            // "YYYY-MM-DD"
	CourseName     string `json:"course_name"`
	Result         string `json:"result"`          // "pass" | "fail"
	CertificateURL string `json:"certificate_url"`
}

// EmployeeProfile — clean domain model for an employee profile.
// Corresponds to the employee_profiles table.
type EmployeeProfile struct {
	ID                  uint
	UserID              uint
	EmployeeCode        string
	FullName            string
	CorporateEmail      string
	Phone               string
	Position            string
	Department          string
	BirthDate           time.Time
	AvatarURL           string
	HireDate            time.Time
	DismissalDate       *time.Time
	// New ERP fields:
	MedicalBookScanURL  string
	SpecialZoneAccess   bool
	GDPTrainingHistory  []GDPTrainingRecord
}

// UpdateEmployeeProfileInput carries only the fields the admin wants to change.
// Pointer semantics: nil means "don't touch this field".
type UpdateEmployeeProfileInput struct {
	FullName            *string             `json:"full_name"`
	CorporateEmail      *string             `json:"corporate_email"`
	Phone               *string             `json:"phone"`
	Position            *string             `json:"position"`
	Department          *string             `json:"department"`
	BirthDate           *time.Time          `json:"birth_date"`
	AvatarURL           *string             `json:"avatar_url"`
	HireDate            *time.Time          `json:"hire_date"`
	DismissalDate       *time.Time          `json:"dismissal_date"`
	MedicalBookScanURL  *string             `json:"medical_book_scan_url"`
	SpecialZoneAccess   *bool               `json:"special_zone_access"`
	GDPTrainingHistory  json.RawMessage     `json:"gdp_training_history"`
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
