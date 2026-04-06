-- name: AddToCart :one
insert into cart(user_id,variant_id,quantity)
values ($1,$2,$3)
returning *;

-- name: ExistsProductId :one
SELECT
    (EXISTS (
        SELECT 1
        FROM cart
        WHERE variant_id = @variant_id::int
            AND user_id = @user_id::int
    )) AS exists_;

-- name: UpdateCart :one
update cart
set quantity = quantity + @quantity::int
where user_id = @user_id::int
and variant_id = @variant_id::int
returning *;


-- name: GetCartsByUserId :many
SELECT
    c.id AS cart_id,
    c.quantity,
    pv.id AS variant_id,
    pv.price,
    pv.sku,
    pv.attributes,
    p.id AS product_id,
    p.name,
    p.thumbnail
FROM cart c
         JOIN product_variants pv ON c.variant_id = pv.id
         JOIN products p ON pv.product_id = p.id
WHERE c.user_id = @user_id::int
ORDER BY c.created_at DESC;

-- name: DeleteItemInCart :exec
DELETE FROM cart
WHERE user_id = @user_id::int
  AND variant_id = @variant_id::int;


-- name: UpdateAllCart :one
update cart
set quantity = @quantity::int
where user_id = @user_id::int
and variant_id = @variant_id::int
returning *;