ALTER TABLE pending_orders
    ADD COLUMN otp_attempts int NOT NULL DEFAULT 0,
  ADD COLUMN expires_at timestamptz;

UPDATE pending_orders
SET expires_at = COALESCE(otp_expires_at, created_at, now());

ALTER TABLE pending_orders ALTER COLUMN expires_at SET NOT NULL;

ALTER TABLE pending_orders
    ALTER COLUMN status SET DEFAULT 'pending',
ALTER COLUMN status SET NOT NULL,
  ADD CONSTRAINT chk_pending_orders_status
    CHECK (status IN ('pending', 'confirmed', 'cancelled', 'expired'));

ALTER TABLE pending_orders ALTER COLUMN otp TYPE varchar(64);

CREATE INDEX idx_pending_orders_expiry
    ON pending_orders (expires_at) WHERE status = 'pending';

CREATE INDEX idx_pending_orders_user_pending
    ON pending_orders (user_id) WHERE status = 'pending';

