-- Drop trigger
DROP TRIGGER IF EXISTS trg_update_categories_updated_at ON categories;

-- Drop function
DROP FUNCTION IF EXISTS update_categories_updated_at();

-- Drop index
DROP INDEX IF EXISTS idx_categories_name;

-- Drop table
DROP TABLE IF EXISTS categories;