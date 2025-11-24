CREATE TABLE order_items (
                             id SERIAL PRIMARY KEY,

                             order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

                             product_id INT NOT NULL REFERENCES products(id),

                             quantity INT NOT NULL CHECK (quantity > 0),

    -- giá tại thời điểm đặt, để cố định khi chạy flash sale, thay đổi giá, v.v.
                             price NUMERIC(10,2) NOT NULL,

    -- snapshot thông tin sản phẩm tại thời điểm đặt
                             product_name VARCHAR(255) NOT NULL,
                             product_thumbnail TEXT,
                             created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
