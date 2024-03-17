-- name: FindLatestActivityForDevice :many
select *
from device_state_logs
where deleted_at is null
and device_id = ?
order by created_at desc
limit ?;

-- name: FindDeviceStateLogsListView :many
select *
from device_state_logs_view
limit ?, ?;

-- name: FindDeviceStateLogsListViewForDevice :many
select *
from device_state_logs_view
where device_id = ?
limit ?, ?;

-- name: CountDeviceStateLogsListView :one
select count(*)
from device_state_logs_view;

-- name: CountDeviceStateLogsListViewForDevice :one
select count(*)
from device_state_logs_view
where device_id = ?;

-- name: MigrateDeviceStateLogs :exec
update device_state_logs
set device_id = sqlc.arg(to_device_id)
where device_id = sqlc.arg(from_device_id);

-- name: CountDeviceStateLogsForTest :one
select count(*)
from device_state_logs
where deleted_at is null
and device_id = ?
and state = ?;
