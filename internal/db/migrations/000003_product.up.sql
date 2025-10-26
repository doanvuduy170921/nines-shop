-- create table products
CREATE TABLE products (
                          id BIGSERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          slug VARCHAR(255) UNIQUE NOT NULL,
                          sku VARCHAR(100) UNIQUE NOT NULL,
                          brand_id INT,
                          category_id INT,
                          description TEXT,
                          short_description VARCHAR(500),
                          price DECIMAL(10,2) DEFAULT 0.00,
                          discount_price DECIMAL(10,2) DEFAULT 0.00,
                          stock_quantity INT DEFAULT 0,
                          status VARCHAR(20) CHECK (status IN ('draft', 'active', 'inactive', 'archived')) DEFAULT 'active',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Create indexes for better query performance
CREATE INDEX idx_product_name ON products (name);
CREATE INDEX idx_product_slug ON products (slug);
CREATE INDEX idx_product_sku ON products (sku);


-- Create function to auto-update updated_at
CREATE OR REPLACE FUNCTION update_to_updated_at_on_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for products table
CREATE TRIGGER update_product_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_to_updated_at_on_column();
