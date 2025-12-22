-- name: CreateOrderStatusHistory :one
insert into order_status_history(order_id,status,note)
values ($1,$2,$3)
returning *;

-- name: GetAllStatusByOrderId :many
select o.id,
       o.status as order_status,
       o.amount_item,
       o.created_at as order_created_at,
       osh.status as order_history_status,
       osh.note,
       osh.created_at as order_history_created_at,
       oi.product_thumbnail
from order_status_history osh
        join orders o on o.id = osh.order_id
join order_items oi on oi.order_id = osh.order_id
where osh.order_id =sqlc.arg(order_id);

-- name: GetAllStatusByOrderIdV2 :many
select id,
       status,
       note,
       created_at
from order_status_history
where order_id = sqlc.arg(id);


