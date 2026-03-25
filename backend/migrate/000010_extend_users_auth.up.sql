-- Add social auth and verification fields to users
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email          VARCHAR(255) UNIQUE,
    ADD COLUMN IF NOT EXISTS google_id      VARCHAR(255) UNIQUE,
    ADD COLUMN IF NOT EXISTS telegram_id    BIGINT UNIQUE,
    ADD COLUMN IF NOT EXISTS status         VARCHAR(50) NOT NULL DEFAULT 'unverified';

-- Make password_hash nullable to support social-only users
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Add 'unverified' to the user_role ENUM
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'unverified';
-- Убрать пароль и просто пересмотреть особенно user_role