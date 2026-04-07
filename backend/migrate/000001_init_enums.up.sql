-- ============================================================
-- 000001_init_enums.up.sql
-- All application-level ENUM types in a single migration.
-- Replaces the original partial 000001 + scattered ALTER TYPE calls.
-- ============================================================

-- Roles
CREATE TYPE user_role AS ENUM (
    'admin',
    'qp',               -- Уполномоченное лицо / QP
    'warehouse_manager',
    'storekeeper',
    'pharmacist'
);

-- Batch (series) status
CREATE TYPE batch_status AS ENUM (
    'quarantine',  -- accepted, awaiting inspection
    'available',   -- released for distribution
    'rejected',    -- rejected at acceptance
    'blocked'      -- blocked (STOP signal / Roszdravnadzor)
);

-- Warehouse zone type
CREATE TYPE zone_type AS ENUM (
    'general',
    'cold_chain',
    'flammable',
    'safe_strong'   -- safe/narcotic (НС/ПВ)
);

-- Shift for environment journal
CREATE TYPE shift_type AS ENUM ('morning', 'evening');

-- Claim type
CREATE TYPE claim_type AS ENUM (
    'recall',               -- Roszdravnadzor withdrawal
    'return_from_pharmacy', -- return from pharmacy
    'return_to_supplier',   -- return to supplier
    'defect'                -- manufacturing defect
);

-- Purchase type
CREATE TYPE purchase_type AS ENUM ('direct', 'tender', 'state');

-- Order type
CREATE TYPE order_type AS ENUM ('regular', 'cito');

-- Order status
CREATE TYPE order_status AS ENUM ('new', 'assembling', 'ready', 'shipped', 'cancelled');

-- Inventory status
CREATE TYPE inventory_status AS ENUM ('draft', 'in_progress', 'completed', 'cancelled');

-- GDP training result
CREATE TYPE training_result AS ENUM ('pass', 'fail');
