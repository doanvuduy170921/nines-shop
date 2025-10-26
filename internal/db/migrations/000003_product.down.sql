-- Drop trigger
DROP TRIGGER IF EXISTS update_product_updated_at ON products;

-- Drop function
DROP FUNCTION IF EXISTS update_to_updated_at_on_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_product_name;
DROP INDEX IF EXISTS idx_product_slug;
DROP INDEX IF EXISTS idx_product_sku;

-- Drop table
DROP TABLE IF EXISTS products;