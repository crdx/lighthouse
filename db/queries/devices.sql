-- name: FindDevicesSorted :many
select *
from devices
where deleted_at is null
order by name asc;

-- name: FindScannableDevices :many
select *
from devices
where deleted_at is null
and state = 'online'
and origin = false;

-- name: FindPingableDevices :many
select *
from devices
where deleted_at is null
and state = 'online'
and origin = false
and ping = 1;

-- name: ResetOriginDevices :exec
update devices set origin = false;

-- name: UpdateDevice :exec
update devices
set
    name = ?,
    icon = ?,
    notes = ?,
    grace_period = ?,
    watch = ?,
    ping = ?,
    `limit` = ?
where id = ?;

-- name: UpdateOriginDevice :exec
update devices
set
    name = ?,
    icon = ?,
    notes = ?
where id = ?;
