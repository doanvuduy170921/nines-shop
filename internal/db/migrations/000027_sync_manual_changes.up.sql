-- 1. Sửa kiểu dữ liệu và ràng buộc của pending_order_items & pending_orders
ALTER TABLE pending_order_items ALTER COLUMN price TYPE numeric(15,2);

ALTER TABLE pending_orders ALTER COLUMN otp_expires_at DROP NOT NULL;

-- 2. Thêm Index mới
CREATE INDEX idx_user_filter_sort ON public.users USING btree (role, is_active, updated_at DESC);

-- 3. Cập nhật / Thêm các Foreign Key Constraints
ALTER TABLE pending_orders DROP CONSTRAINT IF EXISTS fk_payment_method;
ALTER TABLE pending_orders ADD CONSTRAINT fk_payment_method FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id) ON DELETE RESTRICT;

ALTER TABLE pending_orders DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE pending_orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

ALTER TABLE pending_order_items DROP CONSTRAINT IF EXISTS pending_order_items_pending_order_id_fkey;
ALTER TABLE pending_order_items DROP CONSTRAINT IF EXISTS fk_pending_order;
ALTER TABLE pending_order_items ADD CONSTRAINT fk_pending_order FOREIGN KEY (pending_order_id) REFERENCES public.pending_orders(id) ON DELETE CASCADE;