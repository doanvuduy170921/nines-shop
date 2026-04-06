-- Drop index
DROP INDEX IF EXISTS idx_variant_attrs;

-- Drop tables (child → parent)
DROP TABLE IF EXISTS product_specifications;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS attribute_configs;