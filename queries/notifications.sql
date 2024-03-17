-- name: FindNotificationsListView :many
select *
from notifications
where deleted_at is null
order by created_at
desc limit ?, ?;

-- name: CountNotificationsListView :one
select count(*)
from notifications
where deleted_at is null;
