CREATE TABLE pending_order_items (
                                     id SERIAL PRIMARY KEY,
                                     pending_order_id INT NOT NULL REFERENCES pending_orders(id) ON DELETE CASCADE,
                                     product_id INT NOT NULL REFERENCES products(id),
                                     quantity INT NOT NULL CHECK (quantity > 0),
                                     price DECIMAL(10,2) NOT NULL                      -- giá tại thời điểm đặt (để cố định giá)
);