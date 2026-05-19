package dto

import (
	"time"
)

// AssignRoleRequest — body for PATCH /admin/users/{id}/role
type AssignRoleRequest struct {
	Role string `json:"role"`
}

// SetBlockedRequest — body for PATCH /admin/users/{id}/blocked
type SetBlockedRequest struct {
	Blocked bool `json:"blocked"`
}

// UserResponse — public user fields returned by the API
type UserResponse struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	NsPvAccess  bool   `json:"ns_pv_access"`
	UkepBound   bool   `json:"ukep_bound"`
	IsBlocked   bool   `json:"is_blocked"`
}

type UserProfileResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	NsPvAccess   bool   `json:"ns_pv_access"`
	UkepBound    bool   `json:"ukep_bound"`
	IsBlocked    bool   `json:"is_blocked"`
	EmployeeCode string `json:"employee_code,omitempty"`
	FullName     string `json:"full_name,omitempty"`
	Position     string `json:"position,omitempty"`
	Department   string `json:"department,omitempty"`
}

type UserListQuery struct {
	Q     string `schema:"q"`
	Role  string `schema:"role"`
	Page  int    `schema:"page"`
	Limit int    `schema:"limit"`
}

type PatchUserRequest struct {
	Role              *string `json:"role,omitempty"`
	NsPvAccess        *bool   `json:"ns_pv_access,omitempty"`
	SpecialZoneAccess *bool   `json:"special_zone_access,omitempty"`
}

type PatchMeRequest struct {
	FullName            *string         `json:"full_name"`
	Phone               *string         `json:"phone"`
	CorporateEmail      *string         `json:"corporate_email"`
	CorporateMail       *string         `json:"corporatemail"`
	BirthDate           *time.Time      `json:"birth_date"`
	BirthdayDate        *time.Time      `json:"birthday_date"`
	AvatarURL           *string         `json:"avatar_url"`
	AvatarURL2          *string         `json:"avatarurl"`
	MedicalBookScanURL  *string         `json:"medical_book_scan_url"`
	GDPTrainingHistory  interface{} `json:"gdp_training_history"`
}
