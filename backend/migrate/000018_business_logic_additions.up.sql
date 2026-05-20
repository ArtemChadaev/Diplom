-- ============================================================
-- 000018_business_logic_additions.up.sql
-- Add fields required for the remaining business logic.
-- ============================================================

CREATE TYPE ven_category AS ENUM ('V', 'E', 'N');

ALTER TABLE products 
  ADD COLUMN ven_category ven_category NOT NULL DEFAULT 'N',
  ADD COLUMN lead_time_days INT NOT NULL DEFAULT 14,
  ADD COLUMN safety_stock_qty INT NOT NULL DEFAULT 0,
  ADD COLUMN max_stock_qty INT NOT NULL DEFAULT 1000;

ALTER TABLE orders 
  ADD COLUMN order_type order_type NOT NULL DEFAULT 'regular';

ALTER TABLE order_items 
  ADD COLUMN batch_id UUID REFERENCES batches(id),
  ADD COLUMN mos_blocked BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE claims 
  ADD COLUMN type claim_type NOT NULL DEFAULT 'defect',
  ADD COLUMN product_id UUID REFERENCES products(id);
