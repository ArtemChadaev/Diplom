# Database Schema — Full SQL Migrations

← [Back to Main README](../../README.md) | [Backend Docs →](../backend.md)

<!-- All ENUMs and CREATE TABLE statements for the ERP system -->

all sql:
```sql
-- ============================================================
-- ОБЩИЕ ТИПЫ ДАННЫХ (ENUMS)
-- ============================================================

CREATE TYPE user_role AS ENUM (
    'admin',
    'qp',
    'warehouse_manager',
    'storekeeper',
    'pharmacist'
);

CREATE TYPE batch_status AS ENUM (
    'quarantine',
    'available',
    'rejected',
    'blocked'
);

CREATE TYPE zone_type AS ENUM (
    'general',
    'cold_chain',
    'flammable',
    'safe_strong'
);

CREATE TYPE shift_type AS ENUM ('morning', 'evening');

CREATE TYPE claim_type AS ENUM (
    'recall',
    'return_from_pharmacy',
    'return_to_supplier',
    'defect'
);

CREATE TYPE purchase_type AS ENUM ('direct', 'tender', 'state');

CREATE TYPE order_type AS ENUM ('regular', 'cito');

CREATE TYPE order_status AS ENUM ('new', 'assembling', 'ready', 'shipped', 'cancelled');

CREATE TYPE inventory_status AS ENUM ('draft', 'in_progress', 'completed', 'cancelled');

CREATE TYPE training_result AS ENUM ('pass', 'fail');

-- ============================================================
-- БЛОК «СПРАВОЧНИКИ И ПЕРСОНАЛ»
-- ============================================================

CREATE TABLE users (
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(255) UNIQUE NOT NULL,
    google_id    VARCHAR(255) UNIQUE,
    telegram_id  BIGINT UNIQUE,
    role         user_role NOT NULL DEFAULT 'pharmacist',
    ns_pv_access BOOLEAN NOT NULL DEFAULT false,
    ukep_bound   BOOLEAN NOT NULL DEFAULT false,
    is_blocked   BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

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
    medical_book_scan_url TEXT,
    special_zone_access   BOOLEAN NOT NULL DEFAULT false,
    gdp_training_history  JSONB NOT NULL DEFAULT '[]'
);

CREATE TABLE products (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_name           VARCHAR(255) NOT NULL,
    mnn                  VARCHAR(255) NOT NULL,
    sku                  VARCHAR(100) UNIQUE,
    barcode              VARCHAR(50),
    datamatrix_gtin      VARCHAR(50),
    ru_number            VARCHAR(100) NOT NULL,
    atc_codes            TEXT[] NOT NULL DEFAULT '{}',
    dosage_form          VARCHAR(100),
    dosage               VARCHAR(100),
    package_multiplicity INT NOT NULL DEFAULT 1,
    is_jnvlp             BOOLEAN NOT NULL DEFAULT false,
    is_mdlp              BOOLEAN NOT NULL DEFAULT false,
    is_ns_pv             BOOLEAN NOT NULL DEFAULT false,
    cold_chain           BOOLEAN NOT NULL DEFAULT false,
    temp_min             NUMERIC(5,2),
    temp_max             NUMERIC(5,2),
    humidity_max         NUMERIC(5,2),
    weight_g             NUMERIC(10,3),
    width_cm             NUMERIC(8,2),
    height_cm            NUMERIC(8,2),
    depth_cm             NUMERIC(8,2),
    description          TEXT,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE INDEX idx_products_mnn       ON products(mnn);
CREATE INDEX idx_products_sku       ON products(sku);
CREATE INDEX idx_products_ru_number ON products(ru_number);

CREATE TABLE product_photos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE suppliers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255) NOT NULL,
    inn            VARCHAR(12) UNIQUE NOT NULL,
    license_number VARCHAR(100),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- БЛОК «СКЛАДСКАЯ ЛОГИСТИКА И ОТГРУЗКА»
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

CREATE TABLE inbound_receipts (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id          UUID NOT NULL REFERENCES suppliers(id),
    purchase_type        purchase_type NOT NULL,
    invoice_number       VARCHAR(100) NOT NULL,
    country_of_origin    VARCHAR(3) NOT NULL,
    manufacturer         VARCHAR(255) NOT NULL,
    vat_rate             SMALLINT NOT NULL,
    is_jnvlp_controlled  BOOLEAN NOT NULL DEFAULT false,
    jnvlp_markup         NUMERIC(5,2),
    qp_user_id           INT REFERENCES users(id),
    inspection_date      DATE,
    inspection_result    VARCHAR(20),
    inspection_notes     TEXT,
    photo_urls           TEXT[] DEFAULT '{}',
    digital_signature_id UUID,
    created_by           INT NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE inbound_positions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbound_id UUID NOT NULL REFERENCES inbound_receipts(id) ON DELETE CASCADE,
    batch_id   UUID NOT NULL REFERENCES batches(id)
);

CREATE TABLE orders (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type                  order_type NOT NULL DEFAULT 'regular',
    status                order_status NOT NULL DEFAULT 'new',
    destination_id        UUID,
    customer_signature_id UUID,
    outbound_signature_id UUID,
    assembled_by          INT REFERENCES users(id),
    shipped_by            INT REFERENCES users(id),
    shipped_at            TIMESTAMPTZ,
    ttn_url               TEXT,
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
    batch_id      UUID REFERENCES batches(id),
    assembled_qty INT,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending'
);

-- ============================================================
-- БЛОК «КОНТРОЛЬ КАЧЕСТВА И МОНИТОРИНГ»
-- ============================================================

CREATE TABLE claims (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type                 claim_type NOT NULL,
    batch_id             UUID REFERENCES batches(id),
    product_id           UUID NOT NULL REFERENCES products(id),
    status               VARCHAR(20) NOT NULL DEFAULT 'open',
    digital_signature_id UUID,
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

CREATE TABLE inventory_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status       inventory_status NOT NULL DEFAULT 'draft',
    zone_id      UUID REFERENCES warehouse_zones(id),
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
    expected_qty INT NOT NULL,
    actual_qty   INT,
    discrepancy  INT GENERATED ALWAYS AS (actual_qty - expected_qty) STORED
);

CREATE TABLE environment_logs (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id              UUID NOT NULL REFERENCES warehouse_zones(id) ON DELETE CASCADE,
    digital_signature_id UUID,
    recorded_by          INT NOT NULL REFERENCES users(id),
    shift                shift_type NOT NULL,
    temperature          NUMERIC(5,2) NOT NULL,
    humidity             NUMERIC(5,2) NOT NULL,
    notes                TEXT,
    recorded_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (zone_id, (recorded_at::DATE), shift)
);

CREATE INDEX idx_env_logs_zone_date ON environment_logs(zone_id, recorded_at);
```
