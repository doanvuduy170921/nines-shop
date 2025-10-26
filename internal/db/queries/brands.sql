-- name: GetAll :many
select *
from brand
where created_at is not null;