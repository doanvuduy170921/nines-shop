
-- name: CreateProduct :one
insert into products(name,slug,sku,brand_id,category_id,description,short_description,price,discount_price,stock_quantity)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
returning *;


-- name: GetAllProductByFilter :many
select p.id,
       p.name,
       p.sku,
       p.slug,
       p.description,
       p.short_description,
       p.price,
       p.discount_price,
       p.stock_quantity,
       p.thumbnail,
       p.status,
        b.name as brand_name,
        c.name as category_name

from products p
left join brand b on b.id = p.brand_id
left join categories c on c.id = p.category_id
where p.created_at is not null
and(
    @search::text = ''
    or p.name ILIKE '%' || @search::text ||'%'
    or p.sku ILIKE '%' || @search::text ||'%'
    or p.slug ILIKE '%' || @search::text ||'%'
    )

and (
    p.price >= coalesce(@min_price::int,0)
and p.price <= coalesce(@max_price::int,99999999)
    )
and(
    @status::text ='' OR p.status = @status::text
    )

and(
     @category_id::int = 0
    or p.category_id = @category_id::int
    )
order by p.updated_at desc
limit $1
offset $2;


-- name: CountProduct :one
select count(*)
from products
where created_at is not null
and (
    @search::text = ''
    or name ILIKE '%' || @search::text ||'%'
    or sku ILIKE '%' || @search::text ||'%'
    or slug ILIKE '%' || @search::text ||'%'
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

-- name: UpdateThumbnail :one
update products
set thumbnail = @thumbnail::text
where id = @id::int
and created_at is not null
returning *;

-- name: GetProductByCategoryId :many
select p.id,
       p.name,
       p.sku,
       p.slug,
       p.description,
       p.short_description,
       p.price,
       p.discount_price,
       p.stock_quantity,
       p.thumbnail,
       p.status,
       b.name as brand_name,
       c.name as category_name
from products p
left join brand b on b.id = p.brand_id
left join categories c on c.id = p.category_id
where p.created_at is not null
and(
    p.category_id = @id::int
);

-- name: GetProductBySlug :one
select p.id,
       p.name,
       p.sku,
       p.slug,
       p.description,
       p.short_description,
       p.price,
       p.discount_price,
       p.stock_quantity,
       p.thumbnail,
       p.status,
       b.name as brand_name,
       c.name as category_name
from products p
         left join brand b on b.id = p.brand_id
         left join categories c on c.id = p.category_id
where p.created_at is not null
  and(
    p.slug = @slug::text
    );

