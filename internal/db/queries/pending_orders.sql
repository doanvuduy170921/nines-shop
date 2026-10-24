-- name: CreatePendingOrder :one
INSERT INTO pending_orders(
    user_id, name, email, phone, payment_method_id, address,
    otp, otp_expires_at, expires_at, subtotal, total_amount, shipping_price,
    tax, amount_item
)
VALUES (
           @user_id, @name, @email, @phone, @payment_method_id, @address,
           @otp, @otp_expires_at, @expires_at::timestamptz,
           @subtotal::NUMERIC(15,2),
           @total_amount::NUMERIC(15,2),
           @shipping_price::NUMERIC(15,2),
           @tax::NUMERIC(15,2),
           @amount_item
       )
    RETURNING *;


-- name: GetPOrderById :one
select *
from pending_orders
where id = @id::int;

-- name: IncrementOtpAttempts :one
UPDATE pending_orders
SET otp_attempts = otp_attempts + 1
WHERE id = @id::int AND status = 'pending'
RETURNING otp_attempts;

-- name: ConfirmPendingOrder :one
UPDATE pending_orders
SET status = 'confirmed'
WHERE id = @id::int
  AND user_id = @user_id::int
  AND status = 'pending'
  AND expires_at > now()
RETURNING id;

-- name: ListExpiredPendingOrderIDs :many
SELECT id FROM pending_orders
WHERE status = 'pending' AND expires_at < now()
ORDER BY expires_at
    LIMIT 100;
