CREATE TABLE product_specifications (
                                        id SERIAL PRIMARY KEY,
                                        product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                                        spec_key VARCHAR(100) NOT NULL,    -- Ví dụ: 'Màn hình'
                                        spec_value TEXT NOT NULL,         -- Ví dụ: '6.7 inch, OLED'
                                        display_order INTEGER DEFAULT 0   -- Để sắp xếp thứ tự hiển thị trên UI
);

CREATE TABLE product_variants (
                                  id SERIAL PRIMARY KEY,
                                  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                                  sku VARCHAR(100) UNIQUE,           -- Mã định danh SKU
                                  attributes JSONB NOT NULL,         -- Lưu dạng: {"color": "Red", "storage": "128GB"}
                                  price NUMERIC(12,2) NOT NULL,      -- Giá bán biến thể
                                  stock_quantity INTEGER NOT NULL DEFAULT 0,
                                  images TEXT,                       -- Lưu 1 link ảnh đại diện cho biến thể
                                  is_active BOOLEAN DEFAULT TRUE,
                                  created_at TIMESTAMP DEFAULT NOW()
);

-- Index quan trọng để tìm kiếm biến thể theo thuộc tính (màu, size) cực nhanh
CREATE INDEX idx_variant_attrs ON product_variants USING GIN (attributes);


CREATE TABLE attribute_configs (
                                   id SERIAL PRIMARY KEY,
                                   attr_name VARCHAR(50) NOT NULL,   -- Ví dụ: 'color'
                                   attr_value VARCHAR(100) NOT NULL, -- Ví dụ: 'Red'
                                   display_label VARCHAR(100),       -- Ví dụ: 'Màu Đỏ'
                                   color_code VARCHAR(7)             -- Lưu mã màu hex nếu là thuộc tính màu sắc
);