-- name: CreateCategory :one
insert into categories(name)
values ($1)
returning *;


-- name: GetAllCategories :many
select * from categories;