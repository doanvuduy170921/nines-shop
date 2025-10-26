DROP TRIGGER IF EXISTS updated_on_column ON brand;

-- Drop function
DROP FUNCTION IF EXISTS update_on_updated_at();

-- Drop index
DROP INDEX IF EXISTS idx_brand_name;

-- Drop table
DROP TABLE IF EXISTS brand;