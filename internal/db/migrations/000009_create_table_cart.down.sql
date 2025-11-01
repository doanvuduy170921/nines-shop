-- drop trigger before dropping function or table
DROP TRIGGER IF EXISTS trg_cart_updated_at ON cart;

-- drop function
DROP FUNCTION IF EXISTS update_cart_updated_at();

-- drop indexes (optional, will be dropped automatically with table,
-- nhưng viết rõ ra cho chắc chắn)
DROP INDEX IF EXISTS idx_cart_user_id;
DROP INDEX IF EXISTS idx_cart_product_id;

-- drop table
DROP TABLE IF EXISTS cart;