package dto

import (
	"encoding/json"
	"time"
)

// EmployeeProfileDTO — DTO for employee profile responses
type EmployeeProfileDTO struct {
	ID                 uint            `json:"id"`
	UserID             uint            `json:"user_id"`
	EmployeeCode       string          `json:"employee_code"`
	FullName           string          `json:"full_name"`
	CorporateEmail     string          `json:"corporate_email"`
	Phone              string          `json:"phone"`
	Position           string          `json:"position"`
	Department         string          `json:"department"`
	BirthDate          time.Time       `json:"birth_date"`
	AvatarURL          string          `json:"avatar_url"`
	HireDate           time.Time       `json:"hire_date"`
	DismissalDate      *time.Time      `json:"dismissal_date,omitempty"`
	MedicalBookScanURL string          `json:"medical_book_scan_url"`
	SpecialZoneAccess  bool            `json:"special_zone_access"`
	GDPTrainingHistory json.RawMessage `json:"gdp_training_history"`
}
