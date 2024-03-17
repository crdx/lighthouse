-- name: FindUnprocessedDeviceLimitNotifications :many
select *
from device_limit_notifications
where deleted_at is null
and processed = 0
order by created_at asc;

-- name: CountPreviousDeviceLimitNotifications :one
select count(*)
from device_limit_notifications
where deleted_at is null
and device_id = ?
and processed = 0;
