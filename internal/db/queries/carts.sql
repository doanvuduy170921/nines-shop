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