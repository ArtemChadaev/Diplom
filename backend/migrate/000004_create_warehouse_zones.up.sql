-- ============================================================
-- 000004_create_warehouse_zones.up.sql
-- Warehouse zones.
-- ============================================================

CREATE TABLE warehouse_zones (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) UNIQUE NOT NULL,
    type         zone_type NOT NULL DEFAULT 'ambient',
    description  TEXT,
    temp_min     NUMERIC(5,2),
    temp_max     NUMERIC(5,2),
    capacity     INT NOT NULL DEFAULT 0,
    is_active    BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);
