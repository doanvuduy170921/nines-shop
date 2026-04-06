-- 1. Thêm cột đánh dấu có biến thể hay không
ALTER TABLE products ADD COLUMN has_variants BOOLEAN DEFAULT FALSE;

-- 2. Loại bỏ các cột đã chuyển sang bảng product_variants
-- Lưu ý: Hãy chắc chắn bạn đã backup hoặc không cần dữ liệu cũ trong các cột này
ALTER TABLE products DROP COLUMN IF EXISTS sku CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS price CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS discount_price CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS stock_quantity CASCADE;