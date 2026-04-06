-- 1. Xóa các trigger trước
DROP TRIGGER IF EXISTS update_product_updated_at ON products;

-- 2. Xóa các index (Mặc dù xóa bảng sẽ tự xóa index,
-- nhưng viết rõ ràng giúp script minh bạch hơn)
DROP INDEX IF EXISTS idx_variant_attrs;
DROP INDEX IF EXISTS idx_product_slug;
DROP INDEX IF EXISTS idx_product_name;

-- 3. Xóa các bảng theo thứ tự (Bảng con xóa trước để tránh lỗi Foreign Key)
DROP TABLE IF EXISTS attribute_configs;
DROP TABLE IF EXISTS product_specifications;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;

-- 4. Lưu ý: Không xóa function update_to_updated_at_on_column()
-- nếu nó được dùng chung cho các bảng khác (như categories, orders...)
-- Nếu bạn muốn xóa sạch cả function đó, hãy bỏ comment dòng dưới:
-- DROP FUNCTION IF EXISTS update_to_updated_at_on_column;