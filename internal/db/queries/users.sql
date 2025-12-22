-- name: GetAllUser :many
SELECT * FROM users;


-- name: CreateUser :one
INSERT INTO users (name,email,password,phone,address,role,is_active)
values ($1,$2,$3,$4,$5,$6,$7)
returning *;

-- name: FindByEmail :one
select * from users
where email = $1;

-- name: GetAllUserV2 :many
SELECT u.id,
       u.name,
       u.email,
       u.address,
       u.phone,
       u.is_active,
       u.role,
       u.created_at,
       u.updated_at,
       u.user_uuid
FROM users u
WHERE created_at IS NOT NULL
  AND (
    sqlc.narg('search')::text IS NULL
    OR sqlc.narg('search')::text = ''
    OR u.name_search LIKE '%' || LOWER(sqlc.narg('search')::text) || '%'
    OR u.email_search LIKE '%' || LOWER(sqlc.narg('search')::text) || '%'
    )
  AND (
    sqlc.narg('role')::text IS NULL
    OR sqlc.narg('role')::text = ''
    OR u.role = sqlc.narg('role')::text
    )
  AND (
    sqlc.narg('is_active_filter')::text IS NULL
    OR sqlc.narg('is_active_filter')::text = ''
    OR u.is_active = sqlc.narg('is_active')::bool
    )
ORDER BY u.updated_at DESC, u.id DESC
    LIMIT $1 OFFSET $2;

-- name: CountUser :one
SELECT COUNT(*)
FROM users u
WHERE created_at IS NOT NULL
  AND (
    sqlc.narg('search')::text IS NULL
    OR sqlc.narg('search')::text = ''
    OR u.name_search LIKE '%' || LOWER(sqlc.narg('search')::text) || '%'
    OR u.email_search LIKE '%' || LOWER(sqlc.narg('search')::text) || '%'
    )
  AND (
    sqlc.narg('role')::text IS NULL
    OR sqlc.narg('role')::text = ''
    OR u.role = sqlc.narg('role')::text
    )
  AND (
    sqlc.narg('is_active_filter')::text IS NULL
    OR sqlc.narg('is_active_filter')::text = ''
    OR u.is_active = sqlc.narg('is_active')::bool
    );

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


-- name: ActiveUser :exec
update users
set is_active = true
where user_uuid = sqlc.arg(user_uuid);
