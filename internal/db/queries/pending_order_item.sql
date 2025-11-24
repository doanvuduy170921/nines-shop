-- name: CreatePendingOrderItem :one
insert into pending_order_items(pending_order_id,product_id,quantity,price)
values ($1,$2,$3,$4)
returning *;

-- name: GetByPOrderItemId :many
select *
from pending_order_items
where pending_order_id = @id::int;