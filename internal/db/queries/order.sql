-- name: CreateOrder :one
INSERT INTO orders(
    user_id, name, email, phone, payment_method_id, address,
    subtotal, total_amount, shipping_price, tax, status, amount_item,
    transaction_id, payment_status
)
VALUES (
           @user_id, @name, @email, @phone, @payment_method_id, @address,
           @subtotal::NUMERIC(10,2),
           @total_amount::NUMERIC(10,2),
           @shipping_price::NUMERIC(10,2),
           @tax::NUMERIC(10,2),
           @status,
           @amount_item,
           @transaction_id,
           @payment_status
       )
    RETURNING *;


-- name: UpdateStatusForUser :exec
UPDATE orders
SET status = @status
WHERE id = @order_id;

-- name: GetOrderById :one
select *
from orders
where id = @order_id;


-- name: GetAllOrders :many
select o.id,
       o.name,
       o.created_at,
       o.status,
       o.total_amount,
       o.amount_item
from orders o
order by o.created_at desc
;

-- name: GetOrderDetailById :many
select o.name as customer_name,
       o.id,
       o.email,
       o.address,
       o.phone,
       p.name as payment_method_name,
       o.subtotal,
       o.shipping_price,
       o.tax,
       o.total_amount,
       o.status,
       o.created_at as order_created_at,
       oi.product_name,
       oi.quantity,
       oi.price as item_price,
       o.transaction_id,
       o.payment_status,
       oi.product_thumbnail
from orders o
join order_items oi on oi.order_id = o.id
join payment_methods p on p.id = o.payment_method_id
where o.id =sqlc.arg(id)
;



-- name: UpdateOrderPayment :exec
UPDATE orders
SET
    payment_status = @payment_status::text,
    transaction_id = @transaction_id,
    updated_at = now()
WHERE id = @id;

-- name: GetListOrderByOrderId :many
select
    o.id,
    o.status,
    o.created_at,
    o.amount_item,
    o.total_amount,
    (
        select array_agg(oi.product_thumbnail)
        from order_items oi
        where oi.order_id = o.id
    ) as preview_thumbnails
from orders o
where
    o.user_id = @user_id::int
and (
    sqlc.narg(search)::int is null
    or o.id = sqlc.narg(search)::int
)
and (
    sqlc.narg(status)::text is null
    or o.status = sqlc.narg(status)::text
)
order by o.created_at desc
    limit $1
offset $2;

-- name: ViewDetailForMyOrder :many
select p.name as payment_method_name,
       o.shipping_price,
       oi.product_name as product_name,
       oi.quantity,
       oi.price,
       o.amount_item,
       o.subtotal,
       o.tax,
       o.total_amount,
       o.name as user_name,
       o.address,
       o.phone,
       oi.product_thumbnail
from orders o
join payment_methods p on p.id = o.payment_method_id
join order_items oi on oi.order_id = o.id
where o.id = sqlc.arg(order_id);
