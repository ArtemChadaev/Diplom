package domain

import "time"

// EmployeeProfile — чистая доменная модель профиля сотрудника.
// Связана с таблицей employee_profiles.
type EmployeeProfile struct {
	ID             uint
	UserID         uint
	EmployeeCode   string
	FullName       string
	CorporateEmail string
	Phone          string
	TelegramHandle string
	EmergencyContact string
	Position       string
	Department     string
	BirthDate      time.Time
	AvatarURL      string
	HireDate       time.Time
	DismissalDate  *time.Time
}
