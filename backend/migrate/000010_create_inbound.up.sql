-- ============================================================
-- 000010_create_inbound.up.sql
-- Inbound receipts + items.
-- ============================================================

CREATE TYPE inbound_status AS ENUM ('draft', 'received', 'completed', 'cancelled');

CREATE TABLE inbound_receivings (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number       VARCHAR(100) NOT NULL,
    invoice_date         DATE NOT NULL,
    supplier_id          UUID NOT NULL REFERENCES suppliers(id),
    status               inbound_status NOT NULL DEFAULT 'draft',
    total_amount         NUMERIC(10, 2) NOT NULL DEFAULT 0,
    vat_amount           NUMERIC(10, 2) NOT NULL DEFAULT 0,
    notes                TEXT,
    received_by          INT, 
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE inbound_items (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbound_id       UUID NOT NULL REFERENCES inbound_receivings(id) ON DELETE CASCADE,
    product_id       UUID NOT NULL REFERENCES products(id),
    batch_number     VARCHAR(100) NOT NULL,
    expiration_date  DATE NOT NULL,
    quantity         INT NOT NULL,
    price_netto      NUMERIC(10, 2) NOT NULL,
    vat_rate         NUMERIC(5, 2) NOT NULL,
    price_brutto     NUMERIC(10, 2) NOT NULL,
    cert_number      VARCHAR(255),
    zone_id          UUID NOT NULL REFERENCES warehouse_zones(id),
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);
