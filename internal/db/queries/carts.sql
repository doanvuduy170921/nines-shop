-- name: AddToCart :one
insert into cart(user_id,product_id,quantity)
values ($1,$2,$3)
returning *;

-- name: ExistsProductId :one
SELECT
    (EXISTS (
        SELECT 1
        FROM cart
        WHERE product_id = @product_id::int
            AND user_id = @user_id::int
    )) AS exists_;

-- name: UpdateCart :one
update cart
set quantity = quantity + @quantity::int
where user_id = @user_id::int
and product_id = @product_id::int
returning *;


-- name: GetCartsByUserId :many
select p.name,
       p.id,
       p.thumbnail,
       p.price,
       p.stock_quantity,
       c.quantity
from cart c
left join products p on p.id = c.product_id
where user_id = @user_id::int;

-- name: DeleteItemInCart :exec
DELETE FROM cart
WHERE user_id = @user_id::int
  AND product_id = @product_id::int;


-- name: UpdateAllCart :one
update cart
set quantity = @quantity::int
where user_id = @user_id::int
and product_id = @product_id::int
returning *;