-- name: FindServicesForDevice :many
select *
from services
where deleted_at is null
and device_id = ?
order by port asc;

-- name: CountPreviousServices :one
select count(*)
from services
where deleted_at is null
and device_id = ?
and port = ?
and id != ?;

-- name: UpdateService :exec
update services
set name = ?
where id = ?;
