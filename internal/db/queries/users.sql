-- name: GetAllUser :many
SELECT * FROM users;


-- name: CreateUser :one
INSERT INTO users (name,email,password,phone,address,role)
values ($1,$2,$3,$4,$5,$6)
returning *;

-- name: FindByEmail :one
select * from users
where email = $1;

-- name: GetAllUserV2 :many
SELECT *
FROM users
WHERE created_at IS NOT NULL
  AND (
    @search::text = ''
    OR name ILIKE '%' || @search || '%'
    OR email ILIKE '%' || @search || '%'
    )
  AND (
    @role::text = '' OR role = @role
    )
  AND (
    @is_active_filter::text = ''
    OR is_active = @is_active::bool
    )
order by updated_at desc 
LIMIT $1
OFFSET $2;

-- name: CountUser :one
SELECT count(*)
FROM users
WHERE created_at IS NOT NULL
  AND (
    @search::text = ''
    OR name ILIKE '%' || @search || '%'
    OR email ILIKE '%' || @search || '%'
    )
  AND (
    @role::text = '' OR role = @role
    )
  AND (
    @is_active_filter::text = ''
    OR is_active = @is_active::bool
    )
    ;

-- name: SoftDeleteUser :one
update users
SET is_active = false
where user_uuid = @user_uuid::uuid
returning *;


-- name: UpdateUser :one
update users
SET name = @name::text,
    email = @email::text,
    phone = @phone::text,
    address = @address::text,
    role = @role,
    is_active = @is_active::boolean
where user_uuid = @user_uuid::uuid
returning *;

-- name: GetByUuid :one
select *
from users
where user_uuid = @user_uuid::uuid;
