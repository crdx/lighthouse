-- name: FindUnprocessedDeviceStateNotifications :many
select *
from device_state_notifications
where deleted_at is null
and processed = 0
order by created_at asc;

-- name: MigrateDeviceStateNotifications :exec
update device_state_notifications
set device_id = sqlc.arg(to_device_id)
where device_id = sqlc.arg(from_device_id);

-- name: CountDeviceStateNotificationsForTest :one
select count(*)
from device_state_notifications
where deleted_at is null
and device_id = ?
and state = ?
and processed = 0;
