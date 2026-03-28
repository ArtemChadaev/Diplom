-- ============================================================
-- 000002_create_users.up.sql
-- Users table — redesigned per ERP spec.
-- Removes: login, password_hash, status column.
-- Adds:    ns_pv_access, ukep_bound, updated_at.
-- ============================================================

CREATE TABLE users (
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(255) UNIQUE NOT NULL,
    google_id    VARCHAR(255) UNIQUE,
    telegram_id  BIGINT UNIQUE,
    role         user_role NOT NULL DEFAULT 'pharmacist',
    ns_pv_access BOOLEAN NOT NULL DEFAULT false,  -- access to narcotic/psychotropic (НС/ПВ)
    ukep_bound   BOOLEAN NOT NULL DEFAULT false,  -- qualified electronic signature linked
    is_blocked   BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);
