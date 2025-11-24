-- name: CreateOrderStatusHistory :one
insert into order_status_history(order_id,status,note)
values ($1,$2,$3)
returning *;