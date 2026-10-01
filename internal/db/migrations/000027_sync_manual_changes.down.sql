ALTER TABLE pending_order_items ALTER COLUMN price TYPE numeric(10,2);
ALTER TABLE pending_orders ALTER COLUMN otp_expires_at SET NOT NULL;

DROP INDEX IF EXISTS idx_user_filter_sort;

ALTER TABLE pending_orders DROP CONSTRAINT IF EXISTS fk_payment_method;
ALTER TABLE pending_orders DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE pending_order_items DROP CONSTRAINT IF EXISTS fk_pending_order;
ALTER TABLE pending_order_items ADD CONSTRAINT pending_order_items_pending_order_id_fkey FOREIGN KEY (pending_order_id) REFERENCES public.pending_orders(id) ON DELETE CASCADE;