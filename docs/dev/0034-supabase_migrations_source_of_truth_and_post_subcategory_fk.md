# Supabase Migrations Source of Truth and Post Subcategory FK

Date: 2026-03-01

## Summary
Updated docs/tooling to treat `supabase/migrations/` as the active migration source of truth, removed legacy `migrations/README.md`, and added a migration that enforces `post.subcategory_id` referential integrity.

## What Changed

### 1. Added post subcategory FK migration
- Added `supabase/migrations/20260301001000_add_post_subcategory_fk.sql`:
  - sets orphaned `post.subcategory_id` references to `NULL`
  - adds `post_subcategory_id_fkey` from `public.post(subcategory_id)` to `public.subcategory(id)`
  - uses `ON UPDATE CASCADE`, `ON DELETE SET NULL`
  - adds the FK as `NOT VALID` and then validates it.

### 2. Updated migration command guidance in Makefile
- Updated `Makefile` `migrate` target:
  - replaced printed `psql migrations/*.sql` guidance
  - now prints `supabase db push --db-url "$$DATABASE_URL"` guidance.

### 3. Updated README migration/documentation references
- Updated `README.md`:
  - project tree now documents `supabase/migrations/` as SQL schema + migration history source
  - production migration section now points to applying `supabase/migrations/*.sql`.

### 4. Removed legacy migration README and preserved directory contract
- Deleted `migrations/README.md`:
  - removed instructions tied to the old `migrations/*.sql` location.
- Added `migrations/.gitkeep`:
  - keeps the top-level `migrations/` directory present for project-structure contract tests.

### 5. Expanded env example for Supabase DB usage
- Updated `.env.example`:
  - added `SUPABASE_DB_PASSWORD` example variable for database command workflows.

## Why This Matters
- Aligns local docs and commands with Supabase CLI-driven migration workflows.
- Adds missing FK protection for `post.subcategory_id` data quality.
- Reduces ambiguity by removing outdated migration location instructions.

## Files in This Increment
- `.env.example`
- `Makefile`
- `README.md`
- `migrations/.gitkeep`
- `supabase/migrations/20260301001000_add_post_subcategory_fk.sql`
- `docs/dev/0034-supabase_migrations_source_of_truth_and_post_subcategory_fk.md`
- deleted: `migrations/README.md`
