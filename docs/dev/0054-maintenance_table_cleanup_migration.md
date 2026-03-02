# Maintenance Table Cleanup Migration

Date: 2026-03-02

## Summary
Documented and added a cleanup migration that removes temporary maintenance/archive tables no longer required in the current schema path.

## What Changed

### 1. Added cleanup migration for obsolete internal tables
- Added `supabase/migrations/20260301014000_drop_unused_checkpoint_and_archive_tables.sql`.
- Drops the maintenance trigger `trg_backfill_checkpoint_set_updated_at` if present.
- Drops `app_private.backfill_checkpoint` and `app_private.message_orphan_archive` tables if present.

## Why This Matters
- Keeps schema surface area focused on active runtime tables.
- Removes transitional artifacts from earlier backfill/archive workflows.
- Uses `if exists` guards so migration remains safe across environments with different prior state.

## Files in This Increment
- `docs/dev/0054-maintenance_table_cleanup_migration.md`
- `supabase/migrations/20260301014000_drop_unused_checkpoint_and_archive_tables.sql`
