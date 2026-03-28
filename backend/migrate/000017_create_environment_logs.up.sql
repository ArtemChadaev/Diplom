-- ============================================================
-- 000017_create_environment_logs.up.sql
-- Environment monitoring journal — Phase 3.
-- One record per shift per day per zone.
-- ============================================================

CREATE TABLE environment_logs (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id              UUID NOT NULL REFERENCES warehouse_zones(id) ON DELETE CASCADE,
    digital_signature_id UUID,  -- employee digital signature (not implemented)
    recorded_by          INT NOT NULL REFERENCES users(id),
    shift                shift_type NOT NULL,
    temperature          NUMERIC(5,2) NOT NULL,
    humidity             NUMERIC(5,2) NOT NULL,
    notes                TEXT,
    recorded_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (zone_id, (recorded_at::DATE), shift)  -- one record per shift per day
);

CREATE INDEX idx_env_logs_zone_date ON environment_logs(zone_id, recorded_at);
