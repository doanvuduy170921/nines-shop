CREATE TABLE pending_orders (
                                id SERIAL PRIMARY KEY,
                                user_id INTEGER NOT NULL,
                                name VARCHAR(255) NOT NULL,
                                email VARCHAR(255) NOT NULL,
                                phone VARCHAR(20) NOT NULL,
                                payment_method_id INTEGER NOT NULL,
                                address TEXT NOT NULL,
                                otp VARCHAR(6) NOT NULL,
                                otp_expires_at TIMESTAMP NOT NULL,
                                subtotal NUMERIC(10,2) DEFAULT 0.00 NOT NULL,
                                total_amount NUMERIC(10,2) NOT NULL,
                                shipping_price NUMERIC(10,2) DEFAULT 0.00 NOT NULL,
                                tax NUMERIC(10,2) DEFAULT 0.00 NOT NULL,
                                status VARCHAR(50),
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);