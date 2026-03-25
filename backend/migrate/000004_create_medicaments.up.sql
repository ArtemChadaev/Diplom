CREATE TABLE medicaments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    release_form VARCHAR(100),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- TODO: Сделать чтобы был штрихкод, так же условия хранения и м.б. многие ко многим со складами