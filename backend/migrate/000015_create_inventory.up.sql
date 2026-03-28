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
