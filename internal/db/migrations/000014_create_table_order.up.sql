CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id INTEGER NOT NULL,

                        name VARCHAR(255) NOT NULL,
                        email VARCHAR(255) NOT NULL,
                        phone VARCHAR(20) NOT NULL,
                        address TEXT NOT NULL,

                        payment_method_id INTEGER NOT NULL,

                        subtotal NUMERIC(10,2) NOT NULL,
                        shipping_price NUMERIC(10,2) NOT NULL,
                        tax NUMERIC(10,2) NOT NULL,
                        total_amount NUMERIC(10,2) NOT NULL,

                        status VARCHAR(50) NOT NULL DEFAULT 'PENDING_PAYMENT',

                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);