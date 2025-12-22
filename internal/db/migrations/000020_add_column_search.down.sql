ALTER TABLE users
DROP COLUMN IF EXISTS name_search,
    DROP COLUMN IF EXISTS email_search;

-- Không drop extension vì có thể được dùng ở nơi khác
-- DROP EXTENSION IF EXISTS pg_trgm;