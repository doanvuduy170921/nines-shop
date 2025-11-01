-- name: SaveAndUploadImg :one
insert into product_images (product_id,image_url)
values ($1,$2)
returning *;


-- name: GetImagesByProductId :many
select p.image_url
from product_images p
where product_id = @id::int;