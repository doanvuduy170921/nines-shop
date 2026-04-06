-- 1. Xóa ràng buộc UNIQUE và Index của variant_id
ALTER TABLE cart DROP CONSTRAINT IF EXISTS unique_user_variant;
DROP INDEX IF EXISTS idx_cart_variant_id;

-- 2. Thêm lại cột product_id
ALTER TABLE cart ADD COLUMN product_id INTEGER REFERENCES products(id) ON DELETE CASCADE;

-- 3. Xóa cột variant_id
ALTER TABLE cart DROP COLUMN variant_id;

-- 4. Tạo lại Index cũ cho product_id
CREATE INDEX idx_cart_product_id ON cart (product_id);