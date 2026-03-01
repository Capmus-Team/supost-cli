create table if not exists app_private.backfill_checkpoint (
  job_name text not null,
  last_post_id bigint not null default 0,
  processed_posts bigint not null default 0,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),
  constraint backfill_checkpoint_pkey primary key (job_name),
  constraint backfill_checkpoint_last_post_id_nonnegative check (last_post_id >= 0),
  constraint backfill_checkpoint_processed_posts_nonnegative check (processed_posts >= 0)
);

create trigger trg_backfill_checkpoint_set_updated_at
before update on app_private.backfill_checkpoint
for each row
execute function set_updated_at();
