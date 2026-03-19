package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// SessionDAO — DAO-слепок таблицы refresh_tokens для GORM.
type SessionDAO struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    int            `gorm:"not null"`
	TokenHash string         `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `gorm:"not null"`
	UserAgent string         `gorm:"type:text"`
	IPAddress string         `gorm:"type:inet"`
	Metadata  datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	RevokedAt *time.Time
}

func (SessionDAO) TableName() string { return "refresh_tokens" }
