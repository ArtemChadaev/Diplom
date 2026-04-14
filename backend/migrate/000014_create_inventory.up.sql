-- ============================================================
-- 000014_create_inventory.up.sql
-- Inventory sessions, items.
-- ============================================================

CREATE TABLE inventory_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status       inventory_status NOT NULL DEFAULT 'draft',
    zone_id      UUID,                 -- NULL = entire warehouse
    started_by   INT NOT NULL,
    started_at   TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE TABLE inventory_items (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id     UUID NOT NULL REFERENCES inventory_sessions(id) ON DELETE CASCADE,
    product_id     UUID NOT NULL REFERENCES products(id),
    batch_number   VARCHAR(100) NOT NULL,
    system_qty     INT NOT NULL,
    physical_qty   INT NOT NULL,
    reason         VARCHAR(255),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);
