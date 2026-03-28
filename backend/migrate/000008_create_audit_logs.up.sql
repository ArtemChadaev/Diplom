-- ============================================================
-- 000008_create_audit_logs.up.sql
-- Audit log with immutable hash chain support.
-- ============================================================

CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     INT REFERENCES users(id),
    action      VARCHAR(255) NOT NULL,
    entity      VARCHAR(100) NOT NULL,
    entity_id   VARCHAR(100),
    old_values  JSONB,
    new_values  JSONB,
    ip_address  INET,
    -- Immutability chain (Phase 3, currently stored as NULL / mock):
    prev_hash   TEXT,    -- hash of previous log entry in chain
    log_hash    TEXT,    -- SHA-256(prev_hash + user_id + action + entity_id + new_values + created_at)
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_entity  ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);