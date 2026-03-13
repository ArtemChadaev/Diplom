package dao

import "time"

// UserDAO — DAO-слепок таблицы users для GORM.
// Здесь живут gorm-теги; domain-модель тегов не знает.
type UserDAO struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Login        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Role         string    `gorm:"type:user_role;default:employee"`
	IsBlocked    bool      `gorm:"column:is_blocked;default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// TableName явно задаёт имя таблицы, чтобы GORM не пытался угадать.
func (UserDAO) TableName() string { return "users" }
