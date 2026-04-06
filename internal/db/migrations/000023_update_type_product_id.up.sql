-- Sửa lại kiểu dữ liệu cho bảng variants
ALTER TABLE product_variants
ALTER COLUMN product_id TYPE bigint,
  ALTER COLUMN images TYPE jsonb USING images::jsonb; -- Lưu [url1, url2]

-- Sửa lại kiểu dữ liệu cho bảng specifications
ALTER TABLE product_specifications
ALTER COLUMN product_id TYPE bigint;