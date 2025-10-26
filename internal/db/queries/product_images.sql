-- name: SaveAndUploadImg :one
insert into product_images (product_id,image_url)
values ($1,$2)
returning *;