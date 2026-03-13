CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'employee',
    is_blocked BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
