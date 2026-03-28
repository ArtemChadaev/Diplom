-- ============================================================
-- 000009_create_employee_profiles.up.sql
-- Employee profiles — extended with GDP training history,
-- medical book scan, special zone access.
-- Removes: telegram_handle, emergency_contact
--   (those were pre-ERP fields not in the spec).
-- ============================================================

CREATE TABLE employee_profiles (
    id                    SERIAL PRIMARY KEY,
    user_id               INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_code         VARCHAR(100) UNIQUE NOT NULL,
    full_name             VARCHAR(255) NOT NULL,
    corporate_email       VARCHAR(255) UNIQUE NOT NULL,
    phone                 VARCHAR(20) UNIQUE NOT NULL,
    position              VARCHAR(255) NOT NULL,
    department            VARCHAR(255) NOT NULL,
    birth_date            DATE NOT NULL,
    avatar_url            TEXT,
    hire_date             DATE NOT NULL,
    dismissal_date        DATE,
    -- New ERP fields:
    medical_book_scan_url TEXT,                           -- medical book scan URL
    special_zone_access   BOOLEAN NOT NULL DEFAULT false, -- access to special zones
    gdp_training_history  JSONB NOT NULL DEFAULT '[]'    -- GDP training history
    -- JSONB element structure:
    -- { "date": "YYYY-MM-DD", "course_name": "...", "result": "pass"|"fail", "certificate_url": "..." }
);
