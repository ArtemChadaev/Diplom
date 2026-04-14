-- ============================================================
-- 000005_create_suppliers.up.sql
-- Supplier registry.
-- ============================================================

CREATE TABLE suppliers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255) NOT NULL,
    inn            VARCHAR(12) UNIQUE NOT NULL,
    kpp            VARCHAR(9),
    contact_name   VARCHAR(255),
    phone          VARCHAR(50),
    email          VARCHAR(100),
    address        TEXT,
    is_active      BOOLEAN NOT NULL DEFAULT true,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);
