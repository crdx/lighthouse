-- name: FindAdaptersForDevice :many
select *
from adapters
where deleted_at is null
and device_id = ?
order by last_seen_at desc;

-- name: UpdateAdapter :exec
update adapters
set
    name = ?,
    vendor = ?
where id = ?;
