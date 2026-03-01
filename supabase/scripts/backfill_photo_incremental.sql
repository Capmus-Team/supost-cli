-- Incremental, resumable backfill from legacy post photo columns into public.photo.
--
-- This script does NOT modify data in public.post.
-- It only reads public.post and writes public.photo + app_private.backfill_checkpoint.
--
-- Usage examples:
--   psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f supabase/scripts/backfill_photo_incremental.sql
--   psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -v batch_size=2000 -v job_name=photo_backfill_v1 -f supabase/scripts/backfill_photo_incremental.sql
--
-- Re-run until selected_post_count = 0.

\if :{?job_name}
\else
\set job_name photo_backfill_v1
\endif

\if :{?batch_size}
\else
\set batch_size 1000
\endif

begin;

insert into app_private.backfill_checkpoint (job_name)
values (:'job_name')
on conflict (job_name) do nothing;

create temp table tmp_target_posts on commit drop as
select p.id
from public.post p
where p.id > (
  select c.last_post_id
  from app_private.backfill_checkpoint c
  where c.job_name = :'job_name'
)
and (
  nullif(trim(coalesce(p.image_source1, '')), '') is not null or
  nullif(trim(coalesce(p.image_source2, '')), '') is not null or
  nullif(trim(coalesce(p.image_source3, '')), '') is not null or
  nullif(trim(coalesce(p.image_source4, '')), '') is not null or
  nullif(trim(coalesce(p.photo1_file_name, '')), '') is not null or
  nullif(trim(coalesce(p.photo2_file_name, '')), '') is not null or
  nullif(trim(coalesce(p.photo3_file_name, '')), '') is not null or
  nullif(trim(coalesce(p.photo4_file_name, '')), '') is not null
)
order by p.id asc
limit :batch_size;

select 'selected_post_count' as label, count(*)::text as value
from tmp_target_posts;

with deleted as (
  delete from public.photo
  where post_id in (select id from tmp_target_posts)
  returning 1
)
select 'deleted_photo_rows' as label, count(*)::text as value
from deleted;

with expanded as (
  select
    p.id as post_id,
    v.position,
    v.raw_value
  from public.post p
  join tmp_target_posts t on t.id = p.id
  cross join lateral (
    values
      (0, coalesce(nullif(trim(p.image_source1), ''), nullif(trim(p.photo1_file_name), ''))),
      (1, coalesce(nullif(trim(p.image_source2), ''), nullif(trim(p.photo2_file_name), ''))),
      (2, coalesce(nullif(trim(p.image_source3), ''), nullif(trim(p.photo3_file_name), ''))),
      (3, coalesce(nullif(trim(p.image_source4), ''), nullif(trim(p.photo4_file_name), '')))
  ) as v(position, raw_value)
  where v.raw_value is not null
),
normalized as (
  select
    e.post_id,
    e.position,
    regexp_replace(e.raw_value, '^.*/', '') as base_value
  from expanded e
),
prepared as (
  select
    n.post_id,
    n.position,
    format(
      'posts/%s/%s',
      n.post_id,
      case
        when n.base_value like 'post_%' then n.base_value
        when n.base_value like 'ticker_%' then regexp_replace(n.base_value, '^ticker_', 'post_')
        else 'post_' || n.base_value
      end
    ) as s3_key
  from normalized n
),
dedup as (
  select
    post_id,
    s3_key,
    min(position) as position
  from prepared
  group by post_id, s3_key
),
ins as (
  insert into public.photo (post_id, s3_key, ticker_s3_key, position)
  select
    d.post_id,
    d.s3_key,
    regexp_replace(d.s3_key, '/post_', '/ticker_') as ticker_s3_key,
    d.position
  from dedup d
  order by d.post_id asc, d.position asc
  returning 1
)
select 'inserted_photo_rows' as label, count(*)::text as value
from ins;

update app_private.backfill_checkpoint c
set
  last_post_id = coalesce((select max(id) from tmp_target_posts), c.last_post_id),
  processed_posts = c.processed_posts + (select count(*) from tmp_target_posts),
  updated_at = now()
where c.job_name = :'job_name';

select
  c.job_name,
  c.last_post_id,
  c.processed_posts,
  c.updated_at
from app_private.backfill_checkpoint c
where c.job_name = :'job_name';

commit;
