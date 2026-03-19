CREATE TABLE stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id INT REFERENCES warehouses(id) ON DELETE CASCADE,
    medicament_id UUID REFERENCES medicaments(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity >= 0),
    series VARCHAR(100),
    expiry_date DATE NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
