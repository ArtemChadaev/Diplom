-- ============================================================
-- 000013_create_orders.up.sql
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
