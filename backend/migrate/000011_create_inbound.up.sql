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
