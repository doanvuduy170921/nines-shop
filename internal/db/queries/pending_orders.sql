-- name: CreatePendingOrder :one
INSERT INTO pending_orders(user_id,name,email,phone,payment_method_id,address,otp,otp_expires_at,subtotal,total_amount,shipping_price,tax,status,amount_item)
VALUES (@user_id,@name,@email,@phone,@payment_method_id,@address,@otp,
           @otp_expires_at,
           @subtotal::NUMERIC(10,2),
           @total_amount::NUMERIC(10,2),
           @shipping_price::NUMERIC(10,2),
           @tax::NUMERIC(10,2),
           @status,
            @amount_item
       )
    RETURNING *;


-- name: GetPOrderById :one
select *
from pending_orders
where id = @id::int;
