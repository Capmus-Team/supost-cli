-- Preserve orphaned message rows before enforcing non-null post_id.
create table if not exists app_private.message_orphan_archive as
select m.*, now()::timestamptz as archived_at
from app_private.message m
where false;

insert into app_private.message_orphan_archive (
  id,
  message,
  post_id,
  ip,
  email,
  created_at,
  updated_at,
  source,
  status,
  raw_email,
  user_agent,
  account_id,
  scammed,
  archived_at
)
select
  m.id,
  m.message,
  m.post_id,
  m.ip,
  m.email,
  m.created_at,
  m.updated_at,
  m.source,
  m.status,
  m.raw_email,
  m.user_agent,
  m.account_id,
  m.scammed,
  now()
from app_private.message m
where m.post_id is null;

delete from app_private.message
where post_id is null;

alter table app_private.message
  alter column post_id set not null;

alter table app_private.message
  drop constraint if exists message_post_id_fkey;

alter table app_private.message
  add constraint message_post_id_fkey
  foreign key (post_id)
  references public.post(id)
  on update cascade
  on delete restrict
  not valid;

alter table app_private.message
  validate constraint message_post_id_fkey;
