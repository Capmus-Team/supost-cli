# Create Public Photo Table With Constraints, Indexes, and Trigger

Date: 2026-03-01

## Summary
Added a Supabase migration that introduces `public.photo` for post-attached images with FK integrity, uniqueness constraints, ordering support, and `updated_at` trigger behavior.

## What Changed

### 1. Added `public.photo` table migration
- Added `supabase/migrations/20260301010000_create_public_photo_table.sql`:
  - creates `public.photo` with identity PK `id`
  - links each photo to a post via `post_id bigint not null`.

### 2. Added integrity constraints for photo records
- Same migration adds:
  - FK `photo_post_id_fkey` to `public.post(id)` with `ON UPDATE CASCADE`, `ON DELETE CASCADE`
  - check constraint `position >= 0`
  - check constraint requiring non-blank `s3_key`
  - uniqueness on `(post_id, position)`
  - uniqueness on `(post_id, s3_key)`.

### 3. Added access/index and update timestamp behavior
- Same migration adds:
  - index `photo_post_id_position_idx` on `(post_id, position)`
  - trigger `trg_photo_set_updated_at` using `set_updated_at()` on row updates.

## Why This Matters
- Establishes normalized storage for multi-photo posts.
- Enforces per-post photo ordering and key uniqueness at the DB layer.
- Keeps photo rows consistent with parent-post lifecycle via cascading delete.

## Files in This Increment
- `supabase/migrations/20260301010000_create_public_photo_table.sql`
- `docs/dev/0045-create_public_photo_table_with_constraints_indexes_and_trigger.md`
