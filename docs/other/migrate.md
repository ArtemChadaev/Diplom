all sql:
```sql
-- File: backend/migrate/000001_init_enums.down.sql
DROP TYPE IF EXISTS training_result;
DROP TYPE IF EXISTS inventory_status;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS order_type;
DROP TYPE IF EXISTS purchase_type;
DROP TYPE IF EXISTS claim_type;
DROP TYPE IF EXISTS shift_type;
DROP TYPE IF EXISTS zone_type;
DROP TYPE IF EXISTS batch_status;
DROP TYPE IF EXISTS user_role;
-- File: backend/migrate/000001_init_enums.up.sql
-- ============================================================
-- 000001_init_enums.up.sql
-- All application-level ENUM types in a single migration.
-- Replaces the original partial 000001 + scattered ALTER TYPE calls.
-- ============================================================

-- Roles
CREATE TYPE user_role AS ENUM (
    'admin',
    'qp',               -- Уполномоченное лицо / QP
    'warehouse_manager',
    'storekeeper',
    'pharmacist'
);

-- Batch (series) status
CREATE TYPE batch_status AS ENUM (
    'quarantine',  -- accepted, awaiting inspection
    'available',   -- released for distribution
    'rejected',    -- rejected at acceptance
    'blocked'      -- blocked (STOP signal / Roszdravnadzor)
);

-- Warehouse zone type
CREATE TYPE zone_type AS ENUM (
    'general',
    'cold_chain',
    'flammable',
    'safe_strong'   -- safe/narcotic (НС/ПВ)
);

-- Shift for environment journal
CREATE TYPE shift_type AS ENUM ('morning', 'evening');

-- Claim type
CREATE TYPE claim_type AS ENUM (
    'recall',               -- Roszdravnadzor withdrawal
    'return_from_pharmacy', -- return from pharmacy
    'return_to_supplier',   -- return to supplier
    'defect'                -- manufacturing defect
);

-- Purchase type
CREATE TYPE purchase_type AS ENUM ('direct', 'tender', 'state');

-- Order type
CREATE TYPE order_type AS ENUM ('regular', 'cito');

-- Order status
CREATE TYPE order_status AS ENUM ('new', 'assembling', 'ready', 'shipped', 'cancelled');

-- Inventory status
CREATE TYPE inventory_status AS ENUM ('draft', 'in_progress', 'completed', 'cancelled');

-- GDP training result
CREATE TYPE training_result AS ENUM ('pass', 'fail');
-- File: backend/migrate/000002_create_users.down.sql
DROP TABLE IF EXISTS users;
-- File: backend/migrate/000002_create_users.up.sql
-- ============================================================
-- 000002_create_users.up.sql
-- Users table — redesigned per ERP spec.
-- Removes: login, password_hash, status column.
-- Adds:    ns_pv_access, ukep_bound, updated_at.
-- ============================================================

CREATE TABLE users (
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(255) UNIQUE NOT NULL,
    google_id    VARCHAR(255) UNIQUE,
    telegram_id  BIGINT UNIQUE,
    role         user_role NOT NULL DEFAULT 'pharmacist',
    ns_pv_access BOOLEAN NOT NULL DEFAULT false,  -- access to narcotic/psychotropic (НС/ПВ)
    ukep_bound   BOOLEAN NOT NULL DEFAULT false,  -- qualified electronic signature linked
    is_blocked   BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);
-- File: backend/migrate/000003_create_categories.down.sql
DROP TABLE IF EXISTS categories;
-- File: backend/migrate/000003_create_categories.up.sql
-- ============================================================
-- 000003_create_categories.up.sql
-- Product categories (unchanged from original).
-- ============================================================

CREATE TABLE categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);
-- File: backend/migrate/000004_create_products.down.sql
DROP TABLE IF EXISTS products;
-- File: backend/migrate/000004_create_products.up.sql
-- ============================================================
-- 000004_create_products.up.sql
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
-- File: backend/migrate/000005_create_warehouse_zones.down.sql
DROP TABLE IF EXISTS warehouse_zones;
-- File: backend/migrate/000005_create_warehouse_zones.up.sql
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
-- File: backend/migrate/000006_create_suppliers.down.sql
DROP TABLE IF EXISTS suppliers;
-- File: backend/migrate/000006_create_suppliers.up.sql
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
-- File: backend/migrate/000007_create_batches.down.sql
DROP TABLE IF EXISTS batches;
-- File: backend/migrate/000007_create_batches.up.sql
-- ============================================================
-- 000007_create_batches.up.sql
-- Batches (series) — replaces old stock_items + stock_operations.
-- Tracks inventory by product series with quarantine workflow.
-- ============================================================

CREATE TABLE batches (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id       UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    zone_id          UUID REFERENCES warehouse_zones(id) ON DELETE SET NULL,
    serial_number    VARCHAR(100) NOT NULL,
    manufacture_date DATE NOT NULL,
    expiry_date      DATE NOT NULL,
    quantity         INT NOT NULL CHECK (quantity >= 0),
    status           batch_status NOT NULL DEFAULT 'quarantine',
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (product_id, serial_number)
);

CREATE INDEX idx_batches_product_id  ON batches(product_id);
CREATE INDEX idx_batches_expiry_date ON batches(expiry_date);
CREATE INDEX idx_batches_status      ON batches(status);
-- File: backend/migrate/000008_create_audit_logs.down.sql
DROP TABLE IF EXISTS audit_logs;
-- File: backend/migrate/000008_create_audit_logs.up.sql
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
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);-- File: backend/migrate/000009_create_employee_profiles.down.sql
DROP TABLE IF EXISTS employee_profiles;
-- File: backend/migrate/000009_create_employee_profiles.up.sql
-- ============================================================
-- 000009_create_employee_profiles.up.sql
-- Employee profiles — extended with GDP training history,
-- medical book scan, special zone access.
-- Removes: telegram_handle, emergency_contact
--   (those were pre-ERP fields not in the spec).
-- ============================================================

CREATE TABLE employee_profiles (
    id                    SERIAL PRIMARY KEY,
    user_id               INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_code         VARCHAR(100) UNIQUE NOT NULL,
    full_name             VARCHAR(255) NOT NULL,
    corporate_email       VARCHAR(255) UNIQUE NOT NULL,
    phone                 VARCHAR(20) UNIQUE NOT NULL,
    position              VARCHAR(255) NOT NULL,
    department            VARCHAR(255) NOT NULL,
    birth_date            DATE NOT NULL,
    avatar_url            TEXT,
    hire_date             DATE NOT NULL,
    dismissal_date        DATE,
    -- New ERP fields:
    medical_book_scan_url TEXT,                           -- medical book scan URL
    special_zone_access   BOOLEAN NOT NULL DEFAULT false, -- access to special zones
    gdp_training_history  JSONB NOT NULL DEFAULT '[]'    -- GDP training history
    -- JSONB element structure:
    -- { "date": "YYYY-MM-DD", "course_name": "...", "result": "pass"|"fail", "certificate_url": "..." }
);
-- File: backend/migrate/000010_create_refresh_tokens.down.sql
DROP TABLE IF EXISTS refresh_tokens;
-- File: backend/migrate/000010_create_refresh_tokens.up.sql
-- ============================================================
-- 000010_create_refresh_tokens.up.sql
-- Session / refresh token storage (unchanged from spec).
-- ============================================================

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    user_agent  TEXT,
    ip_address  INET,
    metadata    JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);

CREATE INDEX idx_refresh_tokens_user_id    ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
-- File: backend/migrate/000011_create_inbound.down.sql
DROP TABLE IF EXISTS inbound_positions;
DROP TABLE IF EXISTS inbound_receipts;
-- File: backend/migrate/000011_create_inbound.up.sql
-- ============================================================
-- 000011_create_inbound.up.sql
-- Inbound receipts + positions — new tables for acceptance workflow.
-- ============================================================

CREATE TABLE inbound_receipts (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id          UUID NOT NULL REFERENCES suppliers(id),
    purchase_type        purchase_type NOT NULL,
    invoice_number       VARCHAR(100) NOT NULL,
    country_of_origin    VARCHAR(3) NOT NULL,          -- ISO country code
    manufacturer         VARCHAR(255) NOT NULL,
    vat_rate             SMALLINT NOT NULL,             -- 0, 10, 20
    is_jnvlp_controlled  BOOLEAN NOT NULL DEFAULT false,
    jnvlp_markup         NUMERIC(5,2),
    -- Acceptance protocol
    qp_user_id           INT REFERENCES users(id),
    inspection_date      DATE,
    inspection_result    VARCHAR(20),                   -- 'approved' | 'rejected'
    inspection_notes     TEXT,
    -- Attachments
    photo_urls           TEXT[] DEFAULT '{}',           -- photos of damaged packaging
    digital_signature_id UUID,                          -- detached УКЭП signature reference
    -- Meta
    created_by           INT NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE inbound_positions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbound_id UUID NOT NULL REFERENCES inbound_receipts(id) ON DELETE CASCADE,
    batch_id   UUID NOT NULL REFERENCES batches(id)
);
-- File: backend/migrate/000012_create_product_photos.down.sql
DROP TABLE IF EXISTS product_photos;
-- File: backend/migrate/000012_create_product_photos.up.sql
-- ============================================================
-- 000012_create_product_photos.up.sql
-- Product photos — new table.
-- ============================================================

CREATE TABLE product_photos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- File: backend/migrate/000013_create_settings.down.sql
DROP TABLE IF EXISTS settings;
-- File: backend/migrate/000013_create_settings.up.sql
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
-- File: backend/migrate/000014_create_orders.down.sql
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
-- File: backend/migrate/000014_create_orders.up.sql
-- ============================================================
-- 000014_create_orders.up.sql
-- Orders + order items — Phase 2 (FEFO-based warehouse dispatch).
-- ============================================================

CREATE TABLE orders (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type                  order_type NOT NULL DEFAULT 'regular',
    status                order_status NOT NULL DEFAULT 'new',
    destination_id        UUID,              -- FK → destinations (future table)
    -- Signatures (not implemented, stored as NULL)
    customer_signature_id UUID,
    outbound_signature_id UUID,
    assembled_by          INT REFERENCES users(id),
    shipped_by            INT REFERENCES users(id),
    shipped_at            TIMESTAMPTZ,
    ttn_url               TEXT,              -- consignment note URL
    created_by            INT NOT NULL REFERENCES users(id),
    created_at            TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_type   ON orders(type);

CREATE TABLE order_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id    UUID NOT NULL REFERENCES products(id),
    requested_qty INT NOT NULL CHECK (requested_qty > 0),
    batch_id      UUID REFERENCES batches(id),  -- filled during assembly (FEFO)
    assembled_qty INT,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending'
    -- 'pending', 'mos_blocked', 'assembled'
);
-- File: backend/migrate/000015_create_inventory.down.sql
DROP TABLE IF EXISTS inventory_samples;
DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS inventory_sessions;
-- File: backend/migrate/000015_create_inventory.up.sql
-- ============================================================
-- 000015_create_inventory.up.sql
-- Inventory sessions, items, and control samples — Phase 2.
-- ============================================================

CREATE TABLE inventory_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status       inventory_status NOT NULL DEFAULT 'draft',
    zone_id      UUID REFERENCES warehouse_zones(id),  -- NULL = entire warehouse
    started_by   INT NOT NULL REFERENCES users(id),
    completed_by INT REFERENCES users(id),
    started_at   TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE TABLE inventory_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id   UUID NOT NULL REFERENCES inventory_sessions(id) ON DELETE CASCADE,
    product_id   UUID NOT NULL REFERENCES products(id),
    batch_id     UUID REFERENCES batches(id),
    expected_qty INT NOT NULL,    -- from system, hidden until completion
    actual_qty   INT,             -- entered by employee
    discrepancy  INT GENERATED ALWAYS AS (actual_qty - expected_qty) STORED
);

CREATE TABLE inventory_samples (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES inventory_sessions(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    batch_id   UUID REFERENCES batches(id),
    qty        INT NOT NULL CHECK (qty > 0),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- File: backend/migrate/000016_create_claims.down.sql
DROP TABLE IF EXISTS recalled_batches;
DROP TABLE IF EXISTS claim_photos;
DROP TABLE IF EXISTS claims;
-- File: backend/migrate/000016_create_claims.up.sql
-- ============================================================
-- 000016_create_claims.up.sql
-- Claims (рекламации), claim photos, recalled batches — Phase 3.
-- ============================================================

CREATE TABLE claims (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type                 claim_type NOT NULL,
    batch_id             UUID REFERENCES batches(id),
    product_id           UUID NOT NULL REFERENCES products(id),
    status               VARCHAR(20) NOT NULL DEFAULT 'open',
    -- 'open', 'blocked', 'closed'
    digital_signature_id UUID,   -- signature approving block/return (not implemented)
    source               TEXT,
    notes                TEXT,
    resolution           TEXT,
    created_by           INT NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    closed_at            TIMESTAMPTZ
);

CREATE TABLE claim_photos (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id    UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    url         TEXT NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- Recalled batches (Roszdravnadzor sync — MOCK data).
-- ============================================================

CREATE TABLE recalled_batches (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial_number VARCHAR(100) NOT NULL,
    product_name  VARCHAR(255),
    ru_number     VARCHAR(100),
    recall_reason TEXT,
    issued_at     DATE,
    synced_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_recalled_batches_serial ON recalled_batches(serial_number);
-- File: backend/migrate/000017_create_environment_logs.down.sql
DROP TABLE IF EXISTS environment_logs;
-- File: backend/migrate/000017_create_environment_logs.up.sql
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
```
