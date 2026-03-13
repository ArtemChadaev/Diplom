package dao

import "time"

// EmployeeProfileDAO — DAO для таблицы employee_profiles
type EmployeeProfileDAO struct {
	ID               uint       `gorm:"primaryKey;autoIncrement"`
	UserID           uint       `gorm:"uniqueIndex;not null"`
	EmployeeCode     string     `gorm:"uniqueIndex;not null;type:varchar(100)"`
	FullName         string     `gorm:"not null;type:varchar(255)"`
	CorporateEmail   string     `gorm:"uniqueIndex;type:varchar(255)"`
	Phone            string     `gorm:"type:varchar(20)"`
	TelegramHandle   string     `gorm:"type:varchar(100)"`
	EmergencyContact string     `gorm:"type:varchar(255)"`
	Position         string     `gorm:"type:varchar(255)"`
	Department       string     `gorm:"type:varchar(255)"`
	BirthDate        time.Time  `gorm:"type:date"`
	AvatarURL        string     `gorm:"type:text"`
	HireDate         time.Time  `gorm:"not null;type:date"`
	DismissalDate    *time.Time `gorm:"type:date"`
}

func (EmployeeProfileDAO) TableName() string { return "employee_profiles" }
