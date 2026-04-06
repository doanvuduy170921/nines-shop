
-- name: AddProduct :one
    insert into products(name,slug,brand_id,category_id,description,short_description,status,thumbnail,has_variants)
    values ($1,$2,$3,$4,$5,$6,$7,$8,$9)
returning *;

-- name: AddProductSpec :one
insert into product_specifications(product_id,spec_key,spec_value,display_order)
values ($1,$2,$3,$4)
    returning *;


-- name: AddProductVariant :one
    insert into product_variants(product_id,sku,attributes,price,stock_quantity,images,is_active)
    values ($1,$2,$3,$4,$5,$6,$7)
    returning *;

-- name: AddAttributesConf :one
    insert into attribute_configs(attr_name,attr_value,display_label)
    values ($1,$2,$3)
    returning *;

-- name: GetAllProductByFilter :many
SELECT
    p.id,
    p.name,
    p.status,
    p.thumbnail,
    c.name AS category_name,
    b.name AS brand_name,
    COALESCE(SUM(pv.stock_quantity), 0) AS total_stock,
    MIN(pv.price) AS min_price,
    MAX(pv.price) AS max_price,
    COUNT(pv.id) AS variant_count,
    p.created_at
FROM products p
         JOIN categories c ON c.id = p.category_id
         JOIN brand b ON b.id = p.brand_id
         LEFT JOIN product_variants pv ON pv.product_id = p.id
GROUP BY p.id, c.name, b.name;


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

-- name: GetListVariantByPid :many
select pv.id,
       pv.sku,
       pv.attributes,
       pv.price,
       pv.stock_quantity,
       pv.images::jsonb as images,
       pv.is_active
from product_variants pv
where product_id =sqlc.arg(product_id);


-- name: GetTop3Thumbnail :many
select  products.thumbnail
from products
where category_id =8
    limit 3
offset 3;


-- name: GetTop3Trending :many
SELECT p.name,
       p.thumbnail,
       MIN(pv.price) AS min_price, -- Lấy giá thấp nhất trong các biến thể
       b.name AS brand_name,
       p.slug
FROM products p
         LEFT JOIN product_variants pv ON p.id = pv.product_id
         LEFT JOIN brand b ON b.id = p.brand_id
WHERE p.category_id = sqlc.arg(cate_id)
GROUP BY p.id, p.name, p.thumbnail, b.name
    limit  4
offset 2;

-- name: GetProductBySlug :one
SELECT
    p.id,
    p.name,
    p.thumbnail,
    p.status,
    b.name AS brand_name,
    c.name AS category_name,
    p.description,
    -- Gộp tất cả biến thể vào một mảng JSON
    jsonb_agg(
            jsonb_build_object(
                    'variant_id', pv.id,
                    'attributes', pv.attributes,
                    'price', pv.price,
                    'stock_quantity', pv.stock_quantity,
                    'sku', pv.sku
            )
    ) AS variants,
    (select jsonb_agg(pi.image_url) from product_images pi where pi.product_id = p.id ) as gallery,
    -- Gộp Technical Specs
    (SELECT jsonb_agg(jsonb_build_object('key', ps.spec_key, 'value', ps.spec_value))
     FROM product_specifications ps WHERE ps.product_id = p.id) AS technical_specs
FROM products p
         LEFT JOIN categories c ON c.id = p.category_id
         LEFT JOIN brand b ON b.id = p.brand_id
         LEFT JOIN product_variants pv ON pv.product_id = p.id
WHERE p.slug = sqlc.arg(slug)
GROUP BY p.id, b.id, c.id;


-- name: GetVariantById :one
SELECT
    v.id,
    v.price,
    v.stock_quantity,
    v.images as product_variant_img,
    p.name as product_name,
    p.thumbnail as product_thumbnail
FROM product_variants v
         JOIN products p ON v.product_id = p.id where v.id = @id::int;

-- name: GetVariantForUpdate :one
SELECT
    pv.id,
    pv.stock_quantity,
    pv.price,
    pv.is_active,
    p.name as product_name,
    p.thumbnail as product_thumbnail
FROM product_variants pv
         JOIN products p ON pv.product_id = p.id
WHERE pv.id = $1
    LIMIT 1
FOR UPDATE;


-- name: DecreaseStock :execresult
update product_variants
set stock_quantity = stock_quantity - sqlc.arg(stock)::int
where id = sqlc.arg(variant_id)
and stock_quantity >= sqlc.arg(stock)::int
and is_active = true;


-- name: GetListProducts :many
SELECT
    p.id,
    p.name,
    p.slug,
    p.thumbnail,
    b.name AS brand_name,
    c.name AS cate_name,
    MIN(pv.price) AS price_product
FROM products p
         LEFT JOIN brand b ON b.id = p.brand_id
         LEFT JOIN categories c ON c.id = p.category_id
         LEFT JOIN product_variants pv ON pv.product_id = p.id
WHERE
    (p.name ILIKE '%' || sqlc.arg(search_name) || '%' OR sqlc.arg(search_name) = '')
  AND (c.name ILIKE '%' || sqlc.arg(cate_name) || '%' OR sqlc.arg(cate_name) = '')
  AND (b.name ILIKE '%' || sqlc.arg(brand_name) || '%' OR sqlc.arg(brand_name) = '')
GROUP BY
    p.id, p.name, p.slug, p.thumbnail, b.name, c.name
HAVING
    MIN(pv.price) >= sqlc.arg(min_price)
   AND MIN(pv.price) <= sqlc.arg(max_price)
ORDER BY
    CASE WHEN sqlc.arg(sort_desc)::boolean THEN MIN(pv.price) END DESC,
    CASE WHEN NOT sqlc.arg(sort_desc)::boolean THEN MIN(pv.price) END ASC
    LIMIT $1 OFFSET $2;

-- name: CountGetListProducts :one
SELECT COUNT(*)
FROM (
         SELECT p.id
         FROM products p -- Đảm bảo tên bảng là products (có s) đồng nhất
                  LEFT JOIN brand b ON b.id = p.brand_id
                  LEFT JOIN categories c ON c.id = p.category_id
                  LEFT JOIN product_variants pv ON pv.product_id = p.id
         WHERE
             (p.name ILIKE '%' || sqlc.arg(search_name) || '%' OR sqlc.arg(search_name) = '')
           AND (c.name ILIKE '%' || sqlc.arg(cate_name) || '%' OR sqlc.arg(cate_name) = '')
           AND (b.name ILIKE '%' || sqlc.arg(brand_name) || '%' OR sqlc.arg(brand_name) = '')
         GROUP BY
             p.id -- Chỉ cần Group By id là đủ để tính MIN price
         HAVING
             MIN(pv.price) >= sqlc.arg(min_price)
            AND MIN(pv.price) <= sqlc.arg(max_price)
     ) AS filtered_product;