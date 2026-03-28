-- ============================================================
-- 000013_create_settings.up.sql
-- System settings key-value store.
-- ============================================================

CREATE TABLE settings (
    key        VARCHAR(100) PRIMARY KEY,
    value      TEXT NOT NULL,
    updated_by INT REFERENCES users(id),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Initial data
INSERT INTO settings (key, value) VALUES ('mos_percent', '60');
