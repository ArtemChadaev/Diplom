package dto

import "time"

// EmployeeProfileDTO — DTO для профиля сотрудника
type EmployeeProfileDTO struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	EmployeeCode     string     `json:"employee_code"`
	FullName         string     `json:"full_name"`
	CorporateEmail   string     `json:"corporate_email"`
	Phone            string     `json:"phone"`
	TelegramHandle   string     `json:"telegram_handle"`
	EmergencyContact string     `json:"emergency_contact"`
	Position         string     `json:"position"`
	Department       string     `json:"department"`
	BirthDate        time.Time  `json:"birth_date"`
	AvatarURL        string     `json:"avatar_url"`
	HireDate         time.Time  `json:"hire_date"`
	DismissalDate    *time.Time `json:"dismissal_date"`
}
