-- name: UpdateUser :exec
update users
set
    password_hash = ?,
    role = ?
where id = ?;

-- name: FindUsersSorted :many
select *
from users
where deleted_at is null
order by id asc;
