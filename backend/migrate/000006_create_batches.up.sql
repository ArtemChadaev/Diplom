-- ============================================================
-- 000006_create_batches.up.sql
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
