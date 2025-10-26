
CREATE TABLE IF NOT EXISTS brand (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


CREATE INDEX IF NOT EXISTS idx_brand_name ON brand(name);


CREATE OR REPLACE FUNCTION update_on_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER updated_on_column
BEFORE UPDATE ON brand
FOR EACH ROW
EXECUTE FUNCTION update_on_updated_at();
