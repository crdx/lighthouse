-- name: FindUnprocessedDeviceServiceNotifications :many
select *
from device_service_notifications
where deleted_at is null
and processed = 0
order by created_at asc
