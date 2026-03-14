ALTER TABLE users
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS google_id,
    DROP COLUMN IF EXISTS telegram_id,
    DROP COLUMN IF EXISTS status;

ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;

-- Note: Postgres does not support dropping ENUM values easily (without recreating the type).
-- We'll leave the 'unverified' enum value alone in the down migration.
