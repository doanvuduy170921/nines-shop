CREATE TABLE payment_methods (
                                 id SERIAL PRIMARY KEY,
                                 name VARCHAR(100) NOT NULL,       -- Ví dụ: 'COD', 'VNPAY'
                                 description TEXT,
                                 is_active BOOLEAN DEFAULT TRUE
);