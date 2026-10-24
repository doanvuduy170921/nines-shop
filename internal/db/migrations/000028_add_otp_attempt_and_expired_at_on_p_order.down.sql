-- down
DROP INDEX IF EXISTS idx_pending_orders_user_pending;
DROP INDEX IF EXISTS idx_pending_orders_expiry;
ALTER TABLE pending_orders ALTER COLUMN otp TYPE varchar(6);
ALTER TABLE pending_orders
DROP CONSTRAINT IF EXISTS chk_pending_orders_status,
  ALTER COLUMN status DROP NOT NULL,
  ALTER COLUMN status DROP DEFAULT,
  DROP COLUMN IF EXISTS expires_at,
  DROP COLUMN IF EXISTS otp_attempts;