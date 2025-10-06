-- name: FindAdaptersForDevice :many
select *
from adapters
where deleted_at is null
and device_id = ?
order by last_seen_at desc;

-- name: FindAdapterByIPAddressLatest :one
select *
from adapters
where deleted_at is null
and ip_address = ?
order by last_seen_at desc
limit 1;

-- name: UpdateAdapter :exec
update adapters
set
    name = ?,
    vendor = ?
where id = ?;
