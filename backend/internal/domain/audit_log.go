package domain

import "time"

// AuditLog — clean domain model for a single audit log entry.
// Reflects a recorded change in the system (create / update / delete).
// No gorm/json tags — those live in dao.AuditLogDAO.
type AuditLog struct {
	ID        string    // UUID primary key
	UserID    *int      // nil when the action originates from a system/background process
	Action    string    // "create" | "update" | "delete"
	Entity    string    // table name, e.g. "users", "batches"
	EntityID  string    // primary key value of the affected row
	OldValues []byte    // JSON snapshot before change; nil on Create
	NewValues []byte    // JSON snapshot after change; nil on Delete
	IPAddress string    // originating IP; empty when not available
	CreatedAt time.Time
}
