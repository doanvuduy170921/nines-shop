-- 1. Thêm lại các cột đã bị drop (khôi phục schema cũ)
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS sku VARCHAR(100),
    ADD COLUMN IF NOT EXISTS price NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS discount_price NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS stock_quantity INTEGER;

-- 2. Xoá cột has_variants
ALTER TABLE products
DROP COLUMN IF EXISTS has_variants;