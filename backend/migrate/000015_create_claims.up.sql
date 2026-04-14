-- ============================================================
-- 000015_create_claims.up.sql
-- Claims.
-- ============================================================

CREATE TYPE claim_status AS ENUM ('open', 'in_progress', 'resolved', 'closed');

CREATE TABLE claims (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title                VARCHAR(255) NOT NULL,
    description          TEXT,
    inbound_id           UUID REFERENCES inbound_receivings(id),
    order_id             UUID REFERENCES orders(id),
    status               claim_status NOT NULL DEFAULT 'open',
    created_by           INT NOT NULL,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW()
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
