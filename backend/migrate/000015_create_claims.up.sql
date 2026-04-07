-- ============================================================
-- 000015_create_claims.up.sql
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
