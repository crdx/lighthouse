-- name: FindAuditLogsSorted :many
select *
from audit_logs
where deleted_at is null
order by created_at desc;
