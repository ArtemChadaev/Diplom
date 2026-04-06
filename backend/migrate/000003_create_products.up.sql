-- ============================================================
-- 000003_create_products.up.sql
-- Products table — replaces old "medicaments" table.
-- Full ERP spec: MNN, ATC codes, MDLP, JNVLP, НС/ПВ,
-- cold chain, dimensions, soft delete.
-- ============================================================

CREATE TABLE products (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_name           VARCHAR(255) NOT NULL,           -- trade name
    mnn                  VARCHAR(255) NOT NULL,           -- INN (МНН)
    sku                  VARCHAR(100) UNIQUE,             -- article number
    barcode              VARCHAR(50),                     -- barcode
    datamatrix_gtin      VARCHAR(50),                     -- Честный Знак GTIN
    ru_number            VARCHAR(100) NOT NULL,           -- registration certificate №
    atc_codes            TEXT[] NOT NULL DEFAULT '{}',   -- ATC code array
    dosage_form          VARCHAR(100),                    -- dosage form
    dosage               VARCHAR(100),                    -- dosage
    package_multiplicity INT NOT NULL DEFAULT 1,          -- units per package
    -- Flags
    is_jnvlp             BOOLEAN NOT NULL DEFAULT false, -- essential medicines list
    is_mdlp              BOOLEAN NOT NULL DEFAULT false,  -- subject to labelling
    is_ns_pv             BOOLEAN NOT NULL DEFAULT false,  -- narcotic/psychotropic (ПКУ)
    cold_chain           BOOLEAN NOT NULL DEFAULT false,  -- requires cold chain
    -- Storage conditions
    temp_min             NUMERIC(5,2),
    temp_max             NUMERIC(5,2),
    humidity_max         NUMERIC(5,2),
    -- Dimensions
    weight_g             NUMERIC(10,3),
    width_cm             NUMERIC(8,2),
    height_cm            NUMERIC(8,2),
    depth_cm             NUMERIC(8,2),
    -- Meta
    description          TEXT,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ  -- soft delete
);

CREATE INDEX idx_products_mnn       ON products(mnn);
CREATE INDEX idx_products_sku       ON products(sku);
CREATE INDEX idx_products_ru_number ON products(ru_number);
