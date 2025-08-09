-- name: FindUnprocessedDeviceDiscoveryNotifications :many
select *
from device_discovery_notifications
where deleted_at is null
and processed = 0
order by created_at asc;

-- name: MigrateDeviceDiscoveryNotifications :exec
update device_discovery_notifications
set device_id = sqlc.arg(to_device_id)
where device_id = sqlc.arg(from_device_id);
