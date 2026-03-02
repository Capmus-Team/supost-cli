-- Drop no-longer-needed maintenance/archive tables.
-- Keep this additive: do not modify historical migrations.

drop trigger if exists trg_backfill_checkpoint_set_updated_at on app_private.backfill_checkpoint;

drop table if exists app_private.backfill_checkpoint;
drop table if exists app_private.message_orphan_archive;
