-- 1. Xóa các ràng buộc và index cũ liên quan đến product_id
ALTER TABLE cart DROP CONSTRAINT IF EXISTS cart_product_id_fkey;
DROP INDEX IF EXISTS idx_cart_product_id;

-- 2. Thêm cột variant_id
ALTER TABLE cart ADD COLUMN variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE;

-- 3. Xóa cột product_id (vì đã có variant_id để truy xuất ngược lại)
ALTER TABLE cart DROP COLUMN product_id;

-- 4. Tạo Index cho variant_id để tối ưu tốc độ truy vấn giỏ hàng
CREATE INDEX idx_cart_variant_id ON cart (variant_id);

-- 5. Thêm ràng buộc UNIQUE để đảm bảo mỗi người dùng chỉ có 1 dòng cho mỗi biến thể
-- Giúp logic UPSERT (Insert or Update) hoạt động chính xác
ALTER TABLE cart ADD CONSTRAINT unique_user_variant UNIQUE (user_id, variant_id);