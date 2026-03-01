# Message Post ID Not-Null and Orphan Archive Migration

Date: 2026-03-01

## Summary
Added a Supabase migration to enforce `app_private.message.post_id` integrity by archiving/deleting orphaned rows, then applying `NOT NULL` and a restrictive foreign key to `public.post(id)`.

## What Changed

### 1. Added orphan archive table seed pattern
- Added `supabase/migrations/20260301005000_message_post_id_not_null.sql`:
  - creates `app_private.message_orphan_archive` (if missing) using a shape-copy pattern from `app_private.message`
  - includes `archived_at` timestamp column for archival provenance.

### 2. Archived existing orphaned message records
- Same migration:
  - inserts all `app_private.message` rows where `post_id is null` into archive table with `archived_at = now()`.

### 3. Removed orphaned rows from primary table
- Same migration:
  - deletes `app_private.message` rows where `post_id is null` before adding stricter constraints.

### 4. Enforced non-null and restrictive FK for message→post relationship
- Same migration:
  - sets `app_private.message.post_id` to `NOT NULL`
  - drops/recreates `message_post_id_fkey`
  - references `public.post(id)` with `ON UPDATE CASCADE` and `ON DELETE RESTRICT`
  - creates FK as `NOT VALID` then validates.

## Why This Matters
- Prevents future orphaned response-message records without a backing post.
- Preserves prior orphaned data for audit/recovery before cleanup.
- Strengthens referential integrity for message lifecycle and moderation workflows.

## Files in This Increment
- `supabase/migrations/20260301005000_message_post_id_not_null.sql`
- `docs/dev/0038-message_post_id_not_null_and_orphan_archive_migration.md`
