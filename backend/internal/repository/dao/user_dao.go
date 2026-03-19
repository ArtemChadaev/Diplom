package dao

import "time"

// UserDAO — DAO-слепок таблицы users для GORM.
// Здесь живут gorm-теги; domain-модель тегов не знает.
type UserDAO struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Login        string    `gorm:"uniqueIndex;not null"`
	Email        *string   `gorm:"uniqueIndex"`
	GoogleID     *string   `gorm:"uniqueIndex"`
	TelegramID   *int64    `gorm:"uniqueIndex"`
	PasswordHash *string   `gorm:"column:password_hash"`
	Role         string    `gorm:"type:user_role;default:employee"`
	Status       string    `gorm:"type:varchar(50);default:unverified;not null"`
	IsBlocked    bool      `gorm:"column:is_blocked;default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// TableName явно задаёт имя таблицы, чтобы GORM не пытался угадать.
func (UserDAO) TableName() string { return "users" }
