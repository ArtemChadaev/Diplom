-- ============================================================
-- 000013_create_orders.up.sql
-- Orders + order items.
-- ============================================================

CREATE TABLE orders (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number          VARCHAR(100) UNIQUE NOT NULL,
    customer_name         VARCHAR(255) NOT NULL,
    status                order_status NOT NULL DEFAULT 'new',
    priority              INT NOT NULL DEFAULT 1,
    created_by            INT NOT NULL,
    created_at            TIMESTAMPTZ DEFAULT NOW(),
    updated_at            TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_number ON orders(order_number);

CREATE TABLE order_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id    UUID NOT NULL REFERENCES products(id),
    quantity      INT NOT NULL CHECK (quantity > 0),
    picked_qty    INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);
