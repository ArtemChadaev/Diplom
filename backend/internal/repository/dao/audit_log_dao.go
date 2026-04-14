package dao

import (
	"encoding/json"
	"time"
)

// AuditLogDAO — GORM snapshot of the audit_logs table.
// gorm tags live here; the domain model (domain.AuditLog) has no ORM dependency.
type AuditLogDAO struct {
	ID        string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    *int            `gorm:"column:user_id"`
	Action    string          `gorm:"column:action;type:varchar(255);not null"`
	Entity    string          `gorm:"column:entity;type:varchar(100);not null"`
	EntityID  string          `gorm:"column:entity_id;type:varchar(100)"`
	OldValues json.RawMessage `gorm:"column:old_values;type:jsonb"`
	NewValues json.RawMessage `gorm:"column:new_values;type:jsonb"`
	IPAddress *string         `gorm:"column:ip_address;type:inet"`
	// Immutability chain fields (Phase 3, currently NULL / mock):
	PrevHash string `gorm:"column:prev_hash;type:text"`
	LogHash  string `gorm:"column:log_hash;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName explicitly maps this DAO to the audit_logs table.
func (AuditLogDAO) TableName() string { return "audit_logs" }
