-- ============================================================
-- 000003_create_categories.up.sql
-- Product categories (unchanged from original).
-- ============================================================

CREATE TABLE categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);
