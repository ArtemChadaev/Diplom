package dao

import "time"

// UserDAO — GORM snapshot of the users table.
// gorm tags live here; domain model has no ORM dependency.
type UserDAO struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Email       string    `gorm:"uniqueIndex;not null;type:varchar(255)"`
	GoogleID    *string   `gorm:"uniqueIndex;type:varchar(255)"`
	TelegramID  *int64    `gorm:"uniqueIndex"`
	Role        string    `gorm:"type:user_role;default:pharmacist;not null"`
	NsPvAccess  bool      `gorm:"column:ns_pv_access;default:false;not null"`
	UkepBound   bool      `gorm:"column:ukep_bound;default:false;not null"`
	IsBlocked   bool      `gorm:"column:is_blocked;default:false;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName explicitly sets the table name so GORM does not try to guess.
func (UserDAO) TableName() string { return "users" }
