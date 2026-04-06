-- HOÀN TÁC BẢNG PENDING_ORDER_ITEMS
ALTER TABLE pending_order_items DROP CONSTRAINT IF EXISTS fk_pending_order_items_variant;
ALTER TABLE pending_order_items RENAME COLUMN variant_id TO product_id;
ALTER TABLE pending_order_items
    ADD CONSTRAINT pending_order_items_product_id_fkey
        FOREIGN KEY (product_id) REFERENCES products(id);


-- HOÀN TÁC BẢNG ORDER_ITEMS
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS fk_order_items_variant;
ALTER TABLE order_items RENAME COLUMN variant_id TO product_id;
ALTER TABLE order_items
    ADD CONSTRAINT order_items_product_id_fkey
        FOREIGN KEY (product_id) REFERENCES products(id);