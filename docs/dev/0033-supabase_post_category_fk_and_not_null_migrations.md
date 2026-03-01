# Supabase Post Category FK and Not-Null Migrations

Date: 2026-03-01

## Summary
Added Supabase migration files to harden `public.post.category_id` integrity (cleanup, foreign key, and non-null enforcement), plus updated `.env.example` with direct Supabase connection guidance.

## What Changed

### 1. Added placeholder migration to align remote migration history
- Added `supabase/migrations/20260227204936_remote_schema.sql`:
  - placeholder comment-only migration used to align local migration history with remote state.

### 2. Added cleanup migration for invalid post category references
- Added `supabase/migrations/20260228235859_cleanup_invalid_post_category_id.sql`:
  - sets `public.post.category_id` to `NULL` for rows referencing missing categories.

### 3. Added post.category_id foreign key migration
- Added `supabase/migrations/20260228235900_add_post_category_fk.sql`:
  - adds `post_category_id_fkey` from `public.post(category_id)` to `public.category(id)`
  - uses `ON UPDATE CASCADE` and `ON DELETE SET NULL`
  - creates constraint as `NOT VALID`, then validates it.

### 4. Added category backfill and not-null enforcement migration
- Added `supabase/migrations/20260301000000_post_category_not_null.sql`:
  - backfills `post.category_id` from `subcategory.category_id` when available
  - asserts fallback category id `9` exists
  - assigns unresolved rows to fallback category `9`
  - sets `public.post.category_id` to `NOT NULL`.

### 5. Updated environment example with Supabase direct connection details
- Updated `.env.example`:
  - added `SUPABASE_DIRECT_CONNECTION` example value
  - added note with a working `supabase db push --db-url ...` command example.

## Why This Matters
- Prevents orphaned `post.category_id` references.
- Establishes relational integrity with an explicit foreign key.
- Enables stricter data guarantees by enforcing non-null category ownership for posts.
- Documents direct connection examples to support operational migration workflows.

## Files in This Increment
- `.env.example`
- `supabase/migrations/20260227204936_remote_schema.sql`
- `supabase/migrations/20260228235859_cleanup_invalid_post_category_id.sql`
- `supabase/migrations/20260228235900_add_post_category_fk.sql`
- `supabase/migrations/20260301000000_post_category_not_null.sql`
- `docs/dev/0033-supabase_post_category_fk_and_not_null_migrations.md`
