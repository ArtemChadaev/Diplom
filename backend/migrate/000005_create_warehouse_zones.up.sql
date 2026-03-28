-- ============================================================
-- 000005_create_warehouse_zones.up.sql
-- Warehouse zones — replaces old flat "warehouses" table.
-- Supports zoning: general, cold_chain, flammable, safe_strong.
-- ============================================================

CREATE TABLE warehouse_zones (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    type         zone_type NOT NULL DEFAULT 'general',
    description  TEXT,
    temp_min     NUMERIC(5,2),
    temp_max     NUMERIC(5,2),
    humidity_max NUMERIC(5,2),
    created_at   TIMESTAMPTZ DEFAULT NOW()
);
