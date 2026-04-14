-- ============================================================
-- 000003_create_products.up.sql
-- Products table 
-- ============================================================

CREATE TABLE products (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku                  VARCHAR(100) UNIQUE,
    name                 VARCHAR(255) NOT NULL,
    generic_name         VARCHAR(255) NOT NULL,
    atc_code             VARCHAR(255),
    dosage_form          VARCHAR(100),
    strength             VARCHAR(100),
    package_size         INT NOT NULL DEFAULT 1,
    is_jnvlp             BOOLEAN NOT NULL DEFAULT false,
    manufacturer_id      UUID,
    storage_conditions   TEXT,
    photo_url            TEXT,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ  -- soft delete
);

CREATE INDEX idx_products_generic_name ON products(generic_name);
CREATE INDEX idx_products_sku ON products(sku);
