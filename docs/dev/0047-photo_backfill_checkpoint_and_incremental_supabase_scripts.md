# Photo Backfill Checkpoint and Incremental Supabase Scripts

Date: 2026-03-01

## Summary
Added migration support for resumable backfill checkpoints and introduced Supabase SQL scripts for one-shot and incremental backfill of `public.photo` from legacy `public.post` image columns.

## What Changed

### 1. Added backfill checkpoint table migration
- Added `supabase/migrations/20260301012000_create_photo_backfill_checkpoint.sql`:
  - creates `app_private.backfill_checkpoint`
  - tracks `job_name`, `last_post_id`, and `processed_posts`
  - enforces non-negative checkpoint counters
  - adds `trg_backfill_checkpoint_set_updated_at` trigger.

### 2. Added one-shot backfill script from legacy post fields
- Added `supabase/scripts/backfill_photo_from_legacy.sql`:
  - selects recent posts with legacy image values (`image_source*` / `photo*_file_name`)
  - clears existing `public.photo` rows for selected posts
  - normalizes source values into canonical `posts/<post_id>/post_*` keys
  - inserts into `public.photo` with derived `ticker_s3_key`.

### 3. Added incremental/resumable backfill script
- Added `supabase/scripts/backfill_photo_incremental.sql`:
  - runs batch-oriented backfill by ascending post id
  - uses `app_private.backfill_checkpoint` state to resume from prior progress
  - updates checkpoint counters each run
  - repeat until `selected_post_count = 0`.

### 4. Added script usage guidance
- Added `supabase/scripts/README.md`:
  - documents script purpose
  - includes `psql` run commands and checkpoint query.

## Why This Matters
- Enables safe, resumable migration of legacy photo references into normalized `public.photo`.
- Keeps backfill operations repeatable and observable via explicit checkpoint state.
- Reduces risk when processing large post volumes by supporting incremental batches.

## Files in This Increment
- `supabase/migrations/20260301012000_create_photo_backfill_checkpoint.sql`
- `supabase/scripts/README.md`
- `supabase/scripts/backfill_photo_from_legacy.sql`
- `supabase/scripts/backfill_photo_incremental.sql`
- `docs/dev/0047-photo_backfill_checkpoint_and_incremental_supabase_scripts.md`
