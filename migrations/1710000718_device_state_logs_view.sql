create view device_state_logs_view as
select
  dsl.created_at,
  dsl.device_id,
  d.name,
  d.icon,
  d.deleted_at,
  dsl.state
from device_state_logs dsl
join devices d on d.id = dsl.device_id
where d.deleted_at is null
and dsl.deleted_at is null
order by dsl.created_at desc;
