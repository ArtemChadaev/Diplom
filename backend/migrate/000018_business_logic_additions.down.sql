-- ============================================================
-- 000018_business_logic_additions.down.sql
-- Rollback additions for business logic.
-- ============================================================

ALTER TABLE claims 
  DROP COLUMN IF EXISTS type,
  DROP COLUMN IF EXISTS product_id;

ALTER TABLE order_items 
  DROP COLUMN IF EXISTS batch_id,
  DROP COLUMN IF EXISTS mos_blocked;

ALTER TABLE orders 
  DROP COLUMN IF EXISTS order_type;

ALTER TABLE products 
  DROP COLUMN IF EXISTS ven_category,
  DROP COLUMN IF EXISTS lead_time_days,
  DROP COLUMN IF EXISTS safety_stock_qty,
  DROP COLUMN IF EXISTS max_stock_qty;

DROP TYPE IF EXISTS ven_category;
