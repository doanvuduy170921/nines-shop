-- Nâng cấp bảng pending_orders
ALTER TABLE pending_orders ALTER COLUMN subtotal TYPE numeric(15, 2);
ALTER TABLE pending_orders ALTER COLUMN total_amount TYPE numeric(15, 2);
ALTER TABLE pending_orders ALTER COLUMN shipping_price TYPE numeric(15, 2);
ALTER TABLE pending_orders ALTER COLUMN tax TYPE numeric(15, 2);

-- Nâng cấp bảng orders (Nếu bạn đã có)
ALTER TABLE orders ALTER COLUMN subtotal TYPE numeric(15, 2);
ALTER TABLE orders ALTER COLUMN total_amount TYPE numeric(15, 2);
ALTER TABLE orders ALTER COLUMN shipping_price TYPE numeric(15, 2);
ALTER TABLE orders ALTER COLUMN tax TYPE numeric(15, 2);

-- Nâng cấp bảng order_items và product_variants
ALTER TABLE order_items ALTER COLUMN price TYPE numeric(15, 2);
ALTER TABLE product_variants ALTER COLUMN price TYPE numeric(15, 2);