-- name: CreatePendingOrderItem :one
insert into pending_order_items(pending_order_id,variant_id,quantity,price)
values ($1,$2,$3,$4)
returning *;

-- name: GetByPOrderItemId :many
select *
from pending_order_items
where pending_order_id = @id::int;

-- name: GetCountItems :one
select count(*) as count
from pending_order_items
where pending_order_id = @id
group by pending_order_id;