
-- name: CreateProduct :one
insert into products(name,slug,sku,brand_id,category_id,description,short_description,price,discount_price,stock_quantity)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
returning *;


-- name: GetAllProductByFilter :many
select *
from products
where created_at is not null
and(
    @search::text = ''
    or @search::text ilike '%'||sku||'%'
     or @search::text ilike '%'||slug||'%'
     or @search::text ilike '%'||name||'%'
    )

and (
    price >= coalesce(@min_price::int,0)
and price <= coalesce(@max_price::int,99999999)
    )
and(
    @status::text ='' OR status = @status::text
    )

and(
     @category_id::int = 0
    or category_id = @category_id::int
    )
order by updated_at desc
limit $1
offset $2;


-- name: CountProduct :one
select count(*)
from products
where created_at is not null
and (
    @search::text = ''
    or name ilike '%' || @search::text ||'%'
    or sku ilike '%' || @search::text ||'%'
    or slug ilike '%' || @search::text ||'%'
    )
and
(
    price >= coalesce(@min_price::int,0)
    and price <= coalesce(@max_price::int,99999999)
    )
    and(
    @status::text =''
    OR status = @status::text
    )
and(
    @category_id::int = 0
    or category_id = @category_id::int
    );

-- name: GetProductById :one
select *
from products
where created_at is not null
and id = @id::int;