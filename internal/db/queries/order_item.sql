-- name: AddOrderItem :one
insert into order_items(order_id,product_id,quantity,price,product_name,product_thumbnail)
values ($1,$2,$3,$4,$5,$6)
returning *;

-- name: GetOrderItemByUserId :many
SELECT
    -- order info
    o.id AS order_id,
    o.name AS customer_name,
    o.phone,
    o.address,
    o.subtotal,
    o.shipping_price,
    o.tax,
    o.total_amount,
    o.status,
    o.created_at AS order_created_at,

    -- item info
    oi.product_id,
    oi.quantity,
    oi.price AS item_price,
    oi.product_name,
    oi.product_thumbnail,
    oi.created_at AS item_created_at,

    -- payment method
    p.name AS payment_method_name

FROM order_items oi
         JOIN orders o ON oi.order_id = o.id
         JOIN payment_methods p ON p.id = o.payment_method_id
WHERE o.user_id = @id::int
ORDER BY oi.created_at DESC;