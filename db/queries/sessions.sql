-- name: FindOtherSessions :many
select k, v
from sessions
where k != ?;

-- name: DeleteSession :exec
delete from sessions
where k = ?;
