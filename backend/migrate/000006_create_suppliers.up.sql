-- ============================================================
-- 000006_create_suppliers.up.sql
-- Supplier registry — new table.
-- ============================================================

CREATE TABLE suppliers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255) NOT NULL,
    inn            VARCHAR(12) UNIQUE NOT NULL,  -- taxpayer ID
    license_number VARCHAR(100),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);
