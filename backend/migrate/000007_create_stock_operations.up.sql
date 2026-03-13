CREATE TABLE stock_operations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT REFERENCES users(id),
    warehouse_id INT REFERENCES warehouses(id),
    medicament_id UUID REFERENCES medicaments(id),
    type operation_type NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
