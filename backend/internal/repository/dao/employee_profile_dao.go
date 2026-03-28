package dao

import (
	"encoding/json"
	"time"
)

// EmployeeProfileDAO — GORM snapshot of the employee_profiles table.
type EmployeeProfileDAO struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	UserID             uint       `gorm:"uniqueIndex;not null"`
	EmployeeCode       string     `gorm:"uniqueIndex;not null;type:varchar(100)"`
	FullName           string     `gorm:"not null;type:varchar(255)"`
	CorporateEmail     string     `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Phone              string     `gorm:"uniqueIndex;not null;type:varchar(20)"`
	Position           string     `gorm:"not null;type:varchar(255)"`
	Department         string     `gorm:"not null;type:varchar(255)"`
	BirthDate          time.Time  `gorm:"not null;type:date"`
	AvatarURL          string     `gorm:"type:text"`
	HireDate           time.Time  `gorm:"not null;type:date"`
	DismissalDate      *time.Time `gorm:"type:date"`
	// New ERP fields:
	MedicalBookScanURL string          `gorm:"column:medical_book_scan_url;type:text"`
	SpecialZoneAccess  bool            `gorm:"column:special_zone_access;default:false;not null"`
	GDPTrainingHistory json.RawMessage `gorm:"column:gdp_training_history;type:jsonb;default:'[]';not null"`
}

func (EmployeeProfileDAO) TableName() string { return "employee_profiles" }
