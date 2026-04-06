-- CẬP NHẬT BẢNG PENDING_ORDER_ITEMS
-- 1. Xóa ràng buộc khóa ngoại cũ trỏ tới bảng products
ALTER TABLE pending_order_items DROP CONSTRAINT IF EXISTS pending_order_items_product_id_fkey;

-- 2. Đổi tên cột product_id thành variant_id
ALTER TABLE pending_order_items RENAME COLUMN product_id TO variant_id;

-- 3. Thêm ràng buộc khóa ngoại mới trỏ tới bảng product_variants
ALTER TABLE pending_order_items
    ADD CONSTRAINT fk_pending_order_items_variant
        FOREIGN KEY (variant_id) REFERENCES product_variants(id);


-- CẬP NHẬT BẢNG ORDER_ITEMS
-- 1. Xóa ràng buộc khóa ngoại cũ trỏ tới bảng products
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS order_items_product_id_fkey;

-- 2. Đổi tên cột product_id thành variant_id
ALTER TABLE order_items RENAME COLUMN product_id TO variant_id;

-- 3. Thêm ràng buộc khóa ngoại mới trỏ tới bảng product_variants
ALTER TABLE order_items
    ADD CONSTRAINT fk_order_items_variant
        FOREIGN KEY (variant_id) REFERENCES product_variants(id);