// Package hooks provides a GORM plugin that automatically writes an audit-log
// row to the audit_logs table for every CREATE, UPDATE and DELETE operation.
//
// Usage (register once after gorm.Open):
//
//	if err := db.Use(&hooks.AuditPlugin{}); err != nil { ... }
package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ima/diplom-backend/internal/pkg/logger"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

// skipTables contains table names that must never be audited.
var skipTables = map[string]bool{
	"audit_logs":       true, // prevent infinite recursion
	"refresh_tokens":   true, // high-frequency token churn, no business audit value
	"environment_logs": true, // sensor data stream
}

// settingOldValues is the Statement.Settings key used to pass the pre-update
// JSON snapshot from beforeUpdate to afterUpdate.
const settingOldValues = "audit:old_values"

// ---------------------------------------------------------------------------
// AuditPlugin
// ---------------------------------------------------------------------------

// AuditPlugin is a GORM plugin. Register it once with db.Use(&AuditPlugin{}).
type AuditPlugin struct{}

// Name implements the gorm.Plugin interface.
func (p *AuditPlugin) Name() string { return "audit_plugin" }

// Initialize registers all audit callbacks on db.
func (p *AuditPlugin) Initialize(db *gorm.DB) error {
	cbs := []struct {
		scope  string
		pos    string // "before:<name>" or "after:<name>"
		name   string
		fn     func(*gorm.DB)
	}{
		{"update", "before", "audit:before_update", beforeUpdate},
		{"create", "after", "audit:after_create", afterCreate},
		{"update", "after", "audit:after_update", afterUpdate},
		{"delete", "after", "audit:after_delete", afterDelete},
	}

	for _, cb := range cbs {
		var err error
		switch cb.scope {
		case "create":
			switch cb.pos {
			case "after":
				err = db.Callback().Create().After("gorm:create").Register(cb.name, cb.fn)
			}
		case "update":
			switch cb.pos {
			case "before":
				err = db.Callback().Update().Before("gorm:update").Register(cb.name, cb.fn)
			case "after":
				err = db.Callback().Update().After("gorm:update").Register(cb.name, cb.fn)
			}
		case "delete":
			switch cb.pos {
			case "after":
				err = db.Callback().Delete().After("gorm:delete").Register(cb.name, cb.fn)
			}
		}
		if err != nil {
			return fmt.Errorf("audit_plugin: register %s: %w", cb.name, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Callbacks
// ---------------------------------------------------------------------------

// beforeUpdate captures the current DB row as old_values and stores it in
// Statement.Settings so afterUpdate can persist it alongside new_values.
func beforeUpdate(db *gorm.DB) {
	if db.Statement == nil || shouldSkip(db) || db.Statement.Model == nil {
		return
	}

	cond := primaryKeyCondition(db)
	if cond == "" {
		return
	}

	var snapshot map[string]interface{}
	if err := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).
		Table(db.Statement.Table).
		Where(cond).
		Take(&snapshot).Error; err != nil {
		return // non-fatal: old_values will be NULL in the audit log
	}

	if b, err := json.Marshal(snapshot); err == nil {
		db.Statement.Settings.Store(settingOldValues, b)
	}
}

// afterCreate writes a "create" audit entry.
func afterCreate(db *gorm.DB) {
	if db.Statement == nil || shouldSkip(db) {
		return
	}
	insertAuditLog(db, "create", nil, marshalModel(db))
}

// afterUpdate writes an "update" audit entry with before/after snapshots.
func afterUpdate(db *gorm.DB) {
	if db.Statement == nil || shouldSkip(db) {
		return
	}

	var oldValues json.RawMessage
	if v, ok := db.Statement.Settings.Load(settingOldValues); ok {
		if b, ok := v.([]byte); ok {
			oldValues = b
		}
	}

	insertAuditLog(db, "update", oldValues, marshalModel(db))
}

// afterDelete writes a "delete" audit entry (old = deleted record, new = nil).
func afterDelete(db *gorm.DB) {
	if db.Statement == nil || shouldSkip(db) {
		return
	}
	insertAuditLog(db, "delete", marshalModel(db), nil)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// shouldSkip returns true when the statement targets a skipped table or has
// already encountered an error (we do not audit failed operations).
func shouldSkip(db *gorm.DB) bool {
	if db.Error != nil {
		return true
	}
	return skipTables[db.Statement.Table]
}

// insertAuditLog builds and persists one audit_logs row.
func insertAuditLog(db *gorm.DB, action string, oldValues, newValues json.RawMessage) {
	ctx := stmtContext(db)

	entry := dao.AuditLogDAO{
		Action:    action,
		Entity:    db.Statement.Table,
		EntityID:  extractPKValue(db),
		OldValues: oldValues,
		NewValues: newValues,
		IPAddress: logger.IPAddressFromContext(ctx),
	}

	if uid := logger.UserIDFromContext(ctx); uid != 0 {
		entry.UserID = &uid
	}

	// Use a fresh session with a background context so this insert is
	// independent of the parent transaction and never recurses into audit hooks.
	db.Session(&gorm.Session{
		NewDB:     true,
		SkipHooks: true,
		Context:   context.Background(),
	}).Create(&entry)
}

// stmtContext returns the context bound to the current statement, or Background.
func stmtContext(db *gorm.DB) context.Context {
	if db.Statement != nil && db.Statement.Context != nil {
		return db.Statement.Context
	}
	return context.Background()
}

// primaryKeyCondition builds a "col = 'val'" WHERE clause for the primary key
// of the statement model. Returns empty string when the schema is unavailable.
func primaryKeyCondition(db *gorm.DB) string {
	if db.Statement.Schema == nil {
		return ""
	}
	var parts []string
	modelVal := reflect.ValueOf(db.Statement.Model)
	for _, f := range db.Statement.Schema.PrimaryFields {
		val, zero := f.ValueOf(db.Statement.Context, modelVal)
		if zero {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s = '%v'", f.DBName, val))
	}
	return strings.Join(parts, " AND ")
}

// extractPKValue returns the primary key value as a string (for entity_id).
func extractPKValue(db *gorm.DB) string {
	if db.Statement.Schema == nil || db.Statement.Model == nil {
		return ""
	}
	modelVal := reflect.ValueOf(db.Statement.Model)
	for _, f := range db.Statement.Schema.PrimaryFields {
		val, _ := f.ValueOf(db.Statement.Context, modelVal)
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// marshalModel serialises db.Statement.Dest (preferred) or .Model to JSON.
func marshalModel(db *gorm.DB) json.RawMessage {
	target := db.Statement.Dest
	if target == nil {
		target = db.Statement.Model
	}
	if target == nil {
		return nil
	}
	b, _ := json.Marshal(target)
	return b
}
